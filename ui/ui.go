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
	inventoryItems         []InventorySlotUi
	craftButton            basics.FloatRectUI
	craftingBench          *crafting.CraftingBench
	inventoryBgSprite      *ebiten.Image
	equipmentBgSprite      *ebiten.Image
	materialsBgSprite      *ebiten.Image
	craftButtonPressed     *ebiten.Image
	craftButtonUnpressed   *ebiten.Image
	craftButtonUnavailable *ebiten.Image
	craftPressedCounter    float64
	rodEquip               EquippableSlot
	reelEquip              EquippableSlot
	lineEquip              EquippableSlot
	magEquip               EquippableSlot
	bootEquip              EquippableSlot
	elecEquip              EquippableSlot
	repEquip               EquippableSlot
}

const (
	invX, invY                             = 50, 50
	equX, equY                             = (globals.ScreenWidth / 3), 50
	matX, matY                             = (globals.ScreenWidth/3 + globals.ScreenWidth/3 + 50), 50
	headingOffsetX, headingOffsetY         = 30, 5
	matItemOffsetY                         = 55
	matItemListOffsetX, matItemListOffsetY = 32, 91
	matTextSize                            = 50
	cbX, cbY, cbW, cbH                     = 118, 590, 119, 72
	craftPressedDuration                   = 0.25
	rodX, rodY                             = 233, 165
	reelX, reelY                           = 171, 317
	lineX, lineY                           = 381, 299
	magX, magY                             = 350, 418
	bootX, bootY                           = 103, 534
	elecX, elecY                           = 29, 316
	repX, repY                             = 69, 243
	invSlotW, invSlotH                     = 62, 62
)

