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
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Path struct {
	nodeImpl
	SubPaths []SubPath
}

type SubPath struct {
	Commands       []PathCommand
	startX, startY float64 // The start of the path
	endX, endY     float64 // The end of the path
}

type PathCommand struct {
	Command        byte      // The single character command
	Params         []float64 // The parameters to the command
	startX, startY float64   // The start point
	endX, endY     float64   // The point that the command ends on
}

func init() {
	RegisterNodeCreator("path", func() Node { return NewPath() })
}

func NewPath() *Path {
	p := &Path{}
	p.nodeImpl.name = "path"
	p.nodeImpl.onMarshalAttrs = p.marshalAttrs
	p.nodeImpl.onUnmarshalAttrs = p.unmarshalAttrs
	return p
}

func (p *Path) marshalAttrs(am AttrMap) {
	am["d"] = SavePathString(p.SubPaths)
}

func (p *Path) unmarshalAttrs(am AttrMap) error {
	var err error
	p.SubPaths, err = ParsePathString(am["d"])
	if err != nil {
		return err
	}
	delete(am, "d")

	return nil
}

var commandLengths map[byte]int = map[byte]int{
	'm': 2, // move to
	'M': 2,
	'z': 0, // close path
	'Z': 0,
	'l': 2, // line to
	'L': 2,
	'h': 1, // hori line to
	'H': 1,
	'v': 1, // vert line to
	'V': 1,
	'c': 6, // cubic bezier curve to
	'C': 6,
	's': 4, // smooth cubic bezier curve to
	'S': 4,
	'q': 4, // quad bezier curve to
	'Q': 4,
	't': 2, // smooth quad bezier curve to
	'T': 2,
	'a': 7, // arc to
	'A': 7,
}

func isEOD(t pathToken) bool {
	_, ok := t.(eodToken)
	return ok
}

type pathParser struct {
	currX, currY float64
	subPaths     []SubPath
	currSubPath  *SubPath
}

// ParsePathString parses an SVG path string. Each command represents a single
// instance of the command. This is to enable easier further processing.
func ParsePathString(d string) ([]SubPath, error) {
	pp := pathParser{}
	return pp.parse(d)
}

func SavePathString(sps []SubPath) string {
	var buf bytes.Buffer
	for _, sp := range sps {
		for i := 0; i < len(sp.Commands); i++ {
			c := sp.Commands[i]
			buf.WriteByte(c.Command)
			for j, p := range c.Params {
				buf.WriteString(strconv.FormatFloat(p, 'f', -1, 64))
				if j != len(c.Params)-1 {
					buf.WriteByte(' ')
				}
			}
		}
	}

	return buf.String()
}

func (pp *pathParser) parse(d string) ([]SubPath, error) {
	cmds, err := parsePathCommands(d)
	if err != nil {
		return nil, err
	}

	for _, c := range cmds {
		// If we see a start path or don't have a current subpath then we start a
		// new subpath
		if pp.currSubPath == nil || c.Command == 'm' || c.Command == 'M' {
			pp.subPaths = append(pp.subPaths, SubPath{})
			pp.currSubPath = &pp.subPaths[len(pp.subPaths)-1]

			// If the current command isn't 'M' then we inject one
			if c.Command != 'm' && c.Command != 'M' {
				pp.currSubPath.Commands = append(pp.currSubPath.Commands,
					PathCommand{Command: 'M', Params: []float64{pp.currX, pp.currY}})
				pp.updatePositions(&pp.currSubPath.Commands[0])
			}
		}

		pp.currSubPath.Commands = append(pp.currSubPath.Commands, c)
		pp.updatePositions(&pp.currSubPath.Commands[len(pp.currSubPath.Commands)-1])

		// If we see a Z then we close out the subpath
		if c.Command == 'z' || c.Command == 'Z' {
			pp.currSubPath = nil
		}
	}

	return pp.subPaths, nil
}

