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
	"sort"
	"strings"

	"github.com/sparrow-community/sparrow/bpmn/element"
)

// ProcessStartKind classifies a process-level startEvent.
type ProcessStartKind string

const (
	ProcessStartNone        ProcessStartKind = "none"
	ProcessStartMessage     ProcessStartKind = "message"
	ProcessStartTimer       ProcessStartKind = "timer"
	ProcessStartSignal      ProcessStartKind = "signal"
	ProcessStartConditional ProcessStartKind = "conditional"
)

// CreateInstanceEntryID returns the none start or first instantiate entry.
// Typed-only processes return an error — use PublishMessage / PublishSignal / FireDue / EvaluateConditionalStarts.
// Prefer CreateInstanceEntryIDs when multiple instantiate entries must all be armed.
func CreateInstanceEntryID(proc *element.Process) (string, error) {
	ids, err := CreateInstanceEntryIDs(proc)
	if err != nil {
		return "", err
	}
	return ids[0], nil
}

// CreateInstanceEntryIDs returns none start (single) or all instantiate entry element ids.
func CreateInstanceEntryIDs(proc *element.Process) ([]string, error) {
	if proc == nil {
		return nil, fmt.Errorf("no process")
	}
	if id, ok := noneStartEventID(proc); ok {
		return []string{id}, nil
	}
	if len(proc.StartEvents) == 0 {
		return instantiateEntryIDs(proc)
	}
	return nil, fmt.Errorf("INVALID_ARGUMENT: process has no none start; use typed start triggers")
}

// CreateInstanceEntryID returns the CreateInstance entry for this deployment.
func (d *Deployment) CreateInstanceEntryID() (string, error) {
	ids, err := d.CreateInstanceEntryIDs()
	if err != nil {
		return "", err
	}
	return ids[0], nil
}

// CreateInstanceEntryIDs returns all CreateInstance entry element ids.
func (d *Deployment) CreateInstanceEntryIDs() ([]string, error) {
	if d.noneStartID != "" {
		return []string{d.noneStartID}, nil
	}
	if len(d.Process.StartEvents) == 0 {
		return instantiateEntryIDs(&d.Process)
	}
	return nil, fmt.Errorf("INVALID_ARGUMENT: process has no none start; use typed start triggers")
}

func noneStartEventID(proc *element.Process) (string, bool) {
	for _, e := range proc.StartEvents {
		if isNoneStart(e) {
			return e.ID, true
		}
	}
	return "", false
}

func isNoneStart(e element.StartEvent) bool {
	if len(e.TimerEventDefinitions) > 0 {
		return false
	}
	return extraCatchDefinitions(e.EventDefinitions) == 0
}

func (d *Deployment) indexProcessLevelStarts(proc *element.Process, messages []element.Message, signals []element.Signal) error {
	d.noneStartID = ""
	d.messageStarts = make(map[string][]string)
	d.signalStarts = make(map[string][]string)
	d.timerStarts = make(map[string]timerCatch)
	d.conditionalStarts = make(map[string]string)

	for _, e := range proc.StartEvents {
		kind, err := classifyProcessStart(e, messages, signals)
		if err != nil {
			return err
		}
		switch kind {
		case ProcessStartNone:
			if d.noneStartID != "" {
				return fmt.Errorf("UNSUPPORTED_ELEMENT: process has multiple none startEvents (%q and %q)", d.noneStartID, e.ID)
			}
			d.noneStartID = e.ID
		case ProcessStartMessage:
			mc, err := messageCatchFromDefs(e.ID, e.EventDefinitions, messages, strings.TrimSpace(e.Name))
			if err != nil {
				return fmt.Errorf("UNSUPPORTED_ELEMENT: startEvent %q: %v", e.ID, err)
			}
			d.messageStarts[mc.Name] = append(d.messageStarts[mc.Name], e.ID)
		case ProcessStartSignal:
			name, err := signalCatchFromDefs(e.ID, e.EventDefinitions, signals, strings.TrimSpace(e.Name))
			if err != nil {
				return fmt.Errorf("UNSUPPORTED_ELEMENT: startEvent %q: %v", e.ID, err)
			}
			d.signalStarts[name] = append(d.signalStarts[name], e.ID)
		case ProcessStartTimer:
			tc, err := timerCatchFromDefs(e.ID, e.EventDefinitions)
			if err != nil {
				return fmt.Errorf("UNSUPPORTED_ELEMENT: startEvent %q: %v", e.ID, err)
			}
			d.timerStarts[e.ID] = tc
			d.timerCatch[e.ID] = tc
		case ProcessStartConditional:
			text, err := conditionalStartText(e)
			if err != nil {
				return err
			}
			d.conditionalStarts[e.ID] = text
		}
	}
	return nil
}

