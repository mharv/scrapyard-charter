package scenes

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mharv/scrapyard-charter/basics"
	"github.com/mharv/scrapyard-charter/entities"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/solarlune/resolv"
)

type OverworldScene struct {
	entityManager entities.EntityManager
	menuBtn       bool
	physSpace     *resolv.Space
	background    *ebiten.Image
	spawnZone     basics.FloatRect
}

func (o *OverworldScene) Init() {
	o.physSpace = resolv.NewSpace(globals.ScreenWidth, globals.ScreenHeight, 16, 16)
	o.entityManager.Init()

	img, _, err := ebitenutil.NewImageFromFile("images/overworldTerrainPlaceholderGrass.png")
	if err != nil {
		log.Fatal(err)
	} else {
		o.background = img
	}

	o.spawnZone.Width = globals.ScreenWidth / 4
	o.spawnZone.Height = globals.ScreenHeight / 4
	o.spawnZone.X = 0
	o.spawnZone.Y = 0

	// Create player

	p := &entities.OverworldPlayerObject{}
	p.Init("images/placeholderOverworldPlayerAssetTransparent.png")
	o.physSpace.Add(p.GetPhysObj())
	p.SetPosition(basics.Vector2f{X: o.spawnZone.X, Y: (o.spawnZone.Y + p.GetPhysObj().H)})
	o.entityManager.AddEntity(p)
}

func (o *OverworldScene) ReadInput() {
	o.entityManager.ReadInput()

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		o.menuBtn = true
	} else {
		o.menuBtn = false
	}
}

func (o *OverworldScene) Update(state *GameState, deltaTime float64) error {
	o.entityManager.Update(deltaTime)

	if o.menuBtn {
		t := &TitleScene{}
		state.SceneManager.GoTo(t, transitionTime)
	}

	return nil
}

func (o *OverworldScene) Draw(screen *ebiten.Image) {
	// screen.Fill(color.RGBA{255, 0, 255, 255})

	options := &ebiten.DrawImageOptions{}

	screen.DrawImage(o.background, options)
	o.entityManager.Draw(screen)
}
