package processing_test

import (
	"context"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
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

func containsStr(ss []string, want string) bool {
	for _, s := range ss {
		if s == want {
			return true
		}
	}
	return false
}
