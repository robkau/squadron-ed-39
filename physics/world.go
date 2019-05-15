package physics

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/gobuffalo/packr"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"log"
	"math"
	"time"
)

type World interface {
	Update(dt float64, mp pixel.Vec)
	Draw(imd *imdraw.IMDraw)
	ProcessInput(button pixelgl.Button, mp pixel.Vec)

	// todo: collapse world state into one struct
	NumBullets() int
	NumPlatforms() int
	EnergyCount() int
	CheckLoseCondition() bool
}

type world struct {
	shooters      []*BulletSpawner
	bullets       []*bullet
	platforms     []*platform
	collectors    []*collector
	colliders     []collideable
	BulletPool    *bulletPool
	iteration     int
	bulletCounter int
	energyCount   int
	atlas         *text.Atlas
	deadBullet    bool
	deadPlatform  bool
	size          float64
}

func NewWorld(size int) World {
	w := &world{
		shooters:    make([]*BulletSpawner, 0),
		bullets:     make([]*bullet, 0),
		platforms:   make([]*platform, 0),
		collectors:  make([]*collector, 0),
		colliders:   make([]collideable, 0),
		BulletPool:  newPool(BulletPoolSize),
		atlas:       text.NewAtlas(basicfont.Face7x13, text.ASCII),
		energyCount: 40,
		size:        float64(size),
	}

	w.addCollector(pixel.Vec{X: 0, Y: -350})

	// play the background song
	go playBackgroundMusic()

	return w
}

func (world *world) ProcessInput(button pixelgl.Button, mp pixel.Vec) {
	// move shooters towards mouse location on left click or mouse scroll
	switch button {
	case pixelgl.MouseButtonLeft:
		if world.EnergyCount() >= 20 {
			world.addShooter(mp)
			world.subEnergy(20)
		}
	case pixelgl.MouseButtonRight:
		if world.EnergyCount() >= 15 {
			world.bulletSpray(mp)
			world.subEnergy(15)
		}
	}
}

func (world *world) addPlatform(pos pixel.Rect, dest pixel.Vec, health int) {
	dir := dest.Sub(pos.Center())
	pVel := dir.Scaled(PlatformSpeed / dir.Len())
	world.addPlatformWithV(pos, dest, pVel, health)
}

func (world *world) addPlatformWithV(pos pixel.Rect, dest pixel.Vec, vel pixel.Vec, health int) {
	p := &platform{linearRectMovingStrategy: linearRectMovingStrategy{r: pos, dst: dest, v: vel},
		Health: health, Color: colornames.Lightseagreen,
		UniqueId: randomHex(16)}
	world.platforms = append(world.platforms, p)
	world.colliders = append(world.colliders, p)
}

func (world *world) addCollector(pos pixel.Vec) {
	cl := &collector{linearRectMovingStrategy: linearRectMovingStrategy{r: pixel.Rect{Min: pos.Sub(pixel.Vec{X: 25, Y: 25}), Max: pos.Add(pixel.Vec{X: 25, Y: 25})}},
		uniqueId: randomHex(16)}
	world.collectors = append(world.collectors, cl)
	world.colliders = append(world.colliders, cl)
}

func (world *world) addShooter(pos pixel.Vec) {
	sh := &BulletSpawner{&linearPointMovingStrategy{p: pos, stopAtDest: true}}
	world.shooters = append(world.shooters, sh)
}

func (world *world) NumBullets() int {
	return world.bulletCounter
}

func (world *world) NumPlatforms() int {
	return len(world.platforms)
}

func (world *world) EnergyCount() int {
	return world.energyCount
}

func (world *world) subEnergy(e int) {
	world.energyCount -= e
}

func (world *world) lowestPlatform() *platform {
	lowest := math.Inf(1)
	lowestIndex := -1
	for i, p := range world.platforms {
		if p.pos().Y < lowest {
			lowest = p.pos().Y
			lowestIndex = i
		}
	}
	return world.platforms[lowestIndex]
}

func (world *world) CheckLoseCondition() bool {
	for _, pl := range world.platforms {
		if pl.pos().Y < -world.size {
			return true
		}
	}
	return false
}

func randomHex(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		panic("failed to read random bytes for hex ID")
	}
	return hex.EncodeToString(bytes)
}

func playBackgroundMusic() {
	//extract music assets packed into the binary and play them in the background
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
		Volume:   -2,
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
