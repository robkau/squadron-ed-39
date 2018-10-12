package main

import (
	"github.com/faiface/pixel"
	"math"
	"math/rand"
)

func randomNiceColor() pixel.RGBA {
again:
	r := rand.Float64()
	g := rand.Float64()
	b := rand.Float64()
	leng := math.Sqrt(r*r + g*g + b*b)
	if leng == 0 {
		goto again
	}
	return pixel.RGB(r/leng, g/leng, b/leng)
}

func redColor() pixel.RGBA {
	r := 1.0
	g := 0.0
	b := 0.0
	return pixel.RGB(r, g, b)
}
