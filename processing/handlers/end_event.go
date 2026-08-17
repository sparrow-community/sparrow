package handlers

import (
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type EndEventHandler struct{}

func (EndEventHandler) Type() eventv1.Element_Type { return eventv1.Element_TYPE_END_EVENT }

func (EndEventHandler) OnEnter(in EnterInput) (*Effect, error) {
	return &Effect{
		Records:            InstantLifecycle(in.Type, in.ElementID, in.TokenID, nil),
		TryCompleteProcess: true,
	}, nil
}

func (EndEventHandler) OnComplete(CompleteInput) (*Effect, error) {
	return nil, errUnsupportedComplete(eventv1.Element_TYPE_END_EVENT)
}
