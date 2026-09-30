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

package deploy

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// ParseISO8601Duration parses a BPMN timeDuration subset:
//
//	PTnHnMnS
//	P[nW][nD][TnHnMnS]
//
// Day and week components use fixed 24h units (no calendar/DST adjustment).
// Months and years are rejected until calendar arithmetic is specified.
func ParseISO8601Duration(s string) (time.Duration, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	if !strings.HasPrefix(s, "P") {
		return 0, fmt.Errorf("INVALID_CONDITION: duration %q must be ISO-8601 (PnD / PTnHnMnS)", s)
	}
	rest := s[1:]
	if rest == "" {
		return 0, fmt.Errorf("INVALID_CONDITION: duration %q has no components", s)
	}

	var total time.Duration
	saw := false

	// Date part before optional T: nW and/or nD.
	if !strings.HasPrefix(rest, "T") {
		i := 0
		for i < len(rest) && rest[i] != 'T' {
			start := i
			for i < len(rest) && (unicode.IsDigit(rune(rest[i])) || rest[i] == '.') {
				i++
			}
			if i == start || i >= len(rest) || rest[i] == 'T' {
				return 0, fmt.Errorf("INVALID_CONDITION: duration %q is not ISO-8601 (PnD / PTnHnMnS)", s)
			}
			n, err := strconv.ParseFloat(rest[start:i], 64)
			if err != nil {
				return 0, fmt.Errorf("INVALID_CONDITION: duration %q: %w", s, err)
			}
			unit := rest[i]
			i++
			switch unit {
			case 'W':
				total += time.Duration(n * float64(7*24*time.Hour))
			case 'D':
				total += time.Duration(n * float64(24*time.Hour))
			case 'Y', 'M':
				return 0, fmt.Errorf("INVALID_CONDITION: duration %q has unsupported calendar unit %q", s, string(unit))
			default:
				return 0, fmt.Errorf("INVALID_CONDITION: duration %q has unsupported unit %q", s, string(unit))
			}
			saw = true
		}
		rest = rest[i:]
	}

	if rest == "" {
		if !saw {
			return 0, fmt.Errorf("INVALID_CONDITION: duration %q has no components", s)
		}
		return total, nil
	}
	if !strings.HasPrefix(rest, "T") {
		return 0, fmt.Errorf("INVALID_CONDITION: duration %q is not ISO-8601 (PnD / PTnHnMnS)", s)
	}
	rest = rest[1:]
	if rest == "" {
		if !saw {
			return 0, fmt.Errorf("INVALID_CONDITION: duration %q has no time components", s)
		}
		return total, nil
	}

	i := 0
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
		return 0, fmt.Errorf("INVALID_CONDITION: duration %q has no components", s)
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
//	R[n]/P…duration…
//	R[n]/<timeDate>/P…duration…
//	R[n]/P…duration…/<timeDate>
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
	return strings.HasPrefix(u, "P")
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

func (c cycleSpec) NextDue(prev time.Time) (time.Time, bool) {
	if prev.IsZero() || c.Period < 0 {
		return time.Time{}, false
	}
	if c.Repeat == 1 {
		return time.Time{}, false
	}
	next := prev.Add(c.Period)
	if !c.End.IsZero() && next.After(c.End) {
		return time.Time{}, false
	}
	return next, true
}

func (c cycleSpec) Rearmed() (cycleSpec, bool) {
	if c.Repeat == 1 {
		return cycleSpec{}, false
	}
	if c.Repeat > 1 {
		c.Repeat--
	}
	return c, true
}

func (c cycleSpec) String() string {
	repeat := "R"
	if c.Repeat > 0 {
		repeat = fmt.Sprintf("R%d", c.Repeat)
	}
	if !c.Start.IsZero() {
		return repeat + "/" + c.Start.UTC().Format(time.RFC3339) + "/" + formatISO8601Duration(c.Period)
	}
	if !c.End.IsZero() {
		return repeat + "/" + formatISO8601Duration(c.Period) + "/" + c.End.UTC().Format(time.RFC3339)
	}
	return repeat + "/" + formatISO8601Duration(c.Period)
}

func formatISO8601Duration(d time.Duration) string {
	if d == 0 {
		return "PT0S"
	}
	if d < 0 {
		d = -d
	}
	// Prefer day form when the duration is an exact multiple of 24h and ≥ 1 day.
	if d%time.Hour == 0 {
		hours := d / time.Hour
		if hours >= 24 && hours%24 == 0 {
			days := hours / 24
			if days%7 == 0 && days >= 7 {
				return "P" + strconv.FormatInt(int64(days/7), 10) + "W"
			}
			return "P" + strconv.FormatInt(int64(days), 10) + "D"
		}
	}
	var b strings.Builder
	b.WriteString("PT")
	h := d / time.Hour
	if h > 0 {
		b.WriteString(strconv.FormatInt(int64(h), 10))
		b.WriteByte('H')
		d -= h * time.Hour
	}
	m := d / time.Minute
	if m > 0 {
		b.WriteString(strconv.FormatInt(int64(m), 10))
		b.WriteByte('M')
		d -= m * time.Minute
	}
	s := d / time.Second
	if s > 0 || b.Len() == 2 {
		b.WriteString(strconv.FormatInt(int64(s), 10))
		b.WriteByte('S')
	}
	return b.String()
}

func NextCycleTimer(text string, prevDue time.Time) (nextText string, nextDue time.Time, ok bool, err error) {
	spec, err := ParseISO8601Cycle(text)
	if err != nil {
		return "", time.Time{}, false, err
	}
	nextDue, ok = spec.NextDue(prevDue)
	if !ok {
		return "", time.Time{}, false, nil
	}
	spec, ok = spec.Rearmed()
	if !ok {
		return "", time.Time{}, false, nil
	}
	return spec.String(), nextDue, true, nil
}
