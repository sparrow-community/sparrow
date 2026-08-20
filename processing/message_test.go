package processing_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func TestMessageCatchWaitAndComplete(t *testing.T) {
	xml := readTestdataMessage(t, "m2_message_catch.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()

	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}

	inst, ok := eng.GetInstance(instanceID)
	if !ok || inst.Status != projection.StatusActive {
		t.Fatalf("inst=%v ok=%v", inst, ok)
	}

	elementID, tokenID := waitingAt(inst)
	if elementID != "MessageCatch_1" {
		t.Fatalf("expected wait at MessageCatch_1, tokens=%#v", inst.Tokens)
	}
	tok := inst.Tokens[tokenID]
	if tok == nil || tok.Status != projection.TokenWaiting {
		t.Fatalf("token not waiting: %#v", tok)
	}
	if tok.DueUnixMs != 0 {
		t.Fatalf("message catch due should be 0, due=%d", tok.DueUnixMs)
	}
	if tok.MessageName != "order.confirmed" {
		t.Fatalf("message name=%q", tok.MessageName)
	}

	// Ensure ACTIVATED event for the message catch carries message_name.
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	var sawActivated bool
	for _, ev := range events {
		if ev.GetRecordType() != eventv1.Event_RECORD_TYPE_EVENT {
			continue
		}
		el := ev.GetElement()
		if el == nil {
			continue
		}
		if el.GetType() == eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT &&
			el.GetIntent() == eventv1.Element_INTENT_ACTIVATED &&
			el.GetId() == "MessageCatch_1" {
			sawActivated = true
			p := el.GetEventPayload()
			if p == nil || p.GetMessageName() != "order.confirmed" {
				t.Fatalf("event payload=%v", p)
			}
			if p.GetDueUnixMs() != 0 || p.GetDuration() != "" {
				t.Fatalf("unexpected timer fields on message payload: due=%d duration=%q", p.GetDueUnixMs(), p.GetDuration())
			}
		}
	}
	if !sawActivated {
		t.Fatal("missing INTERMEDIATE_CATCH_EVENT ACTIVATED for MessageCatch_1")
	}

	if err := eng.Complete(ctx, instanceID, elementID, tokenID, nil); err != nil {
		t.Fatalf("Complete: %v", err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestMessageCatchRecoverThenComplete(t *testing.T) {
	xml := readTestdataMessage(t, "m2_message_catch.bpmn")
	dir := t.TempDir()
	ctx := context.Background()

	eng1, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	defer eng1.Close()

	dep, err := eng1.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng1.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}

	inst, _ := eng1.GetInstance(instanceID)
	elementID, tokenID := waitingAt(inst)
	if elementID != "MessageCatch_1" {
		t.Fatalf("tokens=%#v", inst.Tokens)
	}
	if inst.Tokens[tokenID].DueUnixMs != 0 {
		t.Fatalf("due should be 0, got %d", inst.Tokens[tokenID].DueUnixMs)
	}
	if inst.Tokens[tokenID].MessageName != "order.confirmed" {
		t.Fatalf("message name=%q", inst.Tokens[tokenID].MessageName)
	}

	if err := eng1.Close(); err != nil {
		t.Fatal(err)
	}

	eng2, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	defer eng2.Close()

	inst2, ok := eng2.GetInstance(instanceID)
	if !ok {
		t.Fatal("instance missing after recover")
	}
	elementID2, tokenID2 := waitingAt(inst2)
	if elementID2 != "MessageCatch_1" {
		t.Fatalf("tokens=%#v", inst2.Tokens)
	}
	if inst2.Tokens[tokenID2].DueUnixMs != 0 {
		t.Fatalf("due after recover should be 0, got %d", inst2.Tokens[tokenID2].DueUnixMs)
	}
	if inst2.Tokens[tokenID2].MessageName != "order.confirmed" {
		t.Fatalf("message name after recover=%q", inst2.Tokens[tokenID2].MessageName)
	}

	if err := eng2.Complete(ctx, instanceID, elementID2, tokenID2, nil); err != nil {
		t.Fatalf("Complete: %v", err)
	}
	inst2, _ = eng2.GetInstance(instanceID)
	if inst2.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst2.Status)
	}
}

func TestPublishMessageCompletesCatch(t *testing.T) {
	xml := readTestdataMessage(t, "m2_message_catch.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}

	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{
		Name:      "order.confirmed",
		Variables: map[string]any{"payload": "ok"},
	})
	if err != nil {
		t.Fatalf("PublishMessage: %v", err)
	}
	if n != 1 {
		t.Fatalf("delivered=%d want 1", n)
	}
	inst, _ := eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
	if inst.Variables["payload"] == "" {
		t.Fatalf("vars=%#v", inst.Variables)
	}
}

