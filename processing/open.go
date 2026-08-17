package processing

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sparrow-community/sparrow/processing/deploy"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// Open opens a durable engine under dataDir:
//
//	dataDir/events.log          append-only event ledger
//	dataDir/deployments/<id>.bpmn
//
// It reloads deployments and rebuilds instance projections by replaying EVENT records.
func Open(dataDir string) (*Engine, error) {
	if dataDir == "" {
		return nil, fmt.Errorf("dataDir is required")
	}
	depDir := filepath.Join(dataDir, "deployments")
	if err := os.MkdirAll(depDir, 0o755); err != nil {
		return nil, fmt.Errorf("mkdir deployments: %w", err)
	}
	fl, err := eventlog.OpenFile(filepath.Join(dataDir, "events.log"))
	if err != nil {
		return nil, err
	}
	e := NewEngine(fl)
	e.dataDir = dataDir
	if err := e.loadDeployments(); err != nil {
		_ = fl.Close()
		return nil, err
	}
	if err := e.recover(context.Background()); err != nil {
		_ = fl.Close()
		return nil, err
	}
	return e, nil
}

// Close releases durable resources (file event log). Safe on memory engines.
func (e *Engine) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if c, ok := e.log.(interface{ Close() error }); ok {
		return c.Close()
	}
	return nil
}

func (e *Engine) loadDeployments() error {
	dir := filepath.Join(e.dataDir, "deployments")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, ent := range entries {
		if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".bpmn") {
			continue
		}
		id := strings.TrimSuffix(ent.Name(), ".bpmn")
		xml, err := os.ReadFile(filepath.Join(dir, ent.Name()))
		if err != nil {
			return err
		}
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
	if e.dataDir == "" {
		return nil
	}
	path := filepath.Join(e.dataDir, "deployments", id+".bpmn")
	return os.WriteFile(path, bpmnXML, 0o644)
}

func (e *Engine) recover(ctx context.Context) error {
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
