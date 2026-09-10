package processing

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/sparrow-community/sparrow/processing/deploy"
	"github.com/sparrow-community/sparrow/processing/handlers"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	"github.com/sparrow-community/sparrow/processing/runtime"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// Engine is the single-node BPMN execution facade.
// Element semantics live in handlers; definitions come from bpmn via deploy.Deployment.
// Persistence is injected: EventLog (behavior ledger), optional deploy.Store (BPMN XML),
// and optional runtime.Store (job leases + message buffer).
type Engine struct {
	log          eventlog.EventLog
	store        deploy.Store
	runtimeStore runtime.Store
	executor     *Executor

	mu          sync.Mutex
	deployments map[string]*deploy.Deployment
	revisions   map[string][]deploy.Revision // process id -> ordered revisions
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
		revisions:   make(map[string][]deploy.Revision),
		instances:   make(map[string]*projection.Instance),
		instMu:      make(map[string]*sync.Mutex),
		leases:      make(map[string]jobLease),
		jobWake:     make(chan struct{}, 1),
		seen:        make(map[string]struct{}),
	}
	e.executor = NewExecutor(handlers.DefaultRegistry())
	e.executor.Now = e.now
	e.wireCalleeResolver()
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
	e.revisions = deploy.AssignProcessVersions(e.deployments)
	e.mu.Unlock()
	return id, nil
}

// CreateInstanceRequest selects a deployment by id and/or process revision.
type CreateInstanceRequest struct {
	DeploymentID   string
	ProcessID      string
	ProcessVersion int32 // with ProcessID; 0 means latest
	Variables      map[string]any
}

func (e *Engine) CreateInstance(ctx context.Context, deploymentID string, vars map[string]any) (string, error) {
	return e.CreateInstanceRequest(ctx, CreateInstanceRequest{
		DeploymentID: deploymentID,
		Variables:    vars,
	})
}

func (e *Engine) CreateInstanceRequest(ctx context.Context, req CreateInstanceRequest) (string, error) {
	e.mu.Lock()
	dep, err := e.resolveDeploymentLocked(req.DeploymentID, req.ProcessID, req.ProcessVersion)
	e.mu.Unlock()
	if err != nil {
		return "", err
	}
	deploymentID := dep.ID

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

	pv, err := projection.VariablesFromMap(req.Variables)
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
	inst.ProcessID = dep.ProcessID()
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
	if err := emitEventSubProcessStartArms(dep, inst, dep.ProcessID(), e.now(), emit); err != nil {
		lock.Unlock()
		return "", err
	}
	pubs, err := e.executor.Enter(ctx, dep, inst, tokenID, startID, emit)
	if err != nil {
		e.rejectEnterFailure(ctx, dep, inst, err)
		lock.Unlock()
		return "", err
	}
	lock.Unlock()
	if err := e.flushPublications(ctx, pubs); err != nil {
		return "", err
	}
	if err := e.tryDeliverBuffered(ctx, instanceID); err != nil {
		return "", err
	}
	return instanceID, nil
}

func (e *Engine) resolveDeploymentLocked(deploymentID, processID string, version int32) (*deploy.Deployment, error) {
	if deploymentID != "" {
		dep := e.deployments[deploymentID]
		if dep == nil {
			return nil, fmt.Errorf("NOT_FOUND: deployment %q", deploymentID)
		}
		if processID != "" && dep.ProcessID() != processID {
			return nil, fmt.Errorf("INVALID_ARGUMENT: deployment %q process_id %q does not match %q", deploymentID, dep.ProcessID(), processID)
		}
		if version != 0 && dep.Version != version {
			return nil, fmt.Errorf("INVALID_ARGUMENT: deployment %q process_version %d does not match %d", deploymentID, dep.Version, version)
		}
		return dep, nil
	}
	if processID == "" {
		return nil, fmt.Errorf("INVALID_ARGUMENT: deployment_id or process_id is required")
	}
	id, ok := deploy.ResolveRevision(e.revisions, processID, version)
	if !ok {
		if version > 0 {
			return nil, fmt.Errorf("NOT_FOUND: process %q version %d", processID, version)
		}
		return nil, fmt.Errorf("NOT_FOUND: process %q", processID)
	}
	dep := e.deployments[id]
	if dep == nil {
		return nil, fmt.Errorf("NOT_FOUND: deployment %q", id)
	}
	return dep, nil
}

// GetDeployment returns a compiled deployment by id.
func (e *Engine) GetDeployment(deploymentID string) (*deploy.Deployment, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	dep, ok := e.deployments[deploymentID]
	return dep, ok
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
	pubs, err := e.completeLocked(ctx, dep, inst, instanceID, elementID, tokenID, vars)
	lock.Unlock()
	if err != nil {
		return err
	}
	if err := e.flushPublications(ctx, pubs); err != nil {
		return err
	}
	return e.tryDeliverBuffered(ctx, instanceID)
}

