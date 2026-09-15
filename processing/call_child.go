package processing

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/sparrow-community/sparrow/processing/deploy"
	"github.com/sparrow-community/sparrow/processing/handlers"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func init() {
	handlers.SetCallChildIDAllocator(NextID)
	handlers.SetCalleeResolver(func(processID string) (string, error) {
		if calleeResolver == nil {
			return "", fmt.Errorf("NOT_FOUND: process %q", processID)
		}
		return calleeResolver(processID)
	})
}

// calleeResolver is set per Engine via wireCalleeResolver.
var calleeResolver func(processID string) (string, error)

func (e *Engine) wireCalleeResolver() {
	calleeResolver = func(processID string) (string, error) {
		e.mu.Lock()
		defer e.mu.Unlock()
		dep, err := e.resolveDeploymentLocked("", processID, 0)
		if err != nil {
			return "", err
		}
		return dep.ID, nil
	}
}

func (e *Engine) startCalledInstance(ctx context.Context, p handlers.Publication) error {
	parentID := p.ParentInstanceID
	e.mu.Lock()
	callerDep := e.deployments[p.DeploymentID]
	parent := e.instances[parentID]
	if parent == nil && p.HostTokenID != "" {
		for id, inst := range e.instances {
			if inst == nil {
				continue
			}
			tok := inst.Tokens[p.HostTokenID]
			if tok != nil && tok.ElementID == p.CallActivityID && tok.CalledProcessInstanceID == p.ChildInstanceID {
				parent = inst
				parentID = id
				break
			}
		}
	}
	calleeDepID := p.CalledDeploymentID
	if calleeDepID == "" {
		calleeDepID = p.DeploymentID
	}
	calleeDep := e.deployments[calleeDepID]
	e.mu.Unlock()
	if callerDep == nil {
		return fmt.Errorf("NOT_FOUND: deployment %q", p.DeploymentID)
	}
	if calleeDep == nil {
		return fmt.Errorf("NOT_FOUND: deployment %q", calleeDepID)
	}
	if parent == nil {
		return fmt.Errorf("NOT_FOUND: parent instance for callActivity %q", p.CallActivityID)
	}
	call, ok := callerDep.CallActivitySpec(p.CallActivityID)
	if !ok {
		return fmt.Errorf("NOT_FOUND: callActivity %q", p.CallActivityID)
	}
	calledProcessID := p.CalledProcessID
	if calledProcessID == "" {
		calledProcessID = call.CalledProcessID
	}

	pv := variablesFromJSONStrings(applyMappings(parent.Variables, call.Inputs))
	if len(p.ChildVariables) > 0 {
		pv = p.ChildVariables
	}

	instanceID := p.ChildInstanceID
	var err error
	if instanceID == "" {
		instanceID, err = NextID()
		if err != nil {
			return err
		}
	}
	tokenID, err := NextID()
	if err != nil {
		return err
	}
	cmdID, err := NextID()
	if err != nil {
		return err
	}

	cmd := &eventv1.Event{
		Id:                cmdID,
		Timestamp:         nowMillis(),
		RecordType:        eventv1.Event_RECORD_TYPE_COMMAND,
		DeploymentId:      calleeDep.ID,
		ProcessInstanceId: instanceID,
		ProcessVersion:    calleeDep.Version,
		Element: &eventv1.Element{
			Intent:  eventv1.Element_INTENT_ACTIVATING,
			Type:    eventv1.Element_TYPE_PROCESS,
			Id:      calledProcessID,
			TokenId: tokenID,
			Payload: &eventv1.Element_ProcessPayload{
				ProcessPayload: &eventv1.ProcessPayload{
					Variables:               pv,
					ParentProcessInstanceId: parentID,
					ParentElementId:         p.CallActivityID,
					ParentTokenId:           p.HostTokenID,
				},
			},
		},
	}
	if _, err := e.log.Append(ctx, cmd); err != nil {
		return err
	}

	inst := projection.NewInstance(instanceID, calleeDep.ID, calleeDep.Version)
	inst.ProcessID = calledProcessID
	inst.ParentProcessInstanceID = parentID
	inst.ParentElementID = p.CallActivityID
	inst.ParentTokenID = p.HostTokenID

	lock := &sync.Mutex{}
	e.mu.Lock()
	e.instances[instanceID] = inst
	e.instMu[instanceID] = lock
	e.mu.Unlock()

	startEventID := call.StartEventID
	if startEventID == "" {
		startEventID, err = calleeDep.StartEventID()
		if err != nil {
			return err
		}
	}

	lock.Lock()
	emit := e.emitter(ctx, inst, cmdID)
	for _, rec := range handlers.ProcessStartRecordsWithParent(calledProcessID, pv, parentID, p.CallActivityID, p.HostTokenID) {
		if err := emit(rec); err != nil {
			lock.Unlock()
			return err
		}
	}
	if err := emitEventSubProcessStartArms(calleeDep, inst, calledProcessID, e.now(), emit); err != nil {
		lock.Unlock()
		return err
	}
	pubs, err := e.executor.Enter(ctx, calleeDep, inst, tokenID, startEventID, emit)
	lock.Unlock()
	if err != nil {
		return err
	}
	if err := e.flushPublications(ctx, pubs); err != nil {
		return err
	}
	return e.tryDeliverBuffered(ctx, instanceID)
}

