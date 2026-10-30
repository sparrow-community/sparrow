package processing

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/sparrow-community/sparrow/processing/deploy"
	"github.com/sparrow-community/sparrow/processing/handlers"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// Emitter appends an Element behavior record as an EVENT and updates projection.
type Emitter func(el *eventv1.Element) error

// Executor drives token movement using ElementHandler semantics over bpmn definitions.
type Executor struct {
	Handlers *handlers.Registry
	Now      func() time.Time

	// intervention is the Engine control plane (nil-safe / default-off).
	intervention *interventionCtl
	// activeSourceCmdID / activeInstanceID are set by Engine for the in-flight COMMAND.
	activeSourceCmdID string
	activeInstanceID  string
}

func NewExecutor(reg *handlers.Registry) *Executor {
	if reg == nil {
		reg = handlers.DefaultRegistry()
	}
	return &Executor{Handlers: reg, intervention: &interventionCtl{}}
}

func (x *Executor) now() time.Time {
	if x != nil && x.Now != nil {
		return x.Now()
	}
	return time.Now()
}

func spawnScopeChildToken(inst *projection.Instance, hostTokenID string) (string, error) {
	childTokenID, err := NextID()
	if err != nil {
		return "", err
	}
	child := inst.Tokens[childTokenID]
	if child == nil {
		child = &projection.Token{ID: childTokenID, LoopInstanceIndex: -1}
		inst.Tokens[childTokenID] = child
	}
	child.ScopeHostTokenID = hostTokenID
	return childTokenID, nil
}

