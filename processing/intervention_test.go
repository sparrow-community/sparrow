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

package processing_test

import (
	"context"
	"strings"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	"github.com/sparrow-community/sparrow/processing/deploy"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// T1: default-off Complete through exclusive gateway matches continuous trail.
func TestInterventionT1_DefaultOffExclusiveGateway(t *testing.T) {
	ctx := context.Background()
	xml := readTestdata(t, "m1_simple.bpmn")

	engOff := processing.NewEngine(eventlog.NewMemory())
	depOff, err := engOff.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	idOff, err := engOff.CreateInstance(ctx, depOff, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst := mustInstance(t, engOff, idOff)
	elem, tok := waitingAt(inst)
	if err := engOff.Complete(ctx, idOff, elem, tok, map[string]any{"approved": true}); err != nil {
		t.Fatal(err)
	}
	eventsOff, err := engOff.ListEvents(ctx, idOff)
	if err != nil {
		t.Fatal(err)
	}

	engOn := processing.NewEngine(eventlog.NewMemory())
	depOn, err := engOn.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	idOn, err := engOn.CreateInstance(ctx, depOn, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := engOn.EnableIntervention(ctx, processing.EnableInterventionRequest{
		InstanceID: idOn,
		Policy:     processing.RunPolicyContinuous,
	}); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, engOn, idOn)
	elem, tok = waitingAt(inst)
	if err := engOn.Complete(ctx, idOn, elem, tok, map[string]any{"approved": true}); err != nil {
		t.Fatal(err)
	}
	eventsOn, err := engOn.ListEvents(ctx, idOn)
	if err != nil {
		t.Fatal(err)
	}

	if trailFingerprint(eventsOff) != trailFingerprint(eventsOn) {
		t.Fatalf("continuous intervention trail diverged from default-off\noff=%s\non=%s",
			trailFingerprint(eventsOff), trailFingerprint(eventsOn))
	}
	st, err := engOn.GetInterventionState(ctx, idOn)
	if err != nil {
		t.Fatal(err)
	}
	if st.Paused {
		t.Fatal("continuous policy must not barrier-pause")
	}
}

// T2: Enable + Step; drive to exclusive gateway barrier.
func TestInterventionT2_StepPauseOnExclusiveGateway(t *testing.T) {
	ctx := context.Background()
	eng := processing.NewEngine(eventlog.NewMemory())
	dep, err := eng.Deploy(ctx, readTestdata(t, "m1_simple.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := eng.EnableIntervention(ctx, processing.EnableInterventionRequest{
		InstanceID: instanceID,
		Policy:     processing.RunPolicyStep,
	}); err != nil {
		t.Fatal(err)
	}
	inst := mustInstance(t, eng, instanceID)
	elem, tok := waitingAt(inst)
	if elem != "UserTask_1" {
		t.Fatalf("want UserTask_1, got %s", elem)
	}
	if err := eng.Complete(ctx, instanceID, elem, tok, map[string]any{"approved": true}); err != nil {
		t.Fatal(err)
	}
	st := mustStepToElement(t, eng, instanceID, "Gateway_1")
	if !st.Paused || st.PauseElementID != "Gateway_1" {
		t.Fatalf("want paused at Gateway_1, got %#v", st)
	}
	if st.Pending == nil || st.Pending.Kind != processing.PendingLeave {
		t.Fatalf("want pending leave, got %#v", st.Pending)
	}
	if len(st.Pending.TakenFlowIDs) != 1 {
		t.Fatalf("want one taken flow, got %#v", st.Pending.TakenFlowIDs)
	}

	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_EXCLUSIVE_GATEWAY, "Gateway_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected Gateway_1 COMPLETED")
	}
	if sawSequenceFlowTakenFrom(events, "Gateway_1") {
		t.Fatal("must not take gateway outgoing before Continue")
	}
	inst = mustInstance(t, eng, instanceID)
	tokID := tokenAnyAtElement(inst, "Gateway_1")
	if tokID == "" {
		t.Fatalf("token should remain on gateway, tokens=%#v", inst.Tokens)
	}
	if inst.Tokens[tokID].Status != projection.TokenActive {
		t.Fatalf("gateway token status=%s want active", inst.Tokens[tokID].Status)
	}
}

