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

type goal struct {
	pos    pixel.Vec
	radius float64
	step   float64

	counter float64
	cols    [2]pixel.RGBA
}

func draw(imd *imdraw.IMDraw, phys *objects) {
	imd.Color = redColor()
	for _, b := range phys.bullets {
		imd.Push(b.rect.Min, b.rect.Max)
		imd.Rectangle(0)
	}
}

func (g *goal) update(dt float64) {
	g.counter += dt
	for g.counter > g.step {
		g.counter -= g.step
		for i := len(g.cols) - 2; i >= 0; i-- {
			g.cols[i+1] = g.cols[i]
		}

	}
}

func (g *goal) draw(imd *imdraw.IMDraw) {
	for i := len(g.cols) - 1; i >= 0; i-- {

		imd.Color = g.cols[i]
		imd.Push(g.pos)
		imd.Circle(float64(i+1)*g.radius/float64(len(g.cols)), 0)
	}
}

func main() {
	pixelgl.Run(run)
}
