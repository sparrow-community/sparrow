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

func TestC_8_0_roundtrip(t *testing.T) {
	// create test use ./test/C.8.0-roundtrip.bpmn
	path := "./test/C.8.0-roundtrip.bpmn"
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
			ID:              "Definitions__a9e3a63e-db87-4fa1-9091-328acae4f490",
			Name:            "Vacation Request - (colors)",
			TargetNamespace: "http://www.softwareag.com/aris/bpmn2",
			Exporter:        "Workflow Modeler",
			ExporterVersion: "7.14.1.202203211419",
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
				},
				Processes: []element.Process{
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID:                "VacationRequestProcess",
									Documentation:     []element.Documentation{{Value: "Vacation Request - BPMN MIWG demo for 2022"}},
									ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil, nil}},
								},
							},
							Name: "Vacation Request - (i18n)",
							IoSpecification: element.IoSpecification{
								InputOutputSpecification: element.InputOutputSpecification{
									DataInputs: []element.DataInput{
										{
											BaseElement: element.BaseElement{
												ID:                "dataInputdataInput_8b9aa28f-5974-4087-9895-0467c25635dc",
												ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil}},
											},
											Name: "Employee Badge Number",
										},
										{
											BaseElement: element.BaseElement{
												ID:                "dataInputdataInput_2bed4da9-f468-4fe0-a37c-ae5d4c17e323",
												ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil}},
											},
											Name: "From",
										},
										{
											BaseElement: element.BaseElement{
												ID:                "dataInputdataInput_103fd816-b422-420d-aec9-c28b54f620db",
												ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil}},
											},
											Name: "To",
										},
									},
									DataOutputs: []element.DataOutput{
										{
											BaseElement: element.BaseElement{
												ID:                "dataOutputdataOutput_dac8ee76-f637-4fd9-8357-6a87fd11ef41",
												ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil}},
											},
											Name: "Reason",
										},
										{
											BaseElement: element.BaseElement{
												ID:                "dataOutputdataOutput_c2d108b2-e6c0-47e6-a3cb-aa5321bcaeaf",
												ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil}},
											},
											Name: "Vacation Approval",
										},
									},
									InputSets: []element.InputSet{
										{
											BaseElement: element.BaseElement{
												ID: "_4d1f4479-b14c-4dd0-9675-607384e9dd81",
											},
											DataInputRefs: []string{"dataInputdataInput_8b9aa28f-5974-4087-9895-0467c25635dc", "dataInputdataInput_2bed4da9-f468-4fe0-a37c-ae5d4c17e323", "dataInputdataInput_103fd816-b422-420d-aec9-c28b54f620db"},
										},
									},
									OutputSets: []element.OutputSet{
										{
											BaseElement: element.BaseElement{
												ID: "_10b372b4-3cdd-4ddd-aa58-129b5df7b2f8",
											},
											DataOutputRefs: []string{"dataOutputdataOutput_dac8ee76-f637-4fd9-8357-6a87fd11ef41", "dataOutputdataOutput_c2d108b2-e6c0-47e6-a3cb-aa5321bcaeaf"},
										},
									},
								},
							},
						},
						ProcessType:  "None",
						IsExecutable: false,
						FlowElements: element.FlowElements{
							StartEvents: []element.StartEvent{
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID:                "_b1625a52-aaf0-4694-86cb-7af891212ac6",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil, nil, nil, nil, nil, nil}},
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
											ID:                "_b02b5bfa-5629-4af5-b1bc-ac2a86d88adb",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil}},
										},
									},
									SourceRef: "_b1625a52-aaf0-4694-86cb-7af891212ac6",
									TargetRef: "_2b960d84-feb1-46a9-a1a1-c300dd996b99",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_81ed1b89-5a4d-4028-a245-4fd38b866b6d",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil}},
										},
									},
									SourceRef: "_5e16a4e0-0f23-482a-be47-d3edbc5741ba",
									TargetRef: "_1cd5fe29-b3ec-4f21-a1aa-57773f0729ca",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_2a5ee73c-c495-48a0-945b-3dabab28ca91",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil, nil, nil}},
										},
										Name: "Refused",
									},
									SourceRef: "_42367c5f-d084-44ee-90c7-960d1ab02a3b",
									TargetRef: "_9ed61a6a-7cc1-4ed1-86d8-03482b0983c9",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_d79c991e-446c-47d1-ac9d-9d0113e35b93",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil}},
										},
										Name: "",
									},
									SourceRef: "_2b960d84-feb1-46a9-a1a1-c300dd996b99",
									TargetRef: "_1a818a94-ba6f-413b-a7e8-6f8fd2a11e32",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_93ac178c-3efc-48ae-8a99-8d29a0c4aaf6",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil}},
										},
									},
									SourceRef: "_02232e32-c3d2-473c-a15d-9c5dca00eadc",
									TargetRef: "_3ae826ca-5f38-43c4-be3a-35d1157aa27f",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_130d903a-f211-4b3e-9711-9061f5978d94",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil}},
										},
									},
									SourceRef: "_9ed61a6a-7cc1-4ed1-86d8-03482b0983c9",
									TargetRef: "_1688f604-5edf-4187-ad9e-18f74bfede53",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_d3a575d8-cd41-4246-ad7a-7a068de6679f",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil, nil, nil}},
										},
										Name: "Refused",
									},
									SourceRef: "_64bb8b55-d348-41d6-9ea8-f333b1d8cb69",
									TargetRef: "_02232e32-c3d2-473c-a15d-9c5dca00eadc",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_f2b0da63-d841-4457-ad85-7d86c8b5c1d2",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil, nil, nil}},
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
													ID: "_f2b0da63-d841-4457-ad85-7d86c8b5c1d2condExpr",
												},
											},
											Value: `Vacation Approval = "Approved"`,
										},
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_a5441c65-06f8-4972-8158-e6f61f904841",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil}},
										},
									},
									SourceRef: "_a97c1a48-faba-447b-bfa6-7aa81a6fe0a0",
									TargetRef: "_5e16a4e0-0f23-482a-be47-d3edbc5741ba",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_188d38f6-b73d-4c62-964b-affc6fe040c1",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil}},
										},
									},
									SourceRef: "_93ec9873-edf1-4549-b052-961994ec8234",
									TargetRef: "_4b72053b-8ebb-4ae6-99c6-7c93cf1c1d1b",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_325973e7-0bc8-4136-b6df-be1e681d8608",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil, nil, nil}},
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
													ID: "_325973e7-0bc8-4136-b6df-be1e681d8608condExpr",
												},
											},
											Value: `Vacation Approval = "Approved"`,
										},
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_3c50d7bb-5ecc-4d34-82a9-a5fa41d9c018",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil}},
										},
										Name: "",
									},
									SourceRef: "_79523269-7444-4b01-90e9-e23957a9d020",
									TargetRef: "_64bb8b55-d348-41d6-9ea8-f333b1d8cb69",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_f250e125-23eb-45cd-9c61-191d172df093",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil}},
										},
									},
									SourceRef: "_f8fcb377-3d7d-4138-9a7e-6ab58b97e29d",
									TargetRef: "_b4d636eb-b501-4462-93c8-04652db10307",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_68ef27a7-011e-4ec0-ad45-41f52093de68",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil}},
										},
										Name: "",
									},
									SourceRef: "_4b72053b-8ebb-4ae6-99c6-7c93cf1c1d1b",
									TargetRef: "_6677ef80-82df-4951-919d-1f36123b681b",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_0a1c4f20-509f-4aeb-baf9-acc762f4fdf9",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil, nil, nil}},
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
													ID: "_0a1c4f20-509f-4aeb-baf9-acc762f4fdf9condExpr",
												},
											},
											Value: `Vacation Approval = "Manual Validation Required"`,
										},
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_ac1d82df-9bf6-404b-be70-b052e311f8b7",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil}},
										},
									},
									SourceRef: "_1a818a94-ba6f-413b-a7e8-6f8fd2a11e32",
									TargetRef: "_42367c5f-d084-44ee-90c7-960d1ab02a3b",
								},
							},
							ServiceTasks: []element.ServiceTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_5e16a4e0-0f23-482a-be47-d3edbc5741ba",
														ExtensionElements: element.ExtensionElements{
															Any: []xml.Token{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
														},
														Any: nil,
													},
													Name: "Update Remaining Vacation",
												},
												Incoming: []string{"_a5441c65-06f8-4972-8158-e6f61f904841"},
												Outgoing: []string{"_81ed1b89-5a4d-4028-a245-4fd38b866b6d"},
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
														ID: "_4b72053b-8ebb-4ae6-99c6-7c93cf1c1d1b",
														ExtensionElements: element.ExtensionElements{
															Any: []xml.Token{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
														},
													},
													Name: "Update Remaining Vacation",
												},
												Incoming: []string{"_188d38f6-b73d-4c62-964b-affc6fe040c1"},
												Outgoing: []string{"_68ef27a7-011e-4ec0-ad45-41f52093de68"},
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
														ID: "_2b960d84-feb1-46a9-a1a1-c300dd996b99",
														ExtensionElements: element.ExtensionElements{
															Any: []xml.Token{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
														},
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
																ID: "DataInput__2b960d84-feb1-46a9-a1a1-c300dd996b99",
															},
														},
													},
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "DataOutput__2b960d84-feb1-46a9-a1a1-c300dd996b99",
															},
														},
													},
													InputSets: []element.InputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_74206bda-6edc-4a98-a433-af5446be8809",
															},
															DataInputRefs: []string{"DataInput__2b960d84-feb1-46a9-a1a1-c300dd996b99"},
														},
													},
													OutputSets: []element.OutputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_cd959c80-ff41-401a-9282-b76dcf95aaf2",
															},
															DataOutputRefs: []string{"DataOutput__2b960d84-feb1-46a9-a1a1-c300dd996b99"},
														},
													},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID:                "_f4f41761-cdd8-4b90-a956-9182345df488",
															ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
														},
														SourceRef: "dataInputdataInput_8b9aa28f-5974-4087-9895-0467c25635dc",
														TargetRef: "DataInput__2b960d84-feb1-46a9-a1a1-c300dd996b99",
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID:                "_40d3cb58-31bb-47a4-9591-032a38011de3",
															ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
														},
														SourceRef: "DataOutput__2b960d84-feb1-46a9-a1a1-c300dd996b99",
														TargetRef: "Reference__775c93ab-82e5-4a62-97f0-d305eda76c92",
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
														ID:                "_f8fcb377-3d7d-4138-9a7e-6ab58b97e29d",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
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
																ID: "Definition__f8fcb377-3d7d-4138-9a7e-6ab58b97e29d",
															},
														},
													},
												},
											},
										},
									},
									AttachedToRef: "_2b960d84-feb1-46a9-a1a1-c300dd996b99",
								},
							},
							DataObjectReferenes: []element.DataObjectReference{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "Reference__775c93ab-82e5-4a62-97f0-d305eda76c92",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
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
											ID:                "_775c93ab-82e5-4a62-97f0-d305eda76c92",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil}},
										},
										Name: "Current Vacation Status",
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
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil}},
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
																ID:                "DataInput__1a818a94-ba6f-413b-a7e8-6f8fd2a11e32_1",
																ExtensionElements: element.ExtensionElements{},
															},
															Name: "",
														},
														{
															BaseElement: element.BaseElement{
																ID:                "DataInput__1a818a94-ba6f-413b-a7e8-6f8fd2a11e32",
																ExtensionElements: element.ExtensionElements{},
															},
															Name: "",
														},
													},
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "DataOutput__1a818a94-ba6f-413b-a7e8-6f8fd2a11e32_1",
															},
															Name: "",
														},
														{
															BaseElement: element.BaseElement{
																ID: "DataOutput__1a818a94-ba6f-413b-a7e8-6f8fd2a11e32",
															},
															Name: "",
														},
													},
													InputSets: []element.InputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_7522b1b6-fa6b-4e86-8bd7-d77f8d0db55a",
															},
															DataInputRefs: []string{"DataInput__1a818a94-ba6f-413b-a7e8-6f8fd2a11e32_1", "DataInput__1a818a94-ba6f-413b-a7e8-6f8fd2a11e32"},
														},
													},
													OutputSets: []element.OutputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_7821d757-2e54-4340-8286-6b138f27ca51",
															},
															DataOutputRefs: []string{"DataOutput__1a818a94-ba6f-413b-a7e8-6f8fd2a11e32_1", "DataOutput__1a818a94-ba6f-413b-a7e8-6f8fd2a11e32"},
														},
													},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID:                "_d7aa7a34-bb93-49ae-87d7-f3ccd25add3f",
															ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
														},
														SourceRef: "dataInputdataInput_2bed4da9-f468-4fe0-a37c-ae5d4c17e323",
														TargetRef: "DataInput__1a818a94-ba6f-413b-a7e8-6f8fd2a11e32_1",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID:                "_bfe6b00e-a1ce-4b12-b868-f445ebc77d1d",
															ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
														},
														SourceRef: "dataInputdataInput_103fd816-b422-420d-aec9-c28b54f620db",
														TargetRef: "DataInput__1a818a94-ba6f-413b-a7e8-6f8fd2a11e32",
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID:                "_3dcf2a50-5a30-4a47-956f-7962eb907737",
															ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
														},
														SourceRef: "DataOutput__1a818a94-ba6f-413b-a7e8-6f8fd2a11e32_1",
														TargetRef: "dataOutputdataOutput_c2d108b2-e6c0-47e6-a3cb-aa5321bcaeaf",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID:                "_66f6f24d-1cad-4c93-b9b3-b93193015a18",
															ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
														},
														SourceRef: "DataOutput__1a818a94-ba6f-413b-a7e8-6f8fd2a11e32",
														TargetRef: "dataOutputdataOutput_dac8ee76-f637-4fd9-8357-6a87fd11ef41",
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
													ID:                "_42367c5f-d084-44ee-90c7-960d1ab02a3b",
													ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil}},
												},
												Name: "",
											},
											Incoming: []string{"_ac1d82df-9bf6-404b-be70-b052e311f8b7"},
											Outgoing: []string{"_2a5ee73c-c495-48a0-945b-3dabab28ca91", "_325973e7-0bc8-4136-b6df-be1e681d8608", "_0a1c4f20-509f-4aeb-baf9-acc762f4fdf9"},
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
													ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil}},
												},
												Name: "",
											},
											Incoming: []string{"_3c50d7bb-5ecc-4d34-82a9-a5fa41d9c018"},
											Outgoing: []string{"_d3a575d8-cd41-4246-ad7a-7a068de6679f", "_f2b0da63-d841-4457-ad85-7d86c8b5c1d2"},
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
														ID:                "_a97c1a48-faba-447b-bfa6-7aa81a6fe0a0",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil}},
													},
													Name: "Notify Employee of Approval",
												},
												Incoming: []string{"_f2b0da63-d841-4457-ad85-7d86c8b5c1d2"},
												Outgoing: []string{"_a5441c65-06f8-4972-8158-e6f61f904841"},
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
														ID:                "_02232e32-c3d2-473c-a15d-9c5dca00eadc",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil}},
													},
													Name: "Notify Employee of Refusal",
												},
												Incoming: []string{"_d3a575d8-cd41-4246-ad7a-7a068de6679f"},
												Outgoing: []string{"_93ac178c-3efc-48ae-8a99-8d29a0c4aaf6"},
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
														ID:                "_9ed61a6a-7cc1-4ed1-86d8-03482b0983c9",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil}},
													},
													Name: "Notify Employee of Refusal",
												},
												Incoming: []string{"_2a5ee73c-c495-48a0-945b-3dabab28ca91"},
												Outgoing: []string{"_130d903a-f211-4b3e-9711-9061f5978d94"},
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
														ID:                "_93ec9873-edf1-4549-b052-961994ec8234",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil}},
													},
													Name: "Notify Employee of Approval",
												},
												Incoming: []string{"_325973e7-0bc8-4136-b6df-be1e681d8608"},
												Outgoing: []string{"_188d38f6-b73d-4c62-964b-affc6fe040c1"},
											},
										},
									},
									Implementation: "##WebService",
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
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil, nil, nil, nil}},
													},
													Name: "Manually Approve Vacation",
												},
												Incoming: []string{"_0a1c4f20-509f-4aeb-baf9-acc762f4fdf9"},
												Outgoing: []string{"_3c50d7bb-5ecc-4d34-82a9-a5fa41d9c018"},
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
														ID:                "_1cd5fe29-b3ec-4f21-a1aa-57773f0729ca",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil}},
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
														ID:                "_6677ef80-82df-4951-919d-1f36123b681b",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil}},
													},
													Name: "Vacation Approved Automatically",
												},
												Incoming: []string{"_68ef27a7-011e-4ec0-ad45-41f52093de68"},
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
														ID:                "_3ae826ca-5f38-43c4-be3a-35d1157aa27f",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil}},
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
														ID:                "_b4d636eb-b501-4462-93c8-04652db10307",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil}},
													},
													Name: "Employee not found",
												},
												Incoming: []string{"_f250e125-23eb-45cd-9c61-191d172df093"},
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
														ID:                "_1688f604-5edf-4187-ad9e-18f74bfede53",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil, nil}},
													},
													Name: "Vacation Refused Automatically",
												},
												Incoming: []string{"_130d903a-f211-4b3e-9711-9061f5978d94"},
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
