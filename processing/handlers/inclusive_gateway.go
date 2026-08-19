package handlers

import (
	"fmt"

	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type InclusiveGatewayHandler struct{}

func (InclusiveGatewayHandler) Type() eventv1.Element_Type {
	return eventv1.Element_TYPE_INCLUSIVE_GATEWAY
}

func (InclusiveGatewayHandler) OnEnter(in EnterInput) (*Effect, error) {
	incoming := in.Deployment.Incoming(in.ElementID)
	outgoing := in.Deployment.Outgoing(in.ElementID)
	if len(incoming) > 1 {
		return inclusiveJoinEnter(in)
	}
	if len(outgoing) > 1 {
		return inclusiveSplit(in)
	}
	if len(outgoing) == 0 {
		return nil, fmt.Errorf("NO_OUTGOING_FLOW: inclusive gateway %q", in.ElementID)
	}
	return &Effect{
		Records:        InstantLifecycle(in.Type, in.ElementID, in.TokenID, nil),
		TakeOutgoing:   true,
		OutgoingFlowID: outgoing[0],
	}, nil
}

func inclusiveSplit(in EnterInput) (*Effect, error) {
	var vars map[string]string
	if in.Instance != nil {
		vars = in.Instance.Variables
	}
	flows, err := in.Deployment.ChooseInclusiveOutgoing(in.ElementID, vars)
	if err != nil {
		return nil, err
	}
	if len(flows) == 1 {
		return &Effect{
			Records:        InstantLifecycle(in.Type, in.ElementID, in.TokenID, nil),
			TakeOutgoing:   true,
			OutgoingFlowID: flows[0],
		}, nil
	}
	return &Effect{
		Records: InstantLifecycle(in.Type, in.ElementID, in.TokenID, nil),
		Fork:    flows,
	}, nil
}

// inclusiveJoinEnter waits for all tokens that were actually forked to arrive.
// It counts how many incoming sequence flows have an active/waiting token
// somewhere upstream that could still reach this gateway.
// Simplified approach: wait for all tokens currently waiting at this element
// to equal the number of incoming flows that have been "activated" (i.e. have
// a token that arrived or is on the way). We use the same strategy as parallel
// join: count arrived tokens. The difference is that we only need to wait for
// tokens that actually exist (were forked), not ALL incoming edges.
func inclusiveJoinEnter(in EnterInput) (*Effect, error) {
	incoming := in.Deployment.Incoming(in.ElementID)
	arrived := waitingAtElement(in.Instance, in.ElementID) + 1
	records := []*eventv1.Element{
		{Intent: eventv1.Element_INTENT_ACTIVATING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
		{Intent: eventv1.Element_INTENT_ACTIVATED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
	}

	expected := countExpectedIncoming(in.Instance, in.ElementID, in.TokenID, incoming)

	if arrived < expected {
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

// countExpectedIncoming determines how many tokens the join should wait for.
// Strategy: total active tokens in the instance (status Active or Waiting)
// that are NOT already waiting at this join element. Each such token represents
// one incoming path. But we cap at len(incoming) since that's the max.
// If the only active tokens are all at this join, we're ready.
// countExpectedIncoming determines how many tokens the join should wait for.
// It counts tokens already at the join + the arriving token + any other active
// tokens elsewhere (excluding the arriving token). Capped at len(incoming).
func countExpectedIncoming(inst *projection.Instance, joinID, arrivingTokenID string, incoming []string) int {
	if inst == nil {
		return 1
	}
	activeElsewhere := 0
	atJoin := 0
	for tid, tok := range inst.Tokens {
		if tok == nil || tid == arrivingTokenID {
			continue
		}
		switch tok.Status {
		case projection.TokenActive, projection.TokenWaiting:
			if tok.ElementID == joinID {
				atJoin++
			} else {
				activeElsewhere++
			}
		}
	}
	expected := atJoin + 1 + activeElsewhere
	if expected > len(incoming) {
		expected = len(incoming)
	}
	return expected
}

func (InclusiveGatewayHandler) OnComplete(CompleteInput) (*Effect, error) {
	return nil, errUnsupportedComplete(eventv1.Element_TYPE_INCLUSIVE_GATEWAY)
}
