package deploy

import (
	"fmt"
	"strings"

	"github.com/sparrow-community/sparrow/bpmn/element"
)

// EventSubProcess describes a triggeredByEvent subProcess and its start trigger.
type EventSubProcess struct {
	ID             string
	ParentScopeID  string // process id or embedding subProcess id
	StartEventID   string
	Interrupting   bool
	Kind           CatchKind
	MessageName    string
	SignalName     string
	ErrorCode      string // empty = catch-all when Kind == CatchKindError
	EscalationCode string // empty = catch-all when Kind == CatchKindEscalation
	Condition      string // set when Kind == CatchKindConditional
	// Timer facts live in timerCatch[StartEventID] when Kind == CatchKindTimer.
}

func validateEventSubProcess(sp *element.SubProcess, messages []element.Message, signals []element.Signal, errors []element.Error, escalations []element.Escalation) error {
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
	if _, err := eventSubProcessStartSpec(sp.ID, start, messages, signals, errors, escalations); err != nil {
		return err
	}
	return nil
}

func eventSubProcessStartSpec(subProcessID string, start element.StartEvent, messages []element.Message, signals []element.Signal, errors []element.Error, escalations []element.Escalation) (EventSubProcess, error) {
	spec := EventSubProcess{
		ID:           subProcessID,
		StartEventID: start.ID,
		// BPMN default for event sub-process start is interrupting when attribute absent.
		Interrupting: start.IsInterrupting(),
	}
	// Prefer message, then signal, then timer, then error, then escalation.
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
	if len(start.ErrorEventDefinitions) > 0 {
		code, err := errorStartCatchFromDefs(start.ID, start.EventDefinitions, errors)
		if err != nil {
			return EventSubProcess{}, fmt.Errorf("UNSUPPORTED_ELEMENT: event subProcess %q startEvent: %v", subProcessID, err)
		}
		spec.Kind = CatchKindError
		spec.ErrorCode = code
		return spec, nil
	}
	if len(start.EscalationEventDefinitions) > 0 {
		code, err := escalationStartCatchFromDefs(start.ID, start.EventDefinitions, escalations)
		if err != nil {
			return EventSubProcess{}, fmt.Errorf("UNSUPPORTED_ELEMENT: event subProcess %q startEvent: %v", subProcessID, err)
		}
		spec.Kind = CatchKindEscalation
		spec.EscalationCode = code
		return spec, nil
	}
	if len(start.CompensateEventDefinitions) > 0 {
		if err := compensateEventSubProcessStartFromDefs(start.ID, start.EventDefinitions); err != nil {
			return EventSubProcess{}, fmt.Errorf("UNSUPPORTED_ELEMENT: event subProcess %q startEvent: %v", subProcessID, err)
		}
		spec.Kind = CatchKindCompensate
		// Compensation event sub-process starts are never interrupting live waits.
		spec.Interrupting = false
		return spec, nil
	}
	if len(start.ConditionalEventDefinitions) > 0 {
		text, err := conditionalCatchFromDefs(start.ID, start.EventDefinitions)
		if err != nil {
			return EventSubProcess{}, fmt.Errorf("UNSUPPORTED_ELEMENT: event subProcess %q startEvent: %v", subProcessID, err)
		}
		spec.Kind = CatchKindConditional
		spec.Condition = text
		return spec, nil
	}
	return EventSubProcess{}, fmt.Errorf("UNSUPPORTED_ELEMENT: event subProcess %q startEvent must be message, signal, timer, error, escalation, compensate, or conditional", subProcessID)
}

func compensateEventSubProcessStartFromDefs(startEventID string, defs element.EventDefinitions) error {
	if len(defs.CompensateEventDefinitions) != 1 {
		return fmt.Errorf("startEvent %q must have exactly one compensateEventDefinition", startEventID)
	}
	other := extraCatchDefinitions(defs) - len(defs.CompensateEventDefinitions)
	if other > 0 {
		return fmt.Errorf("startEvent %q must be a compensate start", startEventID)
	}
	if ref := strings.TrimSpace(defs.CompensateEventDefinitions[0].ActivityRef); ref != "" {
		return fmt.Errorf("startEvent %q compensateEventDefinition activityRef is not supported", startEventID)
	}
	return nil
}

func errorStartCatchFromDefs(startEventID string, defs element.EventDefinitions, errors []element.Error) (string, error) {
	if len(defs.ErrorEventDefinitions) != 1 {
		return "", fmt.Errorf("startEvent %q must have exactly one errorEventDefinition", startEventID)
	}
	other := extraCatchDefinitions(defs) - len(defs.ErrorEventDefinitions)
	if other > 0 {
		return "", fmt.Errorf("startEvent %q must be an error start", startEventID)
	}
	return resolveErrorCode(defs.ErrorEventDefinitions[0].ErrorRef, errors), nil
}

// MatchErrorEventSubProcess returns an armed error event sub-process in scopeID that catches errorCode.
// Prefers an exact code match over a catch-all (empty ErrorCode).
func (d *Deployment) MatchErrorEventSubProcess(scopeID, errorCode string, isArmed func(eventSubProcessID string) bool) (string, bool) {
	var catchAll string
	for _, sp := range d.EventSubProcessesInScope(scopeID) {
		if sp.Kind != CatchKindError {
			continue
		}
		if isArmed != nil && !isArmed(sp.ID) {
			continue
		}
		if sp.ErrorCode == errorCode {
			return sp.ID, true
		}
		if sp.ErrorCode == "" && catchAll == "" {
			catchAll = sp.ID
		}
	}
	if catchAll != "" {
		return catchAll, true
	}
	return "", false
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
