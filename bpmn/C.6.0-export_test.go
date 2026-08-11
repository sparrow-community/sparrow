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

func TestC_6_0_export(t *testing.T) {
	// create test use ./test/C.6.0-export.bpmn
	path := "./test/C.6.0-export.bpmn"
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
			ID:              "Definitions_13oah7q",
			TargetNamespace: "http://bpmn.io/schema/bpmn",
			Exporter:        "Camunda Modeler",
			ExporterVersion: "5.23.0",
			RootElemnts: element.RootElemnts{
				Processes: []element.Process{
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "Process_19noqni",
								},
							},
						},
						IsExecutable: true,
						FlowElements: element.FlowElements{
							StartEvents: []element.StartEvent{
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "StartEvent_1",
													},
													Name: "Receive Travel Request",
												},
												Outgoing: []string{"Flow_1xi649k"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "MessageEventDefinition_0ein9t6",
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
											ID: "Flow_1xi649k",
										},
									},
									SourceRef: "StartEvent_1",
									TargetRef: "Activity_1qdxrgj",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1p58yis",
										},
									},
									SourceRef: "Activity_1qdxrgj",
									TargetRef: "Gateway_1ersh6n",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_128dw9w",
										},
									},
									SourceRef: "Gateway_1ersh6n",
									TargetRef: "Event_0w821nf",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_15l6xor",
										},
									},
									SourceRef: "Gateway_1ersh6n",
									TargetRef: "Event_1gu9t77",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1a3h5v6",
										},
									},
									SourceRef: "Gateway_1ersh6n",
									TargetRef: "Event_19meht8",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1uw2kdx",
										},
									},
									SourceRef: "Event_0w821nf",
									TargetRef: "Activity_1bidfcm",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1dytiuq",
										},
									},
									SourceRef: "Event_1gu9t77",
									TargetRef: "Activity_1g04x9h",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1domhhx",
										},
										Name: "24 Hours",
									},
									SourceRef: "Event_0isfp1w",
									TargetRef: "Activity_1g04x9h",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0gx0bqr",
										},
									},
									SourceRef: "Event_19meht8",
									TargetRef: "Activity_084p2mw",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0yqzl55",
										},
									},
									SourceRef: "Activity_1g04x9h",
									TargetRef: "Event_0lsnfz2",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1aavu4q",
										},
									},
									SourceRef: "Activity_084p2mw",
									TargetRef: "Event_1yyb73n",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_075xnug",
										},
									},
									SourceRef: "Activity_1bidfcm",
									TargetRef: "Activity_0p5xveb",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0ufn18k",
										},
									},
									SourceRef: "Event_0gv16hd",
									TargetRef: "Activity_05lmtrf",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0vi45od",
										},
									},
									SourceRef: "Activity_0p5xveb",
									TargetRef: "Activity_039ic8d",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1p554j5",
										},
									},
									SourceRef: "Event_0wnb2z5",
									TargetRef: "Event_0qemotd",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1y9j47t",
										},
									},
									SourceRef: "Event_0qemotd",
									TargetRef: "Activity_13l7j32",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0ld9uyo",
										},
									},
									SourceRef: "Activity_039ic8d",
									TargetRef: "Activity_1a4r8ya",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1fz8fw8",
										},
									},
									SourceRef: "Activity_13l7j32",
									TargetRef: "Event_0aabyn5",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1wfqz94",
										},
									},
									SourceRef: "Activity_1a4r8ya",
									TargetRef: "Event_0l6vh4o",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0d43f6n",
										},
									},
									SourceRef: "Activity_05lmtrf",
									TargetRef: "Event_0rk6pwf",
								},
							},
							SendTasks: []element.SendTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_1qdxrgj",
													},
													Name: "Make Flights and Hotel Offer",
												},
												Incoming: []string{"Flow_1xi649k"},
												Outgoing: []string{"Flow_1p58yis"},
											},
											Properties: []element.Property{
												{
													BaseElement: element.BaseElement{
														ID: "Property_0jmwx6a",
													},
													Name: "__targetRef_placeholder",
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "DataInputAssociation_17xmvz3",
														},
														SourceRef: "DataObjectReference_0vc6qe4",
														TargetRef: "Property_0jmwx6a",
													},
												},
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
														ID: "Activity_1bidfcm",
													},
													Name: "Request Credit Card Information",
												},
												Incoming: []string{"Flow_1uw2kdx"},
												Outgoing: []string{"Flow_075xnug"},
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
														ID: "Activity_1g04x9h",
													},
													Name: "Notify Customer Offer Expired",
												},
												Incoming: []string{"Flow_1dytiuq", "Flow_1domhhx"},
												Outgoing: []string{"Flow_0yqzl55"},
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
														ID: "Activity_05lmtrf",
													},
													Name: "Notify Failed Booking",
												},
												Incoming: []string{"Flow_0ufn18k"},
												Outgoing: []string{"Flow_0d43f6n"},
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
														ID: "Activity_13l7j32",
													},
													Name: "Notify Failed Credit Transaction",
												},
												Incoming: []string{"Flow_1y9j47t"},
												Outgoing: []string{"Flow_1fz8fw8"},
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
														ID: "Activity_1a4r8ya",
													},
													Name: "Confirm Booking",
												},
												Incoming: []string{"Flow_0ld9uyo"},
												Outgoing: []string{"Flow_1wfqz94"},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "DataOutputAssociation_0yhidap",
														},
														TargetRef: "DataObjectReference_0hy72ip",
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
											ID: "DataObjectReference_0vc6qe4",
										},
										Name: "Travel Request",
									},
									DataObjectRef: "DataObject_0acgqpe",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObjectReference_0hy72ip",
										},
										Name: "Itinerary",
									},
									DataObjectRef: "DataObject_1h6jn3y",
								},
							},
							DataObjects: []element.DataObject{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObject_0acgqpe",
										},
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObject_1h6jn3y",
										},
									},
								},
							},
							EventBasedGatewaies: []element.EventBasedGateway{
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Gateway_1ersh6n",
												},
											},
											Incoming: []string{"Flow_1p58yis"},
											Outgoing: []string{"Flow_128dw9w", "Flow_15l6xor", "Flow_1a3h5v6"},
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
														ID: "Event_0w821nf",
													},
													Name: "Offer Approved",
												},
												Incoming: []string{"Flow_128dw9w"},
												Outgoing: []string{"Flow_1uw2kdx"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "MessageEventDefinition_12s4ok2",
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
														ID: "Event_1gu9t77",
													},
													Name: "24 Hours",
												},
												Incoming: []string{"Flow_15l6xor"},
												Outgoing: []string{"Flow_1dytiuq"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											TimerEventDefinitions: []element.TimerEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "TimerEventDefinition_0gimu3s",
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
														ID: "Event_19meht8",
													},
													Name: "Cancel Request",
												},
												Incoming: []string{"Flow_1a3h5v6"},
												Outgoing: []string{"Flow_0gx0bqr"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											MessageEventDefinitions: []element.MessageEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "MessageEventDefinition_0ylxz6n",
															},
														},
													},
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
														ID: "Event_0qemotd",
													},
													Name: "Booking",
												},
												Incoming: []string{"Flow_1p554j5"},
												Outgoing: []string{"Flow_1y9j47t"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											CompensateEventDefinitions: []element.CompensateEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "CompensateEventDefinition_1nra8kx",
															},
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
														ID: "Event_0isfp1w",
													},
												},
												Outgoing: []string{"Flow_1domhhx"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											TimerEventDefinitions: []element.TimerEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "TimerEventDefinition_0ry7jh8",
															},
														},
													},
												},
											},
										},
									},
									AttachedToRef: "Activity_1bidfcm",
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_0gv16hd",
													},
												},
												Outgoing: []string{"Flow_0ufn18k"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											ErrorEventDefinitions: []element.ErrorEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "ErrorEventDefinition_10wgwwl",
															},
														},
													},
												},
											},
										},
									},
									AttachedToRef: "Activity_0p5xveb",
								},
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_0wnb2z5",
													},
												},
												Outgoing: []string{"Flow_1p554j5"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											ErrorEventDefinitions: []element.ErrorEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "ErrorEventDefinition_1tcdwud",
															},
														},
													},
												},
											},
										},
									},
									AttachedToRef: "Activity_039ic8d",
								},
							},
							ServiceTasks: []element.ServiceTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_084p2mw",
													},
													Name: "Update Customer Record",
												},
												Incoming: []string{"Flow_0gx0bqr"},
												Outgoing: []string{"Flow_1aavu4q"},
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
														ID: "Activity_039ic8d",
													},
													Name: "Charge Credit Card",
												},
												Incoming: []string{"Flow_0vi45od"},
												Outgoing: []string{"Flow_0ld9uyo"},
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
														ID: "Event_0lsnfz2",
													},
													Name: "Offer Expired",
												},
												Incoming: []string{"Flow_0yqzl55"},
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
														ID: "Event_1yyb73n",
													},
													Name: "Request Cancelled",
												},
												Incoming: []string{"Flow_1aavu4q"},
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
														ID: "Event_0aabyn5",
													},
													Name: "Failed Credit Transaction",
												},
												Incoming: []string{"Flow_1fz8fw8"},
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
														ID: "Event_0l6vh4o",
													},
													Name: "Booking Confirmed",
												},
												Incoming: []string{"Flow_1wfqz94"},
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
														ID: "Event_0rk6pwf",
													},
													Name: "Failed Booking",
												},
												Incoming: []string{"Flow_0d43f6n"},
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
													ID: "Activity_0p5xveb",
												},
												Name: "Make Booking",
											},
											Incoming: []string{"Flow_075xnug"},
											Outgoing: []string{"Flow_0vi45od"},
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
																	ID: "Event_1fywxat",
																},
															},
															Outgoing: []string{"Flow_0wb651w"},
														},
													},
												},
											},
										},
										SequenceFlows: []element.SequenceFlow{
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_0wb651w",
													},
												},
												SourceRef: "Event_1fywxat",
												TargetRef: "Gateway_14mg5pf",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_01h7e6d",
													},
												},
												SourceRef: "Gateway_14mg5pf",
												TargetRef: "Activity_13sg203",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_04o2v8w",
													},
												},
												SourceRef: "Activity_13sg203",
												TargetRef: "Gateway_1mo08sa",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_1h0awz9",
													},
												},
												SourceRef: "Gateway_14mg5pf",
												TargetRef: "Activity_0qz49yv",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_1l3h2oz",
													},
												},
												SourceRef: "Activity_0qz49yv",
												TargetRef: "Gateway_1mo08sa",
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Flow_0stizar",
													},
												},
												SourceRef: "Gateway_1mo08sa",
												TargetRef: "Event_1c52ias",
											},
										},
										ParallelGatewaies: []element.ParallelGateway{
											{
												Gateway: element.Gateway{
													FlowNode: element.FlowNode{
														FlowElement: element.FlowElement{
															BaseElement: element.BaseElement{
																ID: "Gateway_14mg5pf",
															},
														},
														Incoming: []string{"Flow_0wb651w"},
														Outgoing: []string{"Flow_01h7e6d", "Flow_1h0awz9"},
													},
												},
											},
											{
												Gateway: element.Gateway{
													FlowNode: element.FlowNode{
														FlowElement: element.FlowElement{
															BaseElement: element.BaseElement{
																ID: "Gateway_1mo08sa",
															},
														},
														Incoming: []string{"Flow_04o2v8w", "Flow_1l3h2oz"},
														Outgoing: []string{"Flow_0stizar"},
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
																	ID: "Activity_13sg203",
																},
																Name: "Book Flight",
															},
															Incoming: []string{"Flow_01h7e6d"},
															Outgoing: []string{"Flow_04o2v8w"},
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
																	ID: "Activity_0hgj2bs",
																},
																Name: "Cancel Flight",
															},
														},
														IsForCompensation: true,
													},
												},
											},
											{
												Task: element.Task{
													Activity: element.Activity{
														FlowNode: element.FlowNode{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "Activity_0qz49yv",
																},
																Name: "Book Hotel",
															},
															Incoming: []string{"Flow_1h0awz9"},
															Outgoing: []string{"Flow_1l3h2oz"},
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
																	ID: "Activity_1n0lwxw",
																},
																Name: "Cancel Hotel",
															},
														},
														IsForCompensation: true,
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
																	ID: "Event_0oxjqip",
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
																		BaseElement: element.BaseElement{
																			ID: "CompensateEventDefinition_0l7f698",
																		},
																	},
																},
															},
														},
													},
												},
												AttachedToRef: "Activity_13sg203",
											},
											{
												CatchEvent: element.CatchEvent{
													Event: element.Event{
														FlowNode: element.FlowNode{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "Event_1gsyz0h",
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
																		BaseElement: element.BaseElement{
																			ID: "CompensateEventDefinition_11wluxa",
																		},
																	},
																},
															},
														},
													},
												},
												AttachedToRef: "Activity_0qz49yv",
											},
										},
										EndEvents: []element.EndEvent{
											{
												ThrowEvent: element.ThrowEvent{
													Event: element.Event{
														FlowNode: element.FlowNode{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "Event_1c52ias",
																},
																Name: "Travel Booked",
															},
															Incoming: []string{"Flow_0stizar"},
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
																ID: "Activity_1t020b1",
															},
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
																				ID: "Event_0hlskm4",
																			},
																		},
																		Outgoing: []string{"Flow_01qcqu0"},
																	},
																},
																EventDefinitions: element.EventDefinitions{
																	CompensateEventDefinitions: []element.CompensateEventDefinition{
																		{
																			EventDefinition: element.EventDefinition{
																				RootElement: element.RootElement{
																					BaseElement: element.BaseElement{
																						ID: "CompensateEventDefinition_11z58bi",
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
																	ID: "Flow_01qcqu0",
																},
															},
															SourceRef: "Event_0hlskm4",
															TargetRef: "Gateway_0azwu61",
														},
														{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "Flow_0ck9rq3",
																},
															},
															SourceRef: "Gateway_0azwu61",
															TargetRef: "Event_1o7y58x",
														},
														{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "Flow_168c72o",
																},
															},
															SourceRef: "Gateway_0azwu61",
															TargetRef: "Event_17sn5te",
														},
														{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "Flow_1f1d6rn",
																},
															},
															SourceRef: "Event_17sn5te",
															TargetRef: "Gateway_1udyyri",
														},
														{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "Flow_170itut",
																},
															},
															SourceRef: "Event_1o7y58x",
															TargetRef: "Gateway_1udyyri",
														},
														{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID: "Flow_0hvvmqj",
																},
															},
															SourceRef: "Gateway_1udyyri",
															TargetRef: "Event_0bfvt7c",
														},
													},
													ParallelGatewaies: []element.ParallelGateway{
														{
															Gateway: element.Gateway{
																FlowNode: element.FlowNode{
																	FlowElement: element.FlowElement{
																		BaseElement: element.BaseElement{
																			ID: "Gateway_0azwu61",
																		},
																	},
																	Incoming: []string{"Flow_01qcqu0"},
																	Outgoing: []string{"Flow_0ck9rq3", "Flow_168c72o"},
																},
															},
														},
														{
															Gateway: element.Gateway{
																FlowNode: element.FlowNode{
																	FlowElement: element.FlowElement{
																		BaseElement: element.BaseElement{
																			ID: "Gateway_1udyyri",
																		},
																	},
																	Incoming: []string{"Flow_1f1d6rn", "Flow_170itut"},
																	Outgoing: []string{"Flow_0hvvmqj"},
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
																				ID: "Event_1o7y58x",
																			},
																			Name: "Flight",
																		},
																		Incoming: []string{"Flow_0ck9rq3"},
																		Outgoing: []string{"Flow_170itut"},
																	},
																},
																EventDefinitions: element.EventDefinitions{
																	CompensateEventDefinitions: []element.CompensateEventDefinition{
																		{
																			EventDefinition: element.EventDefinition{
																				RootElement: element.RootElement{
																					BaseElement: element.BaseElement{
																						ID: "CompensateEventDefinition_1wju1u7",
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
																				ID: "Event_17sn5te",
																			},
																			Name: "Hotel",
																		},
																		Incoming: []string{"Flow_168c72o"},
																		Outgoing: []string{"Flow_1f1d6rn"},
																	},
																},
																EventDefinitions: element.EventDefinitions{
																	CompensateEventDefinitions: []element.CompensateEventDefinition{
																		{
																			EventDefinition: element.EventDefinition{
																				RootElement: element.RootElement{
																					BaseElement: element.BaseElement{
																						ID: "CompensateEventDefinition_08lflw7",
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
																				ID: "Event_0bfvt7c",
																			},
																		},
																		Incoming: []string{"Flow_0hvvmqj"},
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
														ID: "Association_03vh68e",
													},
												},
												AssociationDirection: "One",
												SourceRef:            "Event_0oxjqip",
												TargetRef:            "Activity_0hgj2bs",
											},
											{
												Artifact: element.Artifact{
													BaseElement: element.BaseElement{
														ID: "Association_0srffpa",
													},
												},
												AssociationDirection: "One",
												SourceRef:            "Event_1gsyz0h",
												TargetRef:            "Activity_1n0lwxw",
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
