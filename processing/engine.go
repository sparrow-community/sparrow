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
// Persistence is injected: EventLog (behavior ledger) and optional deploy.Store (BPMN XML).
type Engine struct {
	log      eventlog.EventLog
	store    deploy.Store
	executor *Executor

	mu          sync.Mutex
	deployments map[string]*deploy.Deployment
	instances   map[string]*projection.Instance
	instMu      map[string]*sync.Mutex

	jobMu   sync.Mutex
	leases  map[string]jobLease
	jobWake chan struct{}

	msgMu  sync.Mutex
	msgBuf []bufferedMessage

	seenMu sync.Mutex
	seen   map[string]struct{}

	nowFn func() time.Time
}

func NewEngine(l eventlog.EventLog) *Engine {
	if l == nil {
		l = eventlog.NewMemory()
	}
	e := &Engine{
		log:         l,
		deployments: make(map[string]*deploy.Deployment),
		instances:   make(map[string]*projection.Instance),
		instMu:      make(map[string]*sync.Mutex),
		leases:      make(map[string]jobLease),
		jobWake:     make(chan struct{}, 1),
		seen:        make(map[string]struct{}),
	}
	e.executor = NewExecutor(handlers.DefaultRegistry())
	e.executor.Now = e.now
	return e
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
	if err := e.persistDeployment(id, bpmnXML); err != nil {
		return "", fmt.Errorf("persist deployment: %w", err)
	}

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

	startID, err := dep.StartEventID()
	if err != nil {
		return "", err
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

	pv, err := projection.VariablesFromMap(vars)
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
			Intent:  eventv1.Element_INTENT_ACTIVATING,
			Type:    eventv1.Element_TYPE_PROCESS,
			Id:      dep.ProcessID(),
			TokenId: tokenID,
			Payload: &eventv1.Element_ProcessPayload{
				ProcessPayload: &eventv1.ProcessPayload{Variables: pv},
			},
		},
	}
	if _, err := e.log.Append(ctx, cmd); err != nil {
		return "", err
	}

	inst := projection.NewInstance(instanceID, deploymentID, dep.Version)
	lock := &sync.Mutex{}
	e.mu.Lock()
	e.instances[instanceID] = inst
	e.instMu[instanceID] = lock
	e.mu.Unlock()

	lock.Lock()
	emit := e.emitter(ctx, inst, cmdID)
	for _, rec := range handlers.ProcessStartRecords(dep.ProcessID(), pv) {
		if err := emit(rec); err != nil {
			lock.Unlock()
			return "", err
		}
	}
	err = e.executor.Enter(ctx, dep, inst, tokenID, startID, emit)
	lock.Unlock()
	if err != nil {
		return "", err
	}
	if err := e.tryDeliverBuffered(ctx, instanceID); err != nil {
		return "", err
	}
	return instanceID, nil
}

func (e *Engine) Complete(ctx context.Context, instanceID, elementID, tokenID string, vars map[string]any) error {
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
	err := e.completeLocked(ctx, dep, inst, instanceID, elementID, tokenID, vars)
	lock.Unlock()
	if err != nil {
		return err
	}
	return e.tryDeliverBuffered(ctx, instanceID)
}

func (e *Engine) completeLocked(ctx context.Context, dep *deploy.Deployment, inst *projection.Instance, instanceID, elementID, tokenID string, vars map[string]any) error {
	typ, typeErr := dep.TypeOf(elementID)
	if typeErr != nil {
		typ = eventv1.Element_TYPE_UNSPECIFIED
	}

	tok := inst.Tokens[tokenID]
	if tok == nil || tok.ElementID != elementID || tok.Status != projection.TokenWaiting {
		return e.reject(ctx, inst, elementID, tokenID, typ, eventv1.Element_INTENT_COMPLETING, "INVALID_STATE", "element is not waiting for completion")
	}
	if typeErr != nil {
		return e.reject(ctx, inst, elementID, tokenID, typ, eventv1.Element_INTENT_COMPLETING, "NOT_FOUND", "element not found")
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
			Type:    typ,
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

	if err := e.executor.Complete(ctx, dep, inst, tokenID, elementID, pv, e.emitter(ctx, inst, cmdID)); err != nil {
		return err
	}
	e.releaseLease(instanceID, tokenID)
	return nil
}

func (e *Engine) GetInstance(instanceID string) (*projection.Instance, bool) {
	e.mu.Lock()
	inst, ok := e.instances[instanceID]
	lock := e.instMu[instanceID]
	e.mu.Unlock()
	if !ok || inst == nil || lock == nil {
		return nil, false
	}
	lock.Lock()
	defer lock.Unlock()
	return inst.Clone(), true
}

func (e *Engine) ListEvents(ctx context.Context, processInstanceID string) ([]*eventv1.Event, error) {
	return e.log.ReadByInstance(ctx, processInstanceID)
}

func (e *Engine) emitter(ctx context.Context, inst *projection.Instance, sourceCmdID string) Emitter {
	return func(el *eventv1.Element) error {
		if e.alreadySeen(sourceCmdID, el) {
			return nil
		}
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
		e.markSeen(sourceCmdID, el)
		inst.ApplyEvent(ev)
		if el.GetType() == eventv1.Element_TYPE_SERVICE_TASK &&
			el.GetIntent() == eventv1.Element_INTENT_ACTIVATED {
			e.notifyJobs()
		}
		return nil
	}
}

func intentKey(sourceCmdID string, el *eventv1.Element) string {
	if el == nil {
		return sourceCmdID
	}
	return fmt.Sprintf("%s/%d/%d/%s/%s", sourceCmdID, int32(el.GetType()), int32(el.GetIntent()), el.GetId(), el.GetTokenId())
}

func (e *Engine) alreadySeen(sourceCmdID string, el *eventv1.Element) bool {
	key := intentKey(sourceCmdID, el)
	e.seenMu.Lock()
	defer e.seenMu.Unlock()
	_, ok := e.seen[key]
	return ok
}

func (e *Engine) markSeen(sourceCmdID string, el *eventv1.Element) {
	key := intentKey(sourceCmdID, el)
	e.seenMu.Lock()
	e.seen[key] = struct{}{}
	e.seenMu.Unlock()
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

func (e *Engine) now() time.Time {
	if e != nil && e.nowFn != nil {
		return e.nowFn()
	}
	return time.Now()
}
