package physics

func EnforceMinBulletSpeed(b *Bullet) {
	if b.Vel.Len() < BulletMinSpeed {
		b.Vel = b.Vel.Scaled(BulletMinSpeed / b.Vel.Len())
	}
}

func deleteBullets(bullets *[]*Bullet) {
	// todo: profile and record heap allocations
	// next: profile again with bullet pool
	for i, b := range *bullets {
		if b != nil && !b.collided || (Abs(b.Pos.X) < MAX_BULLET_BOUND && Abs(b.Pos.Y) < MAX_BULLET_BOUND) {
			(*bullets)[i] = (*bullets)[len(*bullets)-1]
			(*bullets)[len(*bullets)-1] = nil
			*bullets = (*bullets)[:len(*bullets)-1]
		}
	}
}
