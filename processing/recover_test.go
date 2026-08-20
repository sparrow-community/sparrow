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
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func TestRecoverRedrivesCreateInstanceCommand(t *testing.T) {
	xml := readM1(t)
	ctx := context.Background()
	memLog := eventlog.NewMemory()
	store := deploy.NewMemoryStore()
	eng1, err := processing.Recover(ctx, memLog, store, nil)
	if err != nil {
		t.Fatal(err)
	}
	depID, err := eng1.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}

	instanceID := processing.MustNextID()
	tokenID := processing.MustNextID()
	cmdID := processing.MustNextID()
	if _, err := memLog.Append(ctx, &eventv1.Event{
		Id:                cmdID,
		Timestamp:         1,
		RecordType:        eventv1.Event_RECORD_TYPE_COMMAND,
		DeploymentId:      depID,
		ProcessInstanceId: instanceID,
		ProcessVersion:    1,
		Element: &eventv1.Element{
			Intent:  eventv1.Element_INTENT_ACTIVATING,
			Type:    eventv1.Element_TYPE_PROCESS,
			Id:      "Process_m1",
			TokenId: tokenID,
		},
	}); err != nil {
		t.Fatal(err)
	}

	eng2, err := processing.Recover(ctx, memLog, store, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, ok := eng2.GetInstance(instanceID)
	if !ok || inst.Status != projection.StatusActive {
		t.Fatalf("inst=%v ok=%v", inst, ok)
	}
	if waitingElement(inst) != "UserTask_1" {
		t.Fatalf("tokens=%#v", inst.Tokens)
	}
}

