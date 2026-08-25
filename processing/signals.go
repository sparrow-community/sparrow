package processing

import (
	"context"
	"fmt"
	"strings"

	"github.com/sparrow-community/sparrow/processing/handlers"
	"github.com/sparrow-community/sparrow/processing/projection"
)

// PublishSignalRequest broadcasts a BPMN signal to waiting signal catches.
// Unlike messages, signals are not buffered when no waiter matches.
type PublishSignalRequest struct {
	Name              string
	ProcessInstanceID string // empty: all matching waiters
	Variables         map[string]any
}

// PublishSignal completes waiting signal-catch tokens whose SignalName matches.
// Late signals with no waiter are dropped (not buffered).
func (e *Engine) PublishSignal(ctx context.Context, req PublishSignalRequest) (int, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return 0, fmt.Errorf("INVALID_ARGUMENT: signal name is required")
	}
	instanceID := strings.TrimSpace(req.ProcessInstanceID)
	if instanceID != "" {
		e.mu.Lock()
		_, ok := e.instances[instanceID]
		e.mu.Unlock()
		if !ok {
			return 0, fmt.Errorf("NOT_FOUND: instance %q", instanceID)
		}
	}
	waiters := e.collectSignalWaiters(name, instanceID)
	scopeWaiters := e.collectScopeSignalWaiters(name, instanceID)
	eventSubProcessWaiters := e.collectEventSubProcessSignalArms(name, instanceID)
	if len(waiters) == 0 && len(scopeWaiters) == 0 && len(eventSubProcessWaiters) == 0 {
		return 0, nil
	}
	delivered := 0
	var first error
	for _, w := range waiters {
		if err := ctx.Err(); err != nil {
			return delivered, err
		}
		if err := e.Complete(ctx, w.instanceID, w.elementID, w.tokenID, req.Variables); err != nil {
			if strings.HasPrefix(err.Error(), "INVALID_STATE:") {
				continue
			}
			if first == nil {
				first = err
			}
			continue
		}
		delivered++
	}
	for _, sw := range scopeWaiters {
		if err := ctx.Err(); err != nil {
			return delivered, err
		}
		if err := e.completeScopeBoundary(ctx, sw.instanceID, sw.boundaryID); err != nil {
			if first == nil {
				first = err
			}
			continue
		}
		delivered++
	}
	for _, ew := range eventSubProcessWaiters {
		if err := ctx.Err(); err != nil {
			return delivered, err
		}
		if err := e.triggerEventSubProcess(ctx, ew.instanceID, ew.eventSubProcessElementID, req.Variables); err != nil {
			if first == nil {
				first = err
			}
			continue
		}
		delivered++
	}
	if delivered == 0 && first != nil {
		return 0, first
	}
	return delivered, first
}

func (e *Engine) collectSignalWaiters(name, instanceID string) []dueWait {
	e.mu.Lock()
	ids := make([]string, 0, len(e.instances))
	if instanceID != "" {
		if _, ok := e.instances[instanceID]; ok {
			ids = append(ids, instanceID)
		}
	} else {
		for id := range e.instances {
			ids = append(ids, id)
		}
	}
	e.mu.Unlock()

	var waiters []dueWait
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
			if tok.SignalName != "" && tok.SignalName == name {
				waiters = append(waiters, dueWait{instanceID: iid, elementID: signalWaiterElementID(tok), tokenID: tok.ID})
			}
		}
		lock.Unlock()
	}
	return waiters
}

func (e *Engine) collectScopeSignalWaiters(name, instanceID string) []scopeDueWait {
	e.mu.Lock()
	ids := make([]string, 0, len(e.instances))
	if instanceID != "" {
		if _, ok := e.instances[instanceID]; ok {
			ids = append(ids, instanceID)
		}
	} else {
		for id := range e.instances {
			ids = append(ids, id)
		}
	}
	e.mu.Unlock()

	var waiters []scopeDueWait
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
			if sb.SignalName != "" && sb.SignalName == name {
				waiters = append(waiters, scopeDueWait{instanceID: iid, boundaryID: sb.BoundaryID})
			}
		}
		lock.Unlock()
	}
	return waiters
}

// flushPublications delivers deferred throws after the instance lock is released.
func (e *Engine) flushPublications(ctx context.Context, pubs []handlers.Publication) error {
	for _, p := range pubs {
		if err := ctx.Err(); err != nil {
			return err
		}
		name := strings.TrimSpace(p.Name)
		if name == "" {
			continue
		}
		switch p.Kind {
		case handlers.PublicationMessage:
			if _, err := e.PublishMessage(ctx, PublishMessageRequest{Name: name}); err != nil {
				return err
			}
		case handlers.PublicationSignal:
			if _, err := e.PublishSignal(ctx, PublishSignalRequest{Name: name}); err != nil {
				return err
			}
		default:
			return fmt.Errorf("UNSUPPORTED_ELEMENT: unknown publication kind %q", p.Kind)
		}
	}
	return nil
}
