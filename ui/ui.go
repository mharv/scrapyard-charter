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

	// init ui (which happens on overworld scene init) clears equipped items

	if globals.GetPlayerData().CheckKeyItemTypeSlotIfOccupied("Magnet") {

		u.magEquip.InitEquibbaleSlot(equX+magX, equY+magY, invSlotW, invSlotH, "Magnet")
		u.magEquip.KeyItem, _ = globals.GetPlayerData().GetEquippedItem("Magnet")
	} else {

		u.magEquip = EquippableSlot{}
		u.magEquip.InitEquibbaleSlot(equX+magX, equY+magY, invSlotW, invSlotH, "Magnet")
	}

	if globals.GetPlayerData().CheckKeyItemTypeSlotIfOccupied("Rod") {

		u.rodEquip.InitEquibbaleSlot(equX+rodX, equY+rodY, invSlotW, invSlotH, "Rod")
		u.rodEquip.KeyItem, _ = globals.GetPlayerData().GetEquippedItem("Rod")
	} else {

		u.rodEquip = EquippableSlot{}
		u.rodEquip.InitEquibbaleSlot(equX+rodX, equY+rodY, invSlotW, invSlotH, "Rod")
	}

	if globals.GetPlayerData().CheckKeyItemTypeSlotIfOccupied("Reel") {

		u.reelEquip.InitEquibbaleSlot(equX+reelX, equY+reelY, invSlotW, invSlotH, "Reel")
		u.reelEquip.KeyItem, _ = globals.GetPlayerData().GetEquippedItem("Reel")
	} else {

		u.reelEquip = EquippableSlot{}
		u.reelEquip.InitEquibbaleSlot(equX+reelX, equY+reelY, invSlotW, invSlotH, "Reel")
	}

	if globals.GetPlayerData().CheckKeyItemTypeSlotIfOccupied("Line") {

		u.lineEquip.InitEquibbaleSlot(equX+lineX, equY+lineY, invSlotW, invSlotH, "Line")
		u.lineEquip.KeyItem, _ = globals.GetPlayerData().GetEquippedItem("Line")
	} else {

		u.lineEquip = EquippableSlot{}
		u.lineEquip.InitEquibbaleSlot(equX+lineX, equY+lineY, invSlotW, invSlotH, "Line")
	}

	if globals.GetPlayerData().CheckKeyItemTypeSlotIfOccupied("Boot") {

		u.bootEquip.InitEquibbaleSlot(equX+bootX, equY+bootY, invSlotW, invSlotH, "Boot")
		u.bootEquip.KeyItem, _ = globals.GetPlayerData().GetEquippedItem("Boot")
	} else {

		u.bootEquip = EquippableSlot{}
		u.bootEquip.InitEquibbaleSlot(equX+bootX, equY+bootY, invSlotW, invSlotH, "Boot")
	}

	if globals.GetPlayerData().CheckKeyItemTypeSlotIfOccupied("Electromagnet") {

		u.elecEquip.InitEquibbaleSlot(equX+elecX, equY+elecY, invSlotW, invSlotH, "Electromagnet")
		u.elecEquip.KeyItem, _ = globals.GetPlayerData().GetEquippedItem("Electromagnet")
	} else {

		u.elecEquip = EquippableSlot{}
		u.elecEquip.InitEquibbaleSlot(equX+elecX, equY+elecY, invSlotW, invSlotH, "Electromagnet")
	}

	if globals.GetPlayerData().CheckKeyItemTypeSlotIfOccupied("Repulsor") {

		u.repEquip.InitEquibbaleSlot(equX+repX, equY+repY, invSlotW, invSlotH, "Repulsor")
		u.repEquip.KeyItem, _ = globals.GetPlayerData().GetEquippedItem("Repulsor")
	} else {

		u.repEquip = EquippableSlot{}
		u.repEquip.InitEquibbaleSlot(equX+repX, equY+repY, invSlotW, invSlotH, "Repulsor")
	}

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
		if !u.openButton {
			globals.GetPlayerData().GetInventory().NewBootsAcquired = false
			globals.GetPlayerData().GetInventory().NewElecAcquired = false
			globals.GetPlayerData().GetInventory().NewLineAcquired = false
			globals.GetPlayerData().GetInventory().NewMagnetAcquired = false
			globals.GetPlayerData().GetInventory().NewReelAcquired = false
			globals.GetPlayerData().GetInventory().NewRepAcquired = false
			globals.GetPlayerData().GetInventory().NewRodAcquired = false
		}
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
		if v.SalvageOneButton.IsClicked(u.cursorClickPos) && u.mouseClick && u.openButton {
			globals.GetPlayerData().GetInventory().SalvageOneItem(v.ItemName)
			u.mouseClick = false
		}
	}

	for _, v := range u.inventoryItems {
		if v.SalvageAllButton.IsClicked(u.cursorClickPos) && u.mouseClick && u.openButton {
			globals.GetPlayerData().GetInventory().SalvageAllItems(v.ItemName)
			u.mouseClick = false
		}
	}

	if u.rodEquip.OpenKeyItemListButton.IsClicked(u.cursorClickPos) && u.mouseClick && u.openButton {
		fmt.Printf("%s equipment slot has been pressed\n", u.rodEquip.ItemName)

		equipKeyItem("Rod", &u.rodEquip)
		u.mouseClick = false
	}

	if u.lineEquip.OpenKeyItemListButton.IsClicked(u.cursorClickPos) && u.mouseClick && u.openButton {
		fmt.Printf("%s equipment slot has been pressed\n", u.lineEquip.ItemName)

		equipKeyItem("Line", &u.lineEquip)
		u.mouseClick = false
	}

	if u.reelEquip.OpenKeyItemListButton.IsClicked(u.cursorClickPos) && u.mouseClick && u.openButton {
		fmt.Printf("%s equipment slot has been pressed\n", u.reelEquip.ItemName)

		equipKeyItem("Reel", &u.reelEquip)
		u.mouseClick = false
	}

	if u.magEquip.OpenKeyItemListButton.IsClicked(u.cursorClickPos) && u.mouseClick && u.openButton {
		fmt.Printf("%s equipment slot has been pressed\n", u.magEquip.ItemName)
		fmt.Println(globals.GetPlayerData().GetInventory().GetKeyItemsByType("Magnet"))

		// do some logic here if rightmousebutton pressed to go back one.

		equipKeyItem("Magnet", &u.magEquip)

		u.mouseClick = false
	}

	if u.bootEquip.OpenKeyItemListButton.IsClicked(u.cursorClickPos) && u.mouseClick && u.openButton {
		fmt.Printf("%s equipment slot has been pressed\n", u.bootEquip.ItemName)

		equipKeyItem("Boots", &u.bootEquip)
		u.mouseClick = false
	}

	if u.elecEquip.OpenKeyItemListButton.IsClicked(u.cursorClickPos) && u.mouseClick && u.openButton {
		fmt.Printf("%s equipment slot has been pressed\n", u.elecEquip.ItemName)

		equipKeyItem("Electromagnet", &u.elecEquip)
		u.mouseClick = false
	}

	if u.repEquip.OpenKeyItemListButton.IsClicked(u.cursorClickPos) && u.mouseClick && u.openButton {
		fmt.Printf("%s equipment slot has been pressed\n", u.repEquip.ItemName)
		equipKeyItem("Repulsor", &u.repEquip)
		u.mouseClick = false
	}

	if u.craftButton.IsClicked(u.cursorClickPos) && globals.GetPlayerData().CheckIfInCraftZone() && u.mouseClick && u.openButton {
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

		// debuggin key items
		// for _, v := range globals.GetPlayerData().GetInventory().GetKeyItems() {
		// 	u.txtRenderer.Draw(fmt.Sprintf("%s", v.GetKeyItemName()), int(globals.ScreenWidth/3), int(0+float64(u.yOffset)))
		// }

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

		// draws key item image

		if globals.GetPlayerData().CheckKeyItemTypeSlotIfOccupied("Rod") {

			KeyItemImage := &ebiten.DrawImageOptions{}
			KeyItemImage.GeoM.Translate(u.rodEquip.X, u.rodEquip.Y)
			screen.DrawImage(u.rodEquip.KeyItem.GetKeyItemImage(), KeyItemImage)
		}

		if globals.GetPlayerData().CheckKeyItemTypeSlotIfOccupied("Reel") {

			KeyItemImage := &ebiten.DrawImageOptions{}
			KeyItemImage.GeoM.Translate(u.reelEquip.X, u.reelEquip.Y)
			screen.DrawImage(u.reelEquip.KeyItem.GetKeyItemImage(), KeyItemImage)
		}
		if globals.GetPlayerData().CheckKeyItemTypeSlotIfOccupied("Line") {

			KeyItemImage := &ebiten.DrawImageOptions{}
			KeyItemImage.GeoM.Translate(u.lineEquip.X, u.lineEquip.Y)
			screen.DrawImage(u.lineEquip.KeyItem.GetKeyItemImage(), KeyItemImage)
		}
		if globals.GetPlayerData().CheckKeyItemTypeSlotIfOccupied("Magnet") {

			KeyItemImage := &ebiten.DrawImageOptions{}
			KeyItemImage.GeoM.Translate(u.magEquip.X, u.magEquip.Y)
			screen.DrawImage(u.magEquip.KeyItem.GetKeyItemImage(), KeyItemImage)
		}

		if globals.GetPlayerData().CheckKeyItemTypeSlotIfOccupied("Electromagnet") {

			KeyItemImage := &ebiten.DrawImageOptions{}
			KeyItemImage.GeoM.Translate(u.elecEquip.X, u.elecEquip.Y)
			screen.DrawImage(u.elecEquip.KeyItem.GetKeyItemImage(), KeyItemImage)
		}

		if globals.GetPlayerData().CheckKeyItemTypeSlotIfOccupied("Repulsor") {

			KeyItemImage := &ebiten.DrawImageOptions{}
			KeyItemImage.GeoM.Translate(u.repEquip.X, u.repEquip.Y)
			screen.DrawImage(u.repEquip.KeyItem.GetKeyItemImage(), KeyItemImage)
		}

		if globals.GetPlayerData().CheckKeyItemTypeSlotIfOccupied("Boots") {

			KeyItemImage := &ebiten.DrawImageOptions{}
			KeyItemImage.GeoM.Translate(u.bootEquip.X, u.bootEquip.Y)
			screen.DrawImage(u.bootEquip.KeyItem.GetKeyItemImage(), KeyItemImage)
		}

		// Draws the new key item indicators

		if globals.GetPlayerData().GetInventory().NewMagnetAcquired {

			indicatorDrawColor := color.RGBA{255, 100, 0, 255}
			// should show u.magEquip.KeyItem.Image
			ebitenutil.DrawRect(screen, u.magEquip.X, u.magEquip.Y, 8, 8, indicatorDrawColor)
		}
		if globals.GetPlayerData().GetInventory().NewRodAcquired {

			indicatorDrawColor := color.RGBA{255, 100, 0, 255}
			ebitenutil.DrawRect(screen, u.rodEquip.X, u.rodEquip.Y, 8, 8, indicatorDrawColor)
		}
		if globals.GetPlayerData().GetInventory().NewReelAcquired {

			indicatorDrawColor := color.RGBA{255, 100, 0, 255}
			ebitenutil.DrawRect(screen, u.reelEquip.X, u.reelEquip.Y, 8, 8, indicatorDrawColor)
		}
		if globals.GetPlayerData().GetInventory().NewLineAcquired {

			indicatorDrawColor := color.RGBA{255, 100, 0, 255}
			ebitenutil.DrawRect(screen, u.lineEquip.X, u.lineEquip.Y, 8, 8, indicatorDrawColor)
		}
		if globals.GetPlayerData().GetInventory().NewElecAcquired {

			indicatorDrawColor := color.RGBA{255, 100, 0, 255}
			ebitenutil.DrawRect(screen, u.elecEquip.X, u.elecEquip.Y, 8, 8, indicatorDrawColor)
		}
		if globals.GetPlayerData().GetInventory().NewRepAcquired {

			indicatorDrawColor := color.RGBA{255, 100, 0, 255}
			ebitenutil.DrawRect(screen, u.repEquip.X, u.repEquip.Y, 8, 8, indicatorDrawColor)
		}
		if globals.GetPlayerData().GetInventory().NewBootsAcquired {

			indicatorDrawColor := color.RGBA{255, 100, 0, 255}
			ebitenutil.DrawRect(screen, u.bootEquip.X, u.bootEquip.Y, 8, 8, indicatorDrawColor)
		}

		// draws the Hover info for key items

		if u.openButton {

			if u.reelEquip.OpenKeyItemListButton.IsHoveredOver(u.cursorPos) {

				drawHover("Reel", screen, &u.reelEquip, u.cursorPos, u.txtRenderer)

			}
			if u.rodEquip.OpenKeyItemListButton.IsHoveredOver(u.cursorPos) {

				drawHover("Rod", screen, &u.rodEquip, u.cursorPos, u.txtRenderer)
			}
			if u.lineEquip.OpenKeyItemListButton.IsHoveredOver(u.cursorPos) {

				drawHover("Line", screen, &u.lineEquip, u.cursorPos, u.txtRenderer)
			}
			if u.magEquip.OpenKeyItemListButton.IsHoveredOver(u.cursorPos) {

				// key item stats display on hover
				drawHover("Magnet", screen, &u.magEquip, u.cursorPos, u.txtRenderer)

			}
			if u.bootEquip.OpenKeyItemListButton.IsHoveredOver(u.cursorPos) {

				drawHover("Boots", screen, &u.bootEquip, u.cursorPos, u.txtRenderer)
			}
			if u.elecEquip.OpenKeyItemListButton.IsHoveredOver(u.cursorPos) {

				drawHover("Electromagnet", screen, &u.elecEquip, u.cursorPos, u.txtRenderer)
			}
			if u.repEquip.OpenKeyItemListButton.IsHoveredOver(u.cursorPos) {

				drawHover("Repulsor", screen, &u.repEquip, u.cursorPos, u.txtRenderer)
			}

		}

	}

}

