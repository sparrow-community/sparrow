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

import "encoding/xml"

type BaseElementInterface interface {
	GetID() string
	GetDocumentation() []Documentation
	GetExtensionElements() ExtensionElements
}

type BaseElement struct {
	ID                string            `xml:"id,attr,omitempty"`
	Documentation     []Documentation   `xml:"documentation,omitempty"`
	ExtensionElements ExtensionElements `xml:"extensionElements,omitempty"`
	Any               []xml.Token       `xml:",any"`
}

func (b BaseElement) GetID() string {
	return b.ID
}

func (b BaseElement) GetDocumentation() []Documentation {
	return b.Documentation
}

func (b BaseElement) GetExtensionElements() ExtensionElements {
	return b.ExtensionElements
}
