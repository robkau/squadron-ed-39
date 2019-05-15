package physics

import "github.com/faiface/pixel"

type BulletSpawner struct {
	moveable
	// todo: control spawn rate with struct var
}

func (world *world) spawnBullet(pos pixel.Vec, dest pixel.Vec) {
	if world.bulletCounter > maxBullets {
		return
	}
	v := pixel.Lerp(pos, dest, BulletSpeedFactor)
	v = v.Sub(pos) // rebase velocity calculation to origin
	b := world.BulletPool.get()
	b.setPos(pos)
	b.setDest(dest)
	b.setVel(v)
	enforceMinBulletSpeed(b)
	world.bullets = append(world.bullets, b)
	world.bulletCounter += 1
}

func (bsp *BulletSpawner) shoot(world *world) {
	world.spawnBullet(bsp.pos(), bsp.pickTarget(world))
}

func (bsp *BulletSpawner) pickTarget(world *world) pixel.Vec {
	if len(world.platforms) == 0 {
		if len(world.collectors) != 0 {
			return world.collectors[0].pos()
		}
		return pixel.ZV
	}

	// todo: make accurate, this can fail for large or small cases
	targetPlatform := world.lowestPlatform()
	timeToImpact := targetPlatform.pos().Sub(bsp.pos()).Len() / BulletMinSpeed
	firingTarget := targetPlatform.pos().Add(targetPlatform.vel().Scaled(timeToImpact))
	return firingTarget
}

func (world *world) bulletSpray(dest pixel.Vec) {
	for _, sh := range world.shooters {
		sh.bulletSpray(world, dest)
	}
}

func (bsp *BulletSpawner) bulletSpray(world *world, dest pixel.Vec) {
	// calculate measurements for firing arc
	firingLine := dest.Sub(bsp.pos())
	firingSpread := 1.0 / 6 //rad
	numProjectiles := 5
	firingSpreadIncrement := 2 * (firingSpread / (float64(numProjectiles)))

	// left arc
	for j := -firingSpread + firingSpreadIncrement/2; j < 0; j += firingSpreadIncrement {
		world.spawnBullet(bsp.pos(), bsp.pos().Add(firingLine.Rotated(-j)))
	}
	// right arc
	for j := firingSpread - firingSpreadIncrement/2; j > 0; j -= firingSpreadIncrement {
		world.spawnBullet(bsp.pos(), bsp.pos().Add(firingLine.Rotated(-j)))
	}
	// center
	if numProjectiles%2 != 0 {
		world.spawnBullet(bsp.pos(), bsp.pos().Add(firingLine))
	}
}

func (world *world) setShooterDestination(dest pixel.Vec) {
	for _, sh := range world.shooters {
		sh.setDest(dest)
		walkDir := dest.Sub(sh.pos()).Unit()
		sh.setVel(walkDir.Scaled(BulletSpawnerMoveSpeed))
	}
}
