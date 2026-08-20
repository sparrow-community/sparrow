package processing_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sparrow-community/sparrow/processing"
	"github.com/sparrow-community/sparrow/processing/projection"
)

func TestJobLeaseSurvivesOpenRecover(t *testing.T) {
	xml, err := os.ReadFile(filepath.Join("testdata", "m2_service_task.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
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
	first, err := eng1.Activate(ctx, processing.ActivateRequest{
		JobType:      "work.v1",
		WorkerID:     "worker-a",
		LockDuration: time.Minute,
	})
	if err != nil || len(first) != 1 {
		t.Fatalf("first=%#v err=%v", first, err)
	}
	if err := eng1.Close(); err != nil {
		t.Fatal(err)
	}

	eng2, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	defer eng2.Close()

	again, err := eng2.Activate(ctx, processing.ActivateRequest{
		JobType:  "work.v1",
		WorkerID: "worker-b",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(again) != 0 {
		t.Fatalf("lease should survive recover: %#v", again)
	}

	if err := eng2.Complete(ctx, instanceID, first[0].ElementID, first[0].TokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ := eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestMessageBufferDroppedWhenInstanceCompleted(t *testing.T) {
	xml := readTestdataRuntime(t, "m2_message_catch.bpmn")
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
	elementID, tokenID := waitingAt(inst)
	if err := eng1.Complete(ctx, instanceID, elementID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng1.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}

	if _, err := eng1.PublishMessage(ctx, processing.PublishMessageRequest{
		Name:              "order.confirmed",
		ProcessInstanceID: instanceID,
	}); err != nil {
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

	n, err := eng2.PublishMessage(ctx, processing.PublishMessageRequest{Name: "order.confirmed"})
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("completed instance buffer should be dropped on reload; delivered=%d", n)
	}
}

func readTestdataRuntime(t *testing.T, name string) []byte {
	t.Helper()
	xml, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return xml
}