// updatePositions updates any begin/end positions on the command and subpath
// based on current context.
func (pp *pathParser) updatePositions(c *PathCommand) {
	c.startX, c.startY = pp.currX, pp.currY

	switch c.Command {
	case 'm':
		pp.currX, pp.currY = (pp.currX + c.Params[0]), (pp.currY + c.Params[1])
		pp.currSubPath.startX, pp.currSubPath.startY = pp.currX, pp.currY
	case 'M':
		pp.currX, pp.currY = c.Params[0], c.Params[1]
		pp.currSubPath.startX, pp.currSubPath.startY = pp.currX, pp.currY
	case 'z':
		fallthrough
	case 'Z':
		pp.currX, pp.currY = pp.currSubPath.startX, pp.currSubPath.startY
	case 'l':
		pp.currX, pp.currY = (pp.currX + c.Params[0]), (pp.currY + c.Params[1])
	case 'L':
		pp.currX, pp.currY = c.Params[0], c.Params[1]
	case 'h':
		pp.currX = pp.currX + c.Params[0]
	case 'H':
		pp.currX = c.Params[0]
	case 'v':
		pp.currY = pp.currY + c.Params[0]
	case 'V':
		pp.currY = c.Params[0]
	case 'c':
		pp.currX, pp.currY = (pp.currX + c.Params[4]), (pp.currY + c.Params[5])
	case 'C':
		pp.currX, pp.currY = c.Params[4], c.Params[5]
	case 's':
		pp.currX, pp.currY = (pp.currX + c.Params[2]), (pp.currY + c.Params[3])
	case 'S':
		pp.currX, pp.currY = c.Params[2], c.Params[3]
	case 'q':
		pp.currX, pp.currY = (pp.currX + c.Params[2]), (pp.currY + c.Params[3])
	case 'Q':
		pp.currX, pp.currY = c.Params[2], c.Params[3]
	case 't':
		pp.currX, pp.currY = (pp.currX + c.Params[0]), (pp.currY + c.Params[1])
	case 'T':
		pp.currX, pp.currY = c.Params[0], c.Params[1]
	case 'a':
		pp.currX, pp.currY = (pp.currX + c.Params[5]), (pp.currY + c.Params[6])
	case 'A':
		pp.currX, pp.currY = c.Params[5], c.Params[6]
	}

	c.endX, c.endY = pp.currX, pp.currY
	pp.currSubPath.endX, pp.currSubPath.endY = pp.currX, pp.currY
}

func parsePathCommands(d string) ([]PathCommand, error) {
	r := []PathCommand{}

	pts, err := pathTokenize(d)
	if err != nil {
		return nil, err
	}

	t, pts := pts[0], pts[1:]
	var currentCommand byte = 'L'

	// If we have a blank input then just exit now
	if isEOD(t) {
		return r, nil
	}

	for !isEOD(t) {
		// See if this is a command token.  If so, then it becomes our current
		// command.
		ct, ok := t.(commandToken)
		if ok {
			currentCommand = ct.command
			t, pts = pts[0], pts[1:]
		}

		c := PathCommand{Command: currentCommand}

		commandLength := commandLengths[currentCommand]
		if len(pts) < commandLength {
			return nil, errors.Errorf("Path command (%c) doesn't have enough parameters (%d)", currentCommand, len(pts))
		}

		for i := 0; i < commandLength; i++ {
			n, ok := t.(numberToken)
			if !ok {
				return nil, errors.Errorf("Path command (%c) doesn't have enough parameters", currentCommand)
			}
			c.Params = append(c.Params, n.number)
			t, pts = pts[0], pts[1:]
		}

		r = append(r, c)

		// Any extra parameters after M are the corresponding L
		if currentCommand == 'm' {
			currentCommand = 'l'
		}
		if currentCommand == 'M' {
			currentCommand = 'L'
		}
	}

	return r, nil
}

// Path tokenizer

type pathToken interface{}
type commandToken struct {
	command byte
}
type numberToken struct {
	number float64
}
type eodToken struct{}

var floatRE *regexp.Regexp = regexp.MustCompile(`^[-+]?[0-9]*\.?[0-9]+(?:[eE][-+]?[0-9]+)?`)

func pathTokenize(d string) ([]pathToken, error) {
	var r []pathToken
	for len(d) > 0 {
		// Trim any whitespace and commas
		d = strings.TrimLeft(d, " \t\r\n,")
		if len(d) == 0 {
			break
		}

		// Look for a float
		loc := floatRE.FindStringIndex(d)
		if loc != nil {
			f, err := strconv.ParseFloat(d[loc[0]:loc[1]], 64)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			d = d[loc[1]:]
			r = append(r, numberToken{f})
		} else {
			// Assume it is a command
			c := d[0]
			if _, ok := commandLengths[c]; !ok {
				return nil, errors.Errorf("Unknown command in path: %c", c)
			}
			d = d[1:]
			r = append(r, commandToken{c})
		}
	}
	r = append(r, eodToken{})
	return r, nil
}