// T3: Continue from T2 leaves gateway under Continuous.
func TestInterventionT3_ContinueFromGateway(t *testing.T) {
	ctx := context.Background()
	eng, instanceID := setupPausedAtGateway(t)
	resp, err := eng.Continue(ctx, processing.ContinueRequest{InstanceID: instanceID})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Paused {
		t.Fatalf("continuous Continue should not re-pause, state=%#v", resp.State)
	}
	inst := mustInstance(t, eng, instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawSequenceFlowTakenFrom(events, "Gateway_1") {
		t.Fatal("expected SEQUENCE_FLOW_TAKEN after Continue")
	}
}

// T4: StepInto from gateway pauses again or stops at wait/end.
func TestInterventionT4_StepIntoFromGateway(t *testing.T) {
	ctx := context.Background()
	eng, instanceID := setupPausedAtGateway(t)
	resp, err := eng.StepInto(ctx, processing.StepIntoRequest{InstanceID: instanceID})
	if err != nil {
		t.Fatal(err)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawSequenceFlowTakenFrom(events, "Gateway_1") {
		t.Fatal("StepInto must take gateway outgoing")
	}
	inst := mustInstance(t, eng, instanceID)
	// Next element is EndEvent (instant) — may pause before process complete or finish.
	if resp.Paused {
		if resp.State.PauseElementID == "Gateway_1" {
			t.Fatal("must not remain paused on the same gateway")
		}
	} else if inst.Status != projection.StatusCompleted {
		t.Fatalf("unpaused but status=%s tokens=%#v", inst.Status, inst.Tokens)
	}
}

// T5: Parallel fork pauses once before fork; Continue enters both branches.
func TestInterventionT5_ParallelForkPauseOnce(t *testing.T) {
	ctx := context.Background()
	eng := processing.NewEngine(eventlog.NewMemory())
	dep, err := eng.Deploy(ctx, readTestdata(t, "m3_parallel_fork_join.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	// Bootstrap instance so EnableIntervention has a focus; CreateInstance rebinds
	// focus when session is enabled and not paused.
	bootstrap, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := eng.EnableIntervention(ctx, processing.EnableInterventionRequest{
		InstanceID: bootstrap,
		Policy:     processing.RunPolicyStep,
	}); err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	st := mustStepToPendingKind(t, eng, instanceID, processing.PendingFork)
	if len(st.Pending.TakenFlowIDs) < 2 {
		t.Fatalf("want >=2 fork flows, got %#v", st.Pending.TakenFlowIDs)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if sawSequenceFlowTakenFrom(events, "Gateway_fork") {
		t.Fatal("must not take fork outgoings before Continue")
	}

	resp, err := eng.Continue(ctx, processing.ContinueRequest{InstanceID: instanceID})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Paused {
		t.Fatalf("Continue from fork under continuous should reach waits, state=%#v", resp.State)
	}
	inst := mustInstance(t, eng, instanceID)
	waiting := allWaiting(inst)
	if !containsStr(waiting, "UserTask_a") || !containsStr(waiting, "UserTask_b") {
		t.Fatalf("want both branches waiting, got %v tokens=%#v", waiting, inst.Tokens)
	}
}

// T6: pause → Continue → Recover matches non-intervention path.
func TestInterventionT6_RecoverAfterContinue(t *testing.T) {
	ctx := context.Background()
	xml := readTestdata(t, "m1_simple.bpmn")
	memLog := eventlog.NewMemory()
	store := deploy.NewMemoryStore()
	eng1, err := processing.Recover(ctx, memLog, store, nil)
	if err != nil {
		t.Fatal(err)
	}
	dep, err := eng1.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng1.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := eng1.EnableIntervention(ctx, processing.EnableInterventionRequest{
		InstanceID: instanceID,
		Policy:     processing.RunPolicyStep,
	}); err != nil {
		t.Fatal(err)
	}
	inst := mustInstance(t, eng1, instanceID)
	elem, tok := waitingAt(inst)
	if err := eng1.Complete(ctx, instanceID, elem, tok, map[string]any{"approved": true}); err != nil {
		t.Fatal(err)
	}
	_ = mustStepToElement(t, eng1, instanceID, "Gateway_1")
	if _, err := eng1.Continue(ctx, processing.ContinueRequest{InstanceID: instanceID}); err != nil {
		t.Fatal(err)
	}
	before, err := eng1.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	instBefore := mustInstance(t, eng1, instanceID)

	eng2, err := processing.Recover(ctx, memLog, store, nil)
	if err != nil {
		t.Fatal(err)
	}
	st, err := eng2.GetInterventionState(ctx, "")
	if err != nil {
		t.Fatal(err)
	}
	if st.Enabled || st.Paused {
		t.Fatalf("Recover must leave intervention off, got %#v", st)
	}
	after, err := eng2.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if trailFingerprint(before) != trailFingerprint(after) {
		t.Fatalf("trail mismatch after Recover\nbefore=%s\nafter=%s", trailFingerprint(before), trailFingerprint(after))
	}
	instAfter := mustInstance(t, eng2, instanceID)
	if instBefore.Status != instAfter.Status {
		t.Fatalf("status before=%s after=%s", instBefore.Status, instAfter.Status)
	}
}

// T7: Call Activity enter is natural wait; child starts; Step can pause before parent leave.
func TestInterventionT7_CallActivityNoBarrierOnEnter(t *testing.T) {
	ctx := context.Background()
	eng := processing.NewEngine(eventlog.NewMemory())
	dep, err := eng.Deploy(ctx, readTestdataCall(t, "m4_call_activity.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	callerID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := eng.EnableIntervention(ctx, processing.EnableInterventionRequest{
		InstanceID: callerID,
		Policy:     processing.RunPolicyStep,
	}); err != nil {
		t.Fatal(err)
	}
	caller := mustInstance(t, eng, callerID)
	hostElem, hostTok := waitingAt(caller)
	if hostElem != "CallActivity_1" {
		t.Fatalf("expected CallActivity_1 wait, got %s", hostElem)
	}
	st, err := eng.GetInterventionState(ctx, callerID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Paused {
		t.Fatal("Call Activity enter must not barrier-pause")
	}
	childID := caller.Tokens[hostTok].CalledProcessInstanceID
	if childID == "" {
		t.Fatal("child must have started (publication flushed on wait)")
	}
	child := mustInstance(t, eng, childID)
	ce, ct := waitingAt(child)
	if err := eng.Complete(ctx, childID, ce, ct, nil); err != nil {
		t.Fatal(err)
	}
	// Parent resumes Complete → may barrier before leave under Step.
	st, err = eng.GetInterventionState(ctx, callerID)
	if err != nil {
		t.Fatal(err)
	}
	if !st.Paused {
		// Child complete resumes parent via publication; focus is caller — Step may pause on leave.
		caller = mustInstance(t, eng, callerID)
		if caller.Status != projection.StatusCompleted {
			t.Fatalf("expected pause before parent leave or completed, status=%s paused=%v tokens=%#v",
				caller.Status, st.Paused, caller.Tokens)
		}
	} else if st.PauseElementID != "CallActivity_1" {
		t.Fatalf("want pause at CallActivity_1 leave, got %#v", st)
	}
}

// T8: natural wait unchanged with Enabled+Continuous.
func TestInterventionT8_NaturalWaitUnchanged(t *testing.T) {
	ctx := context.Background()
	eng := processing.NewEngine(eventlog.NewMemory())
	dep, err := eng.Deploy(ctx, readTestdata(t, "m1_simple.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := eng.EnableIntervention(ctx, processing.EnableInterventionRequest{
		InstanceID: instanceID,
		Policy:     processing.RunPolicyContinuous,
	}); err != nil {
		t.Fatal(err)
	}
	st, err := eng.GetInterventionState(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Paused {
		t.Fatal("natural wait must not set barrier Paused")
	}
	inst := mustInstance(t, eng, instanceID)
	elem, tok := waitingAt(inst)
	if elem != "UserTask_1" {
		t.Fatalf("want UserTask_1, got %s", elem)
	}
	if err := eng.Complete(ctx, instanceID, elem, tok, map[string]any{"approved": true}); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s", inst.Status)
	}
}

// T9: Disable while paused rejects.
func TestInterventionT9_DisableWhilePausedRejects(t *testing.T) {
	ctx := context.Background()
	eng, instanceID := setupPausedAtGateway(t)
	_, err := eng.DisableIntervention(ctx, processing.DisableInterventionRequest{})
	if err == nil || !strings.Contains(err.Error(), "INVALID_STATE") {
		t.Fatalf("want INVALID_STATE, got %v", err)
	}
	st, err := eng.GetInterventionState(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !st.Paused {
		t.Fatal("must still be paused after rejected Disable")
	}
}

// T10: SetBreakpoints store round-trip (hit is P1).
func TestInterventionT10_SetBreakpointsStore(t *testing.T) {
	ctx := context.Background()
	eng := processing.NewEngine(eventlog.NewMemory())
	dep, err := eng.Deploy(ctx, readTestdata(t, "m1_simple.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := eng.EnableIntervention(ctx, processing.EnableInterventionRequest{
		InstanceID: instanceID,
		Policy:     processing.RunPolicyContinuous,
	}); err != nil {
		t.Fatal(err)
	}
	resp, err := eng.SetBreakpoints(ctx, processing.SetBreakpointsRequest{
		InstanceID: instanceID,
		ElementIDs: []string{"Gateway_1", "UserTask_1", "Gateway_1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.ElementIDs) != 2 {
		t.Fatalf("want 2 unique ids, got %v", resp.ElementIDs)
	}
	st, err := eng.GetInterventionState(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !containsStr(st.Breakpoints, "Gateway_1") || !containsStr(st.Breakpoints, "UserTask_1") {
		t.Fatalf("breakpoints=%v", st.Breakpoints)
	}
}

func TestIntervention_RejectCompleteWhilePaused(t *testing.T) {
	ctx := context.Background()
	eng, instanceID := setupPausedAtGateway(t)
	err := eng.Complete(ctx, instanceID, "UserTask_1", "nope", nil)
	if err == nil || !strings.Contains(err.Error(), "instance paused") {
		t.Fatalf("want instance paused, got %v", err)
	}
}

func setupPausedAtGateway(t *testing.T) (*processing.Engine, string) {
	t.Helper()
	ctx := context.Background()
	eng := processing.NewEngine(eventlog.NewMemory())
	dep, err := eng.Deploy(ctx, readTestdata(t, "m1_simple.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := eng.EnableIntervention(ctx, processing.EnableInterventionRequest{
		InstanceID: instanceID,
		Policy:     processing.RunPolicyStep,
	}); err != nil {
		t.Fatal(err)
	}
	inst := mustInstance(t, eng, instanceID)
	elem, tok := waitingAt(inst)
	if err := eng.Complete(ctx, instanceID, elem, tok, map[string]any{"approved": true}); err != nil {
		t.Fatal(err)
	}
	_ = mustStepToElement(t, eng, instanceID, "Gateway_1")
	return eng, instanceID
}

func mustStepToElement(t *testing.T, eng *processing.Engine, instanceID, elementID string) *processing.InterventionState {
	t.Helper()
	ctx := context.Background()
	for i := 0; i < 16; i++ {
		st, err := eng.GetInterventionState(ctx, instanceID)
		if err != nil {
			t.Fatal(err)
		}
		if st.Paused && st.PauseElementID == elementID {
			return st
		}
		if !st.Paused {
			t.Fatalf("not paused while seeking %s (state=%#v)", elementID, st)
		}
		if _, err := eng.StepInto(ctx, processing.StepIntoRequest{InstanceID: instanceID}); err != nil {
			t.Fatalf("StepInto seeking %s: %v (state=%#v)", elementID, err, st)
		}
	}
	t.Fatalf("did not reach pause at %s", elementID)
	return nil
}

func mustStepToPendingKind(t *testing.T, eng *processing.Engine, instanceID string, kind processing.PendingKind) *processing.InterventionState {
	t.Helper()
	ctx := context.Background()
	for i := 0; i < 16; i++ {
		st, err := eng.GetInterventionState(ctx, instanceID)
		if err != nil {
			t.Fatal(err)
		}
		if st.Paused && st.Pending != nil && st.Pending.Kind == kind {
			return st
		}
		if !st.Paused {
			t.Fatalf("not paused while seeking pending %s (state=%#v)", kind, st)
		}
		if _, err := eng.StepInto(ctx, processing.StepIntoRequest{InstanceID: instanceID}); err != nil {
			t.Fatalf("StepInto seeking %s: %v", kind, err)
		}
	}
	t.Fatalf("did not reach pending kind %s", kind)
	return nil
}

func sawSequenceFlowTakenFrom(events []*eventv1.Event, sourceElementID string) bool {
	for _, ev := range events {
		if ev.GetRecordType() != eventv1.Event_RECORD_TYPE_EVENT {
			continue
		}
		el := ev.GetElement()
		if el == nil || el.GetType() != eventv1.Element_TYPE_SEQUENCE_FLOW {
			continue
		}
		if el.GetIntent() != eventv1.Element_INTENT_SEQUENCE_FLOW_TAKEN {
			continue
		}
		if p := el.GetSequenceFlowPayload(); p != nil && p.GetSourceId() == sourceElementID {
			return true
		}
	}
	return false
}

func trailFingerprint(events []*eventv1.Event) string {
	var b strings.Builder
	for _, ev := range events {
		if ev.GetRecordType() != eventv1.Event_RECORD_TYPE_EVENT {
			continue
		}
		el := ev.GetElement()
		if el == nil {
			continue
		}
		b.WriteString(el.GetType().String())
		b.WriteByte(':')
		b.WriteString(el.GetId())
		b.WriteByte(':')
		b.WriteString(el.GetIntent().String())
		b.WriteByte(';')
	}
	return b.String()
}
