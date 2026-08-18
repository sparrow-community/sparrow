package deploy

import (
	"fmt"
	"strings"
	"time"

	"github.com/sparrow-community/sparrow/bpmn/element"
)

// CatchKind classifies an intermediateCatchEvent for runtime handlers.
type CatchKind string

const (
	CatchKindTimer   CatchKind = "timer"
	CatchKindMessage CatchKind = "message"
)

type timerCatch struct {
	Duration time.Duration
	Date     time.Time
	Cycle    *cycleSpec
	Text     string
}

type messageCatch struct {
	Name string
}

func timerCatchSpec(ev element.IntermediateCatchEvent) (timerCatch, error) {
	if n := otherCatchDefinitions(ev); n > 0 {
		return timerCatch{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateCatchEvent %q is not a timer catch", ev.ID)
	}
	if len(ev.TimerEventDefinitions) != 1 {
		return timerCatch{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateCatchEvent %q needs exactly one timerEventDefinition", ev.ID)
	}
	def := ev.TimerEventDefinitions[0]
	dateText := expressionText(def.TimeDate)
	durText := expressionText(def.TimeDuration)
	cycleText := expressionText(def.TimeCycle)
	n := 0
	if dateText != "" {
		n++
	}
	if durText != "" {
		n++
	}
	if cycleText != "" {
		n++
	}
	if n > 1 {
		return timerCatch{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateCatchEvent %q has multiple timer expressions", ev.ID)
	}
	if cycleText != "" {
		cyc, err := ParseISO8601Cycle(cycleText)
		if err != nil {
			return timerCatch{}, err
		}
		return timerCatch{Cycle: &cyc, Text: cycleText}, nil
	}
	if dateText != "" {
		at, err := ParseISO8601Date(dateText)
		if err != nil {
			return timerCatch{}, err
		}
		return timerCatch{Date: at, Text: dateText}, nil
	}
	if durText == "" {
		return timerCatch{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateCatchEvent %q needs timeDuration, timeDate, or timeCycle", ev.ID)
	}
	dur, err := ParseISO8601Duration(durText)
	if err != nil {
		return timerCatch{}, err
	}
	return timerCatch{Duration: dur, Text: durText}, nil
}

func messageCatchSpec(ev element.IntermediateCatchEvent, messages []element.Message) (messageCatch, error) {
	if len(ev.TimerEventDefinitions) > 0 {
		return messageCatch{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateCatchEvent %q is not a message catch (has timerEventDefinition)", ev.ID)
	}
	if len(ev.MessageEventDefinitions) != 1 {
		return messageCatch{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateCatchEvent %q needs exactly one messageEventDefinition", ev.ID)
	}
	other := otherCatchDefinitions(ev) - len(ev.MessageEventDefinitions)
	if other > 0 {
		return messageCatch{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateCatchEvent %q is not a message catch", ev.ID)
	}
	name := resolveMessageName(ev.MessageEventDefinitions[0].MessageRef, messages)
	if name == "" {
		name = strings.TrimSpace(ev.Name)
	}
	if name == "" {
		name = ev.ID
	}
	return messageCatch{Name: name}, nil
}

func resolveMessageName(ref string, messages []element.Message) string {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return ""
	}
	for _, m := range messages {
		if m.ID != ref {
			continue
		}
		if n := strings.TrimSpace(m.Name); n != "" {
			return n
		}
		return m.ID
	}
	return ref
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

// TimerDuration returns the parsed ISO-8601 wait and original text for a duration timer catch.
func (d *Deployment) TimerDuration(id string) (time.Duration, string, error) {
	spec, ok := d.timerCatch[id]
	if !ok {
		return 0, "", fmt.Errorf("NOT_FOUND: timer catch %q", id)
	}
	if !spec.Date.IsZero() || spec.Cycle != nil {
		return 0, "", fmt.Errorf("NOT_FOUND: timer catch %q is not timeDuration", id)
	}
	return spec.Duration, spec.Text, nil
}

// TimerDue returns the absolute due time for a timer catch.
// timeDate uses the parsed instant; timeDuration is now + duration;
// timeCycle uses the first occurrence at or after now (intermediate catch does not re-arm).
func (d *Deployment) TimerDue(id string, now time.Time) (dueUnixMs int64, text string, err error) {
	spec, ok := d.timerCatch[id]
	if !ok {
		return 0, "", fmt.Errorf("NOT_FOUND: timer catch %q", id)
	}
	if now.IsZero() {
		now = time.Now()
	}
	if spec.Cycle != nil {
		due, err := spec.Cycle.FirstDue(now)
		if err != nil {
			return 0, "", err
		}
		return due.UnixMilli(), spec.Text, nil
	}
	if !spec.Date.IsZero() {
		return spec.Date.UnixMilli(), spec.Text, nil
	}
	return now.Add(spec.Duration).UnixMilli(), spec.Text, nil
}

// MessageName returns the BPMN message name a catch is waiting for.
func (d *Deployment) MessageName(id string) (string, error) {
	name, ok := d.messageCatch[id]
	if !ok || name == "" {
		return "", fmt.Errorf("NOT_FOUND: message catch %q", id)
	}
	return name, nil
}

func (d *Deployment) CatchKind(id string) (CatchKind, error) {
	kind, ok := d.catchKinds[id]
	if !ok {
		return "", fmt.Errorf("NOT_FOUND: intermediate catch %q", id)
	}
	return kind, nil
}
