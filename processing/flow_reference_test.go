package processing_test

import (
	"context"
	"strings"
	"testing"

	"github.com/sparrow-community/sparrow/processing"
	eventlog "github.com/sparrow-community/sparrow/processing/log"
)

func TestDeployRejectsFlowToUnmodeledElement(t *testing.T) {
	ctx := context.Background()
	for _, tc := range []struct {
		fixture string
		want    []string
	}{
		{
			fixture: "m37_flow_to_unmodeled_element.bpmn",
			want:    []string{"sequenceFlow", "Choreo_1", "not a supported flow element"},
		},
		{
			fixture: "m37_flow_to_unmodeled_in_subprocess.bpmn",
			want:    []string{"sequenceFlow", "Choreo_1", "not a supported flow element"},
		},
		{
			fixture: "m37_dangling_outgoing.bpmn",
			want:    []string{"Task_A", "outgoing", "Flow_missing"},
		},
	} {
		t.Run(tc.fixture, func(t *testing.T) {
			xml := readTestdata(t, tc.fixture)
			eng := processing.NewEngine(eventlog.NewMemory())
			_, err := eng.Deploy(ctx, xml)
			if err == nil {
				t.Fatal("expected deploy rejection")
			}
			if !strings.Contains(err.Error(), "UNSUPPORTED_ELEMENT") {
				t.Fatalf("missing stable code: %v", err)
			}
			for _, want := range tc.want {
				if !strings.Contains(err.Error(), want) {
					t.Fatalf("error %q missing %q", err.Error(), want)
				}
			}
		})
	}
}

func TestDeployRejectionIsDeterministic(t *testing.T) {
	ctx := context.Background()
	xml := readTestdata(t, "m37_flow_to_unmodeled_element.bpmn")
	var first string
	for i := 0; i < 5; i++ {
		eng := processing.NewEngine(eventlog.NewMemory())
		_, err := eng.Deploy(ctx, xml)
		if err == nil {
			t.Fatal("expected deploy rejection")
		}
		if i == 0 {
			first = err.Error()
			continue
		}
		if err.Error() != first {
			t.Fatalf("rejection not stable: %q vs %q", first, err.Error())
		}
	}
}
