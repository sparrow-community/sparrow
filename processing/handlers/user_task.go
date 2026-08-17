package handlers

import (
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type UserTaskHandler struct{}

func (UserTaskHandler) Type() eventv1.Element_Type { return eventv1.Element_TYPE_USER_TASK }

func (UserTaskHandler) OnEnter(in EnterInput) (*Effect, error) {
	return &Effect{
		Records: []*eventv1.Element{
			{Intent: eventv1.Element_INTENT_ACTIVATING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
			{Intent: eventv1.Element_INTENT_ACTIVATED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
		},
		Wait: true,
	}, nil
}

func (UserTaskHandler) OnComplete(in CompleteInput) (*Effect, error) {
	completing := &eventv1.Element{
		Intent:  eventv1.Element_INTENT_COMPLETING,
		Type:    in.Type,
		Id:      in.ElementID,
		TokenId: in.TokenID,
	}
	if len(in.Variables) > 0 {
		completing.Payload = &eventv1.Element_ActivityPayload{
			ActivityPayload: &eventv1.ActivityPayload{Variables: in.Variables},
		}
	}
	return &Effect{
		Records: []*eventv1.Element{
			completing,
			{Intent: eventv1.Element_INTENT_COMPLETED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
		},
		TakeOutgoing: true,
	}, nil
}
