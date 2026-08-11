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

func TestB_1_0_export(t *testing.T) {
	path := "./test/B.1.0-export.bpmn"
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
			ID:              "Definitions_1w2r0ua",
			TargetNamespace: "http://bpmn.io/schema/bpmn",
			Exporter:        "Camunda Modeler",
			ExporterVersion: "5.23.0",
			RootElemnts: element.RootElemnts{
				Collaborations: []element.Collaboration{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "Collaboration_04c879m",
							},
						},
						Participants: []element.Participant{
							{
								BaseElement: element.BaseElement{
									ID: "Participant_0cgb2hh",
								},
								Name:       "Participant",
								ProcessRef: "Process_1iam7fk",
							},
							{
								BaseElement: element.BaseElement{
									ID: "Participant_1ycx4f0",
								},
								Name:       "Pool",
								ProcessRef: "Process_1ek277i",
							},
						},
						MessageFlows: []element.MessageFlow{
							{
								BaseElement: element.BaseElement{
									ID: "Flow_1t3v3bu",
								},
								Name:      "Message Flow 1",
								SourceRef: "Activity_1l21tbi",
								TargetRef: "Event_0sv3tdm",
							},
							{
								BaseElement: element.BaseElement{
									ID: "Flow_0f3qdjk",
								},
								Name:      "Message Flow 2",
								SourceRef: "Event_1qd5ifd",
								TargetRef: "Activity_0yyznkb",
							},
						},
						Artifacts: element.Artifacts{
							Groups: []element.Group{
								{
									Artifact: element.Artifact{
										BaseElement: element.BaseElement{
											ID: "Group_1dpx6wg",
										},
									},
									CategoryValueRef: "CategoryValue_1srow05",
								},
							},
							TextAnnotations: []element.TextAnnotation{
								{
									Artifact: element.Artifact{
										BaseElement: element.BaseElement{
											ID: "TextAnnotation_13xc6qx",
										},
									},
									Text: "Text Annotation",
								},
							},
							Associations: []element.Association{
								{
									Artifact: element.Artifact{
										BaseElement: element.BaseElement{
											ID: "Association_0eu7lou",
										},
									},
									SourceRef: "Activity_00ifb1p",
									TargetRef: "TextAnnotation_13xc6qx",
								},
							},
						},
					},
				},
				Processes: []element.Process{
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "Process_1iam7fk",
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
														ID: "StartEvent_1",
													},
													Name: "Start Timer Event",
												},
												Outgoing: []string{"Flow_0qh645u"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											TimerEventDefinitions: []element.TimerEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "TimerEventDefinition_1v56mia",
															},
														},
													},
												},
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
													ID: "Activity_1l21tbi",
												},
												Name: "Abstract Task 1",
											},
											Incoming: []string{"Flow_0qh645u"},
											Outgoing: []string{"Flow_1x6yjls"},
										},
									},
								},
							},
							UserTasks: []element.UserTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_1923g7r",
													},
													Name: "User Task 2",
												},
												Incoming: []string{"Flow_1x6yjls"},
												Outgoing: []string{"Flow_1e8jnkr"},
											},
										},
									},
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0qh645u",
										},
									},
									SourceRef: "StartEvent_1",
									TargetRef: "Activity_1l21tbi",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1x6yjls",
										},
									},
									SourceRef: "Activity_1l21tbi",
									TargetRef: "Activity_1923g7r",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1e8jnkr",
										},
									},
									SourceRef: "Activity_1923g7r",
									TargetRef: "Activity_0yyznkb",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0hivzfq",
										},
									},
									SourceRef: "Activity_0yyznkb",
									TargetRef: "Event_0eo53c5",
								},
							},
							ServiceTasks: []element.ServiceTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_0yyznkb",
													},
													Name: "Service Task 3",
												},
												Incoming: []string{"Flow_1e8jnkr"},
												Outgoing: []string{"Flow_0hivzfq"},
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
														ID: "Event_0eo53c5",
													},
													Name: "End Event None 1",
												},
												Incoming: []string{"Flow_0hivzfq"},
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
									ID: "Process_1ek277i",
								},
							},
						},
						IsExecutable: false,
						LaneSets: []element.LaneSet{
							{
								BaseElement: element.BaseElement{
									ID: "LaneSet_0c1we21",
								},
								Lanes: []element.Lane{
									{
										BaseElement: element.BaseElement{
											ID: "Lane_0ii62sr",
										},
										Name: "Lane 2",
										FlowNodeRefs: []string{
											"Activity_15s9oor",
											"Gateway_1x7jv2q",
											"Activity_0349q9v",
											"Activity_0puge1w",
											"Activity_0kxionv",
											"Gateway_1pwtt6e",
											"Event_1c6ino0",
										},
									},
									{
										BaseElement: element.BaseElement{
											ID: "Lane_1tzdy12",
										},
										Name: "Lane 1",
										FlowNodeRefs: []string{
											"Gateway_0u1bqqo",
											"Gateway_1tvwnx5",
											"Activity_00ifb1p",
											"Gateway_02v3rbb",
											"Activity_0wox4hb",
											"Activity_0d5gl6d",
											"Event_1qd5ifd",
											"Event_0sv3tdm",
										},
									},
								},
							},
						},
						FlowElements: element.FlowElements{
							ParallelGatewaies: []element.ParallelGateway{
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Gateway_0u1bqqo",
												},
												Name: "Parallel Gateway Divergence",
											},
											Incoming: []string{"Flow_0fua85o"},
											Outgoing: []string{"Flow_1ra0g8d", "Flow_1ao20of"},
										},
									},
								},
							},
							ExclusiveGatewaies: []element.ExclusiveGateway{
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Gateway_1tvwnx5",
												},
												Name: "Exclusive Gateway Divergence",
											},
											Incoming: []string{"Flow_1ra0g8d"},
											Outgoing: []string{"Flow_0defg8c", "Flow_0bakxd3"},
										},
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Gateway_1x7jv2q",
												},
												Name: "Exclusive Gateway Divergence 2",
											},
											Incoming: []string{"Flow_00q21zs"},
											Outgoing: []string{"Flow_197wmo5", "Flow_0btclp9"},
										},
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Gateway_02v3rbb",
												},
												Name: "Exclusive Gateway Convergence 1",
											},
											Incoming: []string{"Flow_0t2kr8a", "Flow_16e2ast"},
											Outgoing: []string{"Flow_18eyfje"},
										},
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Gateway_1pwtt6e",
												},
												Name: "Exclusive Gateway Convergence 2",
											},
											Incoming: []string{"Flow_1paju09", "Flow_1eeuiyj"},
											Outgoing: []string{"Flow_0ats6eb"},
										},
									},
								},
							},
							CallActivities: []element.CallActivity{
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_00ifb1p",
												},
												Name: "Call Activity Collapsed",
											},
											Incoming: []string{"Flow_0defg8c"},
											Outgoing: []string{"Flow_00rxphe"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_0wox4hb",
												},
												Name: "Call Activity Calling a Global Task",
											},
											Incoming: []string{"Flow_0bakxd3"},
											Outgoing: []string{"Flow_16e2ast"},
										},
									},
								},
							},
							UserTasks: []element.UserTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_15s9oor",
													},
													Name: "User Task 5",
												},
												Incoming: []string{"Flow_1ao20of"},
												Outgoing: []string{"Flow_00q21zs"},
											},
										},
									},
								},
							},
							ServiceTasks: []element.ServiceTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_0349q9v",
													},
													Name: "Service Task 7",
												},
												Incoming: []string{"Flow_0btclp9"},
												Outgoing: []string{"Flow_1eeuiyj"},
											},
											Properties: []element.Property{
												{
													BaseElement: element.BaseElement{
														ID: "Property_1f28g0r",
													},
													Name: "__targetRef_placeholder",
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "DataInputAssociation_0eftry3",
														},
														SourceRef: "DataObjectReference_1w6co9d",
														TargetRef: "Property_1f28g0r",
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "DataOutputAssociation_0o650y5",
														},
														TargetRef: "DataStoreReference_0qlcobt",
													},
												},
											},
										},
									},
								},
							},
							DataStoreReferences: []element.DataStoreReference{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataStoreReference_0qlcobt",
										},
										Name: "Data Store Reference",
									},
								},
							},
							SubProcesses: []element.SubProcess{
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_0puge1w",
												},
												Name: "Collapsed Sub-Process",
											},
											Incoming: []string{"Flow_197wmo5"},
											Outgoing: []string{"Flow_0sx30ud"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_0d5gl6d",
												},
												Name: "Call Activity - Expanded",
											},
											Incoming: []string{"Flow_00rxphe"},
											Outgoing: []string{"Flow_0t2kr8a"},
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
																	ID: "Event_1anrpw6",
																},
																Name: "Start Event None 1",
															},
															Outgoing: []string{"Flow_1inrdrc"},
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
																ID: "Activity_15sig42",
															},
															Name: "Abstract Task 4",
														},
														Incoming: []string{"Flow_1inrdrc"},
														Outgoing: []string{"Flow_10ww89o"},
													},
												},
											},
										},
										SequenceFlows: []element.SequenceFlow{
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_1inrdrc",
													},
												},
												SourceRef: "Event_1anrpw6",
												TargetRef: "Activity_15sig42",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_10ww89o",
													},
												},
												SourceRef: "Activity_15sig42",
												TargetRef: "Event_0sh5bsa",
											},
										},
										EndEvents: []element.EndEvent{
											{
												ThrowEvent: element.ThrowEvent{
													Event: element.Event{
														FlowNode: element.FlowNode{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "Event_0sh5bsa",
																},
																Name: "End Event None 2",
															},
															Incoming: []string{"Flow_10ww89o"},
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
													ID: "Activity_0kxionv",
												},
												Name: "Sub Process - Expanded",
											},
											Incoming: []string{"Flow_0sx30ud"},
											Outgoing: []string{"Flow_1paju09"},
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
																	ID: "Event_11sri5l",
																},
																Name: "Start Event None 2",
															},
															Outgoing: []string{"Flow_07pghkc"},
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
																ID: "Activity_0brlj8k",
															},
															Name: "Abstract Task 6",
														},
														Incoming: []string{"Flow_07pghkc"},
														Outgoing: []string{"Flow_11pggpx"},
													},
												},
											},
										},
										SequenceFlows: []element.SequenceFlow{
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_07pghkc",
													},
												},
												SourceRef: "Event_11sri5l",
												TargetRef: "Activity_0brlj8k",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_11pggpx",
													},
												},
												SourceRef: "Activity_0brlj8k",
												TargetRef: "Event_0axyohg",
											},
										},
										EndEvents: []element.EndEvent{
											{
												ThrowEvent: element.ThrowEvent{
													Event: element.Event{
														FlowNode: element.FlowNode{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "Event_0axyohg",
																},
																Name: "End Event None 3",
															},
															Incoming: []string{"Flow_11pggpx"},
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
											ID: "Flow_16e2ast",
										},
									},
									SourceRef: "Activity_0wox4hb",
									TargetRef: "Gateway_02v3rbb",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_18eyfje",
										},
									},
									SourceRef: "Gateway_02v3rbb",
									TargetRef: "Event_1qd5ifd",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0t2kr8a",
										},
									},
									SourceRef: "Activity_0d5gl6d",
									TargetRef: "Gateway_02v3rbb",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_00rxphe",
										},
									},
									SourceRef: "Activity_00ifb1p",
									TargetRef: "Activity_0d5gl6d",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0ats6eb",
										},
									},
									SourceRef: "Gateway_1pwtt6e",
									TargetRef: "Event_1c6ino0",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1eeuiyj",
										},
									},
									SourceRef: "Activity_0349q9v",
									TargetRef: "Gateway_1pwtt6e",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1paju09",
										},
									},
									SourceRef: "Activity_0kxionv",
									TargetRef: "Gateway_1pwtt6e",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0sx30ud",
										},
									},
									SourceRef: "Activity_0puge1w",
									TargetRef: "Activity_0kxionv",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0btclp9",
										},
									},
									SourceRef: "Gateway_1x7jv2q",
									TargetRef: "Activity_0349q9v",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_197wmo5",
										},
									},
									SourceRef: "Gateway_1x7jv2q",
									TargetRef: "Activity_0puge1w",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_00q21zs",
										},
									},
									SourceRef: "Activity_15s9oor",
									TargetRef: "Gateway_1x7jv2q",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1ao20of",
										},
									},
									SourceRef: "Gateway_0u1bqqo",
									TargetRef: "Activity_15s9oor",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0bakxd3",
										},
									},
									SourceRef: "Gateway_1tvwnx5",
									TargetRef: "Activity_0wox4hb",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0defg8c",
										},
									},
									SourceRef: "Gateway_1tvwnx5",
									TargetRef: "Activity_00ifb1p",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1ra0g8d",
										},
									},
									SourceRef: "Gateway_0u1bqqo",
									TargetRef: "Gateway_1tvwnx5",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0fua85o",
										},
									},
									SourceRef: "Event_0sv3tdm",
									TargetRef: "Gateway_0u1bqqo",
								},
							},
							DataObjectReferenes: []element.DataObjectReference{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObjectReference_1w6co9d",
										},
										Name: "Data Object",
									},
									DataObjectRef: "DataObject_1kcb5rw",
								},
							},
							DataObjects: []element.DataObject{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObject_1kcb5rw",
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
														ID: "Event_1qd5ifd",
													},
													Name: "End Event Message",
												},
												Incoming: []string{"Flow_18eyfje"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "MessageEventDefinition_1jqojp3",
															},
														},
													},
												},
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
														ID: "Event_1c6ino0",
													},
													Name: "End Event Terminate",
												},
												Incoming: []string{"Flow_0ats6eb"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											TerminateEventDefinitions: []element.TerminateEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "TerminateEventDefinition_1p6bo3l",
															},
														},
													},
												},
											},
										},
									},
								},
							},
							StartEvents: []element.StartEvent{
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_0sv3tdm",
													},
													Name: "Start Event Message",
												},
												Outgoing: []string{"Flow_0fua85o"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "MessageEventDefinition_0pp34ra",
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
					},
				},
				Categories: []element.Category{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "Category_08lc9k6",
							},
						},
						CategoryValues: []element.CategoryValue{
							{
								BaseElement: element.BaseElement{
									ID: "CategoryValue_1srow05",
								},
								Value: "Group",
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
