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
	MultiInstanceBehaviorAll MultiInstanceBehavior = "All"
	MultiInstanceBehaviorOne MultiInstanceBehavior = "One"
)

// MultiInstanceSpec is compiled loop characteristics for an activity or embedded sub-process.
type MultiInstanceSpec struct {
	ElementID          string
	Sequential         bool
	CardinalityExpr    string // loopCardinality text; empty when collection-driven
	InputCollectionVar string // loopDataInputRef
	InputElementVar    string // inputDataItem name
	OutputCollectionVar string
	OutputElementVar   string
	CompletionExpr     string // empty → default from Behavior
	Behavior           MultiInstanceBehavior
}

func compileMultiInstance(elementID string, lc element.LoopCharacteristicsElements) (MultiInstanceSpec, error) {
	if len(lc.StandardLoopCharacteristics) > 0 {
		return MultiInstanceSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: %q standardLoopCharacteristics not supported", elementID)
	}
	if len(lc.MultielementLoopCharacteristics) == 0 {
		return MultiInstanceSpec{}, fmt.Errorf("NOT_FOUND")
	}
	if len(lc.MultielementLoopCharacteristics) > 1 {
		return MultiInstanceSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: %q has multiple multiInstanceLoopCharacteristics", elementID)
	}
	mi := lc.MultielementLoopCharacteristics[0]
	if len(mi.ComplexBehaviorDefinitions) > 0 {
		return MultiInstanceSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: %q complexBehaviorDefinition not supported", elementID)
	}
	if mi.OneBehaviorEventRef != "" || mi.NoneBehaviorEventRef != "" {
		return MultiInstanceSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: %q none/one behavior event refs not supported", elementID)
	}
	behavior := MultiInstanceBehaviorAll
	switch mi.Behavior {
	case element.MultielementFlowConditionNone, element.MultielementFlowConditionComplex:
		return MultiInstanceSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: %q behavior %q not supported", elementID, mi.Behavior)
	case element.MultielementFlowConditionOne:
		behavior = MultiInstanceBehaviorOne
	case element.MultielementFlowConditionAll, "":
		behavior = MultiInstanceBehaviorAll
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
	}, nil
}

func indexMultiInstance(d *Deployment, id string, lc element.LoopCharacteristicsElements) error {
	spec, err := compileMultiInstance(id, lc)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			return nil
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
			exprText = "nrOfCompletedInstances >= nrOfInstances"
		}
	}
	merged := loopCounterVars(loop, loopCounter)
	for k, v := range vars {
		merged[k] = v
	}
	return expr.Eval(exprText, merged)
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
