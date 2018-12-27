package physics

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

const (
	Dt                             = 0.1 // global simulation timestep
	MaxBulletBound         float64 = 525
	MaxWindowBound                 = 500
	BulletMinSpeed         float64 = 30
	BulletPoolSize                 = 2000
	BulletSpawnModulo              = 25
	BulletSpawnerMoveSpeed         = 25
	BulletSpeedFactor      float64 = 0.06
	SlowdownFactor                 = 8
	FpsTarget                      = 60
	maxBullets                     = 1000
)

func (world *world) Update(dt float64, mp pixel.Vec) {

	world.movePlatforms(dt)
	//world.checkPlatformCollisions()
	world.moveShooters(dt)
	//world.checkShooterCollisions()
	world.spawnBullets(mp)
	world.moveBullets(dt)
	world.checkBulletCollisions()

	world.iteration += 1
}

func (world *world) movePlatforms(dt float64) {
	for _, pl := range world.platforms {
		pl.move(dt)
	}
}

func (world *world) moveShooters(dt float64) {
	for _, sh := range world.shooters {
		sh.move(dt)
	}
}

func (world *world) spawnBullets(mp pixel.Vec) {
	if world.iteration%BulletSpawnModulo == 0 {
		for _, sh := range world.shooters {
			sh.shoot(world)
		}
	}
}

func (world *world) moveBullets(dt float64) {

	for _, b := range world.bullets {
		// update bullet position
		b.move(dt)
	}

}

func (world *world) checkBulletCollisions() {
	deadBullet := false
	for _, b := range world.bullets {
		if b == nil || b.collided {
			continue
		}

		if b.Vel().X == 0 && b.Vel().Y == 0 {
			b.collided = true
			deadBullet = true
			continue
		}

		if b.outOfBounds() {
			deadBullet = true
			continue
		}

		// collision detection
		// todo: quadtree instead of brute force
		deadPlatform := false
		for _, p := range world.platforms {
			if p.rect.Contains(b.Pos()) && p.Health > 0 {
				b.collided = true
				deadBullet = true
				p.Health -= 1
				if p.Health <= 0 {
					deadPlatform = true
				}
				break
			}
		}
		// clear dead platforms if needed
		if deadPlatform {
			deletePlatforms(&world.platforms)
		}
		// clear dead bullets if needed
		if deadBullet {
			world.bulletCounter -= deleteBullets(&world.bullets, world.BulletPool)
		}
	}
}

func (world *world) Draw(imd *imdraw.IMDraw) {
	imd.Color = redColor()
	for _, b := range world.bullets {
		if b == nil || b.collided {
			continue
		}
		imd.Push(b.Pos())
	}
	imd.Circle(3, 0)

	for _, p := range world.platforms {
		imd.Color = p.Color
		imd.Push(p.rect.Min, p.rect.Max)
		imd.Rectangle(0)
	}

	imd.Color = pixel.RGB(0, 0.5, 1)
	for _, sh := range world.shooters {
		imd.Push(sh.Pos())
	}
	imd.Circle(5, 0)
}
