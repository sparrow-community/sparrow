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
	if code, ok := in.Deployment.EscalationEndCode(in.ElementID); ok {
		return &Effect{
			Records: []*eventv1.Element{
				{Intent: eventv1.Element_INTENT_ACTIVATING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
				{
					Intent:  eventv1.Element_INTENT_ACTIVATED,
					Type:    in.Type,
					Id:      in.ElementID,
					TokenId: in.TokenID,
					Payload: &eventv1.Element_EventPayload{
						EventPayload: &eventv1.EventPayload{EscalationCode: code},
					},
				},
			},
			ThrowEscalation:    &ThrowEscalationEffect{EscalationCode: code},
			TryCompleteProcess: true,
		}, nil
	}
	if in.Deployment.IsCompensateEnd(in.ElementID) {
		return &Effect{
			Records: []*eventv1.Element{
				{Intent: eventv1.Element_INTENT_ACTIVATING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
				{
					Intent:  eventv1.Element_INTENT_ACTIVATED,
					Type:    in.Type,
					Id:      in.ElementID,
					TokenId: in.TokenID,
					Payload: &eventv1.Element_EventPayload{
						EventPayload: &eventv1.EventPayload{TokenWait: true},
					},
				},
			},
			Wait:                true,
			TriggerCompensation: true,
		}, nil
	}
	if in.Deployment != nil && in.Deployment.IsTerminateEnd(in.ElementID) {
		return &Effect{
			Records:            InstantLifecycle(in.Type, in.ElementID, in.TokenID, nil),
			TerminateScope:     true,
			TryCompleteProcess: true,
		}, nil
	}
	return &Effect{
		Records:            InstantLifecycle(in.Type, in.ElementID, in.TokenID, nil),
		TryCompleteProcess: true,
	}, nil
}

func (EndEventHandler) OnComplete(in CompleteInput) (*Effect, error) {
	if in.Deployment == nil || !in.Deployment.IsCompensateEnd(in.ElementID) {
		return nil, errUnsupportedComplete(eventv1.Element_TYPE_END_EVENT)
	}
	return &Effect{
		Records: []*eventv1.Element{
			{Intent: eventv1.Element_INTENT_COMPLETING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
			{Intent: eventv1.Element_INTENT_COMPLETED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
		},
		TryCompleteProcess: true,
	}, nil
}
