package scenes

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mharv/scrapyard-charter/basics"
	"github.com/mharv/scrapyard-charter/entities"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/solarlune/resolv"
	"github.com/tinne26/etxt"
)

type ScavengeScene struct {
	entityManager           entities.EntityManager
	physSpace               *resolv.Space
	menuBtn                 bool
	spawnZone               basics.FloatRect
	distanceOfOverworldCast float64
	txtRenderer             *etxt.Renderer
	countdownTimer          float64
}

const (
	spawnZoneEdgeBorder = 128
	defaultTimerStart   = 60
	timerPositionOffset = 100
)

func (s *ScavengeScene) Init() {
	s.physSpace = resolv.NewSpace(globals.ScreenWidth, globals.ScreenHeight, 16, 16)

	s.entityManager.Init()

	s.countdownTimer = defaultTimerStart

	fontLib := etxt.NewFontLibrary()

	_, _, err := fontLib.ParseDirFonts("fonts")
	if err != nil {
		log.Fatal(err)
	}

	if !fontLib.HasFont("Rajdhani Regular") {
		log.Fatal("missing font Rajdhani-Regular.ttf")
	}

	s.txtRenderer = etxt.NewStdRenderer()
	glyphsCache := etxt.NewDefaultCache(10 * 1024 * 1024) // 10MB
	s.txtRenderer.SetCacheHandler(glyphsCache.NewHandler())
	s.txtRenderer.SetFont(fontLib.GetFont("Rajdhani Regular"))
	s.txtRenderer.SetAlign(etxt.Top, etxt.Left)
	s.txtRenderer.SetSizePx(24)

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
	p.SetPosition(basics.Vector2f{X: s.spawnZone.X, Y: (s.spawnZone.Y - p.GetPhysObj().H)})
	s.entityManager.AddEntity(p)

	r := &entities.ScavRodObject{}
	r.Init("images/rodSection.png")
	r.SetRoot(p.GetFishingRodStartPoint())
	r.SetTip(p.GetFishingRodEndPoint())
	r.SetMagnetPosition(m.GetMagnetPos())
	s.entityManager.AddEntity(r)

	m.SetMagnetStartPos(p.GetFishingRodEndPoint())

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

	s.countdownTimer -= deltaTime
	if s.countdownTimer <= 0 {
		o := &OverworldScene{}
		state.SceneManager.GoTo(o, transitionTime)
	}

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

	timer := fmt.Sprintf("%.2f", s.countdownTimer)

	s.entityManager.Draw(screen)
	s.txtRenderer.SetTarget(screen)
	if s.countdownTimer > 10 {
		s.txtRenderer.SetColor(color.RGBA{255, 255, 255, 255})
	} else {
		s.txtRenderer.SetColor(color.RGBA{255, 0, 0, 255})
	}
	s.txtRenderer.Draw(timer, globals.ScreenWidth-timerPositionOffset, 0)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("cast available: %f", s.distanceOfOverworldCast))
}