func (e *Engine) resumeParentCall(ctx context.Context, p handlers.Publication) error {
	if p.ParentInstanceID == "" || p.CallActivityID == "" || p.HostTokenID == "" {
		return fmt.Errorf("INVALID_ARGUMENT: resume parent requires parent instance, call activity, and host token")
	}
	e.mu.Lock()
	parent := e.instances[p.ParentInstanceID]
	lock := e.instMu[p.ParentInstanceID]
	var callerDep *deploy.Deployment
	if parent != nil {
		callerDep = e.deployments[parent.DeploymentID]
	}
	child := e.instances[p.ChildInstanceID]
	e.mu.Unlock()
	if parent == nil || callerDep == nil {
		return fmt.Errorf("NOT_FOUND: parent instance %q", p.ParentInstanceID)
	}

	if !p.Completed {
		lock.Lock()
		tok := parent.Tokens[p.HostTokenID]
		if tok == nil || tok.ElementID != p.CallActivityID || tok.Status != projection.TokenWaiting {
			// Host may already be gone; still continue parent unfinished-call compensation.
			pubs, err := e.continueParentAfterUnfinishedCallChild(ctx, parent, callerDep, lock)
			lock.Unlock()
			if err != nil {
				return err
			}
			return e.flushPublications(ctx, pubs)
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
			DeploymentId:      parent.DeploymentID,
			ProcessInstanceId: parent.ID,
			ProcessVersion:    parent.Version,
			Element: &eventv1.Element{
				Intent:  eventv1.Element_INTENT_TERMINATING,
				Type:    eventv1.Element_TYPE_CALL_ACTIVITY,
				Id:      p.CallActivityID,
				TokenId: p.HostTokenID,
			},
		}
		if _, err := e.log.Append(ctx, cmd); err != nil {
			lock.Unlock()
			return err
		}
		emit := e.emitter(ctx, parent, cmdID)
		for _, intent := range []eventv1.Element_Intent{
			eventv1.Element_INTENT_TERMINATING,
			eventv1.Element_INTENT_TERMINATED,
		} {
			if err := emit(&eventv1.Element{
				Intent:  intent,
				Type:    eventv1.Element_TYPE_CALL_ACTIVITY,
				Id:      p.CallActivityID,
				TokenId: p.HostTokenID,
				Payload: &eventv1.Element_ActivityPayload{
					ActivityPayload: &eventv1.ActivityPayload{CalledProcessInstanceId: p.ChildInstanceID},
				},
			}); err != nil {
				lock.Unlock()
				return err
			}
		}
		pubs, err := e.continueParentAfterUnfinishedCallChild(ctx, parent, callerDep, lock)
		lock.Unlock()
		if err != nil {
			return err
		}
		return e.flushPublications(ctx, pubs)
	}

	var outVars map[string]any
	if child != nil {
		if call, ok := callerDep.CallActivitySpec(p.CallActivityID); ok && len(call.Outputs) > 0 {
			outVars = anyMapFromJSONStrings(applyMappings(child.Variables, call.Outputs))
		}
	}
	return e.Complete(ctx, p.ParentInstanceID, p.CallActivityID, p.HostTokenID, outVars)
}

