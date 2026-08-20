package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// FileStore persists runtime state as a single JSON document.
type FileStore struct {
	path string
	mu   sync.Mutex
}

type persistedState struct {
	Leases   []JobLease        `json:"leases,omitempty"`
	Messages []BufferedMessage `json:"messages,omitempty"`
}

func OpenFileStore(dir string) (*FileStore, error) {
	if dir == "" {
		return nil, fmt.Errorf("runtime store dir is required")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("mkdir runtime: %w", err)
	}
	return &FileStore{path: filepath.Join(dir, "state.json")}, nil
}

func (s *FileStore) PutLease(lease JobLease) error {
	if lease.InstanceID == "" || lease.TokenID == "" {
		return fmt.Errorf("lease instance_id and token_id are required")
	}
	return s.update(func(st *persistedState) {
		key := LeaseKey(lease.InstanceID, lease.TokenID)
		found := false
		for i := range st.Leases {
			if LeaseKey(st.Leases[i].InstanceID, st.Leases[i].TokenID) == key {
				st.Leases[i] = lease
				found = true
				break
			}
		}
		if !found {
			st.Leases = append(st.Leases, lease)
		}
	})
}

func (s *FileStore) DeleteLease(instanceID, tokenID string) error {
	key := LeaseKey(instanceID, tokenID)
	return s.update(func(st *persistedState) {
		out := st.Leases[:0]
		for _, l := range st.Leases {
			if LeaseKey(l.InstanceID, l.TokenID) == key {
				continue
			}
			out = append(out, l)
		}
		st.Leases = out
	})
}

func (s *FileStore) LoadLeases() ([]JobLease, error) {
	st, err := s.read()
	if err != nil {
		return nil, err
	}
	return append([]JobLease(nil), st.Leases...), nil
}

func (s *FileStore) EnqueueMessage(msg BufferedMessage) error {
	if msg.ID == "" {
		return fmt.Errorf("buffered message id is required")
	}
	return s.update(func(st *persistedState) {
		found := false
		for i := range st.Messages {
			if st.Messages[i].ID == msg.ID {
				st.Messages[i] = msg
				found = true
				break
			}
		}
		if !found {
			st.Messages = append(st.Messages, msg)
		}
	})
}

func (s *FileStore) DeleteMessage(id string) error {
	return s.update(func(st *persistedState) {
		out := st.Messages[:0]
		for _, m := range st.Messages {
			if m.ID == id {
				continue
			}
			out = append(out, m)
		}
		st.Messages = out
	})
}

func (s *FileStore) LoadMessages() ([]BufferedMessage, error) {
	st, err := s.read()
	if err != nil {
		return nil, err
	}
	return append([]BufferedMessage(nil), st.Messages...), nil
}

func (s *FileStore) update(fn func(*persistedState)) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	st, err := s.readLocked()
	if err != nil {
		return err
	}
	fn(&st)
	return s.writeLocked(st)
}

func (s *FileStore) read() (persistedState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.readLocked()
}

func (s *FileStore) readLocked() (persistedState, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return persistedState{}, nil
		}
		return persistedState{}, err
	}
	var st persistedState
	if len(data) == 0 {
		return persistedState{}, nil
	}
	if err := json.Unmarshal(data, &st); err != nil {
		return persistedState{}, fmt.Errorf("decode runtime state: %w", err)
	}
	return st, nil
}

func (s *FileStore) writeLocked(st persistedState) error {
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
