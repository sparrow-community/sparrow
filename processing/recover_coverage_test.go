package processing_test

import (
	"context"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// reopenEngine closes eng and reopens the same data dir, standing in for a
// process restart.
func reopenEngine(t *testing.T, eng *processing.Engine, dir string) *processing.Engine {
	t.Helper()
	if err := eng.Close(); err != nil {
		t.Fatal(err)
	}
	next, err := processing.Open(context.Background(), dir)
	if err != nil {
		t.Fatal(err)
	}
	return next
}

func TestInclusiveGatewayJoinRecover(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	eng, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	dep, err := eng.Deploy(ctx, readTestdata(t, "m2_inclusive_gateway.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, map[string]any{"path_a": "true", "path_c": "true"})
	if err != nil {
		t.Fatal(err)
	}
	inst := mustInstance(t, eng, instanceID)
	tokenA := tokenAtElement(inst, "Task_A")
	if tokenA == "" {
		t.Fatalf("expected Task_A waiting, got %v", allWaiting(inst))
	}
	if err := eng.Complete(ctx, instanceID, "Task_A", tokenA, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, instanceID)
	if tokenAtElement(inst, "Join_1") == "" {
		t.Fatalf("expected a token waiting at Join_1, got %v", allWaiting(inst))
	}

	eng = reopenEngine(t, eng, dir)
	defer eng.Close()

	inst = mustInstance(t, eng, instanceID)
	if tokenAtElement(inst, "Join_1") == "" {
		t.Fatalf("Join_1 wait not rebuilt after recover, waiting=%v", allWaiting(inst))
	}
	tokenC := tokenAtElement(inst, "Task_C")
	if tokenC == "" {
		t.Fatalf("expected Task_C waiting after recover, got %v", allWaiting(inst))
	}
	if err := eng.Complete(ctx, instanceID, "Task_C", tokenC, nil); err != nil {
		t.Fatal(err)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_INCLUSIVE_GATEWAY, "Join_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected Join_1 COMPLETED after recover + second arrival")
	}
	if inst := mustInstance(t, eng, instanceID); inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s waiting=%v", inst.Status, allWaiting(inst))
	}
}

func TestMultiInstanceBehaviorEventRecover(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	eng, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	catchDep, err := eng.Deploy(ctx, readTestdata(t, "m25_catch_mi_each.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	miDep, err := eng.Deploy(ctx, readTestdata(t, "m25_mi_none_behavior_event.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	catchFirst, err := eng.CreateInstance(ctx, catchDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	miID, err := eng.CreateInstance(ctx, miDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inners := waitingTokensAt(mustInstance(t, eng, miID), "UserTask_mi")
	if len(inners) != 2 {
		t.Fatalf("inners=%d want 2", len(inners))
	}
	if err := eng.Complete(ctx, miID, "UserTask_mi", inners[0], nil); err != nil {
		t.Fatal(err)
	}
	if mustInstance(t, eng, catchFirst).Status != projection.StatusCompleted {
		t.Fatal("first catcher must complete on the first inner completion")
	}

	eng = reopenEngine(t, eng, dir)
	defer eng.Close()

	remaining := waitingTokensAt(mustInstance(t, eng, miID), "UserTask_mi")
	if len(remaining) != 1 {
		t.Fatalf("remaining inners after recover=%d want 1", len(remaining))
	}
	catchSecond, err := eng.CreateInstance(ctx, catchDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := eng.Complete(ctx, miID, "UserTask_mi", remaining[0], nil); err != nil {
		t.Fatal(err)
	}
	if inst := mustInstance(t, eng, miID); inst.Status != projection.StatusCompleted {
		t.Fatalf("multi-instance status=%s after recover", inst.Status)
	}
	if inst := mustInstance(t, eng, catchSecond); inst.Status != projection.StatusCompleted {
		t.Fatalf("behavior event did not publish after recover, catcher status=%s", inst.Status)
	}
}

func TestInstantiateReceiveTaskEntryRecover(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	eng, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	dep, err := eng.Deploy(ctx, readTestdata(t, "m26_instantiate_receive.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	if el := waitingElement(mustInstance(t, eng, instanceID)); el != "ReceiveTask_1" {
		t.Fatalf("waiting=%s want ReceiveTask_1", el)
	}

	eng = reopenEngine(t, eng, dir)
	defer eng.Close()

	if el := waitingElement(mustInstance(t, eng, instanceID)); el != "ReceiveTask_1" {
		t.Fatalf("waiting=%s after recover want ReceiveTask_1", el)
	}
	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{
		Name:              "order.confirmed",
		ProcessInstanceID: instanceID,
	})
	if err != nil || n != 1 {
		t.Fatalf("publish n=%d err=%v", n, err)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_RECEIVE_TASK, "ReceiveTask_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected ReceiveTask_1 COMPLETED after recover + publish")
	}
	if inst := mustInstance(t, eng, instanceID); inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestProcessVersionCoexistenceRecover(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	eng, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	dep1, err := eng.Deploy(ctx, readTestdata(t, "m5_version_v1.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	instanceV1, err := eng.CreateInstance(ctx, dep1, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := eng.Deploy(ctx, readTestdata(t, "m5_version_v2.bpmn")); err != nil {
		t.Fatal(err)
	}
	instanceV2, err := eng.CreateInstanceRequest(ctx, processing.CreateInstanceRequest{ProcessID: "Process_version"})
	if err != nil {
		t.Fatal(err)
	}

	eng = reopenEngine(t, eng, dir)
	defer eng.Close()

	instV1 := mustInstance(t, eng, instanceV1)
	if instV1.Version != 1 || instV1.DeploymentID != dep1 || waitingElement(instV1) != "TaskA" {
		t.Fatalf("v1 instance after recover: version=%d dep=%s waiting=%s", instV1.Version, instV1.DeploymentID, waitingElement(instV1))
	}
	instV2 := mustInstance(t, eng, instanceV2)
	if instV2.Version != 2 || waitingElement(instV2) != "TaskB" {
		t.Fatalf("v2 instance after recover: version=%d waiting=%s", instV2.Version, waitingElement(instV2))
	}

	// Latest still resolves to revision 2, explicit version still resolves to 1.
	latest, err := eng.CreateInstanceRequest(ctx, processing.CreateInstanceRequest{ProcessID: "Process_version"})
	if err != nil {
		t.Fatal(err)
	}
	if inst := mustInstance(t, eng, latest); inst.Version != 2 || waitingElement(inst) != "TaskB" {
		t.Fatalf("latest after recover: version=%d waiting=%s", inst.Version, waitingElement(inst))
	}
	pinned, err := eng.CreateInstanceRequest(ctx, processing.CreateInstanceRequest{
		ProcessID:      "Process_version",
		ProcessVersion: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if inst := mustInstance(t, eng, pinned); inst.Version != 1 || waitingElement(inst) != "TaskA" {
		t.Fatalf("pinned v1 after recover: version=%d waiting=%s", inst.Version, waitingElement(inst))
	}

	// Both in-flight instances still complete on their own revision path.
	for _, tc := range []struct {
		instanceID string
		elementID  string
	}{
		{instanceV1, "TaskA"},
		{instanceV2, "TaskB"},
	} {
		inst := mustInstance(t, eng, tc.instanceID)
		tokenID := tokenAtElement(inst, tc.elementID)
		if tokenID == "" {
			t.Fatalf("no token at %s", tc.elementID)
		}
		if err := eng.Complete(ctx, tc.instanceID, tc.elementID, tokenID, nil); err != nil {
			t.Fatalf("complete %s: %v", tc.elementID, err)
		}
		if inst := mustInstance(t, eng, tc.instanceID); inst.Status != projection.StatusCompleted {
			t.Fatalf("%s status=%s want completed", tc.elementID, inst.Status)
		}
	}
}

func TestErrorBoundaryRecoverThenThrowError(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	eng, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	dep, err := eng.Deploy(ctx, readTestdata(t, "m4_error_boundary.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	if el := waitingElement(mustInstance(t, eng, instanceID)); el != "UserTask_1" {
		t.Fatalf("waiting=%s want UserTask_1", el)
	}

	eng = reopenEngine(t, eng, dir)
	defer eng.Close()

	inst := mustInstance(t, eng, instanceID)
	tokenID := tokenAtElement(inst, "UserTask_1")
	if tokenID == "" {
		t.Fatalf("UserTask_1 not waiting after recover, waiting=%v", allWaiting(inst))
	}
	if err := eng.ThrowError(ctx, instanceID, "UserTask_1", tokenID, "BUSINESS_ERROR"); err != nil {
		t.Fatal(err)
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "ErrorBoundary_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected ErrorBoundary_1 COMPLETED after recover + fresh ThrowError")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "EndEvent_err", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected the error path end event after recover")
	}
	if inst := mustInstance(t, eng, instanceID); inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}
