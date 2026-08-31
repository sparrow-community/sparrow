package deploy

import (
	"fmt"
	"strings"

	"github.com/sparrow-community/sparrow/bpmn/element"
)

// VariableMapping is a name-to-name copy between caller and called instance variables.
type VariableMapping struct {
	Source string
	Target string
}

// CallActivity links a callActivity element to a called process in the same definitions.
type CallActivity struct {
	ID              string
	CalledProcessID string
	StartEventID    string
	Inputs          []VariableMapping
	Outputs         []VariableMapping
}

func validateCallActivity(ca element.CallActivity, catalog map[string]*element.Process) (CallActivity, error) {
	called := strings.TrimSpace(ca.CalledElement)
	if called == "" {
		return CallActivity{}, fmt.Errorf("UNSUPPORTED_ELEMENT: callActivity %q needs calledElement", ca.ID)
	}
	proc, ok := catalog[called]
	if !ok {
		return CallActivity{}, fmt.Errorf("UNSUPPORTED_ELEMENT: callActivity %q calledElement %q not found in definitions", ca.ID, called)
	}
	startID, err := StartEventID(proc)
	if err != nil {
		return CallActivity{}, fmt.Errorf("UNSUPPORTED_ELEMENT: callActivity %q called process %q: %v", ca.ID, called, err)
	}
	if len(ca.MultielementLoopCharacteristics) > 0 {
		return CallActivity{}, fmt.Errorf("UNSUPPORTED_ELEMENT: multi-instance callActivity %q not supported", ca.ID)
	}
	inputs, err := compileInputMappings(ca.DataInputAssociations, ca.ID)
	if err != nil {
		return CallActivity{}, err
	}
	outputs, err := compileOutputMappings(ca.DataOutputAssociations, ca.ID)
	if err != nil {
		return CallActivity{}, err
	}
	return CallActivity{
		ID:              ca.ID,
		CalledProcessID: called,
		StartEventID:    startID,
		Inputs:          inputs,
		Outputs:         outputs,
	}, nil
}

func compileInputMappings(assocs []element.DataInputAssociation, callID string) ([]VariableMapping, error) {
	if len(assocs) == 0 {
		return nil, nil
	}
	out := make([]VariableMapping, 0, len(assocs))
	for _, a := range assocs {
		src := strings.TrimSpace(a.SourceRef)
		tgt := strings.TrimSpace(a.TargetRef)
		if src == "" || tgt == "" {
			return nil, fmt.Errorf("UNSUPPORTED_ELEMENT: callActivity %q input association needs sourceRef and targetRef", callID)
		}
		if strings.TrimSpace(a.Transformation.Value) != "" || len(a.Assignments) > 0 {
			return nil, fmt.Errorf("UNSUPPORTED_ELEMENT: callActivity %q input association transformation/assignment not supported", callID)
		}
		out = append(out, VariableMapping{Source: src, Target: tgt})
	}
	return out, nil
}

func compileOutputMappings(assocs []element.DataOutputAssociation, callID string) ([]VariableMapping, error) {
	if len(assocs) == 0 {
		return nil, nil
	}
	out := make([]VariableMapping, 0, len(assocs))
	for _, a := range assocs {
		src := strings.TrimSpace(a.SourceRef)
		tgt := strings.TrimSpace(a.TargetRef)
		if src == "" || tgt == "" {
			return nil, fmt.Errorf("UNSUPPORTED_ELEMENT: callActivity %q output association needs sourceRef and targetRef", callID)
		}
		if strings.TrimSpace(a.Transformation.Value) != "" || len(a.Assignments) > 0 {
			return nil, fmt.Errorf("UNSUPPORTED_ELEMENT: callActivity %q output association transformation/assignment not supported", callID)
		}
		out = append(out, VariableMapping{Source: src, Target: tgt})
	}
	return out, nil
}

// CallActivitySpec returns the compiled call activity, if any.
func (d *Deployment) CallActivitySpec(id string) (CallActivity, bool) {
	c, ok := d.callActivities[id]
	return c, ok
}

// CallActivityForCalledProcess returns a CallActivity that invokes processID, if any.
// With multiple CallActivities naming the same process, the first compiled wins.
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
	_, ok := d.calledProcesses[id]
	return ok
}

// CalledProcess returns the sibling process definition invoked by CallActivity.
func (d *Deployment) CalledProcess(id string) (element.Process, bool) {
	p, ok := d.calledProcesses[id]
	return p, ok
}
