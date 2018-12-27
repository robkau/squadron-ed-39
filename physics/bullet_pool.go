package physics

import "github.com/faiface/pixel"

type BulletPool struct {
	pool chan *Bullet
}

func NewPool(max int) *BulletPool {
	return &BulletPool{
		pool: make(chan *Bullet, max),
	}
}

func (p *BulletPool) Get() *Bullet {
	var b *Bullet
	select {
	case b = <-p.pool:
	default:
		b = &Bullet{LinearPointMovingStrategy: LinearPointMovingStrategy{}}
	}

	return b
}

func (p *BulletPool) Put(b *Bullet) {
	b.collided = false
	b.SetDest(pixel.ZV)
	b.SetPos(pixel.ZV)
	b.SetVel(pixel.ZV)
	select {
	case p.pool <- b:
	default:
	}
}
