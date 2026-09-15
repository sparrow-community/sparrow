package processing

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/sparrow-community/sparrow/processing/deploy"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

const (
	defaultMaxJobs      = 1
	defaultLockDuration = 5 * time.Minute
)

// Job is a worker-facing snapshot of a waiting job-backed task.
// It is a runtime projection, not an EventLog record.
type Job struct {
	JobType           string
	ProcessInstanceID string
	DeploymentID      string
	ElementID         string
	TokenID           string
	Variables         map[string]string
	WorkerID          string
	LockDeadline      time.Time
	// ScriptFormat / Script are set for Script Task jobs when present on the definition.
	ScriptFormat string
	Script       string
}

// ActivateRequest is the pull-worker claim for jobs of one type.
// Wait is the long-poll budget (0 returns immediately). LockDuration is the
// exclusive lease; 0 uses a 5-minute default. With a runtime.Store, leases
// survive restart until expiry or release.
type ActivateRequest struct {
	JobType      string
	MaxJobs      int
	Wait         time.Duration
	WorkerID     string
	LockDuration time.Duration
}

type jobLease struct {
	workerID string
	deadline time.Time
}

func leaseKey(instanceID, tokenID string) string {
	return instanceID + "/" + tokenID
}

// Activate long-polls for waiting ServiceTask jobs of req.JobType.
// Claimed jobs are leased so another Activate will not return them until the
// lock expires or Complete succeeds. Completion stays on Engine.Complete.
func (e *Engine) Activate(ctx context.Context, req ActivateRequest) ([]Job, error) {
	if req.JobType == "" {
		return nil, fmt.Errorf("INVALID_ARGUMENT: job_type is required")
	}
	if req.MaxJobs < 0 {
		return nil, fmt.Errorf("INVALID_ARGUMENT: max_jobs must be >= 0")
	}
	if req.MaxJobs == 0 {
		req.MaxJobs = defaultMaxJobs
	}
	if req.LockDuration <= 0 {
		req.LockDuration = defaultLockDuration
	}

	var waitUntil time.Time
	if req.Wait > 0 {
		waitUntil = time.Now().Add(req.Wait)
	}

	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		jobs, nextRetry := e.claimJobs(req)
		if len(jobs) > 0 {
			return jobs, nil
		}
		if req.Wait <= 0 {
			return nil, nil
		}
		now := time.Now()
		if !waitUntil.After(now) {
			return nil, nil
		}
		wakeAt := waitUntil
		if !nextRetry.IsZero() && nextRetry.Before(wakeAt) {
			wakeAt = nextRetry
		}
		if err := e.waitForJobs(ctx, wakeAt); err != nil {
			return nil, err
		}
	}
}

func (e *Engine) claimJobs(req ActivateRequest) ([]Job, time.Time) {
	now := time.Now()
	deadline := now.Add(req.LockDuration)

	e.mu.Lock()
	ids := make([]string, 0, len(e.instances))
	for id := range e.instances {
		ids = append(ids, id)
	}
	e.mu.Unlock()
	slices.Sort(ids)

	var out []Job
	var nextRetry time.Time
	for _, iid := range ids {
		if len(out) >= req.MaxJobs {
			break
		}
		e.mu.Lock()
		inst := e.instances[iid]
		lock := e.instMu[iid]
		e.mu.Unlock()
		if inst == nil || lock == nil {
			continue
		}

		lock.Lock()
		tokenIDs := make([]string, 0, len(inst.Tokens))
		for tid := range inst.Tokens {
			tokenIDs = append(tokenIDs, tid)
		}
		slices.Sort(tokenIDs)

		e.jobMu.Lock()
		for _, tid := range tokenIDs {
			if len(out) >= req.MaxJobs {
				break
			}
			tok := inst.Tokens[tid]
			if tok == nil || tok.Status != projection.TokenWaiting || tok.JobType != req.JobType {
				continue
			}
			key := leaseKey(iid, tid)
			if lease, ok := e.leases[key]; ok && now.Before(lease.deadline) {
				if nextRetry.IsZero() || lease.deadline.Before(nextRetry) {
					nextRetry = lease.deadline
				}
				continue
			}
			e.leases[key] = jobLease{workerID: req.WorkerID, deadline: deadline}
			e.persistLease(iid, tid, req.WorkerID, deadline)
			job := Job{
				JobType:           tok.JobType,
				ProcessInstanceID: inst.ID,
				DeploymentID:      inst.DeploymentID,
				ElementID:         tok.ElementID,
				TokenID:           tok.ID,
				Variables:         cloneStringMap(inst.Variables),
				WorkerID:          req.WorkerID,
				LockDeadline:      deadline,
			}
			e.mu.Lock()
			dep := e.deployments[inst.DeploymentID]
			e.mu.Unlock()
			if dep != nil {
				if st, err := dep.ScriptTask(tok.ElementID); err == nil {
					job.ScriptFormat = st.ScriptFormat
					job.Script = st.Script
				}
			}
			out = append(out, job)
		}
		e.jobMu.Unlock()
		lock.Unlock()
	}
	return out, nextRetry
}

