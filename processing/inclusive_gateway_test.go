package processing_test

import (
	"context"
	"sort"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
)

func TestInclusiveGatewayDeploy(t *testing.T) {
	xml := readTestdata(t, "m2_inclusive_gateway.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	_, err := eng.Deploy(context.Background(), xml)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInclusiveGatewaySplitAll(t *testing.T) {
	xml := readTestdata(t, "m2_inclusive_gateway.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	// Both conditions true → A and C taken (default not taken)
	vars := map[string]any{"path_a": "true", "path_c": "true"}
	instanceID, err := eng.CreateInstance(ctx, dep, vars)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	waiting := allWaiting(inst)
	sort.Strings(waiting)
	if len(waiting) != 2 {
		t.Fatalf("expected 2 waiting tokens, got %v", waiting)
	}
	expected := []string{"Task_A", "Task_C"}
	// Complete all 3
	for _, elem := range expected {
		inst, _ = eng.GetInstance(instanceID)
		tokenID := tokenAtElement(inst, elem)
		if tokenID == "" {
			t.Fatalf("no token at %s", elem)
		}
		if err := eng.Complete(ctx, instanceID, elem, tokenID, nil); err != nil {
			t.Fatalf("complete %s: %v", elem, err)
		}
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestInclusiveGatewaySplitPartial(t *testing.T) {
	xml := readTestdata(t, "m2_inclusive_gateway.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	// Only path_a true → only A taken (default not taken since one condition matched)
	vars := map[string]any{"path_a": "true", "path_c": "false"}
	instanceID, err := eng.CreateInstance(ctx, dep, vars)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	waiting := allWaiting(inst)
	sort.Strings(waiting)
	if len(waiting) != 1 {
		t.Fatalf("expected 1 waiting token, got %v", waiting)
	}
	expected := []string{"Task_A"}
	for i, e := range expected {
		if waiting[i] != e {
			t.Fatalf("waiting[%d]=%s want %s", i, waiting[i], e)
		}
	}
	// Complete — join should fire with 1 arrival (only token)
	for _, elem := range expected {
		inst, _ = eng.GetInstance(instanceID)
		tokenID := tokenAtElement(inst, elem)
		if err := eng.Complete(ctx, instanceID, elem, tokenID, nil); err != nil {
			t.Fatalf("complete %s: %v", elem, err)
		}
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestInclusiveGatewayDefaultOnly(t *testing.T) {
	xml := readTestdata(t, "m2_inclusive_gateway.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	// No conditions true → only default (B)
	vars := map[string]any{"path_a": "false", "path_c": "false"}
	instanceID, err := eng.CreateInstance(ctx, dep, vars)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	waiting := allWaiting(inst)
	if len(waiting) != 1 || waiting[0] != "Task_B" {
		t.Fatalf("expected [Task_B], got %v", waiting)
	}
	tokenID := tokenAtElement(inst, "Task_B")
	if err := eng.Complete(ctx, instanceID, "Task_B", tokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func allWaiting(inst *projection.Instance) []string {
	var elems []string
	for _, tok := range inst.Tokens {
		if tok != nil && tok.Status == projection.TokenWaiting {
			elems = append(elems, tok.ElementID)
		}
	}
	return elems
}

func tokenAtElement(inst *projection.Instance, elementID string) string {
	for tid, tok := range inst.Tokens {
		if tok != nil && tok.ElementID == elementID && tok.Status == projection.TokenWaiting {
			return tid
		}
	}
	return ""
}
