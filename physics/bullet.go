package physics

type bullet struct {
	linearPointMovingStrategy
	collided bool
}

func enforceMinBulletSpeed(b *bullet) {
	b.setVel(b.vel().Scaled(BulletMinSpeed / b.vel().Len()))
	// all same speed for now
	/*if b.Vel().Len() < BulletMinSpeed {
		b.SetVel(b.Vel().Scaled(BulletMinSpeed / b.Vel().Len()))
	}
	*/
}

func deleteBullets(bullets *[]*bullet, bp *bulletPool) int {
	// todo: profile and record heap allocations
	// next: profile again with bullet pool
	removed := 0
	for i, b := range *bullets {
		if b != nil && (b.collided || b.outOfBounds()) {
			// move dead bullet to end of slice
			(*bullets)[i] = (*bullets)[len(*bullets)-1]
			// dereference dead bullet pointer
			bp.put(b)
			(*bullets)[len(*bullets)-1] = nil
			// shrink slice by 1
			*bullets = (*bullets)[:len(*bullets)-1]
			removed += 1
		}
	}
	return removed
}

func (b *bullet) outOfBounds() bool {
	return Abs(b.pos().X) > MaxBulletBound || Abs(b.pos().Y) > MaxBulletBound
}
