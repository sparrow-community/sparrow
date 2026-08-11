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

func TestC_6_0_roundtrip(t *testing.T) {
	// create test use ./test/C.6.0-roundtrip.bpmn
	path := "./test/C.6.0-roundtrip.bpmn"
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
			ID:              "_af3cbf6a-9c69-460b-9a1d-276711d30213",
			Name:            "Travel Booking with Event Subprocess",
			TargetNamespace: "http://www.trisotech.com/definitions/_af3cbf6a-9c69-460b-9a1d-276711d30213",
			Exporter:        "BPMN Modeler",
			ExporterVersion: "6.2.9.1",
			RootElemnts: element.RootElemnts{
				Processes: []element.Process{
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "_898aa942-9a96-4405-ae71-22b5e2e3d235",
								},
							},
							Name: "Simple Travel Booking",
							IoSpecification: element.IoSpecification{
								InputOutputSpecification: element.InputOutputSpecification{
									DataInputs: []element.DataInput{
										{
											BaseElement: element.BaseElement{
												ID: "_e8347f3a-e1a1-41d1-b2e1-4089c02028bc",
											},
											Name: "Travel Request",
										},
									},
									DataOutputs: []element.DataOutput{
										{
											BaseElement: element.BaseElement{
												ID: "_1a713fab-bf18-4975-afe3-d2bd2820023b",
											},
											Name: "Itinerary",
										},
									},
									InputSets: []element.InputSet{
										{
											BaseElement: element.BaseElement{
												ID: "_41b1bf64-dba0-4ab1-a834-0ee931fb4cab",
											},
											DataInputRefs: []string{"_e8347f3a-e1a1-41d1-b2e1-4089c02028bc"},
										},
									},
									OutputSets: []element.OutputSet{
										{
											BaseElement: element.BaseElement{
												ID: "_2d5ecd07-500d-4693-bb58-8682897d59e2",
											},
											DataOutputRefs: []string{"_1a713fab-bf18-4975-afe3-d2bd2820023b"},
										},
									},
								},
							},
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
														ID: "_44e3f1fa-42cd-40b7-9980-a51ac49d5fa3",
													},
													Name: "Receive Travel Request",
												},
												Outgoing: []string{"_9ec5ad1d-0662-42a4-9118-bca0bf744422"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "_1087adb7-2560-4fec-a13c-5a91f257819c",
															},
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
											ID: "_c7c2f8e1-a93f-44d6-9a94-cb92d7466d73",
										},
									},
									SourceRef: "_15fef309-6718-4352-9b71-f757bcd8c023",
									TargetRef: "_e839800f-ad4f-4bcc-aaf2-d38fe4a32bcd",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_16857297-f0fb-4027-87f6-633a76267c84",
										},
									},
									SourceRef: "_7ab6dbdf-f55b-4be6-bb41-d99793135c1d",
									TargetRef: "_87baeef0-f32e-4a93-b802-fdd588aaf729",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_54044d90-3a1e-4a49-99a3-1497045c7235",
										},
									},
									SourceRef: "_7ab6dbdf-f55b-4be6-bb41-d99793135c1d",
									TargetRef: "_15fef309-6718-4352-9b71-f757bcd8c023",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_2403da0e-5b49-43fa-ba00-ec6c1f3e4790",
										},
									},
									SourceRef: "_7ab6dbdf-f55b-4be6-bb41-d99793135c1d",
									TargetRef: "_e5c69e92-6f98-47c8-bc22-b75d38620f95",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_b792d502-a08e-4a3d-bdc1-1fc33a2022f6",
										},
									},
									SourceRef: "_6db9a77f-189f-4c07-b33b-e0c88f09e0db",
									TargetRef: "_6a5cdbbf-2618-496e-b728-955dc215ef9d",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_c69565cd-3a08-4532-b466-8f9a9ebdc1e1",
										},
									},
									SourceRef: "_22612d45-65ca-4a74-a6eb-53af7ebcb5ff",
									TargetRef: "_42e03d0f-6c6b-4493-971f-c6928eb563b0",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_8eac4fb5-c746-42e4-ab0b-80cdfaf76007",
										},
									},
									SourceRef: "_54d6e34f-4613-48d4-9143-b90927448adb",
									TargetRef: "_c8fa5253-dde9-471e-b933-58b00e8f374c",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_ec8c7430-02e1-4298-9b00-33c4d56d0fac",
										},
									},
									SourceRef: "_c8fa5253-dde9-471e-b933-58b00e8f374c",
									TargetRef: "_a68b0941-b7e3-4791-8e13-fd9622e4448c",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_65293cb4-b859-4ff7-bcd3-ba4e38cf8763",
										},
									},
									SourceRef: "_614d6469-2bb8-4ad6-a20a-db5db6321c6b",
									TargetRef: "_22612d45-65ca-4a74-a6eb-53af7ebcb5ff",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_f81560f5-e1c5-497d-a159-dce350f2283c",
										},
									},
									SourceRef: "_c38139c7-a2d1-47c7-b75a-19e14c7212c8",
									TargetRef: "_614d6469-2bb8-4ad6-a20a-db5db6321c6b",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_9ec5ad1d-0662-42a4-9118-bca0bf744422",
										},
									},
									SourceRef: "_44e3f1fa-42cd-40b7-9980-a51ac49d5fa3",
									TargetRef: "_9cc2ac34-f12c-49e0-b37c-144e5a84fd92",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_bcb8805c-5e93-41ec-991b-575d26ed3a9c",
										},
									},
									SourceRef: "_9cc2ac34-f12c-49e0-b37c-144e5a84fd92",
									TargetRef: "_7ab6dbdf-f55b-4be6-bb41-d99793135c1d",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_f072c57e-dcc6-4a16-8f28-13a2cd6e43ed",
										},
									},
									SourceRef: "_e5c69e92-6f98-47c8-bc22-b75d38620f95",
									TargetRef: "_8afc49f0-42c2-4da9-8e79-e08dbe349776",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_e396fc69-5792-41e3-8a92-ee6b83cd0909",
										},
									},
									SourceRef: "_8afc49f0-42c2-4da9-8e79-e08dbe349776",
									TargetRef: "_7eb87eb8-0d7a-445b-b768-90d754a938ed",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_649072fc-31c7-4e04-bc58-d4442a06f076",
										},
									},
									SourceRef: "_87baeef0-f32e-4a93-b802-fdd588aaf729",
									TargetRef: "_de7e721e-a073-4857-b8c3-c6ae886dbb46",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_51e34f49-b9eb-4f10-ac0a-7af81e1370ff",
										},
									},
									SourceRef: "_de7e721e-a073-4857-b8c3-c6ae886dbb46",
									TargetRef: "_88247168-b457-4663-aad0-753a0236c8df",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_a0c99000-31e1-4de9-baf7-45c2d09ebbff",
										},
									},
									SourceRef: "_32c4138c-74ae-484a-a7e5-0609370d7080",
									TargetRef: "_de7e721e-a073-4857-b8c3-c6ae886dbb46",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_c72535b8-43f7-4906-9f18-691e03050c24",
										},
									},
									SourceRef: "_6a5cdbbf-2618-496e-b728-955dc215ef9d",
									TargetRef: "_2d6586cf-81fc-4e2a-83ec-6cfff5b34bb0",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_94e0bd86-6659-4669-baee-ce671c7f33d6",
										},
									},
									SourceRef: "_2d6586cf-81fc-4e2a-83ec-6cfff5b34bb0",
									TargetRef: "_babdfa54-b55f-463f-9341-424b42db9760",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_7efd61ea-979e-44f4-91c8-e212ae6b22c7",
										},
									},
									SourceRef: "_e839800f-ad4f-4bcc-aaf2-d38fe4a32bcd",
									TargetRef: "_c38139c7-a2d1-47c7-b75a-19e14c7212c8",
								},
							},
							SendTasks: []element.SendTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_e839800f-ad4f-4bcc-aaf2-d38fe4a32bcd",
													},
													Name: "Request\nCredit Card Information",
												},
												Incoming: []string{"_c7c2f8e1-a93f-44d6-9a94-cb92d7466d73"},
												Outgoing: []string{"_7efd61ea-979e-44f4-91c8-e212ae6b22c7"},
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
														ID: "_22612d45-65ca-4a74-a6eb-53af7ebcb5ff",
													},
													Name: "Confirm Booking",
												},
												Incoming: []string{"_65293cb4-b859-4ff7-bcd3-ba4e38cf8763"},
												Outgoing: []string{"_c69565cd-3a08-4532-b466-8f9a9ebdc1e1"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_29b8c649-e2a0-4dd3-804b-567e8cc71718",
															},
															Name: "Itinerary",
														},
													},
													InputSets: []element.InputSet{{}},
													OutputSets: []element.OutputSet{{
														BaseElement: element.BaseElement{
															ID: "_f4813b3e-2820-45ec-9506-0f2a4c7bd086",
														},
														DataOutputRefs: []string{"_29b8c649-e2a0-4dd3-804b-567e8cc71718"},
													}},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_ef945053-2c76-4796-8d8b-afc6be8bbfec",
														},
														SourceRef: "_29b8c649-e2a0-4dd3-804b-567e8cc71718",
														TargetRef: "_1a713fab-bf18-4975-afe3-d2bd2820023b",
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
														ID: "_c8fa5253-dde9-471e-b933-58b00e8f374c",
													},
													Name: "Notify Failed Booking",
												},
												Incoming: []string{"_8eac4fb5-c746-42e4-ab0b-80cdfaf76007"},
												Outgoing: []string{"_ec8c7430-02e1-4298-9b00-33c4d56d0fac"},
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
														ID: "_9cc2ac34-f12c-49e0-b37c-144e5a84fd92",
													},
													Name: "Make Flights and Hotel Offer",
												},
												Incoming: []string{"_9ec5ad1d-0662-42a4-9118-bca0bf744422"},
												Outgoing: []string{"_bcb8805c-5e93-41ec-991b-575d26ed3a9c"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_9628422b-85a6-4857-8c14-7289b9fd9a8a",
															},
															Name: "Travel Request",
														},
													},
													InputSets: []element.InputSet{{
														BaseElement: element.BaseElement{
															ID: "_5a361a87-dcba-4107-a625-026abada856d",
														},
														DataInputRefs: []string{"_9628422b-85a6-4857-8c14-7289b9fd9a8a"},
													}},
													OutputSets: []element.OutputSet{{}},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_82006a49-2cc8-4814-83e2-28c5a85a4c4a",
														},
														SourceRef: "_e8347f3a-e1a1-41d1-b2e1-4089c02028bc",
														TargetRef: "_9628422b-85a6-4857-8c14-7289b9fd9a8a",
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
														ID: "_de7e721e-a073-4857-b8c3-c6ae886dbb46",
													},
													Name: "Notify Customer\nOffer Expired",
												},
												Incoming: []string{"_649072fc-31c7-4e04-bc58-d4442a06f076", "_a0c99000-31e1-4de9-baf7-45c2d09ebbff"},
												Outgoing: []string{"_51e34f49-b9eb-4f10-ac0a-7af81e1370ff"},
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
														ID: "_2d6586cf-81fc-4e2a-83ec-6cfff5b34bb0",
													},
													Name: "Notify Failed Credit Transaction",
												},
												Incoming: []string{"_c72535b8-43f7-4906-9f18-691e03050c24"},
												Outgoing: []string{"_94e0bd86-6659-4669-baee-ce671c7f33d6"},
											},
										},
									},
									Implementation: "##WebService",
								},
							},
							EventBasedGatewaies: []element.EventBasedGateway{
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_7ab6dbdf-f55b-4be6-bb41-d99793135c1d",
												},
											},
											Incoming: []string{"_bcb8805c-5e93-41ec-991b-575d26ed3a9c"},
											Outgoing: []string{"_16857297-f0fb-4027-87f6-633a76267c84", "_54044d90-3a1e-4a49-99a3-1497045c7235", "_2403da0e-5b49-43fa-ba00-ec6c1f3e4790"},
										},
										GatewayDirection: "Diverging",
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
														ID: "_15fef309-6718-4352-9b71-f757bcd8c023",
													},
													Name: "Offer Approved",
												},
												Incoming: []string{"_54044d90-3a1e-4a49-99a3-1497045c7235"},
												Outgoing: []string{"_c7c2f8e1-a93f-44d6-9a94-cb92d7466d73"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{},
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
														ID: "_e5c69e92-6f98-47c8-bc22-b75d38620f95",
													},
													Name: "Cancel Request",
												},
												Incoming: []string{"_2403da0e-5b49-43fa-ba00-ec6c1f3e4790"},
												Outgoing: []string{"_f072c57e-dcc6-4a16-8f28-13a2cd6e43ed"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{{}},
										},
									},
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_87baeef0-f32e-4a93-b802-fdd588aaf729",
													},
													Name: "24 Hours",
												},
												Incoming: []string{"_16857297-f0fb-4027-87f6-633a76267c84"},
												Outgoing: []string{"_649072fc-31c7-4e04-bc58-d4442a06f076"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											TimerEventDefinitions: []element.TimerEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{},
														},
													},
													TimeDate: element.ExpressionUnMarshal{ExpressionSubstitution: &element.Expression{}},
												},
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
														ID: "_6a5cdbbf-2618-496e-b728-955dc215ef9d",
													},
													Name: "Booking",
												},
												Incoming: []string{"_b792d502-a08e-4a3d-bdc1-1fc33a2022f6"},
												Outgoing: []string{"_c72535b8-43f7-4906-9f18-691e03050c24"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											CompensateEventDefinitions: []element.CompensateEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{},
														},
													},
												},
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
														ID: "_6db9a77f-189f-4c07-b33b-e0c88f09e0db",
													},
												},
												Outgoing: []string{"_b792d502-a08e-4a3d-bdc1-1fc33a2022f6"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											ErrorEventDefinitions: []element.ErrorEventDefinition{{}},
										},
									},
									AttachedToRef: "_614d6469-2bb8-4ad6-a20a-db5db6321c6b",
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_54d6e34f-4613-48d4-9143-b90927448adb",
													},
												},
												Outgoing: []string{"_8eac4fb5-c746-42e4-ab0b-80cdfaf76007"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											ErrorEventDefinitions: []element.ErrorEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{},
														},
													},
												},
											},
										},
									},
									AttachedToRef: "_c38139c7-a2d1-47c7-b75a-19e14c7212c8",
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_32c4138c-74ae-484a-a7e5-0609370d7080",
													},
													Name: "24 Hours",
												},
												Outgoing: []string{"_a0c99000-31e1-4de9-baf7-45c2d09ebbff"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											TimerEventDefinitions: []element.TimerEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{},
														},
													},
													TimeDate: element.ExpressionUnMarshal{
														ExpressionSubstitution: &element.Expression{},
													},
												},
											},
										},
									},
									AttachedToRef: "_e839800f-ad4f-4bcc-aaf2-d38fe4a32bcd",
								},
							},
							ServiceTasks: []element.ServiceTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_614d6469-2bb8-4ad6-a20a-db5db6321c6b",
													},
													Name: "Charge Credit Card",
												},
												Incoming: []string{"_f81560f5-e1c5-497d-a159-dce350f2283c"},
												Outgoing: []string{"_65293cb4-b859-4ff7-bcd3-ba4e38cf8763"},
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
														ID: "_8afc49f0-42c2-4da9-8e79-e08dbe349776",
													},
													Name: "Update Customer Record",
												},
												Incoming: []string{"_f072c57e-dcc6-4a16-8f28-13a2cd6e43ed"},
												Outgoing: []string{"_e396fc69-5792-41e3-8a92-ee6b83cd0909"},
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
														ID: "_88247168-b457-4663-aad0-753a0236c8df",
													},
													Name: "Offer Expired",
												},
												Incoming: []string{"_51e34f49-b9eb-4f10-ac0a-7af81e1370ff"},
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
														ID: "_7eb87eb8-0d7a-445b-b768-90d754a938ed",
													},
													Name: "Request Cancelled",
												},
												Incoming: []string{"_e396fc69-5792-41e3-8a92-ee6b83cd0909"},
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
														ID: "_42e03d0f-6c6b-4493-971f-c6928eb563b0",
													},
													Name: "Booking \nConfirmed",
												},
												Incoming: []string{"_c69565cd-3a08-4532-b466-8f9a9ebdc1e1"},
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
														ID: "_babdfa54-b55f-463f-9341-424b42db9760",
													},
													Name: "Failed Credit\nTransaction",
												},
												Incoming: []string{"_94e0bd86-6659-4669-baee-ce671c7f33d6"},
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
														ID: "_a68b0941-b7e3-4791-8e13-fd9622e4448c",
													},
													Name: "Failed \nBooking",
												},
												Incoming: []string{"_ec8c7430-02e1-4298-9b00-33c4d56d0fac"},
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
													ID: "_c38139c7-a2d1-47c7-b75a-19e14c7212c8",
												},
												Name: "Make Booking",
											},
											Incoming: []string{"_7efd61ea-979e-44f4-91c8-e212ae6b22c7"},
											Outgoing: []string{"_f81560f5-e1c5-497d-a159-dce350f2283c"},
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
																	ID: "_31a01c78-9a86-4b53-a485-e8a973ba6383",
																},
															},
															Outgoing: []string{"_d51b76e8-b096-4376-ba01-7ca63ae51f70"},
														},
													},
												},
											},
										},
										SequenceFlows: []element.SequenceFlow{
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_d51b76e8-b096-4376-ba01-7ca63ae51f70",
													},
												},
												SourceRef: "_31a01c78-9a86-4b53-a485-e8a973ba6383",
												TargetRef: "_749dd603-40f5-40fb-89b4-0e305b29892c",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_eba304c2-42e8-4d65-a9dc-e2483fd8ad86",
													},
												},
												SourceRef: "_6a68d4b4-7549-42ce-b903-9da8b2024d31",
												TargetRef: "_6ff2b954-2017-46dd-941e-4badd9326eac",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_018e5307-4c8d-46fd-ba84-7b75409b3124",
													},
												},
												SourceRef: "_749dd603-40f5-40fb-89b4-0e305b29892c",
												TargetRef: "_ea5cc55d-bfce-49c6-8a1a-a8a41a85da12",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_b08efa4c-c27c-4654-9901-279a8c456a4f",
													},
												},
												SourceRef: "_ea5cc55d-bfce-49c6-8a1a-a8a41a85da12",
												TargetRef: "_6a68d4b4-7549-42ce-b903-9da8b2024d31",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_e6011c12-800f-44e6-8ae8-51866664afcf",
													},
												},
												SourceRef: "_b595ec43-0769-4864-8f2e-403c405c8217",
												TargetRef: "_6a68d4b4-7549-42ce-b903-9da8b2024d31",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_b2d48080-6bb1-4ae9-85c9-916d87536a49",
													},
												},
												SourceRef: "_749dd603-40f5-40fb-89b4-0e305b29892c",
												TargetRef: "_b595ec43-0769-4864-8f2e-403c405c8217",
											},
										},
										ParallelGatewaies: []element.ParallelGateway{
											{
												Gateway: element.Gateway{
													FlowNode: element.FlowNode{
														FlowElement: element.FlowElement{
															BaseElement: element.BaseElement{
																ID: "_749dd603-40f5-40fb-89b4-0e305b29892c",
															},
														},
														Incoming: []string{"_d51b76e8-b096-4376-ba01-7ca63ae51f70"},
														Outgoing: []string{"_018e5307-4c8d-46fd-ba84-7b75409b3124", "_b2d48080-6bb1-4ae9-85c9-916d87536a49"},
													},
													GatewayDirection: "Diverging",
												},
											},
											{
												Gateway: element.Gateway{
													FlowNode: element.FlowNode{
														FlowElement: element.FlowElement{
															BaseElement: element.BaseElement{
																ID: "_6a68d4b4-7549-42ce-b903-9da8b2024d31",
															},
														},
														Incoming: []string{"_b08efa4c-c27c-4654-9901-279a8c456a4f", "_e6011c12-800f-44e6-8ae8-51866664afcf"},
														Outgoing: []string{"_eba304c2-42e8-4d65-a9dc-e2483fd8ad86"},
													},
													GatewayDirection: "Converging",
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
																	ID: "_3a2f133c-3ae1-4e21-94b5-6e8cf51acd74",
																},
																Name: "Cancel Hotel",
															},
														},
														IsForCompensation: true,
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
																	ID: "_b595ec43-0769-4864-8f2e-403c405c8217",
																},
																Name: "Book Hotel",
															},
															Incoming: []string{"_b2d48080-6bb1-4ae9-85c9-916d87536a49"},
															Outgoing: []string{"_e6011c12-800f-44e6-8ae8-51866664afcf"},
														},
														IsForCompensation: false,
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
																	ID: "_ea5cc55d-bfce-49c6-8a1a-a8a41a85da12",
																},
																Name: "Book Flight",
															},
															Incoming: []string{"_018e5307-4c8d-46fd-ba84-7b75409b3124"},
															Outgoing: []string{"_b08efa4c-c27c-4654-9901-279a8c456a4f"},
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
																	ID: "_0198160d-b56c-4919-9920-db5f32d16b3f",
																},
																Name: "Cancel Flight",
															},
														},
														IsForCompensation: true,
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
																	ID: "_b25ecc7c-4eff-4a70-96f2-6b2f94cf19b1",
																},
																Name: "Hotel",
															},
														},
													},
													EventDefinitions: element.EventDefinitions{
														CompensateEventDefinitions: []element.CompensateEventDefinition{
															{
																EventDefinition: element.EventDefinition{
																	RootElement: element.RootElement{
																		BaseElement: element.BaseElement{},
																	},
																},
															},
														},
													},
												},
												AttachedToRef: "_b595ec43-0769-4864-8f2e-403c405c8217",
											},
											{
												CatchEvent: element.CatchEvent{
													Event: element.Event{
														FlowNode: element.FlowNode{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "_fe3f9094-097b-416d-adeb-4b7e7e753f3c",
																},
																Name: "Flight",
															},
														},
													},
													EventDefinitions: element.EventDefinitions{
														CompensateEventDefinitions: []element.CompensateEventDefinition{
															{
																EventDefinition: element.EventDefinition{
																	RootElement: element.RootElement{
																		BaseElement: element.BaseElement{},
																	},
																},
															},
														},
													},
												},
												AttachedToRef: "_ea5cc55d-bfce-49c6-8a1a-a8a41a85da12",
											},
										},
										EndEvents: []element.EndEvent{
											{
												ThrowEvent: element.ThrowEvent{
													Event: element.Event{
														FlowNode: element.FlowNode{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "_6ff2b954-2017-46dd-941e-4badd9326eac",
																},
																Name: "Travel Booked",
															},
															Incoming: []string{"_eba304c2-42e8-4d65-a9dc-e2483fd8ad86"},
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
																ID: "_e880bf53-84ca-4776-aa75-d1bf53172240",
															},
															Name: "Handle Compensation",
														},
													},
												},
												TriggeredByEvent: true,
												FlowElements: element.FlowElements{
													StartEvents: []element.StartEvent{
														{
															CatchEvent: element.CatchEvent{
																Event: element.Event{
																	FlowNode: element.FlowNode{
																		FlowElement: element.FlowElement{
																			BaseElement: element.BaseElement{
																				ID: "_8af17ed4-6e13-463b-8333-d397b3002c65",
																			},
																			Name: "Booking",
																		},
																		Outgoing: []string{"_e0c52135-c57f-4ee6-b488-22bbab06dccc"},
																	},
																},
																EventDefinitions: element.EventDefinitions{
																	CompensateEventDefinitions: []element.CompensateEventDefinition{
																		{
																			EventDefinition: element.EventDefinition{
																				RootElement: element.RootElement{
																					BaseElement: element.BaseElement{
																						ID: "_434dc71e-6d59-4940-8e27-a911bf891cdf",
																					},
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
																	ID: "_b3d5a547-87a4-460d-b7c1-53b6506127e9",
																},
															},
															SourceRef: "_c4ddacd9-6e12-49e3-b058-e71c0af08c68",
															TargetRef: "_fc4826b1-1e63-49f6-8670-7cc8104e45ea",
														},
														{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "_8d0dadf1-129e-4ba8-b829-3e7e176bf775",
																},
															},
															SourceRef: "_99bf4db9-3616-4ed1-a0f8-b8175c3fd46f",
															TargetRef: "_c4ddacd9-6e12-49e3-b058-e71c0af08c68",
														},
														{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "_b62ef735-9d04-494a-aeb4-95fd4862cc45",
																},
															},
															SourceRef: "_e4b9fa74-efd8-409f-a2e4-ad917df767b4",
															TargetRef: "_c4ddacd9-6e12-49e3-b058-e71c0af08c68",
														},
														{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "_64d60a4e-90cf-4a9a-afbb-631739b6b704",
																},
															},
															SourceRef: "_08cd4ce0-bf7c-4a05-bacf-546daa51ce9c",
															TargetRef: "_99bf4db9-3616-4ed1-a0f8-b8175c3fd46f",
														},
														{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "_3cd5068b-6df1-4a62-a06a-369d063f6160",
																},
															},
															SourceRef: "_08cd4ce0-bf7c-4a05-bacf-546daa51ce9c",
															TargetRef: "_e4b9fa74-efd8-409f-a2e4-ad917df767b4",
														},
														{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "_e0c52135-c57f-4ee6-b488-22bbab06dccc",
																},
															},
															SourceRef: "_8af17ed4-6e13-463b-8333-d397b3002c65",
															TargetRef: "_08cd4ce0-bf7c-4a05-bacf-546daa51ce9c",
														},
													},
													ParallelGatewaies: []element.ParallelGateway{
														{
															Gateway: element.Gateway{
																FlowNode: element.FlowNode{
																	FlowElement: element.FlowElement{
																		BaseElement: element.BaseElement{
																			ID: "_c4ddacd9-6e12-49e3-b058-e71c0af08c68",
																		},
																	},
																	Incoming: []string{"_8d0dadf1-129e-4ba8-b829-3e7e176bf775", "_b62ef735-9d04-494a-aeb4-95fd4862cc45"},
																	Outgoing: []string{"_b3d5a547-87a4-460d-b7c1-53b6506127e9"},
																},
																GatewayDirection: "Converging",
															},
														},
														{
															Gateway: element.Gateway{
																FlowNode: element.FlowNode{
																	FlowElement: element.FlowElement{
																		BaseElement: element.BaseElement{
																			ID: "_08cd4ce0-bf7c-4a05-bacf-546daa51ce9c",
																		},
																	},
																	Incoming: []string{"_e0c52135-c57f-4ee6-b488-22bbab06dccc"},
																	Outgoing: []string{"_64d60a4e-90cf-4a9a-afbb-631739b6b704", "_3cd5068b-6df1-4a62-a06a-369d063f6160"},
																},
																GatewayDirection: "Diverging",
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
																				ID: "_99bf4db9-3616-4ed1-a0f8-b8175c3fd46f",
																			},
																			Name: "Hotel",
																		},
																		Incoming: []string{"_64d60a4e-90cf-4a9a-afbb-631739b6b704"},
																		Outgoing: []string{"_8d0dadf1-129e-4ba8-b829-3e7e176bf775"},
																	},
																},
																EventDefinitions: element.EventDefinitions{
																	CompensateEventDefinitions: []element.CompensateEventDefinition{
																		{
																			EventDefinition: element.EventDefinition{
																				RootElement: element.RootElement{
																					BaseElement: element.BaseElement{},
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
																				ID: "_e4b9fa74-efd8-409f-a2e4-ad917df767b4",
																			},
																			Name: "Flight",
																		},
																		Incoming: []string{"_3cd5068b-6df1-4a62-a06a-369d063f6160"},
																		Outgoing: []string{"_b62ef735-9d04-494a-aeb4-95fd4862cc45"},
																	},
																},
																EventDefinitions: element.EventDefinitions{
																	CompensateEventDefinitions: []element.CompensateEventDefinition{
																		{
																			EventDefinition: element.EventDefinition{
																				RootElement: element.RootElement{
																					BaseElement: element.BaseElement{},
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
																				ID: "_fc4826b1-1e63-49f6-8670-7cc8104e45ea",
																			},
																		},
																		Incoming: []string{"_b3d5a547-87a4-460d-b7c1-53b6506127e9"},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									Artifacts: element.Artifacts{
										Associations: []element.Association{
											{
												Artifact: element.Artifact{
													BaseElement: element.BaseElement{
														ID: "_651344ad-d784-4ef2-9655-4bc6393ac323",
													},
												},
												AssociationDirection: "One",
												SourceRef:            "_b25ecc7c-4eff-4a70-96f2-6b2f94cf19b1",
												TargetRef:            "_3a2f133c-3ae1-4e21-94b5-6e8cf51acd74",
											},
											{
												Artifact: element.Artifact{
													BaseElement: element.BaseElement{
														ID: "_98b2bd26-e113-4d19-b6a8-bcda781a823b",
													},
												},
												AssociationDirection: "One",
												SourceRef:            "_fe3f9094-097b-416d-adeb-4b7e7e753f3c",
												TargetRef:            "_0198160d-b56c-4919-9920-db5f32d16b3f",
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

	if diff := cmp.Diff(excepted, modelelement); diff != "" {
		t.Errorf("Mismatch in %s (-want +got):\n%s", path, diff)
	}

}
