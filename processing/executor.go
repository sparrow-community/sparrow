package processing

import (
	"context"
	"fmt"

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
}

func NewExecutor(reg *handlers.Registry) *Executor {
	if reg == nil {
		reg = handlers.DefaultRegistry()
	}
	return &Executor{Handlers: reg}
}

func (x *Executor) Enter(
	_ context.Context,
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

		tok := inst.Tokens[tokenID]
		if tok == nil {
			tok = &projection.Token{ID: tokenID}
			inst.Tokens[tokenID] = tok
		}
		tok.ElementID = elementID
		tok.Status = projection.TokenActive

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
		})
		if err != nil {
			return err
		}
		for _, rec := range effect.Records {
			if err := emit(rec); err != nil {
				return err
			}
		}
		if effect.Wait {
			return nil
		}
		if effect.TryCompleteProcess {
			return x.tryCompleteProcess(dep, inst, emit)
		}
		if !effect.TakeOutgoing {
			return nil
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
	if effect.Wait || !effect.TakeOutgoing {
		return nil
	}
	if tok := inst.Tokens[tokenID]; tok != nil {
		tok.Status = projection.TokenActive
	}

	next, err := x.takeOutgoing(dep, tokenID, elementID, effect.OutgoingFlowID, emit)
	if err != nil {
		return err
	}
	return x.Enter(ctx, dep, inst, tokenID, next, emit)
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

func (x *Executor) tryCompleteProcess(dep *deploy.Deployment, inst *projection.Instance, emit Emitter) error {
	for _, tok := range inst.Tokens {
		if tok.Status == projection.TokenWaiting {
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
