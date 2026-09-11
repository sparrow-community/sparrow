package handlers

import (
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type ReceiveTaskHandler struct{}

func (ReceiveTaskHandler) Type() eventv1.Element_Type { return eventv1.Element_TYPE_RECEIVE_TASK }

func (ReceiveTaskHandler) OnEnter(in EnterInput) (*Effect, error) {
	base := &eventv1.ActivityPayload{}
	if in.Deployment != nil {
		if name, err := in.Deployment.MessageName(in.ElementID); err == nil {
			base.MessageName = name
		}
	}
	return waitingTaskEnter(in, base)
}

func (ReceiveTaskHandler) OnComplete(in CompleteInput) (*Effect, error) {
	return waitingTaskComplete(in)
}
