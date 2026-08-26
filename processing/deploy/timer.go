package deploy

import (
	"fmt"
	"strings"
	"time"

	"github.com/sparrow-community/sparrow/bpmn/element"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// CatchKind classifies an intermediateCatchEvent for runtime handlers.
type CatchKind string

const (
	CatchKindTimer   CatchKind = "timer"
	CatchKindMessage CatchKind = "message"
	CatchKindSignal  CatchKind = "signal"
	CatchKindError   CatchKind = "error"
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
	Catch        timerCatch
	AttachedTo   string
	Interrupting bool
}

type messageBoundary struct {
	Name         string
	AttachedTo   string
	Interrupting bool
}

func requireBoundaryAttach(ev element.BoundaryEvent) error {
	if strings.TrimSpace(ev.AttachedToRef) == "" {
		return fmt.Errorf("UNSUPPORTED_ELEMENT: boundaryEvent %q needs attachedToRef", ev.ID)
	}
	return nil
}

func timerBoundarySpec(ev element.BoundaryEvent) (timerBoundary, error) {
	if err := requireBoundaryAttach(ev); err != nil {
		return timerBoundary{}, err
	}
	if extraCatchDefinitions(ev.EventDefinitions) > 0 || len(ev.TimerEventDefinitions) == 0 {
		return timerBoundary{}, fmt.Errorf("UNSUPPORTED_ELEMENT: boundaryEvent %q must be a timer boundary", ev.ID)
	}
	catch, err := timerCatchFromDefs(ev.ID, ev.EventDefinitions)
	if err != nil {
		return timerBoundary{}, err
	}
	return timerBoundary{Catch: catch, AttachedTo: ev.AttachedToRef, Interrupting: ev.CancelActivity}, nil
}

func messageBoundarySpec(ev element.BoundaryEvent, messages []element.Message) (messageBoundary, error) {
	if err := requireBoundaryAttach(ev); err != nil {
		return messageBoundary{}, err
	}
	if len(ev.TimerEventDefinitions) > 0 {
		return messageBoundary{}, fmt.Errorf("UNSUPPORTED_ELEMENT: boundaryEvent %q must be a message boundary", ev.ID)
	}
	spec, err := messageCatchFromDefs(ev.ID, ev.EventDefinitions, messages, strings.TrimSpace(ev.Name))
	if err != nil {
		return messageBoundary{}, err
	}
	return messageBoundary{Name: spec.Name, AttachedTo: ev.AttachedToRef, Interrupting: ev.CancelActivity}, nil
}

type signalBoundary struct {
	Name         string
	AttachedTo   string
	Interrupting bool
}

func signalBoundarySpec(ev element.BoundaryEvent, signals []element.Signal) (signalBoundary, error) {
	if err := requireBoundaryAttach(ev); err != nil {
		return signalBoundary{}, err
	}
	if len(ev.TimerEventDefinitions) > 0 {
		return signalBoundary{}, fmt.Errorf("UNSUPPORTED_ELEMENT: boundaryEvent %q must be a signal boundary", ev.ID)
	}
	name, err := signalCatchFromDefs(ev.ID, ev.EventDefinitions, signals, strings.TrimSpace(ev.Name))
	if err != nil {
		return signalBoundary{}, err
	}
	return signalBoundary{Name: name, AttachedTo: ev.AttachedToRef, Interrupting: ev.CancelActivity}, nil
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
	return messageCatchFromDefs(ev.ID, ev.EventDefinitions, messages, strings.TrimSpace(ev.Name))
}

func messageCatchFromDefs(id string, defs element.EventDefinitions, messages []element.Message, fallbackName string) (messageCatch, error) {
	if len(defs.TimerEventDefinitions) > 0 {
		return messageCatch{}, fmt.Errorf("UNSUPPORTED_ELEMENT: %q is not a message catch (has timerEventDefinition)", id)
	}
	if len(defs.MessageEventDefinitions) != 1 {
		return messageCatch{}, fmt.Errorf("UNSUPPORTED_ELEMENT: %q needs exactly one messageEventDefinition", id)
	}
	other := extraCatchDefinitions(defs) - len(defs.MessageEventDefinitions)
	if other > 0 {
		return messageCatch{}, fmt.Errorf("UNSUPPORTED_ELEMENT: %q is not a message catch", id)
	}
	name := resolveMessageName(defs.MessageEventDefinitions[0].MessageRef, messages)
	if name == "" {
		name = fallbackName
	}
	if name == "" {
		name = id
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
	typ, err := d.TypeOf(id)
	if err != nil || typ != eventv1.Element_TYPE_INTERMEDIATE_CATCH_EVENT {
		return "", fmt.Errorf("NOT_FOUND: intermediate catch %q", id)
	}
	if _, ok := d.timerCatch[id]; ok {
		return CatchKindTimer, nil
	}
	if _, ok := d.messageCatch[id]; ok {
		return CatchKindMessage, nil
	}
	if _, ok := d.signalCatch[id]; ok {
		return CatchKindSignal, nil
	}
	return "", fmt.Errorf("NOT_FOUND: intermediate catch %q", id)
}

// TimerBoundary returns the timer boundary attached to an activity.
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

// MessageBoundary returns the message boundary attached to an activity.
func (d *Deployment) MessageBoundary(activityID string) (string, bool) {
	for _, e := range d.Process.BoundaryEvents {
		if e.AttachedToRef != activityID {
			continue
		}
		if _, ok := d.messageCatch[e.ID]; ok {
			return e.ID, true
		}
	}
	return "", false
}

// AttachedBoundary returns the single boundary attached to an activity.
func (d *Deployment) AttachedBoundary(activityID string) (string, bool) {
	if id, ok := d.TimerBoundary(activityID); ok {
		return id, true
	}
	if id, ok := d.MessageBoundary(activityID); ok {
		return id, true
	}
	return d.SignalBoundary(activityID)
}

// SignalBoundary returns the signal boundary attached to an activity.
func (d *Deployment) SignalBoundary(activityID string) (string, bool) {
	for _, e := range d.Process.BoundaryEvents {
		if e.AttachedToRef != activityID {
			continue
		}
		if _, ok := d.signalCatch[e.ID]; ok {
			return e.ID, true
		}
	}
	return "", false
}

// InterruptingBoundary returns the attached boundary when it is interrupting.
func (d *Deployment) InterruptingBoundary(activityID string) (string, bool) {
	id, ok := d.AttachedBoundary(activityID)
	if !ok || d.BoundaryInterrupting(id) {
		return id, ok
	}
	return "", false
}

// BoundaryInterrupting reports whether a boundary cancels its host activity when it fires.
func (d *Deployment) BoundaryInterrupting(boundaryID string) bool {
	for _, e := range d.Process.BoundaryEvents {
		if e.ID == boundaryID {
			return e.CancelActivity
		}
	}
	return true
}

// AttachedActivity returns the activity a timer or message boundary is attached to.
func (d *Deployment) AttachedActivity(boundaryID string) (string, bool) {
	for _, e := range d.Process.BoundaryEvents {
		if e.ID == boundaryID && e.AttachedToRef != "" {
			return e.AttachedToRef, true
		}
	}
	return "", false
}
