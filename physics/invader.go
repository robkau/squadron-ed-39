package physics

import (
	"github.com/faiface/pixel"
	"image/color"
)

type Invader struct {
	Rect   pixel.Rect
	Color  color.Color
	Health int
}
