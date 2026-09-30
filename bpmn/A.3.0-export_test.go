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

func TestA_3_0_export(t *testing.T) {
	path := "./test/A.3.0-export.bpmn"
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
			ID:              "Definitions_06decf0",
			TargetNamespace: "http://bpmn.io/schema/bpmn",
			Exporter:        "Camunda Modeler",
			ExporterVersion: "5.23.0",
			RootElemnts: element.RootElemnts{
				Processes: []element.Process{
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "Process_1qh1mjw",
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
											Outgoing: []string{"Flow_0su0msi"},
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
													ID: "Activity_164hjmz",
												},
												Name: "Task 1",
											},
											Incoming: []string{"Flow_0su0msi"},
											Outgoing: []string{"Flow_0iagd7s"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_0mu4iin",
												},
												Name: "Task 2",
											},
											Incoming: []string{"Flow_0dukxt7"},
											Outgoing: []string{"Flow_0fg1rbu"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_04ib995",
												},
												Name: "Task 4",
											},
											Incoming: []string{"Flow_0hvhb63"},
											Outgoing: []string{"Flow_12zl0ro"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_12jlku6",
												},
												Name: "Task 3",
											},
											Incoming: []string{"Flow_0g8ko8k"},
											Outgoing: []string{"Flow_02rz41z"},
										},
									},
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0su0msi",
										},
									},
									SourceRef: "StartEvent_1",
									TargetRef: "Activity_164hjmz",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0iagd7s",
										},
									},
									SourceRef: "Activity_164hjmz",
									TargetRef: "Activity_1j4b29j",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0dukxt7",
										},
									},
									SourceRef: "Activity_1j4b29j",
									TargetRef: "Activity_0mu4iin",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0hvhb63",
										},
									},
									SourceRef: "Event_1bgdnfg",
									TargetRef: "Activity_04ib995",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_12zl0ro",
										},
									},
									SourceRef: "Activity_04ib995",
									TargetRef: "Event_1179od7",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0fg1rbu",
										},
									},
									SourceRef: "Activity_0mu4iin",
									TargetRef: "Event_1jn5oqm",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0g8ko8k",
										},
									},
									SourceRef: "Event_1uez1gc",
									TargetRef: "Activity_12jlku6",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_02rz41z",
										},
									},
									SourceRef: "Activity_12jlku6",
									TargetRef: "Event_1jn5oqm",
								},
							},
							SubProcesses: []element.SubProcess{
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_1j4b29j",
												},
												Name: "Collapsed Sub-Process",
											},
											Incoming: []string{"Flow_0iagd7s"},
											Outgoing: []string{"Flow_0dukxt7"},
										},
									},
								},
							},
							BoundaryEvents: []element.BoundaryEvent{
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_1uez1gc",
													},
													Name: "Boundary Intermediate Event Non-Interrupting Message",
												},
												Outgoing: []string{"Flow_0g8ko8k"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "MessageEventDefinition_1eelr04",
															},
														},
													},
												},
											},
										},
									},
									CancelActivity: false,
									AttachedToRef:  "Activity_1j4b29j",
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_1bgdnfg",
													},
													Name: "Boundary Intermediate Event Interrupting Escalation",
												},
												Outgoing: []string{"Flow_0hvhb63"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											EscalationEventDefinitions: []element.EscalationEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "EscalationEventDefinition_1u8izyv",
															},
														},
													},
												},
											},
										},
									},
									CancelActivity: true,
									AttachedToRef: "Activity_1j4b29j",
								},
							},
							EndEvents: []element.EndEvent{
								{
									ThrowEvent: element.ThrowEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_1179od7",
													},
													Name: "End Event 2",
												},
												Incoming: []string{"Flow_12zl0ro"},
											},
										},
									},
								},
								{
									ThrowEvent: element.ThrowEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_1jn5oqm",
													},
													Name: "End Event 1",
												},
												Incoming: []string{"Flow_0fg1rbu", "Flow_02rz41z"},
											},
										},
									},
								},
							},
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
