package processing_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	"github.com/sparrow-community/sparrow/processing/deploy"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func startIncidentService(t *testing.T) (*processing.Engine, string) {
	t.Helper()
	xml, err := os.ReadFile(filepath.Join("testdata", "m7_incident_service.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	depID, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, depID, nil)
	if err != nil {
		t.Fatal(err)
	}
	return eng, instanceID
}

func waitingServiceJob(t *testing.T, eng *processing.Engine, instanceID string) processing.Job {
	t.Helper()
	jobs, err := eng.Activate(context.Background(), processing.ActivateRequest{
		JobType:  "work.v1",
		WorkerID: "w1",
	})
	if err != nil || len(jobs) != 1 {
		t.Fatalf("activate=%#v err=%v", jobs, err)
	}
	return jobs[0]
}

func serviceToken(t *testing.T, eng *processing.Engine, instanceID string) *projection.Token {
	t.Helper()
	inst, ok := eng.GetInstance(instanceID)
	if !ok {
		t.Fatal("instance not found")
	}
	for _, tok := range inst.Tokens {
		if tok.ElementID == "ServiceTask_1" {
			return tok
		}
	}
	t.Fatalf("tokens=%#v", inst.Tokens)
	return nil
}

func openIncidentOnServiceTask(t *testing.T, eng *processing.Engine, instanceID string) (elementID, tokenID string) {
	t.Helper()
	job := waitingServiceJob(t, eng, instanceID)
	for i := 0; i < deploy.DefaultIncidentThreshold; i++ {
		if i > 0 {
			job = waitingServiceJob(t, eng, instanceID)
		}
		if err := eng.Fail(context.Background(), job.ProcessInstanceID, job.ElementID, job.TokenID, "err", false); err != nil {
			t.Fatal(err)
		}
	}
	tok := serviceToken(t, eng, instanceID)
	if tok.Status != projection.TokenBlocked {
		t.Fatalf("status=%s want blocked", tok.Status)
	}
	return tok.ElementID, tok.ID
}

func countIncidentIntent(t *testing.T, eng *processing.Engine, instanceID string, typ eventv1.Element_Type, intent eventv1.Element_Intent) int {
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

func TestIncidentSubThresholdFailRetriable(t *testing.T) {
	eng, instanceID := startIncidentService(t)
	ctx := context.Background()
	job := waitingServiceJob(t, eng, instanceID)

	if err := eng.Fail(ctx, job.ProcessInstanceID, job.ElementID, job.TokenID, "temporary", false); err != nil {
		t.Fatal(err)
	}
	tok := serviceToken(t, eng, instanceID)
	if tok.Status != projection.TokenWaiting || tok.JobFailCount != 1 {
		t.Fatalf("token=%#v", tok)
	}

	job = waitingServiceJob(t, eng, instanceID)
	if err := eng.Complete(ctx, job.ProcessInstanceID, job.ElementID, job.TokenID, map[string]any{"ok": true}); err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s", inst.Status)
	}
}

func TestIncidentOpensAtThreshold(t *testing.T) {
	eng, instanceID := startIncidentService(t)
	ctx := context.Background()
	elementID, tokenID := openIncidentOnServiceTask(t, eng, instanceID)

	tok := serviceToken(t, eng, instanceID)
	if tok.IncidentErrorMessage != "err" || tok.JobFailCount != int32(deploy.DefaultIncidentThreshold) {
		t.Fatalf("token=%#v", tok)
	}
	if n := countIncidentIntent(t, eng, instanceID, eventv1.Element_TYPE_SERVICE_TASK, eventv1.Element_INTENT_FAILED); n != deploy.DefaultIncidentThreshold {
		t.Fatalf("FAILED count=%d", n)
	}
	if n := countIncidentIntent(t, eng, instanceID, eventv1.Element_TYPE_SERVICE_TASK, eventv1.Element_INTENT_INCIDENT_OPENED); n != 1 {
		t.Fatalf("INCIDENT_OPENED count=%d", n)
	}

	jobs, err := eng.Activate(ctx, processing.ActivateRequest{JobType: "work.v1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs) != 0 {
		t.Fatalf("blocked job activated: %#v", jobs)
	}

	err = eng.Complete(ctx, instanceID, elementID, tokenID, nil)
	if err == nil || !strings.Contains(err.Error(), "INCIDENT_OPEN") {
		t.Fatalf("complete err=%v", err)
	}
	tok = serviceToken(t, eng, instanceID)
	if tok.Status != projection.TokenBlocked {
		t.Fatalf("token changed after rejected complete: %#v", tok)
	}
}

func TestIncidentNoRetryImmediate(t *testing.T) {
	eng, instanceID := startIncidentService(t)
	ctx := context.Background()
	job := waitingServiceJob(t, eng, instanceID)

	if err := eng.Fail(ctx, job.ProcessInstanceID, job.ElementID, job.TokenID, "fatal", true); err != nil {
		t.Fatal(err)
	}
	tok := serviceToken(t, eng, instanceID)
	if tok.Status != projection.TokenBlocked || tok.JobFailCount != 1 {
		t.Fatalf("token=%#v", tok)
	}
	if n := countIncidentIntent(t, eng, instanceID, eventv1.Element_TYPE_SERVICE_TASK, eventv1.Element_INTENT_INCIDENT_OPENED); n != 1 {
		t.Fatalf("INCIDENT_OPENED count=%d", n)
	}
}

func TestIncidentResolveAndComplete(t *testing.T) {
	eng, instanceID := startIncidentService(t)
	ctx := context.Background()
	elementID, tokenID := openIncidentOnServiceTask(t, eng, instanceID)

	if err := eng.ResolveIncident(ctx, instanceID, elementID, tokenID); err != nil {
		t.Fatal(err)
	}
	tok := serviceToken(t, eng, instanceID)
	if tok.Status != projection.TokenWaiting || tok.JobFailCount != 0 || tok.IncidentErrorMessage != "" {
		t.Fatalf("token=%#v", tok)
	}
	if n := countIncidentIntent(t, eng, instanceID, eventv1.Element_TYPE_SERVICE_TASK, eventv1.Element_INTENT_INCIDENT_RESOLVED); n != 1 {
		t.Fatalf("INCIDENT_RESOLVED count=%d", n)
	}

	job := waitingServiceJob(t, eng, instanceID)
	if err := eng.Complete(ctx, job.ProcessInstanceID, job.ElementID, job.TokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s", inst.Status)
	}
}

func TestIncidentResolveWithoutOpenRejected(t *testing.T) {
	eng, instanceID := startIncidentService(t)
	ctx := context.Background()
	job := waitingServiceJob(t, eng, instanceID)

	err := eng.ResolveIncident(ctx, job.ProcessInstanceID, job.ElementID, job.TokenID)
	if err == nil || !strings.Contains(err.Error(), "NO_INCIDENT") {
		t.Fatalf("err=%v", err)
	}
}

func TestIncidentFailOnUserTaskNoIncident(t *testing.T) {
	xml, err := os.ReadFile(filepath.Join("testdata", "m1_simple.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	depID, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, depID, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	var elementID, tokenID string
	for _, tok := range inst.Tokens {
		if tok.Status == projection.TokenWaiting {
			elementID, tokenID = tok.ElementID, tok.ID
			break
		}
	}
	err = eng.Fail(ctx, instanceID, elementID, tokenID, "nope", false)
	if err == nil || !strings.Contains(err.Error(), "INVALID_STATE") {
		t.Fatalf("err=%v", err)
	}
}

func TestIncidentResolveOnCompletedInstanceRejected(t *testing.T) {
	eng, instanceID := startIncidentService(t)
	ctx := context.Background()
	job := waitingServiceJob(t, eng, instanceID)
	if err := eng.Complete(ctx, job.ProcessInstanceID, job.ElementID, job.TokenID, nil); err != nil {
		t.Fatal(err)
	}
	err := eng.ResolveIncident(ctx, instanceID, job.ElementID, job.TokenID)
	if err == nil || !strings.Contains(err.Error(), "INVALID_STATE") {
		t.Fatalf("err=%v", err)
	}
}

func TestIncidentBoundarySupersedesBlockedHost(t *testing.T) {
	xml, err := os.ReadFile(filepath.Join("testdata", "m7_incident_boundary.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	depID, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, depID, nil)
	if err != nil {
		t.Fatal(err)
	}

	job := waitingServiceJob(t, eng, instanceID)
	for i := 0; i < deploy.DefaultIncidentThreshold; i++ {
		if i > 0 {
			job = waitingServiceJob(t, eng, instanceID)
		}
		if err := eng.Fail(ctx, job.ProcessInstanceID, job.ElementID, job.TokenID, "err", false); err != nil {
			t.Fatal(err)
		}
	}
	tok := serviceToken(t, eng, instanceID)
	if tok.Status != projection.TokenBlocked {
		t.Fatalf("token=%#v", tok)
	}

	if err := eng.FireDue(ctx); err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s tokens=%#v", inst.Status, inst.Tokens)
	}
	for _, tok := range inst.Tokens {
		if tok.Status == projection.TokenBlocked {
			t.Fatalf("stray blocked token=%#v", tok)
		}
	}
}
