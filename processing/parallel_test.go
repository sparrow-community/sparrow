package processing_test

import (
	"context"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func TestParallelForkJoinTwoUserTasks(t *testing.T) {
	xml := readTestdata(t, "m3_parallel_fork_join.bpmn")
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
	if len(inst.Tokens) != 2 {
		t.Fatalf("want 2 tokens after fork, got %d: %#v", len(inst.Tokens), inst.Tokens)
	}
	var tasks []string
	for _, tok := range inst.Tokens {
		if tok.Status != projection.TokenWaiting {
			t.Fatalf("token not waiting: %#v", tok)
		}
		tasks = append(tasks, tok.ElementID)
	}
	if !containsStr(tasks, "UserTask_a") || !containsStr(tasks, "UserTask_b") {
		t.Fatalf("tasks=%v", tasks)
	}

	for _, tok := range inst.Tokens {
		if err := eng.Complete(ctx, instanceID, tok.ElementID, tok.ID, nil); err != nil {
			t.Fatalf("Complete %s: %v", tok.ElementID, err)
		}
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestParallelTimerBoundaryFireDueTargetsOnlyAttachedBranch(t *testing.T) {
	xml := readTestdata(t, "m3_parallel_timer_boundary.bpmn")
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
	if countWaitingAt(inst, "UserTask_a") != 1 || countWaitingAt(inst, "UserTask_b") != 1 {
		t.Fatalf("expected waits on both branches, tokens=%#v", inst.Tokens)
	}

	if err := eng.FireDue(ctx); err != nil {
		t.Fatalf("FireDue: %v", err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusActive {
		t.Fatalf("status=%s want active", inst.Status)
	}
	if countWaitingAt(inst, "UserTask_a") != 0 || countWaitingAt(inst, "UserTask_b") != 1 {
		t.Fatalf("timer boundary should only clear branch A wait, tokens=%#v", inst.Tokens)
	}

	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "TimerBoundary_a", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected TimerBoundary_a COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "UserTask_a", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected UserTask_a TERMINATED")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "UserTask_b", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("UserTask_b must stay waiting")
	}

	for _, tok := range inst.Tokens {
		if tok.Status == projection.TokenWaiting && tok.ElementID == "UserTask_b" {
			if err := eng.Complete(ctx, instanceID, tok.ElementID, tok.ID, nil); err != nil {
				t.Fatalf("Complete UserTask_b: %v", err)
			}
		}
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestParallelMessageBoundaryPublishTargetsOnlyAttachedBranch(t *testing.T) {
	xml := readTestdata(t, "m3_parallel_message_boundary.bpmn")
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
	if countWaitingAt(inst, "UserTask_a") != 1 || countWaitingAt(inst, "UserTask_b") != 1 {
		t.Fatalf("expected waits on both branches, tokens=%#v", inst.Tokens)
	}

	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "order.cancelled"})
	if err != nil {
		t.Fatalf("PublishMessage: %v", err)
	}
	if n != 1 {
		t.Fatalf("delivered=%d want 1", n)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusActive {
		t.Fatalf("status=%s want active", inst.Status)
	}
	if countWaitingAt(inst, "UserTask_a") != 0 || countWaitingAt(inst, "UserTask_b") != 1 {
		t.Fatalf("message boundary should only clear branch A wait, tokens=%#v", inst.Tokens)
	}

	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "MessageBoundary_a", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected MessageBoundary_a COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "UserTask_a", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected UserTask_a TERMINATED")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "UserTask_b", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("UserTask_b must stay waiting")
	}

	for _, tok := range inst.Tokens {
		if tok.Status == projection.TokenWaiting && tok.ElementID == "UserTask_b" {
			if err := eng.Complete(ctx, instanceID, tok.ElementID, tok.ID, nil); err != nil {
				t.Fatalf("Complete UserTask_b: %v", err)
			}
		}
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestParallelTimerBoundaryRecoverThenFireDue(t *testing.T) {
	xml := readTestdata(t, "m3_parallel_timer_boundary.bpmn")
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
	if countWaitingAt(inst, "UserTask_a") != 1 || countWaitingAt(inst, "UserTask_b") != 1 {
		t.Fatalf("expected waits on both branches after recover, tokens=%#v", inst.Tokens)
	}
	if err := eng2.FireDue(ctx); err != nil {
		t.Fatalf("FireDue: %v", err)
	}
	inst, _ = eng2.GetInstance(instanceID)
	if countWaitingAt(inst, "UserTask_a") != 0 || countWaitingAt(inst, "UserTask_b") != 1 {
		t.Fatalf("timer boundary should only clear branch A wait after recover, tokens=%#v", inst.Tokens)
	}
	for _, tok := range inst.Tokens {
		if tok.Status == projection.TokenWaiting && tok.ElementID == "UserTask_b" {
			if err := eng2.Complete(ctx, instanceID, tok.ElementID, tok.ID, nil); err != nil {
				t.Fatalf("Complete UserTask_b: %v", err)
			}
		}
	}
	inst, _ = eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestParallelMessageBoundaryRecoverThenPublish(t *testing.T) {
	xml := readTestdata(t, "m3_parallel_message_boundary.bpmn")
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
	if countWaitingAt(inst, "UserTask_a") != 1 || countWaitingAt(inst, "UserTask_b") != 1 {
		t.Fatalf("expected waits on both branches after recover, tokens=%#v", inst.Tokens)
	}
	n, err := eng2.PublishMessage(ctx, processing.PublishMessageRequest{Name: "order.cancelled"})
	if err != nil {
		t.Fatalf("PublishMessage: %v", err)
	}
	if n != 1 {
		t.Fatalf("delivered=%d want 1", n)
	}
	inst, _ = eng2.GetInstance(instanceID)
	if countWaitingAt(inst, "UserTask_a") != 0 || countWaitingAt(inst, "UserTask_b") != 1 {
		t.Fatalf("message boundary should only clear branch A wait after recover, tokens=%#v", inst.Tokens)
	}
	for _, tok := range inst.Tokens {
		if tok.Status == projection.TokenWaiting && tok.ElementID == "UserTask_b" {
			if err := eng2.Complete(ctx, instanceID, tok.ElementID, tok.ID, nil); err != nil {
				t.Fatalf("Complete UserTask_b: %v", err)
			}
		}
	}
	inst, _ = eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func countWaitingAt(inst *projection.Instance, elementID string) int {
	n := 0
	for _, tok := range inst.Tokens {
		if tok.Status == projection.TokenWaiting && tok.ElementID == elementID {
			n++
		}
	}
	return n
}

func containsStr(ss []string, want string) bool {
	for _, s := range ss {
		if s == want {
			return true
		}
	}
	return false
}
