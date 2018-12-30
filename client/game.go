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

	StatsText = "--------------------\n%d joules available\n\n%d advancing enemies\n\n%d bullets in flight\n--------------------"
	StoryText = "The earth is under attack\nby strange platforms,\nyou were sent to defend.\n\nBuild turrets.\nDestroy moving platforms.\nGreen square collects.\nTurrets prioritize enemies.\n\nBuild quickly,\nmore enemies soon."
	HelpText  = "Left Click:\n  Place new turret\n  (costs 20 joules)\n\nRight Click:\n  Order bullet spray\n  (costs 15 joules)\n\nr:\n  Reset world"
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
	win.SetMatrix(pixel.IM)

	gameCanvas := pixelgl.NewCanvas(pixel.R(-physics.MaxWindowBound, -physics.MaxWindowBound, physics.MaxWindowBound, physics.MaxWindowBound))

	imd := imdraw.New(nil)
	imd.Precision = 32

	fps := time.Tick(time.Second / physics.FpsTarget)

	frames := 0
	second := time.Tick(time.Second)

	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	txt := text.New(pixel.V(physics.MaxWindowBound+25, physics.MaxWindowBound-50), atlas)
	storyTxt := text.New(pixel.V(physics.MaxWindowBound+10, physics.MaxWindowBound-200), atlas)
	helpTxt := text.New(pixel.V(physics.MaxWindowBound+10, -physics.MaxWindowBound+200), atlas)
	txt.Color = colornames.Lightgrey
	storyTxt.Color = colornames.Lightgrey
	helpTxt.Color = colornames.Lightgrey
	numBullets := 0
	lastNumBullets := -999
	fmt.Fprintf(txt, StatsText, 0, 0, numBullets)
	fmt.Fprint(storyTxt, StoryText)
	fmt.Fprint(helpTxt, HelpText)

	// todo: own function
	// play the background song
	go playBackgroundMusic()

	var world = physics.NewWorld()

	for !win.Closed() {
		dt := physics.Dt
		mp := win.MousePosition()
		if mp.X > physics.MaxWindowBound {
			mp.X = physics.MaxWindowBound
		}

		// slow motion with tab
		if win.Pressed(pixelgl.KeyTab) {
			dt /= physics.SlowdownFactor
		}

		// reset world with r
		if win.JustPressed(pixelgl.KeyR) {
			world = physics.NewWorld()
		}

		// move shooters towards mouse location on left click or mouse scroll
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			if world.EnergyCount() >= 20 {
				world.AddShooter(mp)
				world.SubEnergy(20)
			}
		}

		// bullet spray with right click
		if win.JustPressed(pixelgl.MouseButtonRight) {
			if world.EnergyCount() >= 15 {
				world.BulletSpray(mp)
				world.SubEnergy(15)
			}
		}

		// step physics forward
		world.Update(dt, mp)

		if world.CheckLoseCondition() {
			world = physics.NewWorld()
		}

		// draw updated scene
		win.Clear(colornames.Darkslateblue)
		imd.Clear()
		gameCanvas.Clear(colornames.Black)
		world.Draw(imd)
		imd.Draw(gameCanvas)
		gameCanvas.Draw(win, pixel.IM)

		numBullets = world.NumBullets()
		if numBullets != lastNumBullets {
			txt.Clear()
			fmt.Fprintf(txt, StatsText, world.EnergyCount(), world.NumPlatforms(), numBullets)
			lastNumBullets = numBullets
		}
		txt.Draw(win, pixel.IM)
		storyTxt.Draw(win, pixel.IM)
		helpTxt.Draw(win, pixel.IM)

		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		case <-fps:
		}

		// save memory dump with f9 when debug env var is set
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

func playBackgroundMusic() {
	//return
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
		Volume:   -3,
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
}
