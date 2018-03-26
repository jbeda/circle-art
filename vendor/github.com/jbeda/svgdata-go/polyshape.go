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

import "github.com/jbeda/geom"

// Polyshape is an SVG element is a shape specified with a list of straight
// lines.
type Polyshape struct {
	nodeImpl

	// TODO: Handle Points
	Points []geom.Coord
}

func init() {
	RegisterNodeCreator("polygon", CreatePolygon)
	RegisterNodeCreator("polyline", CreatePolyline)
}

func CreatePolygon() Node {
	p := &Polyshape{}
	p.nodeImpl.name = "polygon"
	p.nodeImpl.onMarshalAttrs = p.marshalAttrs
	p.nodeImpl.onUnmarshalAttrs = p.unmarshalAttrs

	return p
}

func CreatePolyline() Node {
	p := &Polyshape{}
	p.nodeImpl.name = "polygon"
	p.nodeImpl.onMarshalAttrs = p.marshalAttrs
	p.nodeImpl.onUnmarshalAttrs = p.unmarshalAttrs

	return p
}

func (p *Polyshape) marshalAttrs(am AttrMap) {
}

func (p *Polyshape) unmarshalAttrs(am AttrMap) error {
	return nil
}
