package processing

import (
	"context"
	"fmt"
	"strings"

	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// PublishMessageRequest correlates a BPMN message to waiting message catches.
// If no waiter matches, the message is kept in an in-memory buffer until a
// matching catch activates. The buffer is not an EventLog record and is empty
// after Recover.
type PublishMessageRequest struct {
	Name              string
	ProcessInstanceID string // empty: all matching waiters
	// CorrelationKeys match instance variables (JSON-encoded, same as CreateInstance).
	// Empty means name (and optional ProcessInstanceID) only.
	CorrelationKeys map[string]any
	Variables       map[string]any
}

type bufferedMessage struct {
	id         string
	name       string
	instanceID string
	keys       []*eventv1.Variable
	vars       map[string]any
}

// PublishMessage completes waiting message-catch tokens whose MessageName matches.
// Each match is finished through Complete (same waiting story as UserTask / timer).
// With no matching waiter, the message is buffered and delivered when a catch waits.
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
		if err := e.enqueueBuffered(name, instanceID, keys, req.Variables); err != nil {
			return 0, err
		}
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
	if delivered == 0 && first != nil {
		return 0, first
	}
	if delivered == 0 {
		if err := e.enqueueBuffered(name, instanceID, keys, req.Variables); err != nil {
			return 0, err
		}
		return 0, nil
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
				waiters = append(waiters, dueWait{instanceID: iid, elementID: messageWaiterElementID(tok), tokenID: tok.ID, MessageName: tok.MessageName})
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

func (e *Engine) enqueueBuffered(name, instanceID string, keys []*eventv1.Variable, vars map[string]any) error {
	id, err := NextID()
	if err != nil {
		return err
	}
	e.msgMu.Lock()
	e.msgBuf = append(e.msgBuf, bufferedMessage{
		id:         id,
		name:       name,
		instanceID: instanceID,
		keys:       keys,
		vars:       cloneAnyMap(vars),
	})
	e.msgMu.Unlock()
	return nil
}

func (e *Engine) takeBuffered(name, instanceID string, vars map[string]string) (bufferedMessage, bool) {
	e.msgMu.Lock()
	defer e.msgMu.Unlock()
	for i, m := range e.msgBuf {
		if m.name != name {
			continue
		}
		if m.instanceID != "" && m.instanceID != instanceID {
			continue
		}
		if !correlationKeysMatch(vars, m.keys) {
			continue
		}
		e.msgBuf = append(e.msgBuf[:i], e.msgBuf[i+1:]...)
		return m, true
	}
	return bufferedMessage{}, false
}

func (e *Engine) prependBuffered(m bufferedMessage) {
	e.msgMu.Lock()
	e.msgBuf = append([]bufferedMessage{m}, e.msgBuf...)
	e.msgMu.Unlock()
}

func (e *Engine) tryDeliverBuffered(ctx context.Context, instanceID string) error {
	for {
		w, vars, ok := e.messageWaiterOn(instanceID)
		if !ok {
			return nil
		}
		msg, ok := e.takeBuffered(w.MessageName, instanceID, vars)
		if !ok {
			return nil
		}
		if err := e.Complete(ctx, instanceID, w.elementID, w.tokenID, msg.vars); err != nil {
			e.prependBuffered(msg)
			if strings.HasPrefix(err.Error(), "INVALID_STATE:") {
				return nil
			}
			return err
		}
	}
}

func (e *Engine) messageWaiterOn(instanceID string) (dueWait, map[string]string, bool) {
	e.mu.Lock()
	inst := e.instances[instanceID]
	lock := e.instMu[instanceID]
	e.mu.Unlock()
	if inst == nil || lock == nil {
		return dueWait{}, nil, false
	}
	lock.Lock()
	defer lock.Unlock()
	for _, tok := range inst.Tokens {
		if tok == nil || tok.Status != projection.TokenWaiting || tok.MessageName == "" {
			continue
		}
		return dueWait{instanceID: instanceID, elementID: messageWaiterElementID(tok), tokenID: tok.ID, MessageName: tok.MessageName}, cloneStringMap(inst.Variables), true
	}
	return dueWait{}, nil, false
}

func cloneAnyMap(m map[string]any) map[string]any {
	if len(m) == 0 {
		return nil
	}
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
