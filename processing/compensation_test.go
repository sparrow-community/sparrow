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

func TestCompensateEndRunsHandler(t *testing.T) {
	xml := readTestdataCompensation(t, "m4_compensate_end.bpmn")
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
	if err := eng.Complete(ctx, instanceID, elemID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, instanceID)
	var undoTok string
	for tid, tok := range inst.Tokens {
		if tok != nil && tok.ElementID == "Undo_A" && tok.Status == projection.TokenWaiting {
			undoTok = tid
		}
	}
	if undoTok == "" {
		t.Fatalf("missing Undo_A waiting tokens=%#v", inst.Tokens)
	}
	if err := eng.Complete(ctx, instanceID, "Undo_A", undoTok, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s tokens=%#v", inst.Status, inst.Tokens)
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "CompensateEnd_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected CompensateEnd_1 COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "Undo_A", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected Undo_A COMPLETED")
	}
}

func TestCompensateSubProcessBoundary(t *testing.T) {
	xml := readTestdataCompensation(t, "m4_compensate_subprocess.bpmn")
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
	if elemID != "Task_inner" {
		t.Fatalf("expected Task_inner, got %s", elemID)
	}
	if err := eng.Complete(ctx, instanceID, elemID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, instanceID)
	var undoTok string
	for tid, tok := range inst.Tokens {
		if tok != nil && tok.ElementID == "Undo_SP" && tok.Status == projection.TokenWaiting {
			undoTok = tid
		}
	}
	if undoTok == "" {
		t.Fatalf("missing Undo_SP waiting tokens=%#v", inst.Tokens)
	}
	if err := eng.Complete(ctx, instanceID, "Undo_SP", undoTok, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s tokens=%#v", inst.Status, inst.Tokens)
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "CompensateBoundary_SP", eventv1.Element_INTENT_ACTIVATED) {
		t.Fatal("expected SubProcess compensation boundary ACTIVATED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "Undo_SP", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected Undo_SP COMPLETED")
	}
}

func TestCompensationDeploy(t *testing.T) {
	xml := readTestdataCompensation(t, "m3_compensation.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	if _, err := eng.Deploy(context.Background(), xml); err != nil {
		t.Fatal(err)
	}
}

func TestCallActivityCompensationThrowRunsHandler(t *testing.T) {
	xml := readTestdataCall(t, "m9_call_compensate.bpmn")
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
	hostElem, hostTok := waitingAt(caller)
	if hostElem != "CallActivity_1" {
		t.Fatalf("expected CallActivity_1, got %s", hostElem)
	}
	childID := caller.Tokens[hostTok].CalledProcessInstanceID
	elemID, tokenID := waitingAt(mustInstance(t, eng, childID))
	if err := eng.Complete(ctx, childID, elemID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	caller = mustInstance(t, eng, callerID)
	waiting := allWaiting(caller)
	if len(waiting) != 2 {
		t.Fatalf("expected throw + Undo_CA waiting, got %v tokens=%#v", waiting, caller.Tokens)
	}
	var undoTok string
	for tid, tok := range caller.Tokens {
		if tok != nil && tok.ElementID == "Undo_CA" && tok.Status == projection.TokenWaiting {
			undoTok = tid
		}
	}
	if undoTok == "" {
		t.Fatal("missing Undo_CA waiting token")
	}
	if err := eng.Complete(ctx, callerID, "Undo_CA", undoTok, nil); err != nil {
		t.Fatal(err)
	}
	caller = mustInstance(t, eng, callerID)
	if caller.Status != projection.StatusCompleted {
		t.Fatalf("status=%s tokens=%#v", caller.Status, caller.Tokens)
	}
	events, _ := eng.ListEvents(ctx, callerID)
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "Undo_CA", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected Undo_CA COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_INTERMEDIATE_THROW_EVENT, "CompensateThrow_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected CompensateThrow_1 COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "CompensateBoundary_CA", eventv1.Element_INTENT_ACTIVATED) {
		t.Fatal("expected compensation boundary ACTIVATED after CallActivity")
	}
}

func TestCompensateUnfinishedSubProcessBroadcast(t *testing.T) {
	xml := readTestdataCompensation(t, "m12_compensate_unfinished_sp_parallel.bpmn")
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
	var taskATok, sibTok string
	for tid, tok := range inst.Tokens {
		if tok == nil || tok.Status != projection.TokenWaiting {
			continue
		}
		switch tok.ElementID {
		case "Task_A":
			taskATok = tid
		case "Task_sib":
			sibTok = tid
		}
	}
	if taskATok == "" || sibTok == "" {
		t.Fatalf("expected Task_A + Task_sib waiting, got %v tokens=%#v", allWaiting(inst), inst.Tokens)
	}
	if err := eng.Complete(ctx, instanceID, "Task_A", taskATok, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, instanceID)
	if waitingElement(inst) != "Task_B" && !containsWaiting(inst, "Task_B") {
		t.Fatalf("expected Task_B waiting after Task_A, waiting=%v tokens=%#v", allWaiting(inst), inst.Tokens)
	}
	if err := eng.Complete(ctx, instanceID, "Task_sib", sibTok, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, instanceID)
	if containsWaiting(inst, "Task_B") {
		t.Fatalf("Task_B should be terminated, waiting=%v", allWaiting(inst))
	}
	var undoTok string
	for tid, tok := range inst.Tokens {
		if tok != nil && tok.ElementID == "Undo_A" && tok.Status == projection.TokenWaiting {
			undoTok = tid
		}
	}
	if undoTok == "" {
		t.Fatalf("missing Undo_A waiting tokens=%#v", inst.Tokens)
	}
	if err := eng.Complete(ctx, instanceID, "Undo_A", undoTok, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s tokens=%#v", inst.Status, inst.Tokens)
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "Task_B", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected Task_B TERMINATED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "Undo_A", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected Undo_A COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_INTERMEDIATE_THROW_EVENT, "CompensateThrow_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected CompensateThrow_1 COMPLETED")
	}
}

