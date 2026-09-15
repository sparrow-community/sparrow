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

func TestTransactionSuccessPath(t *testing.T) {
	xml := readTestdata(t, "m33_transaction_cancel.bpmn")
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
	elem, tok := waitingAt(inst)
	if elem != "Task_A" {
		t.Fatalf("expected Task_A, got %s", elem)
	}
	if err := eng.Complete(ctx, instanceID, elem, tok, map[string]any{"cancel": false}); err != nil {
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
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_success", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected End_success")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_cancelled", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("cancel end must not complete on success path")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "Undo_A", eventv1.Element_INTENT_ACTIVATED) {
		t.Fatal("Undo_A must not run on success path")
	}
}

func TestTransactionCancelRunsCompensation(t *testing.T) {
	xml := readTestdata(t, "m33_transaction_cancel.bpmn")
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
	elem, tok := waitingAt(inst)
	if err := eng.Complete(ctx, instanceID, elem, tok, map[string]any{"cancel": true}); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(instanceID)
	waiting := allWaiting(inst)
	if !containsStr(waiting, "Undo_A") || !containsStr(waiting, "Cancel_End") {
		t.Fatalf("expected Cancel_End + Undo_A waiting, got %v", waiting)
	}
	undoTok := tokenAtElement(inst, "Undo_A")
	if undoTok == "" {
		t.Fatal("missing Undo_A waiting token")
	}
	if err := eng.Complete(ctx, instanceID, "Undo_A", undoTok, nil); err != nil {
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
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_cancelled", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected End_cancelled")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "Cancel_Boundary", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected Cancel_Boundary completed")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_TRANSACTION, "Transaction_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected Transaction_1 TERMINATED")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_success", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("success end must not complete on cancel path")
	}
}

func TestTransactionCancelRecoverThenCompleteUndo(t *testing.T) {
	xml := readTestdata(t, "m33_transaction_cancel.bpmn")
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
	inst, _ := eng1.GetInstance(instanceID)
	elem, tok := waitingAt(inst)
	if err := eng1.Complete(ctx, instanceID, elem, tok, map[string]any{"cancel": true}); err != nil {
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
	undoTok := tokenAtElement(inst, "Undo_A")
	if undoTok == "" {
		t.Fatalf("expected Undo_A after recover, waiting=%v", allWaiting(inst))
	}
	if err := eng2.Complete(ctx, instanceID, "Undo_A", undoTok, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestTransactionDeployRejectsNested(t *testing.T) {
	xml := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<definitions xmlns="http://www.omg.org/spec/BPMN/20100524/MODEL"
             id="Definitions_nested_tx"
             targetNamespace="http://sparrow.example/nested-tx">
  <process id="Process_nested_tx" isExecutable="true">
    <startEvent id="S"><outgoing>F1</outgoing></startEvent>
    <transaction id="Tx_outer">
      <incoming>F1</incoming>
      <outgoing>F2</outgoing>
      <startEvent id="S_outer"><outgoing>F_inner</outgoing></startEvent>
      <transaction id="Tx_inner">
        <incoming>F_inner</incoming>
        <outgoing>F_inner_end</outgoing>
        <startEvent id="S_inner"><outgoing>F_ie</outgoing></startEvent>
        <endEvent id="E_inner"><incoming>F_ie</incoming></endEvent>
        <sequenceFlow id="F_ie" sourceRef="S_inner" targetRef="E_inner"/>
      </transaction>
      <endEvent id="E_outer"><incoming>F_inner_end</incoming></endEvent>
      <sequenceFlow id="F_inner" sourceRef="S_outer" targetRef="Tx_inner"/>
      <sequenceFlow id="F_inner_end" sourceRef="Tx_inner" targetRef="E_outer"/>
    </transaction>
    <endEvent id="E"><incoming>F2</incoming></endEvent>
    <sequenceFlow id="F1" sourceRef="S" targetRef="Tx_outer"/>
    <sequenceFlow id="F2" sourceRef="Tx_outer" targetRef="E"/>
  </process>
</definitions>`)
	eng := processing.NewEngine(eventlog.NewMemory())
	_, err := eng.Deploy(context.Background(), xml)
	if err == nil || !strings.Contains(err.Error(), "UNSUPPORTED_ELEMENT") {
		t.Fatalf("expected UNSUPPORTED_ELEMENT for nested transaction, got %v", err)
	}
}

func TestTransactionDeployRejectsCancelWithoutBoundary(t *testing.T) {
	xml := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<definitions xmlns="http://www.omg.org/spec/BPMN/20100524/MODEL"
             id="Definitions_tx_noboundary"
             targetNamespace="http://sparrow.example/tx-noboundary">
  <process id="Process_tx_noboundary" isExecutable="true">
    <startEvent id="S"><outgoing>F1</outgoing></startEvent>
    <transaction id="Tx_1">
      <incoming>F1</incoming>
      <outgoing>F2</outgoing>
      <startEvent id="S_tx"><outgoing>Fc</outgoing></startEvent>
      <endEvent id="Cancel_End"><incoming>Fc</incoming>
        <cancelEventDefinition/>
      </endEvent>
      <sequenceFlow id="Fc" sourceRef="S_tx" targetRef="Cancel_End"/>
    </transaction>
    <endEvent id="E"><incoming>F2</incoming></endEvent>
    <sequenceFlow id="F1" sourceRef="S" targetRef="Tx_1"/>
    <sequenceFlow id="F2" sourceRef="Tx_1" targetRef="E"/>
  </process>
</definitions>`)
	eng := processing.NewEngine(eventlog.NewMemory())
	_, err := eng.Deploy(context.Background(), xml)
	if err == nil || !strings.Contains(err.Error(), "UNSUPPORTED_ELEMENT") {
		t.Fatalf("expected UNSUPPORTED_ELEMENT, got %v", err)
	}
}
