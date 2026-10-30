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

func TestC_3_0_roundtrip(t *testing.T) {
	// create test use ./test/C.3.0-roundtrip.bpmn
	path := "./test/C.3.0-roundtrip.bpmn"
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
			ID:              "_5033f8b0-0396-494b-a52c-787034420443",
			Name:            "C.3.0",
			TargetNamespace: "http://www.itp-commerce.com",
			Exporter:        "W4_BPMN_Composer",
			ExporterVersion: "9.2.0.0",
			RootElemnts: element.RootElemnts{
				Processes: []element.Process{
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "_8170787a-3207-434d-9bea-4787059f444f",
									ExtensionElements: element.ExtensionElements{
										Any: []xml.Token{nil, nil},
									},
								},
							},
							Name: "Fridge Repair Process",
						},
						ProcessType:  "Private",
						IsExecutable: true,
						FlowElements: element.FlowElements{
							StartEvents: []element.StartEvent{
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID:                "_cc9778bd-edd8-4df2-ba15-56c310f90e62",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
													},
													Name: "Receive customer request",
												},
												Outgoing: []string{"_a5af06ae-bd69-464d-bbaf-3d7418702d77"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "Bpmn_MessageEventDefinition_vaRFsAO9EeWi8fS3WiOizw",
															},
														},
													},
													MessageRef: "_6840c92f-f02c-46b3-9565-f5e5b1792ba0",
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
														ID: "_c73a5f4a-72f1-4e11-bb40-2f98da75fb9a",
														Documentation: []element.Documentation{
															{
																Id:    "Bpmn_Documentation_jcmzwe_bEeSGoscwBjzAjw",
																Value: "First contact with customer",
															},
														},
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
													},
													Name: "Analyse customer request",
												},
												Incoming: []string{"_a5af06ae-bd69-464d-bbaf-3d7418702d77"},
												Outgoing: []string{"_acb2aca3-8851-48f0-b127-7b3c9db5e18d"},
											},
											StartQuantity:      2,
											CompletionQuantity: 2,
											PotentialOwners: []element.PotentialOwner{
												{
													HumanPerformer: element.HumanPerformer{
														Performer: element.Performer{
															ResourceRole: element.ResourceRole{
																BaseElement: element.BaseElement{
																	ID:  "Bpmn_PotentialOwner_H_12wBqHEeWDuOtG0oS24A",
																	Any: nil,
																},
																Name:        "Potential Owner",
																ResourceRef: "Bpmn_Resource__7wrkBqGEeWDuOtG0oS24A",
															},
														},
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
														ID: "_a92069f7-377b-4dbd-a1fd-1da071aabf6d",
														Documentation: []element.Documentation{
															{
																Id:    "Bpmn_Documentation_jcna0O_bEeSGoscwBjzAjw",
																Value: "Fridge replaced if still under warranty.",
															},
														},
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
													},
													Name: "Replace fridge",
												},
												Incoming: []string{"_b99800c3-c340-460c-a43e-098014a365d0", "_3fb323d5-2c59-487a-af63-804208f6c5cb"},
												Outgoing: []string{"_c5756bb9-6e6f-42d1-8799-c2d673499eb8"},
											},
										},
									},
									Implementation: "WebService",
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_d034722f-751d-4f37-a3d7-47993822e979",
														Documentation: []element.Documentation{
															{
																Id:    "Bpmn_Documentation_jcoo8O_bEeSGoscwBjzAjw",
																Value: "Standard service level is common for most customer support.\nThis level guarantees the fridge will be repaired in a reasonable amount of time.\nCustomers may decide to raise the service level to Premium while waiting for the standard service to solve the issue.",
															},
														},
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
													},
													Name: "Perform repair (standard level)",
												},
												Incoming: []string{"_cddf9325-a85b-4347-8c57-8b909fa77ae9"},
												Outgoing: []string{"_435d9320-bbf4-48ad-aa56-16cb5483e95b"},
											},
										},
									},
									Implementation: "WebService",
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_6a34496f-8cf7-42e5-88a9-d1af98cc3cba",
														Documentation: []element.Documentation{
															{
																Id:    "Bpmn_Documentation_jcoo8e_bEeSGoscwBjzAjw",
																Value: "These customers get special service in case of damage.\nThis service level guarantees less than two days to solve the issue.",
															},
														},
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
													},
													Name: "Perform repair (premium level)",
												},
												Incoming: []string{"_be893987-caec-4605-b078-bd96b7cd6c12", "Bpmn_SequenceFlow_-dQQwBqHEeWDuOtG0oS24A"},
												Outgoing: []string{"_d54c5707-af7a-4b36-a059-46681bbf0004"},
											},
										},
									},
									Implementation: "WebService",
								},
							},
							EndEvents: []element.EndEvent{
								{
									ThrowEvent: element.ThrowEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID:                "_177bd313-c6c9-4df5-8f82-313beb30d2eb",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
													},
													Name: "Fridge replaced",
												},
												Incoming: []string{"_c5756bb9-6e6f-42d1-8799-c2d673499eb8"},
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
														ID:                "_b3dc1906-d4d3-40c5-aaf6-5a74148ae887",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
													},
													Name: "Emergency repair completed",
												},
												Incoming: []string{"_cf380e47-1401-4e7e-b710-193b626e49eb"},
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
														ID:                "_dcee5c64-3010-4ee5-b480-bce856e6f29c",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
													},
													Name: "Repair completed",
												},
												Incoming: []string{"_435d9320-bbf4-48ad-aa56-16cb5483e95b", "_d54c5707-af7a-4b36-a059-46681bbf0004"},
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
													ID:                "_604be023-654c-44df-a64c-365254a100cd",
													ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
												},
												Name: "Service type",
											},
											Incoming: []string{"_acb2aca3-8851-48f0-b127-7b3c9db5e18d"},
											Outgoing: []string{"_b99800c3-c340-460c-a43e-098014a365d0", "_437e5969-1e61-4cb9-aa76-4f8854f32eeb", "_ada039b6-94dd-4a15-a6b1-c7fe662c64ee"},
										},
										GatewayDirection: "Diverging",
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID:                "_2b6bc88e-d3be-4704-87d1-c264bf704589",
													ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
												},
												Name: "Service level",
											},
											Incoming: []string{"_ada039b6-94dd-4a15-a6b1-c7fe662c64ee"},
											Outgoing: []string{"_cddf9325-a85b-4347-8c57-8b909fa77ae9", "_be893987-caec-4605-b078-bd96b7cd6c12"},
										},
										GatewayDirection: "Diverging",
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID:                "_936c0bfa-5ebf-4546-8d1a-cce556148788",
													ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
												},
												Name: "Successful?",
											},
											Incoming: []string{"_b22c01a4-1eef-4f52-9b16-a201a9621619"},
											Outgoing: []string{"_cf380e47-1401-4e7e-b710-193b626e49eb", "_3fb323d5-2c59-487a-af63-804208f6c5cb"},
										},
										GatewayDirection: "Diverging",
									},
								},
							},
							SubProcesses: []element.SubProcess{
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_cd6f230f-13c3-4027-aa3e-57de601a1ab2",
													Documentation: []element.Documentation{
														{
															Id:    "Bpmn_Documentation_jcpQAO_bEeSGoscwBjzAjw",
															Value: "This is a special case for handling emergencies for clients such as Hospitals or Commercial fridges.",
														},
													},
													ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
												},
												Name: "Perform emergency repair",
											},
											Incoming: []string{"_437e5969-1e61-4cb9-aa76-4f8854f32eeb", "Bpmn_SequenceFlow_tcbpgBqGEeWDuOtG0oS24A"},
											Outgoing: []string{"_b22c01a4-1eef-4f52-9b16-a201a9621619"},
										},
									},
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_acb2aca3-8851-48f0-b127-7b3c9db5e18d",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
										},
									},
									SourceRef: "_c73a5f4a-72f1-4e11-bb40-2f98da75fb9a",
									TargetRef: "_604be023-654c-44df-a64c-365254a100cd",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_b99800c3-c340-460c-a43e-098014a365d0",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
										},
										Name: "Warranty",
									},
									SourceRef: "_604be023-654c-44df-a64c-365254a100cd",
									TargetRef: "_a92069f7-377b-4dbd-a1fd-1da071aabf6d",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_437e5969-1e61-4cb9-aa76-4f8854f32eeb",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
										},
										Name: "Emergency service",
									},
									SourceRef: "_604be023-654c-44df-a64c-365254a100cd",
									TargetRef: "_cd6f230f-13c3-4027-aa3e-57de601a1ab2",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_ada039b6-94dd-4a15-a6b1-c7fe662c64ee",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
										},
										Name: "Regular repair service",
									},
									SourceRef: "_604be023-654c-44df-a64c-365254a100cd",
									TargetRef: "_2b6bc88e-d3be-4704-87d1-c264bf704589",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_cddf9325-a85b-4347-8c57-8b909fa77ae9",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
										},
										Name: "Standard",
									},
									SourceRef: "_2b6bc88e-d3be-4704-87d1-c264bf704589",
									TargetRef: "_d034722f-751d-4f37-a3d7-47993822e979",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_be893987-caec-4605-b078-bd96b7cd6c12",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
										},
										Name: "Premium",
									},
									SourceRef: "_2b6bc88e-d3be-4704-87d1-c264bf704589",
									TargetRef: "_6a34496f-8cf7-42e5-88a9-d1af98cc3cba",
									ConditionExpression: element.ExpressionUnMarshal{
										Type: "tFormalExpression",
										ExpressionSubstitution: &element.FormalExpression{
											Expression: element.Expression{
												BaseElementWithMixedContent: element.BaseElementWithMixedContent{
													ID: "Bpmn_FormalExpression_ipW2kBqIEeWDuOtG0oS24A",
												},
											},
											Value: "Service Level == 'Premium'",
										},
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_435d9320-bbf4-48ad-aa56-16cb5483e95b",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
										},
									},
									SourceRef: "_d034722f-751d-4f37-a3d7-47993822e979",
									TargetRef: "_dcee5c64-3010-4ee5-b480-bce856e6f29c",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_d54c5707-af7a-4b36-a059-46681bbf0004",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
										},
									},
									SourceRef: "_6a34496f-8cf7-42e5-88a9-d1af98cc3cba",
									TargetRef: "_dcee5c64-3010-4ee5-b480-bce856e6f29c",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_b22c01a4-1eef-4f52-9b16-a201a9621619",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
										},
									},
									SourceRef: "_cd6f230f-13c3-4027-aa3e-57de601a1ab2",
									TargetRef: "_936c0bfa-5ebf-4546-8d1a-cce556148788",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_c5756bb9-6e6f-42d1-8799-c2d673499eb8",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
										},
									},
									SourceRef: "_a92069f7-377b-4dbd-a1fd-1da071aabf6d",
									TargetRef: "_177bd313-c6c9-4df5-8f82-313beb30d2eb",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_cf380e47-1401-4e7e-b710-193b626e49eb",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
										},
										Name: "yes",
									},
									SourceRef: "_936c0bfa-5ebf-4546-8d1a-cce556148788",
									TargetRef: "_b3dc1906-d4d3-40c5-aaf6-5a74148ae887",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_3fb323d5-2c59-487a-af63-804208f6c5cb",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
										},
										Name: "no",
									},
									SourceRef: "_936c0bfa-5ebf-4546-8d1a-cce556148788",
									TargetRef: "_a92069f7-377b-4dbd-a1fd-1da071aabf6d",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_a5af06ae-bd69-464d-bbaf-3d7418702d77",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
										},
									},
									SourceRef: "_cc9778bd-edd8-4df2-ba15-56c310f90e62",
									TargetRef: "_c73a5f4a-72f1-4e11-bb40-2f98da75fb9a",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "Bpmn_SequenceFlow_tcbpgBqGEeWDuOtG0oS24A",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
										},
									},
									IsImmediate: true,
									SourceRef:   "Bpmn_BoundaryEvent_sS9gABqGEeWDuOtG0oS24A",
									TargetRef:   "_cd6f230f-13c3-4027-aa3e-57de601a1ab2",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "Bpmn_SequenceFlow_-dQQwBqHEeWDuOtG0oS24A",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
										},
									},
									IsImmediate: true,
									SourceRef:   "Bpmn_BoundaryEvent_LwKtwhqHEeWDuOtG0oS24A",
									TargetRef:   "_6a34496f-8cf7-42e5-88a9-d1af98cc3cba",
								},
							},
							BoundaryEvents: []element.BoundaryEvent{
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID:                "Bpmn_BoundaryEvent_sS9gABqGEeWDuOtG0oS24A",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
													},
													Name: "2 hours",
												},
												Outgoing: []string{"Bpmn_SequenceFlow_tcbpgBqGEeWDuOtG0oS24A"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											TimerEventDefinitions: []element.TimerEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "Bpmn_TimerEventDefinition_88uHkBqGEeWDuOtG0oS24A",
															},
														},
													},
													TimeDuration: element.ExpressionUnMarshal{
														Type: element.ExpressionTypeFormal,
														ExpressionSubstitution: &element.FormalExpression{
															Expression: element.Expression{
																BaseElementWithMixedContent: element.BaseElementWithMixedContent{
																	ID: "Bpmn_FormalExpression_8pPK4BqGEeWDuOtG0oS24A",
																},
															},
															Language: "http://www.w3.org/1999/XPath",
															Value:    "PT2H",
														},
													},
												},
											},
										},
									},
									CancelActivity: true,
									AttachedToRef: "_6a34496f-8cf7-42e5-88a9-d1af98cc3cba",
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID:                "Bpmn_BoundaryEvent_LwKtwhqHEeWDuOtG0oS24A",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
													},
												},
												Outgoing: []string{"Bpmn_SequenceFlow_-dQQwBqHEeWDuOtG0oS24A"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "Bpmn_MessageEventDefinition_zhREUBqHEeWDuOtG0oS24A",
															},
														},
													},
													MessageRef: "_6840c92f-f02c-46b3-9565-f5e5b1792ba0",
												},
											},
										},
									},
									CancelActivity: true,
									AttachedToRef: "_d034722f-751d-4f37-a3d7-47993822e979",
								},
							},
						},
					},
				},
				Messages: []element.Message{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID:                "_6840c92f-f02c-46b3-9565-f5e5b1792ba0",
								ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
							},
						},
						Name: "Service Level",
					},
				},
				Resources: []element.Resource{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "Bpmn_Resource__7wrkBqGEeWDuOtG0oS24A",
							},
						},
						Name: "User",
					},
				},
			},
		},
	}

	if diff := cmp.Diff(excepted, modelelement); diff != "" {
		t.Errorf("Mismatch in %s (-want +got):\n%s", path, diff)
	}
}
