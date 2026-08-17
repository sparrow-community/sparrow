package processing

import (
	"fmt"

	"github.com/google/uuid"
)

// NextID returns a new UUIDv7 in canonical string form.
func NextID() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("generate uuid v7: %w", err)
	}
	return id.String(), nil
}

// MustNextID is like NextID but panics on failure.
func MustNextID() string {
	id, err := NextID()
	if err != nil {
		panic(err)
	}
	return id
}
