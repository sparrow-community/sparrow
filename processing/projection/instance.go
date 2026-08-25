package projection

import (
	"encoding/json"

	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

type InstanceStatus string

const (
	StatusActive     InstanceStatus = "active"
	StatusCompleted  InstanceStatus = "completed"
	StatusTerminated InstanceStatus = "terminated"
)

type TokenStatus string

const (
	TokenActive  TokenStatus = "active"
	TokenWaiting TokenStatus = "waiting"
)

type Token struct {
	ID        string
	ElementID string
	Status    TokenStatus
	// JobType is the worker subscription key copied from SERVICE_TASK ACTIVATED.
	// Empty for user tasks and non-waiting tokens.
	JobType string
	// DueUnixMs is the timer catch due time copied from INTERMEDIATE_CATCH_EVENT ACTIVATED.
	// Zero when the token is not waiting on a timer.
	DueUnixMs int64
	// TimerText stores the original timer expression text for waiting timer catches
	// and timer boundaries so cycle-based waits can re-arm across recovery.
	TimerText string
	// MessageName is the BPMN message name copied from a message catch ACTIVATED
	// or a message boundary on a waiting activity.
	// Empty when the token is not waiting on a message.
	MessageName string
	// SignalName is the BPMN signal name copied from a signal catch ACTIVATED
	// or a signal boundary on a waiting activity.
	// Empty when the token is not waiting on a signal.
	SignalName string
	// BoundaryID is the first armed waiting boundary on a waiting activity
	// (timer, or message/signal when that is the only waiting boundary).
	BoundaryID string
	// MessageBoundaryID is the message boundary armed on a waiting activity.
	// Empty when no message boundary is armed. When both timer and message
	// boundaries exist, BoundaryID holds the timer and MessageBoundaryID holds the message.
	MessageBoundaryID string
	// SignalBoundaryID is the signal boundary when another waiting boundary already
	// occupies BoundaryID. Empty when only a signal boundary is armed (then BoundaryID holds it).
	SignalBoundaryID string
}

// ScopeBoundary tracks a boundary armed on a SubProcess scope.
type ScopeBoundary struct {
	BoundaryID  string
	ScopeID     string // the SubProcess element ID
	TokenID     string // the token that entered the SubProcess
	DueUnixMs   int64  // timer due (zero if message/signal boundary)
	TimerText   string // original timer expression
	MessageName string // message name (empty if timer/signal boundary)
	SignalName  string // signal name (empty if timer/message boundary)
}

// EventSubProcessArm is a process-scoped subscription for a triggeredByEvent subProcess.
// It is not an EventLog record; Engine re-syncs from the deployment while the instance is active.
type EventSubProcessArm struct {
	SubProcessID  string
	StartEventID  string
	ParentScopeID string
	Interrupting  bool
	MessageName   string
	SignalName    string
	DueUnixMs     int64
	TimerText     string
}

type Instance struct {
	ID           string
	DeploymentID string
	Version      int32
	Status       InstanceStatus
	Variables    map[string]string // name -> json_value
	Tokens       map[string]*Token
	// elementID -> last intent seen (for validation helpers)
	ElementIntent map[string]eventv1.Element_Intent
	// ScopeBoundaries tracks boundaries armed on SubProcess scopes.
	// Key is the boundary element ID.
	ScopeBoundaries map[string]*ScopeBoundary
	// EventSubProcesses tracks armed event sub-process starts. Key is subProcess id.
	EventSubProcesses map[string]*EventSubProcessArm
	// CompensationSubs tracks compensation subscriptions after host activities complete.
	// Key is compensation boundary id.
	CompensationSubs map[string]*CompensationSub
	// PendingCompensation is set while a compensate throw waits for handlers.
	PendingCompensation *PendingCompensation
}

// CompensationSub is created when a host activity COMPLETED and a compensation
// boundary ACTIVATED (ledger-backed via EventPayload.compensation_handler_id).
type CompensationSub struct {
	BoundaryID   string
	ActivityID   string
	HandlerID    string
	Seq          int64 // subscription order (later compensated first)
	HostTokenID  string
}

// PendingCompensation tracks an in-flight compensate throw waiting for handlers.
type PendingCompensation struct {
	ThrowTokenID   string
	ThrowElementID string
	Queue          []string // remaining handler element ids (reverse execution order)
	ActiveTokenID  string
	ActiveHandler  string
	Consumed       []string // boundary ids consumed by this throw
}

func NewInstance(id, deploymentID string, version int32) *Instance {
	return &Instance{
		ID:                id,
		DeploymentID:      deploymentID,
		Version:           version,
		Status:            StatusActive,
		Variables:         make(map[string]string),
		Tokens:            make(map[string]*Token),
		ElementIntent:     make(map[string]eventv1.Element_Intent),
		ScopeBoundaries:   make(map[string]*ScopeBoundary),
		EventSubProcesses: make(map[string]*EventSubProcessArm),
		CompensationSubs:  make(map[string]*CompensationSub),
	}
}

func (inst *Instance) ApplyEvent(e *eventv1.Event) {
	if e == nil || e.GetRecordType() != eventv1.Event_RECORD_TYPE_EVENT {
		return
	}
	el := e.GetElement()
	if el == nil {
		return
	}

	inst.ElementIntent[el.GetId()] = el.GetIntent()
	mergeVariables(inst, el)
	applyToken(inst, el)
	applyProcessLifecycle(inst, el)
}

// Clone returns a deep copy for read-only snapshots (GetInstance).
func (inst *Instance) Clone() *Instance {
	if inst == nil {
		return nil
	}
	out := &Instance{
		ID:                inst.ID,
		DeploymentID:      inst.DeploymentID,
		Version:           inst.Version,
		Status:            inst.Status,
		Variables:         make(map[string]string, len(inst.Variables)),
		Tokens:            make(map[string]*Token, len(inst.Tokens)),
		ElementIntent:     make(map[string]eventv1.Element_Intent, len(inst.ElementIntent)),
		ScopeBoundaries:   make(map[string]*ScopeBoundary, len(inst.ScopeBoundaries)),
		EventSubProcesses: make(map[string]*EventSubProcessArm, len(inst.EventSubProcesses)),
		CompensationSubs:  make(map[string]*CompensationSub, len(inst.CompensationSubs)),
	}
	for k, v := range inst.Variables {
		out.Variables[k] = v
	}
	for k, tok := range inst.Tokens {
		if tok == nil {
			continue
		}
		cp := *tok
		out.Tokens[k] = &cp
	}
	for k, v := range inst.ElementIntent {
		out.ElementIntent[k] = v
	}
	for k, sb := range inst.ScopeBoundaries {
		cp := *sb
		out.ScopeBoundaries[k] = &cp
	}
	for k, arm := range inst.EventSubProcesses {
		cp := *arm
		out.EventSubProcesses[k] = &cp
	}
	for k, sub := range inst.CompensationSubs {
		cp := *sub
		out.CompensationSubs[k] = &cp
	}
	if inst.PendingCompensation != nil {
		pc := *inst.PendingCompensation
		pc.Queue = append([]string{}, inst.PendingCompensation.Queue...)
		pc.Consumed = append([]string{}, inst.PendingCompensation.Consumed...)
		out.PendingCompensation = &pc
	}
	return out
}

func mergeVariables(inst *Instance, el *eventv1.Element) {
	switch p := el.GetPayload().(type) {
	case *eventv1.Element_ProcessPayload:
		for _, v := range p.ProcessPayload.GetVariables() {
			inst.Variables[v.GetName()] = v.GetJsonValue()
		}
	case *eventv1.Element_ActivityPayload:
		for _, v := range p.ActivityPayload.GetVariables() {
			inst.Variables[v.GetName()] = v.GetJsonValue()
		}
	case *eventv1.Element_EventPayload:
		for _, v := range p.EventPayload.GetVariables() {
			inst.Variables[v.GetName()] = v.GetJsonValue()
		}
	}
}

func applyToken(inst *Instance, el *eventv1.Element) {
	tokenID := el.GetTokenId()
	if tokenID == "" {
		return
	}
	tok, ok := inst.Tokens[tokenID]
	if !ok {
		tok = &Token{ID: tokenID}
		inst.Tokens[tokenID] = tok
	}
	hostElementID := tok.ElementID
	if boundaryDisarmOnWaitingHost(hostElementID, el, tok) {
		if el.GetIntent() == eventv1.Element_INTENT_TERMINATED {
			bid := el.GetId()
			if bid == tok.BoundaryID {
				tok.DueUnixMs = 0
				tok.TimerText = ""
				tok.BoundaryID = ""
			}
			if bid == tok.MessageBoundaryID || (tok.MessageBoundaryID == "" && bid != tok.BoundaryID) {
				tok.MessageName = ""
				tok.MessageBoundaryID = ""
			}
			if bid == tok.SignalBoundaryID || (tok.SignalBoundaryID == "" && tok.SignalName != "" && bid != tok.MessageBoundaryID) {
				tok.SignalName = ""
				tok.SignalBoundaryID = ""
			}
		}
		return
	}
	if boundaryRearmOnWaitingHost(hostElementID, el, tok) {
		if p := el.GetActivityPayload(); p != nil {
			tok.DueUnixMs = p.GetDueUnixMs()
			tok.TimerText = p.GetDuration()
			tok.BoundaryID = p.GetBoundaryId()
			if p.GetMessageName() != "" {
				tok.MessageName = p.GetMessageName()
			}
			if p.GetMessageBoundaryId() != "" {
				tok.MessageBoundaryID = p.GetMessageBoundaryId()
			}
			if p.GetSignalName() != "" {
				tok.SignalName = p.GetSignalName()
			}
			if p.GetSignalBoundaryId() != "" {
				tok.SignalBoundaryID = p.GetSignalBoundaryId()
			}
		}
		return
	}
	if el.GetType() == eventv1.Element_TYPE_BOUNDARY_EVENT {
		if el.GetIntent() == eventv1.Element_INTENT_ACTIVATED {
			if p := el.GetEventPayload(); p != nil && p.GetCompensationHandlerId() != "" {
				inst.CompensationSubs[el.GetId()] = &CompensationSub{
					BoundaryID:  el.GetId(),
					HandlerID:   p.GetCompensationHandlerId(),
					HostTokenID: tokenID,
					Seq:         int64(len(inst.CompensationSubs)) + 1,
				}
				return
			}
		}
		if el.GetIntent() == eventv1.Element_INTENT_TERMINATED || el.GetIntent() == eventv1.Element_INTENT_COMPLETED {
			if _, ok := inst.CompensationSubs[el.GetId()]; ok {
				delete(inst.CompensationSubs, el.GetId())
				return
			}
		}
	}
	tok.ElementID = el.GetId()

	// Tokens are updated only from EVENT records (not by the executor).
	switch el.GetIntent() {
	case eventv1.Element_INTENT_ACTIVATED:
		if waitingActivation(el.GetType()) {
			tok.Status = TokenWaiting
		} else {
			tok.Status = TokenActive
		}
		tok.JobType = ""
		tok.DueUnixMs = 0
		tok.TimerText = ""
		tok.MessageName = ""
		tok.SignalName = ""
		tok.BoundaryID = ""
		tok.MessageBoundaryID = ""
		tok.SignalBoundaryID = ""
		if p := el.GetActivityPayload(); p != nil {
			if el.GetType() == eventv1.Element_TYPE_SUB_PROCESS {
				inst.applyScopeBoundary(el.GetId(), tokenID, p)
			} else {
				tok.JobType = p.GetJobType()
				tok.DueUnixMs = p.GetDueUnixMs()
				tok.TimerText = p.GetDuration()
				tok.BoundaryID = p.GetBoundaryId()
				tok.MessageName = p.GetMessageName()
				tok.MessageBoundaryID = p.GetMessageBoundaryId()
				tok.SignalName = p.GetSignalName()
				tok.SignalBoundaryID = p.GetSignalBoundaryId()
			}
		}
		if p := el.GetEventPayload(); p != nil {
			if p.GetDueUnixMs() != 0 {
				tok.DueUnixMs = p.GetDueUnixMs()
			}
			if p.GetDuration() != "" {
				tok.TimerText = p.GetDuration()
			}
			if p.GetMessageName() != "" {
				tok.MessageName = p.GetMessageName()
			}
			if p.GetSignalName() != "" {
				tok.SignalName = p.GetSignalName()
			}
		}
	case eventv1.Element_INTENT_COMPLETED, eventv1.Element_INTENT_TERMINATED:
		if el.GetIntent() == eventv1.Element_INTENT_TERMINATED &&
			(el.GetType() == eventv1.Element_TYPE_PARALLEL_GATEWAY ||
				el.GetType() == eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT) {
			delete(inst.Tokens, tokenID)
			return
		}
		tok.Status = TokenActive
		tok.JobType = ""
		tok.DueUnixMs = 0
		tok.TimerText = ""
		tok.MessageName = ""
		tok.SignalName = ""
		tok.BoundaryID = ""
		tok.MessageBoundaryID = ""
		tok.SignalBoundaryID = ""
	case eventv1.Element_INTENT_FAILED:
		// Job failure does not complete the activity; worker may retry.
		tok.Status = TokenWaiting
		if p := el.GetActivityPayload(); p != nil && p.GetJobType() != "" {
			tok.JobType = p.GetJobType()
		}
	case eventv1.Element_INTENT_SEQUENCE_FLOW_TAKEN:
		tok.Status = TokenActive
		tok.JobType = ""
		tok.DueUnixMs = 0
		tok.TimerText = ""
		tok.MessageName = ""
		tok.SignalName = ""
		tok.BoundaryID = ""
		tok.MessageBoundaryID = ""
		tok.SignalBoundaryID = ""
		if sp := el.GetSequenceFlowPayload(); sp != nil && sp.GetTargetId() != "" {
			tok.ElementID = sp.GetTargetId()
		}
	}
}

func boundaryDisarmOnWaitingHost(hostElementID string, el *eventv1.Element, tok *Token) bool {
	if el.GetType() != eventv1.Element_TYPE_BOUNDARY_EVENT || tok.Status != TokenWaiting {
		return false
	}
	if hostElementID == "" || hostElementID == el.GetId() {
		return false
	}
	switch el.GetIntent() {
	case eventv1.Element_INTENT_TERMINATING, eventv1.Element_INTENT_TERMINATED:
		return true
	default:
		return false
	}
}

func boundaryRearmOnWaitingHost(hostElementID string, el *eventv1.Element, tok *Token) bool {
	if el.GetType() != eventv1.Element_TYPE_BOUNDARY_EVENT || tok.Status != TokenWaiting {
		return false
	}
	if hostElementID == "" || hostElementID == el.GetId() {
		return false
	}
	switch el.GetIntent() {
	case eventv1.Element_INTENT_ACTIVATING, eventv1.Element_INTENT_ACTIVATED:
		return true
	default:
		return false
	}
}

func waitingActivation(t eventv1.Element_Type) bool {
	switch t {
	case eventv1.Element_TYPE_USER_TASK, eventv1.Element_TYPE_SERVICE_TASK, eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT, eventv1.Element_TYPE_INTERMEDIATE_THROW_EVENT, eventv1.Element_TYPE_END_EVENT, eventv1.Element_TYPE_PARALLEL_GATEWAY, eventv1.Element_TYPE_INCLUSIVE_GATEWAY:
		return true
	default:
		return false
	}
}

func applyProcessLifecycle(inst *Instance, el *eventv1.Element) {
	if el.GetType() != eventv1.Element_TYPE_PROCESS {
		return
	}
	switch el.GetIntent() {
	case eventv1.Element_INTENT_COMPLETED:
		inst.Status = StatusCompleted
		inst.Tokens = make(map[string]*Token)
		inst.EventSubProcesses = make(map[string]*EventSubProcessArm)
		inst.ScopeBoundaries = make(map[string]*ScopeBoundary)
		inst.CompensationSubs = make(map[string]*CompensationSub)
		inst.PendingCompensation = nil
	case eventv1.Element_INTENT_TERMINATED:
		inst.Status = StatusTerminated
		inst.Tokens = make(map[string]*Token)
		inst.EventSubProcesses = make(map[string]*EventSubProcessArm)
		inst.ScopeBoundaries = make(map[string]*ScopeBoundary)
		inst.CompensationSubs = make(map[string]*CompensationSub)
		inst.PendingCompensation = nil
	case eventv1.Element_INTENT_ACTIVATED:
		inst.Status = StatusActive
	}
}

// VariablesFromMap encodes Go values as JSON text variables.
func (inst *Instance) applyScopeBoundary(scopeID, tokenID string, p *eventv1.ActivityPayload) {
	if p.GetBoundaryId() != "" {
		inst.ScopeBoundaries[p.GetBoundaryId()] = &ScopeBoundary{
			BoundaryID:  p.GetBoundaryId(),
			ScopeID:     scopeID,
			TokenID:     tokenID,
			DueUnixMs:   p.GetDueUnixMs(),
			TimerText:   p.GetDuration(),
			MessageName: p.GetMessageName(),
			SignalName:  p.GetSignalName(),
		}
	}
	if p.GetMessageBoundaryId() != "" {
		inst.ScopeBoundaries[p.GetMessageBoundaryId()] = &ScopeBoundary{
			BoundaryID:  p.GetMessageBoundaryId(),
			ScopeID:     scopeID,
			TokenID:     tokenID,
			MessageName: p.GetMessageName(),
		}
	}
	if p.GetSignalBoundaryId() != "" {
		inst.ScopeBoundaries[p.GetSignalBoundaryId()] = &ScopeBoundary{
			BoundaryID:  p.GetSignalBoundaryId(),
			ScopeID:     scopeID,
			TokenID:     tokenID,
			SignalName:  p.GetSignalName(),
		}
	}
}

func (inst *Instance) RemoveScopeBoundary(boundaryID string) {
	delete(inst.ScopeBoundaries, boundaryID)
}

func (inst *Instance) RemoveScopeBoundariesForScope(scopeID string) {
	for id, sb := range inst.ScopeBoundaries {
		if sb.ScopeID == scopeID {
			delete(inst.ScopeBoundaries, id)
		}
	}
}

func (inst *Instance) RemoveEventSubProcessesInScope(scopeID string) {
	for id, arm := range inst.EventSubProcesses {
		if arm.ParentScopeID == scopeID {
			delete(inst.EventSubProcesses, id)
		}
	}
}

func (inst *Instance) RemoveEventSubProcess(subProcessID string) {
	delete(inst.EventSubProcesses, subProcessID)
}

func VariablesFromMap(vars map[string]any) ([]*eventv1.Variable, error) {
	if len(vars) == 0 {
		return nil, nil
	}
	out := make([]*eventv1.Variable, 0, len(vars))
	for k, v := range vars {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		out = append(out, &eventv1.Variable{Name: k, JsonValue: string(b)})
	}
	return out, nil
}
