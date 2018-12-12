package main

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/gobuffalo/packr"
	"github.com/robkau/squadron-ed-39/physics"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"log"
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

	var world = physics.NewWorld()

	canvas := pixelgl.NewCanvas(pixel.R(-800, -500, 800, 500))
	imd := imdraw.New(nil)
	imd.Precision = 32

	fps := time.Tick(time.Second / 144)

	frames := 0
	second := time.Tick(time.Second)

	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	txt := text.New(pixel.V(-750, 450), atlas)
	txt.Color = colornames.Lightgrey

	// play the background song
	go func() {
		return // this experiment is disabled for now
		// it will extract music assets packed into the binary and play them as background sound
		box := packr.NewBox("./assets")

		// Decode the packed .mp3 file
		f, err := box.Open("song.mp3")
		if err != nil {
			log.Fatal(err)
		}
		s, format, _ := mp3.Decode(f)

		// Init the Speaker with the SampleRate of the format and a buffer size of 1/10s
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

		// Initiate control channel
		playing := make(chan struct{})

		// Play the sound
		speaker.Play(beep.Seq(s, beep.Callback(func() {
			// Callback after the stream Ends
			close(playing)
		})))
		<-playing
	}()

	// spawn bullets from cursor
	bulletSpawner := time.NewTicker(physics.BulletSpawnInterval)
	go func() {
		// spawn new bullet towards mouse
		for range bulletSpawner.C {
			mp := win.MousePosition()
			v := pixel.Lerp(pixel.ZV, mp, physics.BulletSpeedFactor)
			b := &physics.Bullet{
				Pos:  pixel.ZV,
				Dest: mp,
				Vel:  v,
			}
			physics.EnforceMinBulletSpeed(b)
			world.SpawnBullet(b)
		}
	}()

	for !win.Closed() {
		dt := physics.Dt

		//canvas.SetMatrix(pixel.Matrix{1, 0, 0, 1, 0, 0})

		// slow motion with tab
		if win.Pressed(pixelgl.KeyTab) {
			dt /= physics.SlowdownFactor
		}

		ctrl := pixel.ZV
		// update the physics and animation
		world.Update(dt, ctrl)

		// draw the scene to the canvas using IMDraw
		canvas.Clear(colornames.Black)
		imd.Clear()

		world.Draw(imd)
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
		txt.Draw(win, pixel.Matrix{1, 0, 0, 1, 0, 0})

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
