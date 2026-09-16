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
	"github.com/sparrow-community/sparrow/processing/expr"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// AdHocSubProcessSpec holds compiled Ad-Hoc SubProcess facts.
type AdHocSubProcessSpec struct {
	ID string
	// CompletionCondition is the expression that finishes the scope.
	CompletionCondition string
	// Sequential enables one inner activity at a time in document order.
	Sequential bool
	// CancelRemaining terminates inner activities still running once the
	// completion condition is true (BPMN default true).
	CancelRemaining bool
	// InnerActivityIDs are the inner activities in document order.
	InnerActivityIDs []string
}

func adHocSubProcessSpec(a element.AdHocSubProcess) AdHocSubProcessSpec {
	return AdHocSubProcessSpec{
		ID:                  a.ID,
		CompletionCondition: expressionText(a.CompletionCondition),
		Sequential:          strings.TrimSpace(a.Ordering) == element.AdHocOrderingSequential,
		CancelRemaining:     a.CancelRemaining(),
		InnerActivityIDs:    adHocInnerActivityIDs(a),
	}
}

// adHocInnerActivityIDs lists supported inner activities in document order.
// Ordering follows the BPMN document, so an id maps to a stable position for
// Sequential ordering across Recover.
func adHocInnerActivityIDs(a element.AdHocSubProcess) []string {
	fe := a.FlowElements
	ids := make([]string, 0,
		len(fe.Tasks)+len(fe.UserTasks)+len(fe.ManualTasks)+len(fe.ServiceTasks)+
			len(fe.ReceiveTasks)+len(fe.BusinessRuleTasks)+len(fe.ScriptTasks))
	for _, e := range fe.Tasks {
		ids = append(ids, e.ID)
	}
	for _, e := range fe.UserTasks {
		ids = append(ids, e.ID)
	}
	for _, e := range fe.ManualTasks {
		ids = append(ids, e.ID)
	}
	for _, e := range fe.ServiceTasks {
		ids = append(ids, e.ID)
	}
	for _, e := range fe.ReceiveTasks {
		ids = append(ids, e.ID)
	}
	for _, e := range fe.BusinessRuleTasks {
		ids = append(ids, e.ID)
	}
	for _, e := range fe.ScriptTasks {
		ids = append(ids, e.ID)
	}
	return ids
}

func validateAdHocSubProcess(a element.AdHocSubProcess) error {
	if a.TriggeredByEvent {
		return fmt.Errorf("UNSUPPORTED_ELEMENT: adHocSubProcess %q must not be triggeredByEvent", a.ID)
	}
	if len(a.LoopCharacteristicsElements.MultielementLoopCharacteristics) > 0 ||
		len(a.LoopCharacteristicsElements.StandardLoopCharacteristics) > 0 {
		return fmt.Errorf("UNSUPPORTED_ELEMENT: multi-instance / standard loop on adHocSubProcess %q is not supported", a.ID)
	}
	switch strings.TrimSpace(a.Ordering) {
	case "", element.AdHocOrderingParallel, element.AdHocOrderingSequential:
	default:
		return fmt.Errorf("UNSUPPORTED_ELEMENT: adHocSubProcess %q ordering %q not supported (Parallel or Sequential)", a.ID, a.Ordering)
	}
	if expressionText(a.CompletionCondition) == "" {
		return fmt.Errorf("UNSUPPORTED_ELEMENT: adHocSubProcess %q needs a completionCondition", a.ID)
	}
	if err := validateAdHocBody(a); err != nil {
		return err
	}
	if len(adHocInnerActivityIDs(a)) == 0 {
		return fmt.Errorf("UNSUPPORTED_ELEMENT: adHocSubProcess %q needs at least one inner activity", a.ID)
	}
	return nil
}

