package expr

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
)

const notFn = "__sparrow_not"

// Eval evaluates a condition against instance variables (name → json_value).
// BPMN `${...}` wrappers are stripped. Evaluation uses github.com/expr-lang/expr
// so comparisons, boolean ops, and property access are available.
// Single-quoted strings are accepted as in typical BPMN conditions.
// Missing variables are treated as nil/false, so `!missing` is true.
func Eval(text string, vars map[string]string) (bool, error) {
	s := unwrap(text)
	if s == "" {
		return false, fmt.Errorf("empty expression")
	}
	s = normalizeSingleQuotes(s)

	env := envFrom(vars)
	program, err := expr.Compile(s,
		expr.Env(env),
		expr.AllowUndefinedVariables(),
		expr.Patch(notPatcher{}),
	)
	if err != nil {
		return false, fmt.Errorf("compile %q: %w", text, err)
	}
	out, err := expr.Run(program, env)
	if err != nil {
		return false, fmt.Errorf("eval %q: %w", text, err)
	}
	return asBool(out)
}

// EvalJSON evaluates an expression against instance variables and returns the
// result as a JSON value string suitable for process variables.
func EvalJSON(text string, vars map[string]string) (string, error) {
	s := unwrap(text)
	if s == "" {
		return "", fmt.Errorf("empty expression")
	}
	s = normalizeSingleQuotes(s)

	env := envFrom(vars)
	program, err := expr.Compile(s,
		expr.Env(env),
		expr.AllowUndefinedVariables(),
		expr.Patch(notPatcher{}),
	)
	if err != nil {
		return "", fmt.Errorf("compile %q: %w", text, err)
	}
	out, err := expr.Run(program, env)
	if err != nil {
		return "", fmt.Errorf("eval %q: %w", text, err)
	}
	b, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("marshal %q: %w", text, err)
	}
	return string(b), nil
}

type notPatcher struct{}

func (notPatcher) Visit(node *ast.Node) {
	n, ok := (*node).(*ast.UnaryNode)
	if !ok || (n.Operator != "!" && n.Operator != "not") {
		return
	}
	ast.Patch(node, &ast.CallNode{
		Callee:    &ast.IdentifierNode{Value: notFn},
		Arguments: []ast.Node{n.Node},
	})
}

func sparrowNot(v any) bool {
	return !isTruthy(v)
}

func unwrap(text string) string {
	s := strings.TrimSpace(text)
	if strings.HasPrefix(s, "${") && strings.HasSuffix(s, "}") {
		return strings.TrimSpace(s[2 : len(s)-1])
	}
	return s
}

// normalizeSingleQuotes turns 'alice' into "alice" without touching double-quoted spans.
func normalizeSingleQuotes(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	inDouble := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if inDouble {
			b.WriteByte(c)
			if c == '"' && (i == 0 || s[i-1] != '\\') {
				inDouble = false
			}
			continue
		}
		if c == '"' {
			inDouble = true
			b.WriteByte(c)
			continue
		}
		if c == '\'' {
			j := i + 1
			for j < len(s) && s[j] != '\'' {
				j++
			}
			if j >= len(s) {
				b.WriteByte(c)
				continue
			}
			b.WriteByte('"')
			b.WriteString(s[i+1 : j])
			b.WriteByte('"')
			i = j
			continue
		}
		b.WriteByte(c)
	}
	return b.String()
}

func envFrom(vars map[string]string) map[string]any {
	out := map[string]any{notFn: sparrowNot}
	for k, raw := range vars {
		if k == notFn {
			continue
		}
		var v any
		if err := json.Unmarshal([]byte(raw), &v); err != nil {
			out[k] = raw
			continue
		}
		out[k] = v
	}
	return out
}

func isTruthy(v any) bool {
	switch t := v.(type) {
	case nil:
		return false
	case bool:
		return t
	case float64:
		return t != 0
	case int:
		return t != 0
	case int64:
		return t != 0
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

func asBool(v any) (bool, error) {
	switch t := v.(type) {
	case bool:
		return t, nil
	case nil:
		return false, nil
	default:
		return false, fmt.Errorf("condition did not return a boolean (got %T)", v)
	}
}
