package processing

import (
	"context"
	"fmt"

	"github.com/sparrow-community/sparrow/processing/deploy"
	"github.com/sparrow-community/sparrow/processing/handlers"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func (x *Executor) runMultiInstanceStart(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	start *handlers.MultiInstanceStart,
	emit Emitter,
) ([]handlers.Publication, error) {
	if start == nil {
		return nil, nil
	}
	spec, ok := dep.MultiInstanceSpec(start.ElementID)
	if !ok {
		return nil, nil
	}
	if loop := inst.MultiInstanceLoops[start.ElementID]; loop != nil {
		initMIBehaviorState(loop, spec)
	}
	var pubs []handlers.Publication
	for _, idx := range start.InnerIndices {
		tid, err := NextID()
		if err != nil {
			return pubs, err
		}
		if extra := spec.CollectionElementVariables(inst.Variables, int(idx)); extra != nil {
			for k, v := range extra {
				inst.Variables[k] = v
			}
		}
		more, err := x.enterWithLoopIndex(ctx, dep, inst, tid, start.ElementID, idx, emit)
		pubs = append(pubs, more...)
		if err != nil {
			return pubs, err
		}
	}
	return pubs, nil
}

func (x *Executor) enterWithLoopIndex(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	tokenID, elementID string,
	loopIndex int32,
	emit Emitter,
) ([]handlers.Publication, error) {
	typ, err := dep.TypeOf(elementID)
	if err != nil {
		return nil, fmt.Errorf("UNSUPPORTED_ELEMENT: %v", err)
	}
	h, err := x.Handlers.Get(typ)
	if err != nil {
		return nil, err
	}
	effect, err := h.OnEnter(handlers.EnterInput{
		Deployment:        dep,
		Instance:          inst,
		ElementID:         elementID,
		Type:              typ,
		TokenID:           tokenID,
		Now:               x.now(),
		LoopInstanceIndex: loopIndex,
	})
	if err != nil {
		return nil, err
	}
	var pubs []handlers.Publication
	for _, rec := range effect.Records {
		if err := emit(rec); err != nil {
			return pubs, err
		}
	}
	if effect.Publish != nil {
		pubs = append(pubs, *effect.Publish)
	}
	if effect.EnterChild != "" {
		childTokenID := tokenID
		if effect.SpawnChildToken {
			childTokenID, err = spawnScopeChildToken(inst, tokenID)
			if err != nil {
				return pubs, err
			}
		}
		if typ == eventv1.Element_TYPE_SUB_PROCESS {
			if err := emitEventSubProcessStartArms(dep, inst, elementID, x.now(), emit); err != nil {
				return pubs, err
			}
		}
		more, err := x.Enter(ctx, dep, inst, childTokenID, effect.EnterChild, emit)
		pubs = append(pubs, more...)
		return pubs, err
	}
	if effect.Wait {
		return pubs, nil
	}
	return pubs, fmt.Errorf("UNSUPPORTED_ELEMENT: multi-instance inner enter on %q did not wait", elementID)
}

func (x *Executor) runMultiInstanceInnerComplete(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	elementID, innerTokenID string,
	innerIndex int32,
	emit Emitter,
) ([]handlers.Publication, error) {
	spec, ok := dep.MultiInstanceSpec(elementID)
	if !ok {
		return nil, nil
	}
	delete(inst.Tokens, innerTokenID)
	loop := inst.MultiInstanceLoops[elementID]
	if loop == nil {
		return nil, nil
	}
	if spec.OutputElementVar != "" {
		if raw, ok := inst.Variables[spec.OutputElementVar]; ok {
			loop.OutputItems = append(loop.OutputItems, raw)
		}
	}
	behaviorPubs, err := x.fireMIBehaviorEvents(spec, loop, inst.Variables, innerIndex)
	if err != nil {
		return nil, err
	}
	met, err := spec.CompletionMet(loop, inst.Variables, innerIndex)
	if err != nil {
		return behaviorPubs, err
	}
	if !met {
		if spec.Sequential && loop.NextIndex < loop.TotalInstances {
			tid, err := NextID()
			if err != nil {
				return behaviorPubs, err
			}
			if extra := spec.CollectionElementVariables(inst.Variables, int(loop.NextIndex)); extra != nil {
				for k, v := range extra {
					inst.Variables[k] = v
				}
			}
			more, err := x.enterWithLoopIndex(ctx, dep, inst, tid, elementID, loop.NextIndex, emit)
			return append(behaviorPubs, more...), err
		}
		return behaviorPubs, nil
	}
	more, err := x.completeMultiInstanceHost(ctx, dep, inst, spec, elementID, loop, emit)
	return append(behaviorPubs, more...), err
}

func (x *Executor) completeMultiInstanceHost(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	spec deploy.MultiInstanceSpec,
	elementID string,
	loop *projection.MultiInstanceLoop,
	emit Emitter,
) ([]handlers.Publication, error) {
	pubs, err := x.terminateMultiInstanceInners(dep, inst, elementID, loop.HostTokenID, emit)
	if err != nil {
		return pubs, err
	}
	AppendOutputCollection(inst, spec, loop)
	hostID := loop.HostTokenID
	typ, err := dep.TypeOf(elementID)
	if err != nil {
		return pubs, err
	}
	hostPayload := &eventv1.ActivityPayload{
		LoopInstanceIndex:      -1,
		LoopTotalInstances:     loop.TotalInstances,
		LoopCompletedInstances: loop.CompletedInstances,
	}
	for _, intent := range []eventv1.Element_Intent{
		eventv1.Element_INTENT_COMPLETING,
		eventv1.Element_INTENT_COMPLETED,
	} {
		el := &eventv1.Element{
			Intent:  intent,
			Type:    typ,
			Id:      elementID,
			TokenId: hostID,
			Payload: &eventv1.Element_ActivityPayload{ActivityPayload: hostPayload},
		}
		if err := emit(el); err != nil {
			return pubs, err
		}
	}
	for _, rec := range handlers.CancelAttachedBoundaries(dep, elementID, hostID) {
		if err := emit(rec); err != nil {
			return pubs, err
		}
	}
	for _, rec := range handlers.SubscribeCompensation(dep, elementID, hostID) {
		if err := emit(rec); err != nil {
			return pubs, err
		}
	}
	delete(inst.MultiInstanceLoops, elementID)
	next, err := x.takeOutgoing(dep, inst, hostID, elementID, "", emit)
	if err != nil {
		return pubs, err
	}
	more, err := x.Enter(ctx, dep, inst, hostID, next, emit)
	pubs = append(pubs, more...)
	return pubs, err
}

func (x *Executor) terminateMultiInstanceInners(
	dep *deploy.Deployment,
	inst *projection.Instance,
	elementID, hostTokenID string,
	emit Emitter,
) ([]handlers.Publication, error) {
	elTyp, err := dep.TypeOf(elementID)
	if err != nil {
		return nil, err
	}
	var pubs []handlers.Publication
	for tid, tok := range inst.Tokens {
		if tid == hostTokenID || tok == nil {
			continue
		}
		if tok.ElementID != elementID || tok.LoopInstanceIndex < 0 {
			continue
		}
		if tok.CalledProcessInstanceID != "" {
			pubs = append(pubs, handlers.Publication{
				Kind:            handlers.PublicationTerminateChild,
				ChildInstanceID: tok.CalledProcessInstanceID,
				CallActivityID:  elementID,
				HostTokenID:     tid,
			})
		}
		if tok.Status == projection.TokenWaiting {
			for _, intent := range []eventv1.Element_Intent{
				eventv1.Element_INTENT_TERMINATING,
				eventv1.Element_INTENT_TERMINATED,
			} {
				if err := emit(&eventv1.Element{
					Intent:  intent,
					Type:    elTyp,
					Id:      elementID,
					TokenId: tid,
					Payload: &eventv1.Element_ActivityPayload{
						ActivityPayload: &eventv1.ActivityPayload{LoopInstanceIndex: tok.LoopInstanceIndex},
					},
				}); err != nil {
					return pubs, err
				}
			}
		}
		delete(inst.Tokens, tid)
	}
	return pubs, nil
}

func (x *Executor) cancelMultiInstanceActivity(
	dep *deploy.Deployment,
	inst *projection.Instance,
	elementID string,
	emit Emitter,
) ([]handlers.Publication, error) {
	elTyp, err := dep.TypeOf(elementID)
	if err != nil {
		return nil, err
	}
	var pubs []handlers.Publication
	var innerHosts []string
	for tid, tok := range inst.Tokens {
		if tok == nil || tok.ElementID != elementID {
			continue
		}
		if tok.LoopInstanceIndex >= 0 && !tok.MultiInstanceHost {
			innerHosts = append(innerHosts, tid)
			if tok.CalledProcessInstanceID != "" {
				pubs = append(pubs, handlers.Publication{
					Kind:            handlers.PublicationTerminateChild,
					ChildInstanceID: tok.CalledProcessInstanceID,
					CallActivityID:  elementID,
					HostTokenID:     tid,
				})
			}
		}
	}
	terminateToken := func(tid string, idx int32) error {
		if inst.Tokens[tid] == nil {
			return nil
		}
		for _, intent := range []eventv1.Element_Intent{
			eventv1.Element_INTENT_TERMINATING,
			eventv1.Element_INTENT_TERMINATED,
		} {
			if err := emit(&eventv1.Element{
				Intent:  intent,
				Type:    elTyp,
				Id:      elementID,
				TokenId: tid,
				Payload: &eventv1.Element_ActivityPayload{
					ActivityPayload: &eventv1.ActivityPayload{LoopInstanceIndex: idx},
				},
			}); err != nil {
				return err
			}
		}
		delete(inst.Tokens, tid)
		return nil
	}
	for tid, tok := range inst.Tokens {
		if tok == nil {
			continue
		}
		for _, host := range innerHosts {
			if tid == host || tok.ScopeHostTokenID == host {
				if err := terminateToken(tid, tok.LoopInstanceIndex); err != nil {
					return pubs, err
				}
				break
			}
		}
	}
	for tid, tok := range inst.Tokens {
		if tok == nil || tok.ElementID != elementID {
			continue
		}
		if err := terminateToken(tid, tok.LoopInstanceIndex); err != nil {
			return pubs, err
		}
	}
	delete(inst.MultiInstanceLoops, elementID)
	inst.RemoveScopeBoundariesForScope(elementID)
	return pubs, nil
}

func initMIBehaviorState(loop *projection.MultiInstanceLoop, spec deploy.MultiInstanceSpec) {
	if loop == nil {
		return
	}
	loop.BehaviorInitialized = true
	loop.OneEventFired = false
	if len(spec.ComplexBehaviors) > 0 {
		loop.ComplexFired = make([]bool, len(spec.ComplexBehaviors))
	} else {
		loop.ComplexFired = nil
	}
}

func (x *Executor) fireMIBehaviorEvents(
	spec deploy.MultiInstanceSpec,
	loop *projection.MultiInstanceLoop,
	vars map[string]string,
	innerIndex int32,
) ([]handlers.Publication, error) {
	if loop == nil {
		return nil, nil
	}
	if err := seedMIBehaviorStateAfterRecover(spec, loop, vars, innerIndex); err != nil {
		return nil, err
	}
	var pubs []handlers.Publication
	appendThrow := func(th *deploy.MIBehaviorThrow) {
		if th == nil || th.Name == "" {
			return
		}
		switch th.Kind {
		case deploy.ThrowKindSignal:
			pubs = append(pubs, handlers.Publication{Kind: handlers.PublicationSignal, Name: th.Name})
		case deploy.ThrowKindMessage:
			pubs = append(pubs, handlers.Publication{Kind: handlers.PublicationMessage, Name: th.Name})
		}
	}

	if spec.NoneEvent != nil {
		appendThrow(spec.NoneEvent)
	}
	if spec.OneEvent != nil && !loop.OneEventFired && loop.CompletedInstances >= 1 {
		appendThrow(spec.OneEvent)
		loop.OneEventFired = true
	}

	for i, cb := range spec.ComplexBehaviors {
		if i < len(loop.ComplexFired) && loop.ComplexFired[i] {
			continue
		}
		met, err := spec.ComplexConditionMet(loop, vars, innerIndex, cb.Condition)
		if err != nil {
			return pubs, err
		}
		if !met {
			continue
		}
		th := cb.Throw
		appendThrow(&th)
		if len(loop.ComplexFired) < len(spec.ComplexBehaviors) {
			n := make([]bool, len(spec.ComplexBehaviors))
			copy(n, loop.ComplexFired)
			loop.ComplexFired = n
		}
		loop.ComplexFired[i] = true
	}
	return pubs, nil
}

func seedMIBehaviorStateAfterRecover(
	spec deploy.MultiInstanceSpec,
	loop *projection.MultiInstanceLoop,
	vars map[string]string,
	innerIndex int32,
) error {
	if loop.BehaviorInitialized {
		return nil
	}
	loop.BehaviorInitialized = true
	if len(spec.ComplexBehaviors) > 0 && loop.ComplexFired == nil {
		loop.ComplexFired = make([]bool, len(spec.ComplexBehaviors))
	}
	// After Recover with prior completions, avoid re-publishing one/complex milestones.
	if loop.CompletedInstances > 1 {
		loop.OneEventFired = true
		saved := loop.CompletedInstances
		loop.CompletedInstances = saved - 1
		for i, cb := range spec.ComplexBehaviors {
			met, err := spec.ComplexConditionMet(loop, vars, innerIndex, cb.Condition)
			if err != nil {
				loop.CompletedInstances = saved
				return err
			}
			if met {
				loop.ComplexFired[i] = true
			}
		}
		loop.CompletedInstances = saved
	}
	return nil
}
