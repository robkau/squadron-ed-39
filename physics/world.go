package physics

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

type world struct {
	shooters      []*BulletSpawner
	bullets       []*Bullet
	platforms     []*platform
	BulletPool    *BulletPool
	iteration     int
	bulletCounter int
	atlas         *text.Atlas
}

func NewWorld() *world {
	platforms := make([]*platform, 0)

	// todo: append platforms elsewhere
	platforms = append(platforms, &platform{LinearRectMovingStrategy: LinearRectMovingStrategy{rect: pixel.Rect{Min: pixel.Vec{X: -300, Y: -500}, Max: pixel.Vec{X: 300, Y: -450}}, dest: pixel.Vec{X: 50, Y: 300}, vel: pixel.Vec{X: 3, Y: 10}}, Health: 50, Color: pixel.RGB(0.1, 0.5, 0.8)})

	sh := make([]*BulletSpawner, 0)
	sh = append(sh, &BulletSpawner{moveable: &LinearPointMovingStrategy{stopAtDest: true}})
	sh = append(sh, &BulletSpawner{moveable: &LinearPointMovingStrategy{stopAtDest: false, pos: pixel.Vec{X: -100, Y: 0}}})

	return &world{
		shooters:   sh,
		bullets:    make([]*Bullet, 0),
		platforms:  platforms,
		BulletPool: NewPool(BulletPoolSize),
		atlas:      text.NewAtlas(basicfont.Face7x13, text.ASCII),
	}
}

func (world *world) NumBullets() int {
	return world.bulletCounter
}
