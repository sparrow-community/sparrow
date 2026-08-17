package deploy

import (
	"fmt"
	"sync"
)

// Store persists BPMN definition bytes keyed by deployment id.
// Memory and directory implementations ship in-tree; callers can plug in their own
// (database, object store, etc.).
type Store interface {
	Put(id string, bpmnXML []byte) error
	LoadAll() (map[string][]byte, error)
}

// MemoryStore keeps definitions in process memory.
type MemoryStore struct {
	mu   sync.Mutex
	docs map[string][]byte
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{docs: make(map[string][]byte)}
}

func (s *MemoryStore) Put(id string, bpmnXML []byte) error {
	if id == "" {
		return fmt.Errorf("deployment id is empty")
	}
	cp := make([]byte, len(bpmnXML))
	copy(cp, bpmnXML)
	s.mu.Lock()
	s.docs[id] = cp
	s.mu.Unlock()
	return nil
}

func (s *MemoryStore) LoadAll() (map[string][]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make(map[string][]byte, len(s.docs))
	for id, xml := range s.docs {
		cp := make([]byte, len(xml))
		copy(cp, xml)
		out[id] = cp
	}
	return out, nil
}
