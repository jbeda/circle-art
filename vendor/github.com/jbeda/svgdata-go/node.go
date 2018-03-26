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

import (
	"encoding/xml"
	"fmt"
)

const (
	SvgNs = "http://www.w3.org/2000/svg"
)

type onMarshalAttrsFunc func(attr AttrMap)
type onUnmarshalAttrsFunc func(attr AttrMap) error

type Node interface {
	xml.Unmarshaler
	xml.Marshaler

	Name() string

	Attrs() AttrMap

	Children() *[]Node
	AddChild(n Node)

	GetText() string
	SetText(t string)
}

type nodeImpl struct {
	name     string
	attrs    AttrMap
	children []Node
	text     string

	// onMarshalAddrs is called during marshalling with a copy of the Attrs that the can be modified before marshalling.
	onMarshalAttrs onMarshalAttrsFunc
	// onUmarshalAttrs is called during unmarshalling.  The calling function can modify the Attrs as necessary and
	// changes will be stored with the node.
	onUnmarshalAttrs onUnmarshalAttrsFunc
}

var _ Node = (*nodeImpl)(nil)

func (n *nodeImpl) Name() string {
	return n.name
}

func (n *nodeImpl) Attrs() AttrMap {
	if n.attrs == nil {
		n.attrs = AttrMap{}
	}
	return n.attrs
}

func (n *nodeImpl) Children() *[]Node {
	return &n.children
}

func (n *nodeImpl) AddChild(c Node) {
	n.children = append(n.children, c)
}

func (n *nodeImpl) GetText() string {
	return n.text
}

func (n *nodeImpl) SetText(t string) {
	n.text = t
}

func (n *nodeImpl) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var err error

	if start.Name.Space != SvgNs {
		return fmt.Errorf("parsing non-SVG element: %v", start.Name)
	}
	n.name = start.Name.Local
	n.attrs = makeAttrMap(start.Attr)

	if n.onUnmarshalAttrs != nil {
		err = n.onUnmarshalAttrs(n.attrs)
		if err != nil {
			return err
		}
	}

	n.children, n.text, err = readChildren(d, &start)
	return err
}

func (n *nodeImpl) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	am := copyAttrMap(n.attrs)

	if n.onMarshalAttrs != nil {
		n.onMarshalAttrs(am)
	}

	se := MakeStartElement(n.name, am)

	// Do the namespace thing for the root. This is a total hack.
	if n.name == "svg" {
		se.Name.Space = SvgNs
	}

	err := e.EncodeToken(se)
	if err != nil {
		return err
	}

	if len(n.children) != 0 {
		for _, c := range n.children {
			err = e.Encode(c)
			if err != nil {
				return err
			}
		}
	} else if n.text != "" {
		tt := xml.CharData(n.text)
		err = e.EncodeToken(tt)
		if err != nil {
			return err
		}
	}

	err = e.EncodeToken(se.End())
	if err != nil {
		return err
	}

	return nil
}
