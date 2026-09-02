package handlers

import (
	"fmt"

	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type CallActivityHandler struct{}

func (CallActivityHandler) Type() eventv1.Element_Type { return eventv1.Element_TYPE_CALL_ACTIVITY }

func (CallActivityHandler) OnEnter(in EnterInput) (*Effect, error) {
	call, ok := in.Deployment.CallActivitySpec(in.ElementID)
	if !ok {
		return nil, fmt.Errorf("NOT_FOUND: callActivity %q", in.ElementID)
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
	return &Effect{
		Records: []*eventv1.Element{
			{Intent: eventv1.Element_INTENT_ACTIVATING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
			{
				Intent:  eventv1.Element_INTENT_ACTIVATED,
				Type:    in.Type,
				Id:      in.ElementID,
				TokenId: in.TokenID,
				Payload: &eventv1.Element_ActivityPayload{
					ActivityPayload: &eventv1.ActivityPayload{CalledProcessInstanceId: childID},
				},
			},
		},
		Wait: true,
		Publish: &Publication{
			Kind:                 PublicationStartChild,
			ChildInstanceID:      childID,
			ParentInstanceID:     in.Instance.ID,
			CallActivityID:       in.ElementID,
			HostTokenID:          in.TokenID,
			CalledProcessID:      call.CalledProcessID,
			DeploymentID:         in.Deployment.ID,
			CalledDeploymentID:   calleeDepID,
		},
	}, nil
}

func (CallActivityHandler) OnComplete(in CompleteInput) (*Effect, error) {
	vars := in.Variables
	completed := &eventv1.Element{
		Intent:  eventv1.Element_INTENT_COMPLETED,
		Type:    in.Type,
		Id:      in.ElementID,
		TokenId: in.TokenID,
	}
	if len(vars) > 0 {
		completed.Payload = &eventv1.Element_ActivityPayload{
			ActivityPayload: &eventv1.ActivityPayload{Variables: vars},
		}
	}
	return &Effect{
		Records: []*eventv1.Element{
			{Intent: eventv1.Element_INTENT_COMPLETING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
			completed,
		},
		TakeOutgoing: true,
	}, nil
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
