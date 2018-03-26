package main

import (
	"image"

	"image/color"

	"github.com/disintegration/imaging"
)

type ImageContent struct {
	w, h       int
	src, small image.Image
}

func NewImageContent(fn string) (*ImageContent, error) {
	src, err := imaging.Open(fn)
	if err != nil {
		return nil, err
	}

	src = imaging.Grayscale(src)

	if src.Bounds().Dx() < src.Bounds().Dy() {
		src = imaging.Rotate90(src)
	}

	return &ImageContent{src: src}, nil
}

func (ic *ImageContent) SetSize(w, h int) {
	ic.w, ic.h = w, h
	ic.small = imaging.Invert(imaging.Fill(ic.src, ic.w, ic.h, imaging.Center, imaging.Lanczos))
}

func (ic *ImageContent) GetValue(x, y int) float64 {
	c := color.GrayModel.Convert(ic.small.At(x, y)).(color.Gray)
	return float64(c.Y) / 255.0
}