func (e *Engine) terminateCalledInstance(ctx context.Context, childInstanceID string) error {
	if childInstanceID == "" {
		return nil
	}
	e.mu.Lock()
	inst := e.instances[childInstanceID]
	lock := e.instMu[childInstanceID]
	var dep *deploy.Deployment
	if inst != nil {
		dep = e.deployments[inst.DeploymentID]
	}
	e.mu.Unlock()
	if inst == nil || dep == nil || inst.Status != projection.StatusActive {
		return nil
	}

	lock.Lock()
	cmdID, err := NextID()
	if err != nil {
		lock.Unlock()
		return err
	}
	pid := inst.ProcessID
	if pid == "" {
		pid = dep.ProcessID()
	}
	cmd := &eventv1.Event{
		Id:                cmdID,
		Timestamp:         nowMillis(),
		RecordType:        eventv1.Event_RECORD_TYPE_COMMAND,
		DeploymentId:      inst.DeploymentID,
		ProcessInstanceId: childInstanceID,
		ProcessVersion:    inst.Version,
		Element: &eventv1.Element{
			Intent: eventv1.Element_INTENT_TERMINATING,
			Type:   eventv1.Element_TYPE_PROCESS,
			Id:     pid,
		},
	}
	if _, err := e.log.Append(ctx, cmd); err != nil {
		lock.Unlock()
		return err
	}
	emit := e.emitter(ctx, inst, cmdID)
	savedParent := inst.ParentProcessInstanceID
	inst.ParentProcessInstanceID = "" // suppress ResumeParent
	err = e.executor.terminateInstance(ctx, dep, inst, emit)
	inst.ParentProcessInstanceID = savedParent
	lock.Unlock()
	return err
}

func applyMappings(src map[string]string, maps []deploy.VariableMapping) map[string]string {
	out := make(map[string]string, len(maps))
	for _, m := range maps {
		if v, ok := src[m.Source]; ok {
			out[m.Target] = v
		}
	}
	return out
}

func variablesFromJSONStrings(m map[string]string) []*eventv1.Variable {
	if len(m) == 0 {
		return nil
	}
	out := make([]*eventv1.Variable, 0, len(m))
	for k, v := range m {
		out = append(out, &eventv1.Variable{Name: k, JsonValue: v})
	}
	return out
}

func anyMapFromJSONStrings(m map[string]string) map[string]any {
	if len(m) == 0 {
		return nil
	}
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = json.RawMessage(v)
	}
	return out
}

func terminateChildPub(tok *projection.Token, callActivityID string) *handlers.Publication {
	if tok == nil || tok.CalledProcessInstanceID == "" {
		return nil
	}
	return &handlers.Publication{
		Kind:            handlers.PublicationTerminateChild,
		ChildInstanceID: tok.CalledProcessInstanceID,
		CallActivityID:  callActivityID,
		HostTokenID:     tok.ID,
	}
}

// continueParentAfterUnfinishedCallChild decrements WaitingChildren and advances
// parent compensation when all unfinished Call Activity children are done.
// Caller must hold parent instance lock.
func (e *Engine) continueParentAfterUnfinishedCallChild(
	ctx context.Context,
	parent *projection.Instance,
	callerDep *deploy.Deployment,
	_ *sync.Mutex,
) ([]handlers.Publication, error) {
	if parent == nil || callerDep == nil {
		return nil, nil
	}
	if parent.PendingCompensation == nil {
		parent.PendingCompensation = rebuildPendingCompensation(parent, callerDep)
	}
	pc := parent.PendingCompensation
	if pc == nil {
		return nil, nil
	}
	if pc.WaitingChildren > 0 {
		pc.WaitingChildren--
	}
	if pc.WaitingChildren > 0 {
		return nil, nil
	}
	cmdID, err := NextID()
	if err != nil {
		return nil, err
	}
	cmd := &eventv1.Event{
		Id:                cmdID,
		Timestamp:         nowMillis(),
		RecordType:        eventv1.Event_RECORD_TYPE_COMMAND,
		DeploymentId:      parent.DeploymentID,
		ProcessInstanceId: parent.ID,
		ProcessVersion:    parent.Version,
		Element: &eventv1.Element{
			Intent:  eventv1.Element_INTENT_COMPLETING,
			Type:    eventv1.Element_TYPE_INTERMEDIATE_THROW_EVENT,
			Id:      pc.ThrowElementID,
			TokenId: pc.ThrowTokenID,
		},
	}
	if _, err := e.log.Append(ctx, cmd); err != nil {
		return nil, err
	}
	emit := e.emitter(ctx, parent, cmdID)
	return e.executor.advanceCompensation(ctx, callerDep, parent, emit)
}

