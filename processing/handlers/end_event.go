package handlers

import (
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type EndEventHandler struct{}

func (EndEventHandler) Type() eventv1.Element_Type { return eventv1.Element_TYPE_END_EVENT }

func (EndEventHandler) OnEnter(in EnterInput) (*Effect, error) {
	if code, ok := in.Deployment.ErrorEndCode(in.ElementID); ok {
		return &Effect{
			Records: []*eventv1.Element{
				{Intent: eventv1.Element_INTENT_ACTIVATING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
				{Intent: eventv1.Element_INTENT_ACTIVATED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
			},
			ThrowError: &ThrowErrorEffect{ErrorCode: code},
		}, nil
	}
	return &Effect{
		Records:            InstantLifecycle(in.Type, in.ElementID, in.TokenID, nil),
		TryCompleteProcess: true,
	}, nil
}

func (EndEventHandler) OnComplete(CompleteInput) (*Effect, error) {
	return nil, errUnsupportedComplete(eventv1.Element_TYPE_END_EVENT)
}
