package physics

import (
	"github.com/faiface/pixel"
	"image/color"
)

type platform struct {
	Rect   pixel.Rect
	Color  color.Color
	Health int
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
