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

type escalationBoundary struct {
	EscalationCode string
	AttachedTo     string
	Interrupting   bool
}

func resolveEscalationCode(escalationRef string, escalations []element.Escalation) string {
	if escalationRef == "" {
		return ""
	}
	for _, e := range escalations {
		if e.ID != escalationRef {
			continue
		}
		if e.EscalationCode != "" {
			return e.EscalationCode
		}
		if e.Name != "" {
			return e.Name
		}
		return e.ID
	}
	return escalationRef
}

func escalationMatches(boundaryCode, thrownCode string) bool {
	if boundaryCode == "" {
		return true
	}
	return boundaryCode == thrownCode
}

func escalationBoundarySpec(ev element.BoundaryEvent, escalations []element.Escalation) (escalationBoundary, error) {
	if err := requireBoundaryAttach(ev); err != nil {
		return escalationBoundary{}, err
	}
	if len(ev.EscalationEventDefinitions) != 1 {
		return escalationBoundary{}, fmt.Errorf("UNSUPPORTED_ELEMENT: boundaryEvent %q must be an escalation boundary", ev.ID)
	}
	other := extraCatchDefinitions(ev.EventDefinitions) - len(ev.EscalationEventDefinitions)
	if other > 0 {
		return escalationBoundary{}, fmt.Errorf("UNSUPPORTED_ELEMENT: boundaryEvent %q must be an escalation boundary", ev.ID)
	}
	code := resolveEscalationCode(ev.EscalationEventDefinitions[0].EscalationRef, escalations)
	return escalationBoundary{
		EscalationCode: code,
		AttachedTo:     ev.AttachedToRef,
		Interrupting:   ev.CancelActivity,
	}, nil
}

func escalationEndSpec(ev element.EndEvent, escalations []element.Escalation) (string, error) {
	if len(ev.EscalationEventDefinitions) != 1 {
		return "", fmt.Errorf("not an escalation end")
	}
	other := extraCatchDefinitions(ev.EventDefinitions) - len(ev.EscalationEventDefinitions)
	if other > 0 {
		return "", fmt.Errorf("not an escalation end")
	}
	return resolveEscalationCode(ev.EscalationEventDefinitions[0].EscalationRef, escalations), nil
}

func escalationThrowSpec(ev element.IntermediateThrowEvent, escalations []element.Escalation) (throwSpec, error) {
	if len(ev.TimerEventDefinitions) > 0 {
		return throwSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateThrowEvent %q cannot have timerEventDefinition", ev.ID)
	}
	if len(ev.EscalationEventDefinitions) != 1 {
		return throwSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateThrowEvent %q needs exactly one escalationEventDefinition", ev.ID)
	}
	other := extraCatchDefinitions(ev.EventDefinitions) - len(ev.EscalationEventDefinitions)
	if other > 0 {
		return throwSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateThrowEvent %q is not an escalation throw", ev.ID)
	}
	code := resolveEscalationCode(ev.EscalationEventDefinitions[0].EscalationRef, escalations)
	if code == "" {
		code = strings.TrimSpace(ev.Name)
	}
	if code == "" {
		code = ev.ID
	}
	return throwSpec{Kind: ThrowKindEscalation, Name: code}, nil
}

func escalationStartCatchFromDefs(startEventID string, defs element.EventDefinitions, escalations []element.Escalation) (string, error) {
	if len(defs.EscalationEventDefinitions) != 1 {
		return "", fmt.Errorf("startEvent %q must have exactly one escalationEventDefinition", startEventID)
	}
	other := extraCatchDefinitions(defs) - len(defs.EscalationEventDefinitions)
	if other > 0 {
		return "", fmt.Errorf("startEvent %q must be an escalation start", startEventID)
	}
	return resolveEscalationCode(defs.EscalationEventDefinitions[0].EscalationRef, escalations), nil
}

// MatchEscalationBoundary returns a boundary on activityID that catches escalationCode.
func (d *Deployment) MatchEscalationBoundary(activityID, escalationCode string) (string, bool) {
	for _, bid := range d.escalationBoundaries[activityID] {
		if escalationMatches(d.escalationCatch[bid], escalationCode) {
			return bid, true
		}
	}
	return "", false
}

// EscalationEndCode reports whether an end event throws an escalation and its code.
func (d *Deployment) EscalationEndCode(endEventID string) (string, bool) {
	code, ok := d.escalationEnds[endEventID]
	return code, ok
}

// IsEscalationBoundary reports whether id is an escalation boundary event.
func (d *Deployment) IsEscalationBoundary(boundaryID string) bool {
	_, ok := d.escalationCatch[boundaryID]
	return ok
}

// EscalationBoundaryCode returns the escalation code a boundary catches (empty = catch-all).
func (d *Deployment) EscalationBoundaryCode(boundaryID string) (string, bool) {
	code, ok := d.escalationCatch[boundaryID]
	return code, ok
}

// MatchEscalationEventSubProcess returns an armed escalation ESP in scopeID that catches code.
func (d *Deployment) MatchEscalationEventSubProcess(scopeID, escalationCode string, isArmed func(eventSubProcessID string) bool) (string, bool) {
	var catchAll string
	for _, sp := range d.EventSubProcessesInScope(scopeID) {
		if sp.Kind != CatchKindEscalation {
			continue
		}
		if isArmed != nil && !isArmed(sp.ID) {
			continue
		}
		if sp.EscalationCode == escalationCode {
			return sp.ID, true
		}
		if sp.EscalationCode == "" && catchAll == "" {
			catchAll = sp.ID
		}
	}
	if catchAll != "" {
		return catchAll, true
	}
	return "", false
}
