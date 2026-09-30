package deploy

import (
	"fmt"
	"sort"
	"strings"

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
	// SourceXML is the BPMN bytes that compiled this deployment (for consumers).
	SourceXML []byte
	Process   element.Process

	timerCatch           map[string]timerCatch // catch or interrupting timer boundary id
	messageCatch         map[string]string     // intermediate catch or interrupting message boundary id -> name
	signalCatch          map[string]string     // intermediate signal catch id -> name
	throwEvents          map[string]throwSpec  // intermediate throw id -> spec
	eventSubProcesses    map[string]EventSubProcess
	compensations        map[string]Compensation // activity id -> compensation
	errorCatch           map[string]string       // error boundary id -> error code (empty = catch-all)
	errorBoundaries      map[string][]string     // activity id -> error boundary ids
	errorEnds            map[string]string       // error end event id -> error code
	compensateEnds       map[string]string       // compensate end event id -> optional activityRef
	cancelEnds           map[string]bool         // cancel end event ids
	cancelBoundaries     map[string]string       // transaction id -> cancel boundary id
	terminateEnds        map[string]bool         // terminate end event ids
	messageEnds          map[string]string       // message end id -> message name
	signalEnds           map[string]string       // signal end id -> signal name
	escalationCatch      map[string]string       // escalation boundary id -> code (empty = catch-all)
	escalationBoundaries map[string][]string     // activity id -> escalation boundary ids
	escalationEnds       map[string]string       // escalation end event id -> code
	linkCatch            map[string]string       // intermediate link catch id -> link name
	callActivities       map[string]CallActivity // callActivity id -> spec
	calledProcessOwner   map[string]string       // called process id -> callActivity id
	calledProcesses      map[string]element.Process
	multiInstances       map[string]MultiInstanceSpec
	standardLoops        map[string]StandardLoopSpec
	incidentThresholds   map[string]int
	complexGateways      map[string]ComplexGatewaySpec
	adHocSubProcesses    map[string]AdHocSubProcessSpec
	elements             map[string]*elemEntry // flat index of all elements (recursive into subprocesses)
	seqFlows             map[string]*seqFlowEntry

	noneStartID       string
	messageStarts     map[string][]string // message name -> start event ids
	signalStarts      map[string][]string // signal name -> start event ids
	timerStarts       map[string]timerCatch
	conditionalStarts map[string]string // start event id -> condition text
	conditionalCatch  map[string]string // intermediate/boundary conditional id -> condition text

	messages []element.Message
	signals  []element.Signal
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
	if err := validateProcess(proc, model.Definitions.Messages, model.Definitions.Signals, model.Definitions.Errors, model.Definitions.Escalations); err != nil {
		return nil, err
	}
	if _, err := StartEventID(proc); err != nil {
		return nil, err
	}

	d := &Deployment{Version: 1, Process: *proc}
	if err := d.compile(model.Definitions.Messages, model.Definitions.Signals, model.Definitions.Errors, model.Definitions.Escalations, catalog); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Deployment) compile(messages []element.Message, signals []element.Signal, errors []element.Error, escalations []element.Escalation, catalog map[string]*element.Process) error {
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
	d.cancelEnds = make(map[string]bool)
	d.cancelBoundaries = make(map[string]string)
	d.terminateEnds = make(map[string]bool)
	d.messageEnds = make(map[string]string)
	d.signalEnds = make(map[string]string)
	d.escalationCatch = make(map[string]string)
	d.escalationBoundaries = make(map[string][]string)
	d.escalationEnds = make(map[string]string)
	d.linkCatch = make(map[string]string)
	d.callActivities = make(map[string]CallActivity)
	d.calledProcessOwner = make(map[string]string)
	d.calledProcesses = make(map[string]element.Process)
	d.multiInstances = make(map[string]MultiInstanceSpec)
	d.complexGateways = make(map[string]ComplexGatewaySpec)
	d.adHocSubProcesses = make(map[string]AdHocSubProcessSpec)
	d.elements = make(map[string]*elemEntry)
	d.seqFlows = make(map[string]*seqFlowEntry)
	d.messageStarts = make(map[string][]string)
	d.signalStarts = make(map[string][]string)
	d.timerStarts = make(map[string]timerCatch)
	d.conditionalStarts = make(map[string]string)
	d.conditionalCatch = make(map[string]string)
	d.messages = messages
	d.signals = signals

	if err := d.indexScope(&p.FlowElements, p.ID, messages, signals, errors, escalations, collectAssociations(p)); err != nil {
		return err
	}
	if err := d.indexProcessLevelStarts(p, messages, signals); err != nil {
		return err
	}
	d.elements[p.ID] = &elemEntry{Type: eventv1.Element_TYPE_PROCESS, ScopeID: ""}
	if err := d.validateLinkPairs(); err != nil {
		return err
	}
	if err := d.validateCancelSemantics(); err != nil {
		return err
	}
	if err := d.indexCallActivities(&p.FlowElements, catalog, messages, signals, errors, escalations); err != nil {
		return err
	}
	return d.validateFlowReferences()
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
		for i := range fe.Transactions {
			tx := &fe.Transactions[i]
			out = append(out, tx.Associations...)
			walk(&tx.FlowElements)
		}
	}
	walk(&p.FlowElements)
	return out
}

