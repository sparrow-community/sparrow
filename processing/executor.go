package processing

import (
	"context"
	"fmt"
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
}

func NewExecutor(reg *handlers.Registry) *Executor {
	if reg == nil {
		reg = handlers.DefaultRegistry()
	}
	return &Executor{Handlers: reg}
}

func (x *Executor) now() time.Time {
	if x != nil && x.Now != nil {
		return x.Now()
	}
	return time.Now()
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
			Deployment: dep,
			Instance:   inst,
			ElementID:  elementID,
			Type:       typ,
			TokenID:    tokenID,
			Now:        x.now(),
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
		if effect.EnterChild != "" {
			childTokenID := tokenID
			if effect.SpawnChildToken {
				var err error
				childTokenID, err = NextID()
				if err != nil {
					return pubs, err
				}
			}
			if typ == eventv1.Element_TYPE_SUB_PROCESS && !dep.IsEventSubProcess(elementID) {
				if err := emitEventSubProcessStartArms(dep, inst, elementID, x.now(), emit); err != nil {
					return pubs, err
				}
			}
			elementID = effect.EnterChild
			tokenID = childTokenID
			continue
		}
		if len(effect.Fork) > 0 {
			for i, flowID := range effect.Fork {
				tid := tokenID
				if i > 0 {
					var err error
					tid, err = NextID()
					if err != nil {
						return pubs, err
					}
				}
				next, err := x.takeOutgoing(dep, tid, elementID, flowID, emit)
				if err != nil {
					return pubs, err
				}
				more, err := x.Enter(ctx, dep, inst, tid, next, emit)
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
		if effect.TryCompleteProcess {
			more, err := x.tryCompleteScope(ctx, dep, inst, tokenID, elementID, emit)
			pubs = append(pubs, more...)
			return pubs, err
		}
		if !effect.TakeOutgoing {
			return pubs, nil
		}
		if effect.TerminateJoinPeers != "" {
			if err := x.terminateJoinPeers(inst, tokenID, effect.TerminateJoinPeers, emit); err != nil {
				return pubs, err
			}
		}

		next, err := x.takeOutgoing(dep, tokenID, elementID, effect.OutgoingFlowID, emit)
		if err != nil {
			return pubs, err
		}
		elementID = next
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
	var pubs []handlers.Publication
	for _, rec := range effect.Records {
		if err := emit(rec); err != nil {
			return pubs, err
		}
	}
	if effect.Publish != nil {
		pubs = append(pubs, *effect.Publish)
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
		next, err := x.takeOutgoing(dep, spawnID, spawn.ElementID, spawn.OutgoingFlowID, emit)
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

	next, err := x.takeOutgoing(dep, tokenID, elementID, effect.OutgoingFlowID, emit)
	if err != nil {
		return pubs, err
	}
	more, err := x.Enter(ctx, dep, inst, tokenID, next, emit)
	pubs = append(pubs, more...)
	return pubs, err
}

func (x *Executor) terminateJoinPeers(inst *projection.Instance, survivorTokenID, joinElementID string, emit Emitter) error {
	typ := eventv1.Element_TYPE_PARALLEL_GATEWAY
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
	tokenID, fromElementID, preferredFlowID string,
	emit Emitter,
) (string, error) {
	flowID := preferredFlowID
	if flowID == "" {
		outs := dep.Outgoing(fromElementID)
		if len(outs) == 0 {
			return "", fmt.Errorf("NO_OUTGOING_FLOW: element %q", fromElementID)
		}
		flowID = outs[0]
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
	// SubProcess scope: ignore parked host; require children at EndEvents.
	hostTokenID := ""
	for tid, tok := range inst.Tokens {
		if tok == nil {
			continue
		}
		if tok.ElementID == scopeID {
			hostTokenID = tid
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
	for tid, tok := range inst.Tokens {
		if tok == nil || tid == completeTokenID {
			continue
		}
		if tok.ElementID == scopeID {
			continue
		}
		tokScope, _ := dep.ScopeOf(tok.ElementID)
		if tokScope == scopeID || isInScope(dep, tokScope, scopeID) {
			delete(inst.Tokens, tid)
		}
	}
	// All tokens in this scope are at EndEvents; disarm scope boundaries and complete
	inst.RemoveScopeBoundariesForScope(scopeID)
	if err := emitEventSubProcessStartDisarmInScope(dep, scopeID, inst, emit); err != nil {
		return nil, err
	}

	completeID := scopeID
	completeType := eventv1.Element_TYPE_SUB_PROCESS
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
	if effect.DiscardToken {
		delete(inst.Tokens, completeTokenID)
	}
	if dep.IsEventSubProcess(scopeID) {
		if spec, ok := dep.EventSubProcessSpec(scopeID); ok && !spec.Interrupting {
			if err := emit(eventSubProcessStartActivated(dep, spec, x.now())); err != nil {
				return pubs, err
			}
		}
	}
	if effect.TryCompleteProcess {
		more, err := x.tryCompleteProcessScope(dep, inst, emit)
		pubs = append(pubs, more...)
		return pubs, err
	}
	if effect.TakeOutgoing {
		next, err := x.takeOutgoing(dep, completeTokenID, outgoingFrom, effect.OutgoingFlowID, emit)
		if err != nil {
			return pubs, err
		}
		more, err := x.Enter(ctx, dep, inst, completeTokenID, next, emit)
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

	type item struct {
		boundaryID string
		handlerID  string
		seq        int64
	}
	var items []item
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
		items = append(items, item{boundaryID: bid, handlerID: c.HandlerID, seq: sub.Seq})
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
			return nil, err
		}
		if err := emit(&eventv1.Element{
			Intent:  eventv1.Element_INTENT_COMPLETED,
			Type:    eventv1.Element_TYPE_BOUNDARY_EVENT,
			Id:      bid,
			TokenId: throwTokenID,
		}); err != nil {
			return nil, err
		}
	}
	inst.PendingCompensation = &projection.PendingCompensation{
		ThrowTokenID:   throwTokenID,
		ThrowElementID: throwElementID,
		Queue:          queue,
		Consumed:       consumed,
	}
	return x.advanceCompensation(ctx, dep, inst, emit)
}

func (x *Executor) advanceCompensation(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	emit Emitter,
) ([]handlers.Publication, error) {
	pc := inst.PendingCompensation
	if pc == nil {
		return nil, nil
	}
	if len(pc.Queue) == 0 {
		throwTok, throwEl := pc.ThrowTokenID, pc.ThrowElementID
		inst.PendingCompensation = nil
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
