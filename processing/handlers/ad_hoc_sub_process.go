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

package handlers

import (
	"fmt"

	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// AdHocSubProcessHandler opens an Ad-Hoc SubProcess scope: the host token parks
// on the scope while the executor enables inner activities.
type AdHocSubProcessHandler struct{}

func (AdHocSubProcessHandler) Type() eventv1.Element_Type {
	return eventv1.Element_TYPE_AD_HOC_SUB_PROCESS
}

func (AdHocSubProcessHandler) OnEnter(in EnterInput) (*Effect, error) {
	spec, ok := in.Deployment.AdHocSubProcessSpecOf(in.ElementID)
	if !ok {
		return nil, fmt.Errorf("NOT_FOUND: adHocSubProcess %q", in.ElementID)
	}
	initial := spec.AdHocInitialActivities()
	if len(initial) == 0 {
		return nil, fmt.Errorf("UNSUPPORTED_ELEMENT: adHocSubProcess %q has no inner activity", in.ElementID)
	}
	activated := &eventv1.Element{Intent: eventv1.Element_INTENT_ACTIVATED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID}
	p, err := attachScopeBoundary(in.Deployment, in.ElementID, in.Now)
	if err != nil {
		return nil, err
	}
	if p != nil {
		activated.Payload = &eventv1.Element_ActivityPayload{ActivityPayload: p}
	}
	return &Effect{
		Records: []*eventv1.Element{
			{Intent: eventv1.Element_INTENT_ACTIVATING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
			activated,
		},
		AdHocEnable: &AdHocEnable{ElementID: in.ElementID, ActivityIDs: initial},
	}, nil
}

func (AdHocSubProcessHandler) OnComplete(in CompleteInput) (*Effect, error) {
	records := []*eventv1.Element{
		{Intent: eventv1.Element_INTENT_COMPLETING, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
		{Intent: eventv1.Element_INTENT_COMPLETED, Type: in.Type, Id: in.ElementID, TokenId: in.TokenID},
	}
	records = append(records, cancelAttachedBoundary(in.Deployment, in.ElementID, in.TokenID)...)
	records = append(records, subscribeCompensation(in.Deployment, in.ElementID, in.TokenID)...)
	return &Effect{Records: records, TakeOutgoing: true}, nil
}
