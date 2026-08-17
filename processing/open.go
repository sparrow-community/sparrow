package processing

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/sparrow-community/sparrow/processing/deploy"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// Recover builds an Engine from caller-supplied persistence:
// EventLog for the behavior ledger, Store for BPMN definitions.
// Both may be in-memory, files, or a custom implementation.
func Recover(ctx context.Context, l eventlog.EventLog, store deploy.Store) (*Engine, error) {
	if l == nil {
		return nil, fmt.Errorf("event log is required")
	}
	e := NewEngine(l)
	e.store = store
	if err := e.loadDeployments(); err != nil {
		return nil, err
	}
	if err := e.replay(ctx); err != nil {
		return nil, err
	}
	return e, nil
}

// Open is a file-backed Recover:
//
//	dataDir/events.log
//	dataDir/deployments/<id>.bpmn
func Open(ctx context.Context, dataDir string) (*Engine, error) {
	if dataDir == "" {
		return nil, fmt.Errorf("dataDir is required")
	}
	fl, err := eventlog.OpenFile(filepath.Join(dataDir, "events.log"))
	if err != nil {
		return nil, err
	}
	ds, err := deploy.OpenDirStore(filepath.Join(dataDir, "deployments"))
	if err != nil {
		_ = fl.Close()
		return nil, err
	}
	e, err := Recover(ctx, fl, ds)
	if err != nil {
		_ = fl.Close()
		return nil, err
	}
	return e, nil
}

// Close releases resources if the EventLog implements Close. Safe for memory backends.
func (e *Engine) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if c, ok := e.log.(interface{ Close() error }); ok {
		return c.Close()
	}
	return nil
}

func (e *Engine) loadDeployments() error {
	if e.store == nil {
		return nil
	}
	docs, err := e.store.LoadAll()
	if err != nil {
		return fmt.Errorf("load deployments: %w", err)
	}
	for id, xml := range docs {
		dep, err := deploy.Compile(xml)
		if err != nil {
			return fmt.Errorf("load deployment %q: %w", id, err)
		}
		dep.ID = id
		e.deployments[id] = dep
	}
	return nil
}

func (e *Engine) persistDeployment(id string, bpmnXML []byte) error {
	if e.store == nil {
		return nil
	}
	return e.store.Put(id, bpmnXML)
}

func (e *Engine) replay(ctx context.Context) error {
	events, err := e.log.ReadAll(ctx)
	if err != nil {
		return fmt.Errorf("read event log: %w", err)
	}
	for _, ev := range events {
		iid := ev.GetProcessInstanceId()
		if iid == "" {
			continue
		}
		e.mu.Lock()
		inst, ok := e.instances[iid]
		if !ok {
			inst = projection.NewInstance(iid, ev.GetDeploymentId(), ev.GetProcessVersion())
			e.instances[iid] = inst
			e.instMu[iid] = &sync.Mutex{}
		}
		e.mu.Unlock()
		if ev.GetRecordType() == eventv1.Event_RECORD_TYPE_EVENT {
			inst.ApplyEvent(ev)
		}
	}
	return nil
}
