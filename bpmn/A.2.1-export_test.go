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
	"encoding/xml"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sparrow-community/sparrow/bpmn/element"
)

func TestA_2_1_export(t *testing.T) {
	// create test using ./test/A.2.1-export.bpmn
	path := "./test/A.2.1-export.bpmn"
	modelelement, err := BpmnModelelementFromFile(path)
	if err != nil {
		t.Fatalf("BpmnModelelementFromFile error %s: %s", path, err)
	}
	if modelelement == nil {
		t.Fatalf("modelelement is nil, %s", path)
	}

	expected := &Bpmn{
		Definitions: &element.Definitions{
			XMLName: xml.Name{
				Space: "http://www.omg.org/spec/BPMN/20100524/MODEL",
				Local: "definitions",
			},
			ID:              "Definitions_1mv2yrn",
			TargetNamespace: "http://bpmn.io/schema/bpmn",
			Exporter:        "Camunda Modeler",
			ExporterVersion: "5.23.0",
			RootElemnts: element.RootElemnts{
				Processes: []element.Process{
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "Process_05abo3f",
								},
							},
						},
						IsExecutable: true,
						FlowElements: element.FlowElements{
							StartEvents: []element.StartEvent{{
								CatchEvent: element.CatchEvent{
									Event: element.Event{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "StartEvent_1",
												},
												Name: "Start Event",
											},
											Outgoing: []string{"Flow_0lacnkf"},
										},
									},
								},
							}},
							Tasks: []element.Task{
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_0ahdk3x",
												},
												Name: "Task 1",
											},
											Incoming: []string{"Flow_0lacnkf"},
											Outgoing: []string{"Flow_187opjq"},
										},
									},
								},
								{
									Activity: element.Activity{
										Default: "Flow_1sd2kpj",
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_172ndxy",
												},
												Name: "Task 2",
											},
											Incoming: []string{"Flow_194jx6p"},
											Outgoing: []string{"Flow_1sd2kpj", "Flow_0zwjy3h"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_1lz0l07",
												},
												Name: "Task 3",
											},
											Incoming: []string{"Flow_1xhc3bf", "Flow_1sd2kpj", "Flow_11m9fcl"},
											Outgoing: []string{"Flow_1a24ve1"},
										},
									},
								},
								{
									Activity: element.Activity{
										Default: "Flow_11m9fcl",
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_171qefk",
												},
												Name: "Task 4",
											},
											Incoming: []string{"Flow_19m0ydj"},
											Outgoing: []string{"Flow_11m9fcl", "Flow_01ckxme"},
										},
									},
								},
							},
							ExclusiveGatewaies: []element.ExclusiveGateway{
								{
									Default: "Flow_194jx6p",
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Gateway_107rogi",
												},
												Name: "Gateway (Split Flow)",
											},
											Incoming: []string{"Flow_187opjq"},
											Outgoing: []string{"Flow_194jx6p", "Flow_1xhc3bf", "Flow_19m0ydj"},
										},
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Gateway_140ec76",
												},
												Name: "Gateway (Merge Flows)",
											},
											Incoming: []string{"Flow_01ckxme", "Flow_1a24ve1"},
											Outgoing: []string{"Flow_1chzh1q"},
										},
									},
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0lacnkf",
										},
									},
									SourceRef: "StartEvent_1",
									TargetRef: "Activity_0ahdk3x",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_187opjq",
										},
									},
									SourceRef: "Activity_0ahdk3x",
									TargetRef: "Gateway_107rogi",
								},
								{
									FlowElement: element.FlowElement{
										Name: "Default",
										BaseElement: element.BaseElement{
											ID: "Flow_194jx6p",
										},
									},
									SourceRef: "Gateway_107rogi",
									TargetRef: "Activity_172ndxy",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1xhc3bf",
										},
									},
									SourceRef: "Gateway_107rogi",
									TargetRef: "Activity_1lz0l07",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1sd2kpj",
										},
									},
									SourceRef: "Activity_172ndxy",
									TargetRef: "Activity_1lz0l07",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_19m0ydj",
										},
									},
									SourceRef: "Gateway_107rogi",
									TargetRef: "Activity_171qefk",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_11m9fcl",
										},
									},
									SourceRef: "Activity_171qefk",
									TargetRef: "Activity_1lz0l07",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_01ckxme",
										},
										Name: "condition",
									},
									SourceRef: "Activity_171qefk",
									TargetRef: "Gateway_140ec76",
									ConditionExpression: element.ExpressionUnMarshal{
										Type: element.ExpressionTypeFormal,
										ExpressionSubstitution: &element.FormalExpression{
											Expression: element.Expression{
												BaseElementWithMixedContent: element.BaseElementWithMixedContent{},
											},
										},
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1a24ve1",
										},
									},
									SourceRef: "Activity_1lz0l07",
									TargetRef: "Gateway_140ec76",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0zwjy3h",
										},
										Name: "Condition",
									},
									SourceRef: "Activity_172ndxy",
									TargetRef: "Event_1wqqwdz",
									ConditionExpression: element.ExpressionUnMarshal{
										Type: element.ExpressionTypeFormal,
										ExpressionSubstitution: &element.FormalExpression{
											Expression: element.Expression{
												BaseElementWithMixedContent: element.BaseElementWithMixedContent{},
											},
										},
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1chzh1q",
										},
									},
									SourceRef: "Gateway_140ec76",
									TargetRef: "Event_1wqqwdz",
								},
							},
							EndEvents: []element.EndEvent{{
								ThrowEvent: element.ThrowEvent{
									Event: element.Event{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Event_1wqqwdz",
												},
												Name: "End Event",
											},
											Incoming: []string{"Flow_0zwjy3h", "Flow_1chzh1q"},
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	}
	if diff := cmp.Diff(expected, modelelement); diff != "" {
		t.Errorf("Unexpected result (-want, +got):\n%s", diff)
	}
}
