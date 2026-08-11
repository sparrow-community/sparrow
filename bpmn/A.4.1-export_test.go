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

func TestA_4_1_export(t *testing.T) {
	path := "./test/A.4.1-export.bpmn"
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
			ID:              "Definitions_18qzd9b",
			TargetNamespace: "http://bpmn.io/schema/bpmn",
			Exporter:        "Camunda Modeler",
			ExporterVersion: "5.23.0",
			RootElemnts: element.RootElemnts{
				Collaborations: []element.Collaboration{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "Collaboration_0weizb5",
							},
						},
						Participants: []element.Participant{
							{
								BaseElement: element.BaseElement{
									ID: "Participant_08dproj",
								},
								Name:       "Pool 1",
								ProcessRef: "Process_0h42ymn",
							},
							{
								BaseElement: element.BaseElement{
									ID: "Participant_1cs40k3",
								},
								Name:       "Pool 1",
								ProcessRef: "Process_18nmg48",
							},
						},
						MessageFlows: []element.MessageFlow{
							{
								BaseElement: element.BaseElement{
									ID: "Flow_0xbmmt1",
								},
								Name:      "Message Flow 2",
								SourceRef: "Activity_0xw7dvh",
								TargetRef: "Activity_1n3qt39",
							},
							{
								BaseElement: element.BaseElement{
									ID: "Flow_0xqa8km",
								},
								Name:      "Message Flow 1",
								SourceRef: "Activity_1n97zgn",
								TargetRef: "Activity_0mmksew",
							},
						},
					},
				},
				Processes: []element.Process{
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "Process_0h42ymn",
								},
							},
						},
						IsExecutable: true,
						LaneSets: []element.LaneSet{
							{
								BaseElement: element.BaseElement{
									ID: "LaneSet_0k0qmv1",
								},
								Lanes: []element.Lane{
									{
										BaseElement: element.BaseElement{
											ID: "Lane_198nahk",
										},
										Name:         "Lane 1",
										FlowNodeRefs: []string{"StartEvent_1", "Activity_1n97zgn", "Activity_1n3qt39", "Event_1a4shvj"},
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
														ID: "StartEvent_1",
													},
													Name: "Start Event 1",
												},
												Outgoing: []string{"Flow_1hlsvck"},
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
													ID: "Activity_1n97zgn",
												},
												Name: "Task 1",
											},
											Incoming: []string{"Flow_1hlsvck"},
											Outgoing: []string{"Flow_1gcj53o"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_1n3qt39",
												},
												Name: "Task 2",
											},
											Incoming: []string{"Flow_1gcj53o"},
											Outgoing: []string{"Flow_1dzj0ew"},
										},
									},
								},
							},
							EndEvents: []element.EndEvent{
								{
									ThrowEvent: element.ThrowEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_1a4shvj",
													},
													Name: "End Event 1",
												},
												Incoming: []string{"Flow_1dzj0ew"},
											},
										},
									},
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1hlsvck",
										},
									},
									SourceRef: "StartEvent_1",
									TargetRef: "Activity_1n97zgn",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1gcj53o",
										},
									},
									SourceRef: "Activity_1n97zgn",
									TargetRef: "Activity_1n3qt39",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1dzj0ew",
										},
									},
									SourceRef: "Activity_1n3qt39",
									TargetRef: "Event_1a4shvj",
								},
							},
						},
					},
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "Process_18nmg48",
								},
							},
						},
						IsExecutable: false,
						LaneSets: []element.LaneSet{
							{
								BaseElement: element.BaseElement{
									ID: "LaneSet_11j5v67",
								},
								Lanes: []element.Lane{
									{
										BaseElement: element.BaseElement{
											ID: "Lane_1pmrgfe",
										},
										Name:         "Lane 4",
										FlowNodeRefs: []string{"Event_099731c", "Activity_0mmksew"},
									},
									{
										BaseElement: element.BaseElement{
											ID: "Lane_0da6uqu",
										},
										Name:         "Lane 3",
										FlowNodeRefs: []string{"Activity_0xw7dvh", "Event_0zje79t", "Event_0wacqbg", "Activity_0hqqmqa", "Activity_1svfico"},
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
														ID: "Event_099731c",
													},
													Name: "Start Event 2",
												},
												Outgoing: []string{"Flow_1dburd6"},
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
													ID: "Activity_0mmksew",
												},
												Name: "Task 3",
											},
											Incoming: []string{"Flow_1dburd6"},
											Outgoing: []string{"Flow_1jziro2", "Flow_0bzve1s"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_0xw7dvh",
												},
												Name: "Task 5",
											},
											Incoming: []string{"Flow_04a6y1c"},
											Outgoing: []string{"Flow_1u6nbqs"},
										},
									},
								},
							},
							EndEvents: []element.EndEvent{
								{
									ThrowEvent: element.ThrowEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_0zje79t",
													},
													Name: "End Event 2",
												},
												Incoming: []string{"Flow_1u6nbqs"},
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
														ID: "Event_0wacqbg",
													},
													Name: "End Event 5",
												},
												Incoming: []string{"Flow_0yduxnw"},
											},
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
													ID: "Activity_0hqqmqa",
												},
												Name: "Expanded Sub-Process 2",
											},
											Incoming: []string{"Flow_1jziro2"},
											Outgoing: []string{"Flow_0yduxnw"},
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
																	ID: "Event_1v26oan",
																},
																Name: "Start Event 4",
															},
															Outgoing: []string{"Flow_079eynm"},
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
																ID: "Activity_0sqf0ju",
															},
															Name: "Task 6",
														},
														Incoming: []string{"Flow_079eynm"},
														Outgoing: []string{"Flow_1palkux"},
													},
												},
											},
										},
										SequenceFlows: []element.SequenceFlow{
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_079eynm",
													},
												},
												SourceRef: "Event_1v26oan",
												TargetRef: "Activity_0sqf0ju",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_1palkux",
													},
												},
												SourceRef: "Activity_0sqf0ju",
												TargetRef: "Event_1r1pop3",
											},
										},
										EndEvents: []element.EndEvent{
											{
												ThrowEvent: element.ThrowEvent{
													Event: element.Event{
														FlowNode: element.FlowNode{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "Event_1r1pop3",
																},
																Name: "End Event 4",
															},
															Incoming: []string{"Flow_1palkux"},
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
													ID: "Activity_1svfico",
												},
												Name: "Expanded Sub-Process 1",
											},
											Incoming: []string{"Flow_0bzve1s"},
											Outgoing: []string{"Flow_04a6y1c"},
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
																	ID: "Event_0xx2ic8",
																},
																Name: "Start Event 3",
															},
															Outgoing: []string{"Flow_0a4cvz3"},
														},
													},
												},
											},
										},
										EndEvents: []element.EndEvent{
											{
												ThrowEvent: element.ThrowEvent{
													Event: element.Event{
														FlowNode: element.FlowNode{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "Event_0gpo0qc",
																},
																Name: "End Event 3",
															},
															Incoming: []string{"Flow_031u0ap"},
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
																ID: "Activity_01cj9um",
															},
															Name: "Task 4",
														},
														Incoming: []string{"Flow_0a4cvz3"},
														Outgoing: []string{"Flow_031u0ap"},
													},
												},
											},
										},
										SequenceFlows: []element.SequenceFlow{
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_0a4cvz3",
													},
												},
												SourceRef: "Event_0xx2ic8",
												TargetRef: "Activity_01cj9um",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_031u0ap",
													},
												},
												SourceRef: "Activity_01cj9um",
												TargetRef: "Event_0gpo0qc",
											},
										},
									},
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1dburd6",
										},
									},
									SourceRef: "Event_099731c",
									TargetRef: "Activity_0mmksew",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1jziro2",
										},
									},
									SourceRef: "Activity_0mmksew",
									TargetRef: "Activity_0hqqmqa",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0bzve1s",
										},
									},
									SourceRef: "Activity_0mmksew",
									TargetRef: "Activity_1svfico",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_04a6y1c",
										},
									},
									SourceRef: "Activity_1svfico",
									TargetRef: "Activity_0xw7dvh",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1u6nbqs",
										},
									},
									SourceRef: "Activity_0xw7dvh",
									TargetRef: "Event_0zje79t",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0yduxnw",
										},
									},
									SourceRef: "Activity_0hqqmqa",
									TargetRef: "Event_0wacqbg",
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
