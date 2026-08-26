package processing_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
)

func TestProcessVersionCoexistence(t *testing.T) {
	v1 := readTestdataVersion(t, "m5_version_v1.bpmn")
	v2 := readTestdataVersion(t, "m5_version_v2.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()

	dep1, err := eng.Deploy(ctx, v1)
	if err != nil {
		t.Fatal(err)
	}
	d1, ok := eng.GetDeployment(dep1)
	if !ok || d1.Version != 1 || d1.ProcessID() != "Process_version" {
		t.Fatalf("dep1 version=%d process=%s", d1.Version, d1.ProcessID())
	}

	instanceA, err := eng.CreateInstance(ctx, dep1, nil)
	if err != nil {
		t.Fatal(err)
	}
	instA := mustInstance(t, eng, instanceA)
	if waitingElement(instA) != "TaskA" {
		t.Fatalf("A waiting=%s", waitingElement(instA))
	}
	if instA.Version != 1 {
		t.Fatalf("A version=%d", instA.Version)
	}

	dep2, err := eng.Deploy(ctx, v2)
	if err != nil {
		t.Fatal(err)
	}
	d2, ok := eng.GetDeployment(dep2)
	if !ok || d2.Version != 2 {
		t.Fatalf("dep2 version=%d", d2.Version)
	}

	instanceB, err := eng.CreateInstanceRequest(ctx, processing.CreateInstanceRequest{
		ProcessID: "Process_version",
	})
	if err != nil {
		t.Fatal(err)
	}
	instB := mustInstance(t, eng, instanceB)
	if waitingElement(instB) != "TaskB" {
		t.Fatalf("B (latest) waiting=%s want TaskB", waitingElement(instB))
	}
	if instB.Version != 2 {
		t.Fatalf("B version=%d want 2", instB.Version)
	}

	instA = mustInstance(t, eng, instanceA)
	if waitingElement(instA) != "TaskA" {
		t.Fatalf("in-flight A must stay on TaskA, got %s", waitingElement(instA))
	}
	if instA.Version != 1 || instA.DeploymentID != dep1 {
		t.Fatalf("in-flight A binding changed: version=%d dep=%s", instA.Version, instA.DeploymentID)
	}

	instanceC, err := eng.CreateInstanceRequest(ctx, processing.CreateInstanceRequest{
		ProcessID:      "Process_version",
		ProcessVersion: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	instC := mustInstance(t, eng, instanceC)
	if waitingElement(instC) != "TaskA" {
		t.Fatalf("C (explicit v1) waiting=%s want TaskA", waitingElement(instC))
	}
	if instC.Version != 1 {
		t.Fatalf("C version=%d want 1", instC.Version)
	}

	_, err = eng.CreateInstanceRequest(ctx, processing.CreateInstanceRequest{
		ProcessID:      "Process_version",
		ProcessVersion: 99,
	})
	if err == nil || !strings.Contains(err.Error(), "NOT_FOUND") {
		t.Fatalf("unknown version err=%v", err)
	}

	// deployment_id path still works for latest snapshot
	instanceD, err := eng.CreateInstance(ctx, dep2, nil)
	if err != nil {
		t.Fatal(err)
	}
	instD := mustInstance(t, eng, instanceD)
	if waitingElement(instD) != "TaskB" || instD.Status != projection.StatusActive {
		t.Fatalf("D tokens=%#v status=%s", instD.Tokens, instD.Status)
	}
}

func readTestdataVersion(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}
