// Copyright 2026 The Sparrow community and contributors
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

//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

func argString(args []js.Value, i int, name string) (string, error) {
	if len(args) <= i || args[i].IsUndefined() || args[i].IsNull() {
		return "", fmt.Errorf("INVALID_ARGUMENT: %s is required", name)
	}
	if args[i].Type() != js.TypeString {
		return "", fmt.Errorf("INVALID_ARGUMENT: %s must be a string", name)
	}
	return args[i].String(), nil
}

func argObject(args []js.Value, i int) (map[string]any, error) {
	if len(args) <= i || args[i].IsUndefined() || args[i].IsNull() {
		return nil, fmt.Errorf("INVALID_ARGUMENT: request object is required")
	}
	raw, err := jsValueToJSON(args[i])
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("INVALID_ARGUMENT: %w", err)
	}
	return m, nil
}

func jsValueToJSON(v js.Value) ([]byte, error) {
	jsonObj := js.Global().Get("JSON")
	s := jsonObj.Call("stringify", v).String()
	return []byte(s), nil
}

func toJS(v any) any {
	if v == nil {
		return js.Null()
	}
	switch x := v.(type) {
	case bool, int, int32, int64, float64, string:
		return x
	case []string:
		arr := js.Global().Get("Array").New(len(x))
		for i, s := range x {
			arr.SetIndex(i, s)
		}
		return arr
	default:
		b, err := json.Marshal(x)
		if err != nil {
			return js.Null()
		}
		return js.Global().Get("JSON").Call("parse", string(b))
	}
}

func asString(v any) string {
	if v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	case float64:
		return fmt.Sprintf("%v", x)
	case bool:
		return fmt.Sprintf("%v", x)
	default:
		return fmt.Sprintf("%v", x)
	}
}

func asBool(v any) bool {
	b, ok := v.(bool)
	return ok && b
}

func asInt(v any) int {
	return int(asInt64(v))
}

func asInt32(v any) int32 {
	return int32(asInt64(v))
}

func asInt64(v any) int64 {
	if v == nil {
		return 0
	}
	switch x := v.(type) {
	case float64:
		return int64(x)
	case int64:
		return x
	case int:
		return int64(x)
	case json.Number:
		n, _ := x.Int64()
		return n
	default:
		return 0
	}
}

func asAnyMap(v any) (map[string]any, error) {
	if v == nil {
		return nil, nil
	}
	m, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected object")
	}
	return m, nil
}
