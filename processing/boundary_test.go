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

func TestNonInterruptingTimerBoundaryFireDueThenComplete(t *testing.T) {
	xml := readTestdata(t, "m2_timer_boundary_non_interrupt.bpmn")
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
		t.Fatalf("status=%s want active after boundary fire", inst.Status)
	}
	if len(inst.Tokens) != 2 {
		t.Fatalf("want 2 tokens after non-interrupting boundary, got %d: %#v", len(inst.Tokens), inst.Tokens)
	}
	elementID, tokenID := waitingAt(inst)
	if elementID != "UserTask_1" {
		t.Fatalf("activity token should still wait at UserTask, tokens=%#v", inst.Tokens)
	}
	tok := inst.Tokens[tokenID]
	if tok.DueUnixMs != 0 || tok.BoundaryID != "" {
		t.Fatalf("boundary should be disarmed on activity token, token=%#v", tok)
	}

	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "UserTask_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("UserTask must not be terminated by non-interrupting boundary")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "TimerBoundary_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected boundary COMPLETED on spawned token")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_timeout", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected timeout end from spawned token")
	}

	if err := eng.Complete(ctx, instanceID, elementID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, err = eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_ok", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected happy-path end after activity completes")
	}
}

func TestNonInterruptingTimerBoundaryCompleteBeforeFireDue(t *testing.T) {
	xml := readTestdata(t, "m2_timer_boundary_non_interrupt.bpmn")
	xmlStr := string(xml)
	xmlStr = strings.Replace(xmlStr, ">PT0S<", ">PT1H<", 1)
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, []byte(xmlStr))
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
	if sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_timeout", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("timeout end must not complete")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "TimerBoundary_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected boundary cancelled")
	}
}

func TestNonInterruptingTimerBoundaryTimeCycleRearms(t *testing.T) {
	xml := readTestdata(t, "m2_timer_boundary_non_interrupt_cycle.bpmn")
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
		t.Fatalf("first FireDue: %v", err)
	}
	inst, _ := eng.GetInstance(instanceID)
	elementID, tokenID := waitingAt(inst)
	if elementID != "UserTask_1" {
		t.Fatalf("activity should still wait after first fire, tokens=%#v", inst.Tokens)
	}
	tok := inst.Tokens[tokenID]
	if tok.BoundaryID != "TimerBoundary_1" || tok.DueUnixMs == 0 || tok.TimerText != "R1/PT0S" {
		t.Fatalf("expected re-armed cycle after first fire, token=%#v", tok)
	}

	if err := eng.FireDue(ctx); err != nil {
		t.Fatalf("second FireDue: %v", err)
	}
	inst, _ = eng.GetInstance(instanceID)
	elementID, tokenID = waitingAt(inst)
	if elementID != "UserTask_1" {
		t.Fatalf("activity should still wait after second fire, tokens=%#v", inst.Tokens)
	}
	tok = inst.Tokens[tokenID]
	if tok.BoundaryID != "" || tok.DueUnixMs != 0 || tok.TimerText != "" {
		t.Fatalf("final cycle fire should disarm boundary, token=%#v", tok)
	}

	if err := eng.FireDue(ctx); err != nil {
		t.Fatalf("third FireDue: %v", err)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if got := countElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "TimerBoundary_1", eventv1.Element_INTENT_COMPLETED); got != 2 {
		t.Fatalf("boundary completions=%d want 2", got)
	}
	if err := eng.Complete(ctx, instanceID, elementID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestNonInterruptingTimerBoundaryTimeCycleRecoverThenRearm(t *testing.T) {
	xml := readTestdata(t, "m2_timer_boundary_non_interrupt_cycle.bpmn")
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
	if err := eng1.FireDue(ctx); err != nil {
		t.Fatalf("first FireDue: %v", err)
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
	elementID, tokenID := waitingAt(inst)
	if elementID != "UserTask_1" {
		t.Fatalf("activity should still wait after recover, tokens=%#v", inst.Tokens)
	}
	tok := inst.Tokens[tokenID]
	if tok.TimerText != "R1/PT0S" || tok.DueUnixMs == 0 {
		t.Fatalf("expected remaining cycle after recover, token=%#v", tok)
	}
	if err := eng2.FireDue(ctx); err != nil {
		t.Fatalf("second FireDue: %v", err)
	}
	events, err := eng2.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if got := countElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "TimerBoundary_1", eventv1.Element_INTENT_COMPLETED); got != 2 {
		t.Fatalf("boundary completions=%d want 2", got)
	}
}

func TestNonInterruptingMessageBoundaryPublishThenComplete(t *testing.T) {
	xml := readTestdata(t, "m2_message_boundary_non_interrupt.bpmn")
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
	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "order.cancelled"})
	if err != nil {
		t.Fatalf("PublishMessage: %v", err)
	}
	if n != 1 {
		t.Fatalf("delivered=%d want 1", n)
	}
	inst, _ := eng.GetInstance(instanceID)
	if inst.Status != projection.StatusActive {
		t.Fatalf("status=%s want active", inst.Status)
	}
	if len(inst.Tokens) != 2 {
		t.Fatalf("want 2 tokens, got %d: %#v", len(inst.Tokens), inst.Tokens)
	}
	elementID, tokenID := waitingAt(inst)
	if elementID != "UserTask_1" {
		t.Fatalf("activity token should still wait, tokens=%#v", inst.Tokens)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "UserTask_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("UserTask must not be terminated")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_msg", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected message end from spawned token")
	}
	if err := eng.Complete(ctx, instanceID, elementID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestDeployAcceptsNonInterruptingTimerBoundary(t *testing.T) {
	xml := readTestdata(t, "m2_timer_boundary_non_interrupt.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	if _, err := eng.Deploy(context.Background(), xml); err != nil {
		t.Fatalf("Deploy: %v", err)
	}
}

func TestDeployAcceptsNonInterruptingMessageBoundary(t *testing.T) {
	xml := readTestdata(t, "m2_message_boundary_non_interrupt.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	if _, err := eng.Deploy(context.Background(), xml); err != nil {
		t.Fatalf("Deploy: %v", err)
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

func countElementIntent(events []*eventv1.Event, typ eventv1.Element_Type, id string, intent eventv1.Element_Intent) int {
	n := 0
	for _, ev := range events {
		if ev.GetRecordType() != eventv1.Event_RECORD_TYPE_EVENT {
			continue
		}
		el := ev.GetElement()
		if el != nil && el.GetType() == typ && el.GetId() == id && el.GetIntent() == intent {
			n++
		}
	}
	return n
}

// --- Dual boundary (timer + message on same activity) ---

func TestDualBoundaryDeploy(t *testing.T) {
	xml := readTestdata(t, "m2_dual_boundary.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	_, err := eng.Deploy(context.Background(), xml)
	if err != nil {
		t.Fatal("should accept dual boundary:", err)
	}
}

func TestDualBoundaryTimerFires(t *testing.T) {
	xml := readTestdata(t, "m2_dual_boundary.bpmn")
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
	if tok.BoundaryID != "TimerBoundary_1" {
		t.Fatalf("expected timer boundary armed, got BoundaryID=%q", tok.BoundaryID)
	}
	if tok.MessageName != "escalate" {
		t.Fatalf("expected message boundary armed, got MessageName=%q", tok.MessageName)
	}
	if tok.MessageBoundaryID != "MessageBoundary_1" {
		t.Fatalf("expected MessageBoundaryID=%q, got %q", "MessageBoundary_1", tok.MessageBoundaryID)
	}

	if err := eng.FireDue(ctx); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "UserTask_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected UserTask TERMINATED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "TimerBoundary_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected timer boundary COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "MessageBoundary_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected message boundary TERMINATED")
	}
}

func TestDualBoundaryMessageFires(t *testing.T) {
	xml := readTestdata(t, "m2_dual_boundary.bpmn")
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
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("delivered=%d want 1", n)
	}
	inst, _ := eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "UserTask_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected UserTask TERMINATED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "MessageBoundary_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected message boundary COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "TimerBoundary_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected timer boundary TERMINATED")
	}
}

