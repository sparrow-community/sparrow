// Copyright 2025 The Sparrow community and contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package processing_test

import (
	"context"
	"testing"
	"time"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
)

func TestReceiveWithTimerBoundariesAcceptsMessage(t *testing.T) {
	xml := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<definitions xmlns="http://www.omg.org/spec/BPMN/20100524/MODEL" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" id="Defs" targetNamespace="http://sparrow.local">
  <message id="Msg1" name="doc.received"/>
  <process id="P" isExecutable="true">
    <startEvent id="Start"><outgoing>f1</outgoing></startEvent>
    <receiveTask id="Recv" messageRef="Msg1"><incoming>f1</incoming><outgoing>f2</outgoing></receiveTask>
    <boundaryEvent id="Daily" cancelActivity="false" attachedToRef="Recv">
      <outgoing>f3</outgoing>
      <timerEventDefinition><timeCycle xsi:type="tFormalExpression">R2/P1D</timeCycle></timerEventDefinition>
    </boundaryEvent>
    <boundaryEvent id="Week" attachedToRef="Recv">
      <outgoing>f4</outgoing>
      <timerEventDefinition><timeDuration xsi:type="tFormalExpression">P7D</timeDuration></timerEventDefinition>
    </boundaryEvent>
    <endEvent id="End"><incoming>f2</incoming></endEvent>
    <endEvent id="RemindEnd"><incoming>f3</incoming></endEvent>
    <endEvent id="EscalationEnd"><incoming>f4</incoming></endEvent>
    <sequenceFlow id="f1" sourceRef="Start" targetRef="Recv"/>
    <sequenceFlow id="f2" sourceRef="Recv" targetRef="End"/>
    <sequenceFlow id="f3" sourceRef="Daily" targetRef="RemindEnd"/>
    <sequenceFlow id="f4" sourceRef="Week" targetRef="EscalationEnd"/>
  </process>
</definitions>`)
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	now := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	eng.SetNow(func() time.Time { return now })
	depID, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	id, err := eng.CreateInstance(ctx, depID, nil)
	if err != nil {
		t.Fatal(err)
	}
	inst := mustInstance(t, eng, id)
	el, tok := waitingAt(inst)
	if el != "Recv" {
		t.Fatalf("wait=%s", el)
	}
	if inst.Tokens[tok].MessageName != "doc.received" {
		t.Fatalf("message=%q", inst.Tokens[tok].MessageName)
	}
	n, err := eng.PublishMessage(ctx, processing.PublishMessageRequest{Name: "doc.received", ProcessInstanceID: id})
	if err != nil || n != 1 {
		t.Fatalf("n=%d err=%v", n, err)
	}
	inst = mustInstance(t, eng, id)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s", inst.Status)
	}
}
