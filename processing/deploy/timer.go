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
	return timerCatchFromDefs(ev.ID, ev.EventDefinitions)
}

type timerBoundary struct {
	Catch      timerCatch
	AttachedTo string
}

func timerBoundarySpec(ev element.BoundaryEvent) (timerBoundary, error) {
	if !ev.CancelActivity {
		return timerBoundary{}, fmt.Errorf("UNSUPPORTED_ELEMENT: boundaryEvent %q must be interrupting (cancelActivity=true)", ev.ID)
	}
	if strings.TrimSpace(ev.AttachedToRef) == "" {
		return timerBoundary{}, fmt.Errorf("UNSUPPORTED_ELEMENT: boundaryEvent %q needs attachedToRef", ev.ID)
	}
	if extraCatchDefinitions(ev.EventDefinitions) > 0 || len(ev.TimerEventDefinitions) == 0 {
		return timerBoundary{}, fmt.Errorf("UNSUPPORTED_ELEMENT: boundaryEvent %q must be an interrupting timer boundary", ev.ID)
	}
	catch, err := timerCatchFromDefs(ev.ID, ev.EventDefinitions)
	if err != nil {
		return timerBoundary{}, err
	}
	return timerBoundary{Catch: catch, AttachedTo: ev.AttachedToRef}, nil
}

func timerCatchFromDefs(id string, defs element.EventDefinitions) (timerCatch, error) {
	if extraCatchDefinitions(defs) > 0 {
		return timerCatch{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateCatchEvent %q is not a timer catch", id)
	}
	if len(defs.TimerEventDefinitions) != 1 {
		return timerCatch{}, fmt.Errorf("UNSUPPORTED_ELEMENT: %q needs exactly one timerEventDefinition", id)
	}
	def := defs.TimerEventDefinitions[0]
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
		return timerCatch{}, fmt.Errorf("UNSUPPORTED_ELEMENT: %q has multiple timer expressions", id)
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
		return timerCatch{}, fmt.Errorf("UNSUPPORTED_ELEMENT: %q needs timeDuration, timeDate, or timeCycle", id)
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

func extraCatchDefinitions(d element.EventDefinitions) int {
	n := 0
	n += len(d.MessageEventDefinitions)
	n += len(d.EscalationEventDefinitions)
	n += len(d.TerminateEventDefinitions)
	n += len(d.SignalEventDefinitions)
	n += len(d.ConditionalEventDefinitions)
	n += len(d.ErrorEventDefinitions)
	n += len(d.LinkEventDefinitions)
	n += len(d.CompensateEventDefinitions)
	return n
}

func otherCatchDefinitions(ev element.IntermediateCatchEvent) int {
	return extraCatchDefinitions(ev.EventDefinitions)
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
	for _, e := range d.Process.IntermediateCatchEvents {
		if e.ID != id {
			continue
		}
		if _, err := timerCatchSpec(e); err == nil {
			return CatchKindTimer, nil
		}
		return CatchKindMessage, nil
	}
	return "", fmt.Errorf("NOT_FOUND: intermediate catch %q", id)
}

// TimerBoundary returns the interrupting timer boundary attached to an activity.
func (d *Deployment) TimerBoundary(activityID string) (string, bool) {
	for _, e := range d.Process.BoundaryEvents {
		if e.AttachedToRef != activityID {
			continue
		}
		if _, err := timerBoundarySpec(e); err != nil {
			continue
		}
		return e.ID, true
	}
	return "", false
}

// AttachedActivity returns the activity a timer boundary is attached to.
func (d *Deployment) AttachedActivity(boundaryID string) (string, bool) {
	for _, e := range d.Process.BoundaryEvents {
		if e.ID == boundaryID && e.AttachedToRef != "" {
			return e.AttachedToRef, true
		}
	}
	return "", false
}
