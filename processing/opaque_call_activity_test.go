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

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
	"github.com/sparrow-community/sparrow/processing/projection"
)

func TestOpaqueCallActivityMissingCalledElement(t *testing.T) {
	xml := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<definitions xmlns="http://www.omg.org/spec/BPMN/20100524/MODEL" id="Defs" targetNamespace="http://sparrow.local">
  <process id="P" isExecutable="true">
    <startEvent id="Start"><outgoing>f1</outgoing></startEvent>
    <callActivity id="StubCA" name="Missing callee"><incoming>f1</incoming><outgoing>f2</outgoing></callActivity>
    <endEvent id="End"><incoming>f2</incoming></endEvent>
    <sequenceFlow id="f1" sourceRef="Start" targetRef="StubCA"/>
    <sequenceFlow id="f2" sourceRef="StubCA" targetRef="End"/>
  </process>
</definitions>`)
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
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
	if el != "StubCA" {
		t.Fatalf("wait=%s", el)
	}
	if err := eng.Complete(ctx, id, el, tok, nil); err != nil {
		t.Fatal(err)
	}
	inst = mustInstance(t, eng, id)
	if inst.Status != projection.StatusCompleted {
		t.Fatalf("status=%s", inst.Status)
	}
}
