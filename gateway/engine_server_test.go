package gateway_test

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/sparrow-community/sparrow/gateway"
	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	enginev1 "github.com/sparrow-community/sparrow/protocol/gen/go/engine/v1"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func TestEngineServiceDeployCreateComplete(t *testing.T) {
	eng, conn, stop := startGRPC(t)
	defer stop()
	_ = eng
	ctx := context.Background()
	client := enginev1.NewEngineServiceClient(conn)

	xml, err := os.ReadFile(filepath.Join("..", "processing", "testdata", "m1_simple.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	dep, err := client.Deploy(ctx, &enginev1.DeployRequest{BpmnXml: xml})
	if err != nil {
		t.Fatal(err)
	}
	if dep.GetDeploymentId() == "" {
		t.Fatal("empty deployment id")
	}

	created, err := client.CreateInstance(ctx, &enginev1.CreateInstanceRequest{
		DeploymentId: dep.GetDeploymentId(),
		Variables:    map[string]string{"initiator": `"alice"`},
	})
	if err != nil {
		t.Fatal(err)
	}
	got, err := client.GetInstance(ctx, &enginev1.GetInstanceRequest{
		ProcessInstanceId: created.GetProcessInstanceId(),
	})
	if err != nil {
		t.Fatal(err)
	}
	inst := got.GetInstance()
	if inst.GetStatus() != string(projection.StatusActive) || inst.GetVariables()["initiator"] == "" {
		t.Fatalf("instance=%v", inst)
	}
	var elementID, tokenID string
	for _, tok := range inst.GetTokens() {
		if tok.GetStatus() == string(projection.TokenWaiting) {
			elementID, tokenID = tok.GetElementId(), tok.GetId()
			break
		}
	}
	if elementID != "UserTask_1" {
		t.Fatalf("waiting=%s tokens=%v", elementID, inst.GetTokens())
	}

	if _, err := client.Complete(ctx, &enginev1.CompleteRequest{
		ProcessInstanceId: created.GetProcessInstanceId(),
		ElementId:         elementID,
		TokenId:           tokenID,
		Variables:         map[string]string{"approved": "true"},
	}); err != nil {
		t.Fatal(err)
	}
	got, err = client.GetInstance(ctx, &enginev1.GetInstanceRequest{
		ProcessInstanceId: created.GetProcessInstanceId(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.GetInstance().GetStatus() != string(projection.StatusCompleted) {
		t.Fatalf("status=%s", got.GetInstance().GetStatus())
	}

	events, err := client.ListEvents(ctx, &enginev1.ListEventsRequest{
		ProcessInstanceId: created.GetProcessInstanceId(),
	})
	if err != nil {
		t.Fatal(err)
	}
	var sawProcessCompleted bool
	for _, ev := range events.GetEvents() {
		el := ev.GetElement()
		if ev.GetRecordType() == eventv1.Event_RECORD_TYPE_EVENT &&
			el.GetType() == eventv1.Element_TYPE_PROCESS &&
			el.GetIntent() == eventv1.Element_INTENT_COMPLETED {
			sawProcessCompleted = true
		}
	}
	if !sawProcessCompleted {
		t.Fatal("missing PROCESS COMPLETED")
	}
}

func TestEngineServiceNotFoundAndHealth(t *testing.T) {
	_, conn, stop := startGRPC(t)
	defer stop()
	ctx := context.Background()
	client := enginev1.NewEngineServiceClient(conn)

	_, err := client.GetInstance(ctx, &enginev1.GetInstanceRequest{ProcessInstanceId: "missing"})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("code=%v err=%v", status.Code(err), err)
	}

	hs := grpc_health_v1.NewHealthClient(conn)
	resp, err := hs.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if resp.GetStatus() != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Fatalf("health=%v", resp.GetStatus())
	}
}

func startGRPC(t *testing.T) (*processing.Engine, *grpc.ClientConn, func()) {
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
	return eng, conn, stop
}
