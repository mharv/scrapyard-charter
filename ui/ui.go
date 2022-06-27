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
	"github.com/mharv/scrapyard-charter/crafting"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/tinne26/etxt"
)

type Ui struct {
	openButton             bool
	xOffset, yOffset       int
	txtRenderer            *etxt.Renderer
	headingTxt             *etxt.Renderer
	characterOffset        int
	itemsByCount           map[string]int
	currentItemStoreLength int
	sortedItemKeys         []string
	cursorPos              basics.Vector2f
	cursorClickPos         basics.Vector2f
	mouseClick             bool
	open                   bool
	inventoryItems         []InventorySlotUi
	craftButton            basics.FloatRectUI
	craftingBench          *crafting.CraftingBench
	inventoryBgSprite      *ebiten.Image
	equipmentBgSprite      *ebiten.Image
	materialsBgSprite      *ebiten.Image
	craftButtonPressed     *ebiten.Image
	craftButtonUnpressed   *ebiten.Image
	craftButtonUnavailable *ebiten.Image
	salvageOnePressed      *ebiten.Image
	salvageOneUnpressed    *ebiten.Image
	salvageAllPressed      *ebiten.Image
	salvageAllUnpressed    *ebiten.Image
	craftPressedCounter    float64
}

const (
	invX, invY                             = 50, 50
	equX, equY                             = (globals.ScreenWidth / 3), 50
	matX, matY                             = (globals.ScreenWidth/3 + globals.ScreenWidth/3 + 50), 50
	headingOffsetX, headingOffsetY         = 30, 5
	invItemOffsetY                         = 41
	invItemListOffsetX, invItemListOffsetY = 30, 100
	matItemOffsetY                         = 55
	matItemListOffsetX, matItemListOffsetY = 32, 91
	salvageOffsetX, salvageOffsetY         = 327, 97
	matTextSize, invTextSize               = 50, 25
	cbX, cbY, cbW, cbH                     = 118, 590, 119, 72
	craftPressedDuration                   = 0.25
	rodX, rodY                             = 233, 165
	reelX, reelY                           = 171, 317
	lineX, lineY                           = 381, 299
	magX, magY                             = 350, 418
	bootX, bootY                           = 103, 534
	elecX, elecY                           = 29, 361
	repX, repY                             = 69, 243
	invSlotW, invSlotH                     = 62, 62
	salvageSize                            = 36
)

func (u *Ui) IsOpen() bool {
	return u.open
}

func (u *Ui) Init() {
	u.xOffset = 50
	u.yOffset = 50
	u.characterOffset = 20
	u.open = false

	u.itemsByCount = make(map[string]int)
	u.currentItemStoreLength = 0
	u.sortedItemKeys = []string{}
	u.inventoryItems = []InventorySlotUi{}
	u.craftButton = basics.FloatRectUI{
		Name:   "CRAFT",
		X:      matX + cbX,
		Y:      matY + cbY,
		Width:  cbW,
		Height: cbH,
	}

	u.craftingBench = &crafting.CraftingBench{}
	u.craftingBench.Init()

	u.LoadImages()

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
	u.txtRenderer.SetColor(color.RGBA{197, 204, 184, 255})

	u.headingTxt = etxt.NewStdRenderer()
	u.headingTxt.SetCacheHandler(glyphsCache.NewHandler())
	u.headingTxt.SetFont(fontLib.GetFont("Rajdhani Regular"))
	u.headingTxt.SetAlign(etxt.Top, etxt.Left)
	u.headingTxt.SetSizePx(70)
	u.headingTxt.SetColor(color.RGBA{110, 105, 98, 255})
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
		u.open = !u.open
	}
}

func (u *Ui) Update(deltaTime float64) error {

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

		i := 0
		for k := range u.itemsByCount {
			u.sortedItemKeys = append(u.sortedItemKeys, k)
			tempInvItem := InventorySlotUi{}

			tempInvItem.InitSlot(invX+salvageOffsetX-(salvageSize*2), invY+salvageOffsetY, salvageSize, salvageSize, invItemOffsetY, salvageOffsetX, salvageOffsetY, i, u.itemsByCount[k], k)
			u.inventoryItems = append(u.inventoryItems, tempInvItem)
			i++
		}
		sort.Strings(u.sortedItemKeys)

		u.currentItemStoreLength = len(globals.GetPlayerData().GetInventory().GetItems())
	}

	for _, v := range u.inventoryItems {
		v.SalvageOnePressed = false
		if v.SalvageOneButton.IsClicked(u.cursorClickPos) && u.mouseClick {
			globals.GetPlayerData().GetInventory().SalvageOneItem(v.ItemName)
			v.SalvageOnePressed = true
			u.mouseClick = false
		}
	}

	for _, v := range u.inventoryItems {
		v.SalvageAllPressed = false
		if v.SalvageAllButton.IsClicked(u.cursorClickPos) && u.mouseClick {
			globals.GetPlayerData().GetInventory().SalvageAllItems(v.ItemName)
			v.SalvageAllPressed = true
			u.mouseClick = false
		}
	}

	if u.craftButton.IsClicked(u.cursorClickPos) && globals.GetPlayerData().CheckIfInCraftZone() && u.mouseClick {
		u.craftPressedCounter = craftPressedDuration
		u.craftingBench.CraftItem()
		u.mouseClick = false
	}

	if u.craftPressedCounter > 0 {
		u.craftPressedCounter -= deltaTime
	}

	return nil
}

