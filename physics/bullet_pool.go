package physics

import "github.com/faiface/pixel"

type bulletPool struct {
	// preallocate and reuse bullets, they are often destroyed and recreated
	pool chan *bullet
}

func newPool(max int) *bulletPool {
	return &bulletPool{
		pool: make(chan *bullet, max),
	}
}

func (p *bulletPool) get() *bullet {
	var b *bullet
	select {
	case b = <-p.pool:
	default:
		b = &bullet{linearPointMovingStrategy: linearPointMovingStrategy{}}
	}

	return b
}

func (p *bulletPool) put(b *bullet) {
	b.collided = false
	b.setDest(pixel.ZV)
	b.setPos(pixel.ZV)
	b.setVel(pixel.ZV)
	select {
	case p.pool <- b:
	default:
	}
}
