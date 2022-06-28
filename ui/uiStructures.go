package ui

import (
	"github.com/mharv/scrapyard-charter/basics"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/mharv/scrapyard-charter/inventory"
)

type InventorySlotUi struct {
	X                float64
	Y                float64
	Width            float64
	Height           float64
	SalvageOneButton basics.FloatRectUI
	SalvageAllButton basics.FloatRectUI
	ItemCount        int
	ItemName         string
	SlotNumber       int
}

const (
	yOffset     = 10
	xOffset     = 10
	titleOffset = 100
)

var (
	salvageButtonOffset = globals.ScreenWidth/3.0 - (50 * 4) // width
)

func (i *InventorySlotUi) InitSlot(invX, invY, invW, invH float64, slotNo, itemCount int, itemName string) {
	i.X = invX
	i.Y = invY*float64(slotNo) + titleOffset
	i.Width = invW
	i.Height = invH
	i.ItemCount = itemCount
	i.ItemName = itemName

	i.SalvageOneButton = basics.FloatRectUI{
		Name:   itemName,
		X:      invX + salvageButtonOffset,
		Y:      i.Y,
		Height: invH,
		Width:  invW,
	}
	i.SalvageAllButton = basics.FloatRectUI{
		Name:   itemName,
		X:      invX + invW + salvageButtonOffset,
		Y:      i.Y,
		Height: invH,
		Width:  invW,
	}
}

// func (i *InventorySlotUi) GetPosition()

type EquippableSlot struct {
	X                     float64
	Y                     float64
	Width                 float64
	Height                float64
	OpenKeyItemListButton basics.FloatRectUI
	ItemName              string
	KeyItem               inventory.KeyItem
}

func (e *EquippableSlot) InitEquibbaleSlot(invX, invY, invW, invH float64, itemName string) {
	e.X = invX
	e.Y = invY
	e.Width = invW
	e.Height = invH
	e.ItemName = itemName

	e.OpenKeyItemListButton = basics.FloatRectUI{
		Name:   itemName,
		X:      e.X,
		Y:      e.Y,
		Height: e.Height,
		Width:  e.Width,
	}
}
