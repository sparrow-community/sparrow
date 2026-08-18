package processing_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func TestTimerCatchPT0SFireDue(t *testing.T) {
	xml := readTestdata(t, "m2_timer_catch.bpmn")
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
	inst, ok := eng.GetInstance(instanceID)
	if !ok || inst.Status != projection.StatusActive {
		t.Fatalf("inst=%v ok=%v", inst, ok)
	}
	elementID, tokenID := waitingAt(inst)
	if elementID != "TimerCatch_1" {
		t.Fatalf("expected wait at TimerCatch_1, tokens=%#v", inst.Tokens)
	}
	tok := inst.Tokens[tokenID]
	if tok.DueUnixMs == 0 {
		t.Fatal("expected due_unix_ms on waiting timer token")
	}

	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	var sawDue bool
	for _, ev := range events {
		if ev.GetRecordType() != eventv1.Event_RECORD_TYPE_EVENT {
			continue
		}
		el := ev.GetElement()
		if el != nil && el.GetType() == eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT &&
			el.GetIntent() == eventv1.Element_INTENT_ACTIVATED {
			p := el.GetEventPayload()
			if p == nil || p.GetDueUnixMs() == 0 || p.GetDuration() != "PT0S" {
				t.Fatalf("event payload=%v", p)
			}
			sawDue = true
		}
	}
	if !sawDue {
		t.Fatal("missing INTERMEDIATE_CATCH_EVENT ACTIVATED with due")
	}

	if err := eng.FireDue(ctx); err != nil {
		t.Fatalf("FireDue: %v", err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestTimerCatchNotYetDue(t *testing.T) {
	xml := readTestdata(t, "m2_timer_catch_1h.bpmn")
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
	if inst.Status != projection.StatusActive {
		t.Fatalf("status=%s want still active", inst.Status)
	}
	elementID, tokenID := waitingAt(inst)
	if elementID != "TimerCatch_1" {
		t.Fatalf("tokens=%#v", inst.Tokens)
	}
	if err := eng.Complete(ctx, instanceID, elementID, tokenID, nil); err != nil {
		t.Fatalf("Complete: %v", err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestTimerCatchRecoverThenFireDue(t *testing.T) {
	xml := readTestdata(t, "m2_timer_catch.bpmn")
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
	_, tokenID := waitingAt(inst)
	due := inst.Tokens[tokenID].DueUnixMs
	if due == 0 {
		t.Fatal("due missing before close")
	}
	if err := eng1.Close(); err != nil {
		t.Fatal(err)
	}

	eng2, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	defer eng2.Close()

	inst2, ok := eng2.GetInstance(instanceID)
	if !ok {
		t.Fatal("instance missing after recover")
	}
	elementID, token2 := waitingAt(inst2)
	if elementID != "TimerCatch_1" {
		t.Fatalf("tokens=%#v", inst2.Tokens)
	}
	if inst2.Tokens[token2].DueUnixMs != due {
		t.Fatalf("due after recover=%d want %d", inst2.Tokens[token2].DueUnixMs, due)
	}
	if err := eng2.FireDue(ctx); err != nil {
		t.Fatalf("FireDue: %v", err)
	}
	inst2, _ = eng2.GetInstance(instanceID)
	if inst2.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst2.Status)
	}
}

func TestDeployRejectsNonTimerCatch(t *testing.T) {
	xml := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<definitions xmlns="http://www.omg.org/spec/BPMN/20100524/MODEL"
             id="Definitions_msg"
             targetNamespace="http://sparrow.example/msg">
  <process id="Process_msg" isExecutable="true">
    <startEvent id="StartEvent_1">
      <outgoing>Flow_1</outgoing>
    </startEvent>
    <intermediateCatchEvent id="Catch_1">
      <incoming>Flow_1</incoming>
      <outgoing>Flow_2</outgoing>
      <messageEventDefinition id="MsgDef_1"/>
    </intermediateCatchEvent>
    <endEvent id="EndEvent_1">
      <incoming>Flow_2</incoming>
    </endEvent>
    <sequenceFlow id="Flow_1" sourceRef="StartEvent_1" targetRef="Catch_1"/>
    <sequenceFlow id="Flow_2" sourceRef="Catch_1" targetRef="EndEvent_1"/>
  </process>
</definitions>`)
	eng := processing.NewEngine(eventlog.NewMemory())
	_, err := eng.Deploy(context.Background(), xml)
	if err == nil || !strings.Contains(err.Error(), "UNSUPPORTED_ELEMENT") {
		t.Fatalf("err=%v", err)
	}
}

func readTestdata(t *testing.T, name string) []byte {
	t.Helper()
	xml, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return xml
}
