package physics

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"github.com/gobuffalo/packr"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"log"
	"math"
	"time"
)

type world struct {
	shooters      []*BulletSpawner
	bullets       []*Bullet
	platforms     []*platform
	collectors    []*collector
	colliders     []collideable
	BulletPool    *BulletPool
	iteration     int
	bulletCounter int
	energyCount   int
	atlas         *text.Atlas
	deadBullet    bool
	deadPlatform  bool
}

func NewWorld() *world {
	w := &world{
		shooters:    make([]*BulletSpawner, 0),
		bullets:     make([]*Bullet, 0),
		platforms:   make([]*platform, 0),
		collectors:  make([]*collector, 0),
		colliders:   make([]collideable, 0),
		BulletPool:  NewPool(BulletPoolSize),
		atlas:       text.NewAtlas(basicfont.Face7x13, text.ASCII),
		energyCount: 40,
	}

	w.AddCollector(pixel.Vec{X: 0, Y: -350})

	// play the background song
	go playBackgroundMusic()

	return w
}

func (world *world) AddPlatform(pos pixel.Rect, dest pixel.Vec, health int) {
	dir := dest.Sub(pos.Center())
	pVel := dir.Scaled(PlatformSpeed / dir.Len())
	world.AddPlatformWithV(pos, dest, pVel, health)
}

func (world *world) AddPlatformWithV(pos pixel.Rect, dest pixel.Vec, vel pixel.Vec, health int) {
	p := &platform{LinearRectMovingStrategy: LinearRectMovingStrategy{rect: pos, dest: dest, vel: vel},
		Health: health, Color: colornames.Lightseagreen,
		UniqueId: randomHex(16)}
	world.platforms = append(world.platforms, p)
	world.colliders = append(world.colliders, p)
}

func (world *world) AddCollector(pos pixel.Vec) {
	cl := &collector{LinearRectMovingStrategy: LinearRectMovingStrategy{rect: pixel.Rect{Min: pos.Sub(pixel.Vec{X: 25, Y: 25}), Max: pos.Add(pixel.Vec{X: 25, Y: 25})}},
		UniqueId: randomHex(16)}
	world.collectors = append(world.collectors, cl)
	world.colliders = append(world.colliders, cl)
}

func (world *world) AddShooter(pos pixel.Vec) {
	sh := &BulletSpawner{&LinearPointMovingStrategy{pos: pos, stopAtDest: true}}
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

func (world *world) SubEnergy(e int) {
	world.energyCount -= e
}

func (world *world) lowestPlatform() *platform {
	lowest := math.Inf(1)
	lowestIndex := -1
	for i, p := range world.platforms {
		if p.Pos().Y < lowest {
			lowest = p.Pos().Y
			lowestIndex = i
		}
	}
	return world.platforms[lowestIndex]
}

func (world *world) CheckLoseCondition() bool {
	for _, pl := range world.platforms {
		if pl.Pos().Y < -MaxWindowBound {
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
