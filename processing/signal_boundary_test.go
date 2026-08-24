package processing_test

import (
	"context"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func TestSignalBoundaryPublishInterrupts(t *testing.T) {
	xml := readTestdata(t, "m3_signal_boundary.bpmn")
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
	elementID, tokenID := waitingAt(inst)
	if elementID != "UserTask_1" {
		t.Fatalf("tokens=%#v", inst.Tokens)
	}
	tok := inst.Tokens[tokenID]
	if tok.SignalName != "go.ahead" || tok.BoundaryID != "SignalBoundary_1" {
		t.Fatalf("armed signal token=%#v", tok)
	}

	n, err := eng.PublishSignal(ctx, processing.PublishSignalRequest{Name: "go.ahead"})
	if err != nil || n != 1 {
		t.Fatalf("publish n=%d err=%v", n, err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "UserTask_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected UserTask TERMINATED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "SignalBoundary_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected boundary COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_sig", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected signal end")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_ok", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("happy-path end must not complete")
	}
}

func TestSignalBoundaryCompleteCancelsWait(t *testing.T) {
	xml := readTestdata(t, "m3_signal_boundary.bpmn")
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
	elementID, tokenID := waitingAt(inst)
	if err := eng.Complete(ctx, instanceID, elementID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_ok", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected happy-path end")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_sig", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("signal end must not complete")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "SignalBoundary_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected boundary cancelled")
	}
}

func TestSignalBoundaryRecoverThenPublish(t *testing.T) {
	xml := readTestdata(t, "m3_signal_boundary.bpmn")
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
	_, tokenID := waitingAt(inst)
	tok := inst.Tokens[tokenID]
	if tok.SignalName != "go.ahead" || tok.BoundaryID != "SignalBoundary_1" {
		t.Fatalf("signal after recover token=%#v", tok)
	}
	n, err := eng2.PublishSignal(ctx, processing.PublishSignalRequest{Name: "go.ahead"})
	if err != nil || n != 1 {
		t.Fatalf("publish n=%d err=%v", n, err)
	}
	inst, _ = eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestSignalBoundaryNotBuffered(t *testing.T) {
	xml := readTestdata(t, "m3_signal_boundary.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	n, err := eng.PublishSignal(ctx, processing.PublishSignalRequest{Name: "go.ahead"})
	if err != nil || n != 0 {
		t.Fatalf("late signal n=%d err=%v", n, err)
	}
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	if inst.Status != projection.StatusActive {
		t.Fatalf("status=%s want active (signal is not buffered)", inst.Status)
	}
	elem, _ := waitingAt(inst)
	if elem != "UserTask_1" {
		t.Fatalf("still waiting at %q", elem)
	}
}

func TestNonInterruptingSignalBoundaryPublishThenComplete(t *testing.T) {
	xml := readTestdata(t, "m3_signal_boundary_non_interrupt.bpmn")
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
	n, err := eng.PublishSignal(ctx, processing.PublishSignalRequest{Name: "go.ahead"})
	if err != nil || n != 1 {
		t.Fatalf("publish n=%d err=%v", n, err)
	}
	inst, _ := eng.GetInstance(instanceID)
	if inst.Status != projection.StatusActive {
		t.Fatalf("status=%s want active", inst.Status)
	}
	if len(inst.Tokens) != 2 {
		t.Fatalf("want 2 tokens, got %d: %#v", len(inst.Tokens), inst.Tokens)
	}
	elementID, tokenID := waitingAt(inst)
	if elementID != "UserTask_1" {
		t.Fatalf("activity token should still wait, tokens=%#v", inst.Tokens)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "UserTask_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("UserTask must not be terminated")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_sig", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected signal end from spawned token")
	}
	if err := eng.Complete(ctx, instanceID, elementID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestTimerSignalBoundarySignalWins(t *testing.T) {
	xml := readTestdata(t, "m3_timer_signal_boundary.bpmn")
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
	_, tokenID := waitingAt(inst)
	tok := inst.Tokens[tokenID]
	if tok.BoundaryID != "TimerBoundary_1" {
		t.Fatalf("expected timer in BoundaryID, token=%#v", tok)
	}
	if tok.SignalName != "go.ahead" || tok.SignalBoundaryID != "SignalBoundary_1" {
		t.Fatalf("expected signal armed, token=%#v", tok)
	}

	n, err := eng.PublishSignal(ctx, processing.PublishSignalRequest{Name: "go.ahead"})
	if err != nil || n != 1 {
		t.Fatalf("publish n=%d err=%v", n, err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "SignalBoundary_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected signal boundary COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "TimerBoundary_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected timer boundary TERMINATED")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_timeout", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("timeout end must not complete")
	}
}

func TestTimerSignalBoundaryTimerWins(t *testing.T) {
	xml := readTestdata(t, "m3_timer_signal_boundary.bpmn")
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
	if err := eng.FireDue(ctx); err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "TimerBoundary_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected timer boundary COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "SignalBoundary_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected signal boundary TERMINATED")
	}
}

func TestSignalThrowWakesSignalBoundary(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	boundDep, err := eng.Deploy(ctx, readTestdata(t, "m3_signal_boundary.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	throwDep, err := eng.Deploy(ctx, readTestdata(t, "m3_signal_throw.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	boundID, err := eng.CreateInstance(ctx, boundDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := eng.CreateInstance(ctx, throwDep, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(boundID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed after throw", inst.Status)
	}
}

func TestSubProcessSignalBoundaryPublish(t *testing.T) {
	xml := readTestdata(t, "m3_subprocess_signal_boundary.bpmn")
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
	elemID, _ := waitingAt(inst)
	if elemID != "Sub_Task" {
		t.Fatalf("expected waiting at Sub_Task, got %q", elemID)
	}
	if len(inst.ScopeBoundaries) == 0 {
		t.Fatal("expected scope boundary to be armed")
	}
	n, err := eng.PublishSignal(ctx, processing.PublishSignalRequest{Name: "go.ahead"})
	if err != nil || n != 1 {
		t.Fatalf("publish n=%d err=%v", n, err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_SUB_PROCESS, "SubProcess_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected SubProcess TERMINATED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "SignalBoundary_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected boundary COMPLETED")
	}
}
