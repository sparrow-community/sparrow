package processing_test

import (
	"context"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func TestSubProcessDeploy(t *testing.T) {
	xml := readTestdata(t, "m2_subprocess.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	_, err := eng.Deploy(context.Background(), xml)
	if err != nil {
		t.Fatal("should accept subprocess:", err)
	}
}

func TestSubProcessEnterAndComplete(t *testing.T) {
	xml := readTestdata(t, "m2_subprocess.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	elemID, tokenID := waitingAt(inst)
	if elemID != "Sub_UserTask_1" {
		t.Fatalf("expected waiting at Sub_UserTask_1, got %q", elemID)
	}

	if err := eng.Complete(ctx, instanceID, elemID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_SUB_PROCESS, "SubProcess_1", eventv1.Element_INTENT_ACTIVATED) {
		t.Fatal("expected SubProcess ACTIVATED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_SUB_PROCESS, "SubProcess_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected SubProcess COMPLETED")
	}
}

func TestSubProcessInternalParallel(t *testing.T) {
	xml := readTestdata(t, "m2_subprocess_parallel.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	// Should have 2 waiting tokens (TaskA and TaskB)
	var waitingTokens []struct{ elem, token string }
	for tid, tok := range inst.Tokens {
		if tok.Status == projection.TokenWaiting {
			waitingTokens = append(waitingTokens, struct{ elem, token string }{tok.ElementID, tid})
		}
	}
	if len(waitingTokens) != 2 {
		t.Fatalf("expected 2 waiting tokens, got %d: %+v", len(waitingTokens), waitingTokens)
	}

	// Complete both
	for _, w := range waitingTokens {
		if err := eng.Complete(ctx, instanceID, w.elem, w.token, nil); err != nil {
			t.Fatalf("Complete %s: %v", w.elem, err)
		}
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		for tid, tok := range inst.Tokens {
			t.Logf("TOKEN: id=%s elem=%s status=%s", tid, tok.ElementID, tok.Status)
		}
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestSubProcessTimerBoundaryFires(t *testing.T) {
	xml := readTestdata(t, "m2_subprocess_timer_boundary.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	// Token should be waiting at Sub_Task inside the subprocess
	inst, _ := eng.GetInstance(instanceID)
	elemID, _ := waitingAt(inst)
	if elemID != "Sub_Task" {
		t.Fatalf("expected waiting at Sub_Task, got %q", elemID)
	}
	// Scope boundary should be armed
	if len(inst.ScopeBoundaries) == 0 {
		t.Fatal("expected scope boundary to be armed")
	}

	// Fire due — timer is PT0S so already expired
	if err := eng.FireDue(ctx); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		for tid, tok := range inst.Tokens {
			t.Logf("TOKEN: id=%s elem=%s status=%s", tid, tok.ElementID, tok.Status)
		}
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_SUB_PROCESS, "SubProcess_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected SubProcess TERMINATED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "TimerBoundary_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected boundary COMPLETED")
	}
}

func TestSubProcessTimerBoundaryActivityCompletesFirst(t *testing.T) {
	xml := readTestdata(t, "m2_subprocess_timer_boundary.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	elemID, tokenID := waitingAt(inst)
	// Complete inner task first
	if err := eng.Complete(ctx, instanceID, elemID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	// Scope boundary should have been removed
	if len(inst.ScopeBoundaries) != 0 {
		t.Fatalf("expected scope boundaries cleared, got %d", len(inst.ScopeBoundaries))
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_SUB_PROCESS, "SubProcess_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected SubProcess COMPLETED")
	}
}

func TestSubProcessRecoverThenComplete(t *testing.T) {
	xml := readTestdata(t, "m2_subprocess.bpmn")
	dir := t.TempDir()
	ctx := context.Background()

	eng1, err := processing.Open(ctx, dir)
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
	if err := eng1.Close(); err != nil {
		t.Fatal(err)
	}

	eng2, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	defer eng2.Close()

	inst, ok := eng2.GetInstance(instanceID)
	if !ok {
		t.Fatal("instance missing after recover")
	}
	elemID, tokenID := waitingAt(inst)
	if elemID != "Sub_UserTask_1" {
		t.Fatalf("expected waiting at Sub_UserTask_1 after recover, got %q", elemID)
	}
	if err := eng2.Complete(ctx, instanceID, elemID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestSubProcessTimerBoundaryRecoverThenFire(t *testing.T) {
	xml := readTestdata(t, "m2_subprocess_timer_boundary.bpmn")
	dir := t.TempDir()
	ctx := context.Background()

	eng1, err := processing.Open(ctx, dir)
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
	if err := eng1.Close(); err != nil {
		t.Fatal(err)
	}

	eng2, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	defer eng2.Close()

	inst, ok := eng2.GetInstance(instanceID)
	if !ok {
		t.Fatal("instance missing after recover")
	}
	if len(inst.ScopeBoundaries) == 0 {
		t.Fatal("expected scope boundary to be armed after recover")
	}
	if err := eng2.FireDue(ctx); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}
