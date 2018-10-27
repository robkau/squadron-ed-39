package main

import (
	"fmt"
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
		Title:  "Squadron E.D. 39",
		Bounds: pixel.R(-800, -500, 800, 500),
		VSync:  false,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	nb := 10
	var bullets = make([]*bullet, nb*40000)

	for i := -10; i < nb; i++ {
		bullets = append(bullets, &bullet{
			pos:  pixel.Vec{X: 0, Y: 0},
			dest: pixel.Vec{X: -20 * float64(i), Y: 480},
			vel:  pixel.Lerp(pixel.ZV, pixel.Vec{X: -20 * float64(i), Y: 480}, 0.05),
		})
	}

	var phys = new(objects)
	phys.bullets = bullets

	// hardcoded level
	platforms := []*platform{
		{rect: pixel.R(-1024/2, 420, 1024/2, 440)},
		{rect: pixel.R(-1099/2, 450, 1099/2, 480)},
		{rect: pixel.R(-1024/2, -300, 1524/2, -240)},
	}
	for i := range platforms {
		platforms[i].color = randomNiceColor()
	}

	canvas := pixelgl.NewCanvas(pixel.R(-800, -500, 800, 500))
	imd := imdraw.New(nil)
	imd.Precision = 32

	fps := time.Tick(time.Second / 60)
	bulletSpawn := time.Tick(time.Second / 30)
	var bulletSpeedFactor float64 = 0.15
	frames := 0
	second := time.Tick(time.Second)
	last := time.Now()

	for !win.Closed() {

		dt := time.Since(last).Seconds()
		last = time.Now()

		// spawn new bullets towards mouse
		select {
		case <-bulletSpawn:
			mp := win.MousePosition()
			phys.bullets = append(phys.bullets, &bullet{
				pos:  pixel.ZV,
				dest: mp,
				vel:  pixel.Lerp(pixel.ZV, mp, bulletSpeedFactor),
			})
		default:
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

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		case <-fps:
		}
	}
}
