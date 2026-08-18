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
	// MessageName is the BPMN message name copied from a message catch ACTIVATED.
	// Empty when the token is not waiting on a message.
	MessageName string
	// BoundaryID is the interrupting timer boundary armed on a waiting activity.
	// Empty when the token is not waiting on an attached timer.
	BoundaryID string
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
}

func NewInstance(id, deploymentID string, version int32) *Instance {
	return &Instance{
		ID:            id,
		DeploymentID:  deploymentID,
		Version:       version,
		Status:        StatusActive,
		Variables:     make(map[string]string),
		Tokens:        make(map[string]*Token),
		ElementIntent: make(map[string]eventv1.Element_Intent),
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
		ID:            inst.ID,
		DeploymentID:  inst.DeploymentID,
		Version:       inst.Version,
		Status:        inst.Status,
		Variables:     make(map[string]string, len(inst.Variables)),
		Tokens:        make(map[string]*Token, len(inst.Tokens)),
		ElementIntent: make(map[string]eventv1.Element_Intent, len(inst.ElementIntent)),
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
		tok.MessageName = ""
		tok.BoundaryID = ""
		if p := el.GetActivityPayload(); p != nil {
			tok.JobType = p.GetJobType()
			tok.DueUnixMs = p.GetDueUnixMs()
			tok.BoundaryID = p.GetBoundaryId()
		}
		if p := el.GetEventPayload(); p != nil {
			if p.GetDueUnixMs() != 0 {
				tok.DueUnixMs = p.GetDueUnixMs()
			}
			tok.MessageName = p.GetMessageName()
		}
	case eventv1.Element_INTENT_COMPLETED, eventv1.Element_INTENT_TERMINATED:
		tok.Status = TokenActive
		tok.JobType = ""
		tok.DueUnixMs = 0
		tok.MessageName = ""
		tok.BoundaryID = ""
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
		tok.MessageName = ""
		tok.BoundaryID = ""
		if sp := el.GetSequenceFlowPayload(); sp != nil && sp.GetTargetId() != "" {
			tok.ElementID = sp.GetTargetId()
		}
	}
}

func waitingActivation(t eventv1.Element_Type) bool {
	switch t {
	case eventv1.Element_TYPE_USER_TASK, eventv1.Element_TYPE_SERVICE_TASK, eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT:
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
	case eventv1.Element_INTENT_TERMINATED:
		inst.Status = StatusTerminated
		inst.Tokens = make(map[string]*Token)
	case eventv1.Element_INTENT_ACTIVATED:
		inst.Status = StatusActive
	}
}

// VariablesFromMap encodes Go values as JSON text variables.
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
