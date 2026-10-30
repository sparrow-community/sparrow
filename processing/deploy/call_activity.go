package deploy

import (
	"fmt"
	"strings"

	"github.com/sparrow-community/sparrow/bpmn/element"
)

// MappingAssignment evaluates From against source variables and writes to To.
type MappingAssignment struct {
	From string
	To   string
}

// VariableMapping maps variables between caller and called instance.
// Priority: Assignments, then Transformation into Target, else Source→Target copy.
type VariableMapping struct {
	Source         string
	Target         string
	Transformation string
	Assignments    []MappingAssignment
}

// CallActivity links a callActivity element to a called process.
type CallActivity struct {
	ID              string
	CalledProcessID string
	StartEventID    string
	Inputs          []VariableMapping
	Outputs         []VariableMapping
	// ExternalCallee is true when calledElement is not embedded in the same definitions.
	ExternalCallee bool
	// Opaque is true when calledElement is missing: the Call Activity is a modeling
	// stub on the control-flow path and runs as wait → Complete (like empty SubProcess).
	// Missing callee must not block Deploy.
	Opaque bool
}

func validateCallActivity(ca element.CallActivity, catalog map[string]*element.Process) (CallActivity, error) {
	called := strings.TrimSpace(ca.CalledElement)
	if called == "" {
		// Modeling stub: incomplete IO associations to Data Objects are documentary —
		// they must not block Deploy when there is no callee to map into.
		return CallActivity{ID: ca.ID, Opaque: true}, nil
	}
	inputs, err := compileInputMappings(ca.DataInputAssociations, ca.ID)
	if err != nil {
		return CallActivity{}, err
	}
	outputs, err := compileOutputMappings(ca.DataOutputAssociations, ca.ID)
	if err != nil {
		return CallActivity{}, err
	}
	spec := CallActivity{
		ID:              ca.ID,
		CalledProcessID: called,
		Inputs:          inputs,
		Outputs:         outputs,
	}
	proc, ok := catalog[called]
	if !ok {
		spec.ExternalCallee = true
		return spec, nil
	}
	startID, err := CreateInstanceEntryID(proc)
	if err != nil {
		return CallActivity{}, fmt.Errorf("UNSUPPORTED_ELEMENT: callActivity %q called process %q must have a none start or instantiate event-based gateway: %v", ca.ID, called, err)
	}
	spec.StartEventID = startID
	return spec, nil
}

func compileInputMappings(assocs []element.DataInputAssociation, callID string) ([]VariableMapping, error) {
	if len(assocs) == 0 {
		return nil, nil
	}
	out := make([]VariableMapping, 0, len(assocs))
	for _, a := range assocs {
		m, err := compileAssociationMapping(callID, "input", a.SourceRef, a.TargetRef, a.Transformation.Value, a.Assignments)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, nil
}

func compileOutputMappings(assocs []element.DataOutputAssociation, callID string) ([]VariableMapping, error) {
	if len(assocs) == 0 {
		return nil, nil
	}
	out := make([]VariableMapping, 0, len(assocs))
	for _, a := range assocs {
		m, err := compileAssociationMapping(callID, "output", a.SourceRef, a.TargetRef, a.Transformation.Value, a.Assignments)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, nil
}

func compileAssociationMapping(callID, kind, sourceRef, targetRef, transform string, assignments []element.Assignment) (VariableMapping, error) {
	src := strings.TrimSpace(sourceRef)
	tgt := strings.TrimSpace(targetRef)
	transform = strings.TrimSpace(transform)

	assigns := make([]MappingAssignment, 0, len(assignments))
	for _, a := range assignments {
		from := expressionText(a.From)
		to := mappingVarName(expressionText(a.To))
		if from == "" {
			return VariableMapping{}, fmt.Errorf("UNSUPPORTED_ELEMENT: callActivity %q %s association assignment needs from", callID, kind)
		}
		if to == "" {
			to = tgt
		}
		if to == "" {
			return VariableMapping{}, fmt.Errorf("UNSUPPORTED_ELEMENT: callActivity %q %s association assignment needs to or targetRef", callID, kind)
		}
		assigns = append(assigns, MappingAssignment{From: from, To: to})
	}

	if len(assigns) > 0 {
		if tgt == "" {
			return VariableMapping{}, fmt.Errorf("UNSUPPORTED_ELEMENT: callActivity %q %s association needs targetRef", callID, kind)
		}
		return VariableMapping{Source: src, Target: tgt, Assignments: assigns}, nil
	}
	if transform != "" {
		if tgt == "" {
			return VariableMapping{}, fmt.Errorf("UNSUPPORTED_ELEMENT: callActivity %q %s association transformation needs targetRef", callID, kind)
		}
		return VariableMapping{Source: src, Target: tgt, Transformation: transform}, nil
	}
	if src == "" || tgt == "" {
		return VariableMapping{}, fmt.Errorf("UNSUPPORTED_ELEMENT: callActivity %q %s association needs sourceRef and targetRef", callID, kind)
	}
	return VariableMapping{Source: src, Target: tgt}, nil
}

func mappingVarName(text string) string {
	s := strings.TrimSpace(text)
	if strings.HasPrefix(s, "${") && strings.HasSuffix(s, "}") {
		s = strings.TrimSpace(s[2 : len(s)-1])
	}
	return s
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
