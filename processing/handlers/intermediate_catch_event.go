package handlers

import (
	"time"

	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type IntermediateCatchEventHandler struct{}

func (IntermediateCatchEventHandler) Type() eventv1.Element_Type {
	return eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT
}

func (IntermediateCatchEventHandler) OnEnter(in EnterInput) (*Effect, error) {
	dur, text, err := in.Deployment.TimerDuration(in.ElementID)
	if err != nil {
		return nil, err
	}
	now := in.Now
	if now.IsZero() {
		now = time.Now()
	}
	return &Effect{
		Records: []*eventv1.Element{
			{Intent: eventv1.Element_INTENT_ACTIVATING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
			{
				Intent:  eventv1.Element_INTENT_ACTIVATED,
				Type:    in.Type,
				Id:      in.ElementID,
				TokenId: in.TokenID,
				Payload: &eventv1.Element_EventPayload{
					EventPayload: &eventv1.EventPayload{
						DueUnixMs: now.Add(dur).UnixMilli(),
						Duration:  text,
					},
				},
			},
		},
		Wait: true,
	}, nil
}

func (IntermediateCatchEventHandler) OnComplete(in CompleteInput) (*Effect, error) {
	return &Effect{
		Records: []*eventv1.Element{
			{Intent: eventv1.Element_INTENT_COMPLETING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
			{Intent: eventv1.Element_INTENT_COMPLETED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
		},
		TakeOutgoing: true,
	}, nil
}
