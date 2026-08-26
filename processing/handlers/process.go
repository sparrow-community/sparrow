package handlers

import (
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type ProcessHandler struct{}

func (ProcessHandler) Type() eventv1.Element_Type { return eventv1.Element_TYPE_PROCESS }

func (ProcessHandler) OnEnter(in EnterInput) (*Effect, error) {
	pid := in.Deployment.ProcessID()
	if in.Instance != nil && in.Instance.ProcessID != "" {
		pid = in.Instance.ProcessID
	}
	return &Effect{Records: ProcessStartRecords(pid, nil)}, nil
}

func (ProcessHandler) OnComplete(in CompleteInput) (*Effect, error) {
	pid := in.Deployment.ProcessID()
	if in.Instance != nil && in.Instance.ProcessID != "" {
		pid = in.Instance.ProcessID
	}
	records := []*eventv1.Element{
		{Intent: eventv1.Element_INTENT_COMPLETING, Type: eventv1.Element_TYPE_PROCESS, Id: pid},
		{Intent: eventv1.Element_INTENT_COMPLETED, Type: eventv1.Element_TYPE_PROCESS, Id: pid},
	}
	effect := &Effect{Records: records}
	if in.Instance != nil && in.Instance.ParentProcessInstanceID != "" {
		pp := &eventv1.ProcessPayload{
			ParentProcessInstanceId: in.Instance.ParentProcessInstanceID,
			ParentElementId:         in.Instance.ParentElementID,
			ParentTokenId:           in.Instance.ParentTokenID,
		}
		records[0].Payload = &eventv1.Element_ProcessPayload{ProcessPayload: pp}
		records[1].Payload = &eventv1.Element_ProcessPayload{ProcessPayload: pp}
		effect.Publish = &Publication{
			Kind:             PublicationResumeParent,
			ParentInstanceID: in.Instance.ParentProcessInstanceID,
			CallActivityID:   in.Instance.ParentElementID,
			HostTokenID:      in.Instance.ParentTokenID,
			ChildInstanceID:  in.Instance.ID,
			Completed:        true,
		}
	}
	return effect, nil
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

// ProcessStartRecordsWithParent starts a called process instance linked to a Call Activity host.
func ProcessStartRecordsWithParent(processID string, vars []*eventv1.Variable, parentInstanceID, parentElementID, parentTokenID string) []*eventv1.Element {
	pp := &eventv1.ProcessPayload{
		Variables:               vars,
		ParentProcessInstanceId: parentInstanceID,
		ParentElementId:         parentElementID,
		ParentTokenId:           parentTokenID,
	}
	return []*eventv1.Element{
		{
			Intent:  eventv1.Element_INTENT_ACTIVATING,
			Type:    eventv1.Element_TYPE_PROCESS,
			Id:      processID,
			Payload: &eventv1.Element_ProcessPayload{ProcessPayload: pp},
		},
		{
			Intent:  eventv1.Element_INTENT_ACTIVATED,
			Type:    eventv1.Element_TYPE_PROCESS,
			Id:      processID,
			Payload: &eventv1.Element_ProcessPayload{ProcessPayload: pp},
		},
	}
}
