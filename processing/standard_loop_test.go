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

func TestStandardLoopDoWhile(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m24_standard_loop_dowhile.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	id, err := eng.CreateInstance(ctx, dep, map[string]any{"count": 0})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		inst, _ := eng.GetInstance(id)
		el, tok := waitingAt(inst)
		if el != "UserTask_1" {
			t.Fatalf("iter %d wait=%s", i, el)
		}
		if err := eng.Complete(ctx, id, el, tok, map[string]any{"count": i + 1}); err != nil {
			t.Fatal(err)
		}
	}
	inst, _ := eng.GetInstance(id)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s tokens=%#v", inst.Status, inst.Tokens)
	}
	if n := countUserTaskActivated(mustEvents(t, eng, id), "UserTask_1"); n != 3 {
		t.Fatalf("activated=%d want 3", n)
	}
}

func countUserTaskActivated(events []*eventv1.Event, id string) int {
	n := 0
	for _, ev := range events {
		if ev.GetRecordType() != eventv1.Event_RECORD_TYPE_EVENT {
			continue
		}
		el := ev.GetElement()
		if el == nil {
			continue
		}
		if el.GetType() == eventv1.Element_TYPE_USER_TASK && el.GetId() == id && el.GetIntent() == eventv1.Element_INTENT_ACTIVATED {
			n++
		}
	}
	return n
}

func TestStandardLoopTestBeforeSkip(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m24_standard_loop_test_before.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	id, err := eng.CreateInstance(ctx, dep, map[string]any{"run": false})
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(id)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s tokens=%#v", inst.Status, inst.Tokens)
	}
	for _, tok := range inst.Tokens {
		if tok != nil && tok.Status == projection.TokenWaiting {
			t.Fatalf("unexpected wait at %s", tok.ElementID)
		}
	}
}

func TestStandardLoopMaximum(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m24_standard_loop_max.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	id, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		inst, _ := eng.GetInstance(id)
		el, tok := waitingAt(inst)
		if err := eng.Complete(ctx, id, el, tok, nil); err != nil {
			t.Fatal(err)
		}
	}
	inst, _ := eng.GetInstance(id)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s", inst.Status)
	}
	if n := countUserTaskActivated(mustEvents(t, eng, id), "UserTask_1"); n != 2 {
		t.Fatalf("activated=%d want 2", n)
	}
}

func TestStandardLoopRecoverMid(t *testing.T) {
	dir := t.TempDir()
	ctx := context.Background()
	eng1, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	dep, err := eng1.Deploy(ctx, readTestdata(t, "m24_standard_loop_dowhile.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	id, err := eng1.CreateInstance(ctx, dep, map[string]any{"count": 0})
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng1.GetInstance(id)
	el, tok := waitingAt(inst)
	if err := eng1.Complete(ctx, id, el, tok, map[string]any{"count": 1}); err != nil {
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
	inst, _ = eng2.GetInstance(id)
	el, tok = waitingAt(inst)
	if el != "UserTask_1" {
		t.Fatalf("wait=%s", el)
	}
	if err := eng2.Complete(ctx, id, el, tok, map[string]any{"count": 2}); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng2.GetInstance(id)
	el, tok = waitingAt(inst)
	if err := eng2.Complete(ctx, id, el, tok, map[string]any{"count": 3}); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng2.GetInstance(id)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s", inst.Status)
	}
}

func TestStandardLoopRejectsWithMultiInstance(t *testing.T) {
	body := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<definitions xmlns="http://www.omg.org/spec/BPMN/20100524/MODEL" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" id="d" targetNamespace="http://sparrow.example/x">
  <process id="P" isExecutable="true">
    <startEvent id="S"><outgoing>F1</outgoing></startEvent>
    <userTask id="U"><incoming>F1</incoming><outgoing>F2</outgoing>
      <standardLoopCharacteristics testBefore="false">
        <loopCondition xsi:type="tFormalExpression">true</loopCondition>
      </standardLoopCharacteristics>
      <multiInstanceLoopCharacteristics>
        <loopCardinality xsi:type="tFormalExpression">2</loopCardinality>
      </multiInstanceLoopCharacteristics>
    </userTask>
    <endEvent id="E"><incoming>F2</incoming></endEvent>
    <sequenceFlow id="F1" sourceRef="S" targetRef="U"/>
    <sequenceFlow id="F2" sourceRef="U" targetRef="E"/>
  </process>
</definitions>`)
	eng := processing.NewEngine(eventlog.NewMemory())
	_, err := eng.Deploy(context.Background(), body)
	if err == nil || !strings.Contains(err.Error(), "UNSUPPORTED_ELEMENT") {
		t.Fatalf("err=%v", err)
	}
}

func mustEvents(t *testing.T, eng *processing.Engine, id string) []*eventv1.Event {
	t.Helper()
	events, err := eng.ListEvents(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}
	return events
}
