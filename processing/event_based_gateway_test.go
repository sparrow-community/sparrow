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

func TestEventBasedGatewayDeploy(t *testing.T) {
	xml := readTestdata(t, "m2_event_based_gateway.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	if _, err := eng.Deploy(context.Background(), xml); err != nil {
		t.Fatal(err)
	}
}

func TestEventBasedGatewayTimerWins(t *testing.T) {
	xml := readTestdata(t, "m2_event_based_gateway.bpmn")
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
	waiting := allWaiting(inst)
	if len(waiting) != 2 {
		t.Fatalf("expected 2 waiting catches, got %v", waiting)
	}

	if err := eng.FireDue(ctx); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT, "TimerCatch_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected timer catch COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT, "MessageCatch_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected message catch TERMINATED as sibling")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_timer", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected End_timer")
	}
}

func TestEventBasedGatewayMessageWins(t *testing.T) {
	xml := readTestdata(t, "m2_event_based_gateway.bpmn")
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

	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "escalate"})
	if err != nil || n != 1 {
		t.Fatalf("PublishMessage n=%d err=%v", n, err)
	}
	inst, _ := eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT, "MessageCatch_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected message catch COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT, "TimerCatch_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected timer catch TERMINATED as sibling")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_msg", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected End_msg")
	}
}

func TestEventBasedGatewayRecoverThenMessageWins(t *testing.T) {
	xml := readTestdata(t, "m2_event_based_gateway.bpmn")
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
	if len(allWaiting(inst)) != 2 {
		t.Fatalf("expected 2 waiting catches after recover, got %v", allWaiting(inst))
	}
	n, err := eng2.PublishMessage(ctx, processing.PublishMessageRequest{Name: "escalate"})
	if err != nil || n != 1 {
		t.Fatalf("PublishMessage n=%d err=%v", n, err)
	}
	inst, _ = eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestParallelEventBasedGatewayKeepsSiblings(t *testing.T) {
	xml := readTestdata(t, "m3_parallel_event_based_gateway.bpmn")
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
	waiting := allWaiting(inst)
	if len(waiting) != 2 {
		t.Fatalf("expected 2 waiting catches, got %v", waiting)
	}

	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "event.a"})
	if err != nil || n != 1 {
		t.Fatalf("PublishMessage A n=%d err=%v", n, err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusActive {
		t.Fatalf("status=%s want active after first event", inst.Status)
	}
	waiting = allWaiting(inst)
	if len(waiting) != 1 || waiting[0] != "MessageCatch_B" {
		t.Fatalf("expected only MessageCatch_B waiting, got %v", waiting)
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if sawElementIntent(events, eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT, "MessageCatch_B", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("parallel EBG must not terminate sibling catch")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_A", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected End_A after event.a")
	}

	n, err = eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "event.b"})
	if err != nil || n != 1 {
		t.Fatalf("PublishMessage B n=%d err=%v", n, err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed after both events", inst.Status)
	}
	events, _ = eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT, "MessageCatch_B", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected MessageCatch_B COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_B", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected End_B")
	}
}

func TestInstantiateEventBasedGatewayDeploy(t *testing.T) {
	xml := readTestdata(t, "m11_instantiate_ebg_timer_message.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	if _, err := eng.Deploy(context.Background(), xml); err != nil {
		t.Fatal(err)
	}
}

func TestInstantiateEventBasedGatewayTimerWins(t *testing.T) {
	xml := readTestdata(t, "m11_instantiate_ebg_timer_message.bpmn")
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
	waiting := allWaiting(inst)
	if len(waiting) != 2 {
		t.Fatalf("expected 2 waiting catches, got %v", waiting)
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_EVENT_BASED_GATEWAY, "EBG_inst", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected instantiate EBG completed as entry")
	}

	if err := eng.FireDue(ctx); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, _ = eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT, "TimerCatch_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected timer catch COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT, "MessageCatch_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected message catch TERMINATED as sibling")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_timer", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected End_timer")
	}
}

func TestInstantiateEventBasedGatewayMessageWins(t *testing.T) {
	xml := readTestdata(t, "m11_instantiate_ebg_timer_message.bpmn")
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

	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "escalate"})
	if err != nil || n != 1 {
		t.Fatalf("PublishMessage n=%d err=%v", n, err)
	}
	inst, _ := eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT, "MessageCatch_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected message catch COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT, "TimerCatch_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected timer catch TERMINATED as sibling")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_msg", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected End_msg")
	}
}

func TestInstantiateEventBasedGatewayRejectInvalid(t *testing.T) {
	cases := []struct {
		file string
		sub  string
	}{
		{"m11_instantiate_ebg_incoming.bpmn", "incoming"},
		{"m11_instantiate_ebg_parallel.bpmn", "parallel"},
		{"m11_instantiate_ebg_with_start.bpmn", "combine startEvent"},
		{"m11_instantiate_ebg_one_catch.bpmn", "at least two"},
	}
	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			xml := readTestdata(t, tc.file)
			eng := processing.NewEngine(eventlog.NewMemory())
			_, err := eng.Deploy(context.Background(), xml)
			if err == nil {
				t.Fatal("expected deploy error")
			}
			if !strings.Contains(err.Error(), "UNSUPPORTED_ELEMENT") {
				t.Fatalf("err=%v want UNSUPPORTED_ELEMENT (hint %q)", err, tc.sub)
			}
		})
	}
}

func TestInstantiateEventBasedGatewayRecoverThenMessageWins(t *testing.T) {
	xml := readTestdata(t, "m11_instantiate_ebg_timer_message.bpmn")
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
	if len(allWaiting(inst)) != 2 {
		t.Fatalf("expected 2 waiting catches after recover, got %v", allWaiting(inst))
	}
	n, err := eng2.PublishMessage(ctx, processing.PublishMessageRequest{Name: "escalate"})
	if err != nil || n != 1 {
		t.Fatalf("PublishMessage n=%d err=%v", n, err)
	}
	inst, _ = eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, _ := eng2.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT, "MessageCatch_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected message catch COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT, "TimerCatch_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected timer catch TERMINATED as sibling")
	}
}
