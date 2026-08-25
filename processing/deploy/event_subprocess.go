package deploy

import (
	"fmt"
	"strings"

	"github.com/sparrow-community/sparrow/bpmn/element"
)

// EventSubProcess describes a triggeredByEvent subProcess and its start trigger.
type EventSubProcess struct {
	ID            string
	ParentScopeID string // process id or embedding subProcess id
	StartEventID  string
	Interrupting  bool
	Kind          CatchKind
	MessageName   string
	SignalName    string
	// Timer facts live in timerCatch[StartEventID] when Kind == CatchKindTimer.
}

func validateEventSubProcess(sp *element.SubProcess, messages []element.Message, signals []element.Signal) error {
	if !sp.TriggeredByEvent {
		return nil
	}
	if len(sp.Incoming) > 0 || len(sp.Outgoing) > 0 {
		return fmt.Errorf("UNSUPPORTED_ELEMENT: event subProcess %q must not have sequence flow connections", sp.ID)
	}
	if len(sp.StartEvents) != 1 {
		return fmt.Errorf("UNSUPPORTED_ELEMENT: event subProcess %q needs exactly one startEvent", sp.ID)
	}
	start := sp.StartEvents[0]
	if _, err := eventSubProcessStartSpec(sp.ID, start, messages, signals); err != nil {
		return err
	}
	return nil
}

func eventSubProcessStartSpec(subProcessID string, start element.StartEvent, messages []element.Message, signals []element.Signal) (EventSubProcess, error) {
	spec := EventSubProcess{
		ID:           subProcessID,
		StartEventID: start.ID,
		// BPMN default for event sub-process start is interrupting when attribute absent.
		Interrupting: start.IsInterrupting(),
	}
	// Prefer message, then signal, then timer.
	if len(start.MessageEventDefinitions) > 0 {
		mc, err := messageCatchFromDefs(start.ID, start.EventDefinitions, messages, strings.TrimSpace(start.Name))
		if err != nil {
			return EventSubProcess{}, fmt.Errorf("UNSUPPORTED_ELEMENT: event subProcess %q startEvent: %v", subProcessID, err)
		}
		spec.Kind = CatchKindMessage
		spec.MessageName = mc.Name
		return spec, nil
	}
	if len(start.SignalEventDefinitions) > 0 {
		name, err := signalCatchFromDefs(start.ID, start.EventDefinitions, signals, strings.TrimSpace(start.Name))
		if err != nil {
			return EventSubProcess{}, fmt.Errorf("UNSUPPORTED_ELEMENT: event subProcess %q startEvent: %v", subProcessID, err)
		}
		spec.Kind = CatchKindSignal
		spec.SignalName = name
		return spec, nil
	}
	if len(start.TimerEventDefinitions) > 0 {
		if _, err := timerCatchFromDefs(start.ID, start.EventDefinitions); err != nil {
			return EventSubProcess{}, fmt.Errorf("UNSUPPORTED_ELEMENT: event subProcess %q startEvent: %v", subProcessID, err)
		}
		spec.Kind = CatchKindTimer
		return spec, nil
	}
	return EventSubProcess{}, fmt.Errorf("UNSUPPORTED_ELEMENT: event subProcess %q startEvent must be message, signal, or timer", subProcessID)
}

// EventSubProcessByStartEvent returns the event sub-process whose start event is startEventID.
func (d *Deployment) EventSubProcessByStartEvent(startEventID string) (EventSubProcess, bool) {
	for _, sp := range d.eventSubProcesses {
		if sp.StartEventID == startEventID {
			return sp, true
		}
	}
	return EventSubProcess{}, false
}

// IsEventSubProcess reports whether id is a triggeredByEvent subProcess.
func (d *Deployment) IsEventSubProcess(id string) bool {
	_, ok := d.eventSubProcesses[id]
	return ok
}

// EventSubProcessSpec returns the compiled event sub-process, if any.
func (d *Deployment) EventSubProcessSpec(id string) (EventSubProcess, bool) {
	sp, ok := d.eventSubProcesses[id]
	return sp, ok
}

// EventSubProcessesInScope returns event sub-processes whose parent is scopeID.
func (d *Deployment) EventSubProcessesInScope(scopeID string) []EventSubProcess {
	var out []EventSubProcess
	for _, sp := range d.eventSubProcesses {
		if sp.ParentScopeID == scopeID {
			out = append(out, sp)
		}
	}
	return out
}
