package deploy

import (
	"fmt"

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
