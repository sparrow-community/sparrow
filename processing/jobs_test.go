package processing_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sparrow-community/sparrow/processing"
	"github.com/sparrow-community/sparrow/processing/deploy"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func TestActivateClaimsServiceTaskJob(t *testing.T) {
	eng, instanceID := startServiceTask(t, map[string]any{"seed": "in"})
	ctx := context.Background()

	jobs, err := eng.Activate(ctx, processing.ActivateRequest{
		JobType:  "work.v1",
		WorkerID: "worker-a",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs) != 1 {
		t.Fatalf("jobs=%d want 1", len(jobs))
	}
	job := jobs[0]
	if job.ProcessInstanceID != instanceID || job.ElementID != "ServiceTask_1" || job.TokenID == "" {
		t.Fatalf("job=%#v", job)
	}
	if job.JobType != "work.v1" || job.WorkerID != "worker-a" || job.LockDeadline.IsZero() {
		t.Fatalf("job=%#v", job)
	}
	if job.Variables["seed"] == "" {
		t.Fatalf("vars=%#v", job.Variables)
	}

	inst, _ := eng.GetInstance(instanceID)
	var tok *projection.Token
	for _, t0 := range inst.Tokens {
		if t0.ID == job.TokenID {
			tok = t0
			break
		}
	}
	if tok == nil || tok.JobType != "work.v1" || tok.Status != projection.TokenWaiting {
		t.Fatalf("token=%#v", tok)
	}

	again, err := eng.Activate(ctx, processing.ActivateRequest{JobType: "work.v1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(again) != 0 {
		t.Fatalf("leased job returned again: %#v", again)
	}

	if err := eng.Complete(ctx, job.ProcessInstanceID, job.ElementID, job.TokenID, map[string]any{"result": "ok"}); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s", inst.Status)
	}
	after, err := eng.Activate(ctx, processing.ActivateRequest{JobType: "work.v1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(after) != 0 {
		t.Fatalf("completed job still listed: %#v", after)
	}
}

func TestActivateIgnoresUserTaskAndWrongType(t *testing.T) {
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
	if _, err := eng.CreateInstance(ctx, dep, nil); err != nil {
		t.Fatal(err)
	}
	jobs, err := eng.Activate(ctx, processing.ActivateRequest{JobType: "work.v1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs) != 0 {
		t.Fatalf("user task must not be a job: %#v", jobs)
	}

	eng2, _ := startServiceTask(t, nil)
	jobs, err = eng2.Activate(ctx, processing.ActivateRequest{JobType: "other.v1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs) != 0 {
		t.Fatalf("wrong type: %#v", jobs)
	}
}

func TestActivateRequiresJobType(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	_, err := eng.Activate(context.Background(), processing.ActivateRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestActivateReclaimsExpiredLease(t *testing.T) {
	eng, instanceID := startServiceTask(t, nil)
	ctx := context.Background()
	first, err := eng.Activate(ctx, processing.ActivateRequest{
		JobType:      "work.v1",
		WorkerID:     "w1",
		LockDuration: 20 * time.Millisecond,
	})
	if err != nil || len(first) != 1 {
		t.Fatalf("first=%#v err=%v", first, err)
	}
	time.Sleep(40 * time.Millisecond)
	second, err := eng.Activate(ctx, processing.ActivateRequest{
		JobType:  "work.v1",
		WorkerID: "w2",
	})
	if err != nil || len(second) != 1 {
		t.Fatalf("second=%#v err=%v", second, err)
	}
	if second[0].ProcessInstanceID != instanceID || second[0].WorkerID != "w2" {
		t.Fatalf("job=%#v", second[0])
	}
	if second[0].TokenID != first[0].TokenID {
		t.Fatalf("token changed first=%s second=%s", first[0].TokenID, second[0].TokenID)
	}
}

func TestActivateLongPollWakesOnServiceTask(t *testing.T) {
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

	errCh := make(chan error, 1)
	go func() {
		jobs, err := eng.Activate(ctx, processing.ActivateRequest{
			JobType: "work.v1",
			Wait:    2 * time.Second,
		})
		if err != nil {
			errCh <- err
			return
		}
		if len(jobs) != 1 || jobs[0].ElementID != "ServiceTask_1" {
			errCh <- fmt.Errorf("jobs=%#v", jobs)
			return
		}
		errCh <- nil
	}()

	time.Sleep(30 * time.Millisecond)
	if _, err := eng.CreateInstance(ctx, dep, nil); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Activate did not return")
	}
}

func TestRecoverThenActivateServiceTask(t *testing.T) {
	xml, err := os.ReadFile(filepath.Join("testdata", "m2_service_task.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	memLog := eventlog.NewMemory()
	memStore := deploy.NewMemoryStore()
	eng1, err := processing.Recover(ctx, memLog, memStore)
	if err != nil {
		t.Fatal(err)
	}
	dep, err := eng1.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng1.CreateInstance(ctx, dep, map[string]any{"k": "v"})
	if err != nil {
		t.Fatal(err)
	}

	eng2, err := processing.Recover(ctx, memLog, memStore)
	if err != nil {
		t.Fatal(err)
	}
	jobs, err := eng2.Activate(ctx, processing.ActivateRequest{JobType: "work.v1", WorkerID: "after-replay"})
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs) != 1 {
		t.Fatalf("jobs=%#v", jobs)
	}
	if jobs[0].ProcessInstanceID != instanceID || jobs[0].Variables["k"] == "" {
		t.Fatalf("job=%#v", jobs[0])
	}
	if err := eng2.Complete(ctx, jobs[0].ProcessInstanceID, jobs[0].ElementID, jobs[0].TokenID, nil); err != nil {
		t.Fatal(err)
	}
}

func TestFailReleasesLeaseAndKeepsWaiting(t *testing.T) {
	eng, instanceID := startServiceTask(t, nil)
	ctx := context.Background()
	jobs, err := eng.Activate(ctx, processing.ActivateRequest{JobType: "work.v1", WorkerID: "w1"})
	if err != nil || len(jobs) != 1 {
		t.Fatalf("activate=%#v err=%v", jobs, err)
	}
	job := jobs[0]
	if err := eng.Fail(ctx, job.ProcessInstanceID, job.ElementID, job.TokenID, "boom"); err != nil {
		t.Fatal(err)
	}

	inst, ok := eng.GetInstance(instanceID)
	if !ok || inst.Status != projection.StatusActive {
		t.Fatalf("inst=%v ok=%v", inst, ok)
	}
	tok := inst.Tokens[job.TokenID]
	if tok == nil || tok.Status != projection.TokenWaiting || tok.JobType != "work.v1" {
		t.Fatalf("token after fail=%#v", tok)
	}

	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	var sawFail bool
	for _, ev := range events {
		if ev.GetRecordType() != eventv1.Event_RECORD_TYPE_EVENT {
			continue
		}
		el := ev.GetElement()
		if el != nil && el.GetIntent() == eventv1.Element_INTENT_FAILED {
			if el.GetActivityPayload().GetErrorMessage() != "boom" {
				t.Fatalf("error_message=%q", el.GetActivityPayload().GetErrorMessage())
			}
			sawFail = true
		}
	}
	if !sawFail {
		t.Fatal("missing FAILED event")
	}

	again, err := eng.Activate(ctx, processing.ActivateRequest{JobType: "work.v1", WorkerID: "w2"})
	if err != nil || len(again) != 1 {
		t.Fatalf("reclaim=%#v err=%v", again, err)
	}
	if again[0].TokenID != job.TokenID || again[0].WorkerID != "w2" {
		t.Fatalf("job=%#v", again[0])
	}
	if err := eng.Complete(ctx, again[0].ProcessInstanceID, again[0].ElementID, again[0].TokenID, nil); err != nil {
		t.Fatal(err)
	}
}

func TestFailRejectsUserTask(t *testing.T) {
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
	var elementID, tokenID string
	for _, tok := range inst.Tokens {
		if tok.Status == projection.TokenWaiting {
			elementID, tokenID = tok.ElementID, tok.ID
			break
		}
	}
	err = eng.Fail(ctx, instanceID, elementID, tokenID, "nope")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFailWakesWaitingActivate(t *testing.T) {
	eng, _ := startServiceTask(t, nil)
	ctx := context.Background()
	held, err := eng.Activate(ctx, processing.ActivateRequest{JobType: "work.v1", WorkerID: "holder"})
	if err != nil || len(held) != 1 {
		t.Fatalf("held=%#v err=%v", held, err)
	}

	errCh := make(chan error, 1)
	go func() {
		jobs, err := eng.Activate(ctx, processing.ActivateRequest{
			JobType:  "work.v1",
			Wait:     2 * time.Second,
			WorkerID: "waiter",
		})
		if err != nil {
			errCh <- err
			return
		}
		if len(jobs) != 1 || jobs[0].TokenID != held[0].TokenID || jobs[0].WorkerID != "waiter" {
			errCh <- fmt.Errorf("jobs=%#v", jobs)
			return
		}
		errCh <- nil
	}()

	time.Sleep(30 * time.Millisecond)
	if err := eng.Fail(ctx, held[0].ProcessInstanceID, held[0].ElementID, held[0].TokenID, "retry"); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Activate did not wake after Fail")
	}
}

func TestHeartbeatExtendsLease(t *testing.T) {
	eng, _ := startServiceTask(t, nil)
	ctx := context.Background()
	jobs, err := eng.Activate(ctx, processing.ActivateRequest{
		JobType:      "work.v1",
		WorkerID:     "w1",
		LockDuration: 40 * time.Millisecond,
	})
	if err != nil || len(jobs) != 1 {
		t.Fatalf("activate=%#v err=%v", jobs, err)
	}
	if err := eng.Heartbeat(ctx, jobs[0].ProcessInstanceID, jobs[0].TokenID, "w1", time.Second); err != nil {
		t.Fatal(err)
	}
	time.Sleep(60 * time.Millisecond)
	again, err := eng.Activate(ctx, processing.ActivateRequest{JobType: "work.v1", WorkerID: "w2"})
	if err != nil {
		t.Fatal(err)
	}
	if len(again) != 0 {
		t.Fatalf("lease should still be held: %#v", again)
	}
	if err := eng.Heartbeat(ctx, jobs[0].ProcessInstanceID, jobs[0].TokenID, "w2", time.Second); err == nil {
		t.Fatal("expected other worker heartbeat to fail")
	}
}

func startServiceTask(t *testing.T, vars map[string]any) (*processing.Engine, string) {
	t.Helper()
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
	instanceID, err := eng.CreateInstance(ctx, dep, vars)
	if err != nil {
		t.Fatal(err)
	}
	return eng, instanceID
}
