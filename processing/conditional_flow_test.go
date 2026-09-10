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

func TestConditionalFlowActivityTakesCondition(t *testing.T) {
	xml := readTestdata(t, "m17_conditional_activity.bpmn")
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
	inst, _ := eng.GetInstance(instanceID)
	elementID, tokenID := waitingAt(inst)
	if err := eng.Complete(ctx, instanceID, elementID, tokenID, map[string]any{"approved": true}); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_yes", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected End_yes")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_default", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("default end must not complete")
	}
}

func TestConditionalFlowActivityTakesDefault(t *testing.T) {
	xml := readTestdata(t, "m17_conditional_activity.bpmn")
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
	inst, _ := eng.GetInstance(instanceID)
	elementID, tokenID := waitingAt(inst)
	if err := eng.Complete(ctx, instanceID, elementID, tokenID, map[string]any{"approved": false}); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_default", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected End_default")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_yes", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("yes end must not complete")
	}
}

func TestConditionalFlowNoMatchNoDefault(t *testing.T) {
	xml := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<definitions xmlns="http://www.omg.org/spec/BPMN/20100524/MODEL"
             xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
             id="Definitions_cond_nomatch"
             targetNamespace="http://sparrow.example/cond-nomatch">
  <process id="Process_cond_nomatch" isExecutable="true">
    <startEvent id="Start_1"><outgoing>Flow_start_to_task</outgoing></startEvent>
    <userTask id="UserTask_1">
      <incoming>Flow_start_to_task</incoming>
      <outgoing>Flow_task_to_end</outgoing>
    </userTask>
    <endEvent id="End_1"><incoming>Flow_task_to_end</incoming></endEvent>
    <sequenceFlow id="Flow_start_to_task" sourceRef="Start_1" targetRef="UserTask_1"/>
    <sequenceFlow id="Flow_task_to_end" sourceRef="UserTask_1" targetRef="End_1">
      <conditionExpression xsi:type="tFormalExpression">${approved == true}</conditionExpression>
    </sequenceFlow>
  </process>
</definitions>`)
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
	inst, _ := eng.GetInstance(instanceID)
	elementID, tokenID := waitingAt(inst)
	err = eng.Complete(ctx, instanceID, elementID, tokenID, map[string]any{"approved": false})
	if err == nil {
		t.Fatal("expected NO_OUTGOING_FLOW")
	}
	if !strings.Contains(err.Error(), "NO_OUTGOING_FLOW") {
		t.Fatalf("err=%v", err)
	}
}

func TestConditionalFlowRecoverThenComplete(t *testing.T) {
	xml := readTestdata(t, "m17_conditional_activity.bpmn")
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
	if err := eng1.Close(); err != nil {
		t.Fatal(err)
	}

	eng2, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	defer eng2.Close()

	inst, ok := eng2.GetInstance(instanceID)
	if !ok {
		t.Fatal("instance missing after recover")
	}
	elementID, tokenID := waitingAt(inst)
	if err := eng2.Complete(ctx, instanceID, elementID, tokenID, map[string]any{"approved": true}); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, err := eng2.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_yes", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected End_yes after recover")
	}
}
