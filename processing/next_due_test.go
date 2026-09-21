// Copyright 2026 The Sparrow community and contributors
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
)

func TestNextDueUnixMsAfterTimerCatch(t *testing.T) {
	xml := readTestdata(t, "m2_timer_catch_1h.bpmn")
	eng := processing.NewEngine(eventlog.NewMemory())
	ctx := context.Background()
	dep, err := eng.Deploy(ctx, xml)
	if err != nil {
		t.Fatal(err)
	}
	if eng.NextDueUnixMs() != 0 {
		t.Fatalf("expected no due before instance, got %d", eng.NextDueUnixMs())
	}
	instanceID, err := eng.CreateInstance(ctx, dep, nil)
	if err != nil {
		t.Fatal(err)
	}
	due := eng.NextDueUnixMs()
	if due == 0 {
		t.Fatal("expected next due after timer catch wait")
	}
	inst, ok := eng.GetInstance(instanceID)
	if !ok {
		t.Fatal("missing instance")
	}
	_, tokenID := waitingAt(inst)
	tok := inst.Tokens[tokenID]
	if tok.DueUnixMs != due {
		t.Fatalf("NextDueUnixMs=%d token due=%d", due, tok.DueUnixMs)
	}
}
