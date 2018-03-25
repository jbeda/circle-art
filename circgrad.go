package main

import "github.com/jbeda/geom"

type CircularGradient struct {
	w, h    int
	center  geom.Coord
	maxDist float64
}

func (c *CircularGradient) SetSize(w, h int) {
	c.w, c.h = w, h
	c.center = geom.Coord{float64(c.w) / 2.0, float64(c.h) / 2.0}
	c.maxDist = geom.Coord{0, 0}.DistanceFrom(c.center)
}

func (c *CircularGradient) GetValue(x, y int) float64 {
	p := geom.Coord{float64(x), float64(y)}
	dist := p.DistanceFrom(c.center)
	return dist / c.maxDist
}
