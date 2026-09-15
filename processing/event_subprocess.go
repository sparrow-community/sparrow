package processing

import (
	"context"
	"fmt"
	"time"

	"github.com/sparrow-community/sparrow/processing/deploy"
	"github.com/sparrow-community/sparrow/processing/handlers"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func eventSubProcessRunning(dep *deploy.Deployment, inst *projection.Instance, eventSubProcessID string) bool {
	for _, tok := range inst.Tokens {
		if tok == nil {
			continue
		}
		if tok.ElementID == eventSubProcessID {
			return true
		}
		scope, _ := dep.ScopeOf(tok.ElementID)
		if scope == eventSubProcessID || isInScope(dep, scope, eventSubProcessID) {
			return true
		}
	}
	return false
}

func emitEventSubProcessStartArms(
	dep *deploy.Deployment,
	inst *projection.Instance,
	scopeID string,
	now time.Time,
	emit Emitter,
) error {
	if dep == nil || inst == nil || emit == nil {
		return nil
	}
	for _, spec := range dep.EventSubProcessesInScope(scopeID) {
		if spec.Kind == deploy.CatchKindCompensate {
			continue
		}
		if eventSubProcessRunning(dep, inst, spec.ID) {
			continue
		}
		if _, ok := inst.EventSubProcesses[spec.ID]; ok {
			continue
		}
		if err := emit(eventSubProcessStartActivated(dep, spec, now)); err != nil {
			return err
		}
	}
	return nil
}

func eventSubProcessStartActivated(dep *deploy.Deployment, spec deploy.EventSubProcess, now time.Time) *eventv1.Element {
	payload := &eventv1.EventPayload{EventSubProcessElementId: spec.ID}
	switch spec.Kind {
	case deploy.CatchKindMessage:
		payload.MessageName = spec.MessageName
	case deploy.CatchKindSignal:
		payload.SignalName = spec.SignalName
	case deploy.CatchKindError:
		payload.ErrorCode = spec.ErrorCode
	case deploy.CatchKindTimer:
		if now.IsZero() {
			now = time.Now()
		}
		if due, text, err := dep.TimerDue(spec.StartEventID, now); err == nil {
			payload.DueUnixMs = due
			payload.Duration = text
		}
	case deploy.CatchKindConditional:
		payload.Duration = spec.Condition
		payload.TokenWait = true
	}
	return &eventv1.Element{
		Intent:  eventv1.Element_INTENT_ACTIVATED,
		Type:    eventv1.Element_TYPE_START_EVENT,
		Id:      spec.StartEventID,
		Payload: &eventv1.Element_EventPayload{EventPayload: payload},
	}
}

func emitEventSubProcessStartDisarmInScope(
	dep *deploy.Deployment,
	scopeID string,
	inst *projection.Instance,
	emit Emitter,
) error {
	for eventSubProcessID := range inst.EventSubProcesses {
		spec, ok := dep.EventSubProcessSpec(eventSubProcessID)
		if !ok || spec.ParentScopeID != scopeID {
			continue
		}
		if err := emitEventSubProcessStartDisarm(dep, eventSubProcessID, emit); err != nil {
			return err
		}
	}
	return nil
}

func emitEventSubProcessStartDisarm(dep *deploy.Deployment, eventSubProcessElementID string, emit Emitter) error {
	spec, ok := dep.EventSubProcessSpec(eventSubProcessElementID)
	if !ok {
		return nil
	}
	ep := &eventv1.Element_EventPayload{
		EventPayload: &eventv1.EventPayload{EventSubProcessElementId: eventSubProcessElementID},
	}
	for _, intent := range []eventv1.Element_Intent{
		eventv1.Element_INTENT_TERMINATING,
		eventv1.Element_INTENT_TERMINATED,
	} {
		if err := emit(&eventv1.Element{
			Intent:  intent,
			Type:    eventv1.Element_TYPE_START_EVENT,
			Id:      spec.StartEventID,
			Payload: ep,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) triggerEventSubProcess(ctx context.Context, instanceID, eventSubProcessElementID string, vars map[string]any) error {
	e.mu.Lock()
	inst := e.instances[instanceID]
	lock := e.instMu[instanceID]
	var dep *deploy.Deployment
	if inst != nil {
		dep = e.deployments[inst.DeploymentID]
	}
	e.mu.Unlock()
	if inst == nil || dep == nil || lock == nil {
		return fmt.Errorf("NOT_FOUND: instance %q", instanceID)
	}

	lock.Lock()
	pubs, err := e.triggerEventSubProcessLocked(ctx, dep, inst, instanceID, eventSubProcessElementID, vars)
	lock.Unlock()
	if err != nil {
		return err
	}
	if err := e.flushPublications(ctx, pubs); err != nil {
		return err
	}
	return e.tryDeliverBuffered(ctx, instanceID)
}

func (e *Engine) triggerEventSubProcessLocked(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	instanceID, eventSubProcessElementID string,
	vars map[string]any,
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
	cmdID, err := NextID()
	if err != nil {
		return nil, err
	}
	pv, err := projection.VariablesFromMap(vars)
	if err != nil {
		return nil, err
	}

	cmd := &eventv1.Event{
		Id:                cmdID,
		Timestamp:         nowMillis(),
		RecordType:        eventv1.Event_RECORD_TYPE_COMMAND,
		DeploymentId:      inst.DeploymentID,
		ProcessInstanceId: instanceID,
		ProcessVersion:    inst.Version,
		Element: &eventv1.Element{
			Intent:  eventv1.Element_INTENT_ACTIVATING,
			Type:    eventv1.Element_TYPE_SUB_PROCESS,
			Id:      eventSubProcessElementID,
			TokenId: tokenID,
		},
	}
	if len(pv) > 0 {
		cmd.Element.Payload = &eventv1.Element_ActivityPayload{
			ActivityPayload: &eventv1.ActivityPayload{Variables: pv},
		}
	}
	if _, err := e.log.Append(ctx, cmd); err != nil {
		return nil, err
	}
	emit := e.emitter(ctx, inst, cmdID)

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
	}

	return e.executor.Enter(ctx, dep, inst, tokenID, eventSubProcessElementID, emit)
}

func (e *Engine) collectEventSubProcessMessageArms(name, instanceID string, keys []*eventv1.Variable) []eventSubProcessWait {
	return e.collectEventSubProcessArms(instanceID, keys, func(_ string, arm *projection.EventSubProcessArm) bool {
		return arm.MessageName != "" && arm.MessageName == name
	})
}

func (e *Engine) collectEventSubProcessSignalArms(name, instanceID string) []eventSubProcessWait {
	return e.collectEventSubProcessArms(instanceID, nil, func(_ string, arm *projection.EventSubProcessArm) bool {
		return arm.SignalName != "" && arm.SignalName == name
	})
}

func (e *Engine) collectEventSubProcessTimerDue(nowUnixMs int64) []eventSubProcessWait {
	e.mu.Lock()
	ids := make([]string, 0, len(e.instances))
	for id := range e.instances {
		ids = append(ids, id)
	}
	e.mu.Unlock()

	var due []eventSubProcessWait
	for _, iid := range ids {
		e.mu.Lock()
		inst := e.instances[iid]
		lock := e.instMu[iid]
		e.mu.Unlock()
		if inst == nil || lock == nil {
			continue
		}
		lock.Lock()
		for eventSubProcessID, arm := range inst.EventSubProcesses {
			if arm != nil && arm.DueUnixMs > 0 && arm.DueUnixMs <= nowUnixMs {
				due = append(due, eventSubProcessWait{instanceID: iid, eventSubProcessElementID: eventSubProcessID})
			}
		}
		lock.Unlock()
	}
	return due
}

type eventSubProcessWait struct {
	instanceID               string
	eventSubProcessElementID string
}

func (e *Engine) collectEventSubProcessArms(
	instanceID string,
	keys []*eventv1.Variable,
	match func(eventSubProcessElementID string, arm *projection.EventSubProcessArm) bool,
) []eventSubProcessWait {
	e.mu.Lock()
	ids := make([]string, 0, len(e.instances))
	if instanceID != "" {
		if _, ok := e.instances[instanceID]; ok {
			ids = append(ids, instanceID)
		}
	} else {
		for id := range e.instances {
			ids = append(ids, id)
		}
	}
	e.mu.Unlock()

	var out []eventSubProcessWait
	for _, iid := range ids {
		e.mu.Lock()
		inst := e.instances[iid]
		lock := e.instMu[iid]
		e.mu.Unlock()
		if inst == nil || lock == nil {
			continue
		}
		lock.Lock()
		if keys != nil && !correlationKeysMatch(inst.Variables, keys) {
			lock.Unlock()
			continue
		}
		for eventSubProcessID, arm := range inst.EventSubProcesses {
			if arm != nil && match(eventSubProcessID, arm) {
				out = append(out, eventSubProcessWait{instanceID: iid, eventSubProcessElementID: eventSubProcessID})
			}
		}
		lock.Unlock()
	}
	return out
}
