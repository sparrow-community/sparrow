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

package element

// StartEvent is a catch event that starts a process or event sub-process.
type StartEvent struct {
	CatchEvent
	// Interrupting maps BPMN isInterrupting. Nil means attribute absent
	// (BPMN default true for event sub-process starts).
	Interrupting *bool `xml:"isInterrupting,attr"`
}

// IsInterrupting returns the effective interrupting flag (default true).
func (s StartEvent) IsInterrupting() bool {
	if s.Interrupting == nil {
		return true
	}
	return *s.Interrupting
}
