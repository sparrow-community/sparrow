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

package deploy

import (
	"fmt"
	"strings"

	"github.com/sparrow-community/sparrow/bpmn/element"
)

func conditionalCatchFromDefs(id string, defs element.EventDefinitions) (string, error) {
	if len(defs.TimerEventDefinitions) > 0 {
		return "", fmt.Errorf("UNSUPPORTED_ELEMENT: %q is not a conditional catch (has timerEventDefinition)", id)
	}
	if len(defs.ConditionalEventDefinitions) != 1 {
		return "", fmt.Errorf("UNSUPPORTED_ELEMENT: %q needs exactly one conditionalEventDefinition", id)
	}
	other := extraCatchDefinitions(defs) - len(defs.ConditionalEventDefinitions)
	if other > 0 {
		return "", fmt.Errorf("UNSUPPORTED_ELEMENT: %q is not a conditional catch", id)
	}
	text := strings.TrimSpace(expressionText(defs.ConditionalEventDefinitions[0].Condition))
	if text == "" {
		// Empty condition stub (common in MIWG interchange): treat as always-true,
		// parallel to empty timer expression → PT0S.
		return "true", nil
	}
	return text, nil
}

func conditionalCatchSpec(ev element.IntermediateCatchEvent) (string, error) {
	if len(ev.TimerEventDefinitions) > 0 {
		return "", fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateCatchEvent %q is not a conditional catch", ev.ID)
	}
	return conditionalCatchFromDefs(ev.ID, ev.EventDefinitions)
}

type conditionalBoundary struct {
	Condition    string
	AttachedTo   string
	Interrupting bool
}

func conditionalBoundarySpec(ev element.BoundaryEvent) (conditionalBoundary, error) {
	if err := requireBoundaryAttach(ev); err != nil {
		return conditionalBoundary{}, err
	}
	if len(ev.TimerEventDefinitions) > 0 {
		return conditionalBoundary{}, fmt.Errorf("UNSUPPORTED_ELEMENT: boundaryEvent %q must be a conditional boundary", ev.ID)
	}
	text, err := conditionalCatchFromDefs(ev.ID, ev.EventDefinitions)
	if err != nil {
		return conditionalBoundary{}, err
	}
	return conditionalBoundary{
		Condition:    text,
		AttachedTo:   strings.TrimSpace(ev.AttachedToRef),
		Interrupting: ev.CancelActivity,
	}, nil
}

// ConditionalCatchExpression returns the condition text for a conditional catch or boundary.
func (d *Deployment) ConditionalCatchExpression(id string) (string, bool) {
	if d == nil {
		return "", false
	}
	text, ok := d.conditionalCatch[id]
	return text, ok && text != ""
}

// ConditionalBoundaries returns conditional boundary ids attached to activityID.
func (d *Deployment) ConditionalBoundaries(activityID string) []string {
	var ids []string
	for _, e := range d.Process.BoundaryEvents {
		if e.AttachedToRef != activityID {
			continue
		}
		if _, ok := d.conditionalCatch[e.ID]; ok {
			ids = append(ids, e.ID)
		}
	}
	return ids
}
