package deploy

import (
	"fmt"
	"strconv"

	"github.com/sparrow-community/sparrow/bpmn/element"
	"github.com/sparrow-community/sparrow/processing/expr"
)

// ComplexGatewaySpec holds compiled complex gateway facts.
type ComplexGatewaySpec struct {
	ID                   string
	ActivationCondition  string
	Default              string
}

func complexGatewaySpec(g element.ComplexGateway) ComplexGatewaySpec {
	return ComplexGatewaySpec{
		ID:                  g.ID,
		ActivationCondition: expressionText(g.ActivationCondition),
		Default:             g.Default,
	}
}

// ComplexGatewayOf returns the compiled complex gateway, if any.
func (d *Deployment) ComplexGatewayOf(id string) (ComplexGatewaySpec, bool) {
	if d == nil {
		return ComplexGatewaySpec{}, false
	}
	s, ok := d.complexGateways[id]
	return s, ok
}

// EvalComplexActivation evaluates activationCondition with arrived/incoming helpers.
// Empty condition means wait-for-all (caller treats false until arrived == incoming).
func (d *Deployment) EvalComplexActivation(gatewayID string, arrived, incoming int, vars map[string]string) (bool, error) {
	spec, ok := d.ComplexGatewayOf(gatewayID)
	if !ok {
		return false, fmt.Errorf("NOT_FOUND: complex gateway %q", gatewayID)
	}
	if spec.ActivationCondition == "" {
		return arrived >= incoming && incoming > 0, nil
	}
	env := make(map[string]string, len(vars)+3)
	for k, v := range vars {
		env[k] = v
	}
	env["arrived"] = strconv.Itoa(arrived)
	env["nrOfInstances"] = strconv.Itoa(arrived)
	env["incoming"] = strconv.Itoa(incoming)
	return expr.Eval(spec.ActivationCondition, env)
}

// ChooseComplexOutgoing selects outgoings like inclusive gateway (all matching / unconditional; else default).
func (d *Deployment) ChooseComplexOutgoing(gatewayID string, vars map[string]string) ([]string, error) {
	spec, ok := d.ComplexGatewayOf(gatewayID)
	if !ok {
		return nil, fmt.Errorf("NOT_FOUND: complex gateway %q", gatewayID)
	}
	var taken []string
	for _, flowID := range d.Outgoing(gatewayID) {
		if spec.Default != "" && flowID == spec.Default {
			continue
		}
		flow, err := d.SequenceFlow(flowID)
		if err != nil {
			return nil, err
		}
		text := ConditionText(flow)
		if text == "" {
			taken = append(taken, flowID)
			continue
		}
		match, err := expr.Eval(text, vars)
		if err != nil {
			return nil, fmt.Errorf("INVALID_CONDITION: flow %s: %w", flowID, err)
		}
		if match {
			taken = append(taken, flowID)
		}
	}
	if len(taken) == 0 {
		if spec.Default != "" {
			if _, err := d.SequenceFlow(spec.Default); err == nil {
				return []string{spec.Default}, nil
			}
		}
		return nil, fmt.Errorf("NO_OUTGOING_FLOW: complex gateway %q", gatewayID)
	}
	return taken, nil
}

func findComplexGatewayIn(fe *element.FlowElements, id string) *element.ComplexGateway {
	for i := range fe.ComplexGatewaies {
		if fe.ComplexGatewaies[i].ID == id {
			return &fe.ComplexGatewaies[i]
		}
	}
	for _, child := range childScopes(fe) {
		if g := findComplexGatewayIn(child, id); g != nil {
			return g
		}
	}
	return nil
}
