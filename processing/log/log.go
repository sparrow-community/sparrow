package log

import (
	"context"
	"fmt"
	"sync"

	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
	"google.golang.org/protobuf/proto"
)

// EventLog is the append-only behavior ledger.
type EventLog interface {
	Append(ctx context.Context, e *eventv1.Event) (position int64, err error)
	ReadByInstance(ctx context.Context, processInstanceID string) ([]*eventv1.Event, error)
	ReadAll(ctx context.Context) ([]*eventv1.Event, error)
}

// Memory is an in-memory EventLog suitable for tests and M1.
type Memory struct {
	mu      sync.Mutex
	records []*eventv1.Event
}

func NewMemory() *Memory {
	return &Memory{records: make([]*eventv1.Event, 0, 64)}
}

func (m *Memory) Append(_ context.Context, e *eventv1.Event) (int64, error) {
	if e == nil {
		return 0, fmt.Errorf("event is nil")
	}
	clone, ok := proto.Clone(e).(*eventv1.Event)
	if !ok || clone == nil {
		return 0, fmt.Errorf("clone event")
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.records = append(m.records, clone)
	return int64(len(m.records)), nil
}

func (m *Memory) ReadByInstance(_ context.Context, processInstanceID string) ([]*eventv1.Event, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	out := make([]*eventv1.Event, 0)
	for _, e := range m.records {
		if e.GetProcessInstanceId() == processInstanceID {
			out = append(out, proto.Clone(e).(*eventv1.Event))
		}
	}
	return out, nil
}

func (m *Memory) ReadAll(_ context.Context) ([]*eventv1.Event, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	out := make([]*eventv1.Event, 0, len(m.records))
	for _, e := range m.records {
		out = append(out, proto.Clone(e).(*eventv1.Event))
	}
	return out, nil
}