// ResolveIncident closes an open Service Task incident and restores waiting job semantics.
func (e *Engine) ResolveIncident(ctx context.Context, instanceID, elementID, tokenID string) error {
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
	if inst.Status == projection.StatusCompleted || inst.Status == projection.StatusTerminated {
		return fmt.Errorf("INVALID_STATE: instance is not active")
	}

	lock.Lock()
	defer lock.Unlock()

	typ, typeErr := dep.TypeOf(elementID)
	if typeErr != nil {
		typ = eventv1.Element_TYPE_UNSPECIFIED
	}

	tok := inst.Tokens[tokenID]
	if tok == nil || tok.ElementID != elementID {
		return e.reject(ctx, inst, elementID, tokenID, typ, eventv1.Element_INTENT_INCIDENT_RESOLVED, "NO_INCIDENT", "no open incident on token")
	}
	if tok.Status != projection.TokenBlocked {
		return e.reject(ctx, inst, elementID, tokenID, typ, eventv1.Element_INTENT_INCIDENT_RESOLVED, "NO_INCIDENT", "no open incident on token")
	}
	if typeErr != nil {
		return e.reject(ctx, inst, elementID, tokenID, typ, eventv1.Element_INTENT_INCIDENT_RESOLVED, "NOT_FOUND", "element not found")
	}

	cmdID, err := NextID()
	if err != nil {
		return err
	}
	resolvePayload := &eventv1.ActivityPayload{JobType: tok.JobType}
	cmd := &eventv1.Event{
		Id:                cmdID,
		Timestamp:         nowMillis(),
		RecordType:        eventv1.Event_RECORD_TYPE_COMMAND,
		DeploymentId:      inst.DeploymentID,
		ProcessInstanceId: instanceID,
		ProcessVersion:    inst.Version,
		Element: &eventv1.Element{
			Intent:  eventv1.Element_INTENT_INCIDENT_RESOLVED,
			Type:    typ,
			Id:      elementID,
			TokenId: tokenID,
			Payload: &eventv1.Element_ActivityPayload{ActivityPayload: resolvePayload},
		},
	}
	if _, err := e.log.Append(ctx, cmd); err != nil {
		return err
	}

	emit := e.emitter(ctx, inst, cmdID)
	if err := emit(&eventv1.Element{
		Intent:  eventv1.Element_INTENT_INCIDENT_RESOLVED,
		Type:    typ,
		Id:      elementID,
		TokenId: tokenID,
		Payload: &eventv1.Element_ActivityPayload{ActivityPayload: resolvePayload},
	}); err != nil {
		return err
	}

	e.notifyJobs()
	return nil
}

