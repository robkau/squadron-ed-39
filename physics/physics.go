package physics

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"math/rand"
)

const (
	Dt                             = 0.1 // global simulation timestep
	MaxBulletBound         float64 = 525
	MaxWindowBound                 = 500
	BulletMinSpeed         float64 = 45
	PlatformSpeed          float64 = 5
	BulletPoolSize                 = 2000
	BulletSpawnModulo              = 25
	BulletSpawnerMoveSpeed         = 25
	BulletSpeedFactor      float64 = 0.06
	SlowdownFactor                 = 8
	FpsTarget                      = 60
	maxBullets                     = 2000
)

func (world *world) Update(dt float64, mp pixel.Vec) {
	world.movePlatforms(dt)
	world.moveShooters(dt)
	world.spawnBullets(mp)
	world.moveBullets(dt)
	world.checkBulletCollisions()

	// spawn progressively harder enemies
	if world.iteration%150 == 0 && world.iteration >= 650 {
		xPos := rand.Float64()*MaxWindowBound*1.5 - (MaxWindowBound*1.5)/2
		world.AddPlatform(pixel.Rect{Min: pixel.Vec{X: xPos, Y: 500}, Max: pixel.Vec{X: xPos + 50, Y: 525}}, pixel.Vec{X: xPos, Y: -1000}, 20)
	}
	if world.iteration%80 == 0 && world.iteration >= 2000 {
		xPos := rand.Float64()*MaxWindowBound*1.5 - (MaxWindowBound*1.5)/2
		world.AddPlatformWithV(pixel.Rect{Min: pixel.Vec{X: xPos, Y: 500}, Max: pixel.Vec{X: xPos + 50, Y: 525}}, pixel.Vec{X: xPos, Y: -1000}, pixel.Vec{X: 0, Y: -20}, 5)
	}
	if world.iteration%4 == 0 && world.iteration >= 3500 {
		xPos := rand.Float64()*MaxWindowBound*1.5 - (MaxWindowBound*1.5)/2
		world.AddPlatformWithV(pixel.Rect{Min: pixel.Vec{X: xPos, Y: 500}, Max: pixel.Vec{X: xPos + 50, Y: 525}}, pixel.Vec{X: xPos, Y: -1000}, pixel.Vec{X: 0, Y: -20}, 5)
	}

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
	// check each bullet for collisions
	for _, b := range world.bullets {
		if b == nil || b.collided {
			continue
		}

		if b.Vel().X == 0 && b.Vel().Y == 0 {
			b.collided = true
			world.deadBullet = true
			continue
		}

		if b.outOfBounds() {
			world.deadBullet = true
			continue
		}

		// collision detection
		// todo: quadtree instead of brute force
		for _, p := range world.colliders {
			if p.Contains(b) {
				p.Collide(b, world)
				break
			}

		}

		// clear dead platforms if needed
		if world.deadPlatform {
			deletePlatforms(&world.platforms, &world.colliders)
			world.deadPlatform = false
		}
		// clear dead bullets if needed
		if world.deadBullet {
			world.bulletCounter -= deleteBullets(&world.bullets, world.BulletPool)
			world.deadBullet = false
		}
	}
}

func (world *world) Draw(imd *imdraw.IMDraw) {
	// translate game world into drawn shapes
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

	imd.Color = collectorColor()
	for _, cl := range world.collectors {
		imd.Push(cl.Rect().Min, cl.Rect().Max)
		imd.Rectangle(3)
	}
}
