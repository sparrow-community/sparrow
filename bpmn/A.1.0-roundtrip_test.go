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

func Test_A_1_0_roundtrip(t *testing.T) {
	// create test use ./test/A.1.0-roundtrip.bpmn
	path := "./test/A.1.0-roundtrip.bpmn"
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
			ID:              "_1373649849716",
			Name:            "A.1.0",
			TargetNamespace: "http://www.trisotech.com/definitions/_1373649849716",
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
							StartEvents: []element.StartEvent{
								{
									CatchEvent: element.CatchEvent{
										Event: element.Event{
											FlowNode: element.FlowNode{
												FlowElement: element.FlowElement{
													BaseElement: element.BaseElement{
														ID: "_93c466ab-b271-4376-a427-f4c353d55ce8",
													},
													Name: "Start Event",
												},
												Outgoing: []string{"_e16564d7-0c4c-413e-95f6-f668a3f851fb"},
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
													ID: "_ec59e164-68b4-4f94-98de-ffb1c58a84af",
												},
												Name: "Task 1",
											},
											Incoming: []string{"_e16564d7-0c4c-413e-95f6-f668a3f851fb"},
											Outgoing: []string{"_d77dd5ec-e4e7-420e-bbe7-8ac9cd1df599"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_820c21c0-45f3-473b-813f-06381cc637cd",
												},
												Name: "Task 2",
											},
											Incoming: []string{"_d77dd5ec-e4e7-420e-bbe7-8ac9cd1df599"},
											Outgoing: []string{"_2aa47410-1b0e-4f8b-ad54-d6f798080cb4"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID: "_e70a6fcb-913c-4a7b-a65d-e83adc73d69c",
												},
												Name: "Task 3",
											},
											Incoming: []string{"_2aa47410-1b0e-4f8b-ad54-d6f798080cb4"},
											Outgoing: []string{"_8e8fe679-eb3b-4c43-a4d6-891e7087ff80"},
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
														ID: "_a47df184-085b-49f7-bb82-031c84625821",
													},
													Name: "End Event",
												},
												Incoming: []string{"_8e8fe679-eb3b-4c43-a4d6-891e7087ff80"},
											},
										},
									},
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_e16564d7-0c4c-413e-95f6-f668a3f851fb",
										},
										Name: "",
									},
									SourceRef: "_93c466ab-b271-4376-a427-f4c353d55ce8",
									TargetRef: "_ec59e164-68b4-4f94-98de-ffb1c58a84af",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_d77dd5ec-e4e7-420e-bbe7-8ac9cd1df599",
										},
									},
									SourceRef: "_ec59e164-68b4-4f94-98de-ffb1c58a84af",
									TargetRef: "_820c21c0-45f3-473b-813f-06381cc637cd",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_2aa47410-1b0e-4f8b-ad54-d6f798080cb4",
										},
									},
									SourceRef: "_820c21c0-45f3-473b-813f-06381cc637cd",
									TargetRef: "_e70a6fcb-913c-4a7b-a65d-e83adc73d69c",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID: "_8e8fe679-eb3b-4c43-a4d6-891e7087ff80",
										},
									},
									SourceRef: "_e70a6fcb-913c-4a7b-a65d-e83adc73d69c",
									TargetRef: "_a47df184-085b-49f7-bb82-031c84625821",
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
