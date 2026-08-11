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
	"strings"
)

type ExpressionUnMarshal struct {
	Type                   ExpressionType              `xml:"http://www.w3.org/2001/XMLSchema-element type,attr"` // xsi:type
	ExpressionSubstitution ExpressionSubstitutionGroup `xml:",omitempty"`
}

func (e *ExpressionUnMarshal) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local == "type" {
			parts := strings.Split(attr.Value, ":")
			if len(parts) > 1 {
				e.Type = ExpressionType(parts[1])
			} else {
				e.Type = ExpressionType(attr.Value)
			}
		}
	}
	if len(e.Type) == 0 {
		e.Type = ""
	}

	t := GetExpressionSubstitutionGroup(e.Type)
	if err := d.DecodeElement(t, &start); err != nil {
		return err
	}
	e.ExpressionSubstitution = t
	return nil
}
