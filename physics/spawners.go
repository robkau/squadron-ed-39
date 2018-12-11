package physics

func (world *world) SpawnBullet(b *Bullet) {
	world.bullets = append(world.bullets, b)
}
