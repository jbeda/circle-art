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
	unitsPerInch       = 72
	strokeWidth        = 0.01
)

func initColors(n int) []string {
	r := []string{}

	for i := 0; i < n; i++ {
		r = append(r, fmt.Sprintf("#00%02x00", int(scaleToRange(float64(i), float64(n), 32, 255))))
	}
	return r
}

// Scales a number between 0 and inMax proportionally to outMin and outMax
func scaleToRange(in, inMax, outMin, outMax float64) float64 {
	return outMin + (outMax-outMin)*(in/inMax)
}

func scaleValue(f float64) float64 {
	return f * unitsPerInch
}

func scaleCoord(c geom.Coord) geom.Coord {
	return geom.Coord{scaleValue(c.X), scaleValue(c.Y)}
}

func style(color string) string {
	return fmt.Sprintf("fill: none; stroke: %s; stroke-width: %g", color, scaleValue(strokeWidth))
}

func main() {
	xNum := int(math.Floor(canvasInsideWidth / cSpace))
	xSpace := canvasInsideWidth / float64(xNum)
	yNum := int(math.Floor(canvasInsideHeight / cSpace))
	ySpace := canvasInsideHeight / float64(yNum)

	xPerChunk := int(xNum / numChunks)

	canvasCenter := scaleCoord(geom.Coord{canvasWidth / 2.0, canvasHeight / 2.0})
	maxDistance := geom.Coord{0, 0}.DistanceFrom(canvasCenter)

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
			outline.Attrs()["transform"] = fmt.Sprint("translate(%g, %g)",
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
			chunkEnd = xNum
		}

		for x := chunkStart; x < chunkEnd; x++ {

			for y := 0; y < yNum; y++ {
				c := scaleCoord(geom.Coord{
					canvasMargin + cSpace/2 + float64(x)*xSpace,
					canvasMargin + cSpace/2 + float64(y)*ySpace,
				})
				rRaw := c.DistanceFrom(canvasCenter)
				rad := scaleValue(scaleToRange(rRaw, maxDistance, cMinRadius, cMaxRadius))
				g.AddChild(svgdata.NewCircle(c, rad))
			}
		}

		d, _ := svgdata.Marshal(r, true)
		ioutil.WriteFile(fmt.Sprintf("test-%03d.svg", chunk), d, 0644)
	}
}
