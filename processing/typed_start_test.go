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

func TestMessageStartCreateInstanceRejected(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m19_message_start.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = eng.CreateInstance(ctx, dep, nil)
	if err == nil || !strings.Contains(err.Error(), "INVALID_ARGUMENT") {
		t.Fatalf("CreateInstance err=%v want INVALID_ARGUMENT", err)
	}
}

func TestMessageStartPublishCreatesInstance(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	if _, err := eng.Deploy(ctx, readTestdata(t, "m19_message_start.bpmn")); err != nil {
		t.Fatal(err)
	}
	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "order.created"})
	if err != nil || n != 1 {
		t.Fatalf("PublishMessage n=%d err=%v", n, err)
	}
	ids := eng.ListInstanceIDs()
	if len(ids) != 1 {
		t.Fatalf("instances=%v", ids)
	}
	inst, ok := eng.GetInstance(ids[0])
	if !ok || inst.Status != projection.StatusActive {
		t.Fatalf("status=%v ok=%v", inst, ok)
	}
	el, tok := waitingAt(inst)
	if el != "UserTask_1" {
		t.Fatalf("wait=%q want UserTask_1", el)
	}
	if err := eng.Complete(ctx, ids[0], el, tok, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(ids[0])
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%v", inst.Status)
	}
}

func TestTimerStartFireDueCreatesInstance(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m19_timer_start.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := eng.CreateInstance(ctx, dep, nil); err == nil {
		t.Fatal("expected CreateInstance reject")
	}
	if err := eng.FireDue(ctx); err != nil {
		t.Fatal(err)
	}
	ids := eng.ListInstanceIDs()
	if len(ids) != 1 {
		t.Fatalf("instances=%v", ids)
	}
	inst, _ := eng.GetInstance(ids[0])
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%v", inst.Status)
	}
}

func TestSignalStartPublishCreatesInstance(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m19_signal_start.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := eng.CreateInstance(ctx, dep, nil); err == nil {
		t.Fatal("expected CreateInstance reject")
	}
	n, err := eng.PublishSignal(ctx, processing.PublishSignalRequest{Name: "go.ahead"})
	if err != nil || n != 1 {
		t.Fatalf("PublishSignal n=%d err=%v", n, err)
	}
	ids := eng.ListInstanceIDs()
	if len(ids) != 1 {
		t.Fatalf("instances=%v", ids)
	}
	inst, _ := eng.GetInstance(ids[0])
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%v", inst.Status)
	}
}

func TestConditionalStartEvaluate(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m19_conditional_start.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := eng.CreateInstance(ctx, dep, nil); err == nil {
		t.Fatal("expected CreateInstance reject")
	}
	n, err := eng.EvaluateConditionalStarts(ctx, processing.EvaluateConditionalStartsRequest{
		DeploymentID: dep,
		Variables:    map[string]any{"approved": false},
	})
	if err != nil || n != 0 {
		t.Fatalf("false eval n=%d err=%v", n, err)
	}
	if len(eng.ListInstanceIDs()) != 0 {
		t.Fatal("expected no instance")
	}
	n, err = eng.EvaluateConditionalStarts(ctx, processing.EvaluateConditionalStartsRequest{
		DeploymentID: dep,
		Variables:    map[string]any{"approved": true},
	})
	if err != nil || n != 1 {
		t.Fatalf("true eval n=%d err=%v", n, err)
	}
	inst, _ := eng.GetInstance(eng.ListInstanceIDs()[0])
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%v", inst.Status)
	}
}

func TestNoneAndMessageAlternativeStarts(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m19_none_and_message_start.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	noneID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, ok := eng.GetInstance(noneID)
	if !ok || inst.Status != projection.StatusCompleted {
		t.Fatalf("none start status=%v ok=%v", inst, ok)
	}
	events, err := eng.ListEvents(ctx, noneID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_START_EVENT, "StartEvent_none", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected StartEvent_none COMPLETED")
	}

	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "order.created"})
	if err != nil || n != 1 {
		t.Fatalf("PublishMessage n=%d err=%v", n, err)
	}
	var msgID string
	for _, id := range eng.ListInstanceIDs() {
		if id == noneID {
			continue
		}
		msgID = id
	}
	if msgID == "" {
		t.Fatal("expected message-start instance")
	}
	inst, _ = eng.GetInstance(msgID)
	el, _ := waitingAt(inst)
	if el != "UserTask_1" {
		t.Fatalf("wait=%q", el)
	}
}

func TestProcessLevelErrorStartRejected(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	_, err := eng.Deploy(context.Background(), readTestdata(t, "m19_error_start_reject.bpmn"))
	if err == nil || !strings.Contains(err.Error(), "UNSUPPORTED_ELEMENT") {
		t.Fatalf("Deploy err=%v want UNSUPPORTED_ELEMENT", err)
	}
}

func TestMessageStartRecoverThenComplete(t *testing.T) {
	dir := t.TempDir()
	ctx := context.Background()
	eng1, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := eng1.Deploy(ctx, readTestdata(t, "m19_message_start.bpmn")); err != nil {
		t.Fatal(err)
	}
	n, err := eng1.PublishMessage(ctx, processing.PublishMessageRequest{Name: "order.created"})
	if err != nil || n != 1 {
		t.Fatalf("PublishMessage n=%d err=%v", n, err)
	}
	instanceID := eng1.ListInstanceIDs()[0]
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
		t.Fatal("missing instance after recover")
	}
	el, tok := waitingAt(inst)
	if err := eng2.Complete(ctx, instanceID, el, tok, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%v", inst.Status)
	}
}

func TestTimerStartRecoverThenFireDue(t *testing.T) {
	dir := t.TempDir()
	ctx := context.Background()
	eng1, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := eng1.Deploy(ctx, readTestdata(t, "m19_timer_start.bpmn")); err != nil {
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
	if err := eng2.FireDue(ctx); err != nil {
		t.Fatal(err)
	}
	ids := eng2.ListInstanceIDs()
	if len(ids) != 1 {
		t.Fatalf("instances=%v", ids)
	}
	inst, _ := eng2.GetInstance(ids[0])
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%v", inst.Status)
	}
}

func TestSignalStartRecoverThenComplete(t *testing.T) {
	dir := t.TempDir()
	ctx := context.Background()
	eng1, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := eng1.Deploy(ctx, readTestdata(t, "m19_signal_start.bpmn")); err != nil {
		t.Fatal(err)
	}
	n, err := eng1.PublishSignal(ctx, processing.PublishSignalRequest{Name: "go.ahead"})
	if err != nil || n != 1 {
		t.Fatalf("PublishSignal n=%d err=%v", n, err)
	}
	instanceID := eng1.ListInstanceIDs()[0]
	if err := eng1.Close(); err != nil {
		t.Fatal(err)
	}

	eng2, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	defer eng2.Close()
	inst, ok := eng2.GetInstance(instanceID)
	if !ok || inst.Status != projection.StatusCompleted {
		t.Fatalf("after recover status=%v ok=%v", inst, ok)
	}
}

func TestInstantiateEBGStillNeedsCreateInstance(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m11_instantiate_ebg_timer_message.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "escalate"})
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("PublishMessage created %d without CreateInstance", n)
	}
	if len(eng.ListInstanceIDs()) != 0 {
		t.Fatal("instantiate EBG must not auto-create")
	}
	if _, err := eng.CreateInstance(ctx, dep, nil); err != nil {
		t.Fatal(err)
	}
}
