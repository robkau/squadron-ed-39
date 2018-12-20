package physics

type BulletPool struct {
	pool chan *Bullet
}

// NewPool creates a new pool of Clients.
func NewPool(max int) *BulletPool {
	return &BulletPool{
		pool: make(chan *Bullet, max),
	}
}

// Borrow a Client from the pool.
func (p *BulletPool) Borrow() *Bullet {
	var b *Bullet
	select {
	case b = <-p.pool:
	default:
		b = &Bullet{}
	}

	return b
}

// Return returns a Client to the pool.
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
