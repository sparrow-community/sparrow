// Copyright 2026 The Sparrow community and contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package processing

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/sparrow-community/sparrow/processing/deploy"
	"github.com/sparrow-community/sparrow/processing/handlers"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// RunPolicy controls when Enter-settled barriers pause an intervention session.
type RunPolicy string

const (
	RunPolicyContinuous  RunPolicy = "continuous"  // only stop on Wait / end / error (and P1 breakpoints)
	RunPolicyStep        RunPolicy = "step"        // stop at every Enter-settled barrier
	RunPolicyBreakpoints RunPolicy = "breakpoints" // P1: stop only on BP hit (and Wait)
)

// PauseReason explains why the session is barrier-paused.
type PauseReason string

const (
	PauseReasonStep       PauseReason = "step"
	PauseReasonBreakpoint PauseReason = "breakpoint" // P1
	PauseReasonManual     PauseReason = "manual"     // optional; not required for P0
)

// PendingKind classifies the transit held at an Enter-settled barrier.
type PendingKind string

const (
	PendingLeave      PendingKind = "leave"
	PendingFork       PendingKind = "fork"
	PendingEnterChild PendingKind = "enterChild"
	PendingLink       PendingKind = "link"
	// PendingDecide: exclusive gateway ACTIVATED; choose on resume.
	PendingDecide PendingKind = "decide"
)

// PendingTransition is ephemeral resume work; never written to the EventLog.
type PendingTransition struct {
	Kind            PendingKind
	FromElementID   string
	TokenID         string
	TakenFlowIDs    []string
	NextElementIDs  []string
	OutgoingFlowID  string
	EnterChildID    string
	SpawnChildToken bool
	LinkCatchIDs    []string
}

// InterventionSession is the opt-in control plane for pause/step/continue.
// It is not ledger state; Recover ignores it.
type InterventionSession struct {
	Enabled         bool
	FocusInstanceID string
	Breakpoints     map[string]struct{}
	RunPolicy       RunPolicy

	Paused         bool
	PauseReason    PauseReason
	PauseElementID string
	PauseTokenID   string

	// SourceCmdID is the interrupted business COMMAND id; resume EVENTs reuse it.
	SourceCmdID string

	Pending      *PendingTransition
	DeferredPubs []handlers.Publication
	SuppressKey  string
}

// InterventionState is the read model returned by GetInterventionState.
type InterventionState struct {
	Enabled         bool
	FocusInstanceID string
	Policy          RunPolicy
	Breakpoints     []string
	Paused          bool
	PauseReason     PauseReason
	PauseElementID  string
	PauseTokenID    string
	Pending         *PendingTransition
}

type EnableInterventionRequest struct {
	InstanceID string
	Policy     RunPolicy
}

type EnableInterventionResponse struct {
	OK bool
}

type DisableInterventionRequest struct{}

type DisableInterventionResponse struct {
	OK bool
}

type SetBreakpointsRequest struct {
	InstanceID string
	ElementIDs []string
}

type SetBreakpointsResponse struct {
	OK         bool
	ElementIDs []string
}

type ContinueRequest struct {
	InstanceID string
}

type ContinueResponse struct {
	OK     bool
	Paused bool
	State  InterventionState
}

type StepIntoRequest struct {
	InstanceID string
}

type StepIntoResponse struct {
	OK     bool
	Paused bool
	State  InterventionState
}

type StepOverRequest struct {
	InstanceID string
}

type StepOverResponse struct {
	OK     bool
	Paused bool
	State  InterventionState
}

// interventionCtl holds the session and fork-suppress depth for the executor.
type interventionCtl struct {
	mu                   sync.Mutex
	session              InterventionSession
	forkSuppressDepth    int
	decideFinalizeDepth  int // suppress leave barrier while finalizing PendingDecide
}

func (c *interventionCtl) snapshot() InterventionState {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.snapshotLocked()
}

