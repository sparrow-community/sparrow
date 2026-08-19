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
	return first
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
			if tok == nil || tok.Status != projection.TokenWaiting {
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
