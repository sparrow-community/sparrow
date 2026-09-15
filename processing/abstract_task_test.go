package processing_test

import (
	"context"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func TestAbstractTaskWaitAndComplete(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m23_abstract_task.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	elementID, tokenID := waitingAt(inst)
	if elementID != "Task_1" {
		t.Fatalf("wait at %s tokens=%#v", elementID, inst.Tokens)
	}
	if err := eng.Complete(ctx, instanceID, elementID, tokenID, map[string]any{"done": true}); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s", inst.Status)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_TASK, "Task_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected TASK COMPLETED")
	}
}

func TestAbstractTaskRecoverThenComplete(t *testing.T) {
	dir := t.TempDir()
	ctx := context.Background()
	eng1, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	dep, err := eng1.Deploy(ctx, readTestdata(t, "m23_abstract_task.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	id, err := eng1.CreateInstance(ctx, dep, nil)
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
	inst, _ := eng2.GetInstance(id)
	el, tok := waitingAt(inst)
	if el != "Task_1" {
		t.Fatalf("wait=%s", el)
	}
	if err := eng2.Complete(ctx, id, el, tok, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng2.GetInstance(id)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%v", inst.Status)
	}
}
