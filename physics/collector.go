package physics

type collector struct {
	linearRectMovingStrategy
	uniqueId string
}

func (cl *collector) collide(b *bullet, world *world) {
	b.collided = true
	world.deadBullet = true
	world.energyCount += 1
}

func (cl *collector) id() string {
	return cl.uniqueId
}