func drawHover(keyItemType string, screen *ebiten.Image, slot *EquippableSlot, cursorPosition basics.Vector2f, txtRenderer *etxt.Renderer) {
	if globals.GetPlayerData().CheckKeyItemTypeSlotIfOccupied(keyItemType) {

		buttonDrawColor := color.RGBA{12, 159, 7, 255}
		ebitenutil.DrawRect(screen, cursorPosition.X+10, cursorPosition.Y-100, 100, 100, buttonDrawColor)
		// draw u.magEquip key item name and modifier here
		txtRenderer.Draw(fmt.Sprintf("%s", slot.KeyItem.GetKeyItemName()), int(cursorPosition.X+10), int(cursorPosition.Y-100))
		if len(globals.GetPlayerData().GetInventory().GetKeyItemsByType(keyItemType)) > 0 {
			txtRenderer.Draw(
				fmt.Sprintf(
					"%d/%d",
					globals.GetPlayerData().GetIndexOfEquippedKeyItem(slot.KeyItem)+1,
					len(globals.GetPlayerData().GetInventory().GetKeyItemsByType(keyItemType)),
				),
				int(cursorPosition.X+10),
				int(cursorPosition.Y-70),
			)
			txtRenderer.Draw(
				fmt.Sprintf(
					"%s: %0.f", slot.KeyItem.GetKeyItemModifiers().ModifierName,
					slot.KeyItem.GetKeyItemModifiers().ModifierValue,
				),
				int(cursorPosition.X+10),
				int(cursorPosition.Y-30),
			)
		}
	}
}

