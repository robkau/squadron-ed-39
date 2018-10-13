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
		b.vel.Y = 25
		b.rect = b.rect.Moved(b.vel.Scaled(dt))
	}

	// check collisions against each platform
	/*
		if gp.vel.Y <= 0 {
			for _, p := range platforms {
				if gp.rect.Max.X <= p.rect.Min.X || gp.rect.Min.X >= p.rect.Max.X {
					continue
				}
				if gp.rect.Min.Y > p.rect.Max.Y || gp.rect.Min.Y < p.rect.Max.Y+gpc.vel.Y*dt {
					continue
				}
				gp.vel.Y = 0
				gp.rect = gp.rect.Moved(pixel.V(0, p.rect.Max.Y-gp.rect.Min.Y))
				gp.ground = true
			}
		}
	*/
}

//https://developer.mozilla.org/en-US/docs/Games/Techniques/2D_collision_detection
