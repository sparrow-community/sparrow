package gateway_test

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sparrow-community/sparrow/gateway"
	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	jobv1 "github.com/sparrow-community/sparrow/protocol/gen/go/job/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

func TestJobServiceActivateComplete(t *testing.T) {
	eng, client, stop := startJobClient(t)
	defer stop()
	ctx := context.Background()

	xml, err := os.ReadFile(filepath.Join("..", "processing", "testdata", "m2_service_task.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, map[string]any{"seed": 1})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.ActivateJobs(ctx, &jobv1.ActivateJobsRequest{
		JobType:  "work.v1",
		WorkerId: "gw-worker",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.GetJobs()) != 1 {
		t.Fatalf("jobs=%v", resp.GetJobs())
	}
	job := resp.GetJobs()[0]
	if job.GetProcessInstanceId() != instanceID || job.GetElementId() != "ServiceTask_1" {
		t.Fatalf("job=%v", job)
	}
	if job.GetVariables()["seed"] == "" || job.GetLockDeadlineUnixMs() == 0 {
		t.Fatalf("job=%v", job)
	}

	_, err = client.CompleteJob(ctx, &jobv1.CompleteJobRequest{
		ProcessInstanceId: job.GetProcessInstanceId(),
		ElementId:         job.GetElementId(),
		TokenId:           job.GetTokenId(),
		Variables:         map[string]string{"result": `"ok"`},
	})
	if err != nil {
		t.Fatal(err)
	}
	inst, ok := eng.GetInstance(instanceID)
	if !ok || inst.Status != projection.StatusCompleted {
		t.Fatalf("inst=%v ok=%v", inst, ok)
	}
}

func TestJobServiceFailThenActivate(t *testing.T) {
	eng, client, stop := startJobClient(t)
	defer stop()
	ctx := context.Background()

	xml, err := os.ReadFile(filepath.Join("..", "processing", "testdata", "m2_service_task.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := eng.CreateInstance(ctx, dep, nil); err != nil {
		t.Fatal(err)
	}

	first, err := client.ActivateJobs(ctx, &jobv1.ActivateJobsRequest{JobType: "work.v1", WorkerId: "a"})
	if err != nil || len(first.GetJobs()) != 1 {
		t.Fatalf("first=%v err=%v", first, err)
	}
	job := first.GetJobs()[0]
	if _, err := client.FailJob(ctx, &jobv1.FailJobRequest{
		ProcessInstanceId: job.GetProcessInstanceId(),
		ElementId:         job.GetElementId(),
		TokenId:           job.GetTokenId(),
		ErrorMessage:      "remote down",
	}); err != nil {
		t.Fatal(err)
	}
	second, err := client.ActivateJobs(ctx, &jobv1.ActivateJobsRequest{JobType: "work.v1", WorkerId: "b"})
	if err != nil || len(second.GetJobs()) != 1 {
		t.Fatalf("second=%v err=%v", second, err)
	}
	if second.GetJobs()[0].GetTokenId() != job.GetTokenId() || second.GetJobs()[0].GetWorkerId() != "b" {
		t.Fatalf("job=%v", second.GetJobs()[0])
	}
}

func TestJobServiceHeartbeatAndErrors(t *testing.T) {
	_, client, stop := startJobClient(t)
	defer stop()
	ctx := context.Background()

	_, err := client.ActivateJobs(ctx, &jobv1.ActivateJobsRequest{WorkerId: "x"})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("code=%v err=%v", status.Code(err), err)
	}

	xml, err := os.ReadFile(filepath.Join("..", "processing", "testdata", "m2_service_task.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	eng, client2, stop2 := startJobClient(t)
	defer stop2()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := eng.CreateInstance(ctx, dep, nil); err != nil {
		t.Fatal(err)
	}
	held, err := client2.ActivateJobs(ctx, &jobv1.ActivateJobsRequest{
		JobType:        "work.v1",
		WorkerId:       "w1",
		LockDurationMs: 40,
	})
	if err != nil || len(held.GetJobs()) != 1 {
		t.Fatalf("held=%v err=%v", held, err)
	}
	_, err = client2.Heartbeat(ctx, &jobv1.HeartbeatRequest{
		ProcessInstanceId: held.GetJobs()[0].GetProcessInstanceId(),
		TokenId:           held.GetJobs()[0].GetTokenId(),
		WorkerId:          "w1",
		LockDurationMs:    1000,
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(60 * time.Millisecond)
	again, err := client2.ActivateJobs(ctx, &jobv1.ActivateJobsRequest{JobType: "work.v1", WorkerId: "w2"})
	if err != nil {
		t.Fatal(err)
	}
	if len(again.GetJobs()) != 0 {
		t.Fatalf("lease should hold: %v", again.GetJobs())
	}
	_, err = client2.Heartbeat(ctx, &jobv1.HeartbeatRequest{
		ProcessInstanceId: held.GetJobs()[0].GetProcessInstanceId(),
		TokenId:           held.GetJobs()[0].GetTokenId(),
		WorkerId:          "w2",
	})
	if status.Code(err) != codes.FailedPrecondition {
		t.Fatalf("code=%v err=%v", status.Code(err), err)
	}
}

func startJobClient(t *testing.T) (*processing.Engine, jobv1.JobServiceClient, func()) {
	t.Helper()
	eng := processing.NewEngine(eventlog.NewMemory())
	lis := bufconn.Listen(bufSize)
	srv := gateway.NewServer(eng)
	go func() {
		_ = srv.Serve(lis)
	}()
	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}
	stop := func() {
		_ = conn.Close()
		srv.Stop()
		_ = lis.Close()
	}
	return eng, jobv1.NewJobServiceClient(conn), stop
}
