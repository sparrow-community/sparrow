package processing_test

import (
	"context"
	"strings"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
)

func TestImplicitThrowEventRejectedAtDeploy(t *testing.T) {
	ctx := context.Background()
	for _, fixture := range []string{
		"m36_implicit_throw_process.bpmn",
		"m36_implicit_throw_subprocess.bpmn",
	} {
		t.Run(fixture, func(t *testing.T) {
			xml := readTestdata(t, fixture)
			eng := processing.NewEngine(eventlog.NewMemory())
			_, err := eng.Deploy(ctx, xml)
			if err == nil {
				t.Fatal("expected deploy rejection")
			}
			if !strings.Contains(err.Error(), "UNSUPPORTED_ELEMENT") ||
				!strings.Contains(err.Error(), "implicitThrowEvent") ||
				!strings.Contains(err.Error(), "ImplicitThrow_1") {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestMultiInstanceBehaviorEventStillDeploys(t *testing.T) {
	ctx := context.Background()
	for _, fixture := range []string{
		"m25_mi_complex_behavior.bpmn",
		"m25_mi_none_behavior_event.bpmn",
		"m25_mi_one_behavior_event.bpmn",
	} {
		t.Run(fixture, func(t *testing.T) {
			xml := readTestdata(t, fixture)
			eng := processing.NewEngine(eventlog.NewMemory())
			if _, err := eng.Deploy(ctx, xml); err != nil {
				t.Fatalf("multi-instance behavior event must stay deployable: %v", err)
			}
		})
	}
}
