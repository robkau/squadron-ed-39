package physics

import (
	"github.com/faiface/pixel"
	"image/color"
)

type platform struct {
	LinearRectMovingStrategy
	Color  color.Color
	Health int
}

func (pl *platform) Rect() *pixel.Rect {
	return pl.Rect()

}

func deletePlatforms(platforms *[]*platform) {
	for i, p := range *platforms {
		if p != nil && p.Health <= 0 {
			(*platforms)[i] = (*platforms)[len(*platforms)-1]
			// dereference dead platform pointer
			(*platforms)[len(*platforms)-1] = nil
			// shrink slice by 1
			*platforms = (*platforms)[:len(*platforms)-1]
		}
	}
}
