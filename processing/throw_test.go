package processing_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func TestNoneThrowCompletesInstance(t *testing.T) {
	xml := readTestdataThrow(t, "m3_none_throw.bpmn")
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
	if !ok || inst.Status != projection.StatusCompleted {
		t.Fatalf("inst=%v ok=%v", inst, ok)
	}

	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	var sawThrow bool
	for _, ev := range events {
		if ev.GetRecordType() != eventv1.Event_RECORD_TYPE_EVENT {
			continue
		}
		el := ev.GetElement()
		if el != nil && el.GetType() == eventv1.Element_TYPE_INTERMEDIATE_THROW_EVENT &&
			el.GetIntent() == eventv1.Element_INTENT_COMPLETED && el.GetId() == "NoneThrow_1" {
			sawThrow = true
		}
	}
	if !sawThrow {
		t.Fatal("missing NoneThrow_1 COMPLETED")
	}
}

func TestMessageThrowCompletesAlone(t *testing.T) {
	xml := readTestdataThrow(t, "m3_message_throw.bpmn")
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
	if !ok || inst.Status != projection.StatusCompleted {
		t.Fatalf("inst=%v ok=%v", inst, ok)
	}

	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	var sawActivated bool
	for _, ev := range events {
		el := ev.GetElement()
		if el == nil || el.GetType() != eventv1.Element_TYPE_INTERMEDIATE_THROW_EVENT {
			continue
		}
		if el.GetIntent() == eventv1.Element_INTENT_ACTIVATED && el.GetId() == "MessageThrow_1" {
			sawActivated = true
			p := el.GetEventPayload()
			if p == nil || p.GetMessageName() != "order.confirmed" {
				t.Fatalf("payload=%v", p)
			}
		}
	}
	if !sawActivated {
		t.Fatal("missing MessageThrow_1 ACTIVATED")
	}
}

func TestMessageThrowWakesCatchAcrossInstances(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()

	catchDep, err := eng.Deploy(ctx, readTestdataThrow(t, "m2_message_catch.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	throwDep, err := eng.Deploy(ctx, readTestdataThrow(t, "m3_message_throw.bpmn"))
	if err != nil {
		t.Fatal(err)
	}

	catchID, err := eng.CreateInstance(ctx, catchDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, ok := eng.GetInstance(catchID)
	if !ok {
		t.Fatal("missing catch instance")
	}
	if waitingElement(inst) != "MessageCatch_1" {
		t.Fatalf("expected MessageCatch_1, got %s", waitingElement(inst))
	}

	throwID, err := eng.CreateInstance(ctx, throwDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	throwInst, ok := eng.GetInstance(throwID)
	if !ok || throwInst.Status != projection.StatusCompleted {
		t.Fatalf("throw inst=%v ok=%v", throwInst, ok)
	}

	catchInst, ok := eng.GetInstance(catchID)
	if !ok || catchInst.Status != projection.StatusCompleted {
		t.Fatalf("catch should complete after throw, inst=%v ok=%v", catchInst, ok)
	}
}

func TestParallelMessageThrowWakesSiblingCatch(t *testing.T) {
	xml := readTestdataThrow(t, "m3_parallel_message_throw_catch.bpmn")
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
	if !ok {
		t.Fatal("missing instance")
	}
	var taskToken string
	var catchWaiting bool
	for _, tok := range inst.Tokens {
		if tok == nil || tok.Status != projection.TokenWaiting {
			continue
		}
		switch tok.ElementID {
		case "UserTask_1":
			taskToken = tok.ID
		case "MessageCatch_1":
			catchWaiting = true
		}
	}
	if taskToken == "" || !catchWaiting {
		t.Fatalf("expected UserTask_1 and MessageCatch_1 waiting, tokens=%#v", inst.Tokens)
	}

	if err := eng.Complete(ctx, instanceID, "UserTask_1", taskToken, nil); err != nil {
		t.Fatal(err)
	}
	inst, ok = eng.GetInstance(instanceID)
	if !ok || inst.Status != projection.StatusCompleted {
		t.Fatalf("expected completed after throw wakes catch, inst=%v ok=%v", inst, ok)
	}
}

func TestSignalCatchWaitAndPublish(t *testing.T) {
	xml := readTestdataThrow(t, "m3_signal_catch.bpmn")
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
	if !ok {
		t.Fatal("missing instance")
	}
	elementID, tokenID := waitingAt(inst)
	if elementID != "SignalCatch_1" {
		t.Fatalf("expected SignalCatch_1, got %s", elementID)
	}
	tok := inst.Tokens[tokenID]
	if tok == nil || tok.SignalName != "go.ahead" {
		t.Fatalf("signal name=%v", tok)
	}
	if tok.MessageName != "" {
		t.Fatalf("message name should be empty, got %q", tok.MessageName)
	}

	n, err := eng.PublishSignal(ctx, processing.PublishSignalRequest{Name: "go.ahead"})
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("delivered=%d", n)
	}
	inst, ok = eng.GetInstance(instanceID)
	if !ok || inst.Status != projection.StatusCompleted {
		t.Fatalf("inst=%v ok=%v", inst, ok)
	}
}

func TestSignalThrowWakesCatchAcrossInstances(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()

	catchDep, err := eng.Deploy(ctx, readTestdataThrow(t, "m3_signal_catch.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	throwDep, err := eng.Deploy(ctx, readTestdataThrow(t, "m3_signal_throw.bpmn"))
	if err != nil {
		t.Fatal(err)
	}

	catchID, err := eng.CreateInstance(ctx, catchDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	if waitingElement(mustInstance(t, eng, catchID)) != "SignalCatch_1" {
		t.Fatal("expected signal catch waiting")
	}

	throwID, err := eng.CreateInstance(ctx, throwDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	if mustInstance(t, eng, throwID).Status != projection.StatusCompleted {
		t.Fatal("throw should complete")
	}
	if mustInstance(t, eng, catchID).Status != projection.StatusCompleted {
		t.Fatal("catch should complete after signal throw")
	}
}

func TestPublishSignalNotBuffered(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()

	n, err := eng.PublishSignal(ctx, processing.PublishSignalRequest{Name: "go.ahead"})
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("delivered=%d", n)
	}

	dep, err := eng.Deploy(ctx, readTestdataThrow(t, "m3_signal_catch.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	if waitingElement(mustInstance(t, eng, instanceID)) != "SignalCatch_1" {
		t.Fatal("late signal must not wake future catch")
	}
}

func mustInstance(t *testing.T, eng *processing.Engine, id string) *projection.Instance {
	t.Helper()
	inst, ok := eng.GetInstance(id)
	if !ok || inst == nil {
		t.Fatalf("missing instance %q", id)
	}
	return inst
}

func readTestdataThrow(t *testing.T, name string) []byte {
	t.Helper()
	xml, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return xml
}
