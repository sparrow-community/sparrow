package processing_test

import (
	"context"
	"strings"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func TestCompensateEventSubProcessRunsHandler(t *testing.T) {
	xml := readTestdataCompensation(t, "m27_compensate_event_subprocess.bpmn")
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
		if tok != nil && tok.ElementID == "Undo_ESP" && tok.Status == projection.TokenWaiting {
			undoTok = tid
		}
	}
	if undoTok == "" {
		t.Fatalf("missing Undo_ESP waiting tokens=%#v waiting=%v", inst.Tokens, allWaiting(inst))
	}
	if _, ok := inst.CompensationSubs["Start_Compensate"]; !ok {
		// Subscription consumed when throw started the event sub-process.
		events, _ := eng.ListEvents(ctx, instanceID)
		if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "Start_Compensate", eventv1.Element_INTENT_ACTIVATED) {
			t.Fatal("expected compensation subscription ACTIVATED for Start_Compensate")
		}
	}
	if err := eng.Complete(ctx, instanceID, "Undo_ESP", undoTok, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s waiting=%v tokens=%#v", inst.Status, allWaiting(inst), inst.Tokens)
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_SUB_PROCESS, "Event_SubProcess_Compensate", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected Event_SubProcess_Compensate COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_INTERMEDIATE_THROW_EVENT, "CompensateThrow_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected CompensateThrow_1 COMPLETED")
	}
}

func TestCompensateEventSubProcessNestedThrows(t *testing.T) {
	xml := readTestdataCompensation(t, "m27_compensate_event_subprocess_nested.bpmn")
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
	elemID, tokenID = waitingAt(inst)
	if elemID != "Task_B" {
		t.Fatalf("expected Task_B, got %s", elemID)
	}
	if err := eng.Complete(ctx, instanceID, elemID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, instanceID)
	// Reverse order: Undo_B then Undo_A.
	var undoTok string
	var undoID string
	for tid, tok := range inst.Tokens {
		if tok == nil || tok.Status != projection.TokenWaiting {
			continue
		}
		if tok.ElementID == "Undo_B" || tok.ElementID == "Undo_A" {
			undoTok = tid
			undoID = tok.ElementID
			break
		}
	}
	if undoID != "Undo_B" {
		t.Fatalf("expected Undo_B first, got %q waiting=%v", undoID, allWaiting(inst))
	}
	if err := eng.Complete(ctx, instanceID, undoID, undoTok, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, instanceID)
	undoTok = ""
	for tid, tok := range inst.Tokens {
		if tok != nil && tok.ElementID == "Undo_A" && tok.Status == projection.TokenWaiting {
			undoTok = tid
		}
	}
	if undoTok == "" {
		t.Fatalf("missing Undo_A waiting=%v tokens=%#v", allWaiting(inst), inst.Tokens)
	}
	if err := eng.Complete(ctx, instanceID, "Undo_A", undoTok, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s waiting=%v tokens=%#v", inst.Status, allWaiting(inst), inst.Tokens)
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "Undo_B", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected Undo_B COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "Undo_A", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected Undo_A COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_SUB_PROCESS, "Event_SubProcess_Compensate", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected Event_SubProcess_Compensate COMPLETED")
	}
}

func TestCompensateEventSubProcessRecover(t *testing.T) {
	xml := readTestdataCompensation(t, "m27_compensate_event_subprocess.bpmn")
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
	elemID, tokenID := waitingAt(inst)
	if err := eng1.Complete(ctx, instanceID, elemID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng1, instanceID)
	var undoTok string
	for tid, tok := range inst.Tokens {
		if tok != nil && tok.ElementID == "Undo_ESP" && tok.Status == projection.TokenWaiting {
			undoTok = tid
		}
	}
	if undoTok == "" {
		t.Fatalf("missing Undo_ESP before recover waiting=%v", allWaiting(inst))
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
		if tok != nil && tok.ElementID == "Undo_ESP" && tok.Status == projection.TokenWaiting {
			undoTok = tid
		}
	}
	if undoTok == "" {
		t.Fatalf("missing Undo_ESP after recover waiting=%v", allWaiting(inst))
	}
	if err := eng2.Complete(ctx, instanceID, "Undo_ESP", undoTok, nil); err != nil {
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
		t.Fatalf("status=%s waiting=%v tokens=%v pending=%#v", inst.Status, allWaiting(inst), parts, inst.PendingCompensation)
	}
}

