package scenes

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mharv/scrapyard-charter/basics"
	"github.com/mharv/scrapyard-charter/entities"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/solarlune/resolv"
)

type OverworldScene struct {
	entityManager                   entities.EntityManager
	menuBtn, castBtn, castAvailable bool
	physSpace                       *resolv.Space
	background                      *ebiten.Image
	spawnZone                       basics.FloatRect
	player                          entities.OverworldPlayerObject
	castDistance                    float64
}

const (
	cellSize = 8
)

func (o *OverworldScene) Init() {
	o.physSpace = resolv.NewSpace(globals.ScreenWidth, globals.ScreenHeight, cellSize, cellSize)
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
		o.AddTags("scrap", "solid")
	}

	img, _, err := ebitenutil.NewImageFromFile("images/overworldTerrainPlaceholderGrass.png")
	if err != nil {
		log.Fatal(err)
	} else {
		o.background = img
	}

	o.spawnZone.Width = globals.ScreenWidth
	o.spawnZone.Height = globals.ScreenHeight
	o.spawnZone.X = o.spawnZone.Width / 2
	o.spawnZone.Y = o.spawnZone.Height / 2

	// create homeBase

	t := &entities.HomeBaseObject{}
	t.Init("images/homeBase.png")
	t.GetPhysObj().AddTags("home", "solid")
	t.SetPosition(basics.Vector2f{X: o.spawnZone.X, Y: o.spawnZone.Y})
	o.physSpace.Add(t.GetPhysObj())
	o.entityManager.AddEntity(t)
	// Create player

	// NOTE: the below argument for CastDistanceLimit should be pulled from player stats struct
	p := &entities.OverworldPlayerObject{CastDistanceLimit: 200.0}
	p.Init("images/placeholderOverworldPlayerAssetTransparent.png")
	o.physSpace.Add(p.GetPhysObj())
	p.SetPosition(basics.Vector2f{X: o.spawnZone.X, Y: (o.spawnZone.Y + t.GetPhysObj().H)})
	o.entityManager.AddEntity(p)
	o.player = *p
}

func (o *OverworldScene) ReadInput() {
	o.entityManager.ReadInput()

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		o.menuBtn = true
	} else {
		o.menuBtn = false
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		inpututil.IsKeyJustPressed(ebiten.Key(ebiten.KeyE)) {
		o.castBtn = true
	} else {
		o.castBtn = false
	}
}

func (o *OverworldScene) Update(state *GameState, deltaTime float64) error {
	o.entityManager.Update(deltaTime)

	if o.menuBtn {
		t := &TitleScene{}
		state.SceneManager.GoTo(t, transitionTime)
	}

	if o.castAvailable && o.castBtn && o.castDistance < o.player.CastDistanceLimit {
		s := &ScavengeScene{distanceOfOverworldCast: o.castDistance}
		state.SceneManager.GoTo(s, transitionTime)
	}

	return nil
}

func (o *OverworldScene) Draw(screen *ebiten.Image) {

	// options := &ebiten.DrawImageOptions{}

	mouseX, mouseY := ebiten.CursorPosition()
	mx, my := o.physSpace.WorldToSpace(float64(mouseX), float64(mouseY))
	cx, cy := o.player.GetCellPosition()
	o.castDistance = math.Sqrt(math.Pow((float64(mx)-float64(cx))*8, 2) + math.Pow((float64(my)-float64(cy))*8, 2))
	drawColor := color.RGBA{255, 0, 0, 255}

	cellAtMouse := o.physSpace.Cell(mx, my)
	if cellAtMouse != nil {
		if cellAtMouse.ContainsTags("scrap") && o.castDistance < o.player.CastDistanceLimit {
			drawColor = color.RGBA{0, 255, 0, 255}
			o.castAvailable = true
		} else {
			o.castAvailable = false
		}
	}
	ebitenutil.DrawLine(screen, float64(cx)*cellSize, float64(cy)*cellSize, float64(mx)*cellSize, float64(my)*cellSize, drawColor)

	// screen.DrawImage(o.background, options)

	for _, o := range o.physSpace.Objects() {
		drawColor := color.RGBA{60, 60, 60, 255}
		if !o.HasTags("player") {
			ebitenutil.DrawRect(screen, o.X, o.Y, o.W, o.H, drawColor)
		}
	}

	o.entityManager.Draw(screen)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("cast available: %t", o.castAvailable))
}
