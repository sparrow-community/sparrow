package deploy

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// ParseISO8601Duration parses a BPMN timeDuration subset: PTnHnMnS.
// Days, weeks, months, and years are not accepted.
func ParseISO8601Duration(s string) (time.Duration, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	if !strings.HasPrefix(s, "PT") {
		return 0, fmt.Errorf("INVALID_CONDITION: duration %q must be ISO-8601 time (PTnHnMnS)", s)
	}
	rest := s[2:]
	if rest == "" {
		return 0, fmt.Errorf("INVALID_CONDITION: duration %q has no time components", s)
	}

	var total time.Duration
	i := 0
	saw := false
	for i < len(rest) {
		start := i
		for i < len(rest) && (unicode.IsDigit(rune(rest[i])) || rest[i] == '.') {
			i++
		}
		if i == start || i >= len(rest) {
			return 0, fmt.Errorf("INVALID_CONDITION: duration %q is not ISO-8601 time (PTnHnMnS)", s)
		}
		n, err := strconv.ParseFloat(rest[start:i], 64)
		if err != nil {
			return 0, fmt.Errorf("INVALID_CONDITION: duration %q: %w", s, err)
		}
		unit := rest[i]
		i++
		var d time.Duration
		switch unit {
		case 'H':
			d = time.Duration(n * float64(time.Hour))
		case 'M':
			d = time.Duration(n * float64(time.Minute))
		case 'S':
			d = time.Duration(n * float64(time.Second))
		default:
			return 0, fmt.Errorf("INVALID_CONDITION: duration %q has unsupported unit %q", s, string(unit))
		}
		total += d
		saw = true
	}
	if !saw {
		return 0, fmt.Errorf("INVALID_CONDITION: duration %q has no time components", s)
	}
	return total, nil
}

// ParseISO8601Date parses a BPMN timeDate subset: RFC3339, or date-only (UTC).
func ParseISO8601Date(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, fmt.Errorf("INVALID_CONDITION: empty timeDate")
	}
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	var last error
	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t.UTC(), nil
		}
		last = err
	}
	return time.Time{}, fmt.Errorf("INVALID_CONDITION: timeDate %q is not ISO-8601 datetime: %w", s, last)
}

type cycleSpec struct {
	Repeat int // -1 unlimited; >= 1 finite (intermediate catch only uses first due)
	Start  time.Time
	End    time.Time
	Period time.Duration
}

// ParseISO8601Cycle parses a BPMN timeCycle subset:
//
//	R[n]/PTnHnMnS
//	R[n]/<timeDate>/PTnHnMnS
//	R[n]/PTnHnMnS/<timeDate>
//
// Intermediate catch uses only the first due instant; it does not re-arm.
func ParseISO8601Cycle(s string) (cycleSpec, error) {
	raw := strings.TrimSpace(s)
	if raw == "" {
		return cycleSpec{}, fmt.Errorf("INVALID_CONDITION: empty timeCycle")
	}
	if len(raw) < 1 || (raw[0] != 'R' && raw[0] != 'r') {
		return cycleSpec{}, fmt.Errorf("INVALID_CONDITION: timeCycle %q must start with R", s)
	}
	rest := raw[1:]
	i := 0
	for i < len(rest) && rest[i] >= '0' && rest[i] <= '9' {
		i++
	}
	repeat := -1
	if i > 0 {
		n, err := strconv.Atoi(rest[:i])
		if err != nil || n <= 0 {
			return cycleSpec{}, fmt.Errorf("INVALID_CONDITION: timeCycle %q has invalid repeat count", s)
		}
		repeat = n
		rest = rest[i:]
	}
	if !strings.HasPrefix(rest, "/") {
		return cycleSpec{}, fmt.Errorf("INVALID_CONDITION: timeCycle %q is not ISO-8601 repeating interval", s)
	}
	parts := strings.Split(rest[1:], "/")
	if len(parts) == 0 || len(parts) > 2 || parts[0] == "" {
		return cycleSpec{}, fmt.Errorf("INVALID_CONDITION: timeCycle %q is not ISO-8601 repeating interval", s)
	}
	out := cycleSpec{Repeat: repeat}
	if len(parts) == 1 {
		dur, err := ParseISO8601Duration(parts[0])
		if err != nil {
			return cycleSpec{}, fmt.Errorf("INVALID_CONDITION: timeCycle %q: %w", s, err)
		}
		out.Period = dur
		return out, nil
	}
	a, b := parts[0], parts[1]
	if isDurationPart(a) {
		dur, err := ParseISO8601Duration(a)
		if err != nil {
			return cycleSpec{}, err
		}
		end, err := ParseISO8601Date(b)
		if err != nil {
			return cycleSpec{}, fmt.Errorf("INVALID_CONDITION: timeCycle %q: %w", s, err)
		}
		out.Period = dur
		out.End = end
		return out, nil
	}
	start, err := ParseISO8601Date(a)
	if err != nil {
		return cycleSpec{}, fmt.Errorf("INVALID_CONDITION: timeCycle %q: %w", s, err)
	}
	dur, err := ParseISO8601Duration(b)
	if err != nil {
		return cycleSpec{}, fmt.Errorf("INVALID_CONDITION: timeCycle %q: %w", s, err)
	}
	out.Start = start
	out.Period = dur
	return out, nil
}

func isDurationPart(s string) bool {
	u := strings.ToUpper(strings.TrimSpace(s))
	return strings.HasPrefix(u, "PT")
}

func (c cycleSpec) FirstDue(now time.Time) (time.Time, error) {
	if now.IsZero() {
		now = time.Now()
	}
	var due time.Time
	if !c.Start.IsZero() {
		due = c.Start
		if due.Before(now) {
			if c.Period <= 0 {
				due = now
			} else {
				n := now.Sub(c.Start) / c.Period
				due = c.Start.Add(n * c.Period)
				if due.Before(now) {
					due = due.Add(c.Period)
				}
			}
		}
	} else {
		due = now.Add(c.Period)
	}
	if !c.End.IsZero() && due.After(c.End) {
		return time.Time{}, fmt.Errorf("INVALID_CONDITION: timeCycle has no occurrence after now")
	}
	return due, nil
}
