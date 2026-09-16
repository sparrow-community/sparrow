package processing_test

import (
	"context"
	"strings"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func completeAdHocInner(t *testing.T, eng *processing.Engine, instanceID, elementID string, vars map[string]any) {
	t.Helper()
	inst, ok := eng.GetInstance(instanceID)
	if !ok {
		t.Fatalf("missing instance %s", instanceID)
	}
	for _, tok := range inst.Tokens {
		if tok.ElementID == elementID && tok.Status == projection.TokenWaiting {
			if err := eng.Complete(context.Background(), instanceID, elementID, tok.ID, vars); err != nil {
				t.Fatalf("Complete %s: %v", elementID, err)
			}
			return
		}
	}
	t.Fatalf("no waiting token at %s, waiting=%v", elementID, allWaiting(inst))
}

func TestAdHocSubProcessParallelEnablesAllInnerActivities(t *testing.T) {
	xml := readTestdata(t, "m35_ad_hoc_parallel.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, map[string]any{"done": 0})
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	waiting := allWaiting(inst)
	for _, id := range []string{"Task_A", "Task_B", "Task_C"} {
		if !containsStr(waiting, id) {
			t.Fatalf("expected %s enabled, waiting=%v", id, waiting)
		}
	}

	completeAdHocInner(t, eng, instanceID, "Task_B", map[string]any{"done": 1})
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status == projection.StatusCompleted {
		t.Fatal("scope must not complete while condition is false")
	}
	if containsStr(allWaiting(inst), "Task_B") {
		t.Fatalf("Task_B should be done, waiting=%v", allWaiting(inst))
	}

	completeAdHocInner(t, eng, instanceID, "Task_A", map[string]any{"done": 2})
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s waiting=%v", inst.Status, allWaiting(inst))
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "Task_C", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected Task_C terminated by cancelRemainingInstances")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_AD_HOC_SUB_PROCESS, "AdHoc_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected AdHoc_1 COMPLETED")
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_END_EVENT, "End_1", eventv1.Element_INTENT_COMPLETED) {
		t.Fatal("expected End_1 after the scope completed")
	}
}

func TestAdHocSubProcessSequentialOrdering(t *testing.T) {
	xml := readTestdata(t, "m35_ad_hoc_sequential.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, map[string]any{"done": 0})
	if err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	waiting := allWaiting(inst)
	if !containsStr(waiting, "Task_A") {
		t.Fatalf("expected Task_A enabled, waiting=%v", waiting)
	}
	if containsStr(waiting, "Task_B") || containsStr(waiting, "Task_C") {
		t.Fatalf("Sequential must enable one activity at a time, waiting=%v", waiting)
	}

	completeAdHocInner(t, eng, instanceID, "Task_A", map[string]any{"done": 1})
	inst, _ = eng.GetInstance(instanceID)
	waiting = allWaiting(inst)
	if !containsStr(waiting, "Task_B") {
		t.Fatalf("expected Task_B enabled next, waiting=%v", waiting)
	}
	if containsStr(waiting, "Task_C") {
		t.Fatalf("Task_C must wait its turn, waiting=%v", waiting)
	}

	completeAdHocInner(t, eng, instanceID, "Task_B", map[string]any{"done": 2})
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s waiting=%v", inst.Status, allWaiting(inst))
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "Task_C", eventv1.Element_INTENT_ACTIVATED) {
		t.Fatal("Task_C must not be enabled after the condition became true")
	}
}

func TestAdHocSubProcessKeepsRemainingInstances(t *testing.T) {
	xml := readTestdata(t, "m35_ad_hoc_keep_remaining.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, map[string]any{"done": 0})
	if err != nil {
		t.Fatal(err)
	}
	completeAdHocInner(t, eng, instanceID, "Task_A", map[string]any{"done": 1})
	inst, _ := eng.GetInstance(instanceID)
	if inst.Status == projection.StatusCompleted {
		t.Fatal("scope must wait for Task_B when cancelRemainingInstances is false")
	}
	if !containsStr(allWaiting(inst), "Task_B") {
		t.Fatalf("expected Task_B still enabled, waiting=%v", allWaiting(inst))
	}

	completeAdHocInner(t, eng, instanceID, "Task_B", map[string]any{"done": 2})
	inst, _ = eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s waiting=%v", inst.Status, allWaiting(inst))
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "Task_B", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("Task_B must not be terminated with cancelRemainingInstances=false")
	}
}

func TestAdHocSubProcessCompletesWhenActivitiesExhausted(t *testing.T) {
	xml := readTestdata(t, "m35_ad_hoc_exhaust.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, map[string]any{"done": 0})
	if err != nil {
		t.Fatal(err)
	}
	completeAdHocInner(t, eng, instanceID, "Task_A", map[string]any{"done": 1})
	completeAdHocInner(t, eng, instanceID, "Task_B", map[string]any{"done": 2})
	inst, _ := eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s waiting=%v", inst.Status, allWaiting(inst))
	}
}

func TestAdHocSubProcessJobBackedInnerActivity(t *testing.T) {
	xml := readTestdata(t, "m35_ad_hoc_job.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	instanceID, err := eng.CreateInstance(ctx, dep, map[string]any{"done": 0})
	if err != nil {
		t.Fatal(err)
	}
	jobs, err := eng.Activate(ctx, processing.ActivateRequest{JobType: "work.v1", WorkerID: "w1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs) != 1 || jobs[0].ElementID != "Service_B" {
		t.Fatalf("jobs=%#v", jobs)
	}
	if err := eng.Complete(ctx, instanceID, jobs[0].ElementID, jobs[0].TokenID, map[string]any{"done": 1}); err != nil {
		t.Fatal(err)
	}
	inst, _ := eng.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s waiting=%v", inst.Status, allWaiting(inst))
	}
	events, err := eng.ListEvents(ctx, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if !sawElementIntent(events, eventv1.Element_TYPE_USER_TASK, "Task_A", eventv1.Element_INTENT_TERMINATED) {
		t.Fatal("expected Task_A terminated when the job completed the scope")
	}
}

func TestAdHocSubProcessRecover(t *testing.T) {
	xml := readTestdata(t, "m35_ad_hoc_parallel.bpmn")
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
	instanceID, err := eng1.CreateInstance(ctx, dep, map[string]any{"done": 0})
	if err != nil {
		t.Fatal(err)
	}
	completeAdHocInner(t, eng1, instanceID, "Task_A", map[string]any{"done": 1})
	if err := eng1.Close(); err != nil {
		t.Fatal(err)
	}

	eng2, err := processing.Open(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	defer eng2.Close()

	inst, ok := eng2.GetInstance(instanceID)
	if !ok {
		t.Fatal("missing instance after recover")
	}
	waiting := allWaiting(inst)
	if !containsStr(waiting, "Task_B") || !containsStr(waiting, "Task_C") {
		t.Fatalf("expected Task_B and Task_C after recover, waiting=%v", waiting)
	}
	var host *projection.Token
	for _, tok := range inst.Tokens {
		if tok.ElementID == "AdHoc_1" {
			host = tok
		}
	}
	if host == nil {
		t.Fatalf("expected parked host token on AdHoc_1, tokens=%v", waiting)
	}

	completeAdHocInner(t, eng2, instanceID, "Task_B", map[string]any{"done": 2})
	inst, _ = eng2.GetInstance(instanceID)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s waiting=%v", inst.Status, allWaiting(inst))
	}
}

func TestAdHocSubProcessDeployRejections(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		name string
		body string
		want string
	}{
		{
			name: "missing completion condition",
			body: `<adHocSubProcess id="AdHoc_1"><userTask id="Task_A"/></adHocSubProcess>`,
			want: "completionCondition",
		},
		{
			name: "bad ordering",
			body: `<adHocSubProcess id="AdHoc_1" ordering="Random"><userTask id="Task_A"/>` +
				`<completionCondition xsi:type="tFormalExpression">${done &gt;= 1}</completionCondition></adHocSubProcess>`,
			want: "ordering",
		},
		{
			name: "sequence flow inside",
			body: `<adHocSubProcess id="AdHoc_1"><userTask id="Task_A"/><userTask id="Task_B"/>` +
				`<sequenceFlow id="Flow_inner" sourceRef="Task_A" targetRef="Task_B"/>` +
				`<completionCondition xsi:type="tFormalExpression">${done &gt;= 1}</completionCondition></adHocSubProcess>`,
			want: "sequence flows",
		},
		{
			name: "no inner activity",
			body: `<adHocSubProcess id="AdHoc_1">` +
				`<completionCondition xsi:type="tFormalExpression">${done &gt;= 1}</completionCondition></adHocSubProcess>`,
			want: "at least one inner activity",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			eng := processing.NewEngine(eventlog.NewMemory())
			if _, err := eng.Deploy(ctx, []byte(adHocRejectionXML(tc.body))); err == nil {
				t.Fatal("expected deploy rejection")
			} else if !strings.Contains(err.Error(), "UNSUPPORTED_ELEMENT") || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func adHocRejectionXML(body string) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<definitions xmlns="http://www.omg.org/spec/BPMN/20100524/MODEL"
             xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
             id="Definitions_ad_hoc_reject"
             targetNamespace="http://sparrow.example/adhoc">
  <process id="Process_ad_hoc_reject" isExecutable="true">
    <startEvent id="Start_1">
      <outgoing>Flow_start_to_adhoc</outgoing>
    </startEvent>
    ` + body + `
    <endEvent id="End_1">
      <incoming>Flow_adhoc_to_end</incoming>
    </endEvent>
    <sequenceFlow id="Flow_start_to_adhoc" sourceRef="Start_1" targetRef="AdHoc_1"/>
    <sequenceFlow id="Flow_adhoc_to_end" sourceRef="AdHoc_1" targetRef="End_1"/>
  </process>
</definitions>`
}
