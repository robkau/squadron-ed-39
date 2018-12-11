package physics

func EnforceMinBulletSpeed(b *Bullet) {
	if b.Vel.Len() < BulletMinSpeed {
		b.Vel = b.Vel.Scaled(BulletMinSpeed / b.Vel.Len())
	}
}

func deleteBullet(i int, b []*Bullet) {
	copy(b[i:], b[i+1:])
	b[len(b)-1] = nil // or the zero value of T
	b = b[:len(b)-1]
}
