// Copyright 2026 The Sparrow community and contributors
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

import "github.com/sparrow-community/sparrow/processing/projection"

// NextDueUnixMs returns the earliest armed timer due time across instances,
// scope boundaries, event sub-process timer arms, and process-level timer starts.
// Returns 0 when nothing is scheduled. Past-due arms are included so a host
// scheduler can call FireDue immediately.
func (e *Engine) NextDueUnixMs() int64 {
	if e == nil {
		return 0
	}
	var next int64

	consider := func(due int64) {
		if due <= 0 {
			return
		}
		if next == 0 || due < next {
			next = due
		}
	}

	e.mu.Lock()
	ids := make([]string, 0, len(e.instances))
	for id := range e.instances {
		ids = append(ids, id)
	}
	e.mu.Unlock()

	for _, iid := range ids {
		e.mu.Lock()
		inst := e.instances[iid]
		lock := e.instMu[iid]
		e.mu.Unlock()
		if inst == nil || lock == nil {
			continue
		}
		lock.Lock()
		for _, tok := range inst.Tokens {
			if tok == nil || (tok.Status != projection.TokenWaiting && tok.Status != projection.TokenBlocked) {
				continue
			}
			if len(tok.BoundaryWaits) > 0 {
				for _, w := range tok.BoundaryWaits {
					if w.Kind == "timer" {
						consider(w.DueUnixMs)
					}
				}
			} else {
				consider(tok.DueUnixMs)
			}
		}
		for _, sb := range inst.ScopeBoundaries {
			if sb != nil {
				consider(sb.DueUnixMs)
			}
		}
		for _, arm := range inst.EventSubProcesses {
			if arm != nil {
				consider(arm.DueUnixMs)
			}
		}
		lock.Unlock()
	}

	e.timerStartMu.Lock()
	for _, arm := range e.timerStartArms {
		consider(arm.dueUnixMs)
	}
	e.timerStartMu.Unlock()

	return next
}
