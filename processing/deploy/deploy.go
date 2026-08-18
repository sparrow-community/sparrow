package deploy

import (
	"fmt"

	"github.com/sparrow-community/sparrow/bpmn"
	"github.com/sparrow-community/sparrow/bpmn/element"
	"github.com/sparrow-community/sparrow/processing/expr"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// Deployment is an immutable, validated BPMN process ready for execution.
// It holds the bpmn element.Process directly — no parallel Node/Flow model.
type Deployment struct {
	ID      string
	Version int32
	Process element.Process

	types        map[string]eventv1.Element_Type
	outgoing     map[string][]string
	flows        map[string]element.SequenceFlow
	xor          map[string]element.ExclusiveGateway
	serviceTasks map[string]element.ServiceTask
	timerCatch   map[string]timerCatch
	messageCatch map[string]string // element id -> BPMN message name
	catchKinds   map[string]CatchKind
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

	d := &Deployment{Version: 1, Process: *proc}
	d.buildIndex(model.Definitions.Messages)
	return d, nil
}

func (d *Deployment) buildIndex(messages []element.Message) {
	p := &d.Process
	d.types = map[string]eventv1.Element_Type{p.ID: eventv1.Element_TYPE_PROCESS}
	d.outgoing = make(map[string][]string)
	d.flows = make(map[string]element.SequenceFlow, len(p.SequenceFlows))
	d.xor = make(map[string]element.ExclusiveGateway, len(p.ExclusiveGatewaies))
	d.serviceTasks = make(map[string]element.ServiceTask, len(p.ServiceTasks))
	d.timerCatch = make(map[string]timerCatch, len(p.IntermediateCatchEvents))
	d.messageCatch = make(map[string]string, len(p.IntermediateCatchEvents))
	d.catchKinds = make(map[string]CatchKind, len(p.IntermediateCatchEvents))

	indexNode := func(id string, typ eventv1.Element_Type, outs []string) {
		d.types[id] = typ
		if len(outs) > 0 {
			d.outgoing[id] = append([]string{}, outs...)
		}
	}
	for _, e := range p.StartEvents {
		indexNode(e.ID, eventv1.Element_TYPE_START_EVENT, e.Outgoing)
	}
	for _, e := range p.EndEvents {
		indexNode(e.ID, eventv1.Element_TYPE_END_EVENT, e.Outgoing)
	}
	for _, e := range p.UserTasks {
		indexNode(e.ID, eventv1.Element_TYPE_USER_TASK, e.Outgoing)
	}
	for _, e := range p.ServiceTasks {
		indexNode(e.ID, eventv1.Element_TYPE_SERVICE_TASK, e.Outgoing)
		d.serviceTasks[e.ID] = e
	}
	for _, g := range p.ExclusiveGatewaies {
		indexNode(g.ID, eventv1.Element_TYPE_EXCLUSIVE_GATEWAY, g.Outgoing)
		d.xor[g.ID] = g
	}
	for _, e := range p.IntermediateCatchEvents {
		indexNode(e.ID, eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT, e.Outgoing)
		if spec, err := timerCatchSpec(e); err == nil {
			d.timerCatch[e.ID] = spec
			d.catchKinds[e.ID] = CatchKindTimer
			continue
		}
		if spec, err := messageCatchSpec(e, messages); err == nil {
			d.messageCatch[e.ID] = spec.Name
			d.catchKinds[e.ID] = CatchKindMessage
			continue
		}
	}
	for _, f := range p.SequenceFlows {
		d.types[f.ID] = eventv1.Element_TYPE_SEQUENCE_FLOW
		d.flows[f.ID] = f
		if _, ok := d.outgoing[f.SourceRef]; !ok {
			d.outgoing[f.SourceRef] = append(d.outgoing[f.SourceRef], f.ID)
		}
	}
}

func validateM1(proc *element.Process) error {
	unsupported := 0
	unsupported += len(proc.Tasks) + len(proc.ManualTasks)
	unsupported += len(proc.SendTasks) + len(proc.ReceiveTasks) + len(proc.BusinessRuleTasks)
	unsupported += len(proc.ParallelGatewaies) + len(proc.InclusiveGatewaies) + len(proc.EventBasedGatewaies)
	unsupported += len(proc.SubProcesses) + len(proc.CallActivities) + len(proc.BoundaryEvents)
	unsupported += len(proc.IntermediateThrowEvents)
	if unsupported > 0 {
		return fmt.Errorf("UNSUPPORTED_ELEMENT: process contains elements outside M1 subset")
	}
	for _, e := range proc.IntermediateCatchEvents {
		if len(e.TimerEventDefinitions) > 0 {
			if _, err := timerCatchSpec(e); err != nil {
				return err
			}
			continue
		}
		if _, err := messageCatchSpec(e, nil); err != nil {
			return fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateCatchEvent %q must be either timer catch or message catch", e.ID)
		}
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

func (d *Deployment) TypeOf(id string) (eventv1.Element_Type, error) {
	if t, ok := d.types[id]; ok {
		return t, nil
	}
	return eventv1.Element_TYPE_UNSPECIFIED, fmt.Errorf("NOT_FOUND: element %q", id)
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
	return append([]string{}, d.outgoing[elementID]...)
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
	for _, e := range proc.ServiceTasks {
		if e.ID == id {
			return append([]string{}, e.Outgoing...)
		}
	}
	for _, e := range proc.ExclusiveGatewaies {
		if e.ID == id {
			return append([]string{}, e.Outgoing...)
		}
	}
	for _, e := range proc.IntermediateCatchEvents {
		if e.ID == id {
			return append([]string{}, e.Outgoing...)
		}
	}
	return nil
}

func (d *Deployment) SequenceFlow(id string) (element.SequenceFlow, error) {
	f, ok := d.flows[id]
	if !ok {
		return element.SequenceFlow{}, fmt.Errorf("NOT_FOUND: sequence flow %q", id)
	}
	return f, nil
}

func (d *Deployment) ServiceTask(id string) (element.ServiceTask, error) {
	st, ok := d.serviceTasks[id]
	if !ok {
		return element.ServiceTask{}, fmt.Errorf("NOT_FOUND: service task %q", id)
	}
	return st, nil
}

// ConditionText returns the sequence flow condition body, or empty if none.
func ConditionText(f element.SequenceFlow) string {
	switch e := f.ConditionExpression.ExpressionSubstitution.(type) {
	case *element.FormalExpression:
		return e.Value
	default:
		return ""
	}
}

// ChooseExclusiveOutgoing picks the first matching non-default condition, else default.
func (d *Deployment) ChooseExclusiveOutgoing(gatewayID string, vars map[string]string) (string, error) {
	g, ok := d.xor[gatewayID]
	if !ok {
		return "", fmt.Errorf("NOT_FOUND: exclusive gateway %q", gatewayID)
	}
	for _, flowID := range d.Outgoing(gatewayID) {
		if g.Default != "" && flowID == g.Default {
			continue
		}
		flow, err := d.SequenceFlow(flowID)
		if err != nil {
			return "", err
		}
		text := ConditionText(flow)
		if text == "" {
			return flowID, nil
		}
		match, err := expr.Eval(text, vars)
		if err != nil {
			return "", fmt.Errorf("INVALID_CONDITION: flow %s: %w", flowID, err)
		}
		if match {
			return flowID, nil
		}
	}
	if g.Default != "" {
		if _, err := d.SequenceFlow(g.Default); err == nil {
			return g.Default, nil
		}
	}
	return "", fmt.Errorf("NO_OUTGOING_FLOW")
}
