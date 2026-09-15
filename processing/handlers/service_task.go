package handlers

import (
	"strings"

	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type ServiceTaskHandler struct{}

func (ServiceTaskHandler) Type() eventv1.Element_Type { return eventv1.Element_TYPE_SERVICE_TASK }

func (ServiceTaskHandler) OnEnter(in EnterInput) (*Effect, error) {
	jobType := in.ElementID
	if st, err := in.Deployment.ServiceTask(in.ElementID); err == nil {
		if t := strings.TrimSpace(st.Implementation); t != "" && t != "##unspecified" {
			jobType = t
		} else if strings.TrimSpace(st.Name) != "" {
			jobType = st.Name
		}
	}
	return jobTaskEnter(in, jobType)
}

func (ServiceTaskHandler) OnComplete(in CompleteInput) (*Effect, error) {
	return waitingTaskComplete(in)
}
