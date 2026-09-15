package handlers

import (
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func standardLoopSkipEnter(in EnterInput) (*Effect, error) {
	records := InstantLifecycle(in.Type, in.ElementID, in.TokenID, nil)
	return &Effect{Records: records, TakeOutgoing: true}, nil
}

// maybeStandardLoopEnter returns (effect, err, handled). When handled and effect is nil,
// the caller should proceed with a normal first-iteration enter.
func maybeStandardLoopEnter(in EnterInput) (*Effect, error, bool) {
	if in.Deployment == nil {
		return nil, nil, false
	}
	spec, ok := in.Deployment.StandardLoopSpec(in.ElementID)
	if !ok {
		return nil, nil, false
	}
	if in.Instance != nil {
		if tok := in.Instance.Tokens[in.TokenID]; tok != nil && tok.StandardLoopIteration > 0 && tok.ElementID == in.ElementID {
			return nil, nil, false
		}
	}
	vars := map[string]string{}
	if in.Instance != nil {
		vars = in.Instance.Variables
	}
	enter, err := spec.ShouldEnterFirst(vars)
	if err != nil {
		return nil, err, true
	}
	if !enter {
		eff, err := standardLoopSkipEnter(in)
		return eff, err, true
	}
	return nil, nil, false
}

func standardLoopComplete(in CompleteInput, records []*eventv1.Element) (*Effect, error) {
	if in.Deployment == nil {
		return nil, nil
	}
	spec, ok := in.Deployment.StandardLoopSpec(in.ElementID)
	if !ok {
		return nil, nil
	}
	counter := 0
	if in.Token != nil {
		counter = in.Token.StandardLoopIteration
	}
	vars := map[string]string{}
	if in.Instance != nil {
		for k, v := range in.Instance.Variables {
			vars[k] = v
		}
	}
	for _, v := range in.Variables {
		if v == nil || v.GetName() == "" {
			continue
		}
		vars[v.GetName()] = v.GetJsonValue()
	}
	cont, err := spec.ShouldContinueAfter(vars, counter)
	if err != nil {
		return nil, err
	}
	records = append(records, cancelAttachedBoundary(in.Deployment, in.ElementID, in.TokenID)...)
	if cont {
		return &Effect{Records: records, ReEnter: true}, nil
	}
	if in.Deployment.IsCompensationHandler(in.ElementID) {
		return &Effect{
			Records:             records,
			DiscardToken:        true,
			AdvanceCompensation: true,
		}, nil
	}
	records = append(records, subscribeCompensation(in.Deployment, in.ElementID, in.TokenID)...)
	return &Effect{Records: records, TakeOutgoing: true}, nil
}
