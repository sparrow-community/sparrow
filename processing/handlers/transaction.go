package handlers

import (
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// TransactionHandler runs BPMN transaction SubProcesses (success path like SubProcess).
type TransactionHandler struct{}

func (TransactionHandler) Type() eventv1.Element_Type { return eventv1.Element_TYPE_TRANSACTION }

func (TransactionHandler) OnEnter(in EnterInput) (*Effect, error) {
	startID, err := in.Deployment.TransactionStartEventID(in.ElementID)
	if err != nil {
		return nil, err
	}
	activated := &eventv1.Element{Intent: eventv1.Element_INTENT_ACTIVATED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID}
	idx := in.LoopInstanceIndex
	if idx < 0 && in.Instance != nil {
		if tok := in.Instance.Tokens[in.TokenID]; tok != nil {
			idx = tok.LoopInstanceIndex
		}
	}
	p, err := attachScopeBoundary(in.Deployment, in.ElementID, in.Now)
	if err != nil {
		return nil, err
	}
	if p == nil {
		p = &eventv1.ActivityPayload{}
	}
	p = activityPayloadWithIndex(p, idx)
	hasBoundary := p.GetBoundaryId() != "" || p.GetMessageBoundaryId() != "" || p.GetSignalBoundaryId() != "" ||
		p.GetDueUnixMs() != 0 || len(p.GetWaitingBoundaries()) > 0
	if !hasBoundary && idx < 0 {
		p = nil
	} else if idx >= 0 && !hasBoundary {
		p = activityPayloadWithIndex(nil, idx)
	}
	if p != nil {
		activated.Payload = &eventv1.Element_ActivityPayload{ActivityPayload: p}
	}
	return &Effect{
		Records: []*eventv1.Element{
			{Intent: eventv1.Element_INTENT_ACTIVATING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
			activated,
		},
		EnterChild:       startID,
		SpawnChildToken: true,
	}, nil
}

func (TransactionHandler) OnComplete(in CompleteInput) (*Effect, error) {
	effect := &Effect{
		Records: []*eventv1.Element{
			{Intent: eventv1.Element_INTENT_COMPLETING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
			{Intent: eventv1.Element_INTENT_COMPLETED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
		},
	}
	effect.Records = append(effect.Records, cancelAttachedBoundary(in.Deployment, in.ElementID, in.TokenID)...)
	effect.Records = append(effect.Records, subscribeCompensation(in.Deployment, in.ElementID, in.TokenID)...)
	effect.TakeOutgoing = true
	return effect, nil
}
