package physics

import "github.com/faiface/pixel"

type moveable interface {
	move(float64) bool
	Pos() pixel.Vec
	SetPos(pixel.Vec)
	Vel() pixel.Vec
	SetVel(pixel.Vec)
	Dest() pixel.Vec
	SetDest(pixel.Vec)
}

type LinearPointMovingStrategy struct {
	pos        pixel.Vec
	vel        pixel.Vec
	dest       pixel.Vec
	stopAtDest bool
}

func (point *LinearPointMovingStrategy) Pos() pixel.Vec {
	return point.pos
}

func (point *LinearPointMovingStrategy) SetPos(pos pixel.Vec) {
	point.pos = pos
}

func (point *LinearPointMovingStrategy) Vel() pixel.Vec {
	return point.vel
}

func (point *LinearPointMovingStrategy) SetVel(vel pixel.Vec) {
	point.vel = vel
}

func (point *LinearPointMovingStrategy) Dest() pixel.Vec {
	return point.dest
}

func (point *LinearPointMovingStrategy) SetDest(dest pixel.Vec) {
	point.dest = dest
}

func (point *LinearPointMovingStrategy) move(dt float64) bool {
	arrivedOneAxis := false

	// sitting at destination
	if point.pos == point.dest {
		return true
	}

	if point.vel != pixel.ZV {
		// move toward destination
		point.pos = point.pos.Add(point.vel.Scaled(dt))

		if point.stopAtDest {
			// arrived perfectly at destination
			if point.pos == point.dest {
				point.vel = pixel.ZV
				return true
			}

			// overshot destination
			if point.vel.X >= 0 {
				if point.pos.X > point.dest.X {
					point.pos.X = point.dest.X
					point.vel.X = 0
					arrivedOneAxis = true
				}
			} else {
				if point.pos.X < point.dest.X {
					point.pos.X = point.dest.X
					point.vel.X = 0
					arrivedOneAxis = true
				}
			}

			if point.vel.Y >= 0 {
				if point.pos.Y > point.dest.Y {
					point.pos.Y = point.dest.Y
					point.vel.Y = 0
					if arrivedOneAxis {
						return true
					}
				}
			} else {
				if point.pos.Y < point.dest.Y {
					point.pos.Y = point.dest.Y
					point.vel.Y = 0
					if arrivedOneAxis {
						return true
					}
				}
			}
		}
		// still moving
		return false
	}
	return false
}

type LinearRectMovingStrategy struct {
	rect pixel.Rect
	vel  pixel.Vec
	dest pixel.Vec
}

func (rect *LinearRectMovingStrategy) Rect() pixel.Rect {
	return rect.rect
}

func (rect *LinearRectMovingStrategy) Pos() pixel.Vec {
	return rect.rect.Center()
}

func (rect *LinearRectMovingStrategy) SetPos(pos pixel.Vec) {
	dx := rect.rect.W() / 2
	dy := rect.rect.H() / 2
	rect.rect.Min.X = pos.X - dx
	rect.rect.Min.Y = pos.Y - dy
	rect.rect.Max.X = pos.X + dx
	rect.rect.Max.Y = pos.Y + dy
}

func (rect *LinearRectMovingStrategy) Contains(b *Bullet) bool {
	return rect.rect.Contains(b.Pos())
}

func (rect *LinearRectMovingStrategy) Vel() pixel.Vec {
	return rect.vel
}

func (rect *LinearRectMovingStrategy) SetVel(vel pixel.Vec) {
	rect.vel = vel
}

func (rect *LinearRectMovingStrategy) Dest() pixel.Vec {
	return rect.dest
}

func (rect *LinearRectMovingStrategy) SetDest(dest pixel.Vec) {
	rect.dest = dest
}

func (rect *LinearRectMovingStrategy) move(dt float64) (arrived bool) {
	arrivedOneAxis := false

	// sitting at destination
	if rect.rect.Center() == rect.dest {
		return true
	}

	if rect.vel != pixel.ZV {
		// move toward destination
		rect.SetPos(rect.rect.Center().Add(rect.vel.Scaled(dt)))

		// arrived perfectly at destination
		if rect.rect.Center() == rect.dest {
			rect.vel = pixel.ZV
			return true
		}

		// overshot destination
		if rect.vel.X >= 0 {
			if rect.rect.Center().X > rect.dest.X {
				rect.SetPos(pixel.Vec{X: rect.dest.X, Y: rect.rect.Center().Y})
				rect.vel.X = 0
				arrivedOneAxis = true
			}
		} else {
			if rect.rect.Center().X < rect.dest.X {
				rect.SetPos(pixel.Vec{X: rect.dest.X, Y: rect.rect.Center().Y})
				rect.vel.X = 0
				arrivedOneAxis = true
			}
		}

		if rect.vel.Y >= 0 {
			if rect.rect.Center().Y > rect.dest.Y {
				rect.SetPos(pixel.Vec{X: rect.rect.Center().X, Y: rect.dest.Y})
				rect.vel.Y = 0
				if arrivedOneAxis {
					return true
				}
			}
		} else {
			if rect.rect.Center().Y < rect.dest.Y {
				rect.SetPos(pixel.Vec{X: rect.rect.Center().X, Y: rect.dest.Y})
				rect.vel.Y = 0
				if arrivedOneAxis {
					return true
				}
			}
		}
		// still moving
		return false
	}
	return false
}
