package handlers

import (
	"fmt"
	"time"

	"github.com/sparrow-community/sparrow/processing/deploy"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func attachInterruptingTimer(dep *deploy.Deployment, activityID string, now time.Time, p *eventv1.ActivityPayload) (*eventv1.ActivityPayload, error) {
	if dep == nil {
		return p, nil
	}
	bid, ok := dep.TimerBoundary(activityID)
	if !ok {
		return p, nil
	}
	if now.IsZero() {
		now = time.Now()
	}
	due, text, err := dep.TimerDue(bid, now)
	if err != nil {
		return nil, err
	}
	if p == nil {
		p = &eventv1.ActivityPayload{}
	}
	p.DueUnixMs = due
	p.Duration = text
	p.BoundaryId = bid
	return p, nil
}

func cancelAttachedTimer(dep *deploy.Deployment, activityID, tokenID string) []*eventv1.Element {
	if dep == nil {
		return nil
	}
	bid, ok := dep.TimerBoundary(activityID)
	if !ok {
		return nil
	}
	return []*eventv1.Element{
		{Intent: eventv1.Element_INTENT_TERMINATING, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: bid, TokenId: tokenID},
		{Intent: eventv1.Element_INTENT_TERMINATED, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: bid, TokenId: tokenID},
	}
}

type BoundaryEventHandler struct{}

func (BoundaryEventHandler) Type() eventv1.Element_Type {
	return eventv1.Element_TYPE_BOUNDARY_EVENT
}

func (BoundaryEventHandler) OnEnter(EnterInput) (*Effect, error) {
	return nil, fmt.Errorf("UNSUPPORTED_ELEMENT: boundary events are entered by interrupting the attached activity")
}

func (BoundaryEventHandler) OnComplete(in CompleteInput) (*Effect, error) {
	attached, ok := in.Deployment.AttachedActivity(in.ElementID)
	if !ok {
		return nil, fmt.Errorf("NOT_FOUND: boundary %q has no attached activity", in.ElementID)
	}
	typ, err := in.Deployment.TypeOf(attached)
	if err != nil {
		return nil, err
	}
	return &Effect{
		Records: []*eventv1.Element{
			{Intent: eventv1.Element_INTENT_TERMINATING, Type: typ, Id: attached, TokenId: in.TokenID},
			{Intent: eventv1.Element_INTENT_TERMINATED, Type: typ, Id: attached, TokenId: in.TokenID},
			{Intent: eventv1.Element_INTENT_COMPLETING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
			{Intent: eventv1.Element_INTENT_COMPLETED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
		},
		TakeOutgoing: true,
	}, nil
}
