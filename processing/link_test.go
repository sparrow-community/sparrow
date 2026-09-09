package processing_test

import (
	"context"
	"strings"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func TestLinkThrowToCatchThenComplete(t *testing.T) {
	xml := readTestdata(t, "m16_link_basic.bpmn")
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
	if !sawElementIntent(events, eventv1.Element_TYPE_INTERMEDIATE_THROW_EVENT, "LinkThrow_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected link throw COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT, "LinkCatch_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected link catch COMPLETED")
	}
}

func TestLinkThrowFansOutToMultipleCatches(t *testing.T) {
	xml := readTestdata(t, "m16_link_multi.bpmn")
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
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_a", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected End_a")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_b", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected End_b")
	}
}

func TestLinkOrphanThrowRejected(t *testing.T) {
	xml := readTestdata(t, "m16_link_orphan.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	_, err := eng.Deploy(context.Background(), xml)
	if err == nil {
		t.Fatal("expected orphan link throw deploy error")
	}
	if !strings.Contains(err.Error(), "no catch") {
		t.Fatalf("err=%v", err)
	}
}

func TestLinkRecoverThenComplete(t *testing.T) {
	xml := readTestdata(t, "m16_link_basic.bpmn")
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
	if elementID != "UserTask_1" {
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
