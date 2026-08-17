package handlers

import (
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type StartEventHandler struct{}

func (StartEventHandler) Type() eventv1.Element_Type { return eventv1.Element_TYPE_START_EVENT }

func (StartEventHandler) OnEnter(in EnterInput) (*Effect, error) {
	return &Effect{
		Records:      InstantLifecycle(in.Type, in.ElementID, in.TokenID, nil),
		TakeOutgoing: true,
	}, nil
}

func (StartEventHandler) OnComplete(CompleteInput) (*Effect, error) {
	return nil, errUnsupportedComplete(eventv1.Element_TYPE_START_EVENT)
}
