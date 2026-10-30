package processing

import (
	"context"
	"fmt"

	"github.com/sparrow-community/sparrow/processing/deploy"
	"github.com/sparrow-community/sparrow/processing/handlers"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// ThrowError records ERROR_THROWN on a waiting activity and propagates the error
// to a matching boundary or terminates the instance when unhandled.
func (e *Engine) ThrowError(ctx context.Context, instanceID, elementID, tokenID, errorCode string) error {
	if err := e.errIfInstancePaused(instanceID); err != nil {
		return err
	}
	e.mu.Lock()
	inst := e.instances[instanceID]
	lock := e.instMu[instanceID]
	var dep *deploy.Deployment
	if inst != nil {
		dep = e.deployments[inst.DeploymentID]
	}
	e.mu.Unlock()
	if inst == nil || dep == nil {
		return fmt.Errorf("NOT_FOUND: instance %q", instanceID)
	}

	lock.Lock()

	typ, typeErr := dep.TypeOf(elementID)
	if typeErr != nil {
		typ = eventv1.Element_TYPE_UNSPECIFIED
	}
	tok := inst.Tokens[tokenID]
	if tok == nil || tok.ElementID != elementID || tok.Status != projection.TokenWaiting {
		lock.Unlock()
		return e.reject(ctx, inst, elementID, tokenID, typ, eventv1.Element_INTENT_ERROR_THROWN, "INVALID_STATE", "element is not waiting for completion")
	}
	if typeErr != nil {
		lock.Unlock()
		return e.reject(ctx, inst, elementID, tokenID, typ, eventv1.Element_INTENT_ERROR_THROWN, "NOT_FOUND", "element not found")
	}
	switch typ {
	case eventv1.Element_TYPE_USER_TASK, eventv1.Element_TYPE_SERVICE_TASK, eventv1.Element_TYPE_TASK, eventv1.Element_TYPE_MANUAL_TASK, eventv1.Element_TYPE_RECEIVE_TASK, eventv1.Element_TYPE_BUSINESS_RULE_TASK, eventv1.Element_TYPE_SCRIPT_TASK, eventv1.Element_TYPE_CALL_ACTIVITY:
	default:
		lock.Unlock()
		return e.reject(ctx, inst, elementID, tokenID, typ, eventv1.Element_INTENT_ERROR_THROWN, "INVALID_STATE", "element cannot throw an error")
	}

	cmdID, err := NextID()
	if err != nil {
		lock.Unlock()
		return err
	}
	cmd := &eventv1.Event{
		Id:                cmdID,
		Timestamp:         nowMillis(),
		RecordType:        eventv1.Event_RECORD_TYPE_COMMAND,
		DeploymentId:      inst.DeploymentID,
		ProcessInstanceId: instanceID,
		ProcessVersion:    inst.Version,
		Element: &eventv1.Element{
			Intent:  eventv1.Element_INTENT_ERROR_THROWN,
			Type:    typ,
			Id:      elementID,
			TokenId: tokenID,
			Payload: &eventv1.Element_EventPayload{
				EventPayload: &eventv1.EventPayload{ErrorCode: errorCode},
			},
		},
	}
	if _, err := e.log.Append(ctx, cmd); err != nil {
		lock.Unlock()
		return err
	}

	emit := e.emitter(ctx, inst, cmdID)
	if err := emit(&eventv1.Element{
		Intent:  eventv1.Element_INTENT_ERROR_THROWN,
		Type:    typ,
		Id:      elementID,
		TokenId: tokenID,
		Payload: &eventv1.Element_EventPayload{
			EventPayload: &eventv1.EventPayload{ErrorCode: errorCode},
		},
	}); err != nil {
		lock.Unlock()
		return err
	}
	if (typ == eventv1.Element_TYPE_SERVICE_TASK || typ == eventv1.Element_TYPE_BUSINESS_RULE_TASK || typ == eventv1.Element_TYPE_SCRIPT_TASK) && tok.JobType != "" {
		e.releaseLease(instanceID, tokenID)
		e.notifyJobs()
	}
	pubs, err := e.executor.propagateError(ctx, dep, inst, elementID, tokenID, errorCode, emit)
	lock.Unlock()
	if err != nil {
		return err
	}
	if err := e.flushPublications(ctx, pubs); err != nil {
		return err
	}
	return e.tryDeliverBuffered(ctx, instanceID)
}

func (x *Executor) propagateError(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	throwElementID, throwTokenID, errorCode string,
	emit Emitter,
) ([]handlers.Publication, error) {
	if throwTyp, err := dep.TypeOf(throwElementID); err == nil {
		switch throwTyp {
		case eventv1.Element_TYPE_USER_TASK, eventv1.Element_TYPE_SERVICE_TASK, eventv1.Element_TYPE_TASK, eventv1.Element_TYPE_MANUAL_TASK, eventv1.Element_TYPE_RECEIVE_TASK, eventv1.Element_TYPE_BUSINESS_RULE_TASK, eventv1.Element_TYPE_SCRIPT_TASK, eventv1.Element_TYPE_CALL_ACTIVITY:
			if bid, ok := dep.MatchErrorBoundary(throwElementID, errorCode); ok {
				return x.fireActivityErrorBoundary(ctx, dep, inst, throwElementID, throwTokenID, bid, emit)
			}
		}
	}

	scope, ok := dep.ScopeOf(throwElementID)
	if !ok {
		return nil, fmt.Errorf("NOT_FOUND: scope for element %q", throwElementID)
	}
	processID := inst.ProcessID
	if processID == "" {
		processID = dep.ProcessID()
	}

	for {
		if eventSubProcessID, ok := dep.MatchErrorEventSubProcess(scope, errorCode, func(id string) bool {
			_, armed := inst.EventSubProcesses[id]
			return armed
		}); ok {
			return x.fireErrorEventSubProcess(ctx, dep, inst, throwElementID, throwTokenID, eventSubProcessID, emit)
		}
		if scope != processID {
			if bid, ok := dep.MatchErrorBoundary(scope, errorCode); ok {
				hostTokenID, ok := findScopeHostToken(dep, inst, scope)
				if !ok {
					return nil, fmt.Errorf("INVALID_STATE: no token for scope %q", scope)
				}
				return x.fireScopeErrorBoundary(ctx, dep, inst, scope, hostTokenID, bid, emit)
			}
			if err := x.terminateScope(ctx, dep, inst, scope, emit); err != nil {
				return nil, err
			}
			parent, ok := dep.ScopeOf(scope)
			if !ok {
				break
			}
			scope = parent
			continue
		}
		return x.terminateInstancePubs(ctx, dep, inst, emit)
	}
	return x.terminateInstancePubs(ctx, dep, inst, emit)
}

func (x *Executor) fireErrorEventSubProcess(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	throwElementID, throwTokenID, eventSubProcessElementID string,
	emit Emitter,
) ([]handlers.Publication, error) {
	if _, ok := inst.EventSubProcesses[eventSubProcessElementID]; !ok {
		return nil, nil
	}
	spec, ok := dep.EventSubProcessSpec(eventSubProcessElementID)
	if !ok {
		return nil, fmt.Errorf("NOT_FOUND: event subProcess %q", eventSubProcessElementID)
	}

	tokenID, err := NextID()
	if err != nil {
		return nil, err
	}

	if spec.Interrupting {
		if err := emitEventSubProcessStartDisarmInScope(dep, spec.ParentScopeID, inst, emit); err != nil {
			return nil, err
		}
		if err := terminateScopeTokens(dep, inst, spec.ParentScopeID, emit, scopeTerminateOpts{
			IncludeHost: true,
			DropTokens:  true,
		}); err != nil {
			return nil, err
		}
		inst.RemoveScopeBoundariesForScope(spec.ParentScopeID)
	} else {
		if err := emitEventSubProcessStartDisarm(dep, eventSubProcessElementID, emit); err != nil {
			return nil, err
		}
		// Consume the thrower so it does not leave a dangling waiting/active token.
		if throwTokenID != "" && inst.Tokens[throwTokenID] != nil {
			typ, typeErr := dep.TypeOf(throwElementID)
			if typeErr == nil {
				for _, intent := range []eventv1.Element_Intent{
					eventv1.Element_INTENT_TERMINATING,
					eventv1.Element_INTENT_TERMINATED,
				} {
					if err := emit(&eventv1.Element{
						Intent:  intent,
						Type:    typ,
						Id:      throwElementID,
						TokenId: throwTokenID,
					}); err != nil {
						return nil, err
					}
				}
			}
		}
	}

	return x.Enter(ctx, dep, inst, tokenID, eventSubProcessElementID, emit)
}

func findScopeHostToken(dep *deploy.Deployment, inst *projection.Instance, scopeID string) (string, bool) {
	if hostID, _, ok := scopeHostElement(dep, scopeID); ok {
		for tid, tok := range inst.Tokens {
			if tok != nil && tok.ElementID == hostID {
				return tid, true
			}
		}
	}
	// Fallback: any token inside the scope (pre-host-token ledgers).
	var fallback string
	for tid, tok := range inst.Tokens {
		if tok == nil {
			continue
		}
		tokScope, ok := dep.ScopeOf(tok.ElementID)
		if !ok {
			continue
		}
		if tokScope == scopeID || isInScope(dep, tokScope, scopeID) {
			if fallback == "" {
				fallback = tid
			}
		}
	}
	if fallback != "" {
		return fallback, true
	}
	return "", false
}

func (x *Executor) fireActivityErrorBoundary(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	activityID, activityTokenID, boundaryID string,
	emit Emitter,
) ([]handlers.Publication, error) {
	typ, err := dep.TypeOf(activityID)
	if err != nil {
		return nil, err
	}
	var pubs []handlers.Publication
	if tok := inst.Tokens[activityTokenID]; tok != nil {
		if pub := terminateChildPub(tok, activityID); pub != nil {
			pubs = append(pubs, *pub)
		}
	}
	records := append(handlers.CancelAttachedBoundaries(dep, activityID, activityTokenID),
		&eventv1.Element{Intent: eventv1.Element_INTENT_TERMINATING, Type: typ, Id: activityID, TokenId: activityTokenID},
		&eventv1.Element{Intent: eventv1.Element_INTENT_TERMINATED, Type: typ, Id: activityID, TokenId: activityTokenID},
		&eventv1.Element{Intent: eventv1.Element_INTENT_COMPLETING, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: boundaryID, TokenId: activityTokenID},
		&eventv1.Element{Intent: eventv1.Element_INTENT_COMPLETED, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: boundaryID, TokenId: activityTokenID},
	)
	for _, rec := range records {
		if err := emit(rec); err != nil {
			return pubs, err
		}
	}
	next, err := x.takeOutgoing(dep, inst, activityTokenID, boundaryID, "", emit)
	if err != nil {
		return pubs, err
	}
	more, err := x.Enter(ctx, dep, inst, activityTokenID, next, emit)
	pubs = append(pubs, more...)
	return pubs, err
}

func (x *Executor) fireScopeErrorBoundary(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	scopeID, hostTokenID, boundaryID string,
	emit Emitter,
) ([]handlers.Publication, error) {
	if err := terminateScopeTokens(dep, inst, scopeID, emit, scopeTerminateOpts{
		IncludeHost: false,
		DropTokens:  true,
	}); err != nil {
		return nil, err
	}
	// Keep token_id: boundary outgoing continues on this token.
	if err := terminateScopeHost(dep, scopeID, hostTokenID, emit); err != nil {
		return nil, err
	}
	if err := emitEventSubProcessStartDisarmInScope(dep, scopeID, inst, emit); err != nil {
		return nil, err
	}
	inst.RemoveScopeBoundariesForScope(scopeID)
	if err := emit(&eventv1.Element{
		Intent:  eventv1.Element_INTENT_COMPLETING,
		Type:    eventv1.Element_TYPE_BOUNDARY_EVENT,
		Id:      boundaryID,
		TokenId: hostTokenID,
	}); err != nil {
		return nil, err
	}
	if err := emit(&eventv1.Element{
		Intent:  eventv1.Element_INTENT_COMPLETED,
		Type:    eventv1.Element_TYPE_BOUNDARY_EVENT,
		Id:      boundaryID,
		TokenId: hostTokenID,
	}); err != nil {
		return nil, err
	}
	next, err := x.takeOutgoing(dep, inst, hostTokenID, boundaryID, "", emit)
	if err != nil {
		return nil, err
	}
	return x.Enter(ctx, dep, inst, hostTokenID, next, emit)
}

func (x *Executor) terminateScope(ctx context.Context, dep *deploy.Deployment, inst *projection.Instance, scopeID string, emit Emitter) error {
	_, hadHost := findScopeHostToken(dep, inst, scopeID)
	if err := terminateScopeTokens(dep, inst, scopeID, emit, scopeTerminateOpts{
		IncludeHost: true,
		DropTokens:  true,
	}); err != nil {
		return err
	}
	// IncludeHost already TERMINATED a parked host (SubProcess or CallActivity).
	// Audit-only when no host token existed (legacy SubProcess without host).
	if !hadHost {
		if err := terminateScopeHost(dep, scopeID, "", emit); err != nil {
			return err
		}
	}
	if err := emitEventSubProcessStartDisarmInScope(dep, scopeID, inst, emit); err != nil {
		return err
	}
	inst.RemoveScopeBoundariesForScope(scopeID)
	return nil
}

func (x *Executor) terminateInstance(ctx context.Context, dep *deploy.Deployment, inst *projection.Instance, emit Emitter) error {
	pubs, err := x.terminateInstancePubs(ctx, dep, inst, emit)
	_ = pubs
	return err
}

func (x *Executor) terminateInstancePubs(ctx context.Context, dep *deploy.Deployment, inst *projection.Instance, emit Emitter) ([]handlers.Publication, error) {
	var pubs []handlers.Publication
	for tid, tok := range inst.Tokens {
		if tok == nil {
			continue
		}
		if pub := terminateChildPub(tok, tok.ElementID); pub != nil {
			pubs = append(pubs, *pub)
		}
		typ, err := dep.TypeOf(tok.ElementID)
		if err != nil {
			continue
		}
		for _, intent := range []eventv1.Element_Intent{
			eventv1.Element_INTENT_TERMINATING,
			eventv1.Element_INTENT_TERMINATED,
		} {
			if err := emit(&eventv1.Element{
				Intent:  intent,
				Type:    typ,
				Id:      tok.ElementID,
				TokenId: tid,
			}); err != nil {
				return pubs, err
			}
		}
	}
	pid := inst.ProcessID
	if pid == "" {
		pid = dep.ProcessID()
	}
	pp := &eventv1.ProcessPayload{}
	if inst.ParentProcessInstanceID != "" {
		pp.ParentProcessInstanceId = inst.ParentProcessInstanceID
		pp.ParentElementId = inst.ParentElementID
		pp.ParentTokenId = inst.ParentTokenID
	}
	for _, intent := range []eventv1.Element_Intent{
		eventv1.Element_INTENT_TERMINATING,
		eventv1.Element_INTENT_TERMINATED,
	} {
		el := &eventv1.Element{
			Intent: intent,
			Type:   eventv1.Element_TYPE_PROCESS,
			Id:     pid,
		}
		if inst.ParentProcessInstanceID != "" {
			el.Payload = &eventv1.Element_ProcessPayload{ProcessPayload: pp}
		}
		if err := emit(el); err != nil {
			return pubs, err
		}
	}
	if inst.ParentProcessInstanceID != "" {
		pubs = append(pubs, handlers.Publication{
			Kind:             handlers.PublicationResumeParent,
			ParentInstanceID: inst.ParentProcessInstanceID,
			CallActivityID:   inst.ParentElementID,
			HostTokenID:      inst.ParentTokenID,
			ChildInstanceID:  inst.ID,
			Completed:        false,
		})
	}
	return pubs, nil
}