func (c *interventionCtl) snapshotLocked() InterventionState {
	s := c.session
	bps := make([]string, 0, len(s.Breakpoints))
	for id := range s.Breakpoints {
		bps = append(bps, id)
	}
	sort.Strings(bps)
	var pending *PendingTransition
	if s.Pending != nil {
		cp := *s.Pending
		pending = &cp
		pending.TakenFlowIDs = append([]string(nil), s.Pending.TakenFlowIDs...)
		pending.NextElementIDs = append([]string(nil), s.Pending.NextElementIDs...)
		pending.LinkCatchIDs = append([]string(nil), s.Pending.LinkCatchIDs...)
	}
	return InterventionState{
		Enabled:         s.Enabled,
		FocusInstanceID: s.FocusInstanceID,
		Policy:          s.RunPolicy,
		Breakpoints:     bps,
		Paused:          s.Paused,
		PauseReason:     s.PauseReason,
		PauseElementID:  s.PauseElementID,
		PauseTokenID:    s.PauseTokenID,
		Pending:         pending,
	}
}

func (c *interventionCtl) clearToDefaultLocked() {
	c.session = InterventionSession{}
	c.forkSuppressDepth = 0
	c.decideFinalizeDepth = 0
}

func barrierKey(p *PendingTransition) string {
	if p == nil {
		return ""
	}
	return p.TokenID + "|" + p.FromElementID + "|" + string(p.Kind)
}

func (c *interventionCtl) shouldBarrier(instanceID, elementID, tokenID, sourceCmdID string, pending *PendingTransition) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	s := &c.session
	if !s.Enabled || s.Paused {
		return false
	}
	if s.FocusInstanceID != "" && s.FocusInstanceID != instanceID {
		return false
	}
	if c.forkSuppressDepth > 0 || c.decideFinalizeDepth > 0 {
		return false
	}
	key := barrierKey(pending)
	if s.SuppressKey != "" && s.SuppressKey == key {
		s.SuppressKey = ""
		return false
	}
	switch s.RunPolicy {
	case RunPolicyStep:
		c.armPauseLocked(pending, elementID, tokenID, PauseReasonStep, sourceCmdID)
		return true
	case RunPolicyBreakpoints, RunPolicyContinuous:
		if _, hit := s.Breakpoints[elementID]; hit {
			c.armPauseLocked(pending, elementID, tokenID, PauseReasonBreakpoint, sourceCmdID)
			return true
		}
		return false
	default:
		return false
	}
}

func (c *interventionCtl) armPauseLocked(pending *PendingTransition, elementID, tokenID string, reason PauseReason, sourceCmdID string) {
	s := &c.session
	s.Paused = true
	s.PauseReason = reason
	s.PauseElementID = elementID
	s.PauseTokenID = tokenID
	s.SourceCmdID = sourceCmdID
	if pending != nil {
		cp := *pending
		cp.TakenFlowIDs = append([]string(nil), pending.TakenFlowIDs...)
		cp.NextElementIDs = append([]string(nil), pending.NextElementIDs...)
		cp.LinkCatchIDs = append([]string(nil), pending.LinkCatchIDs...)
		s.Pending = &cp
	} else {
		s.Pending = nil
	}
}

func (c *interventionCtl) isPausedFocus(instanceID string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.session.Enabled && c.session.Paused && c.session.FocusInstanceID == instanceID
}

func (c *interventionCtl) stashPubs(pubs []handlers.Publication) {
	if len(pubs) == 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.session.DeferredPubs = append(c.session.DeferredPubs, pubs...)
}

func (c *interventionCtl) takeDeferredPubs() []handlers.Publication {
	c.mu.Lock()
	defer c.mu.Unlock()
	pubs := c.session.DeferredPubs
	c.session.DeferredPubs = nil
	return pubs
}

func (e *Engine) errIfInstancePaused(instanceID string) error {
	if e == nil || e.executor == nil || e.executor.intervention == nil {
		return nil
	}
	if e.executor.intervention.isPausedFocus(instanceID) {
		return fmt.Errorf("INVALID_STATE: instance paused")
	}
	return nil
}

// afterUnlockedCommand flushes publications unless the focused instance is mid-barrier-pause.
func (e *Engine) afterUnlockedCommand(ctx context.Context, instanceID string, pubs []handlers.Publication) error {
	ctl := e.executor.intervention
	if ctl != nil && ctl.isPausedFocus(instanceID) {
		ctl.stashPubs(pubs)
		return nil
	}
	var all []handlers.Publication
	if ctl != nil {
		all = append(all, ctl.takeDeferredPubs()...)
	}
	all = append(all, pubs...)
	if err := e.flushPublications(ctx, all); err != nil {
		return err
	}
	return e.tryDeliverBuffered(ctx, instanceID)
}

