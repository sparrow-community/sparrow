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

	if err := eng.Complete(ctx, instanceID, waitingElement, waitingToken, map[string]any{"approved": true}); err != nil {
		t.Fatalf("Complete: %v", err)
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

func TestCompleteInvalidState(t *testing.T) {
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
	err = eng.Complete(ctx, instanceID, "UserTask_1", "missing-token", nil)
	if err == nil {
		t.Fatal("expected error")
	}

	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	var sawRejection bool
	for _, ev := range events {
		if ev.GetRecordType() == eventv1.Event_RECORD_TYPE_REJECTION {
			sawRejection = true
			if ev.GetRejection() == nil || ev.GetRejection().GetCode() != "INVALID_STATE" {
				t.Fatalf("rejection=%v", ev.GetRejection())
			}
		}
	}
	if !sawRejection {
		t.Fatal("expected REJECTION record in event log")
	}
}

func TestDeployRejectsUnsupportedElement(t *testing.T) {
	xml, err := os.ReadFile(filepath.Join("testdata", "m1_unsupported_inclusive.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	eng := processing.NewEngine(eventlog.NewMemory())
	_, err = eng.Deploy(context.Background(), xml)
	if err == nil {
		t.Fatal("expected deploy to reject unsupported element")
	}
}

func TestExclusiveGatewayTakesDefault(t *testing.T) {
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
	inst, _ := eng.GetInstance(instanceID)
	var waitingToken, waitingElement string
	for _, tok := range inst.Tokens {
		if tok.Status == projection.TokenWaiting {
			waitingToken, waitingElement = tok.ID, tok.ElementID
			break
		}
	}
	if err := eng.Complete(ctx, instanceID, waitingElement, waitingToken, nil); err != nil {
		t.Fatal(err)
	}

	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	var tookDefault, reachedOK bool
	for _, ev := range events {
		if ev.GetRecordType() != eventv1.Event_RECORD_TYPE_EVENT {
			continue
		}
		el := ev.GetElement()
		if el == nil {
			continue
		}
		if el.GetType() == eventv1.Element_TYPE_EXCLUSIVE_GATEWAY &&
			el.GetIntent() == eventv1.Element_INTENT_COMPLETED {
			gp := el.GetGatewayPayload()
			if gp == nil || gp.GetTakenSequenceFlowId() != "Flow_gw_to_end_ok" {
				t.Fatalf("XOR taken=%v want Flow_gw_to_end_ok", gp)
			}
			tookDefault = true
		}
		if el.GetType() == eventv1.Element_TYPE_END_EVENT &&
			el.GetId() == "EndEvent_ok" &&
			el.GetIntent() == eventv1.Element_INTENT_COMPLETED {
			reachedOK = true
		}
	}
	if !tookDefault || !reachedOK {
		t.Fatalf("default route missing: tookDefault=%v reachedOK=%v", tookDefault, reachedOK)
	}
}

func TestExclusiveGatewayTakesCondition(t *testing.T) {
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
	inst, _ := eng.GetInstance(instanceID)
	var waitingToken, waitingElement string
	for _, tok := range inst.Tokens {
		if tok.Status == projection.TokenWaiting {
			waitingToken, waitingElement = tok.ID, tok.ElementID
			break
		}
	}
	if err := eng.Complete(ctx, instanceID, waitingElement, waitingToken, map[string]any{"approved": false}); err != nil {
		t.Fatal(err)
	}

	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	var tookOther, reachedOther bool
	for _, ev := range events {
		if ev.GetRecordType() != eventv1.Event_RECORD_TYPE_EVENT {
			continue
		}
		el := ev.GetElement()
		if el == nil {
			continue
		}
		if el.GetType() == eventv1.Element_TYPE_EXCLUSIVE_GATEWAY &&
			el.GetIntent() == eventv1.Element_INTENT_COMPLETED {
			gp := el.GetGatewayPayload()
			if gp == nil || gp.GetTakenSequenceFlowId() != "Flow_gw_to_end_other" {
				t.Fatalf("XOR taken=%v want Flow_gw_to_end_other", gp)
			}
			tookOther = true
		}
		if el.GetType() == eventv1.Element_TYPE_END_EVENT &&
			el.GetId() == "EndEvent_other" &&
			el.GetIntent() == eventv1.Element_INTENT_COMPLETED {
			reachedOther = true
		}
	}
	if !tookOther || !reachedOther {
		t.Fatalf("conditioned route missing: tookOther=%v reachedOther=%v", tookOther, reachedOther)
	}
}

func TestGetInstanceReturnsSnapshot(t *testing.T) {
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
	a, ok := eng.GetInstance(instanceID)
	if !ok {
		t.Fatal("missing")
	}
	a.Status = projection.StatusTerminated
	for _, tok := range a.Tokens {
		tok.Status = "mutated"
	}
	b, _ := eng.GetInstance(instanceID)
	if b.Status != projection.StatusActive {
		t.Fatalf("snapshot leaked, status=%s", b.Status)
	}
	for _, tok := range b.Tokens {
		if tok.Status == "mutated" {
			t.Fatal("token snapshot leaked")
		}
	}
}

func TestServiceTaskWaitAndComplete(t *testing.T) {
	xml, err := os.ReadFile(filepath.Join("testdata", "m2_service_task.bpmn"))
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
	inst, ok := eng.GetInstance(instanceID)
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
	if waitingElement != "ServiceTask_1" {
		t.Fatalf("expected wait at ServiceTask_1, tokens=%#v", inst.Tokens)
	}
	if inst.Tokens[waitingToken].JobType != "work.v1" {
		t.Fatalf("token job_type=%q", inst.Tokens[waitingToken].JobType)
	}

	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	var sawJobType bool
	for _, ev := range events {
		if ev.GetRecordType() != eventv1.Event_RECORD_TYPE_EVENT {
			continue
		}
		el := ev.GetElement()
		if el != nil && el.GetType() == eventv1.Element_TYPE_SERVICE_TASK &&
			el.GetIntent() == eventv1.Element_INTENT_ACTIVATED {
			if el.GetActivityPayload().GetJobType() != "work.v1" {
				t.Fatalf("job_type=%q", el.GetActivityPayload().GetJobType())
			}
			sawJobType = true
		}
	}
	if !sawJobType {
		t.Fatal("missing SERVICE_TASK ACTIVATED with job_type")
	}

	if err := eng.Complete(ctx, instanceID, waitingElement, waitingToken, map[string]any{"result": "ok"}); err != nil {
		t.Fatalf("Complete: %v", err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	if inst.Variables["result"] == "" {
		t.Fatalf("vars=%#v", inst.Variables)
	}
}
