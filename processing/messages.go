package processing

import (
	"context"
	"fmt"
	"strings"

	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// PublishMessageRequest correlates a BPMN message to waiting message catches.
// Messages that arrive before a waiter exists are not buffered.
type PublishMessageRequest struct {
	Name              string
	ProcessInstanceID string // empty: all matching waiters
	// CorrelationKeys match instance variables (JSON-encoded, same as CreateInstance).
	// Empty means name (and optional ProcessInstanceID) only.
	CorrelationKeys map[string]any
	Variables       map[string]any
}

// PublishMessage completes waiting message-catch tokens whose MessageName matches.
// Each match is finished through Complete (same waiting story as UserTask / timer).
func (e *Engine) PublishMessage(ctx context.Context, req PublishMessageRequest) (int, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return 0, fmt.Errorf("INVALID_ARGUMENT: message name is required")
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
	keys, err := projection.VariablesFromMap(req.CorrelationKeys)
	if err != nil {
		return 0, fmt.Errorf("INVALID_ARGUMENT: correlation_keys: %w", err)
	}
	waiters := e.collectMessageWaiters(name, instanceID, keys)
	if len(waiters) == 0 {
		if instanceID != "" {
			return 0, fmt.Errorf("NOT_FOUND: no waiting message catch %q on instance %q", name, instanceID)
		}
		return 0, fmt.Errorf("NOT_FOUND: no waiting message catch %q", name)
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
	if delivered == 0 && first != nil {
		return 0, first
	}
	if delivered == 0 {
		return 0, fmt.Errorf("NOT_FOUND: no waiting message catch %q", name)
	}
	return delivered, first
}

func (e *Engine) collectMessageWaiters(name, instanceID string, keys []*eventv1.Variable) []dueWait {
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
		if !correlationKeysMatch(inst.Variables, keys) {
			lock.Unlock()
			continue
		}
		for _, tok := range inst.Tokens {
			if tok == nil || tok.Status != projection.TokenWaiting {
				continue
			}
			if tok.MessageName != "" && tok.MessageName == name {
				waiters = append(waiters, dueWait{instanceID: iid, elementID: tok.ElementID, tokenID: tok.ID})
			}
		}
		lock.Unlock()
	}
	return waiters
}

func correlationKeysMatch(vars map[string]string, keys []*eventv1.Variable) bool {
	for _, k := range keys {
		if vars[k.GetName()] != k.GetJsonValue() {
			return false
		}
	}
	return true
}
