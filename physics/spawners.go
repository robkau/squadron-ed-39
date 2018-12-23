package physics

import "github.com/faiface/pixel"

// todo: interface for Moveable()
type BulletSpawner struct {
	Pos  pixel.Vec
	Dest pixel.Vec
	Vel  pixel.Vec
	// todo: control spawn rate with var here
}

func (world *world) SpawnBullet(dest pixel.Vec) {
	v := pixel.Lerp(world.shooter.Pos, dest, BulletSpeedFactor)
	v = v.Sub(world.shooter.Pos) // rebase velocity calculation to origin
	b := world.BulletPool.Get()
	b.Pos = world.shooter.Pos
	b.Dest.X = dest.X
	b.Dest.Y = dest.Y
	b.Vel.X = v.X
	b.Vel.Y = v.Y
	EnforceMinBulletSpeed(b)
	// apply ship momentum to launched bullets
	//b.Vel = b.Vel.Add(world.shooter.Vel.Scaled(10))
	world.bullets = append(world.bullets, b)
}

func (world *world) BulletSpray(dest pixel.Vec) {
	// todo: implement for odd number of bullets
	firingLine := dest.Sub(world.shooter.Pos)
	firingSpread := 1.0 / 6 //rad
	numProjectiles := 6
	firingSpreadIncrement := 2 * (firingSpread / (float64(numProjectiles)))

	for n := 0; n < numProjectiles; n++ {
		// n even
		if numProjectiles%2 == 0 {
			// left arc
			for j := -firingSpread + firingSpreadIncrement/2; j < 0; j += firingSpreadIncrement {
				world.SpawnBullet(world.shooter.Pos.Add(firingLine.Rotated(-j)))
			}
			// right arc
			for j := firingSpread - firingSpreadIncrement/2; j > 0; j -= firingSpreadIncrement {
				world.SpawnBullet(world.shooter.Pos.Add(firingLine.Rotated(-j)))
			}
			return
		}
		// n odd
		// left arc
		// center
		// right arc
	}
}

func (world *world) MoveShooter(dest pixel.Vec) {
	world.shooter.Dest = dest
	walkDir := dest.Sub(world.shooter.Pos).Unit()
	world.shooter.Vel = walkDir.Scaled(BulletSpawnerMoveSpeed)
}

func (bsp *BulletSpawner) Walk() {
	if bsp.Vel != pixel.ZV {
		// walk forward
		bsp.Pos = bsp.Pos.Add(bsp.Vel)

		// arrived perfectly at destination
		if bsp.Pos == bsp.Dest {
			bsp.Vel = pixel.ZV
			return
		}

		// overshot destination
		if bsp.Vel.X >= 0 {
			if bsp.Pos.X > bsp.Dest.X {
				bsp.Pos = bsp.Dest
				bsp.Vel = pixel.ZV
				return
			}
		} else {
			if bsp.Pos.X < bsp.Dest.X {
				bsp.Pos = bsp.Dest
				bsp.Vel = pixel.ZV
				return
			}
		}

		if bsp.Vel.Y >= 0 {
			if bsp.Pos.Y > bsp.Dest.Y {
				bsp.Pos = bsp.Dest
				bsp.Vel = pixel.ZV
				return
			}
		} else {
			if bsp.Pos.Y < bsp.Dest.Y {
				bsp.Pos = bsp.Dest
				bsp.Vel = pixel.ZV
				return
			}
		}

		// still walking
	}
}
