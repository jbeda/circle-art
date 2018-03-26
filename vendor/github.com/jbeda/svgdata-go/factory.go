// Copyright 2018 Joe Beda
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package svgdata

import "encoding/xml"

type nodeCreator func() Node

var factoryMap map[string]nodeCreator = map[string]nodeCreator{}
var unknownCreator nodeCreator

func CreateNodeFromName(n xml.Name) Node {
	if n.Space != SvgNs {
		return nil
	}

	creator, ok := factoryMap[n.Local]
	if !ok {
		creator = unknownCreator
	}

	return creator()
}

func RegisterNodeCreator(n string, c nodeCreator) {
	factoryMap[n] = c
}
