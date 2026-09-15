package handlers

import (
	"time"

	"github.com/sparrow-community/sparrow/processing/deploy"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

const loopHostIndex int32 = -1

// MultiInstanceStart tells the executor to spawn inner instances after the host ACTIVATED.
type MultiInstanceStart struct {
	HostTokenID  string
	ElementID    string
	ElementType  eventv1.Element_Type
	Total        int
	Sequential   bool
	InnerIndices []int32 // indices to activate now
}

func isMultiInstanceInner(in EnterInput) bool {
	if in.Deployment != nil {
		if _, ok := in.Deployment.MultiInstanceSpec(in.ElementID); !ok {
			return false
		}
	}
	if in.LoopInstanceIndex >= 0 {
		return true
	}
	return isMultiInstanceInnerToken(in.Instance, in.TokenID)
}

func isMultiInstanceInnerToken(inst *projection.Instance, tokenID string) bool {
	if inst == nil {
		return false
	}
	tok := inst.Tokens[tokenID]
	return tok != nil && tok.LoopInstanceIndex >= 0
}

func multiInstanceHostEnter(in EnterInput, spec deploy.MultiInstanceSpec, total int, buildActivated func(loopIndex int32, p *eventv1.ActivityPayload) *eventv1.Element) (*Effect, error) {
	if total <= 0 {
		records := InstantLifecycle(in.Type, in.ElementID, in.TokenID, nil)
		return &Effect{Records: records, TakeOutgoing: true}, nil
	}
	hostPayload := &eventv1.ActivityPayload{
		LoopInstanceIndex:  loopHostIndex,
		LoopTotalInstances: int32(total),
		LoopSequential:     spec.Sequential,
	}
	p, err := attachBoundary(in.Deployment, in.ElementID, in.Now, hostPayload)
	if err != nil {
		return nil, err
	}
	activated := buildActivated(loopHostIndex, p)
	if activated == nil {
		activated = &eventv1.Element{
			Intent:  eventv1.Element_INTENT_ACTIVATED,
			Type:    in.Type,
			Id:      in.ElementID,
			TokenId: in.TokenID,
			Payload: &eventv1.Element_ActivityPayload{ActivityPayload: p},
		}
	}
	indices := make([]int32, 0, total)
	if spec.Sequential {
		indices = append(indices, 0)
	} else {
		for i := 0; i < total; i++ {
			indices = append(indices, int32(i))
		}
	}
	return &Effect{
		Records: []*eventv1.Element{
			{Intent: eventv1.Element_INTENT_ACTIVATING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
			activated,
		},
		Wait: true,
		MultiInstanceStart: &MultiInstanceStart{
			HostTokenID:  in.TokenID,
			ElementID:    in.ElementID,
			ElementType:  in.Type,
			Total:        total,
			Sequential:   spec.Sequential,
			InnerIndices: indices,
		},
	}, nil
}

func activityPayloadWithIndex(base *eventv1.ActivityPayload, index int32) *eventv1.ActivityPayload {
	p := &eventv1.ActivityPayload{}
	if base != nil {
		*p = *base
	}
	p.LoopInstanceIndex = index
	return p
}

func attachBoundaryAt(in EnterInput, elementID string, now time.Time, p *eventv1.ActivityPayload) (*eventv1.ActivityPayload, error) {
	return attachBoundary(in.Deployment, elementID, now, p)
}

func multiInstanceInnerComplete(in CompleteInput, records []*eventv1.Element) *Effect {
	if in.Deployment != nil && in.Deployment.IsCompensationHandler(in.ElementID) {
		return &Effect{
			Records:             records,
			DiscardToken:        true,
			AdvanceCompensation: true,
		}
	}
	records = append(records, cancelAttachedBoundary(in.Deployment, in.ElementID, in.TokenID)...)
	records = append(records, subscribeCompensation(in.Deployment, in.ElementID, in.TokenID)...)
	return &Effect{
		Records:                    records,
		MultiInstanceInnerComplete: true,
	}
}
