package deploy

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// ParseISO8601Duration parses a BPMN timeDuration subset: PTnHnMnS.
// Days, weeks, months, years, and timeDate/timeCycle are not accepted.
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
