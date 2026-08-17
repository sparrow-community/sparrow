package expr

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Eval evaluates an M1 condition against instance variables (name → json_value).
// Supported forms (optional ${...} wrapper):
//
//	ident
//	!ident
//	ident == literal
//	ident != literal
//
// Literals: true, false, null, numbers, 'strings' or "strings".
func Eval(text string, vars map[string]string) (bool, error) {
	s := unwrap(text)
	if s == "" {
		return false, fmt.Errorf("empty expression")
	}

	neg := false
	if strings.HasPrefix(s, "!") {
		neg = true
		s = strings.TrimSpace(s[1:])
	}

	ident, rest := splitIdent(s)
	if ident == "" {
		return false, fmt.Errorf("expected identifier in %q", text)
	}
	rest = strings.TrimSpace(rest)
	if rest == "" {
		ok := truthy(vars, ident)
		if neg {
			return !ok, nil
		}
		return ok, nil
	}
	if neg {
		return false, fmt.Errorf("! cannot be combined with comparison in %q", text)
	}

	op := ""
	switch {
	case strings.HasPrefix(rest, "=="):
		op = "=="
		rest = strings.TrimSpace(rest[2:])
	case strings.HasPrefix(rest, "!="):
		op = "!="
		rest = strings.TrimSpace(rest[2:])
	default:
		return false, fmt.Errorf("unsupported expression %q", text)
	}

	want, err := parseLiteral(rest)
	if err != nil {
		return false, fmt.Errorf("literal in %q: %w", text, err)
	}
	got, ok := lookup(vars, ident)
	eq := ok && equal(got, want)
	if op == "!=" {
		return !eq, nil
	}
	return eq, nil
}

func unwrap(text string) string {
	s := strings.TrimSpace(text)
	if strings.HasPrefix(s, "${") && strings.HasSuffix(s, "}") {
		return strings.TrimSpace(s[2 : len(s)-1])
	}
	return s
}

func splitIdent(s string) (ident, rest string) {
	if s == "" {
		return "", ""
	}
	if !isIdentStart(rune(s[0])) {
		return "", s
	}
	i := 1
	for i < len(s) && isIdentPart(rune(s[i])) {
		i++
	}
	return s[:i], s[i:]
}

func isIdentStart(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func isIdentPart(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func parseLiteral(s string) (any, error) {
	s = strings.TrimSpace(s)
	switch s {
	case "true":
		return true, nil
	case "false":
		return false, nil
	case "null":
		return nil, nil
	}
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1], nil
		}
	}
	if n, err := strconv.ParseFloat(s, 64); err == nil {
		return n, nil
	}
	return nil, fmt.Errorf("bad literal %q", s)
}

func lookup(vars map[string]string, name string) (any, bool) {
	if vars == nil {
		return nil, false
	}
	raw, ok := vars[name]
	if !ok {
		return nil, false
	}
	var v any
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return raw, true
	}
	return v, true
}

func truthy(vars map[string]string, name string) bool {
	v, ok := lookup(vars, name)
	if !ok {
		return false
	}
	switch t := v.(type) {
	case nil:
		return false
	case bool:
		return t
	case float64:
		return t != 0
	case json.Number:
		n, _ := t.Float64()
		return n != 0
	case string:
		return t != ""
	case []any:
		return len(t) > 0
	case map[string]any:
		return len(t) > 0
	default:
		return true
	}
}

func equal(a, b any) bool {
	af, aok := asFloat(a)
	bf, bok := asFloat(b)
	if aok && bok {
		return af == bf
	}
	return fmt.Sprint(a) == fmt.Sprint(b) && typeKind(a) == typeKind(b)
}

func asFloat(v any) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case json.Number:
		n, err := t.Float64()
		return n, err == nil
	case int:
		return float64(t), true
	default:
		return 0, false
	}
}

func typeKind(v any) string {
	switch v.(type) {
	case nil:
		return "null"
	case bool:
		return "bool"
	case string:
		return "string"
	default:
		if _, ok := asFloat(v); ok {
			return "number"
		}
		return fmt.Sprintf("%T", v)
	}
}
