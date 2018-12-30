package physics

type collector struct {
	LinearRectMovingStrategy
	UniqueId string
}

func (cl *collector) Collide(b *Bullet, world *world) {
	b.collided = true
	world.deadBullet = true
	world.energyCount += 1
}

func (cl *collector) Id() string {
	return cl.UniqueId
}
