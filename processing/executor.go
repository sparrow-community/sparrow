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

// Emitter appends an Element behavior record as an EVENT and updates projection.
type Emitter func(el *eventv1.Element) error

// Executor drives token movement using ElementHandler semantics over bpmn definitions.
type Executor struct {
	Handlers *handlers.Registry
	Now      func() time.Time
}

func NewExecutor(reg *handlers.Registry) *Executor {
	if reg == nil {
		reg = handlers.DefaultRegistry()
	}
	return &Executor{Handlers: reg}
}

func (x *Executor) now() time.Time {
	if x != nil && x.Now != nil {
		return x.Now()
	}
	return time.Now()
}

func (x *Executor) Enter(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	tokenID, elementID string,
	emit Emitter,
) error {
	for {
		typ, err := dep.TypeOf(elementID)
		if err != nil {
			return fmt.Errorf("UNSUPPORTED_ELEMENT: %v", err)
		}

		h, err := x.Handlers.Get(typ)
		if err != nil {
			return err
		}
		effect, err := h.OnEnter(handlers.EnterInput{
			Deployment: dep,
			Instance:   inst,
			ElementID:  elementID,
			Type:       typ,
			TokenID:    tokenID,
			Now:        x.now(),
		})
		if err != nil {
			return err
		}
		for _, rec := range effect.Records {
			if err := emit(rec); err != nil {
				return err
			}
		}
		if effect.EnterChild != "" {
			elementID = effect.EnterChild
			continue
		}
		if len(effect.Fork) > 0 {
			for i, flowID := range effect.Fork {
				tid := tokenID
				if i > 0 {
					var err error
					tid, err = NextID()
					if err != nil {
						return err
					}
				}
				next, err := x.takeOutgoing(dep, tid, elementID, flowID, emit)
				if err != nil {
					return err
				}
				if err := x.Enter(ctx, dep, inst, tid, next, emit); err != nil {
					return err
				}
			}
			return nil
		}
		if effect.Wait {
			return nil
		}
		if effect.TryCompleteProcess {
			return x.tryCompleteScope(ctx, dep, inst, tokenID, elementID, emit)
		}
		if !effect.TakeOutgoing {
			return nil
		}
		if effect.TerminateJoinPeers != "" {
			if err := x.terminateJoinPeers(inst, tokenID, effect.TerminateJoinPeers, emit); err != nil {
				return err
			}
		}

		next, err := x.takeOutgoing(dep, tokenID, elementID, effect.OutgoingFlowID, emit)
		if err != nil {
			return err
		}
		elementID = next
	}
}

func (x *Executor) Complete(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	tokenID, elementID string,
	vars []*eventv1.Variable,
	emit Emitter,
) error {
	typ, err := dep.TypeOf(elementID)
	if err != nil {
		return err
	}
	h, err := x.Handlers.Get(typ)
	if err != nil {
		return err
	}
	effect, err := h.OnComplete(handlers.CompleteInput{
		Deployment: dep,
		Instance:   inst,
		Token:      inst.Tokens[tokenID],
		ElementID:  elementID,
		Type:       typ,
		TokenID:    tokenID,
		Variables:  vars,
	})
	if err != nil {
		return err
	}
	for _, rec := range effect.Records {
		if err := emit(rec); err != nil {
			return err
		}
	}
	if effect.SpawnOutgoing != nil {
		spawn := effect.SpawnOutgoing
		spawnID, err := NextID()
		if err != nil {
			return err
		}
		for _, rec := range handlers.InstantLifecycle(spawn.Type, spawn.ElementID, spawnID, nil) {
			if err := emit(rec); err != nil {
				return err
			}
		}
		next, err := x.takeOutgoing(dep, spawnID, spawn.ElementID, spawn.OutgoingFlowID, emit)
		if err != nil {
			return err
		}
		if err := x.Enter(ctx, dep, inst, spawnID, next, emit); err != nil {
			return err
		}
		return x.tryCompleteProcessScope(dep, inst, emit)
	}
	if effect.Wait || !effect.TakeOutgoing {
		return nil
	}

	next, err := x.takeOutgoing(dep, tokenID, elementID, effect.OutgoingFlowID, emit)
	if err != nil {
		return err
	}
	return x.Enter(ctx, dep, inst, tokenID, next, emit)
}

