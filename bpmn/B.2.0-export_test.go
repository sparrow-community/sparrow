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

func TestB_2_0_export(t *testing.T) {
	path := "./test/B.2.0-export.bpmn"
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
			ID:              "Definitions_0irsq5b",
			TargetNamespace: "http://bpmn.io/schema/bpmn",
			Exporter:        "Camunda Modeler",
			ExporterVersion: "5.23.0",
			RootElemnts: element.RootElemnts{
				Collaborations: []element.Collaboration{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "Collaboration_1yp4w51",
							},
						},
						Participants: []element.Participant{
							{
								BaseElement: element.BaseElement{
									ID: "Participant_0d1io6b",
								},
								Name:       "Participant",
								ProcessRef: "Process_0nca5ry",
							},
							{
								BaseElement: element.BaseElement{
									ID: "Participant_0s1hgls",
								},
								Name:       "Pool",
								ProcessRef: "Process_1xz7va4",
							},
						},
						MessageFlows: []element.MessageFlow{
							{
								BaseElement: element.BaseElement{
									ID: "Flow_1gv8gte",
								},
								Name:      "Message Flow 1",
								SourceRef: "Activity_0pn331p",
								TargetRef: "Event_0yxib88",
							},
							{
								BaseElement: element.BaseElement{
									ID: "Flow_12bkpzd",
								},
								Name:      "Message Flow 2",
								SourceRef: "Event_0dwi53b",
								TargetRef: "Activity_0javi3i",
							},
						},
						Artifacts: element.Artifacts{
							Groups: []element.Group{
								{
									Artifact: element.Artifact{
										BaseElement: element.BaseElement{
											ID: "Group_1mitxsy",
										},
									},
									CategoryValueRef: "CategoryValue_02um9y9",
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
									ID: "Process_0nca5ry",
								},
							},
						},
						IsExecutable: true,
						FlowElements: element.FlowElements{
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0gzdgm8",
										},
									},
									SourceRef: "StartEvent_1",
									TargetRef: "Activity_1uu0yo4",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0xky4tu",
										},
									},
									SourceRef: "Activity_1uu0yo4",
									TargetRef: "Activity_0pn331p",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1qqutfv",
										},
									},
									SourceRef: "Activity_0pn331p",
									TargetRef: "Activity_07vdczu",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0vyan0i",
										},
									},
									SourceRef: "Activity_07vdczu",
									TargetRef: "Gateway_0iz1sti",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1fmiuuz",
										},
										Name: "Conditional Sequence Flow",
									},
									SourceRef: "Gateway_0iz1sti",
									TargetRef: "Activity_00kk0w1",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0uffa0a",
										},
									},
									SourceRef: "Activity_00kk0w1",
									TargetRef: "Event_128e9tk",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_10j3phf",
										},
									},
									SourceRef: "Event_128e9tk",
									TargetRef: "Activity_03q0xwc",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_02vtnpn",
										},
									},
									SourceRef: "Activity_03q0xwc",
									TargetRef: "Activity_0thkytf",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0t9173h",
										},
									},
									SourceRef: "Activity_0thkytf",
									TargetRef: "Gateway_0y5g78s",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1exfm6h",
										},
									},
									SourceRef: "Event_1rsgo7a",
									TargetRef: "Activity_0ayd8ln",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_19165ws",
										},
									},
									SourceRef: "Gateway_0y5g78s",
									TargetRef: "Event_0403g4k",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1lz1d6w",
										},
										Name: "Default Sequence Flow 1",
									},
									SourceRef: "Gateway_0iz1sti",
									TargetRef: "Activity_0qnc8vy",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1iyrc7z",
										},
									},
									SourceRef: "Activity_0qnc8vy",
									TargetRef: "Activity_0ffbdg4",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1md83y0",
										},
									},
									SourceRef: "Activity_0ayd8ln",
									TargetRef: "Gateway_0y5g78s",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1b1ivbn",
										},
									},
									SourceRef: "Activity_0ffbdg4",
									TargetRef: "Activity_0s6eb2u",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0i7jvk9",
										},
									},
									SourceRef: "Event_1f35b4w",
									TargetRef: "Activity_0wdnxu3",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0rs1jsj",
										},
									},
									SourceRef: "Activity_0wdnxu3",
									TargetRef: "Event_0pi17ux",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1luuks1",
										},
									},
									SourceRef: "Event_0pi17ux",
									TargetRef: "Activity_0javi3i",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1vs1xz8",
										},
									},
									SourceRef: "Activity_0javi3i",
									TargetRef: "Event_1cb9pew",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0uy5u27",
										},
									},
									SourceRef: "Activity_0s6eb2u",
									TargetRef: "Gateway_0y5g78s",
								},
							},
							SubProcesses: []element.SubProcess{
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_0ffbdg4",
												},
												Name: "Expanded Sub-Process 1",
											},
											Incoming: []string{"Flow_1iyrc7z"},
											Outgoing: []string{"Flow_1b1ivbn"},
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
																	ID: "Event_1umydnx",
																},
																Name: "Start Event 2",
															},
															Outgoing: []string{"Flow_0pqufnb"},
														},
													},
												},
											},
										},
										SequenceFlows: []element.SequenceFlow{
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_0pqufnb",
													},
												},
												SourceRef: "Event_1umydnx",
												TargetRef: "Activity_1n1hhyt",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_1qzdagc",
													},
												},
												SourceRef: "Activity_1n1hhyt",
												TargetRef: "Event_1bycw4y",
											},
										},
										UserTasks: []element.UserTask{
											{
												Task: element.Task{
													Activity: element.Activity{
														FlowNode: element.FlowNode{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "Activity_1n1hhyt",
																},
																Name: "User Task 7 Standard Loop",
															},
															Incoming: []string{"Flow_0pqufnb"},
															Outgoing: []string{"Flow_1qzdagc"},
														},
														LoopCharacteristicsElements: element.LoopCharacteristicsElements{
															StandardLoopCharacteristics: []element.StandardLoopCharacteristics{{}},
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
																	ID: "Event_1bycw4y",
																},
																Name: "End Event 2",
															},
															Incoming: []string{"Flow_1qzdagc"},
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
													ID: "Activity_03q0xwc",
												},
												Name: "Collapsed Sub-Process 1 Multi-instances",
											},
											Incoming: []string{"Flow_10j3phf"},
											Outgoing: []string{"Flow_02vtnpn"},
										},
										LoopCharacteristicsElements: element.LoopCharacteristicsElements{
											MultielementLoopCharacteristics: []element.MultiInstanceLoopCharacteristics{{}},
										},
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
														ID: "Event_128e9tk",
													},
													Name: "Intermediate Event Signal Throw 1",
												},
												Incoming: []string{"Flow_0uffa0a"},
												Outgoing: []string{"Flow_10j3phf"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											SignalEventDefinitions: []element.SignalEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "SignalEventDefinition_13bzj0y",
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
														ID: "Event_0pi17ux",
													},
													Name: "Intermediate Event Conditional Catch",
												},
												Incoming: []string{"Flow_0rs1jsj"},
												Outgoing: []string{"Flow_1luuks1"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											ConditionalEventDefinitions: []element.ConditionalEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "ConditionalEventDefinition_1hzmt1w",
															},
														},
													},
													Condition: element.ExpressionUnMarshal{
														Type:                   "tFormalExpression",
														ExpressionSubstitution: &element.FormalExpression{},
													},
												},
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
														ID: "Activity_00kk0w1",
													},
													Name: "Service Task 4",
												},
												Incoming: []string{"Flow_1fmiuuz"},
												Outgoing: []string{"Flow_0uffa0a"},
											},
										},
									},
								},
							},
							InclusiveGatewaies: []element.InclusiveGateway{
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Gateway_0iz1sti",
												},
												Name: "Inclusive Gateway",
											},
											Incoming: []string{"Flow_0vyan0i"},
											Outgoing: []string{"Flow_1fmiuuz", "Flow_1lz1d6w"},
										},
									},
									Default: "Flow_1lz1d6w",
								},
							},
							UserTasks: []element.UserTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_07vdczu",
													},
													Name: "User Task 3",
												},
												Incoming: []string{"Flow_1qqutfv"},
												Outgoing: []string{"Flow_0vyan0i"},
											},
											Properties: []element.Property{
												{
													BaseElement: element.BaseElement{
														ID: "Property_0qjltqr",
													},
													Name: "__targetRef_placeholder",
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "DataInputAssociation_1442h6j",
														},
														SourceRef: "DataObjectReference_1rkw6dp",
														TargetRef: "Property_0qjltqr",
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
														ID: "Activity_0s6eb2u",
													},
													Name: "User Task 8",
												},
												Incoming: []string{"Flow_1b1ivbn"},
												Outgoing: []string{"Flow_0uy5u27"},
											},
										},
									},
								},
							},
							SendTasks: []element.SendTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_0pn331p",
													},
													Name: "Send Task 2",
												},
												Incoming: []string{"Flow_0xky4tu"},
												Outgoing: []string{"Flow_1qqutfv"},
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
													ID: "Activity_1uu0yo4",
												},
												Name: "Abstract Task 1",
											},
											Incoming: []string{"Flow_0gzdgm8"},
											Outgoing: []string{"Flow_0xky4tu"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_0thkytf",
												},
												Name: "Task 5",
											},
											Incoming: []string{"Flow_02vtnpn"},
											Outgoing: []string{"Flow_0t9173h"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_0ayd8ln",
												},
												Name: "Task 6",
											},
											Incoming: []string{"Flow_1exfm6h"},
											Outgoing: []string{"Flow_1md83y0"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_0wdnxu3",
												},
												Name: "Task 9",
											},
											Incoming: []string{"Flow_0i7jvk9"},
											Outgoing: []string{"Flow_0rs1jsj"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_0javi3i",
												},
												Name: "Task 10",
											},
											Incoming: []string{"Flow_1luuks1"},
											Outgoing: []string{"Flow_1vs1xz8"},
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
														ID: "StartEvent_1",
													},
													Name: "Start Event 1 Timer",
												},
												Outgoing: []string{"Flow_0gzdgm8"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											TimerEventDefinitions: []element.TimerEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "TimerEventDefinition_1wg06e7",
															},
														},
													},
												},
											},
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
													ID: "Activity_0qnc8vy",
												},
												Name: "Call Activity calling a Global User Task",
											},
											Incoming: []string{"Flow_1lz1d6w"},
											Outgoing: []string{"Flow_1iyrc7z"},
										},
									},
								},
							},
							DataObjectReferenes: []element.DataObjectReference{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObjectReference_1rkw6dp",
										},
										Name: "Data Object",
									},
									DataObjectRef: "DataObject_0v253gd",
								},
							},
							DataObjects: []element.DataObject{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObject_0v253gd",
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
														ID: "Event_1f35b4w",
													},
												},
												Outgoing: []string{"Flow_0i7jvk9"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											EscalationEventDefinitions: []element.EscalationEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "EscalationEventDefinition_1e7kboz",
															},
														},
													},
												},
											},
										},
									},
									CancelActivity: false,
									AttachedToRef:  "Activity_0s6eb2u",
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_1rsgo7a",
													},
													Name: "Boundary Intermediate Event Non-Interrupting Conditional",
												},
												Outgoing: []string{"Flow_1exfm6h"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											ConditionalEventDefinitions: []element.ConditionalEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "ConditionalEventDefinition_1uszs8x",
															},
														},
													},
													Condition: element.ExpressionUnMarshal{
														Type:                   "tFormalExpression",
														ExpressionSubstitution: &element.FormalExpression{},
													},
												},
											},
										},
									},
									CancelActivity: false,
									AttachedToRef:  "Activity_0thkytf",
								},
							},
							EndEvents: []element.EndEvent{
								{
									ThrowEvent: element.ThrowEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_0403g4k",
													},
													Name: "End Event 1 Message",
												},
												Incoming: []string{"Flow_19165ws"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "MessageEventDefinition_00eayw1",
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
														ID: "Event_1cb9pew",
													},
													Name: "End Event 3 Signal",
												},
												Incoming: []string{"Flow_1vs1xz8"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											SignalEventDefinitions: []element.SignalEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "SignalEventDefinition_1v2uxma",
															},
														},
													},
												},
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
													ID: "Gateway_0y5g78s",
												},
												Name: "Parallel Gateway 2",
											},
											Incoming: []string{"Flow_0t9173h", "Flow_1md83y0", "Flow_0uy5u27"},
											Outgoing: []string{"Flow_19165ws"},
										},
									},
								},
							},
						},
						Artifacts: element.Artifacts{
							TextAnnotations: []element.TextAnnotation{
								{
									Artifact: element.Artifact{
										BaseElement: element.BaseElement{
											ID: "TextAnnotation_06uwg6s",
										},
									},
									Text: "Annotation",
								},
							},
							Associations: []element.Association{
								{
									Artifact: element.Artifact{
										BaseElement: element.BaseElement{
											ID: "Association_0ahr1rq",
										},
									},
									SourceRef: "Activity_0ffbdg4",
									TargetRef: "TextAnnotation_06uwg6s",
								},
							},
						},
					},
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "Process_1xz7va4",
								},
							},
						},
						IsExecutable: false,
						LaneSets: []element.LaneSet{
							{
								BaseElement: element.BaseElement{
									ID: "LaneSet_02xvcik",
								},
								Lanes: []element.Lane{
									{
										BaseElement: element.BaseElement{
											ID: "Lane_07fps3z",
										},
										Name: "Lane 1",
										FlowNodeRefs: []string{
											"Event_0n9dn9q",
											"Gateway_1n5jtse",
											"Activity_0dbod77",
											"Event_0yxib88",
											"Event_0dwi53b",
											"Activity_1i3ap44",
											"Event_0b4vgd7",
											"Activity_1xwdlgt",
											"Activity_0u4otp8",
											"Activity_1gsxnu3",
											"Activity_1kayloh",
											"Event_1f1af1t",
											"Activity_01h1c3z",
											"Activity_0wxomwg",
											"Event_1qwowa6",
											"Gateway_156l12u",
											"Activity_1fag12l",
											"Activity_1a394xz",
											"Activity_17ocq3z",
											"Activity_1ig20jv",
											"Event_1gbx6rp",
											"Event_1brbeoj",
											"Event_1e5rzkm",
											"Event_087hfz9",
											"Event_0drv2q5",
											"Event_0gdj0lz",
											"Event_1qgdkdf",
											"Activity_1qy7xgp",
											"Event_009wd9m",
											"Event_1lph52r",
											"Event_1r9ahil",
										},
									},
									{
										BaseElement: element.BaseElement{
											ID: "Lane_1eeat7r",
										},
										Name: "Lane 2",
										FlowNodeRefs: []string{
											"Activity_1c0usbk",
											"Activity_19dq1xp",
											"Activity_11l0l4d",
											"Activity_1cwg4vs",
											"Event_1rx2a8f",
											"Activity_1rcfmjq",
											"Activity_06qihxs",
											"Event_15p2cvq",
											"Activity_0xe0m5g",
											"Event_14urgdq",
											"Activity_15jh1tp",
											"Gateway_0ynnz6a",
											"Activity_0ahqya2",
											"Event_1qd721q",
											"Event_1ofkxhs",
											"Gateway_0gp6gne",
											"Gateway_198yhv3",
											"Event_0sedygx",
											"Event_0kp736q",
										},
									},
								},
							},
						},
						FlowElements: element.FlowElements{
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0m8w25k",
										},
									},
									SourceRef: "Activity_0dbod77",
									TargetRef: "Gateway_1n5jtse",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_13lye0k",
										},
									},
									SourceRef: "Gateway_1n5jtse",
									TargetRef: "Event_0n9dn9q",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0xjmqod",
										},
									},
									SourceRef: "Gateway_1n5jtse",
									TargetRef: "Event_0gdj0lz",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_051lggk",
										},
									},
									SourceRef: "Gateway_1n5jtse",
									TargetRef: "Event_1qgdkdf",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_138p74x",
										},
									},
									SourceRef: "Event_0yxib88",
									TargetRef: "Activity_0dbod77",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1j9rxo8",
										},
									},
									SourceRef: "Event_0n9dn9q",
									TargetRef: "Activity_1fag12l",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1cxmbe8",
										},
									},
									SourceRef: "Event_1qgdkdf",
									TargetRef: "Activity_1a394xz",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0xx9p19",
										},
									},
									SourceRef: "Activity_1fag12l",
									TargetRef: "Gateway_156l12u",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_16g4olt",
										},
										Name: "Default Sequence Flow 2",
									},
									SourceRef: "Gateway_156l12u",
									TargetRef: "Event_1qwowa6",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0vv68hn",
										},
									},
									SourceRef: "Event_1qwowa6",
									TargetRef: "Activity_0wxomwg",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0dk5zr9",
										},
									},
									SourceRef: "Event_087hfz9",
									TargetRef: "Activity_17ocq3z",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0prsfbf",
										},
									},
									SourceRef: "Gateway_156l12u",
									TargetRef: "Activity_17ocq3z",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0npffus",
										},
									},
									SourceRef: "Activity_01h1c3z",
									TargetRef: "Activity_1qy7xgp",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_10q4nrz",
										},
									},
									SourceRef: "Event_009wd9m",
									TargetRef: "Event_1lph52r",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0rx4vqb",
										},
									},
									SourceRef: "Activity_1qy7xgp",
									TargetRef: "Activity_1i3ap44",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0qvdwb3",
										},
									},
									SourceRef: "Activity_1i3ap44",
									TargetRef: "Event_0dwi53b",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0yk73ad",
										},
									},
									SourceRef: "Activity_0wxomwg",
									TargetRef: "Activity_1xwdlgt",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_07zdeu5",
										},
									},
									SourceRef: "Activity_1xwdlgt",
									TargetRef: "Event_0b4vgd7",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_01k51sm",
										},
									},
									SourceRef: "Activity_1ig20jv",
									TargetRef: "Activity_0u4otp8",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0zgmk02",
										},
									},
									SourceRef: "Activity_17ocq3z",
									TargetRef: "Event_1f1af1t",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0zh9gvr",
										},
									},
									SourceRef: "Event_1f1af1t",
									TargetRef: "Activity_0u4otp8",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0twe4nl",
										},
									},
									SourceRef: "Event_1r9ahil",
									TargetRef: "Activity_1gsxnu3",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0qd2cjd",
										},
									},
									SourceRef: "Activity_1gsxnu3",
									TargetRef: "Event_0b4vgd7",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1ef3aqe",
										},
									},
									SourceRef: "Activity_0u4otp8",
									TargetRef: "Event_0b4vgd7",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1p5rcqt",
										},
									},
									SourceRef: "Event_0drv2q5",
									TargetRef: "Activity_1kayloh",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0vzv0ky",
										},
									},
									SourceRef: "Activity_1a394xz",
									TargetRef: "Activity_1ig20jv",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_02c98k2",
										},
									},
									SourceRef: "Event_1qd721q",
									TargetRef: "Activity_0ahqya2",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0jmy9bt",
										},
									},
									SourceRef: "Activity_0ahqya2",
									TargetRef: "Gateway_0ynnz6a",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_163ttxe",
										},
									},
									SourceRef: "Gateway_0ynnz6a",
									TargetRef: "Activity_15jh1tp",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1yz425v",
										},
									},
									SourceRef: "Activity_15jh1tp",
									TargetRef: "Event_14urgdq",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_11zr1a3",
										},
									},
									SourceRef: "Event_14urgdq",
									TargetRef: "Gateway_0gp6gne",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0uaxzdx",
										},
									},
									SourceRef: "Event_1e5rzkm",
									TargetRef: "Activity_1c0usbk",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_00uq1g7",
										},
									},
									SourceRef: "Activity_1c0usbk",
									TargetRef: "Gateway_0gp6gne",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0ax4o0j",
										},
									},
									SourceRef: "Gateway_0gp6gne",
									TargetRef: "Activity_0xe0m5g",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1a59vjg",
										},
									},
									SourceRef: "Event_1gbx6rp",
									TargetRef: "Activity_1rcfmjq",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1k2of3s",
										},
									},
									SourceRef: "Activity_1rcfmjq",
									TargetRef: "Event_1rx2a8f",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0oyt612",
										},
									},
									SourceRef: "Event_1brbeoj",
									TargetRef: "Activity_06qihxs",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_07hgdqa",
										},
									},
									SourceRef: "Activity_06qihxs",
									TargetRef: "Event_0sedygx",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0dh0nnk",
										},
									},
									SourceRef: "Activity_1kayloh",
									TargetRef: "Event_0sedygx",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0fop3s1",
										},
									},
									SourceRef: "Gateway_0ynnz6a",
									TargetRef: "Activity_19dq1xp",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0xqi2l5",
										},
									},
									SourceRef: "Event_1ofkxhs",
									TargetRef: "Activity_19dq1xp",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0ox4w4d",
										},
									},
									SourceRef: "Activity_19dq1xp",
									TargetRef: "Activity_11l0l4d",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1i60rky",
										},
									},
									SourceRef: "Event_15p2cvq",
									TargetRef: "Activity_1cwg4vs",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0r5tyx6",
										},
									},
									SourceRef: "Activity_1cwg4vs",
									TargetRef: "Event_0kp736q",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0ly8iho",
										},
									},
									SourceRef: "Activity_11l0l4d",
									TargetRef: "Gateway_198yhv3",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0rkhydx",
										},
									},
									SourceRef: "Activity_0xe0m5g",
									TargetRef: "Gateway_198yhv3",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1tqpq63",
										},
									},
									SourceRef: "Gateway_198yhv3",
									TargetRef: "Event_0sedygx",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0l823y4",
										},
									},
									SourceRef: "Event_0gdj0lz",
									TargetRef: "Activity_01h1c3z",
								},
							},
							IntermediateThrowEvents: []element.IntermediateThrowEvent{
								{
									ThrowEvent: element.ThrowEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_1qwowa6",
													},
													Name: "Intermediate Event Message Throw",
												},
												Incoming: []string{"Flow_16g4olt"},
												Outgoing: []string{"Flow_0vv68hn"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "MessageEventDefinition_0io8u5r",
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
														ID: "Event_1f1af1t",
													},
												},
												Incoming: []string{"Flow_0zgmk02"},
												Outgoing: []string{"Flow_0zh9gvr"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											EscalationEventDefinitions: []element.EscalationEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "EscalationEventDefinition_185ama0",
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
														ID: "Event_1lph52r",
													},
													Name: "Intermediate End Event",
												},
												Incoming: []string{"Flow_10q4nrz"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											LinkEventDefinitions: []element.LinkEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "LinkEventDefinition_0a59ub9",
															},
														},
													},
													Name: "",
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
														ID: "Event_14urgdq",
													},
													Name: "Intermediate Event Signal Throw 2",
												},
												Incoming: []string{"Flow_1yz425v"},
												Outgoing: []string{"Flow_11zr1a3"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											SignalEventDefinitions: []element.SignalEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "SignalEventDefinition_0jhkcy9",
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
														ID: "Event_0n9dn9q",
													},
													Name: "Intermediate Event Message Catch",
												},
												Incoming: []string{"Flow_13lye0k"},
												Outgoing: []string{"Flow_1j9rxo8"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "MessageEventDefinition_0ox4djh",
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
														ID: "Event_0gdj0lz",
													},
													Name: "Intermediate Event Timer Catch",
												},
												Incoming: []string{"Flow_0xjmqod"},
												Outgoing: []string{"Flow_0l823y4"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											TimerEventDefinitions: []element.TimerEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "TimerEventDefinition_1vvsiz9",
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
														ID: "Event_1qgdkdf",
													},
													Name: "Intermediate Event Message Catch 2",
												},
												Incoming: []string{"Flow_051lggk"},
												Outgoing: []string{"Flow_1cxmbe8"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "MessageEventDefinition_1nec7v4",
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
														ID: "Event_1ofkxhs",
													},
													Name: "Intermediate Event Link",
												},
												Outgoing: []string{"Flow_0xqi2l5"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											LinkEventDefinitions: []element.LinkEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "LinkEventDefinition_0p5xwz6",
															},
														},
													},
													Name: "",
												},
											},
										},
									},
								},
							},
							EventBasedGatewaies: []element.EventBasedGateway{
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Gateway_1n5jtse",
												},
												Name: "Event Base Gateway 3",
											},
											Incoming: []string{"Flow_0m8w25k"},
											Outgoing: []string{"Flow_13lye0k", "Flow_0xjmqod", "Flow_051lggk"},
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
													ID: "Activity_0dbod77",
												},
												Name: "Task 11",
											},
											Incoming: []string{"Flow_138p74x"},
											Outgoing: []string{"Flow_0m8w25k"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_1a394xz",
												},
												Name: "Task 21",
											},
											Incoming: []string{"Flow_1cxmbe8"},
											Outgoing: []string{"Flow_0vzv0ky"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_17ocq3z",
												},
												Name: "Task 18",
											},
											Incoming: []string{"Flow_0dk5zr9", "Flow_0prsfbf"},
											Outgoing: []string{"Flow_0zgmk02"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_1c0usbk",
												},
												Name: "Task 27",
											},
											Incoming: []string{"Flow_0uaxzdx"},
											Outgoing: []string{"Flow_00uq1g7"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_11l0l4d",
												},
												Name: "Task 32",
											},
											Incoming: []string{"Flow_0ox4w4d"},
											Outgoing: []string{"Flow_0ly8iho"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_1cwg4vs",
												},
												Name: "Task 33",
											},
											Incoming: []string{"Flow_1i60rky"},
											Outgoing: []string{"Flow_0r5tyx6"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_1rcfmjq",
												},
												Name: "Task 29",
											},
											Incoming: []string{"Flow_1a59vjg"},
											Outgoing: []string{"Flow_1k2of3s"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_06qihxs",
												},
												Name: "Task 30",
											},
											Incoming: []string{"Flow_0oyt612"},
											Outgoing: []string{"Flow_07hgdqa"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_1xwdlgt",
												},
												Name: "Task 17",
											},
											Incoming: []string{"Flow_0yk73ad"},
											Outgoing: []string{"Flow_07zdeu5"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_0u4otp8",
												},
												Name: "Task 23",
											},
											Incoming: []string{"Flow_01k51sm", "Flow_0zh9gvr"},
											Outgoing: []string{"Flow_1ef3aqe"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_1gsxnu3",
												},
												Name: "Task 19",
											},
											Incoming: []string{"Flow_0twe4nl"},
											Outgoing: []string{"Flow_0qd2cjd"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_1kayloh",
												},
												Name: "Task 24",
											},
											Incoming: []string{"Flow_1p5rcqt"},
											Outgoing: []string{"Flow_0dh0nnk"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_0xe0m5g",
												},
												Name: "Task 28",
											},
											Incoming: []string{"Flow_0ax4o0j"},
											Outgoing: []string{"Flow_0rkhydx"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_15jh1tp",
												},
												Name: "Task 26",
											},
											Incoming: []string{"Flow_163ttxe"},
											Outgoing: []string{"Flow_1yz425v"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_0ahqya2",
												},
												Name: "Task 25",
											},
											Incoming: []string{"Flow_02c98k2"},
											Outgoing: []string{"Flow_0jmy9bt"},
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
														ID: "Event_0yxib88",
													},
													Name: "Start Event 2 Message",
												},
												Outgoing: []string{"Flow_138p74x"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "MessageEventDefinition_1gmfphr",
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
														ID: "Event_1qd721q",
													},
													Name: "Start Event 6 Signal",
												},
												Outgoing: []string{"Flow_02c98k2"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											SignalEventDefinitions: []element.SignalEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "SignalEventDefinition_02lr7ru",
															},
														},
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
											ID: "DataStoreReference_0t1ift6",
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
													ID: "Activity_01h1c3z",
												},
												Name: "Expanded Call Activity",
											},
											Incoming: []string{"Flow_0l823y4"},
											Outgoing: []string{"Flow_0npffus"},
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
																	ID: "Event_05w6wy8",
																},
																Name: "Start Event 3",
															},
															Outgoing: []string{"Flow_18dr9jm"},
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
																	ID: "Event_1wq0sy2",
																},
																Name: "Start Event 4 Conditional",
															},
															Outgoing: []string{"Flow_1sl2jua"},
														},
													},
													EventDefinitions: element.EventDefinitions{
														ConditionalEventDefinitions: []element.ConditionalEventDefinition{
															{
																EventDefinition: element.EventDefinition{
																	RootElement: element.RootElement{
																		BaseElement: element.BaseElement{
																			ID: "ConditionalEventDefinition_0cq4mx7",
																		},
																	},
																},
																Condition: element.ExpressionUnMarshal{
																	Type:                   "tFormalExpression",
																	ExpressionSubstitution: &element.FormalExpression{},
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
														ID: "Flow_18dr9jm",
													},
												},
												SourceRef: "Event_05w6wy8",
												TargetRef: "Activity_1lrpf00",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_1sl2jua",
													},
												},
												SourceRef: "Event_1wq0sy2",
												TargetRef: "Activity_1lrpf00",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_13o472r",
													},
												},
												SourceRef: "Activity_1lrpf00",
												TargetRef: "Activity_0g3cmxi",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_0agposr",
													},
												},
												SourceRef: "Activity_0g3cmxi",
												TargetRef: "Activity_0pyeguq",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_0e3d4iw",
													},
												},
												SourceRef: "Event_0mg1ari",
												TargetRef: "Event_0a5qvbg",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_10obnjq",
													},
												},
												SourceRef: "Activity_0pyeguq",
												TargetRef: "Event_153rv56",
											},
										},
										UserTasks: []element.UserTask{
											{
												Task: element.Task{
													Activity: element.Activity{
														FlowNode: element.FlowNode{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "Activity_1lrpf00",
																},
																Name: "User Task 12 Multi-inst Seq.",
															},
															Incoming: []string{"Flow_18dr9jm", "Flow_1sl2jua"},
															Outgoing: []string{"Flow_13o472r"},
														},
														LoopCharacteristicsElements: element.LoopCharacteristicsElements{
															MultielementLoopCharacteristics: []element.MultiInstanceLoopCharacteristics{{
																IsSequential: true,
															}},
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
																	ID: "Activity_0g3cmxi",
																},
																Name: "User Task 13",
															},
															Incoming: []string{"Flow_13o472r"},
															Outgoing: []string{"Flow_0agposr"},
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
																	ID: "Activity_0pyeguq",
																},
																Name: "Service Task 14",
															},
															Incoming: []string{"Flow_0agposr"},
															Outgoing: []string{"Flow_10obnjq"},
														},
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
																	ID: "Event_0mg1ari",
																},
																Name: "Boundary Intermediate Event Interrupting Message",
															},
															Outgoing: []string{"Flow_0e3d4iw"},
														},
													},
													EventDefinitions: element.EventDefinitions{
														MessageEventDefinitions: []element.MessageEventDefinition{
															{
																EventDefinition: element.EventDefinition{
																	RootElement: element.RootElement{
																		BaseElement: element.BaseElement{
																			ID: "MessageEventDefinition_0latlfm",
																		},
																	},
																},
															},
														},
													},
												},
												CancelActivity: true,
												AttachedToRef: "Activity_0g3cmxi",
											},
										},
										EndEvents: []element.EndEvent{
											{
												ThrowEvent: element.ThrowEvent{
													Event: element.Event{
														FlowNode: element.FlowNode{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "Event_0a5qvbg",
																},
																Name: "End Event 5 Terminate",
															},
															Incoming: []string{"Flow_0e3d4iw"},
														},
													},
													EventDefinitions: element.EventDefinitions{
														TerminateEventDefinitions: []element.TerminateEventDefinition{
															{
																EventDefinition: element.EventDefinition{
																	RootElement: element.RootElement{
																		BaseElement: element.BaseElement{
																			ID: "TerminateEventDefinition_14p6uf0",
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
																	ID: "Event_153rv56",
																},
																Name: "End Event 4",
															},
															Incoming: []string{"Flow_10obnjq"},
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
													ID: "Activity_1fag12l",
												},
												Name: "Collapsed Sub-Process 2",
											},
											Incoming: []string{"Flow_1j9rxo8"},
											Outgoing: []string{"Flow_0xx9p19"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_19dq1xp",
												},
												Name: "Expanded Sub-Process 3",
											},
											Incoming: []string{"Flow_0fop3s1", "Flow_0xqi2l5"},
											Outgoing: []string{"Flow_0ox4w4d"},
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
																	ID: "Event_0wm5nle",
																},
																Name: "Start Event 7",
															},
															Outgoing: []string{"Flow_01v892m"},
														},
													},
												},
											},
										},
										SequenceFlows: []element.SequenceFlow{
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_01v892m",
													},
												},
												SourceRef: "Event_0wm5nle",
												TargetRef: "Event_16rqalg",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_0shj0yz",
													},
												},
												SourceRef: "Event_16rqalg",
												TargetRef: "Activity_03j7y8k",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_1nnas72",
													},
												},
												SourceRef: "Activity_03j7y8k",
												TargetRef: "Gateway_10qt5ud",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_0bwg7yc",
													},
												},
												SourceRef: "Gateway_10qt5ud",
												TargetRef: "Event_0ivkjpo",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_1ierw66",
													},
												},
												SourceRef: "Gateway_10qt5ud",
												TargetRef: "Event_1wir0h9",
											},
										},
										IntermediateCatchEvents: []element.IntermediateCatchEvent{
											{
												CatchEvent: element.CatchEvent{
													Event: element.Event{
														FlowNode: element.FlowNode{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "Event_16rqalg",
																},
																Name: "Intermediate Event Signal Catch",
															},
															Incoming: []string{"Flow_01v892m"},
															Outgoing: []string{"Flow_0shj0yz"},
														},
													},
													EventDefinitions: element.EventDefinitions{
														SignalEventDefinitions: []element.SignalEventDefinition{
															{
																EventDefinition: element.EventDefinition{
																	RootElement: element.RootElement{
																		BaseElement: element.BaseElement{
																			ID: "SignalEventDefinition_0gdrh6v",
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
																ID: "Activity_03j7y8k",
															},
															Name: "Task 31",
														},
														Incoming: []string{"Flow_0shj0yz"},
														Outgoing: []string{"Flow_1nnas72"},
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
																ID: "Gateway_10qt5ud",
															},
															Name: "Exclusive Gateway",
														},
														Incoming: []string{"Flow_1nnas72"},
														Outgoing: []string{"Flow_0bwg7yc", "Flow_1ierw66"},
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
																	ID: "Event_0ivkjpo",
																},
																Name: "End Event 12",
															},
															Incoming: []string{"Flow_0bwg7yc"},
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
																	ID: "Event_1wir0h9",
																},
																Name: "End Event 13 Error",
															},
															Incoming: []string{"Flow_1ierw66"},
														},
													},
													EventDefinitions: element.EventDefinitions{
														ErrorEventDefinitions: []element.ErrorEventDefinition{
															{
																EventDefinition: element.EventDefinition{
																	RootElement: element.RootElement{
																		BaseElement: element.BaseElement{
																			ID: "ErrorEventDefinition_13if476",
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
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_1ig20jv",
												},
												Name: "Expanded Sub-Process 2",
											},
											Incoming: []string{"Flow_0vzv0ky"},
											Outgoing: []string{"Flow_01k51sm"},
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
																	ID: "Event_05uva8p",
																},
																Name: "Start Event 5 None",
															},
															Outgoing: []string{"Flow_0ted34p"},
														},
													},
												},
											},
										},
										SequenceFlows: []element.SequenceFlow{
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_0ted34p",
													},
												},
												SourceRef: "Event_05uva8p",
												TargetRef: "Activity_1k00lfn",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_07jjywp",
													},
												},
												SourceRef: "Activity_1k00lfn",
												TargetRef: "Event_0o1k45w",
											},
										},
										ServiceTasks: []element.ServiceTask{
											{
												Task: element.Task{
													Activity: element.Activity{
														FlowNode: element.FlowNode{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "Activity_1k00lfn",
																},
																Name: "Service Task 22",
															},
															Incoming: []string{"Flow_0ted34p"},
															Outgoing: []string{"Flow_07jjywp"},
														},
														LoopCharacteristicsElements: element.LoopCharacteristicsElements{
															MultielementLoopCharacteristics: []element.MultiInstanceLoopCharacteristics{{}},
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
																	ID: "Event_0o1k45w",
																},
																Name: "End Event 8 None",
															},
															Incoming: []string{"Flow_07jjywp"},
														},
													},
												},
											},
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
													ID: "Activity_0wxomwg",
												},
												Name: "Collapsed Call Activity",
											},
											Incoming: []string{"Flow_0vv68hn"},
											Outgoing: []string{"Flow_0yk73ad"},
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
														ID: "Event_087hfz9",
													},
													Name: "Boundary Intermediate Event Non-Interrupting Escalation",
												},
												Outgoing: []string{"Flow_0dk5zr9"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											EscalationEventDefinitions: []element.EscalationEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "EscalationEventDefinition_0pw5off",
															},
														},
													},
												},
											},
										},
									},
									CancelActivity: false,
									AttachedToRef:  "Activity_0wxomwg",
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_1e5rzkm",
													},
													Name: "Boundary Intermediate Event Interrupting Timer",
												},
												Outgoing: []string{"Flow_0uaxzdx"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											TimerEventDefinitions: []element.TimerEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "TimerEventDefinition_13e8rvq",
															},
														},
													},
												},
											},
										},
									},
									CancelActivity: true,
									AttachedToRef: "Activity_1a394xz",
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_15p2cvq",
													},
												},
												Outgoing: []string{"Flow_1i60rky"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											SignalEventDefinitions: []element.SignalEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "SignalEventDefinition_0e6ibwj",
															},
														},
													},
												},
											},
										},
									},
									CancelActivity: true,
									AttachedToRef: "Activity_11l0l4d",
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_1brbeoj",
													},
													Name: "Boundary Intermediate Event Non-Interrupting Timer",
												},
												Outgoing: []string{"Flow_0oyt612"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											TimerEventDefinitions: []element.TimerEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "TimerEventDefinition_0d4yvld",
															},
														},
													},
												},
											},
										},
									},
									CancelActivity: false,
									AttachedToRef:  "Activity_1ig20jv",
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_1gbx6rp",
													},
													Name: "Boundary Intermediate Event Interrupting Error",
												},
												Outgoing: []string{"Flow_1a59vjg"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											ErrorEventDefinitions: []element.ErrorEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "ErrorEventDefinition_17g5j05",
															},
														},
													},
												},
											},
										},
									},
									CancelActivity: true,
									AttachedToRef: "Activity_1ig20jv",
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_0drv2q5",
													},
													Name: "Boundary Intermediate Event Non-Interrupting Signal",
												},
												Outgoing: []string{"Flow_1p5rcqt"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											SignalEventDefinitions: []element.SignalEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "SignalEventDefinition_168in84",
															},
														},
													},
												},
											},
										},
									},
									CancelActivity: false,
									AttachedToRef:  "Activity_0u4otp8",
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_009wd9m",
													},
													Name: "Boundary Intermediate Event Interrupting Conditional",
												},
												Outgoing: []string{"Flow_10q4nrz"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											ConditionalEventDefinitions: []element.ConditionalEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "ConditionalEventDefinition_14j1wsr",
															},
														},
													},
													Condition: element.ExpressionUnMarshal{
														Type:                   "tFormalExpression",
														ExpressionSubstitution: &element.FormalExpression{},
													},
												},
											},
										},
									},
									CancelActivity: true,
									AttachedToRef: "Activity_1qy7xgp",
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_1r9ahil",
													},
													Name: "Boundary Intermediate Event Non-Interrupting Message",
												},
												Outgoing: []string{"Flow_0twe4nl"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "MessageEventDefinition_1qa9e6d",
															},
														},
													},
												},
											},
										},
									},
									CancelActivity: false,
									AttachedToRef: "Activity_1xwdlgt",
								},
							},
							ExclusiveGatewaies: []element.ExclusiveGateway{
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Gateway_156l12u",
												},
												Name: "Exclusive Gateway 4",
											},
											Incoming: []string{"Flow_0xx9p19"},
											Outgoing: []string{"Flow_16g4olt", "Flow_0prsfbf"},
										},
									},
									Default: "Flow_16g4olt",
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Gateway_198yhv3",
												},
												Name: "Parallel Gateway 7",
											},
											Incoming: []string{"Flow_0ly8iho", "Flow_0rkhydx"},
											Outgoing: []string{"Flow_1tqpq63"},
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
														ID: "Event_1rx2a8f",
													},
													Name: "End Event 10",
												},
												Incoming: []string{"Flow_1k2of3s"},
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
														ID: "Event_0dwi53b",
													},
													Name: "End Event 6 Message",
												},
												Incoming: []string{"Flow_0qvdwb3"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "MessageEventDefinition_1tw5clq",
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
														ID: "Event_0b4vgd7",
													},
													Name: "End Event 7 None",
												},
												Incoming: []string{"Flow_07zdeu5", "Flow_0qd2cjd", "Flow_1ef3aqe"},
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
														ID: "Event_0sedygx",
													},
													Name: "End Event 11 Escalation",
												},
												Incoming: []string{"Flow_07hgdqa", "Flow_0dh0nnk", "Flow_1tqpq63"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											EscalationEventDefinitions: []element.EscalationEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "EscalationEventDefinition_0frigdc",
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
														ID: "Event_0kp736q",
													},
													Name: "End Event 14",
												},
												Incoming: []string{"Flow_0r5tyx6"},
											},
										},
									},
								},
							},
							ReceiveTasks: []element.ReceiveTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_1i3ap44",
													},
													Name: "Receive Task 16",
												},
												Incoming: []string{"Flow_0rx4vqb"},
												Outgoing: []string{"Flow_0qvdwb3"},
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
														ID: "Activity_1qy7xgp",
													},
													Name: "Service Task 15",
												},
												Incoming: []string{"Flow_0npffus"},
												Outgoing: []string{"Flow_0rx4vqb"},
											},
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
													ID: "Gateway_0ynnz6a",
												},
												Name: "Parallel Gateway 5",
											},
											Incoming: []string{"Flow_0jmy9bt"},
											Outgoing: []string{"Flow_163ttxe", "Flow_0fop3s1"},
										},
									},
								},
							},
							InclusiveGatewaies: []element.InclusiveGateway{
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Gateway_0gp6gne",
												},
												Name: "Inclusive Gateway",
											},
											Incoming: []string{"Flow_11zr1a3", "Flow_00uq1g7"},
											Outgoing: []string{"Flow_0ax4o0j"},
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
								ID: "Category_10r7ibr",
							},
						},
						CategoryValues: []element.CategoryValue{
							{
								BaseElement: element.BaseElement{
									ID: "CategoryValue_02um9y9",
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
