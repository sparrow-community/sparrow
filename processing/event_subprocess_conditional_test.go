package processing_test

import (
	"context"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func TestEventSubProcessConditionalInterrupting(t *testing.T) {
	xml := readTestdata(t, "m30_event_subprocess_conditional.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()

	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, map[string]any{"ready": false})
	if err != nil {
		t.Fatal(err)
	}
	inst := mustInstance(t, eng, instanceID)
	if waitingElement(inst) != "UserTask_1" {
		t.Fatalf("expected UserTask_1 waiting=%v", allWaiting(inst))
	}
	if _, ok := inst.EventSubProcesses["Event_SubProcess_1"]; !ok {
		t.Fatalf("expected conditional event sub-process armed arms=%#v", inst.EventSubProcesses)
	}

	n, err := eng.EvaluateConditions(ctx, processing.EvaluateConditionsRequest{
		ProcessInstanceID: instanceID,
		Variables:         map[string]any{"ready": false},
	})
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("false condition should fire 0, got %d", n)
	}
	inst = mustInstance(t, eng, instanceID)
	if waitingElement(inst) != "UserTask_1" {
		t.Fatalf("host must still wait after false eval waiting=%v", allWaiting(inst))
	}

	n, err = eng.EvaluateConditions(ctx, processing.EvaluateConditionsRequest{
		ProcessInstanceID: instanceID,
		Variables:         map[string]any{"ready": true},
	})
	if err != nil || n != 1 {
		t.Fatalf("EvaluateConditions n=%d err=%v", n, err)
	}
	inst = mustInstance(t, eng, instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s waiting=%v", inst.Status, allWaiting(inst))
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_SUB_PROCESS, "Event_SubProcess_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected Event_SubProcess_1 COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "UserTask_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected UserTask_1 TERMINATED")
	}
}

func TestEventSubProcessConditionalNonInterrupting(t *testing.T) {
	xml := readTestdata(t, "m30_event_subprocess_conditional_ni.bpmn")
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
	n, err := eng.EvaluateConditions(ctx, processing.EvaluateConditionsRequest{
		ProcessInstanceID: instanceID,
		Variables:         map[string]any{"ready": true},
	})
	if err != nil || n != 1 {
		t.Fatalf("EvaluateConditions n=%d err=%v", n, err)
	}
	inst := mustInstance(t, eng, instanceID)
	if !containsWaiting(inst, "UserTask_1") {
		t.Fatalf("host UserTask_1 must remain waiting waiting=%v", allWaiting(inst))
	}
	if !containsWaiting(inst, "UserTask_ESP") {
		t.Fatalf("expected UserTask_ESP waiting=%v", allWaiting(inst))
	}
	var espTok, hostTok string
	for tid, tok := range inst.Tokens {
		if tok == nil || tok.Status != projection.TokenWaiting {
			continue
		}
		switch tok.ElementID {
		case "UserTask_ESP":
			espTok = tid
		case "UserTask_1":
			hostTok = tid
		}
	}
	if err := eng.Complete(ctx, instanceID, "UserTask_ESP", espTok, nil); err != nil {
		t.Fatal(err)
	}
	if err := eng.Complete(ctx, instanceID, "UserTask_1", hostTok, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s waiting=%v", inst.Status, allWaiting(inst))
	}
}

func TestEventSubProcessConditionalRecover(t *testing.T) {
	xml := readTestdata(t, "m30_event_subprocess_conditional.bpmn")
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
	instanceID, err := eng1.CreateInstance(ctx, dep, map[string]any{"ready": false})
	if err != nil {
		t.Fatal(err)
	}
	inst := mustInstance(t, eng1, instanceID)
	if _, ok := inst.EventSubProcesses["Event_SubProcess_1"]; !ok {
		t.Fatal("expected armed conditional event sub-process before recover")
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
	if _, ok := inst.EventSubProcesses["Event_SubProcess_1"]; !ok {
		t.Fatalf("expected armed after recover arms=%#v", inst.EventSubProcesses)
	}
	n, err := eng2.EvaluateConditions(ctx, processing.EvaluateConditionsRequest{
		ProcessInstanceID: instanceID,
		Variables:         map[string]any{"ready": true},
	})
	if err != nil || n != 1 {
		t.Fatalf("EvaluateConditions n=%d err=%v", n, err)
	}
	inst = mustInstance(t, eng2, instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s waiting=%v", inst.Status, allWaiting(inst))
	}
}
