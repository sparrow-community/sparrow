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

func TestOpenReplayResumeUserTask(t *testing.T) {
	xml, err := os.ReadFile(filepath.Join("testdata", "m1_simple.bpmn"))
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
	instanceID, err := eng1.CreateInstance(ctx, dep, map[string]any{"initiator": "alice"})
	if err != nil {
		t.Fatal(err)
	}
	inst, ok := eng1.GetInstance(instanceID)
	if !ok || inst.Status != projection.StatusActive {
		t.Fatalf("inst=%v ok=%v", inst, ok)
	}
	var waitingToken, waitingElement string
	for _, tok := range inst.Tokens {
		if tok.Status == projection.TokenWaiting {
			waitingToken, waitingElement = tok.ID, tok.ElementID
			break
		}
	}
	if waitingElement != "UserTask_1" {
		t.Fatalf("tokens=%#v", inst.Tokens)
	}
	if err := eng1.Close(); err != nil {
		t.Fatal(err)
	}

	eng2, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	defer eng2.Close()

	inst2, ok := eng2.GetInstance(instanceID)
	if !ok {
		t.Fatal("instance missing after replay")
	}
	if inst2.Status != projection.StatusActive {
		t.Fatalf("status=%s", inst2.Status)
	}
	if inst2.Variables["initiator"] == "" {
		t.Fatalf("vars=%#v", inst2.Variables)
	}
	var tok2 *projection.Token
	for _, tok := range inst2.Tokens {
		if tok.ID == waitingToken {
			tok2 = tok
			break
		}
	}
	if tok2 == nil || tok2.ElementID != waitingElement || tok2.Status != projection.TokenWaiting {
		t.Fatalf("token after replay=%#v tokens=%#v", tok2, inst2.Tokens)
	}

	if err := eng2.Complete(ctx, instanceID, waitingElement, waitingToken, map[string]any{"approved": true}); err != nil {
		t.Fatalf("Complete after replay: %v", err)
	}
	inst2, _ = eng2.GetInstance(instanceID)
	if inst2.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst2.Status)
	}
}

func TestRecoverReplayWithMemoryBackends(t *testing.T) {
	xml, err := os.ReadFile(filepath.Join("testdata", "m1_simple.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	memLog := eventlog.NewMemory()
	memStore := deploy.NewMemoryStore()

	eng1, err := processing.Recover(ctx, memLog, memStore, nil)
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

	eng2, err := processing.Recover(ctx, memLog, memStore, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, ok := eng2.GetInstance(instanceID)
	if !ok || inst.Status != projection.StatusActive {
		t.Fatalf("after recover inst=%v ok=%v", inst, ok)
	}
	var waitingToken, waitingElement string
	for _, tok := range inst.Tokens {
		if tok.Status == projection.TokenWaiting {
			waitingToken, waitingElement = tok.ID, tok.ElementID
			break
		}
	}
	if waitingElement != "UserTask_1" {
		t.Fatalf("tokens=%#v", inst.Tokens)
	}
	if err := eng2.Complete(ctx, instanceID, waitingElement, waitingToken, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s", inst.Status)
	}
}
