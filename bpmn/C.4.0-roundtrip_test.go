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

func TestC_4_0_roundtrip(t *testing.T) {
	// create test use ./test/C.4.0-roundtrip.bpmn
	path := "./test/C.4.0-roundtrip.bpmn"
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
			ID:              "_e7ccc12c-2987-48ae-92aa-08112101d9c9",
			Name:            "C.4.0",
			TargetNamespace: "http://www.trisotech.com/definitions/_e7ccc12c-2987-48ae-92aa-08112101d9c9",
			Exporter:        "BPMN Modeler",
			ExporterVersion: "6.2.0",
			RootElemnts: element.RootElemnts{
				ItemDefinitions: []element.ItemDefinition{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_triso-default-bpmnItemDefinition-string_id",
							},
						},
						StructureRef: "string",
					},
				},
				Signals: []element.Signal{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_bf41cec3-da5e-4fc7-87f9-7520a720192b",
							},
						},
						Name: "New employee hired",
					},
				},
				Collaborations: []element.Collaboration{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_085241a5-fb5c-44d3-8844-5366f865e353",
							},
						},
						Participants: []element.Participant{
							{
								BaseElement: element.BaseElement{
									ID: "_06745dd2-cbc6-4742-b11c-69281e5dcbdb",
								},
								Name:       "Money Bank",
								ProcessRef: "_42cba3a9-a8ab-40b5-b9a4-2e8f32be364e",
							},
						},
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_674b5898-d454-4114-a8eb-66c790366e94",
							},
						},
						Participants: []element.Participant{
							{
								BaseElement: element.BaseElement{
									ID: "_b5172765-a1c1-4c6e-8e57-9eabaf8cecdf",
								},
								Name:       "IT",
								ProcessRef: "_f0035388-f829-470c-b82b-0b15c3da3399",
							},
						},
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_bbdf487b-c3a3-4bfa-8ba3-79be59f08203",
							},
						},
						Participants: []element.Participant{
							{
								BaseElement: element.BaseElement{
									ID: "_1422af8d-518b-49f1-9563-9b6939de2818",
								},
								Name:       "Payroll",
								ProcessRef: "_da743a6f-d9e5-4fcf-8a96-d2fd5cfb73d4",
							},
						},
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_8d4bb408-4175-4fbb-8131-18c32b7aa022",
							},
						},
						Participants: []element.Participant{
							{
								BaseElement: element.BaseElement{
									ID: "_fab40827-7808-470a-93b6-4a3318bf0c0e",
								},
								Name:       "Facilities",
								ProcessRef: "_3486bf55-0a7f-4ff1-be15-1555669f58ad",
							},
						},
					},
				},
				DataStores: []element.DataStore{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_d367181e-c25a-407e-a5fb-817a27698a30",
							},
						},
						ItemSubjectRef: "_triso-default-bpmnItemDefinition-string_id",
						Name:           "Employee Details",
						IsUnlimited:    false,
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_7c793f30-d450-4a8e-becf-f168b14317ea",
							},
						},
						Name:        "User Management",
						IsUnlimited: false,
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_2addcc9a-638e-4496-b84f-5b0e28fe6536",
							},
						},
						Name:        "Payroll system",
						IsUnlimited: false,
					},
				},
				Processes: []element.Process{
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "_42cba3a9-a8ab-40b5-b9a4-2e8f32be364e",
								},
							},
							Name: "Money Bank - Process",
						},
						IsClosed: false,
						LaneSets: []element.LaneSet{
							{
								Lanes: []element.Lane{
									{
										BaseElement: element.BaseElement{
											ID: "_ff7ff8f6-a4f1-4e93-84e1-01cdb85eb755",
										},
										Name: "HR Department",
										FlowNodeRefs: []string{
											"_a4220c17-364f-4a08-ae9c-757a6468b295",
											"_f8973a92-3d84-4672-a1a3-b0df154121e1",
											"_f9e3cd76-809a-48b5-be1c-e84fc4324268",
											"_aa275782-c989-49ba-bf94-c58916ca7bb5",
											"_987b9b74-333a-4043-a72a-daadf667acc7",
											"_305ddf53-49a8-4105-ad06-70272a2332aa",
											"_0e71ed63-93f9-44b6-a89d-da9628652926",
											"_eba690b9-34ef-49e4-b265-1411809d9302",
											"_67944b4c-4950-45a2-a131-1c4679c6b433",
											"_4c95f4a0-f4ec-45ed-9fdb-7b236155d6f5",
											"_82da02ca-ee9a-4403-9f3b-aad030e089b9",
										},
									},
									{
										BaseElement: element.BaseElement{
											ID: "_937b5086-463f-4c8c-837d-f5eee5cbc1f4",
										},
										Name: "Responsible Department",
										FlowNodeRefs: []string{
											"_986cf801-0780-49d3-91cd-2cc6d3c1aac3",
											"_855451b0-5298-48b2-a81d-84ecbcca0a85",
											"_72da5cee-0456-4c3c-ba8d-6dd085d6f52d",
											"_e3d3ac43-74a3-48ff-9a02-e64b1358cc34",
											"_80f70d22-fb42-403f-8bdb-6805e9467bb7",
											"_fe77c2f2-278f-4752-9d03-aa0c8a12af1e",
											"_db9147a9-7fbc-4657-a506-15e777f2cfd9",
											"_74e2cc7b-99ca-426b-ad53-ad70a56506aa",
											"_19808f32-dfb5-462d-aaa6-e662f9932dba",
											"_351b058e-c37c-4fb7-9d32-24075f53ce02",
											"_52401cbb-02b8-4eaf-84f1-1edbc0854a4a",
											"_36baf139-fb74-43ef-8936-d490238c2825",
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
														ID: "_a4220c17-364f-4a08-ae9c-757a6468b295",
													},
													Name: "Candidate accepted offer",
												},
												Outgoing: []string{"_6d1a9d5a-127d-477a-8dbe-268f125f20d8"},
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
														ID: "_f8973a92-3d84-4672-a1a3-b0df154121e1",
													},
													Name: "Send \ncandidate Contract",
												},
												Incoming: []string{"_6d1a9d5a-127d-477a-8dbe-268f125f20d8", "_9ca6c9e6-4aff-4ec8-9fd3-2743c06483c3"},
												Outgoing: []string{"_4d3df0ce-e5be-443b-a379-8dcd4aae4de9"},
											},
										},
									},
									Implementation: "##unspecified",
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_aa275782-c989-49ba-bf94-c58916ca7bb5",
													},
													Name: "Get signature on contract and notify responsible department",
												},
												Incoming: []string{"_237c8380-5449-446e-a323-aad80181176d"},
												Outgoing: []string{"_41254af9-ca6b-49ae-a58e-c48f4c6120d2"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_d0960287-a964-4916-b4fe-b1238d9e3e75",
															},
															ItemSubjectRef: "_triso-default-bpmnItemDefinition-string_id",
															Name:           "Employee Details",
														},
													},
													InputSets: []element.InputSet{{}},
													OutputSets: []element.OutputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_2850d02b-5efe-4fa8-82c9-eac458b2e4d8",
															},
															DataOutputRefs: []string{"_d0960287-a964-4916-b4fe-b1238d9e3e75"},
														},
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_5c5f0068-874e-4dbb-ad30-ae8389921645",
														},
														SourceRef: "_d0960287-a964-4916-b4fe-b1238d9e3e75",
														TargetRef: "_f00a95f6-043e-4ec1-8731-153a707e950c",
													},
												},
											},
										},
									},
									Implementation: "##unspecified",
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_987b9b74-333a-4043-a72a-daadf667acc7",
													},
													Name: "Review terms of contract",
												},
												Incoming: []string{"_7e9d8b8b-faa9-4264-858b-7454702c4ec2"},
												Outgoing: []string{"_9ca6c9e6-4aff-4ec8-9fd3-2743c06483c3"},
											},
										},
									},
									Implementation: "##unspecified",
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_0e71ed63-93f9-44b6-a89d-da9628652926",
													},
													Name: "Inform employee of company policies",
												},
												Incoming: []string{"_d1c906b1-8f68-4485-94ea-00e732fba73e"},
												Outgoing: []string{"_5ae80bdd-cda7-491b-b3f7-bce2c4b35a83"},
											},
										},
									},
									Implementation: "##unspecified",
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_eba690b9-34ef-49e4-b265-1411809d9302",
													},
													Name: "Introduce employee to company Mission, Vision and Values",
												},
												Incoming: []string{"_5ae80bdd-cda7-491b-b3f7-bce2c4b35a83"},
												Outgoing: []string{"_15f3cde0-eeb5-4ad2-820a-114815aa782d"},
											},
										},
									},
									Implementation: "##unspecified",
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_67944b4c-4950-45a2-a131-1c4679c6b433",
													},
													Name: "Perform training for time reports sick leave and holidays",
												},
												Incoming: []string{"_15f3cde0-eeb5-4ad2-820a-114815aa782d"},
												Outgoing: []string{"_2d93c5db-2e41-463b-a1a6-63a05616d912"},
											},
										},
									},
									Implementation: "##unspecified",
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_4c95f4a0-f4ec-45ed-9fdb-7b236155d6f5",
													},
													Name: "Register for medical insurance",
												},
												Incoming: []string{"_2d93c5db-2e41-463b-a1a6-63a05616d912"},
												Outgoing: []string{"_3104b7bb-a417-4557-8c1b-d95b1c7e36a4"},
											},
										},
									},
									Implementation: "##unspecified",
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_986cf801-0780-49d3-91cd-2cc6d3c1aac3",
													},
													Name: "Request preparations for a new employee",
												},
												Incoming: []string{"_0811a8cc-1fde-4319-a718-7fb3ce64ae0b"},
												Outgoing: []string{"_3bb98d19-0d0b-4f47-a097-30c60183e334"},
											},
										},
									},
									Implementation: "##unspecified",
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_72da5cee-0456-4c3c-ba8d-6dd085d6f52d",
													},
													Name: "Introduce new employee to the team",
												},
												Incoming: []string{"_cad69ba8-8b26-41f9-87a1-372a899dcbf3"},
												Outgoing: []string{"_2ef25e1e-bbdd-4cba-b411-6091ec8669d6"},
											},
										},
									},
									Implementation: "##unspecified",
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_e3d3ac43-74a3-48ff-9a02-e64b1358cc34",
													},
													Name: "Perform training for position",
												},
												Incoming: []string{"_2ef25e1e-bbdd-4cba-b411-6091ec8669d6"},
												Outgoing: []string{"_240df34c-86c2-49fe-b1ec-58ea25d35c1f"},
											},
										},
									},
									Implementation: "##unspecified",
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_351b058e-c37c-4fb7-9d32-24075f53ce02",
													},
													Name: "Compile welcome package",
												},
												Incoming: []string{"_ece7bc06-ccf7-49f5-b5f1-99e1103cf538"},
												Outgoing: []string{"_c8454490-8090-4cf7-a164-d0aded7bc4da"},
											},
										},
									},
									Implementation: "##unspecified",
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_52401cbb-02b8-4eaf-84f1-1edbc0854a4a",
													},
													Name: "Give employee welcome package",
												},
												Incoming: []string{"_c8454490-8090-4cf7-a164-d0aded7bc4da"},
												Outgoing: []string{"_0fdd6987-6236-4972-a7be-3431ed6654ec"},
											},
										},
									},
									Implementation: "##unspecified",
								},
							},
							ExclusiveGatewaies: []element.ExclusiveGateway{
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_f9e3cd76-809a-48b5-be1c-e84fc4324268",
												},
												Name: "Contract terms accepted ?",
											},
											Incoming: []string{"_4d3df0ce-e5be-443b-a379-8dcd4aae4de9"},
											Outgoing: []string{"_237c8380-5449-446e-a323-aad80181176d", "_7e9d8b8b-faa9-4264-858b-7454702c4ec2"},
										},
										GatewayDirection: "Diverging",
									},
								},
							},
							ParallelGatewaies: []element.ParallelGateway{
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_305ddf53-49a8-4105-ad06-70272a2332aa",
												},
												Name: "Non-exclusive Gateway",
											},
											Incoming: []string{"_41254af9-ca6b-49ae-a58e-c48f4c6120d2"},
											Outgoing: []string{"_d1c906b1-8f68-4485-94ea-00e732fba73e", "_0811a8cc-1fde-4319-a718-7fb3ce64ae0b"},
										},
										GatewayDirection: "Diverging",
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_82da02ca-ee9a-4403-9f3b-aad030e089b9",
												},
												Name: "Non-exclusive Gateway",
											},
											Incoming: []string{"_3104b7bb-a417-4557-8c1b-d95b1c7e36a4", "_71b2c3f1-0d34-4cfb-956c-cedf83660dfb"},
											Outgoing: []string{"_cad69ba8-8b26-41f9-87a1-372a899dcbf3"},
										},
										GatewayDirection: "Converging",
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_80f70d22-fb42-403f-8bdb-6805e9467bb7",
												},
												Name: "Non-exclusive Gateway",
											},
											Incoming: []string{"_240df34c-86c2-49fe-b1ec-58ea25d35c1f"},
											Outgoing: []string{"_2be942a3-d921-4ad9-bce7-33e3411d98b4", "_4b186d86-fae1-4f55-a3b6-b90fa27affc1", "_c1307108-5b3d-48af-ae53-11b30dd31bed"},
										},
										GatewayDirection: "Diverging",
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_19808f32-dfb5-462d-aaa6-e662f9932dba",
												},
												Name: "Non-exclusive Gateway",
											},
											Incoming: []string{"_e82e6ee3-24d6-419a-96c0-a147f1943c23", "_16e4630e-39c0-4e94-ab65-a8a622245adf", "_c8da08ab-0b59-4c21-9386-f857d624c471"},
											Outgoing: []string{"_ece7bc06-ccf7-49f5-b5f1-99e1103cf538"},
										},
										GatewayDirection: "Converging",
									},
								},
							},
							DataStoreReferences: []element.DataStoreReference{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_f00a95f6-043e-4ec1-8731-153a707e950c",
										},
										Name: "Employee Details",
									},
									DataStoreRef: "_d367181e-c25a-407e-a5fb-817a27698a30",
								},
							},
							IntermediateThrowEvents: []element.IntermediateThrowEvent{
								{
									ThrowEvent: element.ThrowEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_855451b0-5298-48b2-a81d-84ecbcca0a85",
													},
													Name: "New employee in department X",
												},
												Incoming: []string{"_3bb98d19-0d0b-4f47-a097-30c60183e334"},
												Outgoing: []string{"_71b2c3f1-0d34-4cfb-956c-cedf83660dfb"},
											},
										},
										DataInputs: []element.DataInput{
											{
												BaseElement: element.BaseElement{
													ID: "_d5083b64-f0cd-479b-a773-b477a20408b5",
												},
												ItemSubjectRef: "_triso-default-bpmnItemDefinition-string_id",
												Name:           "Employee Details",
											},
										},
										DataInputAssociations: []element.DataInputAssociation{
											{
												DataAssociation: element.DataAssociation{
													BaseElement: element.BaseElement{
														ID: "_7a5b52f3-ad7e-4c5c-9ac2-432c8318b821",
													},
													SourceRef: "_f00a95f6-043e-4ec1-8731-153a707e950c",
													TargetRef: "_d5083b64-f0cd-479b-a773-b477a20408b5",
												},
											},
										},
										InputSet: element.InputSet{
											BaseElement: element.BaseElement{
												ID: "_e304d0f8-6eaa-4a79-a6ff-b08e993d0d06",
											},
											DataInputRefs: []string{"_d5083b64-f0cd-479b-a773-b477a20408b5"},
										},
										EventDefinitions: element.EventDefinitions{
											SignalEventDefinitions: []element.SignalEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "_e027a3d0-02f7-40a6-b84c-3b49282070ba",
															},
														},
													},
													SignalRef: "_bf41cec3-da5e-4fc7-87f9-7520a720192b",
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
														ID: "_fe77c2f2-278f-4752-9d03-aa0c8a12af1e",
													},
													Name: "Input from Payroll ready",
												},
												Incoming: []string{"_4b186d86-fae1-4f55-a3b6-b90fa27affc1"},
												Outgoing: []string{"_e82e6ee3-24d6-419a-96c0-a147f1943c23"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "_9e9bd4fc-0a94-4956-ac98-daf2ddeb001c",
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
														ID: "_db9147a9-7fbc-4657-a506-15e777f2cfd9",
													},
													Name: "Input from Facilities ready",
												},
												Incoming: []string{"_c1307108-5b3d-48af-ae53-11b30dd31bed"},
												Outgoing: []string{"_c8da08ab-0b59-4c21-9386-f857d624c471"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "_3454f552-d717-4f8c-ae58-971bef888af0",
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
														ID: "_74e2cc7b-99ca-426b-ad53-ad70a56506aa",
													},
													Name: "Input from IT ready",
												},
												Incoming: []string{"_2be942a3-d921-4ad9-bce7-33e3411d98b4"},
												Outgoing: []string{"_16e4630e-39c0-4e94-ab65-a8a622245adf"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "_940a88a5-974a-4583-afed-4510a94e5249",
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
														ID: "_36baf139-fb74-43ef-8936-d490238c2825",
													},
													Name: "End Event",
												},
												Incoming: []string{"_0fdd6987-6236-4972-a7be-3431ed6654ec"},
											},
										},
									},
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_6d1a9d5a-127d-477a-8dbe-268f125f20d8",
										},
									},
									SourceRef: "_a4220c17-364f-4a08-ae9c-757a6468b295",
									TargetRef: "_f8973a92-3d84-4672-a1a3-b0df154121e1",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_4d3df0ce-e5be-443b-a379-8dcd4aae4de9",
										},
									},
									SourceRef: "_f8973a92-3d84-4672-a1a3-b0df154121e1",
									TargetRef: "_f9e3cd76-809a-48b5-be1c-e84fc4324268",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_237c8380-5449-446e-a323-aad80181176d",
										},
										Name: "Yes",
									},
									SourceRef: "_f9e3cd76-809a-48b5-be1c-e84fc4324268",
									TargetRef: "_aa275782-c989-49ba-bf94-c58916ca7bb5",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_7e9d8b8b-faa9-4264-858b-7454702c4ec2",
										},
										Name: "No",
									},
									SourceRef: "_f9e3cd76-809a-48b5-be1c-e84fc4324268",
									TargetRef: "_987b9b74-333a-4043-a72a-daadf667acc7",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_41254af9-ca6b-49ae-a58e-c48f4c6120d2",
										},
									},
									SourceRef: "_aa275782-c989-49ba-bf94-c58916ca7bb5",
									TargetRef: "_305ddf53-49a8-4105-ad06-70272a2332aa",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_d1c906b1-8f68-4485-94ea-00e732fba73e",
										},
									},
									SourceRef: "_305ddf53-49a8-4105-ad06-70272a2332aa",
									TargetRef: "_0e71ed63-93f9-44b6-a89d-da9628652926",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_5ae80bdd-cda7-491b-b3f7-bce2c4b35a83",
										},
									},
									SourceRef: "_0e71ed63-93f9-44b6-a89d-da9628652926",
									TargetRef: "_eba690b9-34ef-49e4-b265-1411809d9302",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_15f3cde0-eeb5-4ad2-820a-114815aa782d",
										},
									},
									SourceRef: "_eba690b9-34ef-49e4-b265-1411809d9302",
									TargetRef: "_67944b4c-4950-45a2-a131-1c4679c6b433",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_2d93c5db-2e41-463b-a1a6-63a05616d912",
										},
									},
									SourceRef: "_67944b4c-4950-45a2-a131-1c4679c6b433",
									TargetRef: "_4c95f4a0-f4ec-45ed-9fdb-7b236155d6f5",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_3104b7bb-a417-4557-8c1b-d95b1c7e36a4",
										},
									},
									SourceRef: "_4c95f4a0-f4ec-45ed-9fdb-7b236155d6f5",
									TargetRef: "_82da02ca-ee9a-4403-9f3b-aad030e089b9",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_0811a8cc-1fde-4319-a718-7fb3ce64ae0b",
										},
									},
									SourceRef: "_305ddf53-49a8-4105-ad06-70272a2332aa",
									TargetRef: "_986cf801-0780-49d3-91cd-2cc6d3c1aac3",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_3bb98d19-0d0b-4f47-a097-30c60183e334",
										},
									},
									SourceRef: "_986cf801-0780-49d3-91cd-2cc6d3c1aac3",
									TargetRef: "_855451b0-5298-48b2-a81d-84ecbcca0a85",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_71b2c3f1-0d34-4cfb-956c-cedf83660dfb",
										},
									},
									SourceRef: "_855451b0-5298-48b2-a81d-84ecbcca0a85",
									TargetRef: "_82da02ca-ee9a-4403-9f3b-aad030e089b9",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_cad69ba8-8b26-41f9-87a1-372a899dcbf3",
										},
									},
									SourceRef: "_82da02ca-ee9a-4403-9f3b-aad030e089b9",
									TargetRef: "_72da5cee-0456-4c3c-ba8d-6dd085d6f52d",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_2ef25e1e-bbdd-4cba-b411-6091ec8669d6",
										},
									},
									SourceRef: "_72da5cee-0456-4c3c-ba8d-6dd085d6f52d",
									TargetRef: "_e3d3ac43-74a3-48ff-9a02-e64b1358cc34",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_240df34c-86c2-49fe-b1ec-58ea25d35c1f",
										},
									},
									SourceRef: "_e3d3ac43-74a3-48ff-9a02-e64b1358cc34",
									TargetRef: "_80f70d22-fb42-403f-8bdb-6805e9467bb7",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_2be942a3-d921-4ad9-bce7-33e3411d98b4",
										},
									},
									SourceRef: "_80f70d22-fb42-403f-8bdb-6805e9467bb7",
									TargetRef: "_74e2cc7b-99ca-426b-ad53-ad70a56506aa",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_4b186d86-fae1-4f55-a3b6-b90fa27affc1",
										},
									},
									SourceRef: "_80f70d22-fb42-403f-8bdb-6805e9467bb7",
									TargetRef: "_fe77c2f2-278f-4752-9d03-aa0c8a12af1e",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_c1307108-5b3d-48af-ae53-11b30dd31bed",
										},
									},
									SourceRef: "_80f70d22-fb42-403f-8bdb-6805e9467bb7",
									TargetRef: "_db9147a9-7fbc-4657-a506-15e777f2cfd9",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_e82e6ee3-24d6-419a-96c0-a147f1943c23",
										},
									},
									SourceRef: "_fe77c2f2-278f-4752-9d03-aa0c8a12af1e",
									TargetRef: "_19808f32-dfb5-462d-aaa6-e662f9932dba",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_16e4630e-39c0-4e94-ab65-a8a622245adf",
										},
									},
									SourceRef: "_74e2cc7b-99ca-426b-ad53-ad70a56506aa",
									TargetRef: "_19808f32-dfb5-462d-aaa6-e662f9932dba",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_c8da08ab-0b59-4c21-9386-f857d624c471",
										},
									},
									SourceRef: "_db9147a9-7fbc-4657-a506-15e777f2cfd9",
									TargetRef: "_19808f32-dfb5-462d-aaa6-e662f9932dba",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_ece7bc06-ccf7-49f5-b5f1-99e1103cf538",
										},
									},
									SourceRef: "_19808f32-dfb5-462d-aaa6-e662f9932dba",
									TargetRef: "_351b058e-c37c-4fb7-9d32-24075f53ce02",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_c8454490-8090-4cf7-a164-d0aded7bc4da",
										},
									},
									SourceRef: "_351b058e-c37c-4fb7-9d32-24075f53ce02",
									TargetRef: "_52401cbb-02b8-4eaf-84f1-1edbc0854a4a",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_0fdd6987-6236-4972-a7be-3431ed6654ec",
										},
									},
									SourceRef: "_52401cbb-02b8-4eaf-84f1-1edbc0854a4a",
									TargetRef: "_36baf139-fb74-43ef-8936-d490238c2825",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_9ca6c9e6-4aff-4ec8-9fd3-2743c06483c3",
										},
									},
									SourceRef: "_987b9b74-333a-4043-a72a-daadf667acc7",
									TargetRef: "_f8973a92-3d84-4672-a1a3-b0df154121e1",
								},
							},
						},
					},
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "_f0035388-f829-470c-b82b-0b15c3da3399",
								},
							},
							Name: "IT - Process",
						},
						IsClosed: false,
						FlowElements: element.FlowElements{
							StartEvents: []element.StartEvent{
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_e9306b3f-3a77-42e1-b53e-2ed8ee45486d",
													},
													Name: "New \nemployee\nhired",
												},
												Outgoing: []string{"_495f04a4-ed17-4e7b-a952-11c06c2230d1"},
											},
										},
										DataOutputs: []element.DataOutput{
											{
												BaseElement: element.BaseElement{
													ID: "_e96a6b5c-b24a-4221-8d35-cd3a4e68641e",
												},
												ItemSubjectRef: "_triso-default-bpmnItemDefinition-string_id",
												Name:           "Employee Details",
											},
										},
										DataOutputAssociations: []element.DataOutputAssociation{
											{
												DataAssociation: element.DataAssociation{
													BaseElement: element.BaseElement{
														ID: "_ed210028-7993-4aa4-b2f7-951ab85243be",
													},
													SourceRef: "_e96a6b5c-b24a-4221-8d35-cd3a4e68641e",
													TargetRef: "_e355b511-e181-4e7b-a9b1-30cb48dd6cab",
												},
											},
										},
										OutputSet: element.OutputSet{
											BaseElement: element.BaseElement{
												ID: "_2d74cb72-2fcc-4869-b5f0-b867a7f2f7ec",
											},
											DataOutputRefs: []string{"_e96a6b5c-b24a-4221-8d35-cd3a4e68641e"},
										},
										EventDefinitions: element.EventDefinitions{
											SignalEventDefinitions: []element.SignalEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "_5ec0186c-9f7d-43b1-ba20-571cfe42bc5c",
															},
														},
													},
													SignalRef: "_bf41cec3-da5e-4fc7-87f9-7520a720192b",
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
														ID: "_7e9d2e5a-21f7-493b-9ae4-03245aa33a5c",
													},
													Name: "Create domain account",
												},
												Incoming: []string{"_495f04a4-ed17-4e7b-a952-11c06c2230d1"},
												Outgoing: []string{"_0bccc7f6-0f56-49bc-98df-59d70433c043"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_b1d40770-825b-4d23-a026-db95f815d520",
															},
															Name:           "Employee Details",
															ItemSubjectRef: "_triso-default-bpmnItemDefinition-string_id",
														},
													},
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_18dcd171-c602-488e-b53b-0289e2fae114",
															},
															Name: "User Management",
														},
													},
													InputSets: []element.InputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_c228336b-c661-4e3a-920f-8e13f4844656",
															},
															DataInputRefs: []string{"_b1d40770-825b-4d23-a026-db95f815d520"},
														},
													},
													OutputSets: []element.OutputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_3c50eff0-7da3-42bf-96be-d7e22f4a2bb5",
															},
															DataOutputRefs: []string{"_18dcd171-c602-488e-b53b-0289e2fae114"},
														},
													},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_c4eff1b8-8817-418b-9ef9-5d3738d43659",
														},
														SourceRef: "_e355b511-e181-4e7b-a9b1-30cb48dd6cab",
														TargetRef: "_b1d40770-825b-4d23-a026-db95f815d520",
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_0282dbd7-a11b-4cac-84e2-0f26bf86b853",
														},
														SourceRef: "_18dcd171-c602-488e-b53b-0289e2fae114",
														TargetRef: "_b60be020-1f04-46c2-b1a9-19bfc66ead82",
													},
												},
											},
										},
									},
									Implementation: "##unspecified",
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_912f731a-fb6b-499b-a3f4-4c1632d606bd",
													},
													Name: "Assign required applications and permissions",
												},
												Incoming: []string{"_dcec1142-90eb-41a8-b809-458b713ab94c"},
												Outgoing: []string{"_be72815e-ed18-4d77-9252-1ab0a7438dc9"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_4b0aa8de-e86c-492a-a5da-fd80464b71bd",
															},
															Name: "User Management",
														},
													},
													InputSets: []element.InputSet{{}},
													OutputSets: []element.OutputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_7809a0c3-5349-4e66-bf42-fa19126dd198",
															},
															DataOutputRefs: []string{"_4b0aa8de-e86c-492a-a5da-fd80464b71bd"},
														},
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_17c38035-fb11-4abb-b2b6-1cbd14acdab0",
														},
														SourceRef: "_4b0aa8de-e86c-492a-a5da-fd80464b71bd",
														TargetRef: "_b60be020-1f04-46c2-b1a9-19bfc66ead82",
													},
												},
											},
										},
									},
									Implementation: "##unspecified",
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_ad173aff-cfe3-4098-8c65-02f783ad9e1f",
													},
													Name: "Prepare IT part of welcome package",
												},
												Incoming: []string{"_ba635734-5672-4f1b-b753-3b08c778c64e"},
												Outgoing: []string{"_c904a36b-64f6-46e3-8ef1-9f2ecc3faaff"},
											},
										},
									},
									Implementation: "##unspecified",
								},
							},
							ManualTasks: []element.ManualTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_c29af228-0768-4dfe-945a-17755e173674",
													},
													Name: "Prepare workstation",
												},
												Incoming: []string{"_0bccc7f6-0f56-49bc-98df-59d70433c043"},
												Outgoing: []string{"_dcec1142-90eb-41a8-b809-458b713ab94c"},
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
														ID: "_9db2d136-aa33-4de2-be76-554e7843363d",
													},
													Name: "Configure workstation",
												},
												Incoming: []string{"_be72815e-ed18-4d77-9252-1ab0a7438dc9"},
												Outgoing: []string{"_ba635734-5672-4f1b-b753-3b08c778c64e"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_4606e5ae-aaa4-4416-9576-68f29181c5b2",
															},
															Name: "User Management",
														},
													},
													InputSets: []element.InputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_3b4c5d5b-b5bc-4155-8ff6-adb9f629de15",
															},
															DataInputRefs: []string{"_4606e5ae-aaa4-4416-9576-68f29181c5b2"},
														},
													},
													OutputSets: []element.OutputSet{{}},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_de54f4cb-0603-4389-ab3a-9edb9cbfec2a",
														},
														SourceRef: "_b60be020-1f04-46c2-b1a9-19bfc66ead82",
														TargetRef: "_4606e5ae-aaa4-4416-9576-68f29181c5b2",
													},
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
														ID: "_c82dd8eb-ce54-4aa7-b8c4-b8d3e8fd654e",
													},
													Name: "Workstation and permissions ready",
												},
												Incoming: []string{"_c904a36b-64f6-46e3-8ef1-9f2ecc3faaff"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "_2f246656-2c54-430d-ad1b-9c50428d1b9d",
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
											ID: "_e355b511-e181-4e7b-a9b1-30cb48dd6cab",
										},
										Name: "Employee Details",
									},
									DataStoreRef: "_d367181e-c25a-407e-a5fb-817a27698a30",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_b60be020-1f04-46c2-b1a9-19bfc66ead82",
										},
										Name: "User Management",
									},
									DataStoreRef: "_7c793f30-d450-4a8e-becf-f168b14317ea",
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_495f04a4-ed17-4e7b-a952-11c06c2230d1",
										},
									},
									SourceRef: "_e9306b3f-3a77-42e1-b53e-2ed8ee45486d",
									TargetRef: "_7e9d2e5a-21f7-493b-9ae4-03245aa33a5c",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_0bccc7f6-0f56-49bc-98df-59d70433c043",
										},
									},
									SourceRef: "_7e9d2e5a-21f7-493b-9ae4-03245aa33a5c",
									TargetRef: "_c29af228-0768-4dfe-945a-17755e173674",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_dcec1142-90eb-41a8-b809-458b713ab94c",
										},
									},
									SourceRef: "_c29af228-0768-4dfe-945a-17755e173674",
									TargetRef: "_912f731a-fb6b-499b-a3f4-4c1632d606bd",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_be72815e-ed18-4d77-9252-1ab0a7438dc9",
										},
									},
									SourceRef: "_912f731a-fb6b-499b-a3f4-4c1632d606bd",
									TargetRef: "_9db2d136-aa33-4de2-be76-554e7843363d",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_ba635734-5672-4f1b-b753-3b08c778c64e",
										},
									},
									SourceRef: "_9db2d136-aa33-4de2-be76-554e7843363d",
									TargetRef: "_ad173aff-cfe3-4098-8c65-02f783ad9e1f",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_c904a36b-64f6-46e3-8ef1-9f2ecc3faaff",
										},
									},
									SourceRef: "_ad173aff-cfe3-4098-8c65-02f783ad9e1f",
									TargetRef: "_c82dd8eb-ce54-4aa7-b8c4-b8d3e8fd654e",
								},
							},
						},
						Artifacts: element.Artifacts{
							TextAnnotations: []element.TextAnnotation{
								{
									Artifact: element.Artifact{
										BaseElement: element.BaseElement{
											ID: "_246c3cf7-fb4b-4c05-8ece-55c77fb4dc1e",
										},
									},
									Text: "With PowerShell",
								},
							},
							Associations: []element.Association{
								{
									Artifact: element.Artifact{
										BaseElement: element.BaseElement{
											ID: "_fa0de585-aa40-40bb-a3a5-965d956e884d",
										},
									},
									AssociationDirection: "None",
									SourceRef:            "_9db2d136-aa33-4de2-be76-554e7843363d",
									TargetRef:            "_246c3cf7-fb4b-4c05-8ece-55c77fb4dc1e",
								},
							},
						},
					},
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "_da743a6f-d9e5-4fcf-8a96-d2fd5cfb73d4",
								},
							},
							Name: "Payroll - Process",
						},
						IsClosed: false,
						FlowElements: element.FlowElements{
							StartEvents: []element.StartEvent{
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_3d4130c6-48c9-47fe-8e95-2eeb56060e2b",
													},
													Name: "New \nemployee\nhired",
												},
												Outgoing: []string{"_6b24f6f7-bc44-4d73-ab54-1ac8b638423e"},
											},
										},
										DataOutputs: []element.DataOutput{
											{
												BaseElement: element.BaseElement{
													ID: "_06fff5da-cd41-4538-a6b8-cd73c6368aa3",
												},
												Name:           "Employee Details",
												ItemSubjectRef: "_triso-default-bpmnItemDefinition-string_id",
											},
										},
										DataOutputAssociations: []element.DataOutputAssociation{
											{
												DataAssociation: element.DataAssociation{
													BaseElement: element.BaseElement{
														ID: "_fba68566-2c5e-49e2-8f1e-2feff2873742",
													},
													SourceRef: "_06fff5da-cd41-4538-a6b8-cd73c6368aa3",
													TargetRef: "_16498c9e-c33f-4f92-9033-714b89582a7c",
												},
											},
										},
										OutputSet: element.OutputSet{
											BaseElement: element.BaseElement{
												ID: "_0b28282c-3c90-4b70-920c-a6898a3fc18a",
											},
											DataOutputRefs: []string{"_06fff5da-cd41-4538-a6b8-cd73c6368aa3"},
										},
										EventDefinitions: element.EventDefinitions{
											SignalEventDefinitions: []element.SignalEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "_1f790746-cb6c-4eaa-85a8-a0a6c0fa7602",
															},
														},
													},
													SignalRef: "_bf41cec3-da5e-4fc7-87f9-7520a720192b",
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
														ID: "_ae47ce79-bd91-452b-be68-47a2ea589e75",
													},
													Name: "Validate provided information",
												},
												Incoming: []string{"_6b24f6f7-bc44-4d73-ab54-1ac8b638423e"},
												Outgoing: []string{"_d110be18-ff17-4375-818b-6463eff4e2d6"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_b3935a94-0b1b-4cd7-8a3b-a10f5b8b5e34",
															},
															ItemSubjectRef: "_triso-default-bpmnItemDefinition-string_id",
															Name:           "Employee Details",
														},
													},
													InputSets: []element.InputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_4dbfa005-dede-45b3-adfd-f81e5900e5a6",
															},
															DataInputRefs: []string{"_b3935a94-0b1b-4cd7-8a3b-a10f5b8b5e34"},
														},
													},
													OutputSets: []element.OutputSet{{}},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_3ba0a27a-2a49-458b-8680-be1330f46f3a",
														},
														SourceRef: "_16498c9e-c33f-4f92-9033-714b89582a7c",
														TargetRef: "_b3935a94-0b1b-4cd7-8a3b-a10f5b8b5e34",
													},
												},
											},
										},
									},
									Implementation: "##unspecified",
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_9dbd92a5-5c0a-4039-b741-bf4ede54ccf0",
													},
													Name: "Update payroll system",
												},
												Incoming: []string{"_12e2f256-ec10-4a8d-934c-329143c588d8", "_7b347349-7441-4e5c-ab24-2e614051f222"},
												Outgoing: []string{"_46c8bd85-1b7f-45fb-81c6-23d22300c0cd"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_de11bfcc-129c-4570-b4f0-2b395f932f16",
															},
															Name: "Payroll system",
														},
													},
													InputSets: []element.InputSet{{}},
													OutputSets: []element.OutputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_468ab3e6-4b63-4f54-9d07-16c8f1320ef1",
															},
															DataOutputRefs: []string{"_de11bfcc-129c-4570-b4f0-2b395f932f16"},
														},
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_079db308-7d1e-497d-962e-23499236c981",
														},
														SourceRef: "_de11bfcc-129c-4570-b4f0-2b395f932f16",
														TargetRef: "_a9cffac5-6e03-41de-83bf-d34304f28b1a",
													},
												},
											},
										},
									},
									Implementation: "##unspecified",
								},
							},
							ExclusiveGatewaies: []element.ExclusiveGateway{
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_fa14ca2d-ea97-49a2-b75e-72e7d27d6fd1",
												},
												Name: "All necessary data available?",
											},
											Incoming: []string{"_d110be18-ff17-4375-818b-6463eff4e2d6"},
											Outgoing: []string{"_ca6f904d-e30d-4777-a7dd-661650e1e3a2", "_12e2f256-ec10-4a8d-934c-329143c588d8"},
										},
										GatewayDirection: "Diverging",
									},
								},
							},
							ManualTasks: []element.ManualTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_788443d9-65f0-43a4-96a8-63e8d6f380a7",
													},
													Name: "Clarify missing points",
												},
												Incoming: []string{"_ca6f904d-e30d-4777-a7dd-661650e1e3a2"},
												Outgoing: []string{"_7b347349-7441-4e5c-ab24-2e614051f222"},
											},
											LoopCharacteristicsElements: element.LoopCharacteristicsElements{
												StandardLoopCharacteristics: []element.StandardLoopCharacteristics{
													{
														LoopCharacteristics: element.LoopCharacteristics{
															BaseElement: element.BaseElement{
																ID: "_be244352-1e67-4664-9a10-5d088542f02e",
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
														ID: "_efbd0983-76cd-4a4c-acf3-6dde71d7c760",
													},
													Name: "Payroll ready",
												},
												Incoming: []string{"_46c8bd85-1b7f-45fb-81c6-23d22300c0cd"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "_0be7d261-5be3-40e4-a6fd-954d9c836607",
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
											ID: "_16498c9e-c33f-4f92-9033-714b89582a7c",
										},
										Name: "Employee Details",
									},
									DataStoreRef: "_d367181e-c25a-407e-a5fb-817a27698a30",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_a9cffac5-6e03-41de-83bf-d34304f28b1a",
										},
										Name: "Payroll system",
									},
									DataStoreRef: "_2addcc9a-638e-4496-b84f-5b0e28fe6536",
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_6b24f6f7-bc44-4d73-ab54-1ac8b638423e",
										},
									},
									SourceRef: "_3d4130c6-48c9-47fe-8e95-2eeb56060e2b",
									TargetRef: "_ae47ce79-bd91-452b-be68-47a2ea589e75",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_d110be18-ff17-4375-818b-6463eff4e2d6",
										},
									},
									SourceRef: "_ae47ce79-bd91-452b-be68-47a2ea589e75",
									TargetRef: "_fa14ca2d-ea97-49a2-b75e-72e7d27d6fd1",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_ca6f904d-e30d-4777-a7dd-661650e1e3a2",
										},
										Name: "No",
									},
									SourceRef: "_fa14ca2d-ea97-49a2-b75e-72e7d27d6fd1",
									TargetRef: "_788443d9-65f0-43a4-96a8-63e8d6f380a7",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_12e2f256-ec10-4a8d-934c-329143c588d8",
										},
										Name: "Yes",
									},
									SourceRef: "_fa14ca2d-ea97-49a2-b75e-72e7d27d6fd1",
									TargetRef: "_9dbd92a5-5c0a-4039-b741-bf4ede54ccf0",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_7b347349-7441-4e5c-ab24-2e614051f222",
										},
									},
									SourceRef: "_788443d9-65f0-43a4-96a8-63e8d6f380a7",
									TargetRef: "_9dbd92a5-5c0a-4039-b741-bf4ede54ccf0",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_46c8bd85-1b7f-45fb-81c6-23d22300c0cd",
										},
									},
									SourceRef: "_9dbd92a5-5c0a-4039-b741-bf4ede54ccf0",
									TargetRef: "_efbd0983-76cd-4a4c-acf3-6dde71d7c760",
								},
							},
						},
					},
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "_3486bf55-0a7f-4ff1-be15-1555669f58ad",
								},
							},
							Name: "Facilities - Process",
						},
						IsClosed: false,
						FlowElements: element.FlowElements{
							StartEvents: []element.StartEvent{
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_94a62738-dc7a-49f6-81d8-f5642f7ae850",
													},
													Name: "New \nemployee \nhired",
												},
												Outgoing: []string{"_ea0e1abf-e9f6-406a-96fb-12fcdf765634"},
											},
										},
										DataOutputs: []element.DataOutput{
											{
												BaseElement: element.BaseElement{
													ID: "_811f33cd-dd3b-4ecd-aa83-342ee7ccb86c",
												},
												ItemSubjectRef: "_triso-default-bpmnItemDefinition-string_id",
												Name:           "Employee Details",
											},
										},
										DataOutputAssociations: []element.DataOutputAssociation{
											{
												DataAssociation: element.DataAssociation{
													BaseElement: element.BaseElement{
														ID: "_2ae32762-1124-45f5-a01b-7292f9721c7d",
													},
													SourceRef: "_811f33cd-dd3b-4ecd-aa83-342ee7ccb86c",
													TargetRef: "_9931ba94-ede8-4f53-8081-878c549ba1c8",
												},
											},
										},
										OutputSet: element.OutputSet{
											BaseElement: element.BaseElement{
												ID: "_e0d2138c-bd9b-400c-9501-7fdfc1b41f5f",
											},
											DataOutputRefs: []string{"_811f33cd-dd3b-4ecd-aa83-342ee7ccb86c"},
										},
										EventDefinitions: element.EventDefinitions{
											SignalEventDefinitions: []element.SignalEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "_3acc62e3-0cfd-400d-8691-3f05b1589d95",
															},
														},
													},
													SignalRef: "_bf41cec3-da5e-4fc7-87f9-7520a720192b",
												},
											},
										},
									},
								},
							},
							ManualTasks: []element.ManualTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_2bf94039-15a1-44bb-9d14-81358777466c",
													},
													Name: "Prepare\naccess \ncard",
												},
												Incoming: []string{"_ea0e1abf-e9f6-406a-96fb-12fcdf765634"},
												Outgoing: []string{"_4436759d-6848-4400-93e5-3971c69b2f67"},
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
														ID: "_737503c8-10bc-483f-8871-5461d822b469",
													},
													Name: "Configure access details",
												},
												Incoming: []string{"_4436759d-6848-4400-93e5-3971c69b2f67"},
												Outgoing: []string{"_f405a593-5f87-4f62-862e-bcf9f776bc95"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_946f3a81-8249-481e-afc4-07647f204192",
															},
															ItemSubjectRef: "_triso-default-bpmnItemDefinition-string_id",
															Name:           "Employee Details",
														},
													},
													InputSets: []element.InputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_9e2a6ff5-7f57-4b9a-873c-e3905310f0c7",
															},
															DataInputRefs: []string{"_946f3a81-8249-481e-afc4-07647f204192"},
														},
													},
													OutputSets: []element.OutputSet{{}},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_ad35d446-d8b5-46d9-b2b0-61178b689754",
														},
														SourceRef: "_9931ba94-ede8-4f53-8081-878c549ba1c8",
														TargetRef: "_946f3a81-8249-481e-afc4-07647f204192",
													},
												},
											},
										},
									},
									Implementation: "##unspecified",
								},
							},
							EndEvents: []element.EndEvent{
								{
									ThrowEvent: element.ThrowEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_5ee09fe4-f38f-454d-b6e4-1c3703a6a239",
													},
													Name: "Access \ncard \nready",
												},
												Incoming: []string{"_f405a593-5f87-4f62-862e-bcf9f776bc95"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "_a4a343f5-62b7-4939-95c4-93b2520d5fa6",
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
											ID: "_9931ba94-ede8-4f53-8081-878c549ba1c8",
										},
										Name: "Employee Details",
									},
									DataStoreRef: "_d367181e-c25a-407e-a5fb-817a27698a30",
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_ea0e1abf-e9f6-406a-96fb-12fcdf765634",
										},
									},
									SourceRef: "_94a62738-dc7a-49f6-81d8-f5642f7ae850",
									TargetRef: "_2bf94039-15a1-44bb-9d14-81358777466c",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_4436759d-6848-4400-93e5-3971c69b2f67",
										},
									},
									SourceRef: "_2bf94039-15a1-44bb-9d14-81358777466c",
									TargetRef: "_737503c8-10bc-483f-8871-5461d822b469",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_f405a593-5f87-4f62-862e-bcf9f776bc95",
										},
									},
									SourceRef: "_737503c8-10bc-483f-8871-5461d822b469",
									TargetRef: "_5ee09fe4-f38f-454d-b6e4-1c3703a6a239",
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
