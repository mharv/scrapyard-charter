package scenes

import (
	"image/color"
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
	o.physSpace = resolv.NewSpace(globals.ScreenWidth, globals.ScreenHeight, 8, 8)
	o.entityManager.Init()

	// Construct geometry
	geometry := []*resolv.Object{
		// left wall
		// resolv.NewObject(0, 0, 16, globals.ScreenHeight),
		// right wall
		resolv.NewObject(globals.ScreenWidth-16, 0, 16, globals.ScreenHeight),
		// top wall
		resolv.NewObject(0, 0, globals.ScreenWidth, 16),
		// bottom wall
		resolv.NewObject(0, globals.ScreenHeight-16, globals.ScreenWidth, 16),

		// middle wall
		resolv.NewObject(400, 0, 80, globals.ScreenHeight/2),
	}

	o.physSpace.Add(geometry...)

	for _, o := range o.physSpace.Objects() {
		o.AddTags("solid")
	}

	img, _, err := ebitenutil.NewImageFromFile("images/overworldTerrainPlaceholderGrass.png")
	if err != nil {
		log.Fatal(err)
	} else {
		o.background = img
	}

	o.spawnZone.Width = globals.ScreenWidth
	o.spawnZone.Height = globals.ScreenHeight
	o.spawnZone.X = 0
	o.spawnZone.Y = 0

	// create homeBase

	t := &entities.HomeBaseObject{}
	t.Init("images/homeBase.png")
	t.GetPhysObj().AddTags("solid")
	t.SetPosition(basics.Vector2f{X: o.spawnZone.Width / 2, Y: o.spawnZone.Height / 2})
	o.physSpace.Add(t.GetPhysObj())
	o.entityManager.AddEntity(t)
	// Create player

	p := &entities.OverworldPlayerObject{}
	p.Init("images/placeholderOverworldPlayerAssetTransparent.png")
	o.physSpace.Add(p.GetPhysObj())
	p.SetPosition(basics.Vector2f{X: o.spawnZone.X + p.GetPhysObj().W, Y: (o.spawnZone.Y + p.GetPhysObj().H)})
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

	// options := &ebiten.DrawImageOptions{}

	// screen.DrawImage(o.background, options)

	for _, o := range o.physSpace.Objects() {
		drawColor := color.RGBA{60, 60, 60, 255}
		if !o.HasTags("player") {
			ebitenutil.DrawRect(screen, o.X, o.Y, o.W, o.H, drawColor)
		}

	}

	o.entityManager.Draw(screen)
}
