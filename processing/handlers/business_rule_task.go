package handlers

import (
	"strings"

	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type BusinessRuleTaskHandler struct{}

func (BusinessRuleTaskHandler) Type() eventv1.Element_Type {
	return eventv1.Element_TYPE_BUSINESS_RULE_TASK
}

func (BusinessRuleTaskHandler) OnEnter(in EnterInput) (*Effect, error) {
	jobType := in.ElementID
	if st, err := in.Deployment.BusinessRuleTask(in.ElementID); err == nil {
		jobType = jobTypeFrom(st.Implementation, st.Name, st.ID)
	}
	return jobTaskEnter(in, jobType)
}

func (BusinessRuleTaskHandler) OnComplete(in CompleteInput) (*Effect, error) {
	return waitingTaskComplete(in)
}

func jobTypeFrom(implementation, name, id string) string {
	if t := strings.TrimSpace(implementation); t != "" && t != "##unspecified" {
		return t
	}
	if t := strings.TrimSpace(name); t != "" {
		return t
	}
	return id
}