func TestDualBoundaryActivityCompletes(t *testing.T) {
	xml := readTestdata(t, "m2_dual_boundary.bpmn")
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
	elemID, tokenID := waitingAt(inst)
	if err := eng.Complete(ctx, instanceID, elemID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, _ := eng.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "TimerBoundary_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected timer boundary TERMINATED on complete")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "MessageBoundary_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected message boundary TERMINATED on complete")
	}
}

func TestDualBoundaryRecoverThenTimerFires(t *testing.T) {
	xml := readTestdata(t, "m2_dual_boundary.bpmn")
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
	if tok.BoundaryID != "TimerBoundary_1" {
		t.Fatalf("timer boundary not restored: BoundaryID=%q", tok.BoundaryID)
	}
	if tok.MessageBoundaryID != "MessageBoundary_1" {
		t.Fatalf("message boundary not restored: MessageBoundaryID=%q", tok.MessageBoundaryID)
	}
	if tok.MessageName != "escalate" {
		t.Fatalf("message name not restored: MessageName=%q", tok.MessageName)
	}
	if tok.DueUnixMs == 0 {
		t.Fatal("DueUnixMs not restored")
	}

	if err := eng2.FireDue(ctx); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestDualBoundaryRecoverThenMessageFires(t *testing.T) {
	xml := readTestdata(t, "m2_dual_boundary.bpmn")
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
	_ = instanceID
	if err := eng1.Close(); err != nil {
		t.Fatal(err)
	}

	eng2, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	defer eng2.Close()

	n, err := eng2.PublishMessage(ctx, processing.PublishMessageRequest{Name: "escalate"})
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("delivered=%d want 1", n)
	}
	inst, _ := eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	events, _ := eng2.ListEvents(ctx, instanceID)
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "MessageBoundary_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected message boundary COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "TimerBoundary_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected timer boundary TERMINATED")
	}
}