func (d *Deployment) indexScope(fe *element.FlowElements, scopeID string, messages []element.Message, signals []element.Signal, errors []element.Error, escalations []element.Escalation, associations []element.Association) error {
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
		} else if code, err := escalationEndSpec(e, escalations); err == nil {
			d.escalationEnds[e.ID] = code
		} else if spec, err := compensateEndSpec(e); err == nil {
			d.compensateEnds[e.ID] = spec.ActivityRef
		} else if err := cancelEndSpec(e); err == nil {
			d.cancelEnds[e.ID] = true
		} else if err := terminateEndSpec(e); err == nil {
			d.terminateEnds[e.ID] = true
		} else if name, err := messageEndSpec(e, messages); err == nil {
			d.messageEnds[e.ID] = name
		} else if name, err := signalEndSpec(e, signals); err == nil {
			d.signalEnds[e.ID] = name
		} else if endHasEventDefinitions(e) {
			return fmt.Errorf("UNSUPPORTED_ELEMENT: endEvent %q has unsupported event definitions", e.ID)
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
	for _, e := range fe.Tasks {
		reg(e.ID, eventv1.Element_TYPE_TASK, e.Outgoing, e.Incoming)
		if err := indexMultiInstance(d, e.ID, e.LoopCharacteristicsElements); err != nil {
			return err
		}
	}
	for _, e := range fe.ManualTasks {
		reg(e.ID, eventv1.Element_TYPE_MANUAL_TASK, e.Outgoing, e.Incoming)
		if err := indexMultiInstance(d, e.ID, e.LoopCharacteristicsElements); err != nil {
			return err
		}
	}
	for _, e := range fe.ReceiveTasks {
		if e.Instantiate && scopeID != d.Process.ID {
			return fmt.Errorf("UNSUPPORTED_ELEMENT: receiveTask %q instantiate inside subProcess is not supported", e.ID)
		}
		reg(e.ID, eventv1.Element_TYPE_RECEIVE_TASK, e.Outgoing, e.Incoming)
		if err := indexMultiInstance(d, e.ID, e.LoopCharacteristicsElements); err != nil {
			return err
		}
		d.messageCatch[e.ID] = resolveTaskMessageName(e.MessageRef, e.Name, e.ID, messages)
	}
	for _, e := range fe.SendTasks {
		reg(e.ID, eventv1.Element_TYPE_SEND_TASK, e.Outgoing, e.Incoming)
		if err := indexMultiInstance(d, e.ID, e.LoopCharacteristicsElements); err != nil {
			return err
		}
		d.throwEvents[e.ID] = throwSpec{Kind: ThrowKindMessage, Name: resolveTaskMessageName(e.MessageRef, e.Name, e.ID, messages)}
	}
	for _, e := range fe.BusinessRuleTasks {
		reg(e.ID, eventv1.Element_TYPE_BUSINESS_RULE_TASK, e.Outgoing, e.Incoming)
		if err := indexMultiInstance(d, e.ID, e.LoopCharacteristicsElements); err != nil {
			return err
		}
		if err := indexIncidentThreshold(d, e.ID, e.ExtensionElements); err != nil {
			return err
		}
	}
	for _, e := range fe.ScriptTasks {
		reg(e.ID, eventv1.Element_TYPE_SCRIPT_TASK, e.Outgoing, e.Incoming)
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
	for _, e := range fe.ComplexGatewaies {
		reg(e.ID, eventv1.Element_TYPE_COMPLEX_GATEWAY, e.Outgoing, e.Incoming)
		d.complexGateways[e.ID] = complexGatewaySpec(e)
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
		} else if name, err := linkCatchSpec(e); err == nil {
			d.linkCatch[e.ID] = name
		} else if text, err := conditionalCatchSpec(e); err == nil {
			d.conditionalCatch[e.ID] = text
		} else if code, err := escalationCatchSpec(e, escalations); err == nil {
			d.escalationCatch[e.ID] = code
		}
	}
	for _, e := range fe.IntermediateThrowEvents {
		reg(e.ID, eventv1.Element_TYPE_INTERMEDIATE_THROW_EVENT, e.Outgoing, e.Incoming)
		if spec, err := throwEventSpec(e, messages, signals, escalations); err == nil {
			d.throwEvents[e.ID] = spec
		}
	}
	if len(fe.ImplicitThrowEvents) > 0 {
		// implicitThrowEvent has no token semantics as a flow element; BPMN uses it
		// for multi-instance behavior events inside multiInstanceLoopCharacteristics.
		return fmt.Errorf("UNSUPPORTED_ELEMENT: implicitThrowEvent %q is not supported as a flow element (only inside multi-instance behavior definitions)", fe.ImplicitThrowEvents[0].ID)
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
		} else if attached, err := cancelBoundarySpec(e); err == nil {
			if prev, ok := d.cancelBoundaries[attached]; ok {
				return fmt.Errorf("UNSUPPORTED_ELEMENT: transaction %q already has cancel boundary %q", attached, prev)
			}
			d.cancelBoundaries[attached] = e.ID
		} else if spec, err := errorBoundarySpec(e, errors); err == nil {
			d.errorCatch[e.ID] = spec.ErrorCode
			d.errorBoundaries[spec.AttachedTo] = append(d.errorBoundaries[spec.AttachedTo], e.ID)
		} else if spec, err := escalationBoundarySpec(e, escalations); err == nil {
			d.escalationCatch[e.ID] = spec.EscalationCode
			d.escalationBoundaries[spec.AttachedTo] = append(d.escalationBoundaries[spec.AttachedTo], e.ID)
		} else if spec, err := conditionalBoundarySpec(e); err == nil {
			d.conditionalCatch[e.ID] = spec.Condition
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
			spec, err := eventSubProcessStartSpec(sp.ID, sp.StartEvents[0], messages, signals, errors, escalations)
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
				case CatchKindCompensate:
					if err := d.registerCompensationEventSubProcess(scopeID, spec, d.Process.ID); err != nil {
						return err
					}
				case CatchKindConditional:
					d.conditionalCatch[spec.StartEventID] = spec.Condition
				}
			}
		}
		if err := d.indexScope(&sp.FlowElements, sp.ID, messages, signals, errors, escalations, associations); err != nil {
			return err
		}
	}
	for i := range fe.Transactions {
		tx := &fe.Transactions[i]
		if err := validateTransaction(*tx, d.IsTransaction(scopeID)); err != nil {
			return err
		}
		// Nested transaction: parent scope is itself a transaction.
		if d.IsTransaction(scopeID) {
			return fmt.Errorf("UNSUPPORTED_ELEMENT: nested transaction %q is not supported", tx.ID)
		}
		reg(tx.ID, eventv1.Element_TYPE_TRANSACTION, tx.Outgoing, tx.Incoming)
		if err := d.indexScope(&tx.FlowElements, tx.ID, messages, signals, errors, escalations, associations); err != nil {
			return err
		}
	}
	for i := range fe.AdHocSubProcesses {
		ah := &fe.AdHocSubProcesses[i]
		if err := validateAdHocSubProcess(*ah); err != nil {
			return err
		}
		reg(ah.ID, eventv1.Element_TYPE_AD_HOC_SUB_PROCESS, ah.Outgoing, ah.Incoming)
		d.adHocSubProcesses[ah.ID] = adHocSubProcessSpec(*ah)
		if err := d.indexScope(&ah.FlowElements, ah.ID, messages, signals, errors, escalations, associations); err != nil {
			return err
		}
	}
	return nil
}

