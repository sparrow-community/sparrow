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

func TestMessageEndCompletesAlone(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m20_message_end.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	id, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, ok := eng.GetInstance(id)
	if !ok || inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%v ok=%v", inst, ok)
	}
	events, err := eng.ListEvents(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "MessageEnd_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected MessageEnd_1 COMPLETED")
	}
}

func TestMessageEndDeliversToCatch(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	catchDep, err := eng.Deploy(ctx, readTestdata(t, "m2_message_catch.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	endDep, err := eng.Deploy(ctx, readTestdata(t, "m20_message_end.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	catchID, err := eng.CreateInstance(ctx, catchDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	if waitingElement(mustInstance(t, eng, catchID)) != "MessageCatch_1" {
		t.Fatal("expected wait at MessageCatch_1")
	}
	endID, err := eng.CreateInstance(ctx, endDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	endInst, _ := eng.GetInstance(endID)
	if endInst.Status != projection.StatusCompleted {
		t.Fatalf("ender status=%v", endInst.Status)
	}
	catchInst, _ := eng.GetInstance(catchID)
	if catchInst.Status != projection.StatusCompleted {
		t.Fatalf("catcher status=%v", catchInst.Status)
	}
}

func TestSignalEndDeliversToCatch(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	catchDep, err := eng.Deploy(ctx, readTestdata(t, "m3_signal_catch.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	endDep, err := eng.Deploy(ctx, readTestdata(t, "m20_signal_end.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	catchID, err := eng.CreateInstance(ctx, catchDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	endID, err := eng.CreateInstance(ctx, endDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	endInst, _ := eng.GetInstance(endID)
	if endInst.Status != projection.StatusCompleted {
		t.Fatalf("ender status=%v", endInst.Status)
	}
	catchInst, _ := eng.GetInstance(catchID)
	if catchInst.Status != projection.StatusCompleted {
		t.Fatalf("catcher status=%v", catchInst.Status)
	}
}

func TestSubProcessMessageEnd(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m20_subprocess_message_end.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	id, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(id)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%v", inst.Status)
	}
}

func TestSubProcessSignalEnd(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m20_subprocess_signal_end.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	id, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(id)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%v", inst.Status)
	}
}

func TestMixedMessageTerminateEndRejected(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	_, err := eng.Deploy(context.Background(), readTestdata(t, "m20_mixed_message_terminate_end.bpmn"))
	if err == nil || !strings.Contains(err.Error(), "UNSUPPORTED_ELEMENT") {
		t.Fatalf("Deploy err=%v want UNSUPPORTED_ELEMENT", err)
	}
}

func TestMessageEndRecoverAfterDelivery(t *testing.T) {
	dir := t.TempDir()
	ctx := context.Background()
	eng1, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	catchDep, err := eng1.Deploy(ctx, readTestdata(t, "m2_message_catch.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	endDep, err := eng1.Deploy(ctx, readTestdata(t, "m20_message_end.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	catchID, err := eng1.CreateInstance(ctx, catchDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	endID, err := eng1.CreateInstance(ctx, endDep, nil)
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
	catchInst, ok := eng2.GetInstance(catchID)
	if !ok || catchInst.Status != projection.StatusCompleted {
		t.Fatalf("catch after recover status=%v ok=%v", catchInst, ok)
	}
	endInst, ok := eng2.GetInstance(endID)
	if !ok || endInst.Status != projection.StatusCompleted {
		t.Fatalf("end after recover status=%v ok=%v", endInst, ok)
	}
}

func TestSignalEndRecoverAfterDelivery(t *testing.T) {
	dir := t.TempDir()
	ctx := context.Background()
	eng1, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	catchDep, err := eng1.Deploy(ctx, readTestdata(t, "m3_signal_catch.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	endDep, err := eng1.Deploy(ctx, readTestdata(t, "m20_signal_end.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	catchID, err := eng1.CreateInstance(ctx, catchDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	endID, err := eng1.CreateInstance(ctx, endDep, nil)
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
	catchInst, ok := eng2.GetInstance(catchID)
	if !ok || catchInst.Status != projection.StatusCompleted {
		t.Fatalf("catch after recover status=%v ok=%v", catchInst, ok)
	}
	endInst, ok := eng2.GetInstance(endID)
	if !ok || endInst.Status != projection.StatusCompleted {
		t.Fatalf("end after recover status=%v ok=%v", endInst, ok)
	}
}
