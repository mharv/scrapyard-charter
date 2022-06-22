package ui

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/mharv/scrapyard-charter/inventory"
	"github.com/tinne26/etxt"
)

type Ui struct {
	openButton       bool
	xOffset, yOffset int
	txtRenderer      *etxt.Renderer
}

func (u *Ui) Init() {
	u.xOffset = 50
	u.yOffset = 50

	tempItem := inventory.Item{}

	tempItem.Init()
	tempItem.SetName("Iron Pipe")
	tempItem.AddRawMaterial("Iron", 25)

	globals.GetPlayerData().GetInventory().AddItem(tempItem)

	fontLib := etxt.NewFontLibrary()

	_, _, err := fontLib.ParseDirFonts("fonts")
	if err != nil {
		log.Fatal(err)
	}

	if !fontLib.HasFont("Rajdhani Regular") {
		log.Fatal("missing font Rajdhani-Regular.ttf")
	}

	u.txtRenderer = etxt.NewStdRenderer()
	glyphsCache := etxt.NewDefaultCache(10 * 1024 * 1024) // 10MB
	u.txtRenderer.SetCacheHandler(glyphsCache.NewHandler())
	u.txtRenderer.SetFont(fontLib.GetFont("Rajdhani Regular"))
	u.txtRenderer.SetAlign(etxt.YCenter, etxt.XCenter)
	u.txtRenderer.SetSizePx(24)
}

func (u *Ui) ReadInput() {

	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		u.openButton = !u.openButton
	}
}

func (u *Ui) Update() error {
	return nil
}

func (u *Ui) Draw(screen *ebiten.Image) {
	if u.openButton {
		screen.Fill(color.RGBA{0, 0, 0, 255})

		// draw item list

		drawColor := color.RGBA{222, 130, 22, 255}
		// inventory
		ebitenutil.DrawRect(screen, (0 + float64(u.xOffset)), (0 + float64(u.yOffset)), (globals.ScreenWidth/3)-float64(2*u.xOffset), globals.ScreenHeight-float64(u.yOffset*2), drawColor)
		// character equip
		ebitenutil.DrawRect(screen, (globals.ScreenWidth / 3), (0 + float64(u.yOffset)), (globals.ScreenWidth / 3), globals.ScreenHeight-float64(u.yOffset*2), drawColor)
		// materials
		ebitenutil.DrawRect(screen, (globals.ScreenWidth/3 + globals.ScreenWidth/3 + float64(u.xOffset)), (0 + float64(u.yOffset)), (globals.ScreenWidth/3)-float64(2*u.xOffset), globals.ScreenHeight-float64(u.yOffset*2), drawColor)

		// display inv items

		u.txtRenderer.SetTarget(screen)
		u.txtRenderer.SetColor(color.RGBA{255, 255, 255, 255})
		// u.txtRenderer.Draw("globals.GetPlayerData().GetInventory().GetItems()[0].GetName()", globals.ScreenWidth/2, globals.ScreenHeight/2)
		u.txtRenderer.Draw(globals.GetPlayerData().GetInventory().GetItems()[0].GetName(), globals.ScreenWidth/2, globals.ScreenHeight/2)

	}

}