func classifyProcessStart(e element.StartEvent, messages []element.Message, signals []element.Signal) (ProcessStartKind, error) {
	_ = messages
	_ = signals
	if isNoneStart(e) {
		return ProcessStartNone, nil
	}
	if len(e.ErrorEventDefinitions) > 0 {
		return "", fmt.Errorf("UNSUPPORTED_ELEMENT: process-level startEvent %q error start is not supported (Event Sub-Process only)", e.ID)
	}
	if len(e.EscalationEventDefinitions) > 0 {
		return "", fmt.Errorf("UNSUPPORTED_ELEMENT: process-level startEvent %q escalation start is not supported (Event Sub-Process only)", e.ID)
	}
	if len(e.CompensateEventDefinitions) > 0 || len(e.LinkEventDefinitions) > 0 || len(e.TerminateEventDefinitions) > 0 {
		return "", fmt.Errorf("UNSUPPORTED_ELEMENT: startEvent %q has unsupported event definitions", e.ID)
	}
	nMsg := len(e.MessageEventDefinitions)
	nSig := len(e.SignalEventDefinitions)
	nCond := len(e.ConditionalEventDefinitions)
	nTimer := len(e.TimerEventDefinitions)
	kinds := 0
	if nMsg > 0 {
		kinds++
	}
	if nSig > 0 {
		kinds++
	}
	if nCond > 0 {
		kinds++
	}
	if nTimer > 0 {
		kinds++
	}
	if kinds != 1 {
		return "", fmt.Errorf("UNSUPPORTED_ELEMENT: startEvent %q must have exactly one event definition kind", e.ID)
	}
	if nMsg > 0 {
		if nMsg != 1 || extraCatchDefinitions(e.EventDefinitions)-nMsg > 0 {
			return "", fmt.Errorf("UNSUPPORTED_ELEMENT: startEvent %q must have exactly one messageEventDefinition", e.ID)
		}
		return ProcessStartMessage, nil
	}
	if nSig > 0 {
		if nSig != 1 || extraCatchDefinitions(e.EventDefinitions)-nSig > 0 {
			return "", fmt.Errorf("UNSUPPORTED_ELEMENT: startEvent %q must have exactly one signalEventDefinition", e.ID)
		}
		return ProcessStartSignal, nil
	}
	if nTimer > 0 {
		if nTimer != 1 || extraCatchDefinitions(e.EventDefinitions) > 0 {
			return "", fmt.Errorf("UNSUPPORTED_ELEMENT: startEvent %q must have exactly one timerEventDefinition", e.ID)
		}
		return ProcessStartTimer, nil
	}
	if nCond > 0 {
		if nCond != 1 || extraCatchDefinitions(e.EventDefinitions)-nCond > 0 {
			return "", fmt.Errorf("UNSUPPORTED_ELEMENT: startEvent %q must have exactly one conditionalEventDefinition", e.ID)
		}
		return ProcessStartConditional, nil
	}
	return "", fmt.Errorf("UNSUPPORTED_ELEMENT: startEvent %q has unsupported event definitions", e.ID)
}

func conditionalStartText(e element.StartEvent) (string, error) {
	if len(e.ConditionalEventDefinitions) != 1 {
		return "", fmt.Errorf("UNSUPPORTED_ELEMENT: startEvent %q must have exactly one conditionalEventDefinition", e.ID)
	}
	text := expressionText(e.ConditionalEventDefinitions[0].Condition)
	if text == "" {
		return "", fmt.Errorf("UNSUPPORTED_ELEMENT: startEvent %q conditional start needs a condition expression", e.ID)
	}
	return text, nil
}

// MessageStartIDs returns process-level message start event ids for name.
func (d *Deployment) MessageStartIDs(name string) []string {
	return append([]string(nil), d.messageStarts[name]...)
}

// MessageStartNames returns sorted process-level message start names.
func (d *Deployment) MessageStartNames() []string {
	out := make([]string, 0, len(d.messageStarts))
	for name := range d.messageStarts {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// SignalStartIDs returns process-level signal start event ids for name.
func (d *Deployment) SignalStartIDs(name string) []string {
	return append([]string(nil), d.signalStarts[name]...)
}

// SignalStartNames returns sorted process-level signal start names.
func (d *Deployment) SignalStartNames() []string {
	out := make([]string, 0, len(d.signalStarts))
	for name := range d.signalStarts {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// TimerStartIDs returns process-level timer start event ids.
func (d *Deployment) TimerStartIDs() []string {
	out := make([]string, 0, len(d.timerStarts))
	for id := range d.timerStarts {
		out = append(out, id)
	}
	return out
}

// ConditionalStartIDs returns process-level conditional start event ids.
func (d *Deployment) ConditionalStartIDs() []string {
	out := make([]string, 0, len(d.conditionalStarts))
	for id := range d.conditionalStarts {
		out = append(out, id)
	}
	return out
}

// ConditionalStartExpression returns the condition text for a conditional start.
func (d *Deployment) ConditionalStartExpression(startEventID string) (string, bool) {
	text, ok := d.conditionalStarts[startEventID]
	return text, ok
}

// HasTypedStarts reports whether the deployment has any typed process-level starts.
func (d *Deployment) HasTypedStarts() bool {
	return len(d.messageStarts) > 0 || len(d.signalStarts) > 0 || len(d.timerStarts) > 0 || len(d.conditionalStarts) > 0
}

// TimerStartIsCycle reports whether a process-level timer start uses timeCycle.
func (d *Deployment) TimerStartIsCycle(startEventID string) bool {
	tc, ok := d.timerStarts[startEventID]
	return ok && tc.Cycle != nil
}
