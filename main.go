package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
)

type platform struct {
	rect  pixel.Rect
	color color.Color
}

func (p *platform) draw(imd *imdraw.IMDraw) {
	imd.Color = p.color
	imd.Push(p.rect.Min, p.rect.Max)
	imd.Rectangle(0)
}

func draw(imd *imdraw.IMDraw, phys *objects) {
	imd.Color = redColor()
	for _, b := range phys.bullets {
		imd.Push(b.rect.Min, b.rect.Max)
		imd.Rectangle(0)
	}
}

func main() {
	pixelgl.Run(run)
}
