package processing

import (
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type Processor interface {
	Handler(msg *eventv1.Event) error
}
