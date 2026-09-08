package processing_test

import (
	"context"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func TestEscalationUncaughtIntermediateThrowCompletes(t *testing.T) {
	xml := readTestdata(t, "m15_escalation_uncaught.bpmn")
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
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_INTERMEDIATE_THROW_EVENT, "Throw_esc", eventv1.Element_INTENT_ESCALATION_THROWN) {
		t.Fatal("expected ESCALATION_THROWN")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_ok", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected End_ok")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_PROCESS, "Process_escalation_uncaught", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("uncaught escalation must not terminate process")
	}
}

func TestEscalationBoundaryInterruptsSubProcess(t *testing.T) {
	xml := readTestdata(t, "m15_escalation_boundary.bpmn")
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
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "EscalationBoundary_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected escalation boundary COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_SUB_PROCESS, "SubProcess_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected SubProcess TERMINATED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_esc", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected End_esc")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_ok", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("happy-path end must not complete")
	}
}

func TestEscalationEventSubProcessInterrupting(t *testing.T) {
	xml := readTestdata(t, "m15_escalation_esp.bpmn")
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
	if elementID != "UserTask_EventSubProcess_1" {
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
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_esc", eventv1.Element_INTENT_ESCALATION_THROWN) {
		t.Fatal("expected escalation end ESCALATION_THROWN")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "UserTask_EventSubProcess_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected ESP handler completed")
	}
}

func TestEscalationEventSubProcessRecoverThenComplete(t *testing.T) {
	xml := readTestdata(t, "m15_escalation_esp.bpmn")
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
	elementID, tokenID := waitingAt(inst)
	if elementID != "UserTask_EventSubProcess_1" {
		t.Fatalf("after recover tokens=%#v", inst.Tokens)
	}
	if err := eng2.Complete(ctx, instanceID, elementID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}
