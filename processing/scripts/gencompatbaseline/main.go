// Command gencompatbaseline writes an Open()-shaped data directory and
// expected.json for the compatibility golden test.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sparrow-community/sparrow/processing"
	"github.com/sparrow-community/sparrow/processing/projection"
)

func main() {
	out := flag.String("out", "", "output directory (Open data dir)")
	flag.Parse()
	if *out == "" {
		fmt.Fprintln(os.Stderr, "-out is required")
		os.Exit(2)
	}
	if err := generate(*out); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type expected struct {
	InstanceCount int    `json:"instance_count"`
	Status        string `json:"status"`
	WaitingElement string `json:"waiting_element"`
	ProcessID     string `json:"process_id"`
	VariableKeys  []string `json:"variable_keys"`
}

func generate(out string) error {
	ctx := context.Background()
	if err := os.MkdirAll(out, 0o755); err != nil {
		return err
	}
	eng, err := processing.Open(ctx, out)
	if err != nil {
		return err
	}
	defer eng.Close()

	bpmn, err := os.ReadFile(filepath.Join("testdata", "m1_simple.bpmn"))
	if err != nil {
		return err
	}
	dep, err := eng.Deploy(ctx, bpmn)
	if err != nil {
		return fmt.Errorf("deploy: %w", err)
	}
	iid, err := eng.CreateInstance(ctx, dep, map[string]any{"approved": true})
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	inst, ok := eng.GetInstance(iid)
	if !ok || inst == nil {
		return fmt.Errorf("missing instance %s", iid)
	}
	waiting := ""
	for _, tok := range inst.Tokens {
		if tok != nil && tok.Status == projection.TokenWaiting {
			waiting = tok.ElementID
			break
		}
	}
	if waiting != "UserTask_1" {
		return fmt.Errorf("expected wait at UserTask_1, got %q status=%s", waiting, inst.Status)
	}
	keys := make([]string, 0, len(inst.Variables))
	for k := range inst.Variables {
		keys = append(keys, k)
	}
	exp := expected{
		InstanceCount:  1,
		Status:         string(inst.Status),
		WaitingElement: waiting,
		ProcessID:      inst.ProcessID,
		VariableKeys:   keys,
	}
	b, err := json.MarshalIndent(exp, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(out, "expected.json"), append(b, '\n'), 0o644); err != nil {
		return err
	}
	return nil
}
