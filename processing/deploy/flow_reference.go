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

package deploy

import (
	"fmt"
	"sort"

	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// validateFlowReferences rejects flow wiring that cannot be executed: a sequence
// flow whose endpoint is not an indexed element (an unsupported or non-executable
// flow element the parser dropped), or a declared incoming/outgoing that names no
// sequence flow. Without this, such a definition deploys and only fails when a
// token reaches the gap. Ids are visited in sorted order so the first failure is
// the same on every deploy.
func (d *Deployment) validateFlowReferences() error {
	flowIDs := make([]string, 0, len(d.seqFlows))
	for id := range d.seqFlows {
		flowIDs = append(flowIDs, id)
	}
	sort.Strings(flowIDs)
	for _, flowID := range flowIDs {
		f := d.seqFlows[flowID]
		for _, ref := range []struct {
			kind string
			id   string
		}{
			{"sourceRef", f.SourceID},
			{"targetRef", f.TargetID},
		} {
			if ref.id == "" {
				return fmt.Errorf("UNSUPPORTED_ELEMENT: sequenceFlow %q has no %s", flowID, ref.kind)
			}
			if _, ok := d.elements[ref.id]; !ok {
				return fmt.Errorf("UNSUPPORTED_ELEMENT: sequenceFlow %q %s %q is not a supported flow element", flowID, ref.kind, ref.id)
			}
		}
	}

	elementIDs := make([]string, 0, len(d.elements))
	for id := range d.elements {
		elementIDs = append(elementIDs, id)
	}
	sort.Strings(elementIDs)
	for _, elementID := range elementIDs {
		e := d.elements[elementID]
		for _, ref := range []struct {
			kind string
			ids  []string
		}{
			{"outgoing", e.Outgoing},
			{"incoming", e.Incoming},
		} {
			for _, flowID := range ref.ids {
				target, ok := d.elements[flowID]
				if !ok || target.Type != eventv1.Element_TYPE_SEQUENCE_FLOW {
					return fmt.Errorf("UNSUPPORTED_ELEMENT: element %q declares %s %q which is not a sequence flow", elementID, ref.kind, flowID)
				}
			}
		}
	}
	return nil
}
