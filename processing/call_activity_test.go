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

func TestCallActivityCompletesChildThenCaller(t *testing.T) {
	xml := readTestdataCall(t, "m4_call_activity.bpmn")
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
	if elemID != "Task_called" {
		t.Fatalf("expected Task_called inside called process, got %s tokens=%#v", elemID, inst.Tokens)
	}
	if err := eng.Complete(ctx, instanceID, elemID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s tokens=%#v", inst.Status, inst.Tokens)
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_CALL_ACTIVITY, "CallActivity_1", eventv1.Element_INTENT_ACTIVATED) {
		t.Fatal("expected CallActivity ACTIVATED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_CALL_ACTIVITY, "CallActivity_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected CallActivity COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "Task_called", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected called task COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_caller", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected caller end COMPLETED")
	}
}

func readTestdataCall(t *testing.T, name string) []byte {
	t.Helper()
	xml, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return xml
}
