package deploy

import (
	"fmt"
	"strings"

	"github.com/sparrow-community/sparrow/bpmn/element"
)

func terminateEndSpec(ev element.EndEvent) error {
	if len(ev.TerminateEventDefinitions) != 1 {
		return fmt.Errorf("not a terminate end")
	}
	other := extraCatchDefinitions(ev.EventDefinitions) - len(ev.TerminateEventDefinitions)
	if other > 0 || len(ev.TimerEventDefinitions) > 0 {
		return fmt.Errorf("not a terminate end")
	}
	return nil
}

func messageEndSpec(ev element.EndEvent, messages []element.Message) (string, error) {
	if len(ev.TimerEventDefinitions) > 0 {
		return "", fmt.Errorf("not a message end")
	}
	if len(ev.MessageEventDefinitions) != 1 {
		return "", fmt.Errorf("not a message end")
	}
	other := extraCatchDefinitions(ev.EventDefinitions) - len(ev.MessageEventDefinitions)
	if other > 0 {
		return "", fmt.Errorf("not a message end")
	}
	name := resolveMessageName(ev.MessageEventDefinitions[0].MessageRef, messages)
	if name == "" {
		name = strings.TrimSpace(ev.Name)
	}
	if name == "" {
		name = ev.ID
	}
	return name, nil
}

func signalEndSpec(ev element.EndEvent, signals []element.Signal) (string, error) {
	if len(ev.TimerEventDefinitions) > 0 {
		return "", fmt.Errorf("not a signal end")
	}
	if len(ev.SignalEventDefinitions) != 1 {
		return "", fmt.Errorf("not a signal end")
	}
	other := extraCatchDefinitions(ev.EventDefinitions) - len(ev.SignalEventDefinitions)
	if other > 0 {
		return "", fmt.Errorf("not a signal end")
	}
	name := resolveSignalName(ev.SignalEventDefinitions[0].SignalRef, signals)
	if name == "" {
		name = strings.TrimSpace(ev.Name)
	}
	if name == "" {
		name = ev.ID
	}
	return name, nil
}

func endHasEventDefinitions(ev element.EndEvent) bool {
	if extraCatchDefinitions(ev.EventDefinitions) > 0 {
		return true
	}
	return len(ev.TimerEventDefinitions) > 0
}

// IsTerminateEnd reports whether id is a terminate end event.
func (d *Deployment) IsTerminateEnd(id string) bool {
	return d != nil && d.terminateEnds[id]
}

// MessageEndName returns the message name for a message end event.
func (d *Deployment) MessageEndName(id string) (string, bool) {
	if d == nil {
		return "", false
	}
	name, ok := d.messageEnds[id]
	return name, ok && name != ""
}

// SignalEndName returns the signal name for a signal end event.
func (d *Deployment) SignalEndName(id string) (string, bool) {
	if d == nil {
		return "", false
	}
	name, ok := d.signalEnds[id]
	return name, ok && name != ""
}
