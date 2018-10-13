package main

import "github.com/faiface/pixel"

type objects struct {
	bullets []*bullet
}

type bullet struct {
	rect pixel.Rect
	vel  pixel.Vec
}

func (objects *objects) update(dt float64, ctrl pixel.Vec, platforms []platform) {
	for _, b := range objects.bullets {
		b.vel.X = 0
		b.vel.Y = 250
		for _, p := range platforms {
			if b.rect.Max.X <= p.rect.Min.X || b.rect.Min.X >= p.rect.Max.X {
				b.rect = b.rect.Moved(b.vel.Scaled(dt))
				continue
			}
			if b.rect.Max.Y+b.vel.Y*dt < p.rect.Min.Y {
				b.rect = b.rect.Moved(b.vel.Scaled(dt))
				continue
			}
			b.vel.Y = 0
			b.rect = b.rect.Moved(pixel.V(0, 0))
		}
	}

	// check collisions against each platform

}

//https://developer.mozilla.org/en-US/docs/Games/Techniques/2D_collision_detection
