package handlers

import (
	"fmt"

	"github.com/sparrow-community/sparrow/processing/deploy"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type IntermediateThrowEventHandler struct{}

func (IntermediateThrowEventHandler) Type() eventv1.Element_Type {
	return eventv1.Element_TYPE_INTERMEDIATE_THROW_EVENT
}

func (IntermediateThrowEventHandler) OnEnter(in EnterInput) (*Effect, error) {
	kind, err := in.Deployment.ThrowKind(in.ElementID)
	if err != nil {
		return nil, err
	}
	var payload *eventv1.EventPayload
	var pub *Publication
	switch kind {
	case deploy.ThrowKindNone:
		// Milestone / none throw: instantaneous pass-through.
	case deploy.ThrowKindMessage:
		name, err := in.Deployment.ThrowName(in.ElementID)
		if err != nil {
			return nil, err
		}
		payload = &eventv1.EventPayload{MessageName: name}
		pub = &Publication{Kind: PublicationMessage, Name: name}
	case deploy.ThrowKindSignal:
		name, err := in.Deployment.ThrowName(in.ElementID)
		if err != nil {
			return nil, err
		}
		payload = &eventv1.EventPayload{SignalName: name}
		pub = &Publication{Kind: PublicationSignal, Name: name}
	default:
		return nil, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateThrowEvent %q has unsupported kind %q", in.ElementID, string(kind))
	}

	records := InstantLifecycle(in.Type, in.ElementID, in.TokenID, nil)
	if payload != nil {
		for _, rec := range records {
			if rec.GetIntent() == eventv1.Element_INTENT_ACTIVATED {
				rec.Payload = &eventv1.Element_EventPayload{EventPayload: payload}
				break
			}
		}
	}
	return &Effect{
		Records:      records,
		TakeOutgoing: true,
		Publish:      pub,
	}, nil
}

func (IntermediateThrowEventHandler) OnComplete(CompleteInput) (*Effect, error) {
	return nil, errUnsupportedComplete(eventv1.Element_TYPE_INTERMEDIATE_THROW_EVENT)
}
