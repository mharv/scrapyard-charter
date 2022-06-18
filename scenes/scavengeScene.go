package scenes

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mharv/scrapyard-charter/basics"
	"github.com/mharv/scrapyard-charter/entities"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/solarlune/resolv"
)

type ScavengeScene struct {
	entityManager entities.EntityManager
	physSpace     *resolv.Space
	menuBtn       bool
	spawnZone     basics.FloatRect
}

const spawnZoneEdgeBorder = 128

func (s *ScavengeScene) Init() {
	s.physSpace = resolv.NewSpace(globals.ScreenWidth, globals.ScreenHeight, 16, 16)

	s.entityManager.Init()

	s.spawnZone.Width = globals.ScreenWidth - (spawnZoneEdgeBorder * 4)
	s.spawnZone.Height = globals.ScreenHeight - (spawnZoneEdgeBorder * 2)
	s.spawnZone.X = spawnZoneEdgeBorder * 2
	s.spawnZone.Y = spawnZoneEdgeBorder * 1.5

	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)

	// Create junk
	for i := 0; i < 10; i++ {
		j := &entities.JunkObject{}
		j.Init("images/oldpc.png")
		s.physSpace.Add(j.GetPhysObj())

		x := (rnd.Float64() * (s.spawnZone.Width))
		x += s.spawnZone.X - float64(j.GetPhysObj().X)
		y := (rnd.Float64() * (s.spawnZone.Height))
		y += s.spawnZone.Y - float64(j.GetPhysObj().Y)

		j.SetPosition(basics.Vector2f{X: x, Y: y})
		s.entityManager.AddEntity(j)
	}

	// Create magnet
	m := &entities.MagnetObject{}
	m.Init("images/magnet.png")
	s.physSpace.Add(m.GetPhysObj())
	s.physSpace.Add(m.GetFieldPhysObj())
	s.entityManager.AddEntity(m)

	// Create player
	p := &entities.ScavPlayerObject{}
	p.Init("images/player.png")
	p.SetMagnet(m)
	s.physSpace.Add(p.GetPhysObj())
	s.physSpace.Add(p.GetPhysObj())
	p.SetPosition(basics.Vector2f{X: s.spawnZone.X, Y: (s.spawnZone.Y - p.GetPhysObj().H)})
	s.entityManager.AddEntity(p)

	m.SetMagnetStartPos(p.GetFishingLinePoint())

	s.menuBtn = false
}

func (s *ScavengeScene) ReadInput() {
	s.entityManager.ReadInput()

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.menuBtn = true
	} else {
		s.menuBtn = false
	}
}

func (s *ScavengeScene) Update(state *GameState, deltaTime float64) error {
	s.entityManager.Update(deltaTime)

	if s.menuBtn {
		t := &TitleScene{}
		state.SceneManager.GoTo(t, transitionTime)
	}

	return nil
}

func (s *ScavengeScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	if globals.Debug {
		ebitenutil.DrawRect(screen, s.spawnZone.X, s.spawnZone.Y, s.spawnZone.Width, s.spawnZone.Height, color.RGBA{128, 96, 64, 255})
	}

	s.entityManager.Draw(screen)
}
