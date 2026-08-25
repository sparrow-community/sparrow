package deploy

import (
	"fmt"
	"strings"

	"github.com/sparrow-community/sparrow/bpmn/element"
)

// CallActivity links a callActivity element to a called process in the same definitions.
type CallActivity struct {
	ID              string
	CalledProcessID string
	StartEventID    string
}

func validateCallActivity(ca element.CallActivity, catalog map[string]*element.Process, claimed map[string]string) (CallActivity, error) {
	called := strings.TrimSpace(ca.CalledElement)
	if called == "" {
		return CallActivity{}, fmt.Errorf("UNSUPPORTED_ELEMENT: callActivity %q needs calledElement", ca.ID)
	}
	proc, ok := catalog[called]
	if !ok {
		return CallActivity{}, fmt.Errorf("UNSUPPORTED_ELEMENT: callActivity %q calledElement %q not found in definitions", ca.ID, called)
	}
	if prev, ok := claimed[called]; ok {
		return CallActivity{}, fmt.Errorf("UNSUPPORTED_ELEMENT: process %q is already called by %q (one CallActivity per called process in v1)", called, prev)
	}
	startID, err := StartEventID(proc)
	if err != nil {
		return CallActivity{}, fmt.Errorf("UNSUPPORTED_ELEMENT: callActivity %q called process %q: %v", ca.ID, called, err)
	}
	if len(proc.CallActivities) > 0 {
		return CallActivity{}, fmt.Errorf("UNSUPPORTED_ELEMENT: called process %q must not contain callActivity (no recursive calls yet)", called)
	}
	claimed[called] = ca.ID
	return CallActivity{ID: ca.ID, CalledProcessID: called, StartEventID: startID}, nil
}

// CallActivitySpec returns the compiled call activity, if any.
func (d *Deployment) CallActivitySpec(id string) (CallActivity, bool) {
	c, ok := d.callActivities[id]
	return c, ok
}

// CallActivityForCalledProcess returns the CallActivity that invokes processID.
func (d *Deployment) CallActivityForCalledProcess(processID string) (CallActivity, bool) {
	id, ok := d.calledProcessOwner[processID]
	if !ok {
		return CallActivity{}, false
	}
	return d.CallActivitySpec(id)
}

// CalledProcessStartEventID returns the start event of the process invoked by callActivityID.
func (d *Deployment) CalledProcessStartEventID(callActivityID string) (string, error) {
	c, ok := d.callActivities[callActivityID]
	if !ok {
		return "", fmt.Errorf("NOT_FOUND: callActivity %q", callActivityID)
	}
	return c.StartEventID, nil
}

// IsCalledProcess reports whether id is a process invoked via CallActivity (not the root).
func (d *Deployment) IsCalledProcess(id string) bool {
	_, ok := d.calledProcessOwner[id]
	return ok
}
