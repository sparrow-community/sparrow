package processing

import (
	"context"
	"fmt"

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
			// Uncaught in this scope: bubble without terminating (unlike error).
			parent, ok := dep.ScopeOf(scope)
			if !ok {
				break
			}
			scope = parent
			continue
		}
		// Process scope: optional process-level activity boundary is not applicable for
		// throw/end; uncaught → no-op.
		return nil, nil
	}
	return nil, nil
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
