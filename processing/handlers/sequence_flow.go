package handlers

import (
	"fmt"

	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type SequenceFlowHandler struct{}

func (SequenceFlowHandler) Type() eventv1.Element_Type { return eventv1.Element_TYPE_SEQUENCE_FLOW }

func (SequenceFlowHandler) OnEnter(EnterInput) (*Effect, error) {
	return nil, fmt.Errorf("sequence flow is taken via transit, not OnEnter")
}

func (SequenceFlowHandler) OnComplete(CompleteInput) (*Effect, error) {
	return nil, errUnsupportedComplete(eventv1.Element_TYPE_SEQUENCE_FLOW)
}

func SequenceFlowTaken(flowID, sourceID, targetID, tokenID string) *eventv1.Element {
	return &eventv1.Element{
		Intent:  eventv1.Element_INTENT_SEQUENCE_FLOW_TAKEN,
		Type:    eventv1.Element_TYPE_SEQUENCE_FLOW,
		Id:      flowID,
		TokenId: tokenID,
		Payload: &eventv1.Element_SequenceFlowPayload{
			SequenceFlowPayload: &eventv1.SequenceFlowPayload{SourceId: sourceID, TargetId: targetID},
		},
	}
}