func (u *Ui) Init() {
	u.xOffset = 50
	u.yOffset = 50
	u.characterOffset = 20

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

	u.rodEquip = EquippableSlot{}
	u.reelEquip = EquippableSlot{}
	u.lineEquip = EquippableSlot{}
	u.magEquip = EquippableSlot{}
	u.bootEquip = EquippableSlot{}
	u.elecEquip = EquippableSlot{}
	u.repEquip = EquippableSlot{}

	u.rodEquip.InitEquibbaleSlot(equX+rodX, equY+rodY, invSlotW, invSlotH, "Rod")
	u.reelEquip.InitEquibbaleSlot(equX+reelX, equY+reelY, invSlotW, invSlotH, "Reel")
	u.lineEquip.InitEquibbaleSlot(equX+lineX, equY+lineY, invSlotW, invSlotH, "Line")
	u.magEquip.InitEquibbaleSlot(equX+magX, equY+magY, invSlotW, invSlotH, "Magnet")
	u.bootEquip.InitEquibbaleSlot(equX+bootX, equY+bootY, invSlotW, invSlotH, "Boot")
	u.elecEquip.InitEquibbaleSlot(equX+elecX, equY+elecY, invSlotW, invSlotH, "Electromagnet")
	u.repEquip.InitEquibbaleSlot(equX+repX, equY+repY, invSlotW, invSlotH, "Reverse Polarity")

	u.craftingBench = &crafting.CraftingBench{}
	u.craftingBench.Init()

	invbg, _, inverr := ebitenutil.NewImageFromFile("images/inventorypanel.png")
	if inverr != nil {
		log.Fatal(inverr)
	} else {
		u.inventoryBgSprite = invbg
	}

	equbg, _, equerr := ebitenutil.NewImageFromFile("images/inventorycenterpanel.png")
	if equerr != nil {
		log.Fatal(equerr)
	} else {
		u.equipmentBgSprite = equbg
	}

	matbg, _, materr := ebitenutil.NewImageFromFile("images/materialspanel.png")
	if materr != nil {
		log.Fatal(materr)
	} else {
		u.materialsBgSprite = matbg
	}

	cbu, _, cbuerr := ebitenutil.NewImageFromFile("images/craftButtonUnavailable.png")
	if cbuerr != nil {
		log.Fatal(cbuerr)
	} else {
		u.craftButtonUnavailable = cbu
	}

	cbup, _, cbuperr := ebitenutil.NewImageFromFile("images/craftButtonUnpressed.png")
	if cbuperr != nil {
		log.Fatal(cbuperr)
	} else {
		u.craftButtonUnpressed = cbup
	}

	cbp, _, cbperr := ebitenutil.NewImageFromFile("images/craftButtonPressed.png")
	if cbperr != nil {
		log.Fatal(cbperr)
	} else {
		u.craftButtonPressed = cbp
	}

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
			globals.GetPlayerData().GetInventory().SalvageAllItems(v.ItemName)
			u.mouseClick = false
		}
	}

	if u.rodEquip.OpenKeyItemListButton.IsClicked(u.cursorClickPos) && u.mouseClick {
		fmt.Printf("%s equipment slot has been pressed\n", u.rodEquip.ItemName)
		u.mouseClick = false
	}

	if u.lineEquip.OpenKeyItemListButton.IsClicked(u.cursorClickPos) && u.mouseClick {
		fmt.Printf("%s equipment slot has been pressed\n", u.lineEquip.ItemName)
		u.mouseClick = false
	}

	if u.reelEquip.OpenKeyItemListButton.IsClicked(u.cursorClickPos) && u.mouseClick {
		fmt.Printf("%s equipment slot has been pressed\n", u.reelEquip.ItemName)
		u.mouseClick = false
	}

	if u.magEquip.OpenKeyItemListButton.IsClicked(u.cursorClickPos) && u.mouseClick {
		fmt.Printf("%s equipment slot has been pressed\n", u.magEquip.ItemName)
		u.mouseClick = false
	}

	if u.bootEquip.OpenKeyItemListButton.IsClicked(u.cursorClickPos) && u.mouseClick {
		fmt.Printf("%s equipment slot has been pressed\n", u.bootEquip.ItemName)
		u.mouseClick = false
	}

	if u.rodEquip.OpenKeyItemListButton.IsClicked(u.cursorClickPos) && u.mouseClick {
		fmt.Printf("%s equipment slot has been pressed\n", u.rodEquip.ItemName)
		u.mouseClick = false
	}

	if u.elecEquip.OpenKeyItemListButton.IsClicked(u.cursorClickPos) && u.mouseClick {
		fmt.Printf("%s equipment slot has been pressed\n", u.elecEquip.ItemName)
		u.mouseClick = false
	}

	if u.repEquip.OpenKeyItemListButton.IsClicked(u.cursorClickPos) && u.mouseClick {
		fmt.Printf("%s equipment slot has been pressed\n", u.repEquip.ItemName)
		u.mouseClick = false
	}

	if u.craftButton.IsClicked(u.cursorClickPos) && globals.GetPlayerData().CheckIfInCraftZone() && u.mouseClick {
		u.craftPressedCounter = craftPressedDuration
		u.craftingBench.CraftItem()
		u.mouseClick = false
	}

	if u.craftPressedCounter > 0 {
		u.craftPressedCounter -= deltaTime
	}

	// fmt.Println(u.craftPressedCounter)

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
		for _, v := range u.inventoryItems {
			// draw salvageOne buttons
			buttonDrawColor := color.RGBA{240, 17, 17, 255}
			ebitenutil.DrawRect(screen, v.SalvageOneButton.X, v.SalvageOneButton.Y, v.SalvageOneButton.Width, v.SalvageOneButton.Height, buttonDrawColor)

			// draw salvageAll buttons
			buttonDrawColor = color.RGBA{12, 159, 7, 255}
			ebitenutil.DrawRect(screen, v.SalvageAllButton.X, v.SalvageAllButton.Y, v.SalvageAllButton.Width, v.SalvageAllButton.Height, buttonDrawColor)

			// draw button labels
			// u.txtRenderer.Draw(fmt.Sprintf("%s", v.ItemName), int(v.SalvageOneButton.X), int(v.SalvageOneButton.Y))

			// draw button count and name
			u.txtRenderer.Draw(fmt.Sprintf("%d x %s", v.ItemCount, v.ItemName), int(v.X), int(v.Y))
		}

		u.headingTxt.SetTarget(screen)
		u.headingTxt.Draw("INVENTORY", invX+headingOffsetX, invY+headingOffsetY)
		u.headingTxt.Draw("EQUIPMENT", equX+headingOffsetX, equY+headingOffsetY)
		u.headingTxt.Draw("MATERIALS", matX+headingOffsetX, matY+headingOffsetY)

		for _, v := range globals.GetPlayerData().GetInventory().GetKeyItems() {
			u.txtRenderer.Draw(fmt.Sprintf("%s", v.GetKeyItemName()), int(globals.ScreenWidth/3), int(0+float64(u.yOffset)))
		}

		for i, v := range globals.MaterialNamesList {
			tempVal := 0

			if val, ok := globals.GetPlayerData().GetInventory().GetMaterials()[v]; ok {
				tempVal = val
			}

			u.txtRenderer.SetSizePx(matTextSize)
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

		// for debugging equippable slot buttons
		// buttonDrawColor := color.RGBA{12, 159, 7, 255}
		// ebitenutil.DrawRect(screen, u.rodEquip.X, u.rodEquip.Y, u.rodEquip.Width, u.rodEquip.Height, buttonDrawColor)
		// ebitenutil.DrawRect(screen, u.reelEquip.X, u.reelEquip.Y, u.reelEquip.Width, u.reelEquip.Height, buttonDrawColor)
		// ebitenutil.DrawRect(screen, u.lineEquip.X, u.lineEquip.Y, u.lineEquip.Width, u.lineEquip.Height, buttonDrawColor)
		// ebitenutil.DrawRect(screen, u.magEquip.X, u.magEquip.Y, u.magEquip.Width, u.magEquip.Height, buttonDrawColor)
		// ebitenutil.DrawRect(screen, u.elecEquip.X, u.elecEquip.Y, u.elecEquip.Width, u.elecEquip.Height, buttonDrawColor)
		// ebitenutil.DrawRect(screen, u.repEquip.X, u.repEquip.Y, u.repEquip.Width, u.repEquip.Height, buttonDrawColor)
		// ebitenutil.DrawRect(screen, u.bootEquip.X, u.bootEquip.Y, u.bootEquip.Width, u.bootEquip.Height, buttonDrawColor)

		// u.txtRenderer.Draw("globals.GetPlayerData().GetInventory().GetItems()[0].GetName()", globals.ScreenWidth/2, globals.ScreenHeight/2)

	}

}
