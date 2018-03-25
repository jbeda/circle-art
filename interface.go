package main

type GridContent interface {
	SetSize(w, h int)

	// The GridContent should return a value between 0 and 1 for this "pixel"
	GetValue(x, y int) float64
}
