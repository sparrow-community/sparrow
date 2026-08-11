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

func TestC_5_0_roundtrip(t *testing.T) {
	// create test use ./test/C.5.0-roundtrip.bpmn
	path := "./test/C.5.0-roundtrip.bpmn"
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
			ID:              "_5f5fbc02-46c1-4615-b2bd-89b2d938ac0b",
			Name:            "C.5.0",
			TargetNamespace: "http://www.trisotech.com/definitions/_5f5fbc02-46c1-4615-b2bd-89b2d938ac0b",
			Exporter:        "BPMN Modeler",
			ExporterVersion: "6.1.19.201902051502",
			RootElemnts: element.RootElemnts{
				Collaborations: []element.Collaboration{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_906eeac9-47e3-41c3-a8db-b8abb8fd95e6",
							},
						},
						Participants: []element.Participant{
							{
								BaseElement: element.BaseElement{
									ID: "_f200a9f1-5554-4868-b895-707186915d62",
								},
								Name:       "Bank",
								ProcessRef: "_3d1ef204-2d4c-4643-8fc5-c319cc032ec0",
							},
						},
					},
				},
				DataStores: []element.DataStore{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_776bd6ca-5f18-432f-a748-f1930bfdf16e",
							},
						},
						Name: "Bank System",
					},
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_46aee7ee-fab2-4735-80cc-fa78092ef92e",
							},
						},
						Name: "Customer Data (temporary storage)",
					},
				},
				Processes: []element.Process{
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "_3d1ef204-2d4c-4643-8fc5-c319cc032ec0",
								},
							},
							Name: "Bank - Process",
						},
						IsExecutable: false,
						LaneSets: []element.LaneSet{
							{
								Lanes: []element.Lane{
									{
										BaseElement: element.BaseElement{
											ID: "_2935f981-e194-4a1b-bb22-846ad3c0f72c",
										},
										Name: "Private Customer Account Manager",
										FlowNodeRefs: []string{
											"_0254d83d-d943-466f-8b62-20e87cdfda4e",
											"_945cd271-46b6-4d71-83a1-530e445af820",
											"_17db66a1-badd-4942-9ebd-02bc5595cdde",
											"_138f9ebc-0211-4051-b7c0-1c55695d5246",
											"_664f14a9-c1f1-490a-bbec-1f66ba4e7fe4",
											"_d22de266-6170-4783-91f9-40832e4cc58d",
											"_a4936291-3787-404c-bec7-8a3f3c5fd6e5",
											"_87785f46-7026-4d3c-b2c0-6a9468da67f6",
											"_a73027a7-615e-4a4d-95ee-c4cd78ab30c4",
											"_2b156883-2852-4665-aba0-d9bc57c7c225",
											"_9c5d383f-df57-4012-b490-fa36f9f90eed",
											"_3355cffe-aab4-4a05-8388-becf8ad599ae",
											"_be6ea91a-4f8e-4240-86e8-f85036aee96f",
											"_000a0565-911b-4f71-9993-1177021edd97",
											"_f006114d-c7cb-4ce0-9bfe-f0938c36a53e",
											"_b9338c62-a257-47dd-8c2e-88b80b73c330",
											"_b360104e-8410-4b99-827a-776e2083fb96",
											"_8055ae64-cafd-4fd0-be36-2216e3b02e37",
											"_2fd5c7d3-797d-45a5-a0d8-dfa60654ba5e",
											"_09074897-556d-4fd2-afb6-2f6c774e1820",
											"_54d66428-417b-447e-89d5-e726c1f12659",
											"_29b4f749-037a-4199-b33f-3cd3a3c7805e",
											"_3f3a831c-9b08-4827-92b3-3877a749e3df",
										},
									},
									{
										BaseElement: element.BaseElement{
											ID: "_af0d417a-6492-48d6-ace2-f70b807564a8",
										},
										Name: "Corporate Account Manager",
										FlowNodeRefs: []string{
											"_f0422f0d-396b-4ee7-ad83-fdd34a8bab71",
											"_a393c065-4412-4b16-b04f-52d07e76b518",
											"_05a1a66a-9308-41c7-a611-4fc57627a058",
											"_acfa265a-a449-4e30-baa0-1ba47dbea418",
										},
									},
									{
										BaseElement: element.BaseElement{
											ID: "_1c6c313d-4950-47a0-b6bf-c51a4b2ea7ed",
										},
										Name: "Head of Market Service",
										FlowNodeRefs: []string{
											"_1fc87527-9cad-4f8e-b9c7-ebe106cbe98d",
											"_5f56934b-8a7e-4c35-b9f7-bf2605711bfd",
											"_1da34f39-8338-4ecb-a93f-90349fa10260",
											"_1cf552d4-5152-4595-9218-84f31533bc70",
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
														ID: "_0254d83d-d943-466f-8b62-20e87cdfda4e",
													},
													Name: "Customer interested in Bank offer",
												},
												Outgoing: []string{"_bac89998-7084-4b2e-a2ea-dff839336d74"},
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
														ID: "_945cd271-46b6-4d71-83a1-530e445af820",
													},
													Name: "Interview customer",
												},
												Incoming: []string{"_bac89998-7084-4b2e-a2ea-dff839336d74"},
												Outgoing: []string{"_c132a484-8828-4f03-9a57-8201f4e9e864"},
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
														ID: "_17db66a1-badd-4942-9ebd-02bc5595cdde",
													},
													Name: "Prove/Provide identity",
												},
												Incoming: []string{"_c132a484-8828-4f03-9a57-8201f4e9e864"},
												Outgoing: []string{"_14c352ad-cd81-4a92-89e6-85f6e658e7bb"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_e75ca33b-71b4-405b-82b7-26b036834ad0",
															},
															Name: "ID document",
														},
													},
													InputSets: []element.InputSet{{}},
													OutputSets: []element.OutputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_2fd4dffa-bee0-48b8-a33e-a6fe462224d7",
															},
															DataOutputRefs: []string{"_e75ca33b-71b4-405b-82b7-26b036834ad0"},
														},
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_e5e5b852-0c5f-4f5f-8c6c-a1d8245024b4",
														},
														SourceRef: "_e75ca33b-71b4-405b-82b7-26b036834ad0",
														TargetRef: "_dca1cc5b-74a4-453f-8ad2-61609dde35ff",
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
														ID: "_664f14a9-c1f1-490a-bbec-1f66ba4e7fe4",
													},
													Name: "Obtain supporting data and documents of the customer",
												},
												Incoming: []string{"_b6414e79-5e61-4a8e-a584-65201c4ea9a9"},
												Outgoing: []string{"_8e88b32b-90d5-4ca4-a969-00b667a92d29"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_28d09a7b-8274-45b9-8ae1-66995807edf0",
															},
															Name: "ID document",
														},
													},
													InputSets: []element.InputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_91229d44-17f5-41db-ba5a-6041c238ad6a",
															},
															DataInputRefs: []string{"_28d09a7b-8274-45b9-8ae1-66995807edf0"},
														},
													},
													OutputSets: []element.OutputSet{{}},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_8a1de917-0cd0-4cb3-a449-5d10733575ee",
														},
														SourceRef: "_dca1cc5b-74a4-453f-8ad2-61609dde35ff",
														TargetRef: "_28d09a7b-8274-45b9-8ae1-66995807edf0",
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
														ID: "_d22de266-6170-4783-91f9-40832e4cc58d",
													},
													Name: "Check customer documents",
												},
												Incoming: []string{"_8e88b32b-90d5-4ca4-a969-00b667a92d29"},
												Outgoing: []string{"_c2f47722-9bb7-416c-b6bc-cff4e3a6eb65"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_f129e42a-e180-4068-970d-1895015a5a4e",
															},
															Name: "ID document",
														},
													},
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_59c497a7-8a58-41f7-ab3f-e823fc8843e6",
															},
															Name: "ID document",
														},
													},
													InputSets: []element.InputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_12650346-51b8-4431-807b-987777273d0b",
															},
															DataInputRefs: []string{"_f129e42a-e180-4068-970d-1895015a5a4e"},
														},
													},
													OutputSets: []element.OutputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_adf81738-a454-46bf-96b3-736985455e4c",
															},
															DataOutputRefs: []string{"_59c497a7-8a58-41f7-ab3f-e823fc8843e6"},
														},
													},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_e4d5f0c6-49e2-4467-80da-4818237bbbed",
														},
														SourceRef: "_dca1cc5b-74a4-453f-8ad2-61609dde35ff",
														TargetRef: "_f129e42a-e180-4068-970d-1895015a5a4e",
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_1cfb93bd-7a02-481d-a42c-c69c614a4b00",
														},
														SourceRef: "_59c497a7-8a58-41f7-ab3f-e823fc8843e6",
														TargetRef: "_e0b6d1fe-4ddf-4e22-b84c-6d676afb389a",
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
														ID: "_87785f46-7026-4d3c-b2c0-6a9468da67f6",
													},
													Name: "Copy, sign, and scan documents",
												},
												Incoming: []string{"_60066421-3da9-4318-8df2-d8b16b3011bb"},
												Outgoing: []string{"_65fba7be-af41-4ebe-b4e8-a8f7a830a702"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_da43ba9b-0964-4360-bba3-a63d040f0a76",
															},
															Name: "ID document",
														},
													},
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_581f63c2-d888-4329-9161-760518a3668d",
															},
															Name: "ID document",
														},
													},
													InputSets: []element.InputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_5b1eeb3a-06d7-4ab4-a50b-3507eaf6a96d",
															},
															DataInputRefs: []string{"_da43ba9b-0964-4360-bba3-a63d040f0a76"},
														},
													},
													OutputSets: []element.OutputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_b96d5da8-9e1e-4ba2-ad01-e6616577cd91",
															},
															DataOutputRefs: []string{"_581f63c2-d888-4329-9161-760518a3668d"},
														},
													},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_24b56549-d0fc-44d1-9fcf-b136ef8f420c",
														},
														SourceRef: "_e0b6d1fe-4ddf-4e22-b84c-6d676afb389a",
														TargetRef: "_da43ba9b-0964-4360-bba3-a63d040f0a76",
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_fc3badce-d027-4ab8-91a5-bfcf3be7daa5",
														},
														SourceRef: "_581f63c2-d888-4329-9161-760518a3668d",
														TargetRef: "_4f982ccc-c9e4-43b1-832e-b82992d4d400",
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
														ID: "_a73027a7-615e-4a4d-95ee-c4cd78ab30c4",
													},
													Name: "File documents in customer file",
												},
												Incoming: []string{"_65fba7be-af41-4ebe-b4e8-a8f7a830a702"},
												Outgoing: []string{"_0cf86177-f481-4c25-8801-cd18d91ddec6"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_0fed0a5a-5936-4ce7-aca5-472138e30238",
															},
															Name: "ID document",
														},
													},
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_4c2bb575-fa47-4f1c-9aed-e10b9633884f",
															},
															Name: "Customer Data (temporary storage)",
														},
													},
													InputSets: []element.InputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_a0daf95f-d885-47dc-b92b-bab63009296e",
															},
															DataInputRefs: []string{"_0fed0a5a-5936-4ce7-aca5-472138e30238"},
														},
													},
													OutputSets: []element.OutputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_e4533bfc-2668-4df5-b387-d4cf6576aa51",
															},
															DataOutputRefs: []string{"_4c2bb575-fa47-4f1c-9aed-e10b9633884f"},
														},
													},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_6518fef7-dc30-4737-84ee-f33ceba4a844",
														},
														SourceRef: "_4f982ccc-c9e4-43b1-832e-b82992d4d400",
														TargetRef: "_0fed0a5a-5936-4ce7-aca5-472138e30238",
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_ed10a8d7-01db-4221-a281-03174f633376",
														},
														SourceRef: "_4c2bb575-fa47-4f1c-9aed-e10b9633884f",
														TargetRef: "_c5285566-0657-47c5-a7ad-d09e144377f3",
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
														ID: "_9c5d383f-df57-4012-b490-fa36f9f90eed",
													},
													Name: "Add personal data",
												},
												Incoming: []string{"_8bd8a9d5-84da-4e96-b2af-9b8e003ad5a3"},
												Outgoing: []string{"_c7c63b1c-bdd1-443b-903b-7c31da1e6eb9"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_0ff75fcb-3aa9-4a0c-9bbf-64e31ffec23a",
															},
															Name: "Customer data",
														},
														{
															BaseElement: element.BaseElement{
																ID: "_ffecd813-8488-4289-8a5b-60ca006a7e6c",
															},
															Name: "Customer Data (temporary storage)",
														},
													},
													InputSets: []element.InputSet{{}},
													OutputSets: []element.OutputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_613918ad-e459-4b79-a85f-b34c1e875097",
															},
															DataOutputRefs: []string{"_0ff75fcb-3aa9-4a0c-9bbf-64e31ffec23a", "_ffecd813-8488-4289-8a5b-60ca006a7e6c"},
														},
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_a1577e31-32ee-46db-be07-c31f1f047a5e",
														},
														SourceRef: "_0ff75fcb-3aa9-4a0c-9bbf-64e31ffec23a",
														TargetRef: "_be9fc840-d916-4495-9061-989af6035030",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_aa4bda11-11a7-441a-b50f-7296b83542e4",
														},
														SourceRef: "_ffecd813-8488-4289-8a5b-60ca006a7e6c",
														TargetRef: "_c5285566-0657-47c5-a7ad-d09e144377f3",
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
														ID: "_be6ea91a-4f8e-4240-86e8-f85036aee96f",
													},
													Name: "Perform risk assessment of the customer",
												},
												Incoming: []string{"_1a933b26-7ac6-47a8-85de-07e658241e1e"},
												Outgoing: []string{"_40c01bbe-095a-4937-a334-95fa5b5225cc"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_721254f5-d2c9-4bc5-b1eb-2499e2fea042",
															},
															Name: "Customer data",
														},
													},
													DataOutputs: []element.DataOutput{

														{
															BaseElement: element.BaseElement{
																ID: "_7c351bcf-ed3d-481b-a696-254d32fda486",
															},
															Name: "Customer Data (temporary storage)",
														},
													},
													InputSets: []element.InputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_d0f0b35b-fc63-4e62-a45d-925fa79845cf",
															},
															DataInputRefs: []string{"_721254f5-d2c9-4bc5-b1eb-2499e2fea042"},
														},
													},
													OutputSets: []element.OutputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_e1621c49-620d-475b-9a23-3fc9c1800ee9",
															},
															DataOutputRefs: []string{"_7c351bcf-ed3d-481b-a696-254d32fda486"},
														},
													},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_9f68656f-415d-4288-be20-fd75bcd92f6a",
														},
														SourceRef: "_4257e72e-9c36-4832-8134-b68b8919de7e",
														TargetRef: "_721254f5-d2c9-4bc5-b1eb-2499e2fea042",
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_2fdc6abd-cfec-425d-b386-fe8b61371495",
														},
														SourceRef: "_7c351bcf-ed3d-481b-a696-254d32fda486",
														TargetRef: "_c5285566-0657-47c5-a7ad-d09e144377f3",
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
														ID: "_f006114d-c7cb-4ce0-9bfe-f0938c36a53e",
													},
													Name: "Document risk assessment",
												},
												Incoming: []string{"_7ab3e406-ce00-40e0-88f7-c64bcb0b864f"},
												Outgoing: []string{"_f2b37ee6-7727-4106-bdd4-d93c4db15703"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_7e58a64a-6396-4158-9a76-a0597a566f67",
															},
															Name: "Customer data",
														},
													},
													DataOutputs: []element.DataOutput{

														{
															BaseElement: element.BaseElement{
																ID: "_e3b88a3e-9f69-4bd4-ada3-e277fd53656e",
															},
															Name: "Customer Data (temporary storage)",
														},
													},
													InputSets: []element.InputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_b0706d89-92f6-481e-89f0-4a9ad5a1b1d5",
															},
															DataInputRefs: []string{"_7e58a64a-6396-4158-9a76-a0597a566f67"},
														},
													},
													OutputSets: []element.OutputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_f0d63836-16aa-4f31-a122-01dc1fd556cb",
															},
															DataOutputRefs: []string{"_e3b88a3e-9f69-4bd4-ada3-e277fd53656e"},
														},
													},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_a8b81881-5b94-4d11-8b5c-5e9b76878dcd",
														},
														SourceRef: "_4257e72e-9c36-4832-8134-b68b8919de7e",
														TargetRef: "_7e58a64a-6396-4158-9a76-a0597a566f67",
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_c7006cb5-a9f4-4528-8fca-9083fc7559ed",
														},
														SourceRef: "_e3b88a3e-9f69-4bd4-ada3-e277fd53656e",
														TargetRef: "_c5285566-0657-47c5-a7ad-d09e144377f3",
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
														ID: "_b360104e-8410-4b99-827a-776e2083fb96",
													},
													Name: "Create customer in the system",
												},
												Incoming: []string{"_f974945f-f2a6-4547-9d35-0d26049dcd6f"},
												Outgoing: []string{"_0b3ae343-49aa-49dc-8b88-fe97b4ada3eb"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_42e1d578-55c2-4ea6-9b6f-a32ee680ec72",
															},
															Name: "Customer data",
														},
														{
															BaseElement: element.BaseElement{
																ID: "_e4a8fda6-25e5-4d3a-af66-6bdc77e82aed",
															},
															Name: "Customer Data (temporary storage)",
														},
													},
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_174d4bae-bb40-4877-86d7-6e1f88b9353e",
															},
															Name: "Bank System",
														},
													},
													InputSets: []element.InputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_2a408dec-555c-4fa0-a9cc-08bad960c9c8",
															},
															DataInputRefs: []string{"_42e1d578-55c2-4ea6-9b6f-a32ee680ec72", "_e4a8fda6-25e5-4d3a-af66-6bdc77e82aed"},
														},
													},
													OutputSets: []element.OutputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_41ca4899-93f5-4f71-b5c0-6aa3f91ed081",
															},
															DataOutputRefs: []string{"_174d4bae-bb40-4877-86d7-6e1f88b9353e"},
														},
													},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_bd748ab9-bd4e-47d1-8f69-cc8e451e0b2b",
														},
														SourceRef: "_dd653cc9-c874-4cfe-ba6c-146726809ea5",
														TargetRef: "_42e1d578-55c2-4ea6-9b6f-a32ee680ec72",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_48edea23-f132-4ec7-b296-35c091e1755b",
														},
														SourceRef: "_66f4a1ff-fd8a-4002-ade7-7a4d1b1d529c",
														TargetRef: "_e4a8fda6-25e5-4d3a-af66-6bdc77e82aed",
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_8eaa28a5-123c-4224-ad2a-218525c6c79e",
														},
														SourceRef: "_174d4bae-bb40-4877-86d7-6e1f88b9353e",
														TargetRef: "_ffd44905-4f54-4fe0-a156-0312c2a7a3cc",
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
														ID: "_2fd5c7d3-797d-45a5-a0d8-dfa60654ba5e",
													},
													Name: "Complete data and documents",
												},
												Incoming: []string{"_664f3a71-3efa-4c94-8797-815f2e377cf7"},
												Outgoing: []string{"_52e3c106-b055-40a0-a9f9-db9529adfb38"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_f0fd8c00-f695-4115-87b6-6aa5eaa1f7c4",
															},
															Name: "ID document",
														},
													},
													InputSets: []element.InputSet{{}},
													OutputSets: []element.OutputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_9d7b0978-2803-4fef-b491-bec904315432",
															},
															DataOutputRefs: []string{"_f0fd8c00-f695-4115-87b6-6aa5eaa1f7c4"},
														},
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_a290cb0f-07ce-4e95-8d51-eb98307b3beb",
														},
														SourceRef: "_f0fd8c00-f695-4115-87b6-6aa5eaa1f7c4",
														TargetRef: "_e0b6d1fe-4ddf-4e22-b84c-6d676afb389a",
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
														ID: "_09074897-556d-4fd2-afb6-2f6c774e1820",
													},
													Name: "Perform know your customer (KYC) activities",
												},
												Incoming: []string{"_f47c4e3f-3339-4c96-b8b0-829022bada78"},
												Outgoing: []string{"_a4eaa7b6-f9ef-458a-ba70-e1fd3ff0595a"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_3fcd7c6a-65c8-434d-9998-5c8722c6810b",
															},
															Name: "Customer data",
														},
													},
													DataOutputs: []element.DataOutput{
														{
															BaseElement: element.BaseElement{
																ID: "_a130c3b1-aa52-47f6-8818-fc567f559ce2",
															},
															Name: "Customer data",
														},
														{
															BaseElement: element.BaseElement{
																ID: "_b6f553cf-dc70-4169-a2a0-0943cb92953a",
															},
															Name: "Customer Data (temporary storage)",
														},
													},
													InputSets: []element.InputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_3d521f50-0d38-4e83-9992-f8910130aa81",
															},
															DataInputRefs: []string{"_3fcd7c6a-65c8-434d-9998-5c8722c6810b"},
														},
													},
													OutputSets: []element.OutputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_0e1a0a39-c100-4eab-91fc-bef562e23e68",
															},
															DataOutputRefs: []string{"_a130c3b1-aa52-47f6-8818-fc567f559ce2", "_b6f553cf-dc70-4169-a2a0-0943cb92953a"},
														},
													},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_b181f888-d411-4665-aeea-89569a373e98",
														},
														SourceRef: "_be9fc840-d916-4495-9061-989af6035030",
														TargetRef: "_3fcd7c6a-65c8-434d-9998-5c8722c6810b",
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_8234f225-632d-4b63-b7b4-9ee05ebd4971",
														},
														SourceRef: "_a130c3b1-aa52-47f6-8818-fc567f559ce2",
														TargetRef: "_4257e72e-9c36-4832-8134-b68b8919de7e",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_9f70cdff-8fb7-432d-84f0-29db168610fa",
														},
														SourceRef: "_b6f553cf-dc70-4169-a2a0-0943cb92953a",
														TargetRef: "_c5285566-0657-47c5-a7ad-d09e144377f3",
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
														ID: "_f0422f0d-396b-4ee7-ad83-fdd34a8bab71",
													},
													Name: "Document the identity of the economic owner",
												},
												Incoming: []string{"_6912e5ed-5cae-4798-b6e6-8d367ddb3f60"},
												Outgoing: []string{"_a8da83e1-4e09-4a7b-a09f-1f869b372176"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_9d9659ca-4d33-4d6a-903a-dfeed15cee9a",
															},
															Name: "ID document",
														},
													},
													InputSets: []element.InputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_0ded1261-c1bd-49de-8f5b-3d46c4d5a095",
															},
															DataInputRefs: []string{"_9d9659ca-4d33-4d6a-903a-dfeed15cee9a"},
														},
													},
													OutputSets: []element.OutputSet{{}},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_e470d57d-3556-40dd-b518-5ed19b3f63d9",
														},
														SourceRef: "_dca1cc5b-74a4-453f-8ad2-61609dde35ff",
														TargetRef: "_9d9659ca-4d33-4d6a-903a-dfeed15cee9a",
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
														ID: "_05a1a66a-9308-41c7-a611-4fc57627a058",
													},
													Name: "End business relation",
												},
												Incoming: []string{"_470d1a12-e1fa-4471-b753-cfd7dbe20d80"},
												Outgoing: []string{"_461976fc-e457-4a39-baf3-785885c58896"},
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
														ID: "_1fc87527-9cad-4f8e-b9c7-ebe106cbe98d",
													},
													Name: "Check risk and decide about approval",
												},
												Incoming: []string{"_e88d64c7-3aaf-4a5f-9787-4e5ba696312b"},
												Outgoing: []string{"_c3e90e5c-e5cf-4cfb-b15d-bdedcf0719ca"},
											},
											IoSpecification: element.IoSpecification{
												InputOutputSpecification: element.InputOutputSpecification{
													DataInputs: []element.DataInput{
														{
															BaseElement: element.BaseElement{
																ID: "_85c0085f-9a81-472b-a54c-51943f4c247c",
															},
															Name: "Customer data",
														},
													},
													InputSets: []element.InputSet{
														{
															BaseElement: element.BaseElement{
																ID: "_312af7f4-e30b-4606-aeb3-dd5008632579",
															},
															DataInputRefs: []string{"_85c0085f-9a81-472b-a54c-51943f4c247c"},
														},
													},
													OutputSets: []element.OutputSet{{}},
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "_f4536257-490a-4f11-8174-59c44fe6d51d",
														},
														SourceRef: "_4257e72e-9c36-4832-8134-b68b8919de7e",
														TargetRef: "_85c0085f-9a81-472b-a54c-51943f4c247c",
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
													ID: "_138f9ebc-0211-4051-b7c0-1c55695d5246",
												},
												Name: "Legal entity or individual?",
											},
											Incoming: []string{"_14c352ad-cd81-4a92-89e6-85f6e658e7bb"},
											Outgoing: []string{"_fcb09e30-bfe6-46b9-af01-6777c60026f2", "_6912e5ed-5cae-4798-b6e6-8d367ddb3f60"},
										},
										GatewayDirection: "Diverging",
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_a4936291-3787-404c-bec7-8a3f3c5fd6e5",
												},
												Name: "Data complete?",
											},
											Incoming: []string{"_c2f47722-9bb7-416c-b6bc-cff4e3a6eb65"},
											Outgoing: []string{"_eec53048-291a-4802-8fec-012857f845eb", "_664f3a71-3efa-4c94-8797-815f2e377cf7"},
										},
										GatewayDirection: "Diverging",
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_000a0565-911b-4f71-9993-1177021edd97",
												},
												Name: "Subject to approval?",
											},
											Incoming: []string{"_40c01bbe-095a-4937-a334-95fa5b5225cc"},
											Outgoing: []string{"_1887adee-8291-4640-a5fa-3dc5f33d5a34", "_e88d64c7-3aaf-4a5f-9787-4e5ba696312b"},
										},
										GatewayDirection: "Diverging",
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_54d66428-417b-447e-89d5-e726c1f12659",
												},
											},
											Incoming: []string{"_fcb09e30-bfe6-46b9-af01-6777c60026f2", "_afc35597-a482-432a-8ace-17adaf609572"},
											Outgoing: []string{"_b6414e79-5e61-4a8e-a584-65201c4ea9a9"},
										},
										GatewayDirection: "Converging",
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_29b4f749-037a-4199-b33f-3cd3a3c7805e",
												},
											},
											Incoming: []string{"_eec53048-291a-4802-8fec-012857f845eb", "_52e3c106-b055-40a0-a9f9-db9529adfb38"},
											Outgoing: []string{"_60066421-3da9-4318-8df2-d8b16b3011bb"},
										},
										GatewayDirection: "Converging",
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_3f3a831c-9b08-4827-92b3-3877a749e3df",
												},
											},
											Incoming: []string{"_1887adee-8291-4640-a5fa-3dc5f33d5a34", "_8a77d7f6-320a-47ff-a155-57ee200df478"},
											Outgoing: []string{"_7ab3e406-ce00-40e0-88f7-c64bcb0b864f"},
										},
										GatewayDirection: "Converging",
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_a393c065-4412-4b16-b04f-52d07e76b518",
												},
												Name: "Identity of the economic owner certifiable?",
											},
											Incoming: []string{"_a8da83e1-4e09-4a7b-a09f-1f869b372176"},
											Outgoing: []string{"_afc35597-a482-432a-8ace-17adaf609572", "_470d1a12-e1fa-4471-b753-cfd7dbe20d80"},
										},
										GatewayDirection: "Diverging",
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_5f56934b-8a7e-4c35-b9f7-bf2605711bfd",
												},
												Name: "Approval?",
											},
											Incoming: []string{"_c3e90e5c-e5cf-4cfb-b15d-bdedcf0719ca"},
											Outgoing: []string{"_8a77d7f6-320a-47ff-a155-57ee200df478", "_ae94e0c1-aa12-433f-bd44-a4d8a70a089e"},
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
													ID: "_2b156883-2852-4665-aba0-d9bc57c7c225",
												},
											},
											Incoming: []string{"_0cf86177-f481-4c25-8801-cd18d91ddec6"},
											Outgoing: []string{"_8bd8a9d5-84da-4e96-b2af-9b8e003ad5a3", "_f47c4e3f-3339-4c96-b8b0-829022bada78"},
										},
										GatewayDirection: "Diverging",
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_3355cffe-aab4-4a05-8388-becf8ad599ae",
												},
											},
											Incoming: []string{"_c7c63b1c-bdd1-443b-903b-7c31da1e6eb9", "_a4eaa7b6-f9ef-458a-ba70-e1fd3ff0595a"},
											Outgoing: []string{"_1a933b26-7ac6-47a8-85de-07e658241e1e"},
										},
										GatewayDirection: "Converging",
									},
								},
							},
							CallActivities: []element.CallActivity{
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_b9338c62-a257-47dd-8c2e-88b80b73c330",
												},
												Name: "Check for connected clients",
											},
											Incoming: []string{"_f2b37ee6-7727-4106-bdd4-d93c4db15703"},
											Outgoing: []string{"_f974945f-f2a6-4547-9d35-0d26049dcd6f"},
										},
										IoSpecification: element.IoSpecification{
											InputOutputSpecification: element.InputOutputSpecification{
												DataInputs: []element.DataInput{
													{
														BaseElement: element.BaseElement{
															ID: "_32c34c33-47cf-4481-8310-401ce5b04a35",
														},
														Name: "Customer data",
													},
												},
												DataOutputs: []element.DataOutput{

													{
														BaseElement: element.BaseElement{
															ID: "_f87c9123-3f07-420c-ae9f-4e29e3af5130",
														},
														Name: "Customer data",
													},
												},
												InputSets: []element.InputSet{
													{
														BaseElement: element.BaseElement{
															ID: "_2d5efe68-8501-481d-96d4-6ab41ed2d4a1",
														},
														DataInputRefs: []string{"_32c34c33-47cf-4481-8310-401ce5b04a35"},
													},
												},
												OutputSets: []element.OutputSet{
													{
														BaseElement: element.BaseElement{
															ID: "_7ab57802-2e41-4a0e-90de-cfd521600034",
														},
														DataOutputRefs: []string{"_f87c9123-3f07-420c-ae9f-4e29e3af5130"},
													},
												},
											},
										},
										DataInputAssociations: []element.DataInputAssociation{
											{
												DataAssociation: element.DataAssociation{
													BaseElement: element.BaseElement{
														ID: "_10466056-a9bd-4d46-b6ca-60cc6c5f6b13",
													},
													SourceRef: "_4257e72e-9c36-4832-8134-b68b8919de7e",
													TargetRef: "_32c34c33-47cf-4481-8310-401ce5b04a35",
												},
											},
										},
										DataOutputAssociations: []element.DataOutputAssociation{
											{
												DataAssociation: element.DataAssociation{
													BaseElement: element.BaseElement{
														ID: "_39450d6b-b008-4af2-bb98-35d65093922e",
													},
													SourceRef: "_f87c9123-3f07-420c-ae9f-4e29e3af5130",
													TargetRef: "_dd653cc9-c874-4cfe-ba6c-146726809ea5",
												},
											},
										},
									},
									CalledElement: "_774bc005-0917-43d5-ab70-0f9fe123fbd1",
								},
							},
							EndEvents: []element.EndEvent{
								{
									ThrowEvent: element.ThrowEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_8055ae64-cafd-4fd0-be36-2216e3b02e37",
													},
													Name: "Identity determined and new customer created",
												},
												Incoming: []string{"_0b3ae343-49aa-49dc-8b88-fe97b4ada3eb"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											SignalEventDefinitions: []element.SignalEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "_8c6f817f-0d1e-464a-96bc-e32dbfa8110d",
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
														ID: "_acfa265a-a449-4e30-baa0-1ba47dbea418",
													},
													Name: "Business relation ended",
												},
												Incoming: []string{"_461976fc-e457-4a39-baf3-785885c58896"},
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
														ID: "_1cf552d4-5152-4595-9218-84f31533bc70",
													},
													Name: "No business relation created",
												},
												Incoming: []string{"_f7820843-8922-430e-ab4a-974a1e35d89b"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											SignalEventDefinitions: []element.SignalEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "_8374dba6-01d8-45b8-9703-103b2fdc79a7",
															},
														},
													},
												},
											},
										},
									},
								},
							},
							DataObjectReferenes: []element.DataObjectReference{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_dca1cc5b-74a4-453f-8ad2-61609dde35ff",
										},
										Name: "ID document [for analysis]",
									},
									DataObjectRef: "_0fa2af52-a3d6-4be6-9ecc-9fad8ca4d448",
									DataState: element.DataState{
										Name: "for analysis",
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_e0b6d1fe-4ddf-4e22-b84c-6d676afb389a",
										},
										Name: "ID document [analysed]",
									},
									DataObjectRef: "_0fa2af52-a3d6-4be6-9ecc-9fad8ca4d448",
									DataState: element.DataState{
										Name: "analysed",
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_4f982ccc-c9e4-43b1-832e-b82992d4d400",
										},
										Name: "ID document [scanned]",
									},
									DataObjectRef: "_0fa2af52-a3d6-4be6-9ecc-9fad8ca4d448",
									DataState: element.DataState{
										Name: "scanned",
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_be9fc840-d916-4495-9061-989af6035030",
										},
										Name: "Customer data [for KYC analysis]",
									},
									DataObjectRef: "_79993b0e-60c2-487b-a12e-0cbecc2ef5c1",
									DataState: element.DataState{
										Name: "for KYC analysis",
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_4257e72e-9c36-4832-8134-b68b8919de7e",
										},
										Name: "Customer data [for risk assessment]",
									},
									DataObjectRef: "_79993b0e-60c2-487b-a12e-0cbecc2ef5c1",
									DataState: element.DataState{
										Name: "for risk assessment",
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_dd653cc9-c874-4cfe-ba6c-146726809ea5",
										},
										Name: "Customer data [with connected clients analysis]",
									},
									DataObjectRef: "_79993b0e-60c2-487b-a12e-0cbecc2ef5c1",
									DataState: element.DataState{
										Name: "with connected clients analysis",
									},
								},
							},
							DataStoreReferences: []element.DataStoreReference{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_ffd44905-4f54-4fe0-a156-0312c2a7a3cc",
										},
										Name: "Bank System",
									},
									DataStoreRef: "_776bd6ca-5f18-432f-a748-f1930bfdf16e",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_66f4a1ff-fd8a-4002-ade7-7a4d1b1d529c",
										},
										Name: "Customer Data (temporary storage)",
									},
									DataStoreRef: "_46aee7ee-fab2-4735-80cc-fa78092ef92e",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_c5285566-0657-47c5-a7ad-d09e144377f3",
										},
										Name: "Customer Data (temporary storage)",
									},
									DataStoreRef: "_46aee7ee-fab2-4735-80cc-fa78092ef92e",
								},
							},
							Tasks: []element.Task{
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_1da34f39-8338-4ecb-a93f-90349fa10260",
												},
												Name: "Reject customer request",
											},
											Incoming: []string{"_ae94e0c1-aa12-433f-bd44-a4d8a70a089e"},
											Outgoing: []string{"_f7820843-8922-430e-ab4a-974a1e35d89b"},
										},
									},
								},
							},
							DataObjects: []element.DataObject{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_0fa2af52-a3d6-4be6-9ecc-9fad8ca4d448",
										},
										Name: "ID document",
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_79993b0e-60c2-487b-a12e-0cbecc2ef5c1",
										},
										Name: "Customer data",
									},
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_bac89998-7084-4b2e-a2ea-dff839336d74",
										},
									},
									SourceRef: "_0254d83d-d943-466f-8b62-20e87cdfda4e",
									TargetRef: "_945cd271-46b6-4d71-83a1-530e445af820",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_c132a484-8828-4f03-9a57-8201f4e9e864",
										},
									},
									SourceRef: "_945cd271-46b6-4d71-83a1-530e445af820",
									TargetRef: "_17db66a1-badd-4942-9ebd-02bc5595cdde",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_14c352ad-cd81-4a92-89e6-85f6e658e7bb",
										},
									},
									SourceRef: "_17db66a1-badd-4942-9ebd-02bc5595cdde",
									TargetRef: "_138f9ebc-0211-4051-b7c0-1c55695d5246",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_fcb09e30-bfe6-46b9-af01-6777c60026f2",
										},
										Name: "Individual Person",
									},
									SourceRef: "_138f9ebc-0211-4051-b7c0-1c55695d5246",
									TargetRef: "_54d66428-417b-447e-89d5-e726c1f12659",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_8e88b32b-90d5-4ca4-a969-00b667a92d29",
										},
									},
									SourceRef: "_664f14a9-c1f1-490a-bbec-1f66ba4e7fe4",
									TargetRef: "_d22de266-6170-4783-91f9-40832e4cc58d",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_c2f47722-9bb7-416c-b6bc-cff4e3a6eb65",
										},
									},
									SourceRef: "_d22de266-6170-4783-91f9-40832e4cc58d",
									TargetRef: "_a4936291-3787-404c-bec7-8a3f3c5fd6e5",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_eec53048-291a-4802-8fec-012857f845eb",
										},
										Name: "Yes",
									},
									SourceRef: "_a4936291-3787-404c-bec7-8a3f3c5fd6e5",
									TargetRef: "_29b4f749-037a-4199-b33f-3cd3a3c7805e",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_65fba7be-af41-4ebe-b4e8-a8f7a830a702",
										},
									},
									SourceRef: "_87785f46-7026-4d3c-b2c0-6a9468da67f6",
									TargetRef: "_a73027a7-615e-4a4d-95ee-c4cd78ab30c4",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_0cf86177-f481-4c25-8801-cd18d91ddec6",
										},
									},
									SourceRef: "_a73027a7-615e-4a4d-95ee-c4cd78ab30c4",
									TargetRef: "_2b156883-2852-4665-aba0-d9bc57c7c225",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_8bd8a9d5-84da-4e96-b2af-9b8e003ad5a3",
										},
									},
									SourceRef: "_2b156883-2852-4665-aba0-d9bc57c7c225",
									TargetRef: "_9c5d383f-df57-4012-b490-fa36f9f90eed",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_c7c63b1c-bdd1-443b-903b-7c31da1e6eb9",
										},
									},
									SourceRef: "_9c5d383f-df57-4012-b490-fa36f9f90eed",
									TargetRef: "_3355cffe-aab4-4a05-8388-becf8ad599ae",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_1a933b26-7ac6-47a8-85de-07e658241e1e",
										},
									},
									SourceRef: "_3355cffe-aab4-4a05-8388-becf8ad599ae",
									TargetRef: "_be6ea91a-4f8e-4240-86e8-f85036aee96f",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_40c01bbe-095a-4937-a334-95fa5b5225cc",
										},
									},
									SourceRef: "_be6ea91a-4f8e-4240-86e8-f85036aee96f",
									TargetRef: "_000a0565-911b-4f71-9993-1177021edd97",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_1887adee-8291-4640-a5fa-3dc5f33d5a34",
										},
										Name: "No",
									},
									SourceRef: "_000a0565-911b-4f71-9993-1177021edd97",
									TargetRef: "_3f3a831c-9b08-4827-92b3-3877a749e3df",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_f2b37ee6-7727-4106-bdd4-d93c4db15703",
										},
									},
									SourceRef: "_f006114d-c7cb-4ce0-9bfe-f0938c36a53e",
									TargetRef: "_b9338c62-a257-47dd-8c2e-88b80b73c330",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_f974945f-f2a6-4547-9d35-0d26049dcd6f",
										},
									},
									SourceRef: "_b9338c62-a257-47dd-8c2e-88b80b73c330",
									TargetRef: "_b360104e-8410-4b99-827a-776e2083fb96",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_0b3ae343-49aa-49dc-8b88-fe97b4ada3eb",
										},
									},
									SourceRef: "_b360104e-8410-4b99-827a-776e2083fb96",
									TargetRef: "_8055ae64-cafd-4fd0-be36-2216e3b02e37",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_664f3a71-3efa-4c94-8797-815f2e377cf7",
										},
										Name: "No",
									},
									SourceRef: "_a4936291-3787-404c-bec7-8a3f3c5fd6e5",
									TargetRef: "_2fd5c7d3-797d-45a5-a0d8-dfa60654ba5e",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_52e3c106-b055-40a0-a9f9-db9529adfb38",
										},
									},
									SourceRef: "_2fd5c7d3-797d-45a5-a0d8-dfa60654ba5e",
									TargetRef: "_29b4f749-037a-4199-b33f-3cd3a3c7805e",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_f47c4e3f-3339-4c96-b8b0-829022bada78",
										},
									},
									SourceRef: "_2b156883-2852-4665-aba0-d9bc57c7c225",
									TargetRef: "_09074897-556d-4fd2-afb6-2f6c774e1820",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_a4eaa7b6-f9ef-458a-ba70-e1fd3ff0595a",
										},
									},
									SourceRef: "_09074897-556d-4fd2-afb6-2f6c774e1820",
									TargetRef: "_3355cffe-aab4-4a05-8388-becf8ad599ae",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_6912e5ed-5cae-4798-b6e6-8d367ddb3f60",
										},
										Name: "Legal Entity",
									},
									SourceRef: "_138f9ebc-0211-4051-b7c0-1c55695d5246",
									TargetRef: "_f0422f0d-396b-4ee7-ad83-fdd34a8bab71",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_a8da83e1-4e09-4a7b-a09f-1f869b372176",
										},
									},
									SourceRef: "_f0422f0d-396b-4ee7-ad83-fdd34a8bab71",
									TargetRef: "_a393c065-4412-4b16-b04f-52d07e76b518",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_afc35597-a482-432a-8ace-17adaf609572",
										},
										Name: "Yes",
									},
									SourceRef: "_a393c065-4412-4b16-b04f-52d07e76b518",
									TargetRef: "_54d66428-417b-447e-89d5-e726c1f12659",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_470d1a12-e1fa-4471-b753-cfd7dbe20d80",
										},
										Name: "No",
									},
									SourceRef: "_a393c065-4412-4b16-b04f-52d07e76b518",
									TargetRef: "_05a1a66a-9308-41c7-a611-4fc57627a058",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_461976fc-e457-4a39-baf3-785885c58896",
										},
									},
									SourceRef: "_05a1a66a-9308-41c7-a611-4fc57627a058",
									TargetRef: "_acfa265a-a449-4e30-baa0-1ba47dbea418",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_e88d64c7-3aaf-4a5f-9787-4e5ba696312b",
										},
										Name: "Yes",
									},
									SourceRef: "_000a0565-911b-4f71-9993-1177021edd97",
									TargetRef: "_1fc87527-9cad-4f8e-b9c7-ebe106cbe98d",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_c3e90e5c-e5cf-4cfb-b15d-bdedcf0719ca",
										},
									},
									SourceRef: "_1fc87527-9cad-4f8e-b9c7-ebe106cbe98d",
									TargetRef: "_5f56934b-8a7e-4c35-b9f7-bf2605711bfd",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_8a77d7f6-320a-47ff-a155-57ee200df478",
										},
										Name: "Yes",
									},
									SourceRef: "_5f56934b-8a7e-4c35-b9f7-bf2605711bfd",
									TargetRef: "_3f3a831c-9b08-4827-92b3-3877a749e3df",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_ae94e0c1-aa12-433f-bd44-a4d8a70a089e",
										},
										Name: "No",
									},
									SourceRef: "_5f56934b-8a7e-4c35-b9f7-bf2605711bfd",
									TargetRef: "_1da34f39-8338-4ecb-a93f-90349fa10260",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_f7820843-8922-430e-ab4a-974a1e35d89b",
										},
									},
									SourceRef: "_1da34f39-8338-4ecb-a93f-90349fa10260",
									TargetRef: "_1cf552d4-5152-4595-9218-84f31533bc70",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_b6414e79-5e61-4a8e-a584-65201c4ea9a9",
										},
									},
									SourceRef: "_54d66428-417b-447e-89d5-e726c1f12659",
									TargetRef: "_664f14a9-c1f1-490a-bbec-1f66ba4e7fe4",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_60066421-3da9-4318-8df2-d8b16b3011bb",
										},
									},
									SourceRef: "_29b4f749-037a-4199-b33f-3cd3a3c7805e",
									TargetRef: "_87785f46-7026-4d3c-b2c0-6a9468da67f6",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_7ab3e406-ce00-40e0-88f7-c64bcb0b864f",
										},
									},
									SourceRef: "_3f3a831c-9b08-4827-92b3-3877a749e3df",
									TargetRef: "_f006114d-c7cb-4ce0-9bfe-f0938c36a53e",
								},
							},
						},
					},
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "_774bc005-0917-43d5-ab70-0f9fe123fbd1",
								},
							},
							Name: "Check for connected clients",
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
														ID: "_d8214574-bb4c-42ff-aabb-398eb95b2f2a",
													},
													Name: "Check for connected clients",
												},
												Outgoing: []string{"_6c121fa1-4bab-4e82-98a6-672c6a6bb6f3"},
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
														ID: "_8b104885-149e-4af6-a459-d924dacd81b3",
													},
													Name: "Check if group of connected clients exists",
												},
												Incoming: []string{"_6c121fa1-4bab-4e82-98a6-672c6a6bb6f3"},
												Outgoing: []string{"_4e00ba20-d2c5-47d6-b53f-b50d802855fd"},
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
														ID: "_7507ae41-a1fa-405c-b4ea-85ed920eace5",
													},
													Name: "Document group of connected clients according to Capital Requirements Regulation (CRR)",
												},
												Incoming: []string{"_e25b9c2d-690e-470b-8993-112002994fc2"},
												Outgoing: []string{"_20c68a0e-6aba-445b-94a3-6bb588dfa2d8"},
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
													ID: "_080399c9-3c91-44c6-b510-80367e23a5af",
												},
												Name: "Group of connected Clients existing?",
											},
											Incoming: []string{"_4e00ba20-d2c5-47d6-b53f-b50d802855fd"},
											Outgoing: []string{"_e25b9c2d-690e-470b-8993-112002994fc2", "_c329c58d-4a71-471f-acb9-b7f7113d8547"},
										},
										GatewayDirection: "Diverging",
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_956bb101-c9f7-467d-b3e2-198fa1d3e12b",
												},
											},
											Incoming: []string{"_20c68a0e-6aba-445b-94a3-6bb588dfa2d8", "_c329c58d-4a71-471f-acb9-b7f7113d8547"},
											Outgoing: []string{"_90d759bb-90dc-48b7-8c06-ddf8664c36ff"},
										},
										GatewayDirection: "Converging",
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
														ID: "_f7ce4bda-22c2-4ef9-aad9-5203dff18538",
													},
												},
												Incoming: []string{"_90d759bb-90dc-48b7-8c06-ddf8664c36ff"},
											},
										},
									},
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_6c121fa1-4bab-4e82-98a6-672c6a6bb6f3",
										},
									},
									SourceRef: "_d8214574-bb4c-42ff-aabb-398eb95b2f2a",
									TargetRef: "_8b104885-149e-4af6-a459-d924dacd81b3",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_4e00ba20-d2c5-47d6-b53f-b50d802855fd",
										},
									},
									SourceRef: "_8b104885-149e-4af6-a459-d924dacd81b3",
									TargetRef: "_080399c9-3c91-44c6-b510-80367e23a5af",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_e25b9c2d-690e-470b-8993-112002994fc2",
										},
										Name: "Yes",
									},
									SourceRef: "_080399c9-3c91-44c6-b510-80367e23a5af",
									TargetRef: "_7507ae41-a1fa-405c-b4ea-85ed920eace5",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_20c68a0e-6aba-445b-94a3-6bb588dfa2d8",
										},
									},
									SourceRef: "_7507ae41-a1fa-405c-b4ea-85ed920eace5",
									TargetRef: "_956bb101-c9f7-467d-b3e2-198fa1d3e12b",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_90d759bb-90dc-48b7-8c06-ddf8664c36ff",
										},
									},
									SourceRef: "_956bb101-c9f7-467d-b3e2-198fa1d3e12b",
									TargetRef: "_f7ce4bda-22c2-4ef9-aad9-5203dff18538",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_c329c58d-4a71-471f-acb9-b7f7113d8547",
										},
										Name: "No",
									},
									SourceRef: "_080399c9-3c91-44c6-b510-80367e23a5af",
									TargetRef: "_956bb101-c9f7-467d-b3e2-198fa1d3e12b",
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
