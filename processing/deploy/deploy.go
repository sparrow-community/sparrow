package deploy

import (
	"fmt"
	"sort"

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

	timerCatch        map[string]timerCatch // catch or interrupting timer boundary id
	messageCatch      map[string]string     // intermediate catch or interrupting message boundary id -> name
	signalCatch       map[string]string     // intermediate signal catch id -> name
	throwEvents       map[string]throwSpec  // intermediate throw id -> spec
	eventSubProcesses map[string]EventSubProcess
	compensations     map[string]Compensation // activity id -> compensation
	errorCatch        map[string]string       // error boundary id -> error code (empty = catch-all)
	errorBoundaries   map[string][]string     // activity id -> error boundary ids
	errorEnds         map[string]string       // error end event id -> error code
	compensateEnds    map[string]string       // compensate end event id -> optional activityRef
	callActivities      map[string]CallActivity // callActivity id -> spec
	calledProcessOwner  map[string]string      // called process id -> callActivity id
	calledProcesses     map[string]element.Process
	multiInstances      map[string]MultiInstanceSpec
	incidentThresholds  map[string]int
	elements            map[string]*elemEntry // flat index of all elements (recursive into subprocesses)
	seqFlows            map[string]*seqFlowEntry
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

// Compile parses BPMN XML and keeps the first executable process as the root.
// Sibling processes in the same definitions may be invoked via CallActivity.
func Compile(bpmnXML []byte) (*Deployment, error) {
	model, err := bpmn.BpmnModelelementFromBytes(bpmnXML)
	if err != nil {
		return nil, fmt.Errorf("parse bpmn: %w", err)
	}
	if model.Definitions == nil || len(model.Definitions.Processes) == 0 {
		return nil, fmt.Errorf("no process in definitions")
	}

	catalog := make(map[string]*element.Process, len(model.Definitions.Processes))
	for i := range model.Definitions.Processes {
		p := &model.Definitions.Processes[i]
		catalog[p.ID] = p
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
	if err := validateM1(proc, model.Definitions.Errors); err != nil {
		return nil, err
	}
	if _, err := StartEventID(proc); err != nil {
		return nil, err
	}

	d := &Deployment{Version: 1, Process: *proc}
	if err := d.compile(model.Definitions.Messages, model.Definitions.Signals, model.Definitions.Errors, catalog); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Deployment) compile(messages []element.Message, signals []element.Signal, errors []element.Error, catalog map[string]*element.Process) error {
	p := &d.Process
	d.timerCatch = make(map[string]timerCatch, len(p.IntermediateCatchEvents)+len(p.BoundaryEvents))
	d.messageCatch = make(map[string]string, len(p.IntermediateCatchEvents)+len(p.BoundaryEvents))
	d.signalCatch = make(map[string]string, len(p.IntermediateCatchEvents)+len(p.BoundaryEvents))
	d.throwEvents = make(map[string]throwSpec, len(p.IntermediateThrowEvents))
	d.eventSubProcesses = make(map[string]EventSubProcess)
	d.compensations = make(map[string]Compensation)
	d.errorCatch = make(map[string]string)
	d.errorBoundaries = make(map[string][]string)
	d.errorEnds = make(map[string]string)
	d.compensateEnds = make(map[string]string)
	d.callActivities = make(map[string]CallActivity)
	d.calledProcessOwner = make(map[string]string)
	d.calledProcesses = make(map[string]element.Process)
	d.multiInstances = make(map[string]MultiInstanceSpec)
	d.elements = make(map[string]*elemEntry)
	d.seqFlows = make(map[string]*seqFlowEntry)

	if err := d.indexScope(&p.FlowElements, p.ID, messages, signals, errors, collectAssociations(p)); err != nil {
		return err
	}
	d.elements[p.ID] = &elemEntry{Type: eventv1.Element_TYPE_PROCESS, ScopeID: ""}
	return d.indexCallActivities(&p.FlowElements, catalog, messages, signals, errors)
}

func collectAssociations(p *element.Process) []element.Association {
	out := append([]element.Association{}, p.Associations...)
	var walk func(fe *element.FlowElements)
	walk = func(fe *element.FlowElements) {
		for i := range fe.SubProcesses {
			sp := &fe.SubProcesses[i]
			out = append(out, sp.Associations...)
			walk(&sp.FlowElements)
		}
	}
	walk(&p.FlowElements)
	return out
}

func (d *Deployment) indexScope(fe *element.FlowElements, scopeID string, messages []element.Message, signals []element.Signal, errors []element.Error, associations []element.Association) error {
	reg := func(id string, typ eventv1.Element_Type, outgoing, incoming []string) {
		d.elements[id] = &elemEntry{Type: typ, ScopeID: scopeID, Outgoing: outgoing, Incoming: incoming}
	}
	for _, e := range fe.StartEvents {
		reg(e.ID, eventv1.Element_TYPE_START_EVENT, e.Outgoing, e.Incoming)
	}
	for _, e := range fe.EndEvents {
		reg(e.ID, eventv1.Element_TYPE_END_EVENT, e.Outgoing, e.Incoming)
		if code, err := errorEndSpec(e, errors); err == nil {
			d.errorEnds[e.ID] = code
		} else if spec, err := compensateEndSpec(e); err == nil {
			d.compensateEnds[e.ID] = spec.ActivityRef
		}
	}
	for _, e := range fe.UserTasks {
		reg(e.ID, eventv1.Element_TYPE_USER_TASK, e.Outgoing, e.Incoming)
		if err := indexMultiInstance(d, e.ID, e.LoopCharacteristicsElements); err != nil {
			return err
		}
	}
	for _, e := range fe.ServiceTasks {
		reg(e.ID, eventv1.Element_TYPE_SERVICE_TASK, e.Outgoing, e.Incoming)
		if err := indexMultiInstance(d, e.ID, e.LoopCharacteristicsElements); err != nil {
			return err
		}
		if err := indexIncidentThreshold(d, e.ID, e.ExtensionElements); err != nil {
			return err
		}
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
		} else if name, err := signalCatchSpec(e, signals); err == nil {
			d.signalCatch[e.ID] = name
		}
	}
	for _, e := range fe.IntermediateThrowEvents {
		reg(e.ID, eventv1.Element_TYPE_INTERMEDIATE_THROW_EVENT, e.Outgoing, e.Incoming)
		if spec, err := throwEventSpec(e, messages, signals); err == nil {
			d.throwEvents[e.ID] = spec
		}
	}
	for _, e := range fe.CallActivities {
		reg(e.ID, eventv1.Element_TYPE_CALL_ACTIVITY, e.Outgoing, e.Incoming)
		if err := indexMultiInstance(d, e.ID, e.LoopCharacteristicsElements); err != nil {
			return err
		}
	}
	for _, e := range fe.BoundaryEvents {
		reg(e.ID, eventv1.Element_TYPE_BOUNDARY_EVENT, e.Outgoing, e.Incoming)
		if spec, err := timerBoundarySpec(e); err == nil {
			d.timerCatch[e.ID] = spec.Catch
		} else if spec, err := messageBoundarySpec(e, messages); err == nil {
			d.messageCatch[e.ID] = spec.Name
		} else if spec, err := signalBoundarySpec(e, signals); err == nil {
			d.signalCatch[e.ID] = spec.Name
		} else if spec, err := compensationBoundarySpec(e, associations); err == nil {
			d.compensations[spec.ActivityID] = spec
		} else if spec, err := errorBoundarySpec(e, errors); err == nil {
			d.errorCatch[e.ID] = spec.ErrorCode
			d.errorBoundaries[spec.AttachedTo] = append(d.errorBoundaries[spec.AttachedTo], e.ID)
		}
	}
	for _, e := range fe.SequenceFlows {
		d.seqFlows[e.ID] = &seqFlowEntry{ScopeID: scopeID, SourceID: e.SourceRef, TargetID: e.TargetRef}
		d.elements[e.ID] = &elemEntry{Type: eventv1.Element_TYPE_SEQUENCE_FLOW, ScopeID: scopeID}
	}
	for i := range fe.SubProcesses {
		sp := &fe.SubProcesses[i]
		reg(sp.ID, eventv1.Element_TYPE_SUB_PROCESS, sp.Outgoing, sp.Incoming)
		if err := indexMultiInstance(d, sp.ID, sp.LoopCharacteristicsElements); err != nil {
			return err
		}
		if sp.TriggeredByEvent {
			spec, err := eventSubProcessStartSpec(sp.ID, sp.StartEvents[0], messages, signals, errors)
			if err == nil {
				spec.ParentScopeID = scopeID
				d.eventSubProcesses[sp.ID] = spec
				switch spec.Kind {
				case CatchKindMessage:
					d.messageCatch[spec.StartEventID] = spec.MessageName
				case CatchKindSignal:
					d.signalCatch[spec.StartEventID] = spec.SignalName
				case CatchKindTimer:
					if tc, err := timerCatchFromDefs(spec.StartEventID, sp.StartEvents[0].EventDefinitions); err == nil {
						d.timerCatch[spec.StartEventID] = tc
					}
				}
			}
		}
		if err := d.indexScope(&sp.FlowElements, sp.ID, messages, signals, errors, associations); err != nil {
			return err
		}
	}
	return nil
}

func (d *Deployment) indexCallActivities(fe *element.FlowElements, catalog map[string]*element.Process, messages []element.Message, signals []element.Signal, errors []element.Error) error {
	var walk func(*element.FlowElements) error
	walk = func(fe *element.FlowElements) error {
		for _, ca := range fe.CallActivities {
			spec, err := validateCallActivity(ca, catalog)
			if err != nil {
				return err
			}
			if spec.CalledProcessID == d.Process.ID {
				return fmt.Errorf("UNSUPPORTED_ELEMENT: callActivity %q cannot call the root process", ca.ID)
			}
			d.callActivities[spec.ID] = spec
			if spec.ExternalCallee {
				continue
			}
			called := catalog[spec.CalledProcessID]
			if err := validateM1(called, errors); err != nil {
				return fmt.Errorf("UNSUPPORTED_ELEMENT: called process %q: %v", spec.CalledProcessID, err)
			}
			if _, exists := d.calledProcessOwner[spec.CalledProcessID]; !exists {
				d.calledProcessOwner[spec.CalledProcessID] = spec.ID
			}
			if _, exists := d.calledProcesses[spec.CalledProcessID]; !exists {
				if _, exists := d.elements[spec.CalledProcessID]; exists {
					return fmt.Errorf("UNSUPPORTED_ELEMENT: called process id %q collides with an existing element", spec.CalledProcessID)
				}
				d.calledProcesses[spec.CalledProcessID] = *called
				if err := d.indexScope(&called.FlowElements, called.ID, messages, signals, errors, collectAssociations(called)); err != nil {
					return err
				}
				d.elements[called.ID] = &elemEntry{Type: eventv1.Element_TYPE_PROCESS, ScopeID: ""}
			}
			// Index CallActivities nested inside the called process (once per process).
			if err := walk(&called.FlowElements); err != nil {
				return err
			}
		}
		for i := range fe.SubProcesses {
			if err := walk(&fe.SubProcesses[i].FlowElements); err != nil {
				return err
			}
		}
		return nil
	}
	return walk(fe)
}

func validateM1(proc *element.Process, errors []element.Error) error {
	unsupported := 0
	unsupported += len(proc.Tasks) + len(proc.ManualTasks)
	unsupported += len(proc.SendTasks) + len(proc.ReceiveTasks) + len(proc.BusinessRuleTasks)
	if unsupported > 0 {
		return fmt.Errorf("UNSUPPORTED_ELEMENT: process contains elements outside M1 subset")
	}
	if err := validateSubProcesses(&proc.FlowElements, errors); err != nil {
		return err
	}
	if err := validateEventBasedGateways(&proc.FlowElements); err != nil {
		return err
	}
	if err := validateCatchAndThrow(&proc.FlowElements); err != nil {
		return err
	}
	type seenKey struct {
		activity string
		kind     string // "timer", "message", "signal", or "compensate"
	}
	seenAttach := make(map[seenKey]string, len(proc.BoundaryEvents))
	assocs := collectAssociations(proc)
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
		case len(e.SignalEventDefinitions) > 0:
			spec, err := signalBoundarySpec(e, nil)
			if err != nil {
				return err
			}
			attached = spec.AttachedTo
			kind = "signal"
		case len(e.CompensateEventDefinitions) > 0:
			spec, err := compensationBoundarySpec(e, assocs)
			if err != nil {
				return err
			}
			if err := validateCompensationHandler(&proc.FlowElements, spec.HandlerID); err != nil {
				return err
			}
			attached = spec.ActivityID
			kind = "compensate"
		case len(e.ErrorEventDefinitions) > 0:
			spec, err := errorBoundarySpec(e, nil)
			if err != nil {
				return err
			}
			attached = spec.AttachedTo
			kind = "error:" + spec.ErrorCode
		default:
			return fmt.Errorf("UNSUPPORTED_ELEMENT: boundaryEvent %q must be a timer, message, signal, error, or compensation boundary", e.ID)
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
		if _, err := instantiateEntryID(proc); err != nil {
			return fmt.Errorf("no startEvent in process")
		}
		return nil
	}
	for _, g := range proc.EventBasedGatewaies {
		if g.Instantiate {
			return fmt.Errorf("UNSUPPORTED_ELEMENT: process cannot combine startEvent with instantiate eventBasedGateway %q", g.ID)
		}
	}
	return nil
}

