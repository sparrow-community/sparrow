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

func TestTimerBoundaryPT0SFireDue(t *testing.T) {
	xml := readTestdata(t, "m2_timer_boundary.bpmn")
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
	if elementID != "UserTask_1" {
		t.Fatalf("tokens=%#v", inst.Tokens)
	}
	tok := inst.Tokens[tokenID]
	if tok.DueUnixMs == 0 || tok.BoundaryID != "TimerBoundary_1" {
		t.Fatalf("armed timer token=%#v", tok)
	}

	if err := eng.FireDue(ctx); err != nil {
		t.Fatalf("FireDue: %v", err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "UserTask_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected UserTask TERMINATED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "TimerBoundary_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected boundary COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_timeout", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected timeout end")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_ok", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("happy-path end must not complete")
	}
}

func TestTimerBoundaryCompleteCancelsTimer(t *testing.T) {
	xml := readTestdata(t, "m2_timer_boundary_1h.bpmn")
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
		t.Fatalf("FireDue: %v", err)
	}
	inst, _ := eng.GetInstance(instanceID)
	elementID, tokenID := waitingAt(inst)
	if elementID != "UserTask_1" {
		t.Fatalf("tokens=%#v", inst.Tokens)
	}
	if err := eng.Complete(ctx, instanceID, elementID, tokenID, nil); err != nil {
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
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_ok", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected happy-path end")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_timeout", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("timeout end must not complete")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "TimerBoundary_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected boundary cancelled")
	}
}

