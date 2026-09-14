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
	"sort"

	"github.com/sparrow-community/sparrow/processing/deploy"
)

type timerStartArm struct {
	deploymentID string
	startEventID string
	dueUnixMs    int64
}

func timerStartArmKey(deploymentID, startEventID string) string {
	return deploymentID + "/" + startEventID
}

func (e *Engine) armProcessTimerStarts(dep *deploy.Deployment) {
	if e == nil || dep == nil {
		return
	}
	now := e.now()
	ids := dep.TimerStartIDs()
	sort.Strings(ids)
	e.timerStartMu.Lock()
	defer e.timerStartMu.Unlock()
	if e.timerStartArms == nil {
		e.timerStartArms = make(map[string]timerStartArm)
	}
	for _, startID := range ids {
		due, _, err := dep.TimerDue(startID, now)
		if err != nil {
			continue
		}
		key := timerStartArmKey(dep.ID, startID)
		e.timerStartArms[key] = timerStartArm{
			deploymentID: dep.ID,
			startEventID: startID,
			dueUnixMs:    due,
		}
	}
}

func (e *Engine) rearmAllProcessTimerStarts() {
	for _, dep := range e.snapshotDeployments() {
		e.armProcessTimerStarts(dep)
	}
}

func (e *Engine) collectDueTimerStarts(nowUnixMs int64) []timerStartArm {
	e.timerStartMu.Lock()
	defer e.timerStartMu.Unlock()
	out := make([]timerStartArm, 0)
	for _, arm := range e.timerStartArms {
		if arm.dueUnixMs > 0 && arm.dueUnixMs <= nowUnixMs {
			out = append(out, arm)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].dueUnixMs != out[j].dueUnixMs {
			return out[i].dueUnixMs < out[j].dueUnixMs
		}
		if out[i].deploymentID != out[j].deploymentID {
			return out[i].deploymentID < out[j].deploymentID
		}
		return out[i].startEventID < out[j].startEventID
	})
	return out
}

func (e *Engine) consumeOrRearmTimerStart(arm timerStartArm) {
	e.mu.Lock()
	dep := e.deployments[arm.deploymentID]
	e.mu.Unlock()
	key := timerStartArmKey(arm.deploymentID, arm.startEventID)
	e.timerStartMu.Lock()
	defer e.timerStartMu.Unlock()
	delete(e.timerStartArms, key)
	if dep == nil || !dep.TimerStartIsCycle(arm.startEventID) {
		return
	}
	due, _, err := dep.TimerDue(arm.startEventID, e.now())
	if err != nil {
		return
	}
	e.timerStartArms[key] = timerStartArm{
		deploymentID: arm.deploymentID,
		startEventID: arm.startEventID,
		dueUnixMs:    due,
	}
}
