package handlers

import (
	"fmt"
	"time"

	"github.com/sparrow-community/sparrow/processing/deploy"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// Effect describes element behavior outcomes to be applied by the executor.
type Effect struct {
	Records            []*eventv1.Element
	Wait               bool
	OutgoingFlowID     string
	TakeOutgoing       bool
	TryCompleteProcess bool
	// Fork lists outgoing sequence flow ids for a parallel split. The entering
	// token takes Fork[0]; additional tokens are minted for Fork[1:].
	Fork []string
	// TerminateJoinPeers ends other waiting tokens still at this parallel/inclusive join.
	TerminateJoinPeers string
	// TerminateWaitingAt ends waiting tokens at the listed element IDs (other tokens).
	// Used by exclusive event-based gateway: winning catch cancels sibling catches.
	TerminateWaitingAt []string
	// SpawnOutgoing mints a new token to follow the boundary outgoing flow while the
	// activity token keeps waiting (non-interrupting boundary).
	SpawnOutgoing *SpawnOutgoingEffect
	// EnterChild tells the executor to enter the given element on the same token
	// (used by SubProcess to enter the internal start event after ACTIVATED).
	EnterChild string
	// Publish is a deferred message/signal throw. The engine delivers it after
	// the instance lock is released (avoids re-entrant Complete under the same lock).
	Publish *Publication
	// DiscardToken removes the token after records are applied (event sub-process
	// completion has no outgoing sequence flow).
	DiscardToken bool
}

// PublicationKind is message or signal broadcast from an intermediate throw.
type PublicationKind string

const (
	PublicationMessage PublicationKind = "message"
	PublicationSignal  PublicationKind = "signal"
)

// Publication is delivered by the engine after Enter/Complete unlocks the instance.
type Publication struct {
	Kind PublicationKind
	Name string
}

// SpawnOutgoingEffect describes a boundary path taken on a newly minted token.
type SpawnOutgoingEffect struct {
	ElementID      string
	Type           eventv1.Element_Type
	OutgoingFlowID string
}

// EnterInput is the context when a token arrives at an element.
type EnterInput struct {
	Deployment *deploy.Deployment
	Instance   *projection.Instance
	ElementID  string
	Type       eventv1.Element_Type
	TokenID    string
	Now        time.Time
}

// CompleteInput is the context for an external completion command.
type CompleteInput struct {
	Deployment *deploy.Deployment
	Instance   *projection.Instance
	Token      *projection.Token
	ElementID  string
	Type       eventv1.Element_Type
	TokenID    string
	Variables  []*eventv1.Variable
}

// ElementHandler encapsulates BPMN element-type behavior.
type ElementHandler interface {
	Type() eventv1.Element_Type
	OnEnter(in EnterInput) (*Effect, error)
	OnComplete(in CompleteInput) (*Effect, error)
}

// Registry maps Element.Type to handlers.
type Registry struct {
	byType map[eventv1.Element_Type]ElementHandler
}

func NewRegistry(hs ...ElementHandler) *Registry {
	r := &Registry{byType: make(map[eventv1.Element_Type]ElementHandler, len(hs))}
	for _, h := range hs {
		r.byType[h.Type()] = h
	}
	return r
}

func DefaultRegistry() *Registry {
	return NewRegistry(
		ProcessHandler{},
		StartEventHandler{},
		EndEventHandler{},
		UserTaskHandler{},
		ServiceTaskHandler{},
		ExclusiveGatewayHandler{},
		ParallelGatewayHandler{},
		IntermediateCatchEventHandler{},
		BoundaryEventHandler{},
		SequenceFlowHandler{},
		SubProcessHandler{},
		InclusiveGatewayHandler{},
		EventBasedGatewayHandler{},
		IntermediateThrowEventHandler{},
	)
}

func (r *Registry) Get(t eventv1.Element_Type) (ElementHandler, error) {
	h, ok := r.byType[t]
	if !ok {
		return nil, fmt.Errorf("UNSUPPORTED_ELEMENT: no handler for type %v", t)
	}
	return h, nil
}

func errUnsupportedComplete(t eventv1.Element_Type) error {
	return fmt.Errorf("UNSUPPORTED_ELEMENT: OnComplete not supported for %v", t)
}

func InstantLifecycle(typ eventv1.Element_Type, elementID, tokenID string, onCompleted func(*eventv1.Element)) []*eventv1.Element {
	intents := []eventv1.Element_Intent{
		eventv1.Element_INTENT_ACTIVATING,
		eventv1.Element_INTENT_ACTIVATED,
		eventv1.Element_INTENT_COMPLETING,
		eventv1.Element_INTENT_COMPLETED,
	}
	out := make([]*eventv1.Element, 0, len(intents))
	for _, intent := range intents {
		el := &eventv1.Element{
			Intent:  intent,
			Type:    typ,
			Id:      elementID,
			TokenId: tokenID,
		}
		if intent == eventv1.Element_INTENT_COMPLETED && onCompleted != nil {
			onCompleted(el)
		}
		out = append(out, el)
	}
	return out
}
