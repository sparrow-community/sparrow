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
	xml := readTestdataEventSubProcess(t, "m3_event_subprocess_message.bpmn")
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
		t.Fatalf("expected event sub-process armed, arms=%#v", inst.EventSubProcesses)
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
	if !sawEventSubProcessStartActivated(events, "Event_SubProcess_1") {
		t.Fatal("expected event sub-process start ACTIVATED in event log")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_SUB_PROCESS, "Event_SubProcess_1", eventv1.Element_INTENT_ACTIVATED) {
		t.Fatal("expected event sub-process ACTIVATED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_SUB_PROCESS, "Event_SubProcess_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected event sub-process COMPLETED")
	}
}

func TestEventSubProcessMessageNonInterrupting(t *testing.T) {
	xml := readTestdataEventSubProcess(t, "m3_event_subprocess_message_ni.bpmn")
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
	xml := readTestdataEventSubProcess(t, "m3_event_subprocess_timer.bpmn")
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
	xml := readTestdataEventSubProcess(t, "m3_event_subprocess_message.bpmn")
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
		t.Fatalf("buffered message should trigger event sub-process, status=%s", inst.Status)
	}
}

func TestNestedEventSubProcessMessageInterrupting(t *testing.T) {
	xml := readTestdataEventSubProcess(t, "m4_nested_event_subprocess.bpmn")
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
	hostCount := 0
	for _, tok := range inst.Tokens {
		if tok != nil && tok.ScopeHost {
			hostCount++
			if tok.ElementID != "SubProcess_1" {
				t.Fatalf("scope host on %s, want SubProcess_1", tok.ElementID)
			}
		}
	}
	if hostCount != 1 {
		t.Fatalf("expected one SubProcess host token, got %d tokens=%#v", hostCount, inst.Tokens)
	}
	if _, ok := inst.EventSubProcesses["Event_SubProcess_nested"]; !ok {
		t.Fatalf("expected nested event sub-process armed while inside SubProcess, arms=%#v", inst.EventSubProcesses)
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
		t.Fatal("expected embedding SubProcess TERMINATED by interrupting nested event sub-process")
	}
	if tid := elementTokenID(events, eventv1.Element_TYPE_SUB_PROCESS, "SubProcess_1", eventv1.Element_INTENT_TERMINATED); tid == "" {
		t.Fatal("expected host token_id on SubProcess TERMINATED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_SUB_PROCESS, "Event_SubProcess_nested", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected nested event sub-process COMPLETED")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_ok", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("happy-path end must not complete after interrupting nested event sub-process")
	}
}

func TestRecoverRearmsNestedEventSubProcess(t *testing.T) {
	xml := readTestdataEventSubProcess(t, "m4_nested_event_subprocess.bpmn")
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
		t.Fatalf("expected nested event sub-process armed before recover, arms=%#v", inst.EventSubProcesses)
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
		t.Fatalf("expected nested event sub-process re-armed after Recover, arms=%#v", inst.EventSubProcesses)
	}
	n, err := eng2.PublishMessage(ctx, processing.PublishMessageRequest{Name: "cancel.inner"})
	if err != nil || n != 1 {
		t.Fatalf("publish after recover n=%d err=%v", n, err)
	}
	inst = mustInstance(t, eng2, instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s after nested event sub-process trigger post-recover", inst.Status)
	}
}

func TestNestedEventSubProcessDisarmedAfterComplete(t *testing.T) {
	xml := readTestdataEventSubProcess(t, "m4_nested_event_subprocess.bpmn")
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
		t.Fatalf("nested event sub-process should be disarmed after SubProcess completes, arms=%#v", inst.EventSubProcesses)
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_ok", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected EndEvent_ok COMPLETED")
	}
}

func readTestdataEventSubProcess(t *testing.T, name string) []byte {
	t.Helper()
	xml, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return xml
}

func sawEventSubProcessStartActivated(events []*eventv1.Event, eventSubProcessElementID string) bool {
	for _, ev := range events {
		if ev.GetRecordType() != eventv1.Event_RECORD_TYPE_EVENT {
			continue
		}
		el := ev.GetElement()
		if el == nil || el.GetType() != eventv1.Element_TYPE_START_EVENT {
			continue
		}
		if el.GetIntent() != eventv1.Element_INTENT_ACTIVATED {
			continue
		}
		if p := el.GetEventPayload(); p != nil && p.GetEventSubProcessElementId() == eventSubProcessElementID {
			return true
		}
	}
	return false
}
