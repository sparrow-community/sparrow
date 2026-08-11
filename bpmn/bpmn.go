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
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"os"

	"github.com/sparrow-community/sparrow/bpmn/element"
	"golang.org/x/net/html/charset"
)

type Bpmn struct {
	Definitions *element.Definitions
}

// unmarshal unmarshals the given byte array into the given value, different charsets are supported
func unmarshal(data []byte, v interface{}) error {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = charset.NewReaderLabel
	return decoder.Decode(v)
}

// BpmnModelelementFromBytes reads a BPMN model element from a byte array
func BpmnModelelementFromBytes(data []byte) (*Bpmn, error) {
	definitions := &element.Definitions{}
	err := unmarshal(data, definitions)
	if err != nil {
		return nil, err
	}

	return &Bpmn{
		Definitions: definitions,
	}, nil
}

// BpmnModelelementFromFile reads a BPMN model element from a file
func BpmnModelelementFromFile(filePath string) (*Bpmn, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return BpmnModelelementFromBytes(data)
}

// BpmnModelelementFromStream reads a BPMN model element from a stream
func BpmnModelelementFromStream(stream io.Reader) (*Bpmn, error) {
	data, err := io.ReadAll(stream)
	if err != nil {
		return nil, err
	}
	return BpmnModelelementFromBytes(data)
}

// FindElementByID finds an element by its ID in the BPMN model
func (b *Bpmn) FindElementByID(id string) (element.BaseElementInterface, error) {
	if b.Definitions == nil {
		return nil, errors.New("BPMN definitions are nil")
	}

	for _, process := range b.Definitions.RootElemnts.Processes {
		if process.ID == id {
			return &process, nil
		}
		for _, startEvents := range process.StartEvents {
			if startEvents.ID == id {
				return &startEvents, nil
			}
		}
		for _, sequenceFlows := range process.SequenceFlows {
			if sequenceFlows.ID == id {
				return &sequenceFlows, nil
			}
		}
		for _, endEvents := range process.EndEvents {
			if endEvents.ID == id {
				return &endEvents, nil
			}
		}
		for _, tasks := range process.Tasks {
			if tasks.ID == id {
				return &tasks, nil
			}
		}
		for _, manualTasks := range process.ManualTasks {
			if manualTasks.ID == id {
				return &manualTasks, nil
			}
		}
		for _, userTasks := range process.UserTasks {
			if userTasks.ID == id {
				return &userTasks, nil
			}
		}
		for _, serviceTasks := range process.ServiceTasks {
			if serviceTasks.ID == id {
				return &serviceTasks, nil
			}
		}
		for _, sendTasks := range process.SendTasks {
			if sendTasks.ID == id {
				return &sendTasks, nil
			}
		}
		for _, receiveTasks := range process.ReceiveTasks {
			if receiveTasks.ID == id {
				return &receiveTasks, nil
			}
		}
		for _, businessRuleTasks := range process.BusinessRuleTasks {
			if businessRuleTasks.ID == id {
				return &businessRuleTasks, nil
			}
		}
		for _, dataStoreReferences := range process.DataStoreReferences {
			if dataStoreReferences.ID == id {
				return &dataStoreReferences, nil
			}
		}
		for _, parallelGateways := range process.ParallelGatewaies {
			if parallelGateways.ID == id {
				return &parallelGateways, nil
			}
		}
		for _, exclusiveGateways := range process.ExclusiveGatewaies {
			if exclusiveGateways.ID == id {
				return &exclusiveGateways, nil
			}
		}
		for _, inclusiveGateways := range process.InclusiveGatewaies {
			if inclusiveGateways.ID == id {
				return &inclusiveGateways, nil
			}
		}
		for _, eventBasedGateways := range process.EventBasedGatewaies {
			if eventBasedGateways.ID == id {
				return &eventBasedGateways, nil
			}
		}
		for _, subProcesses := range process.SubProcesses {
			subChild, err := b.FindElementByID(id)
			if err == nil {
				return subChild, nil
			}
			if subProcesses.ID == id {
				return &subProcesses, nil
			}
		}
		for _, boundaryEvents := range process.BoundaryEvents {
			if boundaryEvents.ID == id {
				return &boundaryEvents, nil
			}
		}
		for _, callActivities := range process.CallActivities {
			if callActivities.ID == id {
				return &callActivities, nil
			}
		}
		for _, dataObjectReferences := range process.DataObjectReferenes {
			if dataObjectReferences.ID == id {
				return &dataObjectReferences, nil
			}
		}
		for _, dataObjects := range process.DataObjects {
			if dataObjects.ID == id {
				return &dataObjects, nil
			}
		}
		for _, intermediateThrowEvents := range process.IntermediateThrowEvents {
			if intermediateThrowEvents.ID == id {
				return &intermediateThrowEvents, nil
			}
		}
		for _, intermediateCatchEvents := range process.IntermediateCatchEvents {
			if intermediateCatchEvents.ID == id {
				return &intermediateCatchEvents, nil
			}
		}
		for _, implicitThrowEvents := range process.ImplicitThrowEvents {
			if implicitThrowEvents.ID == id {
				return &implicitThrowEvents, nil
			}
		}
	}

	return nil, errors.New("element not found")
}
