package deploy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DirStore stores each deployment as <dir>/<id>.bpmn.
type DirStore struct {
	dir string
}

func OpenDirStore(dir string) (*DirStore, error) {
	if dir == "" {
		return nil, fmt.Errorf("deployment store dir is required")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("mkdir deployments: %w", err)
	}
	return &DirStore{dir: dir}, nil
}

func (s *DirStore) Put(id string, bpmnXML []byte) error {
	if id == "" {
		return fmt.Errorf("deployment id is empty")
	}
	return os.WriteFile(filepath.Join(s.dir, id+".bpmn"), bpmnXML, 0o644)
}

func (s *DirStore) LoadAll() (map[string][]byte, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}
	out := make(map[string][]byte)
	for _, ent := range entries {
		if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".bpmn") {
			continue
		}
		id := strings.TrimSuffix(ent.Name(), ".bpmn")
		xml, err := os.ReadFile(filepath.Join(s.dir, ent.Name()))
		if err != nil {
			return nil, err
		}
		out[id] = xml
	}
	return out, nil
}
