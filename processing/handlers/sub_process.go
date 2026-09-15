package handlers

import (
	"time"

	"github.com/sparrow-community/sparrow/processing/deploy"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type SubProcessHandler struct{}

func (SubProcessHandler) Type() eventv1.Element_Type { return eventv1.Element_TYPE_SUB_PROCESS }

func (SubProcessHandler) OnEnter(in EnterInput) (*Effect, error) {
	if spec, ok := in.Deployment.MultiInstanceSpec(in.ElementID); ok && !in.Deployment.IsEventSubProcess(in.ElementID) && !isMultiInstanceInner(in) {
		total, err := spec.InstanceCount(in.Instance.Variables)
		if err != nil {
			return nil, err
		}
		return multiInstanceHostEnter(in, spec, total, func(_ int32, p *eventv1.ActivityPayload) *eventv1.Element {
			activated := &eventv1.Element{Intent: eventv1.Element_INTENT_ACTIVATED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID}
			if p != nil {
				activated.Payload = &eventv1.Element_ActivityPayload{ActivityPayload: p}
			}
			return activated
		})
	}
	var startID string
	var err error
	if spec, ok := in.Deployment.EventSubProcessSpec(in.ElementID); ok {
		startID = spec.StartEventID
	} else {
		startID, err = in.Deployment.SubProcessStartEventID(in.ElementID)
		if err != nil {
			return nil, err
		}
	}
	activated := &eventv1.Element{Intent: eventv1.Element_INTENT_ACTIVATED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID}
	idx := in.LoopInstanceIndex
	if idx < 0 && in.Instance != nil {
		if tok := in.Instance.Tokens[in.TokenID]; tok != nil {
			idx = tok.LoopInstanceIndex
		}
	}
	p, err := attachScopeBoundary(in.Deployment, in.ElementID, in.Now)
	if err != nil {
		return nil, err
	}
	if p == nil {
		p = &eventv1.ActivityPayload{}
	}
	p = activityPayloadWithIndex(p, idx)
	hasBoundary := p.GetBoundaryId() != "" || p.GetMessageBoundaryId() != "" || p.GetSignalBoundaryId() != "" ||
		p.GetDueUnixMs() != 0 || len(p.GetWaitingBoundaries()) > 0
	if !hasBoundary && idx < 0 {
		p = nil
	} else if idx >= 0 && !hasBoundary {
		p = activityPayloadWithIndex(nil, idx)
	}
	if p != nil {
		activated.Payload = &eventv1.Element_ActivityPayload{ActivityPayload: p}
	}
	effect := &Effect{
		Records: []*eventv1.Element{
			{Intent: eventv1.Element_INTENT_ACTIVATING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
			activated,
		},
		EnterChild: startID,
	}
	// Embedded SubProcess: park host token; mint a child for the internal start.
	// Event Sub-Process keeps the trigger token (no outgoing host).
	if in.Deployment == nil || !in.Deployment.IsEventSubProcess(in.ElementID) {
		effect.SpawnChildToken = true
	}
	return effect, nil
}

func attachScopeBoundary(dep *deploy.Deployment, subProcessID string, now time.Time) (*eventv1.ActivityPayload, error) {
	return attachBoundary(dep, subProcessID, now, nil)
}

func (SubProcessHandler) OnComplete(in CompleteInput) (*Effect, error) {
	effect := &Effect{
		Records: []*eventv1.Element{
			{Intent: eventv1.Element_INTENT_COMPLETING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
			{Intent: eventv1.Element_INTENT_COMPLETED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
		},
	}
	if in.Deployment != nil && in.Deployment.IsEventSubProcess(in.ElementID) {
		// Event sub-process has no outgoing sequence flow; drop the token.
		// Compensation event sub-process advances the pending compensate throw;
		// other event sub-processes try to complete the enclosing scope.
		effect.DiscardToken = true
		if in.Deployment.IsCompensationHandler(in.ElementID) {
			effect.AdvanceCompensation = true
		} else {
			effect.TryCompleteProcess = true
		}
		return effect, nil
	}
	if _, ok := in.Deployment.MultiInstanceSpec(in.ElementID); ok && in.Token != nil && in.Token.ScopeHost && !in.Token.MultiInstanceHost {
		idx := in.Token.LoopInstanceIndex
		ap := activityPayloadWithIndex(nil, idx)
		effect.Records[0].Payload = &eventv1.Element_ActivityPayload{ActivityPayload: ap}
		effect.Records[1].Payload = &eventv1.Element_ActivityPayload{ActivityPayload: ap}
		return multiInstanceInnerComplete(in, effect.Records), nil
	}
	effect.Records = append(effect.Records, cancelAttachedBoundary(in.Deployment, in.ElementID, in.TokenID)...)
	effect.Records = append(effect.Records, subscribeCompensation(in.Deployment, in.ElementID, in.TokenID)...)
	effect.TakeOutgoing = true
	return effect, nil
}
