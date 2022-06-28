package scenes

import (
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mharv/scrapyard-charter/basics"
	"github.com/mharv/scrapyard-charter/entities"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/mharv/scrapyard-charter/mapgen"
	"github.com/mharv/scrapyard-charter/ui"
	"github.com/solarlune/resolv"
)

type Tile struct {
	X, Y  int
	Value float64
}

type TileMap struct {
	Xmax, Ymax int
	Tiles      []Tile
}

func (t *TileMap) GetTile(X, Y int) float64 {
	for _, v := range t.Tiles {
		if v.X == X && v.Y == Y {
			return v.Value
		}
	}
	return -1
}

type OverworldScene struct {
	entityManager                   entities.EntityManager
	menuBtn, castBtn, castAvailable bool
	physSpace                       *resolv.Space
	scrapspritesheet                *ebiten.Image
	overlayspritesheet              *ebiten.Image
	landspritesheet                 *ebiten.Image
	cursorNo                        *ebiten.Image
	cursorYes                       *ebiten.Image
	spawnZone                       basics.FloatRect
	player                          entities.OverworldPlayerObject
	castDistance                    float64
	ui                              ui.Ui
	terrain                         TileMap
}

const (
	cellSize      = 8
	tilesetcellsX = 4
	tilesetcellsY = 4
)

func (o *OverworldScene) Init() {
	globals.GetAudioPlayer().StopAllAudio()

	o.physSpace = resolv.NewSpace(globals.ScreenWidth, globals.ScreenHeight, cellSize, cellSize)
	o.entityManager.Init()
	o.ui = ui.Ui{}
	o.ui.Init()
	// object array
	geometry := []*resolv.Object{}

	// 1366 * 768
	var terrain [globals.ScreenWidth][globals.ScreenHeight]float64

	// create a terrain map L, R, U, D - if true, side is open
	terrain = mapgen.GenerateMap(false, false, false, false)

	// we create 32 x 32 pixel blocks
	tempCellSize := cellSize * 4

	// used for determining if something is scrap or land
	threshold := 0.8

	newTerrain := &TileMap{}

	i, j := 0, 0
	for x := 0; x < len(terrain); x += tempCellSize {
		j = 0
		for y := 0; y < len(terrain[x]); y += tempCellSize {
			newTerrain.Tiles = append(newTerrain.Tiles, Tile{X: i, Y: j, Value: terrain[x][y]})
			j++
		}
		i++
	}

	newTerrain.Xmax = i
	newTerrain.Ymax = j

	// create objects based off smoothed map
	for x := 0; x < i; x++ {
		for y := 0; y < j; y++ {
			randomChanceToAdd := newTerrain.GetTile(x, y)

			up, down, left, right := 0, 0, 0, 0

			if randomChanceToAdd > threshold {
				if x == 0 {
					left = 1
				} else if newTerrain.GetTile(x-1, y) > threshold {
					left = 1
				}
				if x == i-1 {
					right = 1
				} else if newTerrain.GetTile(x+1, y) > threshold {
					right = 1
				}

				if y == 0 {
					up = 1
				} else if newTerrain.GetTile(x, y-1) > threshold {
					up = 1
				}
				if y == j-1 {
					down = 1
				} else if newTerrain.GetTile(x, y+1) > threshold {
					down = 1
				}

				calculatevalue := 1*up + 2*left + 4*right + 8*down

				tempCellObject := resolv.NewObject(float64(x*tempCellSize), float64(y*tempCellSize), float64(tempCellSize), float64(tempCellSize), "scrap", "solid", strconv.Itoa(calculatevalue))

				geometry = append(geometry, tempCellObject)
			}

			if randomChanceToAdd <= threshold {
				tempCellObject := resolv.NewObject(float64(x*tempCellSize), float64(y*tempCellSize), float64(tempCellSize), float64(tempCellSize), "land", strconv.Itoa(rand.Intn(16)))
				geometry = append(geometry, tempCellObject)
			}
		}
	}

	o.terrain = *newTerrain

	o.scrapspritesheet = LoadImage("images/junkTileset.png")
	o.overlayspritesheet = LoadImage("images/junktileset2.png")
	o.landspritesheet = LoadImage("images/dirttileset.png")
	o.cursorNo = LoadImage("images/owCursorNo.png")
	o.cursorYes = LoadImage("images/owCursorYes.png")

	// add generated objects to scene space
	o.physSpace.Add(geometry...)

	o.spawnZone.Width = globals.ScreenWidth
	o.spawnZone.Height = globals.ScreenHeight
	o.spawnZone.X = o.spawnZone.Width/2 + 100
	o.spawnZone.Y = o.spawnZone.Height / 2

	// create homeBase
	t := &entities.HomeBaseObject{}
	t.Init("images/homeBase.png")
	t.GetPhysObj().AddTags("home", "solid")
	t.SetPosition(basics.Vector2f{X: o.spawnZone.X, Y: o.spawnZone.Y})
	o.physSpace.Add(t.GetPhysObj())
	o.entityManager.AddEntity(t)

	// create crafting zone around homebase
	o.physSpace.Add(t.GetCraftZone())

	// Create player
	p := &entities.OverworldPlayerObject{}
	p.Init("images/overworldplayer.png")
	o.physSpace.Add(p.GetPhysObj())
	p.SetPosition(globals.GetPlayerData().GetPlayerPosition())
	o.entityManager.AddEntity(p)
	o.player = *p
}

