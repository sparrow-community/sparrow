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
	p := &eventv1.ActivityPayload{JobType: jobType}
	p, err := attachInterruptingBoundary(in.Deployment, in.ElementID, in.Now, p)
	if err != nil {
		return nil, err
	}
	return &Effect{
		Records: []*eventv1.Element{
			{Intent: eventv1.Element_INTENT_ACTIVATING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
			{
				Intent:  eventv1.Element_INTENT_ACTIVATED,
				Type:    in.Type,
				Id:      in.ElementID,
				TokenId: in.TokenID,
				Payload: &eventv1.Element_ActivityPayload{ActivityPayload: p},
			},
		},
		Wait: true,
	}, nil
}

func (ServiceTaskHandler) OnComplete(in CompleteInput) (*Effect, error) {
	completing := &eventv1.Element{
		Intent:  eventv1.Element_INTENT_COMPLETING,
		Type:    in.Type,
		Id:      in.ElementID,
		TokenId: in.TokenID,
	}
	if len(in.Variables) > 0 {
		completing.Payload = &eventv1.Element_ActivityPayload{
			ActivityPayload: &eventv1.ActivityPayload{Variables: in.Variables},
		}
	}
	records := []*eventv1.Element{
		completing,
		{Intent: eventv1.Element_INTENT_COMPLETED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
	}
	records = append(records, cancelAttachedBoundary(in.Deployment, in.ElementID, in.TokenID)...)
	return &Effect{
		Records:      records,
		TakeOutgoing: true,
	}, nil
}
