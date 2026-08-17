package processing

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sparrow-community/sparrow/processing/deploy"
	"github.com/sparrow-community/sparrow/processing/handlers"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// Engine is the single-node BPMN execution facade.
// Element semantics live in handlers; definitions come from bpmn via deploy.Deployment.
type Engine struct {
	log      eventlog.EventLog
	executor *Executor

	mu          sync.Mutex
	deployments map[string]*deploy.Deployment
	instances   map[string]*projection.Instance
	instMu      map[string]*sync.Mutex
}

func NewEngine(l eventlog.EventLog) *Engine {
	if l == nil {
		l = eventlog.NewMemory()
	}
	return &Engine{
		log:         l,
		executor:    NewExecutor(handlers.DefaultRegistry()),
		deployments: make(map[string]*deploy.Deployment),
		instances:   make(map[string]*projection.Instance),
		instMu:      make(map[string]*sync.Mutex),
	}
}

func (e *Engine) Deploy(_ context.Context, bpmnXML []byte) (string, error) {
	dep, err := deploy.Compile(bpmnXML)
	if err != nil {
		return "", err
	}
	id, err := NextID()
	if err != nil {
		return "", err
	}
	dep.ID = id

	e.mu.Lock()
	e.deployments[id] = dep
	e.mu.Unlock()
	return id, nil
}

func (e *Engine) CreateInstance(ctx context.Context, deploymentID string, vars map[string]any) (string, error) {
	e.mu.Lock()
	dep := e.deployments[deploymentID]
	e.mu.Unlock()
	if dep == nil {
		return "", fmt.Errorf("NOT_FOUND: deployment %q", deploymentID)
	}

	instanceID, err := NextID()
	if err != nil {
		return "", err
	}
	tokenID, err := NextID()
	if err != nil {
		return "", err
	}
	cmdID, err := NextID()
	if err != nil {
		return "", err
	}

	inst := projection.NewInstance(instanceID, deploymentID, dep.Version)
	pv, err := projection.VariablesFromMap(vars)
	if err != nil {
		return "", err
	}

	e.mu.Lock()
	e.instances[instanceID] = inst
	e.instMu[instanceID] = &sync.Mutex{}
	lock := e.instMu[instanceID]
	e.mu.Unlock()

	lock.Lock()
	defer lock.Unlock()

	startID, err := dep.StartEventID()
	if err != nil {
		return "", err
	}

	cmd := &eventv1.Event{
		Id:                cmdID,
		Timestamp:         nowMillis(),
		RecordType:        eventv1.Event_RECORD_TYPE_COMMAND,
		DeploymentId:      deploymentID,
		ProcessInstanceId: instanceID,
		ProcessVersion:    dep.Version,
		Element: &eventv1.Element{
			Intent: eventv1.Element_INTENT_ACTIVATING,
			Type:   eventv1.Element_TYPE_PROCESS,
			Id:     dep.ProcessID(),
			Payload: &eventv1.Element_ProcessPayload{
				ProcessPayload: &eventv1.ProcessPayload{Variables: pv},
			},
		},
	}
	if _, err := e.log.Append(ctx, cmd); err != nil {
		return "", err
	}

	emit := e.emitter(ctx, inst, cmdID)
	for _, rec := range handlers.ProcessStartRecords(dep.ProcessID(), pv) {
		if err := emit(rec); err != nil {
			return "", err
		}
	}

	inst.Tokens[tokenID] = &projection.Token{
		ID:        tokenID,
		ElementID: startID,
		Status:    projection.TokenActive,
	}
	if err := e.executor.Enter(ctx, dep, inst, tokenID, startID, emit); err != nil {
		return "", err
	}
	return instanceID, nil
}

