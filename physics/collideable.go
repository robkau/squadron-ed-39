package physics

type collideable interface {
	contains(*bullet) bool
	collide(*bullet, *world)
	id() string
}
