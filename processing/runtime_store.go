package processing

import (
	"encoding/json"
	"time"

	"github.com/sparrow-community/sparrow/processing/projection"
	"github.com/sparrow-community/sparrow/processing/runtime"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

const defaultMessageBufferTTL = 7 * 24 * time.Hour

func (e *Engine) loadRuntime() error {
	if e.runtimeStore == nil {
		return nil
	}
	now := e.now()

	leases, err := e.runtimeStore.LoadLeases()
	if err != nil {
		return err
	}
	e.jobMu.Lock()
	for _, l := range leases {
		if !e.leaseStillValid(l, now) {
			_ = e.runtimeStore.DeleteLease(l.InstanceID, l.TokenID)
			continue
		}
		e.leases[leaseKey(l.InstanceID, l.TokenID)] = jobLease{
			workerID: l.WorkerID,
			deadline: time.UnixMilli(l.DeadlineUnixMs),
		}
	}
	e.jobMu.Unlock()

	msgs, err := e.runtimeStore.LoadMessages()
	if err != nil {
		return err
	}
	e.msgMu.Lock()
	for _, m := range msgs {
		if !e.bufferStillValid(m, now) {
			_ = e.runtimeStore.DeleteMessage(m.ID)
			continue
		}
		e.msgBuf = append(e.msgBuf, fromRuntimeMessage(m))
	}
	e.msgMu.Unlock()
	return nil
}

func (e *Engine) persistLease(instanceID, tokenID, workerID string, deadline time.Time) {
	if e.runtimeStore == nil {
		return
	}
	_ = e.runtimeStore.PutLease(runtime.JobLease{
		InstanceID:     instanceID,
		TokenID:        tokenID,
		WorkerID:       workerID,
		DeadlineUnixMs: deadline.UnixMilli(),
	})
}

func (e *Engine) persistReleaseLease(instanceID, tokenID string) {
	if e.runtimeStore == nil {
		return
	}
	_ = e.runtimeStore.DeleteLease(instanceID, tokenID)
}

func (e *Engine) persistEnqueueBuffered(m bufferedMessage) {
	if e.runtimeStore == nil {
		return
	}
	_ = e.runtimeStore.EnqueueMessage(toRuntimeMessage(m))
}

func (e *Engine) persistDeleteBuffered(id string) {
	if e.runtimeStore == nil {
		return
	}
	_ = e.runtimeStore.DeleteMessage(id)
}

func (e *Engine) sweepRuntimeForInstance(instanceID string) {
	e.jobMu.Lock()
	for key := range e.leases {
		if leaseInstanceID(key) == instanceID {
			delete(e.leases, key)
		}
	}
	e.jobMu.Unlock()

	e.msgMu.Lock()
	out := e.msgBuf[:0]
	for _, m := range e.msgBuf {
		if m.instanceID != instanceID {
			out = append(out, m)
			continue
		}
		e.persistDeleteBuffered(m.id)
	}
	e.msgBuf = out
	e.msgMu.Unlock()

	if e.runtimeStore == nil {
		return
	}
	leases, _ := e.runtimeStore.LoadLeases()
	for _, l := range leases {
		if l.InstanceID == instanceID {
			_ = e.runtimeStore.DeleteLease(l.InstanceID, l.TokenID)
		}
	}
	msgs, _ := e.runtimeStore.LoadMessages()
	for _, m := range msgs {
		if m.InstanceID == instanceID {
			_ = e.runtimeStore.DeleteMessage(m.ID)
		}
	}
}

func (e *Engine) leaseStillValid(l runtime.JobLease, now time.Time) bool {
	if now.UnixMilli() >= l.DeadlineUnixMs {
		return false
	}
	e.mu.Lock()
	inst := e.instances[l.InstanceID]
	e.mu.Unlock()
	if inst == nil {
		return false
	}
	tok := inst.Tokens[l.TokenID]
	return tok != nil && tok.Status == projection.TokenWaiting && tok.JobType != ""
}

func (e *Engine) bufferStillValid(m runtime.BufferedMessage, now time.Time) bool {
	if m.EnqueuedUnixMs > 0 && now.Sub(time.UnixMilli(m.EnqueuedUnixMs)) > defaultMessageBufferTTL {
		return false
	}
	if m.InstanceID == "" {
		return true
	}
	e.mu.Lock()
	inst := e.instances[m.InstanceID]
	e.mu.Unlock()
	if inst == nil {
		return false
	}
	return inst.Status != projection.StatusCompleted && inst.Status != projection.StatusTerminated
}

func leaseInstanceID(key string) string {
	for i := 0; i < len(key); i++ {
		if key[i] == '/' {
			return key[:i]
		}
	}
	return key
}

func toRuntimeMessage(m bufferedMessage) runtime.BufferedMessage {
	out := runtime.BufferedMessage{
		ID:             m.id,
		Name:           m.name,
		InstanceID:     m.instanceID,
		EnqueuedUnixMs: m.enqueuedUnixMs,
	}
	for _, k := range m.keys {
		out.Keys = append(out.Keys, runtime.CorrelationKey{
			Name:      k.GetName(),
			JsonValue: k.GetJsonValue(),
		})
	}
	if len(m.vars) > 0 {
		out.Vars = make(map[string]string, len(m.vars))
		for k, v := range m.vars {
			b, err := json.Marshal(v)
			if err != nil {
				continue
			}
			out.Vars[k] = string(b)
		}
	}
	return out
}

func fromRuntimeMessage(m runtime.BufferedMessage) bufferedMessage {
	out := bufferedMessage{
		id:             m.ID,
		name:           m.Name,
		instanceID:     m.InstanceID,
		enqueuedUnixMs: m.EnqueuedUnixMs,
	}
	for _, k := range m.Keys {
		out.keys = append(out.keys, &eventv1.Variable{Name: k.Name, JsonValue: k.JsonValue})
	}
	if len(m.Vars) > 0 {
		out.vars = make(map[string]any, len(m.Vars))
		for k, raw := range m.Vars {
			var v any
			if err := json.Unmarshal([]byte(raw), &v); err != nil {
				out.vars[k] = raw
				continue
			}
			out.vars[k] = v
		}
	}
	return out
}
