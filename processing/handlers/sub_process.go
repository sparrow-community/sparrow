package handlers

import (
	"time"

	"github.com/sparrow-community/sparrow/processing/deploy"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type SubProcessHandler struct{}

func (SubProcessHandler) Type() eventv1.Element_Type { return eventv1.Element_TYPE_SUB_PROCESS }

func (SubProcessHandler) OnEnter(in EnterInput) (*Effect, error) {
	startID, err := in.Deployment.SubProcessStartEventID(in.ElementID)
	if err != nil {
		return nil, err
	}
	activated := &eventv1.Element{Intent: eventv1.Element_INTENT_ACTIVATED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID}
	p, err := attachScopeBoundary(in.Deployment, in.ElementID, in.Now)
	if err != nil {
		return nil, err
	}
	if p != nil {
		activated.Payload = &eventv1.Element_ActivityPayload{ActivityPayload: p}
	}
	return &Effect{
		Records: []*eventv1.Element{
			{Intent: eventv1.Element_INTENT_ACTIVATING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
			activated,
		},
		EnterChild: startID,
	}, nil
}

func attachScopeBoundary(dep *deploy.Deployment, subProcessID string, now time.Time) (*eventv1.ActivityPayload, error) {
	if dep == nil {
		return nil, nil
	}
	var p *eventv1.ActivityPayload
	if bid, ok := dep.TimerBoundary(subProcessID); ok {
		if now.IsZero() {
			now = time.Now()
		}
		due, text, err := dep.TimerDue(bid, now)
		if err != nil {
			return nil, err
		}
		p = &eventv1.ActivityPayload{
			DueUnixMs:  due,
			Duration:   text,
			BoundaryId: bid,
		}
	}
	if bid, ok := dep.MessageBoundary(subProcessID); ok {
		name, err := dep.MessageName(bid)
		if err != nil {
			return nil, err
		}
		if p == nil {
			p = &eventv1.ActivityPayload{}
		}
		p.MessageName = name
		if p.BoundaryId == "" {
			p.BoundaryId = bid
		} else {
			p.MessageBoundaryId = bid
		}
	}
	if bid, ok := dep.SignalBoundary(subProcessID); ok {
		name, err := dep.SignalName(bid)
		if err != nil {
			return nil, err
		}
		if p == nil {
			p = &eventv1.ActivityPayload{}
		}
		p.SignalName = name
		if p.BoundaryId == "" {
			p.BoundaryId = bid
		} else {
			p.SignalBoundaryId = bid
		}
	}
	return p, nil
}

func (SubProcessHandler) OnComplete(in CompleteInput) (*Effect, error) {
	effect := &Effect{
		Records: []*eventv1.Element{
			{Intent: eventv1.Element_INTENT_COMPLETING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
			{Intent: eventv1.Element_INTENT_COMPLETED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
		},
	}
	if in.Deployment != nil && in.Deployment.IsEventSubProcess(in.ElementID) {
		// Event sub-process has no outgoing sequence flow; drop the token and
		// try to complete the enclosing process/scope.
		effect.DiscardToken = true
		effect.TryCompleteProcess = true
		return effect, nil
	}
	effect.Records = append(effect.Records, cancelAttachedBoundary(in.Deployment, in.ElementID, in.TokenID)...)
	effect.Records = append(effect.Records, subscribeCompensation(in.Deployment, in.ElementID, in.TokenID)...)
	effect.TakeOutgoing = true
	return effect, nil
}
