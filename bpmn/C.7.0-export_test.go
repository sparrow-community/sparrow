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

func TestC_7_0_export(t *testing.T) {
	// create test use ./test/C.7.0-export.bpmn
	path := "./test/C.7.0-export.bpmn"
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
				Collaborations: []element.Collaboration{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "Collaboration_0ekjvyv",
							},
						},
						Participants: []element.Participant{
							{
								BaseElement: element.BaseElement{
									ID: "Participant_1tm3xzh",
								},
								Name:       "EU Bank",
								ProcessRef: "Process_19noqni",
							},
						},
					},
				},
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
						LaneSets: []element.LaneSet{
							{
								BaseElement: element.BaseElement{
									ID: "LaneSet_1jedytp",
								},
								Lanes: []element.Lane{
									{
										BaseElement: element.BaseElement{
											ID: "Lane_0v9qlya",
										},
										Name: "Hiring manager",
										FlowNodeRefs: []string{
											"Event_0fz6grj",
											"Activity_1bdtuqb",
											"Activity_16s4p3v",
											"Gateway_1i2920t",
										},
									},
									{
										BaseElement: element.BaseElement{
											ID: "Lane_198kgx3",
										},
										Name: "Recruitment",
										FlowNodeRefs: []string{
											"Activity_0abrfr7",
											"Gateway_1g9xocl",
											"Activity_0rnylu0",
											"Activity_1qw4rcj",
											"Activity_05ada8y",
											"Gateway_19qn20n",
											"Event_0eh3ntl",
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
											ID: "DataObjectReference_097t49g",
										},
										Name: "Role required",
									},
									DataObjectRef: "DataObject_1qa62ts",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObjectReference_1s5aok5",
										},
										Name: "Description",
									},
									DataObjectRef: "DataObject_1vyoqxz",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObjectReference_0wqgr90",
										},
										Name: "Advertisement",
									},
									DataObjectRef: "DataObject_1gqt4l3",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObjectReference_0iennvg",
										},
										Name: "Advertisement [Approved]",
									},
									DataObjectRef: "DataObject_1ar5qpp",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObjectReference_13o5hr9",
										},
										Name: "Selected platforms",
									},
									DataObjectRef: "DataObject_1uesvoq",
								},
							},
							DataObjects: []element.DataObject{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObject_1qa62ts",
										},
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObject_1vyoqxz",
										},
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObject_1gqt4l3",
										},
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObject_1ar5qpp",
										},
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "DataObject_1uesvoq",
										},
									},
									IsCollection: false,
								},
							},
							StartEvents: []element.StartEvent{
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "Event_0fz6grj",
													},
													Name: "Job vacancy",
												},
												Outgoing: []string{"Flow_0vycspr"},
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
														ID: "Activity_1bdtuqb",
													},
													Name: "Write description",
												},
												Incoming: []string{"Flow_0vycspr"},
												Outgoing: []string{"Flow_0e56yhs"},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "DataOutputAssociation_1s64n6p",
														},
														TargetRef: "DataObjectReference_1s5aok5",
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
														ID: "Activity_0abrfr7",
													},
													Name: "Complete advertisement",
												},
												Incoming: []string{"Flow_0e56yhs", "Flow_14ytgtt"},
												Outgoing: []string{"Flow_0xsfsj0"},
											},
											Properties: []element.Property{
												{
													BaseElement: element.BaseElement{
														ID: "Property_1nq6xew",
													},
													Name: "__targetRef_placeholder",
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "DataInputAssociation_0bne97i",
														},
														SourceRef: "DataObjectReference_1s5aok5",
														TargetRef: "Property_1nq6xew",
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "DataOutputAssociation_1srlezz",
														},
														TargetRef: "DataObjectReference_0wqgr90",
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
														ID: "Activity_16s4p3v",
													},
													Name: "Approve advertisement",
												},
												Incoming: []string{"Flow_0xsfsj0"},
												Outgoing: []string{"Flow_1phkfiv"},
											},
											Properties: []element.Property{
												{
													BaseElement: element.BaseElement{
														ID: "Property_1k8x8kh",
													},
													Name: "__targetRef_placeholder",
												},
											},
											DataInputAssociations: []element.DataInputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "DataInputAssociation_0b75hsa",
														},
														SourceRef: "DataObjectReference_0wqgr90",
														TargetRef: "Property_1k8x8kh",
													},
												},
											},
											DataOutputAssociations: []element.DataOutputAssociation{
												{
													DataAssociation: element.DataAssociation{
														BaseElement: element.BaseElement{
															ID: "DataOutputAssociation_1jm5pv0",
														},
														TargetRef: "DataObjectReference_0iennvg",
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
													ID: "Gateway_1i2920t",
												},
												Name: "Advertisement approved?",
											},
											Incoming: []string{"Flow_1phkfiv"},
											Outgoing: []string{"Flow_14ytgtt", "Flow_0puyce6"},
										},
									},
								},
							},
							ParallelGatewaies: []element.ParallelGateway{
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Gateway_1g9xocl",
												},
											},
											Incoming: []string{"Flow_0puyce6"},
											Outgoing: []string{"Flow_1phvh3p", "Flow_1reqq1i"},
										},
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "Gateway_19qn20n",
												},
											},
											Incoming: []string{"Flow_0flt1nf", "Flow_0l0n3nh"},
											Outgoing: []string{"Flow_0s3ij0s"},
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
														ID: "Activity_0rnylu0",
													},
													Name: "Publish on homepage",
												},
												Incoming: []string{"Flow_1phvh3p"},
												Outgoing: []string{"Flow_0flt1nf"},
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
														ID: "Activity_05ada8y",
													},
													Name: "Publish on other platforms",
												},
												Incoming: []string{"Flow_0l6dn1e"},
												Outgoing: []string{"Flow_0l0n3nh"},
											},
											LoopCharacteristicsElements: element.LoopCharacteristicsElements{
												MultielementLoopCharacteristics: []element.MultiInstanceLoopCharacteristics{
													{},
												},
											},
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
														ID: "Activity_1qw4rcj",
													},
													Name: "Select other platforms",
												},
												Incoming: []string{"Flow_1reqq1i"},
												Outgoing: []string{"Flow_0l6dn1e"},
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
														ID: "Event_0eh3ntl",
													},
													Name: "Vacancy advertised",
												},
												Incoming: []string{"Flow_0s3ij0s"},
											},
										},
									},
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0vycspr",
										},
									},
									SourceRef: "Event_0fz6grj",
									TargetRef: "Activity_1bdtuqb",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0e56yhs",
										},
									},
									SourceRef: "Activity_1bdtuqb",
									TargetRef: "Activity_0abrfr7",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0xsfsj0",
										},
									},
									SourceRef: "Activity_0abrfr7",
									TargetRef: "Activity_16s4p3v",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1phkfiv",
										},
									},
									SourceRef: "Activity_16s4p3v",
									TargetRef: "Gateway_1i2920t",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_14ytgtt",
										},
										Name: "No",
									},
									SourceRef: "Gateway_1i2920t",
									TargetRef: "Activity_0abrfr7",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0puyce6",
										},
										Name: "Yes",
									},
									SourceRef: "Gateway_1i2920t",
									TargetRef: "Gateway_1g9xocl",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1phvh3p",
										},
									},
									SourceRef: "Gateway_1g9xocl",
									TargetRef: "Activity_0rnylu0",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_1reqq1i",
										},
									},
									SourceRef: "Gateway_1g9xocl",
									TargetRef: "Activity_1qw4rcj",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0l6dn1e",
										},
									},
									SourceRef: "Activity_1qw4rcj",
									TargetRef: "Activity_05ada8y",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0flt1nf",
										},
									},
									SourceRef: "Activity_0rnylu0",
									TargetRef: "Gateway_19qn20n",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0l0n3nh",
										},
									},
									SourceRef: "Activity_05ada8y",
									TargetRef: "Gateway_19qn20n",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "Flow_0s3ij0s",
										},
									},
									SourceRef: "Gateway_19qn20n",
									TargetRef: "Event_0eh3ntl",
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
