package processing

import (
	"context"
	"strings"

	"github.com/sparrow-community/sparrow/processing/projection"
)

type dueWait struct {
	instanceID  string
	elementID   string
	tokenID     string
	MessageName string
}

// FireDue completes waiting timer catch tokens whose due time has been reached.
// Each due token is finished through Complete (same waiting story as UserTask).
func (e *Engine) FireDue(ctx context.Context) error {
	now := e.now().UnixMilli()
	due := e.collectDue(now)
	var first error
	for _, w := range due {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := e.Complete(ctx, w.instanceID, w.elementID, w.tokenID, nil); err != nil {
			if strings.HasPrefix(err.Error(), "INVALID_STATE:") {
				continue
			}
			if first == nil {
				first = err
			}
		}
	}
	// Scope boundaries (SubProcess timer boundaries)
	scopeDue := e.collectScopeDue(now)
	for _, sd := range scopeDue {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := e.completeScopeBoundary(ctx, sd.instanceID, sd.boundaryID); err != nil {
			if first == nil {
				first = err
			}
		}
	}
	for _, ew := range e.collectEventSubProcessTimerDue(now) {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := e.triggerEventSubProcess(ctx, ew.instanceID, ew.eventSubProcessElementID, nil); err != nil {
			if first == nil {
				first = err
			}
		}
	}
	return first
}

type scopeDueWait struct {
	instanceID string
	boundaryID string
}

func (e *Engine) collectScopeDue(nowUnixMs int64) []scopeDueWait {
	e.mu.Lock()
	ids := make([]string, 0, len(e.instances))
	for id := range e.instances {
		ids = append(ids, id)
	}
	e.mu.Unlock()

	var due []scopeDueWait
	for _, iid := range ids {
		e.mu.Lock()
		inst := e.instances[iid]
		lock := e.instMu[iid]
		e.mu.Unlock()
		if inst == nil || lock == nil {
			continue
		}
		lock.Lock()
		for _, sb := range inst.ScopeBoundaries {
			if sb.DueUnixMs > 0 && sb.DueUnixMs <= nowUnixMs {
				due = append(due, scopeDueWait{instanceID: iid, boundaryID: sb.BoundaryID})
			}
		}
		lock.Unlock()
	}
	return due
}

func (e *Engine) collectDue(nowUnixMs int64) []dueWait {
	e.mu.Lock()
	ids := make([]string, 0, len(e.instances))
	for id := range e.instances {
		ids = append(ids, id)
	}
	e.mu.Unlock()

	var due []dueWait
	for _, iid := range ids {
		e.mu.Lock()
		inst := e.instances[iid]
		lock := e.instMu[iid]
		e.mu.Unlock()
		if inst == nil || lock == nil {
			continue
		}
		lock.Lock()
		for _, tok := range inst.Tokens {
			if tok == nil || (tok.Status != projection.TokenWaiting && tok.Status != projection.TokenBlocked) {
				continue
			}
			if len(tok.BoundaryWaits) > 0 {
				for _, w := range tok.BoundaryWaits {
					if w.Kind == "timer" && w.DueUnixMs > 0 && w.DueUnixMs <= nowUnixMs {
						due = append(due, dueWait{instanceID: iid, elementID: w.BoundaryID, tokenID: tok.ID})
					}
				}
				continue
			}
			if tok.DueUnixMs > 0 && tok.DueUnixMs <= nowUnixMs {
				due = append(due, dueWait{instanceID: iid, elementID: waiterElementID(tok), tokenID: tok.ID})
			}
		}
		lock.Unlock()
	}
	return due
}

func waiterElementID(tok *projection.Token) string {
	if tok.BoundaryID != "" {
		return tok.BoundaryID
	}
	return tok.ElementID
}

func messageWaiterElementID(tok *projection.Token) string {
	if tok.MessageBoundaryID != "" {
		return tok.MessageBoundaryID
	}
	if tok.BoundaryID != "" {
		return tok.BoundaryID
	}
	return tok.ElementID
}

func signalWaiterElementID(tok *projection.Token) string {
	if tok.SignalBoundaryID != "" {
		return tok.SignalBoundaryID
	}
	if tok.BoundaryID != "" && tok.SignalName != "" {
		return tok.BoundaryID
	}
	return tok.ElementID
}
