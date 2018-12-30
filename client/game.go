package main

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
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
	"os"
	"runtime/pprof"
	"time"
)

const (
	CpuProfile  = "cpu.pprof"
	DebugEnvVar = "sq39_debug"
	StatUiWidth = 200
)

func startFreePlay(debugSet bool) {
	cfg := pixelgl.WindowConfig{
		Title:  "Squadron E.D. 39",
		Bounds: pixel.R(-physics.MaxWindowBound, -physics.MaxWindowBound, physics.MaxWindowBound+StatUiWidth, physics.MaxWindowBound),
		VSync:  false,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(false)

	GameOffset := pixel.Vec{X: 0}
	StatOffset := pixel.Vec{X: physics.MaxWindowBound + StatUiWidth/2}

	gameCanvas := pixelgl.NewCanvas(pixel.R(-physics.MaxWindowBound, -physics.MaxWindowBound, physics.MaxWindowBound, physics.MaxWindowBound))
	statCanvas := pixelgl.NewCanvas(pixel.R(-StatUiWidth/2, -physics.MaxWindowBound, StatUiWidth/2, physics.MaxWindowBound))

	win.SetMatrix(pixel.IM)

	imd := imdraw.New(nil)
	imd.Precision = 32

	fps := time.Tick(time.Second / physics.FpsTarget)

	frames := 0
	second := time.Tick(time.Second)

	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	txt := text.New(pixel.V(physics.MaxWindowBound+15, physics.MaxWindowBound-50), atlas)
	txt.Color = colornames.Lightgrey
	numBullets := 0
	lastNumBullets := -999
	fmt.Fprintf(txt, "%d bullets\n%d joules", numBullets, 0)

	// todo: own function
	// play the background song
	go func() {
		//return // this experiment is disabled for now
		//it will extract music assets packed into the binary and play them in the background
		box := packr.NewBox("./assets")

		// Decode the packed .mp3 file
		f, err := box.Open("song.mp3")
		if err != nil {
			log.Fatal(err)
		}
		s, format, _ := mp3.Decode(f)
		v := effects.Volume{
			Streamer: s,
			Base:     2,
			Volume:   -4,
		}

		// Init the Speaker with the SampleRate of the format and a buffer size of 1/10s
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

		// Initiate control channel
		playing := make(chan struct{})

		// Play the sound
		speaker.Play(beep.Seq(&v, beep.Callback(func() {
			// Callback after the stream Ends
			close(playing)
		})))
		<-playing
	}()

	var world = physics.NewWorld()

	for !win.Closed() {
		dt := physics.Dt
		mp := win.MousePosition().Sub(GameOffset) // rebase towards game canvas
		if mp.X > physics.MaxWindowBound {
			mp.X = physics.MaxWindowBound
		}

		// slow motion with tab
		if win.Pressed(pixelgl.KeyTab) {
			dt /= physics.SlowdownFactor
		}

		// move shooters towards mouse location on left click or mouse scroll
		if win.JustPressed(pixelgl.MouseButtonLeft) || win.MouseScroll().Y != 0 {
			world.SetShooterDestination(mp)
		}

		// bullet spray with right click
		if win.JustPressed(pixelgl.MouseButtonRight) {
			world.BulletSpray(mp)
		}

		// step physics forward
		world.Update(dt, mp)

		// draw updated scene
		gameCanvas.Clear(colornames.Black)
		statCanvas.Clear(colornames.Darkblue)
		imd.Clear()
		world.Draw(imd)
		imd.Draw(gameCanvas)

		// stretch the canvas to the window
		// todo: stats UI
		win.Clear(colornames.White)

		gameCanvas.Draw(win, pixel.IM.Moved(GameOffset)) //).Moved(canvas.Bounds().Center()))
		statCanvas.Draw(win, pixel.IM.Moved(StatOffset)) //).Moved(canvas.Bounds().Center()))

		numBullets = world.NumBullets()
		if numBullets != lastNumBullets {
			txt.Clear()
			fmt.Fprintf(txt, "%d bullets\n%d joules", numBullets, world.EnergyCount())
			lastNumBullets = numBullets
		}
		txt.Draw(win, pixel.IM)

		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		case <-fps:
		}

		// save memory dump with f9
		if debugSet && win.JustPressed(pixelgl.KeyF9) {
			fp := "mem.mprof"
			f, err := os.Create(fp)
			if err != nil {
				log.Fatal(err)
			}
			pprof.WriteHeapProfile(f)
			f.Close()
		}
	}
}
