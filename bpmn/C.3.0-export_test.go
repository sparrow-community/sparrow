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

func TestC_3_0_export(t *testing.T) {
	// create test use ./test/C.3.0-export.bpmn
	path := "./test/C.3.0-export.bpmn"
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
			Exporter:        "Camunda Modeler",
			ExporterVersion: "5.23.0",
			RootElemnts: element.RootElemnts{
				Processes: []element.Process{
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "_8170787a-3207-434d-9bea-4787059f444f",
								},
							},
							Name: "Fridge Repair Process",
						},
						ProcessType:  "Private",
						IsExecutable: true,
						FlowElements: element.FlowElements{
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_027gqe8",
										},
									},
									SourceRef: "Event_0issfmv",
									TargetRef: "Activity_175emni",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0wv1yt9",
										},
									},
									SourceRef: "Activity_175emni",
									TargetRef: "Gateway_0mgekl4",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0am3e0w",
										},
										Name: "Warranty",
									},
									SourceRef: "Gateway_0mgekl4",
									TargetRef: "Activity_0vku94f",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0wow8xd",
										},
										Name: "Emergency service",
									},
									SourceRef: "Gateway_0mgekl4",
									TargetRef: "Activity_14jt63w",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1fkchcd",
										},
									},
									SourceRef: "Activity_14jt63w",
									TargetRef: "Gateway_0pp15o5",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1m9fllr",
										},
										Name: "no",
									},
									SourceRef: "Gateway_0pp15o5",
									TargetRef: "Activity_0vku94f",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1p0y3gu",
										},
									},
									SourceRef: "Activity_0vku94f",
									TargetRef: "Event_1279qzr",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0pr12q5",
										},
										Name: "yes",
									},
									SourceRef: "Gateway_0pp15o5",
									TargetRef: "Event_12jgnvi",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_042vwgm",
										},
										Name: "Regular repair service",
									},
									SourceRef: "Gateway_0mgekl4",
									TargetRef: "Gateway_1mknm4o",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_07lhkb0",
										},
										Name: "Premium",
									},
									SourceRef: "Gateway_1mknm4o",
									TargetRef: "Activity_0v5jdpu",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1u59cr8",
										},
									},
									SourceRef: "Activity_0v5jdpu",
									TargetRef: "Event_1ji69gn",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_15h599x",
										},
										Name: "Standard",
									},
									SourceRef: "Gateway_1mknm4o",
									TargetRef: "Activity_0ymu39d",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0iypits",
										},
									},
									SourceRef: "Activity_0ymu39d",
									TargetRef: "Event_1ji69gn",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0yrmaoo",
										},
									},
									SourceRef: "Event_1p3ah5d",
									TargetRef: "Activity_14jt63w",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_02hek2e",
										},
									},
									SourceRef: "Event_0ad1gmz",
									TargetRef: "Activity_0v5jdpu",
								},
							},
							UserTasks: []element.UserTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_175emni",
													},
													Name: "Analyse customer request",
												},
												Incoming: []string{"Flow_027gqe8"},
												Outgoing: []string{"Flow_0wv1yt9"},
											},
										},
									},
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_0vku94f",
													},
													Name: "Replace fridge",
												},
												Incoming: []string{"Flow_0am3e0w", "Flow_1m9fllr"},
												Outgoing: []string{"Flow_1p0y3gu"},
											},
										},
									},
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_0v5jdpu",
													},
													Name: "Perform repair (premium level)",
												},
												Incoming: []string{"Flow_07lhkb0", "Flow_02hek2e"},
												Outgoing: []string{"Flow_1u59cr8"},
											},
										},
									},
								},
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_0ymu39d",
													},
													Name: "Perform repair (standard level)",
												},
												Incoming: []string{"Flow_15h599x"},
												Outgoing: []string{"Flow_0iypits"},
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
													ID: "Gateway_0mgekl4",
												},
												Name: "Service type",
											},
											Incoming: []string{"Flow_0wv1yt9"},
											Outgoing: []string{"Flow_0am3e0w", "Flow_0wow8xd", "Flow_042vwgm"},
										},
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Gateway_0pp15o5",
												},
												Name: "Successful?",
											},
											Incoming: []string{"Flow_1fkchcd"},
											Outgoing: []string{"Flow_1m9fllr", "Flow_0pr12q5"},
										},
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Gateway_1mknm4o",
												},
												Name: "Service level",
											},
											Incoming: []string{"Flow_042vwgm"},
											Outgoing: []string{"Flow_07lhkb0", "Flow_15h599x"},
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
													ID: "Activity_14jt63w",
												},
												Name: "Perform emergency repair",
											},
											Incoming: []string{"Flow_0wow8xd", "Flow_0yrmaoo"},
											Outgoing: []string{"Flow_1fkchcd"},
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
														ID: "Event_1279qzr",
													},
													Name: "Fridge replaced",
												},
												Incoming: []string{"Flow_1p0y3gu"},
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
														ID: "Event_12jgnvi",
													},
													Name: "Emergency repair completed",
												},
												Incoming: []string{"Flow_0pr12q5"},
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
														ID: "Event_1ji69gn",
													},
													Name: "Repair completed",
												},
												Incoming: []string{"Flow_1u59cr8", "Flow_0iypits"},
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
														ID: "Event_1p3ah5d",
													},
													Name: "2 hours",
												},
												Outgoing: []string{"Flow_0yrmaoo"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											TimerEventDefinitions: []element.TimerEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "TimerEventDefinition_088mpk0",
															},
														},
													},
												},
											},
										},
									},
									AttachedToRef: "Activity_0v5jdpu",
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_0ad1gmz",
													},
												},
												Outgoing: []string{"Flow_02hek2e"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "MessageEventDefinition_1th22sv",
															},
														},
													},
												},
											},
										},
									},
									AttachedToRef: "Activity_0ymu39d",
								},
							},
							StartEvents: []element.StartEvent{
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_0issfmv",
													},
													Name: "Receive customer request",
												},
												Outgoing: []string{"Flow_027gqe8"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "MessageEventDefinition_12mw7e8",
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
					},
				},
				Messages: []element.Message{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "_6840c92f-f02c-46b3-9565-f5e5b1792ba0",
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