func (x *Executor) Enter(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	tokenID, elementID string,
	emit Emitter,
) ([]handlers.Publication, error) {
	var pubs []handlers.Publication
	for {
		typ, err := dep.TypeOf(elementID)
		if err != nil {
			return pubs, fmt.Errorf("UNSUPPORTED_ELEMENT: %v", err)
		}

		h, err := x.Handlers.Get(typ)
		if err != nil {
			return pubs, err
		}
		effect, err := h.OnEnter(handlers.EnterInput{
			Deployment:        dep,
			Instance:          inst,
			ElementID:         elementID,
			Type:              typ,
			TokenID:           tokenID,
			Now:               x.now(),
			LoopInstanceIndex: -1,
		})
		if err != nil {
			return pubs, err
		}
		for _, rec := range effect.Records {
			if err := emit(rec); err != nil {
				return pubs, err
			}
		}
		if effect.Publish != nil {
			pubs = append(pubs, *effect.Publish)
		}
		if effect.TerminateScope {
			more, err := x.terminateEnclosingScope(dep, inst, tokenID, elementID, emit)
			pubs = append(pubs, more...)
			if err != nil {
				return pubs, err
			}
		}
		if effect.MultiInstanceInnerComplete {
			innerIdx := int32(-1)
			if tok := inst.Tokens[tokenID]; tok != nil {
				innerIdx = tok.LoopInstanceIndex
			}
			more, err := x.runMultiInstanceInnerComplete(ctx, dep, inst, elementID, tokenID, innerIdx, emit)
			pubs = append(pubs, more...)
			return pubs, err
		}
		if effect.MultiInstanceStart != nil {
			more, err := x.runMultiInstanceStart(ctx, dep, inst, effect.MultiInstanceStart, emit)
			pubs = append(pubs, more...)
			if err != nil {
				return pubs, err
			}
		}
		if effect.AdHocEnable != nil {
			more, err := x.runAdHocEnable(ctx, dep, inst, tokenID, effect.AdHocEnable, emit)
			pubs = append(pubs, more...)
			return pubs, err
		}
		if effect.DecideExclusive || effect.DecideInclusive || effect.DecideComplex {
			mode := DecideExclusive
			if effect.DecideInclusive {
				mode = DecideInclusive
			} else if effect.DecideComplex {
				mode = DecideComplex
			}
			pending := &PendingTransition{
				Kind:             PendingDecide,
				FromElementID:    elementID,
				TokenID:          tokenID,
				DecideMode:       mode,
				TerminateJoinPeers: effect.DecideTerminateJoinPeers,
			}
			if x.tryBarrier(inst, elementID, tokenID, pending) {
				return pubs, nil
			}
			more, err := x.finalizeGatewayDecide(ctx, dep, inst, tokenID, elementID, mode, effect.DecideTerminateJoinPeers, emit)
			pubs = append(pubs, more...)
			return pubs, err
		}
		if effect.EnterChild != "" {
			pending := &PendingTransition{
				Kind:            PendingEnterChild,
				FromElementID:   elementID,
				TokenID:         tokenID,
				EnterChildID:    effect.EnterChild,
				SpawnChildToken: effect.SpawnChildToken,
				NextElementIDs:  []string{effect.EnterChild},
			}
			if x.tryBarrier(inst, elementID, tokenID, pending) {
				return pubs, nil
			}
			childTokenID := tokenID
			if effect.SpawnChildToken {
				var err error
				childTokenID, err = spawnScopeChildToken(inst, tokenID)
				if err != nil {
					return pubs, err
				}
			}
			if typ == eventv1.Element_TYPE_SUB_PROCESS || typ == eventv1.Element_TYPE_TRANSACTION {
				if err := emitEventSubProcessStartArms(dep, inst, elementID, x.now(), emit); err != nil {
					return pubs, err
				}
			}
			elementID = effect.EnterChild
			tokenID = childTokenID
			continue
		}
		if effect.TerminateJoinPeers != "" {
			if err := x.terminateJoinPeers(dep, inst, tokenID, effect.TerminateJoinPeers, emit); err != nil {
				return pubs, err
			}
		}
		if len(effect.Fork) > 0 {
			pending, err := pendingFork(dep, tokenID, elementID, effect.Fork)
			if err != nil {
				return pubs, err
			}
			if x.tryBarrier(inst, elementID, tokenID, pending) {
				return pubs, nil
			}
			more, err := x.forkOutgoings(ctx, dep, inst, tokenID, elementID, effect.Fork, emit)
			pubs = append(pubs, more...)
			return pubs, err
		}
		if len(effect.LinkContinue) > 0 {
			pending := &PendingTransition{
				Kind:           PendingLink,
				FromElementID:  elementID,
				TokenID:        tokenID,
				LinkCatchIDs:   append([]string(nil), effect.LinkContinue...),
				NextElementIDs: append([]string(nil), effect.LinkContinue...),
			}
			if x.tryBarrier(inst, elementID, tokenID, pending) {
				return pubs, nil
			}
			for i, catchID := range effect.LinkContinue {
				tid := tokenID
				if i > 0 {
					var err error
					tid, err = NextID()
					if err != nil {
						return pubs, err
					}
				}
				more, err := x.Enter(ctx, dep, inst, tid, catchID, emit)
				pubs = append(pubs, more...)
				if err != nil {
					return pubs, err
				}
			}
			return pubs, nil
		}
		if effect.Wait {
			if effect.TriggerCompensation {
				more, err := x.startCompensation(ctx, dep, inst, tokenID, elementID, emit)
				pubs = append(pubs, more...)
				return pubs, err
			}
			if effect.TriggerTransactionCancel {
				more, err := x.startTransactionCancel(ctx, dep, inst, tokenID, elementID, emit)
				pubs = append(pubs, more...)
				return pubs, err
			}
			// K1: optional breakpoint hit when a wait settles (ACTIVATED).
			pending := &PendingTransition{
				Kind:          PendingWait,
				FromElementID: elementID,
				TokenID:       tokenID,
			}
			if x.tryBarrierWait(inst, elementID, tokenID, pending) {
				return pubs, nil
			}
			return pubs, nil
		}
		if effect.ThrowError != nil {
			throwTyp, err := dep.TypeOf(elementID)
			if err != nil {
				return pubs, fmt.Errorf("UNSUPPORTED_ELEMENT: %v", err)
			}
			if err := emit(&eventv1.Element{
				Intent:  eventv1.Element_INTENT_ERROR_THROWN,
				Type:    throwTyp,
				Id:      elementID,
				TokenId: tokenID,
				Payload: &eventv1.Element_EventPayload{
					EventPayload: &eventv1.EventPayload{ErrorCode: effect.ThrowError.ErrorCode},
				},
			}); err != nil {
				return pubs, err
			}
			more, err := x.propagateError(ctx, dep, inst, elementID, tokenID, effect.ThrowError.ErrorCode, emit)
			pubs = append(pubs, more...)
			return pubs, err
		}
		if effect.ThrowEscalation != nil {
			throwTyp, err := dep.TypeOf(elementID)
			if err != nil {
				return pubs, fmt.Errorf("UNSUPPORTED_ELEMENT: %v", err)
			}
			if err := emit(&eventv1.Element{
				Intent:  eventv1.Element_INTENT_ESCALATION_THROWN,
				Type:    throwTyp,
				Id:      elementID,
				TokenId: tokenID,
				Payload: &eventv1.Element_EventPayload{
					EventPayload: &eventv1.EventPayload{EscalationCode: effect.ThrowEscalation.EscalationCode},
				},
			}); err != nil {
				return pubs, err
			}
			more, err := x.propagateEscalation(ctx, dep, inst, elementID, tokenID, effect.ThrowEscalation.EscalationCode, emit)
			pubs = append(pubs, more...)
			if err != nil {
				return pubs, err
			}
			if inst.Tokens[tokenID] == nil {
				return pubs, nil
			}
			if effect.TryCompleteProcess {
				more, err := x.tryCompleteScope(ctx, dep, inst, tokenID, elementID, emit)
				pubs = append(pubs, more...)
				return pubs, err
			}
			if !effect.TakeOutgoing {
				return pubs, nil
			}
			// Fall through to TakeOutgoing for intermediate escalation throw.
		} else if effect.TryCompleteProcess {
			more, err := x.tryCompleteScope(ctx, dep, inst, tokenID, elementID, emit)
			pubs = append(pubs, more...)
			return pubs, err
		}
		if !effect.TakeOutgoing {
			return pubs, nil
		}

		more, err := x.leaveViaOutgoings(ctx, dep, inst, tokenID, elementID, effect.OutgoingFlowID, emit)
		pubs = append(pubs, more...)
		return pubs, err
	}
}

