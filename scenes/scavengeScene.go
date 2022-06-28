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
	background              *ebiten.Image
	bgwalls                 *ebiten.Image
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
	defaultTimerStart   = 30
	uiXOffset           = 32
	uiYOffset           = 32
	uiGlassXOffset      = 25
	uiGlassYOffset      = 4
	textXOffset         = 35
	textYOffset         = 0
	fontSize            = 50
	iconXOffset         = 184
	iconYOffset         = 66
	textRedLimit        = 10.0
)

func (s *ScavengeScene) Init() {
	s.physSpace = resolv.NewSpace(globals.ScreenWidth, globals.ScreenHeight, 16, 16)
	s.UIPosition = basics.Vector2f{X: globals.ScreenWidth - (uiXOffset + iconXOffset), Y: uiYOffset + iconYOffset}

	bg, _, bgerr := ebitenutil.NewImageFromFile("images/scavbackground.png")
	if bgerr != nil {
		log.Fatal(bgerr)
	}
	s.background = bg

	bgwall, _, bgwerr := ebitenutil.NewImageFromFile("images/bgwalls.png")
	if bgwerr != nil {
		log.Fatal(bgwerr)
	}
	s.bgwalls = bgwall

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

	// Create player
	p := &entities.ScavPlayerObject{}
	p.Init("images/player.png")
	s.physSpace.Add(p.GetPhysObj())
	p.SetPosition(basics.Vector2f{X: s.spawnZone.X, Y: (s.spawnZone.Y - p.GetPhysObj().H)})
	s.entityManager.AddEntity(p)

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
	p.SetMagnet(m)

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
	globals.GetAudioPlayer().PlayFile("audio/scavenge.mp3")

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

	bgop := &ebiten.DrawImageOptions{}
	bgop.GeoM.Translate(0, 0)
	screen.DrawImage(s.background, bgop)

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

	bgwallop := &ebiten.DrawImageOptions{}
	bgwallop.GeoM.Translate(0, 0)
	screen.DrawImage(s.bgwalls, bgwallop)

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
	cog.AddAudioFile("audio/mom1.mp3")
	cog.AddAudioFile("audio/mom2.mp3")
	cog.AddAudioFile("audio/mom3.mp3")
	cog.AddAudioFile("audio/mom4.mp3")
	cog.AddAudioFile("audio/mom12.mp3")
	s.junkList = append(s.junkList, *cog)

	ironPipe := &entities.JunkObject{}
	ironPipe.SetImageFilepath("images/ironpipe.png")
	ironPipe.InitData()
	ironPipe.SetItemDataName("Iron Pipe")
	ironPipe.SetItemDataDepthAndRarity(1, 60, 0.2)
	ironPipe.AddItemDataMaterial("Iron", 15, 30)
	ironPipe.AddAudioFile("audio/mop1.mp3")
	ironPipe.AddAudioFile("audio/mop2.mp3")
	s.junkList = append(s.junkList, *ironPipe)

	tyre := &entities.JunkObject{}
	tyre.SetImageFilepath("images/tyre.png")
	tyre.InitData()
	tyre.SetItemDataName("Tyre")
	tyre.SetItemDataDepthAndRarity(2, 50, 0.2)
	tyre.AddItemDataMaterial("Rubber", 15, 25)
	tyre.AddItemDataMaterial("Iron", 15, 20)
	tyre.AddAudioFile("audio/mor1.mp3")
	tyre.AddAudioFile("audio/mor2.mp3")
	tyre.AddAudioFile("audio/mor3.mp3")
	tyre.AddAudioFile("audio/mor4.mp3")
	s.junkList = append(s.junkList, *tyre)

	steelBikeFrame := &entities.JunkObject{}
	steelBikeFrame.SetImageFilepath("images/steelbikeframe.png")
	steelBikeFrame.InitData()
	steelBikeFrame.SetItemDataName("Steel Bike Frame")
	steelBikeFrame.SetItemDataDepthAndRarity(3, 40, 0.3)
	steelBikeFrame.AddItemDataMaterial("Steel", 15, 20)
	steelBikeFrame.AddItemDataMaterial("Iron", 2, 7)
	steelBikeFrame.AddItemDataMaterial("Plastic", 0, 1)
	steelBikeFrame.AddAudioFile("audio/mom5.mp3")
	steelBikeFrame.AddAudioFile("audio/mom6.mp3")
	steelBikeFrame.AddAudioFile("audio/mom7.mp3")
	steelBikeFrame.AddAudioFile("audio/mom8.mp3")
	steelBikeFrame.AddAudioFile("audio/mom9.mp3")
	s.junkList = append(s.junkList, *steelBikeFrame)

	monitor := &entities.JunkObject{}
	monitor.SetImageFilepath("images/monitor.png")
	monitor.InitData()
	monitor.SetItemDataName("Monitor")
	monitor.SetItemDataDepthAndRarity(4, 30, 0.4)
	monitor.AddItemDataMaterial("Copper", 10, 15)
	monitor.AddItemDataMaterial("Plastic", 5, 10)
	monitor.AddItemDataMaterial("Iron", 0, 2)
	monitor.AddAudioFile("audio/mor5.mp3")
	monitor.AddAudioFile("audio/mor6.mp3")
	s.junkList = append(s.junkList, *monitor)

	toaster := &entities.JunkObject{}
	toaster.SetImageFilepath("images/toaster.png")
	toaster.InitData()
	toaster.SetItemDataName("Toaster")
	toaster.SetItemDataDepthAndRarity(5, 20, 0.5)
	toaster.AddItemDataMaterial("Nickel", 2, 10)
	toaster.AddItemDataMaterial("Iron", 1, 5)
	toaster.AddAudioFile("audio/mom8.mp3")
	toaster.AddAudioFile("audio/mom9.mp3")
	toaster.AddAudioFile("audio/mom10.mp3")
	toaster.AddAudioFile("audio/mom11.mp3")
	s.junkList = append(s.junkList, *toaster)

	steelPipe := &entities.JunkObject{}
	steelPipe.SetImageFilepath("images/steelpipe.png")
	steelPipe.InitData()
	steelPipe.SetItemDataName("Steel Pipe")
	steelPipe.SetItemDataDepthAndRarity(6, 18, 0.6)
	steelPipe.AddItemDataMaterial("Steel", 15, 30)
	steelPipe.AddAudioFile("audio/mop3.mp3")
	steelPipe.AddAudioFile("audio/mop4.mp3")
	s.junkList = append(s.junkList, *steelPipe)

	belt := &entities.JunkObject{}
	belt.SetImageFilepath("images/belt.png")
	belt.InitData()
	belt.SetItemDataName("Belt")
	belt.SetItemDataDepthAndRarity(7, 15, 0.7)
	belt.AddItemDataMaterial("Cobalt", 1, 7)
	belt.AddItemDataMaterial("Rubber", 3, 8)
	belt.AddItemDataMaterial("Plastic", 0, 2)
	belt.AddAudioFile("audio/belt1.mp3")
	belt.AddAudioFile("audio/belt1.mp3")
	s.junkList = append(s.junkList, *belt)

	copperPipe := &entities.JunkObject{}
	copperPipe.SetImageFilepath("images/copperpipe.png")
	copperPipe.InitData()
	copperPipe.SetItemDataName("Copper Pipe")
	copperPipe.SetItemDataDepthAndRarity(8, 13, 0.9)
	copperPipe.AddItemDataMaterial("Copper", 15, 30)
	copperPipe.AddAudioFile("audio/mop5.mp3")
	copperPipe.AddAudioFile("audio/mop6.mp3")
	s.junkList = append(s.junkList, *copperPipe)

	titaniumBikeFrame := &entities.JunkObject{}
	titaniumBikeFrame.SetImageFilepath("images/titaniumbikeframe.png")
	titaniumBikeFrame.InitData()
	titaniumBikeFrame.SetItemDataName("Titanium Bike Frame")
	titaniumBikeFrame.SetItemDataDepthAndRarity(9, 10, 1.2)
	titaniumBikeFrame.AddItemDataMaterial("Titanium", 15, 25)
	titaniumBikeFrame.AddItemDataMaterial("Iron", 0, 2)
	titaniumBikeFrame.AddItemDataMaterial("Plastic", 2, 4)
	titaniumBikeFrame.AddAudioFile("audio/mom9.mp3")
	titaniumBikeFrame.AddAudioFile("audio/mom10.mp3")
	s.junkList = append(s.junkList, *titaniumBikeFrame)

	titaniumPipe := &entities.JunkObject{}
	titaniumPipe.SetImageFilepath("images/titaniumpipe.png")
	titaniumPipe.InitData()
	titaniumPipe.SetItemDataName("Titanium Pipe")
	titaniumPipe.SetItemDataDepthAndRarity(10, 9, 1.4)
	titaniumPipe.AddItemDataMaterial("Titanium", 15, 30)
	titaniumPipe.AddAudioFile("audio/mop6.mp3")
	titaniumPipe.AddAudioFile("audio/mop7.mp3")
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
	oldPC.AddAudioFile("audio/mom7.mp3")
	oldPC.AddAudioFile("audio/mom8.mp3")
	s.junkList = append(s.junkList, *oldPC)

	battery := &entities.JunkObject{}
	battery.SetImageFilepath("images/battery.png")
	battery.InitData()
	battery.SetItemDataName("Battery")
	battery.SetItemDataDepthAndRarity(12, 4, 2.5)
	battery.AddItemDataMaterial("Cobalt", 10, 15)
	battery.AddItemDataMaterial("Nickel", 10, 15)
	battery.AddAudioFile("audio/belt1.mp3")
	battery.AddAudioFile("audio/belt1.mp3")
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
