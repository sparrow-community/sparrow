package handlers

import (
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type SubProcessHandler struct{}

func (SubProcessHandler) Type() eventv1.Element_Type { return eventv1.Element_TYPE_SUB_PROCESS }

func (SubProcessHandler) OnEnter(in EnterInput) (*Effect, error) {
	startID, err := in.Deployment.SubProcessStartEventID(in.ElementID)
	if err != nil {
		return nil, err
	}
	return &Effect{
		Records: []*eventv1.Element{
			{Intent: eventv1.Element_INTENT_ACTIVATING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
			{Intent: eventv1.Element_INTENT_ACTIVATED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
		},
		EnterChild: startID,
	}, nil
}

func (SubProcessHandler) OnComplete(in CompleteInput) (*Effect, error) {
	return &Effect{
		Records: []*eventv1.Element{
			{Intent: eventv1.Element_INTENT_COMPLETING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
			{Intent: eventv1.Element_INTENT_COMPLETED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
		},
		TakeOutgoing: true,
	}, nil
}
