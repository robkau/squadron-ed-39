package physics

import (
	"image/color"
)

type platform struct {
	LinearRectMovingStrategy
	Color    color.Color
	Health   int
	UniqueId string
}

func (pl *platform) Collide(b *Bullet, world *world) {
	b.collided = true
	world.deadBullet = true
	pl.Health -= 1
	if pl.Health <= 0 {
		world.deadPlatform = true
	}
}

func (pl *platform) Id() string {
	return pl.UniqueId

}

func deletePlatforms(platforms *[]*platform, collideables *[]collideable) {
	for i, p := range *platforms {
		if p != nil && p.Health <= 0 {
			pid := p.Id()
			// remove from world platforms slice
			(*platforms)[i] = (*platforms)[len(*platforms)-1]
			// dereference dead platform
			(*platforms)[len(*platforms)-1] = nil
			// shrink slice by 1
			*platforms = (*platforms)[:len(*platforms)-1]

			// remove from world collideables slice
			for j, c := range *collideables {
				if c != nil && pid == c.Id() {
					(*collideables)[j] = (*collideables)[len(*collideables)-1]
					// dereference dead platform
					(*collideables)[len(*collideables)-1] = nil
					// shrink slice by 1
					*collideables = (*collideables)[:len(*collideables)-1]
				}
			}
		}
	}
}
