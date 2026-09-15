package deploy

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/sparrow-community/sparrow/bpmn/element"
	"github.com/sparrow-community/sparrow/processing/expr"
	"github.com/sparrow-community/sparrow/processing/projection"
)

// MultiInstanceBehavior is the default join semantics when no completionCondition is set.
type MultiInstanceBehavior string

const (
	MultiInstanceBehaviorAll     MultiInstanceBehavior = "All"
	MultiInstanceBehaviorOne     MultiInstanceBehavior = "One"
	MultiInstanceBehaviorNone    MultiInstanceBehavior = "None"
	MultiInstanceBehaviorComplex MultiInstanceBehavior = "Complex"
)

// MIBehaviorThrow is a signal/message published at an MI milestone.
type MIBehaviorThrow struct {
	Kind ThrowKind // signal or message
	Name string
}

// MIComplexBehavior is one complexBehaviorDefinition.
type MIComplexBehavior struct {
	Condition string
	Throw     MIBehaviorThrow
}

// MultiInstanceSpec is compiled loop characteristics for an activity or embedded sub-process.
type MultiInstanceSpec struct {
	ElementID           string
	Sequential          bool
	CardinalityExpr     string // loopCardinality text; empty when collection-driven
	InputCollectionVar  string // loopDataInputRef
	InputElementVar     string // inputDataItem name
	OutputCollectionVar string
	OutputElementVar    string
	CompletionExpr      string // empty → default from Behavior
	Behavior            MultiInstanceBehavior
	NoneEvent           *MIBehaviorThrow
	OneEvent            *MIBehaviorThrow
	ComplexBehaviors    []MIComplexBehavior
}

func compileMultiInstance(elementID string, lc element.LoopCharacteristicsElements, messages []element.Message, signals []element.Signal) (MultiInstanceSpec, error) {
	if len(lc.StandardLoopCharacteristics) > 0 && len(lc.MultielementLoopCharacteristics) > 0 {
		return MultiInstanceSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: %q cannot combine standardLoopCharacteristics and multiInstanceLoopCharacteristics", elementID)
	}
	if len(lc.StandardLoopCharacteristics) > 0 {
		return MultiInstanceSpec{}, fmt.Errorf("NOT_FOUND")
	}
	if len(lc.MultielementLoopCharacteristics) == 0 {
		return MultiInstanceSpec{}, fmt.Errorf("NOT_FOUND")
	}
	if len(lc.MultielementLoopCharacteristics) > 1 {
		return MultiInstanceSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: %q has multiple multiInstanceLoopCharacteristics", elementID)
	}
	mi := lc.MultielementLoopCharacteristics[0]
	behavior := MultiInstanceBehaviorAll
	switch mi.Behavior {
	case element.MultielementFlowConditionNone:
		behavior = MultiInstanceBehaviorNone
	case element.MultielementFlowConditionComplex:
		behavior = MultiInstanceBehaviorComplex
	case element.MultielementFlowConditionOne:
		behavior = MultiInstanceBehaviorOne
	case element.MultielementFlowConditionAll, "":
		behavior = MultiInstanceBehaviorAll
	default:
		return MultiInstanceSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: %q behavior %q not supported", elementID, mi.Behavior)
	}
	if behavior == MultiInstanceBehaviorComplex && len(mi.ComplexBehaviorDefinitions) == 0 {
		return MultiInstanceSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: %q behavior Complex requires complexBehaviorDefinition", elementID)
	}
	if behavior != MultiInstanceBehaviorComplex && len(mi.ComplexBehaviorDefinitions) > 0 {
		return MultiInstanceSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: %q complexBehaviorDefinition requires behavior Complex", elementID)
	}

	noneEvent, err := resolveMIBehaviorEventRef(elementID, "noneBehaviorEventRef", mi.NoneBehaviorEventRef, messages, signals)
	if err != nil {
		return MultiInstanceSpec{}, err
	}
	oneEvent, err := resolveMIBehaviorEventRef(elementID, "oneBehaviorEventRef", mi.OneBehaviorEventRef, messages, signals)
	if err != nil {
		return MultiInstanceSpec{}, err
	}
	complex, err := compileComplexBehaviors(elementID, mi.ComplexBehaviorDefinitions, messages, signals)
	if err != nil {
		return MultiInstanceSpec{}, err
	}

	inputVar := strings.TrimSpace(mi.InputDataItem.Name)
	if inputVar == "" {
		inputVar = strings.TrimSpace(mi.InputDataItem.ID)
	}
	outputVar := strings.TrimSpace(mi.OutputDataItem.Name)
	if outputVar == "" {
		outputVar = strings.TrimSpace(mi.OutputDataItem.ID)
	}
	cardExpr := expressionText(mi.LoopCardinality)
	inputColl := strings.TrimSpace(mi.LoopDataInputRef)
	if cardExpr != "" && inputColl != "" {
		return MultiInstanceSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: %q has both loopCardinality and loopDataInputRef", elementID)
	}
	return MultiInstanceSpec{
		ElementID:           elementID,
		Sequential:          mi.IsSequential,
		CardinalityExpr:     cardExpr,
		InputCollectionVar:  inputColl,
		InputElementVar:     inputVar,
		OutputCollectionVar: strings.TrimSpace(mi.LoopDataOutputRef),
		OutputElementVar:    outputVar,
		CompletionExpr:      expressionText(mi.CompletionCondition),
		Behavior:            behavior,
		NoneEvent:           noneEvent,
		OneEvent:            oneEvent,
		ComplexBehaviors:    complex,
	}, nil
}

