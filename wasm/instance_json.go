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

//go:build js && wasm

package main

import (
	"github.com/sparrow-community/sparrow/processing/projection"
)

func instanceJSON(inst *projection.Instance) map[string]any {
	tokens := make(map[string]any, len(inst.Tokens))
	for id, tok := range inst.Tokens {
		if tok == nil {
			continue
		}
		waits := make([]map[string]any, 0, len(tok.BoundaryWaits))
		for _, w := range tok.BoundaryWaits {
			waits = append(waits, map[string]any{
				"boundaryId":  w.BoundaryID,
				"kind":        w.Kind,
				"dueUnixMs":   w.DueUnixMs,
				"timerText":   w.TimerText,
				"messageName": w.MessageName,
				"signalName":  w.SignalName,
			})
		}
		tokens[id] = map[string]any{
			"id":                      tok.ID,
			"elementId":               tok.ElementID,
			"status":                  string(tok.Status),
			"jobType":                 tok.JobType,
			"dueUnixMs":               tok.DueUnixMs,
			"timerText":               tok.TimerText,
			"messageName":             tok.MessageName,
			"signalName":              tok.SignalName,
			"boundaryId":              tok.BoundaryID,
			"messageBoundaryId":       tok.MessageBoundaryID,
			"signalBoundaryId":        tok.SignalBoundaryID,
			"scopeHost":               tok.ScopeHost,
			"calledProcessInstanceId": tok.CalledProcessInstanceID,
			"loopInstanceIndex":       tok.LoopInstanceIndex,
			"multiInstanceHost":       tok.MultiInstanceHost,
			"boundaryWaits":           waits,
			"jobFailCount":            tok.JobFailCount,
			"incidentErrorMessage":    tok.IncidentErrorMessage,
		}
	}
	intents := make(map[string]string, len(inst.ElementIntent))
	for id, intent := range inst.ElementIntent {
		intents[id] = intent.String()
	}
	return map[string]any{
		"instanceId":              inst.ID,
		"deploymentId":            inst.DeploymentID,
		"processId":               inst.ProcessID,
		"version":                 inst.Version,
		"status":                  string(inst.Status),
		"parentProcessInstanceId": inst.ParentProcessInstanceID,
		"parentElementId":         inst.ParentElementID,
		"parentTokenId":           inst.ParentTokenID,
		"variables":               inst.Variables,
		"tokens":                  tokens,
		"elementIntent":           intents,
	}
}
