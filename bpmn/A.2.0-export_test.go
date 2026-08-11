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

func TestA_2_0_export(t *testing.T) {
	// create test use ./test/A.2.0-export.bpmn
	path := "./test/A.2.0-export.bpmn"
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
			TargetNamespace: "http://bpmn.io/bpmn",
			ID:              "sid-38422fae-e03e-43a3-bef4-bd33b32041b2",
			Exporter:        "Camunda Modeler",
			ExporterVersion: "5.23.0",
			RootElemnts: element.RootElemnts{
				Processes: []element.Process{
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "Process_1",
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
													ID: "Event_072o7cv",
												},
												Name: "Start Event",
											},
											Outgoing: []string{"Flow_12fenxe"},
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
													ID: "Activity_0opq70y",
												},
												Name: "Task 1",
											},
											Incoming: []string{"Flow_12fenxe"},
											Outgoing: []string{"Flow_1yupnqc"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_1ljp29t",
												},
												Name: "Task 2",
											},
											Incoming: []string{"Flow_0dd1rck"},
											Outgoing: []string{"Flow_0y0m2k8"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_0jhawx0",
												},
												Name: "Task 3",
											},
											Incoming: []string{"Flow_0x796n6"},
											Outgoing: []string{"Flow_1lk8qao"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Activity_0ddly78",
												},
												Name: "Task 4",
											},
											Incoming: []string{"Flow_1801a2c"},
											Outgoing: []string{"Flow_17lrcjr"},
										},
									},
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_12fenxe",
										},
									},
									SourceRef: "Event_072o7cv",
									TargetRef: "Activity_0opq70y",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1yupnqc",
										},
									},
									SourceRef: "Activity_0opq70y",
									TargetRef: "Gateway_03s9abx",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0dd1rck",
										},
									},
									SourceRef: "Gateway_03s9abx",
									TargetRef: "Activity_1ljp29t",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0x796n6",
										},
									},
									SourceRef: "Gateway_03s9abx",
									TargetRef: "Activity_0jhawx0",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1801a2c",
										},
									},
									SourceRef: "Gateway_03s9abx",
									TargetRef: "Activity_0ddly78",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0y0m2k8",
										},
									},
									SourceRef: "Activity_1ljp29t",
									TargetRef: "Event_1d5wxn1",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1lk8qao",
										},
									},
									SourceRef: "Activity_0jhawx0",
									TargetRef: "Gateway_03haizn",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_17lrcjr",
										},
									},
									SourceRef: "Activity_0ddly78",
									TargetRef: "Gateway_03haizn",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_05zv16r",
										},
									},
									SourceRef: "Gateway_03haizn",
									TargetRef: "Event_1d5wxn1",
								},
							},
							ExclusiveGatewaies: []element.ExclusiveGateway{
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Gateway_03s9abx",
												},
												Name: "Gateway (Split Flow)",
											},
											Incoming: []string{"Flow_1yupnqc"},
											Outgoing: []string{"Flow_0dd1rck", "Flow_0x796n6", "Flow_1801a2c"},
										},
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Gateway_03haizn",
												},
												Name: "Gateway (Merge Flows)",
											},
											Incoming: []string{"Flow_1lk8qao", "Flow_17lrcjr"},
											Outgoing: []string{"Flow_05zv16r"},
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
														ID: "Event_1d5wxn1",
													},
													Name: "End Event",
												},
												Incoming: []string{"Flow_0y0m2k8", "Flow_05zv16r"},
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
