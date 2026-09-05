package processing_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	"github.com/sparrow-community/sparrow/processing/deploy"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func TestCallActivityCompletesChildThenCaller(t *testing.T) {
	xml := readTestdataCall(t, "m4_call_activity.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()

	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	callerID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	caller := mustInstance(t, eng, callerID)
	hostElem, hostTok := waitingAt(caller)
	if hostElem != "CallActivity_1" {
		t.Fatalf("expected CallActivity_1 host, got %s tokens=%#v", hostElem, caller.Tokens)
	}
	host := caller.Tokens[hostTok]
	if host == nil || !host.ScopeHost || host.CalledProcessInstanceID == "" {
		t.Fatalf("expected CallActivity host with child id, tok=%#v", host)
	}
	childID := host.CalledProcessInstanceID
	child := mustInstance(t, eng, childID)
	if child.ParentProcessInstanceID != callerID {
		t.Fatalf("child parent=%q want %q", child.ParentProcessInstanceID, callerID)
	}
	elemID, tokenID := waitingAt(child)
	if elemID != "Task_called" {
		t.Fatalf("expected Task_called on child, got %s tokens=%#v", elemID, child.Tokens)
	}
	if err := eng.Complete(ctx, childID, elemID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	child = mustInstance(t, eng, childID)
	if child.Status != projection.StatusCompleted {
		t.Fatalf("child status=%s", child.Status)
	}
	caller = mustInstance(t, eng, callerID)
	if caller.Status != projection.StatusCompleted {
		t.Fatalf("caller status=%s tokens=%#v", caller.Status, caller.Tokens)
	}
	events, _ := eng.ListEvents(ctx, callerID)
	if !sawElementIntent(events, eventv1.Element_TYPE_CALL_ACTIVITY, "CallActivity_1", eventv1.Element_INTENT_ACTIVATED) {
		t.Fatal("expected CallActivity ACTIVATED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_CALL_ACTIVITY, "CallActivity_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected CallActivity COMPLETED")
	}
	childEvents, _ := eng.ListEvents(ctx, childID)
	if !sawElementIntent(childEvents, eventv1.Element_TYPE_USER_TASK, "Task_called", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected called task COMPLETED on child")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_caller", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected caller end COMPLETED")
	}
}

func TestCallActivityRecoverMidCall(t *testing.T) {
	xml := readTestdataCall(t, "m5_call_instance.bpmn")
	ctx := context.Background()
	memLog := eventlog.NewMemory()
	store := deploy.NewMemoryStore()
	eng1, err := processing.Recover(ctx, memLog, store, nil)
	if err != nil {
		t.Fatal(err)
	}
	depID, err := eng1.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	callerID, err := eng1.CreateInstance(ctx, depID, nil)
	if err != nil {
		t.Fatal(err)
	}
	caller := mustInstance(t, eng1, callerID)
	_, hostTok := waitingAt(caller)
	childID := caller.Tokens[hostTok].CalledProcessInstanceID
	if childID == "" {
		t.Fatal("expected called process instance id before recover")
	}

	eng2, err := processing.Recover(ctx, memLog, store, nil)
	if err != nil {
		t.Fatal(err)
	}
	caller = mustInstance(t, eng2, callerID)
	_, hostTok = waitingAt(caller)
	host := caller.Tokens[hostTok]
	if host == nil || host.CalledProcessInstanceID != childID {
		t.Fatalf("caller link lost after recover: tok=%#v", host)
	}
	child := mustInstance(t, eng2, childID)
	if child.ParentProcessInstanceID != callerID {
		t.Fatalf("child parent=%q after recover", child.ParentProcessInstanceID)
	}
	elemID, tokenID := waitingAt(child)
	if err := eng2.Complete(ctx, childID, elemID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	child = mustInstance(t, eng2, childID)
	caller = mustInstance(t, eng2, callerID)
	if child.Status != projection.StatusCompleted || caller.Status != projection.StatusCompleted {
		t.Fatalf("after recover complete: child=%s caller=%s", child.Status, caller.Status)
	}
}

func TestCallActivityEventSubProcessInterrupting(t *testing.T) {
	xml := readTestdataCall(t, "m4_call_activity_event_subprocess.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()

	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	callerID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	caller := mustInstance(t, eng, callerID)
	_, hostTok := waitingAt(caller)
	childID := caller.Tokens[hostTok].CalledProcessInstanceID
	child := mustInstance(t, eng, childID)
	if waitingElement(child) != "Task_called" {
		t.Fatalf("expected Task_called, got %s tokens=%#v", waitingElement(child), child.Tokens)
	}
	if _, ok := child.EventSubProcesses["Event_SubProcess_called"]; !ok {
		t.Fatalf("expected event sub-process armed in called process, arms=%#v", child.EventSubProcesses)
	}
	if !sawEventSubProcessStartActivated(mustListEvents(t, eng, childID), "Event_SubProcess_called") {
		t.Fatal("expected event sub-process start ACTIVATED in child event log")
	}

	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "cancel.called", ProcessInstanceID: childID})
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("delivered=%d", n)
	}
	child = mustInstance(t, eng, childID)
	if child.Status != projection.StatusCompleted {
		t.Fatalf("child status=%s tokens=%#v", child.Status, child.Tokens)
	}
	caller = mustInstance(t, eng, callerID)
	if caller.Status != projection.StatusCompleted {
		t.Fatalf("caller status=%s tokens=%#v", caller.Status, caller.Tokens)
	}
	childEvents, _ := eng.ListEvents(ctx, childID)
	if !sawElementIntent(childEvents, eventv1.Element_TYPE_SUB_PROCESS, "Event_SubProcess_called", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected event sub-process COMPLETED on child")
	}
	events, _ := eng.ListEvents(ctx, callerID)
	if !sawElementIntent(events, eventv1.Element_TYPE_CALL_ACTIVITY, "CallActivity_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected CallActivity COMPLETED after child Event Sub-Process finishes the called process")
	}
}

func TestCallActivityTwoCallsSameChild(t *testing.T) {
	xml := readTestdataCall(t, "m5_call_instance_twice.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	callerID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	caller := mustInstance(t, eng, callerID)
	_, hostTok := waitingAt(caller)
	child1 := caller.Tokens[hostTok].CalledProcessInstanceID
	c1 := mustInstance(t, eng, child1)
	e1, t1 := waitingAt(c1)
	if err := eng.Complete(ctx, child1, e1, t1, nil); err != nil {
		t.Fatal(err)
	}
	caller = mustInstance(t, eng, callerID)
	_, hostTok = waitingAt(caller)
	if caller.Tokens[hostTok].ElementID != "CallActivity_2" {
		t.Fatalf("expected second call, got %#v", caller.Tokens)
	}
	child2 := caller.Tokens[hostTok].CalledProcessInstanceID
	if child2 == "" || child2 == child1 {
		t.Fatalf("expected distinct child2, got %q (child1=%q)", child2, child1)
	}
	c2 := mustInstance(t, eng, child2)
	e2, t2 := waitingAt(c2)
	if err := eng.Complete(ctx, child2, e2, t2, nil); err != nil {
		t.Fatal(err)
	}
	caller = mustInstance(t, eng, callerID)
	if caller.Status != projection.StatusCompleted {
		t.Fatalf("caller status=%s", caller.Status)
	}
}

func mustListEvents(t *testing.T, eng *processing.Engine, instanceID string) []*eventv1.Event {
	t.Helper()
	events, err := eng.ListEvents(context.Background(), instanceID)
	if err != nil {
		t.Fatal(err)
	}
	return events
}

func readTestdataCall(t *testing.T, name string) []byte {
	t.Helper()
	xml, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return xml
}

func TestCallActivityIOMapping(t *testing.T) {
	xml := readTestdataCall(t, "m5_call_io.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	callerID, err := eng.CreateInstance(ctx, dep, map[string]any{"orderId": "o1", "extra": 1})
	if err != nil {
		t.Fatal(err)
	}
	caller := mustInstance(t, eng, callerID)
	_, hostTok := waitingAt(caller)
	childID := caller.Tokens[hostTok].CalledProcessInstanceID
	child := mustInstance(t, eng, childID)
	if _, ok := child.Variables["id"]; !ok {
		t.Fatalf("expected mapped id on child, vars=%v", child.Variables)
	}
	if _, ok := child.Variables["extra"]; ok {
		t.Fatalf("extra must not leak to child, vars=%v", child.Variables)
	}
	e1, t1 := waitingAt(child)
	if err := eng.Complete(ctx, childID, e1, t1, map[string]any{"total": 9}); err != nil {
		t.Fatal(err)
	}
	caller = mustInstance(t, eng, callerID)
	if caller.Status != projection.StatusCompleted {
		t.Fatalf("caller status=%s", caller.Status)
	}
	if got := caller.Variables["amount"]; got == "" {
		t.Fatalf("expected amount on caller, vars=%v", caller.Variables)
	}
	if _, ok := caller.Variables["extra"]; !ok {
		t.Fatalf("extra must remain on caller, vars=%v", caller.Variables)
	}
}

func TestCrossDeployCallActivitySpawnsCalleeDeployment(t *testing.T) {
	ctx := context.Background()
	eng := processing.NewEngine(eventlog.NewMemory())

	calleeDep, err := eng.Deploy(ctx, readTestdataCall(t, "m8_cross_call_callee.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	callerDep, err := eng.Deploy(ctx, readTestdataCall(t, "m8_cross_call_caller.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	if calleeDep == callerDep {
		t.Fatal("expected distinct deployments")
	}

	callerID, err := eng.CreateInstance(ctx, callerDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	caller := mustInstance(t, eng, callerID)
	if caller.DeploymentID != callerDep {
		t.Fatalf("caller deployment=%q want %q", caller.DeploymentID, callerDep)
	}
	hostElem, hostTok := waitingAt(caller)
	if hostElem != "CallActivity_1" {
		t.Fatalf("expected CallActivity_1 host, got %s", hostElem)
	}
	childID := caller.Tokens[hostTok].CalledProcessInstanceID
	child := mustInstance(t, eng, childID)
	if child.DeploymentID != calleeDep {
		t.Fatalf("child deployment=%q want callee %q", child.DeploymentID, calleeDep)
	}
	if child.ProcessID != "Process_called_cross" {
		t.Fatalf("child process=%q", child.ProcessID)
	}
	if child.ParentProcessInstanceID != callerID {
		t.Fatalf("child parent=%q want %q", child.ParentProcessInstanceID, callerID)
	}
	elemID, tokenID := waitingAt(child)
	if err := eng.Complete(ctx, childID, elemID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	caller = mustInstance(t, eng, callerID)
	if caller.Status != projection.StatusCompleted {
		t.Fatalf("caller status=%s", caller.Status)
	}
}

func TestCrossDeployCallerOnlyDeploySucceeds(t *testing.T) {
	ctx := context.Background()
	eng := processing.NewEngine(eventlog.NewMemory())
	if _, err := eng.Deploy(ctx, readTestdataCall(t, "m8_cross_call_caller.bpmn")); err != nil {
		t.Fatal(err)
	}
}

func TestCrossDeployCalleeNotDeployedRejects(t *testing.T) {
	ctx := context.Background()
	memLog := eventlog.NewMemory()
	eng := processing.NewEngine(memLog)
	callerDep, err := eng.Deploy(ctx, readTestdataCall(t, "m8_cross_call_caller.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = eng.CreateInstance(ctx, callerDep, nil)
	if err == nil {
		t.Fatal("expected CreateInstance error when callee not deployed")
	}
	records, err := memLog.ReadAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var instanceID string
	var sawRejection bool
	for _, ev := range records {
		if ev.GetProcessInstanceId() != "" && instanceID == "" {
			instanceID = ev.GetProcessInstanceId()
		}
		if ev.GetRecordType() == eventv1.Event_RECORD_TYPE_REJECTION && ev.GetRejection().GetCode() == "NOT_FOUND" {
			sawRejection = true
		}
	}
	if !sawRejection {
		t.Fatal("expected NOT_FOUND REJECTION in event log")
	}
	caller, ok := eng.GetInstance(instanceID)
	if !ok {
		t.Fatal("caller instance should exist after failed enter")
	}
	for _, tok := range caller.Tokens {
		if tok != nil && tok.Status == projection.TokenWaiting && tok.ElementID == "CallActivity_1" {
			t.Fatalf("caller must not wait on call when callee missing, tok=%#v", tok)
		}
	}
}

func TestCrossDeployCallActivityIOMapping(t *testing.T) {
	ctx := context.Background()
	eng := processing.NewEngine(eventlog.NewMemory())
	if _, err := eng.Deploy(ctx, readTestdataCall(t, "m8_cross_call_io_callee.bpmn")); err != nil {
		t.Fatal(err)
	}
	callerDep, err := eng.Deploy(ctx, readTestdataCall(t, "m8_cross_call_io_caller.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	callerID, err := eng.CreateInstance(ctx, callerDep, map[string]any{"orderId": "o1", "extra": 1})
	if err != nil {
		t.Fatal(err)
	}
	caller := mustInstance(t, eng, callerID)
	_, hostTok := waitingAt(caller)
	childID := caller.Tokens[hostTok].CalledProcessInstanceID
	child := mustInstance(t, eng, childID)
	if _, ok := child.Variables["id"]; !ok {
		t.Fatalf("expected mapped id on child, vars=%v", child.Variables)
	}
	if _, ok := child.Variables["extra"]; ok {
		t.Fatalf("extra must not leak to child, vars=%v", child.Variables)
	}
	e1, t1 := waitingAt(child)
	if err := eng.Complete(ctx, childID, e1, t1, map[string]any{"total": 9}); err != nil {
		t.Fatal(err)
	}
	caller = mustInstance(t, eng, callerID)
	if caller.Status != projection.StatusCompleted {
		t.Fatalf("caller status=%s", caller.Status)
	}
	if got := caller.Variables["amount"]; got == "" {
		t.Fatalf("expected amount on caller, vars=%v", caller.Variables)
	}
	if _, ok := caller.Variables["extra"]; !ok {
		t.Fatalf("extra must remain on caller, vars=%v", caller.Variables)
	}
}

func TestCallActivityTimerBoundaryPT0SFireDue(t *testing.T) {
	xml := readTestdataCall(t, "m9_call_timer_boundary.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	callerID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	caller := mustInstance(t, eng, callerID)
	hostElem, hostTok := waitingAt(caller)
	if hostElem != "CallActivity_1" {
		t.Fatalf("expected CallActivity_1 host, got %s tokens=%#v", hostElem, caller.Tokens)
	}
	tok := caller.Tokens[hostTok]
	if tok.DueUnixMs == 0 || tok.BoundaryID != "TimerBoundary_1" {
		t.Fatalf("armed timer token=%#v", tok)
	}
	childID := tok.CalledProcessInstanceID
	if childID == "" {
		t.Fatal("expected child instance")
	}
	child := mustInstance(t, eng, childID)
	if waitingElement(child) != "Task_called" {
		t.Fatalf("expected child active, tokens=%#v", child.Tokens)
	}

	if err := eng.FireDue(ctx); err != nil {
		t.Fatal(err)
	}
	child = mustInstance(t, eng, childID)
	if child.Status != projection.StatusTerminated {
		t.Fatalf("child status=%s want terminated", child.Status)
	}
	caller = mustInstance(t, eng, callerID)
	if caller.Status != projection.StatusCompleted {
		t.Fatalf("caller status=%s tokens=%#v", caller.Status, caller.Tokens)
	}
	events, _ := eng.ListEvents(ctx, callerID)
	if !sawElementIntent(events, eventv1.Element_TYPE_CALL_ACTIVITY, "CallActivity_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected CallActivity TERMINATED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "TimerBoundary_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected boundary COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_caller_timeout", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected timeout end")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_caller_ok", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("happy-path end must not complete")
	}
}

func TestCallActivityErrorBoundaryThrowInterrupts(t *testing.T) {
	xml := readTestdataCall(t, "m9_call_error_boundary.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	callerID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	caller := mustInstance(t, eng, callerID)
	hostElem, hostTok := waitingAt(caller)
	if hostElem != "CallActivity_1" {
		t.Fatalf("tokens=%#v", caller.Tokens)
	}
	childID := caller.Tokens[hostTok].CalledProcessInstanceID

	if err := eng.ThrowError(ctx, callerID, hostElem, hostTok, "BUSINESS_ERROR"); err != nil {
		t.Fatal(err)
	}
	child := mustInstance(t, eng, childID)
	if child.Status != projection.StatusTerminated {
		t.Fatalf("child status=%s want terminated", child.Status)
	}
	caller = mustInstance(t, eng, callerID)
	if caller.Status != projection.StatusCompleted {
		t.Fatalf("caller status=%s", caller.Status)
	}
	events, _ := eng.ListEvents(ctx, callerID)
	if !sawElementIntent(events, eventv1.Element_TYPE_CALL_ACTIVITY, "CallActivity_1", eventv1.Element_INTENT_ERROR_THROWN) {
		t.Fatal("expected ERROR_THROWN on call activity")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_CALL_ACTIVITY, "CallActivity_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected CallActivity TERMINATED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_BOUNDARY_EVENT, "ErrorBoundary_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected error boundary COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_caller_err", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected error path end")
	}
}

func TestCallActivityNonInterruptingMessageBoundary(t *testing.T) {
	xml := readTestdataCall(t, "m9_call_message_non_interrupt.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	callerID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	caller := mustInstance(t, eng, callerID)
	hostElem, hostTok := waitingAt(caller)
	if hostElem != "CallActivity_1" {
		t.Fatalf("tokens=%#v", caller.Tokens)
	}
	childID := caller.Tokens[hostTok].CalledProcessInstanceID

	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "order.cancelled"})
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("delivered=%d want 1", n)
	}
	caller = mustInstance(t, eng, callerID)
	if caller.Status != projection.StatusActive {
		t.Fatalf("status=%s want active", caller.Status)
	}
	if len(caller.Tokens) < 2 {
		t.Fatalf("want host + boundary tokens, got %#v", caller.Tokens)
	}
	hostElem, hostTok = waitingAt(caller)
	if hostElem != "CallActivity_1" {
		t.Fatalf("call host must still wait, tokens=%#v", caller.Tokens)
	}
	child := mustInstance(t, eng, childID)
	if child.Status != projection.StatusActive {
		t.Fatalf("child must remain active, status=%s", child.Status)
	}
	events, _ := eng.ListEvents(ctx, callerID)
	if sawElementIntent(events, eventv1.Element_TYPE_CALL_ACTIVITY, "CallActivity_1", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("CallActivity must not be terminated by non-interrupting boundary")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_caller_msg", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected message end from spawned token")
	}
	elemID, tokenID := waitingAt(child)
	if err := eng.Complete(ctx, childID, elemID, tokenID, nil); err != nil {
		t.Fatal(err)
	}
	caller = mustInstance(t, eng, callerID)
	if caller.Status != projection.StatusCompleted {
		t.Fatalf("caller status=%s tokens=%#v", caller.Status, caller.Tokens)
	}
}

func TestDeployAcceptsMultiInstanceCallActivity(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	if _, err := eng.Deploy(context.Background(), readTestdataCall(t, "m10_mi_call_parallel.bpmn")); err != nil {
		t.Fatalf("Deploy: %v", err)
	}
}

func miCallInnerChildren(t *testing.T, eng *processing.Engine, callerID string) []*projection.Instance {
	t.Helper()
	caller := mustInstance(t, eng, callerID)
	var children []*projection.Instance
	seen := map[string]bool{}
	for _, tok := range caller.Tokens {
		if tok == nil || tok.LoopInstanceIndex < 0 || tok.MultiInstanceHost {
			continue
		}
		if tok.CalledProcessInstanceID == "" || seen[tok.CalledProcessInstanceID] {
			continue
		}
		seen[tok.CalledProcessInstanceID] = true
		children = append(children, mustInstance(t, eng, tok.CalledProcessInstanceID))
	}
	return children
}

func TestMultiInstanceCallActivityParallel(t *testing.T) {
	xml := readTestdataCall(t, "m10_mi_call_parallel.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	callerID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	children := miCallInnerChildren(t, eng, callerID)
	if len(children) != 3 {
		t.Fatalf("children=%d want 3", len(children))
	}
	ids := map[string]bool{}
	for _, c := range children {
		if ids[c.ID] {
			t.Fatalf("duplicate child %s", c.ID)
		}
		ids[c.ID] = true
		if c.ParentProcessInstanceID != callerID {
			t.Fatalf("child parent=%q", c.ParentProcessInstanceID)
		}
		elemID, tokenID := waitingAt(c)
		if err := eng.Complete(ctx, c.ID, elemID, tokenID, nil); err != nil {
			t.Fatal(err)
		}
	}
	caller := mustInstance(t, eng, callerID)
	if caller.Status != projection.StatusCompleted {
		t.Fatalf("caller status=%s tokens=%#v", caller.Status, caller.Tokens)
	}
}

func TestMultiInstanceCallActivitySequential(t *testing.T) {
	xml := readTestdataCall(t, "m10_mi_call_sequential.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	callerID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		children := miCallInnerChildren(t, eng, callerID)
		if len(children) != 1 {
			t.Fatalf("step %d: active children=%d want 1 tokens=%#v", i, len(children), mustInstance(t, eng, callerID).Tokens)
		}
		elemID, tokenID := waitingAt(children[0])
		if err := eng.Complete(ctx, children[0].ID, elemID, tokenID, nil); err != nil {
			t.Fatal(err)
		}
	}
	caller := mustInstance(t, eng, callerID)
	if caller.Status != projection.StatusCompleted {
		t.Fatalf("caller status=%s", caller.Status)
	}
}

func TestMultiInstanceCallActivityCollectionIO(t *testing.T) {
	xml := readTestdataCall(t, "m10_mi_call_collection.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	callerID, err := eng.CreateInstance(ctx, dep, map[string]any{"items": []any{"a", "b"}})
	if err != nil {
		t.Fatal(err)
	}
	children := miCallInnerChildren(t, eng, callerID)
	if len(children) != 2 {
		t.Fatalf("children=%d want 2", len(children))
	}
	lines := map[string]bool{}
	for _, c := range children {
		if _, ok := c.Variables["line"]; !ok {
			t.Fatalf("expected mapped line on child vars=%v", c.Variables)
		}
		lines[c.Variables["line"]] = true
		elemID, tokenID := waitingAt(c)
		if err := eng.Complete(ctx, c.ID, elemID, tokenID, map[string]any{"out": "x"}); err != nil {
			t.Fatal(err)
		}
	}
	if len(lines) != 2 {
		t.Fatalf("expected distinct line values, got %v", lines)
	}
	caller := mustInstance(t, eng, callerID)
	if caller.Status != projection.StatusCompleted {
		t.Fatalf("caller status=%s", caller.Status)
	}
	if got := caller.Variables["results"]; got == "" {
		t.Fatalf("expected results output collection, vars=%v", caller.Variables)
	}
}

func TestMultiInstanceCallActivityEmptyCollection(t *testing.T) {
	xml := readTestdataCall(t, "m10_mi_call_collection.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	callerID, err := eng.CreateInstance(ctx, dep, map[string]any{"items": []any{}})
	if err != nil {
		t.Fatal(err)
	}
	caller := mustInstance(t, eng, callerID)
	if caller.Status != projection.StatusCompleted {
		t.Fatalf("caller status=%s tokens=%#v", caller.Status, caller.Tokens)
	}
	if len(miCallInnerChildren(t, eng, callerID)) != 0 {
		t.Fatal("empty collection must not start children")
	}
}

func TestMultiInstanceCallActivityCompletionEarly(t *testing.T) {
	xml := readTestdataCall(t, "m10_mi_call_completion.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	callerID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	children := miCallInnerChildren(t, eng, callerID)
	if len(children) != 5 {
		t.Fatalf("children=%d want 5", len(children))
	}
	for i := 0; i < 2; i++ {
		elemID, tokenID := waitingAt(children[i])
		if err := eng.Complete(ctx, children[i].ID, elemID, tokenID, nil); err != nil {
			t.Fatal(err)
		}
	}
	caller := mustInstance(t, eng, callerID)
	if caller.Status != projection.StatusCompleted {
		t.Fatalf("caller status=%s tokens=%#v", caller.Status, caller.Tokens)
	}
	for i := 2; i < 5; i++ {
		c := mustInstance(t, eng, children[i].ID)
		if c.Status != projection.StatusTerminated {
			t.Fatalf("straggler child %d status=%s want terminated", i, c.Status)
		}
	}
}

func TestMultiInstanceCallActivityTimerBoundary(t *testing.T) {
	xml := readTestdataCall(t, "m10_mi_call_timer_boundary.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	callerID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	children := miCallInnerChildren(t, eng, callerID)
	if len(children) != 2 {
		t.Fatalf("children=%d want 2", len(children))
	}
	if err := eng.FireDue(ctx); err != nil {
		t.Fatal(err)
	}
	for _, c := range children {
		got := mustInstance(t, eng, c.ID)
		if got.Status != projection.StatusTerminated {
			t.Fatalf("child status=%s want terminated", got.Status)
		}
	}
	caller := mustInstance(t, eng, callerID)
	if caller.Status != projection.StatusCompleted {
		t.Fatalf("caller status=%s", caller.Status)
	}
	events, _ := eng.ListEvents(ctx, callerID)
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_caller_timeout", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected timeout end")
	}
	if sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_caller_ok", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("happy path must not complete")
	}
}
