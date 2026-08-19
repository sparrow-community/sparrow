package handlers

import (
	"fmt"

	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type ParallelGatewayHandler struct{}

func (ParallelGatewayHandler) Type() eventv1.Element_Type {
	return eventv1.Element_TYPE_PARALLEL_GATEWAY
}

func (ParallelGatewayHandler) OnEnter(in EnterInput) (*Effect, error) {
	incoming := in.Deployment.Incoming(in.ElementID)
	outgoing := in.Deployment.Outgoing(in.ElementID)
	if len(incoming) > 1 {
		return parallelJoinEnter(in)
	}
	if len(outgoing) > 1 {
		return &Effect{
			Records: InstantLifecycle(in.Type, in.ElementID, in.TokenID, nil),
			Fork:    append([]string{}, outgoing...),
		}, nil
	}
	if len(outgoing) == 0 {
		return nil, fmt.Errorf("NO_OUTGOING_FLOW: parallel gateway %q", in.ElementID)
	}
	return &Effect{
		Records:        InstantLifecycle(in.Type, in.ElementID, in.TokenID, nil),
		TakeOutgoing:   true,
		OutgoingFlowID: outgoing[0],
	}, nil
}

func parallelJoinEnter(in EnterInput) (*Effect, error) {
	need := len(in.Deployment.Incoming(in.ElementID))
	arrived := waitingAtElement(in.Instance, in.ElementID) + 1
	records := []*eventv1.Element{
		&eventv1.Element{Intent: eventv1.Element_INTENT_ACTIVATING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
		&eventv1.Element{Intent: eventv1.Element_INTENT_ACTIVATED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
	}
	if arrived < need {
		return &Effect{Records: records, Wait: true}, nil
	}
	records = append(records,
		&eventv1.Element{Intent: eventv1.Element_INTENT_COMPLETING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
		&eventv1.Element{Intent: eventv1.Element_INTENT_COMPLETED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
	)
	return &Effect{
		Records:            records,
		TakeOutgoing:       true,
		TerminateJoinPeers: in.ElementID,
	}, nil
}

func waitingAtElement(inst *projection.Instance, elementID string) int {
	if inst == nil {
		return 0
	}
	n := 0
	for _, tok := range inst.Tokens {
		if tok != nil && tok.ElementID == elementID && tok.Status == projection.TokenWaiting {
			n++
		}
	}
	return n
}

func (ParallelGatewayHandler) OnComplete(CompleteInput) (*Effect, error) {
	return nil, errUnsupportedComplete(eventv1.Element_TYPE_PARALLEL_GATEWAY)
}
