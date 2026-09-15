package handlers

import (
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type UserTaskHandler struct{}

func (UserTaskHandler) Type() eventv1.Element_Type { return eventv1.Element_TYPE_USER_TASK }

func (UserTaskHandler) OnEnter(in EnterInput) (*Effect, error) {
	return waitingTaskEnter(in, nil)
}

func (UserTaskHandler) OnComplete(in CompleteInput) (*Effect, error) {
	return waitingTaskComplete(in)
}
