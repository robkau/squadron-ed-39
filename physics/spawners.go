package physics

import "github.com/faiface/pixel"

// todo: interface for Moveable()
type BulletSpawner struct {
	Pos  pixel.Vec
	Dest pixel.Vec
	Vel  pixel.Vec
	// todo: control spawn rate with var here
}

func (world *world) SpawnBullet(b *Bullet) {
	world.bullets = append(world.bullets, b)
}

func (world *world) BulletSpray(dest pixel.Vec) {
	// shoot shotgun blast towards dest from bsp
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