func (u *Ui) Draw(screen *ebiten.Image) {
	if u.openButton {
		// inventory
		invop := &ebiten.DrawImageOptions{}
		invop.GeoM.Translate(invX, invY)
		screen.DrawImage(u.inventoryBgSprite, invop)

		// character equip
		equop := &ebiten.DrawImageOptions{}
		equop.GeoM.Translate(equX, equY)
		screen.DrawImage(u.equipmentBgSprite, equop)

		// materials
		matop := &ebiten.DrawImageOptions{}
		matop.GeoM.Translate(matX, matY)
		screen.DrawImage(u.materialsBgSprite, matop)

		// display inv items
		u.txtRenderer.SetTarget(screen)

		// salvage button test

		u.txtRenderer.SetSizePx(invTextSize)
		for i, v := range u.inventoryItems {
			// draw salvageOne buttons

			soop := &ebiten.DrawImageOptions{}
			soop.GeoM.Translate(invX+salvageOffsetX-(salvageSize*2), float64(invY+salvageOffsetY+(invItemOffsetY*i)))
			if v.SalvageOnePressed {
				screen.DrawImage(u.salvageOnePressed, soop)
			} else {
				screen.DrawImage(u.salvageOneUnpressed, soop)
			}

			saop := &ebiten.DrawImageOptions{}
			saop.GeoM.Translate(1+invX+salvageOffsetX-(salvageSize), float64(invY+salvageOffsetY+(invItemOffsetY*i)))
			if v.SalvageAllPressed {
				screen.DrawImage(u.salvageAllPressed, saop)
			} else {
				screen.DrawImage(u.salvageAllUnpressed, saop)
			}

			u.txtRenderer.Draw(fmt.Sprintf("%d x %s", v.ItemCount, v.ItemName), invX+invItemListOffsetX, invY+invItemListOffsetY+(invItemOffsetY*i))
		}

		u.headingTxt.SetTarget(screen)
		u.headingTxt.Draw("INVENTORY", invX+headingOffsetX, invY+headingOffsetY)
		u.headingTxt.Draw("EQUIPMENT", equX+headingOffsetX, equY+headingOffsetY)
		u.headingTxt.Draw("MATERIALS", matX+headingOffsetX, matY+headingOffsetY)

		for _, v := range globals.GetPlayerData().GetInventory().GetKeyItems() {
			u.txtRenderer.Draw(fmt.Sprintf("%s", v.GetKeyItemName()), int(globals.ScreenWidth/3), int(0+float64(u.yOffset)))
		}

		u.txtRenderer.SetSizePx(matTextSize)
		for i, v := range globals.MaterialNamesList {
			tempVal := 0

			if val, ok := globals.GetPlayerData().GetInventory().GetMaterials()[v]; ok {
				tempVal = val
			}

			u.txtRenderer.Draw(fmt.Sprintf("%d x %s", tempVal, v), (matX + matItemListOffsetX), (matY+matItemListOffsetY)+matItemOffsetY*(i))
		}

		cbop := &ebiten.DrawImageOptions{}
		cbop.GeoM.Translate(matX+cbX, matY+cbY)
		// draw craft button
		if globals.GetPlayerData().CheckIfInCraftZone() {
			if u.craftPressedCounter > 0 {
				screen.DrawImage(u.craftButtonPressed, cbop)
			} else {
				screen.DrawImage(u.craftButtonUnpressed, cbop)
			}

		} else {
			screen.DrawImage(u.craftButtonUnavailable, cbop)
		}
		if globals.Debug {
			buttonDrawColor := color.RGBA{12, 159, 7, 255}
			ebitenutil.DrawRect(
				screen,
				u.craftButton.X,
				u.craftButton.Y,
				u.craftButton.Width,
				u.craftButton.Height,
				buttonDrawColor,
			)
		}

		// u.txtRenderer.Draw("globals.GetPlayerData().GetInventory().GetItems()[0].GetName()", globals.ScreenWidth/2, globals.ScreenHeight/2)

	}

}

func (u *Ui) LoadImages() {
	u.inventoryBgSprite = LoadImage("images/inventorypanel.png")
	u.equipmentBgSprite = LoadImage("images/inventorycenterpanel.png")
	u.materialsBgSprite = LoadImage("images/materialspanel.png")
	u.craftButtonUnavailable = LoadImage("images/craftButtonUnavailable.png")
	u.craftButtonUnpressed = LoadImage("images/craftButtonUnpressed.png")
	u.craftButtonPressed = LoadImage("images/craftButtonPressed.png")
	u.salvageAllPressed = LoadImage("images/salvageallpressed.png")
	u.salvageAllUnpressed = LoadImage("images/salvageallunpressed.png")
	u.salvageOnePressed = LoadImage("images/salvageonepressed.png")
	u.salvageOneUnpressed = LoadImage("images/salvageoneunpressed.png")
}

func LoadImage(filepath string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	return img
}