func (x *Executor) terminateJoinPeers(inst *projection.Instance, survivorTokenID, joinElementID string, emit Emitter) error {
	typ := eventv1.Element_TYPE_PARALLEL_GATEWAY
	for tid, tok := range inst.Tokens {
		if tid == survivorTokenID || tok == nil {
			continue
		}
		if tok.ElementID != joinElementID || tok.Status != projection.TokenWaiting {
			continue
		}
		for _, intent := range []eventv1.Element_Intent{
			eventv1.Element_INTENT_TERMINATING,
			eventv1.Element_INTENT_TERMINATED,
		} {
			if err := emit(&eventv1.Element{
				Intent:  intent,
				Type:    typ,
				Id:      joinElementID,
				TokenId: tid,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (x *Executor) takeOutgoing(
	dep *deploy.Deployment,
	tokenID, fromElementID, preferredFlowID string,
	emit Emitter,
) (string, error) {
	flowID := preferredFlowID
	if flowID == "" {
		outs := dep.Outgoing(fromElementID)
		if len(outs) == 0 {
			return "", fmt.Errorf("NO_OUTGOING_FLOW: element %q", fromElementID)
		}
		flowID = outs[0]
	}
	flow, err := dep.SequenceFlow(flowID)
	if err != nil {
		return "", err
	}
	if err := emit(handlers.SequenceFlowTaken(flow.ID, flow.SourceRef, flow.TargetRef, tokenID)); err != nil {
		return "", err
	}
	return flow.TargetRef, nil
}

// tryCompleteScope checks if the scope containing elementID can be completed.
// If elementID is inside a SubProcess, tries to complete that SubProcess.
// If at process level, tries to complete the Process.
func (x *Executor) tryCompleteScope(ctx context.Context, dep *deploy.Deployment, inst *projection.Instance, tokenID, elementID string, emit Emitter) error {
	scopeID, _ := dep.ScopeOf(elementID)
	if scopeID == "" || scopeID == dep.ProcessID() {
		return x.tryCompleteProcessScope(dep, inst, emit)
	}
	// SubProcess scope: check all tokens in this scope are at EndEvents
	for _, tok := range inst.Tokens {
		if tok.Status == projection.TokenWaiting {
			tokScope, _ := dep.ScopeOf(tok.ElementID)
			if tokScope == scopeID {
				return nil
			}
			continue
		}
		tokScope, _ := dep.ScopeOf(tok.ElementID)
		if tokScope == scopeID {
			typ, err := dep.TypeOf(tok.ElementID)
			if err != nil || typ != eventv1.Element_TYPE_END_EVENT {
				return nil
			}
		}
	}
	// All tokens in this scope are at EndEvents; disarm scope boundaries and complete
	inst.RemoveScopeBoundariesForScope(scopeID)
	h, err := x.Handlers.Get(eventv1.Element_TYPE_SUB_PROCESS)
	if err != nil {
		return err
	}
	effect, err := h.OnComplete(handlers.CompleteInput{
		Deployment: dep,
		Instance:   inst,
		ElementID:  scopeID,
		Type:       eventv1.Element_TYPE_SUB_PROCESS,
		TokenID:    tokenID,
	})
	if err != nil {
		return err
	}
	for _, rec := range effect.Records {
		if err := emit(rec); err != nil {
			return err
		}
	}
	if effect.TakeOutgoing {
		next, err := x.takeOutgoing(dep, tokenID, scopeID, effect.OutgoingFlowID, emit)
		if err != nil {
			return err
		}
		return x.Enter(ctx, dep, inst, tokenID, next, emit)
	}
	return nil
}

// tryCompleteProcessScope checks if all tokens are at process-level EndEvents.
func (x *Executor) tryCompleteProcessScope(dep *deploy.Deployment, inst *projection.Instance, emit Emitter) error {
	processID := dep.ProcessID()
	for _, tok := range inst.Tokens {
		if tok.Status == projection.TokenWaiting {
			return nil
		}
		tokScope, _ := dep.ScopeOf(tok.ElementID)
		if tokScope != processID {
			return nil
		}
		typ, err := dep.TypeOf(tok.ElementID)
		if err != nil || typ != eventv1.Element_TYPE_END_EVENT {
			return nil
		}
	}
	h, err := x.Handlers.Get(eventv1.Element_TYPE_PROCESS)
	if err != nil {
		return err
	}
	effect, err := h.OnComplete(handlers.CompleteInput{Deployment: dep, Instance: inst})
	if err != nil {
		return err
	}
	for _, rec := range effect.Records {
		if err := emit(rec); err != nil {
			return err
		}
	}
	return nil
}
