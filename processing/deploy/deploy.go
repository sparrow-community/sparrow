package deploy

import (
	"fmt"

	"github.com/sparrow-community/sparrow/bpmn"
	"github.com/sparrow-community/sparrow/bpmn/element"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// Deployment is an immutable, validated BPMN process ready for execution.
// It holds the bpmn element.Process directly — no parallel Node/Flow model.
type Deployment struct {
	ID      string
	Version int32
	Process element.Process
}

// Compile parses BPMN XML and keeps the first executable process (M1 subset).
func Compile(bpmnXML []byte) (*Deployment, error) {
	model, err := bpmn.BpmnModelelementFromBytes(bpmnXML)
	if err != nil {
		return nil, fmt.Errorf("parse bpmn: %w", err)
	}
	if model.Definitions == nil || len(model.Definitions.Processes) == 0 {
		return nil, fmt.Errorf("no process in definitions")
	}

	var proc *element.Process
	for i := range model.Definitions.Processes {
		p := &model.Definitions.Processes[i]
		if p.IsExecutable || proc == nil {
			proc = p
			if p.IsExecutable {
				break
			}
		}
	}
	if proc == nil {
		return nil, fmt.Errorf("process not found")
	}
	if err := validateM1(proc); err != nil {
		return nil, err
	}
	if _, err := StartEventID(proc); err != nil {
		return nil, err
	}

	return &Deployment{Version: 1, Process: *proc}, nil
}

func validateM1(proc *element.Process) error {
	unsupported := 0
	unsupported += len(proc.Tasks) + len(proc.ServiceTasks) + len(proc.ManualTasks)
	unsupported += len(proc.SendTasks) + len(proc.ReceiveTasks) + len(proc.BusinessRuleTasks)
	unsupported += len(proc.ParallelGatewaies) + len(proc.InclusiveGatewaies) + len(proc.EventBasedGatewaies)
	unsupported += len(proc.SubProcesses) + len(proc.CallActivities) + len(proc.BoundaryEvents)
	unsupported += len(proc.IntermediateCatchEvents) + len(proc.IntermediateThrowEvents)
	if unsupported > 0 {
		return fmt.Errorf("UNSUPPORTED_ELEMENT: process contains elements outside M1 subset")
	}
	if len(proc.StartEvents) == 0 {
		return fmt.Errorf("no startEvent in process")
	}
	return nil
}

func (d *Deployment) ProcessID() string { return d.Process.ID }

func StartEventID(proc *element.Process) (string, error) {
	if len(proc.StartEvents) == 0 {
		return "", fmt.Errorf("no startEvent in process")
	}
	id := proc.StartEvents[0].ID
	outs := Outgoing(proc, id)
	if len(outs) == 0 {
		// fall back to sequenceFlow sourceRef
		for _, f := range proc.SequenceFlows {
			if f.SourceRef == id {
				return id, nil
			}
		}
		return "", fmt.Errorf("startEvent has no outgoing sequence flow")
	}
	return id, nil
}

func (d *Deployment) StartEventID() (string, error) {
	return StartEventID(&d.Process)
}

// TypeOf returns the event.v1 Element.Type for a BPMN element id.
func TypeOf(proc *element.Process, id string) (eventv1.Element_Type, error) {
	for _, e := range proc.StartEvents {
		if e.ID == id {
			return eventv1.Element_TYPE_START_EVENT, nil
		}
	}
	for _, e := range proc.EndEvents {
		if e.ID == id {
			return eventv1.Element_TYPE_END_EVENT, nil
		}
	}
	for _, e := range proc.UserTasks {
		if e.ID == id {
			return eventv1.Element_TYPE_USER_TASK, nil
		}
	}
	for _, e := range proc.ExclusiveGatewaies {
		if e.ID == id {
			return eventv1.Element_TYPE_EXCLUSIVE_GATEWAY, nil
		}
	}
	for _, e := range proc.SequenceFlows {
		if e.ID == id {
			return eventv1.Element_TYPE_SEQUENCE_FLOW, nil
		}
	}
	if proc.ID == id {
		return eventv1.Element_TYPE_PROCESS, nil
	}
	return eventv1.Element_TYPE_UNSPECIFIED, fmt.Errorf("NOT_FOUND: element %q", id)
}

func (d *Deployment) TypeOf(id string) (eventv1.Element_Type, error) {
	return TypeOf(&d.Process, id)
}

// Outgoing returns sequence flow ids leaving the element.
func Outgoing(proc *element.Process, elementID string) []string {
	if outs := flowNodeOutgoing(proc, elementID); len(outs) > 0 {
		return outs
	}
	var ids []string
	for _, f := range proc.SequenceFlows {
		if f.SourceRef == elementID {
			ids = append(ids, f.ID)
		}
	}
	return ids
}

func (d *Deployment) Outgoing(elementID string) []string {
	return Outgoing(&d.Process, elementID)
}

func flowNodeOutgoing(proc *element.Process, id string) []string {
	for _, e := range proc.StartEvents {
		if e.ID == id {
			return append([]string{}, e.Outgoing...)
		}
	}
	for _, e := range proc.EndEvents {
		if e.ID == id {
			return append([]string{}, e.Outgoing...)
		}
	}
	for _, e := range proc.UserTasks {
		if e.ID == id {
			return append([]string{}, e.Outgoing...)
		}
	}
	for _, e := range proc.ExclusiveGatewaies {
		if e.ID == id {
			return append([]string{}, e.Outgoing...)
		}
	}
	return nil
}

// SequenceFlow returns the BPMN sequence flow by id.
func SequenceFlow(proc *element.Process, id string) (element.SequenceFlow, error) {
	for _, f := range proc.SequenceFlows {
		if f.ID == id {
			return f, nil
		}
	}
	return element.SequenceFlow{}, fmt.Errorf("NOT_FOUND: sequence flow %q", id)
}

func (d *Deployment) SequenceFlow(id string) (element.SequenceFlow, error) {
	return SequenceFlow(&d.Process, id)
}

// ExclusiveGateway returns the gateway by id.
func ExclusiveGateway(proc *element.Process, id string) (element.ExclusiveGateway, error) {
	for _, g := range proc.ExclusiveGatewaies {
		if g.ID == id {
			return g, nil
		}
	}
	return element.ExclusiveGateway{}, fmt.Errorf("NOT_FOUND: exclusive gateway %q", id)
}

// ChooseExclusiveOutgoing picks default flow, else first outgoing.
func (d *Deployment) ChooseExclusiveOutgoing(gatewayID string) (string, error) {
	g, err := ExclusiveGateway(&d.Process, gatewayID)
	if err != nil {
		return "", err
	}
	if g.Default != "" {
		if _, err := d.SequenceFlow(g.Default); err == nil {
			return g.Default, nil
		}
	}
	outs := d.Outgoing(gatewayID)
	if len(outs) == 0 {
		return "", fmt.Errorf("NO_OUTGOING_FLOW")
	}
	return outs[0], nil
}
