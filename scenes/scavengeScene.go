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
	timerUIboxSprite        *ebiten.Image
	timerUIglassSprite      *ebiten.Image
	UIPipeSprite            *ebiten.Image
	physSpace               *resolv.Space
	menuBtn                 bool
	spawnZone               basics.FloatRect
	distanceOfOverworldCast float64
	txtRenderer             *etxt.Renderer
	countdownTimer          float64
	junkList                []entities.JunkObject
	UIPosition              basics.Vector2f
}

const (
	spawnZoneEdgeBorder = 128
	defaultTimerStart   = 10
	uiXOffset           = 32
	uiYOffset           = 32
	uiGlassXOffset      = 25
	uiGlassYOffset      = 4
	textXOffset         = 40
	textYOffset         = 0
	fontSize            = 50
	iconXOffset         = 184
	iconYOffset         = 66
	textRedLimit        = 5.0
)

func (s *ScavengeScene) Init() {
	fmt.Println(s.distanceOfOverworldCast)
	s.physSpace = resolv.NewSpace(globals.ScreenWidth, globals.ScreenHeight, 16, 16)
	s.UIPosition = basics.Vector2f{X: globals.ScreenWidth - (uiXOffset + iconXOffset), Y: uiYOffset + iconYOffset}

	img, _, imgerr := ebitenutil.NewImageFromFile("images/timerUIBox.png")
	if imgerr != nil {
		log.Fatal(imgerr)
	}
	s.timerUIboxSprite = img

	glimg, _, glerr := ebitenutil.NewImageFromFile("images/timerUIBoxGlass.png")
	if glerr != nil {
		log.Fatal(glerr)
	}
	s.timerUIglassSprite = glimg

	pipimg, _, piperr := ebitenutil.NewImageFromFile("images/UIPipes.png")
	if piperr != nil {
		log.Fatal(piperr)
	}
	s.UIPipeSprite = pipimg

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
	s.txtRenderer.SetSizePx(fontSize)

	s.spawnZone.Width = globals.ScreenWidth - (spawnZoneEdgeBorder * 4)
	s.spawnZone.Height = globals.ScreenHeight - (spawnZoneEdgeBorder * 2)
	s.spawnZone.X = spawnZoneEdgeBorder * 2
	s.spawnZone.Y = spawnZoneEdgeBorder * 1.5

	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)

	s.InitJunkList()

	junkLookup := make(map[*resolv.Object]*entities.JunkObject)
	// Create junk
	for i := 0; i < 100; i++ {
		j := s.SelectJunk(s.distanceOfOverworldCast)
		depth := j.GetItemData().GetDepth()
		percent := depth/float64(len(s.junkList)) + ((rnd.Float64() * 0.6) - 0.3)
		percent = basics.FloatClamp(percent, 0, 1)
		j.Init(j.GetImageFilepath())
		s.physSpace.Add(j.GetPhysObj())

		x := (rnd.Float64() * (s.spawnZone.Width))
		x += s.spawnZone.X - float64(j.GetPhysObj().X)
		y := (percent * (s.spawnZone.Height))
		y += s.spawnZone.Y - float64(j.GetPhysObj().Y)

		j.SetPosition(basics.Vector2f{X: x, Y: y})
		s.entityManager.AddEntity(&j)
		junkLookup[j.GetPhysObj()] = &j
	}

	// Create magnet
	m := &entities.MagnetObject{}
	m.Init("images/magnet.png")
	m.SetJunkLookup(junkLookup)
	m.SetUIPos(s.UIPosition)
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
	r.SetMagnetOffset(m.GetMagnetOffset())
	s.entityManager.AddEntity(r)
	p.SetFishingRodEndPoint(r.GetTip())

	m.SetMagnetStartPos(p.GetFishingRodEndPoint())

	p.Update(0)
	r.Update(0)
	m.Update(0)
	r.Update(0)

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
		s.countdownTimer = 0
		o := &OverworldScene{}
		state.SceneManager.GoTo(o, transitionTime)
	}

	if s.menuBtn {
		t := &TitleScene{}
		state.SceneManager.GoTo(t, transitionTime)
	}

	s.entityManager.RemoveDead(s.physSpace)

	return nil
}

