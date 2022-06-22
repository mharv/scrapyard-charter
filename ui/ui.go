package ui

import (
	"fmt"
	"image/color"
	"log"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/tinne26/etxt"
)

type Ui struct {
	openButton             bool
	xOffset, yOffset       int
	txtRenderer            *etxt.Renderer
	characterOffset        int
	itemsByCount           map[string]int
	currentItemStoreLength int
	sortedItemKeys         []string
}

func (u *Ui) Init() {
	u.xOffset = 50
	u.yOffset = 50
	u.characterOffset = 20

	u.itemsByCount = make(map[string]int)
	u.currentItemStoreLength = 0
	u.sortedItemKeys = []string{}

	// tempItem := inventory.Item{}

	// for i := 0; i < 10; i++ {
	// 	tempItem.Init()
	// 	tempItem.SetName("Iron Pipe")
	// 	tempItem.AddRawMaterial("Iron", 25)

	// 	globals.GetPlayerData().GetInventory().AddItem(tempItem)

	// }

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
	u.txtRenderer.SetAlign(etxt.Top, etxt.Left)
	u.txtRenderer.SetSizePx(u.characterOffset)
}

func (u *Ui) ReadInput() {

	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		u.openButton = !u.openButton
	}
}

func (u *Ui) Update() error {

	if len(globals.GetPlayerData().GetInventory().GetItems()) != u.currentItemStoreLength {
		for _, v := range globals.GetPlayerData().GetInventory().GetItems() {

			if val, ok := u.itemsByCount[v.GetName()]; ok {
				u.itemsByCount[v.GetName()] = val + 1
			} else {
				u.itemsByCount[v.GetName()] = 1
			}

		}

		for k := range u.itemsByCount {
			u.sortedItemKeys = append(u.sortedItemKeys, k)
		}
		sort.Strings(u.sortedItemKeys)

		u.currentItemStoreLength = len(globals.GetPlayerData().GetInventory().GetItems())
	}

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

		// render items
		// for i, v := range globals.GetPlayerData().GetInventory().GetItems() {

		// 	u.txtRenderer.Draw(v.GetName(), (0 + u.xOffset*2), (0 + u.yOffset*2 + u.characterOffset*(i+2)))
		// }
		for i, k := range u.sortedItemKeys {
			u.txtRenderer.Draw(fmt.Sprintf("%d x %s", u.itemsByCount[k], k), (0 + u.xOffset*2), (0 + u.yOffset*2 + u.characterOffset*(i+2)))
		}

		for i, v := range globals.MaterialNamesList {

			u.txtRenderer.Draw(v, ((globals.ScreenWidth/3 + globals.ScreenWidth/3) + u.xOffset*2), (0 + u.yOffset*2 + u.characterOffset*(i+2)))
		}

		// u.txtRenderer.Draw("globals.GetPlayerData().GetInventory().GetItems()[0].GetName()", globals.ScreenWidth/2, globals.ScreenHeight/2)

	}

}
