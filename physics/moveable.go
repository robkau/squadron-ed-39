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
	pos  pixel.Vec
	vel  pixel.Vec
	dest pixel.Vec
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

func (point *LinearPointMovingStrategy) move(dt float64) (arrived bool) {
	// sitting at destination
	if point.pos == point.dest {
		return true
	}

	if point.vel != pixel.ZV {
		// move toward destination
		point.pos = point.pos.Add(point.vel.Scaled(dt))

		// arrived perfectly at destination
		if point.pos == point.dest {
			point.vel = pixel.ZV
			return true
		}

		// overshot destination
		if point.vel.X >= 0 {
			if point.pos.X > point.dest.X {
				point.pos = point.dest
				point.vel = pixel.ZV
				return true
			}
		} else {
			if point.pos.X < point.dest.X {
				point.pos = point.dest
				point.vel = pixel.ZV
				return true
			}
		}

		if point.vel.Y >= 0 {
			if point.pos.Y > point.dest.Y {
				point.pos = point.dest
				point.vel = pixel.ZV
				return true
			}
		} else {
			if point.pos.Y < point.dest.Y {
				point.pos = point.dest
				point.vel = pixel.ZV
				return true
			}
		}
		// still moving
		return false
	}
	return false
}

type LinearRectMovingStrategy struct {
}
