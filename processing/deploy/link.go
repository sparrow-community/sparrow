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
	"strings"

	"github.com/sparrow-community/sparrow/bpmn/element"
)

func linkNameFromDefs(id string, defs element.EventDefinitions, fallback string) (string, error) {
	if len(defs.LinkEventDefinitions) != 1 {
		return "", fmt.Errorf("UNSUPPORTED_ELEMENT: %q needs exactly one linkEventDefinition", id)
	}
	other := extraCatchDefinitions(defs) - len(defs.LinkEventDefinitions)
	if other > 0 {
		return "", fmt.Errorf("UNSUPPORTED_ELEMENT: %q is not a link event", id)
	}
	name := strings.TrimSpace(defs.LinkEventDefinitions[0].Name)
	if name == "" {
		name = strings.TrimSpace(fallback)
	}
	if name == "" {
		name = id
	}
	return name, nil
}

func linkThrowSpec(ev element.IntermediateThrowEvent) (throwSpec, error) {
	if len(ev.TimerEventDefinitions) > 0 {
		return throwSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateThrowEvent %q cannot have timerEventDefinition", ev.ID)
	}
	name, err := linkNameFromDefs(ev.ID, ev.EventDefinitions, ev.Name)
	if err != nil {
		return throwSpec{}, err
	}
	if len(ev.Outgoing) > 0 {
		return throwSpec{}, fmt.Errorf("UNSUPPORTED_ELEMENT: link throw %q must not have outgoing sequence flows", ev.ID)
	}
	return throwSpec{Kind: ThrowKindLink, Name: name}, nil
}

func linkCatchSpec(ev element.IntermediateCatchEvent) (string, error) {
	if len(ev.TimerEventDefinitions) > 0 {
		return "", fmt.Errorf("UNSUPPORTED_ELEMENT: intermediateCatchEvent %q is not a link catch", ev.ID)
	}
	name, err := linkNameFromDefs(ev.ID, ev.EventDefinitions, ev.Name)
	if err != nil {
		return "", err
	}
	if len(ev.Incoming) > 0 {
		return "", fmt.Errorf("UNSUPPORTED_ELEMENT: link catch %q must not have incoming sequence flows", ev.ID)
	}
	return name, nil
}

// LinkCatches returns catch element ids with the given link name (sorted for determinism).
func (d *Deployment) LinkCatches(name string) []string {
	var ids []string
	for id, n := range d.linkCatch {
		if n == name {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	return ids
}

// LinkCatchName returns the link name for a catch, if any.
func (d *Deployment) LinkCatchName(id string) (string, bool) {
	n, ok := d.linkCatch[id]
	return n, ok
}

func (d *Deployment) validateLinkPairs() error {
	for id, spec := range d.throwEvents {
		if spec.Kind != ThrowKindLink {
			continue
		}
		if len(d.LinkCatches(spec.Name)) == 0 {
			return fmt.Errorf("UNSUPPORTED_ELEMENT: link throw %q has no catch named %q", id, spec.Name)
		}
	}
	return nil
}
