package ui

import (
	"fmt"
	"image/color"
	"log"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mharv/scrapyard-charter/basics"
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
	cursorPos              basics.Vector2f
	cursorClickPos         basics.Vector2f
	mouseClick             bool
	inventoryItems         []InventorySlotUi
}

func (u *Ui) Init() {
	u.xOffset = 50
	u.yOffset = 50
	u.characterOffset = 20

	u.itemsByCount = make(map[string]int)
	u.currentItemStoreLength = 0
	u.sortedItemKeys = []string{}
	u.inventoryItems = []InventorySlotUi{}

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

	x, y := ebiten.CursorPosition()
	u.cursorPos.X = float64(x)
	u.cursorPos.Y = float64(y)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		end := basics.Vector2f{X: u.cursorPos.X, Y: u.cursorPos.Y}
		u.cursorClickPos = end
		u.mouseClick = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		u.openButton = !u.openButton
	}
}

func (u *Ui) Update() error {

	if len(globals.GetPlayerData().GetInventory().GetItems()) != u.currentItemStoreLength {

		// should I just set everything to zero insted of making a new map?
		u.itemsByCount = make(map[string]int)

		for _, v := range globals.GetPlayerData().GetInventory().GetItems() {

			if val, ok := u.itemsByCount[v.GetName()]; ok {
				u.itemsByCount[v.GetName()] = val + 1
			} else {
				u.itemsByCount[v.GetName()] = 1
			}

		}

		u.sortedItemKeys = []string{}

		u.inventoryItems = []InventorySlotUi{}

		i := 1
		for k := range u.itemsByCount {
			u.sortedItemKeys = append(u.sortedItemKeys, k)
			tempInvItem := InventorySlotUi{}
			tempInvItem.InitSlot(70, 40, 40, 40, i, u.itemsByCount[k], k)
			u.inventoryItems = append(u.inventoryItems, tempInvItem)
			i++
		}
		sort.Strings(u.sortedItemKeys)

		u.currentItemStoreLength = len(globals.GetPlayerData().GetInventory().GetItems())
	}

	for _, v := range u.inventoryItems {
		if v.SalvageOneButton.IsClicked(u.cursorClickPos) && u.mouseClick {
			globals.GetPlayerData().GetInventory().SalvageOneItem(v.ItemName)
			u.mouseClick = false
		}
	}

	for _, v := range u.inventoryItems {
		if v.SalvageAllButton.IsClicked(u.cursorClickPos) && u.mouseClick {
			globals.GetPlayerData().GetInventory().SalvageOneItem(v.ItemName)
			u.mouseClick = false
		}
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

		// salvage button test
		for _, v := range u.inventoryItems {
			// draw salvageOne buttons
			buttonDrawColor := color.RGBA{240, 17, 17, 255}
			ebitenutil.DrawRect(screen, v.SalvageOneButton.X, v.SalvageOneButton.Y, v.SalvageOneButton.Width, v.SalvageOneButton.Height, buttonDrawColor)

			// draw salvageAll buttons
			buttonDrawColor = color.RGBA{12, 159, 7, 255}
			ebitenutil.DrawRect(screen, v.SalvageAllButton.X, v.SalvageAllButton.Y, v.SalvageAllButton.Width, v.SalvageAllButton.Height, buttonDrawColor)

			// draw button labels
			u.txtRenderer.Draw(fmt.Sprintf("%s", v.ItemName), int(v.SalvageOneButton.X), int(v.SalvageOneButton.Y))

			// draw button count and name
			u.txtRenderer.Draw(fmt.Sprintf("%d x %s", v.ItemCount, v.ItemName), int(v.X), int(v.Y))
		}

		for i, v := range globals.MaterialNamesList {
			tempVal := 0

			if val, ok := globals.GetPlayerData().GetInventory().GetMaterials()[v]; ok {
				tempVal = val
			}

			u.txtRenderer.Draw(fmt.Sprintf("%d x %s", tempVal, v), ((globals.ScreenWidth/3 + globals.ScreenWidth/3) + u.xOffset*2), (0 + u.yOffset*2 + u.characterOffset*(i+2)))
		}

		// u.txtRenderer.Draw("globals.GetPlayerData().GetInventory().GetItems()[0].GetName()", globals.ScreenWidth/2, globals.ScreenHeight/2)

	}

}