func (e *Engine) EnableIntervention(_ context.Context, req EnableInterventionRequest) (*EnableInterventionResponse, error) {
	if e == nil || e.executor == nil || e.executor.intervention == nil {
		return nil, fmt.Errorf("INVALID_STATE: engine not ready")
	}
	instanceID := req.InstanceID
	if instanceID == "" {
		return nil, fmt.Errorf("INVALID_ARGUMENT: instance_id is required")
	}
	e.mu.Lock()
	_, ok := e.instances[instanceID]
	e.mu.Unlock()
	if !ok {
		return nil, fmt.Errorf("NOT_FOUND: instance %q", instanceID)
	}
	policy := req.Policy
	if policy == "" {
		policy = RunPolicyContinuous
	}
	switch policy {
	case RunPolicyContinuous, RunPolicyStep, RunPolicyBreakpoints:
	default:
		return nil, fmt.Errorf("INVALID_ARGUMENT: unknown run policy %q", policy)
	}

	ctl := e.executor.intervention
	ctl.mu.Lock()
	defer ctl.mu.Unlock()
	if ctl.session.Paused && ctl.session.FocusInstanceID != "" && ctl.session.FocusInstanceID != instanceID {
		return nil, fmt.Errorf("INVALID_STATE: intervene continue or discard")
	}
	if !ctl.session.Paused {
		ctl.session.Pending = nil
		ctl.session.PauseReason = ""
		ctl.session.PauseElementID = ""
		ctl.session.PauseTokenID = ""
		ctl.session.SuppressKey = ""
		ctl.session.DeferredPubs = nil
	}
	ctl.session.Enabled = true
	ctl.session.FocusInstanceID = instanceID
	ctl.session.RunPolicy = policy
	if ctl.session.Breakpoints == nil {
		ctl.session.Breakpoints = make(map[string]struct{})
	}
	return &EnableInterventionResponse{OK: true}, nil
}

func (e *Engine) DisableIntervention(_ context.Context, _ DisableInterventionRequest) (*DisableInterventionResponse, error) {
	if e == nil || e.executor == nil || e.executor.intervention == nil {
		return nil, fmt.Errorf("INVALID_STATE: engine not ready")
	}
	ctl := e.executor.intervention
	ctl.mu.Lock()
	defer ctl.mu.Unlock()
	if ctl.session.Paused && ctl.session.Pending != nil {
		return nil, fmt.Errorf("INVALID_STATE: intervene continue or discard")
	}
	ctl.clearToDefaultLocked()
	return &DisableInterventionResponse{OK: true}, nil
}

func (e *Engine) SetBreakpoints(_ context.Context, req SetBreakpointsRequest) (*SetBreakpointsResponse, error) {
	if e == nil || e.executor == nil || e.executor.intervention == nil {
		return nil, fmt.Errorf("INVALID_STATE: engine not ready")
	}
	ctl := e.executor.intervention
	ctl.mu.Lock()
	defer ctl.mu.Unlock()
	if req.InstanceID != "" {
		if ctl.session.FocusInstanceID != "" && ctl.session.FocusInstanceID != req.InstanceID {
			return nil, fmt.Errorf("INVALID_ARGUMENT: instance_id does not match focus")
		}
		if ctl.session.FocusInstanceID == "" {
			ctl.session.FocusInstanceID = req.InstanceID
		}
	}
	bps := make(map[string]struct{}, len(req.ElementIDs))
	out := make([]string, 0, len(req.ElementIDs))
	seen := make(map[string]struct{}, len(req.ElementIDs))
	for _, id := range req.ElementIDs {
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		bps[id] = struct{}{}
		out = append(out, id)
	}
	sort.Strings(out)
	ctl.session.Breakpoints = bps
	return &SetBreakpointsResponse{OK: true, ElementIDs: out}, nil
}

func (e *Engine) GetInterventionState(_ context.Context, instanceID string) (*InterventionState, error) {
	if e == nil || e.executor == nil || e.executor.intervention == nil {
		return nil, fmt.Errorf("INVALID_STATE: engine not ready")
	}
	st := e.executor.intervention.snapshot()
	if instanceID != "" && st.FocusInstanceID != "" && st.FocusInstanceID != instanceID {
		return nil, fmt.Errorf("INVALID_ARGUMENT: instance_id does not match focus")
	}
	return &st, nil
}

