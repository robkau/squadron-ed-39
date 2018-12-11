package physics

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"image/color"
	"time"
)

const (
	Dt                          = 0.05 // universal time step for simulation
	MAX_BULLET_BOUND    float64 = 1500
	BulletMinSpeed      float64 = 4
	BulletSpawnInterval         = time.Second / 10
	BulletSpeedFactor   float64 = 0.05
	SlowdownFactor              = 8
)

type world struct {
	bullets   []*Bullet
	platforms []*platform
}

type platform struct {
	rect   pixel.Rect
	color  color.Color
	health int
}

type Bullet struct {
	Pos      pixel.Vec
	Vel      pixel.Vec
	Dest     pixel.Vec
	collided bool
}

func NewWorld() *world {
	platforms := make([]*platform, 0)
	platforms = append(platforms, &platform{rect: pixel.Rect{Min: pixel.Vec{X: -300, Y: -500}, Max: pixel.Vec{X: 300, Y: -450}}, color: randomNiceColor()})
	return &world{
		bullets:   make([]*Bullet, 0),
		platforms: platforms,
	}
}

func (world *world) Update(dt float64, ctrl pixel.Vec) {
	for bi, b := range world.bullets {
		if b == nil || b.collided {
			continue
		}

		if b.Vel.X == 0 {

		}

		for _, p := range world.platforms {
			if p.rect.Contains(b.Pos) {
				b.collided = true
				deleteBullet(bi, world.bullets)
			} else {
				if Abs(b.Pos.X) > MAX_BULLET_BOUND || Abs(b.Pos.Y) > MAX_BULLET_BOUND {
					deleteBullet(bi, world.bullets)
					continue
				}
				b.Pos = b.Pos.Add(b.Vel.Scaled(dt))
			}
		}
	}
}

func (world *world) Draw(imd *imdraw.IMDraw) {
	imd.Color = redColor()
	for _, b := range world.bullets {
		if b == nil || b.collided {
			continue
		}
		imd.Push(b.Pos)
	}
	imd.Circle(1, 2)

	for _, p := range world.platforms {
		imd.Color = p.color
		imd.Push(p.rect.Min, p.rect.Max)
		imd.Rectangle(0)
	}
}
