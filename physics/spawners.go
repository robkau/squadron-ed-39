package physics

import "github.com/faiface/pixel"

type BulletSpawner struct {
	moveable
	// todo: control spawn rate with struct var
}

func (world *world) SpawnBullet(pos pixel.Vec, dest pixel.Vec) {
	if world.bulletCounter > maxBullets {
		return
	}
	v := pixel.Lerp(pos, dest, BulletSpeedFactor)
	v = v.Sub(pos) // rebase velocity calculation to origin
	b := world.BulletPool.Get()
	b.SetPos(pos)
	b.SetDest(dest)
	b.SetVel(v)
	EnforceMinBulletSpeed(b)
	world.bullets = append(world.bullets, b)
	world.bulletCounter += 1
}

func (bsp *BulletSpawner) shoot(world *world) {
	world.SpawnBullet(bsp.Pos(), bsp.pickTarget(world))
}

func (bsp *BulletSpawner) pickTarget(world *world) pixel.Vec {
	if len(world.platforms) == 0 {
		if len(world.collectors) != 0 {
			return world.collectors[0].Pos()
		}
		return pixel.ZV
	}

	// todo: make accurate, this can fail for large or small cases
	targetPlatform := world.lowestPlatform()
	timeToImpact := targetPlatform.Pos().Sub(bsp.Pos()).Len() / BulletMinSpeed
	firingTarget := targetPlatform.Pos().Add(targetPlatform.Vel().Scaled(timeToImpact))
	return firingTarget
}

func (world *world) BulletSpray(dest pixel.Vec) {
	for _, sh := range world.shooters {
		sh.BulletSpray(world, dest)
	}
}

func (bsp *BulletSpawner) BulletSpray(world *world, dest pixel.Vec) {
	// calculate measurements for firing arc
	firingLine := dest.Sub(bsp.Pos())
	firingSpread := 1.0 / 6 //rad
	numProjectiles := 5
	firingSpreadIncrement := 2 * (firingSpread / (float64(numProjectiles)))

	// left arc
	for j := -firingSpread + firingSpreadIncrement/2; j < 0; j += firingSpreadIncrement {
		world.SpawnBullet(bsp.Pos(), bsp.Pos().Add(firingLine.Rotated(-j)))
	}
	// right arc
	for j := firingSpread - firingSpreadIncrement/2; j > 0; j -= firingSpreadIncrement {
		world.SpawnBullet(bsp.Pos(), bsp.Pos().Add(firingLine.Rotated(-j)))
	}
	// center
	if numProjectiles%2 != 0 {
		world.SpawnBullet(bsp.Pos(), bsp.Pos().Add(firingLine))
	}
}

func (world *world) SetShooterDestination(dest pixel.Vec) {
	for _, sh := range world.shooters {
		sh.SetDest(dest)
		walkDir := dest.Sub(sh.Pos()).Unit()
		sh.SetVel(walkDir.Scaled(BulletSpawnerMoveSpeed))
	}
}