func (x *Executor) Complete(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	tokenID, elementID string,
	vars []*eventv1.Variable,
	emit Emitter,
) ([]handlers.Publication, error) {
	typ, err := dep.TypeOf(elementID)
	if err != nil {
		return nil, err
	}
	h, err := x.Handlers.Get(typ)
	if err != nil {
		return nil, err
	}
	effect, err := h.OnComplete(handlers.CompleteInput{
		Deployment: dep,
		Instance:   inst,
		Token:      inst.Tokens[tokenID],
		ElementID:  elementID,
		Type:       typ,
		TokenID:    tokenID,
		Variables:  vars,
	})
	if err != nil {
		return nil, err
	}
	innerIdx := int32(-1)
	if tok := inst.Tokens[tokenID]; tok != nil {
		innerIdx = tok.LoopInstanceIndex
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
	if effect.MultiInstanceInnerComplete {
		more, err := x.runMultiInstanceInnerComplete(ctx, dep, inst, elementID, tokenID, innerIdx, emit)
		pubs = append(pubs, more...)
		return pubs, err
	}
	if effect.AdHocInnerComplete != "" {
		more, err := x.advanceAdHoc(ctx, dep, inst, effect.AdHocInnerComplete, tokenID, elementID, emit)
		pubs = append(pubs, more...)
		return pubs, err
	}
	if effect.MultiInstanceCancel != "" {
		more, err := x.cancelMultiInstanceActivity(dep, inst, effect.MultiInstanceCancel, emit)
		pubs = append(pubs, more...)
		if err != nil {
			return pubs, err
		}
	}
	if effect.ReEnter {
		more, err := x.Enter(ctx, dep, inst, tokenID, elementID, emit)
		pubs = append(pubs, more...)
		return pubs, err
	}
	if len(effect.TerminateWaitingAt) > 0 {
		if err := x.terminateWaitingAt(inst, tokenID, effect.TerminateWaitingAt, emit); err != nil {
			return pubs, err
		}
	}
	if effect.SpawnOutgoing != nil {
		spawn := effect.SpawnOutgoing
		spawnID, err := NextID()
		if err != nil {
			return pubs, err
		}
		for _, rec := range handlers.InstantLifecycle(spawn.Type, spawn.ElementID, spawnID, nil) {
			if err := emit(rec); err != nil {
				return pubs, err
			}
		}
		next, err := x.takeOutgoing(dep, inst, spawnID, spawn.ElementID, spawn.OutgoingFlowID, emit)
		if err != nil {
			return pubs, err
		}
		more, err := x.Enter(ctx, dep, inst, spawnID, next, emit)
		pubs = append(pubs, more...)
		if err != nil {
			return pubs, err
		}
		more, err = x.tryCompleteProcessScope(dep, inst, emit)
		pubs = append(pubs, more...)
		return pubs, err
	}
	if effect.DiscardToken {
		delete(inst.Tokens, tokenID)
	}
	if effect.AdvanceCompensation {
		if inst.PendingCompensation == nil && dep.IsCompensationHandler(elementID) && instHasParentCall(inst) {
			inst.PendingCompensation = &projection.PendingCompensation{
				NotifyParentUnfinishedCall: true,
			}
		}
		more, err := x.advanceCompensation(ctx, dep, inst, emit)
		pubs = append(pubs, more...)
		return pubs, err
	}
	if effect.TryCompleteProcess {
		more, err := x.tryCompleteScope(ctx, dep, inst, tokenID, elementID, emit)
		pubs = append(pubs, more...)
		return pubs, err
	}
	if effect.Wait || !effect.TakeOutgoing {
		return pubs, nil
	}

	more, err := x.leaveViaOutgoings(ctx, dep, inst, tokenID, elementID, effect.OutgoingFlowID, emit)
	pubs = append(pubs, more...)
	return pubs, err
}

func (x *Executor) terminateJoinPeers(dep *deploy.Deployment, inst *projection.Instance, survivorTokenID, joinElementID string, emit Emitter) error {
	typ := eventv1.Element_TYPE_PARALLEL_GATEWAY
	if dep != nil {
		if t, err := dep.TypeOf(joinElementID); err == nil {
			typ = t
		}
	}
	for tid, tok := range inst.Tokens {
		if tid == survivorTokenID || tok == nil {
			continue
		}
		if tok.ElementID != joinElementID || tok.Status != projection.TokenWaiting {
			continue
		}
		for _, intent := range []eventv1.Element_Intent{
			eventv1.Element_INTENT_TERMINATING,
			eventv1.Element_INTENT_TERMINATED,
		} {
			if err := emit(&eventv1.Element{
				Intent:  intent,
				Type:    typ,
				Id:      joinElementID,
				TokenId: tid,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (x *Executor) terminateWaitingAt(inst *projection.Instance, survivorTokenID string, elementIDs []string, emit Emitter) error {
	want := make(map[string]bool, len(elementIDs))
	for _, id := range elementIDs {
		want[id] = true
	}
	for tid, tok := range inst.Tokens {
		if tid == survivorTokenID || tok == nil {
			continue
		}
		if !want[tok.ElementID] || tok.Status != projection.TokenWaiting {
			continue
		}
		typ := eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT
		for _, intent := range []eventv1.Element_Intent{
			eventv1.Element_INTENT_TERMINATING,
			eventv1.Element_INTENT_TERMINATED,
		} {
			if err := emit(&eventv1.Element{
				Intent:  intent,
				Type:    typ,
				Id:      tok.ElementID,
				TokenId: tid,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (x *Executor) takeOutgoing(
	dep *deploy.Deployment,
	inst *projection.Instance,
	tokenID, fromElementID, preferredFlowID string,
	emit Emitter,
) (string, error) {
	flowID := preferredFlowID
	if flowID == "" {
		var vars map[string]string
		if inst != nil {
			vars = inst.Variables
		}
		var err error
		flowID, err = dep.ChooseConditionalOutgoing(fromElementID, vars)
		if err != nil {
			return "", err
		}
	}
	flow, err := dep.SequenceFlow(flowID)
	if err != nil {
		return "", err
	}
	if err := emit(handlers.SequenceFlowTaken(flow.ID, flow.SourceRef, flow.TargetRef, tokenID)); err != nil {
		return "", err
	}
	return flow.TargetRef, nil
}

// leaveViaOutgoings takes preferred or chosen outgoings and Enter each target.
// Multiple unconditional flows fork tokens like a parallel split.
func (x *Executor) leaveViaOutgoings(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	tokenID, fromElementID, preferredFlowID string,
	emit Emitter,
) ([]handlers.Publication, error) {
	if preferredFlowID != "" {
		pending, err := pendingLeave(dep, tokenID, fromElementID, preferredFlowID)
		if err != nil {
			return nil, err
		}
		if x.tryBarrier(inst, fromElementID, tokenID, pending) {
			return nil, nil
		}
		return x.takeOutgoingThenEnter(ctx, dep, inst, tokenID, fromElementID, preferredFlowID, emit)
	}
	var vars map[string]string
	if inst != nil {
		vars = inst.Variables
	}
	flows, err := dep.ChooseOutgoingFlows(fromElementID, vars)
	if err != nil {
		return nil, err
	}
	if len(flows) > 1 {
		pending, err := pendingFork(dep, tokenID, fromElementID, flows)
		if err != nil {
			return nil, err
		}
		if x.tryBarrier(inst, fromElementID, tokenID, pending) {
			return nil, nil
		}
		return x.forkOutgoings(ctx, dep, inst, tokenID, fromElementID, flows, emit)
	}
	pending, err := pendingLeave(dep, tokenID, fromElementID, flows[0])
	if err != nil {
		return nil, err
	}
	if x.tryBarrier(inst, fromElementID, tokenID, pending) {
		return nil, nil
	}
	return x.takeOutgoingThenEnter(ctx, dep, inst, tokenID, fromElementID, flows[0], emit)
}

// takeOutgoingThenEnter emits SEQUENCE_FLOW_TAKEN, may barrier on the edge (K2),
// then Enter the target.
func (x *Executor) takeOutgoingThenEnter(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	tokenID, fromElementID, flowID string,
	emit Emitter,
) ([]handlers.Publication, error) {
	next, err := x.takeOutgoing(dep, inst, tokenID, fromElementID, flowID, emit)
	if err != nil {
		return nil, err
	}
	edge := &PendingTransition{
		Kind:           PendingEdge,
		FromElementID:  fromElementID,
		TokenID:        tokenID,
		TakenFlowIDs:   []string{flowID},
		NextElementIDs: []string{next},
		OutgoingFlowID: flowID,
	}
	// Pause element for BP matching is the sequence flow id (edge stop).
	if x.tryBarrier(inst, flowID, tokenID, edge) {
		return nil, nil
	}
	return x.Enter(ctx, dep, inst, tokenID, next, emit)
}

func (x *Executor) forkOutgoings(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	tokenID, elementID string,
	flows []string,
	emit Emitter,
) ([]handlers.Publication, error) {
	if x.intervention != nil {
		x.intervention.mu.Lock()
		x.intervention.forkSuppressDepth++
		x.intervention.mu.Unlock()
		defer func() {
			x.intervention.mu.Lock()
			x.intervention.forkSuppressDepth--
			x.intervention.mu.Unlock()
		}()
	}
	type forkBranch struct {
		tid  string
		next string
	}
	var rest, terms []forkBranch
	var pubs []handlers.Publication
	for i, flowID := range flows {
		tid := tokenID
		if i > 0 {
			var err error
			tid, err = NextID()
			if err != nil {
				return pubs, err
			}
		}
		next, err := x.takeOutgoing(dep, inst, tid, elementID, flowID, emit)
		if err != nil {
			return pubs, err
		}
		b := forkBranch{tid: tid, next: next}
		if dep.IsTerminateEnd(next) {
			terms = append(terms, b)
		} else {
			rest = append(rest, b)
		}
	}
	for _, b := range append(rest, terms...) {
		more, err := x.Enter(ctx, dep, inst, b.tid, b.next, emit)
		pubs = append(pubs, more...)
		if err != nil {
			return pubs, err
		}
	}
	return pubs, nil
}

func (x *Executor) tryBarrier(inst *projection.Instance, elementID, tokenID string, pending *PendingTransition) bool {
	if x == nil || x.intervention == nil || pending == nil {
		return false
	}
	instanceID := x.activeInstanceID
	if instanceID == "" && inst != nil {
		instanceID = inst.ID
	}
	return x.intervention.shouldBarrier(instanceID, elementID, tokenID, x.activeSourceCmdID, pending)
}

// tryBarrierWait arms only on breakpoint hit at wait ACTIVATED (not Step policy).
func (x *Executor) tryBarrierWait(inst *projection.Instance, elementID, tokenID string, pending *PendingTransition) bool {
	if x == nil || x.intervention == nil || pending == nil {
		return false
	}
	instanceID := x.activeInstanceID
	if instanceID == "" && inst != nil {
		instanceID = inst.ID
	}
	return x.intervention.shouldBarrierWait(instanceID, elementID, tokenID, x.activeSourceCmdID, pending)
}

func (x *Executor) resumePending(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	pending *PendingTransition,
	emit Emitter,
) ([]handlers.Publication, error) {
	if pending == nil {
		return nil, fmt.Errorf("INVALID_STATE: no pending transition")
	}
	switch pending.Kind {
	case PendingDecide:
		mode := pending.DecideMode
		if mode == "" {
			mode = DecideExclusive
		}
		return x.finalizeGatewayDecide(ctx, dep, inst, pending.TokenID, pending.FromElementID, mode, pending.TerminateJoinPeers, emit)
	case PendingLeave:
		flowID := pending.OutgoingFlowID
		if flowID == "" && len(pending.TakenFlowIDs) > 0 {
			flowID = pending.TakenFlowIDs[0]
		}
		return x.leaveViaOutgoings(ctx, dep, inst, pending.TokenID, pending.FromElementID, flowID, emit)
	case PendingEdge:
		next := ""
		if len(pending.NextElementIDs) > 0 {
			next = pending.NextElementIDs[0]
		}
		if next == "" {
			return nil, fmt.Errorf("INVALID_STATE: edge pending missing next element")
		}
		return x.Enter(ctx, dep, inst, pending.TokenID, next, emit)
	case PendingWait:
		// Wait BP / manual pause: resume clears pause only; token stays waiting.
		return nil, nil
	case PendingFork:
		return x.forkOutgoings(ctx, dep, inst, pending.TokenID, pending.FromElementID, pending.TakenFlowIDs, emit)
	case PendingEnterChild:
		childTokenID := pending.TokenID
		if pending.SpawnChildToken {
			var err error
			childTokenID, err = spawnScopeChildToken(inst, pending.TokenID)
			if err != nil {
				return nil, err
			}
		}
		typ, err := dep.TypeOf(pending.FromElementID)
		if err == nil && (typ == eventv1.Element_TYPE_SUB_PROCESS || typ == eventv1.Element_TYPE_TRANSACTION) {
			if err := emitEventSubProcessStartArms(dep, inst, pending.FromElementID, x.now(), emit); err != nil {
				return nil, err
			}
		}
		return x.Enter(ctx, dep, inst, childTokenID, pending.EnterChildID, emit)
	case PendingLink:
		var pubs []handlers.Publication
		for i, catchID := range pending.LinkCatchIDs {
			tid := pending.TokenID
			if i > 0 {
				var err error
				tid, err = NextID()
				if err != nil {
					return pubs, err
				}
			}
			more, err := x.Enter(ctx, dep, inst, tid, catchID, emit)
			pubs = append(pubs, more...)
			if err != nil {
				return pubs, err
			}
		}
		return pubs, nil
	default:
		return nil, fmt.Errorf("INVALID_STATE: unknown pending kind %q", pending.Kind)
	}
}

// finalizeGatewayDecide chooses outgoings with current variables, emits
// COMPLETING/COMPLETED, then leaves or forks. Leave barriers are suppressed in
// this burst so StepInto from PendingDecide does not double-stop on the gateway;
// sequence-flow-edge barriers (K2) still apply after SEQUENCE_FLOW_TAKEN.
func (x *Executor) finalizeGatewayDecide(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	tokenID, elementID string,
	mode DecideMode,
	terminateJoinPeers string,
	emit Emitter,
) ([]handlers.Publication, error) {
	var vars map[string]string
	if inst != nil {
		vars = inst.Variables
	}
	var flows []string
	var err error
	switch mode {
	case DecideInclusive:
		flows, err = dep.ChooseInclusiveOutgoing(elementID, vars)
	case DecideComplex:
		flows, err = dep.ChooseComplexOutgoing(elementID, vars)
	default:
		var flowID string
		flowID, err = dep.ChooseExclusiveOutgoing(elementID, vars)
		if err == nil {
			flows = []string{flowID}
		}
	}
	if err != nil {
		return nil, err
	}
	if len(flows) == 0 {
		return nil, fmt.Errorf("NO_OUTGOING_FLOW: gateway %q", elementID)
	}
	typ := eventv1.Element_TYPE_EXCLUSIVE_GATEWAY
	if t, typeErr := dep.TypeOf(elementID); typeErr == nil {
		typ = t
	}
	for _, intent := range []eventv1.Element_Intent{
		eventv1.Element_INTENT_COMPLETING,
		eventv1.Element_INTENT_COMPLETED,
	} {
		el := &eventv1.Element{
			Intent:  intent,
			Type:    typ,
			Id:      elementID,
			TokenId: tokenID,
		}
		if intent == eventv1.Element_INTENT_COMPLETED && mode == DecideExclusive {
			el.Payload = &eventv1.Element_GatewayPayload{
				GatewayPayload: &eventv1.GatewayPayload{TakenSequenceFlowId: flows[0]},
			}
		}
		if err := emit(el); err != nil {
			return nil, err
		}
	}
	if terminateJoinPeers != "" {
		if err := x.terminateJoinPeers(dep, inst, tokenID, terminateJoinPeers, emit); err != nil {
			return nil, err
		}
	}
	if x.intervention != nil {
		x.intervention.mu.Lock()
		x.intervention.decideFinalizeDepth++
		x.intervention.mu.Unlock()
		defer func() {
			x.intervention.mu.Lock()
			x.intervention.decideFinalizeDepth--
			x.intervention.mu.Unlock()
		}()
	}
	if len(flows) == 1 {
		return x.leaveViaOutgoings(ctx, dep, inst, tokenID, elementID, flows[0], emit)
	}
	return x.forkOutgoings(ctx, dep, inst, tokenID, elementID, flows, emit)
}

// tryCompleteScope checks if the scope containing elementID can be completed.
// If elementID is inside a SubProcess, tries to complete that SubProcess.
// If at process level, tries to complete the Process.
func (x *Executor) tryCompleteScope(ctx context.Context, dep *deploy.Deployment, inst *projection.Instance, tokenID, elementID string, emit Emitter) ([]handlers.Publication, error) {
	scopeID, _ := dep.ScopeOf(elementID)
	rootProcessID := inst.ProcessID
	if rootProcessID == "" {
		rootProcessID = dep.ProcessID()
	}
	if scopeID == "" || scopeID == rootProcessID {
		return x.tryCompleteProcessScope(dep, inst, emit)
	}
	_, miScope := dep.MultiInstanceSpec(scopeID)
	iterationHostID := ""
	if trigger := inst.Tokens[tokenID]; trigger != nil {
		iterationHostID = trigger.ScopeHostTokenID
	}
	// SubProcess scope: ignore parked host; require children at EndEvents.
	hostTokenID := ""
	for tid, tok := range inst.Tokens {
		if tok == nil {
			continue
		}
		if tok.ElementID == scopeID {
			if miScope && iterationHostID != "" && tid != iterationHostID {
				continue
			}
			hostTokenID = tid
			continue
		}
		if miScope && iterationHostID != "" && tok.ScopeHostTokenID != iterationHostID {
			continue
		}
		tokScope, _ := dep.ScopeOf(tok.ElementID)
		inScope := tokScope == scopeID || isInScope(dep, tokScope, scopeID)
		if !inScope {
			continue
		}
		if tok.Status == projection.TokenWaiting {
			return nil, nil
		}
		typ, err := dep.TypeOf(tok.ElementID)
		if err != nil || typ != eventv1.Element_TYPE_END_EVENT {
			return nil, nil
		}
	}
	completeTokenID := hostTokenID
	if completeTokenID == "" {
		// Legacy inline without host: completing token is the end-event token.
		completeTokenID = tokenID
	}
	// Drop finished child tokens; host (if any) continues via OnComplete.
	// Emit TERMINATED so Recover replay does not revive EndEvent children.
	for tid, tok := range inst.Tokens {
		if tok == nil || tid == completeTokenID {
			continue
		}
		if tok.ElementID == scopeID {
			continue
		}
		if miScope && iterationHostID != "" && tok.ScopeHostTokenID != iterationHostID {
			continue
		}
		tokScope, _ := dep.ScopeOf(tok.ElementID)
		if tokScope == scopeID || isInScope(dep, tokScope, scopeID) {
			typ, err := dep.TypeOf(tok.ElementID)
			if err != nil {
				return nil, err
			}
			for _, intent := range []eventv1.Element_Intent{
				eventv1.Element_INTENT_TERMINATING,
				eventv1.Element_INTENT_TERMINATED,
			} {
				if err := emit(&eventv1.Element{
					Intent:  intent,
					Type:    typ,
					Id:      tok.ElementID,
					TokenId: tid,
				}); err != nil {
					return nil, err
				}
			}
		}
	}
	// All tokens in this scope are at EndEvents; disarm scope boundaries and complete
	inst.RemoveScopeBoundariesForScope(scopeID)
	if err := emitEventSubProcessStartDisarmInScope(dep, scopeID, inst, emit); err != nil {
		return nil, err
	}

	completeID := scopeID
	completeType := eventv1.Element_TYPE_SUB_PROCESS
	if t, err := dep.TypeOf(scopeID); err == nil {
		completeType = t
	}
	outgoingFrom := scopeID

	h, err := x.Handlers.Get(completeType)
	if err != nil {
		return nil, err
	}
	effect, err := h.OnComplete(handlers.CompleteInput{
		Deployment: dep,
		Instance:   inst,
		ElementID:  completeID,
		Type:       completeType,
		TokenID:    completeTokenID,
		Token:      inst.Tokens[completeTokenID],
	})
	if err != nil {
		return nil, err
	}
	innerIdx := int32(-1)
	if tok := inst.Tokens[completeTokenID]; tok != nil {
		innerIdx = tok.LoopInstanceIndex
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
	if effect.MultiInstanceInnerComplete {
		more, err := x.runMultiInstanceInnerComplete(ctx, dep, inst, completeID, completeTokenID, innerIdx, emit)
		pubs = append(pubs, more...)
		return pubs, err
	}
	if effect.DiscardToken {
		delete(inst.Tokens, completeTokenID)
	}
	if dep.IsEventSubProcess(scopeID) {
		if spec, ok := dep.EventSubProcessSpec(scopeID); ok && !spec.Interrupting && spec.Kind != deploy.CatchKindCompensate {
			if err := emit(eventSubProcessStartActivated(dep, spec, x.now())); err != nil {
				return pubs, err
			}
		}
	}
	if effect.AdvanceCompensation {
		more, err := x.advanceCompensation(ctx, dep, inst, emit)
		pubs = append(pubs, more...)
		return pubs, err
	}
	if effect.TryCompleteProcess {
		more, err := x.tryCompleteProcessScope(dep, inst, emit)
		pubs = append(pubs, more...)
		return pubs, err
	}
	if effect.TakeOutgoing {
		more, err := x.leaveViaOutgoings(ctx, dep, inst, completeTokenID, outgoingFrom, effect.OutgoingFlowID, emit)
		pubs = append(pubs, more...)
		return pubs, err
	}
	return pubs, nil
}

// tryCompleteProcessScope checks if all tokens are at process-level EndEvents.
func (x *Executor) tryCompleteProcessScope(dep *deploy.Deployment, inst *projection.Instance, emit Emitter) ([]handlers.Publication, error) {
	processID := inst.ProcessID
	if processID == "" {
		processID = dep.ProcessID()
	}
	for _, tok := range inst.Tokens {
		if tok.Status == projection.TokenWaiting {
			return nil, nil
		}
		tokScope, _ := dep.ScopeOf(tok.ElementID)
		if tokScope != processID {
			return nil, nil
		}
		typ, err := dep.TypeOf(tok.ElementID)
		if err != nil || typ != eventv1.Element_TYPE_END_EVENT {
			return nil, nil
		}
	}
	h, err := x.Handlers.Get(eventv1.Element_TYPE_PROCESS)
	if err != nil {
		return nil, err
	}
	effect, err := h.OnComplete(handlers.CompleteInput{Deployment: dep, Instance: inst})
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
	return pubs, nil
}

func (x *Executor) startCompensation(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	throwTokenID, throwElementID string,
	emit Emitter,
) ([]handlers.Publication, error) {
	activityRef, err := dep.CompensateActivityRef(throwElementID)
	if err != nil {
		return nil, err
	}
	throwScope, _ := dep.ScopeOf(throwElementID)
	// Compensate throws inside a compensation event sub-process undo work in the
	// enclosing compensated SubProcess, not inside the event sub-process itself.
	if esp, ok := dep.EventSubProcessSpec(throwScope); ok && esp.Kind == deploy.CatchKindCompensate {
		throwScope = esp.ParentScopeID
	}

	type item struct {
		boundaryID string
		handlerID  string
		seq        int64
	}
	seen := make(map[string]bool)
	var items []item
	addItem := func(bid string, sub *projection.CompensationSub, c deploy.Compensation) {
		if sub == nil || seen[bid] {
			return
		}
		seen[bid] = true
		items = append(items, item{boundaryID: bid, handlerID: c.HandlerID, seq: sub.Seq})
	}
	collectSameScope := func() {
		for bid, sub := range inst.CompensationSubs {
			if sub == nil {
				continue
			}
			c, ok := dep.CompensationByBoundary(bid)
			if !ok {
				continue
			}
			if activityRef != "" && c.ActivityID != activityRef {
				continue
			}
			actScope, _ := dep.ScopeOf(c.ActivityID)
			if actScope != throwScope {
				continue
			}
			addItem(bid, sub, c)
		}
	}
	collectUnderSubProcess := func(spID string) {
		for bid, sub := range inst.CompensationSubs {
			if sub == nil {
				continue
			}
			c, ok := dep.CompensationByBoundary(bid)
			if !ok {
				continue
			}
			if !activityInScopeTree(dep, c.ActivityID, spID) {
				continue
			}
			addItem(bid, sub, c)
		}
	}

	var pubs []handlers.Publication
	appendUnfinishedCall := func(callActivityID, hostTokenID, childInstanceID string) {
		if childInstanceID == "" {
			return
		}
		pubs = append(pubs, handlers.Publication{
			Kind:            handlers.PublicationCompensateUnfinishedChild,
			ChildInstanceID: childInstanceID,
			ParentInstanceID: inst.ID,
			CallActivityID:  callActivityID,
			HostTokenID:     hostTokenID,
		})
	}

	if hostTok, childID, ok := unfinishedCallActivityHost(dep, inst, activityRef); ok && scopeParentIs(dep, activityRef, throwScope) {
		appendUnfinishedCall(activityRef, hostTok, childID)
	} else if unfinishedSubProcessHost(dep, inst, activityRef) && scopeParentIs(dep, activityRef, throwScope) {
		if err := x.terminateScope(ctx, dep, inst, activityRef, emit); err != nil {
			return nil, err
		}
		collectUnderSubProcess(activityRef)
	} else if activityRef != "" {
		collectSameScope()
	} else {
		for _, ca := range unfinishedChildCallActivities(dep, inst, throwScope) {
			appendUnfinishedCall(ca.callActivityID, ca.hostTokenID, ca.childInstanceID)
		}
		for _, spID := range unfinishedChildSubProcesses(dep, inst, throwScope) {
			if err := x.terminateScope(ctx, dep, inst, spID, emit); err != nil {
				return nil, err
			}
			collectUnderSubProcess(spID)
		}
		collectSameScope()
	}

	// Reverse order of completion (higher Seq first).
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].seq > items[i].seq {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
	queue := make([]string, 0, len(items))
	consumed := make([]string, 0, len(items))
	for _, it := range items {
		queue = append(queue, it.handlerID)
		consumed = append(consumed, it.boundaryID)
	}
	for _, bid := range consumed {
		if err := emit(&eventv1.Element{
			Intent:  eventv1.Element_INTENT_COMPLETING,
			Type:    eventv1.Element_TYPE_BOUNDARY_EVENT,
			Id:      bid,
			TokenId: throwTokenID,
		}); err != nil {
			return pubs, err
		}
		if err := emit(&eventv1.Element{
			Intent:  eventv1.Element_INTENT_COMPLETED,
			Type:    eventv1.Element_TYPE_BOUNDARY_EVENT,
			Id:      bid,
			TokenId: throwTokenID,
		}); err != nil {
			return pubs, err
		}
	}
	inst.PendingCompensation = &projection.PendingCompensation{
		ThrowTokenID:    throwTokenID,
		ThrowElementID:  throwElementID,
		Queue:           queue,
		Consumed:        consumed,
		Parent:          inst.PendingCompensation,
		WaitingChildren: len(pubs),
	}
	more, err := x.advanceCompensation(ctx, dep, inst, emit)
	pubs = append(pubs, more...)
	return pubs, err
}

// activityInScopeTree reports whether activityID lives under rootScopeID
// (ScopeOf chain includes rootScopeID).
func activityInScopeTree(dep *deploy.Deployment, activityID, rootScopeID string) bool {
	if dep == nil || activityID == "" || rootScopeID == "" {
		return false
	}
	s, ok := dep.ScopeOf(activityID)
	for ok && s != "" {
		if s == rootScopeID {
			return true
		}
		s, ok = dep.ScopeOf(s)
	}
	return false
}

func scopeParentIs(dep *deploy.Deployment, elementID, parentScopeID string) bool {
	if dep == nil {
		return false
	}
	p, ok := dep.ScopeOf(elementID)
	return ok && p == parentScopeID
}

func unfinishedSubProcessHost(dep *deploy.Deployment, inst *projection.Instance, spID string) bool {
	if dep == nil || inst == nil || spID == "" {
		return false
	}
	if !dep.IsEmbeddedScope(spID) {
		return false
	}
	for _, tok := range inst.Tokens {
		if tok != nil && tok.ElementID == spID {
			return true
		}
	}
	return false
}

func unfinishedChildSubProcesses(dep *deploy.Deployment, inst *projection.Instance, throwScope string) []string {
	if dep == nil || inst == nil {
		return nil
	}
	seen := make(map[string]bool)
	var ids []string
	for _, tok := range inst.Tokens {
		if tok == nil {
			continue
		}
		spID := tok.ElementID
		if seen[spID] {
			continue
		}
		if !dep.IsEmbeddedScope(spID) {
			continue
		}
		parent, ok := dep.ScopeOf(spID)
		if !ok || parent != throwScope {
			continue
		}
		seen[spID] = true
		ids = append(ids, spID)
	}
	sort.Strings(ids)
	return ids
}

func (x *Executor) advanceCompensation(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	emit Emitter,
) ([]handlers.Publication, error) {
	pc := inst.PendingCompensation
	if pc == nil {
		pc = rebuildPendingCompensation(inst, dep)
		if pc == nil {
			return nil, nil
		}
		inst.PendingCompensation = pc
	}
	if len(pc.Queue) == 0 {
		if pc.WaitingChildren > 0 {
			return nil, nil
		}
		if pc.NotifyParentUnfinishedCall {
			inst.PendingCompensation = pc.Parent
			return x.terminateInstancePubs(ctx, dep, inst, emit)
		}
		if pc.CancelBoundaryID != "" {
			throwTok, throwEl := pc.ThrowTokenID, pc.ThrowElementID
			boundaryID := pc.CancelBoundaryID
			hostTok := pc.CancelHostTokenID
			txID := pc.CancelTransactionID
			inst.PendingCompensation = pc.Parent
			more, err := x.Complete(ctx, dep, inst, throwTok, throwEl, nil, emit)
			if err != nil {
				return more, err
			}
			more2, err := x.fireCancelBoundary(ctx, dep, inst, txID, hostTok, boundaryID, emit)
			return append(more, more2...), err
		}
		throwTok, throwEl := pc.ThrowTokenID, pc.ThrowElementID
		inst.PendingCompensation = pc.Parent
		return x.Complete(ctx, dep, inst, throwTok, throwEl, nil, emit)
	}
	handler := pc.Queue[0]
	pc.Queue = pc.Queue[1:]
	hid, err := NextID()
	if err != nil {
		return nil, err
	}
	pc.ActiveTokenID = hid
	pc.ActiveHandler = handler
	return x.Enter(ctx, dep, inst, hid, handler, emit)
}

// rebuildPendingCompensation reconstructs in-flight compensate throw state after event replay.
func rebuildPendingCompensation(inst *projection.Instance, dep *deploy.Deployment) *projection.PendingCompensation {
	if inst == nil || dep == nil {
		return nil
	}
	var throwTokenID, throwElementID string
	for tid, tok := range inst.Tokens {
		if tok == nil || tok.Status != projection.TokenWaiting {
			continue
		}
		typ, err := dep.TypeOf(tok.ElementID)
		if err != nil {
			continue
		}
		switch typ {
		case eventv1.Element_TYPE_INTERMEDIATE_THROW_EVENT:
			kind, err := dep.ThrowKind(tok.ElementID)
			if err != nil || kind != deploy.ThrowKindCompensate {
				continue
			}
		case eventv1.Element_TYPE_END_EVENT:
			if !dep.IsCompensateEnd(tok.ElementID) && !dep.IsCancelEnd(tok.ElementID) {
				continue
			}
		default:
			continue
		}
		throwTokenID = tid
		throwElementID = tok.ElementID
		break
	}
	var activeTokenID, activeHandler string
	for tid, tok := range inst.Tokens {
		if tok == nil || tok.Status != projection.TokenWaiting {
			continue
		}
		if !dep.IsCompensationHandler(tok.ElementID) {
			continue
		}
		activeTokenID = tid
		activeHandler = tok.ElementID
		break
	}
	if throwTokenID == "" {
		if !instHasParentCall(inst) || activeHandler == "" {
			return nil
		}
		return &projection.PendingCompensation{
			ActiveTokenID:              activeTokenID,
			ActiveHandler:              activeHandler,
			NotifyParentUnfinishedCall: true,
		}
	}
	waitingChildren := 0
	for _, tok := range inst.Tokens {
		if tok == nil || tok.Status != projection.TokenWaiting {
			continue
		}
		typ, err := dep.TypeOf(tok.ElementID)
		if err != nil || typ != eventv1.Element_TYPE_CALL_ACTIVITY {
			continue
		}
		if tok.CalledProcessInstanceID != "" {
			waitingChildren++
		}
	}
	return &projection.PendingCompensation{
		ThrowTokenID:    throwTokenID,
		ThrowElementID:  throwElementID,
		ActiveTokenID:   activeTokenID,
		ActiveHandler:   activeHandler,
		WaitingChildren: waitingChildren,
		CancelBoundaryID: func() string {
			if dep.IsCancelEnd(throwElementID) {
				if txID, ok := dep.ScopeOf(throwElementID); ok {
					if bid, ok := dep.CancelBoundaryOf(txID); ok {
						return bid
					}
				}
			}
			return ""
		}(),
		CancelTransactionID: func() string {
			if dep.IsCancelEnd(throwElementID) {
				if txID, ok := dep.ScopeOf(throwElementID); ok {
					return txID
				}
			}
			return ""
		}(),
		CancelHostTokenID: func() string {
			if !dep.IsCancelEnd(throwElementID) {
				return ""
			}
			txID, ok := dep.ScopeOf(throwElementID)
			if !ok {
				return ""
			}
			for tid, tok := range inst.Tokens {
				if tok != nil && tok.ElementID == txID {
					return tid
				}
			}
			return ""
		}(),
	}
}

func instHasParentCall(inst *projection.Instance) bool {
	return inst != nil && inst.ParentProcessInstanceID != "" && inst.ParentElementID != "" && inst.ParentTokenID != ""
}

type unfinishedCallHost struct {
	callActivityID  string
	hostTokenID     string
	childInstanceID string
}

func unfinishedCallActivityHost(dep *deploy.Deployment, inst *projection.Instance, callActivityID string) (hostTokenID, childInstanceID string, ok bool) {
	if dep == nil || inst == nil || callActivityID == "" {
		return "", "", false
	}
	typ, err := dep.TypeOf(callActivityID)
	if err != nil || typ != eventv1.Element_TYPE_CALL_ACTIVITY {
		return "", "", false
	}
	for tid, tok := range inst.Tokens {
		if tok == nil || tok.ElementID != callActivityID {
			continue
		}
		if tok.CalledProcessInstanceID == "" {
			continue
		}
		return tid, tok.CalledProcessInstanceID, true
	}
	return "", "", false
}

func unfinishedChildCallActivities(dep *deploy.Deployment, inst *projection.Instance, throwScope string) []unfinishedCallHost {
	if dep == nil || inst == nil {
		return nil
	}
	seen := make(map[string]bool)
	var out []unfinishedCallHost
	for tid, tok := range inst.Tokens {
		if tok == nil || tok.CalledProcessInstanceID == "" {
			continue
		}
		caID := tok.ElementID
		if seen[caID] {
			continue
		}
		typ, err := dep.TypeOf(caID)
		if err != nil || typ != eventv1.Element_TYPE_CALL_ACTIVITY {
			continue
		}
		parent, ok := dep.ScopeOf(caID)
		if !ok || parent != throwScope {
			continue
		}
		seen[caID] = true
		out = append(out, unfinishedCallHost{
			callActivityID:  caID,
			hostTokenID:     tid,
			childInstanceID: tok.CalledProcessInstanceID,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].callActivityID < out[j].callActivityID })
	return out
}
