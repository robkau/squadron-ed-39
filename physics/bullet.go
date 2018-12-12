package physics

func EnforceMinBulletSpeed(b *Bullet) {
	if b.Vel.Len() < BulletMinSpeed {
		b.Vel = b.Vel.Scaled(BulletMinSpeed / b.Vel.Len())
	}
}

func deleteBullets(bullets *[]*Bullet) {
	// todo: debug
	j := 0
	for _, b := range *bullets {
		if !b.collided || (Abs(b.Pos.X) < MAX_BULLET_BOUND && Abs(b.Pos.Y) < MAX_BULLET_BOUND) {
			(*bullets)[j] = b
			j++
		}
	}
	*bullets = (*bullets)[:j]
}
