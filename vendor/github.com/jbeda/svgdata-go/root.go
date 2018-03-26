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

// Root represents the root <svg> element.
type Root struct {
	nodeImpl
}

var _ Node = (*Root)(nil)

func CreateRoot() *Root {
	r := &Root{}
	r.nodeImpl.name = "svg"
	r.nodeImpl.onMarshalAttrs = r.marshalAttrs
	r.nodeImpl.onUnmarshalAttrs = r.unmarshalAttrs
	return r
}

func init() {
	RegisterNodeCreator("svg", func() Node { return CreateRoot() })
}

func (r *Root) marshalAttrs(am AttrMap) {

}

func (r *Root) unmarshalAttrs(am AttrMap) error {
	return nil
}
