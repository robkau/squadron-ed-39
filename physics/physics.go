package physics

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"image/color"
	"time"
)

const (
	Dt                          = 0.05 // global simulation timestep
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
	platforms = append(platforms, &platform{health: 50, rect: pixel.Rect{Min: pixel.Vec{X: -300, Y: -500}, Max: pixel.Vec{X: 300, Y: -450}}, color: pixel.RGB(0.1, 0.5, 0.8)})
	return &world{
		bullets:   make([]*Bullet, 0),
		platforms: platforms,
	}
}

func (world *world) Update(dt float64, ctrl pixel.Vec) {
	for _, b := range world.bullets {
		if b == nil || b.collided {
			continue
		}

		// dead bullet wasting resources?
		if b.Vel.X == 0 && b.Vel.Y == 0 {
			panic(fmt.Sprintf("There is a dead bullet at %f, %f", b.Pos.X, b.Pos.Y))
		}

		deadBullet := false
		if Abs(b.Pos.X) > MAX_BULLET_BOUND || Abs(b.Pos.Y) > MAX_BULLET_BOUND {
			deadBullet = true
			continue
		}

		// collision detection
		deadPlatform := false
		for _, p := range world.platforms {
			if p.rect.Contains(b.Pos) && p.health > 0 { // do not hit already killed platforms before removal
				b.collided = true
				deadBullet = true
				p.health -= 1
				if p.health <= 0 {
					deadPlatform = true
				}
				break
			}
		}
		// clear the platforms list if any were killed
		if deadPlatform {
			deletePlatforms(&world.platforms)
		}

		// clear the bullets list if any collided or flew too far away
		if deadBullet {
			deleteBullets(&world.bullets)
		}

		// update bullet position if it's still in flight
		b.Pos = b.Pos.Add(b.Vel.Scaled(dt))
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
