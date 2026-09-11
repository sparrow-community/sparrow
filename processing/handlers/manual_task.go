package handlers

import (
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type ManualTaskHandler struct{}

func (ManualTaskHandler) Type() eventv1.Element_Type { return eventv1.Element_TYPE_MANUAL_TASK }

func (ManualTaskHandler) OnEnter(in EnterInput) (*Effect, error) {
	return waitingTaskEnter(in, nil)
}

func (ManualTaskHandler) OnComplete(in CompleteInput) (*Effect, error) {
	return waitingTaskComplete(in)
}