func (e *Engine) compensateUnfinishedChild(ctx context.Context, p handlers.Publication) error {
	childID := p.ChildInstanceID
	if childID == "" {
		return e.resumeParentCall(ctx, handlers.Publication{
			Kind:             handlers.PublicationResumeParent,
			ParentInstanceID: p.ParentInstanceID,
			CallActivityID:   p.CallActivityID,
			HostTokenID:      p.HostTokenID,
			ChildInstanceID:  childID,
			Completed:        false,
		})
	}
	e.mu.Lock()
	inst := e.instances[childID]
	lock := e.instMu[childID]
	var dep *deploy.Deployment
	if inst != nil {
		dep = e.deployments[inst.DeploymentID]
	}
	e.mu.Unlock()
	if inst == nil || dep == nil || lock == nil || inst.Status != projection.StatusActive {
		return e.resumeParentCall(ctx, handlers.Publication{
			Kind:             handlers.PublicationResumeParent,
			ParentInstanceID: p.ParentInstanceID,
			CallActivityID:   p.CallActivityID,
			HostTokenID:      p.HostTokenID,
			ChildInstanceID:  childID,
			Completed:        false,
		})
	}

	lock.Lock()
	cmdID, err := NextID()
	if err != nil {
		lock.Unlock()
		return err
	}
	pid := inst.ProcessID
	if pid == "" {
		pid = dep.ProcessID()
	}
	cmd := &eventv1.Event{
		Id:                cmdID,
		Timestamp:         nowMillis(),
		RecordType:        eventv1.Event_RECORD_TYPE_COMMAND,
		DeploymentId:      inst.DeploymentID,
		ProcessInstanceId: childID,
		ProcessVersion:    inst.Version,
		Element: &eventv1.Element{
			Intent: eventv1.Element_INTENT_TERMINATING,
			Type:   eventv1.Element_TYPE_PROCESS,
			Id:     pid,
		},
	}
	if _, err := e.log.Append(ctx, cmd); err != nil {
		lock.Unlock()
		return err
	}
	emit := e.emitter(ctx, inst, cmdID)

	if err := terminateScopeTokens(dep, inst, pid, emit, scopeTerminateOpts{
		IncludeHost: false,
		DropTokens:  true,
	}); err != nil {
		lock.Unlock()
		return err
	}
	inst.RemoveScopeBoundariesForScope(pid)
	if err := emitEventSubProcessStartDisarmInScope(dep, pid, inst, emit); err != nil {
		lock.Unlock()
		return err
	}

	type item struct {
		boundaryID string
		handlerID  string
		seq        int64
	}
	var items []item
	for bid, sub := range inst.CompensationSubs {
		if sub == nil {
			continue
		}
		c, ok := dep.CompensationByBoundary(bid)
		if !ok {
			continue
		}
		items = append(items, item{boundaryID: bid, handlerID: c.HandlerID, seq: sub.Seq})
	}
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].seq > items[i].seq {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
	queue := make([]string, 0, len(items))
	consumed := make([]string, 0, len(items))
	for _, it := range items {
		queue = append(queue, it.handlerID)
		consumed = append(consumed, it.boundaryID)
	}
	for _, bid := range consumed {
		if err := emit(&eventv1.Element{
			Intent: eventv1.Element_INTENT_COMPLETING,
			Type:   eventv1.Element_TYPE_BOUNDARY_EVENT,
			Id:     bid,
		}); err != nil {
			lock.Unlock()
			return err
		}
		if err := emit(&eventv1.Element{
			Intent: eventv1.Element_INTENT_COMPLETED,
			Type:   eventv1.Element_TYPE_BOUNDARY_EVENT,
			Id:     bid,
		}); err != nil {
			lock.Unlock()
			return err
		}
	}

	inst.PendingCompensation = &projection.PendingCompensation{
		Queue:                      queue,
		Consumed:                   consumed,
		NotifyParentUnfinishedCall: true,
	}
	pubs, err := e.executor.advanceCompensation(ctx, dep, inst, emit)
	lock.Unlock()
	if err != nil {
		return err
	}
	return e.flushPublications(ctx, pubs)
}
