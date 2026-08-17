package handlers

import (
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type ProcessHandler struct{}

func (ProcessHandler) Type() eventv1.Element_Type { return eventv1.Element_TYPE_PROCESS }

func (ProcessHandler) OnEnter(in EnterInput) (*Effect, error) {
	return &Effect{Records: ProcessStartRecords(in.Deployment.ProcessID(), nil)}, nil
}

func (ProcessHandler) OnComplete(in CompleteInput) (*Effect, error) {
	pid := in.Deployment.ProcessID()
	return &Effect{
		Records: []*eventv1.Element{
			{Intent: eventv1.Element_INTENT_COMPLETING, Type: eventv1.Element_TYPE_PROCESS, Id: pid},
			{Intent: eventv1.Element_INTENT_COMPLETED, Type: eventv1.Element_TYPE_PROCESS, Id: pid},
		},
	}, nil
}

func ProcessStartRecords(processID string, vars []*eventv1.Variable) []*eventv1.Element {
	activating := &eventv1.Element{
		Intent: eventv1.Element_INTENT_ACTIVATING,
		Type:   eventv1.Element_TYPE_PROCESS,
		Id:     processID,
	}
	if len(vars) > 0 {
		activating.Payload = &eventv1.Element_ProcessPayload{
			ProcessPayload: &eventv1.ProcessPayload{Variables: vars},
		}
	}
	return []*eventv1.Element{
		activating,
		{Intent: eventv1.Element_INTENT_ACTIVATED, Type: eventv1.Element_TYPE_PROCESS, Id: processID},
	}
}
