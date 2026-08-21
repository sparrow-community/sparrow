package handlers

import (
	"fmt"
	"time"

	"github.com/sparrow-community/sparrow/processing/deploy"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type IntermediateCatchEventHandler struct{}

func (IntermediateCatchEventHandler) Type() eventv1.Element_Type {
	return eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT
}

func (IntermediateCatchEventHandler) OnEnter(in EnterInput) (*Effect, error) {
	kind, err := in.Deployment.CatchKind(in.ElementID)
	if err != nil {
		return nil, err
	}
	now := in.Now
	if now.IsZero() {
		now = time.Now()
	}

	var payload *eventv1.EventPayload
	switch kind {
	case deploy.CatchKindTimer:
		dueUnixMs, text, err := in.Deployment.TimerDue(in.ElementID, now)
		if err != nil {
			return nil, err
		}
		payload = &eventv1.EventPayload{
			DueUnixMs: dueUnixMs,
			Duration:  text,
		}
	case deploy.CatchKindMessage:
		name, err := in.Deployment.MessageName(in.ElementID)
		if err != nil {
			return nil, err
		}
		payload = &eventv1.EventPayload{MessageName: name}
	case deploy.CatchKindSignal:
		name, err := in.Deployment.SignalName(in.ElementID)
		if err != nil {
			return nil, err
		}
		payload = &eventv1.EventPayload{SignalName: name}
	default:
		return nil, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateCatchEvent %q has unsupported kind %q", in.ElementID, string(kind))
	}

	return &Effect{
		Records: []*eventv1.Element{
			{Intent: eventv1.Element_INTENT_ACTIVATING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
			{
				Intent:  eventv1.Element_INTENT_ACTIVATED,
				Type:    in.Type,
				Id:      in.ElementID,
				TokenId: in.TokenID,
				Payload: &eventv1.Element_EventPayload{EventPayload: payload},
			},
		},
		Wait: true,
	}, nil
}

func (IntermediateCatchEventHandler) OnComplete(in CompleteInput) (*Effect, error) {
	completing := &eventv1.Element{
		Intent:  eventv1.Element_INTENT_COMPLETING,
		Type:    in.Type,
		Id:      in.ElementID,
		TokenId: in.TokenID,
	}
	if len(in.Variables) > 0 {
		completing.Payload = &eventv1.Element_EventPayload{
			EventPayload: &eventv1.EventPayload{Variables: in.Variables},
		}
	}
	return &Effect{
		Records: []*eventv1.Element{
			completing,
			{Intent: eventv1.Element_INTENT_COMPLETED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
		},
		TakeOutgoing:       true,
		TerminateWaitingAt: eventBasedSiblingCatches(in),
	}, nil
}

func eventBasedSiblingCatches(in CompleteInput) []string {
	if in.Deployment == nil {
		return nil
	}
	_, siblings := in.Deployment.EventBasedSiblings(in.ElementID)
	return siblings
}
