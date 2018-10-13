package main

import (
	"github.com/faiface/pixel"
)

const MAX_BULLET_BOUND float64 = 3000

type objects struct {
	bullets []*bullet
}

type bullet struct {
	rect     pixel.Rect
	vel      pixel.Vec
	dest     pixel.Vec
	collided bool
}

var bulletSlowdownFactor float64 = 8

func (objects *objects) update(dt float64, ctrl pixel.Vec, platforms *[]*platform) {
	for bi, b := range objects.bullets {
		if b == nil || b.collided {
			continue
		}

		b.vel.X = b.dest.X - b.rect.Center().X/bulletSlowdownFactor
		b.vel.Y = b.dest.Y - b.rect.Center().Y/bulletSlowdownFactor

		for _, p := range *platforms {
			if collide := checkBulletCollision(b, p, dt); collide {
				b.rect = b.rect.Moved(pixel.V(0, 0))
				deleteBullet(bi, objects.bullets)
			} else {
				if Abs(b.rect.Center().X) > MAX_BULLET_BOUND || Abs(b.rect.Center().Y) > MAX_BULLET_BOUND {
					deleteBullet(bi, objects.bullets)
					continue
				}
				b.rect = b.rect.Moved(b.vel.Scaled(dt))
			}
		}
	}
}

func checkBulletCollision(b *bullet, p *platform, dt float64) bool {
	if b.rect.Max.X <= p.rect.Min.X || b.rect.Min.X >= p.rect.Max.X {
		return false
	}
	if b.rect.Max.Y+b.vel.Y*dt < p.rect.Min.Y {
		return false
	}
	b.collided = true
	p.color = redColor()
	return true
}

//https://developer.mozilla.org/en-US/docs/Games/Techniques/2D_collision_detection
