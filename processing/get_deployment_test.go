package processing_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
)

func TestGetDeploymentReturnsSourceXML(t *testing.T) {
	xml := readTestdata(t, "m1_simple.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	id, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	dep, ok := eng.GetDeployment(id)
	if !ok || dep == nil {
		t.Fatal("missing deployment")
	}
	if dep.ProcessID() == "" || dep.Version < 1 {
		t.Fatalf("process=%s version=%d", dep.ProcessID(), dep.Version)
	}
	got, err := eng.GetDeploymentXML(id)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, xml) {
		t.Fatalf("xml mismatch len got=%d want=%d", len(got), len(xml))
	}
}

func TestGetDeploymentAfterOpen(t *testing.T) {
	xml := readTestdata(t, "m1_simple.bpmn")
	dir := t.TempDir()
	ctx := context.Background()
	eng1, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	id, err := eng1.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	if err := eng1.Close(); err != nil {
		t.Fatal(err)
	}

	eng2, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	defer eng2.Close()

	got, err := eng2.GetDeploymentXML(id)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, xml) {
		t.Fatalf("xml mismatch after open")
	}
	// Persisted file still present for the store path.
	if _, err := os.Stat(filepath.Join(dir, "deployments", id+".bpmn")); err != nil {
		t.Fatalf("store file: %v", err)
	}
}

func TestGetDeploymentNotFound(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	_, err := eng.GetDeploymentXML("no-such-deployment")
	if err == nil || !strings.Contains(err.Error(), "NOT_FOUND") {
		t.Fatalf("err=%v", err)
	}
}
