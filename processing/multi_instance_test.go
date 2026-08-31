package processing_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	"github.com/sparrow-community/sparrow/processing/deploy"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
)

func TestMultiInstanceParallelUserTask(t *testing.T) {
	eng := newMIEngine(t, "m6_mi_parallel_cardinality.bpmn")
	ctx := context.Background()
	dep := miDeploy(t, eng, "m6_mi_parallel_cardinality.bpmn")
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, ok := eng.GetInstance(instanceID)
	if !ok {
		t.Fatal("missing instance")
	}
	waiting := waitingMItokens(inst, "UserTask_mi")
	if len(waiting) != 3 {
		t.Fatalf("waiting=%d want 3 tokens=%#v", len(waiting), inst.Tokens)
	}
	for _, tok := range waiting {
		if tok.LoopInstanceIndex < 0 {
			t.Fatalf("inner token missing loop index: %#v", tok)
		}
		if err := eng.Complete(ctx, instanceID, "UserTask_mi", tok.ID, nil); err != nil {
			t.Fatal(err)
		}
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestMultiInstanceParallelServiceTaskJobs(t *testing.T) {
	eng := newMIEngine(t, "m6_mi_parallel_service.bpmn")
	ctx := context.Background()
	dep := miDeploy(t, eng, "m6_mi_parallel_service.bpmn")
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	jobs, err := eng.Activate(ctx, processing.ActivateRequest{JobType: "work.v1", WorkerID: "w1", MaxJobs: 3})
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs) != 3 {
		t.Fatalf("jobs=%d want 3", len(jobs))
	}
	for _, job := range jobs {
		if err := eng.Complete(ctx, instanceID, job.ElementID, job.TokenID, nil); err != nil {
			t.Fatal(err)
		}
	}
	inst, _ := eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestMultiInstanceSequentialUserTask(t *testing.T) {
	eng := newMIEngine(t, "m6_mi_sequential_cardinality.bpmn")
	ctx := context.Background()
	dep := miDeploy(t, eng, "m6_mi_sequential_cardinality.bpmn")
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		inst, _ := eng.GetInstance(instanceID)
		waiting := waitingMItokens(inst, "UserTask_mi")
		if len(waiting) != 1 {
			t.Fatalf("step %d waiting=%d want 1", i, len(waiting))
		}
		if err := eng.Complete(ctx, instanceID, "UserTask_mi", waiting[0].ID, nil); err != nil {
			t.Fatal(err)
		}
	}
	inst, _ := eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestMultiInstanceRecoverMidLoop(t *testing.T) {
	xml, err := os.ReadFile(filepath.Join("testdata", "m6_mi_parallel_cardinality.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	log := eventlog.NewMemory()
	store := deploy.NewMemoryStore()
	eng1, err := processing.Recover(ctx, log, store, nil)
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
	waiting := waitingMItokens(inst, "UserTask_mi")
	if err := eng1.Complete(ctx, instanceID, "UserTask_mi", waiting[0].ID, nil); err != nil {
		t.Fatal(err)
	}
	eng2, err := processing.Recover(ctx, log, store, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ = eng2.GetInstance(instanceID)
	waiting = waitingMItokens(inst, "UserTask_mi")
	if len(waiting) != 2 {
		t.Fatalf("after recover waiting=%d want 2", len(waiting))
	}
	for _, tok := range waiting {
		if err := eng2.Complete(ctx, instanceID, "UserTask_mi", tok.ID, nil); err != nil {
			t.Fatal(err)
		}
	}
	inst, _ = eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func waitingMItokens(inst *projection.Instance, elementID string) []*projection.Token {
	var out []*projection.Token
	for _, tok := range inst.Tokens {
		if tok == nil || tok.ElementID != elementID || tok.Status != projection.TokenWaiting {
			continue
		}
		if tok.LoopInstanceIndex < 0 || tok.MultiInstanceHost {
			continue
		}
		out = append(out, tok)
	}
	return out
}

func newMIEngine(t *testing.T, name string) *processing.Engine {
	t.Helper()
	xml, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	eng := processing.NewEngine(eventlog.NewMemory())
	if _, err := eng.Deploy(context.Background(), xml); err != nil {
		t.Fatal(err)
	}
	return eng
}

func miDeploy(t *testing.T, eng *processing.Engine, name string) string {
	t.Helper()
	xml, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	dep, err := eng.Deploy(context.Background(), xml)
	if err != nil {
		t.Fatal(err)
	}
	return dep
}

func TestMultiInstanceCollectionInput(t *testing.T) {
	eng := newMIEngine(t, "m6_mi_collection.bpmn")
	ctx := context.Background()
	dep := miDeploy(t, eng, "m6_mi_collection.bpmn")
	instanceID, err := eng.CreateInstance(ctx, dep, map[string]any{"items": []string{"a", "b"}})
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	if len(waitingMItokens(inst, "UserTask_mi")) != 2 {
		t.Fatalf("tokens=%#v", inst.Tokens)
	}
	for _, tok := range waitingMItokens(inst, "UserTask_mi") {
		if err := eng.Complete(ctx, instanceID, "UserTask_mi", tok.ID, map[string]any{"result": tok.LoopInstanceIndex}); err != nil {
			t.Fatal(err)
		}
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Variables["results"] == "" {
		t.Fatalf("missing results: %#v", inst.Variables)
	}
}

func TestMultiInstanceEmptyCollection(t *testing.T) {
	eng := newMIEngine(t, "m6_mi_collection_empty.bpmn")
	ctx := context.Background()
	dep := miDeploy(t, eng, "m6_mi_collection_empty.bpmn")
	instanceID, err := eng.CreateInstance(ctx, dep, map[string]any{"items": []string{}})
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	if len(waitingMItokens(inst, "UserTask_mi")) != 0 {
		t.Fatalf("expected no waiting work, tokens=%#v", inst.Tokens)
	}
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestMultiInstanceEarlyCompletion(t *testing.T) {
	eng := newMIEngine(t, "m6_mi_completion_early.bpmn")
	ctx := context.Background()
	dep := miDeploy(t, eng, "m6_mi_completion_early.bpmn")
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	waiting := waitingMItokens(inst, "UserTask_mi")
	if len(waiting) != 5 {
		t.Fatalf("waiting=%d want 5", len(waiting))
	}
	for i := 0; i < 2; i++ {
		if err := eng.Complete(ctx, instanceID, "UserTask_mi", waiting[i].ID, nil); err != nil {
			t.Fatal(err)
		}
	}
	inst, _ = eng.GetInstance(instanceID)
	if len(waitingMItokens(inst, "UserTask_mi")) != 0 {
		t.Fatalf("stragglers remain: %#v", inst.Tokens)
	}
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestMultiInstanceSubProcess(t *testing.T) {
	eng := newMIEngine(t, "m6_mi_subprocess.bpmn")
	ctx := context.Background()
	dep := miDeploy(t, eng, "m6_mi_subprocess.bpmn")
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		inst, _ := eng.GetInstance(instanceID)
		var subTask string
		for _, tok := range inst.Tokens {
			if tok != nil && tok.ElementID == "SubTask_1" && tok.Status == projection.TokenWaiting {
				subTask = tok.ID
				break
			}
		}
		if subTask == "" {
			t.Fatalf("step %d: no waiting sub task, tokens=%#v", i, inst.Tokens)
		}
		if err := eng.Complete(ctx, instanceID, "SubTask_1", subTask, nil); err != nil {
			t.Fatal(err)
		}
	}
	inst, _ := eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestMultiInstanceTimerBoundaryCancel(t *testing.T) {
	eng := newMIEngine(t, "m6_mi_timer_boundary.bpmn")
	ctx := context.Background()
	dep := miDeploy(t, eng, "m6_mi_timer_boundary.bpmn")
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	if len(waitingMItokens(inst, "UserTask_mi")) != 3 {
		t.Fatalf("waiting inner tokens=%#v", inst.Tokens)
	}
	if err := eng.FireDue(ctx); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	for _, tok := range inst.Tokens {
		if tok != nil && tok.ElementID == "UserTask_mi" && tok.Status == projection.TokenWaiting {
			t.Fatalf("stray MI token=%#v", tok)
		}
	}
}
