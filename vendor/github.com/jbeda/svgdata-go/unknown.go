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

// Unknown is an SVG element that has no specialized code or representation.
type Unknown struct {
	nodeImpl
}

func init() {
	unknownCreator = func() Node { return createUnknown() }
}

func createUnknown() *Unknown {
	return &Unknown{}
}

func NewGroup() Node {
	u := createUnknown()
	u.name = "g"
	return u
}

func NewStyle() Node {
	u := createUnknown()
	u.name = "style"
	return u
}
