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

func TestC_4_0_export(t *testing.T) {
	// create test use ./test/C.4.0-export.bpmn
	path := "./test/C.4.0-export.bpmn"
	modelelement, err := BpmnModelelementFromFile(path)
	if err != nil {
		t.Fatalf("BpmnModelelementFromFile error %s: %s", path, err)
	}
	if modelelement == nil {
		t.Fatalf("modelelement is nil, %s", path)
	}

	excepted := &Bpmn{
		Definitions: &element.Definitions{
			XMLName: xml.Name{
				Space: "http://www.omg.org/spec/BPMN/20100524/MODEL",
				Local: "definitions",
			},
			ID:              "Definitions_008og7z",
			TargetNamespace: "http://bpmn.io/schema/bpmn",
			Exporter:        "Camunda Modeler",
			ExporterVersion: "5.23.0",
			RootElemnts: element.RootElemnts{
				Collaborations: []element.Collaboration{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "Collaboration_08tzwy8",
							},
						},
						Participants: []element.Participant{
							{
								BaseElement: element.BaseElement{
									ID: "Participant_0d33sio",
								},
								Name:       "Money Bank",
								ProcessRef: "Process_07wr932",
							},
						},
					},
				},
				Processes: []element.Process{
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "Process_07wr932",
								},
							},
						},
						IsExecutable: false,
						LaneSets: []element.LaneSet{
							{
								BaseElement: element.BaseElement{
									ID: "LaneSet_0ujafgq",
								},
								Lanes: []element.Lane{
									{
										BaseElement: element.BaseElement{
											ID: "Lane_0w7h48m",
										},
										Name: "HR Department",
										FlowNodeRefs: []string{
											"Event_0aqeoee",
											"Activity_1vif48i",
											"Gateway_1nitinv",
											"Activity_0gis8ag",
											"Activity_1iegcae",
											"Gateway_13pi9cy",
											"Activity_1fydq08",
											"Activity_1ptc2l6",
											"Activity_0fxxj1x",
											"Gateway_0tmub1o",
											"Activity_0pb64l8",
										},
										ChildLaneSet: element.LaneSet{
											BaseElement: element.BaseElement{
												ID: "LaneSet_0gxwbwu",
											},
										},
									},
									{
										BaseElement: element.BaseElement{
											ID: "Lane_1x43fzg",
										},
										Name: "Responsible Department",
										FlowNodeRefs: []string{
											"Activity_1567dil",
											"Event_1uj58zc",
											"Activity_1uixx0o",
											"Activity_05503d3",
											"Event_0qacxgu",
											"Event_0axou2i",
											"Gateway_1s7by1w",
											"Gateway_0vtzq53",
											"Event_0bjpe63",
											"Activity_0nzk36a",
											"Activity_0yq2h1t",
											"Event_1dg5pp3",
										},
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
														ID: "Event_0aqeoee",
													},
													Name: "Candidate accepted offer",
												},
												Outgoing: []string{"Flow_161isb1"},
											},
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
														ID: "Activity_1vif48i",
													},
													Name: "Send candidate Contract",
												},
												Incoming: []string{"Flow_161isb1", "Flow_0d72bsp"},
												Outgoing: []string{"Flow_1kkgeei"},
											},
										},
									},
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_0gis8ag",
													},
													Name: "Review terms of contract",
												},
												Incoming: []string{"Flow_1k91huu"},
												Outgoing: []string{"Flow_0d72bsp"},
											},
										},
									},
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_1iegcae",
													},
													Name: "Get signature on contract and notify responsible department",
												},
												Incoming: []string{"Flow_02hcma9"},
												Outgoing: []string{"Flow_1xymj9u"},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "DataOutputAssociation_0q6chhu",
														},
														TargetRef: "DataStoreReference_1el2zpd",
													},
												},
											},
										},
									},
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_1fydq08",
													},
													Name: "Inform employee of company policies",
												},
												Incoming: []string{"Flow_0jlumov"},
												Outgoing: []string{"Flow_0gs48ue"},
											},
										},
									},
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_1ptc2l6",
													},
													Name: "Introduce employee to company Mission, Vision and Values",
												},
												Incoming: []string{"Flow_0gs48ue"},
												Outgoing: []string{"Flow_0f14mpt"},
											},
										},
									},
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_1567dil",
													},
													Name: "Request preparations for a new employee",
												},
												Incoming: []string{"Flow_0u0mecw"},
												Outgoing: []string{"Flow_0qufezv"},
											},
										},
									},
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_0fxxj1x",
													},
													Name: "Register for medical insurance",
												},
												Incoming: []string{"Flow_13ztu1l"},
												Outgoing: []string{"Flow_0jag15j"},
											},
										},
									},
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_1uixx0o",
													},
													Name: "Introduce new employee to the team",
												},
												Incoming: []string{"Flow_0pzvylb"},
												Outgoing: []string{"Flow_1fuxj6b"},
											},
										},
									},
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_05503d3",
													},
													Name: "Perform training for position",
												},
												Incoming: []string{"Flow_1fuxj6b"},
												Outgoing: []string{"Flow_1ikbbnd"},
											},
										},
									},
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_0nzk36a",
													},
													Name: "Compile welcome package",
												},
												Incoming: []string{"Flow_1367fco"},
												Outgoing: []string{"Flow_0n1kpzs"},
											},
										},
									},
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_0pb64l8",
													},
													Name: "Perform training for time reports sick leave and holidays",
												},
												Incoming: []string{"Flow_0f14mpt"},
												Outgoing: []string{"Flow_13ztu1l"},
											},
										},
									},
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_0yq2h1t",
													},
													Name: "Give employee welcome package",
												},
												Incoming: []string{"Flow_0n1kpzs"},
												Outgoing: []string{"Flow_1acnq20"},
											},
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
													ID: "Gateway_1nitinv",
												},
												Name: "Contract terms accepted ?",
											},
											Incoming: []string{"Flow_1kkgeei"},
											Outgoing: []string{"Flow_02hcma9", "Flow_1k91huu"},
										},
									},
								},
							},
							ParallelGatewaies: []element.ParallelGateway{
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Gateway_13pi9cy",
												},
												Name: "Non exclusive Gateway",
											},
											Incoming: []string{"Flow_1xymj9u"},
											Outgoing: []string{"Flow_0jlumov", "Flow_0u0mecw"},
										},
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Gateway_0tmub1o",
												},
												Name: "Non exclusive Gateway",
											},
											Incoming: []string{"Flow_0jag15j", "Flow_1jp11xh"},
											Outgoing: []string{"Flow_0pzvylb"},
										},
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Gateway_1s7by1w",
												},
											},
											Incoming: []string{"Flow_01sh6l9", "Flow_0utj2cp", "Flow_1dqygu2"},
											Outgoing: []string{"Flow_1367fco"},
										},
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Gateway_0vtzq53",
												},
												Name: "Non exclusive Gateway",
											},
											Incoming: []string{"Flow_1ikbbnd"},
											Outgoing: []string{"Flow_1ysfo4f", "Flow_07uts0i", "Flow_11nkk44"},
										},
									},
								},
							},
							DataStoreReferences: []element.DataStoreReference{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataStoreReference_1el2zpd",
										},
										Name: "Employee Details",
									},
								},
							},
							IntermediateThrowEvents: []element.IntermediateThrowEvent{
								{
									ThrowEvent: element.ThrowEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_1uj58zc",
													},
													Name: "New employee in department X",
												},
												Incoming: []string{"Flow_0qufezv"},
												Outgoing: []string{"Flow_1jp11xh"},
											},
											Properties: []element.Property{
												{
													BaseElement: element.BaseElement{
														ID: "Property_0bxqyv2",
													},
													Name: "__targetRef_placeholder",
												},
											},
										},
										DataInputAssociations: []element.DataInputAssociation{
											{
												DataAssociation: element.DataAssociation{
													BaseElement: element.BaseElement{
														ID: "DataInputAssociation_0z6q52m",
													},
													SourceRef: "DataStoreReference_1el2zpd",
													TargetRef: "Property_0bxqyv2",
												},
											},
										},
										EventDefinitions: element.EventDefinitions{
											SignalEventDefinitions: []element.SignalEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "SignalEventDefinition_031ap9q",
															},
														},
													},
												},
											},
										},
									},
								},
							},
							IntermediateCatchEvents: []element.IntermediateCatchEvent{
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_0qacxgu",
													},
													Name: "Input from IT ready",
												},
												Incoming: []string{"Flow_1ysfo4f"},
												Outgoing: []string{"Flow_0utj2cp"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "MessageEventDefinition_0mcgs4q",
															},
														},
													},
												},
											},
										},
									},
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_0axou2i",
													},
													Name: "Input from Payroll ready",
												},
												Incoming: []string{"Flow_07uts0i"},
												Outgoing: []string{"Flow_01sh6l9"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "MessageEventDefinition_11piq11",
															},
														},
													},
												},
											},
										},
									},
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_0bjpe63",
													},
													Name: "Input from Facilities ready",
												},
												Incoming: []string{"Flow_11nkk44"},
												Outgoing: []string{"Flow_1dqygu2"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "MessageEventDefinition_0txndb3",
															},
														},
													},
												},
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
														ID: "Event_1dg5pp3",
													},
													Name: "End Event",
												},
												Incoming: []string{"Flow_1acnq20"},
											},
										},
									},
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_161isb1",
										},
									},
									SourceRef: "Event_0aqeoee",
									TargetRef: "Activity_1vif48i",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1kkgeei",
										},
									},
									SourceRef: "Activity_1vif48i",
									TargetRef: "Gateway_1nitinv",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1k91huu",
										},
										Name: "No",
									},
									SourceRef: "Gateway_1nitinv",
									TargetRef: "Activity_0gis8ag",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0d72bsp",
										},
										Name: "",
									},
									SourceRef: "Activity_0gis8ag",
									TargetRef: "Activity_1vif48i",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_02hcma9",
										},
										Name: "Yes",
									},
									SourceRef: "Gateway_1nitinv",
									TargetRef: "Activity_1iegcae",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1xymj9u",
										},
									},
									SourceRef: "Activity_1iegcae",
									TargetRef: "Gateway_13pi9cy",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0jlumov",
										},
									},
									SourceRef: "Gateway_13pi9cy",
									TargetRef: "Activity_1fydq08",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0gs48ue",
										},
									},
									SourceRef: "Activity_1fydq08",
									TargetRef: "Activity_1ptc2l6",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0f14mpt",
										},
									},
									SourceRef: "Activity_1ptc2l6",
									TargetRef: "Activity_0pb64l8",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0jag15j",
										},
									},
									SourceRef: "Activity_0fxxj1x",
									TargetRef: "Gateway_0tmub1o",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0u0mecw",
										},
									},
									SourceRef: "Gateway_13pi9cy",
									TargetRef: "Activity_1567dil",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0qufezv",
										},
									},
									SourceRef: "Activity_1567dil",
									TargetRef: "Event_1uj58zc",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1jp11xh",
										},
									},
									SourceRef: "Event_1uj58zc",
									TargetRef: "Gateway_0tmub1o",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0pzvylb",
										},
									},
									SourceRef: "Gateway_0tmub1o",
									TargetRef: "Activity_1uixx0o",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1fuxj6b",
										},
									},
									SourceRef: "Activity_1uixx0o",
									TargetRef: "Activity_05503d3",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1ikbbnd",
										},
									},
									SourceRef: "Activity_05503d3",
									TargetRef: "Gateway_0vtzq53",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1ysfo4f",
										},
									},
									SourceRef: "Gateway_0vtzq53",
									TargetRef: "Event_0qacxgu",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_07uts0i",
										},
									},
									SourceRef: "Gateway_0vtzq53",
									TargetRef: "Event_0axou2i",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_11nkk44",
										},
									},
									SourceRef: "Gateway_0vtzq53",
									TargetRef: "Event_0bjpe63",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_01sh6l9",
										},
									},
									SourceRef: "Event_0axou2i",
									TargetRef: "Gateway_1s7by1w",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0utj2cp",
										},
									},
									SourceRef: "Event_0qacxgu",
									TargetRef: "Gateway_1s7by1w",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1dqygu2",
										},
									},
									SourceRef: "Event_0bjpe63",
									TargetRef: "Gateway_1s7by1w",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1367fco",
										},
									},
									SourceRef: "Gateway_1s7by1w",
									TargetRef: "Activity_0nzk36a",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0n1kpzs",
										},
									},
									SourceRef: "Activity_0nzk36a",
									TargetRef: "Activity_0yq2h1t",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1acnq20",
										},
									},
									SourceRef: "Activity_0yq2h1t",
									TargetRef: "Event_1dg5pp3",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_13ztu1l",
										},
									},
									SourceRef: "Activity_0pb64l8",
									TargetRef: "Activity_0fxxj1x",
								},
							},
						},
					},
				},
			},
		},
	}

	if diff := cmp.Diff(excepted, modelelement); diff != "" {
		t.Errorf("Mismatch in %s (-want +got):\n%s", path, diff)
	}
}
