package main

import (
	"fmt"

	"math"

	"io/ioutil"

	"github.com/jbeda/geom"
	"github.com/jbeda/svgdata-go"
)

const (
	cSpace             = 0.15
	cMargin            = 0.01
	cMinRadius         = 0.01
	cMaxRadius         = (cSpace - cMargin) / 2.0
	canvasMargin       = 0.2
	canvasWidth        = 15.0
	canvasInsideWidth  = canvasWidth - canvasMargin*2.0
	canvasHeight       = 10.0
	canvasInsideHeight = canvasHeight - canvasMargin*2.0
	boardWidth         = 19.0
	boardHeight        = 11.0
	numChunks          = 4
)

// Scales a number between 0 and inMax proportionally to outMin and outMax
func ScaleNum(in, inMax, outMin, outMax float64) float64 {
	return outMin + (outMax-outMin)*(in/inMax)
}

func initColors(n int) []string {
	r := []string{}

	for i := 0; i < n; i++ {
		r = append(r, fmt.Sprintf("#00%02x00", int(ScaleNum(float64(i), float64(n), 32, 255))))
	}
	return r
}

func main() {
	xNum := int(math.Floor(canvasInsideWidth / cSpace))
	xSpace := canvasInsideWidth / float64(xNum)
	yNum := int(math.Floor(canvasInsideHeight / cSpace))
	ySpace := canvasInsideHeight / float64(yNum)

	xPerChunk := int(xNum / numChunks)

	canvasCenter := geom.Coord{canvasWidth / 2.0, canvasHeight / 2.0}
	maxDistance := geom.Coord{0, 0}.DistanceFrom(canvasCenter)

	colors := initColors(numChunks)

	for chunk := 0; chunk < numChunks; chunk++ {
		// Set up new root for chunk
		r := svgdata.CreateRoot()
		r.Attrs()["viewBox"] = fmt.Sprintf("0 0 %g %g", boardWidth, boardHeight)
		r.Attrs()["width"] = fmt.Sprintf("%gin", boardWidth)
		r.Attrs()["height"] = fmt.Sprintf("%gin", boardHeight)

		if chunk == numChunks-1 {
			outline := svgdata.NewRectXYWH(0, 0, canvasWidth, canvasHeight)
			r.AddChild(outline)
			outline.Attrs()["style"] = "fill: none; stroke: red; stroke-width: 0.001"
			outline.Attrs()["transform"] = fmt.Sprint("translate(%g, %g)", (boardWidth-canvasWidth)/2.0, (boardHeight-canvasHeight)/2.0)
		}

		g := svgdata.NewGroup()
		r.AddChild(g)
		g.Attrs()["style"] = fmt.Sprintf("fill: none; stroke: %s; stroke-width: 0.001", colors[chunk])
		g.Attrs()["transform"] = fmt.Sprintf("translate(%g, %g)", (boardWidth-canvasWidth)/2.0, (boardHeight-canvasHeight)/2.0)

		chunkStart := chunk * xPerChunk
		chunkEnd := (chunk + 1) * xPerChunk
		if chunk == numChunks-1 {
			chunkEnd = xNum
		}

		for x := chunkStart; x < chunkEnd; x++ {

			for y := 0; y < yNum; y++ {
				c := geom.Coord{
					canvasMargin + cSpace/2 + float64(x)*xSpace,
					canvasMargin + cSpace/2 + float64(y)*ySpace,
				}
				rRaw := c.DistanceFrom(canvasCenter)
				rad := ScaleNum(rRaw, maxDistance, cMinRadius, cMaxRadius)
				g.AddChild(svgdata.NewCircle(c, rad))
			}
		}

		d, _ := svgdata.Marshal(r, true)
		ioutil.WriteFile(fmt.Sprintf("output%03d.svg", chunk), d, 0644)
	}
}
