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

import (
	"encoding/xml"
)

// BoundaryEvent is a catch event attached to an activity.
// BPMN XSD default for cancelActivity is true (interrupting); missing attr → true.
type BoundaryEvent struct {
	CatchEvent
	AttachedToRef  string `xml:"attachedToRef,attr"`
	CancelActivity bool   `xml:"-"`
}

func (b *BoundaryEvent) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	b.CancelActivity = true // XSD default
	type alias BoundaryEvent
	aux := struct {
		*alias
		CancelActivity *bool `xml:"cancelActivity,attr"`
	}{alias: (*alias)(b)}
	if err := d.DecodeElement(&aux, &start); err != nil {
		return err
	}
	if aux.CancelActivity != nil {
		b.CancelActivity = *aux.CancelActivity
	}
	return nil
}
