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

func TestA_4_1_roundtrip(t *testing.T) {
	path := "./test/A.4.1-roundtrip.bpmn"
	modelelement, err := BpmnModelelementFromFile(path)
	if err != nil {
		t.Fatalf("BpmnModelelementFromFile error %s: %s", path, err)
	}
	if modelelement == nil {
		t.Fatalf("modelelement is nil, %s", path)
	}

	expected := &Bpmn{
		Definitions: &element.Definitions{
			XMLName: xml.Name{
				Space: "http://www.omg.org/spec/BPMN/20100524/MODEL",
				Local: "definitions",
			},
			ID:              "sid-ad44e239-e96e-4a80-b0e4-cf63b741c3cb",
			TargetNamespace: "http://www.trisotech.com/definitions/_1365195427479",
			Exporter:        "Signavio Process Editor, http://www.signavio.com",
			Name:            "A.4.1",
			ExporterVersion: "7.6.0",
			RootElemnts: element.RootElemnts{
				Collaborations: []element.Collaboration{
					{
						RootElement: element.RootElement{
							BaseElement: element.BaseElement{
								ID: "sid-467b00a2-7f22-4314-bd57-2f84b409dc80",
							},
						},
						Participants: []element.Participant{
							{
								BaseElement: element.BaseElement{
									ID: "sid-66751F1E-EEB9-4BA7-9FDA-7965A1CA9CD1",
									ExtensionElements: element.ExtensionElements{
										Any: []xml.Token{nil},
									},
								},
								Name:       "Pool 1",
								ProcessRef: "sid-34746A54-1D7D-46CA-B219-0C4CEAE51170",
							},
							{
								BaseElement: element.BaseElement{
									ID: "sid-7E61DCD0-0700-4828-8A28-CD65132273D7",
									ExtensionElements: element.ExtensionElements{
										Any: []xml.Token{nil},
									},
								},
								Name:       "Pool 2",
								ProcessRef: "sid-54D696FD-DEDC-45F3-99DB-1404DA433FC4",
							},
						},
						MessageFlows: []element.MessageFlow{
							{
								BaseElement: element.BaseElement{
									ID: "sid-96EF2D8F-C322-42B1-8C08-0DA05524C904",
								},
								Name:      "Message Flow 2 ",
								SourceRef: "sid-485E1184-9951-4B41-9794-A9AFD42A3249",
								TargetRef: "sid-1208A5BA-9E1C-49D2-82E3-5DB2C0E9887D",
							},
							{
								BaseElement: element.BaseElement{
									ID: "sid-D0B859BF-CBFB-4B35-BBC8-BCA308F6455C",
								},
								Name:      "Message Flow 1 ",
								SourceRef: "sid-3D477D07-D669-4A26-9454-12AD775FDE70",
								TargetRef: "sid-34E8C3A5-5C2A-4593-AC67-038B737814D7",
							},
						},
					},
				},
				Processes: []element.Process{
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "sid-34746A54-1D7D-46CA-B219-0C4CEAE51170",
								},
							},
							Name: "Pool 1",
						},
						ProcessType:  element.TProcessTypeNone,
						IsExecutable: false,
						LaneSets: []element.LaneSet{
							{
								BaseElement: element.BaseElement{
									ID: "sid-600cb997-ee14-4230-a60a-fe80e85800b5",
								},
								Lanes: []element.Lane{
									{
										BaseElement: element.BaseElement{
											ID: "sid-4F568BD0-1CB0-4F1C-8729-9DD775B5B37D",
											ExtensionElements: element.ExtensionElements{
												Any: []xml.Token{nil},
											},
										},
										Name:         "Lane 1",
										FlowNodeRefs: []string{"sid-5F0F3508-96EF-4F9B-9182-64AD17334E23", "sid-70D2F83B-77E6-4301-835C-AFF6357344F8", "sid-3D477D07-D669-4A26-9454-12AD775FDE70", "sid-1208A5BA-9E1C-49D2-82E3-5DB2C0E9887D"},
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
														ID:                "sid-70D2F83B-77E6-4301-835C-AFF6357344F8",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
													},
													Name: "Start Event 1 ",
												},
												Outgoing: []string{"sid-576A3375-50D2-4E0B-90AD-CD756E199FB7"},
											},
										},
									},
								},
							},
							Tasks: []element.Task{
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID:                "sid-3D477D07-D669-4A26-9454-12AD775FDE70",
													ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil}},
												},
												Name: "Task 1 ",
											},
											Incoming: []string{"sid-576A3375-50D2-4E0B-90AD-CD756E199FB7"},
											Outgoing: []string{"sid-D1E9B201-87A2-47B9-82A0-1BA208440CAE"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID:                "sid-1208A5BA-9E1C-49D2-82E3-5DB2C0E9887D",
													ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil}},
												},
												Name: "Task 2 ",
											},
											Incoming: []string{"sid-D1E9B201-87A2-47B9-82A0-1BA208440CAE"},
											Outgoing: []string{"sid-F9B17890-98C4-44FA-B7A8-CA940866741B"},
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
														ID:                "sid-5F0F3508-96EF-4F9B-9182-64AD17334E23",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
													},
													Name: "End Event 1 ",
												},
												Incoming: []string{"sid-F9B17890-98C4-44FA-B7A8-CA940866741B"},
											},
										},
									},
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "sid-576A3375-50D2-4E0B-90AD-CD756E199FB7",
										},
									},
									SourceRef:   "sid-70D2F83B-77E6-4301-835C-AFF6357344F8",
									TargetRef:   "sid-3D477D07-D669-4A26-9454-12AD775FDE70",
									IsImmediate: true,
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "sid-D1E9B201-87A2-47B9-82A0-1BA208440CAE",
										},
									},
									SourceRef:   "sid-3D477D07-D669-4A26-9454-12AD775FDE70",
									TargetRef:   "sid-1208A5BA-9E1C-49D2-82E3-5DB2C0E9887D",
									IsImmediate: true,
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "sid-F9B17890-98C4-44FA-B7A8-CA940866741B",
										},
									},
									SourceRef:   "sid-1208A5BA-9E1C-49D2-82E3-5DB2C0E9887D",
									TargetRef:   "sid-5F0F3508-96EF-4F9B-9182-64AD17334E23",
									IsImmediate: true,
								},
							},
						},
					},
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "sid-54D696FD-DEDC-45F3-99DB-1404DA433FC4",
								},
							},
							Name: "Pool 2",
						},
						ProcessType:  element.TProcessTypeNone,
						IsExecutable: false,
						LaneSets: []element.LaneSet{
							{
								BaseElement: element.BaseElement{
									ID: "sid-4e31a46d-899c-4b59-86e2-4d54abea99cd",
								},
								Lanes: []element.Lane{
									{
										BaseElement: element.BaseElement{
											ID:                "sid-FBA8B122-2EFC-4DD5-B714-A13CD36AAA6E",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
										},
										Name:         "Lane 2 ",
										FlowNodeRefs: []string{"sid-78073B2D-35BB-45D5-9CF1-D446602F8E59", "sid-34E8C3A5-5C2A-4593-AC67-038B737814D7", "sid-00A82BF4-1D0A-48DC-8389-C8AAF3E7F754", "sid-485E1184-9951-4B41-9794-A9AFD42A3249", "sid-C189128A-82D2-4E5F-8FB4-F6E21FF27E83"},
									},
									{
										BaseElement: element.BaseElement{
											ID:                "sid-FC452F0B-05C5-4BB2-AA79-F9195F47BD11",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
										},
										Name:         "Lane 3",
										FlowNodeRefs: []string{"sid-93C83C6A-1122-4E0F-9F47-4027C9080456", "sid-645780CC-D61F-4715-8B58-71679305245F"},
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
														ID:                "sid-C189128A-82D2-4E5F-8FB4-F6E21FF27E83",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
													},
													Name: "Start Event 2 ",
												},
												Outgoing: []string{"sid-AD419767-6626-42E7-ADD5-E0EDB9C7975F"},
											},
										},
									},
								},
							},
							Tasks: []element.Task{
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID:                "sid-34E8C3A5-5C2A-4593-AC67-038B737814D7",
													ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil}},
												},
												Name: "Task 3 ",
											},
											Incoming: []string{"sid-AD419767-6626-42E7-ADD5-E0EDB9C7975F"},
											Outgoing: []string{"sid-4052C63C-CB50-4E0C-8901-80D86A1F9759", "sid-0C093502-276D-4B83-A271-2ABE22F335A6"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID:                "sid-485E1184-9951-4B41-9794-A9AFD42A3249",
													ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil}},
												},
												Name: "Task 5 ",
											},
											Incoming: []string{"sid-1DE02844-4989-4A6A-88E7-B75261042119"},
											Outgoing: []string{"sid-B57FC7E5-7709-4E81-A829-2AB8CF5AB3BB"},
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
														ID:                "sid-78073B2D-35BB-45D5-9CF1-D446602F8E59",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
													},
													Name: "End Event 2 ",
												},
												Incoming: []string{"sid-B57FC7E5-7709-4E81-A829-2AB8CF5AB3BB"},
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
														ID:                "sid-93C83C6A-1122-4E0F-9F47-4027C9080456",
														ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
													},
													Name: "End Event 5 ",
												},
												Incoming: []string{"sid-77013C0C-99FE-4BCB-AA8E-1ADDB67DCB6B"},
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
													ID:                "sid-00A82BF4-1D0A-48DC-8389-C8AAF3E7F754",
													ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
												},
												Name: "Expanded Sub-Process 1 ",
											},
											Incoming: []string{"sid-0C093502-276D-4B83-A271-2ABE22F335A6"},
											Outgoing: []string{"sid-1DE02844-4989-4A6A-88E7-B75261042119"},
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
																	ID:                "sid-A9E08E89-FC9E-4519-9A6B-D9347C6AAAAE",
																	ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
																},
																Name: "Start Event 3 ",
															},
															Outgoing: []string{"sid-70CA8C5F-FF45-4403-93C5-44DE37ED60E3"},
														},
													},
												},
											},
										},
										Tasks: []element.Task{
											{
												Activity: element.Activity{
													FlowNode: element.FlowNode{
														FlowElement: element.FlowElement{
															BaseElement: element.BaseElement{
																ID:                "sid-A52AFB6A-43EE-47FE-A95F-057845582F1D",
																ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil}},
															},
															Name: "Task 4 ",
														},
														Incoming: []string{"sid-70CA8C5F-FF45-4403-93C5-44DE37ED60E3"},
														Outgoing: []string{"sid-DE3E0ED7-7F9B-4917-AD34-9C43A6F58918"},
													},
												},
											},
										},
										SequenceFlows: []element.SequenceFlow{
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "sid-70CA8C5F-FF45-4403-93C5-44DE37ED60E3",
													},
												},
												SourceRef:   "sid-A9E08E89-FC9E-4519-9A6B-D9347C6AAAAE",
												TargetRef:   "sid-A52AFB6A-43EE-47FE-A95F-057845582F1D",
												IsImmediate: true,
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "sid-DE3E0ED7-7F9B-4917-AD34-9C43A6F58918",
													},
												},
												SourceRef:   "sid-A52AFB6A-43EE-47FE-A95F-057845582F1D",
												TargetRef:   "sid-E0D38B39-5E32-4FFA-ADC3-5E26F70C7380",
												IsImmediate: true,
											},
										},
										EndEvents: []element.EndEvent{
											{
												ThrowEvent: element.ThrowEvent{
													Event: element.Event{
														FlowNode: element.FlowNode{
															FlowElement: element.FlowElement{
																BaseElement: element.BaseElement{
																	ID:                "sid-E0D38B39-5E32-4FFA-ADC3-5E26F70C7380",
																	ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
																},
																Name: "End Event 3 ",
															},
															Incoming: []string{"sid-DE3E0ED7-7F9B-4917-AD34-9C43A6F58918"},
														},
													},
												},
											},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID:                "sid-645780CC-D61F-4715-8B58-71679305245F",
													ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
												},
												Name: "Expanded Sub-Process 2 ",
											},
											Incoming: []string{"sid-4052C63C-CB50-4E0C-8901-80D86A1F9759"},
											Outgoing: []string{"sid-77013C0C-99FE-4BCB-AA8E-1ADDB67DCB6B"},
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
																	ID:                "sid-1F026F68-099F-44C9-A40E-38A6C9F83D99",
																	ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil}},
																},
																Name: "Start Event 4 ",
															},
															Outgoing: []string{"sid-72E93035-EAF2-4445-AFFE-39C8C0143066"},
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
																	ID:                "sid-46E6675F-8040-45FE-B5C3-B904596F3D4F",
																	ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
																},
																Name: "End Event 4 ",
															},
															Incoming: []string{"sid-4B747910-16CA-4FFD-B92A-8894BB3D7AB6"},
														},
													},
												},
											},
										},
										Tasks: []element.Task{
											{
												Activity: element.Activity{
													FlowNode: element.FlowNode{
														FlowElement: element.FlowElement{
															BaseElement: element.BaseElement{
																ID:                "sid-B414AE83-11A2-4968-B4E4-45833D641928",
																ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil, nil, nil}},
															},
															Name: "Task 6 ",
														},
														Incoming: []string{"sid-72E93035-EAF2-4445-AFFE-39C8C0143066"},
														Outgoing: []string{"sid-4B747910-16CA-4FFD-B92A-8894BB3D7AB6"},
													},
												},
											},
										},
										SequenceFlows: []element.SequenceFlow{
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "sid-72E93035-EAF2-4445-AFFE-39C8C0143066",
													},
												},
												SourceRef:   "sid-1F026F68-099F-44C9-A40E-38A6C9F83D99",
												TargetRef:   "sid-B414AE83-11A2-4968-B4E4-45833D641928",
												IsImmediate: true,
											},
											{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "sid-4B747910-16CA-4FFD-B92A-8894BB3D7AB6",
													},
												},
												SourceRef:   "sid-B414AE83-11A2-4968-B4E4-45833D641928",
												TargetRef:   "sid-46E6675F-8040-45FE-B5C3-B904596F3D4F",
												IsImmediate: true,
											},
										},
									},
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "sid-77013C0C-99FE-4BCB-AA8E-1ADDB67DCB6B",
										},
									},
									SourceRef:   "sid-645780CC-D61F-4715-8B58-71679305245F",
									TargetRef:   "sid-93C83C6A-1122-4E0F-9F47-4027C9080456",
									IsImmediate: true,
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "sid-1DE02844-4989-4A6A-88E7-B75261042119",
										},
									},
									SourceRef:   "sid-00A82BF4-1D0A-48DC-8389-C8AAF3E7F754",
									TargetRef:   "sid-485E1184-9951-4B41-9794-A9AFD42A3249",
									IsImmediate: true,
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "sid-4052C63C-CB50-4E0C-8901-80D86A1F9759",
										},
									},
									SourceRef:   "sid-34E8C3A5-5C2A-4593-AC67-038B737814D7",
									TargetRef:   "sid-645780CC-D61F-4715-8B58-71679305245F",
									IsImmediate: true,
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "sid-AD419767-6626-42E7-ADD5-E0EDB9C7975F",
										},
									},
									SourceRef:   "sid-C189128A-82D2-4E5F-8FB4-F6E21FF27E83",
									TargetRef:   "sid-34E8C3A5-5C2A-4593-AC67-038B737814D7",
									IsImmediate: true,
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "sid-0C093502-276D-4B83-A271-2ABE22F335A6",
										},
									},
									SourceRef:   "sid-34E8C3A5-5C2A-4593-AC67-038B737814D7",
									TargetRef:   "sid-00A82BF4-1D0A-48DC-8389-C8AAF3E7F754",
									IsImmediate: true,
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "sid-B57FC7E5-7709-4E81-A829-2AB8CF5AB3BB",
										},
									},
									SourceRef:   "sid-485E1184-9951-4B41-9794-A9AFD42A3249",
									TargetRef:   "sid-78073B2D-35BB-45D5-9CF1-D446602F8E59",
									IsImmediate: true,
								},
							},
						},
					},
				},
			},
		},
	}

	if diff := cmp.Diff(expected, modelelement); diff != "" {
		t.Errorf("BpmnModelelementFromFile() mismatch (-want +got):\n%s", diff)
	}
}
