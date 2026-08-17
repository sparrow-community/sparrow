package processing

import (
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// Processor handles COMMAND records. Engine currently inlines processing;
// the interface remains for DESIGN compatibility and future stream processors.
type Processor interface {
	Handler(msg *eventv1.Event) error
}