func TestCompensateUnfinishedSubProcessTargeted(t *testing.T) {
	xml := readTestdataCompensation(t, "m12_compensate_unfinished_sp_targeted.bpmn")
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
	var taskATok, taskXTok string
	for tid, tok := range inst.Tokens {
		if tok == nil || tok.Status != projection.TokenWaiting {
			continue
		}
		switch tok.ElementID {
		case "Task_A":
			taskATok = tid
		case "Task_X":
			taskXTok = tid
		}
	}
	if taskATok == "" || taskXTok == "" {
		t.Fatalf("missing tokens tokens=%#v", inst.Tokens)
	}
	if err := eng.Complete(ctx, instanceID, "Task_A", taskATok, nil); err != nil {
		t.Fatal(err)
	}
	if err := eng.Complete(ctx, instanceID, "Task_X", taskXTok, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, instanceID)
	var undoATok string
	for tid, tok := range inst.Tokens {
		if tok != nil && tok.ElementID == "Undo_A" && tok.Status == projection.TokenWaiting {
			undoATok = tid
		}
		if tok != nil && tok.ElementID == "Undo_X" && tok.Status == projection.TokenWaiting {
			t.Fatal("Undo_X must not run for targeted SubProcess_1 compensate")
		}
	}
	if undoATok == "" {
		t.Fatalf("missing Undo_A waiting tokens=%#v", inst.Tokens)
	}
	if _, ok := inst.CompensationSubs["CompensateBoundary_X"]; !ok {
		// Subscription may be keyed by boundary id — verify X was not consumed.
		events, _ := eng.ListEvents(ctx, instanceID)
		if sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "CompensateBoundary_X", eventv1.Element_INTENT_COMPLETED) {
			t.Fatal("CompensateBoundary_X must not be completed by targeted SP throw")
		}
	}
	if err := eng.Complete(ctx, instanceID, "Undo_A", undoATok, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s tokens=%#v", inst.Status, inst.Tokens)
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "Undo_X", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("Undo_X must not complete")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "Undo_A", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected Undo_A COMPLETED")
	}
}

func TestCompensateUnfinishedSubProcessRecover(t *testing.T) {
	xml := readTestdataCompensation(t, "m12_compensate_unfinished_sp_parallel.bpmn")
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
	inst := mustInstance(t, eng1, instanceID)
	var taskATok, sibTok string
	for tid, tok := range inst.Tokens {
		if tok == nil || tok.Status != projection.TokenWaiting {
			continue
		}
		switch tok.ElementID {
		case "Task_A":
			taskATok = tid
		case "Task_sib":
			sibTok = tid
		}
	}
	if err := eng1.Complete(ctx, instanceID, "Task_A", taskATok, nil); err != nil {
		t.Fatal(err)
	}
	if err := eng1.Complete(ctx, instanceID, "Task_sib", sibTok, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng1, instanceID)
	var undoTok string
	for tid, tok := range inst.Tokens {
		if tok != nil && tok.ElementID == "Undo_A" && tok.Status == projection.TokenWaiting {
			undoTok = tid
		}
	}
	if undoTok == "" {
		t.Fatalf("missing Undo_A before recover tokens=%#v", inst.Tokens)
	}
	if err := eng1.Close(); err != nil {
		t.Fatal(err)
	}

	eng2, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	defer eng2.Close()
	inst = mustInstance(t, eng2, instanceID)
	undoTok = ""
	for tid, tok := range inst.Tokens {
		if tok != nil && tok.ElementID == "Undo_A" && tok.Status == projection.TokenWaiting {
			undoTok = tid
		}
	}
	if undoTok == "" {
		var parts []string
		for tid, tok := range inst.Tokens {
			if tok == nil {
				continue
			}
			parts = append(parts, tid+":"+tok.ElementID+":"+string(tok.Status))
		}
		t.Fatalf("missing Undo_A after recover waiting=%v tokens=%v", allWaiting(inst), parts)
	}
	if containsWaiting(inst, "Task_B") || containsWaiting(inst, "SubProcess_1") {
		t.Fatalf("unfinished SP work revived after recover waiting=%v", allWaiting(inst))
	}
	if err := eng2.Complete(ctx, instanceID, "Undo_A", undoTok, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng2, instanceID)
	if inst.Status != projection.StatusCompleted {
		var parts []string
		for tid, tok := range inst.Tokens {
			if tok == nil {
				continue
			}
			parts = append(parts, tid+":"+tok.ElementID+":"+string(tok.Status))
		}
		t.Fatalf("status=%s waiting=%v tokens=%v", inst.Status, allWaiting(inst), parts)
	}
}

func containsWaiting(inst *projection.Instance, elementID string) bool {
	for _, tok := range inst.Tokens {
		if tok != nil && tok.Status == projection.TokenWaiting && tok.ElementID == elementID {
			return true
		}
	}
	return false
}

func readTestdataCompensation(t *testing.T, name string) []byte {
	t.Helper()
	xml, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return xml
}