func (s *ScavengeScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	if globals.Debug {
		ebitenutil.DrawRect(screen, s.spawnZone.X, s.spawnZone.Y, s.spawnZone.Width, s.spawnZone.Height, color.RGBA{128, 96, 64, 255})
	}

	timer := fmt.Sprintf("%.2f", s.countdownTimer)

	uipipeop := &ebiten.DrawImageOptions{}
	uipipeop.GeoM.Translate(globals.ScreenWidth-float64(s.UIPipeSprite.Bounds().Dx()), 0)
	screen.DrawImage(s.UIPipeSprite, uipipeop)

	s.entityManager.Draw(screen)
	s.txtRenderer.SetTarget(screen)

	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(globals.ScreenWidth-(float64(s.timerUIboxSprite.Bounds().Dx())+uiXOffset), uiYOffset)
	screen.DrawImage(s.timerUIboxSprite, options)

	if s.countdownTimer > textRedLimit {
		s.txtRenderer.SetColor(color.RGBA{197, 204, 184, 255})
	} else {
		s.txtRenderer.SetColor(color.RGBA{154, 79, 80, 255})
	}
	s.txtRenderer.Draw(timer, int(globals.ScreenWidth-(float64(s.timerUIboxSprite.Bounds().Dx())+uiXOffset)+textXOffset), int(uiYOffset+textYOffset))

	ebitenutil.DebugPrint(screen, fmt.Sprintf("cast available: %f", s.distanceOfOverworldCast))

	glassop := &ebiten.DrawImageOptions{}
	glassop.GeoM.Translate(globals.ScreenWidth-(float64(s.timerUIboxSprite.Bounds().Dx())+uiXOffset)+uiGlassXOffset, uiYOffset+uiGlassYOffset)
	screen.DrawImage(s.timerUIglassSprite, glassop)
}

