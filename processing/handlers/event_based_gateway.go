package handlers

import (
	"fmt"

	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// EventBasedGatewayHandler implements event-based gateway:
// fork tokens to all outgoing intermediate catches or receive tasks.
// Exclusive (default): the first target to complete cancels sibling waiting targets.
// Parallel: siblings stay armed; each target continues on its own path when its event arrives.
type EventBasedGatewayHandler struct{}

func (EventBasedGatewayHandler) Type() eventv1.Element_Type {
	return eventv1.Element_TYPE_EVENT_BASED_GATEWAY
}

func (EventBasedGatewayHandler) OnEnter(in EnterInput) (*Effect, error) {
	outgoing := in.Deployment.Outgoing(in.ElementID)
	if len(outgoing) < 2 {
		return nil, fmt.Errorf("NO_OUTGOING_FLOW: eventBasedGateway %q needs at least two outgoing flows", in.ElementID)
	}
	return &Effect{
		Records: InstantLifecycle(in.Type, in.ElementID, in.TokenID, nil),
		Fork:    append([]string{}, outgoing...),
	}, nil
}

func (EventBasedGatewayHandler) OnComplete(CompleteInput) (*Effect, error) {
	return nil, errUnsupportedComplete(eventv1.Element_TYPE_EVENT_BASED_GATEWAY)
}
