package deploy

import (
	"fmt"

	"github.com/sparrow-community/sparrow/bpmn/element"
)

type errorBoundary struct {
	ErrorCode  string
	AttachedTo string
}

func resolveErrorCode(errorRef string, errors []element.Error) string {
	if errorRef == "" {
		return ""
	}
	for _, e := range errors {
		if e.ID == errorRef {
			if e.ErrorCode != "" {
				return e.ErrorCode
			}
			return e.Name
		}
	}
	return errorRef
}

func errorMatches(boundaryCode, thrownCode string) bool {
	if boundaryCode == "" {
		return true
	}
	return boundaryCode == thrownCode
}

func errorBoundarySpec(ev element.BoundaryEvent, errors []element.Error) (errorBoundary, error) {
	if err := requireBoundaryAttach(ev); err != nil {
		return errorBoundary{}, err
	}
	if len(ev.ErrorEventDefinitions) != 1 {
		return errorBoundary{}, fmt.Errorf("UNSUPPORTED_ELEMENT: boundaryEvent %q must be an error boundary", ev.ID)
	}
	other := extraCatchDefinitions(ev.EventDefinitions) - len(ev.ErrorEventDefinitions)
	if other > 0 {
		return errorBoundary{}, fmt.Errorf("UNSUPPORTED_ELEMENT: boundaryEvent %q must be an error boundary", ev.ID)
	}
	if !ev.CancelActivity {
		return errorBoundary{}, fmt.Errorf("UNSUPPORTED_ELEMENT: error boundary %q must be interrupting (cancelActivity=true)", ev.ID)
	}
	code := resolveErrorCode(ev.ErrorEventDefinitions[0].ErrorRef, errors)
	return errorBoundary{ErrorCode: code, AttachedTo: ev.AttachedToRef}, nil
}

func errorEndSpec(ev element.EndEvent, errors []element.Error) (string, error) {
	if len(ev.ErrorEventDefinitions) != 1 {
		return "", fmt.Errorf("not an error end")
	}
	other := extraCatchDefinitions(ev.EventDefinitions) - len(ev.ErrorEventDefinitions)
	if other > 0 {
		return "", fmt.Errorf("not an error end")
	}
	return resolveErrorCode(ev.ErrorEventDefinitions[0].ErrorRef, errors), nil
}

// MatchErrorBoundary returns a boundary on activityID that catches errorCode.
func (d *Deployment) MatchErrorBoundary(activityID, errorCode string) (string, bool) {
	for _, bid := range d.errorBoundaries[activityID] {
		if errorMatches(d.errorCatch[bid], errorCode) {
			return bid, true
		}
	}
	return "", false
}

// ErrorEndCode reports whether an end event throws an error and its code.
func (d *Deployment) ErrorEndCode(endEventID string) (string, bool) {
	code, ok := d.errorEnds[endEventID]
	return code, ok
}

// IsErrorBoundary reports whether id is an error boundary event.
func (d *Deployment) IsErrorBoundary(boundaryID string) bool {
	_, ok := d.errorCatch[boundaryID]
	return ok
}

// ErrorBoundaryCode returns the error code a boundary catches (empty = catch-all).
func (d *Deployment) ErrorBoundaryCode(boundaryID string) (string, bool) {
	code, ok := d.errorCatch[boundaryID]
	return code, ok
}
