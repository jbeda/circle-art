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

	colors := initColors(4)

	b := bytes.Buffer{}
	b.WriteString(fmt.Sprintf(".border{fill:none;stroke:red;stroke-width:%g;}\n", scaleValue(strokeWidth)))
	for i := 0; i < 4; i++ {
		b.WriteString(fmt.Sprintf(".c%d{fill:none;stroke:%s;stroke-width:%g;}\n", i, colors[i], scaleValue(strokeWidth)))
	}

	style.SetText(b.String())
	return style
}

func (sg *SVGGrid) CreateRoot() *svgdata.Root {
	r := svgdata.CreateRoot()
	r.Attrs()["viewBox"] = fmt.Sprintf("0 0 %g %g", scaleValue(boardWidth), scaleValue(boardHeight))
	r.Attrs()["version"] = "1.1"
	//r.Attrs()["width"] = fmt.Sprintf("%gin", boardWidth)
	//r.Attrs()["height"] = fmt.Sprintf("%gin", boardHeight)
	r.Attrs()["x"] = "0px"
	r.Attrs()["y"] = "0px"
	r.Attrs()["style"] = fmt.Sprintf("enable-background:new %s;", r.Attrs()["viewBox"])

	r.AddChild(sg.createStyleElement())

	return r
}

func (sg *SVGGrid) RenderGrid(gc GridContent, outputPrefix string) {
	gc.SetSize(sg.xNum, sg.yNum)

	xOffset := (boardWidth - canvasWidth) / 2.0
	yOffset := (boardHeight - canvasHeight) / 2.0

	r := sg.CreateRoot()

	for xSkip := 0; xSkip < 2; xSkip++ {
		for ySkip := 0; ySkip < 2; ySkip++ {

			if xSkip == 1 && ySkip == 1 {
				outline := svgdata.NewRectXYWH(scaleValue(xOffset), scaleValue(yOffset), scaleValue(canvasWidth), scaleValue(canvasHeight))
				r.AddChild(outline)
				outline.Attrs()["class"] = "border"
			}

			g := svgdata.NewGroup()
			r.AddChild(g)

			for x := 0 + xSkip; x < sg.xNum; x += 2 {
				for y := 0 + ySkip; y < sg.yNum; y += 2 {
					c := scaleCoord(geom.Coord{
						xOffset + canvasMargin + cSpace/2 + float64(x)*sg.xSpace,
						yOffset + canvasMargin + cSpace/2 + float64(y)*sg.ySpace,
					})
					rRaw := gc.GetValue(x, y)
					rad := scaleValue(scaleToRange(rRaw, 1.0, cMinRadius, cMaxRadius))
					circle := svgdata.NewCircle(c, rad)
					circle.Attrs()["class"] = fmt.Sprintf("c%d", 2*xSkip+ySkip)
					g.AddChild(circle)
				}
			}
		}
	}

	// Write out the SVG file
	d, _ := svgdata.Marshal(r, true)
	ioutil.WriteFile(fmt.Sprintf("%s.svg", outputPrefix), d, 0644)
}