// validateAdHocBody keeps the body a flat set of waiting / job activities: an
// inner activity is enabled directly, so flows, events, gateways, nested scopes
// and per-activity loops have no defined runtime here.
func validateAdHocBody(a element.AdHocSubProcess) error {
	fe := a.FlowElements
	reject := func(what string) error {
		return fmt.Errorf("UNSUPPORTED_ELEMENT: adHocSubProcess %q must not contain %s", a.ID, what)
	}
	switch {
	case len(fe.SequenceFlows) > 0:
		return reject("sequence flows")
	case len(fe.StartEvents) > 0 || len(fe.EndEvents) > 0:
		return reject("start or end events")
	case len(fe.IntermediateCatchEvents) > 0 || len(fe.IntermediateThrowEvents) > 0 || len(fe.ImplicitThrowEvents) > 0:
		return reject("intermediate events")
	case len(fe.BoundaryEvents) > 0:
		return reject("boundary events")
	case len(fe.ExclusiveGatewaies) > 0 || len(fe.ParallelGatewaies) > 0 ||
		len(fe.InclusiveGatewaies) > 0 || len(fe.ComplexGatewaies) > 0 || len(fe.EventBasedGatewaies) > 0:
		return reject("gateways")
	case len(fe.SubProcesses) > 0 || len(fe.Transactions) > 0 || len(fe.AdHocSubProcesses) > 0 || len(fe.CallActivities) > 0:
		return reject("nested scopes")
	case len(fe.SendTasks) > 0:
		return reject("send tasks")
	}
	for _, e := range fe.Tasks {
		if err := rejectAdHocInnerLoop(a.ID, e.ID, e.LoopCharacteristicsElements); err != nil {
			return err
		}
	}
	for _, e := range fe.UserTasks {
		if err := rejectAdHocInnerLoop(a.ID, e.ID, e.LoopCharacteristicsElements); err != nil {
			return err
		}
	}
	for _, e := range fe.ManualTasks {
		if err := rejectAdHocInnerLoop(a.ID, e.ID, e.LoopCharacteristicsElements); err != nil {
			return err
		}
	}
	for _, e := range fe.ServiceTasks {
		if err := rejectAdHocInnerLoop(a.ID, e.ID, e.LoopCharacteristicsElements); err != nil {
			return err
		}
	}
	for _, e := range fe.ReceiveTasks {
		if err := rejectAdHocInnerLoop(a.ID, e.ID, e.LoopCharacteristicsElements); err != nil {
			return err
		}
	}
	for _, e := range fe.BusinessRuleTasks {
		if err := rejectAdHocInnerLoop(a.ID, e.ID, e.LoopCharacteristicsElements); err != nil {
			return err
		}
	}
	for _, e := range fe.ScriptTasks {
		if err := rejectAdHocInnerLoop(a.ID, e.ID, e.LoopCharacteristicsElements); err != nil {
			return err
		}
	}
	return nil
}

func rejectAdHocInnerLoop(adHocID, activityID string, loop element.LoopCharacteristicsElements) error {
	if len(loop.MultielementLoopCharacteristics) > 0 || len(loop.StandardLoopCharacteristics) > 0 {
		return fmt.Errorf("UNSUPPORTED_ELEMENT: adHocSubProcess %q inner activity %q must not have multi-instance / standard loop", adHocID, activityID)
	}
	return nil
}

// IsAdHocSubProcess reports whether id is a BPMN adHocSubProcess.
func (d *Deployment) IsAdHocSubProcess(id string) bool {
	if d == nil {
		return false
	}
	t, err := d.TypeOf(id)
	return err == nil && t == eventv1.Element_TYPE_AD_HOC_SUB_PROCESS
}

// AdHocSubProcessSpecOf returns the compiled spec for an Ad-Hoc SubProcess.
func (d *Deployment) AdHocSubProcessSpecOf(id string) (AdHocSubProcessSpec, bool) {
	if d == nil {
		return AdHocSubProcessSpec{}, false
	}
	spec, ok := d.adHocSubProcesses[id]
	return spec, ok
}

// AdHocScopeOf returns the enclosing Ad-Hoc SubProcess of an inner activity.
func (d *Deployment) AdHocScopeOf(activityID string) (string, bool) {
	if d == nil {
		return "", false
	}
	scope, ok := d.ScopeOf(activityID)
	if !ok || !d.IsAdHocSubProcess(scope) {
		return "", false
	}
	return scope, true
}

// AdHocInitialActivities lists the inner activities enabled when the scope opens.
func (s AdHocSubProcessSpec) AdHocInitialActivities() []string {
	if len(s.InnerActivityIDs) == 0 {
		return nil
	}
	if s.Sequential {
		return []string{s.InnerActivityIDs[0]}
	}
	return append([]string{}, s.InnerActivityIDs...)
}

// NextSequentialActivity returns the inner activity following activityID in
// document order.
func (s AdHocSubProcessSpec) NextSequentialActivity(activityID string) (string, bool) {
	for i, id := range s.InnerActivityIDs {
		if id != activityID {
			continue
		}
		if i+1 < len(s.InnerActivityIDs) {
			return s.InnerActivityIDs[i+1], true
		}
		return "", false
	}
	return "", false
}

// EvalAdHocCompletion evaluates the completion condition against instance variables.
func (d *Deployment) EvalAdHocCompletion(adHocID string, vars map[string]string) (bool, error) {
	spec, ok := d.AdHocSubProcessSpecOf(adHocID)
	if !ok {
		return false, fmt.Errorf("NOT_FOUND: adHocSubProcess %q", adHocID)
	}
	ok, err := expr.Eval(spec.CompletionCondition, vars)
	if err != nil {
		return false, fmt.Errorf("INVALID_EXPRESSION: adHocSubProcess %q completionCondition: %w", adHocID, err)
	}
	return ok, nil
}
