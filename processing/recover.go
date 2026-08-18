package processing

import (
	"context"
	"fmt"

	"github.com/sparrow-community/sparrow/processing/deploy"
	"github.com/sparrow-community/sparrow/processing/handlers"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type cmdFollow struct {
	event    bool
	rejected bool
}

func (e *Engine) redrive(ctx context.Context) error {
	records, err := e.log.ReadAll(ctx)
	if err != nil {
		return fmt.Errorf("redrive read log: %w", err)
	}
	follow := make(map[string]*cmdFollow)
	cmds := make([]*eventv1.Event, 0)
	for _, rec := range records {
		switch rec.GetRecordType() {
		case eventv1.Event_RECORD_TYPE_COMMAND:
			cmds = append(cmds, rec)
			if follow[rec.GetId()] == nil {
				follow[rec.GetId()] = &cmdFollow{}
			}
		case eventv1.Event_RECORD_TYPE_EVENT:
			id := rec.GetSourceRecordId()
			if follow[id] == nil {
				follow[id] = &cmdFollow{}
			}
			follow[id].event = true
		case eventv1.Event_RECORD_TYPE_REJECTION:
			id := rec.GetSourceRecordId()
			if follow[id] == nil {
				follow[id] = &cmdFollow{}
			}
			follow[id].rejected = true
		}
	}
	for _, cmd := range cmds {
		st := follow[cmd.GetId()]
		if st != nil && st.rejected {
			continue
		}
		hasEvent := st != nil && st.event
		if err := e.maybeRedrive(ctx, cmd, hasEvent); err != nil {
			return fmt.Errorf("redrive command %s: %w", cmd.GetId(), err)
		}
	}
	return nil
}

func (e *Engine) maybeRedrive(ctx context.Context, cmd *eventv1.Event, hasEvent bool) error {
	iid := cmd.GetProcessInstanceId()
	e.mu.Lock()
	inst := e.instances[iid]
	lock := e.instMu[iid]
	dep := e.deployments[cmd.GetDeploymentId()]
	e.mu.Unlock()
	if inst == nil || lock == nil {
		return fmt.Errorf("NOT_FOUND: instance %q", iid)
	}
	if dep == nil {
		return fmt.Errorf("NOT_FOUND: deployment %q", cmd.GetDeploymentId())
	}

	lock.Lock()
	defer lock.Unlock()
	if hasEvent && instanceStable(inst) && !completeCommandUnfinished(inst, cmd) {
		return nil
	}
	return e.redriveCommand(ctx, dep, inst, cmd)
}

func completeCommandUnfinished(inst *projection.Instance, cmd *eventv1.Event) bool {
	el := cmd.GetElement()
	if el == nil || el.GetIntent() != eventv1.Element_INTENT_COMPLETING {
		return false
	}
	tok := inst.Tokens[el.GetTokenId()]
	return tok != nil && tok.ElementID == el.GetId() && tok.Status == projection.TokenWaiting
}

func instanceStable(inst *projection.Instance) bool {
	if inst == nil {
		return false
	}
	if inst.Status == projection.StatusCompleted || inst.Status == projection.StatusTerminated {
		return true
	}
	for _, tok := range inst.Tokens {
		if tok != nil && tok.Status == projection.TokenWaiting {
			return true
		}
	}
	return false
}

func (e *Engine) redriveCommand(ctx context.Context, dep *deploy.Deployment, inst *projection.Instance, cmd *eventv1.Event) error {
	el := cmd.GetElement()
	if el == nil {
		return fmt.Errorf("command %s has no element", cmd.GetId())
	}
	emit := e.emitter(ctx, inst, cmd.GetId())
	switch el.GetIntent() {
	case eventv1.Element_INTENT_ACTIVATING:
		if el.GetType() != eventv1.Element_TYPE_PROCESS {
			return fmt.Errorf("unsupported activating command type %v", el.GetType())
		}
		return e.redriveCreateInstance(ctx, dep, inst, cmd, emit)
	case eventv1.Element_INTENT_COMPLETING:
		return e.redriveComplete(ctx, dep, inst, cmd, emit)
	case eventv1.Element_INTENT_FAILED:
		return e.redriveFail(inst, cmd, emit)
	default:
		return fmt.Errorf("unsupported command intent %v", el.GetIntent())
	}
}

func (e *Engine) redriveCreateInstance(ctx context.Context, dep *deploy.Deployment, inst *projection.Instance, cmd *eventv1.Event, emit Emitter) error {
	startID, err := dep.StartEventID()
	if err != nil {
		return err
	}
	tokenID := cmd.GetElement().GetTokenId()
	if tokenID == "" {
		for _, tok := range inst.Tokens {
			if tok != nil {
				tokenID = tok.ID
				break
			}
		}
	}
	if tokenID == "" {
		tokenID, err = NextID()
		if err != nil {
			return err
		}
	}
	var pv []*eventv1.Variable
	if p := cmd.GetElement().GetProcessPayload(); p != nil {
		pv = p.GetVariables()
	}
	for _, rec := range handlers.ProcessStartRecords(dep.ProcessID(), pv) {
		if err := emit(rec); err != nil {
			return err
		}
	}
	return e.executor.Enter(ctx, dep, inst, tokenID, startID, emit)
}

func (e *Engine) redriveComplete(ctx context.Context, dep *deploy.Deployment, inst *projection.Instance, cmd *eventv1.Event, emit Emitter) error {
	el := cmd.GetElement()
	var pv []*eventv1.Variable
	if p := el.GetActivityPayload(); p != nil {
		pv = p.GetVariables()
	}
	if err := e.executor.Complete(ctx, dep, inst, el.GetTokenId(), el.GetId(), pv, emit); err != nil {
		return err
	}
	e.releaseLease(inst.ID, el.GetTokenId())
	return nil
}

func (e *Engine) redriveFail(inst *projection.Instance, cmd *eventv1.Event, emit Emitter) error {
	el := cmd.GetElement()
	failed := &eventv1.Element{
		Intent:  eventv1.Element_INTENT_FAILED,
		Type:    el.GetType(),
		Id:      el.GetId(),
		TokenId: el.GetTokenId(),
		Payload: el.GetPayload(),
	}
	if err := emit(failed); err != nil {
		return err
	}
	e.releaseLease(inst.ID, el.GetTokenId())
	e.notifyJobs()
	return nil
}
