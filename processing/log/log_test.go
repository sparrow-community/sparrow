package log

import (
	"context"
	"testing"

	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func TestMemoryAppendAndReadByInstance(t *testing.T) {
	mem := NewMemory()
	ctx := context.Background()

	_, err := mem.Append(ctx, &eventv1.Event{
		Id:                "e1",
		ProcessInstanceId: "i1",
		RecordType:        eventv1.Event_RECORD_TYPE_EVENT,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = mem.Append(ctx, &eventv1.Event{
		Id:                "e2",
		ProcessInstanceId: "i2",
		RecordType:        eventv1.Event_RECORD_TYPE_EVENT,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = mem.Append(ctx, &eventv1.Event{
		Id:                "e3",
		ProcessInstanceId: "i1",
		RecordType:        eventv1.Event_RECORD_TYPE_COMMAND,
	})
	if err != nil {
		t.Fatal(err)
	}

	got, err := mem.ReadByInstance(ctx, "i1")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("len=%d want 2", len(got))
	}
	if got[0].GetId() != "e1" || got[1].GetId() != "e3" {
		t.Fatalf("unexpected order: %#v %#v", got[0].GetId(), got[1].GetId())
	}

	// Mutation isolation
	got[0].Id = "mutated"
	again, _ := mem.ReadByInstance(ctx, "i1")
	if again[0].GetId() != "e1" {
		t.Fatalf("log record mutated via returned slice")
	}
}
