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

func deleteBullet(i int, b []*bullet) {
	copy(b[i:], b[i+1:])
	b[len(b)-1] = nil // or the zero value of T
	b = b[:len(b)-1]
}

func Abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