func compileComplexBehaviors(elementID string, defs []element.ComplexBehaviorDefinition, messages []element.Message, signals []element.Signal) ([]MIComplexBehavior, error) {
	out := make([]MIComplexBehavior, 0, len(defs))
	for i, def := range defs {
		cond := expressionText(def.Condition)
		if cond == "" {
			return nil, fmt.Errorf("UNSUPPORTED_ELEMENT: %q complexBehaviorDefinition[%d] needs a condition", elementID, i)
		}
		th, err := implicitThrowToMIBehavior(elementID, i, def.Event, messages, signals)
		if err != nil {
			return nil, err
		}
		out = append(out, MIComplexBehavior{Condition: cond, Throw: th})
	}
	return out, nil
}

func implicitThrowToMIBehavior(elementID string, idx int, ev element.ImplicitThrowEvent, messages []element.Message, signals []element.Signal) (MIBehaviorThrow, error) {
	defs := ev.EventDefinitions
	switch {
	case len(defs.SignalEventDefinitions) == 1 && extraCatchDefinitions(defs) == 1:
		name := resolveSignalName(defs.SignalEventDefinitions[0].SignalRef, signals)
		if name == "" {
			name = strings.TrimSpace(ev.Name)
		}
		if name == "" {
			name = ev.ID
		}
		if name == "" {
			return MIBehaviorThrow{}, fmt.Errorf("UNSUPPORTED_ELEMENT: %q complexBehaviorDefinition[%d] signal throw needs a name", elementID, idx)
		}
		return MIBehaviorThrow{Kind: ThrowKindSignal, Name: name}, nil
	case len(defs.MessageEventDefinitions) == 1 && extraCatchDefinitions(defs) == 1:
		name := resolveMessageName(defs.MessageEventDefinitions[0].MessageRef, messages)
		if name == "" {
			name = strings.TrimSpace(ev.Name)
		}
		if name == "" {
			name = ev.ID
		}
		if name == "" {
			return MIBehaviorThrow{}, fmt.Errorf("UNSUPPORTED_ELEMENT: %q complexBehaviorDefinition[%d] message throw needs a name", elementID, idx)
		}
		return MIBehaviorThrow{Kind: ThrowKindMessage, Name: name}, nil
	default:
		return MIBehaviorThrow{}, fmt.Errorf("UNSUPPORTED_ELEMENT: %q complexBehaviorDefinition[%d] event must be signal or message", elementID, idx)
	}
}

func resolveMIBehaviorEventRef(elementID, attr, ref string, messages []element.Message, signals []element.Signal) (*MIBehaviorThrow, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return nil, nil
	}
	for _, s := range signals {
		if s.ID == ref {
			name := strings.TrimSpace(s.Name)
			if name == "" {
				name = s.ID
			}
			return &MIBehaviorThrow{Kind: ThrowKindSignal, Name: name}, nil
		}
	}
	for _, m := range messages {
		if m.ID == ref {
			name := strings.TrimSpace(m.Name)
			if name == "" {
				name = m.ID
			}
			return &MIBehaviorThrow{Kind: ThrowKindMessage, Name: name}, nil
		}
	}
	return nil, fmt.Errorf("UNSUPPORTED_ELEMENT: %q %s %q must reference a signal or message", elementID, attr, ref)
}

