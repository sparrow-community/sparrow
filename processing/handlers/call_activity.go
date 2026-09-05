package handlers

import (
	"fmt"

	"github.com/sparrow-community/sparrow/processing/deploy"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type CallActivityHandler struct{}

func (CallActivityHandler) Type() eventv1.Element_Type { return eventv1.Element_TYPE_CALL_ACTIVITY }

func (CallActivityHandler) OnEnter(in EnterInput) (*Effect, error) {
	call, ok := in.Deployment.CallActivitySpec(in.ElementID)
	if !ok {
		return nil, fmt.Errorf("NOT_FOUND: callActivity %q", in.ElementID)
	}
	if spec, ok := in.Deployment.MultiInstanceSpec(in.ElementID); ok && !isMultiInstanceInner(in) {
		total, err := spec.InstanceCount(in.Instance.Variables)
		if err != nil {
			return nil, err
		}
		return multiInstanceHostEnter(in, spec, total, func(_ int32, p *eventv1.ActivityPayload) *eventv1.Element {
			return &eventv1.Element{
				Intent:  eventv1.Element_INTENT_ACTIVATED,
				Type:    in.Type,
				Id:      in.ElementID,
				TokenId: in.TokenID,
				Payload: &eventv1.Element_ActivityPayload{ActivityPayload: p},
			}
		})
	}

	calleeDepID := in.Deployment.ID
	if call.ExternalCallee {
		id, err := resolveCalleeDeployment(call.CalledProcessID)
		if err != nil {
			return nil, err
		}
		calleeDepID = id
	}
	childID, err := nextCallChildID()
	if err != nil {
		return nil, err
	}
	idx := in.LoopInstanceIndex
	if idx < 0 && in.Instance != nil {
		if tok := in.Instance.Tokens[in.TokenID]; tok != nil {
			idx = tok.LoopInstanceIndex
		}
	}
	payload := activityPayloadWithIndex(&eventv1.ActivityPayload{CalledProcessInstanceId: childID}, idx)
	if idx < 0 {
		var err error
		payload, err = attachBoundary(in.Deployment, in.ElementID, in.Now, payload)
		if err != nil {
			return nil, err
		}
	}
	activated := &eventv1.Element{
		Intent:  eventv1.Element_INTENT_ACTIVATED,
		Type:    in.Type,
		Id:      in.ElementID,
		TokenId: in.TokenID,
		Payload: &eventv1.Element_ActivityPayload{ActivityPayload: payload},
	}
	pub := &Publication{
		Kind:               PublicationStartChild,
		ChildInstanceID:    childID,
		ParentInstanceID:   in.Instance.ID,
		CallActivityID:     in.ElementID,
		HostTokenID:        in.TokenID,
		CalledProcessID:    call.CalledProcessID,
		DeploymentID:       in.Deployment.ID,
		CalledDeploymentID: calleeDepID,
		ChildVariables:     snapCallChildInputs(call, in.Instance.Variables),
	}
	return &Effect{
		Records: []*eventv1.Element{
			{Intent: eventv1.Element_INTENT_ACTIVATING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
			activated,
		},
		Wait:    true,
		Publish: pub,
	}, nil
}

func (CallActivityHandler) OnComplete(in CompleteInput) (*Effect, error) {
	idx := int32(-1)
	if in.Token != nil {
		idx = in.Token.LoopInstanceIndex
	}
	vars := in.Variables
	completed := &eventv1.Element{
		Intent:  eventv1.Element_INTENT_COMPLETED,
		Type:    in.Type,
		Id:      in.ElementID,
		TokenId: in.TokenID,
		Payload: &eventv1.Element_ActivityPayload{
			ActivityPayload: activityPayloadWithIndex(&eventv1.ActivityPayload{Variables: vars}, idx),
		},
	}
	if len(vars) == 0 {
		completed.Payload = &eventv1.Element_ActivityPayload{
			ActivityPayload: activityPayloadWithIndex(nil, idx),
		}
	}
	records := []*eventv1.Element{
		{
			Intent:  eventv1.Element_INTENT_COMPLETING,
			Type:    in.Type,
			Id:      in.ElementID,
			TokenId: in.TokenID,
			Payload: &eventv1.Element_ActivityPayload{
				ActivityPayload: activityPayloadWithIndex(&eventv1.ActivityPayload{Variables: vars}, idx),
			},
		},
		completed,
	}
	if _, ok := in.Deployment.MultiInstanceSpec(in.ElementID); ok && in.Token != nil && !in.Token.MultiInstanceHost {
		return &Effect{
			Records:                    records,
			MultiInstanceInnerComplete: true,
		}, nil
	}
	records = append(records, cancelAttachedBoundary(in.Deployment, in.ElementID, in.TokenID)...)
	records = append(records, subscribeCompensation(in.Deployment, in.ElementID, in.TokenID)...)
	return &Effect{
		Records:      records,
		TakeOutgoing: true,
	}, nil
}

func snapCallChildInputs(call deploy.CallActivity, src map[string]string) []*eventv1.Variable {
	if len(call.Inputs) == 0 || len(src) == 0 {
		return nil
	}
	out := make([]*eventv1.Variable, 0, len(call.Inputs))
	for _, m := range call.Inputs {
		if v, ok := src[m.Source]; ok {
			out = append(out, &eventv1.Variable{Name: m.Target, JsonValue: v})
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// nextCallChildID is set by the processing package to avoid an import cycle.
var nextCallChildID = func() (string, error) {
	return "", fmt.Errorf("call child id allocator not configured")
}

// resolveCalleeDeployment resolves calledElement to a deployment id at runtime.
var resolveCalleeDeployment = func(processID string) (string, error) {
	return "", fmt.Errorf("callee resolver not configured")
}

// SetCallChildIDAllocator wires UUIDv7 allocation from processing into handlers.
func SetCallChildIDAllocator(fn func() (string, error)) {
	if fn != nil {
		nextCallChildID = fn
	}
}

// SetCalleeResolver wires runtime callee deployment resolution from processing into handlers.
func SetCalleeResolver(fn func(processID string) (deploymentID string, err error)) {
	if fn != nil {
		resolveCalleeDeployment = fn
	}
}
