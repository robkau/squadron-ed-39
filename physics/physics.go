package physics

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

const (
	Dt                             = 0.05 // global simulation timestep
	MAX_BULLET_BOUND       float64 = 1500
	BulletMinSpeed         float64 = 15
	BulletPoolSize                 = 2000
	BulletSpawnModulo              = 10
	BulletSpawnerMoveSpeed         = 2.5
	BulletSpeedFactor      float64 = 0.06
	SlowdownFactor                 = 8
	FpsTarget                      = 60
)

type world struct {
	shooter    *BulletSpawner
	bullets    []*Bullet
	platforms  []*platform
	BulletPool *BulletPool
}

func NewWorld() *world {
	platforms := make([]*platform, 0)
	platforms = append(platforms, &platform{Health: 50, Rect: pixel.Rect{Min: pixel.Vec{X: -300, Y: -500}, Max: pixel.Vec{X: 300, Y: -450}}, Color: pixel.RGB(0.1, 0.5, 0.8)})
	return &world{
		shooter:    &BulletSpawner{},
		bullets:    make([]*Bullet, 0),
		platforms:  platforms,
		BulletPool: NewPool(BulletPoolSize),
	}
}

func (world *world) Update(dt float64, mp pixel.Vec, iteration int) {
	// spawn new bullets
	if iteration%BulletSpawnModulo == 0 {
		v := pixel.Lerp(world.shooter.Pos, mp, BulletSpeedFactor)
		v = v.Sub(world.shooter.Pos) // rebase velocity calculation to origin
		b := world.BulletPool.Get()
		b.Pos = world.shooter.Pos
		b.Dest.X = mp.X
		b.Dest.Y = mp.Y
		b.Vel.X = v.X
		b.Vel.Y = v.Y
		EnforceMinBulletSpeed(b)
		world.SpawnBullet(b)
	}

	// update world
	world.shooter.Walk()

	deadBullet := false
	for _, b := range world.bullets {
		if b == nil || b.collided {
			continue
		}

		if b.Vel.X == 0 && b.Vel.Y == 0 {
			panic(fmt.Sprintf("There is a stuck bullet at %f, %f", b.Pos.X, b.Pos.Y))
		}

		if Abs(b.Pos.X) > MAX_BULLET_BOUND || Abs(b.Pos.Y) > MAX_BULLET_BOUND {
			deadBullet = true
			continue
		}

		// collision detection
		deadPlatform := false
		for _, p := range world.platforms {
			if p.Rect.Contains(b.Pos) && p.Health > 0 {
				b.collided = true
				deadBullet = true
				p.Health -= 1
				if p.Health <= 0 {
					deadPlatform = true
				}
				break
			}
		}
		// clear the platforms list if any were killed
		if deadPlatform {
			deletePlatforms(&world.platforms)
		}

		// update bullet position
		b.Pos = b.Pos.Add(b.Vel.Scaled(dt))
	}
	// clear the bullets list if any collided or flew too far away
	if deadBullet {
		deleteBullets(&world.bullets, world.BulletPool)
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
		imd.Color = p.Color
		imd.Push(p.Rect.Min, p.Rect.Max)
		imd.Rectangle(0)
	}
}
