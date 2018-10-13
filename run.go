package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"math"
	"math/rand"
	"time"
)

func run() {
	rand.Seed(time.Now().UnixNano())

	cfg := pixelgl.WindowConfig{
		Title:  "Platformer",
		Bounds: pixel.R(-800, -500, 800, 500),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	nb := 10
	var bullets = make([]*bullet, nb*2, nb*40)

	for i := -10; i < nb; i++ {
		var x float64
		x = float64(i * 25)
		bullets[i+nb] = &bullet{
			rect: pixel.R(x-5, -7, x+5, 7),
			dest: pixel.Vec{X: -20 * float64(i), Y: 480},
		}
	}

	var phys = new(objects)
	phys.bullets = bullets

	// hardcoded level
	platforms := []*platform{
		{rect: pixel.R(-1024/2, 480, 1024/2, 500)},
	}
	for i := range platforms {
		platforms[i].color = randomNiceColor()
	}

	canvas := pixelgl.NewCanvas(pixel.R(-800, -500, 800, 500))
	imd := imdraw.New(nil)
	imd.Precision = 32

	lastBulletSpawn := 0
	bulletSpawnDiffFrames := 1

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		// spawn new bullets towards mouse
		if lastBulletSpawn >= bulletSpawnDiffFrames {
			phys.bullets = append(phys.bullets, &bullet{
				rect: pixel.R(-5, -7, 5, 7),
				dest: win.MousePosition(),
			})
			lastBulletSpawn = 0
		} else {
			lastBulletSpawn += 1
		}

		//canvas.SetMatrix(pixel.Matrix{1, 0, 0, 1, 0, 0})

		// slow motion with tab
		if win.Pressed(pixelgl.KeyTab) {
			dt /= 8
		}

		// control the gopher with keys
		ctrl := pixel.ZV
		// update the physics and animation
		phys.update(dt, ctrl, &platforms)

		// draw the scene to the canvas using IMDraw
		canvas.Clear(colornames.Black)
		imd.Clear()
		for _, p := range platforms {
			p.draw(imd)
		}
		draw(imd, phys)
		imd.Draw(canvas)

		// stretch the canvas to the window
		win.Clear(colornames.White)
		win.SetMatrix(pixel.IM.Scaled(pixel.ZV,
			math.Min(
				win.Bounds().W()/canvas.Bounds().W(),
				win.Bounds().H()/canvas.Bounds().H(),
			),
		).Moved(win.Bounds().Center()))
		canvas.Draw(win, pixel.IM.Moved(canvas.Bounds().Center()))
		win.Update()
	}
}