func validateCatchAndThrow(fe *element.FlowElements) error {
	for _, e := range fe.IntermediateCatchEvents {
		if len(e.TimerEventDefinitions) > 0 {
			if _, err := timerCatchSpec(e); err != nil {
				return err
			}
			continue
		}
		if _, err := messageCatchSpec(e, nil); err == nil {
			continue
		}
		if _, err := signalCatchSpec(e, nil); err == nil {
			continue
		}
		return fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateCatchEvent %q must be timer, message, or signal catch", e.ID)
	}
	for _, e := range fe.IntermediateThrowEvents {
		if _, err := throwEventSpec(e, nil, nil); err != nil {
			return err
		}
	}
	for i := range fe.SubProcesses {
		if err := validateCatchAndThrow(&fe.SubProcesses[i].FlowElements); err != nil {
			return err
		}
	}
	return nil
}

func (d *Deployment) ProcessID() string { return d.Process.ID }

func StartEventID(proc *element.Process) (string, error) {
	if len(proc.StartEvents) == 0 {
		return instantiateEntryID(proc)
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

// instantiateEntryID returns the sole process-level exclusive instantiate
// event-based gateway when the process has no startEvent.
func instantiateEntryID(proc *element.Process) (string, error) {
	var found string
	for _, g := range proc.EventBasedGatewaies {
		if !g.Instantiate {
			continue
		}
		if g.EventGatewayType == element.EventGatewayTypeParallel {
			return "", fmt.Errorf("UNSUPPORTED_ELEMENT: eventBasedGateway %q parallel instantiate is not supported", g.ID)
		}
		if len(g.Incoming) > 0 {
			return "", fmt.Errorf("UNSUPPORTED_ELEMENT: eventBasedGateway %q instantiate must have no incoming sequence flow", g.ID)
		}
		for _, f := range proc.SequenceFlows {
			if f.TargetRef == g.ID {
				return "", fmt.Errorf("UNSUPPORTED_ELEMENT: eventBasedGateway %q instantiate must have no incoming sequence flow", g.ID)
			}
		}
		outs := g.Outgoing
		if len(outs) == 0 {
			for _, f := range proc.SequenceFlows {
				if f.SourceRef == g.ID {
					outs = append(outs, f.ID)
				}
			}
		}
		if len(outs) < 2 {
			return "", fmt.Errorf("UNSUPPORTED_ELEMENT: eventBasedGateway %q must have at least two outgoing flows", g.ID)
		}
		if found != "" {
			return "", fmt.Errorf("UNSUPPORTED_ELEMENT: process has multiple instantiate eventBasedGateway entries")
		}
		found = g.ID
	}
	if found == "" {
		return "", fmt.Errorf("no startEvent in process")
	}
	return found, nil
}

func (d *Deployment) StartEventID() (string, error) {
	return StartEventID(&d.Process)
}

// SubProcessStartEventID returns the start event of the given subprocess.
func (d *Deployment) SubProcessStartEventID(subProcessID string) (string, error) {
	sp := d.findSubProcess(&d.Process.FlowElements, subProcessID)
	if sp == nil {
		for _, called := range d.calledProcesses {
			fe := called.FlowElements
			if sp = d.findSubProcess(&fe, subProcessID); sp != nil {
				break
			}
		}
	}
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
	sort.Strings(ids)
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
	sort.Strings(ids)
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
	sort.Strings(ids)
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
	if f, err := d.findSequenceFlow(id); err == nil {
		return f, nil
	}
	for _, proc := range d.calledProcesses {
		if f, err := findSequenceFlowIn(&proc.FlowElements, id); err == nil {
			return f, nil
		}
	}
	return element.SequenceFlow{}, fmt.Errorf("NOT_FOUND: sequence flow %q", id)
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
	for _, e := range proc.CallActivities {
		if e.ID == activityID {
			return nil
		}
	}
	return fmt.Errorf("UNSUPPORTED_ELEMENT: boundary must attach to a userTask, serviceTask, subProcess, or callActivity (%q)", activityID)
}

func validateSubProcesses(fe *element.FlowElements, errors []element.Error) error {
	return validateSubProcessesAt(fe, false, false, errors)
}

func validateSubProcessesAt(fe *element.FlowElements, insideEmbedded, insideEventSubProcess bool, errors []element.Error) error {
	for i := range fe.SubProcesses {
		sp := &fe.SubProcesses[i]
		if sp.TriggeredByEvent {
			if insideEventSubProcess {
				return fmt.Errorf("UNSUPPORTED_ELEMENT: event subProcess %q nested in event subProcess is not supported", sp.ID)
			}
			if err := validateEventSubProcess(sp, nil, nil, errors); err != nil {
				return err
			}
			for _, f := range fe.SequenceFlows {
				if f.SourceRef == sp.ID || f.TargetRef == sp.ID {
					return fmt.Errorf("UNSUPPORTED_ELEMENT: event subProcess %q must not have sequence flow connections", sp.ID)
				}
			}
			if err := validateSubProcessesAt(&sp.FlowElements, false, true, errors); err != nil {
				return err
			}
			continue
		}
		if len(sp.StartEvents) == 0 {
			return fmt.Errorf("UNSUPPORTED_ELEMENT: subProcess %q must have a startEvent", sp.ID)
		}
		if err := validateSubProcessesAt(&sp.FlowElements, true, insideEventSubProcess, errors); err != nil {
			return err
		}
	}
	return nil
}

func validateEventBasedGateways(fe *element.FlowElements) error {
	return validateEventBasedGatewaysAt(fe, false)
}

func validateEventBasedGatewaysAt(fe *element.FlowElements, insideSubProcess bool) error {
	catchIDs := make(map[string]bool, len(fe.IntermediateCatchEvents))
	for _, e := range fe.IntermediateCatchEvents {
		catchIDs[e.ID] = true
	}
	for _, g := range fe.EventBasedGatewaies {
		if g.Instantiate {
			if insideSubProcess {
				return fmt.Errorf("UNSUPPORTED_ELEMENT: eventBasedGateway %q instantiate inside subProcess is not supported", g.ID)
			}
			if g.EventGatewayType == element.EventGatewayTypeParallel {
				return fmt.Errorf("UNSUPPORTED_ELEMENT: eventBasedGateway %q parallel instantiate is not supported", g.ID)
			}
			if len(g.Incoming) > 0 {
				return fmt.Errorf("UNSUPPORTED_ELEMENT: eventBasedGateway %q instantiate must have no incoming sequence flow", g.ID)
			}
			for _, f := range fe.SequenceFlows {
				if f.TargetRef == g.ID {
					return fmt.Errorf("UNSUPPORTED_ELEMENT: eventBasedGateway %q instantiate must have no incoming sequence flow", g.ID)
				}
			}
		}
		// Exclusive (default) and Parallel intermediate event-based gateways are supported.
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
		if err := validateEventBasedGatewaysAt(&fe.SubProcesses[i].FlowElements, true); err != nil {
			return err
		}
	}
	return nil
}

// EventBasedSiblings returns sibling catch element IDs when catchID is a target
// of an exclusive event-based gateway. Parallel event-based gateways do not
// cancel siblings (returns empty). The gateway id is returned as the first value.
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
		if d.IsParallelEventBasedGateway(gatewayID) {
			return gatewayID, nil
		}
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

// IsParallelEventBasedGateway reports eventGatewayType="Parallel".
// Absent / empty / Exclusive means exclusive (BPMN default).
func (d *Deployment) IsParallelEventBasedGateway(id string) bool {
	if d == nil {
		return false
	}
	return isParallelEventBasedGatewayIn(&d.Process.FlowElements, id)
}

func isParallelEventBasedGatewayIn(fe *element.FlowElements, id string) bool {
	for _, g := range fe.EventBasedGatewaies {
		if g.ID != id {
			continue
		}
		return g.EventGatewayType == element.EventGatewayTypeParallel
	}
	for i := range fe.SubProcesses {
		if isParallelEventBasedGatewayIn(&fe.SubProcesses[i].FlowElements, id) {
			return true
		}
	}
	return false
}
