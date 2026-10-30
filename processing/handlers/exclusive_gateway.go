package handlers

import (
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type ExclusiveGatewayHandler struct{}

func (ExclusiveGatewayHandler) Type() eventv1.Element_Type {
	return eventv1.Element_TYPE_EXCLUSIVE_GATEWAY
}

func (ExclusiveGatewayHandler) OnEnter(in EnterInput) (*Effect, error) {
	// Two-phase Instant: ACTIVATING/ACTIVATED first; executor finalizes choose
	// (COMPLETING/COMPLETED + leave). Default-off still emits the full lifecycle
	// in one Enter burst. Intervention may barrier at ACTIVATED before choose.
	return &Effect{
		Records: []*eventv1.Element{
			{
				Intent:  eventv1.Element_INTENT_ACTIVATING,
				Type:    in.Type,
				Id:      in.ElementID,
				TokenId: in.TokenID,
			},
			{
				Intent:  eventv1.Element_INTENT_ACTIVATED,
				Type:    in.Type,
				Id:      in.ElementID,
				TokenId: in.TokenID,
			},
		},
		DecideExclusive: true,
	}, nil
}

func (ExclusiveGatewayHandler) OnComplete(CompleteInput) (*Effect, error) {
	return nil, errUnsupportedComplete(eventv1.Element_TYPE_EXCLUSIVE_GATEWAY)
}