func TestPublishMessageNotFoundAndScope(t *testing.T) {
	xml := readTestdataMessage(t, "m2_message_catch.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{}); err == nil {
		t.Fatal("expected INVALID_ARGUMENT")
	}
	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "unknown"})
	if err != nil || n != 0 {
		t.Fatalf("late unknown name should buffer n=%d err=%v", n, err)
	}
	if _, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{
		Name:              "order.confirmed",
		ProcessInstanceID: "missing",
	}); err == nil {
		t.Fatal("expected NOT_FOUND for missing instance")
	}

	n, err = eng.PublishMessage(ctx, processing.PublishMessageRequest{
		Name:              "order.confirmed",
		ProcessInstanceID: instanceID,
	})
	if err != nil || n != 1 {
		t.Fatalf("scoped publish n=%d err=%v", n, err)
	}
}

func TestPublishMessageCorrelationKeys(t *testing.T) {
	xml := readTestdataMessage(t, "m2_message_catch.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	a, err := eng.CreateInstance(ctx, dep, map[string]any{"orderId": "A"})
	if err != nil {
		t.Fatal(err)
	}
	b, err := eng.CreateInstance(ctx, dep, map[string]any{"orderId": "B"})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{
		Name:            "order.confirmed",
		CorrelationKeys: map[string]any{"orderId": "missing"},
	}); err != nil {
		t.Fatalf("unmatched keys should buffer: %v", err)
	}
	instA, _ := eng.GetInstance(a)
	instB, _ := eng.GetInstance(b)
	if instA.Status != projection.StatusActive || instB.Status != projection.StatusActive {
		t.Fatalf("A=%s B=%s want both still waiting", instA.Status, instB.Status)
	}

	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{
		Name:            "order.confirmed",
		CorrelationKeys: map[string]any{"orderId": "A"},
	})
	if err != nil || n != 1 {
		t.Fatalf("correlated A n=%d err=%v", n, err)
	}
	instA, _ = eng.GetInstance(a)
	instB, _ = eng.GetInstance(b)
	if instA.Status != projection.StatusCompleted {
		t.Fatalf("A status=%s want completed", instA.Status)
	}
	if instB.Status != projection.StatusActive {
		t.Fatalf("B status=%s want still active", instB.Status)
	}

	n, err = eng.PublishMessage(ctx, processing.PublishMessageRequest{
		Name:            "order.confirmed",
		CorrelationKeys: map[string]any{"orderId": "B"},
	})
	if err != nil || n != 1 {
		t.Fatalf("correlated B n=%d err=%v", n, err)
	}
	instB, _ = eng.GetInstance(b)
	if instB.Status != projection.StatusCompleted {
		t.Fatalf("B status=%s want completed", instB.Status)
	}
}

func TestPublishMessageWithoutKeysCompletesAll(t *testing.T) {
	xml := readTestdataMessage(t, "m2_message_catch.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := eng.CreateInstance(ctx, dep, map[string]any{"orderId": "A"}); err != nil {
		t.Fatal(err)
	}
	if _, err := eng.CreateInstance(ctx, dep, map[string]any{"orderId": "B"}); err != nil {
		t.Fatal(err)
	}
	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "order.confirmed"})
	if err != nil || n != 2 {
		t.Fatalf("broadcast n=%d err=%v", n, err)
	}
}

func TestPublishMessageCorrelationAfterRecover(t *testing.T) {
	xml := readTestdataMessage(t, "m2_message_catch.bpmn")
	dir := t.TempDir()
	ctx := context.Background()

	eng1, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	dep, err := eng1.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	a, err := eng1.CreateInstance(ctx, dep, map[string]any{"orderId": "A"})
	if err != nil {
		t.Fatal(err)
	}
	b, err := eng1.CreateInstance(ctx, dep, map[string]any{"orderId": "B"})
	if err != nil {
		t.Fatal(err)
	}
	if err := eng1.Close(); err != nil {
		t.Fatal(err)
	}

	eng2, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	defer eng2.Close()

	n, err := eng2.PublishMessage(ctx, processing.PublishMessageRequest{
		Name:            "order.confirmed",
		CorrelationKeys: map[string]any{"orderId": "A"},
	})
	if err != nil || n != 1 {
		t.Fatalf("n=%d err=%v", n, err)
	}
	instA, _ := eng2.GetInstance(a)
	instB, _ := eng2.GetInstance(b)
	if instA.Status != projection.StatusCompleted {
		t.Fatalf("A status=%s", instA.Status)
	}
	if instB.Status != projection.StatusActive {
		t.Fatalf("B status=%s want still waiting", instB.Status)
	}
}

