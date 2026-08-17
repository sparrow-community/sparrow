package processing_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func TestM1DeployCreateComplete(t *testing.T) {
	xml, err := os.ReadFile(filepath.Join("testdata", "m1_simple.bpmn"))
	if err != nil {
		t.Fatal(err)
	}

	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()

	deploymentID, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatalf("Deploy: %v", err)
	}
	if deploymentID == "" {
		t.Fatal("empty deployment id")
	}

	instanceID, err := eng.CreateInstance(ctx, deploymentID, map[string]any{"initiator": "alice"})
	if err != nil {
		t.Fatalf("CreateInstance: %v", err)
	}

	inst, ok := eng.GetInstance(instanceID)
	if !ok {
		t.Fatal("instance missing")
	}
	if inst.Status != projection.StatusActive {
		t.Fatalf("status=%s want active", inst.Status)
	}
	if inst.Variables["initiator"] == "" {
		t.Fatalf("variables not applied: %#v", inst.Variables)
	}

	var waitingToken string
	var waitingElement string
	for _, tok := range inst.Tokens {
		if tok.Status == projection.TokenWaiting {
			waitingToken = tok.ID
			waitingElement = tok.ElementID
			break
		}
	}
	if waitingElement != "UserTask_1" || waitingToken == "" {
		t.Fatalf("expected wait at UserTask_1, tokens=%#v", inst.Tokens)
	}

	if err := eng.CompleteUserTask(ctx, instanceID, waitingElement, waitingToken, map[string]any{"approved": true}); err != nil {
		t.Fatalf("CompleteUserTask: %v", err)
	}

	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	if inst.Variables["approved"] == "" {
		t.Fatalf("complete vars missing: %#v", inst.Variables)
	}

	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) < 5 {
		t.Fatalf("too few events: %d", len(events))
	}

	var sawFlow, sawUserActivated, sawProcessCompleted bool
	for _, ev := range events {
		if ev.GetRecordType() != eventv1.Event_RECORD_TYPE_EVENT {
			continue
		}
		el := ev.GetElement()
		if el == nil {
			continue
		}
		if el.GetType() == eventv1.Element_TYPE_USER_TASK && el.GetIntent() == eventv1.Element_INTENT_ACTIVATED {
			sawUserActivated = true
		}
		if el.GetType() == eventv1.Element_TYPE_SEQUENCE_FLOW && el.GetIntent() == eventv1.Element_INTENT_SEQUENCE_FLOW_TAKEN {
			sawFlow = true
		}
		if el.GetType() == eventv1.Element_TYPE_PROCESS && el.GetIntent() == eventv1.Element_INTENT_COMPLETED {
			sawProcessCompleted = true
		}
	}
	if !sawUserActivated || !sawFlow || !sawProcessCompleted {
		t.Fatalf("missing expected events: activated=%v flow=%v processCompleted=%v", sawUserActivated, sawFlow, sawProcessCompleted)
	}
}

func TestCompleteUserTaskInvalidState(t *testing.T) {
	xml, err := os.ReadFile(filepath.Join("testdata", "m1_simple.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
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
	err = eng.CompleteUserTask(ctx, instanceID, "UserTask_1", "missing-token", nil)
	if err == nil {
		t.Fatal("expected error")
	}
}
