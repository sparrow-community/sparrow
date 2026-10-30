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
	jobv1 "github.com/sparrow-community/sparrow/protocol/gen/go/job/v1"
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
	fetched, err := client.GetDeployment(ctx, &enginev1.GetDeploymentRequest{DeploymentId: dep.GetDeploymentId()})
	if err != nil {
		t.Fatal(err)
	}
	if fetched.GetProcessId() == "" || fetched.GetProcessVersion() < 1 {
		t.Fatalf("deployment meta=%v", fetched)
	}
	if string(fetched.GetBpmnXml()) != string(xml) {
		t.Fatalf("bpmn_xml mismatch")
	}
	if _, err := client.GetDeployment(ctx, &enginev1.GetDeploymentRequest{DeploymentId: "missing"}); status.Code(err) != codes.NotFound {
		t.Fatalf("want NotFound, got %v", err)
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

func TestEngineServicePublishMessage(t *testing.T) {
	_, conn, stop := startGRPC(t)
	defer stop()
	ctx := context.Background()
	client := enginev1.NewEngineServiceClient(conn)

	xml, err := os.ReadFile(filepath.Join("..", "processing", "testdata", "m2_message_catch.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	dep, err := client.Deploy(ctx, &enginev1.DeployRequest{BpmnXml: xml})
	if err != nil {
		t.Fatal(err)
	}
	created, err := client.CreateInstance(ctx, &enginev1.CreateInstanceRequest{DeploymentId: dep.GetDeploymentId()})
	if err != nil {
		t.Fatal(err)
	}
	got, err := client.GetInstance(ctx, &enginev1.GetInstanceRequest{ProcessInstanceId: created.GetProcessInstanceId()})
	if err != nil {
		t.Fatal(err)
	}
	var sawName bool
	for _, tok := range got.GetInstance().GetTokens() {
		if tok.GetStatus() == string(projection.TokenWaiting) && tok.GetMessageName() == "order.confirmed" {
			sawName = true
		}
	}
	if !sawName {
		t.Fatalf("tokens=%v", got.GetInstance().GetTokens())
	}

	pub, err := client.PublishMessage(ctx, &enginev1.PublishMessageRequest{
		Name:      "order.confirmed",
		Variables: map[string]string{"payload": `"ok"`},
	})
	if err != nil {
		t.Fatal(err)
	}
	if pub.GetDelivered() != 1 {
		t.Fatalf("delivered=%d", pub.GetDelivered())
	}
	got, err = client.GetInstance(ctx, &enginev1.GetInstanceRequest{ProcessInstanceId: created.GetProcessInstanceId()})
	if err != nil {
		t.Fatal(err)
	}
	if got.GetInstance().GetStatus() != string(projection.StatusCompleted) {
		t.Fatalf("status=%s", got.GetInstance().GetStatus())
	}

	late, err := client.PublishMessage(ctx, &enginev1.PublishMessageRequest{Name: "order.confirmed"})
	if err != nil {
		t.Fatal(err)
	}
	if late.GetDelivered() != 0 || late.GetBuffered() != 1 {
		t.Fatalf("late publish delivered=%d buffered=%d", late.GetDelivered(), late.GetBuffered())
	}
}

func TestEngineServicePublishMessageCorrelation(t *testing.T) {
	_, conn, stop := startGRPC(t)
	defer stop()
	ctx := context.Background()
	client := enginev1.NewEngineServiceClient(conn)

	xml, err := os.ReadFile(filepath.Join("..", "processing", "testdata", "m2_message_catch.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	dep, err := client.Deploy(ctx, &enginev1.DeployRequest{BpmnXml: xml})
	if err != nil {
		t.Fatal(err)
	}
	createdA, err := client.CreateInstance(ctx, &enginev1.CreateInstanceRequest{
		DeploymentId: dep.GetDeploymentId(),
		Variables:    map[string]string{"orderId": `"A"`},
	})
	if err != nil {
		t.Fatal(err)
	}
	createdB, err := client.CreateInstance(ctx, &enginev1.CreateInstanceRequest{
		DeploymentId: dep.GetDeploymentId(),
		Variables:    map[string]string{"orderId": `"B"`},
	})
	if err != nil {
		t.Fatal(err)
	}

	miss, err := client.PublishMessage(ctx, &enginev1.PublishMessageRequest{
		Name:            "order.confirmed",
		CorrelationKeys: map[string]string{"orderId": `"missing"`},
	})
	if err != nil {
		t.Fatal(err)
	}
	if miss.GetDelivered() != 0 || miss.GetBuffered() != 1 {
		t.Fatalf("unmatched delivered=%d buffered=%d", miss.GetDelivered(), miss.GetBuffered())
	}

	pub, err := client.PublishMessage(ctx, &enginev1.PublishMessageRequest{
		Name:            "order.confirmed",
		CorrelationKeys: map[string]string{"orderId": `"A"`},
	})
	if err != nil {
		t.Fatal(err)
	}
	if pub.GetDelivered() != 1 {
		t.Fatalf("delivered=%d", pub.GetDelivered())
	}
	gotA, err := client.GetInstance(ctx, &enginev1.GetInstanceRequest{ProcessInstanceId: createdA.GetProcessInstanceId()})
	if err != nil {
		t.Fatal(err)
	}
	gotB, err := client.GetInstance(ctx, &enginev1.GetInstanceRequest{ProcessInstanceId: createdB.GetProcessInstanceId()})
	if err != nil {
		t.Fatal(err)
	}
	if gotA.GetInstance().GetStatus() != string(projection.StatusCompleted) {
		t.Fatalf("A status=%s", gotA.GetInstance().GetStatus())
	}
	if gotB.GetInstance().GetStatus() != string(projection.StatusActive) {
		t.Fatalf("B status=%s want still active", gotB.GetInstance().GetStatus())
	}
}

func TestEngineServiceThrowError(t *testing.T) {
	_, conn, stop := startGRPC(t)
	defer stop()
	ctx := context.Background()
	client := enginev1.NewEngineServiceClient(conn)

	xml, err := os.ReadFile(filepath.Join("..", "processing", "testdata", "m4_error_boundary.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	dep, err := client.Deploy(ctx, &enginev1.DeployRequest{BpmnXml: xml})
	if err != nil {
		t.Fatal(err)
	}
	created, err := client.CreateInstance(ctx, &enginev1.CreateInstanceRequest{DeploymentId: dep.GetDeploymentId()})
	if err != nil {
		t.Fatal(err)
	}
	got, err := client.GetInstance(ctx, &enginev1.GetInstanceRequest{ProcessInstanceId: created.GetProcessInstanceId()})
	if err != nil {
		t.Fatal(err)
	}
	var elementID, tokenID string
	for _, tok := range got.GetInstance().GetTokens() {
		if tok.GetStatus() == string(projection.TokenWaiting) {
			elementID, tokenID = tok.GetElementId(), tok.GetId()
			break
		}
	}
	if elementID != "UserTask_1" {
		t.Fatalf("waiting=%s", elementID)
	}
	if _, err := client.ThrowError(ctx, &enginev1.ThrowErrorRequest{
		ProcessInstanceId: created.GetProcessInstanceId(),
		ElementId:         elementID,
		TokenId:           tokenID,
		ErrorCode:         "BUSINESS_ERROR",
	}); err != nil {
		t.Fatal(err)
	}
	got, err = client.GetInstance(ctx, &enginev1.GetInstanceRequest{ProcessInstanceId: created.GetProcessInstanceId()})
	if err != nil {
		t.Fatal(err)
	}
	if got.GetInstance().GetStatus() != string(projection.StatusCompleted) {
		t.Fatalf("status=%s want completed", got.GetInstance().GetStatus())
	}
}

func TestEngineServiceProcessVersionCoexistence(t *testing.T) {
	_, conn, stop := startGRPC(t)
	defer stop()
	ctx := context.Background()
	client := enginev1.NewEngineServiceClient(conn)

	v1, err := os.ReadFile(filepath.Join("..", "processing", "testdata", "m5_version_v1.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	v2, err := os.ReadFile(filepath.Join("..", "processing", "testdata", "m5_version_v2.bpmn"))
	if err != nil {
		t.Fatal(err)
	}

	dep1, err := client.Deploy(ctx, &enginev1.DeployRequest{BpmnXml: v1})
	if err != nil {
		t.Fatal(err)
	}
	if dep1.GetProcessId() != "Process_version" || dep1.GetProcessVersion() != 1 {
		t.Fatalf("dep1=%v", dep1)
	}
	createdA, err := client.CreateInstance(ctx, &enginev1.CreateInstanceRequest{DeploymentId: dep1.GetDeploymentId()})
	if err != nil {
		t.Fatal(err)
	}

	dep2, err := client.Deploy(ctx, &enginev1.DeployRequest{BpmnXml: v2})
	if err != nil {
		t.Fatal(err)
	}
	if dep2.GetProcessVersion() != 2 {
		t.Fatalf("dep2 version=%d", dep2.GetProcessVersion())
	}

	createdB, err := client.CreateInstance(ctx, &enginev1.CreateInstanceRequest{ProcessId: "Process_version"})
	if err != nil {
		t.Fatal(err)
	}
	gotB, err := client.GetInstance(ctx, &enginev1.GetInstanceRequest{ProcessInstanceId: createdB.GetProcessInstanceId()})
	if err != nil {
		t.Fatal(err)
	}
	if gotB.GetInstance().GetProcessVersion() != 2 {
		t.Fatalf("B version=%d", gotB.GetInstance().GetProcessVersion())
	}
	waitingB := ""
	for _, tok := range gotB.GetInstance().GetTokens() {
		if tok.GetStatus() == string(projection.TokenWaiting) {
			waitingB = tok.GetElementId()
		}
	}
	if waitingB != "TaskB" {
		t.Fatalf("B waiting=%s", waitingB)
	}

	gotA, err := client.GetInstance(ctx, &enginev1.GetInstanceRequest{ProcessInstanceId: createdA.GetProcessInstanceId()})
	if err != nil {
		t.Fatal(err)
	}
	if gotA.GetInstance().GetProcessVersion() != 1 {
		t.Fatalf("A version=%d", gotA.GetInstance().GetProcessVersion())
	}
	waitingA := ""
	for _, tok := range gotA.GetInstance().GetTokens() {
		if tok.GetStatus() == string(projection.TokenWaiting) {
			waitingA = tok.GetElementId()
		}
	}
	if waitingA != "TaskA" {
		t.Fatalf("A waiting=%s", waitingA)
	}

	createdC, err := client.CreateInstance(ctx, &enginev1.CreateInstanceRequest{
		ProcessId:      "Process_version",
		ProcessVersion: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	gotC, err := client.GetInstance(ctx, &enginev1.GetInstanceRequest{ProcessInstanceId: createdC.GetProcessInstanceId()})
	if err != nil {
		t.Fatal(err)
	}
	waitingC := ""
	for _, tok := range gotC.GetInstance().GetTokens() {
		if tok.GetStatus() == string(projection.TokenWaiting) {
			waitingC = tok.GetElementId()
		}
	}
	if waitingC != "TaskA" || gotC.GetInstance().GetProcessVersion() != 1 {
		t.Fatalf("C waiting=%s version=%d", waitingC, gotC.GetInstance().GetProcessVersion())
	}

	_, err = client.CreateInstance(ctx, &enginev1.CreateInstanceRequest{
		ProcessId:      "Process_version",
		ProcessVersion: 99,
	})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("unknown version code=%v err=%v", status.Code(err), err)
	}
}

func TestEngineServiceLoopInstanceIndex(t *testing.T) {
	_, conn, stop := startGRPC(t)
	defer stop()
	ctx := context.Background()
	client := enginev1.NewEngineServiceClient(conn)

	xml, err := os.ReadFile(filepath.Join("..", "processing", "testdata", "m6_mi_parallel_cardinality.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	dep, err := client.Deploy(ctx, &enginev1.DeployRequest{BpmnXml: xml})
	if err != nil {
		t.Fatal(err)
	}
	created, err := client.CreateInstance(ctx, &enginev1.CreateInstanceRequest{DeploymentId: dep.GetDeploymentId()})
	if err != nil {
		t.Fatal(err)
	}
	got, err := client.GetInstance(ctx, &enginev1.GetInstanceRequest{ProcessInstanceId: created.GetProcessInstanceId()})
	if err != nil {
		t.Fatal(err)
	}
	seen := map[int32]bool{}
	for _, tok := range got.GetInstance().GetTokens() {
		if tok.GetElementId() != "UserTask_mi" || tok.GetStatus() != string(projection.TokenWaiting) {
			continue
		}
		idx := tok.GetLoopInstanceIndex()
		if idx < 0 {
			continue
		}
		seen[idx] = true
	}
	if len(seen) != 3 {
		t.Fatalf("loop indexes=%v want 0,1,2", seen)
	}
}

func TestEngineServiceResolveIncident(t *testing.T) {
	eng, conn, stop := startGRPC(t)
	defer stop()
	ctx := context.Background()
	engineClient := enginev1.NewEngineServiceClient(conn)
	jobClient := jobv1.NewJobServiceClient(conn)

	xml, err := os.ReadFile(filepath.Join("..", "processing", "testdata", "m7_incident_service.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	dep, err := engineClient.Deploy(ctx, &enginev1.DeployRequest{BpmnXml: xml})
	if err != nil {
		t.Fatal(err)
	}
	created, err := engineClient.CreateInstance(ctx, &enginev1.CreateInstanceRequest{DeploymentId: dep.GetDeploymentId()})
	if err != nil {
		t.Fatal(err)
	}
	instanceID := created.GetProcessInstanceId()

	for i := 0; i < 3; i++ {
		act, err := jobClient.ActivateJobs(ctx, &jobv1.ActivateJobsRequest{JobType: "work.v1", WorkerId: "gw"})
		if err != nil || len(act.GetJobs()) != 1 {
			t.Fatalf("activate i=%d jobs=%v err=%v", i, act.GetJobs(), err)
		}
		job := act.GetJobs()[0]
		if _, err := jobClient.FailJob(ctx, &jobv1.FailJobRequest{
			ProcessInstanceId: job.GetProcessInstanceId(),
			ElementId:         job.GetElementId(),
			TokenId:           job.GetTokenId(),
			ErrorMessage:      "err",
		}); err != nil {
			t.Fatal(err)
		}
	}

	got, err := engineClient.GetInstance(ctx, &enginev1.GetInstanceRequest{ProcessInstanceId: instanceID})
	if err != nil {
		t.Fatal(err)
	}
	var elementID, tokenID string
	for _, tok := range got.GetInstance().GetTokens() {
		if tok.GetStatus() == string(projection.TokenBlocked) {
			elementID, tokenID = tok.GetElementId(), tok.GetId()
			if tok.GetIncidentErrorMessage() != "err" {
				t.Fatalf("token=%v", tok)
			}
			break
		}
	}
	if elementID == "" {
		t.Fatalf("tokens=%v", got.GetInstance().GetTokens())
	}

	if _, err := engineClient.ResolveIncident(ctx, &enginev1.ResolveIncidentRequest{
		ProcessInstanceId: instanceID,
		ElementId:         elementID,
		TokenId:           tokenID,
	}); err != nil {
		t.Fatal(err)
	}

	act, err := jobClient.ActivateJobs(ctx, &jobv1.ActivateJobsRequest{JobType: "work.v1", WorkerId: "gw2"})
	if err != nil || len(act.GetJobs()) != 1 {
		t.Fatalf("activate after resolve=%v err=%v", act.GetJobs(), err)
	}
	job := act.GetJobs()[0]
	if _, err := jobClient.CompleteJob(ctx, &jobv1.CompleteJobRequest{
		ProcessInstanceId: job.GetProcessInstanceId(),
		ElementId:         job.GetElementId(),
		TokenId:           job.GetTokenId(),
	}); err != nil {
		t.Fatal(err)
	}
	got, err = engineClient.GetInstance(ctx, &enginev1.GetInstanceRequest{ProcessInstanceId: instanceID})
	if err != nil {
		t.Fatal(err)
	}
	if got.GetInstance().GetStatus() != string(projection.StatusCompleted) {
		t.Fatalf("status=%s", got.GetInstance().GetStatus())
	}
	_ = eng
}

func TestEngineServiceInterventionStepContinue(t *testing.T) {
	_, conn, stop := startGRPC(t)
	defer stop()
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
	created, err := client.CreateInstance(ctx, &enginev1.CreateInstanceRequest{DeploymentId: dep.GetDeploymentId()})
	if err != nil {
		t.Fatal(err)
	}
	id := created.GetProcessInstanceId()
	if _, err := client.EnableIntervention(ctx, &enginev1.EnableInterventionRequest{
		ProcessInstanceId: id,
		Policy:            "step",
	}); err != nil {
		t.Fatal(err)
	}
	got, err := client.GetInstance(ctx, &enginev1.GetInstanceRequest{ProcessInstanceId: id})
	if err != nil {
		t.Fatal(err)
	}
	var elementID, tokenID string
	for _, tok := range got.GetInstance().GetTokens() {
		if tok.GetStatus() == string(projection.TokenWaiting) {
			elementID, tokenID = tok.GetElementId(), tok.GetId()
			break
		}
	}
	if elementID == "" {
		t.Fatal("no waiting token")
	}
	if _, err := client.Complete(ctx, &enginev1.CompleteRequest{
		ProcessInstanceId: id,
		ElementId:         elementID,
		TokenId:           tokenID,
		Variables:         map[string]string{"approved": "true"},
	}); err != nil {
		t.Fatal(err)
	}
	st, err := client.GetInterventionState(ctx, &enginev1.GetInterventionStateRequest{ProcessInstanceId: id})
	if err != nil {
		t.Fatal(err)
	}
	if !st.GetState().GetPaused() {
		t.Fatalf("want paused after step Complete, got %#v", st.GetState())
	}
	for i := 0; i < 8 && st.GetState().GetPauseElementId() != "Gateway_1"; i++ {
		if _, err := client.StepInto(ctx, &enginev1.StepIntoRequest{ProcessInstanceId: id}); err != nil {
			t.Fatal(err)
		}
		st, err = client.GetInterventionState(ctx, &enginev1.GetInterventionStateRequest{ProcessInstanceId: id})
		if err != nil {
			t.Fatal(err)
		}
	}
	if st.GetState().GetPauseElementId() != "Gateway_1" {
		t.Fatalf("want Gateway_1, got %#v", st.GetState())
	}
	if _, err := client.SetVariables(ctx, &enginev1.SetVariablesRequest{
		ProcessInstanceId: id,
		Variables:         map[string]string{"approved": "false"},
	}); err != nil {
		t.Fatal(err)
	}
	cont, err := client.Continue(ctx, &enginev1.ContinueRequest{ProcessInstanceId: id})
	if err != nil {
		t.Fatal(err)
	}
	if cont.GetPaused() {
		t.Fatalf("want unpaused after Continue, %#v", cont.GetState())
	}
	got, err = client.GetInstance(ctx, &enginev1.GetInstanceRequest{ProcessInstanceId: id})
	if err != nil {
		t.Fatal(err)
	}
	if got.GetInstance().GetStatus() != string(projection.StatusCompleted) {
		t.Fatalf("status=%s", got.GetInstance().GetStatus())
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
