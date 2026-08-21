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

// IsCompensationHandler reports whether id is a compensation handler activity.
func (d *Deployment) IsCompensationHandler(id string) bool {
	for _, c := range d.compensations {
		if c.HandlerID == id {
			return true
		}
	}
	return false
}

// CompensateActivityRef returns optional activityRef on a compensate throw (empty = all in scope).
func (d *Deployment) CompensateActivityRef(throwID string) (string, error) {
	spec, ok := d.throwEvents[throwID]
	if !ok || spec.Kind != ThrowKindCompensate {
		return "", fmt.Errorf("NOT_FOUND: compensate throw %q", throwID)
	}
	return spec.ActivityRef, nil
}
