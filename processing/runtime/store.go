package runtime

import (
	"fmt"
	"sync"
)

// JobLease tracks an in-flight ServiceTask claim. Not an EventLog fact.
type JobLease struct {
	InstanceID     string `json:"instance_id"`
	TokenID        string `json:"token_id"`
	WorkerID       string `json:"worker_id"`
	DeadlineUnixMs int64  `json:"deadline_unix_ms"`
}

// CorrelationKey mirrors instance variable JSON for message correlation.
type CorrelationKey struct {
	Name      string `json:"name"`
	JsonValue string `json:"json_value"`
}

// BufferedMessage is a published message waiting for a catch or message boundary.
type BufferedMessage struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	InstanceID     string            `json:"instance_id,omitempty"`
	Keys           []CorrelationKey  `json:"keys,omitempty"`
	Vars           map[string]string `json:"vars,omitempty"` // name -> json_value
	EnqueuedUnixMs int64             `json:"enqueued_unix_ms"`
}

// Store persists auxiliary runtime state (job leases, message buffer).
// It is optional: nil means in-memory only for the process lifetime.
type Store interface {
	PutLease(lease JobLease) error
	DeleteLease(instanceID, tokenID string) error
	LoadLeases() ([]JobLease, error)

	EnqueueMessage(msg BufferedMessage) error
	DeleteMessage(id string) error
	LoadMessages() ([]BufferedMessage, error)
}

func LeaseKey(instanceID, tokenID string) string {
	return instanceID + "/" + tokenID
}

// MemoryStore keeps runtime state in process memory (tests / custom Recover).
type MemoryStore struct {
	mu       sync.Mutex
	leases   map[string]JobLease
	messages map[string]BufferedMessage
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		leases:   make(map[string]JobLease),
		messages: make(map[string]BufferedMessage),
	}
}

func (s *MemoryStore) PutLease(lease JobLease) error {
	if lease.InstanceID == "" || lease.TokenID == "" {
		return fmt.Errorf("lease instance_id and token_id are required")
	}
	s.mu.Lock()
	s.leases[LeaseKey(lease.InstanceID, lease.TokenID)] = lease
	s.mu.Unlock()
	return nil
}

func (s *MemoryStore) DeleteLease(instanceID, tokenID string) error {
	s.mu.Lock()
	delete(s.leases, LeaseKey(instanceID, tokenID))
	s.mu.Unlock()
	return nil
}

func (s *MemoryStore) LoadLeases() ([]JobLease, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]JobLease, 0, len(s.leases))
	for _, l := range s.leases {
		out = append(out, l)
	}
	return out, nil
}

func (s *MemoryStore) EnqueueMessage(msg BufferedMessage) error {
	if msg.ID == "" {
		return fmt.Errorf("buffered message id is required")
	}
	s.mu.Lock()
	s.messages[msg.ID] = msg
	s.mu.Unlock()
	return nil
}

func (s *MemoryStore) DeleteMessage(id string) error {
	s.mu.Lock()
	delete(s.messages, id)
	s.mu.Unlock()
	return nil
}

func (s *MemoryStore) LoadMessages() ([]BufferedMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]BufferedMessage, 0, len(s.messages))
	for _, m := range s.messages {
		out = append(out, m)
	}
	return out, nil
}
