package handlers

import (
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type SendTaskHandler struct{}

func (SendTaskHandler) Type() eventv1.Element_Type { return eventv1.Element_TYPE_SEND_TASK }

func (SendTaskHandler) OnEnter(in EnterInput) (*Effect, error) {
	if spec, ok := in.Deployment.MultiInstanceSpec(in.ElementID); ok && !isMultiInstanceInner(in) {
		total, err := spec.InstanceCount(in.Instance.Variables)
		if err != nil {
			return nil, err
		}
		return multiInstanceHostEnter(in, spec, total, func(_ int32, p *eventv1.ActivityPayload) *eventv1.Element {
			return &eventv1.Element{
				Intent:  eventv1.Element_INTENT_ACTIVATED,
				Type:    in.Type,
				Id:      in.ElementID,
				TokenId: in.TokenID,
				Payload: &eventv1.Element_ActivityPayload{ActivityPayload: p},
			}
		})
	}
	name, err := in.Deployment.ThrowName(in.ElementID)
	if err != nil {
		return nil, err
	}
	records := InstantLifecycle(in.Type, in.ElementID, in.TokenID, nil)
	for _, rec := range records {
		if rec.GetIntent() == eventv1.Element_INTENT_ACTIVATED {
			rec.Payload = &eventv1.Element_ActivityPayload{
				ActivityPayload: &eventv1.ActivityPayload{MessageName: name},
			}
			break
		}
	}
	effect := &Effect{
		Records: records,
		Publish: &Publication{Kind: PublicationMessage, Name: name},
	}
	if _, ok := in.Deployment.MultiInstanceSpec(in.ElementID); ok && isMultiInstanceInner(in) {
		effect.MultiInstanceInnerComplete = true
		return effect, nil
	}
	effect.TakeOutgoing = true
	return effect, nil
}

func (SendTaskHandler) OnComplete(CompleteInput) (*Effect, error) {
	return nil, errUnsupportedComplete(eventv1.Element_TYPE_SEND_TASK)
}
