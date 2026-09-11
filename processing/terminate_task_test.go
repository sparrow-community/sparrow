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

func TestTerminateEndCancelsParallelWait(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m18_terminate_parallel.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, ok := eng.GetInstance(instanceID)
	if !ok || inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%v ok=%v", inst, ok)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "UserTask_a", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected UserTask_a TERMINATED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "TerminateEnd_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected TerminateEnd_1 COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_PROCESS, "Process_terminate_parallel", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected PROCESS COMPLETED")
	}
}

func TestTerminateEndInsideSubProcess(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m18_terminate_subprocess.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, ok := eng.GetInstance(instanceID)
	if !ok || inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%v ok=%v", inst, ok)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "UserTask_sp", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected UserTask_sp TERMINATED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_SUB_PROCESS, "SubProcess_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected SubProcess COMPLETED")
	}
}

func TestManualTaskWaitAndComplete(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m18_manual_task.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	elementID, tokenID := waitingAt(inst)
	if elementID != "ManualTask_1" {
		t.Fatalf("wait at %s tokens=%#v", elementID, inst.Tokens)
	}
	if err := eng.Complete(ctx, instanceID, elementID, tokenID, map[string]any{"done": true}); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s", inst.Status)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_MANUAL_TASK, "ManualTask_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected MANUAL_TASK COMPLETED")
	}
}

func TestReceiveTaskPublishMessage(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m18_receive_task.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	elementID, tokenID := waitingAt(inst)
	if elementID != "ReceiveTask_1" {
		t.Fatalf("wait at %s", elementID)
	}
	if inst.Tokens[tokenID].MessageName != "order.confirmed" {
		t.Fatalf("message=%q", inst.Tokens[tokenID].MessageName)
	}
	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "order.confirmed"})
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("delivered=%d", n)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s", inst.Status)
	}
}

func TestReceiveTaskInstantiateRejected(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	_, err := eng.Deploy(context.Background(), readTestdata(t, "m18_receive_instantiate.bpmn"))
	if err == nil || !strings.Contains(err.Error(), "UNSUPPORTED_ELEMENT") {
		t.Fatalf("err=%v", err)
	}
}

func TestSendTaskCompletesAndWakesReceive(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	recvDep, err := eng.Deploy(ctx, readTestdata(t, "m18_receive_task.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	sendDep, err := eng.Deploy(ctx, readTestdata(t, "m18_send_task.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	recvID, err := eng.CreateInstance(ctx, recvDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	sendID, err := eng.CreateInstance(ctx, sendDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	sendInst, ok := eng.GetInstance(sendID)
	if !ok || sendInst.Status != projection.StatusCompleted {
		t.Fatalf("send inst=%v ok=%v", sendInst, ok)
	}
	recvInst, ok := eng.GetInstance(recvID)
	if !ok || recvInst.Status != projection.StatusCompleted {
		t.Fatalf("receive should complete after send, inst=%v ok=%v", recvInst, ok)
	}
	events, err := eng.ListEvents(ctx, sendID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_SEND_TASK, "SendTask_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected SEND_TASK COMPLETED")
	}
}

func TestBusinessRuleTaskActivateAndComplete(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m18_business_rule_task.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	jobs, err := eng.Activate(ctx, processing.ActivateRequest{JobType: "decide.v1", WorkerID: "w1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs) != 1 {
		t.Fatalf("jobs=%d", len(jobs))
	}
	if err := eng.Complete(ctx, jobs[0].ProcessInstanceID, jobs[0].ElementID, jobs[0].TokenID, map[string]any{"result": "ok"}); err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s", inst.Status)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BUSINESS_RULE_TASK, "BusinessRuleTask_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected BUSINESS_RULE_TASK COMPLETED")
	}
}

func TestRecoverManualReceiveBusinessRule(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	eng1, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	defer eng1.Close()

	manDep, err := eng1.Deploy(ctx, readTestdata(t, "m18_manual_task.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	recvDep, err := eng1.Deploy(ctx, readTestdata(t, "m18_receive_task.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	brDep, err := eng1.Deploy(ctx, readTestdata(t, "m18_business_rule_task.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	manID, err := eng1.CreateInstance(ctx, manDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	recvID, err := eng1.CreateInstance(ctx, recvDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	brID, err := eng1.CreateInstance(ctx, brDep, nil)
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

	man, _ := eng2.GetInstance(manID)
	el, tok := waitingAt(man)
	if el != "ManualTask_1" {
		t.Fatalf("manual wait=%s", el)
	}
	if err := eng2.Complete(ctx, manID, el, tok, nil); err != nil {
		t.Fatal(err)
	}
	if n, err := eng2.PublishMessage(ctx, processing.PublishMessageRequest{Name: "order.confirmed", ProcessInstanceID: recvID}); err != nil || n != 1 {
		t.Fatalf("publish n=%d err=%v", n, err)
	}
	jobs, err := eng2.Activate(ctx, processing.ActivateRequest{JobType: "decide.v1", WorkerID: "w2"})
	if err != nil || len(jobs) != 1 {
		t.Fatalf("jobs=%d err=%v", len(jobs), err)
	}
	if err := eng2.Complete(ctx, jobs[0].ProcessInstanceID, jobs[0].ElementID, jobs[0].TokenID, nil); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{manID, recvID, brID} {
		inst, ok := eng2.GetInstance(id)
		if !ok || inst.Status != projection.StatusCompleted {
			t.Fatalf("%s status=%v ok=%v", id, inst, ok)
		}
	}
}

func TestAbstractTaskStillRejected(t *testing.T) {
	body := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<definitions xmlns="http://www.omg.org/spec/BPMN/20100524/MODEL" id="d" targetNamespace="http://sparrow.example/abstract">
  <process id="Process_abs" isExecutable="true">
    <startEvent id="StartEvent_1"><outgoing>Flow_1</outgoing></startEvent>
    <task id="Task_1"><incoming>Flow_1</incoming><outgoing>Flow_2</outgoing></task>
    <endEvent id="EndEvent_1"><incoming>Flow_2</incoming></endEvent>
    <sequenceFlow id="Flow_1" sourceRef="StartEvent_1" targetRef="Task_1"/>
    <sequenceFlow id="Flow_2" sourceRef="Task_1" targetRef="EndEvent_1"/>
  </process>
</definitions>`)
	eng := processing.NewEngine(eventlog.NewMemory())
	_, err := eng.Deploy(context.Background(), body)
	if err == nil || !strings.Contains(err.Error(), "UNSUPPORTED_ELEMENT") {
		t.Fatalf("err=%v", err)
	}
}
