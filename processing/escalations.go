package processing

import (
	"context"
	"fmt"
	"sort"

	"github.com/sparrow-community/sparrow/processing/deploy"
	"github.com/sparrow-community/sparrow/processing/handlers"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func (x *Executor) propagateEscalation(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	throwElementID, throwTokenID, escalationCode string,
	emit Emitter,
) ([]handlers.Publication, error) {
	scope, ok := dep.ScopeOf(throwElementID)
	if !ok {
		return nil, fmt.Errorf("NOT_FOUND: scope for element %q", throwElementID)
	}
	processID := inst.ProcessID
	if processID == "" {
		processID = dep.ProcessID()
	}

	for {
		if eventSubProcessID, ok := dep.MatchEscalationEventSubProcess(scope, escalationCode, func(id string) bool {
			_, armed := inst.EventSubProcesses[id]
			return armed
		}); ok {
			return x.fireErrorEventSubProcess(ctx, dep, inst, throwElementID, throwTokenID, eventSubProcessID, emit)
		}
		if scope != processID {
			if bid, ok := dep.MatchEscalationBoundary(scope, escalationCode); ok {
				hostTokenID, ok := findScopeHostToken(dep, inst, scope)
				if !ok {
					return nil, fmt.Errorf("INVALID_STATE: no token for scope %q", scope)
				}
				if dep.BoundaryInterrupting(bid) {
					return x.fireScopeErrorBoundary(ctx, dep, inst, scope, hostTokenID, bid, emit)
				}
				return x.fireNonInterruptingScopeEscalationBoundary(ctx, dep, inst, bid, emit)
			}
		}
		if pubs, handled, err := x.fireWaitingEscalationCatches(ctx, dep, inst, scope, escalationCode, emit); handled || err != nil {
			return pubs, err
		}
		if scope == processID {
			// Uncaught at process scope → no-op (escalation does not terminate).
			return nil, nil
		}
		parent, ok := dep.ScopeOf(scope)
		if !ok {
			break
		}
		scope = parent
	}
	return nil, nil
}

func (x *Executor) fireWaitingEscalationCatches(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	scopeID, escalationCode string,
	emit Emitter,
) ([]handlers.Publication, bool, error) {
	type waiter struct {
		tokenID   string
		elementID string
	}
	var waiters []waiter
	for tid, tok := range inst.Tokens {
		if tok == nil || tok.Status != projection.TokenWaiting {
			continue
		}
		code, ok := dep.EscalationIntermediateCatchCode(tok.ElementID)
		if !ok {
			continue
		}
		tokScope, _ := dep.ScopeOf(tok.ElementID)
		if tokScope != scopeID {
			continue
		}
		if !escalationCatchMatches(code, escalationCode) {
			continue
		}
		waiters = append(waiters, waiter{tokenID: tid, elementID: tok.ElementID})
	}
	if len(waiters) == 0 {
		return nil, false, nil
	}
	sort.Slice(waiters, func(i, j int) bool {
		if waiters[i].elementID == waiters[j].elementID {
			return waiters[i].tokenID < waiters[j].tokenID
		}
		return waiters[i].elementID < waiters[j].elementID
	})
	var pubs []handlers.Publication
	for _, w := range waiters {
		if tok := inst.Tokens[w.tokenID]; tok == nil || tok.Status != projection.TokenWaiting {
			continue
		}
		more, err := x.Complete(ctx, dep, inst, w.tokenID, w.elementID, nil, emit)
		pubs = append(pubs, more...)
		if err != nil {
			return pubs, true, err
		}
	}
	return pubs, true, nil
}

func escalationCatchMatches(catchCode, thrownCode string) bool {
	if catchCode == "" {
		return true
	}
	return catchCode == thrownCode
}

func (x *Executor) fireNonInterruptingScopeEscalationBoundary(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	boundaryID string,
	emit Emitter,
) ([]handlers.Publication, error) {
	spawnID, err := NextID()
	if err != nil {
		return nil, err
	}
	for _, rec := range handlers.InstantLifecycle(eventv1.Element_TYPE_BOUNDARY_EVENT, boundaryID, spawnID, nil) {
		if err := emit(rec); err != nil {
			return nil, err
		}
	}
	next, err := x.takeOutgoing(dep, inst, spawnID, boundaryID, "", emit)
	if err != nil {
		return nil, err
	}
	return x.Enter(ctx, dep, inst, spawnID, next, emit)
}
