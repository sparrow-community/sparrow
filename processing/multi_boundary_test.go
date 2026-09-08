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

func TestMultiMessageBoundaryArmsAllThree(t *testing.T) {
	xml := readTestdata(t, "m14_multi_message_boundary.bpmn")
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
	if len(tok.BoundaryWaits) != 3 {
		t.Fatalf("want 3 BoundaryWaits, got %#v", tok.BoundaryWaits)
	}
	names := map[string]bool{}
	for _, w := range tok.BoundaryWaits {
		if w.Kind != "message" {
			t.Fatalf("kind=%q wait=%#v", w.Kind, w)
		}
		names[w.MessageName] = true
	}
	for _, want := range []string{"cancel.a", "cancel.b", "cancel.c"} {
		if !names[want] {
			t.Fatalf("missing wait %q in %#v", want, tok.BoundaryWaits)
		}
	}

	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "cancel.b"})
	if err != nil || n != 1 {
		t.Fatalf("publish cancel.b n=%d err=%v", n, err)
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
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "MessageBoundary_B", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected MessageBoundary_B COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "MessageBoundary_A", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected sibling MessageBoundary_A TERMINATED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "MessageBoundary_C", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected sibling MessageBoundary_C TERMINATED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_b", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected End_b")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_ok", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("happy-path end must not complete")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_a", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("End_a must not complete")
	}
}

func TestMultiTimerBoundaryShortWins(t *testing.T) {
	xml := readTestdata(t, "m14_multi_timer_boundary.bpmn")
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
	_, tokenID := waitingAt(inst)
	tok := inst.Tokens[tokenID]
	if len(tok.BoundaryWaits) != 2 {
		t.Fatalf("want 2 timer waits, got %#v", tok.BoundaryWaits)
	}
	var shortDue, longDue int64
	for _, w := range tok.BoundaryWaits {
		switch w.BoundaryID {
		case "TimerBoundary_short":
			shortDue = w.DueUnixMs
		case "TimerBoundary_long":
			longDue = w.DueUnixMs
		}
	}
	if shortDue == 0 || longDue == 0 || longDue <= shortDue {
		t.Fatalf("timer dues short=%d long=%d waits=%#v", shortDue, longDue, tok.BoundaryWaits)
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
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "TimerBoundary_short", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected short timer COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "TimerBoundary_long", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected long timer sibling TERMINATED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_short", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected End_short")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_long", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("End_long must not complete")
	}
}

func TestMultiMessageBoundaryRecoverThenPublish(t *testing.T) {
	xml := readTestdata(t, "m14_multi_message_boundary.bpmn")
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
	if len(tok.BoundaryWaits) != 3 {
		t.Fatalf("after recover want 3 waits, got %#v", tok.BoundaryWaits)
	}
	n, err := eng2.PublishMessage(ctx, processing.PublishMessageRequest{Name: "cancel.c"})
	if err != nil || n != 1 {
		t.Fatalf("publish cancel.c n=%d err=%v", n, err)
	}
	inst, _ = eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, err := eng2.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_c", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected End_c after recover publish")
	}
}

func TestMultiMessageBoundaryDuplicateNameRejected(t *testing.T) {
	xml := readTestdata(t, "m14_multi_message_dup.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	_, err := eng.Deploy(context.Background(), xml)
	if err == nil {
		t.Fatal("expected duplicate message-name deploy error")
	}
	if !strings.Contains(err.Error(), "already has message") {
		t.Fatalf("err=%v", err)
	}
}

func TestMultiMessageBoundarySameNameDifferentRefsRejected(t *testing.T) {
	xml := readTestdata(t, "m14_multi_message_same_name.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	_, err := eng.Deploy(context.Background(), xml)
	if err == nil {
		t.Fatal("expected same message-name deploy error")
	}
	if !strings.Contains(err.Error(), "already has message") {
		t.Fatalf("err=%v", err)
	}
}
