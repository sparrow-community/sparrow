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

func TestB_2_0_roundtrip(t *testing.T) {
	path := "./test/B.2.0-roundtrip.bpmn"
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
			ID:              "_1373638079286",
			Name:            "B.2.0",
			TargetNamespace: "http://www.trisotech.com/definitions/_1373638079286",
			Exporter:        "",
			ExporterVersion: "",
			RootElemnts: element.RootElemnts{
				GlobalUserTasks: []element.GlobalUserTask{
					{
						GlobalTask: element.GlobalTask{
							CallableElement: element.CallableElement{
								RootElement: element.RootElement{
									BaseElement: element.BaseElement{
										ID: "_d4afdad7-65a5-40c6-9fee-e13ba4c0bed4",
									},
								},
								Name: "Call Activity calling a Global User Task",
							},
						},
						Implementation: "##WebService",
					},
				},
				DataStores: []element.DataStore{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "DS1373638079677",
							},
						},
						Name:        "Data Store Reference",
						Capacity:    0,
						IsUnlimited: false,
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
													Name: "Start Event 3",
												},
												Outgoing: []string{"_60ed96e6-5954-48de-861b-7d1e3c1fb23e"},
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
														ID: "_cba8fbed-2bb6-40a9-8ac5-83e827ce9d9f",
													},
													Name: "Start Event 4 Conditional",
												},
												Outgoing: []string{"_c72ddf62-21b8-4895-a7bf-7a765f5981eb"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											ConditionalEventDefinitions: []element.ConditionalEventDefinition{
												{
													Condition: element.ExpressionUnMarshal{
														ExpressionSubstitution: &element.Expression{},
													},
												},
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
														ID: "_c57a5344-213f-4834-a6c3-94ce878b413c",
													},
													Name: "User Task 12 Muti-Inst. Seq.",
												},
												Incoming: []string{"_60ed96e6-5954-48de-861b-7d1e3c1fb23e", "_c72ddf62-21b8-4895-a7bf-7a765f5981eb"},
												Outgoing: []string{"_6c6288e8-43f6-4085-87c7-1ff21c38fe17"},
											},
											LoopCharacteristicsElements: element.LoopCharacteristicsElements{
												MultielementLoopCharacteristics: []element.MultiInstanceLoopCharacteristics{
													{
														IsSequential: true,
													},
												},
											},
										},
									},
									Implementation: "##WebService",
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_7f4fe4ea-901f-4c74-bcd4-e933495712fd",
													},
													Name: "User Task 13",
												},
												Incoming: []string{"_6c6288e8-43f6-4085-87c7-1ff21c38fe17"},
												Outgoing: []string{"_088b81bd-2a43-43a7-86c9-34c1cdd5f11b"},
											},
										},
									},
									Implementation: "##WebService",
								},
							},
							BoundaryEvents: []element.BoundaryEvent{
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_86b052b4-225c-424e-b900-bb94bdd77cec",
													},
													Name: "Boundary Intermediate Event Interrupting Message",
												},
												Outgoing: []string{"_99e68e88-586e-4f77-9506-f30903307209"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{},
											},
										},
									},
									AttachedToRef: "_7f4fe4ea-901f-4c74-bcd4-e933495712fd",
								},
							},
							ServiceTasks: []element.ServiceTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_3bfec246-ab94-4807-a79a-3df91ac13800",
													},
													Name: "Service Task 14",
												},
												Incoming: []string{"_088b81bd-2a43-43a7-86c9-34c1cdd5f11b"},
												Outgoing: []string{"_9ff30396-5ef2-42c0-ab4b-5343fbb0af21"},
											},
										},
									},
									Implementation: "##WebService",
								},
							},
							EndEvents: []element.EndEvent{
								{
									ThrowEvent: element.ThrowEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_778ff738-a5af-4373-a8da-0fbbfae9e00a",
													},
													Name: "End Event 5 Terminate",
												},
												Incoming: []string{"_99e68e88-586e-4f77-9506-f30903307209"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											TerminateEventDefinitions: []element.TerminateEventDefinition{
												{},
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
														ID: "_ed405919-9fd6-47d0-bb00-9be7d5467efb",
													},
													Name: "End Event 4",
												},
												Incoming: []string{"_9ff30396-5ef2-42c0-ab4b-5343fbb0af21"},
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
										Name: "",
									},
									SourceRef: "_200f43e7-1385-46e2-a380-3ef16ebe7847",
									TargetRef: "_c57a5344-213f-4834-a6c3-94ce878b413c",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_6c6288e8-43f6-4085-87c7-1ff21c38fe17",
										},
										Name: "",
									},
									SourceRef: "_c57a5344-213f-4834-a6c3-94ce878b413c",
									TargetRef: "_7f4fe4ea-901f-4c74-bcd4-e933495712fd",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_9ff30396-5ef2-42c0-ab4b-5343fbb0af21",
										},
										Name: "",
									},
									SourceRef: "_3bfec246-ab94-4807-a79a-3df91ac13800",
									TargetRef: "_ed405919-9fd6-47d0-bb00-9be7d5467efb",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_99e68e88-586e-4f77-9506-f30903307209",
										},
										Name: "",
									},
									SourceRef: "_86b052b4-225c-424e-b900-bb94bdd77cec",
									TargetRef: "_778ff738-a5af-4373-a8da-0fbbfae9e00a",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_c72ddf62-21b8-4895-a7bf-7a765f5981eb",
										},
										Name: "",
									},
									SourceRef: "_cba8fbed-2bb6-40a9-8ac5-83e827ce9d9f",
									TargetRef: "_c57a5344-213f-4834-a6c3-94ce878b413c",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_088b81bd-2a43-43a7-86c9-34c1cdd5f11b",
										},
									},
									SourceRef: "_7f4fe4ea-901f-4c74-bcd4-e933495712fd",
									TargetRef: "_3bfec246-ab94-4807-a79a-3df91ac13800",
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
							DataObjects: []element.DataObject{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DF1373638080458",
										},
										Name: "Data Object",
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
														ID: "_c9870992-6643-4094-acfd-d76e5e37941b",
													},
													Name: "User Task 8",
												},
												Incoming: []string{"_e6537f9d-e5ea-4abc-a6e8-add13a11b536"},
												Outgoing: []string{"_8321494b-9c3c-483a-b4d2-79e6254211c0"},
											},
										},
									},
									Implementation: "##WebService",
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_0e87da16-736e-45b2-95e5-8f45940f3adf",
													},
													Name: "User Task 3",
												},
												Incoming: []string{"_e1f948a6-db76-4c35-ba36-e430eecb63a4"},
												Outgoing: []string{"_4474c77a-01bb-435e-929a-7eb71bd824a8"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "Din1373638080886",
															},
														},
													},
													InputSets: []element.InputSet{
														{
															DataInputRefs: []string{"Din1373638080886"},
														},
													},
													OutputSets: []element.OutputSet{{}},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_f906ca20-8666-41ff-9d37-b76e09ac4f94",
														},
														SourceRef: "_aa8c769a-276c-4589-b182-7c7bbd0a9e1e",
														TargetRef: "Din1373638080886",
													},
												},
											},
										},
									},
									Implementation: "##WebService",
								},
							},
							BoundaryEvents: []element.BoundaryEvent{
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_708d55c8-684a-4e3b-a69d-69c620cd0ac0",
													},
												},
												Outgoing: []string{"_da47ac6a-4ea1-4bf6-ad3e-8cc60c0ea8a9"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											EscalationEventDefinitions: []element.EscalationEventDefinition{
												{},
											},
										},
									},
									CancelActivity: false,
									AttachedToRef:  "_c9870992-6643-4094-acfd-d76e5e37941b",
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_732c0641-b12f-448b-b9f8-a68b355782e3",
													},
													Name: "Boundary Intermediate Event Non-Interrupting Conditional",
												},
												Outgoing: []string{"_54574511-c29e-4157-bb8f-8b26f15faaf5"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											ConditionalEventDefinitions: []element.ConditionalEventDefinition{
												{
													Condition: element.ExpressionUnMarshal{
														ExpressionSubstitution: &element.Expression{},
													},
												},
											},
										},
									},
									CancelActivity: false,
									AttachedToRef:  "_2a08c361-be51-437e-a86d-c62798c14e83",
								},
							},
							Tasks: []element.Task{
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_2a08c361-be51-437e-a86d-c62798c14e83",
												},
												Name: "Task 5",
											},
											Incoming: []string{"_dbae7e12-67c2-4256-8c50-5811c207ba55"},
											Outgoing: []string{"_b41b9c86-bc41-4b6e-ac32-70e1342e6128"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_8f41085a-3302-4cfc-9231-23bb4b43c7a1",
												},
												Name: "Task 9",
											},
											Incoming: []string{"_da47ac6a-4ea1-4bf6-ad3e-8cc60c0ea8a9"},
											Outgoing: []string{"_5f52bb2a-c49d-43ae-9253-7b89a8251b2b"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_a01498ae-086c-4adc-9229-ec3135bc2bcf",
												},
												Name: "Task 10",
											},
											Incoming: []string{"_e503a157-f034-4a8a-ae12-aa524a8ebeee"},
											Outgoing: []string{"_875031ae-f87c-45a3-ae60-ba0cb0ce0bc1"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_f3698a93-fe0b-4f49-bfa6-68d055936af3",
												},
												Name: "Task 6",
											},
											Incoming: []string{"_54574511-c29e-4157-bb8f-8b26f15faaf5"},
											Outgoing: []string{"_de249cbe-d431-4e1f-bc8c-a90aec438683"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_d84e5824-7bb7-4057-9dba-6c8794f7948c",
												},
												Name: "Abstract Task 1",
											},
											Incoming: []string{"_2c124731-3338-46e0-81b5-b435f34941c8"},
											Outgoing: []string{"_19a7094a-7bdf-4819-af86-be22d2cfdc49"},
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
														ID: "_4e71bf73-1719-401e-a9a2-85dc89fc1150",
													},
													Name: "Start Event 1 Timer",
												},
												Outgoing: []string{"_2c124731-3338-46e0-81b5-b435f34941c8"},
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
														ID: "_5cc02d0f-c090-4e48-8da3-f32cbbca9565",
													},
													Name: "End Event 3 Signal",
												},
												Incoming: []string{"_875031ae-f87c-45a3-ae60-ba0cb0ce0bc1"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											SignalEventDefinitions: []element.SignalEventDefinition{
												{},
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
														ID: "_b67ba682-c8d6-465b-b538-c287db18d1be",
													},
													Name: "End Event 1\nMessage",
												},
												Incoming: []string{"_2752e07c-f2e4-4384-8721-70e4abea9765"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{},
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
														ID: "_ac1fde31-c0cd-4a8a-9728-a5fb49602de7",
													},
													Name: "Service Task 4",
												},
												Incoming: []string{"_a9966baf-d9b9-4be1-a4d9-2906ab5add30"},
												Outgoing: []string{"_b8694c28-408b-4394-b476-41e2a5bcfd94"},
											},
										},
									},
									Implementation: "##WebService",
								},
							},
							InclusiveGatewaies: []element.InclusiveGateway{
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_dec393e7-f182-4d31-b05f-e33ac3a5e35f",
												},
												Name: "Inclusive Gateway 1",
											},
											Incoming: []string{"_4474c77a-01bb-435e-929a-7eb71bd824a8"},
											Outgoing: []string{"_a9966baf-d9b9-4be1-a4d9-2906ab5add30", "_be19c2da-316a-47f6-ad7b-eb6c82bf8609"},
										},
									},
									Default: "_be19c2da-316a-47f6-ad7b-eb6c82bf8609",
								},
							},
							ParallelGatewaies: []element.ParallelGateway{
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_397c783e-ad6a-4cf3-8266-9b41962c83bd",
												},
												Name: "Parallel Gateway 2",
											},
											Incoming: []string{"_8321494b-9c3c-483a-b4d2-79e6254211c0", "_b41b9c86-bc41-4b6e-ac32-70e1342e6128", "_de249cbe-d431-4e1f-bc8c-a90aec438683"},
											Outgoing: []string{"_2752e07c-f2e4-4384-8721-70e4abea9765"},
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
														ID: "_c9cb2415-6a2e-49d6-84b9-27babcde4088",
													},
													Name: "Intermediate Event Conditional Catch",
												},
												Incoming: []string{"_5f52bb2a-c49d-43ae-9253-7b89a8251b2b"},
												Outgoing: []string{"_e503a157-f034-4a8a-ae12-aa524a8ebeee"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											ConditionalEventDefinitions: []element.ConditionalEventDefinition{
												{
													Condition: element.ExpressionUnMarshal{
														ExpressionSubstitution: &element.Expression{},
													},
												},
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
														ID: "_76ee26df-2c95-495b-9d9a-cb806aea6baf",
													},
													Name: "Send Task 2",
												},
												Incoming: []string{"_19a7094a-7bdf-4819-af86-be22d2cfdc49"},
												Outgoing: []string{"_e1f948a6-db76-4c35-ba36-e430eecb63a4"},
											},
										},
									},
									Implementation: "##WebService",
									MessageRef:     "Message_1373638080955",
								},
							},
							SubProcesses: []element.SubProcess{
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_149a6e1d-0385-4d0f-a90c-c2150a291a67",
												},
												Name: "Collapsed Sub-Process 1 Multi-Instances",
											},
											Incoming: []string{"_e0fbb051-0e07-43cc-a8bf-158492541c0f"},
											Outgoing: []string{"_dbae7e12-67c2-4256-8c50-5811c207ba55"},
										},
										LoopCharacteristicsElements: element.LoopCharacteristicsElements{
											MultielementLoopCharacteristics: []element.MultiInstanceLoopCharacteristics{
												{},
											},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_303e68ec-dbb3-4d90-8a96-26e0be44f5f3",
												},
												Name: "Expanded Sub-Process 1",
											},
											Incoming: []string{"_ce63770e-9f4b-4bc0-a56c-88401c59f0bf"},
											Outgoing: []string{"_e6537f9d-e5ea-4abc-a6e8-add13a11b536"},
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
																	ID: "_ef372c4d-6c65-4ad2-bc22-5c25e5d0870e",
																},
																Name: "Start Event 2",
															},
															Outgoing: []string{"_c68194b4-3618-4655-b128-7eec67483c84"},
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
																	ID: "_b9343536-6490-4559-8365-71d5c4cbb7cb",
																},
																Name: "User Task 7 Standard Loop",
															},
															Incoming: []string{"_c68194b4-3618-4655-b128-7eec67483c84"},
															Outgoing: []string{"_87ffa0fa-1a2d-4149-bbe9-04e20bc1014b"},
														},
														LoopCharacteristicsElements: element.LoopCharacteristicsElements{
															StandardLoopCharacteristics: []element.StandardLoopCharacteristics{
																{},
															},
														},
													},
												},
												Implementation: "##WebService",
											},
										},
										EndEvents: []element.EndEvent{
											{
												ThrowEvent: element.ThrowEvent{
													Event: element.Event{
														FlowNode: element.FlowNode{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "_a9b9c08d-377a-49a8-a869-f82308702018",
																},
																Name: "End Event 2",
															},
															Incoming: []string{"_87ffa0fa-1a2d-4149-bbe9-04e20bc1014b"},
														},
													},
												},
											},
										},
										SequenceFlows: []element.SequenceFlow{
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_c68194b4-3618-4655-b128-7eec67483c84",
													},
												},
												SourceRef: "_ef372c4d-6c65-4ad2-bc22-5c25e5d0870e",
												TargetRef: "_b9343536-6490-4559-8365-71d5c4cbb7cb",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_87ffa0fa-1a2d-4149-bbe9-04e20bc1014b",
													},
												},
												SourceRef: "_b9343536-6490-4559-8365-71d5c4cbb7cb",
												TargetRef: "_a9b9c08d-377a-49a8-a869-f82308702018",
											},
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
														ID: "_0326fdf5-7c71-41d9-838c-ab141a1b1ed0",
													},
													Name: "Intermediate Event Signal Throw 1",
												},
												Incoming: []string{"_b8694c28-408b-4394-b476-41e2a5bcfd94"},
												Outgoing: []string{"_e0fbb051-0e07-43cc-a8bf-158492541c0f"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											SignalEventDefinitions: []element.SignalEventDefinition{{}},
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
													ID: "_a74c1d4d-db90-43ff-8920-139a300b39a5",
												},
												Name: "Call Activity calling a Global User Task",
											},
											Incoming: []string{"_be19c2da-316a-47f6-ad7b-eb6c82bf8609"},
											Outgoing: []string{"_ce63770e-9f4b-4bc0-a56c-88401c59f0bf"},
										},
									},
									CalledElement: "_d4afdad7-65a5-40c6-9fee-e13ba4c0bed4",
								},
							},
							DataObjectReferenes: []element.DataObjectReference{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_aa8c769a-276c-4589-b182-7c7bbd0a9e1e",
										},
										Name: "Data Object",
									},
									DataObjectRef: "DF1373638080458",
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_2c124731-3338-46e0-81b5-b435f34941c8",
										},
									},
									SourceRef: "_4e71bf73-1719-401e-a9a2-85dc89fc1150",
									TargetRef: "_d84e5824-7bb7-4057-9dba-6c8794f7948c",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_19a7094a-7bdf-4819-af86-be22d2cfdc49",
										},
									},
									SourceRef: "_d84e5824-7bb7-4057-9dba-6c8794f7948c",
									TargetRef: "_76ee26df-2c95-495b-9d9a-cb806aea6baf",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_e1f948a6-db76-4c35-ba36-e430eecb63a4",
										},
									},
									SourceRef: "_76ee26df-2c95-495b-9d9a-cb806aea6baf",
									TargetRef: "_0e87da16-736e-45b2-95e5-8f45940f3adf",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_dbae7e12-67c2-4256-8c50-5811c207ba55",
										},
									},
									SourceRef: "_149a6e1d-0385-4d0f-a90c-c2150a291a67",
									TargetRef: "_2a08c361-be51-437e-a86d-c62798c14e83",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_4474c77a-01bb-435e-929a-7eb71bd824a8",
										},
									},
									SourceRef: "_0e87da16-736e-45b2-95e5-8f45940f3adf",
									TargetRef: "_dec393e7-f182-4d31-b05f-e33ac3a5e35f",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_a9966baf-d9b9-4be1-a4d9-2906ab5add30",
										},
										Name: "Conditional Sequence Flow",
									},
									SourceRef: "_dec393e7-f182-4d31-b05f-e33ac3a5e35f",
									TargetRef: "_ac1fde31-c0cd-4a8a-9728-a5fb49602de7",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_be19c2da-316a-47f6-ad7b-eb6c82bf8609",
										},
										Name: "Default Sequence Flow 1",
									},
									SourceRef: "_dec393e7-f182-4d31-b05f-e33ac3a5e35f",
									TargetRef: "_a74c1d4d-db90-43ff-8920-139a300b39a5",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_b8694c28-408b-4394-b476-41e2a5bcfd94",
										},
									},
									SourceRef: "_ac1fde31-c0cd-4a8a-9728-a5fb49602de7",
									TargetRef: "_0326fdf5-7c71-41d9-838c-ab141a1b1ed0",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_ce63770e-9f4b-4bc0-a56c-88401c59f0bf",
										},
									},
									SourceRef: "_a74c1d4d-db90-43ff-8920-139a300b39a5",
									TargetRef: "_303e68ec-dbb3-4d90-8a96-26e0be44f5f3",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_e0fbb051-0e07-43cc-a8bf-158492541c0f",
										},
									},
									SourceRef: "_0326fdf5-7c71-41d9-838c-ab141a1b1ed0",
									TargetRef: "_149a6e1d-0385-4d0f-a90c-c2150a291a67",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_2752e07c-f2e4-4384-8721-70e4abea9765",
										},
									},
									SourceRef: "_397c783e-ad6a-4cf3-8266-9b41962c83bd",
									TargetRef: "_b67ba682-c8d6-465b-b538-c287db18d1be",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_e6537f9d-e5ea-4abc-a6e8-add13a11b536",
										},
									},
									SourceRef: "_303e68ec-dbb3-4d90-8a96-26e0be44f5f3",
									TargetRef: "_c9870992-6643-4094-acfd-d76e5e37941b",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_da47ac6a-4ea1-4bf6-ad3e-8cc60c0ea8a9",
										},
									},
									SourceRef: "_708d55c8-684a-4e3b-a69d-69c620cd0ac0",
									TargetRef: "_8f41085a-3302-4cfc-9231-23bb4b43c7a1",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_8321494b-9c3c-483a-b4d2-79e6254211c0",
										},
									},
									SourceRef: "_c9870992-6643-4094-acfd-d76e5e37941b",
									TargetRef: "_397c783e-ad6a-4cf3-8266-9b41962c83bd",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_54574511-c29e-4157-bb8f-8b26f15faaf5",
										},
									},
									SourceRef: "_732c0641-b12f-448b-b9f8-a68b355782e3",
									TargetRef: "_f3698a93-fe0b-4f49-bfa6-68d055936af3",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_b41b9c86-bc41-4b6e-ac32-70e1342e6128",
										},
									},
									SourceRef: "_2a08c361-be51-437e-a86d-c62798c14e83",
									TargetRef: "_397c783e-ad6a-4cf3-8266-9b41962c83bd",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_de249cbe-d431-4e1f-bc8c-a90aec438683",
										},
									},
									SourceRef: "_f3698a93-fe0b-4f49-bfa6-68d055936af3",
									TargetRef: "_397c783e-ad6a-4cf3-8266-9b41962c83bd",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_5f52bb2a-c49d-43ae-9253-7b89a8251b2b",
										},
									},
									SourceRef: "_8f41085a-3302-4cfc-9231-23bb4b43c7a1",
									TargetRef: "_c9cb2415-6a2e-49d6-84b9-27babcde4088",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_875031ae-f87c-45a3-ae60-ba0cb0ce0bc1",
										},
									},
									SourceRef: "_a01498ae-086c-4adc-9229-ec3135bc2bcf",
									TargetRef: "_5cc02d0f-c090-4e48-8da3-f32cbbca9565",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_e503a157-f034-4a8a-ae12-aa524a8ebeee",
										},
									},
									SourceRef: "_c9cb2415-6a2e-49d6-84b9-27babcde4088",
									TargetRef: "_a01498ae-086c-4adc-9229-ec3135bc2bcf",
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
									Text: "Annotation",
								},
							},
							Associations: []element.Association{
								{
									Artifact: element.Artifact{
										BaseElement: element.BaseElement{
											ID: "_5362a7ef-ce7e-4a91-9c38-66c07b1b5f49",
										},
									},
									AssociationDirection: element.AssociationDirectionNone,
									SourceRef:            "_303e68ec-dbb3-4d90-8a96-26e0be44f5f3",
									TargetRef:            "_4815ea6a-ede2-489b-8b37-2cdb2835b02c",
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
									ID: "ls_55bb31e8-9e62-48ea-8f0e-1a748c04bbf6",
								},
								Lanes: []element.Lane{
									{
										BaseElement: element.BaseElement{
											ID: "_4a6df7ac-26d8-4718-ac05-90af463d5e23",
										},
										Name: "Lane 1",
										FlowNodeRefs: []string{
											"_7e6ccf38-e740-4537-a439-a8e984d066de",
											"_fa90f891-fc07-463a-97c9-2ee0812351e1",
											"_1237e756-d53c-4591-a731-dafffbf0b3f9",
											"_73343358-7838-48a1-baf2-2b0266e3a55b",
											"_137281ee-758e-4c36-8942-74c5d807e1b3",
											"_15e9a3bf-53de-40d8-8364-7f534b175ff4",
											"_dbca671f-08b6-4b58-a614-62b98fa36be5",
											"_f2081fdb-3b8a-480b-9f61-fbf683e2018c",
											"_be29f267-9d56-46ef-8bbc-e13513b25fce",
											"_3a19ce2d-e30d-461c-aeaa-b2acb118f9a4",
											"_ba16239e-181e-4b9f-bc5b-0bb2ee973450",
											"_087d0602-ff51-491e-a021-d2e7c940dbd8",
											"_49e94b5f-ce21-4c2b-b78d-3cde5c09c15e",
											"_a38484e2-7bdb-48b1-b62e-139d51d6a147",
											"_c889aa27-e389-48eb-aa18-613ec53614e7",
											"_f27040d5-765c-493c-bbe7-9fb6ad04cbdc",
											"_0263ca9e-2ca0-4f4e-b7dd-86e15dcf2447",
											"_796ccbc5-ad88-465c-849a-87447a0283d3",
											"_05c6bc89-5265-435c-8a9e-533c44a6888b",
											"_511d95ed-38f9-473e-9466-525285a007f5",
											"_187063b6-107e-4ac9-bdb3-e8ce9d83763d",
											"_f07e4bd2-768d-42c6-a8d5-24d1c3bfa3cb",
											"_663e9963-9cf6-4032-9652-3a20f50dcda3",
											"_034907bf-d3d7-4629-818c-14c3e69d5bc6",
										},
									},
									{
										BaseElement: element.BaseElement{
											ID: "_3400f56a-4565-47d1-91db-0ba17b958cb2",
										},
										Name: "Lane 2",
										FlowNodeRefs: []string{
											"_d58753a7-d38b-49cd-914d-14e4cdaa4449",
											"_147b1900-7e7d-4b8f-b243-0155265f1c00",
											"_25beeb17-acc3-4cca-9590-f1cd2f353434",
											"_928cd158-ebe1-4c0a-9c3e-42e77d663aa2",
											"_0ab0e0ac-f88b-4402-9986-e95a72dfc8a3",
											"_242b8e6c-681c-438e-ab34-729255121eff",
											"_25a1f26b-d8b8-4b9b-b768-f14b8d104f9c",
											"_8476a0f7-36b7-4666-a3b2-c18efcc68a94",
											"_10ecbff1-cd15-4a5c-9aa5-6f2a35479416",
											"_4f5e6e50-d9d0-4f97-959a-d1b8e1e32788",
											"_cbebc7f2-9fb5-4fbf-a6dc-13140c784da7",
											"_df7727a0-f509-45eb-bb89-85753f439576",
											"_189118eb-65fc-43e8-8a34-a155b113914f",
											"_d92d850c-37ac-47ea-9344-dc986289de47",
											"_1215d072-524b-4724-99d7-a0a406435904",
											"_dfb273c6-0ad3-4030-9e72-638adf7ca75f",
											"_56ab3884-4b26-4398-be6f-5f0bb5f81be4",
											"_0e99d67a-a88a-4631-85cc-aa1f9cd8cc5e",
										},
									},
								},
							},
						},
						FlowElements: element.FlowElements{
							SubProcesses: []element.SubProcess{
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_7e6ccf38-e740-4537-a439-a8e984d066de",
												},
												Name: "Expanded Sub-Process 2",
											},
											Incoming: []string{"_3e8b97e7-d6a5-44dc-b5ca-6f2c39ba4911"},
											Outgoing: []string{"_70617827-824b-4797-b424-4179b8b6cbdd"},
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
																Name: "Start Event 5 None",
															},
															Outgoing: []string{"_2d1047ce-fdd5-4cb6-9f0c-0ee8d6d3044a"},
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
																	ID: "_6936f794-7bbb-4aa1-ae48-3a35bab4e2f4",
																},
																Name: "Service Task 22",
															},
															Incoming: []string{"_2d1047ce-fdd5-4cb6-9f0c-0ee8d6d3044a"},
															Outgoing: []string{"_062ae395-4aba-408b-ac64-4987752be95b"},
														},
														LoopCharacteristicsElements: element.LoopCharacteristicsElements{
															MultielementLoopCharacteristics: []element.MultiInstanceLoopCharacteristics{
																{},
															},
														},
													},
												},
												Implementation: "##WebService",
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
																Name: "End Event 8 None",
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
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_0263ca9e-2ca0-4f4e-b7dd-86e15dcf2447",
												},
												Name: "Collapsed Sub-Process 2",
											},
											Incoming: []string{"_e062c113-2da3-40f1-845e-6fd48ccc880f"},
											Outgoing: []string{"_8095da9c-0faa-47b9-85d4-2df24e021770"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_189118eb-65fc-43e8-8a34-a155b113914f",
												},
												Name: "Expanded Sub-Process 3",
											},
											Incoming: []string{"_40d118ea-6a9b-4210-a1f1-e093831e0df0", "_c9768243-8e1a-4b09-907b-3f27fa831a6e"},
											Outgoing: []string{"_a3fffd23-7ce0-4199-8918-c9df7a2c8158"},
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
																	ID: "_21976b84-4ddd-4ddc-a5a3-825335550796",
																},
																Name: "Start Event 7",
															},
															Outgoing: []string{"_1f05bcca-f418-4dc6-a78c-a9dde2fab53d"},
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
																	ID: "_e233b5e1-244d-422e-8886-4588b7566122",
																},
																Name: "Intermediate Event Signal Catch",
															},
															Incoming: []string{"_1f05bcca-f418-4dc6-a78c-a9dde2fab53d"},
															Outgoing: []string{"_03b1de69-7605-4fc2-9797-2309271c208c"},
														},
													},
													EventDefinitions: element.EventDefinitions{
														SignalEventDefinitions: []element.SignalEventDefinition{
															{},
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
																ID: "_6d90f706-17c9-4635-87c2-ccab31e9a32d",
															},
															Name: "Task 31",
														},
														Incoming: []string{"_03b1de69-7605-4fc2-9797-2309271c208c"},
														Outgoing: []string{"_955edc35-abcc-4cdb-a890-0d7ec15167ae"},
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
																ID: "_84918a6a-5e98-40c6-9155-c081bb5bdee8",
															},
															Name: "Exclusive Gateway 7",
														},
														Incoming: []string{"_955edc35-abcc-4cdb-a890-0d7ec15167ae"},
														Outgoing: []string{"_8021571a-7a77-426b-b824-545cc61f7334", "_5d70a293-03ea-4a73-bed0-65ad145796d6"},
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
																	ID: "_d60db966-c03d-4f0e-a5bd-945525fa0aaf",
																},
																Name: "End Event 13 Error",
															},
															Incoming: []string{"_5d70a293-03ea-4a73-bed0-65ad145796d6"},
														},
													},
													EventDefinitions: element.EventDefinitions{
														ErrorEventDefinitions: []element.ErrorEventDefinition{
															{},
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
																	ID: "_cd7b1449-fc16-4015-befe-cc9b5aa1df27",
																},
																Name: "End Event 12",
															},
															Incoming: []string{"_8021571a-7a77-426b-b824-545cc61f7334"},
														},
													},
												},
											},
										},
										SequenceFlows: []element.SequenceFlow{
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_1f05bcca-f418-4dc6-a78c-a9dde2fab53d",
													},
												},
												SourceRef: "_21976b84-4ddd-4ddc-a5a3-825335550796",
												TargetRef: "_e233b5e1-244d-422e-8886-4588b7566122",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_03b1de69-7605-4fc2-9797-2309271c208c",
													},
												},
												SourceRef: "_e233b5e1-244d-422e-8886-4588b7566122",
												TargetRef: "_6d90f706-17c9-4635-87c2-ccab31e9a32d",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_955edc35-abcc-4cdb-a890-0d7ec15167ae",
													},
												},
												SourceRef: "_6d90f706-17c9-4635-87c2-ccab31e9a32d",
												TargetRef: "_84918a6a-5e98-40c6-9155-c081bb5bdee8",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_8021571a-7a77-426b-b824-545cc61f7334",
													},
												},
												SourceRef: "_84918a6a-5e98-40c6-9155-c081bb5bdee8",
												TargetRef: "_cd7b1449-fc16-4015-befe-cc9b5aa1df27",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_5d70a293-03ea-4a73-bed0-65ad145796d6",
													},
												},
												SourceRef: "_84918a6a-5e98-40c6-9155-c081bb5bdee8",
												TargetRef: "_d60db966-c03d-4f0e-a5bd-945525fa0aaf",
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
														ID: "_5a6baa94-303a-4750-bde2-e1cd6edace37",
													},
													Name: "Boundary Intermediate Event Non-Interrupting Timer",
												},
												Outgoing: []string{"_e59dbf35-3f4e-4701-bc2a-75e32cbaa732"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											TimerEventDefinitions: []element.TimerEventDefinition{
												{TimeDate: element.ExpressionUnMarshal{
													ExpressionSubstitution: &element.Expression{},
												}},
											},
										},
									},
									CancelActivity: false,
									AttachedToRef:  "_7e6ccf38-e740-4537-a439-a8e984d066de",
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_3c56e6dc-bc87-4d98-b499-462c5b741c5a",
													},
													Name: "Boundary Intermediate Event Interrupting Error",
												},
												Outgoing: []string{"_708324bb-26d0-4358-a96f-a4fabc14069f"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											ErrorEventDefinitions: []element.ErrorEventDefinition{
												{},
											},
										},
									},
									AttachedToRef: "_7e6ccf38-e740-4537-a439-a8e984d066de",
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_68ca1f8b-5028-4079-9e35-619b529f4d71",
													},
													Name: "Boundary Intermediate Event Interrupting Conditional",
												},
												Outgoing: []string{"_78361e03-8ab9-4317-b8f6-2daa90f20f54"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											ConditionalEventDefinitions: []element.ConditionalEventDefinition{
												{
													Condition: element.ExpressionUnMarshal{
														ExpressionSubstitution: &element.Expression{},
													},
												},
											},
										},
									},
									AttachedToRef: "_fa90f891-fc07-463a-97c9-2ee0812351e1",
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_45ceee21-0f15-4bf8-87a9-b3f808173e61",
													},
													Name: "Boundary Intermediate Event Non-Interrupting Escalation",
												},
												Outgoing: []string{"_07b92d12-375f-41c5-b0a9-8961ac773afe"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											EscalationEventDefinitions: []element.EscalationEventDefinition{
												{},
											},
										},
									},
									CancelActivity: false,
									AttachedToRef:  "_1237e756-d53c-4591-a731-dafffbf0b3f9",
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_e454657a-0173-41a4-a4c7-d16ec224f2e1",
													},
													Name: "Boundary Intermediate Event Non-Interrupting Signal",
												},
												Outgoing: []string{"_a639a36d-a1df-43b1-8ed1-e06afa899733"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											SignalEventDefinitions: []element.SignalEventDefinition{
												{},
											},
										},
									},
									CancelActivity: false,
									AttachedToRef:  "_73343358-7838-48a1-baf2-2b0266e3a55b",
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_79341f54-50d4-4c60-85f3-fe8839a7554b",
													},
													Name: "Boundary Intermediate Event Interrupting Timer",
												},
												Outgoing: []string{"_4c3f3102-d31a-4a71-a3d0-b65cb61a94ea"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											TimerEventDefinitions: []element.TimerEventDefinition{
												{TimeDate: element.ExpressionUnMarshal{
													ExpressionSubstitution: &element.Expression{},
												}},
											},
										},
									},
									AttachedToRef: "_137281ee-758e-4c36-8942-74c5d807e1b3",
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_e369fd30-1a71-4d0e-b4d7-2174dd5ba388",
													},
													Name: "Boundary Intermediate Event Non-Interrupting Message",
												},
												Outgoing: []string{"_504b3bd1-54f6-4f20-8244-32f4ec639248"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{},
											},
										},
									},
									CancelActivity: false,
									AttachedToRef:  "_15e9a3bf-53de-40d8-8364-7f534b175ff4",
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_209105e0-96fc-4278-8451-3b2a1dd18ec9",
													},
												},
												Outgoing: []string{"_aa88dfff-5242-4699-b722-a593e54537ce"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											SignalEventDefinitions: []element.SignalEventDefinition{
												{},
											},
										},
									},
									AttachedToRef: "_d58753a7-d38b-49cd-914d-14e4cdaa4449",
								},
							},
							ServiceTasks: []element.ServiceTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_fa90f891-fc07-463a-97c9-2ee0812351e1",
													},
													Name: "Service Task 15",
												},
												Incoming: []string{"_c1931975-c1c3-499e-8a57-89b61affba3a"},
												Outgoing: []string{"_c8c10b02-ea53-4a0c-b605-b3252d5560cd"},
											},
										},
									},
									Implementation: "##WebService",
								},
							},
							CallActivities: []element.CallActivity{
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_1237e756-d53c-4591-a731-dafffbf0b3f9",
												},
												Name: "Collapsed Call Activity",
											},
											Incoming: []string{"_837a8629-1540-473c-b41f-7e13ad80c052"},
											Outgoing: []string{"_b77e4b02-22b9-4a57-b0f9-4248410d0675"},
										},
									},
									CalledElement: "WFP-0-",
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_ba16239e-181e-4b9f-bc5b-0bb2ee973450",
												},
												Name: "Expanded Call Activity",
											},
											Incoming: []string{"_5106fe5e-184d-4069-8c1c-54f81fd577a9"},
											Outgoing: []string{"_c1931975-c1c3-499e-8a57-89b61affba3a"},
										},
									},
									CalledElement: "Process_ba16239e-181e-4b9f-bc5b-0bb2ee973450",
								},
							},
							Tasks: []element.Task{
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_73343358-7838-48a1-baf2-2b0266e3a55b",
												},
												Name: "Task 23",
											},
											Incoming: []string{"_70617827-824b-4797-b424-4179b8b6cbdd", "_2257c019-bc17-4ba9-9e07-a9de7913ada3"},
											Outgoing: []string{"_02f751bb-2e42-489a-bef1-cc5360d5c3d8"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_137281ee-758e-4c36-8942-74c5d807e1b3",
												},
												Name: "Task 21",
											},
											Incoming: []string{"_2a32599c-d1f4-4f2c-bf65-0f0e4f6ac87f"},
											Outgoing: []string{"_3e8b97e7-d6a5-44dc-b5ca-6f2c39ba4911"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_15e9a3bf-53de-40d8-8364-7f534b175ff4",
												},
												Name: "Task 17",
											},
											Incoming: []string{"_b77e4b02-22b9-4a57-b0f9-4248410d0675"},
											Outgoing: []string{"_ce44f766-7c92-43e2-93cf-1e2ba7110c4a"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_dbca671f-08b6-4b58-a614-62b98fa36be5",
												},
												Name: "Task 18",
											},
											Incoming: []string{"_8b98cde1-aec2-46e8-8be2-d9aa244fcda6", "_07b92d12-375f-41c5-b0a9-8961ac773afe"},
											Outgoing: []string{"_bd5730fa-544e-4a99-a238-57171787d52c"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_511d95ed-38f9-473e-9466-525285a007f5",
												},
												Name: "Task 11",
											},
											Incoming: []string{"_a63c8cd6-eee8-4fbe-be5e-f6980b180b52"},
											Outgoing: []string{"_168f4ce9-ccf7-4833-b38b-0140ff601d40"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_187063b6-107e-4ac9-bdb3-e8ce9d83763d",
												},
												Name: "Task 19",
											},
											Incoming: []string{"_504b3bd1-54f6-4f20-8244-32f4ec639248"},
											Outgoing: []string{"_b54d5ab0-66d8-4595-8a89-9d3e255f75ba"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_663e9963-9cf6-4032-9652-3a20f50dcda3",
												},
												Name: "Task 24",
											},
											Incoming: []string{"_a639a36d-a1df-43b1-8ed1-e06afa899733"},
											Outgoing: []string{"_9e4cd50c-cd6b-4523-8dbe-ebece9b823cb"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_d58753a7-d38b-49cd-914d-14e4cdaa4449",
												},
												Name: "Task 32",
											},
											Incoming: []string{"_a3fffd23-7ce0-4199-8918-c9df7a2c8158"},
											Outgoing: []string{"_b9a903b5-4525-42d1-ae5b-24f26d774556"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_928cd158-ebe1-4c0a-9c3e-42e77d663aa2",
												},
												Name: "Task 28",
											},
											Incoming: []string{"_831dbaee-1434-45fc-9c7e-bedd44ad9ea7"},
											Outgoing: []string{"_202c373c-f243-413d-904e-9132a0c0e923"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_0ab0e0ac-f88b-4402-9986-e95a72dfc8a3",
												},
												Name: "Task 33",
											},
											Incoming: []string{"_aa88dfff-5242-4699-b722-a593e54537ce"},
											Outgoing: []string{"_219e1df4-5bfe-4485-a4fd-dcfeeaa7c3d6"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_242b8e6c-681c-438e-ab34-729255121eff",
												},
												Name: "Task 25",
											},
											Incoming: []string{"_7c690c39-4274-40c5-a0d2-b503acabbf93"},
											Outgoing: []string{"_fdd08093-e5b4-4e9e-8088-26887892078a"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_25a1f26b-d8b8-4b9b-b768-f14b8d104f9c",
												},
												Name: "Task 30",
											},
											Incoming: []string{"_e59dbf35-3f4e-4701-bc2a-75e32cbaa732"},
											Outgoing: []string{"_37312691-7037-44d1-8696-fc6f2021e7a6"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_cbebc7f2-9fb5-4fbf-a6dc-13140c784da7",
												},
												Name: "Task 26",
											},
											Incoming: []string{"_022aa2b9-f472-49be-b00b-2eec7f075acf"},
											Outgoing: []string{"_be71b068-7fa6-4f76-b0a3-8760f4a2911a"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_56ab3884-4b26-4398-be6f-5f0bb5f81be4",
												},
												Name: "Task 29",
											},
											Incoming: []string{"_708324bb-26d0-4358-a96f-a4fabc14069f"},
											Outgoing: []string{"_e03319e1-0dff-4337-9443-b053651000a0"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_0e99d67a-a88a-4631-85cc-aa1f9cd8cc5e",
												},
												Name: "Task 27",
											},
											Incoming: []string{"_4c3f3102-d31a-4a71-a3d0-b65cb61a94ea"},
											Outgoing: []string{"_00140039-2e5d-4c58-9197-1f8512b54c99"},
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
														ID: "_f2081fdb-3b8a-480b-9f61-fbf683e2018c",
													},
													Name: "Intermediate Event Message Catch",
												},
												Incoming: []string{"_7cae752c-c2bd-438b-8d24-48196805a4e8"},
												Outgoing: []string{"_e062c113-2da3-40f1-845e-6fd48ccc880f"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{},
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
														ID: "_05c6bc89-5265-435c-8a9e-533c44a6888b",
													},
													Name: "Intermediate Event Message Catch 2",
												},
												Incoming: []string{"_ab34472d-95a4-459c-a13b-5ed8b8b75eca"},
												Outgoing: []string{"_2a32599c-d1f4-4f2c-bf65-0f0e4f6ac87f"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{},
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
														ID: "_034907bf-d3d7-4629-818c-14c3e69d5bc6",
													},
													Name: "Intermediate Event Timer Catch",
												},
												Incoming: []string{"_5853836e-d7ca-45e2-852a-7db8c3c642bb"},
												Outgoing: []string{"_5106fe5e-184d-4069-8c1c-54f81fd577a9"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											TimerEventDefinitions: []element.TimerEventDefinition{
												{TimeDate: element.ExpressionUnMarshal{
													ExpressionSubstitution: &element.Expression{},
												}},
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
														ID: "_4f5e6e50-d9d0-4f97-959a-d1b8e1e32788",
													},
													Name: "Intermediate Event Link",
												},
												Outgoing: []string{"_c9768243-8e1a-4b09-907b-3f27fa831a6e"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											LinkEventDefinitions: []element.LinkEventDefinition{
												{Name: "Intermediate Event Link"},
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
													ID: "_be29f267-9d56-46ef-8bbc-e13513b25fce",
												},
												Name: "Event Base Gateway 3",
											},
											Incoming: []string{"_168f4ce9-ccf7-4833-b38b-0140ff601d40"},
											Outgoing: []string{"_ab34472d-95a4-459c-a13b-5ed8b8b75eca", "_5853836e-d7ca-45e2-852a-7db8c3c642bb", "_7cae752c-c2bd-438b-8d24-48196805a4e8"},
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
														ID: "_3a19ce2d-e30d-461c-aeaa-b2acb118f9a4",
													},
												},
												Incoming: []string{"_bd5730fa-544e-4a99-a238-57171787d52c"},
												Outgoing: []string{"_2257c019-bc17-4ba9-9e07-a9de7913ada3"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											EscalationEventDefinitions: []element.EscalationEventDefinition{
												{},
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
														ID: "_f27040d5-765c-493c-bbe7-9fb6ad04cbdc",
													},
													Name: "Intermediate Event Link",
												},
												Incoming: []string{"_78361e03-8ab9-4317-b8f6-2daa90f20f54"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											LinkEventDefinitions: []element.LinkEventDefinition{
												{Name: "Intermediate Event Link"},
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
														ID: "_796ccbc5-ad88-465c-849a-87447a0283d3",
													},
													Name: "Intermediate Event Message Throw",
												},
												Incoming: []string{"_670ceb69-cd3a-46e8-96a0-a520a8fc589b"},
												Outgoing: []string{"_837a8629-1540-473c-b41f-7e13ad80c052"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{},
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
														ID: "_8476a0f7-36b7-4666-a3b2-c18efcc68a94",
													},
													Name: "Intermediate Event Signal Throw 2",
												},
												Incoming: []string{"_be71b068-7fa6-4f76-b0a3-8760f4a2911a"},
												Outgoing: []string{"_f61be5ab-acb2-4348-a9a0-bdfdde0c42ad"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											SignalEventDefinitions: []element.SignalEventDefinition{
												{},
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
														ID: "_087d0602-ff51-491e-a021-d2e7c940dbd8",
													},
													Name: "End Event 7 None",
												},
												Incoming: []string{
													"_ce44f766-7c92-43e2-93cf-1e2ba7110c4a",
													"_b54d5ab0-66d8-4595-8a89-9d3e255f75ba",
													"_02f751bb-2e42-489a-bef1-cc5360d5c3d8",
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
														ID: "_c889aa27-e389-48eb-aa18-613ec53614e7",
													},
													Name: "End Event 6\nMessage",
												},
												Incoming: []string{"_cc2f1bcd-2a97-48c7-94ae-0899e56402ea"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{MessageRef: "Message_1373638080954"},
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
														ID: "_147b1900-7e7d-4b8f-b243-0155265f1c00",
													},
													Name: "End Event 14",
												},
												Incoming: []string{"_219e1df4-5bfe-4485-a4fd-dcfeeaa7c3d6"},
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
														ID: "_d92d850c-37ac-47ea-9344-dc986289de47",
													},
													Name: "End Event 10",
												},
												Incoming: []string{"_e03319e1-0dff-4337-9443-b053651000a0"},
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
														ID: "_dfb273c6-0ad3-4030-9e72-638adf7ca75f",
													},
													Name: "End Event 11 Escalation",
												},
												Incoming: []string{"_9e4cd50c-cd6b-4523-8dbe-ebece9b823cb", "_37312691-7037-44d1-8696-fc6f2021e7a6", "_0dbfad2a-7b07-4b7b-87b4-3819d28f175d"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											EscalationEventDefinitions: []element.EscalationEventDefinition{
												{},
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
													ID: "_49e94b5f-ce21-4c2b-b78d-3cde5c09c15e",
												},
												Name: "Exclusive Gateway 4",
											},
											Incoming: []string{"_8095da9c-0faa-47b9-85d4-2df24e021770"},
											Outgoing: []string{"_670ceb69-cd3a-46e8-96a0-a520a8fc589b", "_8b98cde1-aec2-46e8-8be2-d9aa244fcda6"},
										},
									},
									Default: "_670ceb69-cd3a-46e8-96a0-a520a8fc589b",
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
													Name: "Start Event 2 Message",
												},
												Outgoing: []string{"_a63c8cd6-eee8-4fbe-be5e-f6980b180b52"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													MessageRef: "Message_1373638080955",
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
														ID: "_25beeb17-acc3-4cca-9590-f1cd2f353434",
													},
													Name: "Start Event 6 Signal",
												},
												Outgoing: []string{"_7c690c39-4274-40c5-a0d2-b503acabbf93"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											SignalEventDefinitions: []element.SignalEventDefinition{
												{},
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
														ID: "_f07e4bd2-768d-42c6-a8d5-24d1c3bfa3cb",
													},
													Name: "Receive Task 16",
												},
												Incoming: []string{"_c8c10b02-ea53-4a0c-b605-b3252d5560cd"},
												Outgoing: []string{"_cc2f1bcd-2a97-48c7-94ae-0899e56402ea"},
											},
										},
									},
									Implementation: "##WebService",
								},
							},
							InclusiveGatewaies: []element.InclusiveGateway{
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_10ecbff1-cd15-4a5c-9aa5-6f2a35479416",
												},
												Name: "Inclusive Gateway 6",
											},
											Incoming: []string{"_00140039-2e5d-4c58-9197-1f8512b54c99", "_f61be5ab-acb2-4348-a9a0-bdfdde0c42ad"},
											Outgoing: []string{"_831dbaee-1434-45fc-9c7e-bedd44ad9ea7"},
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
													ID: "_df7727a0-f509-45eb-bb89-85753f439576",
												},
												Name: "Parallel Gateway 7",
											},
											Incoming: []string{"_202c373c-f243-413d-904e-9132a0c0e923", "_b9a903b5-4525-42d1-ae5b-24f26d774556"},
											Outgoing: []string{"_0dbfad2a-7b07-4b7b-87b4-3819d28f175d"},
										},
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_1215d072-524b-4724-99d7-a0a406435904",
												},
												Name: "Parallel Gateway 5",
											},
											Incoming: []string{"_fdd08093-e5b4-4e9e-8088-26887892078a"},
											Outgoing: []string{"_022aa2b9-f472-49be-b00b-2eec7f075acf", "_40d118ea-6a9b-4210-a1f1-e093831e0df0"},
										},
									},
								},
							},
							DataStoreReferences: []element.DataStoreReference{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_b9385abf-d293-40b7-848b-8add4db48415",
										},
										Name: "Data Store Reference",
									},
									DataStoreRef: "DS1373638079677",
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
									TargetRef: "_511d95ed-38f9-473e-9466-525285a007f5",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_ab34472d-95a4-459c-a13b-5ed8b8b75eca",
										},
									},
									SourceRef: "_be29f267-9d56-46ef-8bbc-e13513b25fce",
									TargetRef: "_05c6bc89-5265-435c-8a9e-533c44a6888b",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_5853836e-d7ca-45e2-852a-7db8c3c642bb",
										},
									},
									SourceRef: "_be29f267-9d56-46ef-8bbc-e13513b25fce",
									TargetRef: "_034907bf-d3d7-4629-818c-14c3e69d5bc6",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_168f4ce9-ccf7-4833-b38b-0140ff601d40",
										},
									},
									SourceRef: "_511d95ed-38f9-473e-9466-525285a007f5",
									TargetRef: "_be29f267-9d56-46ef-8bbc-e13513b25fce",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_7cae752c-c2bd-438b-8d24-48196805a4e8",
										},
									},
									SourceRef: "_be29f267-9d56-46ef-8bbc-e13513b25fce",
									TargetRef: "_f2081fdb-3b8a-480b-9f61-fbf683e2018c",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_5106fe5e-184d-4069-8c1c-54f81fd577a9",
										},
									},
									SourceRef: "_034907bf-d3d7-4629-818c-14c3e69d5bc6",
									TargetRef: "_ba16239e-181e-4b9f-bc5b-0bb2ee973450",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_e062c113-2da3-40f1-845e-6fd48ccc880f",
										},
									},
									SourceRef: "_f2081fdb-3b8a-480b-9f61-fbf683e2018c",
									TargetRef: "_0263ca9e-2ca0-4f4e-b7dd-86e15dcf2447",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_2a32599c-d1f4-4f2c-bf65-0f0e4f6ac87f",
										},
									},
									SourceRef: "_05c6bc89-5265-435c-8a9e-533c44a6888b",
									TargetRef: "_137281ee-758e-4c36-8942-74c5d807e1b3",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_8095da9c-0faa-47b9-85d4-2df24e021770",
										},
									},
									SourceRef: "_0263ca9e-2ca0-4f4e-b7dd-86e15dcf2447",
									TargetRef: "_49e94b5f-ce21-4c2b-b78d-3cde5c09c15e",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_670ceb69-cd3a-46e8-96a0-a520a8fc589b",
										},
										Name: "Default Sequence Flow 2",
									},
									SourceRef: "_49e94b5f-ce21-4c2b-b78d-3cde5c09c15e",
									TargetRef: "_796ccbc5-ad88-465c-849a-87447a0283d3",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_8b98cde1-aec2-46e8-8be2-d9aa244fcda6",
										},
									},
									SourceRef: "_49e94b5f-ce21-4c2b-b78d-3cde5c09c15e",
									TargetRef: "_dbca671f-08b6-4b58-a614-62b98fa36be5",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_837a8629-1540-473c-b41f-7e13ad80c052",
										},
									},
									SourceRef: "_796ccbc5-ad88-465c-849a-87447a0283d3",
									TargetRef: "_1237e756-d53c-4591-a731-dafffbf0b3f9",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_4c3f3102-d31a-4a71-a3d0-b65cb61a94ea",
										},
									},
									SourceRef: "_79341f54-50d4-4c60-85f3-fe8839a7554b",
									TargetRef: "_0e99d67a-a88a-4631-85cc-aa1f9cd8cc5e",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_07b92d12-375f-41c5-b0a9-8961ac773afe",
										},
									},
									SourceRef: "_45ceee21-0f15-4bf8-87a9-b3f808173e61",
									TargetRef: "_dbca671f-08b6-4b58-a614-62b98fa36be5",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_c1931975-c1c3-499e-8a57-89b61affba3a",
										},
									},
									SourceRef: "_ba16239e-181e-4b9f-bc5b-0bb2ee973450",
									TargetRef: "_fa90f891-fc07-463a-97c9-2ee0812351e1",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_7c690c39-4274-40c5-a0d2-b503acabbf93",
										},
									},
									SourceRef: "_25beeb17-acc3-4cca-9590-f1cd2f353434",
									TargetRef: "_242b8e6c-681c-438e-ab34-729255121eff",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_fdd08093-e5b4-4e9e-8088-26887892078a",
										},
									},
									SourceRef: "_242b8e6c-681c-438e-ab34-729255121eff",
									TargetRef: "_1215d072-524b-4724-99d7-a0a406435904",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_022aa2b9-f472-49be-b00b-2eec7f075acf",
										},
									},
									SourceRef: "_1215d072-524b-4724-99d7-a0a406435904",
									TargetRef: "_cbebc7f2-9fb5-4fbf-a6dc-13140c784da7",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_40d118ea-6a9b-4210-a1f1-e093831e0df0",
										},
									},
									SourceRef: "_1215d072-524b-4724-99d7-a0a406435904",
									TargetRef: "_189118eb-65fc-43e8-8a34-a155b113914f",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_c8c10b02-ea53-4a0c-b605-b3252d5560cd",
										},
									},
									SourceRef: "_fa90f891-fc07-463a-97c9-2ee0812351e1",
									TargetRef: "_f07e4bd2-768d-42c6-a8d5-24d1c3bfa3cb",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_78361e03-8ab9-4317-b8f6-2daa90f20f54",
										},
									},
									SourceRef: "_68ca1f8b-5028-4079-9e35-619b529f4d71",
									TargetRef: "_f27040d5-765c-493c-bbe7-9fb6ad04cbdc",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_cc2f1bcd-2a97-48c7-94ae-0899e56402ea",
										},
									},
									SourceRef: "_f07e4bd2-768d-42c6-a8d5-24d1c3bfa3cb",
									TargetRef: "_c889aa27-e389-48eb-aa18-613ec53614e7",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_c9768243-8e1a-4b09-907b-3f27fa831a6e",
										},
									},
									SourceRef: "_4f5e6e50-d9d0-4f97-959a-d1b8e1e32788",
									TargetRef: "_189118eb-65fc-43e8-8a34-a155b113914f",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_be71b068-7fa6-4f76-b0a3-8760f4a2911a",
										},
									},
									SourceRef: "_cbebc7f2-9fb5-4fbf-a6dc-13140c784da7",
									TargetRef: "_8476a0f7-36b7-4666-a3b2-c18efcc68a94",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_3e8b97e7-d6a5-44dc-b5ca-6f2c39ba4911",
										},
									},
									SourceRef: "_137281ee-758e-4c36-8942-74c5d807e1b3",
									TargetRef: "_7e6ccf38-e740-4537-a439-a8e984d066de",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_00140039-2e5d-4c58-9197-1f8512b54c99",
										},
									},
									SourceRef: "_0e99d67a-a88a-4631-85cc-aa1f9cd8cc5e",
									TargetRef: "_10ecbff1-cd15-4a5c-9aa5-6f2a35479416",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_f61be5ab-acb2-4348-a9a0-bdfdde0c42ad",
										},
									},
									SourceRef: "_8476a0f7-36b7-4666-a3b2-c18efcc68a94",
									TargetRef: "_10ecbff1-cd15-4a5c-9aa5-6f2a35479416",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_bd5730fa-544e-4a99-a238-57171787d52c",
										},
									},
									SourceRef: "_dbca671f-08b6-4b58-a614-62b98fa36be5",
									TargetRef: "_3a19ce2d-e30d-461c-aeaa-b2acb118f9a4",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_b77e4b02-22b9-4a57-b0f9-4248410d0675",
										},
									},
									SourceRef: "_1237e756-d53c-4591-a731-dafffbf0b3f9",
									TargetRef: "_15e9a3bf-53de-40d8-8364-7f534b175ff4",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_ce44f766-7c92-43e2-93cf-1e2ba7110c4a",
										},
									},
									SourceRef: "_15e9a3bf-53de-40d8-8364-7f534b175ff4",
									TargetRef: "_087d0602-ff51-491e-a021-d2e7c940dbd8",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_70617827-824b-4797-b424-4179b8b6cbdd",
										},
									},
									SourceRef: "_7e6ccf38-e740-4537-a439-a8e984d066de",
									TargetRef: "_73343358-7838-48a1-baf2-2b0266e3a55b",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_b54d5ab0-66d8-4595-8a89-9d3e255f75ba",
										},
									},
									SourceRef: "_187063b6-107e-4ac9-bdb3-e8ce9d83763d",
									TargetRef: "_087d0602-ff51-491e-a021-d2e7c940dbd8",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_504b3bd1-54f6-4f20-8244-32f4ec639248",
										},
									},
									SourceRef: "_e369fd30-1a71-4d0e-b4d7-2174dd5ba388",
									TargetRef: "_187063b6-107e-4ac9-bdb3-e8ce9d83763d",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_708324bb-26d0-4358-a96f-a4fabc14069f",
										},
									},
									SourceRef: "_3c56e6dc-bc87-4d98-b499-462c5b741c5a",
									TargetRef: "_56ab3884-4b26-4398-be6f-5f0bb5f81be4",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_e03319e1-0dff-4337-9443-b053651000a0",
										},
									},
									SourceRef: "_56ab3884-4b26-4398-be6f-5f0bb5f81be4",
									TargetRef: "_d92d850c-37ac-47ea-9344-dc986289de47",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_e59dbf35-3f4e-4701-bc2a-75e32cbaa732",
										},
									},
									SourceRef: "_5a6baa94-303a-4750-bde2-e1cd6edace37",
									TargetRef: "_25a1f26b-d8b8-4b9b-b768-f14b8d104f9c",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_2257c019-bc17-4ba9-9e07-a9de7913ada3",
										},
									},
									SourceRef: "_3a19ce2d-e30d-461c-aeaa-b2acb118f9a4",
									TargetRef: "_73343358-7838-48a1-baf2-2b0266e3a55b",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_02f751bb-2e42-489a-bef1-cc5360d5c3d8",
										},
									},
									SourceRef: "_73343358-7838-48a1-baf2-2b0266e3a55b",
									TargetRef: "_087d0602-ff51-491e-a021-d2e7c940dbd8",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_a639a36d-a1df-43b1-8ed1-e06afa899733",
										},
									},
									SourceRef: "_e454657a-0173-41a4-a4c7-d16ec224f2e1",
									TargetRef: "_663e9963-9cf6-4032-9652-3a20f50dcda3",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_9e4cd50c-cd6b-4523-8dbe-ebece9b823cb",
										},
									},
									SourceRef: "_663e9963-9cf6-4032-9652-3a20f50dcda3",
									TargetRef: "_dfb273c6-0ad3-4030-9e72-638adf7ca75f",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_a3fffd23-7ce0-4199-8918-c9df7a2c8158",
										},
									},
									SourceRef: "_189118eb-65fc-43e8-8a34-a155b113914f",
									TargetRef: "_d58753a7-d38b-49cd-914d-14e4cdaa4449",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_37312691-7037-44d1-8696-fc6f2021e7a6",
										},
									},
									SourceRef: "_25a1f26b-d8b8-4b9b-b768-f14b8d104f9c",
									TargetRef: "_dfb273c6-0ad3-4030-9e72-638adf7ca75f",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_831dbaee-1434-45fc-9c7e-bedd44ad9ea7",
										},
									},
									SourceRef: "_10ecbff1-cd15-4a5c-9aa5-6f2a35479416",
									TargetRef: "_928cd158-ebe1-4c0a-9c3e-42e77d663aa2",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_202c373c-f243-413d-904e-9132a0c0e923",
										},
									},
									SourceRef: "_928cd158-ebe1-4c0a-9c3e-42e77d663aa2",
									TargetRef: "_df7727a0-f509-45eb-bb89-85753f439576",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_0dbfad2a-7b07-4b7b-87b4-3819d28f175d",
										},
									},
									SourceRef: "_df7727a0-f509-45eb-bb89-85753f439576",
									TargetRef: "_dfb273c6-0ad3-4030-9e72-638adf7ca75f",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_aa88dfff-5242-4699-b722-a593e54537ce",
										},
									},
									SourceRef: "_209105e0-96fc-4278-8451-3b2a1dd18ec9",
									TargetRef: "_0ab0e0ac-f88b-4402-9986-e95a72dfc8a3",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_219e1df4-5bfe-4485-a4fd-dcfeeaa7c3d6",
										},
									},
									SourceRef: "_0ab0e0ac-f88b-4402-9986-e95a72dfc8a3",
									TargetRef: "_147b1900-7e7d-4b8f-b243-0155265f1c00",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_b9a903b5-4525-42d1-ae5b-24f26d774556",
										},
									},
									SourceRef: "_d58753a7-d38b-49cd-914d-14e4cdaa4449",
									TargetRef: "_df7727a0-f509-45eb-bb89-85753f439576",
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
														ID: "_820dcc70-45ac-4a1e-88ae-f1b4ff925ef6",
													},
													Name: "Start Event 8",
												},
												Outgoing: []string{"_1c5e547a-2391-4133-8199-850cdc024971"},
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
													ID: "_13fbe8ab-af64-4b54-8efb-4c91dd6c6c18",
												},
												Name: "Task 34",
											},
											Incoming: []string{"_1c5e547a-2391-4133-8199-850cdc024971"},
											Outgoing: []string{"_af94c58e-db10-449f-978d-03e3b375b5a5"},
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
														ID: "_3cec2a74-8a45-4ef3-a196-690ba64f1b2b",
													},
													Name: "End Event 15",
												},
												Incoming: []string{"_af94c58e-db10-449f-978d-03e3b375b5a5"},
											},
										},
									},
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_1c5e547a-2391-4133-8199-850cdc024971",
										},
									},
									SourceRef: "_820dcc70-45ac-4a1e-88ae-f1b4ff925ef6",
									TargetRef: "_13fbe8ab-af64-4b54-8efb-4c91dd6c6c18",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_af94c58e-db10-449f-978d-03e3b375b5a5",
										},
									},
									SourceRef: "_13fbe8ab-af64-4b54-8efb-4c91dd6c6c18",
									TargetRef: "_3cec2a74-8a45-4ef3-a196-690ba64f1b2b",
								},
							},
						},
					},
				},
				Messages: []element.Message{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "Message_1373638080954",
							},
						},
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "Message_1373638080955",
							},
						},
					},
				},
				Categories: []element.Category{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "Cat1373638080956",
							},
						},
						CategoryValues: []element.CategoryValue{
							{
								BaseElement: element.BaseElement{
									ID: "Value_Cat1373638080956",
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
								ID: "C1373638080953",
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
									ID: "_55bb31e8-9e62-48ea-8f0e-1a748c04bbf6",
								},
								Name:       "Pool",
								ProcessRef: "WFP-6-2",
							},
						},
						MessageFlows: []element.MessageFlow{
							{
								BaseElement: element.BaseElement{
									ID: "_9428f666-fc8a-41be-8a77-9b280e14e7ae",
								},
								Name:       "Message Flow 2",
								SourceRef:  "_c889aa27-e389-48eb-aa18-613ec53614e7",
								TargetRef:  "_a01498ae-086c-4adc-9229-ec3135bc2bcf",
								MessageRef: "Message_1373638080954",
							},
							{
								BaseElement: element.BaseElement{
									ID: "_09e7cb23-4a1b-4165-b93a-cf635c223ee5",
								},
								Name:       "Message Flow 1",
								SourceRef:  "_76ee26df-2c95-495b-9d9a-cb806aea6baf",
								TargetRef:  "_a38484e2-7bdb-48b1-b62e-139d51d6a147",
								MessageRef: "Message_1373638080955",
							},
						},
						Artifacts: element.Artifacts{
							Groups: []element.Group{
								{
									Artifact: element.Artifact{
										BaseElement: element.BaseElement{
											ID: "_48d300c1-487a-409b-a04a-b195e222ef90",
										},
									},
									CategoryValueRef: "Value_Cat1373638080956",
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
