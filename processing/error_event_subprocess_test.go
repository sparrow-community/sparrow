package processing_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	"github.com/sparrow-community/sparrow/processing/deploy"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func TestErrorEventSubProcessInterrupting(t *testing.T) {
	xml := readTestdataErrorEventSubProcess(t, "m5_error_event_subprocess.bpmn")
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
	if waitingElement(inst) != "UserTask_EventSubProcess_1" {
		t.Fatalf("expected UserTask_EventSubProcess_1, got %s tokens=%#v", waitingElement(inst), inst.Tokens)
	}
	if _, ok := inst.EventSubProcesses["Event_SubProcess_1"]; ok {
		t.Fatalf("interrupting Event Sub-Process should be disarmed after trigger, arms=%#v", inst.EventSubProcesses)
	}

	elemID, tokenID := waitingAt(inst)
	if err := eng.Complete(ctx, instanceID, elemID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawEventSubProcessStartActivated(events, "Event_SubProcess_1") {
		t.Fatal("expected error Event Sub-Process start ACTIVATED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_err", eventv1.Element_INTENT_ERROR_THROWN) {
		t.Fatal("expected ERROR_THROWN on error end")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_SUB_PROCESS, "Event_SubProcess_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected error Event Sub-Process COMPLETED")
	}
}

func TestErrorEventSubProcessNonInterrupting(t *testing.T) {
	xml := readTestdataErrorEventSubProcess(t, "m5_error_event_subprocess_ni.bpmn")
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
	if inst.Status != projection.StatusActive {
		t.Fatalf("status=%s", inst.Status)
	}
	if !hasWaitingElement(inst, "UserTask_main") {
		t.Fatalf("main user task should still wait, tokens=%#v", inst.Tokens)
	}
	if !hasWaitingElement(inst, "UserTask_EventSubProcess_1") {
		t.Fatalf("Event Sub-Process handler should wait, tokens=%#v", inst.Tokens)
	}

	handlerElem, handlerTok := waitingTokenAt(inst, "UserTask_EventSubProcess_1")
	if err := eng.Complete(ctx, instanceID, handlerElem, handlerTok, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, instanceID)
	if !hasWaitingElement(inst, "UserTask_main") {
		t.Fatalf("main user task must remain after non-interrupting Event Sub-Process completes, tokens=%#v", inst.Tokens)
	}

	mainElem, mainTok := waitingTokenAt(inst, "UserTask_main")
	if err := eng.Complete(ctx, instanceID, mainElem, mainTok, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestErrorEventSubProcessCodeMismatchTerminates(t *testing.T) {
	xml := readTestdataErrorEventSubProcess(t, "m5_error_event_subprocess_code_miss.bpmn")
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
	if inst.Status != projection.StatusTerminated {
		t.Fatalf("status=%s want terminated for unmatched E2", inst.Status)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if sawElementIntent(events, eventv1.Element_TYPE_SUB_PROCESS, "Event_SubProcess_1", eventv1.Element_INTENT_ACTIVATED) {
		t.Fatal("Event Sub-Process must not start for unmatched error code")
	}
}

func TestNestedErrorEventSubProcess(t *testing.T) {
	xml := readTestdataErrorEventSubProcess(t, "m5_error_event_subprocess_nested.bpmn")
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
	if waitingElement(inst) != "UserTask_sp" {
		t.Fatalf("expected UserTask_sp, got %s", waitingElement(inst))
	}
	if _, ok := inst.EventSubProcesses["Event_SubProcess_nested"]; !ok {
		t.Fatalf("expected nested error Event Sub-Process armed, arms=%#v", inst.EventSubProcesses)
	}

	elemID, tokenID := waitingAt(inst)
	if err := eng.Complete(ctx, instanceID, elemID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, instanceID)
	if waitingElement(inst) != "UserTask_EventSubProcess_1" {
		t.Fatalf("expected UserTask_EventSubProcess_1 after inner error, got %s tokens=%#v", waitingElement(inst), inst.Tokens)
	}

	elemID, tokenID = waitingAt(inst)
	if err := eng.Complete(ctx, instanceID, elemID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed via nested handler", inst.Status)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_SUB_PROCESS, "SubProcess_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected embedding SubProcess TERMINATED by interrupting nested error Event Sub-Process")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_ok", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("happy-path end must not complete after interrupting nested error Event Sub-Process")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_PROCESS, "Process_error_event_subprocess_nested", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("process must complete successfully, not terminate")
	}
}

func TestRecoverRearmsErrorEventSubProcess(t *testing.T) {
	xml := readTestdataErrorEventSubProcess(t, "m5_error_event_subprocess_nested.bpmn")
	ctx := context.Background()
	memLog := eventlog.NewMemory()
	store := deploy.NewMemoryStore()
	eng1, err := processing.Recover(ctx, memLog, store, nil)
	if err != nil {
		t.Fatal(err)
	}
	depID, err := eng1.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng1.CreateInstance(ctx, depID, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst := mustInstance(t, eng1, instanceID)
	if _, ok := inst.EventSubProcesses["Event_SubProcess_nested"]; !ok {
		t.Fatalf("expected nested error Event Sub-Process armed before recover, arms=%#v", inst.EventSubProcesses)
	}

	eng2, err := processing.Recover(ctx, memLog, store, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng2, instanceID)
	if waitingElement(inst) != "UserTask_sp" {
		t.Fatalf("expected UserTask_sp after recover, got %s", waitingElement(inst))
	}
	if _, ok := inst.EventSubProcesses["Event_SubProcess_nested"]; !ok {
		t.Fatalf("expected nested error Event Sub-Process re-armed after Recover, arms=%#v", inst.EventSubProcesses)
	}

	elemID, tokenID := waitingAt(inst)
	if err := eng2.Complete(ctx, instanceID, elemID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng2, instanceID)
	if waitingElement(inst) != "UserTask_EventSubProcess_1" {
		t.Fatalf("expected UserTask_EventSubProcess_1 after error post-recover, got %s", waitingElement(inst))
	}
	elemID, tokenID = waitingAt(inst)
	if err := eng2.Complete(ctx, instanceID, elemID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng2, instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s after nested error Event Sub-Process post-recover", inst.Status)
	}
}

func readTestdataErrorEventSubProcess(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func hasWaitingElement(inst *projection.Instance, elementID string) bool {
	for _, tok := range inst.Tokens {
		if tok != nil && tok.Status == projection.TokenWaiting && tok.ElementID == elementID {
			return true
		}
	}
	return false
}

func waitingTokenAt(inst *projection.Instance, elementID string) (elemID, tokenID string) {
	for tid, tok := range inst.Tokens {
		if tok != nil && tok.Status == projection.TokenWaiting && tok.ElementID == elementID {
			return elementID, tid
		}
	}
	return "", ""
}
