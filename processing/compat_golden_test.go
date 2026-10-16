package processing_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	"github.com/sparrow-community/sparrow/processing/projection"
)

// Compat golden baseline: a committed Open() data directory that must Recover
// to a known projection. Refresh only via scripts/regen-compat-baseline.sh
// after a constitution-reviewed additive change (Principle VII).
func TestCompatGoldenBaseline(t *testing.T) {
	ctx := context.Background()
	src := filepath.Join("testdata", "compat", "baseline")
	expPath := filepath.Join(src, "expected.json")
	expBytes, err := os.ReadFile(expPath)
	if err != nil {
		t.Fatalf("read expected.json: %v (run ./processing/scripts/regen-compat-baseline.sh)", err)
	}
	var exp struct {
		InstanceCount  int      `json:"instance_count"`
		Status         string   `json:"status"`
		WaitingElement string   `json:"waiting_element"`
		ProcessID      string   `json:"process_id"`
		VariableKeys   []string `json:"variable_keys"`
	}
	if err := json.Unmarshal(expBytes, &exp); err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	copyCompatDir(t, src, dir)

	eng, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatalf("Open golden: %v", err)
	}
	defer eng.Close()

	ids := eng.ListInstanceIDs()
	if len(ids) != exp.InstanceCount {
		t.Fatalf("instances=%d want %d", len(ids), exp.InstanceCount)
	}
	inst, ok := eng.GetInstance(ids[0])
	if !ok || inst == nil {
		t.Fatal("missing instance after recover")
	}
	if string(inst.Status) != exp.Status {
		t.Fatalf("status=%s want %s", inst.Status, exp.Status)
	}
	if inst.ProcessID != exp.ProcessID {
		t.Fatalf("process_id=%s want %s", inst.ProcessID, exp.ProcessID)
	}
	waiting := ""
	for _, tok := range inst.Tokens {
		if tok != nil && tok.Status == projection.TokenWaiting {
			waiting = tok.ElementID
			break
		}
	}
	if waiting != exp.WaitingElement {
		t.Fatalf("waiting=%q want %q (projection drift — Principle VII)", waiting, exp.WaitingElement)
	}
	gotKeys := make([]string, 0, len(inst.Variables))
	for k := range inst.Variables {
		gotKeys = append(gotKeys, k)
	}
	sort.Strings(gotKeys)
	wantKeys := append([]string{}, exp.VariableKeys...)
	sort.Strings(wantKeys)
	if len(gotKeys) != len(wantKeys) {
		t.Fatalf("variable keys=%v want %v", gotKeys, wantKeys)
	}
	for i := range wantKeys {
		if gotKeys[i] != wantKeys[i] {
			t.Fatalf("variable keys=%v want %v", gotKeys, wantKeys)
		}
	}
}

func copyCompatDir(t *testing.T, src, dst string) {
	t.Helper()
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "expected.json" {
			return nil
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, b, info.Mode())
	})
	if err != nil {
		t.Fatal(err)
	}
}
