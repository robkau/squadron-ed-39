package physics

import "github.com/faiface/pixel"

// todo: interface for Moveable()
type BulletSpawner struct {
	moveable
	// todo: control spawn rate with var here
}

func (world *world) SpawnBullet(dest pixel.Vec) {
	if world.bulletCounter > maxBullets {
		return
	}
	v := pixel.Lerp(world.shooter.Pos(), dest, BulletSpeedFactor)
	v = v.Sub(world.shooter.Pos()) // rebase velocity calculation to origin
	b := world.BulletPool.Get()
	b.SetPos(world.shooter.Pos())
	b.SetDest(dest)
	b.SetVel(v)
	EnforceMinBulletSpeed(b)
	// apply ship momentum to launched bullets
	//b.Vel = b.Vel.Add(world.shooter.Vel.Scaled(10))
	world.bullets = append(world.bullets, b)
	world.bulletCounter += 1
}

func (world *world) BulletSpray(dest pixel.Vec) {
	// todo: implement for odd number of bullets
	firingLine := dest.Sub(world.shooter.Pos())
	firingSpread := 1.0 / 6 //rad
	numProjectiles := 12
	firingSpreadIncrement := 2 * (firingSpread / (float64(numProjectiles)))

	// left arc
	for j := -firingSpread + firingSpreadIncrement/2; j < 0; j += firingSpreadIncrement {
		world.SpawnBullet(world.shooter.Pos().Add(firingLine.Rotated(-j)))
	}
	// right arc
	for j := firingSpread - firingSpreadIncrement/2; j > 0; j -= firingSpreadIncrement {
		world.SpawnBullet(world.shooter.Pos().Add(firingLine.Rotated(-j)))
	}
	// center
	if numProjectiles%2 != 0 {
		world.SpawnBullet(world.shooter.Pos().Add(firingLine))
	}
}

func (world *world) SetShooterDestination(dest pixel.Vec) {
	world.shooter.SetDest(dest)
	walkDir := dest.Sub(world.shooter.Pos()).Unit()
	world.shooter.SetVel(walkDir.Scaled(BulletSpawnerMoveSpeed))
}
