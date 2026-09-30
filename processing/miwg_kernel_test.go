// Copyright 2025 The Sparrow community and contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package processing_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/sparrow-community/sparrow/processing"
	"github.com/sparrow-community/sparrow/processing/deploy"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
)

// TestMIWGKernelFixtures deploys every bpmn/test MIWG fixture and, when the
// root process is Supported, mints an instance and drives waits until the
// instance completes or no further automatic progress is possible.
//
// Deploy failures that contain UNSUPPORTED_ELEMENT or INVALID_CONDITION are
// intentional skips (Excluded / current kernel limits). Other deploy errors fail.
func TestMIWGKernelFixtures(t *testing.T) {
	dir := filepath.Join("..", "bpmn", "test")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".bpmn") {
			continue
		}
		names = append(names, e.Name())
	}
	sort.Strings(names)
	if len(names) == 0 {
		t.Fatal("no MIWG fixtures")
	}

	for _, name := range names {
		name := name
		t.Run(name, func(t *testing.T) {
			xml, err := os.ReadFile(filepath.Join(dir, name))
			if err != nil {
				t.Fatal(err)
			}
			compiled, err := deploy.Compile(xml)
			if err != nil {
				msg := err.Error()
				if isIntentionalDeploySkip(msg) {
					t.Skipf("intentional skip: %s", firstLine(msg))
				}
				t.Fatalf("Deploy/Compile: %v", err)
			}

			eng := processing.NewEngine(eventlog.NewMemory())
			ctx := context.Background()
			now := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
			eng.SetNow(func() time.Time { return now })

			depID, err := eng.Deploy(ctx, xml)
			if err != nil {
				msg := err.Error()
				if isIntentionalDeploySkip(msg) {
					t.Skipf("intentional skip: %s", firstLine(msg))
				}
				t.Fatalf("Deploy: %v", err)
			}
			dep, ok := eng.GetDeployment(depID)
			if !ok || dep == nil {
				t.Fatal("deployment missing after Deploy")
			}

			instanceID, err := mintMIWGInstance(ctx, eng, dep, depID)
			if err != nil {
				if isIntentionalRunSkip(err.Error()) {
					t.Skipf("intentional skip: %s", firstLine(err.Error()))
				}
				t.Fatalf("mint instance: %v", err)
			}

			outcome := driveMIWGInstance(t, ctx, eng, &now, instanceID)
			t.Logf("process=%s outcome=%s", compiled.Process.ID, outcome)
		})
	}
}

func isIntentionalDeploySkip(msg string) bool {
	return strings.Contains(msg, "UNSUPPORTED_ELEMENT") ||
		strings.Contains(msg, "INVALID_CONDITION")
}

func isIntentionalRunSkip(msg string) bool {
	return strings.Contains(msg, "UNSUPPORTED_ELEMENT") ||
		strings.Contains(msg, "INVALID_CONDITION") ||
		strings.Contains(msg, "INVALID_ARGUMENT") ||
		strings.Contains(msg, "no CreateInstance entry") ||
		strings.Contains(msg, "typed start")
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	if len(s) > 240 {
		return s[:240] + "..."
	}
	return s
}

func mintMIWGInstance(ctx context.Context, eng *processing.Engine, dep *deploy.Deployment, depID string) (string, error) {
	if _, err := dep.CreateInstanceEntryIDs(); err == nil {
		return eng.CreateInstance(ctx, depID, nil)
	}

	for _, name := range dep.MessageStartNames() {
		n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: name})
		if err != nil {
			return "", err
		}
		if n > 0 {
			ids := eng.ListInstanceIDs()
			if len(ids) == 0 {
				return "", fmt.Errorf("INVALID_ARGUMENT: message start %q minted 0 instances", name)
			}
			sort.Strings(ids)
			return ids[len(ids)-1], nil
		}
	}
	for _, name := range dep.SignalStartNames() {
		n, err := eng.PublishSignal(ctx, processing.PublishSignalRequest{Name: name})
		if err != nil {
			return "", err
		}
		if n > 0 {
			ids := eng.ListInstanceIDs()
			if len(ids) == 0 {
				return "", fmt.Errorf("INVALID_ARGUMENT: signal start %q minted 0 instances", name)
			}
			sort.Strings(ids)
			return ids[len(ids)-1], nil
		}
	}
	if len(dep.TimerStartIDs()) > 0 {
		if err := eng.FireDue(ctx); err != nil {
			return "", err
		}
		ids := eng.ListInstanceIDs()
		if len(ids) > 0 {
			sort.Strings(ids)
			return ids[len(ids)-1], nil
		}
	}
	if len(dep.ConditionalStartIDs()) > 0 {
		return "", fmt.Errorf("INVALID_ARGUMENT: typed start (conditional) — MIWG runner does not auto-EvaluateConditionalStarts")
	}
	return "", fmt.Errorf("INVALID_ARGUMENT: no CreateInstance entry and no mintable typed start")
}

