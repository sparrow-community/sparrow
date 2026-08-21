package processing_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func TestCompensationThrowRunsHandler(t *testing.T) {
	xml := readTestdataCompensation(t, "m3_compensation.bpmn")
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
	inst := mustInstance(t, eng, instanceID)
	elemID, tokenID := waitingAt(inst)
	if elemID != "Task_A" {
		t.Fatalf("expected Task_A, got %s", elemID)
	}
	if err := eng.Complete(ctx, instanceID, elemID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, instanceID)
	// After Task_A: compensate throw waits while Undo_A waits.
	waiting := allWaiting(inst)
	if len(waiting) != 2 {
		t.Fatalf("expected throw + Undo_A waiting, got %v tokens=%#v", waiting, inst.Tokens)
	}
	var undoTok string
	for tid, tok := range inst.Tokens {
		if tok != nil && tok.ElementID == "Undo_A" && tok.Status == projection.TokenWaiting {
			undoTok = tid
		}
	}
	if undoTok == "" {
		t.Fatal("missing Undo_A waiting token")
	}
	if err := eng.Complete(ctx, instanceID, "Undo_A", undoTok, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s tokens=%#v", inst.Status, inst.Tokens)
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "Undo_A", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected Undo_A COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_INTERMEDIATE_THROW_EVENT, "CompensateThrow_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected CompensateThrow_1 COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "CompensateBoundary_A", eventv1.Element_INTENT_ACTIVATED) {
		t.Fatal("expected compensation boundary ACTIVATED after Task_A")
	}
}

func TestCompensationDeploy(t *testing.T) {
	xml := readTestdataCompensation(t, "m3_compensation.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	if _, err := eng.Deploy(context.Background(), xml); err != nil {
		t.Fatal(err)
	}
}

func readTestdataCompensation(t *testing.T, name string) []byte {
	t.Helper()
	xml, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return xml
}
