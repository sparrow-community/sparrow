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

package processing

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/sparrow-community/sparrow/processing/deploy"
	"github.com/sparrow-community/sparrow/processing/expr"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// EvaluateConditionsRequest evaluates waiting conditional catches and boundaries.
type EvaluateConditionsRequest struct {
	ProcessInstanceID string
	Variables         map[string]any // overlay for evaluation; applied on Complete when fired
}

// EvaluateConditions completes waiting intermediate conditional catches and fires
// conditional boundaries whose condition is true.
func (e *Engine) EvaluateConditions(ctx context.Context, req EvaluateConditionsRequest) (int, error) {
	overlay, err := variablesJSONMap(req.Variables)
	if err != nil {
		return 0, err
	}
	instanceID := strings.TrimSpace(req.ProcessInstanceID)
	if instanceID != "" {
		e.mu.Lock()
		_, ok := e.instances[instanceID]
		e.mu.Unlock()
		if !ok {
			return 0, fmt.Errorf("NOT_FOUND: instance %q", instanceID)
		}
	}

	type hit struct {
		instanceID      string
		elementID       string
		tokenID         string
		scope           bool
		eventSubProcess bool
	}
	var hits []hit

	ids := e.ListInstanceIDs()
	if instanceID != "" {
		ids = []string{instanceID}
	}
	for _, iid := range ids {
		e.mu.Lock()
		inst := e.instances[iid]
		lock := e.instMu[iid]
		var dep *deploy.Deployment
		if inst != nil {
			dep = e.deployments[inst.DeploymentID]
		}
		e.mu.Unlock()
		if inst == nil || lock == nil || dep == nil {
			continue
		}
		lock.Lock()
		env := mergeVarMaps(inst.Variables, overlay)
		for _, tok := range inst.Tokens {
			if tok == nil || tok.Status != projection.TokenWaiting {
				continue
			}
			typ, err := dep.TypeOf(tok.ElementID)
			if err == nil && typ == eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT {
				kind, err := dep.CatchKind(tok.ElementID)
				if err == nil && kind == deploy.CatchKindConditional {
					text, ok := dep.ConditionalCatchExpression(tok.ElementID)
					if ok && evalTrue(text, env) {
						hits = append(hits, hit{instanceID: iid, elementID: tok.ElementID, tokenID: tok.ID})
					}
				}
			}
			for _, w := range tok.BoundaryWaits {
				if w.Kind != "conditional" {
					continue
				}
				text, ok := dep.ConditionalCatchExpression(w.BoundaryID)
				if !ok {
					text = w.TimerText
				}
				if text == "" {
					continue
				}
				if evalTrue(text, env) {
					hits = append(hits, hit{instanceID: iid, elementID: w.BoundaryID, tokenID: tok.ID})
				}
			}
		}
		for _, sb := range inst.ScopeBoundaries {
			if sb == nil {
				continue
			}
			text, ok := dep.ConditionalCatchExpression(sb.BoundaryID)
			if !ok {
				continue
			}
			if evalTrue(text, env) {
				hits = append(hits, hit{instanceID: iid, elementID: sb.BoundaryID, tokenID: sb.TokenID, scope: true})
			}
		}
		for eventSubProcessID, arm := range inst.EventSubProcesses {
			if arm == nil {
				continue
			}
			spec, ok := dep.EventSubProcessSpec(eventSubProcessID)
			if !ok || spec.Kind != deploy.CatchKindConditional {
				continue
			}
			text := spec.Condition
			if text == "" {
				text, _ = dep.ConditionalCatchExpression(spec.StartEventID)
			}
			if text == "" {
				text = arm.TimerText
			}
			if text == "" || !evalTrue(text, env) {
				continue
			}
			hits = append(hits, hit{instanceID: iid, elementID: eventSubProcessID, tokenID: "", eventSubProcess: true})
		}
		lock.Unlock()
	}

	sort.Slice(hits, func(i, j int) bool {
		if hits[i].instanceID != hits[j].instanceID {
			return hits[i].instanceID < hits[j].instanceID
		}
		return hits[i].elementID < hits[j].elementID
	})

	fired := 0
	var first error
	seenTok := map[string]bool{}
	for _, h := range hits {
		if err := ctx.Err(); err != nil {
			return fired, err
		}
		key := h.instanceID + "/" + h.tokenID + "/" + h.elementID
		if seenTok[key] {
			continue
		}
		seenTok[key] = true
		var err error
		switch {
		case h.eventSubProcess:
			err = e.triggerEventSubProcess(ctx, h.instanceID, h.elementID, req.Variables)
		case h.scope:
			err = e.completeScopeBoundary(ctx, h.instanceID, h.elementID)
		default:
			err = e.Complete(ctx, h.instanceID, h.elementID, h.tokenID, req.Variables)
		}
		if err != nil {
			if strings.HasPrefix(err.Error(), "INVALID_STATE:") {
				continue
			}
			if first == nil {
				first = err
			}
			continue
		}
		fired++
	}
	return fired, first
}

func evalTrue(text string, env map[string]string) bool {
	ok, err := expr.Eval(text, env)
	return err == nil && ok
}

func mergeVarMaps(base, overlay map[string]string) map[string]string {
	out := make(map[string]string, len(base)+len(overlay))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range overlay {
		out[k] = v
	}
	return out
}
