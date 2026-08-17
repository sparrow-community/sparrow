package handlers

import (
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type ExclusiveGatewayHandler struct{}

func (ExclusiveGatewayHandler) Type() eventv1.Element_Type {
	return eventv1.Element_TYPE_EXCLUSIVE_GATEWAY
}

func (ExclusiveGatewayHandler) OnEnter(in EnterInput) (*Effect, error) {
	var vars map[string]string
	if in.Instance != nil {
		vars = in.Instance.Variables
	}
	flowID, err := in.Deployment.ChooseExclusiveOutgoing(in.ElementID, vars)
	if err != nil {
		return nil, err
	}
	return &Effect{
		Records: InstantLifecycle(in.Type, in.ElementID, in.TokenID, func(el *eventv1.Element) {
			el.Payload = &eventv1.Element_GatewayPayload{
				GatewayPayload: &eventv1.GatewayPayload{TakenSequenceFlowId: flowID},
			}
		}),
		TakeOutgoing:   true,
		OutgoingFlowID: flowID,
	}, nil
}

func (ExclusiveGatewayHandler) OnComplete(CompleteInput) (*Effect, error) {
	return nil, errUnsupportedComplete(eventv1.Element_TYPE_EXCLUSIVE_GATEWAY)
}
