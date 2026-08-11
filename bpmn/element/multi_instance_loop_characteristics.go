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

package element

type MultielementFlowCondition string

const (
	MultielementFlowConditionNone    MultielementFlowCondition = "None"
	MultielementFlowConditionOne     MultielementFlowCondition = "One"
	MultielementFlowConditionAll     MultielementFlowCondition = "All"
	MultielementFlowConditionComplex MultielementFlowCondition = "Complex"
)

type MultiInstanceLoopCharacteristics struct {
	LoopCharacteristics
	IsSequential               bool                        `xml:"isSequential,attr"`
	Behavior                   MultielementFlowCondition   `xml:"behavior,attr"`
	OneBehaviorEventRef        string                      `xml:"oneBehaviorEventRef,attr"`
	NoneBehaviorEventRef       string                      `xml:"noneBehaviorEventRef,attr"`
	LoopCardinality            ExpressionUnMarshal         `xml:"loopCardinality"`
	LoopDataInputRef           string                      `xml:"loopDataInputRef,attr"`
	LoopDataOutputRef          string                      `xml:"loopDataOutputRef,attr"`
	InputDataItem              DataInput                   `xml:"inputDataItem"`
	OutputDataItem             DataOutput                  `xml:"outputDataItem"`
	ComplexBehaviorDefinitions []ComplexBehaviorDefinition `xml:"complexBehaviorDefinition"`
	CompletionCondition        ExpressionUnMarshal         `xml:"completionCondition"`
}
