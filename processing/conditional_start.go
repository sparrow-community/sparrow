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
)

// EvaluateConditionalStartsRequest evaluates process-level conditional starts.
type EvaluateConditionalStartsRequest struct {
	DeploymentID   string
	ProcessID      string
	ProcessVersion int32
	Variables      map[string]any
}

// EvaluateConditionalStarts creates instances for matching conditional starts whose condition is true.
func (e *Engine) EvaluateConditionalStarts(ctx context.Context, req EvaluateConditionalStartsRequest) (int, error) {
	var deps []*deploy.Deployment
	if strings.TrimSpace(req.DeploymentID) != "" || strings.TrimSpace(req.ProcessID) != "" {
		e.mu.Lock()
		dep, err := e.resolveDeploymentLocked(req.DeploymentID, req.ProcessID, req.ProcessVersion)
		e.mu.Unlock()
		if err != nil {
			return 0, err
		}
		deps = []*deploy.Deployment{dep}
	} else {
		return 0, fmt.Errorf("INVALID_ARGUMENT: deployment_id or process_id is required")
	}

	vars, err := variablesJSONMap(req.Variables)
	if err != nil {
		return 0, err
	}
	created := 0
	var first error
	for _, dep := range deps {
		starts := dep.ConditionalStartIDs()
		sort.Strings(starts)
		for _, startID := range starts {
			if err := ctx.Err(); err != nil {
				return created, err
			}
			text, ok := dep.ConditionalStartExpression(startID)
			if !ok {
				continue
			}
			match, err := expr.Eval(text, vars)
			if err != nil {
				if first == nil {
					first = fmt.Errorf("INVALID_CONDITION: start %s: %w", startID, err)
				}
				continue
			}
			if !match {
				continue
			}
			if _, err := e.createInstanceAt(ctx, dep, req.Variables, startID); err != nil {
				if first == nil {
					first = err
				}
				continue
			}
			created++
		}
	}
	return created, first
}

func variablesJSONMap(vars map[string]any) (map[string]string, error) {
	pv, err := projection.VariablesFromMap(vars)
	if err != nil {
		return nil, fmt.Errorf("INVALID_ARGUMENT: variables: %w", err)
	}
	out := make(map[string]string, len(pv))
	for _, v := range pv {
		if v == nil {
			continue
		}
		out[v.GetName()] = v.GetJsonValue()
	}
	return out, nil
}
