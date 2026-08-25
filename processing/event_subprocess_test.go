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

func TestEventSubProcessMessageInterrupting(t *testing.T) {
	xml := readTestdataESP(t, "m3_event_subprocess_message.bpmn")
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
	if waitingElement(inst) != "UserTask_1" {
		t.Fatalf("expected UserTask_1, got %s", waitingElement(inst))
	}
	if _, ok := inst.EventSubProcesses["Event_SubProcess_1"]; !ok {
		t.Fatalf("expected ESP armed, arms=%#v", inst.EventSubProcesses)
	}

	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "cancel.order"})
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("delivered=%d", n)
	}
	inst = mustInstance(t, eng, instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s tokens=%#v", inst.Status, inst.Tokens)
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_SUB_PROCESS, "Event_SubProcess_1", eventv1.Element_INTENT_ACTIVATED) {
		t.Fatal("expected ESP ACTIVATED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_SUB_PROCESS, "Event_SubProcess_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected ESP COMPLETED")
	}
}

func TestEventSubProcessMessageNonInterrupting(t *testing.T) {
	xml := readTestdataESP(t, "m3_event_subprocess_message_ni.bpmn")
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

	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "notify"})
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("delivered=%d", n)
	}
	inst := mustInstance(t, eng, instanceID)
	if inst.Status != projection.StatusActive {
		t.Fatalf("status=%s", inst.Status)
	}
	if waitingElement(inst) != "UserTask_1" {
		t.Fatalf("UserTask should still wait, got %s tokens=%#v", waitingElement(inst), inst.Tokens)
	}
	elemID, tokenID := waitingAt(inst)
	if err := eng.Complete(ctx, instanceID, elemID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s", inst.Status)
	}
}

func TestEventSubProcessTimerInterrupting(t *testing.T) {
	xml := readTestdataESP(t, "m3_event_subprocess_timer.bpmn")
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
	if err := eng.FireDue(ctx); err != nil {
		t.Fatal(err)
	}
	inst := mustInstance(t, eng, instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s tokens=%#v", inst.Status, inst.Tokens)
	}
}

func TestEventSubProcessBufferedMessage(t *testing.T) {
	xml := readTestdataESP(t, "m3_event_subprocess_message.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()

	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "cancel.order"})
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("expected buffer, delivered=%d", n)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst := mustInstance(t, eng, instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("buffered message should trigger ESP, status=%s", inst.Status)
	}
}

func TestNestedEventSubProcessMessageInterrupting(t *testing.T) {
	xml := readTestdataESP(t, "m4_nested_esp.bpmn")
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
		t.Fatalf("expected UserTask_sp, got %s tokens=%#v", waitingElement(inst), inst.Tokens)
	}
	if _, ok := inst.EventSubProcesses["Event_SubProcess_nested"]; !ok {
		t.Fatalf("expected nested ESP armed while inside SubProcess, arms=%#v", inst.EventSubProcesses)
	}

	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "cancel.inner"})
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("delivered=%d", n)
	}
	inst = mustInstance(t, eng, instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s tokens=%#v", inst.Status, inst.Tokens)
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_SUB_PROCESS, "SubProcess_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected embedding SubProcess TERMINATED by interrupting nested ESP")
	}
	if tid := elementTokenID(events, eventv1.Element_TYPE_SUB_PROCESS, "SubProcess_1", eventv1.Element_INTENT_TERMINATED); tid != "" {
		t.Fatalf("audit-only SubProcess TERMINATED must have empty token_id, got %q", tid)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_SUB_PROCESS, "Event_SubProcess_nested", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected nested ESP COMPLETED")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_ok", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("happy-path end must not complete after interrupting nested ESP")
	}
}

func TestRecoverRearmsNestedEventSubProcess(t *testing.T) {
	xml := readTestdataESP(t, "m4_nested_esp.bpmn")
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
		t.Fatalf("expected nested ESP armed before recover, arms=%#v", inst.EventSubProcesses)
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
		t.Fatalf("expected nested ESP re-armed after Recover, arms=%#v", inst.EventSubProcesses)
	}
	n, err := eng2.PublishMessage(ctx, processing.PublishMessageRequest{Name: "cancel.inner"})
	if err != nil || n != 1 {
		t.Fatalf("publish after recover n=%d err=%v", n, err)
	}
	inst = mustInstance(t, eng2, instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s after nested ESP trigger post-recover", inst.Status)
	}
}

func TestNestedEventSubProcessDisarmedAfterComplete(t *testing.T) {
	xml := readTestdataESP(t, "m4_nested_esp.bpmn")
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
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s", inst.Status)
	}
	if _, ok := inst.EventSubProcesses["Event_SubProcess_nested"]; ok {
		t.Fatalf("nested ESP should be disarmed after SubProcess completes, arms=%#v", inst.EventSubProcesses)
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_ok", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected EndEvent_ok COMPLETED")
	}
}

func readTestdataESP(t *testing.T, name string) []byte {
	t.Helper()
	xml, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return xml
}
