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
type Rect struct {
	nodeImpl
	R geom.Rect
}

func init() {
	RegisterNodeCreator("rect", func() Node { return createRect() })
}

func createRect() *Rect {
	r := &Rect{}
	r.nodeImpl.name = "rect"
	r.nodeImpl.onMarshalAttrs = r.marshalAttrs
	r.nodeImpl.onUnmarshalAttrs = r.unmarshalAttrs
	return r
}

func NewRect(r geom.Rect) *Rect {
	n := createRect()
	n.R = r
	return n
}

func NewRectXYWH(x, y, w, h float64) *Rect {
	return NewRect(geom.Rect{
		geom.Coord{x, y},
		geom.Coord{x + w, y + h},
	})
}

func (r *Rect) marshalAttrs(am AttrMap) {
	am["x"] = floatToString(r.R.Min.X)
	am["y"] = floatToString(r.R.Min.Y)
	am["width"] = floatToString(r.R.Width())
	am["height"] = floatToString(r.R.Height())
}

func (r *Rect) unmarshalAttrs(am AttrMap) error {
	x, err := am.ExtractValue("x")
	if err != nil {
		return err
	}

	y, err := am.ExtractValue("y")
	if err != nil {
		return err
	}

	width, err := am.ExtractValueNoDefault("width")
	if err != nil {
		return err
	}

	height, err := am.ExtractValueNoDefault("height")
	if err != nil {
		return nil
	}

	r.R.Min = geom.Coord{x, y}
	r.R.Max = geom.Coord{x + width, y + height}

	return nil
}
