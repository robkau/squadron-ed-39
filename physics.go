package main

import (
	"github.com/faiface/pixel"
)

const MAX_BULLET_BOUND float64 = 3000

type objects struct {
	bullets []*bullet
}

type bullet struct {
	pos      pixel.Vec
	vel      pixel.Vec
	dest     pixel.Vec
	collided bool
}

func (objects *objects) update(dt float64, ctrl pixel.Vec, platforms *[]*platform) {
	for bi, b := range objects.bullets {
		if b == nil || b.collided {
			continue
		}

		if b.vel.X == 0 {

		}

		for _, p := range *platforms {
			if p.rect.Contains(b.pos) {
				b.collided = true
				deleteBullet(bi, objects.bullets)
			} else {
				if Abs(b.pos.X) > MAX_BULLET_BOUND || Abs(b.pos.Y) > MAX_BULLET_BOUND {
					deleteBullet(bi, objects.bullets)
					continue
				}
				b.pos = b.pos.Add(b.vel.Scaled(dt))
			}
		}
	}
}
