package physics

type BulletPool struct {
	pool chan *Bullet
}

func NewPool(max int) *BulletPool {
	return &BulletPool{
		pool: make(chan *Bullet, max),
	}
}

// borrow from the pool.
func (p *BulletPool) Borrow() *Bullet {
	var b *Bullet
	select {
	case b = <-p.pool:
	default:
		b = &Bullet{}
	}

	return b
}

// returns to the pool.
func (p *BulletPool) Return(b *Bullet) {
	b.collided = false
	b.Dest.X = 0
	b.Dest.Y = 0
	b.Pos.X = 0
	b.Pos.Y = 0
	b.Vel.X = 0
	b.Vel.Y = 0
	select {
	case p.pool <- b:
	default:
	}
}
