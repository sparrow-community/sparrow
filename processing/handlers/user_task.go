package handlers

import (
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type UserTaskHandler struct{}

func (UserTaskHandler) Type() eventv1.Element_Type { return eventv1.Element_TYPE_USER_TASK }

func (UserTaskHandler) OnEnter(in EnterInput) (*Effect, error) {
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
	idx := in.LoopInstanceIndex
	if idx < 0 && in.Instance != nil {
		if tok := in.Instance.Tokens[in.TokenID]; tok != nil {
			idx = tok.LoopInstanceIndex
		}
	}
	p, err := attachBoundary(in.Deployment, in.ElementID, in.Now, activityPayloadWithIndex(nil, idx))
	if err != nil {
		return nil, err
	}
	activated := &eventv1.Element{
		Intent:  eventv1.Element_INTENT_ACTIVATED,
		Type:    in.Type,
		Id:      in.ElementID,
		TokenId: in.TokenID,
		Payload: &eventv1.Element_ActivityPayload{ActivityPayload: p},
	}
	return &Effect{
		Records: []*eventv1.Element{
			{Intent: eventv1.Element_INTENT_ACTIVATING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
			activated,
		},
		Wait: true,
	}, nil
}

func (UserTaskHandler) OnComplete(in CompleteInput) (*Effect, error) {
	idx := int32(-1)
	if in.Token != nil {
		idx = in.Token.LoopInstanceIndex
	}
	completing := &eventv1.Element{
		Intent:  eventv1.Element_INTENT_COMPLETING,
		Type:    in.Type,
		Id:      in.ElementID,
		TokenId: in.TokenID,
		Payload: &eventv1.Element_ActivityPayload{
			ActivityPayload: activityPayloadWithIndex(&eventv1.ActivityPayload{Variables: in.Variables}, idx),
		},
	}
	if len(in.Variables) == 0 {
		completing.Payload = &eventv1.Element_ActivityPayload{
			ActivityPayload: activityPayloadWithIndex(nil, idx),
		}
	}
	records := []*eventv1.Element{
		completing,
		{
			Intent:  eventv1.Element_INTENT_COMPLETED,
			Type:    in.Type,
			Id:      in.ElementID,
			TokenId: in.TokenID,
			Payload: &eventv1.Element_ActivityPayload{ActivityPayload: activityPayloadWithIndex(nil, idx)},
		},
	}
	if _, ok := in.Deployment.MultiInstanceSpec(in.ElementID); ok && in.Token != nil && !in.Token.MultiInstanceHost {
		return multiInstanceInnerComplete(in, records), nil
	}
	if in.Deployment != nil && in.Deployment.IsCompensationHandler(in.ElementID) {
		return &Effect{
			Records:             records,
			DiscardToken:        true,
			AdvanceCompensation: true,
		}, nil
	}
	records = append(records, cancelAttachedBoundary(in.Deployment, in.ElementID, in.TokenID)...)
	records = append(records, subscribeCompensation(in.Deployment, in.ElementID, in.TokenID)...)
	return &Effect{
		Records:      records,
		TakeOutgoing: true,
	}, nil
}
