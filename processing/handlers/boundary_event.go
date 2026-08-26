package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/sparrow-community/sparrow/processing/deploy"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func attachBoundary(dep *deploy.Deployment, activityID string, now time.Time, p *eventv1.ActivityPayload) (*eventv1.ActivityPayload, error) {
	if dep == nil {
		return p, nil
	}
	if bid, ok := dep.TimerBoundary(activityID); ok {
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
	}
	if bid, ok := dep.MessageBoundary(activityID); ok {
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
	if bid, ok := dep.SignalBoundary(activityID); ok {
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

func CancelAttachedBoundaries(dep *deploy.Deployment, activityID, tokenID string) []*eventv1.Element {
	return cancelAttachedBoundary(dep, activityID, tokenID)
}

func cancelAttachedBoundary(dep *deploy.Deployment, activityID, tokenID string) []*eventv1.Element {
	if dep == nil {
		return nil
	}
	var records []*eventv1.Element
	if bid, ok := dep.TimerBoundary(activityID); ok {
		records = append(records,
			&eventv1.Element{Intent: eventv1.Element_INTENT_TERMINATING, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: bid, TokenId: tokenID},
			&eventv1.Element{Intent: eventv1.Element_INTENT_TERMINATED, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: bid, TokenId: tokenID},
		)
	}
	if bid, ok := dep.MessageBoundary(activityID); ok {
		records = append(records,
			&eventv1.Element{Intent: eventv1.Element_INTENT_TERMINATING, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: bid, TokenId: tokenID},
			&eventv1.Element{Intent: eventv1.Element_INTENT_TERMINATED, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: bid, TokenId: tokenID},
		)
	}
	if bid, ok := dep.SignalBoundary(activityID); ok {
		records = append(records,
			&eventv1.Element{Intent: eventv1.Element_INTENT_TERMINATING, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: bid, TokenId: tokenID},
			&eventv1.Element{Intent: eventv1.Element_INTENT_TERMINATED, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: bid, TokenId: tokenID},
		)
	}
	// Compensation boundaries are subscribed on COMPLETED, not cancelled here.
	return records
}

// subscribeCompensation arms a compensation subscription after the host activity completes.
func subscribeCompensation(dep *deploy.Deployment, activityID, tokenID string) []*eventv1.Element {
	if dep == nil {
		return nil
	}
	comp, ok := dep.CompensationOf(activityID)
	if !ok {
		return nil
	}
	return []*eventv1.Element{
		{Intent: eventv1.Element_INTENT_ACTIVATING, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: comp.BoundaryID, TokenId: tokenID},
		{
			Intent:  eventv1.Element_INTENT_ACTIVATED,
			Type:    eventv1.Element_TYPE_BOUNDARY_EVENT,
			Id:      comp.BoundaryID,
			TokenId: tokenID,
			Payload: &eventv1.Element_EventPayload{
				EventPayload: &eventv1.EventPayload{CompensationHandlerId: comp.HandlerID},
			},
		},
	}
}

func disarmAttachedBoundary(boundaryID, activityTokenID string) []*eventv1.Element {
	return []*eventv1.Element{
		{Intent: eventv1.Element_INTENT_TERMINATING, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: boundaryID, TokenId: activityTokenID},
		{Intent: eventv1.Element_INTENT_TERMINATED, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: boundaryID, TokenId: activityTokenID},
	}
}

func rearmAttachedBoundary(boundaryID, activityTokenID, text string, dueUnixMs int64) []*eventv1.Element {
	return []*eventv1.Element{
		{Intent: eventv1.Element_INTENT_ACTIVATING, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: boundaryID, TokenId: activityTokenID},
		{
			Intent:  eventv1.Element_INTENT_ACTIVATED,
			Type:    eventv1.Element_TYPE_BOUNDARY_EVENT,
			Id:      boundaryID,
			TokenId: activityTokenID,
			Payload: &eventv1.Element_ActivityPayload{
				ActivityPayload: &eventv1.ActivityPayload{
					BoundaryId: boundaryID,
					DueUnixMs:  dueUnixMs,
					Duration:   text,
				},
			},
		},
	}
}

type BoundaryEventHandler struct{}

func (BoundaryEventHandler) Type() eventv1.Element_Type {
	return eventv1.Element_TYPE_BOUNDARY_EVENT
}

func (BoundaryEventHandler) OnEnter(EnterInput) (*Effect, error) {
	return nil, fmt.Errorf("UNSUPPORTED_ELEMENT: boundary events are entered via Complete on the attached activity token")
}

func (BoundaryEventHandler) OnComplete(in CompleteInput) (*Effect, error) {
	attached, ok := in.Deployment.AttachedActivity(in.ElementID)
	if !ok {
		return nil, fmt.Errorf("NOT_FOUND: boundary %q has no attached activity", in.ElementID)
	}
	if !in.Deployment.BoundaryInterrupting(in.ElementID) {
		records := disarmAttachedBoundary(in.ElementID, in.TokenID)
		if in.Token != nil && in.Token.DueUnixMs > 0 && strings.HasPrefix(strings.TrimSpace(in.Token.TimerText), "R") {
			nextText, nextDue, ok, err := deploy.NextCycleTimer(in.Token.TimerText, time.UnixMilli(in.Token.DueUnixMs))
			if err != nil {
				return nil, err
			}
			if ok {
				records = append(records, rearmAttachedBoundary(in.ElementID, in.TokenID, nextText, nextDue.UnixMilli())...)
			}
		}
		return &Effect{
			Records: records,
			SpawnOutgoing: &SpawnOutgoingEffect{
				ElementID: in.ElementID,
				Type:      in.Type,
			},
		}, nil
	}
	typ, err := in.Deployment.TypeOf(attached)
	if err != nil {
		return nil, err
	}
	records := []*eventv1.Element{
		{Intent: eventv1.Element_INTENT_TERMINATING, Type: typ, Id: attached, TokenId: in.TokenID},
		{Intent: eventv1.Element_INTENT_TERMINATED, Type: typ, Id: attached, TokenId: in.TokenID},
	}
	// Terminate sibling boundary (the other boundary on the same activity).
	if bid, ok := in.Deployment.TimerBoundary(attached); ok && bid != in.ElementID {
		records = append(records,
			&eventv1.Element{Intent: eventv1.Element_INTENT_TERMINATING, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: bid, TokenId: in.TokenID},
			&eventv1.Element{Intent: eventv1.Element_INTENT_TERMINATED, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: bid, TokenId: in.TokenID},
		)
	}
	if bid, ok := in.Deployment.MessageBoundary(attached); ok && bid != in.ElementID {
		records = append(records,
			&eventv1.Element{Intent: eventv1.Element_INTENT_TERMINATING, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: bid, TokenId: in.TokenID},
			&eventv1.Element{Intent: eventv1.Element_INTENT_TERMINATED, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: bid, TokenId: in.TokenID},
		)
	}
	if bid, ok := in.Deployment.SignalBoundary(attached); ok && bid != in.ElementID {
		records = append(records,
			&eventv1.Element{Intent: eventv1.Element_INTENT_TERMINATING, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: bid, TokenId: in.TokenID},
			&eventv1.Element{Intent: eventv1.Element_INTENT_TERMINATED, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: bid, TokenId: in.TokenID},
		)
	}
	records = append(records,
		&eventv1.Element{Intent: eventv1.Element_INTENT_COMPLETING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
		&eventv1.Element{Intent: eventv1.Element_INTENT_COMPLETED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
	)
	effect := &Effect{
		Records:      records,
		TakeOutgoing: true,
	}
	if typ == eventv1.Element_TYPE_CALL_ACTIVITY && in.Token != nil && in.Token.CalledProcessInstanceID != "" {
		effect.Publish = &Publication{
			Kind:            PublicationTerminateChild,
			ChildInstanceID: in.Token.CalledProcessInstanceID,
			CallActivityID:  attached,
			HostTokenID:     in.TokenID,
		}
	}
	return effect, nil
}
