package scenes

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mharv/scrapyard-charter/basics"
	"github.com/mharv/scrapyard-charter/entities"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/mharv/scrapyard-charter/mapgen"
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

	// object array
	geometry := []*resolv.Object{}

	// 1366 * 768
	var terrain [globals.ScreenWidth][globals.ScreenHeight]float64

	// create a terrain map L, R, U, D - if true, side is open
	terrain = mapgen.GenerateMap(false, true, false, true)

	// we create 16 x 16 pixel blocks
	tempCellSize := cellSize * 4

	// used for determining if something is scrap or land
	threshold := 0.8
	// create objects based off smoothed map
	for x := 0; x < len(terrain); x += tempCellSize {
		for y := 0; y < len(terrain[x]); y += tempCellSize {
			randomChanceToAdd := terrain[x][y]

			if randomChanceToAdd <= threshold {
				tempCellObject := resolv.NewObject(float64(x), float64(y), float64(tempCellSize), float64(tempCellSize), "land")
				geometry = append(geometry, tempCellObject)
			}
			if randomChanceToAdd > threshold {
				tempCellObject := resolv.NewObject(float64(x), float64(y), float64(tempCellSize), float64(tempCellSize), "scrap", "solid")
				geometry = append(geometry, tempCellObject)
			}
		}
	}

	// add generated objects to scene space
	o.physSpace.Add(geometry...)

	// img, _, err := ebitenutil.NewImageFromFile("images/overworldTerrainPlaceholderGrass.png")
	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	o.background = img
	// }

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

	// cursor to player drawline and cast valid checks
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

	// draws the color depending on the tags for each object belonging to space
	for _, o := range o.physSpace.Objects() {
		if o.HasTags("scrap") {
			drawColor := color.RGBA{60, 60, 60, 255}
			ebitenutil.DrawRect(screen, o.X, o.Y, o.W, o.H, drawColor)
		}
		if o.HasTags("beach") {
			drawColor := color.RGBA{222, 130, 22, 255}
			ebitenutil.DrawRect(screen, o.X, o.Y, o.W, o.H, drawColor)
		}
		if o.HasTags("land") {
			drawColor := color.RGBA{119, 174, 74, 255}
			ebitenutil.DrawRect(screen, o.X, o.Y, o.W, o.H, drawColor)
		}
	}

	o.entityManager.Draw(screen)

	// draw the mouse to character distance check line
	ebitenutil.DrawLine(screen, float64(cx)*cellSize, float64(cy)*cellSize, float64(mx)*cellSize, float64(my)*cellSize, drawColor)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("cast available: %t", o.castAvailable))
}