func (d *Deployment) indexCallActivities(fe *element.FlowElements, catalog map[string]*element.Process, messages []element.Message, signals []element.Signal, errors []element.Error, escalations []element.Escalation) error {
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
			if err := validateProcess(called, messages, signals, errors, escalations); err != nil {
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
				if err := d.indexScope(&called.FlowElements, called.ID, messages, signals, errors, escalations, collectAssociations(called)); err != nil {
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
		for i := range fe.Transactions {
			if err := walk(&fe.Transactions[i].FlowElements); err != nil {
				return err
			}
		}
		return nil
	}
	return walk(fe)
}

// validateProcess checks process-level deploy rules: supported boundaries and
// hosts, duplicate same-kind attachments, and start vs instantiate entry exclusivity.
func validateProcess(proc *element.Process, messages []element.Message, signals []element.Signal, errors []element.Error, escalations []element.Escalation) error {
	if err := validateSubProcesses(&proc.FlowElements, errors, escalations); err != nil {
		return err
	}
	if err := validateEventBasedGateways(&proc.FlowElements); err != nil {
		return err
	}
	if err := validateCatchAndThrow(&proc.FlowElements); err != nil {
		return err
	}
	if err := validateSequenceFlowConditions(&proc.FlowElements); err != nil {
		return err
	}
	type seenKey struct {
		activity string
		kind     string // timer|message|signal|compensate|error:CODE
		disc     string // message/signal name; timer uses boundary id; empty for compensate
	}
	seenAttach := make(map[seenKey]string, len(proc.BoundaryEvents))
	assocs := collectAssociations(proc)
	for _, e := range proc.BoundaryEvents {
		attached := ""
		kind := ""
		disc := ""
		switch {
		case len(e.TimerEventDefinitions) > 0:
			spec, err := timerBoundarySpec(e)
			if err != nil {
				return err
			}
			attached = spec.AttachedTo
			kind = "timer"
			disc = e.ID
		case len(e.MessageEventDefinitions) > 0:
			spec, err := messageBoundarySpec(e, messages)
			if err != nil {
				return err
			}
			attached = spec.AttachedTo
			kind = "message"
			disc = spec.Name
		case len(e.SignalEventDefinitions) > 0:
			spec, err := signalBoundarySpec(e, signals)
			if err != nil {
				return err
			}
			attached = spec.AttachedTo
			kind = "signal"
			disc = spec.Name
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
			spec, err := errorBoundarySpec(e, errors)
			if err != nil {
				return err
			}
			attached = spec.AttachedTo
			kind = "error:" + spec.ErrorCode
		case len(e.EscalationEventDefinitions) > 0:
			spec, err := escalationBoundarySpec(e, escalations)
			if err != nil {
				return err
			}
			attached = spec.AttachedTo
			kind = "escalation:" + spec.EscalationCode
		case len(e.ConditionalEventDefinitions) > 0:
			spec, err := conditionalBoundarySpec(e)
			if err != nil {
				return err
			}
			attached = spec.AttachedTo
			kind = "conditional"
			disc = spec.Condition
		case len(e.CancelEventDefinitions) > 0:
			attachedTo, err := cancelBoundarySpec(e)
			if err != nil {
				return err
			}
			attached = attachedTo
			kind = "cancel"
		default:
			return fmt.Errorf("UNSUPPORTED_ELEMENT: boundaryEvent %q must be a timer, message, signal, error, escalation, compensation, conditional, or cancel boundary", e.ID)
		}
		key := seenKey{attached, kind, disc}
		if prev, ok := seenAttach[key]; ok {
			return fmt.Errorf("UNSUPPORTED_ELEMENT: activity %q already has %s boundary %q", attached, kind, prev)
		}
		if err := validateBoundaryHost(proc, attached); err != nil {
			return err
		}
		seenAttach[key] = e.ID
	}
	if len(proc.StartEvents) == 0 {
		if _, err := instantiateEntryIDs(proc); err != nil {
			return fmt.Errorf("no startEvent in process")
		}
		return nil
	}
	for _, g := range proc.EventBasedGatewaies {
		if g.Instantiate {
			return fmt.Errorf("UNSUPPORTED_ELEMENT: process cannot combine startEvent with instantiate eventBasedGateway %q", g.ID)
		}
	}
	for _, e := range proc.ReceiveTasks {
		if e.Instantiate {
			return fmt.Errorf("UNSUPPORTED_ELEMENT: process cannot combine startEvent with instantiate receiveTask %q", e.ID)
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
		if _, err := linkCatchSpec(e); err == nil {
			continue
		}
		if _, err := conditionalCatchSpec(e); err == nil {
			continue
		}
		if _, err := escalationCatchSpec(e, nil); err == nil {
			continue
		}
		return fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateCatchEvent %q must be timer, message, signal, link, conditional, or escalation catch", e.ID)
	}
	for _, e := range fe.IntermediateThrowEvents {
		if _, err := throwEventSpec(e, nil, nil, nil); err != nil {
			return err
		}
	}
	for _, child := range childScopes(fe) {
		if err := validateCatchAndThrow(child); err != nil {
			return err
		}
	}
	return nil
}

func (d *Deployment) ProcessID() string { return d.Process.ID }

func StartEventID(proc *element.Process) (string, error) {
	if len(proc.StartEvents) == 0 {
		ids, err := instantiateEntryIDs(proc)
		if err != nil {
			return "", err
		}
		return ids[0], nil
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

// instantiateEntryIDs returns process-level instantiate entries: event-based
// gateways and/or receive tasks with instantiate=true and no incoming flows.
func instantiateEntryIDs(proc *element.Process) ([]string, error) {
	var ids []string
	hasIncoming := func(elementID string, declaredIncoming []string) bool {
		if len(declaredIncoming) > 0 {
			return true
		}
		for _, f := range proc.SequenceFlows {
			if f.TargetRef == elementID {
				return true
			}
		}
		return false
	}
	outgoingOf := func(elementID string, declared []string) []string {
		outs := declared
		if len(outs) == 0 {
			for _, f := range proc.SequenceFlows {
				if f.SourceRef == elementID {
					outs = append(outs, f.ID)
				}
			}
		}
		return outs
	}
	for _, g := range proc.EventBasedGatewaies {
		if !g.Instantiate {
			continue
		}
		if hasIncoming(g.ID, g.Incoming) {
			return nil, fmt.Errorf("UNSUPPORTED_ELEMENT: eventBasedGateway %q instantiate must have no incoming sequence flow", g.ID)
		}
		outs := outgoingOf(g.ID, g.Outgoing)
		if len(outs) < 2 {
			return nil, fmt.Errorf("UNSUPPORTED_ELEMENT: eventBasedGateway %q must have at least two outgoing flows", g.ID)
		}
		ids = append(ids, g.ID)
	}
	for _, e := range proc.ReceiveTasks {
		if !e.Instantiate {
			continue
		}
		if hasIncoming(e.ID, e.Incoming) {
			return nil, fmt.Errorf("UNSUPPORTED_ELEMENT: receiveTask %q instantiate must have no incoming sequence flow", e.ID)
		}
		ids = append(ids, e.ID)
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("no startEvent in process")
	}
	return ids, nil
}

// InstantiateEntryIDs returns all CreateInstance entry element ids for instantiate-only processes.
func (d *Deployment) InstantiateEntryIDs() ([]string, error) {
	if d == nil {
		return nil, fmt.Errorf("no deployment")
	}
	return instantiateEntryIDs(&d.Process)
}

// InstantiateAlternativePeers returns other instantiate-entry wait element ids to cancel
// when winningElementID completes as part of an exclusive instantiate race.
// Parallel event-based gateway targets only cancel peers from *other* instantiate entries.
func (d *Deployment) InstantiateAlternativePeers(winningElementID string) []string {
	if d == nil {
		return nil
	}
	entries, err := instantiateEntryIDs(&d.Process)
	if err != nil || len(entries) == 0 {
		return nil
	}
	winningEntry := ""
	for _, entryID := range entries {
		if entryID == winningElementID {
			winningEntry = entryID
			break
		}
		for _, outID := range d.Outgoing(entryID) {
			sf, err := d.SequenceFlow(outID)
			if err != nil {
				continue
			}
			if sf.TargetRef == winningElementID {
				winningEntry = entryID
				break
			}
		}
		if winningEntry != "" {
			break
		}
	}
	if winningEntry == "" {
		return nil
	}
	// Exclusive race across alternative entries: cancel waits on other entries.
	// Same-gateway siblings are handled by EventBasedSiblings.
	var peers []string
	for _, entryID := range entries {
		if entryID == winningEntry {
			continue
		}
		typ, err := d.TypeOf(entryID)
		if err != nil {
			continue
		}
		if typ == eventv1.Element_TYPE_RECEIVE_TASK {
			peers = append(peers, entryID)
			continue
		}
		for _, outID := range d.Outgoing(entryID) {
			sf, err := d.SequenceFlow(outID)
			if err != nil {
				continue
			}
			if sf.TargetRef != "" && sf.TargetRef != winningElementID {
				peers = append(peers, sf.TargetRef)
			}
		}
	}
	return peers
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
		if isOpaqueSubProcessBody(&sp.FlowElements) {
			return "", fmt.Errorf("NOT_FOUND: opaque subProcess %q has no startEvent", subProcessID)
		}
		return "", fmt.Errorf("no startEvent in subProcess %q", subProcessID)
	}
	return sp.StartEvents[0].ID, nil
}

// IsOpaqueSubProcess reports whether id is an empty/collapsed SubProcess with no
// inner flow nodes. Such elements run as wait → Complete (opaque Activity).
func (d *Deployment) IsOpaqueSubProcess(id string) bool {
	sp := d.findSubProcess(&d.Process.FlowElements, id)
	if sp == nil {
		for _, called := range d.calledProcesses {
			fe := called.FlowElements
			if sp = d.findSubProcess(&fe, id); sp != nil {
				break
			}
		}
	}
	if sp == nil || sp.TriggeredByEvent {
		return false
	}
	return len(sp.StartEvents) == 0 && isOpaqueSubProcessBody(&sp.FlowElements)
}

// isOpaqueSubProcessBody is true when the SubProcess has no executable inner
// flow nodes (modeling stub / collapsed empty body). Data objects alone do not
// make the body executable.
func isOpaqueSubProcessBody(fe *element.FlowElements) bool {
	if fe == nil {
		return true
	}
	return len(fe.StartEvents) == 0 &&
		len(fe.EndEvents) == 0 &&
		len(fe.Tasks) == 0 &&
		len(fe.ManualTasks) == 0 &&
		len(fe.UserTasks) == 0 &&
		len(fe.ServiceTasks) == 0 &&
		len(fe.SendTasks) == 0 &&
		len(fe.ReceiveTasks) == 0 &&
		len(fe.BusinessRuleTasks) == 0 &&
		len(fe.ScriptTasks) == 0 &&
		len(fe.ParallelGatewaies) == 0 &&
		len(fe.ExclusiveGatewaies) == 0 &&
		len(fe.InclusiveGatewaies) == 0 &&
		len(fe.ComplexGatewaies) == 0 &&
		len(fe.EventBasedGatewaies) == 0 &&
		len(fe.SubProcesses) == 0 &&
		len(fe.Transactions) == 0 &&
		len(fe.AdHocSubProcesses) == 0 &&
		len(fe.BoundaryEvents) == 0 &&
		len(fe.CallActivities) == 0 &&
		len(fe.IntermediateThrowEvents) == 0 &&
		len(fe.IntermediateCatchEvents) == 0 &&
		len(fe.ImplicitThrowEvents) == 0
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
	for _, child := range childScopes(fe) {
		if f, err := findSequenceFlowIn(child, id); err == nil {
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
	for _, child := range childScopes(fe) {
		if st, err := findServiceTaskIn(child, id); err == nil {
			return st, nil
		}
	}
	return element.ServiceTask{}, fmt.Errorf("NOT_FOUND: service task %q", id)
}

func (d *Deployment) BusinessRuleTask(id string) (element.BusinessRuleTask, error) {
	return findBusinessRuleTaskIn(&d.Process.FlowElements, id)
}

func (d *Deployment) ScriptTask(id string) (element.ScriptTask, error) {
	return findScriptTaskIn(&d.Process.FlowElements, id)
}

func findScriptTaskIn(fe *element.FlowElements, id string) (element.ScriptTask, error) {
	for _, st := range fe.ScriptTasks {
		if st.ID == id {
			return st, nil
		}
	}
	for _, child := range childScopes(fe) {
		if st, err := findScriptTaskIn(child, id); err == nil {
			return st, nil
		}
	}
	return element.ScriptTask{}, fmt.Errorf("NOT_FOUND: script task %q", id)
}

func findBusinessRuleTaskIn(fe *element.FlowElements, id string) (element.BusinessRuleTask, error) {
	for _, st := range fe.BusinessRuleTasks {
		if st.ID == id {
			return st, nil
		}
	}
	for _, child := range childScopes(fe) {
		if st, err := findBusinessRuleTaskIn(child, id); err == nil {
			return st, nil
		}
	}
	return element.BusinessRuleTask{}, fmt.Errorf("NOT_FOUND: business rule task %q", id)
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

// validateSequenceFlowConditions rejects FEEL constructs Sparrow does not evaluate
// (quantifiers / satisfies). Narrow interop sugar lives in processing/expr.
func validateSequenceFlowConditions(fe *element.FlowElements) error {
	if fe == nil {
		return nil
	}
	for _, f := range fe.SequenceFlows {
		text := strings.TrimSpace(ConditionText(f))
		if text == "" {
			continue
		}
		lower := strings.ToLower(text)
		if strings.Contains(lower, " satisfies ") ||
			strings.Contains(lower, " some ") ||
			strings.HasPrefix(lower, "some ") ||
			strings.Contains(lower, " every ") ||
			strings.HasPrefix(lower, "every ") {
			return fmt.Errorf("INVALID_CONDITION: flow %s: FEEL quantifier expressions are not supported", f.ID)
		}
	}
	for i := range fe.SubProcesses {
		if err := validateSequenceFlowConditions(&fe.SubProcesses[i].FlowElements); err != nil {
			return err
		}
	}
	for i := range fe.Transactions {
		if err := validateSequenceFlowConditions(&fe.Transactions[i].FlowElements); err != nil {
			return err
		}
	}
	for i := range fe.AdHocSubProcesses {
		if err := validateSequenceFlowConditions(&fe.AdHocSubProcesses[i].FlowElements); err != nil {
			return err
		}
	}
	return nil
}

// DefaultOutgoing returns the BPMN default sequence flow id for an exclusive
// gateway or activity, or empty if none.
func (d *Deployment) DefaultOutgoing(elementID string) string {
	if g := findExclusiveGatewayIn(&d.Process.FlowElements, elementID); g != nil {
		return g.Default
	}
	return findActivityDefaultIn(&d.Process.FlowElements, elementID)
}

func findActivityDefaultIn(fe *element.FlowElements, id string) string {
	for i := range fe.UserTasks {
		if fe.UserTasks[i].ID == id {
			return fe.UserTasks[i].Default
		}
	}
	for i := range fe.ServiceTasks {
		if fe.ServiceTasks[i].ID == id {
			return fe.ServiceTasks[i].Default
		}
	}
	for i := range fe.Tasks {
		if fe.Tasks[i].ID == id {
			return fe.Tasks[i].Default
		}
	}
	for i := range fe.ManualTasks {
		if fe.ManualTasks[i].ID == id {
			return fe.ManualTasks[i].Default
		}
	}
	for i := range fe.ReceiveTasks {
		if fe.ReceiveTasks[i].ID == id {
			return fe.ReceiveTasks[i].Default
		}
	}
	for i := range fe.SendTasks {
		if fe.SendTasks[i].ID == id {
			return fe.SendTasks[i].Default
		}
	}
	for i := range fe.BusinessRuleTasks {
		if fe.BusinessRuleTasks[i].ID == id {
			return fe.BusinessRuleTasks[i].Default
		}
	}
	for i := range fe.ScriptTasks {
		if fe.ScriptTasks[i].ID == id {
			return fe.ScriptTasks[i].Default
		}
	}
	for i := range fe.CallActivities {
		if fe.CallActivities[i].ID == id {
			return fe.CallActivities[i].Default
		}
	}
	for i := range fe.SubProcesses {
		if fe.SubProcesses[i].ID == id {
			return fe.SubProcesses[i].Default
		}
		if def := findActivityDefaultIn(&fe.SubProcesses[i].FlowElements, id); def != "" {
			return def
		}
	}
	for i := range fe.Transactions {
		if fe.Transactions[i].ID == id {
			return fe.Transactions[i].Default
		}
		if def := findActivityDefaultIn(&fe.Transactions[i].FlowElements, id); def != "" {
			return def
		}
	}
	return ""
}

// ChooseConditionalOutgoing picks an outgoing flow: first matching non-default
// condition, else default. If there are no conditions and no default, returns the
// first outgoing (callers that need take-all should use ChooseOutgoingFlows).
func (d *Deployment) ChooseConditionalOutgoing(elementID string, vars map[string]string) (string, error) {
	outs := d.Outgoing(elementID)
	if len(outs) == 0 {
		return "", fmt.Errorf("NO_OUTGOING_FLOW: element %q", elementID)
	}
	def := d.DefaultOutgoing(elementID)
	hasCondition := false
	for _, flowID := range outs {
		flow, err := d.SequenceFlow(flowID)
		if err != nil {
			return "", err
		}
		if ConditionText(flow) != "" {
			hasCondition = true
			break
		}
	}
	if !hasCondition {
		if def != "" {
			if _, err := d.SequenceFlow(def); err == nil {
				return def, nil
			}
		}
		return outs[0], nil
	}
	for _, flowID := range outs {
		if def != "" && flowID == def {
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
	if def != "" {
		if _, err := d.SequenceFlow(def); err == nil {
			return def, nil
		}
	}
	return "", fmt.Errorf("NO_OUTGOING_FLOW: element %q", elementID)
}

// ChooseOutgoingFlows selects outgoing sequence flows when leaving an element.
// Multiple unconditional outgoings with no default are an implicit parallel split
// (all flows). Otherwise selection is exclusive (one flow via ChooseConditionalOutgoing).
func (d *Deployment) ChooseOutgoingFlows(elementID string, vars map[string]string) ([]string, error) {
	outs := d.Outgoing(elementID)
	if len(outs) == 0 {
		return nil, fmt.Errorf("NO_OUTGOING_FLOW: element %q", elementID)
	}
	def := d.DefaultOutgoing(elementID)
	hasCondition := false
	for _, flowID := range outs {
		flow, err := d.SequenceFlow(flowID)
		if err != nil {
			return nil, err
		}
		if ConditionText(flow) != "" {
			hasCondition = true
			break
		}
	}
	if !hasCondition && def == "" {
		return append([]string{}, outs...), nil
	}
	if !hasCondition && def != "" {
		return []string{def}, nil
	}
	one, err := d.ChooseConditionalOutgoing(elementID, vars)
	if err != nil {
		return nil, err
	}
	return []string{one}, nil
}

// ChooseExclusiveOutgoing picks the first matching non-default condition, else default.
func (d *Deployment) ChooseExclusiveOutgoing(gatewayID string, vars map[string]string) (string, error) {
	if findExclusiveGatewayIn(&d.Process.FlowElements, gatewayID) == nil {
		return "", fmt.Errorf("NOT_FOUND: exclusive gateway %q", gatewayID)
	}
	return d.ChooseConditionalOutgoing(gatewayID, vars)
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
	for _, child := range childScopes(fe) {
		if g := findInclusiveGatewayIn(child, id); g != nil {
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
	for _, child := range childScopes(fe) {
		if g := findExclusiveGatewayIn(child, id); g != nil {
			return g
		}
	}
	return nil
}

func resolveTaskMessageName(messageRef, name, id string, messages []element.Message) string {
	if n := resolveMessageName(messageRef, messages); n != "" {
		return n
	}
	if n := strings.TrimSpace(name); n != "" {
		return n
	}
	return id
}

func validateBoundaryHost(proc *element.Process, activityID string) error {
	if boundaryHostExists(&proc.FlowElements, activityID) {
		return nil
	}
	return fmt.Errorf("UNSUPPORTED_ELEMENT: boundary must attach to a task, userTask, serviceTask, manualTask, receiveTask, sendTask, businessRuleTask, scriptTask, subProcess, transaction, or callActivity (%q)", activityID)
}

func boundaryHostExists(fe *element.FlowElements, activityID string) bool {
	for _, e := range fe.UserTasks {
		if e.ID == activityID {
			return true
		}
	}
	for _, e := range fe.ServiceTasks {
		if e.ID == activityID {
			return true
		}
	}
	for _, e := range fe.Tasks {
		if e.ID == activityID {
			return true
		}
	}
	for _, e := range fe.ManualTasks {
		if e.ID == activityID {
			return true
		}
	}
	for _, e := range fe.ReceiveTasks {
		if e.ID == activityID {
			return true
		}
	}
	for _, e := range fe.SendTasks {
		if e.ID == activityID {
			return true
		}
	}
	for _, e := range fe.BusinessRuleTasks {
		if e.ID == activityID {
			return true
		}
	}
	for _, e := range fe.ScriptTasks {
		if e.ID == activityID {
			return true
		}
	}
	for _, e := range fe.SubProcesses {
		if e.ID == activityID {
			return true
		}
		if boundaryHostExists(&e.FlowElements, activityID) {
			return true
		}
	}
	for _, e := range fe.Transactions {
		if e.ID == activityID {
			return true
		}
		if boundaryHostExists(&e.FlowElements, activityID) {
			return true
		}
	}
	for _, e := range fe.CallActivities {
		if e.ID == activityID {
			return true
		}
	}
	return false
}

func validateSubProcesses(fe *element.FlowElements, errors []element.Error, escalations []element.Escalation) error {
	return validateSubProcessesAt(fe, false, false, errors, escalations)
}

func validateSubProcessesAt(fe *element.FlowElements, insideEmbedded, insideEventSubProcess bool, errors []element.Error, escalations []element.Escalation) error {
	return validateScopesAt(fe, insideEmbedded, insideEventSubProcess, false, errors, escalations)
}

func validateScopesAt(fe *element.FlowElements, insideEmbedded, insideEventSubProcess, insideTransaction bool, errors []element.Error, escalations []element.Escalation) error {
	for i := range fe.SubProcesses {
		sp := &fe.SubProcesses[i]
		if sp.TriggeredByEvent {
			if err := validateEventSubProcess(sp, nil, nil, errors, escalations); err != nil {
				return err
			}
			for _, f := range fe.SequenceFlows {
				if f.SourceRef == sp.ID || f.TargetRef == sp.ID {
					return fmt.Errorf("UNSUPPORTED_ELEMENT: event subProcess %q must not have sequence flow connections", sp.ID)
				}
			}
			if err := validateScopesAt(&sp.FlowElements, false, true, insideTransaction, errors, escalations); err != nil {
				return err
			}
			continue
		}
		if len(sp.StartEvents) == 0 {
			if !isOpaqueSubProcessBody(&sp.FlowElements) {
				return fmt.Errorf("UNSUPPORTED_ELEMENT: subProcess %q must have a startEvent", sp.ID)
			}
			// Collapsed/empty SubProcess (no inner flow nodes): treated as opaque
			// Activity (wait → Complete), not an enterable scope.
			continue
		}
		if err := validateScopesAt(&sp.FlowElements, true, insideEventSubProcess, insideTransaction, errors, escalations); err != nil {
			return err
		}
	}
	for i := range fe.Transactions {
		tx := &fe.Transactions[i]
		if err := validateTransaction(*tx, insideTransaction); err != nil {
			return err
		}
		if err := validateScopesAt(&tx.FlowElements, true, insideEventSubProcess, true, errors, escalations); err != nil {
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
	receiveIDs := make(map[string]bool, len(fe.ReceiveTasks))
	for _, e := range fe.ReceiveTasks {
		receiveIDs[e.ID] = true
	}
	for _, g := range fe.EventBasedGatewaies {
		if g.Instantiate {
			if insideSubProcess {
				return fmt.Errorf("UNSUPPORTED_ELEMENT: eventBasedGateway %q instantiate inside subProcess is not supported", g.ID)
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
		// Exclusive (default) and Parallel intermediate / instantiate event-based gateways are supported.
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
			if !catchIDs[target] && !receiveIDs[target] {
				return fmt.Errorf("UNSUPPORTED_ELEMENT: eventBasedGateway %q must target intermediateCatchEvent or receiveTask (%q)", g.ID, target)
			}
		}
	}
	for _, child := range childScopes(fe) {
		if err := validateEventBasedGatewaysAt(child, true); err != nil {
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
	for _, child := range childScopes(fe) {
		if isParallelEventBasedGatewayIn(child, id) {
			return true
		}
	}
	return false
}
