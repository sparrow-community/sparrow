package processing_test

import (
	"context"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func TestComplexGatewayJoinActivationCondition(t *testing.T) {
	xml := readTestdata(t, "m34_complex_join.bpmn")
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
	if len(inst.Tokens) != 3 {
		t.Fatalf("want 3 tokens after fork, got %d", len(inst.Tokens))
	}
	completed := 0
	for _, tok := range inst.Tokens {
		if completed == 2 {
			break
		}
		if tok.Status != projection.TokenWaiting {
			continue
		}
		if err := eng.Complete(ctx, instanceID, tok.ElementID, tok.ID, nil); err != nil {
			t.Fatal(err)
		}
		completed++
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_COMPLEX_GATEWAY, "Gateway_join", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected complex join COMPLETED after 2 arrivals")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected End_1 after join fired")
	}
	// Third branch may still be waiting — 2-of-3 does not cancel upstream work.
	inst, _ = eng.GetInstance(instanceID)
	if !containsStr(allWaiting(inst), "Task_A") && !containsStr(allWaiting(inst), "Task_B") && !containsStr(allWaiting(inst), "Task_C") {
		// at least one of the three branch tasks should remain if only two completed
		waiting := allWaiting(inst)
		hasBranch := containsStr(waiting, "Task_A") || containsStr(waiting, "Task_B") || containsStr(waiting, "Task_C")
		if !hasBranch && inst.Status != projection.StatusCompleted {
			t.Fatalf("unexpected state waiting=%v status=%s", waiting, inst.Status)
		}
	}
}

func TestComplexGatewayJoinWaitAll(t *testing.T) {
	xml := readTestdata(t, "m34_complex_join_all.bpmn")
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
	var first *projection.Token
	for _, tok := range inst.Tokens {
		if tok.Status == projection.TokenWaiting {
			first = tok
			break
		}
	}
	if first == nil {
		t.Fatal("no waiting token")
	}
	if err := eng.Complete(ctx, instanceID, first.ElementID, first.ID, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status == projection.StatusCompleted {
		t.Fatal("must not complete after only one branch")
	}
	if !containsStr(allWaiting(inst), "Gateway_join") && !containsStr(allWaiting(inst), "Task_A") && !containsStr(allWaiting(inst), "Task_B") {
		// one at join, one still at task
	}
	for _, tok := range inst.Tokens {
		if tok.Status != projection.TokenWaiting {
			continue
		}
		if tok.ElementID == "Gateway_join" {
			continue
		}
		if err := eng.Complete(ctx, instanceID, tok.ElementID, tok.ID, nil); err != nil {
			t.Fatal(err)
		}
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed; waiting=%v", inst.Status, allWaiting(inst))
	}
}

func TestComplexGatewaySplitInclusiveStyle(t *testing.T) {
	xml := readTestdata(t, "m34_complex_split.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, map[string]any{"takeA": true, "takeB": true})
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	waiting := allWaiting(inst)
	if !containsStr(waiting, "Task_A") || !containsStr(waiting, "Task_B") {
		t.Fatalf("expected Task_A and Task_B, got %v", waiting)
	}
	if containsStr(waiting, "Task_Default") {
		t.Fatalf("default must not run when conditions match, got %v", waiting)
	}
	for _, tok := range inst.Tokens {
		if tok.Status != projection.TokenWaiting {
			continue
		}
		if err := eng.Complete(ctx, instanceID, tok.ElementID, tok.ID, nil); err != nil {
			t.Fatal(err)
		}
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s waiting=%v", inst.Status, allWaiting(inst))
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_A", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected End_A")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_B", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected End_B")
	}
}

func TestComplexGatewayJoinRecover(t *testing.T) {
	xml := readTestdata(t, "m34_complex_join.bpmn")
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
	inst, _ := eng1.GetInstance(instanceID)
	var first *projection.Token
	for _, tok := range inst.Tokens {
		if tok.Status == projection.TokenWaiting {
			first = tok
			break
		}
	}
	if err := eng1.Complete(ctx, instanceID, first.ElementID, first.ID, nil); err != nil {
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
		t.Fatal("missing instance after recover")
	}
	completed := 0
	for _, tok := range inst.Tokens {
		if tok.Status != projection.TokenWaiting {
			continue
		}
		if tok.ElementID == "Gateway_join" {
			continue
		}
		if err := eng2.Complete(ctx, instanceID, tok.ElementID, tok.ID, nil); err != nil {
			t.Fatal(err)
		}
		completed++
		if completed == 1 {
			break // second branch → arrived>=2
		}
	}
	events, err := eng2.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_COMPLEX_GATEWAY, "Gateway_join", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected complex join COMPLETED after recover + second arrival")
	}
}
