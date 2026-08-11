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

func TestA_2_0_rountrip(t *testing.T) {
	// create test use ./test/A.2.0-rountrip.bpmn
	path := "./test/A.2.0-rountrip.bpmn"
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
			TargetNamespace: "http://www.trisotech.com/definitions/_1373649889746",
			ID:              "_1373649889746",
			Name:            "A.2.0",
			RootElemnts: element.RootElemnts{
				Processes: []element.Process{
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "WFP-6-",
								},
							},
						},
						IsExecutable: false,
						FlowElements: element.FlowElements{
							StartEvents: []element.StartEvent{{
								CatchEvent: element.CatchEvent{
									Event: element.Event{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_6b5db6a9-037a-49ad-9201-09201e2aaa97",
												},
												Name: "Start Event",
											},
											Outgoing: []string{"_b50f530c-3450-4e1a-b81f-ea346dc6e1cb"},
										},
									},
								},
							}},
							Tasks: []element.Task{
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_5a972b87-735d-454a-b31c-f52fb3afc5c7",
												},
												Name: "Task 1",
											},
											Incoming: []string{"_b50f530c-3450-4e1a-b81f-ea346dc6e1cb"},
											Outgoing: []string{"_fe74c141-8843-4b00-a704-5e5e13be53b0"},
										},
										StartQuantity:      1,
										CompletionQuantity: 1,
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_4f7d62d7-f0e6-46bc-be00-69e02da38f65",
												},
												Name: "Task 2",
											},
											Incoming: []string{"_f1478fb7-98c4-4c01-8c15-68bd04c91535"},
											Outgoing: []string{"_a3d40a56-9b7f-417e-911e-d39e7f18b90c"},
										},
										StartQuantity:      1,
										CompletionQuantity: 1,
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_e6eb725a-34bc-45c7-aed0-9f9596cd7bee",
												},
												Name: "Task 3",
											},
											Incoming: []string{"_a1570a53-28d2-41b1-a3a2-3e50c00d747e"},
											Outgoing: []string{"_e9ebc7c7-995d-46db-86ce-d823bc2b4687"},
										},
										StartQuantity:      1,
										CompletionQuantity: 1,
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_7d399717-1aba-47ac-8d7d-8aaa033255e0",
												},
												Name: "Task 4",
											},
											Incoming: []string{"_20ebb3c1-5178-4c7c-a91d-23e58f2aa73b"},
											Outgoing: []string{"_698b593f-18eb-42ea-b8cd-bcd51e1514cc"},
										},
										StartQuantity:      1,
										CompletionQuantity: 1,
									},
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_b50f530c-3450-4e1a-b81f-ea346dc6e1cb",
										},
									},
									SourceRef: "_6b5db6a9-037a-49ad-9201-09201e2aaa97",
									TargetRef: "_5a972b87-735d-454a-b31c-f52fb3afc5c7",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_fe74c141-8843-4b00-a704-5e5e13be53b0",
										},
									},
									SourceRef: "_5a972b87-735d-454a-b31c-f52fb3afc5c7",
									TargetRef: "_35fe57a7-1302-44e2-bf58-032f11af7ecb",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_f1478fb7-98c4-4c01-8c15-68bd04c91535",
										},
									},
									SourceRef: "_35fe57a7-1302-44e2-bf58-032f11af7ecb",
									TargetRef: "_4f7d62d7-f0e6-46bc-be00-69e02da38f65",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_a3d40a56-9b7f-417e-911e-d39e7f18b90c",
										},
									},
									SourceRef: "_4f7d62d7-f0e6-46bc-be00-69e02da38f65",
									TargetRef: "_258f51eb-b764-4a71-b681-3a01cca14143",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_e9ebc7c7-995d-46db-86ce-d823bc2b4687",
										},
									},
									SourceRef: "_e6eb725a-34bc-45c7-aed0-9f9596cd7bee",
									TargetRef: "_33c66216-391c-49c2-aa19-d8f0b7f5f91d",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_698b593f-18eb-42ea-b8cd-bcd51e1514cc",
										},
									},
									SourceRef: "_7d399717-1aba-47ac-8d7d-8aaa033255e0",
									TargetRef: "_33c66216-391c-49c2-aa19-d8f0b7f5f91d",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_d4ce87c6-1373-45d6-a3b4-fbb2a04ee2e5",
										},
									},
									SourceRef: "_33c66216-391c-49c2-aa19-d8f0b7f5f91d",
									TargetRef: "_258f51eb-b764-4a71-b681-3a01cca14143",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_a1570a53-28d2-41b1-a3a2-3e50c00d747e",
										},
									},
									SourceRef: "_35fe57a7-1302-44e2-bf58-032f11af7ecb",
									TargetRef: "_e6eb725a-34bc-45c7-aed0-9f9596cd7bee",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_20ebb3c1-5178-4c7c-a91d-23e58f2aa73b",
										},
									},
									SourceRef: "_35fe57a7-1302-44e2-bf58-032f11af7ecb",
									TargetRef: "_7d399717-1aba-47ac-8d7d-8aaa033255e0",
								},
							},
							ExclusiveGatewaies: []element.ExclusiveGateway{
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_35fe57a7-1302-44e2-bf58-032f11af7ecb",
												},
												Name: "Gateway\n(Split Flow)",
											},
											Incoming: []string{"_fe74c141-8843-4b00-a704-5e5e13be53b0"},
											Outgoing: []string{"_f1478fb7-98c4-4c01-8c15-68bd04c91535", "_a1570a53-28d2-41b1-a3a2-3e50c00d747e", "_20ebb3c1-5178-4c7c-a91d-23e58f2aa73b"},
										},
										GatewayDirection: element.GatewayDirectionUnspecified,
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_33c66216-391c-49c2-aa19-d8f0b7f5f91d",
												},
												Name: "Gateway\n(Merge Flows)",
											},
											Incoming: []string{"_e9ebc7c7-995d-46db-86ce-d823bc2b4687", "_698b593f-18eb-42ea-b8cd-bcd51e1514cc"},
											Outgoing: []string{"_d4ce87c6-1373-45d6-a3b4-fbb2a04ee2e5"},
										},
										GatewayDirection: element.GatewayDirectionUnspecified,
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
														ID: "_258f51eb-b764-4a71-b681-3a01cca14143",
													},
													Name: "End Event",
												},
												Incoming: []string{"_a3d40a56-9b7f-417e-911e-d39e7f18b90c", "_d4ce87c6-1373-45d6-a3b4-fbb2a04ee2e5"},
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