func equipKeyItem(keyItemType string, slot *EquippableSlot) {

	if globals.GetPlayerData().CheckKeyItemTypeSlotIfOccupied(keyItemType) {
		getIndexOfEquippedItem := globals.GetPlayerData().GetIndexOfEquippedKeyItem(slot.KeyItem)
		if getIndexOfEquippedItem != -1 {
			keyItemTypeListLength := len(globals.GetPlayerData().GetInventory().GetKeyItemsByType(keyItemType))

			if keyItemTypeListLength-1 == getIndexOfEquippedItem {
				slot.KeyItem = globals.GetPlayerData().GetInventory().GetKeyItemsByType(keyItemType)[0]
				globals.GetPlayerData().EquipItem(slot.KeyItem)
			} else {
				slot.KeyItem = globals.GetPlayerData().GetInventory().GetKeyItemsByType(keyItemType)[getIndexOfEquippedItem+1]
				globals.GetPlayerData().EquipItem(slot.KeyItem)
			}
		} else {
			fmt.Println("EQUIPPED ITEM NOT IN KEYITEMTYPE LIST")
			fmt.Println(slot.KeyItem)
			// seems like when a new item is added to the keyitem list, the equipped item gets reset to nothing.
			// this makes the program land in this code block. Quick fix below...
			if len(globals.GetPlayerData().GetInventory().GetKeyItemsByType(keyItemType)) > 0 {
				slot.KeyItem = globals.GetPlayerData().GetInventory().GetKeyItemsByType(keyItemType)[0]
				globals.GetPlayerData().EquipItem(slot.KeyItem)
			}
		}
	} else {
		if len(globals.GetPlayerData().GetInventory().GetKeyItemsByType(keyItemType)) > 0 {
			slot.KeyItem = globals.GetPlayerData().GetInventory().GetKeyItemsByType(keyItemType)[0]
			globals.GetPlayerData().EquipItem(slot.KeyItem)
		}
	}
}