func (o *OverworldScene) ReadInput() {
	o.entityManager.ReadInput()
	o.ui.ReadInput()

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
	globals.GetAudioPlayer().PlayFile("audio/overworld.mp3")

	o.entityManager.Update(deltaTime)
	o.ui.Update(deltaTime)

	// if o.menuBtn {
	// 	if !o.ui.IsOpen() {

	// 		t := &TitleScene{}
	// 		state.SceneManager.GoTo(t, transitionTime)

	// 	} else {

	// 	}
	// }

	if o.castAvailable && o.castBtn && o.castDistance < o.player.CastDistanceLimit {
		s := &ScavengeScene{distanceOfOverworldCast: o.castDistance}
		globals.GetPlayerData().SetPlayerPosition(basics.Vector2f{X: o.player.GetPhysObj().X, Y: o.player.GetPhysObj().Y})
		state.SceneManager.GoTo(s, transitionTime)
	}

	winCondition, err := globals.GetPlayerData().GetEquippedItem("Magnet")
	if err == nil {
		if winCondition.GetKeyItemName() == "GOLDENMAGNET" {
			w := &WinScene{}
			state.SceneManager.GoTo(w, transitionTime)
		}
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
	for _, tile := range o.physSpace.Objects() {
		if tile.HasTags("scrap") {
			index := 0
			if tile.HasTags("0") {
				index = 0
			} else if tile.HasTags("1") {
				index = 1
			} else if tile.HasTags("2") {
				index = 2
			} else if tile.HasTags("3") {
				index = 3
			} else if tile.HasTags("4") {
				index = 4
			} else if tile.HasTags("5") {
				index = 5
			} else if tile.HasTags("6") {
				index = 6
			} else if tile.HasTags("7") {
				index = 7
			} else if tile.HasTags("8") {
				index = 8
			} else if tile.HasTags("9") {
				index = 9
			} else if tile.HasTags("10") {
				index = 10
			} else if tile.HasTags("11") {
				index = 11
			} else if tile.HasTags("12") {
				index = 12
			} else if tile.HasTags("13") {
				index = 13
			} else if tile.HasTags("14") {
				index = 14
			} else if tile.HasTags("15") {
				index = 15
			}

			sx := index % tilesetcellsX
			sy := (index - sx) / tilesetcellsY

			sx *= 32
			sy *= 32

			options := &ebiten.DrawImageOptions{}
			options.GeoM.Translate(tile.X, tile.Y)
			screen.DrawImage(o.scrapspritesheet.SubImage(image.Rect(sx, sy, sx+32, sy+32)).(*ebiten.Image), options)
		}
		if tile.HasTags("land") {
			index := 0
			if tile.HasTags("0") {
				index = 0
			} else if tile.HasTags("1") {
				index = 1
			} else if tile.HasTags("2") {
				index = 2
			} else if tile.HasTags("3") {
				index = 3
			} else if tile.HasTags("4") {
				index = 4
			} else if tile.HasTags("5") {
				index = 5
			} else if tile.HasTags("6") {
				index = 6
			} else if tile.HasTags("7") {
				index = 7
			} else if tile.HasTags("8") {
				index = 8
			} else if tile.HasTags("9") {
				index = 9
			} else if tile.HasTags("10") {
				index = 10
			} else if tile.HasTags("11") {
				index = 11
			} else if tile.HasTags("12") {
				index = 12
			} else if tile.HasTags("13") {
				index = 13
			} else if tile.HasTags("14") {
				index = 14
			} else if tile.HasTags("15") {
				index = 15
			}

			sx := index % tilesetcellsX
			sy := (index - sx) / tilesetcellsY

			sx *= 32
			sy *= 32

			options := &ebiten.DrawImageOptions{}
			options.GeoM.Translate(tile.X, tile.Y)
			screen.DrawImage(o.landspritesheet.SubImage(image.Rect(sx, sy, sx+32, sy+32)).(*ebiten.Image), options)
		}
	}

	o.entityManager.Draw(screen)
	o.DrawOverlay(screen)
	o.ui.Draw(screen)

	if !o.ui.IsOpen() {
		mop := &ebiten.DrawImageOptions{}
		mop.GeoM.Translate(float64(mx)*cellSize, float64(my)*cellSize)
		if o.castAvailable {
			mop.GeoM.Translate(-float64(o.cursorYes.Bounds().Dx())/2, -float64(o.cursorYes.Bounds().Dy())/2)
			screen.DrawImage(o.cursorYes, mop)
		} else {
			mop.GeoM.Translate(-float64(o.cursorNo.Bounds().Dx())/2, -float64(o.cursorNo.Bounds().Dy())/2)
			screen.DrawImage(o.cursorNo, mop)
		}
	}

	// draw the mouse to character distance check line
	if globals.Debug {
		ebitenutil.DrawLine(screen, float64(cx)*cellSize, float64(cy)*cellSize, float64(mx)*cellSize, float64(my)*cellSize, drawColor)
	}
}

func (o *OverworldScene) DrawOverlay(screen *ebiten.Image) {
	src := rand.NewSource(int64(globals.GetPlayerData().GetWorldSeed()))
	rnd := rand.New(src)

	for _, tile := range o.physSpace.Objects() {
		if tile.HasTags("scrap") {
			if rnd.Intn(2) != 0 {
				continue
			} else {
				index := rnd.Intn(13)
				sx := index % tilesetcellsX
				sy := (index - sx) / tilesetcellsY

				sx *= 32
				sy *= 32

				options := &ebiten.DrawImageOptions{}
				options.GeoM.Translate(tile.X, tile.Y-float64(rnd.Intn(16)))
				screen.DrawImage(o.overlayspritesheet.SubImage(image.Rect(sx, sy, sx+32, sy+32)).(*ebiten.Image), options)
			}
		}
	}
}

func LoadImage(filepath string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	return img
}
