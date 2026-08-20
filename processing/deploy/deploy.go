package deploy

import (
	"fmt"

	"github.com/sparrow-community/sparrow/bpmn"
	"github.com/sparrow-community/sparrow/bpmn/element"
	"github.com/sparrow-community/sparrow/processing/expr"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// Deployment is an immutable, validated BPMN process ready for execution.
// Structure lives on Process; only compiled timer/message facts are cached.
type Deployment struct {
	ID      string
	Version int32
	Process element.Process

	timerCatch   map[string]timerCatch // catch or interrupting timer boundary id
	messageCatch map[string]string     // intermediate catch or interrupting message boundary id -> name
	elements     map[string]*elemEntry // flat index of all elements (recursive into subprocesses)
	seqFlows     map[string]*seqFlowEntry
}

type elemEntry struct {
	Type     eventv1.Element_Type
	ScopeID  string   // parent scope (process ID or subprocess ID)
	Outgoing []string // outgoing sequence flow IDs
	Incoming []string // incoming sequence flow IDs
}

type seqFlowEntry struct {
	ScopeID  string
	SourceID string
	TargetID string
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
	d.compile(model.Definitions.Messages)
	return d, nil
}

func (d *Deployment) compile(messages []element.Message) {
	p := &d.Process
	d.timerCatch = make(map[string]timerCatch, len(p.IntermediateCatchEvents)+len(p.BoundaryEvents))
	d.messageCatch = make(map[string]string, len(p.IntermediateCatchEvents)+len(p.BoundaryEvents))
	d.elements = make(map[string]*elemEntry)
	d.seqFlows = make(map[string]*seqFlowEntry)

	d.indexScope(&p.FlowElements, p.ID, messages)
	d.elements[p.ID] = &elemEntry{Type: eventv1.Element_TYPE_PROCESS, ScopeID: ""}
}

func (d *Deployment) indexScope(fe *element.FlowElements, scopeID string, messages []element.Message) {
	reg := func(id string, typ eventv1.Element_Type, outgoing, incoming []string) {
		d.elements[id] = &elemEntry{Type: typ, ScopeID: scopeID, Outgoing: outgoing, Incoming: incoming}
	}
	for _, e := range fe.StartEvents {
		reg(e.ID, eventv1.Element_TYPE_START_EVENT, e.Outgoing, e.Incoming)
	}
	for _, e := range fe.EndEvents {
		reg(e.ID, eventv1.Element_TYPE_END_EVENT, e.Outgoing, e.Incoming)
	}
	for _, e := range fe.UserTasks {
		reg(e.ID, eventv1.Element_TYPE_USER_TASK, e.Outgoing, e.Incoming)
	}
	for _, e := range fe.ServiceTasks {
		reg(e.ID, eventv1.Element_TYPE_SERVICE_TASK, e.Outgoing, e.Incoming)
	}
	for _, e := range fe.ExclusiveGatewaies {
		reg(e.ID, eventv1.Element_TYPE_EXCLUSIVE_GATEWAY, e.Outgoing, e.Incoming)
	}
	for _, e := range fe.ParallelGatewaies {
		reg(e.ID, eventv1.Element_TYPE_PARALLEL_GATEWAY, e.Outgoing, e.Incoming)
	}
	for _, e := range fe.InclusiveGatewaies {
		reg(e.ID, eventv1.Element_TYPE_INCLUSIVE_GATEWAY, e.Outgoing, e.Incoming)
	}
	for _, e := range fe.EventBasedGatewaies {
		reg(e.ID, eventv1.Element_TYPE_EVENT_BASED_GATEWAY, e.Outgoing, e.Incoming)
	}
	for _, e := range fe.IntermediateCatchEvents {
		reg(e.ID, eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT, e.Outgoing, e.Incoming)
		if spec, err := timerCatchSpec(e); err == nil {
			d.timerCatch[e.ID] = spec
		} else if spec, err := messageCatchSpec(e, messages); err == nil {
			d.messageCatch[e.ID] = spec.Name
		}
	}
	for _, e := range fe.BoundaryEvents {
		reg(e.ID, eventv1.Element_TYPE_BOUNDARY_EVENT, e.Outgoing, e.Incoming)
		if spec, err := timerBoundarySpec(e); err == nil {
			d.timerCatch[e.ID] = spec.Catch
		} else if spec, err := messageBoundarySpec(e, messages); err == nil {
			d.messageCatch[e.ID] = spec.Name
		}
	}
	for _, e := range fe.SequenceFlows {
		d.seqFlows[e.ID] = &seqFlowEntry{ScopeID: scopeID, SourceID: e.SourceRef, TargetID: e.TargetRef}
		d.elements[e.ID] = &elemEntry{Type: eventv1.Element_TYPE_SEQUENCE_FLOW, ScopeID: scopeID}
	}
	for i := range fe.SubProcesses {
		sp := &fe.SubProcesses[i]
		reg(sp.ID, eventv1.Element_TYPE_SUB_PROCESS, sp.Outgoing, sp.Incoming)
		d.indexScope(&sp.FlowElements, sp.ID, messages)
	}
}

func validateM1(proc *element.Process) error {
	unsupported := 0
	unsupported += len(proc.Tasks) + len(proc.ManualTasks)
	unsupported += len(proc.SendTasks) + len(proc.ReceiveTasks) + len(proc.BusinessRuleTasks)
	unsupported += len(proc.CallActivities)
	unsupported += len(proc.IntermediateThrowEvents)
	if unsupported > 0 {
		return fmt.Errorf("UNSUPPORTED_ELEMENT: process contains elements outside M1 subset")
	}
	if err := validateSubProcesses(&proc.FlowElements); err != nil {
		return err
	}
	if err := validateEventBasedGateways(&proc.FlowElements); err != nil {
		return err
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
	type seenKey struct {
		activity string
		kind     string // "timer" or "message"
	}
	seenAttach := make(map[seenKey]string, len(proc.BoundaryEvents))
	for _, e := range proc.BoundaryEvents {
		attached := ""
		kind := ""
		switch {
		case len(e.TimerEventDefinitions) > 0:
			spec, err := timerBoundarySpec(e)
			if err != nil {
				return err
			}
			attached = spec.AttachedTo
			kind = "timer"
		case len(e.MessageEventDefinitions) > 0:
			spec, err := messageBoundarySpec(e, nil)
			if err != nil {
				return err
			}
			attached = spec.AttachedTo
			kind = "message"
		default:
			return fmt.Errorf("UNSUPPORTED_ELEMENT: boundaryEvent %q must be a timer or message boundary", e.ID)
		}
		key := seenKey{attached, kind}
		if prev, ok := seenAttach[key]; ok {
			return fmt.Errorf("UNSUPPORTED_ELEMENT: activity %q already has %s boundary %q", attached, kind, prev)
		}
		if err := validateBoundaryHost(proc, attached); err != nil {
			return err
		}
		seenAttach[key] = e.ID
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

// SubProcessStartEventID returns the start event of the given subprocess.
func (d *Deployment) SubProcessStartEventID(subProcessID string) (string, error) {
	sp := d.findSubProcess(&d.Process.FlowElements, subProcessID)
	if sp == nil {
		return "", fmt.Errorf("NOT_FOUND: subProcess %q", subProcessID)
	}
	if len(sp.StartEvents) == 0 {
		return "", fmt.Errorf("no startEvent in subProcess %q", subProcessID)
	}
	return sp.StartEvents[0].ID, nil
}

func (d *Deployment) findSubProcess(fe *element.FlowElements, id string) *element.SubProcess {
	for i := range fe.SubProcesses {
		if fe.SubProcesses[i].ID == id {
			return &fe.SubProcesses[i]
		}
		if sp := d.findSubProcess(&fe.SubProcesses[i].FlowElements, id); sp != nil {
			return sp
		}
	}
	return nil
}

func (d *Deployment) TypeOf(id string) (eventv1.Element_Type, error) {
	if e, ok := d.elements[id]; ok {
		return e.Type, nil
	}
	return eventv1.Element_TYPE_UNSPECIFIED, fmt.Errorf("NOT_FOUND: element %q", id)
}

// ScopeOf returns the parent scope ID (process or subprocess) of an element.
func (d *Deployment) ScopeOf(id string) (string, bool) {
	if e, ok := d.elements[id]; ok {
		return e.ScopeID, true
	}
	return "", false
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
	if e, ok := d.elements[elementID]; ok && len(e.Outgoing) > 0 {
		return append([]string{}, e.Outgoing...)
	}
	var ids []string
	for id, sf := range d.seqFlows {
		if sf.SourceID == elementID {
			ids = append(ids, id)
		}
	}
	return ids
}

// Incoming returns sequence flow ids entering the element.
func Incoming(proc *element.Process, elementID string) []string {
	if ins := flowNodeIncoming(proc, elementID); len(ins) > 0 {
		return ins
	}
	var ids []string
	for _, f := range proc.SequenceFlows {
		if f.TargetRef == elementID {
			ids = append(ids, f.ID)
		}
	}
	return ids
}

func (d *Deployment) Incoming(elementID string) []string {
	if e, ok := d.elements[elementID]; ok && len(e.Incoming) > 0 {
		return append([]string{}, e.Incoming...)
	}
	var ids []string
	for id, sf := range d.seqFlows {
		if sf.TargetID == elementID {
			ids = append(ids, id)
		}
	}
	return ids
}

func flowNodeIncoming(proc *element.Process, id string) []string {
	for _, e := range proc.StartEvents {
		if e.ID == id {
			return append([]string{}, e.Incoming...)
		}
	}
	for _, e := range proc.EndEvents {
		if e.ID == id {
			return append([]string{}, e.Incoming...)
		}
	}
	for _, e := range proc.UserTasks {
		if e.ID == id {
			return append([]string{}, e.Incoming...)
		}
	}
	for _, e := range proc.ServiceTasks {
		if e.ID == id {
			return append([]string{}, e.Incoming...)
		}
	}
	for _, e := range proc.ExclusiveGatewaies {
		if e.ID == id {
			return append([]string{}, e.Incoming...)
		}
	}
	for _, e := range proc.ParallelGatewaies {
		if e.ID == id {
			return append([]string{}, e.Incoming...)
		}
	}
	for _, e := range proc.IntermediateCatchEvents {
		if e.ID == id {
			return append([]string{}, e.Incoming...)
		}
	}
	for _, e := range proc.BoundaryEvents {
		if e.ID == id {
			return append([]string{}, e.Incoming...)
		}
	}
	return nil
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
	for _, e := range proc.ParallelGatewaies {
		if e.ID == id {
			return append([]string{}, e.Outgoing...)
		}
	}
	for _, e := range proc.IntermediateCatchEvents {
		if e.ID == id {
			return append([]string{}, e.Outgoing...)
		}
	}
	for _, e := range proc.BoundaryEvents {
		if e.ID == id {
			return append([]string{}, e.Outgoing...)
		}
	}
	return nil
}

func (d *Deployment) SequenceFlow(id string) (element.SequenceFlow, error) {
	return d.findSequenceFlow(id)
}

func (d *Deployment) findSequenceFlow(id string) (element.SequenceFlow, error) {
	return findSequenceFlowIn(&d.Process.FlowElements, id)
}

func findSequenceFlowIn(fe *element.FlowElements, id string) (element.SequenceFlow, error) {
	for _, f := range fe.SequenceFlows {
		if f.ID == id {
			return f, nil
		}
	}
	for i := range fe.SubProcesses {
		if f, err := findSequenceFlowIn(&fe.SubProcesses[i].FlowElements, id); err == nil {
			return f, nil
		}
	}
	return element.SequenceFlow{}, fmt.Errorf("NOT_FOUND: sequence flow %q", id)
}

func (d *Deployment) ServiceTask(id string) (element.ServiceTask, error) {
	return findServiceTaskIn(&d.Process.FlowElements, id)
}

func findServiceTaskIn(fe *element.FlowElements, id string) (element.ServiceTask, error) {
	for _, st := range fe.ServiceTasks {
		if st.ID == id {
			return st, nil
		}
	}
	for i := range fe.SubProcesses {
		if st, err := findServiceTaskIn(&fe.SubProcesses[i].FlowElements, id); err == nil {
			return st, nil
		}
	}
	return element.ServiceTask{}, fmt.Errorf("NOT_FOUND: service task %q", id)
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
	g := findExclusiveGatewayIn(&d.Process.FlowElements, gatewayID)
	if g == nil {
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

// ChooseInclusiveOutgoing evaluates all outgoing flows; returns all whose condition is true.
// If none match, the default flow is returned (or error if no default).
func (d *Deployment) ChooseInclusiveOutgoing(gatewayID string, vars map[string]string) ([]string, error) {
	g := findInclusiveGatewayIn(&d.Process.FlowElements, gatewayID)
	if g == nil {
		return nil, fmt.Errorf("NOT_FOUND: inclusive gateway %q", gatewayID)
	}
	var taken []string
	for _, flowID := range d.Outgoing(gatewayID) {
		if g.Default != "" && flowID == g.Default {
			continue
		}
		flow, err := d.SequenceFlow(flowID)
		if err != nil {
			return nil, err
		}
		text := ConditionText(flow)
		if text == "" {
			taken = append(taken, flowID)
			continue
		}
		match, err := expr.Eval(text, vars)
		if err != nil {
			return nil, fmt.Errorf("INVALID_CONDITION: flow %s: %w", flowID, err)
		}
		if match {
			taken = append(taken, flowID)
		}
	}
	if len(taken) == 0 {
		if g.Default != "" {
			if _, err := d.SequenceFlow(g.Default); err == nil {
				return []string{g.Default}, nil
			}
		}
		return nil, fmt.Errorf("NO_OUTGOING_FLOW: inclusive gateway %q", gatewayID)
	}
	return taken, nil
}

func findInclusiveGatewayIn(fe *element.FlowElements, id string) *element.InclusiveGateway {
	for i := range fe.InclusiveGatewaies {
		if fe.InclusiveGatewaies[i].ID == id {
			return &fe.InclusiveGatewaies[i]
		}
	}
	for i := range fe.SubProcesses {
		if g := findInclusiveGatewayIn(&fe.SubProcesses[i].FlowElements, id); g != nil {
			return g
		}
	}
	return nil
}

func findExclusiveGatewayIn(fe *element.FlowElements, id string) *element.ExclusiveGateway {
	for i := range fe.ExclusiveGatewaies {
		if fe.ExclusiveGatewaies[i].ID == id {
			return &fe.ExclusiveGatewaies[i]
		}
	}
	for i := range fe.SubProcesses {
		if g := findExclusiveGatewayIn(&fe.SubProcesses[i].FlowElements, id); g != nil {
			return g
		}
	}
	return nil
}

func validateBoundaryHost(proc *element.Process, activityID string) error {
	for _, e := range proc.UserTasks {
		if e.ID == activityID {
			return nil
		}
	}
	for _, e := range proc.ServiceTasks {
		if e.ID == activityID {
			return nil
		}
	}
	for _, e := range proc.SubProcesses {
		if e.ID == activityID {
			return nil
		}
	}
	return fmt.Errorf("UNSUPPORTED_ELEMENT: boundary must attach to a userTask, serviceTask, or subProcess (%q)", activityID)
}

func validateSubProcesses(fe *element.FlowElements) error {
	for _, sp := range fe.SubProcesses {
		if len(sp.StartEvents) == 0 {
			return fmt.Errorf("UNSUPPORTED_ELEMENT: subProcess %q must have a startEvent", sp.ID)
		}
		if err := validateSubProcesses(&sp.FlowElements); err != nil {
			return err
		}
	}
	return nil
}

func validateEventBasedGateways(fe *element.FlowElements) error {
	catchIDs := make(map[string]bool, len(fe.IntermediateCatchEvents))
	for _, e := range fe.IntermediateCatchEvents {
		catchIDs[e.ID] = true
	}
	for _, g := range fe.EventBasedGatewaies {
		if g.Instantiate {
			return fmt.Errorf("UNSUPPORTED_ELEMENT: eventBasedGateway %q instantiate is not supported", g.ID)
		}
		if g.EventGatewayType == element.EventGatewayTypeParallel {
			return fmt.Errorf("UNSUPPORTED_ELEMENT: eventBasedGateway %q parallel type is not supported", g.ID)
		}
		outs := g.Outgoing
		if len(outs) == 0 {
			for _, f := range fe.SequenceFlows {
				if f.SourceRef == g.ID {
					outs = append(outs, f.ID)
				}
			}
		}
		if len(outs) < 2 {
			return fmt.Errorf("UNSUPPORTED_ELEMENT: eventBasedGateway %q must have at least two outgoing flows", g.ID)
		}
		for _, flowID := range outs {
			var target string
			for _, f := range fe.SequenceFlows {
				if f.ID == flowID {
					target = f.TargetRef
					break
				}
			}
			if target == "" {
				return fmt.Errorf("UNSUPPORTED_ELEMENT: eventBasedGateway %q outgoing flow %q not found", g.ID, flowID)
			}
			if !catchIDs[target] {
				return fmt.Errorf("UNSUPPORTED_ELEMENT: eventBasedGateway %q must target intermediateCatchEvent (%q)", g.ID, target)
			}
		}
	}
	for i := range fe.SubProcesses {
		if err := validateEventBasedGateways(&fe.SubProcesses[i].FlowElements); err != nil {
			return err
		}
	}
	return nil
}

// EventBasedSiblings returns sibling catch element IDs when catchID is a target
// of an exclusive event-based gateway. The gateway id is returned as the first value.
func (d *Deployment) EventBasedSiblings(catchID string) (gatewayID string, siblings []string) {
	if d == nil {
		return "", nil
	}
	for _, flowID := range d.Incoming(catchID) {
		sf, err := d.SequenceFlow(flowID)
		if err != nil {
			continue
		}
		typ, err := d.TypeOf(sf.SourceRef)
		if err != nil || typ != eventv1.Element_TYPE_EVENT_BASED_GATEWAY {
			continue
		}
		gatewayID = sf.SourceRef
		for _, outID := range d.Outgoing(gatewayID) {
			out, err := d.SequenceFlow(outID)
			if err != nil || out.TargetRef == catchID {
				continue
			}
			siblings = append(siblings, out.TargetRef)
		}
		return gatewayID, siblings
	}
	return "", nil
}
