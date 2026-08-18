package deploy

import (
	"fmt"
	"strings"
	"time"

	"github.com/sparrow-community/sparrow/bpmn/element"
)

type timerCatch struct {
	Duration time.Duration
	Text     string
}

func timerCatchSpec(ev element.IntermediateCatchEvent) (timerCatch, error) {
	if n := otherCatchDefinitions(ev); n > 0 {
		return timerCatch{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateCatchEvent %q is not a timer catch", ev.ID)
	}
	if len(ev.TimerEventDefinitions) != 1 {
		return timerCatch{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateCatchEvent %q needs exactly one timerEventDefinition", ev.ID)
	}
	def := ev.TimerEventDefinitions[0]
	if text := expressionText(def.TimeDate); text != "" {
		return timerCatch{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateCatchEvent %q uses timeDate", ev.ID)
	}
	if text := expressionText(def.TimeCycle); text != "" {
		return timerCatch{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateCatchEvent %q uses timeCycle", ev.ID)
	}
	text := expressionText(def.TimeDuration)
	if text == "" {
		return timerCatch{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateCatchEvent %q needs timeDuration", ev.ID)
	}
	dur, err := ParseISO8601Duration(text)
	if err != nil {
		return timerCatch{}, err
	}
	return timerCatch{Duration: dur, Text: text}, nil
}

func otherCatchDefinitions(ev element.IntermediateCatchEvent) int {
	n := 0
	n += len(ev.MessageEventDefinitions)
	n += len(ev.EscalationEventDefinitions)
	n += len(ev.TerminateEventDefinitions)
	n += len(ev.SignalEventDefinitions)
	n += len(ev.ConditionalEventDefinitions)
	n += len(ev.ErrorEventDefinitions)
	n += len(ev.LinkEventDefinitions)
	n += len(ev.CompensateEventDefinitions)
	return n
}

func expressionText(e element.ExpressionUnMarshal) string {
	switch x := e.ExpressionSubstitution.(type) {
	case *element.FormalExpression:
		return strings.TrimSpace(x.Value)
	default:
		return ""
	}
}

// TimerDuration returns the parsed ISO-8601 wait and original text for a timer catch.
func (d *Deployment) TimerDuration(id string) (time.Duration, string, error) {
	spec, ok := d.timerCatch[id]
	if !ok {
		return 0, "", fmt.Errorf("NOT_FOUND: timer catch %q", id)
	}
	return spec.Duration, spec.Text, nil
}
