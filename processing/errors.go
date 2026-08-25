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
	defer lock.Unlock()

	typ, typeErr := dep.TypeOf(elementID)
	if typeErr != nil {
		typ = eventv1.Element_TYPE_UNSPECIFIED
	}
	tok := inst.Tokens[tokenID]
	if tok == nil || tok.ElementID != elementID || tok.Status != projection.TokenWaiting {
		return e.reject(ctx, inst, elementID, tokenID, typ, eventv1.Element_INTENT_ERROR_THROWN, "INVALID_STATE", "element is not waiting for completion")
	}
	if typeErr != nil {
		return e.reject(ctx, inst, elementID, tokenID, typ, eventv1.Element_INTENT_ERROR_THROWN, "NOT_FOUND", "element not found")
	}
	switch typ {
	case eventv1.Element_TYPE_USER_TASK, eventv1.Element_TYPE_SERVICE_TASK:
	default:
		return e.reject(ctx, inst, elementID, tokenID, typ, eventv1.Element_INTENT_ERROR_THROWN, "INVALID_STATE", "element cannot throw an error")
	}

	cmdID, err := NextID()
	if err != nil {
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
		return err
	}
	if typ == eventv1.Element_TYPE_SERVICE_TASK && tok.JobType != "" {
		e.releaseLease(instanceID, tokenID)
		e.notifyJobs()
	}
	_, err = e.executor.propagateError(ctx, dep, inst, elementID, tokenID, errorCode, emit)
	return err
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
		case eventv1.Element_TYPE_USER_TASK, eventv1.Element_TYPE_SERVICE_TASK:
			if bid, ok := dep.MatchErrorBoundary(throwElementID, errorCode); ok {
				return x.fireActivityErrorBoundary(ctx, dep, inst, throwElementID, throwTokenID, bid, emit)
			}
		}
	}

	scope, ok := dep.ScopeOf(throwElementID)
	if !ok {
		return nil, fmt.Errorf("NOT_FOUND: scope for element %q", throwElementID)
	}
	processID := dep.ProcessID()

	for {
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
		return nil, x.terminateInstance(ctx, dep, inst, emit)
	}
	return nil, x.terminateInstance(ctx, dep, inst, emit)
}

func findScopeHostToken(dep *deploy.Deployment, inst *projection.Instance, scopeID string) (string, bool) {
	var fallback string
	for tid, tok := range inst.Tokens {
		if tok == nil {
			continue
		}
		if tok.ElementID == scopeID {
			return tid, true
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
	records := append(handlers.CancelAttachedBoundaries(dep, activityID, activityTokenID),
		&eventv1.Element{Intent: eventv1.Element_INTENT_TERMINATING, Type: typ, Id: activityID, TokenId: activityTokenID},
		&eventv1.Element{Intent: eventv1.Element_INTENT_TERMINATED, Type: typ, Id: activityID, TokenId: activityTokenID},
		&eventv1.Element{Intent: eventv1.Element_INTENT_COMPLETING, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: boundaryID, TokenId: activityTokenID},
		&eventv1.Element{Intent: eventv1.Element_INTENT_COMPLETED, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: boundaryID, TokenId: activityTokenID},
	)
	for _, rec := range records {
		if err := emit(rec); err != nil {
			return nil, err
		}
	}
	next, err := x.takeOutgoing(dep, activityTokenID, boundaryID, "", emit)
	if err != nil {
		return nil, err
	}
	return x.Enter(ctx, dep, inst, activityTokenID, next, emit)
}

func (x *Executor) fireScopeErrorBoundary(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	scopeID, hostTokenID, boundaryID string,
	emit Emitter,
) ([]handlers.Publication, error) {
	if err := x.terminateScopeTokens(dep, inst, scopeID, emit); err != nil {
		return nil, err
	}
	if err := emit(&eventv1.Element{
		Intent:  eventv1.Element_INTENT_TERMINATING,
		Type:    eventv1.Element_TYPE_SUB_PROCESS,
		Id:      scopeID,
		TokenId: hostTokenID,
	}); err != nil {
		return nil, err
	}
	if err := emit(&eventv1.Element{
		Intent:  eventv1.Element_INTENT_TERMINATED,
		Type:    eventv1.Element_TYPE_SUB_PROCESS,
		Id:      scopeID,
		TokenId: hostTokenID,
	}); err != nil {
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
	next, err := x.takeOutgoing(dep, hostTokenID, boundaryID, "", emit)
	if err != nil {
		return nil, err
	}
	return x.Enter(ctx, dep, inst, hostTokenID, next, emit)
}

func (x *Executor) terminateScope(ctx context.Context, dep *deploy.Deployment, inst *projection.Instance, scopeID string, emit Emitter) error {
	if err := x.terminateScopeTokens(dep, inst, scopeID, emit); err != nil {
		return err
	}
	hostTokenID, ok := findScopeHostToken(dep, inst, scopeID)
	if !ok {
		inst.RemoveScopeBoundariesForScope(scopeID)
		return nil
	}
	if err := emit(&eventv1.Element{
		Intent:  eventv1.Element_INTENT_TERMINATING,
		Type:    eventv1.Element_TYPE_SUB_PROCESS,
		Id:      scopeID,
		TokenId: hostTokenID,
	}); err != nil {
		return err
	}
	if err := emit(&eventv1.Element{
		Intent:  eventv1.Element_INTENT_TERMINATED,
		Type:    eventv1.Element_TYPE_SUB_PROCESS,
		Id:      scopeID,
		TokenId: hostTokenID,
	}); err != nil {
		return err
	}
	inst.RemoveScopeBoundariesForScope(scopeID)
	return nil
}

func (x *Executor) terminateScopeTokens(dep *deploy.Deployment, inst *projection.Instance, scopeID string, emit Emitter) error {
	for tid, tok := range inst.Tokens {
		if tok == nil {
			continue
		}
		tokScope, _ := dep.ScopeOf(tok.ElementID)
		if tokScope != scopeID && !isInScope(dep, tokScope, scopeID) && tok.ElementID != scopeID {
			continue
		}
		if tok.ElementID == scopeID {
			continue
		}
		tokType, err := dep.TypeOf(tok.ElementID)
		if err != nil {
			continue
		}
		for _, intent := range []eventv1.Element_Intent{
			eventv1.Element_INTENT_TERMINATING,
			eventv1.Element_INTENT_TERMINATED,
		} {
			if err := emit(&eventv1.Element{
				Intent:  intent,
				Type:    tokType,
				Id:      tok.ElementID,
				TokenId: tid,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (x *Executor) terminateInstance(ctx context.Context, dep *deploy.Deployment, inst *projection.Instance, emit Emitter) error {
	for tid, tok := range inst.Tokens {
		if tok == nil {
			continue
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
				return err
			}
		}
	}
	if err := emit(&eventv1.Element{
		Intent: eventv1.Element_INTENT_TERMINATING,
		Type:   eventv1.Element_TYPE_PROCESS,
		Id:     dep.ProcessID(),
	}); err != nil {
		return err
	}
	return emit(&eventv1.Element{
		Intent: eventv1.Element_INTENT_TERMINATED,
		Type:   eventv1.Element_TYPE_PROCESS,
		Id:     dep.ProcessID(),
	})
}
