package processing

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	InitSonyflake(
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		func() (int, error) {
			return 1, nil
		},
	)
	code := m.Run()
	os.Exit(code)
}

func TestNextID(t *testing.T) {
	id := NextID()
	t.Log("sonyflake id:", id)
	t.Log("decompose:", sf.Decompose(id))

	fmt.Println("sonyflake id:", id)
	fmt.Println("decompose:", sf.Decompose(id))
}
