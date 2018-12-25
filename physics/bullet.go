package physics

import "github.com/faiface/pixel"

type Bullet struct {
	Pos      pixel.Vec
	Vel      pixel.Vec
	Dest     pixel.Vec
	collided bool
}

func EnforceMinBulletSpeed(b *Bullet) {
	if b.Vel.Len() < BulletMinSpeed {
		b.Vel = b.Vel.Scaled(BulletMinSpeed / b.Vel.Len())
	}
}

func deleteBullets(bullets *[]*Bullet, bp *BulletPool) int {
	// todo: profile and record heap allocations
	// next: profile again with bullet pool
	removed := 0
	for i, b := range *bullets {
		if b != nil && (b.collided || (Abs(b.Pos.X) > MaxBulletBound || Abs(b.Pos.Y) > MaxBulletBound)) {
			// move dead bullet to end of slice
			(*bullets)[i] = (*bullets)[len(*bullets)-1]
			// dereference dead bullet pointer
			bp.Put(b)
			(*bullets)[len(*bullets)-1] = nil
			// shrink slice by 1
			*bullets = (*bullets)[:len(*bullets)-1]
			removed += 1
		}
	}
	return removed
}