func indexMultiInstance(d *Deployment, id string, lc element.LoopCharacteristicsElements) error {
	spec, err := compileMultiInstance(id, lc, d.messages, d.signals)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			return indexStandardLoop(d, id, lc)
		}
		return err
	}
	if d.multiInstances == nil {
		d.multiInstances = make(map[string]MultiInstanceSpec)
	}
	d.multiInstances[id] = spec
	return nil
}

// MultiInstanceSpec returns compiled loop characteristics for an element id.
func (d *Deployment) MultiInstanceSpec(id string) (MultiInstanceSpec, bool) {
	if d == nil {
		return MultiInstanceSpec{}, false
	}
	spec, ok := d.multiInstances[id]
	return spec, ok
}

// InstanceCount resolves how many inner instances to start.
func (spec MultiInstanceSpec) InstanceCount(vars map[string]string) (int, error) {
	if spec.InputCollectionVar != "" {
		raw, ok := vars[spec.InputCollectionVar]
		if !ok || strings.TrimSpace(raw) == "" {
			return 0, nil
		}
		var arr []json.RawMessage
		if err := json.Unmarshal([]byte(raw), &arr); err != nil {
			return 0, nil
		}
		return len(arr), nil
	}
	if strings.TrimSpace(spec.CardinalityExpr) == "" {
		return 0, nil
	}
	s := strings.TrimSpace(spec.CardinalityExpr)
	if strings.HasPrefix(s, "${") && strings.HasSuffix(s, "}") {
		s = strings.TrimSpace(s[2 : len(s)-1])
	}
	if n, err := strconv.Atoi(s); err == nil {
		if n < 0 {
			return 0, nil
		}
		return n, nil
	}
	if raw, ok := vars[s]; ok {
		var num float64
		if err := json.Unmarshal([]byte(raw), &num); err == nil {
			if num < 0 {
				return 0, nil
			}
			return int(num), nil
		}
	}
	return 0, fmt.Errorf("cardinality %q did not evaluate to an integer", spec.CardinalityExpr)
}

// CompletionMet evaluates the completion condition after an inner instance completes.
func (spec MultiInstanceSpec) CompletionMet(loop *projection.MultiInstanceLoop, vars map[string]string, loopCounter int32) (bool, error) {
	exprText := strings.TrimSpace(spec.CompletionExpr)
	if exprText == "" {
		switch spec.Behavior {
		case MultiInstanceBehaviorOne:
			exprText = "nrOfCompletedInstances >= 1"
		default:
			// None / All / Complex: wait for all unless completionCondition set.
			exprText = "nrOfCompletedInstances >= nrOfInstances"
		}
	}
	merged := loopCounterVars(loop, loopCounter)
	for k, v := range vars {
		merged[k] = v
	}
	return expr.Eval(exprText, merged)
}

// ComplexConditionMet evaluates one complexBehaviorDefinition condition.
func (spec MultiInstanceSpec) ComplexConditionMet(loop *projection.MultiInstanceLoop, vars map[string]string, loopCounter int32, condition string) (bool, error) {
	merged := loopCounterVars(loop, loopCounter)
	for k, v := range vars {
		merged[k] = v
	}
	return expr.Eval(condition, merged)
}

func loopCounterVars(loop *projection.MultiInstanceLoop, loopCounter int32) map[string]string {
	total, active, completed := 0, 0, 0
	if loop != nil {
		total = int(loop.TotalInstances)
		active = int(loop.ActiveInstances)
		completed = int(loop.CompletedInstances)
	}
	b := func(n int) string {
		out, _ := json.Marshal(n)
		return string(out)
	}
	lc, _ := json.Marshal(int(loopCounter))
	return map[string]string{
		"loopCounter":                string(lc),
		"nrOfInstances":              b(total),
		"nrOfActiveInstances":        b(active),
		"nrOfCompletedInstances":     b(completed),
		"numberOfInstances":          b(total),
		"numberOfActiveInstances":    b(active),
		"numberOfCompletedInstances": b(completed),
	}
}

// CollectionElementVariables returns per-iteration input variables.
func (spec MultiInstanceSpec) CollectionElementVariables(vars map[string]string, index int) map[string]string {
	if spec.InputCollectionVar == "" || spec.InputElementVar == "" {
		return nil
	}
	raw, ok := vars[spec.InputCollectionVar]
	if !ok {
		return nil
	}
	var arr []json.RawMessage
	if err := json.Unmarshal([]byte(raw), &arr); err != nil {
		return nil
	}
	if index < 0 || index >= len(arr) {
		return nil
	}
	return map[string]string{spec.InputElementVar: string(arr[index])}
}