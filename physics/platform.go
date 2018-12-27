package physics

import (
	"image/color"
)

type platform struct {
	moveable
	Color  color.Color
	Health int
}

func (pl *platform) move(dt float64) {

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
