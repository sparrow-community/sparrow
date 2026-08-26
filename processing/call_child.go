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
}

func (e *Engine) startCalledInstance(ctx context.Context, p handlers.Publication) error {
	parentID := p.ParentInstanceID
	e.mu.Lock()
	dep := e.deployments[p.DeploymentID]
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
	e.mu.Unlock()
	if dep == nil {
		return fmt.Errorf("NOT_FOUND: deployment %q", p.DeploymentID)
	}
	if parent == nil {
		return fmt.Errorf("NOT_FOUND: parent instance for callActivity %q", p.CallActivityID)
	}
	call, ok := dep.CallActivitySpec(p.CallActivityID)
	if !ok {
		return fmt.Errorf("NOT_FOUND: callActivity %q", p.CallActivityID)
	}
	calledProcessID := p.CalledProcessID
	if calledProcessID == "" {
		calledProcessID = call.CalledProcessID
	}

	pv := variablesFromJSONStrings(applyMappings(parent.Variables, call.Inputs))

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
		DeploymentId:      dep.ID,
		ProcessInstanceId: instanceID,
		ProcessVersion:    dep.Version,
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

	inst := projection.NewInstance(instanceID, dep.ID, dep.Version)
	inst.ProcessID = calledProcessID
	inst.ParentProcessInstanceID = parentID
	inst.ParentElementID = p.CallActivityID
	inst.ParentTokenID = p.HostTokenID

	lock := &sync.Mutex{}
	e.mu.Lock()
	e.instances[instanceID] = inst
	e.instMu[instanceID] = lock
	e.mu.Unlock()

	lock.Lock()
	emit := e.emitter(ctx, inst, cmdID)
	for _, rec := range handlers.ProcessStartRecordsWithParent(calledProcessID, pv, parentID, p.CallActivityID, p.HostTokenID) {
		if err := emit(rec); err != nil {
			lock.Unlock()
			return err
		}
	}
	if err := emitEventSubProcessStartArms(dep, inst, calledProcessID, e.now(), emit); err != nil {
		lock.Unlock()
		return err
	}
	pubs, err := e.executor.Enter(ctx, dep, inst, tokenID, call.StartEventID, emit)
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
	var dep *deploy.Deployment
	if parent != nil {
		dep = e.deployments[parent.DeploymentID]
	}
	child := e.instances[p.ChildInstanceID]
	e.mu.Unlock()
	if parent == nil || dep == nil {
		return fmt.Errorf("NOT_FOUND: parent instance %q", p.ParentInstanceID)
	}

	if !p.Completed {
		lock.Lock()
		defer lock.Unlock()
		tok := parent.Tokens[p.HostTokenID]
		if tok == nil || tok.ElementID != p.CallActivityID || tok.Status != projection.TokenWaiting {
			return nil
		}
		cmdID, err := NextID()
		if err != nil {
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
				return err
			}
		}
		return nil
	}

	var outVars map[string]any
	if child != nil {
		if call, ok := dep.CallActivitySpec(p.CallActivityID); ok && len(call.Outputs) > 0 {
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