func TestTimerBoundaryRecoverThenFireDue(t *testing.T) {
	xml := readTestdata(t, "m2_timer_boundary.bpmn")
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
	_, tokenID := waitingAt(inst)
	tok := inst.Tokens[tokenID]
	if tok.DueUnixMs == 0 || tok.BoundaryID != "TimerBoundary_1" {
		t.Fatalf("due after recover token=%#v", tok)
	}
	if err := eng2.FireDue(ctx); err != nil {
		t.Fatalf("FireDue: %v", err)
	}
	inst, _ = eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestDeployRejectsNonInterruptingTimerBoundary(t *testing.T) {
	xml := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<definitions xmlns="http://www.omg.org/spec/BPMN/20100524/MODEL" id="Definitions_b" targetNamespace="http://sparrow.example/b">
  <process id="Process_b" isExecutable="true">
    <startEvent id="StartEvent_1"><outgoing>Flow_1</outgoing></startEvent>
    <userTask id="UserTask_1"><incoming>Flow_1</incoming><outgoing>Flow_2</outgoing></userTask>
    <boundaryEvent id="TimerBoundary_1" attachedToRef="UserTask_1" cancelActivity="false">
      <outgoing>Flow_3</outgoing>
      <timerEventDefinition id="TimerDef_1">
        <timeDuration xsi:type="tFormalExpression">PT0S</timeDuration>
      </timerEventDefinition>
    </boundaryEvent>
    <endEvent id="EndEvent_1"><incoming>Flow_2</incoming></endEvent>
    <endEvent id="EndEvent_2"><incoming>Flow_3</incoming></endEvent>
    <sequenceFlow id="Flow_1" sourceRef="StartEvent_1" targetRef="UserTask_1"/>
    <sequenceFlow id="Flow_2" sourceRef="UserTask_1" targetRef="EndEvent_1"/>
    <sequenceFlow id="Flow_3" sourceRef="TimerBoundary_1" targetRef="EndEvent_2"/>
  </process>
</definitions>`)
	eng := processing.NewEngine(eventlog.NewMemory())
	_, err := eng.Deploy(context.Background(), xml)
	if err == nil || !strings.Contains(err.Error(), "UNSUPPORTED_ELEMENT") {
		t.Fatalf("err=%v", err)
	}
}

func TestMessageBoundaryPublishInterrupts(t *testing.T) {
	xml := readTestdata(t, "m2_message_boundary.bpmn")
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
	if elementID != "UserTask_1" {
		t.Fatalf("tokens=%#v", inst.Tokens)
	}
	tok := inst.Tokens[tokenID]
	if tok.MessageName != "order.cancelled" || tok.BoundaryID != "MessageBoundary_1" {
		t.Fatalf("armed message token=%#v", tok)
	}
	if tok.DueUnixMs != 0 {
		t.Fatalf("message boundary due should be 0, due=%d", tok.DueUnixMs)
	}
	if err := eng.FireDue(ctx); err != nil {
		t.Fatalf("FireDue: %v", err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusActive {
		t.Fatalf("FireDue must not complete message boundary, status=%s", inst.Status)
	}

	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "order.cancelled"})
	if err != nil || n != 1 {
		t.Fatalf("publish n=%d err=%v", n, err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "UserTask_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected UserTask TERMINATED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "MessageBoundary_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected boundary COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_msg", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected message end")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_ok", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("happy-path end must not complete")
	}
}

func TestMessageBoundaryCompleteCancelsWait(t *testing.T) {
	xml := readTestdata(t, "m2_message_boundary.bpmn")
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
	if err := eng.Complete(ctx, instanceID, elementID, tokenID, nil); err != nil {
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
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_ok", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected happy-path end")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_msg", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("message end must not complete")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "MessageBoundary_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected boundary cancelled")
	}
}

func TestMessageBoundaryRecoverThenPublish(t *testing.T) {
	xml := readTestdata(t, "m2_message_boundary.bpmn")
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
	_, tokenID := waitingAt(inst)
	tok := inst.Tokens[tokenID]
	if tok.MessageName != "order.cancelled" || tok.BoundaryID != "MessageBoundary_1" {
		t.Fatalf("message after recover token=%#v", tok)
	}
	n, err := eng2.PublishMessage(ctx, processing.PublishMessageRequest{Name: "order.cancelled"})
	if err != nil || n != 1 {
		t.Fatalf("publish n=%d err=%v", n, err)
	}
	inst, _ = eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestMessageBoundaryBuffersUntilWait(t *testing.T) {
	xml := readTestdata(t, "m2_message_boundary.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "order.cancelled"})
	if err != nil || n != 0 {
		t.Fatalf("buffer n=%d err=%v", n, err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed from buffer", inst.Status)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_msg", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected message end from buffer")
	}
}

func TestDeployRejectsNonInterruptingMessageBoundary(t *testing.T) {
	xml := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<definitions xmlns="http://www.omg.org/spec/BPMN/20100524/MODEL" id="Definitions_mb" targetNamespace="http://sparrow.example/mb">
  <message id="message_1" name="order.cancelled"/>
  <process id="Process_mb" isExecutable="true">
    <startEvent id="StartEvent_1"><outgoing>Flow_1</outgoing></startEvent>
    <userTask id="UserTask_1"><incoming>Flow_1</incoming><outgoing>Flow_2</outgoing></userTask>
    <boundaryEvent id="MessageBoundary_1" attachedToRef="UserTask_1" cancelActivity="false">
      <outgoing>Flow_3</outgoing>
      <messageEventDefinition id="MsgDef_1" messageRef="message_1"/>
    </boundaryEvent>
    <endEvent id="EndEvent_1"><incoming>Flow_2</incoming></endEvent>
    <endEvent id="EndEvent_2"><incoming>Flow_3</incoming></endEvent>
    <sequenceFlow id="Flow_1" sourceRef="StartEvent_1" targetRef="UserTask_1"/>
    <sequenceFlow id="Flow_2" sourceRef="UserTask_1" targetRef="EndEvent_1"/>
    <sequenceFlow id="Flow_3" sourceRef="MessageBoundary_1" targetRef="EndEvent_2"/>
  </process>
</definitions>`)
	eng := processing.NewEngine(eventlog.NewMemory())
	_, err := eng.Deploy(context.Background(), xml)
	if err == nil || !strings.Contains(err.Error(), "UNSUPPORTED_ELEMENT") {
		t.Fatalf("err=%v", err)
	}
}

func sawElementIntent(events []*eventv1.Event, typ eventv1.Element_Type, id string, intent eventv1.Element_Intent) bool {
	for _, ev := range events {
		if ev.GetRecordType() != eventv1.Event_RECORD_TYPE_EVENT {
			continue
		}
		el := ev.GetElement()
		if el != nil && el.GetType() == typ && el.GetId() == id && el.GetIntent() == intent {
			return true
		}
	}
	return false
}