func TestRecoverRedrivesCompleteCommand(t *testing.T) {
	xml := readM1(t)
	ctx := context.Background()
	memLog := eventlog.NewMemory()
	store := deploy.NewMemoryStore()
	eng1, err := processing.Recover(ctx, memLog, store, nil)
	if err != nil {
		t.Fatal(err)
	}
	depID, err := eng1.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng1.CreateInstance(ctx, depID, map[string]any{"approved": true})
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng1.GetInstance(instanceID)
	elementID, tokenID := waitingAt(inst)

	cmdID := processing.MustNextID()
	if _, err := memLog.Append(ctx, &eventv1.Event{
		Id:                cmdID,
		Timestamp:         1,
		RecordType:        eventv1.Event_RECORD_TYPE_COMMAND,
		DeploymentId:      depID,
		ProcessInstanceId: instanceID,
		ProcessVersion:    1,
		Element: &eventv1.Element{
			Intent:  eventv1.Element_INTENT_COMPLETING,
			Type:    eventv1.Element_TYPE_USER_TASK,
			Id:      elementID,
			TokenId: tokenID,
			Payload: &eventv1.Element_ActivityPayload{
				ActivityPayload: &eventv1.ActivityPayload{},
			},
		},
	}); err != nil {
		t.Fatal(err)
	}

	eng2, err := processing.Recover(ctx, memLog, store, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ = eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestRecoverRedrivesPartialCompleteChain(t *testing.T) {
	xml := readM1(t)
	ctx := context.Background()
	memLog := eventlog.NewMemory()
	store := deploy.NewMemoryStore()
	eng1, err := processing.Recover(ctx, memLog, store, nil)
	if err != nil {
		t.Fatal(err)
	}
	depID, err := eng1.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng1.CreateInstance(ctx, depID, map[string]any{"approved": true})
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng1.GetInstance(instanceID)
	elementID, tokenID := waitingAt(inst)

	cmdID := processing.MustNextID()
	if _, err := memLog.Append(ctx, &eventv1.Event{
		Id:                cmdID,
		Timestamp:         1,
		RecordType:        eventv1.Event_RECORD_TYPE_COMMAND,
		DeploymentId:      depID,
		ProcessInstanceId: instanceID,
		ProcessVersion:    1,
		Element: &eventv1.Element{
			Intent:  eventv1.Element_INTENT_COMPLETING,
			Type:    eventv1.Element_TYPE_USER_TASK,
			Id:      elementID,
			TokenId: tokenID,
		},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := memLog.Append(ctx, &eventv1.Event{
		Id:                processing.MustNextID(),
		Timestamp:         2,
		RecordType:        eventv1.Event_RECORD_TYPE_EVENT,
		DeploymentId:      depID,
		ProcessInstanceId: instanceID,
		ProcessVersion:    1,
		SourceRecordId:    cmdID,
		Element: &eventv1.Element{
			Intent:  eventv1.Element_INTENT_COMPLETING,
			Type:    eventv1.Element_TYPE_USER_TASK,
			Id:      elementID,
			TokenId: tokenID,
		},
	}); err != nil {
		t.Fatal(err)
	}

	eng2, err := processing.Recover(ctx, memLog, store, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ = eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	if n := countIntent(t, eng2, instanceID, eventv1.Element_TYPE_PROCESS, eventv1.Element_INTENT_COMPLETED); n != 1 {
		t.Fatalf("PROCESS COMPLETED count=%d", n)
	}
}

func TestRecoverDoesNotDuplicateFinishedCommands(t *testing.T) {
	xml := readM1(t)
	ctx := context.Background()
	memLog := eventlog.NewMemory()
	store := deploy.NewMemoryStore()
	eng1, err := processing.Recover(ctx, memLog, store, nil)
	if err != nil {
		t.Fatal(err)
	}
	depID, err := eng1.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng1.CreateInstance(ctx, depID, map[string]any{"approved": true})
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng1.GetInstance(instanceID)
	elementID, tokenID := waitingAt(inst)
	if err := eng1.Complete(ctx, instanceID, elementID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	before, err := eng1.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}

	eng2, err := processing.Recover(ctx, memLog, store, nil)
	if err != nil {
		t.Fatal(err)
	}
	after, err := eng2.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if len(after) != len(before) {
		t.Fatalf("events before=%d after=%d", len(before), len(after))
	}
	inst, _ = eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s", inst.Status)
	}
}

func TestRecoverRedrivesFailCommand(t *testing.T) {
	xml, err := os.ReadFile(filepath.Join("testdata", "m2_service_task.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	memLog := eventlog.NewMemory()
	store := deploy.NewMemoryStore()
	eng1, err := processing.Recover(ctx, memLog, store, nil)
	if err != nil {
		t.Fatal(err)
	}
	depID, err := eng1.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng1.CreateInstance(ctx, depID, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng1.GetInstance(instanceID)
	elementID, tokenID := waitingAt(inst)

	cmdID := processing.MustNextID()
	if _, err := memLog.Append(ctx, &eventv1.Event{
		Id:                cmdID,
		Timestamp:         1,
		RecordType:        eventv1.Event_RECORD_TYPE_COMMAND,
		DeploymentId:      depID,
		ProcessInstanceId: instanceID,
		ProcessVersion:    1,
		Element: &eventv1.Element{
			Intent:  eventv1.Element_INTENT_FAILED,
			Type:    eventv1.Element_TYPE_SERVICE_TASK,
			Id:      elementID,
			TokenId: tokenID,
			Payload: &eventv1.Element_ActivityPayload{
				ActivityPayload: &eventv1.ActivityPayload{JobType: "work.v1", ErrorMessage: "down"},
			},
		},
	}); err != nil {
		t.Fatal(err)
	}

	eng2, err := processing.Recover(ctx, memLog, store, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ = eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusActive || waitingElement(inst) != "ServiceTask_1" {
		t.Fatalf("inst=%#v", inst)
	}
	if n := countIntent(t, eng2, instanceID, eventv1.Element_TYPE_SERVICE_TASK, eventv1.Element_INTENT_FAILED); n != 1 {
		t.Fatalf("FAILED count=%d", n)
	}
	jobs, err := eng2.Activate(ctx, processing.ActivateRequest{JobType: "work.v1"})
	if err != nil || len(jobs) != 1 {
		t.Fatalf("activate=%#v err=%v", jobs, err)
	}
}

func readM1(t *testing.T) []byte {
	t.Helper()
	xml, err := os.ReadFile(filepath.Join("testdata", "m1_simple.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	return xml
}

func waitingAt(inst *projection.Instance) (elementID, tokenID string) {
	for _, tok := range inst.Tokens {
		if tok.Status == projection.TokenWaiting {
			return tok.ElementID, tok.ID
		}
	}
	return "", ""
}

func waitingElement(inst *projection.Instance) string {
	el, _ := waitingAt(inst)
	return el
}

func countIntent(t *testing.T, eng *processing.Engine, instanceID string, typ eventv1.Element_Type, intent eventv1.Element_Intent) int {
	t.Helper()
	events, err := eng.ListEvents(context.Background(), instanceID)
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, ev := range events {
		if ev.GetRecordType() != eventv1.Event_RECORD_TYPE_EVENT {
			continue
		}
		el := ev.GetElement()
		if el.GetType() == typ && el.GetIntent() == intent {
			n++
		}
	}
	return n
}