func (e *Engine) Continue(ctx context.Context, req ContinueRequest) (*ContinueResponse, error) {
	return e.resumeIntervention(ctx, req.InstanceID, RunPolicyContinuous)
}

func (e *Engine) StepInto(ctx context.Context, req StepIntoRequest) (*StepIntoResponse, error) {
	resp, err := e.resumeIntervention(ctx, req.InstanceID, RunPolicyStep)
	if err != nil {
		return nil, err
	}
	return &StepIntoResponse{OK: resp.OK, Paused: resp.Paused, State: resp.State}, nil
}

// StepOver is Continue-equivalent in P0 (real stack-aware StepOver is P2).
func (e *Engine) StepOver(ctx context.Context, req StepOverRequest) (*StepOverResponse, error) {
	resp, err := e.Continue(ctx, ContinueRequest{InstanceID: req.InstanceID})
	if err != nil {
		return nil, err
	}
	return &StepOverResponse{OK: resp.OK, Paused: resp.Paused, State: resp.State}, nil
}

// SetVariablesRequest patches instance variables while barrier-paused (ledger COMMAND).
type SetVariablesRequest struct {
	InstanceID string
	Variables  map[string]any
}

type SetVariablesResponse struct {
	OK    bool
	State InterventionState
}

// SetVariables merges variables onto the focused paused instance via a PROCESS
// ACTIVATED EVENT (variable payload only). When paused at PendingDecide, updates
// Pending.TakenFlowIDs with a tentative exclusive choose for observers.
func (e *Engine) SetVariables(ctx context.Context, req SetVariablesRequest) (*SetVariablesResponse, error) {
	if e == nil || e.executor == nil || e.executor.intervention == nil {
		return nil, fmt.Errorf("INVALID_STATE: engine not ready")
	}
	instanceID := req.InstanceID
	if instanceID == "" {
		return nil, fmt.Errorf("INVALID_ARGUMENT: instance_id is required")
	}
	if !e.executor.intervention.isPausedFocus(instanceID) {
		return nil, fmt.Errorf("INVALID_STATE: instance not paused")
	}

	pv, err := projection.VariablesFromMap(req.Variables)
	if err != nil {
		return nil, fmt.Errorf("INVALID_ARGUMENT: variables: %w", err)
	}

	e.mu.Lock()
	inst := e.instances[instanceID]
	lock := e.instMu[instanceID]
	dep := (*deploy.Deployment)(nil)
	if inst != nil {
		dep = e.deployments[inst.DeploymentID]
	}
	e.mu.Unlock()
	if inst == nil || lock == nil {
		return nil, fmt.Errorf("NOT_FOUND: instance %q", instanceID)
	}

	lock.Lock()
	defer lock.Unlock()

	cmdID, err := NextID()
	if err != nil {
		return nil, err
	}
	processID := inst.ProcessID
	cmd := &eventv1.Event{
		Id:                cmdID,
		Timestamp:         nowMillis(),
		RecordType:        eventv1.Event_RECORD_TYPE_COMMAND,
		DeploymentId:      inst.DeploymentID,
		ProcessInstanceId: instanceID,
		ProcessVersion:    inst.Version,
		Element: &eventv1.Element{
			Intent: eventv1.Element_INTENT_ACTIVATED,
			Type:   eventv1.Element_TYPE_PROCESS,
			Id:     processID,
			Payload: &eventv1.Element_ProcessPayload{
				ProcessPayload: &eventv1.ProcessPayload{Variables: pv},
			},
		},
	}
	if _, err := e.log.Append(ctx, cmd); err != nil {
		return nil, err
	}
	if err := e.emitter(ctx, inst, cmdID)(&eventv1.Element{
		Intent: eventv1.Element_INTENT_ACTIVATED,
		Type:   eventv1.Element_TYPE_PROCESS,
		Id:     processID,
		Payload: &eventv1.Element_ProcessPayload{
			ProcessPayload: &eventv1.ProcessPayload{Variables: pv},
		},
	}); err != nil {
		return nil, err
	}

	ctl := e.executor.intervention
	ctl.mu.Lock()
	if dep != nil && ctl.session.Pending != nil && ctl.session.Pending.Kind == PendingDecide {
		if flowID, chooseErr := dep.ChooseExclusiveOutgoing(ctl.session.Pending.FromElementID, inst.Variables); chooseErr == nil {
			taken, next, _ := copyFlowTargets(dep, []string{flowID})
			ctl.session.Pending.TakenFlowIDs = taken
			ctl.session.Pending.NextElementIDs = next
			ctl.session.Pending.OutgoingFlowID = flowID
		}
	}
	st := ctl.snapshotLocked()
	ctl.mu.Unlock()
	return &SetVariablesResponse{OK: true, State: st}, nil
}

