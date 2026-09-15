package processing_test

import (
	"context"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
)

func TestActivityParallelFanout(t *testing.T) {
	xml := readTestdata(t, "m32_activity_parallel_fanout.bpmn")
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
	elem, tok := waitingAt(inst)
	if elem != "UserTask_split" {
		t.Fatalf("expected UserTask_split, got %s", elem)
	}
	if err := eng.Complete(ctx, instanceID, elem, tok, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if len(inst.Tokens) != 2 {
		t.Fatalf("want 2 tokens after fan-out, got %d: %#v", len(inst.Tokens), inst.Tokens)
	}
	var tasks []string
	for _, tkn := range inst.Tokens {
		if tkn.Status != projection.TokenWaiting {
			t.Fatalf("token not waiting: %#v", tkn)
		}
		tasks = append(tasks, tkn.ElementID)
	}
	if !containsStr(tasks, "UserTask_a") || !containsStr(tasks, "UserTask_b") {
		t.Fatalf("tasks=%v", tasks)
	}
	for _, tkn := range inst.Tokens {
		if err := eng.Complete(ctx, instanceID, tkn.ElementID, tkn.ID, nil); err != nil {
			t.Fatalf("Complete %s: %v", tkn.ElementID, err)
		}
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestStartEventParallelFanout(t *testing.T) {
	xml := readTestdata(t, "m32_start_parallel_fanout.bpmn")
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
	if len(inst.Tokens) != 2 {
		t.Fatalf("want 2 tokens after start fan-out, got %d: %#v", len(inst.Tokens), inst.Tokens)
	}
	var tasks []string
	for _, tkn := range inst.Tokens {
		if tkn.Status != projection.TokenWaiting {
			t.Fatalf("token not waiting: %#v", tkn)
		}
		tasks = append(tasks, tkn.ElementID)
	}
	if !containsStr(tasks, "UserTask_a") || !containsStr(tasks, "UserTask_b") {
		t.Fatalf("tasks=%v", tasks)
	}
	for _, tkn := range inst.Tokens {
		if err := eng.Complete(ctx, instanceID, tkn.ElementID, tkn.ID, nil); err != nil {
			t.Fatalf("Complete %s: %v", tkn.ElementID, err)
		}
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestActivityParallelFanoutRecoverThenComplete(t *testing.T) {
	xml := readTestdata(t, "m32_activity_parallel_fanout.bpmn")
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
	elem, tok := waitingAt(inst)
	if elem != "UserTask_split" {
		t.Fatalf("expected UserTask_split after recover, got %s", elem)
	}
	if err := eng2.Complete(ctx, instanceID, elem, tok, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng2.GetInstance(instanceID)
	if len(inst.Tokens) != 2 {
		t.Fatalf("want 2 tokens after fan-out, got %d", len(inst.Tokens))
	}
	for _, tkn := range inst.Tokens {
		if err := eng2.Complete(ctx, instanceID, tkn.ElementID, tkn.ID, nil); err != nil {
			t.Fatalf("Complete %s: %v", tkn.ElementID, err)
		}
	}
	inst, _ = eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}
