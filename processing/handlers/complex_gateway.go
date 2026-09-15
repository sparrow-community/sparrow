package handlers

import (
	"fmt"

	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type ComplexGatewayHandler struct{}

func (ComplexGatewayHandler) Type() eventv1.Element_Type {
	return eventv1.Element_TYPE_COMPLEX_GATEWAY
}

func (ComplexGatewayHandler) OnEnter(in EnterInput) (*Effect, error) {
	incoming := in.Deployment.Incoming(in.ElementID)
	outgoing := in.Deployment.Outgoing(in.ElementID)
	if len(incoming) > 1 {
		return complexJoinEnter(in, incoming, outgoing)
	}
	if len(outgoing) > 1 {
		return complexSplit(in)
	}
	if len(outgoing) == 0 {
		return nil, fmt.Errorf("NO_OUTGOING_FLOW: complex gateway %q", in.ElementID)
	}
	return &Effect{
		Records:        InstantLifecycle(in.Type, in.ElementID, in.TokenID, nil),
		TakeOutgoing:   true,
		OutgoingFlowID: outgoing[0],
	}, nil
}

func complexSplit(in EnterInput) (*Effect, error) {
	var vars map[string]string
	if in.Instance != nil {
		vars = in.Instance.Variables
	}
	flows, err := in.Deployment.ChooseComplexOutgoing(in.ElementID, vars)
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

func complexJoinEnter(in EnterInput, incoming, outgoing []string) (*Effect, error) {
	arrived := waitingAtElement(in.Instance, in.ElementID) + 1
	records := []*eventv1.Element{
		{Intent: eventv1.Element_INTENT_ACTIVATING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
		{Intent: eventv1.Element_INTENT_ACTIVATED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
	}
	var vars map[string]string
	if in.Instance != nil {
		vars = in.Instance.Variables
	}
	fire, err := in.Deployment.EvalComplexActivation(in.ElementID, arrived, len(incoming), vars)
	if err != nil {
		return nil, err
	}
	if !fire {
		return &Effect{Records: records, Wait: true}, nil
	}
	records = append(records,
		&eventv1.Element{Intent: eventv1.Element_INTENT_COMPLETING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
		&eventv1.Element{Intent: eventv1.Element_INTENT_COMPLETED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
	)
	effect := &Effect{
		Records:            records,
		TerminateJoinPeers: in.ElementID,
	}
	if len(outgoing) == 0 {
		return nil, fmt.Errorf("NO_OUTGOING_FLOW: complex gateway %q", in.ElementID)
	}
	if len(outgoing) == 1 {
		effect.TakeOutgoing = true
		effect.OutgoingFlowID = outgoing[0]
		return effect, nil
	}
	flows, err := in.Deployment.ChooseComplexOutgoing(in.ElementID, vars)
	if err != nil {
		return nil, err
	}
	if len(flows) == 1 {
		effect.TakeOutgoing = true
		effect.OutgoingFlowID = flows[0]
		return effect, nil
	}
	effect.Fork = flows
	return effect, nil
}

func (ComplexGatewayHandler) OnComplete(CompleteInput) (*Effect, error) {
	return nil, errUnsupportedComplete(eventv1.Element_TYPE_COMPLEX_GATEWAY)
}
