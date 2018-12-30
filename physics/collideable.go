package physics

type collideable interface {
	Contains(*Bullet) bool
	Collide(*Bullet, *world)
	Id() string
}
