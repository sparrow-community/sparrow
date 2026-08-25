package processing_test

import (
	"context"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func TestErrorBoundaryThrowInterrupts(t *testing.T) {
	xml := readTestdata(t, "m4_error_boundary.bpmn")
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

	if err := eng.ThrowError(ctx, instanceID, elementID, tokenID, "BUSINESS_ERROR"); err != nil {
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
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "UserTask_1", eventv1.Element_INTENT_ERROR_THROWN) {
		t.Fatal("expected ERROR_THROWN on user task")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "UserTask_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected UserTask TERMINATED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "ErrorBoundary_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected error boundary COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_err", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected error path end")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_ok", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("happy-path end must not complete")
	}
}

func TestErrorEndInSubProcessCaughtByBoundary(t *testing.T) {
	xml := readTestdata(t, "m4_error_subprocess.bpmn")
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
	if elementID != "UserTask_sp" {
		t.Fatalf("tokens=%#v", inst.Tokens)
	}

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
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_sp_err", eventv1.Element_INTENT_ERROR_THROWN) {
		t.Fatal("expected ERROR_THROWN on error end")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_SUB_PROCESS, "SubProcess_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected SubProcess TERMINATED")
	}
	if tid := elementTokenID(events, eventv1.Element_TYPE_SUB_PROCESS, "SubProcess_1", eventv1.Element_INTENT_TERMINATED); tid == "" {
		t.Fatal("scope error boundary SubProcess TERMINATED must keep token_id for outgoing")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "ErrorBoundary_sp", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected scope error boundary COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_err", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected error handler end")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_ok", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("happy-path end must not complete")
	}
}

func TestUnhandledErrorTerminatesInstance(t *testing.T) {
	xml := readTestdata(t, "m4_error_unhandled.bpmn")
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

	if err := eng.ThrowError(ctx, instanceID, elementID, tokenID, "FATAL"); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusTerminated {
		t.Fatalf("status=%s want terminated", inst.Status)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_PROCESS, "Process_error_unhandled", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected process TERMINATED")
	}
}

func TestErrorBoundaryCodeMismatchUnhandled(t *testing.T) {
	xml := readTestdata(t, "m4_error_boundary.bpmn")
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

	if err := eng.ThrowError(ctx, instanceID, elementID, tokenID, "OTHER_ERROR"); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusTerminated {
		t.Fatalf("status=%s want terminated for unmatched error code", inst.Status)
	}
}
