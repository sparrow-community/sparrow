package handlers

import (
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type ScriptTaskHandler struct{}

func (ScriptTaskHandler) Type() eventv1.Element_Type { return eventv1.Element_TYPE_SCRIPT_TASK }

func (ScriptTaskHandler) OnEnter(in EnterInput) (*Effect, error) {
	jobType := in.ElementID
	if st, err := in.Deployment.ScriptTask(in.ElementID); err == nil {
		jobType = jobTypeFrom(st.ScriptFormat, st.Name, st.ID)
	}
	return jobTaskEnter(in, jobType)
}

func (ScriptTaskHandler) OnComplete(in CompleteInput) (*Effect, error) {
	return waitingTaskComplete(in)
}
