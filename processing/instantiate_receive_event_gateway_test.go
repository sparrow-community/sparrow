package processing_test

import (
	"context"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func TestInstantiateReceiveTaskEntry(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m26_instantiate_receive.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	id, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst := mustInstance(t, eng, id)
	if waitingElement(inst) != "ReceiveTask_1" {
		t.Fatalf("wait=%s", waitingElement(inst))
	}
	if n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "order.confirmed", ProcessInstanceID: id}); err != nil || n != 1 {
		t.Fatalf("n=%d err=%v", n, err)
	}
	if mustInstance(t, eng, id).Status != projection.StatusCompleted {
		t.Fatal("expected completed")
	}
}

func TestEventBasedGatewayTargetsReceiveTask(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m26_event_based_to_receive.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	id, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst := mustInstance(t, eng, id)
	waiting := map[string]bool{}
	for _, tok := range inst.Tokens {
		if tok != nil && tok.Status == projection.TokenWaiting {
			waiting[tok.ElementID] = true
		}
	}
	if !waiting["ReceiveTask_1"] || !waiting["MessageCatch_1"] {
		t.Fatalf("waiting=%v", waiting)
	}
	if n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "order.a", ProcessInstanceID: id}); err != nil || n != 1 {
		t.Fatalf("n=%d err=%v", n, err)
	}
	events, _ := eng.ListEvents(ctx, id)
	if !sawElementIntent(events, eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT, "MessageCatch_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected sibling catch TERMINATED")
	}
	if mustInstance(t, eng, id).Status != projection.StatusCompleted {
		t.Fatal("expected completed")
	}
}

func TestMultipleInstantiateEventBasedGateways(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m26_multiple_instantiate_event_based.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	id, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst := mustInstance(t, eng, id)
	waiting := 0
	for _, tok := range inst.Tokens {
		if tok != nil && tok.Status == projection.TokenWaiting {
			waiting++
		}
	}
	if waiting != 4 {
		t.Fatalf("waiting=%d want 4", waiting)
	}
	if n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "path.a", ProcessInstanceID: id}); err != nil || n != 1 {
		t.Fatalf("n=%d err=%v", n, err)
	}
	events, _ := eng.ListEvents(ctx, id)
	for _, catchID := range []string{"Catch_a2", "Catch_b1", "Catch_b2"} {
		if !sawElementIntent(events, eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT, catchID, eventv1.Element_INTENT_TERMINATED) {
			t.Fatalf("expected %s TERMINATED", catchID)
		}
	}
	if mustInstance(t, eng, id).Status != projection.StatusCompleted {
		t.Fatal("expected completed")
	}
}