func (e *Engine) completeLocked(ctx context.Context, dep *deploy.Deployment, inst *projection.Instance, instanceID, elementID, tokenID string, vars map[string]any) ([]handlers.Publication, error) {
	typ, typeErr := dep.TypeOf(elementID)
	if typeErr != nil {
		typ = eventv1.Element_TYPE_UNSPECIFIED
	}

	tok := inst.Tokens[tokenID]
	if tok != nil && tok.Status == projection.TokenBlocked && tok.ElementID == elementID {
		return nil, e.reject(ctx, inst, elementID, tokenID, typ, eventv1.Element_INTENT_COMPLETING, "INCIDENT_OPEN", "incident is open")
	}
	waiting := tok != nil && (tok.Status == projection.TokenWaiting || tok.Status == projection.TokenBlocked)
	if waiting && typ == eventv1.Element_TYPE_BOUNDARY_EVENT {
		attached, ok := dep.AttachedActivity(elementID)
		waiting = ok && tok.ElementID == attached && projection.TokenHasArmedBoundary(tok, elementID)
	} else if waiting {
		waiting = tok.ElementID == elementID
	}
	if tok == nil || !waiting {
		return nil, e.reject(ctx, inst, elementID, tokenID, typ, eventv1.Element_INTENT_COMPLETING, "INVALID_STATE", "element is not waiting for completion")
	}
	if typeErr != nil {
		return nil, e.reject(ctx, inst, elementID, tokenID, typ, eventv1.Element_INTENT_COMPLETING, "NOT_FOUND", "element not found")
	}

	cmdID, err := NextID()
	if err != nil {
		return nil, err
	}
	pv, err := projection.VariablesFromMap(vars)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	pubs, err := e.executor.Complete(ctx, dep, inst, tokenID, elementID, pv, e.emitter(ctx, inst, cmdID))
	if err != nil {
		e.rejectEnterFailure(ctx, dep, inst, err)
		return pubs, err
	}
	e.releaseLease(instanceID, tokenID)
	return pubs, nil
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
		prevStatus := inst.Status
		inst.ApplyEvent(ev)
		if el.GetType() == eventv1.Element_TYPE_PROCESS {
			switch el.GetIntent() {
			case eventv1.Element_INTENT_COMPLETED, eventv1.Element_INTENT_TERMINATED:
				if prevStatus != inst.Status {
					e.sweepRuntimeForInstance(inst.ID)
				}
			}
		}
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

func rejectionFromError(err error) (code, message string) {
	if err == nil {
		return "INTERNAL", "unknown error"
	}
	s := err.Error()
	if i := strings.Index(s, ":"); i > 0 {
		return strings.TrimSpace(s[:i]), strings.TrimSpace(s[i+1:])
	}
	return "INTERNAL", s
}

func activeEnterFailureToken(inst *projection.Instance) (elementID, tokenID string, ok bool) {
	for id, tok := range inst.Tokens {
		if tok != nil && tok.Status == projection.TokenActive && tok.ElementID != "" {
			return tok.ElementID, id, true
		}
	}
	return "", "", false
}

func (e *Engine) rejectEnterFailure(ctx context.Context, dep *deploy.Deployment, inst *projection.Instance, cause error) {
	elementID, tokenID, ok := activeEnterFailureToken(inst)
	if !ok {
		return
	}
	typ, err := dep.TypeOf(elementID)
	if err != nil {
		return
	}
	code, msg := rejectionFromError(cause)
	_ = e.reject(ctx, inst, elementID, tokenID, typ, eventv1.Element_INTENT_ACTIVATING, code, msg)
}

func (e *Engine) now() time.Time {
	if e != nil && e.nowFn != nil {
		return e.nowFn()
	}
	return time.Now()
}

// completeScopeBoundary fires a boundary attached to a SubProcess.
// It terminates all tokens inside the scope, terminates the SubProcess,
// completes the boundary, and takes the boundary's outgoing flow.
func (e *Engine) completeScopeBoundary(ctx context.Context, instanceID, boundaryID string) error {
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
	pubs, err := e.completeScopeBoundaryLocked(ctx, dep, inst, instanceID, boundaryID)
	lock.Unlock()
	if err != nil {
		return err
	}
	if err := e.flushPublications(ctx, pubs); err != nil {
		return err
	}
	return e.tryDeliverBuffered(ctx, instanceID)
}

func (e *Engine) completeScopeBoundaryLocked(ctx context.Context, dep *deploy.Deployment, inst *projection.Instance, instanceID, boundaryID string) ([]handlers.Publication, error) {
	sb, ok := inst.ScopeBoundaries[boundaryID]
	if !ok {
		return nil, nil // already disarmed
	}
	scopeID := sb.ScopeID
	tokenID := sb.TokenID

	cmdID, err := NextID()
	if err != nil {
		return nil, err
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
			Type:    eventv1.Element_TYPE_BOUNDARY_EVENT,
			Id:      boundaryID,
			TokenId: tokenID,
		},
	}
	if _, err := e.log.Append(ctx, cmd); err != nil {
		return nil, err
	}

	emit := e.emitter(ctx, inst, cmdID)

	if err := terminateScopeTokens(dep, inst, scopeID, emit, scopeTerminateOpts{
		IncludeHost: false,
		DropTokens:  true,
	}); err != nil {
		return nil, err
	}
	if err := terminateEmbeddedScope(scopeID, tokenID, emit); err != nil {
		return nil, err
	}

	// Remove scope boundaries
	inst.RemoveScopeBoundariesForScope(scopeID)
	if err := emitEventSubProcessStartDisarmInScope(dep, scopeID, inst, emit); err != nil {
		return nil, err
	}

	// Complete the boundary and take outgoing
	if err := emit(&eventv1.Element{
		Intent:  eventv1.Element_INTENT_COMPLETING,
		Type:    eventv1.Element_TYPE_BOUNDARY_EVENT,
		Id:      boundaryID,
		TokenId: tokenID,
	}); err != nil {
		return nil, err
	}
	if err := emit(&eventv1.Element{
		Intent:  eventv1.Element_INTENT_COMPLETED,
		Type:    eventv1.Element_TYPE_BOUNDARY_EVENT,
		Id:      boundaryID,
		TokenId: tokenID,
	}); err != nil {
		return nil, err
	}

	// Take boundary outgoing flow
	next, err := e.executor.takeOutgoing(dep, inst, tokenID, boundaryID, "", emit)
	if err != nil {
		return nil, err
	}
	return e.executor.Enter(ctx, dep, inst, tokenID, next, emit)
}

// isInScope checks if tokScope is inside targetScope (direct or nested).
func isInScope(dep *deploy.Deployment, tokScope, targetScope string) bool {
	for tokScope != "" {
		if tokScope == targetScope {
			return true
		}
		parent, ok := dep.ScopeOf(tokScope)
		if !ok {
			return false
		}
		tokScope = parent
	}
	return false
}