func TestCompensateEventSubProcessDeployRejects(t *testing.T) {
	ctx := context.Background()
	eng := processing.NewEngine(eventlog.NewMemory())

	root := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<definitions xmlns="http://www.omg.org/spec/BPMN/20100524/MODEL" id="Defs" targetNamespace="t">
  <process id="Process_1" isExecutable="true">
    <startEvent id="Start_1"><outgoing>Flow_1</outgoing></startEvent>
    <userTask id="Task_A"><incoming>Flow_1</incoming><outgoing>Flow_2</outgoing></userTask>
    <endEvent id="End_1"><incoming>Flow_2</incoming></endEvent>
    <sequenceFlow id="Flow_1" sourceRef="Start_1" targetRef="Task_A"/>
    <sequenceFlow id="Flow_2" sourceRef="Task_A" targetRef="End_1"/>
    <subProcess id="Event_SubProcess_Compensate" triggeredByEvent="true">
      <startEvent id="Start_Compensate">
        <outgoing>Flow_esp</outgoing>
        <compensateEventDefinition id="CompensateStartDef_1"/>
      </startEvent>
      <endEvent id="End_ESP"><incoming>Flow_esp</incoming></endEvent>
      <sequenceFlow id="Flow_esp" sourceRef="Start_Compensate" targetRef="End_ESP"/>
    </subProcess>
  </process>
</definitions>`)
	if _, err := eng.Deploy(ctx, root); err == nil || !strings.Contains(err.Error(), "UNSUPPORTED_ELEMENT") {
		t.Fatalf("expected UNSUPPORTED_ELEMENT for process-root compensation event sub-process, got %v", err)
	}

	both := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<definitions xmlns="http://www.omg.org/spec/BPMN/20100524/MODEL" id="Defs" targetNamespace="t">
  <process id="Process_1" isExecutable="true">
    <startEvent id="Start_1"><outgoing>Flow_start</outgoing></startEvent>
    <subProcess id="SubProcess_1">
      <incoming>Flow_start</incoming>
      <outgoing>Flow_sp_end</outgoing>
      <startEvent id="Start_sp"><outgoing>Flow_sp</outgoing></startEvent>
      <endEvent id="End_sp"><incoming>Flow_sp</incoming></endEvent>
      <sequenceFlow id="Flow_sp" sourceRef="Start_sp" targetRef="End_sp"/>
      <subProcess id="Event_SubProcess_Compensate" triggeredByEvent="true">
        <startEvent id="Start_Compensate">
          <outgoing>Flow_esp</outgoing>
          <compensateEventDefinition id="CompensateStartDef_1"/>
        </startEvent>
        <endEvent id="End_ESP"><incoming>Flow_esp</incoming></endEvent>
        <sequenceFlow id="Flow_esp" sourceRef="Start_Compensate" targetRef="End_ESP"/>
      </subProcess>
    </subProcess>
    <boundaryEvent id="CompensateBoundary_SP" attachedToRef="SubProcess_1" cancelActivity="false">
      <compensateEventDefinition id="CompensateDef_SP"/>
    </boundaryEvent>
    <userTask id="Undo_SP" isForCompensation="true"/>
    <association id="Assoc_SP" sourceRef="CompensateBoundary_SP" targetRef="Undo_SP" associationDirection="One"/>
    <endEvent id="End_1"><incoming>Flow_sp_end</incoming></endEvent>
    <sequenceFlow id="Flow_start" sourceRef="Start_1" targetRef="SubProcess_1"/>
    <sequenceFlow id="Flow_sp_end" sourceRef="SubProcess_1" targetRef="End_1"/>
  </process>
</definitions>`)
	if _, err := eng.Deploy(ctx, both); err == nil || !strings.Contains(err.Error(), "UNSUPPORTED_ELEMENT") {
		t.Fatalf("expected UNSUPPORTED_ELEMENT when boundary and compensation event sub-process coexist, got %v", err)
	}
}
