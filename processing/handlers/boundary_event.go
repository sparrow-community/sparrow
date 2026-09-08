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
	if now.IsZero() {
		now = time.Now()
	}
	var waits []*eventv1.WaitingBoundary
	for _, bid := range dep.TimerBoundaries(activityID) {
		due, text, err := dep.TimerDue(bid, now)
		if err != nil {
			return nil, err
		}
		waits = append(waits, &eventv1.WaitingBoundary{
			BoundaryId: bid,
			Kind:       "timer",
			DueUnixMs:  due,
			Duration:   text,
		})
		if p == nil {
			p = &eventv1.ActivityPayload{}
		}
		if p.BoundaryId == "" {
			p.DueUnixMs = due
			p.Duration = text
			p.BoundaryId = bid
		}
	}
	for _, bid := range dep.MessageBoundaries(activityID) {
		name, err := dep.MessageName(bid)
		if err != nil {
			return nil, err
		}
		waits = append(waits, &eventv1.WaitingBoundary{
			BoundaryId:  bid,
			Kind:        "message",
			MessageName: name,
		})
		if p == nil {
			p = &eventv1.ActivityPayload{}
		}
		if p.MessageName == "" {
			p.MessageName = name
		}
		if p.BoundaryId == "" {
			p.BoundaryId = bid
		} else if p.MessageBoundaryId == "" && p.BoundaryId != bid {
			p.MessageBoundaryId = bid
		}
	}
	for _, bid := range dep.SignalBoundaries(activityID) {
		name, err := dep.SignalName(bid)
		if err != nil {
			return nil, err
		}
		waits = append(waits, &eventv1.WaitingBoundary{
			BoundaryId: bid,
			Kind:       "signal",
			SignalName: name,
		})
		if p == nil {
			p = &eventv1.ActivityPayload{}
		}
		if p.SignalName == "" {
			p.SignalName = name
		}
		if p.BoundaryId == "" {
			p.BoundaryId = bid
		} else if p.SignalBoundaryId == "" && p.BoundaryId != bid {
			p.SignalBoundaryId = bid
		}
	}
	if len(waits) > 0 {
		if p == nil {
			p = &eventv1.ActivityPayload{}
		}
		p.WaitingBoundaries = waits
	}
	return p, nil
}

func CancelAttachedBoundaries(dep *deploy.Deployment, activityID, tokenID string) []*eventv1.Element {
	return cancelAttachedBoundary(dep, activityID, tokenID)
}

// SubscribeCompensation exports compensation subscription records for host complete paths.
func SubscribeCompensation(dep *deploy.Deployment, activityID, tokenID string) []*eventv1.Element {
	return subscribeCompensation(dep, activityID, tokenID)
}

func cancelAttachedBoundary(dep *deploy.Deployment, activityID, tokenID string) []*eventv1.Element {
	if dep == nil {
		return nil
	}
	var records []*eventv1.Element
	seen := map[string]bool{}
	appendCancel := func(bid string) {
		if bid == "" || seen[bid] {
			return
		}
		seen[bid] = true
		records = append(records,
			&eventv1.Element{Intent: eventv1.Element_INTENT_TERMINATING, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: bid, TokenId: tokenID},
			&eventv1.Element{Intent: eventv1.Element_INTENT_TERMINATED, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: bid, TokenId: tokenID},
		)
	}
	for _, bid := range dep.TimerBoundaries(activityID) {
		appendCancel(bid)
	}
	for _, bid := range dep.MessageBoundaries(activityID) {
		appendCancel(bid)
	}
	for _, bid := range dep.SignalBoundaries(activityID) {
		appendCancel(bid)
	}
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
		{
			Intent:  eventv1.Element_INTENT_ACTIVATING,
			Type:    eventv1.Element_TYPE_BOUNDARY_EVENT,
			Id:      comp.BoundaryID,
			TokenId: tokenID,
			Payload: &eventv1.Element_EventPayload{
				EventPayload: &eventv1.EventPayload{CompensationHandlerId: comp.HandlerID},
			},
		},
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
					WaitingBoundaries: []*eventv1.WaitingBoundary{{
						BoundaryId: boundaryID,
						Kind:       "timer",
						DueUnixMs:  dueUnixMs,
						Duration:   text,
					}},
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
	// Terminate sibling waiting boundaries on the same activity.
	seen := map[string]bool{in.ElementID: true}
	appendSibling := func(bid string) {
		if bid == "" || seen[bid] {
			return
		}
		seen[bid] = true
		records = append(records,
			&eventv1.Element{Intent: eventv1.Element_INTENT_TERMINATING, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: bid, TokenId: in.TokenID},
			&eventv1.Element{Intent: eventv1.Element_INTENT_TERMINATED, Type: eventv1.Element_TYPE_BOUNDARY_EVENT, Id: bid, TokenId: in.TokenID},
		)
	}
	for _, bid := range in.Deployment.TimerBoundaries(attached) {
		appendSibling(bid)
	}
	for _, bid := range in.Deployment.MessageBoundaries(attached) {
		appendSibling(bid)
	}
	for _, bid := range in.Deployment.SignalBoundaries(attached) {
		appendSibling(bid)
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
	if _, ok := in.Deployment.MultiInstanceSpec(attached); ok {
		effect.MultiInstanceCancel = attached
	}
	return effect, nil
}
