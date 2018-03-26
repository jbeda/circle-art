package main

import (
	"fmt"
	"io/ioutil"
	"math"

	"bytes"

	"github.com/jbeda/geom"
	svgdata "github.com/jbeda/svgdata-go"
)

type SVGGrid struct {
	xNum, yNum     int
	xSpace, ySpace float64
}

func NewSVGGrid() *SVGGrid {
	sg := &SVGGrid{}

	sg.xNum = int(math.Floor(canvasInsideWidth / cSpace))
	sg.xSpace = canvasInsideWidth / float64(sg.xNum)
	sg.yNum = int(math.Floor(canvasInsideHeight / cSpace))
	sg.ySpace = canvasInsideHeight / float64(sg.yNum)

	return sg
}

func (sg *SVGGrid) createStyleElement() svgdata.Node {
	style := svgdata.NewStyle()
	style.Attrs()["type"] = "text/css"

	colors := initColors(numChunks)

	b := bytes.Buffer{}
	b.WriteString(fmt.Sprintf(".border{fill:none;stroke:red;stroke-width:%g;}\n", scaleValue(strokeWidth)))
	for i := 0; i < numChunks; i++ {
		b.WriteString(fmt.Sprintf(".c%d{fill:none;stroke:%s;stroke-width:%g;}\n", i, colors[i], scaleValue(strokeWidth)))
	}

	style.SetText(b.String())
	return style
}

func (sg *SVGGrid) RenderGrid(gc GridContent) {
	gc.SetSize(sg.xNum, sg.yNum)

	xPerChunk := int(sg.xNum / numChunks)

	xOffset := (boardWidth - canvasWidth) / 2.0
	yOffset := (boardHeight - canvasHeight) / 2.0

	for chunk := 0; chunk < numChunks; chunk++ {
		// Set up new root for chunk
		r := svgdata.CreateRoot()
		r.Attrs()["viewBox"] = fmt.Sprintf("0 0 %g %g", scaleValue(boardWidth), scaleValue(boardHeight))
		r.Attrs()["version"] = "1.1"
		//r.Attrs()["width"] = fmt.Sprintf("%gin", boardWidth)
		//r.Attrs()["height"] = fmt.Sprintf("%gin", boardHeight)
		r.Attrs()["x"] = "0px"
		r.Attrs()["y"] = "0px"
		r.Attrs()["style"] = fmt.Sprintf("enable-background:new %s;", r.Attrs()["viewBox"])

		r.AddChild(sg.createStyleElement())

		if chunk == numChunks-1 {
			outline := svgdata.NewRectXYWH(xOffset, yOffset, scaleValue(canvasWidth), scaleValue(canvasHeight))
			r.AddChild(outline)
			outline.Attrs()["class"] = "border"
		}

		g := svgdata.NewGroup()
		r.AddChild(g)

		chunkStart := chunk * xPerChunk
		chunkEnd := (chunk + 1) * xPerChunk
		if chunk == numChunks-1 {
			chunkEnd = sg.xNum
		}

		for x := chunkStart; x < chunkEnd; x++ {
			for y := 0; y < sg.yNum; y++ {
				c := scaleCoord(geom.Coord{
					xOffset + canvasMargin + cSpace/2 + float64(x)*sg.xSpace,
					yOffset + canvasMargin + cSpace/2 + float64(y)*sg.ySpace,
				})
				rRaw := gc.GetValue(x, y)
				rad := scaleValue(scaleToRange(rRaw, 1.0, cMinRadius, cMaxRadius))
				circle := svgdata.NewCircle(c, rad)
				circle.Attrs()["class"] = fmt.Sprintf("c%d", chunk)
				g.AddChild(circle)
			}
		}

		d, _ := svgdata.Marshal(r, true)
		ioutil.WriteFile(fmt.Sprintf("test-%03d.svg", chunk), d, 0644)
	}
}
