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

func TestA_4_0_export(t *testing.T) {
	path := "./test/A.4.0-export.bpmn"
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
			ID:              "Definitions_14ccggy",
			TargetNamespace: "http://bpmn.io/schema/bpmn",
			Exporter:        "Camunda Modeler",
			ExporterVersion: "5.23.0",
			RootElemnts: element.RootElemnts{
				Collaborations: []element.Collaboration{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "Collaboration_1eeqkua",
							},
						},
						Participants: []element.Participant{
							{
								BaseElement: element.BaseElement{
									ID: "PoolParticipant",
								},
								Name:       "Pool",
								ProcessRef: "Process_0elb8rq",
							},
							{
								BaseElement: element.BaseElement{
									ID: "Participant_1b8727b",
								},
								ProcessRef: "Process_0wqyt7t",
							},
						},
						MessageFlows: []element.MessageFlow{
							{
								BaseElement: element.BaseElement{
									ID: "MessageFlow1MessageFlow",
								},
								Name:      "Message Flow 1",
								SourceRef: "Task1Task",
								TargetRef: "Task3Task",
							},
							{
								BaseElement: element.BaseElement{
									ID: "MessageFlow2MessageFlow",
								},
								Name:      "Message Flow 2",
								SourceRef: "Task5Task",
								TargetRef: "Task2Task",
							},
						},
					},
				},
				Processes: []element.Process{
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "Process_0elb8rq",
								},
							},
						},
						IsExecutable: true,
						FlowElements: element.FlowElements{
							StartEvents: []element.StartEvent{
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "StartEvent1StartEvent",
													},
													Name: "Start Event 1",
												},
												Outgoing: []string{"Flow_193i5eg"},
											},
										},
									},
								},
							},
							Tasks: []element.Task{
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Task1Task",
												},
												Name: "Task 1",
											},
											Incoming: []string{"Flow_193i5eg"},
											Outgoing: []string{"Flow_01btz5o"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Task2Task",
												},
												Name: "Task 2",
											},
											Incoming: []string{"Flow_01btz5o"},
											Outgoing: []string{"Flow_14uttgb"},
										},
									},
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_193i5eg",
										},
									},
									SourceRef: "StartEvent1StartEvent",
									TargetRef: "Task1Task",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_01btz5o",
										},
									},
									SourceRef: "Task1Task",
									TargetRef: "Task2Task",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_14uttgb",
										},
									},
									SourceRef: "Task2Task",
									TargetRef: "EndEvent1EndEvent",
								},
							},
							EndEvents: []element.EndEvent{
								{
									ThrowEvent: element.ThrowEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "EndEvent1EndEvent",
													},
													Name: "End Event 1",
												},
												Incoming: []string{"Flow_14uttgb"},
											},
										},
									},
								},
							},
						},
					},
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "Process_0wqyt7t",
								},
							},
						},
						IsExecutable: false,
						LaneSets: []element.LaneSet{
							{
								BaseElement: element.BaseElement{
									ID: "LaneSet_1lf7xw1",
								},
								Lanes: []element.Lane{
									{
										BaseElement: element.BaseElement{
											ID: "Lane2Lane",
										},
										Name:         "Lane 2",
										FlowNodeRefs: []string{"ExpandedSubProcess2SubProcess", "EndEvent5EndEvent"},
									},
									{
										BaseElement: element.BaseElement{
											ID: "Lane1Lane",
										},
										Name:         "Lane 1",
										FlowNodeRefs: []string{"StartEvent2StartEvent", "Task3Task", "Task5Task", "EndEvent2EndEvent", "ExpandedSubProcess1SubProcess"},
									},
								},
							},
						},
						FlowElements: element.FlowElements{
							StartEvents: []element.StartEvent{
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "StartEvent2StartEvent",
													},
													Name: "Start Event 2",
												},
												Outgoing: []string{"Flow_0g54c4e"},
											},
										},
									},
								},
							},
							Tasks: []element.Task{
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Task3Task",
												},
												Name: "Task 3",
											},
											Incoming: []string{"Flow_0g54c4e"},
											Outgoing: []string{"Flow_0gi2n31", "Flow_14r4yf8"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Task5Task",
												},
												Name: "Task 5",
											},
											Incoming: []string{"Flow_0fjya97"},
											Outgoing: []string{"Flow_0tx0rfn"},
										},
									},
								},
							},
							SubProcesses: []element.SubProcess{
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "ExpandedSubProcess2SubProcess",
												},
												Name: "Expanded Sub-Process 2",
											},
											Incoming: []string{"Flow_0gi2n31"},
											Outgoing: []string{"Flow_04fun2b"},
										},
									},
									FlowElements: element.FlowElements{
										StartEvents: []element.StartEvent{
											{
												CatchEvent: element.CatchEvent{
													Event: element.Event{
														FlowNode: element.FlowNode{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "StartEvent4StartEvent",
																},
																Name: "Start Event 4",
															},
															Outgoing: []string{"Flow_1ups9iv"},
														},
													},
												},
											},
										},
										Tasks: []element.Task{
											{
												Activity: element.Activity{
													FlowNode: element.FlowNode{
														FlowElement: element.FlowElement{
															BaseElement: element.BaseElement{
																ID: "Task6Task",
															},
															Name: "Task 6",
														},
														Incoming: []string{"Flow_1ups9iv"},
														Outgoing: []string{"Flow_0w343e8"},
													},
												},
											},
										},
										SequenceFlows: []element.SequenceFlow{
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_1ups9iv",
													},
												},
												SourceRef: "StartEvent4StartEvent",
												TargetRef: "Task6Task",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_0w343e8",
													},
												},
												SourceRef: "Task6Task",
												TargetRef: "EndEvent4EndEvent",
											},
										},
										EndEvents: []element.EndEvent{
											{
												ThrowEvent: element.ThrowEvent{
													Event: element.Event{
														FlowNode: element.FlowNode{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "EndEvent4EndEvent",
																},
																Name: "End Event 4",
															},
															Incoming: []string{"Flow_0w343e8"},
														},
													},
												},
											},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "ExpandedSubProcess1SubProcess",
												},
												Name: "Expanded Sub-Process 1",
											},
											Incoming: []string{"Flow_14r4yf8"},
											Outgoing: []string{"Flow_0fjya97"},
										},
									},
									FlowElements: element.FlowElements{
										StartEvents: []element.StartEvent{
											{
												CatchEvent: element.CatchEvent{
													Event: element.Event{
														FlowNode: element.FlowNode{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "StartEvent3StartEvent",
																},
																Name: "Start Event 3",
															},
															Outgoing: []string{"Flow_0iuanwa"},
														},
													},
												},
											},
										},
										Tasks: []element.Task{
											{
												Activity: element.Activity{
													FlowNode: element.FlowNode{
														FlowElement: element.FlowElement{
															BaseElement: element.BaseElement{
																ID: "Task4Task",
															},
															Name: "Task 4",
														},
														Incoming: []string{"Flow_0iuanwa"},
														Outgoing: []string{"Flow_1y4ym2b"},
													},
												},
											},
										},
										SequenceFlows: []element.SequenceFlow{
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_0iuanwa",
													},
												},
												SourceRef: "StartEvent3StartEvent",
												TargetRef: "Task4Task",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_1y4ym2b",
													},
												},
												SourceRef: "Task4Task",
												TargetRef: "EndEvent3EndEvent",
											},
										},
										EndEvents: []element.EndEvent{
											{
												ThrowEvent: element.ThrowEvent{
													Event: element.Event{
														FlowNode: element.FlowNode{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "EndEvent3EndEvent",
																},
																Name: "End Event 3",
															},
															Incoming: []string{"Flow_1y4ym2b"},
														},
													},
												},
											},
										},
									},
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_14r4yf8",
										},
									},
									SourceRef: "Task3Task",
									TargetRef: "ExpandedSubProcess1SubProcess",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0gi2n31",
										},
									},
									SourceRef: "Task3Task",
									TargetRef: "ExpandedSubProcess2SubProcess",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_04fun2b",
										},
									},
									SourceRef: "ExpandedSubProcess2SubProcess",
									TargetRef: "EndEvent5EndEvent",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0tx0rfn",
										},
									},
									SourceRef: "Task5Task",
									TargetRef: "EndEvent2EndEvent",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0fjya97",
										},
									},
									SourceRef: "ExpandedSubProcess1SubProcess",
									TargetRef: "Task5Task",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0g54c4e",
										},
									},
									SourceRef: "StartEvent2StartEvent",
									TargetRef: "Task3Task",
								},
							},
							EndEvents: []element.EndEvent{
								{
									ThrowEvent: element.ThrowEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "EndEvent2EndEvent",
													},
													Name: "End Event 2",
												},
												Incoming: []string{"Flow_0tx0rfn"},
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
														ID: "EndEvent5EndEvent",
													},
													Name: "End Event 5",
												},
												Incoming: []string{"Flow_04fun2b"},
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
		t.Errorf("BpmnModelelementFromFile() mismatch (-want +got):\n%s", diff)
	}
}
