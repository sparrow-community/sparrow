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

type TProcessType string

const (
	TProcessTypeNone    TProcessType = "None"
	TProcessTypePublic  TProcessType = "Public"
	TProcessTypePrivate TProcessType = "Private"
)

type Process struct {
	CallableElement
	ProcessType  TProcessType `xml:"processType,attr,omitempty"`
	IsClosed     bool         `xml:"isClosed,attr,omitempty"`
	IsExecutable bool         `xml:"isExecutable,attr,omitempty"`
	Auditing     Auditing     `xml:"auditing,omitempty"`
	Monitoring   Monitoring   `xml:"monitoring,omitempty"`
	Properties   []Property   `xml:"property,omitempty"`
	LaneSets     []LaneSet    `xml:"laneSet,omitempty"`
	FlowElements
	Artifacts
	ResourceRoles            []ResourceRole            `xml:"resourceRole,omitempty"`
	CorrelationSubscriptions []CorrelationSubscription `xml:"correlationSubscription,omitempty"`
	Supports                 []string                  `xml:"supports,omitempty"`
}