func (e *Engine) resumeIntervention(ctx context.Context, instanceID string, afterPolicy RunPolicy) (*ContinueResponse, error) {
	if e == nil || e.executor == nil || e.executor.intervention == nil {
		return nil, fmt.Errorf("INVALID_STATE: engine not ready")
	}
	ctl := e.executor.intervention
	ctl.mu.Lock()
	s := &ctl.session
	if !s.Enabled || !s.Paused || s.Pending == nil {
		ctl.mu.Unlock()
		return nil, fmt.Errorf("INVALID_STATE: not paused")
	}
	focus := s.FocusInstanceID
	if instanceID == "" {
		instanceID = focus
	}
	if instanceID != focus {
		ctl.mu.Unlock()
		return nil, fmt.Errorf("INVALID_ARGUMENT: instance_id does not match focus")
	}
	pending := *s.Pending
	pending.TakenFlowIDs = append([]string(nil), s.Pending.TakenFlowIDs...)
	pending.NextElementIDs = append([]string(nil), s.Pending.NextElementIDs...)
	pending.LinkCatchIDs = append([]string(nil), s.Pending.LinkCatchIDs...)
	sourceCmdID := s.SourceCmdID
	s.SuppressKey = barrierKey(&pending)
	s.RunPolicy = afterPolicy
	s.Paused = false
	s.PauseReason = ""
	s.PauseElementID = ""
	s.PauseTokenID = ""
	s.Pending = nil
	ctl.mu.Unlock()

	e.mu.Lock()
	inst := e.instances[instanceID]
	lock := e.instMu[instanceID]
	var dep *deploy.Deployment
	if inst != nil {
		dep = e.deployments[inst.DeploymentID]
	}
	e.mu.Unlock()
	if inst == nil || dep == nil || lock == nil {
		return nil, fmt.Errorf("NOT_FOUND: instance %q", instanceID)
	}

	lock.Lock()
	e.executor.activeSourceCmdID = sourceCmdID
	e.executor.activeInstanceID = instanceID
	pubs, err := e.executor.resumePending(ctx, dep, inst, &pending, e.emitter(ctx, inst, sourceCmdID))
	e.executor.activeSourceCmdID = ""
	e.executor.activeInstanceID = ""
	lock.Unlock()
	if err != nil {
		return nil, err
	}
	if err := e.afterUnlockedCommand(ctx, instanceID, pubs); err != nil {
		return nil, err
	}
	st := ctl.snapshot()
	return &ContinueResponse{OK: true, Paused: st.Paused, State: st}, nil
}

func copyFlowTargets(dep *deploy.Deployment, flows []string) ([]string, []string, error) {
	taken := append([]string(nil), flows...)
	next := make([]string, 0, len(flows))
	for _, flowID := range flows {
		sf, err := dep.SequenceFlow(flowID)
		if err != nil {
			return nil, nil, err
		}
		next = append(next, sf.TargetRef)
	}
	return taken, next, nil
}

func pendingLeave(dep *deploy.Deployment, tokenID, fromElementID, flowID string) (*PendingTransition, error) {
	taken, next, err := copyFlowTargets(dep, []string{flowID})
	if err != nil {
		return nil, err
	}
	return &PendingTransition{
		Kind:           PendingLeave,
		FromElementID:  fromElementID,
		TokenID:        tokenID,
		TakenFlowIDs:   taken,
		NextElementIDs: next,
		OutgoingFlowID: flowID,
	}, nil
}

func pendingFork(dep *deploy.Deployment, tokenID, fromElementID string, flows []string) (*PendingTransition, error) {
	taken, next, err := copyFlowTargets(dep, flows)
	if err != nil {
		return nil, err
	}
	return &PendingTransition{
		Kind:           PendingFork,
		FromElementID:  fromElementID,
		TokenID:        tokenID,
		TakenFlowIDs:   taken,
		NextElementIDs: next,
	}, nil
}
