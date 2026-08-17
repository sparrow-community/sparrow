package event_test

import (
	"sync"
	"testing"

	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
	"google.golang.org/protobuf/proto"
)

var eventPool = sync.Pool{
	New: func() interface{} {
		return &eventv1.Event{}
	},
}

const (
	testRecordID   = "01900000-0000-7000-8000-000000000001"
	testDeployID   = "01900000-0000-7000-8000-000000000002"
	testInstanceID = "01900000-0000-7000-8000-000000000003"
	testSourceID   = "01900000-0000-7000-8000-00000000000a"
	testTokenID    = "01900000-0000-7000-8000-000000000064"
)

func TestBuildStartEventEvent(t *testing.T) {
	e, err := BuildStartEventEvent()
	if err != nil {
		t.Fatalf("Error building StartEvent: %v", err)
	}
	evt := &eventv1.Event{}
	if err := proto.Unmarshal(e, evt); err != nil {
		t.Fatalf("Error unmarshaling StartEvent: %v", err)
	}

	if evt.GetId() != testRecordID {
		t.Fatalf("id = %q, want %q", evt.GetId(), testRecordID)
	}
	if evt.GetRecordType() != eventv1.Event_RECORD_TYPE_EVENT {
		t.Fatalf("record_type = %v", evt.GetRecordType())
	}
	if evt.GetSourceRecordId() != testSourceID {
		t.Fatalf("source_record_id = %q, want %q", evt.GetSourceRecordId(), testSourceID)
	}
	if evt.GetElement().GetType() != eventv1.Element_TYPE_START_EVENT {
		t.Fatalf("element.type = %v", evt.GetElement().GetType())
	}
}

func TestBuildUserTaskRejection(t *testing.T) {
	msg := &eventv1.Event{
		Id:                "01900000-0000-7000-8000-000000000002",
		Timestamp:         1234567890,
		RecordType:        eventv1.Event_RECORD_TYPE_REJECTION,
		DeploymentId:      testDeployID,
		ProcessInstanceId: testInstanceID,
		ProcessVersion:    1,
		SourceRecordId:    testRecordID,
		Rejection: &eventv1.Rejection{
			Code:    "INVALID_STATE",
			Message: "task is not activated",
		},
		Element: &eventv1.Element{
			Intent:  eventv1.Element_INTENT_COMPLETING,
			Type:    eventv1.Element_TYPE_USER_TASK,
			Id:      "task-1",
			TokenId: testTokenID,
			Payload: &eventv1.Element_ActivityPayload{
				ActivityPayload: &eventv1.ActivityPayload{},
			},
		},
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := &eventv1.Event{}
	if err := proto.Unmarshal(data, got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.GetElement().GetType() != eventv1.Element_TYPE_USER_TASK {
		t.Fatalf("type = %v", got.GetElement().GetType())
	}
	if got.GetRejection().GetCode() != "INVALID_STATE" {
		t.Fatalf("rejection.code = %q", got.GetRejection().GetCode())
	}
}

func BenchmarkBuildStartEventEvent(b *testing.B) {
	evt := eventPool.Get().(*eventv1.Event)
	defer eventPool.Put(evt)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		evt.Reset()
		evt.Id = testRecordID
		evt.DeploymentId = testDeployID
		evt.ProcessInstanceId = testInstanceID
		evt.ProcessVersion = 1
		evt.Timestamp = 1234567890
		evt.RecordType = eventv1.Event_RECORD_TYPE_COMMAND
		evt.SourceRecordId = ""
		evt.Element = &eventv1.Element{
			Intent:  eventv1.Element_INTENT_ACTIVATING,
			Type:    eventv1.Element_TYPE_START_EVENT,
			Id:      "start-event-1",
			TokenId: testTokenID,
			Payload: &eventv1.Element_EventPayload{
				EventPayload: &eventv1.EventPayload{},
			},
		}

		if _, err := proto.Marshal(evt); err != nil {
			b.Errorf("Error marshaling: %v", err)
		}
	}
}

func BenchmarkParseEventFromBinary(b *testing.B) {
	data, err := BuildStartEventEvent()
	if err != nil {
		b.Fatalf("Error building StartEvent: %v", err)
	}

	evt := eventPool.Get().(*eventv1.Event)
	defer eventPool.Put(evt)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		evt.Reset()
		if err := proto.Unmarshal(data, evt); err != nil {
			b.Errorf("Error unmarshaling StartEvent: %v", err)
		}
		_ = evt.Id
	}
}

func BuildStartEventEvent() ([]byte, error) {
	eventMsg := &eventv1.Event{
		Id:                testRecordID,
		DeploymentId:      testDeployID,
		ProcessInstanceId: testInstanceID,
		ProcessVersion:    1,
		Timestamp:         1234567890,
		RecordType:        eventv1.Event_RECORD_TYPE_EVENT,
		SourceRecordId:    testSourceID,
		Element: &eventv1.Element{
			Intent:  eventv1.Element_INTENT_ACTIVATING,
			Type:    eventv1.Element_TYPE_START_EVENT,
			Id:      "start-event-1",
			TokenId: testTokenID,
			Payload: &eventv1.Element_EventPayload{
				EventPayload: &eventv1.EventPayload{},
			},
		},
	}

	return proto.Marshal(eventMsg)
}
