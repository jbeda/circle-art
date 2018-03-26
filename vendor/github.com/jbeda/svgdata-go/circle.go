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

// Circle is an SVG element that has no specialized code or representation.
type Circle struct {
	nodeImpl
	Center geom.Coord
	Radius float64
}

func init() {
	RegisterNodeCreator("circle", func() Node { return createCircle() })
}

func createCircle() *Circle {
	c := &Circle{}
	c.nodeImpl.name = "circle"
	c.nodeImpl.onMarshalAttrs = c.marshalAttrs
	c.nodeImpl.onUnmarshalAttrs = c.unmarshalAttrs
	return c
}

func NewCircle(c geom.Coord, r float64) *Circle {
	n := createCircle()
	n.Center = c
	n.Radius = r
	return n
}

func (c *Circle) marshalAttrs(am AttrMap) {
	am["cx"] = floatToString(c.Center.X)
	am["cy"] = floatToString(c.Center.Y)
	am["r"] = floatToString(c.Radius)
}

func (c *Circle) unmarshalAttrs(am AttrMap) error {
	cx, err := parseValue(am["cx"])
	if err != nil {
		return err
	}
	delete(am, "cx")

	cy, err := parseValue(am["cy"])
	if err != nil {
		return err
	}
	delete(am, "cy")

	r, err := parseValue(am["r"])
	if err != nil {
		return err
	}
	delete(am, "r")

	c.Center = geom.Coord{X: cx, Y: cy}
	c.Radius = r

	return nil
}
