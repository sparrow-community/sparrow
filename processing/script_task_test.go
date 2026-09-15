package processing_test

import (
	"context"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func TestScriptTaskActivateComplete(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, readTestdata(t, "m22_script_task.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	id, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	jobs, err := eng.Activate(ctx, processing.ActivateRequest{JobType: "javascript", WorkerID: "w1"})
	if err != nil || len(jobs) != 1 {
		t.Fatalf("Activate jobs=%d err=%v", len(jobs), err)
	}
	if jobs[0].ScriptFormat != "javascript" || jobs[0].Script != "return true;" {
		t.Fatalf("job script=%q format=%q", jobs[0].Script, jobs[0].ScriptFormat)
	}
	if err := eng.Complete(ctx, id, jobs[0].ElementID, jobs[0].TokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(id)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%v", inst.Status)
	}
	events, _ := eng.ListEvents(ctx, id)
	if !sawElementIntent(events, eventv1.Element_TYPE_SCRIPT_TASK, "ScriptTask_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected SCRIPT_TASK COMPLETED")
	}
}

func TestScriptTaskRecoverThenActivateComplete(t *testing.T) {
	dir := t.TempDir()
	ctx := context.Background()
	eng1, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	dep, err := eng1.Deploy(ctx, readTestdata(t, "m22_script_task.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	id, err := eng1.CreateInstance(ctx, dep, nil)
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
	jobs, err := eng2.Activate(ctx, processing.ActivateRequest{JobType: "javascript", WorkerID: "w2"})
	if err != nil || len(jobs) != 1 {
		t.Fatalf("Activate jobs=%d err=%v", len(jobs), err)
	}
	if err := eng2.Complete(ctx, id, jobs[0].ElementID, jobs[0].TokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ := eng2.GetInstance(id)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%v", inst.Status)
	}
}
