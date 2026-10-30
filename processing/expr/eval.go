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

package expr

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
)

const notFn = "__sparrow_not"

var (
	reGetDataObject   = regexp.MustCompile(`(?i)bpmn:getDataObject\(\s*'([^']+)'\s*\)`)
	reGetDataObjectDQ = regexp.MustCompile(`(?i)bpmn:getDataObject\(\s*"([^"]+)"\s*\)`)
	reFeelSome        = regexp.MustCompile(`(?i)^some\s+([A-Za-z_][\w]*)\s+in\s+(\S+)\s+satisfies\s+(.+)$`)
	reFeelEvery       = regexp.MustCompile(`(?i)^every\s+([A-Za-z_][\w]*)\s+in\s+(\S+)\s+satisfies\s+(.+)$`)
)

// Eval evaluates a condition against instance variables (name → json_value).
// BPMN `${...}` wrappers are stripped. Evaluation uses github.com/expr-lang/expr
// so comparisons, boolean ops, and property access are available.
// Single-quoted strings are accepted as in typical BPMN conditions.
// Missing variables are treated as nil/false, so `!missing` is true.
//
// Interop sugar (Data Objects stay Excluded as ledger subjects — variables carry data):
//   - bpmn:getDataObject('name') → process variable name
//   - leading FEEL "=" expression marker stripped
//   - FEEL "=" equality rewritten to "=="
//   - FEEL not(x) available as an env function
//   - FEEL some/every … satisfies → expr any/all
func Eval(text string, vars map[string]string) (bool, error) {
	s := unwrap(text)
	if s == "" {
		return false, fmt.Errorf("empty expression")
	}
	s = rewriteInterop(s)
	s, aliases := rewriteSpacedIdentifiers(s, vars)
	s = normalizeSingleQuotes(s)

	env := envFrom(vars)
	for k, v := range aliases {
		env[k] = v
	}
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
	s = rewriteInterop(s)
	s, aliases := rewriteSpacedIdentifiers(s, vars)
	s = normalizeSingleQuotes(s)

	env := envFrom(vars)
	for k, v := range aliases {
		env[k] = v
	}
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

func rewriteInterop(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "=") && !strings.HasPrefix(s, "==") {
		s = strings.TrimSpace(s[1:])
	}
	s = reGetDataObject.ReplaceAllString(s, "$1")
	s = reGetDataObjectDQ.ReplaceAllString(s, "$1")
	s = rewriteFeelQuantifiers(s)
	return rewriteFeelEquality(s)
}

// rewriteFeelQuantifiers maps FEEL `some/every x in coll satisfies pred` to
// expr-lang `any/all(coll, {pred with x → #})`.
func rewriteFeelQuantifiers(s string) string {
	s = strings.TrimSpace(s)
	if m := reFeelSome.FindStringSubmatch(s); m != nil {
		return "any(" + m[2] + ", {" + replaceBoundIdent(m[3], m[1], "#") + "})"
	}
	if m := reFeelEvery.FindStringSubmatch(s); m != nil {
		return "all(" + m[2] + ", {" + replaceBoundIdent(m[3], m[1], "#") + "})"
	}
	return s
}

func replaceBoundIdent(pred, bound, repl string) string {
	re := regexp.MustCompile(`\b` + regexp.QuoteMeta(bound) + `\b`)
	return re.ReplaceAllString(pred, repl)
}

// rewriteSpacedIdentifiers turns `Vacation Approval == "x"` into a safe alias
// when the multi-word name matches a process variable (or is introduced as nil).
var reSpacedIdent = regexp.MustCompile(`\b([A-Za-z_][A-Za-z0-9_]*(?:[ \t]+[A-Za-z_][A-Za-z0-9_]*)+)\b`)

func rewriteSpacedIdentifiers(s string, vars map[string]string) (string, map[string]any) {
	aliases := map[string]any{}
	out := reSpacedIdent.ReplaceAllStringFunc(s, func(name string) string {
		alias := "__spaced_" + strings.ReplaceAll(name, " ", "_")
		if raw, ok := vars[name]; ok {
			var v any
			if err := json.Unmarshal([]byte(raw), &v); err != nil {
				aliases[alias] = raw
			} else {
				aliases[alias] = v
			}
		} else {
			aliases[alias] = nil
		}
		return alias
	})
	return out, aliases
}

func rewriteFeelEquality(s string) string {
	var b strings.Builder
	b.Grow(len(s) + 8)
	inSingle, inDouble := false, false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if inSingle {
			b.WriteByte(c)
			if c == '\'' {
				inSingle = false
			}
			continue
		}
		if inDouble {
			b.WriteByte(c)
			if c == '"' && (i == 0 || s[i-1] != '\\') {
				inDouble = false
			}
			continue
		}
		switch c {
		case '\'':
			inSingle = true
			b.WriteByte(c)
		case '"':
			inDouble = true
			b.WriteByte(c)
		case '!':
			if i+1 < len(s) && s[i+1] == '=' {
				b.WriteString("!=")
				i++
				continue
			}
			b.WriteByte(c)
		case '<', '>':
			b.WriteByte(c)
			if i+1 < len(s) && s[i+1] == '=' {
				b.WriteByte('=')
				i++
			}
		case '=':
			if i+1 < len(s) && s[i+1] == '=' {
				b.WriteString("==")
				i++
				continue
			}
			b.WriteString("==")
		default:
			b.WriteByte(c)
		}
	}
	return b.String()
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
	out := map[string]any{
		notFn: sparrowNot,
		"not": sparrowNot,
	}
	for k, raw := range vars {
		if k == notFn || k == "not" {
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
