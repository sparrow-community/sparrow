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

type EventDefinition struct {
	RootElement
}

type EventDefinitions struct {
	MessageEventDefinitions     []MessageEventDefinition     `xml:"messageEventDefinition"`
	EscalationEventDefinitions  []EscalationEventDefinition  `xml:"escalationEventDefinition"`
	TimerEventDefinitions       []TimerEventDefinition       `xml:"timerEventDefinition"`
	TerminateEventDefinitions   []TerminateEventDefinition   `xml:"terminateEventDefinition"`
	SignalEventDefinitions      []SignalEventDefinition      `xml:"signalEventDefinition"`
	ConditionalEventDefinitions []ConditionalEventDefinition `xml:"conditionalEventDefinition"`
	ErrorEventDefinitions       []ErrorEventDefinition       `xml:"errorEventDefinition"`
	LinkEventDefinitions        []LinkEventDefinition        `xml:"linkEventDefinition"`
	CompensateEventDefinitions  []CompensateEventDefinition  `xml:"compensateEventDefinition"`
	CancelEventDefinitions      []CancelEventDefinition      `xml:"cancelEventDefinition"`
}
