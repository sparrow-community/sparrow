package handlers

import (
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// TaskHandler runs abstract bpmn:task as Manual-equivalent wait → Complete.
type TaskHandler struct{}

func (TaskHandler) Type() eventv1.Element_Type { return eventv1.Element_TYPE_TASK }

func (TaskHandler) OnEnter(in EnterInput) (*Effect, error) {
	return waitingTaskEnter(in, nil)
}

func (TaskHandler) OnComplete(in CompleteInput) (*Effect, error) {
	return waitingTaskComplete(in)
}
