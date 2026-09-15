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
	"sort"

	"github.com/sparrow-community/sparrow/processing/deploy"
)

func (e *Engine) snapshotDeployments() []*deploy.Deployment {
	e.mu.Lock()
	defer e.mu.Unlock()
	ids := make([]string, 0, len(e.deployments))
	for id := range e.deployments {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	out := make([]*deploy.Deployment, 0, len(ids))
	for _, id := range ids {
		out = append(out, e.deployments[id])
	}
	return out
}

func (e *Engine) createMessageStartInstances(ctx context.Context, name string, vars map[string]any) (int, error) {
	created := 0
	var first error
	for _, dep := range e.snapshotDeployments() {
		starts := dep.MessageStartIDs(name)
		sort.Strings(starts)
		for _, startID := range starts {
			if err := ctx.Err(); err != nil {
				return created, err
			}
			if _, err := e.createInstanceAt(ctx, dep, vars, startID); err != nil {
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

func (e *Engine) createSignalStartInstances(ctx context.Context, name string, vars map[string]any) (int, error) {
	created := 0
	var first error
	for _, dep := range e.snapshotDeployments() {
		starts := dep.SignalStartIDs(name)
		sort.Strings(starts)
		for _, startID := range starts {
			if err := ctx.Err(); err != nil {
				return created, err
			}
			if _, err := e.createInstanceAt(ctx, dep, vars, startID); err != nil {
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
