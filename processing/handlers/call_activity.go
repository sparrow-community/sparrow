package handlers

import (
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type CallActivityHandler struct{}

func (CallActivityHandler) Type() eventv1.Element_Type { return eventv1.Element_TYPE_CALL_ACTIVITY }

func (CallActivityHandler) OnEnter(in EnterInput) (*Effect, error) {
	startID, err := in.Deployment.CalledProcessStartEventID(in.ElementID)
	if err != nil {
		return nil, err
	}
	return &Effect{
		Records: []*eventv1.Element{
			{Intent: eventv1.Element_INTENT_ACTIVATING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
			{Intent: eventv1.Element_INTENT_ACTIVATED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
		},
		EnterChild:      startID,
		SpawnChildToken: true,
	}, nil
}

func (CallActivityHandler) OnComplete(in CompleteInput) (*Effect, error) {
	return &Effect{
		Records: []*eventv1.Element{
			{Intent: eventv1.Element_INTENT_COMPLETING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
			{Intent: eventv1.Element_INTENT_COMPLETED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
		},
		TakeOutgoing: true,
	}, nil
}
