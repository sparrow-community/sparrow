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

func TestC_8_1_roundtrip(t *testing.T) {
	// create test use ./test/C.8.1-roundtrip.bpmn
	path := "./test/C.8.1-roundtrip.bpmn"
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
			ID:                 "_d4aecb6e-8641-4d7a-af45-dcfae1d639ea",
			Name:               "Vacation Request",
			TargetNamespace:    "http://www.trisotech.com/definitions/_d4aecb6e-8641-4d7a-af45-dcfae1d639ea",
			ExpressionLanguage: "https://www.omg.org/spec/DMN/20191111/FEEL/",
			TypeLanguage:       "http://www.trisotech.com/2015/triso/modeling/ItemDefinitionType",
			Exporter:           "Workflow Modeler",
			ExporterVersion:    "8.0.0",
			RootElemnts: element.RootElemnts{
				ItemDefinitions: []element.ItemDefinition{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_triso-default-bpmnItemDefinition-string_id",
							},
						},
						StructureRef: "feel:string",
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_24ce0503-1c7a-41eb-bbfd-28f533e0268f",
							},
						},
						StructureRef: "feel:number",
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_7d2aae2a-a2b3-4927-8865-154dcb88e7f1",
							},
						},
						StructureRef: "feel:date",
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_e9f96a27-1ff4-4bc0-a846-9c244d08cae7",
							},
						},
						StructureRef: "feel:number",
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_bb207a19-f088-40e0-8f0a-f151c8002cf5",
							},
						},
						StructureRef: "feel:string",
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_356dbbe9-8dcf-436a-a5c9-26d618dd4b5a",
							},
						},
						StructureRef: "feel:string",
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_ae0f56cf-0cb4-43ec-8c77-16d0ba7b9642",
							},
						},
						StructureRef: "feel:string",
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_6a883e1e-df2f-42e3-bc04-d122469b52a3",
							},
						},
						StructureRef: "feel:string",
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_f6564795-5df9-4a0d-941e-72b937fcb877",
							},
						},
						StructureRef: "feel:number",
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_bbda4f78-e54c-434c-929c-81fbf14af93d",
							},
						},
						StructureRef: "feel:string",
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_68f999b9-915a-40d9-b3b0-78fd2b0d6b39",
							},
						},
						StructureRef: "simon.1.simon.vacation",
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_55a74c87-4810-4c82-a5c1-83cbfaf0d016",
							},
						},
						StructureRef: "simon.1.simon.cvacation",
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_9f8358c3-1ac6-47cf-a768-6e90d4b81e7f",
							},
						},
						StructureRef: "simon.1.simon.input_vacation",
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_35aa1a98-77be-4b78-a590-ba5bafe10a58",
							},
						},
						StructureRef: "ApprovalStatus",
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_52182389-d104-457d-b152-d312b38991f3",
							},
						},
						StructureRef: "Vacation_Approval",
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_8a28ca0b-0650-47a2-a1d1-a1c1b72a0997",
							},
						},
						StructureRef: "tEmployeeInformation",
					},
				},
				Errors: []element.Error{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_bb1e92be-5879-4fbe-8efb-4f0c6ea64187",
							},
						},
						Name:      "Not Found",
						ErrorCode: "404",
					},
				},
				Messages: []element.Message{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID:                "_b4e85145-0a59-4799-bd6c-18a5cc2b041a",
								ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil}},
							},
						},
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_322ee005-6121-4712-ba90-909f359b9b5b",
							},
						},
						Name:    "vacation",
						ItemRef: "_68f999b9-915a-40d9-b3b0-78fd2b0d6b39",
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID:                "_712b7cfb-332e-4df5-87ac-f5a44560c009",
								ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
							},
						},
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID:                "_fc36570a-8e3d-4802-8a9a-09263799c850",
								ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil, nil, nil, nil}},
							},
						},
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_70e2fb1b-c19f-4c53-810c-46aaba76b4dc",
							},
						},
						Name:    "vacation",
						ItemRef: "_55a74c87-4810-4c82-a5c1-83cbfaf0d016",
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_42a34a80-8aa4-445c-820e-974f9b759a82",
							},
						},
						Name:    "OData-EntityId",
						ItemRef: "_triso-default-bpmnItemDefinition-string_id",
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID:                "_9f3c81a6-53b8-4809-b771-badf80c3777c",
								ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
							},
						},
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID:                "_48402d3f-cbaf-4d40-aaa7-c7a4424f26f7",
								ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
							},
						},
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID:                "_298108a4-f565-4023-b8a8-20b0d0b4345a",
								ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
							},
						},
					},
				},
				Collaborations: []element.Collaboration{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "Collaboration_1rz14um",
							},
						},
						Participants: []element.Participant{
							{
								BaseElement: element.BaseElement{
									ID:                "Participant_04dpvds",
									ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
								},
								Name:       "Vacation Request",
								ProcessRef: "VacationRequestProcess",
							},
						},
					},
				},
				Processes: []element.Process{
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID:                "VacationRequestProcess",
									Documentation:     []element.Documentation{{Value: "Vacation Request - BPMN MIWG demo for 2022"}},
									ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
								},
							},
							Name: "Vacation Request",
							IoSpecification: element.IoSpecification{
								InputOutputSpecification: element.InputOutputSpecification{
									DataInputs: []element.DataInput{
										{
											BaseElement: element.BaseElement{
												ID: "_8b9aa28f-5974-4087-9895-0467c25635dc",
											},
											Name:           "Employee Badge Number",
											ItemSubjectRef: "_24ce0503-1c7a-41eb-bbfd-28f533e0268f",
										},
										{
											BaseElement: element.BaseElement{
												ID: "_2bed4da9-f468-4fe0-a37c-ae5d4c17e323",
											},
											Name:           "From",
											ItemSubjectRef: "_7d2aae2a-a2b3-4927-8865-154dcb88e7f1",
										},
										{
											BaseElement: element.BaseElement{
												ID: "_103fd816-b422-420d-aec9-c28b54f620db",
											},
											Name:           "To",
											ItemSubjectRef: "_7d2aae2a-a2b3-4927-8865-154dcb88e7f1",
										},
									},
									DataOutputs: []element.DataOutput{
										{
											BaseElement: element.BaseElement{
												ID: "_c2d108b2-e6c0-47e6-a3cb-aa5321bcaeaf",
											},
											Name:           "Vacation Approval",
											ItemSubjectRef: "_35aa1a98-77be-4b78-a590-ba5bafe10a58",
										},
										{
											BaseElement: element.BaseElement{
												ID: "_dac8ee76-f637-4fd9-8357-6a87fd11ef41",
											},
											Name:           "Reason",
											ItemSubjectRef: "_triso-default-bpmnItemDefinition-string_id",
										},
									},
									InputSets: []element.InputSet{
										{
											BaseElement: element.BaseElement{
												ID: "_8dfe6186-8dd8-478c-b094-4851e2a5d8d7",
											},
											DataInputRefs: []string{"_8b9aa28f-5974-4087-9895-0467c25635dc", "_2bed4da9-f468-4fe0-a37c-ae5d4c17e323", "_103fd816-b422-420d-aec9-c28b54f620db"},
										},
									},
									OutputSets: []element.OutputSet{
										{
											BaseElement: element.BaseElement{
												ID: "_bd2a87b8-c67a-4fb4-aa46-39c677ba9b09",
											},
											DataOutputRefs: []string{"_c2d108b2-e6c0-47e6-a3cb-aa5321bcaeaf", "_dac8ee76-f637-4fd9-8357-6a87fd11ef41"},
										},
									},
								},
							},
						},
						ProcessType:  "",
						IsExecutable: true,
						FlowElements: element.FlowElements{
							StartEvents: []element.StartEvent{
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID:                "_b1625a52-aaf0-4694-86cb-7af891212ac6",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil, nil}},
													},
													Name: "Vacation Request Received",
												},
												Outgoing: []string{"_b02b5bfa-5629-4af5-b1bc-ac2a86d88adb"},
											},
										},
									},
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_b02b5bfa-5629-4af5-b1bc-ac2a86d88adb",
										},
									},
									SourceRef: "_b1625a52-aaf0-4694-86cb-7af891212ac6",
									TargetRef: "_2b960d84-feb1-46a9-a1a1-c300dd996b99",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_ac1d82df-9bf6-404b-be70-b052e311f8b7",
										},
									},
									SourceRef: "_1a818a94-ba6f-413b-a7e8-6f8fd2a11e32",
									TargetRef: "_42367c5f-d084-44ee-90c7-960d1ab02a3b",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_2a5ee73c-c495-48a0-945b-3dabab28ca91",
										},
										Name: "Refused",
									},
									SourceRef: "_42367c5f-d084-44ee-90c7-960d1ab02a3b",
									TargetRef: "_9ed61a6a-7cc1-4ed1-86d8-03482b0983c9",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_130d903a-f211-4b3e-9711-9061f5978d94",
										},
									},
									SourceRef: "_9ed61a6a-7cc1-4ed1-86d8-03482b0983c9",
									TargetRef: "_1688f604-5edf-4187-ad9e-18f74bfede53",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_0a1c4f20-509f-4aeb-baf9-acc762f4fdf9",
										},
										Name: "Manual Validation Required",
									},
									SourceRef: "_42367c5f-d084-44ee-90c7-960d1ab02a3b",
									TargetRef: "_79523269-7444-4b01-90e9-e23957a9d020",
									ConditionExpression: element.ExpressionUnMarshal{
										Type: "tFormalExpression",
										ExpressionSubstitution: &element.FormalExpression{
											Expression: element.Expression{
												BaseElementWithMixedContent: element.BaseElementWithMixedContent{
													ID: "",
												},
											},
											Value:    `Vacation Approval = "Manual Validation Required"`,
											Language: "https://www.omg.org/spec/DMN/20191111/FEEL/",
										},
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_325973e7-0bc8-4136-b6df-be1e681d8608",
										},
										Name: "Approved",
									},
									SourceRef: "_42367c5f-d084-44ee-90c7-960d1ab02a3b",
									TargetRef: "_93ec9873-edf1-4549-b052-961994ec8234",
									ConditionExpression: element.ExpressionUnMarshal{
										Type: "tFormalExpression",
										ExpressionSubstitution: &element.FormalExpression{
											Expression: element.Expression{
												BaseElementWithMixedContent: element.BaseElementWithMixedContent{
													ID: "",
												},
											},
											Value:    `Vacation Approval = "Approved"`,
											Language: "https://www.omg.org/spec/DMN/20191111/FEEL/",
										},
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_188d38f6-b73d-4c62-964b-affc6fe040c1",
										},
									},
									SourceRef: "_93ec9873-edf1-4549-b052-961994ec8234",
									TargetRef: "_4b72053b-8ebb-4ae6-99c6-7c93cf1c1d1b",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_3c50d7bb-5ecc-4d34-82a9-a5fa41d9c018",
										},
										Name: "",
									},
									SourceRef: "_79523269-7444-4b01-90e9-e23957a9d020",
									TargetRef: "_64bb8b55-d348-41d6-9ea8-f333b1d8cb69",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_f2b0da63-d841-4457-ad85-7d86c8b5c1d2",
										},
										Name: "Approved",
									},
									SourceRef: "_64bb8b55-d348-41d6-9ea8-f333b1d8cb69",
									TargetRef: "_a97c1a48-faba-447b-bfa6-7aa81a6fe0a0",
									ConditionExpression: element.ExpressionUnMarshal{
										Type: "tFormalExpression",
										ExpressionSubstitution: &element.FormalExpression{
											Expression: element.Expression{
												BaseElementWithMixedContent: element.BaseElementWithMixedContent{
													ID: "",
												},
											},
											Value:    `Vacation Approval = "Approved"`,
											Language: "https://www.omg.org/spec/DMN/20191111/FEEL/",
										},
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_d3a575d8-cd41-4246-ad7a-7a068de6679f",
										},
										Name: "Refused",
									},
									SourceRef: "_64bb8b55-d348-41d6-9ea8-f333b1d8cb69",
									TargetRef: "_02232e32-c3d2-473c-a15d-9c5dca00eadc",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_93ac178c-3efc-48ae-8a99-8d29a0c4aaf6",
										},
									},
									SourceRef: "_02232e32-c3d2-473c-a15d-9c5dca00eadc",
									TargetRef: "_3ae826ca-5f38-43c4-be3a-35d1157aa27f",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_68ef27a7-011e-4ec0-ad45-41f52093de68",
										},
										Name: "",
									},
									SourceRef: "_4b72053b-8ebb-4ae6-99c6-7c93cf1c1d1b",
									TargetRef: "_6677ef80-82df-4951-919d-1f36123b681b",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_a5441c65-06f8-4972-8158-e6f61f904841",
										},
									},
									SourceRef: "_a97c1a48-faba-447b-bfa6-7aa81a6fe0a0",
									TargetRef: "_5e16a4e0-0f23-482a-be47-d3edbc5741ba",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_81ed1b89-5a4d-4028-a245-4fd38b866b6d",
										},
									},
									SourceRef: "_5e16a4e0-0f23-482a-be47-d3edbc5741ba",
									TargetRef: "_1cd5fe29-b3ec-4f21-a1aa-57773f0729ca",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_f250e125-23eb-45cd-9c61-191d172df093",
										},
									},
									SourceRef: "_f8fcb377-3d7d-4138-9a7e-6ab58b97e29d",
									TargetRef: "_b4d636eb-b501-4462-93c8-04652db10307",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_d79c991e-446c-47d1-ac9d-9d0113e35b93",
										},
										Name: "",
									},
									SourceRef: "_2b960d84-feb1-46a9-a1a1-c300dd996b99",
									TargetRef: "_1a818a94-ba6f-413b-a7e8-6f8fd2a11e32",
								},
							},
							ServiceTasks: []element.ServiceTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_2b960d84-feb1-46a9-a1a1-c300dd996b99",
														ExtensionElements: element.ExtensionElements{
															Any: []xml.Token{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
														},
														Any: nil,
													},
													Name: "Fetch Vacation Information",
												},
												Incoming: []string{"_b02b5bfa-5629-4af5-b1bc-ac2a86d88adb"},
												Outgoing: []string{"_d79c991e-446c-47d1-ac9d-9d0113e35b93"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_8fd43f87-3dd6-47b2-bbe9-39d66e7fe316",
															},
															Name:           "Fixed input",
															ItemSubjectRef: "_24ce0503-1c7a-41eb-bbfd-28f533e0268f",
														},
														{
															BaseElement: element.BaseElement{ID: "DI_Property_10qeslx"},
															Name:        "__targetRef_placeholder",
														},
													},
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_18249f3e-0eeb-4c35-adb2-af035a694a83",
															},
															Name:           "Fixed output",
															ItemSubjectRef: "_68f999b9-915a-40d9-b3b0-78fd2b0d6b39",
														},
														{
															BaseElement: element.BaseElement{ID: "DO_Property_10qeslx"},
															Name:        "__targetRef_placeholder",
														},
													},
													InputSets: []element.InputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_9eed8142-006d-46af-97d7-abdec861abcd",
															},
															DataInputRefs: []string{"_8fd43f87-3dd6-47b2-bbe9-39d66e7fe316", "DI_Property_10qeslx"},
														},
													},
													OutputSets: []element.OutputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_d3f65368-fb14-4eb5-ad89-71af07d4913a",
															},
															DataOutputRefs: []string{"_18249f3e-0eeb-4c35-adb2-af035a694a83", "DO_Property_10qeslx"},
														},
													},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_f4f41761-cdd8-4b90-a956-9182345df488",
														},
														SourceRef: "_8b9aa28f-5974-4087-9895-0467c25635dc",
														TargetRef: "_8fd43f87-3dd6-47b2-bbe9-39d66e7fe316",
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_40d3cb58-31bb-47a4-9591-032a38011de3",
														},
														SourceRef: "_18249f3e-0eeb-4c35-adb2-af035a694a83",
														TargetRef: "_831726fb-e5ba-40aa-b484-1fc65de4980a",
														Transformation: element.FormalExpression{
															Language: "https://www.omg.org/spec/DMN/20191111/MODEL/",
															Value:    `<ns1:literalExpression xmlns:ns1="https://www.omg.org/spec/DMN/20191111/MODEL/" id="_086343dd-32b0-4cd3-9730-9964c61de8d5" expressionLanguage="https://www.omg.org/spec/DMN/20191111/FEEL/"><ns1:text>vacation.value[1]</ns1:text></ns1:literalExpression>`,
														},
													},
												},
											},
										},
									},
									Implementation: "kp:http",
									OperationRef:   "_dd434d56-40a4-4d63-8a22-58d8e5a5edc7",
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_4b72053b-8ebb-4ae6-99c6-7c93cf1c1d1b",
														ExtensionElements: element.ExtensionElements{
															Any: []xml.Token{nil, nil, nil, nil, nil, nil, nil, nil, nil},
														},
													},
													Name: "Update Remaining Vacation",
												},
												Incoming: []string{"_188d38f6-b73d-4c62-964b-affc6fe040c1"},
												Outgoing: []string{"_68ef27a7-011e-4ec0-ad45-41f52093de68"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_803ff266-47ce-45ea-b3e6-019916eddc71",
															},
															Name:           "Fixed input",
															ItemSubjectRef: "_24ce0503-1c7a-41eb-bbfd-28f533e0268f",
														},
													},
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_3eb8a29c-a29f-4cf5-bab6-95b374f1b183",
															},
															Name: "Fixed output",
														},
													},
													InputSets: []element.InputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_a129c4a1-10ff-4835-b36f-392feb105ebc",
															},
															DataInputRefs: []string{"_803ff266-47ce-45ea-b3e6-019916eddc71"},
														},
													},
													OutputSets: []element.OutputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_cc1b25bf-8bc3-4d4c-af06-df2d6c9f7eef",
															},
															DataOutputRefs: []string{"_3eb8a29c-a29f-4cf5-bab6-95b374f1b183"},
														},
													},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_8184912c-03a4-4c1c-a478-13b581231aa2",
														},
														SourceRef: "_8b9aa28f-5974-4087-9895-0467c25635dc",
														TargetRef: "_803ff266-47ce-45ea-b3e6-019916eddc71",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_2c85e650-0096-4780-b6fe-f52a9654af77",
														},
														SourceRef: "_775c93ab-82e5-4a62-97f0-d305eda76c92",
														TargetRef: "_803ff266-47ce-45ea-b3e6-019916eddc71",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_2c7c930e-5b1c-4792-ba24-e4f80b6b1369",
														},
														SourceRef: "_103fd816-b422-420d-aec9-c28b54f620db",
														TargetRef: "_803ff266-47ce-45ea-b3e6-019916eddc71",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_c5c828aa-9a42-4e96-8266-7b42fad69511",
														},
														SourceRef: "_2bed4da9-f468-4fe0-a37c-ae5d4c17e323",
														TargetRef: "_803ff266-47ce-45ea-b3e6-019916eddc71",
													},
												},
											},
										},
									},
									Implementation: "kp:http",
									OperationRef:   "_a0e96e19-a01f-4dde-b613-0b774e8d350b",
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_5e16a4e0-0f23-482a-be47-d3edbc5741ba",
														ExtensionElements: element.ExtensionElements{
															Any: []xml.Token{nil, nil, nil, nil, nil, nil, nil, nil, nil},
														},
													},
													Name: "Update Remaining Vacation",
												},
												Incoming: []string{"_a5441c65-06f8-4972-8158-e6f61f904841"},
												Outgoing: []string{"_81ed1b89-5a4d-4028-a245-4fd38b866b6d"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_751edb63-a514-44c3-a7e0-7ce5c18f09ec",
															},
															Name: "Fixed input",
														},
													},
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_8d84c387-82d2-4d75-ac08-d9f9b5e430c8",
															},
															Name: "Fixed output",
														},
													},
													InputSets: []element.InputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_fc1c7175-8c6b-460a-bb71-c52a9874e7cb",
															},
															DataInputRefs: []string{"_751edb63-a514-44c3-a7e0-7ce5c18f09ec"},
														},
													},
													OutputSets: []element.OutputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_2b9361e2-8253-481a-b773-c78d9dc59cda",
															},
															DataOutputRefs: []string{"_8d84c387-82d2-4d75-ac08-d9f9b5e430c8"},
														},
													},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_95e48019-8c07-4005-a706-2c4b421bed60",
														},
														SourceRef: "_103fd816-b422-420d-aec9-c28b54f620db",
														TargetRef: "_751edb63-a514-44c3-a7e0-7ce5c18f09ec",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_3c85e650-0096-4780-b6fe-f52a9654af78",
														},
														SourceRef: "_775c93ab-82e5-4a62-97f0-d305eda76c92",
														TargetRef: "_751edb63-a514-44c3-a7e0-7ce5c18f09ec",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_6d535768-1e33-4054-bb8c-b141a0098316",
														},
														SourceRef: "_2bed4da9-f468-4fe0-a37c-ae5d4c17e323",
														TargetRef: "_751edb63-a514-44c3-a7e0-7ce5c18f09ec",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_67897da3-a828-4759-bc67-ecdda5e1c550",
														},
														SourceRef: "_8b9aa28f-5974-4087-9895-0467c25635dc",
														TargetRef: "_751edb63-a514-44c3-a7e0-7ce5c18f09ec",
													},
												},
											},
										},
									},
									Implementation: "kp:http",
									OperationRef:   "_a0e96e19-a01f-4dde-b613-0b774e8d350b",
								},
							},
							BoundaryEvents: []element.BoundaryEvent{
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_f8fcb377-3d7d-4138-9a7e-6ab58b97e29d",
													},
												},
												Outgoing: []string{"_f250e125-23eb-45cd-9c61-191d172df093"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											ErrorEventDefinitions: []element.ErrorEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "_2a8bacf9-ed42-4b2c-805e-c15c1ea325af",
															},
														},
													},
													ErrorRef: "_bb1e92be-5879-4fbe-8efb-4f0c6ea64187",
												},
											},
										},
									},
									CancelActivity: true,
									AttachedToRef: "_2b960d84-feb1-46a9-a1a1-c300dd996b99",
								},
							},
							DataObjectReferenes: []element.DataObjectReference{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_831726fb-e5ba-40aa-b484-1fc65de4980a",
										},
										Name: "Current Vacation Status",
									},
									DataObjectRef: "_775c93ab-82e5-4a62-97f0-d305eda76c92",
								},
							},
							DataObjects: []element.DataObject{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_775c93ab-82e5-4a62-97f0-d305eda76c92",
										},
										Name: "Current Vacation Status",
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_4234afd4-3ab0-4e09-8a7c-b475458db18a",
										},
										Name: "Reason",
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "vacationApproved",
										},
										Name: "vacationApproved",
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "vacationRefused",
										},
										Name: "vacationRefused",
									},
								},
							},
							BusinessRuleTasks: []element.BusinessRuleTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID:                "_1a818a94-ba6f-413b-a7e8-6f8fd2a11e32",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil, nil, nil, nil}},
													},
													Name: "Vacation Approval",
												},
												Incoming: []string{"_d79c991e-446c-47d1-ac9d-9d0113e35b93"},
												Outgoing: []string{"_ac1d82df-9bf6-404b-be70-b052e311f8b7"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID:                "_c86de5bd-bc4a-49e1-ad80-ae681712f235",
																ExtensionElements: element.ExtensionElements{},
															},
															Name:           "Days Remaining",
															ItemSubjectRef: "_24ce0503-1c7a-41eb-bbfd-28f533e0268f",
														},
														{
															BaseElement: element.BaseElement{
																ID:                "_06570edb-e883-409d-a645-4392e91f5cdd",
																ExtensionElements: element.ExtensionElements{},
															},
															Name:           "From",
															ItemSubjectRef: "_7d2aae2a-a2b3-4927-8865-154dcb88e7f1",
														},
														{
															BaseElement: element.BaseElement{
																ID: "_efb47006-64e5-42cd-beb2-af36354e3875",
															},
															Name:           "To",
															ItemSubjectRef: "_7d2aae2a-a2b3-4927-8865-154dcb88e7f1",
														},
														{
															BaseElement: element.BaseElement{
																ID: "DI_Property_1jfswd4",
															},
															Name: "__targetRef_placeholder",
														},
													},
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_60979891-0a46-4a66-9043-492cd159a227",
															},
															Name:           "Vacation Approval Status",
															ItemSubjectRef: "_52182389-d104-457d-b152-d312b38991f3",
														},
														{
															BaseElement: element.BaseElement{
																ID: "DO_Property_1jfswd4",
															},
															Name: "__targetRef_placeholder",
														},
													},
													InputSets: []element.InputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_b7c4db9c-d212-4893-b468-e4c275f868bb",
															},
															DataInputRefs: []string{"_c86de5bd-bc4a-49e1-ad80-ae681712f235", "_06570edb-e883-409d-a645-4392e91f5cdd", "_efb47006-64e5-42cd-beb2-af36354e3875", "DI_Property_1jfswd4"},
														},
													},
													OutputSets: []element.OutputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_ebc20a5a-c487-4bd7-9d14-588dbd28961a",
															},
															DataOutputRefs: []string{"_60979891-0a46-4a66-9043-492cd159a227", "DO_Property_1jfswd4"},
														},
													},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_9425c610-3ef7-4a5c-bc39-e4242bd1d1a1",
														},
														SourceRef: "_775c93ab-82e5-4a62-97f0-d305eda76c92",
														TargetRef: "_c86de5bd-bc4a-49e1-ad80-ae681712f235",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_d7aa7a34-bb93-49ae-87d7-f3ccd25add3f",
														},
														SourceRef: "_2bed4da9-f468-4fe0-a37c-ae5d4c17e323",
														TargetRef: "_06570edb-e883-409d-a645-4392e91f5cdd",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_bfe6b00e-a1ce-4b12-b868-f445ebc77d1d",
														},
														SourceRef: "_103fd816-b422-420d-aec9-c28b54f620db",
														TargetRef: "_efb47006-64e5-42cd-beb2-af36354e3875",
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_3dcf2a50-5a30-4a47-956f-7962eb907737",
														},
														SourceRef: "_60979891-0a46-4a66-9043-492cd159a227",
														TargetRef: "_c2d108b2-e6c0-47e6-a3cb-aa5321bcaeaf",
														Transformation: element.FormalExpression{
															Language: "https://www.omg.org/spec/DMN/20191111/MODEL/",
															Value:    `<ns1:literalExpression xmlns:ns1="https://www.omg.org/spec/DMN/20191111/MODEL/" id="_06eb6996-5a49-4de7-a6f5-ac4f72dc253e" expressionLanguage="https://www.omg.org/spec/DMN/20191111/FEEL/"><ns1:text>Vacation Approval Status.Approval</ns1:text></ns1:literalExpression>`,
														},
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_66f6f24d-1cad-4c93-b9b3-b93193015a18",
														},
														SourceRef: "_60979891-0a46-4a66-9043-492cd159a227",
														TargetRef: "_dac8ee76-f637-4fd9-8357-6a87fd11ef41",
														Transformation: element.FormalExpression{
															Language: "https://www.omg.org/spec/DMN/20191111/MODEL/",
															Value:    `<ns1:literalExpression xmlns:ns1="https://www.omg.org/spec/DMN/20191111/MODEL/" id="_1ccc4a25-aca7-475b-bc42-f710be757403" expressionLanguage="https://www.omg.org/spec/DMN/20191111/FEEL/"><ns1:text>Vacation Approval Status.Reason</ns1:text></ns1:literalExpression>`,
														},
													},
												},
											},
										},
									},
									Implementation: "kp:dmn",
								},
							},
							ExclusiveGatewaies: []element.ExclusiveGateway{
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID:                "_42367c5f-d084-44ee-90c7-960d1ab02a3b",
													ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
												},
												Name: "",
											},
											Incoming: []string{"_ac1d82df-9bf6-404b-be70-b052e311f8b7"},
											Outgoing: []string{"_2a5ee73c-c495-48a0-945b-3dabab28ca91", "_0a1c4f20-509f-4aeb-baf9-acc762f4fdf9", "_325973e7-0bc8-4136-b6df-be1e681d8608"},
										},
										GatewayDirection: "Diverging",
									},
									Default: "_2a5ee73c-c495-48a0-945b-3dabab28ca91",
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID:                "_64bb8b55-d348-41d6-9ea8-f333b1d8cb69",
													ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
												},
												Name: "",
											},
											Incoming: []string{"_3c50d7bb-5ecc-4d34-82a9-a5fa41d9c018"},
											Outgoing: []string{"_f2b0da63-d841-4457-ad85-7d86c8b5c1d2", "_d3a575d8-cd41-4246-ad7a-7a068de6679f"},
										},
										GatewayDirection: "Diverging",
									},
									Default: "_d3a575d8-cd41-4246-ad7a-7a068de6679f",
								},
							},
							SendTasks: []element.SendTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID:                "_9ed61a6a-7cc1-4ed1-86d8-03482b0983c9",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil, nil, nil, nil, nil}},
													},
													Name: "Notify Employee of Refusal",
												},
												Incoming: []string{"_2a5ee73c-c495-48a0-945b-3dabab28ca91"},
												Outgoing: []string{"_130d903a-f211-4b3e-9711-9061f5978d94"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_3664e2fd-3dbd-41f2-a931-ff00709d9f14",
															},
															Name: "Fixed input",
														},
													},
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_02243091-1c47-49e2-a23f-405ab423eb36",
															},
															Name: "Fixed output",
														},
													},
													InputSets: []element.InputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_ca96a978-8a8d-48b5-9876-594b8e1be75a",
															},
															DataInputRefs: []string{"_3664e2fd-3dbd-41f2-a931-ff00709d9f14"},
														},
													},
													OutputSets: []element.OutputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_132afdae-6a94-4af3-8227-788f99c5d117",
															},
															DataOutputRefs: []string{"_02243091-1c47-49e2-a23f-405ab423eb36"},
														},
													},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_3f72d33b-e9cb-4dce-b97f-34b1975ba494",
														},
														SourceRef: "_775c93ab-82e5-4a62-97f0-d305eda76c92",
														TargetRef: "_3664e2fd-3dbd-41f2-a931-ff00709d9f14",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_6ae63a82-1bc1-408f-95a3-af645633672a",
														},
														SourceRef: "_2bed4da9-f468-4fe0-a37c-ae5d4c17e323",
														TargetRef: "_3664e2fd-3dbd-41f2-a931-ff00709d9f14",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_7e4b1382-d551-4406-979d-c586b5925e64",
														},
														SourceRef: "_103fd816-b422-420d-aec9-c28b54f620db",
														TargetRef: "_3664e2fd-3dbd-41f2-a931-ff00709d9f14",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_c5b9c7c0-35dc-4efb-93bd-a019fe7d0f5b",
														},
														SourceRef: "_4234afd4-3ab0-4e09-8a7c-b475458db18a",
														TargetRef: "_3664e2fd-3dbd-41f2-a931-ff00709d9f14",
													},
												},
											},
										},
									},
									Implementation: "kp:mail",
									OperationRef:   "_83f1c680-7e06-4a9a-9206-396ed8155a71",
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID:                "_93ec9873-edf1-4549-b052-961994ec8234",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil, nil, nil, nil, nil}},
													},
													Name: "Notify Employee of Approval",
												},
												Incoming: []string{"_325973e7-0bc8-4136-b6df-be1e681d8608"},
												Outgoing: []string{"_188d38f6-b73d-4c62-964b-affc6fe040c1"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_1c26901a-059c-49ee-9458-6f9f0db67197",
															},
															Name: "Fixed input",
														},
													},
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_98d35e58-e638-411a-bf53-9b689a6d1170",
															},
															Name: "Fixed output",
														},
													},
													InputSets: []element.InputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_f1ca5d59-ff49-475b-adcf-88c191d75997",
															},
															DataInputRefs: []string{"_1c26901a-059c-49ee-9458-6f9f0db67197"},
														},
													},
													OutputSets: []element.OutputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_cbe32ab6-5278-4256-a8ab-b27d201a6ad9",
															},
															DataOutputRefs: []string{"_98d35e58-e638-411a-bf53-9b689a6d1170"},
														},
													},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_32df05cc-b9ef-417d-9eca-34c98a381182",
														},
														SourceRef: "_775c93ab-82e5-4a62-97f0-d305eda76c92",
														TargetRef: "_1c26901a-059c-49ee-9458-6f9f0db67197",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_d2df05cc-b9ef-417d-9eca-34c98a381183",
														},
														SourceRef: "_2bed4da9-f468-4fe0-a37c-ae5d4c17e323",
														TargetRef: "_1c26901a-059c-49ee-9458-6f9f0db67197",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_7e7c0711-6a8c-4f30-be77-85f15ea4903c",
														},
														SourceRef: "_103fd816-b422-420d-aec9-c28b54f620db",
														TargetRef: "_1c26901a-059c-49ee-9458-6f9f0db67197",
													},
												},
											},
										},
									},
									Implementation: "kp:mail",
									OperationRef:   "_83f1c680-7e06-4a9a-9206-396ed8155a71",
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID:                "_a97c1a48-faba-447b-bfa6-7aa81a6fe0a0",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil, nil, nil, nil, nil}},
													},
													Name: "Notify Employee of Approval",
												},
												Incoming: []string{"_f2b0da63-d841-4457-ad85-7d86c8b5c1d2"},
												Outgoing: []string{"_a5441c65-06f8-4972-8158-e6f61f904841"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_871b6349-9bc5-440b-9e01-6e3cf700f9e6",
															},
															Name:           "From",
															ItemSubjectRef: "_7d2aae2a-a2b3-4927-8865-154dcb88e7f1",
														},
														{
															BaseElement: element.BaseElement{
																ID: "_34c6f8e2-b981-4ff2-ad77-bf89575b604c",
															},
															Name:           "To",
															ItemSubjectRef: "_7d2aae2a-a2b3-4927-8865-154dcb88e7f1",
														},
														{
															BaseElement: element.BaseElement{
																ID: "_d52e7328-818c-4352-ba80-25432f28d380",
															},
															Name: "Current Vacation Status",
														},
														{
															BaseElement: element.BaseElement{
																ID: "_ef216ef0-3b7f-460f-aa59-52f8b6cf578b",
															},
															Name:           "Fixed input",
															ItemSubjectRef: "_7d2aae2a-a2b3-4927-8865-154dcb88e7f1",
														},
													},
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_351b30d2-99d9-43f3-8110-dfa85429a78a",
															},
															Name: "Fixed output",
														},
													},
													InputSets: []element.InputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_9cfa60fd-b2f0-44f5-b97e-5f7b7384851b",
															},
															DataInputRefs: []string{"_871b6349-9bc5-440b-9e01-6e3cf700f9e6", "_34c6f8e2-b981-4ff2-ad77-bf89575b604c", "_d52e7328-818c-4352-ba80-25432f28d380", "_ef216ef0-3b7f-460f-aa59-52f8b6cf578b"},
														},
													},
													OutputSets: []element.OutputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_27eab781-b10a-4a37-9d23-cdd528f75f9e",
															},
															DataOutputRefs: []string{"_351b30d2-99d9-43f3-8110-dfa85429a78a"},
														},
													},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_fd5f72e9-0220-4cec-9346-e5c69ea60171",
														},
														SourceRef: "_103fd816-b422-420d-aec9-c28b54f620db",
														TargetRef: "_871b6349-9bc5-440b-9e01-6e3cf700f9e6",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_7ee2312c-4d81-4d98-a903-8e2dc4398438",
														},
														SourceRef: "_2bed4da9-f468-4fe0-a37c-ae5d4c17e323",
														TargetRef: "_ef216ef0-3b7f-460f-aa59-52f8b6cf578b",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_76336c7a-b159-4425-94c9-57df39f635c2",
														},
														SourceRef: "_775c93ab-82e5-4a62-97f0-d305eda76c92",
														TargetRef: "_ef216ef0-3b7f-460f-aa59-52f8b6cf578b",
													},
												},
											},
										},
									},
									Implementation: "kp:mail",
									OperationRef:   "_83f1c680-7e06-4a9a-9206-396ed8155a71",
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID:                "_02232e32-c3d2-473c-a15d-9c5dca00eadc",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil, nil, nil, nil, nil}},
													},
													Name: "Notify Employee of Refusal",
												},
												Incoming: []string{"_d3a575d8-cd41-4246-ad7a-7a068de6679f"},
												Outgoing: []string{"_93ac178c-3efc-48ae-8a99-8d29a0c4aaf6"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_e5cc49f3-6c54-4673-8bf4-5a2b39bed65c",
															},
															ItemSubjectRef: "_7d2aae2a-a2b3-4927-8865-154dcb88e7f1",
															Name:           "From",
														},
														{
															BaseElement: element.BaseElement{
																ID: "_3a460194-21af-44b5-a9c9-386dd21ce7a6",
															},
															ItemSubjectRef: "_7d2aae2a-a2b3-4927-8865-154dcb88e7f1",
															Name:           "To",
														},
														{
															BaseElement: element.BaseElement{
																ID: "_4e07b75e-586a-4bcb-92b8-086ab645cb89",
															},
															Name: "Current Vacation Status",
														},
														{
															BaseElement: element.BaseElement{
																ID: "_f9102b34-5286-4f47-8342-196ce9540f62",
															},
															ItemSubjectRef: "_7d2aae2a-a2b3-4927-8865-154dcb88e7f1",
															Name:           "Fixed input",
														},
													},
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_f581a76e-a407-4e29-8ac0-781aedc60ecd",
															},
															Name: "Fixed output",
														},
													},
													InputSets: []element.InputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_7fdd5a06-459c-47cc-8e2a-3a509d6ef0d6",
															},
															DataInputRefs: []string{"_e5cc49f3-6c54-4673-8bf4-5a2b39bed65c", "_3a460194-21af-44b5-a9c9-386dd21ce7a6", "_4e07b75e-586a-4bcb-92b8-086ab645cb89", "_f9102b34-5286-4f47-8342-196ce9540f62"},
														},
													},
													OutputSets: []element.OutputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_c69d2a22-01ed-4886-9d43-52dc28c14fc8",
															},
															DataOutputRefs: []string{"_f581a76e-a407-4e29-8ac0-781aedc60ecd"},
														},
													},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_0b462a60-66c6-4b03-8b35-fa9a0d04904e",
														},
														SourceRef: "_103fd816-b422-420d-aec9-c28b54f620db",
														TargetRef: "_e5cc49f3-6c54-4673-8bf4-5a2b39bed65c",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_912f2b3a-c01f-4e9c-ab7a-20e42ccc4f1f",
														},
														SourceRef: "_4234afd4-3ab0-4e09-8a7c-b475458db18a",
														TargetRef: "_e5cc49f3-6c54-4673-8bf4-5a2b39bed65c",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_0e27350a-1980-4a65-b7e0-d8725f63ffa1",
														},
														SourceRef: "_2bed4da9-f468-4fe0-a37c-ae5d4c17e323",
														TargetRef: "_f9102b34-5286-4f47-8342-196ce9540f62",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_d81e3c4f-3ffc-4321-83ed-979433de0163",
														},
														SourceRef: "_775c93ab-82e5-4a62-97f0-d305eda76c92",
														TargetRef: "_f9102b34-5286-4f47-8342-196ce9540f62",
													},
												},
											},
										},
									},
									Implementation: "kp:mail",
									OperationRef:   "_83f1c680-7e06-4a9a-9206-396ed8155a71",
								},
							},
							UserTasks: []element.UserTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID:                "_79523269-7444-4b01-90e9-e23957a9d020",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil}},
													},
													Name: "Manually Approve Vacation",
												},
												Incoming: []string{"_0a1c4f20-509f-4aeb-baf9-acc762f4fdf9"},
												Outgoing: []string{"_3c50d7bb-5ecc-4d34-82a9-a5fa41d9c018"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_5b7fefeb-4c68-4d6a-98ed-530da72165e6",
															},
															Name:           "Employee",
															ItemSubjectRef: "_triso-default-bpmnItemDefinition-string_id",
														},
														{
															BaseElement: element.BaseElement{
																ID: "_40d4e009-34a4-401c-a8f7-acb2c6477e46",
															},
															Name:           "From",
															ItemSubjectRef: "_7d2aae2a-a2b3-4927-8865-154dcb88e7f1",
														},
														{
															BaseElement: element.BaseElement{
																ID: "_4573edc3-bdd5-4971-af84-5b4e1a5f9fd4",
															},
															Name:           "To",
															ItemSubjectRef: "_7d2aae2a-a2b3-4927-8865-154dcb88e7f1",
														},
														{
															BaseElement: element.BaseElement{
																ID: "_f27bddd0-01a1-4ce9-a36c-f25f0d3085be",
															},
															Name:           "Remaining Vacation Days",
															ItemSubjectRef: "_24ce0503-1c7a-41eb-bbfd-28f533e0268f",
														},
													},
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_415c6b4c-7a19-4205-a306-01266576041b",
															},
															Name:           "Vacation Approval",
															ItemSubjectRef: "_35aa1a98-77be-4b78-a590-ba5bafe10a58",
														},
													},
													InputSets: []element.InputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_4fd72db9-08e5-4aaf-8b1a-31b4f7e038ff",
															},
															DataInputRefs: []string{"_5b7fefeb-4c68-4d6a-98ed-530da72165e6", "_40d4e009-34a4-401c-a8f7-acb2c6477e46", "_4573edc3-bdd5-4971-af84-5b4e1a5f9fd4", "_f27bddd0-01a1-4ce9-a36c-f25f0d3085be"},
														},
													},
													OutputSets: []element.OutputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_5537a8b8-9dbe-47e1-8373-1ccff5bd2004",
															},
															DataOutputRefs: []string{"_415c6b4c-7a19-4205-a306-01266576041b"},
														},
													},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_06dd6266-3084-46f1-b5fd-11ae71afdc68",
														},
														SourceRef: "_775c93ab-82e5-4a62-97f0-d305eda76c92",
														TargetRef: "_5b7fefeb-4c68-4d6a-98ed-530da72165e6",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_760f223b-9115-45c9-a8d3-2088688398b0",
														},
														SourceRef: "_2bed4da9-f468-4fe0-a37c-ae5d4c17e323",
														TargetRef: "_40d4e009-34a4-401c-a8f7-acb2c6477e46",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_1e8d5b09-9085-4230-ba1d-06d53b91e0ca",
														},
														SourceRef: "_103fd816-b422-420d-aec9-c28b54f620db",
														TargetRef: "_4573edc3-bdd5-4971-af84-5b4e1a5f9fd4",
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_166012b2-dc8c-4b7c-b383-d8720d9433df",
														},
														SourceRef: "_415c6b4c-7a19-4205-a306-01266576041b",
														TargetRef: "_c2d108b2-e6c0-47e6-a3cb-aa5321bcaeaf",
													},
												},
											},
											PotentialOwners: []element.PotentialOwner{
												{
													HumanPerformer: element.HumanPerformer{
														Performer: element.Performer{
															ResourceRole: element.ResourceRole{
																BaseElement: element.BaseElement{
																	ID: "_b7c60f92-2ce6-4b15-8237-3558a761c2f2",
																},
																ResourceAssignmentExpression: element.ResourceAssignmentExpression{
																	BaseElement: element.BaseElement{
																		ID: "_ac0fe901-da4a-4b3b-93e4-49df342fec22",
																	},
																},
															},
														},
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
														ID: "_1688f604-5edf-4187-ad9e-18f74bfede53",
													},
													Name: "Vacation Refused Automatically",
												},
												Incoming: []string{"_130d903a-f211-4b3e-9711-9061f5978d94"},
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
														ID: "_6677ef80-82df-4951-919d-1f36123b681b",
													},
													Name: "Vacation Approved Automatically",
												},
												Incoming: []string{"_68ef27a7-011e-4ec0-ad45-41f52093de68"},
											},
										},
										DataInputs: []element.DataInput{
											{
												BaseElement: element.BaseElement{
													ID: "_9068cab3-c967-470e-86cd-180f3a207c19",
												},
												Name:           "From",
												ItemSubjectRef: "_7d2aae2a-a2b3-4927-8865-154dcb88e7f1",
											},
										},
										InputSet: element.InputSet{
											BaseElement: element.BaseElement{
												ID: "_d25026f4-35f0-402f-88a5-50c11e8fc6b7",
											},
											DataInputRefs: []string{"_9068cab3-c967-470e-86cd-180f3a207c19"},
										},
									},
								},
								{
									ThrowEvent: element.ThrowEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_1cd5fe29-b3ec-4f21-a1aa-57773f0729ca",
													},
													Name: "Vacation Approved by Manager",
												},
												Incoming: []string{"_81ed1b89-5a4d-4028-a245-4fd38b866b6d"},
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
														ID: "_3ae826ca-5f38-43c4-be3a-35d1157aa27f",
													},
													Name: "Vacation Refused by Manager",
												},
												Incoming: []string{"_93ac178c-3efc-48ae-8a99-8d29a0c4aaf6"},
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
														ID: "_b4d636eb-b501-4462-93c8-04652db10307",
													},
													Name: "Employee not found",
												},
												Incoming: []string{"_f250e125-23eb-45cd-9c61-191d172df093"},
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
