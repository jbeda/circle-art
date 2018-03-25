package main

import (
	"fmt"
	"io/ioutil"
	"math"

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

func (sg *SVGGrid) RenderGrid(gc GridContent) {
	gc.SetSize(sg.xNum, sg.yNum)

	xPerChunk := int(sg.xNum / numChunks)
	colors := initColors(numChunks)

	for chunk := 0; chunk < numChunks; chunk++ {
		// Set up new root for chunk
		r := svgdata.CreateRoot()
		r.Attrs()["viewBox"] = fmt.Sprintf("0 0 %g %g", scaleValue(boardWidth), scaleValue(boardHeight))
		r.Attrs()["width"] = fmt.Sprintf("%gin", boardWidth)
		r.Attrs()["height"] = fmt.Sprintf("%gin", boardHeight)

		if chunk == numChunks-1 {
			outline := svgdata.NewRectXYWH(0, 0, scaleValue(canvasWidth), scaleValue(canvasHeight))
			r.AddChild(outline)
			outline.Attrs()["style"] = style("red")
			outline.Attrs()["transform"] = fmt.Sprintf("translate(%g, %g)",
				scaleValue((boardWidth-canvasWidth)/2.0),
				scaleValue((boardHeight-canvasHeight)/2.0))
		}

		g := svgdata.NewGroup()
		r.AddChild(g)
		g.Attrs()["style"] = style(colors[chunk])
		g.Attrs()["transform"] = fmt.Sprintf("translate(%g, %g)",
			scaleValue((boardWidth-canvasWidth)/2.0),
			scaleValue((boardHeight-canvasHeight)/2.0))

		chunkStart := chunk * xPerChunk
		chunkEnd := (chunk + 1) * xPerChunk
		if chunk == numChunks-1 {
			chunkEnd = sg.xNum
		}

		for x := chunkStart; x < chunkEnd; x++ {
			for y := 0; y < sg.yNum; y++ {
				c := scaleCoord(geom.Coord{
					canvasMargin + cSpace/2 + float64(x)*sg.xSpace,
					canvasMargin + cSpace/2 + float64(y)*sg.ySpace,
				})
				rRaw := gc.GetValue(x, y)
				rad := scaleValue(scaleToRange(rRaw, 1.0, cMinRadius, cMaxRadius))
				g.AddChild(svgdata.NewCircle(c, rad))
			}
		}

		d, _ := svgdata.Marshal(r, true)
		ioutil.WriteFile(fmt.Sprintf("test-%03d.svg", chunk), d, 0644)
	}
}
