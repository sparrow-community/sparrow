package processing

import (
	"github.com/sparrow-community/sparrow/bpmn/element"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type Process struct {
	element.Process
}

func NewProcessInstance(process element.Process) *Process {
	return &Process{
		Process: process,
	}
}

func (pi *Process) Handler(e *eventv1.Event) error {

	return nil
}
