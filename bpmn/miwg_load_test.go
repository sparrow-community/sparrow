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

package bpmn

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestMIWGFixturesLoad ensures every BPMN MIWG fixture under test/ parses.
// Structural round-trip assertions for individual cases live in *-export_test.go /
// *-roundtrip_test.go; this table covers cases that lack those (and C.10.0).
func TestMIWGFixturesLoad(t *testing.T) {
	dir := filepath.Join("test")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".bpmn") {
			continue
		}
		names = append(names, e.Name())
	}
	if len(names) == 0 {
		t.Fatal("no MIWG fixtures under test/")
	}
	for _, name := range names {
		name := name
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(dir, name)
			model, err := BpmnModelelementFromFile(path)
			if err != nil {
				t.Fatalf("parse %s: %v", path, err)
			}
			if model == nil || model.Definitions == nil {
				t.Fatalf("nil model for %s", path)
			}
			if len(model.Definitions.Processes) == 0 {
				t.Fatalf("%s: no processes", path)
			}
		})
	}
}
