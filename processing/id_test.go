package processing

import (
	"testing"

	"github.com/google/uuid"
)

func TestNextID(t *testing.T) {
	id, err := NextID()
	if err != nil {
		t.Fatalf("NextID: %v", err)
	}
	parsed, err := uuid.Parse(id)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if parsed.Version() != 7 {
		t.Fatalf("version = %d, want 7", parsed.Version())
	}
}

func TestMustNextIDUnique(t *testing.T) {
	a := MustNextID()
	b := MustNextID()
	if a == b {
		t.Fatalf("expected distinct ids, both %q", a)
	}
}
