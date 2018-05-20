package main

const (
	// The total width/height of each "cell"
	cSpace = 0.12
	// The min between circles
	cMargin = 0.025
	// The radius for "white" circles
	cMinRadius = 0.01
	// The radius of "black" circles
	cMaxRadius = (cSpace - cMargin) / 2.0

	// The margin between edge of art and first circle
	canvasMargin = 0.2
	// canvasHeight       = 10.5
	canvasHeight       = 5
	canvasInsideHeight = canvasHeight - canvasMargin*2.0
	canvasWidth        = canvasHeight * 3.0 / 2.0
	canvasInsideWidth  = canvasWidth - canvasMargin*2.0

	// The total size of the board
	boardWidth  = 19.0
	boardHeight = 11.0

	// The units to use when rendering SVG. It shouldn't matter but the GlowForge seems to care.
	unitsPerInch = 96
	strokeWidth  = 0.01
)
