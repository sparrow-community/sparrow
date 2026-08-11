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

func TestC_8_0_export(t *testing.T) {
	// create test use ./test/C.8.0-export.bpmn
	path := "./test/C.8.0-export.bpmn"
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
									ID: "Process_1xl5gyi",
								},
							},
						},
						IsExecutable: false,
						FlowElements: element.FlowElements{
							StartEvents: []element.StartEvent{
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_0q4cbsy",
													},
													Name: "Vacation Request Received",
												},
												Outgoing: []string{"Flow_1916suv"},
											},
										},
									},
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1916suv",
										},
									},
									SourceRef: "Event_0q4cbsy",
									TargetRef: "Activity_0lrykpf",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_05r3ppa",
										},
									},
									SourceRef: "Activity_0lrykpf",
									TargetRef: "Activity_13m56bo",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1229gb2",
										},
									},
									SourceRef: "Activity_13m56bo",
									TargetRef: "Gateway_1m1epyv",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1slje8w",
										},
										Name: "Approved",
									},
									SourceRef: "Gateway_1m1epyv",
									TargetRef: "Activity_0kxqh4y",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_17w7s53",
										},
									},
									SourceRef: "Activity_0kxqh4y",
									TargetRef: "Activity_0h1298a",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_06f2rqz",
										},
										Name: "Manual Validation Required",
									},
									SourceRef: "Gateway_1m1epyv",
									TargetRef: "Activity_14g794r",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1bn59o5",
										},
									},
									SourceRef: "Activity_14g794r",
									TargetRef: "Gateway_0bpdi0b",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_18wfjpz",
										},
										Name: "Approved",
									},
									SourceRef: "Gateway_0bpdi0b",
									TargetRef: "Activity_1exp9va",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1kjzd51",
										},
									},
									SourceRef: "Activity_1exp9va",
									TargetRef: "Activity_0qxgfc9",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1s1ikgt",
										},
									},
									SourceRef: "Activity_0qxgfc9",
									TargetRef: "Event_0fselnt",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0dhwpuk",
										},
									},
									SourceRef: "Activity_0h1298a",
									TargetRef: "Event_0jqpsr9",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_16iauun",
										},
										Name: "Refused",
									},
									SourceRef: "Gateway_0bpdi0b",
									TargetRef: "Activity_02k6h6m",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1uc40yw",
										},
									},
									SourceRef: "Activity_02k6h6m",
									TargetRef: "Event_10bj8rf",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1x3mfjv",
										},
										Name: "Refused",
									},
									SourceRef: "Gateway_1m1epyv",
									TargetRef: "Activity_0wd8gll",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1x4un6h",
										},
									},
									SourceRef: "Activity_0wd8gll",
									TargetRef: "Event_0990itm",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_15v234s",
										},
									},
									SourceRef: "Event_1aoguwi",
									TargetRef: "Event_0c4c8gy",
								},
							},
							ServiceTasks: []element.ServiceTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_0lrykpf",
													},
													Name: "Fetch Vacation Information",
												},
												Incoming: []string{"Flow_1916suv"},
												Outgoing: []string{"Flow_05r3ppa"},
											},
											Properties: []element.Property{
												{
													BaseElement: element.BaseElement{
														ID: "Property_1ffrmni",
													},
													Name: "__targetRef_placeholder",
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "DataInputAssociation_19594jt",
														},
														SourceRef: "DataObjectReference_0s5y8ld",
														TargetRef: "Property_1ffrmni",
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "DataOutputAssociation_1gvweaq",
														},
														TargetRef: "DataObjectReference_0uia8i9",
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
														ID: "Activity_0h1298a",
													},
													Name: "Update Remaining Vacation",
												},
												Incoming: []string{"Flow_17w7s53"},
												Outgoing: []string{"Flow_0dhwpuk"},
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
														ID: "Activity_0qxgfc9",
													},
													Name: "Update Remaining Vacation",
												},
												Incoming: []string{"Flow_1kjzd51"},
												Outgoing: []string{"Flow_1s1ikgt"},
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
														ID: "Event_1aoguwi",
													},
												},
												Outgoing: []string{"Flow_15v234s"},
											},
										},
										EventDefinitions: element.EventDefinitions{
											ErrorEventDefinitions: []element.ErrorEventDefinition{
												{
													EventDefinition: element.EventDefinition{
														RootElement: element.RootElement{
															BaseElement: element.BaseElement{
																ID: "ErrorEventDefinition_1kchpzk",
															},
														},
													},
												},
											},
										},
									},
									AttachedToRef: "Activity_0lrykpf",
								},
							},
							DataObjectReferenes: []element.DataObjectReference{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObjectReference_0s5y8ld",
										},
										Name: "Employee Badge Number",
									},
									DataObjectRef: "DataObject_13kw65u",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObjectReference_0uia8i9",
										},
										Name: "Current Vacation Status",
									},
									DataObjectRef: "DataObject_0ta7zty",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObjectReference_1m7aoov",
										},
										Name: "To",
									},
									DataObjectRef: "DataObject_1grql3s",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObjectReference_06x7a95",
										},
										Name: "From",
									},
									DataObjectRef: "DataObject_1r856yi",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObjectReference_19p7a9l",
										},
										Name: "Vacation Approval",
									},
									DataObjectRef: "DataObject_0msvzum",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObjectReference_0ynczgz",
										},
										Name: "Reason",
									},
									DataObjectRef: "DataObject_0v2smxx",
								},
							},
							DataObjects: []element.DataObject{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObject_13kw65u",
										},
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObject_0ta7zty",
										},
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObject_1grql3s",
										},
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObject_1r856yi",
										},
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObject_0msvzum",
										},
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObject_0v2smxx",
										},
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
														ID: "Activity_13m56bo",
													},
													Name: "Vacation Approval",
												},
												Incoming: []string{"Flow_05r3ppa"},
												Outgoing: []string{"Flow_1229gb2"},
											},
											Properties: []element.Property{
												{
													BaseElement: element.BaseElement{
														ID: "Property_1c2el8a",
													},
													Name: "__targetRef_placeholder",
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "DataInputAssociation_03oginu",
														},
														SourceRef: "DataObjectReference_1m7aoov",
														TargetRef: "Property_1c2el8a",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "DataInputAssociation_1j9yja8",
														},
														SourceRef: "DataObjectReference_06x7a95",
														TargetRef: "Property_1c2el8a",
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "DataOutputAssociation_06bl88p",
														},
														TargetRef: "DataObjectReference_19p7a9l",
													},
												},
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "DataOutputAssociation_0jmelr8",
														},
														TargetRef: "DataObjectReference_0ynczgz",
													},
												},
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
													ID: "Gateway_1m1epyv",
												},
												Name: "",
											},
											Incoming: []string{"Flow_1229gb2"},
											Outgoing: []string{"Flow_1slje8w", "Flow_06f2rqz", "Flow_1x3mfjv"},
										},
									},
									Default: "Flow_1x3mfjv",
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Gateway_0bpdi0b",
												},
												Name: "",
											},
											Incoming: []string{"Flow_1bn59o5"},
											Outgoing: []string{"Flow_18wfjpz", "Flow_16iauun"},
										},
									},
									Default: "Flow_16iauun",
								},
							},
							SendTasks: []element.SendTask{
								{
									Task: element.Task{
										Activity: element.Activity{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Activity_0kxqh4y",
													},
													Name: "Notify Employee of Approval",
												},
												Incoming: []string{"Flow_1slje8w"},
												Outgoing: []string{"Flow_17w7s53"},
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
														ID: "Activity_1exp9va",
													},
													Name: "Notify Employee of Approval",
												},
												Incoming: []string{"Flow_18wfjpz"},
												Outgoing: []string{"Flow_1kjzd51"},
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
														ID: "Activity_02k6h6m",
													},
													Name: "Notify Employee of Refusal",
												},
												Incoming: []string{"Flow_16iauun"},
												Outgoing: []string{"Flow_1uc40yw"},
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
														ID: "Activity_0wd8gll",
													},
													Name: "Notify Employee of Refusal",
												},
												Incoming: []string{"Flow_1x3mfjv"},
												Outgoing: []string{"Flow_1x4un6h"},
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
														ID: "Activity_14g794r",
													},
													Name: "Manually Approve Vacation",
												},
												Incoming: []string{"Flow_06f2rqz"},
												Outgoing: []string{"Flow_1bn59o5"},
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
														ID: "Event_0fselnt",
													},
													Name: "Vacation Approved by Manager",
												},
												Incoming: []string{"Flow_1s1ikgt"},
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
														ID: "Event_0jqpsr9",
													},
													Name: "Vacation Approved Automatically",
												},
												Incoming: []string{"Flow_0dhwpuk"},
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
														ID: "Event_10bj8rf",
													},
													Name: "Vacation Refused by Manager",
												},
												Incoming: []string{"Flow_1uc40yw"},
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
														ID: "Event_0990itm",
													},
													Name: "Vacation Refused Automatically",
												},
												Incoming: []string{"Flow_1x4un6h"},
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
														ID: "Event_0c4c8gy",
													},
													Name: "Employee not found",
												},
												Incoming: []string{"Flow_15v234s"},
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
