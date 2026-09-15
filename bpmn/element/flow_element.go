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

type FlowElement struct {
	BaseElement
	Name              string     `xml:"name,attr"`
	Auditing          Auditing   `xml:"auditing"`
	Monitoring        Monitoring `xml:"monitoring"`
	CategoryValueRefs []string   `xml:"categoryValueRef"`
}

type FlowElements struct {
	StartEvents             []StartEvent             `xml:"startEvent"`
	EndEvents               []EndEvent               `xml:"endEvent"`
	Tasks                   []Task                   `xml:"task"`
	ManualTasks             []ManualTask             `xml:"manualTask"`
	UserTasks               []UserTask               `xml:"userTask"`
	ServiceTasks            []ServiceTask            `xml:"serviceTask"`
	SendTasks               []SendTask               `xml:"sendTask"`
	ReceiveTasks            []ReceiveTask            `xml:"receiveTask"`
	BusinessRuleTasks       []BusinessRuleTask       `xml:"businessRuleTask"`
	ScriptTasks             []ScriptTask             `xml:"scriptTask"`
	SequenceFlows           []SequenceFlow           `xml:"sequenceFlow"`
	DataStoreReferences     []DataStoreReference     `xml:"dataStoreReference"`
	ParallelGatewaies       []ParallelGateway        `xml:"parallelGateway"`
	ExclusiveGatewaies      []ExclusiveGateway       `xml:"exclusiveGateway"`
	InclusiveGatewaies      []InclusiveGateway       `xml:"inclusiveGateway"`
	ComplexGatewaies        []ComplexGateway         `xml:"complexGateway"`
	EventBasedGatewaies     []EventBasedGateway      `xml:"eventBasedGateway"`
	SubProcesses            []SubProcess             `xml:"subProcess"`
	Transactions            []Transaction            `xml:"transaction"`
	BoundaryEvents          []BoundaryEvent          `xml:"boundaryEvent"`
	CallActivities          []CallActivity           `xml:"callActivity"`
	DataObjectReferenes     []DataObjectReference    `xml:"dataObjectReference"`
	DataObjects             []DataObject             `xml:"dataObject"`
	IntermediateThrowEvents []IntermediateThrowEvent `xml:"intermediateThrowEvent"`
	IntermediateCatchEvents []IntermediateCatchEvent `xml:"intermediateCatchEvent"`
	ImplicitThrowEvents     []ImplicitThrowEvent     `xml:"implicitThrowEvent"`
}
