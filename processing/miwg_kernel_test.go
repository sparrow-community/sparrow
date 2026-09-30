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
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// TestMIWGKernelFixtures deploys every bpmn/test MIWG fixture and mints an
// instance, driving waits until completed/terminated.
//
// Compile / Deploy / mint failures fail the test (no intentional skips).
// Decorative / modeling stubs must Deploy and run (opaque / defaults).
// A successful Deploy that does not reach completed/terminated fails the test.
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
				t.Fatalf("Deploy/Compile: %v", err)
			}

			eng := processing.NewEngine(eventlog.NewMemory())
			ctx := context.Background()
			now := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
			eng.SetNow(func() time.Time { return now })

			depID, err := eng.Deploy(ctx, xml)
			if err != nil {
				t.Fatalf("Deploy: %v", err)
			}
			dep, ok := eng.GetDeployment(depID)
			if !ok || dep == nil {
				t.Fatal("deployment missing after Deploy")
			}

			instanceID, err := mintMIWGInstance(ctx, eng, dep, depID)
			if err != nil {
				t.Fatalf("mint instance: %v", err)
			}

			outcome := driveMIWGInstance(t, ctx, eng, dep, &now, instanceID)
			t.Logf("process=%s outcome=%s", compiled.Process.ID, outcome)
			switch outcome {
			case string(projection.StatusCompleted), string(projection.StatusTerminated):
				// ok
			default:
				t.Fatalf("expected completed/terminated, got %s", outcome)
			}
		})
	}
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

// miwgHappyVars seeds common decision variables so exclusive gateways with
// author conditions can leave on the progressing branch. Unlabeled exclusive
// splits (C.4 / C.7) rely on kernel non-back-edge preference instead of fixture edits.
func miwgHappyVars() map[string]any {
	return map[string]any{
		"approved":          true,
		"clarified":         "yes",
		"Vacation Approval": "Approved",
		// C.9.0-roundtrip Risk gateway: non-red/non-all-yellow → default Green path.
		"riskLevels": []string{"green"},
	}
}

func driveMIWGInstance(t *testing.T, ctx context.Context, eng *processing.Engine, dep *deploy.Deployment, now *time.Time, instanceID string) string {
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

		// Prefer job-backed waits: Activate then Complete with happy-path vars.
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
				if err := eng.Complete(ctx, job.ProcessInstanceID, job.ElementID, job.TokenID, miwgHappyVars()); err != nil {
					t.Fatalf("Complete job %s/%s: %v", job.ElementID, job.TokenID, err)
				}
				progressed = true
			}
		}
		if progressed {
			continue
		}

		inst, _ = eng.GetInstance(instanceID)
		// Message / signal catches — prefer over timers so NI timer cycles do not
		// starve a concurrent message receive (C.9.1).
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

		inst, _ = eng.GetInstance(instanceID)
		// Conditional catches / boundaries: evaluate with happy-path vars.
		if n, err := eng.EvaluateConditions(ctx, processing.EvaluateConditionsRequest{
			ProcessInstanceID: instanceID,
			Variables:         miwgHappyVars(),
		}); err != nil {
			t.Fatalf("EvaluateConditions: %v", err)
		} else if n > 0 {
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

		// Waiting activities (user / manual / abstract / opaque SP / receive).
		// Never Complete gateway join waits — they advance when peer tokens arrive.
		inst, _ = eng.GetInstance(instanceID)
		el, tokID := miwgWaitingActivity(dep, inst)
		if el == "" || tokID == "" {
			return "active_no_wait"
		}
		if err := eng.Complete(ctx, instanceID, el, tokID, miwgHappyVars()); err != nil {
			t.Logf("Complete %s blocked: %v", el, err)
			return "active_stuck:" + el
		}
	}
	return "max_steps"
}

func miwgWaitingActivity(dep *deploy.Deployment, inst *projection.Instance) (elementID, tokenID string) {
	var hostElem, hostTok string
	var candidates [][2]string
	for _, tok := range inst.Tokens {
		if tok.Status != projection.TokenWaiting {
			continue
		}
		if tok.ScopeHost {
			if hostTok == "" {
				hostElem, hostTok = tok.ElementID, tok.ID
			}
			continue
		}
		if miwgIsGateway(dep, tok.ElementID) {
			continue
		}
		candidates = append(candidates, [2]string{tok.ElementID, tok.ID})
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i][0] != candidates[j][0] {
			return candidates[i][0] < candidates[j][0]
		}
		return candidates[i][1] < candidates[j][1]
	})
	if len(candidates) > 0 {
		return candidates[0][0], candidates[0][1]
	}
	if hostElem != "" && !miwgIsGateway(dep, hostElem) {
		return hostElem, hostTok
	}
	return "", ""
}

func miwgIsGateway(dep *deploy.Deployment, elementID string) bool {
	if dep == nil {
		return false
	}
	typ, err := dep.TypeOf(elementID)
	if err != nil {
		return false
	}
	switch typ {
	case eventv1.Element_TYPE_EXCLUSIVE_GATEWAY,
		eventv1.Element_TYPE_PARALLEL_GATEWAY,
		eventv1.Element_TYPE_INCLUSIVE_GATEWAY,
		eventv1.Element_TYPE_COMPLEX_GATEWAY,
		eventv1.Element_TYPE_EVENT_BASED_GATEWAY:
		return true
	default:
		return false
	}
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
