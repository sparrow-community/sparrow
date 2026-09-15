package processing

import (
	"context"
	"fmt"

	"github.com/sparrow-community/sparrow/processing/deploy"
	"github.com/sparrow-community/sparrow/processing/handlers"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func (x *Executor) startTransactionCancel(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	cancelTokenID, cancelEndID string,
	emit Emitter,
) ([]handlers.Publication, error) {
	txID, ok := dep.ScopeOf(cancelEndID)
	if !ok || !dep.IsTransaction(txID) {
		return nil, fmt.Errorf("UNSUPPORTED_ELEMENT: cancel end %q is not inside a transaction", cancelEndID)
	}
	boundaryID, ok := dep.CancelBoundaryOf(txID)
	if !ok {
		return nil, fmt.Errorf("UNSUPPORTED_ELEMENT: transaction %q has no cancel boundary", txID)
	}
	hostTok, ok := findScopeHostToken(dep, inst, txID)
	if !ok {
		return nil, fmt.Errorf("NOT_FOUND: transaction host for %q", txID)
	}
	keep := map[string]bool{cancelTokenID: true, hostTok: true}
	if err := terminateScopeTokens(dep, inst, txID, emit, scopeTerminateOpts{
		IncludeHost:  false,
		DropTokens:   true,
		KeepTokenIDs: keep,
	}); err != nil {
		return nil, err
	}
	if err := emitEventSubProcessStartDisarmInScope(dep, txID, inst, emit); err != nil {
		return nil, err
	}
	inst.RemoveScopeBoundariesForScope(txID)

	pubs, err := x.startCompensation(ctx, dep, inst, cancelTokenID, cancelEndID, emit)
	if err != nil {
		return pubs, err
	}
	if pc := inst.PendingCompensation; pc != nil {
		pc.CancelBoundaryID = boundaryID
		pc.CancelTransactionID = txID
		pc.CancelHostTokenID = hostTok
	}
	return pubs, nil
}

func (x *Executor) fireCancelBoundary(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	txID, hostTokenID, boundaryID string,
	emit Emitter,
) ([]handlers.Publication, error) {
	// Terminate any leftover tokens in the transaction except the host (boundary continues on host).
	if err := terminateScopeTokens(dep, inst, txID, emit, scopeTerminateOpts{
		IncludeHost:  false,
		DropTokens:   true,
		KeepTokenIDs: map[string]bool{hostTokenID: true},
	}); err != nil {
		return nil, err
	}
	if err := terminateScopeHost(dep, txID, hostTokenID, emit); err != nil {
		return nil, err
	}
	inst.RemoveScopeBoundariesForScope(txID)
	if err := emit(&eventv1.Element{
		Intent:  eventv1.Element_INTENT_COMPLETING,
		Type:    eventv1.Element_TYPE_BOUNDARY_EVENT,
		Id:      boundaryID,
		TokenId: hostTokenID,
	}); err != nil {
		return nil, err
	}
	if err := emit(&eventv1.Element{
		Intent:  eventv1.Element_INTENT_COMPLETED,
		Type:    eventv1.Element_TYPE_BOUNDARY_EVENT,
		Id:      boundaryID,
		TokenId: hostTokenID,
	}); err != nil {
		return nil, err
	}
	return x.leaveViaOutgoings(ctx, dep, inst, hostTokenID, boundaryID, "", emit)
}