func TestPublishMessageIgnoresTimerCatch(t *testing.T) {
	xml := readTestdata(t, "m2_timer_catch_1h.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "order.confirmed"})
	if err != nil || n != 0 {
		t.Fatalf("timer catch must not consume message n=%d err=%v", n, err)
	}
	inst, _ := eng.GetInstance(instanceID)
	if waitingElement(inst) != "TimerCatch_1" {
		t.Fatalf("tokens=%#v", inst.Tokens)
	}
}

func TestPublishMessageAfterRecover(t *testing.T) {
	xml := readTestdataMessage(t, "m2_message_catch.bpmn")
	dir := t.TempDir()
	ctx := context.Background()

	eng1, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	dep, err := eng1.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng1.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := eng1.Close(); err != nil {
		t.Fatal(err)
	}

	eng2, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	defer eng2.Close()

	n, err := eng2.PublishMessage(ctx, processing.PublishMessageRequest{Name: "order.confirmed"})
	if err != nil || n != 1 {
		t.Fatalf("n=%d err=%v", n, err)
	}
	inst, _ := eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed", inst.Status)
	}
}

func TestPublishMessageBuffersUntilCatch(t *testing.T) {
	xml := readTestdataMessage(t, "m2_message_catch.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{
		Name:      "order.confirmed",
		Variables: map[string]any{"payload": "ok"},
	})
	if err != nil || n != 0 {
		t.Fatalf("buffer n=%d err=%v", n, err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed from buffer", inst.Status)
	}
	if inst.Variables["payload"] == "" {
		t.Fatalf("vars=%#v", inst.Variables)
	}
}

func TestPublishMessageBufferHonorsCorrelation(t *testing.T) {
	xml := readTestdataMessage(t, "m2_message_catch.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{
		Name:            "order.confirmed",
		CorrelationKeys: map[string]any{"orderId": "A"},
	}); err != nil {
		t.Fatal(err)
	}
	b, err := eng.CreateInstance(ctx, dep, map[string]any{"orderId": "B"})
	if err != nil {
		t.Fatal(err)
	}
	instB, _ := eng.GetInstance(b)
	if instB.Status != projection.StatusActive {
		t.Fatalf("B status=%s want still waiting", instB.Status)
	}
	a, err := eng.CreateInstance(ctx, dep, map[string]any{"orderId": "A"})
	if err != nil {
		t.Fatal(err)
	}
	instA, _ := eng.GetInstance(a)
	instB, _ = eng.GetInstance(b)
	if instA.Status != projection.StatusCompleted {
		t.Fatalf("A status=%s want completed from buffer", instA.Status)
	}
	if instB.Status != projection.StatusActive {
		t.Fatalf("B status=%s want still waiting", instB.Status)
	}
}

func TestPublishMessageBufferAfterUserTask(t *testing.T) {
	xml := readTestdataMessage(t, "m2_message_after_task.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	elementID, tokenID := waitingAt(inst)
	if elementID != "UserTask_1" {
		t.Fatalf("tokens=%#v", inst.Tokens)
	}
	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{
		Name:              "order.confirmed",
		ProcessInstanceID: instanceID,
		Variables:         map[string]any{"payload": "late"},
	})
	if err != nil || n != 0 {
		t.Fatalf("buffer while on task n=%d err=%v", n, err)
	}
	if err := eng.Complete(ctx, instanceID, elementID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed after task then buffered message", inst.Status)
	}
	if inst.Variables["payload"] == "" {
		t.Fatalf("vars=%#v", inst.Variables)
	}
}

func TestPublishMessageBufferSurvivesRecover(t *testing.T) {
	xml := readTestdataMessage(t, "m2_message_catch.bpmn")
	dir := t.TempDir()
	ctx := context.Background()

	eng1, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	dep, err := eng1.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := eng1.PublishMessage(ctx, processing.PublishMessageRequest{Name: "order.confirmed"}); err != nil {
		t.Fatal(err)
	}
	if err := eng1.Close(); err != nil {
		t.Fatal(err)
	}

	eng2, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	defer eng2.Close()

	instanceID, err := eng2.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s want completed from persisted buffer", inst.Status)
	}
}

func readTestdataMessage(t *testing.T, name string) []byte {
	t.Helper()
	xml, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return xml
}
