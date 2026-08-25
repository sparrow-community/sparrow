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

// armEventSubProcesses registers event sub-process subscriptions for the process
// scope and any currently active embedded subProcess scopes.
// Arms are projection-only (not EventLog records); Recover re-syncs them.
func (e *Engine) armEventSubProcesses(dep *deploy.Deployment, inst *projection.Instance, now time.Time) {
	if dep == nil || inst == nil {
		return
	}
	if now.IsZero() {
		now = e.now()
	}
	armEventSubProcessesInScope(dep, inst, dep.ProcessID(), now)
	for scopeID := range activeEmbeddedScopes(dep, inst) {
		armEventSubProcessesInScope(dep, inst, scopeID, now)
	}
}

func activeEmbeddedScopes(dep *deploy.Deployment, inst *projection.Instance) map[string]bool {
	scopes := make(map[string]bool)
	for _, tok := range inst.Tokens {
		if tok == nil {
			continue
		}
		if typ, err := dep.TypeOf(tok.ElementID); err == nil &&
			typ == eventv1.Element_TYPE_SUB_PROCESS &&
			!dep.IsEventSubProcess(tok.ElementID) {
			scopes[tok.ElementID] = true
		}
		scope, ok := dep.ScopeOf(tok.ElementID)
		if !ok || scope == "" || scope == dep.ProcessID() {
			continue
		}
		if !dep.IsEventSubProcess(scope) {
			scopes[scope] = true
		}
	}
	return scopes
}

func armEventSubProcessesInScope(dep *deploy.Deployment, inst *projection.Instance, scopeID string, now time.Time) {
	for _, spec := range dep.EventSubProcessesInScope(scopeID) {
		if eventSubProcessRunning(dep, inst, spec.ID) {
			continue
		}
		if _, ok := inst.EventSubProcesses[spec.ID]; ok {
			continue
		}
		arm := &projection.EventSubProcessArm{
			SubProcessID:  spec.ID,
			StartEventID:  spec.StartEventID,
			ParentScopeID: spec.ParentScopeID,
			Interrupting:  spec.Interrupting,
			MessageName:   spec.MessageName,
			SignalName:    spec.SignalName,
		}
		if spec.Kind == deploy.CatchKindTimer {
			due, text, err := dep.TimerDue(spec.StartEventID, now)
			if err == nil {
				arm.DueUnixMs = due
				arm.TimerText = text
			}
		}
		inst.EventSubProcesses[spec.ID] = arm
	}
}

func eventSubProcessRunning(dep *deploy.Deployment, inst *projection.Instance, espID string) bool {
	for _, tok := range inst.Tokens {
		if tok == nil {
			continue
		}
		if tok.ElementID == espID {
			return true
		}
		scope, _ := dep.ScopeOf(tok.ElementID)
		if scope == espID || isInScope(dep, scope, espID) {
			return true
		}
	}
	return false
}

func (e *Engine) rearmEventSubProcesses() {
	e.mu.Lock()
	defer e.mu.Unlock()
	now := e.now()
	for _, inst := range e.instances {
		if inst == nil || inst.Status != projection.StatusActive {
			continue
		}
		dep := e.deployments[inst.DeploymentID]
		if dep == nil {
			continue
		}
		e.armEventSubProcesses(dep, inst, now)
	}
}

func (e *Engine) triggerEventSubProcess(ctx context.Context, instanceID, subProcessID string, vars map[string]any) error {
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
	pubs, err := e.triggerEventSubProcessLocked(ctx, dep, inst, instanceID, subProcessID, vars)
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
	instanceID, subProcessID string,
	vars map[string]any,
) ([]handlers.Publication, error) {
	arm := inst.EventSubProcesses[subProcessID]
	if arm == nil {
		return nil, nil
	}
	if _, ok := dep.EventSubProcessSpec(subProcessID); !ok {
		return nil, fmt.Errorf("NOT_FOUND: event subProcess %q", subProcessID)
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
			Id:      subProcessID,
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

	if arm.Interrupting {
		hostTokenID := ""
		for tid, tok := range inst.Tokens {
			if tok != nil && tokenInOrIsScope(dep, tok, arm.ParentScopeID) {
				hostTokenID = tid
				break
			}
		}
		if err := terminateScopeTokens(dep, inst, arm.ParentScopeID, emit, scopeTerminateOpts{
			IncludeHost: true,
			DropTokens:  true,
		}); err != nil {
			return nil, err
		}
		// Embedded parent scopes have no lingering host token (the entering token
		// already moved inside). Still record SubProcess TERMINATED for audit.
		if arm.ParentScopeID != dep.ProcessID() && hostTokenID != "" {
			for _, intent := range []eventv1.Element_Intent{
				eventv1.Element_INTENT_TERMINATING,
				eventv1.Element_INTENT_TERMINATED,
			} {
				if err := emit(&eventv1.Element{
					Intent:  intent,
					Type:    eventv1.Element_TYPE_SUB_PROCESS,
					Id:      arm.ParentScopeID,
					TokenId: hostTokenID,
				}); err != nil {
					return nil, err
				}
			}
			// applyToken would revive the host token on TERMINATED; drop it —
			// the entering token already lived inside the scope and must not linger.
			delete(inst.Tokens, hostTokenID)
		}
		inst.RemoveScopeBoundariesForScope(arm.ParentScopeID)
		inst.RemoveEventSubProcessesInScope(arm.ParentScopeID)
	} else {
		inst.RemoveEventSubProcess(subProcessID)
	}

	return e.executor.Enter(ctx, dep, inst, tokenID, subProcessID, emit)
}

func (e *Engine) collectESPMessageArms(name, instanceID string, keys []*eventv1.Variable) []espWait {
	return e.collectESPArms(instanceID, keys, func(arm *projection.EventSubProcessArm) bool {
		return arm.MessageName != "" && arm.MessageName == name
	})
}

func (e *Engine) collectESPSignalArms(name, instanceID string) []espWait {
	return e.collectESPArms(instanceID, nil, func(arm *projection.EventSubProcessArm) bool {
		return arm.SignalName != "" && arm.SignalName == name
	})
}

func (e *Engine) collectESPTimerDue(nowUnixMs int64) []espWait {
	e.mu.Lock()
	ids := make([]string, 0, len(e.instances))
	for id := range e.instances {
		ids = append(ids, id)
	}
	e.mu.Unlock()

	var due []espWait
	for _, iid := range ids {
		e.mu.Lock()
		inst := e.instances[iid]
		lock := e.instMu[iid]
		e.mu.Unlock()
		if inst == nil || lock == nil {
			continue
		}
		lock.Lock()
		for _, arm := range inst.EventSubProcesses {
			if arm != nil && arm.DueUnixMs > 0 && arm.DueUnixMs <= nowUnixMs {
				due = append(due, espWait{instanceID: iid, subProcessID: arm.SubProcessID})
			}
		}
		lock.Unlock()
	}
	return due
}

type espWait struct {
	instanceID   string
	subProcessID string
}

func (e *Engine) collectESPArms(instanceID string, keys []*eventv1.Variable, match func(*projection.EventSubProcessArm) bool) []espWait {
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

	var out []espWait
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
		for _, arm := range inst.EventSubProcesses {
			if arm != nil && match(arm) {
				out = append(out, espWait{instanceID: iid, subProcessID: arm.SubProcessID})
			}
		}
		lock.Unlock()
	}
	return out
}