func driveMIWGInstance(t *testing.T, ctx context.Context, eng *processing.Engine, now *time.Time, instanceID string) string {
	t.Helper()
	const maxSteps = 250
	prevFingerprint := ""
	for step := 0; step < maxSteps; step++ {
		inst, ok := eng.GetInstance(instanceID)
		if !ok {
			t.Fatalf("instance %s missing", instanceID)
		}
		switch inst.Status {
		case projection.StatusCompleted, projection.StatusTerminated:
			return string(inst.Status)
		}
		fp := miwgWaitFingerprint(inst)
		if fp != "" && fp == prevFingerprint {
			return "active_no_progress"
		}
		prevFingerprint = fp

		progressed := false

		// Prefer job-backed waits: Activate then Complete.
		jobTypes := map[string]struct{}{}
		for _, tok := range inst.Tokens {
			if tok.Status == projection.TokenWaiting && tok.JobType != "" && !tok.ScopeHost {
				jobTypes[tok.JobType] = struct{}{}
			}
		}
		for jt := range jobTypes {
			jobs, err := eng.Activate(ctx, processing.ActivateRequest{
				JobType:  jt,
				WorkerID: "miwg-runner",
				MaxJobs:  32,
			})
			if err != nil {
				t.Fatalf("Activate %s: %v", jt, err)
			}
			for _, job := range jobs {
				if job.ProcessInstanceID != instanceID {
					continue
				}
				if err := eng.Complete(ctx, job.ProcessInstanceID, job.ElementID, job.TokenID, nil); err != nil {
					t.Fatalf("Complete job %s/%s: %v", job.ElementID, job.TokenID, err)
				}
				progressed = true
			}
		}
		if progressed {
			continue
		}

		inst, _ = eng.GetInstance(instanceID)
		// Timer waits: advance clock past due and FireDue.
		var due int64
		for _, tok := range inst.Tokens {
			if tok.Status != projection.TokenWaiting {
				continue
			}
			if tok.DueUnixMs > 0 && (due == 0 || tok.DueUnixMs < due) {
				due = tok.DueUnixMs
			}
			for _, bw := range tok.BoundaryWaits {
				if bw.DueUnixMs > 0 && (due == 0 || bw.DueUnixMs < due) {
					due = bw.DueUnixMs
				}
			}
		}
		if due > 0 {
			*now = time.UnixMilli(due).Add(time.Second).UTC()
			if err := eng.FireDue(ctx); err != nil {
				t.Fatalf("FireDue: %v", err)
			}
			progressed = true
			continue
		}

		inst, _ = eng.GetInstance(instanceID)
		// Message / signal catches.
		for _, tok := range inst.Tokens {
			if tok.Status != projection.TokenWaiting {
				continue
			}
			if tok.MessageName != "" {
				if _, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{
					Name:              tok.MessageName,
					ProcessInstanceID: instanceID,
				}); err != nil {
					t.Fatalf("PublishMessage %s: %v", tok.MessageName, err)
				}
				progressed = true
				break
			}
			if tok.SignalName != "" {
				if _, err := eng.PublishSignal(ctx, processing.PublishSignalRequest{
					Name:              tok.SignalName,
					ProcessInstanceID: instanceID,
				}); err != nil {
					t.Fatalf("PublishSignal %s: %v", tok.SignalName, err)
				}
				progressed = true
				break
			}
			for _, bw := range tok.BoundaryWaits {
				if bw.MessageName != "" {
					if _, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{
						Name:              bw.MessageName,
						ProcessInstanceID: instanceID,
					}); err != nil {
						t.Fatalf("PublishMessage boundary %s: %v", bw.MessageName, err)
					}
					progressed = true
					break
				}
				if bw.SignalName != "" {
					if _, err := eng.PublishSignal(ctx, processing.PublishSignalRequest{
						Name:              bw.SignalName,
						ProcessInstanceID: instanceID,
					}); err != nil {
						t.Fatalf("PublishSignal boundary %s: %v", bw.SignalName, err)
					}
					progressed = true
					break
				}
			}
			if progressed {
				break
			}
		}
		if progressed {
			continue
		}

		// Waiting activities (user / manual / abstract / receive without message name).
		inst, _ = eng.GetInstance(instanceID)
		el, tokID := waitingAt(inst)
		if el == "" || tokID == "" {
			return "active_no_wait"
		}
		if err := eng.Complete(ctx, instanceID, el, tokID, nil); err != nil {
			// Gateway / join waits may not accept Complete — treat as stuck run.
			t.Logf("Complete %s blocked: %v", el, err)
			return "active_stuck:" + el
		}
	}
	return "max_steps"
}

func miwgWaitFingerprint(inst *projection.Instance) string {
	parts := make([]string, 0, len(inst.Tokens))
	for _, tok := range inst.Tokens {
		if tok.Status != projection.TokenWaiting && tok.Status != projection.TokenBlocked {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s|%s|%s|%s|%d", tok.ID, tok.ElementID, tok.Status, tok.JobType, tok.DueUnixMs))
	}
	sort.Strings(parts)
	return strings.Join(parts, ";")
}