func (s *ScavengeScene) InitJunkList() {
	cog := &entities.JunkObject{}
	cog.SetImageFilepath("images/cog.png")
	cog.InitData()
	cog.SetItemDataName("Cog")
	cog.SetItemDataDepthAndRarity(0, 80, 0.1)
	cog.AddItemDataMaterial("Iron", 5, 10)
	s.junkList = append(s.junkList, *cog)

	ironPipe := &entities.JunkObject{}
	ironPipe.SetImageFilepath("images/ironpipe.png")
	ironPipe.InitData()
	ironPipe.SetItemDataName("Iron Pipe")
	ironPipe.SetItemDataDepthAndRarity(1, 60, 0.2)
	ironPipe.AddItemDataMaterial("Iron", 15, 30)
	s.junkList = append(s.junkList, *ironPipe)

	tyre := &entities.JunkObject{}
	tyre.SetImageFilepath("images/tyre.png")
	tyre.InitData()
	tyre.SetItemDataName("Tyre")
	tyre.SetItemDataDepthAndRarity(2, 50, 0.2)
	tyre.AddItemDataMaterial("Rubber", 15, 25)
	tyre.AddItemDataMaterial("Iron", 15, 20)
	s.junkList = append(s.junkList, *tyre)

	steelBikeFrame := &entities.JunkObject{}
	steelBikeFrame.SetImageFilepath("images/steelbikeframe.png")
	steelBikeFrame.InitData()
	steelBikeFrame.SetItemDataName("Steel Bike Frame")
	steelBikeFrame.SetItemDataDepthAndRarity(3, 40, 0.3)
	steelBikeFrame.AddItemDataMaterial("Steel", 15, 20)
	steelBikeFrame.AddItemDataMaterial("Iron", 2, 7)
	steelBikeFrame.AddItemDataMaterial("Plastic", 0, 1)
	s.junkList = append(s.junkList, *steelBikeFrame)

	monitor := &entities.JunkObject{}
	monitor.SetImageFilepath("images/monitor.png")
	monitor.InitData()
	monitor.SetItemDataName("Monitor")
	monitor.SetItemDataDepthAndRarity(4, 30, 0.4)
	monitor.AddItemDataMaterial("Copper", 10, 15)
	monitor.AddItemDataMaterial("Plastic", 5, 10)
	monitor.AddItemDataMaterial("Iron", 0, 2)
	s.junkList = append(s.junkList, *monitor)

	toaster := &entities.JunkObject{}
	toaster.SetImageFilepath("images/toaster.png")
	toaster.InitData()
	toaster.SetItemDataName("Toaster")
	toaster.SetItemDataDepthAndRarity(5, 20, 0.5)
	toaster.AddItemDataMaterial("Nickel", 2, 10)
	toaster.AddItemDataMaterial("Iron", 1, 5)
	s.junkList = append(s.junkList, *toaster)

	steelPipe := &entities.JunkObject{}
	steelPipe.SetImageFilepath("images/steelpipe.png")
	steelPipe.InitData()
	steelPipe.SetItemDataName("Steel Pipe")
	steelPipe.SetItemDataDepthAndRarity(6, 18, 0.6)
	steelPipe.AddItemDataMaterial("Steel", 15, 30)
	s.junkList = append(s.junkList, *steelPipe)

	belt := &entities.JunkObject{}
	belt.SetImageFilepath("images/belt.png")
	belt.InitData()
	belt.SetItemDataName("Belt")
	belt.SetItemDataDepthAndRarity(7, 15, 0.7)
	belt.AddItemDataMaterial("Cobalt", 1, 7)
	belt.AddItemDataMaterial("Rubber", 3, 8)
	belt.AddItemDataMaterial("Plastic", 0, 2)
	s.junkList = append(s.junkList, *belt)

	copperPipe := &entities.JunkObject{}
	copperPipe.SetImageFilepath("images/copperpipe.png")
	copperPipe.InitData()
	copperPipe.SetItemDataName("Copper Pipe")
	copperPipe.SetItemDataDepthAndRarity(8, 13, 0.9)
	copperPipe.AddItemDataMaterial("Copper", 15, 30)
	s.junkList = append(s.junkList, *copperPipe)

	titaniumBikeFrame := &entities.JunkObject{}
	titaniumBikeFrame.SetImageFilepath("images/titaniumbikeframe.png")
	titaniumBikeFrame.InitData()
	titaniumBikeFrame.SetItemDataName("Titanium Bike Frame")
	titaniumBikeFrame.SetItemDataDepthAndRarity(9, 10, 1.2)
	titaniumBikeFrame.AddItemDataMaterial("Titanium", 15, 25)
	titaniumBikeFrame.AddItemDataMaterial("Iron", 0, 2)
	titaniumBikeFrame.AddItemDataMaterial("Plastic", 2, 4)
	s.junkList = append(s.junkList, *titaniumBikeFrame)

	titaniumPipe := &entities.JunkObject{}
	titaniumPipe.SetImageFilepath("images/titaniumpipe.png")
	titaniumPipe.InitData()
	titaniumPipe.SetItemDataName("Titanium Pipe")
	titaniumPipe.SetItemDataDepthAndRarity(10, 9, 1.4)
	titaniumPipe.AddItemDataMaterial("Titanium", 15, 30)
	s.junkList = append(s.junkList, *titaniumPipe)

	oldPC := &entities.JunkObject{}
	oldPC.SetImageFilepath("images/oldpc.png")
	oldPC.InitData()
	oldPC.SetItemDataName("Old PC")
	oldPC.SetItemDataDepthAndRarity(11, 8, 1.8)
	oldPC.AddItemDataMaterial("Copper", 10, 15)
	oldPC.AddItemDataMaterial("Plastic", 8, 12)
	oldPC.AddItemDataMaterial("Steel", 3, 7)
	oldPC.AddItemDataMaterial("Gold", 1, 3)
	s.junkList = append(s.junkList, *oldPC)

	battery := &entities.JunkObject{}
	battery.SetImageFilepath("images/battery.png")
	battery.InitData()
	battery.SetItemDataName("Battery")
	battery.SetItemDataDepthAndRarity(12, 4, 2.5)
	battery.AddItemDataMaterial("Cobalt", 10, 15)
	battery.AddItemDataMaterial("Nickel", 10, 15)
	s.junkList = append(s.junkList, *battery)
}

func (s *ScavengeScene) SelectJunk(castDistance float64) entities.JunkObject {

	var junkList []float64
	var totalRarity float64

	castPercent := (castDistance / globals.GetPlayerData().GetOverworldCastDistance()) * 100

	for _, v := range s.junkList {
		vRarity := v.GetItemData().GetRarity() * (castPercent * v.GetItemData().GetRarityScale())
		totalRarity += vRarity
		junkList = append(junkList, totalRarity)
	}

	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)

	num := rnd.Intn(int(totalRarity))

	currentChosen := 0
	for i := range junkList {
		if num < int(junkList[i]) {
			currentChosen = i
			break
		}
	}

	return s.junkList[currentChosen]
}
