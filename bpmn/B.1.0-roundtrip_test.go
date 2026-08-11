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

func TestB_1_0_roundtrip(t *testing.T) {
	path := "./test/B.1.0-roundtrip.bpmn"
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
			ID:              "_1373655174418",
			Name:            "B.1.0",
			TargetNamespace: "http://www.trisotech.com/definitions/_1373655174418",
			RootElemnts: element.RootElemnts{
				DataStores: []element.DataStore{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "DS1373655174514",
							},
						},
						Name:        "Data\nStore Reference",
						Capacity:    0,
						IsUnlimited: false,
					},
				},
				GlobalTasks: []element.GlobalTask{
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "global-task",
								},
							},
							Name: "Global Task",
						},
					},
				},
				Processes: []element.Process{
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "Process_ba16239e-181e-4b9f-bc5b-0bb2ee973450",
								},
							},
						},
						IsExecutable: false,
						FlowElements: element.FlowElements{
							StartEvents: []element.StartEvent{
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_200f43e7-1385-46e2-a380-3ef16ebe7847",
													},
													Name: "Start Event None 1",
												},
												Outgoing: []string{"_60ed96e6-5954-48de-861b-7d1e3c1fb23e"},
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
													ID: "_c57a5344-213f-4834-a6c3-94ce878b413c",
												},
												Name: "Abstract Task 4",
											},
											Incoming: []string{"_60ed96e6-5954-48de-861b-7d1e3c1fb23e"},
											Outgoing: []string{"_6c6288e8-43f6-4085-87c7-1ff21c38fe17"},
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
														ID: "_ed405919-9fd6-47d0-bb00-9be7d5467efb",
													},
													Name: "End Event None 2",
												},
												Incoming: []string{"_6c6288e8-43f6-4085-87c7-1ff21c38fe17"},
											},
										},
									},
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_60ed96e6-5954-48de-861b-7d1e3c1fb23e",
										},
									},
									SourceRef: "_200f43e7-1385-46e2-a380-3ef16ebe7847",
									TargetRef: "_c57a5344-213f-4834-a6c3-94ce878b413c",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_6c6288e8-43f6-4085-87c7-1ff21c38fe17",
										},
									},
									SourceRef: "_c57a5344-213f-4834-a6c3-94ce878b413c",
									TargetRef: "_ed405919-9fd6-47d0-bb00-9be7d5467efb",
								},
							},
						},
					},
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "WFP-6-1",
								},
							},
						},
						IsExecutable: false,
						FlowElements: element.FlowElements{
							StartEvents: []element.StartEvent{
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_e314751e-5c3a-41f2-a1ae-4cb99efa0916",
													},
													Name: "Start Event Timer",
												},
												Outgoing: []string{"_3eaa52c9-8d39-43d1-9528-b4047ff7fcdf"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											TimerEventDefinitions: []element.TimerEventDefinition{
												{
													TimeDate: element.ExpressionUnMarshal{
														Type:                   element.ExpressionTypeDefault,
														ExpressionSubstitution: &element.Expression{},
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
														ID: "_94efa7e0-2322-4fc3-a5bf-6c6296488927",
													},
													Name: "End Event None 1",
												},
												Incoming: []string{"_bbb25218-69a3-4401-823f-22f468cbd80d"},
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
													ID: "_219b9ca1-d4c5-497d-a4f7-06a44a6da20e",
												},
												Name: "Abstract Task 1",
											},
											Incoming: []string{"_3eaa52c9-8d39-43d1-9528-b4047ff7fcdf"},
											Outgoing: []string{"_a1505d79-bbc0-42cf-866a-401a2f94b675"},
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
														ID: "_f7eade87-bb98-47d3-85c7-66033a62b124",
													},
													Name: "User\nTask 2",
												},
												Incoming: []string{"_a1505d79-bbc0-42cf-866a-401a2f94b675"},
												Outgoing: []string{"_ba610e14-bf4c-4150-a1b1-460fe6a29f83"},
											},
										},
									},
									Implementation: "##WebService",
								},
							},
							ServiceTasks: []element.ServiceTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_ec919941-53ec-403d-97e1-6a163a063f21",
													},
													Name: "Service Task 3",
												},
												Incoming: []string{"_ba610e14-bf4c-4150-a1b1-460fe6a29f83"},
												Outgoing: []string{"_bbb25218-69a3-4401-823f-22f468cbd80d"},
											},
										},
									},
									Implementation: "##WebService",
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_3eaa52c9-8d39-43d1-9528-b4047ff7fcdf",
										},
									},
									SourceRef: "_e314751e-5c3a-41f2-a1ae-4cb99efa0916",
									TargetRef: "_219b9ca1-d4c5-497d-a4f7-06a44a6da20e",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_a1505d79-bbc0-42cf-866a-401a2f94b675",
										},
									},
									SourceRef: "_219b9ca1-d4c5-497d-a4f7-06a44a6da20e",
									TargetRef: "_f7eade87-bb98-47d3-85c7-66033a62b124",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_ba610e14-bf4c-4150-a1b1-460fe6a29f83",
										},
									},
									SourceRef: "_f7eade87-bb98-47d3-85c7-66033a62b124",
									TargetRef: "_ec919941-53ec-403d-97e1-6a163a063f21",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_bbb25218-69a3-4401-823f-22f468cbd80d",
										},
									},
									SourceRef: "_ec919941-53ec-403d-97e1-6a163a063f21",
									TargetRef: "_94efa7e0-2322-4fc3-a5bf-6c6296488927",
								},
							},
						},
					},
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "WFP-6-2",
								},
							},
						},
						IsExecutable: false,
						LaneSets: []element.LaneSet{
							{
								BaseElement: element.BaseElement{
									ID: "ls_0623a9bd-fd34-462a-b09d-85cb5004be78",
								},
								Lanes: []element.Lane{
									{
										BaseElement: element.BaseElement{
											ID: "_4a6df7ac-26d8-4718-ac05-90af463d5e23",
										},
										Name: "Lane 1",
										FlowNodeRefs: []string{
											"_2ee553a1-cb03-41e3-b285-345c826fc88d",
											"_fa3a8e53-5be0-4f0b-8680-d2498e255209",
											"_ba16239e-181e-4b9f-bc5b-0bb2ee973450",
											"_93021cd0-6f49-485d-966f-209744c748de",
											"_3c8c32c3-089a-4643-bf42-6c37c0dac7e0",
											"_a38484e2-7bdb-48b1-b62e-139d51d6a147",
											"_be29f267-9d56-46ef-8bbc-e13513b25fce",
											"_1237e756-d53c-4591-a731-dafffbf0b3f9",
										},
									},
									{
										BaseElement: element.BaseElement{
											ID: "_3400f56a-4565-47d1-91db-0ba17b958cb2",
										},
										Name: "Lane 2",
										FlowNodeRefs: []string{
											"_1eb62392-1f21-4a63-bbcb-c78880c3165e",
											"_7706e700-2aed-4b94-8070-961f118aab8f",
											"_ad81e6ba-40f5-43c1-9602-47d2e58804c8",
											"_33f30031-2e29-46b6-b080-30a192a36b45",
											"_7e6ccf38-e740-4537-a439-a8e984d066de",
											"_fea1c5af-6c76-403f-809e-26d476d92741",
											"_ae916437-d9aa-4e3d-a7c3-34998c410beb",
										},
									},
								},
							},
						},
						FlowElements: element.FlowElements{
							DataObjects: []element.DataObject{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DF1373655174778",
										},
										Name: "Data Object",
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
														ID: "_2ee553a1-cb03-41e3-b285-345c826fc88d",
													},
													Name: "End Event Message",
												},
												Incoming: []string{"_994697ca-8927-4c84-a9a6-8682f8dee032"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													MessageRef: "Message_1373655174960",
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
														ID: "_ae916437-d9aa-4e3d-a7c3-34998c410beb",
													},
													Name: "End Event Terminate",
												},
												Incoming: []string{"_54f45351-aa18-4c65-b0d0-edc3aa0a140d"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											TerminateEventDefinitions: []element.TerminateEventDefinition{{}},
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
													ID: "_fa3a8e53-5be0-4f0b-8680-d2498e255209",
												},
												Name: "Call Activity Calling a Global Task",
											},
											Incoming: []string{"_eeb6812d-d182-489f-aea2-493ab8732cfd"},
											Outgoing: []string{"_61abe245-5604-46ba-8152-94d6e68ffda4"},
										},
									},
									CalledElement: "global-task",
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_ba16239e-181e-4b9f-bc5b-0bb2ee973450",
												},
												Name: "Call Activity - Expanded",
											},
											Incoming: []string{"_10a16fd5-0d56-4fdb-8529-0a0610a573be"},
											Outgoing: []string{"_f5c5d52a-204f-4f97-b872-817d63cf36ab"},
										},
									},
									CalledElement: "Process_ba16239e-181e-4b9f-bc5b-0bb2ee973450",
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_1237e756-d53c-4591-a731-dafffbf0b3f9",
												},
												Name: "Call Activity Collapsed",
											},
											Incoming: []string{"_9d489bd9-9435-4692-bc98-4cdda4a61569"},
											Outgoing: []string{"_10a16fd5-0d56-4fdb-8529-0a0610a573be"},
										},
									},
									CalledElement: "WFP-0-",
								},
							},
							ExclusiveGatewaies: []element.ExclusiveGateway{
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_93021cd0-6f49-485d-966f-209744c748de",
												},
												Name: "Exclusive Gateway Convergence 1",
											},
											Incoming: []string{"_61abe245-5604-46ba-8152-94d6e68ffda4", "_f5c5d52a-204f-4f97-b872-817d63cf36ab"},
											Outgoing: []string{"_994697ca-8927-4c84-a9a6-8682f8dee032"},
										},
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_3c8c32c3-089a-4643-bf42-6c37c0dac7e0",
												},
												Name: "Exclusive Gateway Divergence 1",
											},
											Incoming: []string{"_d30f7fb3-ec91-4a60-b73e-42419417f3be"},
											Outgoing: []string{"_9d489bd9-9435-4692-bc98-4cdda4a61569", "_eeb6812d-d182-489f-aea2-493ab8732cfd"},
										},
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_ad81e6ba-40f5-43c1-9602-47d2e58804c8",
												},
												Name: "Exclusive Gateway Divergence 2",
											},
											Incoming: []string{"_00d30c39-29a7-4c36-86e3-bc6e893efb42"},
											Outgoing: []string{"_6ee42e88-3d90-4259-83c0-9abd4574a15a", "_6a248585-952e-40ff-82ec-b6d8f410b73a"},
										},
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_33f30031-2e29-46b6-b080-30a192a36b45",
												},
												Name: "Exclusive Gateway Convergence 2",
											},
											Incoming: []string{"_8f68b889-83a4-44ad-9777-d39acdd5415e", "_657f30ba-0690-4a14-8b8e-d8939efcc7bd"},
											Outgoing: []string{"_54f45351-aa18-4c65-b0d0-edc3aa0a140d"},
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
														ID: "_a38484e2-7bdb-48b1-b62e-139d51d6a147",
													},
													Name: "Start Event\nMessage",
												},
												Outgoing: []string{"_a63c8cd6-eee8-4fbe-be5e-f6980b180b52"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													MessageRef: "Message_1373655174959",
												},
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
													ID: "_be29f267-9d56-46ef-8bbc-e13513b25fce",
												},
												Name: "Parallel Gateway Divergence",
											},
											Incoming: []string{"_a63c8cd6-eee8-4fbe-be5e-f6980b180b52"},
											Outgoing: []string{"_ab34472d-95a4-459c-a13b-5ed8b8b75eca", "_d30f7fb3-ec91-4a60-b73e-42419417f3be"},
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
													ID: "_1eb62392-1f21-4a63-bbcb-c78880c3165e",
												},
												Name: "Collapsed Sub-Process",
											},
											Incoming: []string{"_6ee42e88-3d90-4259-83c0-9abd4574a15a"},
											Outgoing: []string{"_d234800f-72d3-46cb-b603-30f1da7b1205"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_7e6ccf38-e740-4537-a439-a8e984d066de",
												},
												Name: "Sub Process - Expanded",
											},
											Incoming: []string{"_d234800f-72d3-46cb-b603-30f1da7b1205"},
											Outgoing: []string{"_8f68b889-83a4-44ad-9777-d39acdd5415e"},
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
																	ID: "_1df01cbc-5d8c-444e-b1db-da3efdee254a",
																},
																Name: "Start Event None 2",
															},
															Outgoing: []string{"_2d1047ce-fdd5-4cb6-9f0c-0ee8d6d3044a"},
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
																ID: "_6936f794-7bbb-4aa1-ae48-3a35bab4e2f4",
															},
															Name: "Abstract Task 6",
														},
														Incoming: []string{"_2d1047ce-fdd5-4cb6-9f0c-0ee8d6d3044a"},
														Outgoing: []string{"_062ae395-4aba-408b-ac64-4987752be95b"},
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
																	ID: "_4f744697-3643-41a9-9d07-84c78e2df64b",
																},
																Name: "End Event None 3",
															},
															Incoming: []string{"_062ae395-4aba-408b-ac64-4987752be95b"},
														},
													},
												},
											},
										},
										SequenceFlows: []element.SequenceFlow{
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_2d1047ce-fdd5-4cb6-9f0c-0ee8d6d3044a",
													},
												},
												SourceRef: "_1df01cbc-5d8c-444e-b1db-da3efdee254a",
												TargetRef: "_6936f794-7bbb-4aa1-ae48-3a35bab4e2f4",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_062ae395-4aba-408b-ac64-4987752be95b",
													},
												},
												SourceRef: "_6936f794-7bbb-4aa1-ae48-3a35bab4e2f4",
												TargetRef: "_4f744697-3643-41a9-9d07-84c78e2df64b",
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
														ID: "_7706e700-2aed-4b94-8070-961f118aab8f",
													},
													Name: "User Task 5",
												},
												Incoming: []string{"_ab34472d-95a4-459c-a13b-5ed8b8b75eca"},
												Outgoing: []string{"_00d30c39-29a7-4c36-86e3-bc6e893efb42"},
											},
										},
									},
									Implementation: "##WebService",
								},
							},
							ServiceTasks: []element.ServiceTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_fea1c5af-6c76-403f-809e-26d476d92741",
													},
													Name: "Service Task 7",
												},
												Incoming: []string{"_6a248585-952e-40ff-82ec-b6d8f410b73a"},
												Outgoing: []string{"_657f30ba-0690-4a14-8b8e-d8939efcc7bd"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "Din1373655174951",
															},
														},
													},
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "Dout1373655174952",
															},
														},
													},
													InputSets: []element.InputSet{
														{
															DataInputRefs: []string{"Din1373655174951"},
														},
													},
													OutputSets: []element.OutputSet{
														{
															DataOutputRefs: []string{"Dout1373655174952"},
														},
													},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_73afd30d-7d54-4897-9350-1f7d301ef1b2",
														},
														SourceRef: "_3d35229f-2c75-4d5d-a066-2d14e46e442e",
														TargetRef: "Din1373655174951",
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_fa10ebaf-7088-4def-8cc3-d959b8876b06",
														},
														SourceRef: "Dout1373655174952",
														TargetRef: "_b9385abf-d293-40b7-848b-8add4db48415",
													},
												},
											},
										},
									},
									Implementation: "##WebService",
								},
							},
							DataObjectReferenes: []element.DataObjectReference{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_3d35229f-2c75-4d5d-a066-2d14e46e442e",
										},
										Name: "Data Object",
									},
									DataObjectRef: "DF1373655174778",
								},
							},
							DataStoreReferences: []element.DataStoreReference{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_b9385abf-d293-40b7-848b-8add4db48415",
										},
										Name: "Data\nStore Reference",
									},
									DataStoreRef: "DS1373655174514",
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_a63c8cd6-eee8-4fbe-be5e-f6980b180b52",
										},
									},
									SourceRef: "_a38484e2-7bdb-48b1-b62e-139d51d6a147",
									TargetRef: "_be29f267-9d56-46ef-8bbc-e13513b25fce",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_ab34472d-95a4-459c-a13b-5ed8b8b75eca",
										},
									},
									SourceRef: "_be29f267-9d56-46ef-8bbc-e13513b25fce",
									TargetRef: "_7706e700-2aed-4b94-8070-961f118aab8f",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_8f68b889-83a4-44ad-9777-d39acdd5415e",
										},
									},
									SourceRef: "_7e6ccf38-e740-4537-a439-a8e984d066de",
									TargetRef: "_33f30031-2e29-46b6-b080-30a192a36b45",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_d30f7fb3-ec91-4a60-b73e-42419417f3be",
										},
									},
									SourceRef: "_be29f267-9d56-46ef-8bbc-e13513b25fce",
									TargetRef: "_3c8c32c3-089a-4643-bf42-6c37c0dac7e0",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_10a16fd5-0d56-4fdb-8529-0a0610a573be",
										},
									},
									SourceRef: "_1237e756-d53c-4591-a731-dafffbf0b3f9",
									TargetRef: "_ba16239e-181e-4b9f-bc5b-0bb2ee973450",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_d234800f-72d3-46cb-b603-30f1da7b1205",
										},
									},
									SourceRef: "_1eb62392-1f21-4a63-bbcb-c78880c3165e",
									TargetRef: "_7e6ccf38-e740-4537-a439-a8e984d066de",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_61abe245-5604-46ba-8152-94d6e68ffda4",
										},
									},
									SourceRef: "_fa3a8e53-5be0-4f0b-8680-d2498e255209",
									TargetRef: "_93021cd0-6f49-485d-966f-209744c748de",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_00d30c39-29a7-4c36-86e3-bc6e893efb42",
										},
									},
									SourceRef: "_7706e700-2aed-4b94-8070-961f118aab8f",
									TargetRef: "_ad81e6ba-40f5-43c1-9602-47d2e58804c8",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_6ee42e88-3d90-4259-83c0-9abd4574a15a",
										},
									},
									SourceRef: "_ad81e6ba-40f5-43c1-9602-47d2e58804c8",
									TargetRef: "_1eb62392-1f21-4a63-bbcb-c78880c3165e",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_6a248585-952e-40ff-82ec-b6d8f410b73a",
										},
									},
									SourceRef: "_ad81e6ba-40f5-43c1-9602-47d2e58804c8",
									TargetRef: "_fea1c5af-6c76-403f-809e-26d476d92741",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_f5c5d52a-204f-4f97-b872-817d63cf36ab",
										},
									},
									SourceRef: "_ba16239e-181e-4b9f-bc5b-0bb2ee973450",
									TargetRef: "_93021cd0-6f49-485d-966f-209744c748de",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_657f30ba-0690-4a14-8b8e-d8939efcc7bd",
										},
									},
									SourceRef: "_fea1c5af-6c76-403f-809e-26d476d92741",
									TargetRef: "_33f30031-2e29-46b6-b080-30a192a36b45",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_54f45351-aa18-4c65-b0d0-edc3aa0a140d",
										},
									},
									SourceRef: "_33f30031-2e29-46b6-b080-30a192a36b45",
									TargetRef: "_ae916437-d9aa-4e3d-a7c3-34998c410beb",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_9d489bd9-9435-4692-bc98-4cdda4a61569",
										},
									},
									SourceRef: "_3c8c32c3-089a-4643-bf42-6c37c0dac7e0",
									TargetRef: "_1237e756-d53c-4591-a731-dafffbf0b3f9",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_eeb6812d-d182-489f-aea2-493ab8732cfd",
										},
									},
									SourceRef: "_3c8c32c3-089a-4643-bf42-6c37c0dac7e0",
									TargetRef: "_fa3a8e53-5be0-4f0b-8680-d2498e255209",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_994697ca-8927-4c84-a9a6-8682f8dee032",
										},
									},
									SourceRef: "_93021cd0-6f49-485d-966f-209744c748de",
									TargetRef: "_2ee553a1-cb03-41e3-b285-345c826fc88d",
								},
							},
						},
						Artifacts: element.Artifacts{
							TextAnnotations: []element.TextAnnotation{
								{
									Artifact: element.Artifact{
										BaseElement: element.BaseElement{
											ID: "_4815ea6a-ede2-489b-8b37-2cdb2835b02c",
										},
									},
									Text: "Text Annotation",
								},
							},
							Associations: []element.Association{
								{
									Artifact: element.Artifact{
										BaseElement: element.BaseElement{
											ID: "_5362a7ef-ce7e-4a91-9c38-66c07b1b5f49",
										},
									},
									AssociationDirection: "None",
									SourceRef:            "_1237e756-d53c-4591-a731-dafffbf0b3f9",
									TargetRef:            "_4815ea6a-ede2-489b-8b37-2cdb2835b02c",
								},
							},
						},
					},
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "WFP-0-",
								},
							},
						},
						IsExecutable: false,
						FlowElements: element.FlowElements{
							StartEvents: []element.StartEvent{
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_18770c5c-c117-4570-aaf2-8c7a6910c34d",
													},
													Name: "Start Event None 3",
												},
												Outgoing: []string{"_107de09e-8506-4d2b-ad00-3341be723dff"},
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
													ID: "_47207d3b-8dc2-4679-bf33-c1e7e677765b",
												},
												Name: "Abstract Task 8",
											},
											Incoming: []string{"_107de09e-8506-4d2b-ad00-3341be723dff"},
											Outgoing: []string{"_2e21b112-d974-4add-9bee-91dafbed0a22"},
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
														ID: "_ab12c75c-eaf3-4ae1-9021-ee556711757f",
													},
													Name: "End Event None 4",
												},
												Incoming: []string{"_2e21b112-d974-4add-9bee-91dafbed0a22"},
											},
										},
									},
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_107de09e-8506-4d2b-ad00-3341be723dff",
										},
									},
									SourceRef: "_18770c5c-c117-4570-aaf2-8c7a6910c34d",
									TargetRef: "_47207d3b-8dc2-4679-bf33-c1e7e677765b",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_2e21b112-d974-4add-9bee-91dafbed0a22",
										},
									},
									SourceRef: "_47207d3b-8dc2-4679-bf33-c1e7e677765b",
									TargetRef: "_ab12c75c-eaf3-4ae1-9021-ee556711757f",
								},
							},
						},
					},
				},
				Messages: []element.Message{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "Message_1373655174959",
							},
						},
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "Message_1373655174960",
							},
						},
					},
				},
				Categories: []element.Category{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "Cat1373655174961",
							},
						},
						CategoryValues: []element.CategoryValue{
							{
								BaseElement: element.BaseElement{
									ID: "Value_Cat1373655174961",
								},
								Value: "Group",
							},
						},
					},
				},
				Collaborations: []element.Collaboration{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "C1373655174958",
							},
						},
						Participants: []element.Participant{
							{
								BaseElement: element.BaseElement{
									ID: "_cde15ee4-b395-43a3-9f5e-9028446f8a52",
								},
								Name:       "Participant",
								ProcessRef: "WFP-6-1",
							},
							{
								BaseElement: element.BaseElement{
									ID: "_0623a9bd-fd34-462a-b09d-85cb5004be78",
								},
								Name:       "Pool",
								ProcessRef: "WFP-6-2",
							},
						},
						MessageFlows: []element.MessageFlow{
							{
								BaseElement: element.BaseElement{
									ID: "_5d195b1c-ffea-4b53-b98f-78d9616a5038",
								},
								Name:       "Message Flow 1",
								SourceRef:  "_219b9ca1-d4c5-497d-a4f7-06a44a6da20e",
								TargetRef:  "_a38484e2-7bdb-48b1-b62e-139d51d6a147",
								MessageRef: "Message_1373655174959",
							},
							{
								BaseElement: element.BaseElement{
									ID: "_9428f666-fc8a-41be-8a77-9b280e14e7ae",
								},
								Name:       "Message\nFlow 2",
								SourceRef:  "_2ee553a1-cb03-41e3-b285-345c826fc88d",
								TargetRef:  "_ec919941-53ec-403d-97e1-6a163a063f21",
								MessageRef: "Message_1373655174960",
							},
						},
						Artifacts: element.Artifacts{
							Groups: []element.Group{
								{
									Artifact: element.Artifact{
										BaseElement: element.BaseElement{
											ID: "_bd04180e-49f6-4cf0-a7d6-da59e2840b4b",
										},
									},
									CategoryValueRef: "Value_Cat1373655174961",
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
