package deploy

import (
	"fmt"
	"strings"

	"github.com/sparrow-community/sparrow/bpmn/element"
)

// ThrowKind classifies an intermediateThrowEvent for runtime handlers.
type ThrowKind string

const (
	ThrowKindNone       ThrowKind = "none"
	ThrowKindMessage    ThrowKind = "message"
	ThrowKindSignal     ThrowKind = "signal"
	ThrowKindCompensate ThrowKind = "compensate"
)

type throwSpec struct {
	Kind        ThrowKind
	Name        string // message or signal name; empty for none/compensate
	ActivityRef string // optional compensate target activity id
}

func messageThrowSpec(ev element.IntermediateThrowEvent, messages []element.Message) (throwSpec, error) {
	if len(ev.TimerEventDefinitions) > 0 {
		return throwSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateThrowEvent %q cannot have timerEventDefinition", ev.ID)
	}
	if len(ev.MessageEventDefinitions) != 1 {
		return throwSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateThrowEvent %q needs exactly one messageEventDefinition", ev.ID)
	}
	other := extraCatchDefinitions(ev.EventDefinitions) - len(ev.MessageEventDefinitions)
	if other > 0 {
		return throwSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateThrowEvent %q is not a message throw", ev.ID)
	}
	name := resolveMessageName(ev.MessageEventDefinitions[0].MessageRef, messages)
	if name == "" {
		name = strings.TrimSpace(ev.Name)
	}
	if name == "" {
		name = ev.ID
	}
	return throwSpec{Kind: ThrowKindMessage, Name: name}, nil
}

func signalThrowSpec(ev element.IntermediateThrowEvent, signals []element.Signal) (throwSpec, error) {
	if len(ev.TimerEventDefinitions) > 0 {
		return throwSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateThrowEvent %q cannot have timerEventDefinition", ev.ID)
	}
	if len(ev.SignalEventDefinitions) != 1 {
		return throwSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateThrowEvent %q needs exactly one signalEventDefinition", ev.ID)
	}
	other := extraCatchDefinitions(ev.EventDefinitions) - len(ev.SignalEventDefinitions)
	if other > 0 {
		return throwSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateThrowEvent %q is not a signal throw", ev.ID)
	}
	name := resolveSignalName(ev.SignalEventDefinitions[0].SignalRef, signals)
	if name == "" {
		name = strings.TrimSpace(ev.Name)
	}
	if name == "" {
		name = ev.ID
	}
	return throwSpec{Kind: ThrowKindSignal, Name: name}, nil
}

func noneThrowSpec(ev element.IntermediateThrowEvent) (throwSpec, error) {
	if extraCatchDefinitions(ev.EventDefinitions) > 0 || len(ev.TimerEventDefinitions) > 0 {
		return throwSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateThrowEvent %q has unsupported event definitions", ev.ID)
	}
	return throwSpec{Kind: ThrowKindNone}, nil
}

func throwEventSpec(ev element.IntermediateThrowEvent, messages []element.Message, signals []element.Signal) (throwSpec, error) {
	switch {
	case len(ev.MessageEventDefinitions) > 0:
		return messageThrowSpec(ev, messages)
	case len(ev.SignalEventDefinitions) > 0:
		return signalThrowSpec(ev, signals)
	case len(ev.CompensateEventDefinitions) > 0:
		return compensateThrowSpec(ev)
	default:
		return noneThrowSpec(ev)
	}
}

func signalCatchSpec(ev element.IntermediateCatchEvent, signals []element.Signal) (string, error) {
	if len(ev.TimerEventDefinitions) > 0 {
		return "", fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateCatchEvent %q is not a signal catch (has timerEventDefinition)", ev.ID)
	}
	return signalCatchFromDefs(ev.ID, ev.EventDefinitions, signals, strings.TrimSpace(ev.Name))
}

func signalCatchFromDefs(id string, defs element.EventDefinitions, signals []element.Signal, fallbackName string) (string, error) {
	if len(defs.TimerEventDefinitions) > 0 {
		return "", fmt.Errorf("UNSUPPORTED_ELEMENT: %q is not a signal catch (has timerEventDefinition)", id)
	}
	if len(defs.SignalEventDefinitions) != 1 {
		return "", fmt.Errorf("UNSUPPORTED_ELEMENT: %q needs exactly one signalEventDefinition", id)
	}
	other := extraCatchDefinitions(defs) - len(defs.SignalEventDefinitions)
	if other > 0 {
		return "", fmt.Errorf("UNSUPPORTED_ELEMENT: %q is not a signal catch", id)
	}
	name := resolveSignalName(defs.SignalEventDefinitions[0].SignalRef, signals)
	if name == "" {
		name = fallbackName
	}
	if name == "" {
		name = id
	}
	return name, nil
}

func resolveSignalName(ref string, signals []element.Signal) string {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return ""
	}
	for _, s := range signals {
		if s.ID != ref {
			continue
		}
		if n := strings.TrimSpace(s.Name); n != "" {
			return n
		}
		return s.ID
	}
	return ref
}

// ThrowKind returns the throw classification for an intermediateThrowEvent.
func (d *Deployment) ThrowKind(id string) (ThrowKind, error) {
	spec, ok := d.throwEvents[id]
	if !ok {
		return "", fmt.Errorf("NOT_FOUND: intermediate throw %q", id)
	}
	return spec.Kind, nil
}

// ThrowName returns the message or signal name for a throw event.
func (d *Deployment) ThrowName(id string) (string, error) {
	spec, ok := d.throwEvents[id]
	if !ok {
		return "", fmt.Errorf("NOT_FOUND: intermediate throw %q", id)
	}
	if spec.Kind == ThrowKindNone {
		return "", nil
	}
	return spec.Name, nil
}

// SignalName returns the BPMN signal name a catch is waiting for.
func (d *Deployment) SignalName(id string) (string, error) {
	name, ok := d.signalCatch[id]
	if !ok || name == "" {
		return "", fmt.Errorf("NOT_FOUND: signal catch %q", id)
	}
	return name, nil
}
