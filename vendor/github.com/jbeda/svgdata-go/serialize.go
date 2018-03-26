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
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"strconv"
)

func Unmarshal(data []byte) (*Root, error) {
	d := xml.NewDecoder(bytes.NewReader(data))

	se, _, err := findNextStart(d, nil)
	if err != nil {
		return nil, err
	}
	if se == nil {
		return nil, fmt.Errorf("no root element found")
	}

	r := &Root{}
	err = r.UnmarshalXML(d, *se)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func Marshal(r *Root, pretty bool) ([]byte, error) {
	var b bytes.Buffer

	b.WriteString("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n")
	//b.WriteString("<!-- Generator: Adobe Illustrator 22.1.0, SVG Export Plug-In . SVG Version: 6.00 Build 0)  -->\n")

	e := xml.NewEncoder(&b)
	if pretty {
		e.Indent("", "  ")
	}

	e.Encode(r)

	return b.Bytes(), nil
}

// findNextStart searches the token stream for the next StartElement.  nil is
// returned if there is no NextElement. An error is returned if an unexpected
// EndElement is found. All other data is ignored/dropped (including CharData).
func findNextStart(d *xml.Decoder, se *xml.StartElement) (*xml.StartElement, string, error) {
	var chardata string
	for {
		t, err := d.Token()
		if err != nil {
			return nil, "", err
		}

		switch tt := t.(type) {
		case xml.StartElement:
			return &tt, chardata, nil
		case xml.EndElement:
			if se != nil && se.End() == tt {
				return nil, chardata, nil
			}
			return nil, "", fmt.Errorf("unexpected EndElement: %r", tt)
		case xml.CharData:
			chardata += string(tt)
		default:
			// Ignore other tokens
			break
		}
	}
}

// readChildren reads a set of SVG Nodes and returns an array
func readChildren(d *xml.Decoder, se *xml.StartElement) ([]Node, string, error) {
	var children []Node
	var chardata string
	for {
		cse, cd, err := findNextStart(d, se)
		if err != nil {
			return nil, "", err
		}
		chardata += cd
		if cse == nil {
			return children, chardata, nil
		}
		child := CreateNodeFromName(cse.Name)
		err = child.UnmarshalXML(d, *cse)
		if err != nil {
			return nil, "", err
		}
		children = append(children, child)
	}
}

func MakeStartElement(name string, am AttrMap) xml.StartElement {
	return xml.StartElement{
		// So XML Namespaces are totally broken in encoding/xml. Gah
		// Name: xml.Name{Space: SvgNs, Local: name},
		Name: xml.Name{Local: name},
		Attr: attrMapSlice(am),
	}
}

func makeAttrMap(attrs []xml.Attr) AttrMap {
	r := make(AttrMap)

	for _, a := range attrs {
		if len(a.Name.Space) != 0 {
			log.Printf("Namespaced attribute, dropping: %s", a.Name)
		}
		_, ok := r[a.Name.Local]
		if ok {
			log.Printf("Repeated attr: %s", a)
		}

		r[a.Name.Local] = a.Value
	}

	return r
}

func attrMapSlice(am AttrMap) []xml.Attr {
	r := make([]xml.Attr, 0, len(am))
	for k, v := range am {
		a := xml.Attr{Name: xml.Name{Local: k}, Value: v}
		r = append(r, a)
	}
	return r
}

func copyAttrMap(am AttrMap) AttrMap {
	r := make(AttrMap, len(am))
	for k, v := range am {
		r[k] = v
	}
	return r
}

func floatToString(f float64) string {
	return strconv.FormatFloat(f, 'g', 4, 64)
}
