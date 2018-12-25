package physics

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

const (
	Dt                             = 0.05 // global simulation timestep
	MaxBulletBound         float64 = 525
	MaxWindowBound                 = 500
	BulletMinSpeed         float64 = 30
	BulletPoolSize                 = 2000
	BulletSpawnModulo              = 25
	BulletSpawnerMoveSpeed         = 2.5
	BulletSpeedFactor      float64 = 0.06
	SlowdownFactor                 = 8
	FpsTarget                      = 60
	maxBullets                     = 1000
)

func (world *world) Update(dt float64, mp pixel.Vec) {
	// spawn new bullets
	if world.iteration%BulletSpawnModulo == 0 {
		world.SpawnBullet(mp)
	}

	// update world
	world.shooter.Walk()

	deadBullet := false
	// todo: allocate nothing inside loop
	for _, b := range world.bullets {
		if b == nil || b.collided {
			continue
		}

		if b.Vel.X == 0 && b.Vel.Y == 0 {
			panic(fmt.Sprintf("There is a stuck bullet at %f, %f", b.Pos.X, b.Pos.Y))
		}

		if Abs(b.Pos.X) > MaxBulletBound || Abs(b.Pos.Y) > MaxBulletBound {
			deadBullet = true
			continue
		}

		// collision detection
		// todo: quadtree instead of brute force
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
		world.bulletCounter -= deleteBullets(&world.bullets, world.BulletPool)
	}

	world.iteration += 1
}

func (world *world) Draw(imd *imdraw.IMDraw) {
	imd.Color = redColor()
	for _, b := range world.bullets {
		if b == nil || b.collided {
			continue
		}
		imd.Push(b.Pos)
	}
	imd.Circle(3, 0)

	for _, p := range world.platforms {
		imd.Color = p.Color
		imd.Push(p.Rect.Min, p.Rect.Max)
		imd.Rectangle(0)
	}

	imd.Color = pixel.RGB(0, 0.5, 1)
	imd.Push(world.shooter.Pos)
	imd.Circle(5, 0)
}
