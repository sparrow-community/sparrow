package processing_test

import (
	"context"
	"strings"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
)

func TestMIOneBehaviorEventRef(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	catchDep, err := eng.Deploy(ctx, readTestdata(t, "m25_catch_mi_first.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	miDep, err := eng.Deploy(ctx, readTestdata(t, "m25_mi_one_behavior_event.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	catchID, err := eng.CreateInstance(ctx, catchDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	miID, err := eng.CreateInstance(ctx, miDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inners := waitingTokensAt(mustInstance(t, eng, miID), "UserTask_mi")
	if len(inners) != 3 {
		t.Fatalf("inners=%d", len(inners))
	}
	if err := eng.Complete(ctx, miID, "UserTask_mi", inners[0], nil); err != nil {
		t.Fatal(err)
	}
	mi := mustInstance(t, eng, miID)
	if mi.Status != projection.StatusCompleted {
		t.Fatalf("mi status=%s", mi.Status)
	}
	catch := mustInstance(t, eng, catchID)
	if catch.Status != projection.StatusCompleted {
		t.Fatalf("catch status=%s", catch.Status)
	}
}

func TestMINoneBehaviorEventRef(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	catchDep, err := eng.Deploy(ctx, readTestdata(t, "m25_catch_mi_each.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	miDep, err := eng.Deploy(ctx, readTestdata(t, "m25_mi_none_behavior_event.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	c1, err := eng.CreateInstance(ctx, catchDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	c2, err := eng.CreateInstance(ctx, catchDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	miID, err := eng.CreateInstance(ctx, miDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inners := waitingTokensAt(mustInstance(t, eng, miID), "UserTask_mi")
	if len(inners) != 2 {
		t.Fatalf("inners=%d", len(inners))
	}
	if err := eng.Complete(ctx, miID, "UserTask_mi", inners[0], nil); err != nil {
		t.Fatal(err)
	}
	if err := eng.Complete(ctx, miID, "UserTask_mi", inners[1], nil); err != nil {
		t.Fatal(err)
	}
	if mustInstance(t, eng, miID).Status != projection.StatusCompleted {
		t.Fatal("mi not completed")
	}
	if mustInstance(t, eng, c1).Status != projection.StatusCompleted {
		t.Fatal("catch1 not completed")
	}
	if mustInstance(t, eng, c2).Status != projection.StatusCompleted {
		t.Fatal("catch2 not completed")
	}
}

func TestMIComplexBehaviorDefinition(t *testing.T) {
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	catchDep, err := eng.Deploy(ctx, readTestdata(t, "m25_catch_mi_half.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	miDep, err := eng.Deploy(ctx, readTestdata(t, "m25_mi_complex_behavior.bpmn"))
	if err != nil {
		t.Fatal(err)
	}
	catchID, err := eng.CreateInstance(ctx, catchDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	miID, err := eng.CreateInstance(ctx, miDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	inners := waitingTokensAt(mustInstance(t, eng, miID), "UserTask_mi")
	if len(inners) != 3 {
		t.Fatalf("inners=%d", len(inners))
	}
	if err := eng.Complete(ctx, miID, "UserTask_mi", inners[0], nil); err != nil {
		t.Fatal(err)
	}
	if mustInstance(t, eng, catchID).Status != projection.StatusCompleted {
		t.Fatal("expected complex signal on first complete")
	}
	// Second catcher should not complete on later publishes (one-fire).
	c2, err := eng.CreateInstance(ctx, catchDep, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := eng.Complete(ctx, miID, "UserTask_mi", inners[1], nil); err != nil {
		t.Fatal(err)
	}
	if err := eng.Complete(ctx, miID, "UserTask_mi", inners[2], nil); err != nil {
		t.Fatal(err)
	}
	if mustInstance(t, eng, miID).Status != projection.StatusCompleted {
		t.Fatal("mi not completed")
	}
	if mustInstance(t, eng, c2).Status == projection.StatusCompleted {
		t.Fatal("complex event must not re-publish")
	}
}

func TestMIBehaviorRejectsBadComplexEvent(t *testing.T) {
	body := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<definitions xmlns="http://www.omg.org/spec/BPMN/20100524/MODEL" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" id="d" targetNamespace="http://sparrow.example/x">
  <process id="P" isExecutable="true">
    <startEvent id="S"><outgoing>F1</outgoing></startEvent>
    <userTask id="U"><incoming>F1</incoming><outgoing>F2</outgoing>
      <multiInstanceLoopCharacteristics behavior="Complex">
        <loopCardinality xsi:type="tFormalExpression">2</loopCardinality>
        <complexBehaviorDefinition>
          <condition xsi:type="tFormalExpression">true</condition>
          <event id="E"><timerEventDefinition><timeDuration xsi:type="tFormalExpression">PT1S</timeDuration></timerEventDefinition></event>
        </complexBehaviorDefinition>
      </multiInstanceLoopCharacteristics>
    </userTask>
    <endEvent id="End"><incoming>F2</incoming></endEvent>
    <sequenceFlow id="F1" sourceRef="S" targetRef="U"/>
    <sequenceFlow id="F2" sourceRef="U" targetRef="End"/>
  </process>
</definitions>`)
	eng := processing.NewEngine(eventlog.NewMemory())
	_, err := eng.Deploy(context.Background(), body)
	if err == nil || !strings.Contains(err.Error(), "UNSUPPORTED_ELEMENT") {
		t.Fatalf("err=%v", err)
	}
}

func waitingTokensAt(inst *projection.Instance, elementID string) []string {
	var out []string
	for id, tok := range inst.Tokens {
		if tok != nil && tok.ElementID == elementID && tok.Status == projection.TokenWaiting && !tok.MultiInstanceHost {
			out = append(out, id)
		}
	}
	return out
}
