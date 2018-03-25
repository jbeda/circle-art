package main

import (
	"fmt"

	"github.com/jbeda/geom"
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
