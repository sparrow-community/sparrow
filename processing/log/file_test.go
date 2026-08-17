package log

import (
	"context"
	"path/filepath"
	"testing"

	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func TestFileAppendReadRestart(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "events.log")
	ctx := context.Background()

	fl, err := OpenFile(path)
	if err != nil {
		t.Fatal(err)
	}
	_, err = fl.Append(ctx, &eventv1.Event{
		Id:                "e1",
		ProcessInstanceId: "i1",
		RecordType:        eventv1.Event_RECORD_TYPE_EVENT,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = fl.Append(ctx, &eventv1.Event{
		Id:                "e2",
		ProcessInstanceId: "i2",
		RecordType:        eventv1.Event_RECORD_TYPE_COMMAND,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := fl.Close(); err != nil {
		t.Fatal(err)
	}

	fl2, err := OpenFile(path)
	if err != nil {
		t.Fatal(err)
	}
	defer fl2.Close()

	all, err := fl2.ReadAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 || all[0].GetId() != "e1" || all[1].GetId() != "e2" {
		t.Fatalf("all=%v", all)
	}
	byInst, err := fl2.ReadByInstance(ctx, "i1")
	if err != nil {
		t.Fatal(err)
	}
	if len(byInst) != 1 || byInst[0].GetId() != "e1" {
		t.Fatalf("byInst=%v", byInst)
	}

	_, err = fl2.Append(ctx, &eventv1.Event{
		Id:                "e3",
		ProcessInstanceId: "i1",
		RecordType:        eventv1.Event_RECORD_TYPE_EVENT,
	})
	if err != nil {
		t.Fatal(err)
	}
	byInst, _ = fl2.ReadByInstance(ctx, "i1")
	if len(byInst) != 2 {
		t.Fatalf("len=%d want 2 after append", len(byInst))
	}
}
