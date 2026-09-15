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

func TestConditionalCatchWaitsThenEvaluate(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m21_conditional_catch.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	id, err := eng.CreateInstance(ctx, dep, map[string]any{"approved": false})
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(id)
	if waitingElement(inst) != "CondCatch_1" {
		t.Fatalf("wait=%q tokens=%#v", waitingElement(inst), inst.Tokens)
	}
	n, err := eng.EvaluateConditions(ctx, processing.EvaluateConditionsRequest{
		ProcessInstanceID: id,
		Variables:         map[string]any{"approved": true},
	})
	if err != nil || n != 1 {
		t.Fatalf("EvaluateConditions n=%d err=%v", n, err)
	}
	inst, _ = eng.GetInstance(id)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%v", inst.Status)
	}
}

func TestConditionalCatchTrueAtEnter(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m21_conditional_catch.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	id, err := eng.CreateInstance(ctx, dep, map[string]any{"approved": true})
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(id)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%v", inst.Status)
	}
}

func TestConditionalBoundaryInterrupting(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m21_conditional_boundary.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	id, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(id)
	el, tokID := waitingAt(inst)
	if el != "UserTask_1" {
		t.Fatalf("wait=%q", el)
	}
	tok := inst.Tokens[tokID]
	found := false
	for _, w := range tok.BoundaryWaits {
		if w.Kind == "conditional" && w.BoundaryID == "CondBoundary_1" {
			found = true
		}
	}
	if !found {
		t.Fatalf("conditional boundary not armed: %#v", tok.BoundaryWaits)
	}
	n, err := eng.EvaluateConditions(ctx, processing.EvaluateConditionsRequest{
		ProcessInstanceID: id,
		Variables:         map[string]any{"approved": true},
	})
	if err != nil || n != 1 {
		t.Fatalf("EvaluateConditions n=%d err=%v", n, err)
	}
	inst, _ = eng.GetInstance(id)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%v", inst.Status)
	}
	events, _ := eng.ListEvents(ctx, id)
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "CondBoundary_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected CondBoundary_1 COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "UserTask_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected UserTask_1 TERMINATED")
	}
}

func TestConditionalBoundaryNonInterrupting(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m21_conditional_boundary_ni.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	id, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	n, err := eng.EvaluateConditions(ctx, processing.EvaluateConditionsRequest{
		ProcessInstanceID: id,
		Variables:         map[string]any{"approved": true},
	})
	if err != nil || n != 1 {
		t.Fatalf("EvaluateConditions n=%d err=%v", n, err)
	}
	inst, _ := eng.GetInstance(id)
	if waitingElement(inst) != "UserTask_1" {
		t.Fatalf("host should still wait, wait=%q status=%v", waitingElement(inst), inst.Status)
	}
	el, tok := waitingAt(inst)
	if err := eng.Complete(ctx, id, el, tok, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(id)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%v", inst.Status)
	}
}

func TestDuplicateConditionalBoundaryRejected(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	_, err := eng.Deploy(context.Background(), readTestdata(t, "m21_duplicate_conditional_boundary.bpmn"))
	if err == nil || !strings.Contains(err.Error(), "UNSUPPORTED_ELEMENT") {
		t.Fatalf("Deploy err=%v", err)
	}
}

func TestConditionalCatchRecoverThenEvaluate(t *testing.T) {
	dir := t.TempDir()
	ctx := context.Background()
	eng1, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	dep, err := eng1.Deploy(ctx, readTestdata(t, "m21_conditional_catch.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	id, err := eng1.CreateInstance(ctx, dep, map[string]any{"approved": false})
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
	n, err := eng2.EvaluateConditions(ctx, processing.EvaluateConditionsRequest{
		ProcessInstanceID: id,
		Variables:         map[string]any{"approved": true},
	})
	if err != nil || n != 1 {
		t.Fatalf("EvaluateConditions n=%d err=%v", n, err)
	}
	inst, _ := eng2.GetInstance(id)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%v", inst.Status)
	}
}
