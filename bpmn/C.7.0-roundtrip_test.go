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

func TestC_7_0_roundtrip(t *testing.T) {
	// create test use ./test/C.7.0-roundtrip.bpmn
	path := "./test/C.7.0-roundtrip.bpmn"
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
			ID:                 "_7e80c232-2bf0-46e0-b77e-29ec586a23a6",
			Name:               "Advertise a job vacancy",
			TargetNamespace:    "http://www.trisotech.com/definitions/_7e80c232-2bf0-46e0-b77e-29ec586a23a6",
			ExpressionLanguage: "http://www.omg.org/spec/DMN/20180521/FEEL/",
			TypeLanguage:       "http://www.trisotech.com/2015/triso/modeling/ItemDefinitionType",
			Exporter:           "Workflow Modeler",
			ExporterVersion:    "6.12.3",
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
								ID: "_bd0e1992-4669-4a5b-98ea-641930d87150",
							},
						},
						StructureRef: "collectionOfString",
					},
				},
				Resources: []element.Resource{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_b5b6808a-be81-426c-98ae-f33f44a2f871",
							},
						},
						Name: "Hiring manager",
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_dc7df8e9-cc60-4953-9ae5-a9ea25fd9c5b",
							},
						},
						Name: "Recruiter",
					},
				},
				Collaborations: []element.Collaboration{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_0322c8c5-b921-44cc-9bf7-261dcb16f257",
							},
						},
						Participants: []element.Participant{
							{
								BaseElement: element.BaseElement{
									ID: "_d3aa8a96-e9df-4336-9b0d-01b17e6587ad",
								},
								Name:       "EU Bank",
								ProcessRef: "_4a690dd7-809a-4fa9-ad63-515ac6685375",
							},
						},
					},
				},
				Processes: []element.Process{
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "_4a690dd7-809a-4fa9-ad63-515ac6685375",
								},
							},
							Name: "EU Bank - Process",
							IoSpecification: element.IoSpecification{
								InputOutputSpecification: element.InputOutputSpecification{
									DataInputs: []element.DataInput{
										{
											BaseElement: element.BaseElement{
												ID: "_d08869ef-4951-4592-bb73-363cee03cb90",
											},
											ItemSubjectRef: "_triso-default-bpmnItemDefinition-string_id",
											Name:           "Role \nrequired",
										},
									},
									DataOutputs: []element.DataOutput{
										{
											BaseElement: element.BaseElement{
												ID:  "_b6464e75-dd3d-45d9-84cd-861c42a3bedf",
												Any: []xml.Token{nil},
											},
											ItemSubjectRef: "_triso-default-bpmnItemDefinition-string_id",
											Name:           "Advertisement",
											DataState: element.DataState{
												Name: "",
											},
										},
									},
									InputSets: []element.InputSet{
										{
											BaseElement: element.BaseElement{
												ID: "_e717d617-bac8-4278-ac70-3a9767e3b27f",
											},
											DataInputRefs: []string{"_d08869ef-4951-4592-bb73-363cee03cb90"},
										},
									},
									OutputSets: []element.OutputSet{
										{
											BaseElement: element.BaseElement{
												ID: "_a8015818-8340-4ba0-a3c9-84ee85a353eb",
											},
											DataOutputRefs: []string{"_b6464e75-dd3d-45d9-84cd-861c42a3bedf"},
										},
									},
								},
							},
						},
						IsExecutable: false,
						LaneSets: []element.LaneSet{
							{
								Lanes: []element.Lane{
									{
										BaseElement: element.BaseElement{
											ID: "_b836aa5e-fb94-4479-af77-64a3a5202451",
										},
										Name: "Hiring manager",
										FlowNodeRefs: []string{
											"_5ba97787-8a90-4002-8277-b0895e45cf1f",
											"_392c86ba-38b5-4dc9-b98d-f97ad4c2add5",
											"_15b00027-5049-4081-8952-fd398e8b722a",
											"_26c40c03-5d1f-46c5-81f1-ddd485868125",
										},
									},
									{
										BaseElement: element.BaseElement{
											ID: "_dd32321b-8e95-4801-8eed-5451399b4378",
										},
										Name: "Recruitment",
										FlowNodeRefs: []string{
											"_d3435084-f2c7-43cc-abcc-c679bc4232ac",
											"_b13d6fa3-fc78-40c7-ae77-609be07493e9",
											"_64eabfe9-6947-43eb-ac45-8d331745f86c",
											"_eae674ce-4d6e-48ac-819c-c79e0868e40d",
											"_0783f019-f40c-43d6-ab40-0f1c81f8d9e7",
											"_c456dbcc-bbe3-4c75-b57d-9427525c0a94",
											"_a36ddf2f-23c1-46c5-86d4-bd2a0eb42535",
										},
									},
								},
							},
						},
						FlowElements: element.FlowElements{
							DataObjectReferenes: []element.DataObjectReference{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_bd7b6a15-4ef8-46a9-8be9-20a5abb32abd",
										},
										Name: "Description",
									},
									DataObjectRef: "_8f2796af-2fbe-4f72-80c1-96933c38990f",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_bac0224d-62f1-46aa-bb4c-5371c3983ffb",
										},
										Name: "Advertisement",
									},
									DataObjectRef: "_f60fe1d9-58bd-462c-9d62-153e530dc79d",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_c68abea8-c5b4-4aef-b1a5-1e81caec0cba",
										},
										Name: "Selected\n platforms",
									},
									DataObjectRef: "_ef29e636-bdfe-4eb0-9633-7d0195a8ae3a",
								},
							},
							DataObjects: []element.DataObject{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_ef29e636-bdfe-4eb0-9633-7d0195a8ae3a",
										},
										Name: "Selected\n platforms",
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_8f2796af-2fbe-4f72-80c1-96933c38990f",
										},
										Name: "Description",
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_f60fe1d9-58bd-462c-9d62-153e530dc79d",
										},
										Name: "Advertisement",
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
														ID: "_5ba97787-8a90-4002-8277-b0895e45cf1f",
													},
													Name: "Job \nvacancy",
												},
												Outgoing: []string{"_a4c93e8a-2b52-4367-8381-a3f78450a075"},
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
														ID: "_392c86ba-38b5-4dc9-b98d-f97ad4c2add5",
														Documentation: []element.Documentation{
															{
																TextFormat: "text/html",
																Value:      "<p>A very specific requirement to the job vacancy has to be written and forwarded to the recruitment department.</p>",
															},
														},
													},
													Name: "Write \ndescription",
												},
												Incoming: []string{"_a4c93e8a-2b52-4367-8381-a3f78450a075"},
												Outgoing: []string{"_8a27a9ee-b8e5-49d8-8b7d-41a59d74b3f3"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_31ca347a-7ad3-4416-9d4a-c080479329c7",
															},
															ItemSubjectRef: "_triso-default-bpmnItemDefinition-string_id",
															Name:           "Description",
														},
													},
													InputSets: []element.InputSet{{}},
													OutputSets: []element.OutputSet{{
														BaseElement: element.BaseElement{
															ID: "_0d01e81d-03c2-4002-ae80-53df972febb6",
														},
														DataOutputRefs: []string{"_31ca347a-7ad3-4416-9d4a-c080479329c7"},
													}},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_e2734375-2aa0-418c-9f0b-8c2ca1022285",
														},
														SourceRef: "_31ca347a-7ad3-4416-9d4a-c080479329c7",
														TargetRef: "_bd7b6a15-4ef8-46a9-8be9-20a5abb32abd",
													},
												},
											},
											Performers: []element.Performer{
												{
													ResourceRole: element.ResourceRole{
														BaseElement: element.BaseElement{
															ID: "_c663bd75-4868-4eb4-8829-0cbbd5ae4580",
														},
														ResourceRef: "_b5b6808a-be81-426c-98ae-f33f44a2f871",
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
														ID: "_15b00027-5049-4081-8952-fd398e8b722a",
														Documentation: []element.Documentation{
															{
																TextFormat: "text/html",
																Value:      "<p>The job description edited to a job advertisement has to be checked and approved.</p>",
															},
														},
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
													},
													Name: "Approve advertisement",
												},
												Incoming: []string{"_842553bd-02e1-45f7-aed0-c78e9dba9e2b"},
												Outgoing: []string{"_7b66938e-e5cb-4026-a863-1e6de6d11a79"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_b02b2e9f-17f7-4663-b12c-255ee4351d42",
															},
															ItemSubjectRef: "_triso-default-bpmnItemDefinition-string_id",
															Name:           "Advertisement",
														},
													},
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_e115b5bb-0dd7-444f-b5ec-20cffffef082",
															},
															ItemSubjectRef: "_triso-default-bpmnItemDefinition-string_id",
															Name:           "Advertisement",
														},
													},
													InputSets: []element.InputSet{{
														BaseElement: element.BaseElement{
															ID: "_0cddfa48-58e5-47f8-92f0-a5153ae17d78",
														},
														DataInputRefs: []string{"_b02b2e9f-17f7-4663-b12c-255ee4351d42"},
													}},
													OutputSets: []element.OutputSet{{
														BaseElement: element.BaseElement{
															ID: "_9492b1d3-f6ff-4c9c-8707-07e2d7b07c89",
														},
														DataOutputRefs: []string{"_e115b5bb-0dd7-444f-b5ec-20cffffef082"},
													}},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_9ce7bb4a-930c-4132-a09d-61b345cce434",
														},
														SourceRef: "_bac0224d-62f1-46aa-bb4c-5371c3983ffb",
														TargetRef: "_b02b2e9f-17f7-4663-b12c-255ee4351d42",
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_adddc8ea-507e-4894-ac4b-92bca7e0a89f",
														},
														SourceRef: "_e115b5bb-0dd7-444f-b5ec-20cffffef082",
														TargetRef: "_b6464e75-dd3d-45d9-84cd-861c42a3bedf",
													},
												},
											},
											Performers: []element.Performer{
												{
													ResourceRole: element.ResourceRole{
														BaseElement: element.BaseElement{
															ID: "_ceb60e9f-c1f4-43bc-aefd-6f346c176a6c",
														},
														ResourceRef: "_b5b6808a-be81-426c-98ae-f33f44a2f871",
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
														ID: "_d3435084-f2c7-43cc-abcc-c679bc4232ac",
														Documentation: []element.Documentation{
															{
																TextFormat: "text/html",
																Value:      "<p>The job description received from the specialist department has to be completed (layout, additional information etc.) to a job advertisement in accordance to the guidelines and approved by the Hiring Manager.</p>",
															},
														},
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
													},
													Name: "Complete advertisement",
												},
												Incoming: []string{"_8a27a9ee-b8e5-49d8-8b7d-41a59d74b3f3", "_d74707c7-6af3-4db7-9403-924bfdf6a7d8"},
												Outgoing: []string{"_842553bd-02e1-45f7-aed0-c78e9dba9e2b"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_c083f111-3c38-4250-9594-3f25b4620db3",
															},
															ItemSubjectRef: "_triso-default-bpmnItemDefinition-string_id",
															Name:           "Description",
														},
													},
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_50a31269-be15-435b-81e0-0a91b57ba10d",
															},
															ItemSubjectRef: "_triso-default-bpmnItemDefinition-string_id",
															Name:           "Advertisement",
														},
													},
													InputSets: []element.InputSet{{
														BaseElement: element.BaseElement{
															ID: "_9c3de1df-92ef-4b69-9585-311d8d45952f",
														},
														DataInputRefs: []string{"_c083f111-3c38-4250-9594-3f25b4620db3"},
													}},
													OutputSets: []element.OutputSet{{
														BaseElement: element.BaseElement{
															ID: "_a4d61918-d67b-4038-b488-95c244824918",
														},
														DataOutputRefs: []string{"_50a31269-be15-435b-81e0-0a91b57ba10d"},
													}},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_5c3fc96e-20d0-4879-8471-d41224632e24",
														},
														SourceRef: "_bd7b6a15-4ef8-46a9-8be9-20a5abb32abd",
														TargetRef: "_c083f111-3c38-4250-9594-3f25b4620db3",
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_4164c380-3ee3-4a6a-8e66-5bc201416108",
														},
														SourceRef: "_50a31269-be15-435b-81e0-0a91b57ba10d",
														TargetRef: "_bac0224d-62f1-46aa-bb4c-5371c3983ffb",
													},
												},
											},
											Performers: []element.Performer{
												{
													ResourceRole: element.ResourceRole{
														BaseElement: element.BaseElement{
															ID: "_41c0e9b3-8ef9-4657-baab-c59b6ba1b1c4",
														},
														ResourceRef: "_dc7df8e9-cc60-4953-9ae5-a9ea25fd9c5b",
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
													ID: "_26c40c03-5d1f-46c5-81f1-ddd485868125",
												},
												Name: "Advertisement approved?",
											},
											Incoming: []string{"_7b66938e-e5cb-4026-a863-1e6de6d11a79"},
											Outgoing: []string{"_d74707c7-6af3-4db7-9403-924bfdf6a7d8", "_1d201a22-d500-4412-a32a-2c7e24ad4d6b"},
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
													ID: "_b13d6fa3-fc78-40c7-ae77-609be07493e9",
												},
											},
											Incoming: []string{"_1d201a22-d500-4412-a32a-2c7e24ad4d6b"},
											Outgoing: []string{"_f476667f-44f5-4fec-9414-bf70d987853f", "_f3187dce-c37e-4d5b-b49c-ed28965abb73"},
										},
										GatewayDirection: "Diverging",
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_0783f019-f40c-43d6-ab40-0f1c81f8d9e7",
												},
											},
											Incoming: []string{"_720cb9a3-20df-4da1-a923-5336b269c104", "_847352f2-ac0c-44be-9e24-f4c7f76bfe7e"},
											Outgoing: []string{"_c43defc5-4470-4bfe-8a8f-4d59ca6abeeb"},
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
														ID: "_64eabfe9-6947-43eb-ac45-8d331745f86c",
														Documentation: []element.Documentation{
															{
																TextFormat: "text/html",
																Value:      "<p>The approved job advertisement has to be published on the homepage.</p>",
															},
														},
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
													},
													Name: "Publish on \nhomepage",
												},
												Incoming: []string{"_f476667f-44f5-4fec-9414-bf70d987853f"},
												Outgoing: []string{"_720cb9a3-20df-4da1-a923-5336b269c104"},
											},
											Performers: []element.Performer{
												{
													ResourceRole: element.ResourceRole{
														BaseElement: element.BaseElement{
															ID: "_2d69e812-dd75-4405-ae41-5f64be302e5f",
														},
														ResourceRef: "_dc7df8e9-cc60-4953-9ae5-a9ea25fd9c5b",
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
														ID: "_a36ddf2f-23c1-46c5-86d4-bd2a0eb42535",
													},
													Name: "Publish on \nother platforms",
												},
												Incoming: []string{"_64ac4473-a38f-4043-973a-69c497adc88a"},
												Outgoing: []string{"_847352f2-ac0c-44be-9e24-f4c7f76bfe7e"},
											},
											LoopCharacteristicsElements: element.LoopCharacteristicsElements{
												MultielementLoopCharacteristics: []element.MultiInstanceLoopCharacteristics{
													{
														LoopCharacteristics: element.LoopCharacteristics{
															BaseElement: element.BaseElement{
																ID: "_970a7f54-893a-495e-8cbc-545fd631f626",
															},
														},
													},
												},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_d130e681-7e04-48d6-be39-2b7cfe24b9d4",
															},
															ItemSubjectRef: "_triso-default-bpmnItemDefinition-string_id",
															Name:           "Selected\n platforms",
															IsCollection:   true,
														},
													},
													InputSets: []element.InputSet{{
														BaseElement: element.BaseElement{
															ID: "_f49af959-916f-43bc-9f37-00ef9c44e967",
														},
														DataInputRefs: []string{"_d130e681-7e04-48d6-be39-2b7cfe24b9d4"},
													}},
													OutputSets: []element.OutputSet{{}},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_505fa8cf-5697-4e35-bc4e-63df8e99093e",
														},
														SourceRef: "_c68abea8-c5b4-4aef-b1a5-1e81caec0cba",
														TargetRef: "_d130e681-7e04-48d6-be39-2b7cfe24b9d4",
													},
												},
											},
										},
									},
									Implementation: "##WebService",
								},
							},
							BusinessRuleTasks: []element.BusinessRuleTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_eae674ce-4d6e-48ac-819c-c79e0868e40d",
														Documentation: []element.Documentation{
															{
																TextFormat: "text/html",
																Value:      "<p>3rd party career platforms have to be selected on which the job advertisement will be published.</p>",
															},
														},
													},
													Name: "Select \nother platforms",
												},
												Incoming: []string{"_f3187dce-c37e-4d5b-b49c-ed28965abb73"},
												Outgoing: []string{"_64ac4473-a38f-4043-973a-69c497adc88a"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_3daad61f-55b2-4fbd-bc79-dff572a23d69",
															},
															ItemSubjectRef: "_triso-default-bpmnItemDefinition-string_id",
															Name:           "Selected\n platforms",
															IsCollection:   true,
														},
													},
													InputSets: []element.InputSet{{}},
													OutputSets: []element.OutputSet{{
														BaseElement: element.BaseElement{
															ID: "_a389083a-36ae-4975-b210-9464fe82c22e",
														},
														DataOutputRefs: []string{"_3daad61f-55b2-4fbd-bc79-dff572a23d69"},
													}},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_d65b0a3d-8c42-4be7-9ac1-7324aa3d31f8",
														},
														SourceRef: "_3daad61f-55b2-4fbd-bc79-dff572a23d69",
														TargetRef: "_c68abea8-c5b4-4aef-b1a5-1e81caec0cba",
													},
												},
											},
											Performers: []element.Performer{
												{
													ResourceRole: element.ResourceRole{
														BaseElement: element.BaseElement{
															ID: "_867590ee-ec0c-4f56-b0f0-32e444b73079",
														},
														ResourceRef: "_dc7df8e9-cc60-4953-9ae5-a9ea25fd9c5b",
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
														ID: "_c456dbcc-bbe3-4c75-b57d-9427525c0a94",
													},
													Name: "Vacancy \nadvertised",
												},
												Incoming: []string{"_c43defc5-4470-4bfe-8a8f-4d59ca6abeeb"},
											},
										},
									},
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_a4c93e8a-2b52-4367-8381-a3f78450a075",
										},
									},
									SourceRef: "_5ba97787-8a90-4002-8277-b0895e45cf1f",
									TargetRef: "_392c86ba-38b5-4dc9-b98d-f97ad4c2add5",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_8a27a9ee-b8e5-49d8-8b7d-41a59d74b3f3",
										},
									},
									SourceRef: "_392c86ba-38b5-4dc9-b98d-f97ad4c2add5",
									TargetRef: "_d3435084-f2c7-43cc-abcc-c679bc4232ac",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_842553bd-02e1-45f7-aed0-c78e9dba9e2b",
										},
									},
									SourceRef: "_d3435084-f2c7-43cc-abcc-c679bc4232ac",
									TargetRef: "_15b00027-5049-4081-8952-fd398e8b722a",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_7b66938e-e5cb-4026-a863-1e6de6d11a79",
										},
									},
									SourceRef: "_15b00027-5049-4081-8952-fd398e8b722a",
									TargetRef: "_26c40c03-5d1f-46c5-81f1-ddd485868125",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_d74707c7-6af3-4db7-9403-924bfdf6a7d8",
										},
										Name: "No",
									},
									SourceRef: "_26c40c03-5d1f-46c5-81f1-ddd485868125",
									TargetRef: "_d3435084-f2c7-43cc-abcc-c679bc4232ac",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_1d201a22-d500-4412-a32a-2c7e24ad4d6b",
										},
										Name: "Yes",
									},
									SourceRef: "_26c40c03-5d1f-46c5-81f1-ddd485868125",
									TargetRef: "_b13d6fa3-fc78-40c7-ae77-609be07493e9",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_f476667f-44f5-4fec-9414-bf70d987853f",
										},
									},
									SourceRef: "_b13d6fa3-fc78-40c7-ae77-609be07493e9",
									TargetRef: "_64eabfe9-6947-43eb-ac45-8d331745f86c",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_f3187dce-c37e-4d5b-b49c-ed28965abb73",
										},
									},
									SourceRef: "_b13d6fa3-fc78-40c7-ae77-609be07493e9",
									TargetRef: "_eae674ce-4d6e-48ac-819c-c79e0868e40d",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_720cb9a3-20df-4da1-a923-5336b269c104",
										},
									},
									SourceRef: "_64eabfe9-6947-43eb-ac45-8d331745f86c",
									TargetRef: "_0783f019-f40c-43d6-ab40-0f1c81f8d9e7",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_64ac4473-a38f-4043-973a-69c497adc88a",
										},
									},
									SourceRef: "_eae674ce-4d6e-48ac-819c-c79e0868e40d",
									TargetRef: "_a36ddf2f-23c1-46c5-86d4-bd2a0eb42535",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_c43defc5-4470-4bfe-8a8f-4d59ca6abeeb",
										},
									},
									SourceRef: "_0783f019-f40c-43d6-ab40-0f1c81f8d9e7",
									TargetRef: "_c456dbcc-bbe3-4c75-b57d-9427525c0a94",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_847352f2-ac0c-44be-9e24-f4c7f76bfe7e",
										},
									},
									SourceRef: "_a36ddf2f-23c1-46c5-86d4-bd2a0eb42535",
									TargetRef: "_0783f019-f40c-43d6-ab40-0f1c81f8d9e7",
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
