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

func TestRecoverRedrivesThrowErrorCommand(t *testing.T) {
	xml, err := os.ReadFile(filepath.Join("testdata", "m4_error_boundary.bpmn"))
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
			Intent:  eventv1.Element_INTENT_ERROR_THROWN,
			Type:    eventv1.Element_TYPE_USER_TASK,
			Id:      elementID,
			TokenId: tokenID,
			Payload: &eventv1.Element_EventPayload{
				EventPayload: &eventv1.EventPayload{ErrorCode: "BUSINESS_ERROR"},
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
		t.Fatalf("status=%s want completed after redrive", inst.Status)
	}
	if n := countIntent(t, eng2, instanceID, eventv1.Element_TYPE_USER_TASK, eventv1.Element_INTENT_ERROR_THROWN); n != 1 {
		t.Fatalf("ERROR_THROWN count=%d", n)
	}
	if n := countIntent(t, eng2, instanceID, eventv1.Element_TYPE_BOUNDARY_EVENT, eventv1.Element_INTENT_COMPLETED); n != 1 {
		t.Fatalf("boundary COMPLETED count=%d", n)
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
	var hostElem, hostTok string
	for _, tok := range inst.Tokens {
		if tok.Status != projection.TokenWaiting {
			continue
		}
		if tok.ScopeHost {
			if hostTok == "" {
				hostElem, hostTok = tok.ElementID, tok.ID
			}
			continue
		}
		return tok.ElementID, tok.ID
	}
	return hostElem, hostTok
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

func TestRecoverRebuildsBlockedIncident(t *testing.T) {
	xml, err := os.ReadFile(filepath.Join("testdata", "m7_incident_service.bpmn"))
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
	inst1, _ := eng1.GetInstance(instanceID)
	elementID, tokenID := waitingAt(inst1)
	for i := 0; i < deploy.DefaultIncidentThreshold; i++ {
		jobs, err := eng1.Activate(ctx, processing.ActivateRequest{JobType: "work.v1", WorkerID: "w1"})
		if err != nil || len(jobs) != 1 {
			t.Fatalf("activate=%#v err=%v", jobs, err)
		}
		if err := eng1.Fail(ctx, jobs[0].ProcessInstanceID, jobs[0].ElementID, jobs[0].TokenID, "err", false); err != nil {
			t.Fatal(err)
		}
	}

	eng2, err := processing.Recover(ctx, memLog, store, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng2.GetInstance(instanceID)
	tok := inst.Tokens[tokenID]
	if tok == nil || tok.Status != projection.TokenBlocked || tok.IncidentErrorMessage != "err" {
		t.Fatalf("token=%#v", tok)
	}

	if err := eng2.ResolveIncident(ctx, instanceID, elementID, tokenID); err != nil {
		t.Fatal(err)
	}
	jobs, err := eng2.Activate(ctx, processing.ActivateRequest{JobType: "work.v1", WorkerID: "w2"})
	if err != nil || len(jobs) != 1 {
		t.Fatalf("activate=%#v err=%v", jobs, err)
	}
	if err := eng2.Complete(ctx, jobs[0].ProcessInstanceID, jobs[0].ElementID, jobs[0].TokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s", inst.Status)
	}
}

func TestRecoverRedrivesFailIncidentOpen(t *testing.T) {
	xml, err := os.ReadFile(filepath.Join("testdata", "m7_incident_service.bpmn"))
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
				ActivityPayload: &eventv1.ActivityPayload{
					JobType:      "work.v1",
					ErrorMessage: "down",
					JobFailCount: int32(deploy.DefaultIncidentThreshold),
				},
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
	tok := inst.Tokens[tokenID]
	if tok == nil || tok.Status != projection.TokenBlocked {
		t.Fatalf("token=%#v", tok)
	}
	if n := countIntent(t, eng2, instanceID, eventv1.Element_TYPE_SERVICE_TASK, eventv1.Element_INTENT_FAILED); n != 1 {
		t.Fatalf("FAILED count=%d", n)
	}
	if n := countIntent(t, eng2, instanceID, eventv1.Element_TYPE_SERVICE_TASK, eventv1.Element_INTENT_INCIDENT_OPENED); n != 1 {
		t.Fatalf("INCIDENT_OPENED count=%d", n)
	}
}
