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

func TestA_2_1_roundtrip(t *testing.T) {
	// create test using ./test/A.2.1-roundtrip.bpmn
	path := "./test/A.2.1-roundtrip.bpmn"
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
			ID:                 "Bpmn_Definitions_--SwsH2BEeWQ6qGdY3x14w",
			Name:               "A.2.1",
			TargetNamespace:    "http://bonitasoft.com/_To9ZoDOCEeSknpIVFCxNIQ",
			ExpressionLanguage: "http://groovy.codehaus.org/",
			Exporter:           "BonitaSoft",
			ExporterVersion:    "6.3.3",
			RootElemnts: element.RootElemnts{
				Processes: []element.Process{
					{
						CallableElement: element.CallableElement{
							RootElement: element.RootElement{
								BaseElement: element.BaseElement{
									ID: "_To9ZoTOCEeSknpIVFCxNIQ",
									ExtensionElements: element.ExtensionElements{
										Any: []xml.Token{nil},
									},
								},
							},
							IoSpecification: element.IoSpecification{
								InputOutputSpecification: element.InputOutputSpecification{
									BaseElement: element.BaseElement{
										ID: "_cVGqYDOCEeSknpIVFCxNIQ",
									},
									InputSets: []element.InputSet{
										{
											BaseElement: element.BaseElement{
												ID: "_cVHRcDOCEeSknpIVFCxNIQ",
											},
										},
									},
									OutputSets: []element.OutputSet{
										{
											BaseElement: element.BaseElement{
												ID: "_cVH4gDOCEeSknpIVFCxNIQ",
											},
										},
									},
								},
							},
							Name: "A.2.1",
						},
						IsExecutable: false,
						ProcessType:  element.TProcessTypeNone,
						FlowElements: element.FlowElements{
							StartEvents: []element.StartEvent{{
								CatchEvent: element.CatchEvent{
									Event: element.Event{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID:                "_To9ZojOCEeSknpIVFCxNIQ",
													ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
												},
												Name: "Start Event",
											},
											Outgoing: []string{"_To9Z5DOCEeSknpIVFCxNIQ"},
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
													ID:                "_To9ZpzOCEeSknpIVFCxNIQ",
													ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
												},
												Name: "Task 1",
											},
											Incoming: []string{"_To9Z5DOCEeSknpIVFCxNIQ"},
											Outgoing: []string{"_To9Z5zOCEeSknpIVFCxNIQ"},
										},
									},
								},
								{
									Activity: element.Activity{
										Default: "Bpmn_SequenceFlow_edepQQbbEealeL5I4Yl3Dw",
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID:                "_To9ZtjOCEeSknpIVFCxNIQ",
													ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
												},
												Name: "Task 2",
											},
											Incoming: []string{"_To9Z6jOCEeSknpIVFCxNIQ"},
											Outgoing: []string{"_To9Z7TOCEeSknpIVFCxNIQ", "Bpmn_SequenceFlow_edepQQbbEealeL5I4Yl3Dw"},
										},
									},
								},
								{
									Activity: element.Activity{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID:                "_To9ZwDOCEeSknpIVFCxNIQ",
													ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
												},
												Name: "Task 3",
											},
											Incoming: []string{"_To9Z-TOCEeSknpIVFCxNIQ", "Bpmn_SequenceFlow_edepQQbbEealeL5I4Yl3Dw", "Bpmn_SequenceFlow_f9nmUQbbEealeL5I4Yl3Dw"},
											Outgoing: []string{"_To9Z8DOCEeSknpIVFCxNIQ"},
										},
									},
								},
								{
									Activity: element.Activity{
										Default: "Bpmn_SequenceFlow_f9nmUQbbEealeL5I4Yl3Dw",
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID:                "_To9ZzzOCEeSknpIVFCxNIQ",
													ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
												},
												Name: "Task 4",
											},
											Incoming: []string{"_To9Z_DOCEeSknpIVFCxNIQ"},
											Outgoing: []string{"_To9Z8zOCEeSknpIVFCxNIQ", "Bpmn_SequenceFlow_f9nmUQbbEealeL5I4Yl3Dw"},
										},
									},
								},
							},
							ExclusiveGatewaies: []element.ExclusiveGateway{
								{
									Default: "_To9Z6jOCEeSknpIVFCxNIQ",
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID:                "_To9ZyjOCEeSknpIVFCxNIQ",
													ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
												},
												Name: "Gateway\n\n(Split Flow)",
											},
											Incoming: []string{"_To9Z5zOCEeSknpIVFCxNIQ"},
											Outgoing: []string{"_To9Z6jOCEeSknpIVFCxNIQ", "_To9Z-TOCEeSknpIVFCxNIQ", "_To9Z_DOCEeSknpIVFCxNIQ"},
										},
									},
								},
								{
									Gateway: element.Gateway{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID:                "_To9Z2TOCEeSknpIVFCxNIQ",
													ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
												},
												Name: "Gateway\n(Merge Flows)",
											},
											Incoming: []string{"_To9Z8DOCEeSknpIVFCxNIQ", "_To9Z8zOCEeSknpIVFCxNIQ"},
											Outgoing: []string{"_To9Z9jOCEeSknpIVFCxNIQ"},
										},
									},
								},
							},
							SequenceFlows: []element.SequenceFlow{
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_To9Z5DOCEeSknpIVFCxNIQ",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
										},
									},
									SourceRef: "_To9ZojOCEeSknpIVFCxNIQ",
									TargetRef: "_To9ZpzOCEeSknpIVFCxNIQ",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_To9Z5zOCEeSknpIVFCxNIQ",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
										},
									},
									SourceRef: "_To9ZpzOCEeSknpIVFCxNIQ",
									TargetRef: "_To9ZyjOCEeSknpIVFCxNIQ",
								},
								{
									FlowElement: element.FlowElement{
										Name: "Default",
										BaseElement: element.BaseElement{
											ID:                "_To9Z6jOCEeSknpIVFCxNIQ",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
										},
									},
									SourceRef: "_To9ZyjOCEeSknpIVFCxNIQ",
									TargetRef: "_To9ZtjOCEeSknpIVFCxNIQ",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_To9Z7TOCEeSknpIVFCxNIQ",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
										},
										Name: "Condition",
									},
									SourceRef: "_To9ZtjOCEeSknpIVFCxNIQ",
									TargetRef: "_To9ZsTOCEeSknpIVFCxNIQ",
									ConditionExpression: element.ExpressionUnMarshal{
										Type: element.ExpressionTypeFormal,
										ExpressionSubstitution: &element.FormalExpression{
											Expression: element.Expression{
												BaseElementWithMixedContent: element.BaseElementWithMixedContent{
													ID: "_cVKUwTOCEeSknpIVFCxNIQ",
												},
											},
											Language: "http://www.w3.org/1999/XPath",
											Value:    "true",
										},
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_To9Z8DOCEeSknpIVFCxNIQ",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
										},
									},
									SourceRef: "_To9ZwDOCEeSknpIVFCxNIQ",
									TargetRef: "_To9Z2TOCEeSknpIVFCxNIQ",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_To9Z8zOCEeSknpIVFCxNIQ",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
										},
										Name: "condition",
									},
									SourceRef: "_To9ZzzOCEeSknpIVFCxNIQ",
									TargetRef: "_To9Z2TOCEeSknpIVFCxNIQ",
									ConditionExpression: element.ExpressionUnMarshal{
										Type: element.ExpressionTypeFormal,
										ExpressionSubstitution: &element.FormalExpression{
											Expression: element.Expression{
												BaseElementWithMixedContent: element.BaseElementWithMixedContent{
													ID: "_cVKUwzOCEeSknpIVFCxNIQ",
												},
											},
											Language: "http://www.w3.org/1999/XPath",
										},
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_To9Z9jOCEeSknpIVFCxNIQ",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
										},
									},
									SourceRef: "_To9Z2TOCEeSknpIVFCxNIQ",
									TargetRef: "_To9ZsTOCEeSknpIVFCxNIQ",
									ConditionExpression: element.ExpressionUnMarshal{
										Type: element.ExpressionTypeFormal,
										ExpressionSubstitution: &element.FormalExpression{
											Expression: element.Expression{
												BaseElementWithMixedContent: element.BaseElementWithMixedContent{
													ID: "_cVKUxDOCEeSknpIVFCxNIQ",
												},
											},
											Language: "http://www.w3.org/1999/XPath",
										},
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_To9Z-TOCEeSknpIVFCxNIQ",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
										},
									},
									SourceRef: "_To9ZyjOCEeSknpIVFCxNIQ",
									TargetRef: "_To9ZwDOCEeSknpIVFCxNIQ",
									ConditionExpression: element.ExpressionUnMarshal{
										Type: element.ExpressionTypeFormal,
										ExpressionSubstitution: &element.FormalExpression{
											Expression: element.Expression{
												BaseElementWithMixedContent: element.BaseElementWithMixedContent{
													ID: "_cVKUxTOCEeSknpIVFCxNIQ",
												},
											},
											Language: "http://www.w3.org/1999/XPath",
										},
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "_To9Z_DOCEeSknpIVFCxNIQ",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
										},
									},
									SourceRef: "_To9ZyjOCEeSknpIVFCxNIQ",
									TargetRef: "_To9ZzzOCEeSknpIVFCxNIQ",
									ConditionExpression: element.ExpressionUnMarshal{
										Type: element.ExpressionTypeFormal,
										ExpressionSubstitution: &element.FormalExpression{
											Expression: element.Expression{
												BaseElementWithMixedContent: element.BaseElementWithMixedContent{
													ID: "_cVKUxjOCEeSknpIVFCxNIQ",
												},
											},
											Language: "http://www.w3.org/1999/XPath",
										},
									},
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "Bpmn_SequenceFlow_edepQQbbEealeL5I4Yl3Dw",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
										},
									},
									IsImmediate: true,
									SourceRef:   "_To9ZtjOCEeSknpIVFCxNIQ",
									TargetRef:   "_To9ZwDOCEeSknpIVFCxNIQ",
								},
								{
									FlowElement: element.FlowElement{
										BaseElement: element.BaseElement{
											ID:                "Bpmn_SequenceFlow_f9nmUQbbEealeL5I4Yl3Dw",
											ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
										},
									},
									IsImmediate: true,
									SourceRef:   "_To9ZzzOCEeSknpIVFCxNIQ",
									TargetRef:   "_To9ZwDOCEeSknpIVFCxNIQ",
								},
							},
							EndEvents: []element.EndEvent{{
								ThrowEvent: element.ThrowEvent{
									Event: element.Event{
										FlowNode: element.FlowNode{
											FlowElement: element.FlowElement{
												BaseElement: element.BaseElement{
													ID:                "_To9ZsTOCEeSknpIVFCxNIQ",
													ExtensionElements: element.ExtensionElements{Any: []xml.Token{nil}},
												},
												Name: "End Event",
											},
											Incoming: []string{"_To9Z7TOCEeSknpIVFCxNIQ", "_To9Z9jOCEeSknpIVFCxNIQ"},
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	}
	if diff := cmp.Diff(expected, modelelement); diff != "" {
		t.Errorf("Unexpected result (-want, +got):\n%s", diff)
	}
}
