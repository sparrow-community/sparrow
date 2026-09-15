package deploy

import (
	"fmt"
	"strings"

	"github.com/sparrow-community/sparrow/bpmn/element"
)

// Compensation links a completed activity to its compensation handler via a
// compensation boundary event and an association.
type Compensation struct {
	BoundaryID string
	ActivityID string
	HandlerID  string
}

func compensateThrowSpec(ev element.IntermediateThrowEvent) (throwSpec, error) {
	if len(ev.TimerEventDefinitions) > 0 {
		return throwSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateThrowEvent %q cannot have timerEventDefinition", ev.ID)
	}
	if len(ev.CompensateEventDefinitions) != 1 {
		return throwSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateThrowEvent %q needs exactly one compensateEventDefinition", ev.ID)
	}
	other := extraCatchDefinitions(ev.EventDefinitions) - len(ev.CompensateEventDefinitions)
	if other > 0 {
		return throwSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateThrowEvent %q is not a compensate throw", ev.ID)
	}
	return throwSpec{
		Kind:        ThrowKindCompensate,
		ActivityRef: strings.TrimSpace(ev.CompensateEventDefinitions[0].ActivityRef),
	}, nil
}

func compensateEndSpec(ev element.EndEvent) (throwSpec, error) {
	if len(ev.CompensateEventDefinitions) != 1 {
		return throwSpec{}, fmt.Errorf("not a compensate end")
	}
	other := extraCatchDefinitions(ev.EventDefinitions) - len(ev.CompensateEventDefinitions)
	if other > 0 {
		return throwSpec{}, fmt.Errorf("not a compensate end")
	}
	return throwSpec{
		Kind:        ThrowKindCompensate,
		ActivityRef: strings.TrimSpace(ev.CompensateEventDefinitions[0].ActivityRef),
	}, nil
}

func compensationBoundarySpec(ev element.BoundaryEvent, associations []element.Association) (Compensation, error) {
	if err := requireBoundaryAttach(ev); err != nil {
		return Compensation{}, err
	}
	if len(ev.CompensateEventDefinitions) != 1 {
		return Compensation{}, fmt.Errorf("UNSUPPORTED_ELEMENT: boundaryEvent %q must be a compensation boundary", ev.ID)
	}
	other := extraCatchDefinitions(ev.EventDefinitions) - len(ev.CompensateEventDefinitions)
	if other > 0 || len(ev.TimerEventDefinitions) > 0 {
		return Compensation{}, fmt.Errorf("UNSUPPORTED_ELEMENT: boundaryEvent %q must be a compensation boundary", ev.ID)
	}
	handler := ""
	for _, a := range associations {
		if a.SourceRef == ev.ID {
			handler = a.TargetRef
			break
		}
	}
	if handler == "" {
		return Compensation{}, fmt.Errorf("UNSUPPORTED_ELEMENT: compensation boundary %q needs an association to a handler", ev.ID)
	}
	return Compensation{
		BoundaryID: ev.ID,
		ActivityID: ev.AttachedToRef,
		HandlerID:  handler,
	}, nil
}

func validateCompensationHandler(fe *element.FlowElements, handlerID string) error {
	for _, t := range fe.UserTasks {
		if t.ID == handlerID {
			return nil
		}
	}
	for _, t := range fe.ServiceTasks {
		if t.ID == handlerID {
			return nil
		}
	}
	for i := range fe.SubProcesses {
		if err := validateCompensationHandler(&fe.SubProcesses[i].FlowElements, handlerID); err == nil {
			return nil
		}
	}
	return fmt.Errorf("UNSUPPORTED_ELEMENT: compensation handler %q must be a userTask or serviceTask", handlerID)
}

// CompensationOf returns the compensation handler linked to activityID, if any.
func (d *Deployment) CompensationOf(activityID string) (Compensation, bool) {
	c, ok := d.compensations[activityID]
	return c, ok
}

// CompensationByBoundary looks up compensation by boundary event id.
func (d *Deployment) CompensationByBoundary(boundaryID string) (Compensation, bool) {
	for _, c := range d.compensations {
		if c.BoundaryID == boundaryID {
			return c, true
		}
	}
	return Compensation{}, false
}

// IsCompensationHandler reports whether id is a compensation handler activity
// (isForCompensation task) or a compensation event sub-process.
func (d *Deployment) IsCompensationHandler(id string) bool {
	for _, c := range d.compensations {
		if c.HandlerID == id {
			return true
		}
	}
	return false
}

// registerCompensationEventSubProcess links an enclosing SubProcess to its
// compensation event sub-process handler (subscription key = compensate start id).
func (d *Deployment) registerCompensationEventSubProcess(parentScopeID string, spec EventSubProcess, processID string) error {
	if parentScopeID == "" || parentScopeID == processID {
		return fmt.Errorf("UNSUPPORTED_ELEMENT: compensation event subProcess %q must be nested in an embedded subProcess", spec.ID)
	}
	if existing, ok := d.compensations[parentScopeID]; ok {
		return fmt.Errorf("UNSUPPORTED_ELEMENT: subProcess %q already has compensation handler %q; cannot also use compensation event subProcess %q", parentScopeID, existing.HandlerID, spec.ID)
	}
	for _, other := range d.eventSubProcesses {
		if other.Kind == CatchKindCompensate && other.ParentScopeID == parentScopeID && other.ID != spec.ID {
			return fmt.Errorf("UNSUPPORTED_ELEMENT: subProcess %q already has compensation event subProcess %q", parentScopeID, other.ID)
		}
	}
	d.compensations[parentScopeID] = Compensation{
		BoundaryID: spec.StartEventID,
		ActivityID: parentScopeID,
		HandlerID:  spec.ID,
	}
	return nil
}

// CompensateActivityRef returns optional activityRef on a compensate throw/end (empty = all in scope).
func (d *Deployment) CompensateActivityRef(throwID string) (string, error) {
	if spec, ok := d.throwEvents[throwID]; ok && spec.Kind == ThrowKindCompensate {
		return spec.ActivityRef, nil
	}
	if ref, ok := d.compensateEnds[throwID]; ok {
		return ref, nil
	}
	return "", fmt.Errorf("NOT_FOUND: compensate throw %q", throwID)
}

// IsCompensateEnd reports whether id is a compensate end event.
func (d *Deployment) IsCompensateEnd(id string) bool {
	_, ok := d.compensateEnds[id]
	return ok
}
