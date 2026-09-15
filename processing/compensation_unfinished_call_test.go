package processing_test

import (
	"context"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func TestCompensateUnfinishedCallActivityBroadcast(t *testing.T) {
	xml := readTestdataCompensation(t, "m28_compensate_unfinished_call.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()

	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	callerID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	caller := mustInstance(t, eng, callerID)
	var hostTok, sibTok string
	for tid, tok := range caller.Tokens {
		if tok == nil || tok.Status != projection.TokenWaiting {
			continue
		}
		switch tok.ElementID {
		case "CallActivity_1":
			hostTok = tid
		case "Task_sib":
			sibTok = tid
		}
	}
	if hostTok == "" || sibTok == "" {
		t.Fatalf("expected CallActivity_1 + Task_sib waiting=%v", allWaiting(caller))
	}
	childID := caller.Tokens[hostTok].CalledProcessInstanceID
	if childID == "" {
		t.Fatal("missing child instance id")
	}
	child := mustInstance(t, eng, childID)
	elemID, tokenID := waitingAt(child)
	if elemID != "Task_A" {
		t.Fatalf("expected Task_A on child, got %s", elemID)
	}
	if err := eng.Complete(ctx, childID, elemID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	child = mustInstance(t, eng, childID)
	if waitingElement(child) != "Task_B" && !containsWaiting(child, "Task_B") {
		t.Fatalf("expected Task_B waiting on child waiting=%v", allWaiting(child))
	}
	if err := eng.Complete(ctx, callerID, "Task_sib", sibTok, nil); err != nil {
		t.Fatal(err)
	}
	child = mustInstance(t, eng, childID)
	var undoTok string
	for tid, tok := range child.Tokens {
		if tok != nil && tok.ElementID == "Undo_A" && tok.Status == projection.TokenWaiting {
			undoTok = tid
		}
	}
	if undoTok == "" {
		t.Fatalf("missing Undo_A on child waiting=%v tokens=%#v status=%s", allWaiting(child), child.Tokens, child.Status)
	}
	if containsWaiting(child, "Task_B") {
		t.Fatal("Task_B should be terminated")
	}
	if err := eng.Complete(ctx, childID, "Undo_A", undoTok, nil); err != nil {
		t.Fatal(err)
	}
	child = mustInstance(t, eng, childID)
	if child.Status != projection.StatusTerminated {
		t.Fatalf("child status=%s want terminated", child.Status)
	}
	caller = mustInstance(t, eng, callerID)
	if caller.Status != projection.StatusCompleted {
		t.Fatalf("caller status=%s waiting=%v tokens=%#v", caller.Status, allWaiting(caller), caller.Tokens)
	}
	events, _ := eng.ListEvents(ctx, childID)
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "Undo_A", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected Undo_A COMPLETED on child")
	}
	callerEvents, _ := eng.ListEvents(ctx, callerID)
	if !sawElementIntent(callerEvents, eventv1.Element_TYPE_INTERMEDIATE_THROW_EVENT, "CompensateThrow_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected CompensateThrow_1 COMPLETED")
	}
}

func TestCompensateUnfinishedCallActivityTargeted(t *testing.T) {
	xml := readTestdataCompensation(t, "m28_compensate_unfinished_call_targeted.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()

	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	callerID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	caller := mustInstance(t, eng, callerID)
	var hostTok, taskXTok string
	for tid, tok := range caller.Tokens {
		if tok == nil || tok.Status != projection.TokenWaiting {
			continue
		}
		switch tok.ElementID {
		case "CallActivity_1":
			hostTok = tid
		case "Task_X":
			taskXTok = tid
		}
	}
	if hostTok == "" || taskXTok == "" {
		t.Fatalf("missing hosts waiting=%v", allWaiting(caller))
	}
	childID := caller.Tokens[hostTok].CalledProcessInstanceID
	elemID, tokenID := waitingAt(mustInstance(t, eng, childID))
	if err := eng.Complete(ctx, childID, elemID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	if err := eng.Complete(ctx, callerID, "Task_X", taskXTok, nil); err != nil {
		t.Fatal(err)
	}
	child := mustInstance(t, eng, childID)
	var undoTok string
	for tid, tok := range child.Tokens {
		if tok != nil && tok.ElementID == "Undo_A" && tok.Status == projection.TokenWaiting {
			undoTok = tid
		}
	}
	if undoTok == "" {
		t.Fatalf("missing Undo_A on child waiting=%v", allWaiting(child))
	}
	caller = mustInstance(t, eng, callerID)
	if containsWaiting(caller, "Undo_X") {
		t.Fatal("Undo_X must not run for targeted CallActivity_1 compensate")
	}
	if _, ok := caller.CompensationSubs["CompensateBoundary_X"]; !ok {
		events, _ := eng.ListEvents(ctx, callerID)
		if sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "CompensateBoundary_X", eventv1.Element_INTENT_COMPLETED) {
			t.Fatal("CompensateBoundary_X must not be completed by targeted call throw")
		}
	}
	if err := eng.Complete(ctx, childID, "Undo_A", undoTok, nil); err != nil {
		t.Fatal(err)
	}
	caller = mustInstance(t, eng, callerID)
	if caller.Status != projection.StatusCompleted {
		t.Fatalf("caller status=%s waiting=%v", caller.Status, allWaiting(caller))
	}
	events, _ := eng.ListEvents(ctx, callerID)
	if sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "Undo_X", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("Undo_X must not complete")
	}
}

func TestCompensateUnfinishedCallActivityRecover(t *testing.T) {
	xml := readTestdataCompensation(t, "m28_compensate_unfinished_call.bpmn")
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
	callerID, err := eng1.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	caller := mustInstance(t, eng1, callerID)
	var hostTok, sibTok string
	for tid, tok := range caller.Tokens {
		if tok == nil || tok.Status != projection.TokenWaiting {
			continue
		}
		switch tok.ElementID {
		case "CallActivity_1":
			hostTok = tid
		case "Task_sib":
			sibTok = tid
		}
	}
	childID := caller.Tokens[hostTok].CalledProcessInstanceID
	elemID, tokenID := waitingAt(mustInstance(t, eng1, childID))
	if err := eng1.Complete(ctx, childID, elemID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	if err := eng1.Complete(ctx, callerID, "Task_sib", sibTok, nil); err != nil {
		t.Fatal(err)
	}
	child := mustInstance(t, eng1, childID)
	var undoTok string
	for tid, tok := range child.Tokens {
		if tok != nil && tok.ElementID == "Undo_A" && tok.Status == projection.TokenWaiting {
			undoTok = tid
		}
	}
	if undoTok == "" {
		t.Fatalf("missing Undo_A before recover waiting=%v", allWaiting(child))
	}
	if err := eng1.Close(); err != nil {
		t.Fatal(err)
	}

	eng2, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	defer eng2.Close()
	child = mustInstance(t, eng2, childID)
	undoTok = ""
	for tid, tok := range child.Tokens {
		if tok != nil && tok.ElementID == "Undo_A" && tok.Status == projection.TokenWaiting {
			undoTok = tid
		}
	}
	if undoTok == "" {
		t.Fatalf("missing Undo_A after recover waiting=%v", allWaiting(child))
	}
	if err := eng2.Complete(ctx, childID, "Undo_A", undoTok, nil); err != nil {
		t.Fatal(err)
	}
	caller = mustInstance(t, eng2, callerID)
	if caller.Status != projection.StatusCompleted {
		t.Fatalf("caller status=%s waiting=%v", caller.Status, allWaiting(caller))
	}
}
