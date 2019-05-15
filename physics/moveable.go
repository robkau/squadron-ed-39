package physics

import "github.com/faiface/pixel"

type moveable interface {
	move(float64) bool
	pos() pixel.Vec
	setPos(pixel.Vec)
	vel() pixel.Vec
	setVel(pixel.Vec)
	dest() pixel.Vec
	setDest(pixel.Vec)
}

type linearPointMovingStrategy struct {
	p          pixel.Vec
	v          pixel.Vec
	dst        pixel.Vec
	stopAtDest bool
}

func (point *linearPointMovingStrategy) pos() pixel.Vec {
	return point.p
}

func (point *linearPointMovingStrategy) setPos(pos pixel.Vec) {
	point.p = pos
}

func (point *linearPointMovingStrategy) vel() pixel.Vec {
	return point.v
}

func (point *linearPointMovingStrategy) setVel(vel pixel.Vec) {
	point.v = vel
}

func (point *linearPointMovingStrategy) dest() pixel.Vec {
	return point.dst
}

func (point *linearPointMovingStrategy) setDest(dest pixel.Vec) {
	point.dst = dest
}

func (point *linearPointMovingStrategy) move(dt float64) bool {
	arrivedOneAxis := false

	// sitting at destination
	if point.p == point.dst {
		return true
	}

	if point.v != pixel.ZV {
		// move toward destination
		point.p = point.p.Add(point.v.Scaled(dt))

		if point.stopAtDest {
			// arrived perfectly at destination
			if point.p == point.dst {
				point.v = pixel.ZV
				return true
			}

			// overshot destination
			if point.v.X >= 0 {
				if point.p.X > point.dst.X {
					point.p.X = point.dst.X
					point.v.X = 0
					arrivedOneAxis = true
				}
			} else {
				if point.p.X < point.dst.X {
					point.p.X = point.dst.X
					point.v.X = 0
					arrivedOneAxis = true
				}
			}

			if point.v.Y >= 0 {
				if point.p.Y > point.dst.Y {
					point.p.Y = point.dst.Y
					point.v.Y = 0
					if arrivedOneAxis {
						return true
					}
				}
			} else {
				if point.p.Y < point.dst.Y {
					point.p.Y = point.dst.Y
					point.v.Y = 0
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

type linearRectMovingStrategy struct {
	r   pixel.Rect
	v   pixel.Vec
	dst pixel.Vec
}

func (rect *linearRectMovingStrategy) rect() pixel.Rect {
	return rect.r
}

func (rect *linearRectMovingStrategy) pos() pixel.Vec {
	return rect.r.Center()
}

func (rect *linearRectMovingStrategy) setPos(pos pixel.Vec) {
	dx := rect.r.W() / 2
	dy := rect.r.H() / 2
	rect.r.Min.X = pos.X - dx
	rect.r.Min.Y = pos.Y - dy
	rect.r.Max.X = pos.X + dx
	rect.r.Max.Y = pos.Y + dy
}

func (rect *linearRectMovingStrategy) contains(b *bullet) bool {
	return rect.r.Contains(b.pos())
}

func (rect *linearRectMovingStrategy) vel() pixel.Vec {
	return rect.v
}

func (rect *linearRectMovingStrategy) setVel(vel pixel.Vec) {
	rect.v = vel
}

func (rect *linearRectMovingStrategy) dest() pixel.Vec {
	return rect.dst
}

func (rect *linearRectMovingStrategy) setDest(dest pixel.Vec) {
	rect.dst = dest
}

func (rect *linearRectMovingStrategy) move(dt float64) (arrived bool) {
	arrivedOneAxis := false

	// sitting at destination
	if rect.r.Center() == rect.dst {
		return true
	}

	if rect.v != pixel.ZV {
		// move toward destination
		rect.setPos(rect.r.Center().Add(rect.v.Scaled(dt)))

		// arrived perfectly at destination
		if rect.r.Center() == rect.dst {
			rect.v = pixel.ZV
			return true
		}

		// overshot destination
		if rect.v.X >= 0 {
			if rect.r.Center().X > rect.dst.X {
				rect.setPos(pixel.Vec{X: rect.dst.X, Y: rect.r.Center().Y})
				rect.v.X = 0
				arrivedOneAxis = true
			}
		} else {
			if rect.r.Center().X < rect.dst.X {
				rect.setPos(pixel.Vec{X: rect.dst.X, Y: rect.r.Center().Y})
				rect.v.X = 0
				arrivedOneAxis = true
			}
		}

		if rect.v.Y >= 0 {
			if rect.r.Center().Y > rect.dst.Y {
				rect.setPos(pixel.Vec{X: rect.r.Center().X, Y: rect.dst.Y})
				rect.v.Y = 0
				if arrivedOneAxis {
					return true
				}
			}
		} else {
			if rect.r.Center().Y < rect.dst.Y {
				rect.setPos(pixel.Vec{X: rect.r.Center().X, Y: rect.dst.Y})
				rect.v.Y = 0
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
