package deploy

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sparrow-community/sparrow/bpmn/element"
	"github.com/sparrow-community/sparrow/processing/expr"
)

// StandardLoopSpec is compiled standardLoopCharacteristics for an activity.
type StandardLoopSpec struct {
	ElementID     string
	TestBefore    bool
	LoopMaximum   int // 0 = unset
	ConditionExpr string
}

func compileStandardLoop(elementID string, lc element.LoopCharacteristicsElements) (StandardLoopSpec, error) {
	if len(lc.StandardLoopCharacteristics) == 0 {
		return StandardLoopSpec{}, fmt.Errorf("NOT_FOUND")
	}
	if len(lc.MultielementLoopCharacteristics) > 0 {
		return StandardLoopSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: %q cannot combine standardLoopCharacteristics and multiInstanceLoopCharacteristics", elementID)
	}
	if len(lc.StandardLoopCharacteristics) > 1 {
		return StandardLoopSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: %q has multiple standardLoopCharacteristics", elementID)
	}
	st := lc.StandardLoopCharacteristics[0]
	return StandardLoopSpec{
		ElementID:     elementID,
		TestBefore:    st.TestBefore,
		LoopMaximum:   st.LoopMaximum,
		ConditionExpr: expressionText(st.LoopCondition),
	}, nil
}

func indexStandardLoop(d *Deployment, id string, lc element.LoopCharacteristicsElements) error {
	spec, err := compileStandardLoop(id, lc)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			return nil
		}
		return err
	}
	if d.standardLoops == nil {
		d.standardLoops = make(map[string]StandardLoopSpec)
	}
	d.standardLoops[id] = spec
	return nil
}

// StandardLoopSpec returns compiled standard loop characteristics for an element id.
func (d *Deployment) StandardLoopSpec(id string) (StandardLoopSpec, bool) {
	if d == nil {
		return StandardLoopSpec{}, false
	}
	spec, ok := d.standardLoops[id]
	return spec, ok
}

// ShouldEnterFirst reports whether the first iteration should run (testBefore gate).
func (spec StandardLoopSpec) ShouldEnterFirst(vars map[string]string) (bool, error) {
	if !spec.TestBefore {
		return true, nil
	}
	if spec.LoopMaximum > 0 && 1 > spec.LoopMaximum {
		return false, nil
	}
	if strings.TrimSpace(spec.ConditionExpr) == "" {
		return true, nil
	}
	return spec.evalCondition(vars, 1)
}

// ShouldContinueAfter reports whether another iteration should run after completedCounter executions.
func (spec StandardLoopSpec) ShouldContinueAfter(vars map[string]string, completedCounter int) (bool, error) {
	if completedCounter <= 0 {
		return false, nil
	}
	if spec.LoopMaximum > 0 && completedCounter >= spec.LoopMaximum {
		return false, nil
	}
	if strings.TrimSpace(spec.ConditionExpr) == "" {
		return spec.LoopMaximum > 0 && completedCounter < spec.LoopMaximum, nil
	}
	loopCounter := completedCounter
	if spec.TestBefore {
		loopCounter = completedCounter + 1
		if spec.LoopMaximum > 0 && loopCounter > spec.LoopMaximum {
			return false, nil
		}
	}
	return spec.evalCondition(vars, loopCounter)
}

func (spec StandardLoopSpec) evalCondition(vars map[string]string, loopCounter int) (bool, error) {
	merged := map[string]string{}
	for k, v := range vars {
		merged[k] = v
	}
	b, _ := json.Marshal(loopCounter)
	merged["loopCounter"] = string(b)
	return expr.Eval(spec.ConditionExpr, merged)
}