func (e *Engine) waitForJobs(ctx context.Context, until time.Time) error {
	delay := time.Until(until)
	if delay <= 0 {
		return nil
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	case <-e.jobWake:
		return nil
	}
}

func (e *Engine) notifyJobs() {
	select {
	case e.jobWake <- struct{}{}:
	default:
	}
}

func (e *Engine) releaseLease(instanceID, tokenID string) {
	e.jobMu.Lock()
	delete(e.leases, leaseKey(instanceID, tokenID))
	e.jobMu.Unlock()
	e.persistReleaseLease(instanceID, tokenID)
}

// Fail records SERVICE_TASK FAILED, keeps the token waiting (or opens an incident),
// and drops the lease so another Activate can claim the job when retriable.
// UserTask and non-waiting elements are rejected. Completing the activity is Engine.Complete.
func (e *Engine) Fail(ctx context.Context, instanceID, elementID, tokenID, message string, noRetry bool) error {
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

	typ, typeErr := dep.TypeOf(elementID)
	if typeErr != nil {
		typ = eventv1.Element_TYPE_UNSPECIFIED
	}

	tok := inst.Tokens[tokenID]
	if tok == nil || tok.ElementID != elementID {
		return e.reject(ctx, inst, elementID, tokenID, typ, eventv1.Element_INTENT_FAILED, "INVALID_STATE", "element is not waiting for completion")
	}
	if tok.Status == projection.TokenBlocked {
		return e.reject(ctx, inst, elementID, tokenID, typ, eventv1.Element_INTENT_FAILED, "INCIDENT_OPEN", "incident is open")
	}
	if tok.Status != projection.TokenWaiting {
		return e.reject(ctx, inst, elementID, tokenID, typ, eventv1.Element_INTENT_FAILED, "INVALID_STATE", "element is not waiting for completion")
	}
	if typeErr != nil {
		return e.reject(ctx, inst, elementID, tokenID, typ, eventv1.Element_INTENT_FAILED, "NOT_FOUND", "element not found")
	}
	if tok.JobType == "" {
		return e.reject(ctx, inst, elementID, tokenID, typ, eventv1.Element_INTENT_FAILED, "INVALID_STATE", "element is not a job")
	}

	failCount := tok.JobFailCount + 1
	cmdID, err := NextID()
	if err != nil {
		return err
	}
	cmdPayload := jobFailActivityPayload(tok, message, failCount, noRetry)
	cmd := &eventv1.Event{
		Id:                cmdID,
		Timestamp:         nowMillis(),
		RecordType:        eventv1.Event_RECORD_TYPE_COMMAND,
		DeploymentId:      inst.DeploymentID,
		ProcessInstanceId: instanceID,
		ProcessVersion:    inst.Version,
		Element: &eventv1.Element{
			Intent:  eventv1.Element_INTENT_FAILED,
			Type:    typ,
			Id:      elementID,
			TokenId: tokenID,
			Payload: &eventv1.Element_ActivityPayload{ActivityPayload: cmdPayload},
		},
	}
	if _, err := e.log.Append(ctx, cmd); err != nil {
		return err
	}

	emit := e.emitter(ctx, inst, cmdID)
	if err := emitJobFailChain(dep, inst, typ, elementID, tokenID, message, failCount, noRetry, emit); err != nil {
		return err
	}

	e.releaseLease(instanceID, tokenID)
	e.notifyJobs()
	return nil
}

// Heartbeat extends an in-memory job lease. It does not write the EventLog.
// workerID must match the Activate holder.
func (e *Engine) Heartbeat(_ context.Context, instanceID, tokenID, workerID string, lockDuration time.Duration) error {
	if lockDuration <= 0 {
		lockDuration = defaultLockDuration
	}

	e.mu.Lock()
	inst := e.instances[instanceID]
	lock := e.instMu[instanceID]
	e.mu.Unlock()
	if inst == nil || lock == nil {
		return fmt.Errorf("NOT_FOUND: instance %q", instanceID)
	}

	lock.Lock()
	defer lock.Unlock()

	tok := inst.Tokens[tokenID]
	if tok == nil || tok.Status != projection.TokenWaiting || tok.JobType == "" {
		return fmt.Errorf("INVALID_STATE: job is not waiting")
	}

	key := leaseKey(instanceID, tokenID)
	now := time.Now()
	e.jobMu.Lock()
	lease, ok := e.leases[key]
	if !ok || !now.Before(lease.deadline) {
		e.jobMu.Unlock()
		return fmt.Errorf("INVALID_STATE: job is not locked")
	}
	if lease.workerID != workerID {
		e.jobMu.Unlock()
		return fmt.Errorf("INVALID_STATE: job is locked by another worker")
	}
	newDeadline := now.Add(lockDuration)
	e.leases[key] = jobLease{workerID: workerID, deadline: newDeadline}
	e.jobMu.Unlock()
	e.persistLease(instanceID, tokenID, workerID, newDeadline)
	return nil
}

func cloneStringMap(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