func (e *Engine) CompleteUserTask(ctx context.Context, instanceID, elementID, tokenID string, vars map[string]any) error {
	e.mu.Lock()
	inst := e.instances[instanceID]
	lock := e.instMu[instanceID]
	var dep *deploy.Deployment
	if inst != nil {
		dep = e.deployments[inst.DeploymentID]
	}
	e.mu.Unlock()
	if inst == nil || dep == nil {
		return fmt.Errorf("NOT_FOUND: instance %q", instanceID)
	}

	lock.Lock()
	defer lock.Unlock()

	tok := inst.Tokens[tokenID]
	if tok == nil || tok.ElementID != elementID || tok.Status != projection.TokenWaiting {
		return e.reject(ctx, inst, elementID, tokenID, eventv1.Element_TYPE_USER_TASK, eventv1.Element_INTENT_COMPLETING, "INVALID_STATE", "user task is not waiting for completion")
	}
	typ, err := dep.TypeOf(elementID)
	if err != nil || typ != eventv1.Element_TYPE_USER_TASK {
		return e.reject(ctx, inst, elementID, tokenID, eventv1.Element_TYPE_USER_TASK, eventv1.Element_INTENT_COMPLETING, "NOT_FOUND", "user task element not found")
	}

	cmdID, err := NextID()
	if err != nil {
		return err
	}
	pv, err := projection.VariablesFromMap(vars)
	if err != nil {
		return err
	}

	cmd := &eventv1.Event{
		Id:                cmdID,
		Timestamp:         nowMillis(),
		RecordType:        eventv1.Event_RECORD_TYPE_COMMAND,
		DeploymentId:      inst.DeploymentID,
		ProcessInstanceId: instanceID,
		ProcessVersion:    inst.Version,
		Element: &eventv1.Element{
			Intent:  eventv1.Element_INTENT_COMPLETING,
			Type:    eventv1.Element_TYPE_USER_TASK,
			Id:      elementID,
			TokenId: tokenID,
			Payload: &eventv1.Element_ActivityPayload{
				ActivityPayload: &eventv1.ActivityPayload{Variables: pv},
			},
		},
	}
	if _, err := e.log.Append(ctx, cmd); err != nil {
		return err
	}

	return e.executor.Complete(ctx, dep, inst, tokenID, elementID, pv, e.emitter(ctx, inst, cmdID))
}

func (e *Engine) GetInstance(instanceID string) (*projection.Instance, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	inst, ok := e.instances[instanceID]
	return inst, ok
}

func (e *Engine) ListEvents(ctx context.Context, processInstanceID string) ([]*eventv1.Event, error) {
	return e.log.ReadByInstance(ctx, processInstanceID)
}

func (e *Engine) emitter(ctx context.Context, inst *projection.Instance, sourceCmdID string) Emitter {
	return func(el *eventv1.Element) error {
		id, err := NextID()
		if err != nil {
			return err
		}
		ev := &eventv1.Event{
			Id:                id,
			Timestamp:         nowMillis(),
			RecordType:        eventv1.Event_RECORD_TYPE_EVENT,
			DeploymentId:      inst.DeploymentID,
			ProcessInstanceId: inst.ID,
			ProcessVersion:    inst.Version,
			SourceRecordId:    sourceCmdID,
			Element:           el,
		}
		if _, err := e.log.Append(ctx, ev); err != nil {
			return err
		}
		inst.ApplyEvent(ev)
		return nil
	}
}

func (e *Engine) reject(ctx context.Context, inst *projection.Instance, elementID, tokenID string, typ eventv1.Element_Type, intent eventv1.Element_Intent, code, message string) error {
	id, err := NextID()
	if err != nil {
		return err
	}
	cmdID, err := NextID()
	if err != nil {
		return err
	}
	cmd := &eventv1.Event{
		Id:                cmdID,
		Timestamp:         nowMillis(),
		RecordType:        eventv1.Event_RECORD_TYPE_COMMAND,
		DeploymentId:      inst.DeploymentID,
		ProcessInstanceId: inst.ID,
		ProcessVersion:    inst.Version,
		Element: &eventv1.Element{
			Intent:  intent,
			Type:    typ,
			Id:      elementID,
			TokenId: tokenID,
		},
	}
	if _, err := e.log.Append(ctx, cmd); err != nil {
		return err
	}
	rej := &eventv1.Event{
		Id:                id,
		Timestamp:         nowMillis(),
		RecordType:        eventv1.Event_RECORD_TYPE_REJECTION,
		DeploymentId:      inst.DeploymentID,
		ProcessInstanceId: inst.ID,
		ProcessVersion:    inst.Version,
		SourceRecordId:    cmdID,
		Rejection:         &eventv1.Rejection{Code: code, Message: message},
		Element: &eventv1.Element{
			Intent:  intent,
			Type:    typ,
			Id:      elementID,
			TokenId: tokenID,
		},
	}
	if _, err := e.log.Append(ctx, rej); err != nil {
		return err
	}
	return fmt.Errorf("%s: %s", code, message)
}

func nowMillis() int64 {
	return time.Now().UnixMilli()
}
