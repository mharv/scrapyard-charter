package ui

import (
	"github.com/mharv/scrapyard-charter/basics"
)

type InventorySlotUi struct {
	X                 float64
	Y                 float64
	Width             float64
	Height            float64
	SalvageOneButton  basics.FloatRectUI
	SalvageAllButton  basics.FloatRectUI
	SalvageOnePressed bool
	SalvageAllPressed bool
	ItemCount         int
	ItemName          string
	SlotNumber        int
}

func (i *InventorySlotUi) InitSlot(invX, invY, invW, invH, titleOffset, salvageButtonOffsetX, salvageButtonOffsetY float64, slotNo, itemCount int, itemName string) {
	i.X = invX
	i.Y = invY + (float64(slotNo) * titleOffset)
	i.Width = invW
	i.Height = invH
	i.ItemCount = itemCount
	i.ItemName = itemName
	i.SalvageOnePressed = false
	i.SalvageAllPressed = false

	i.SalvageOneButton = basics.FloatRectUI{
		Name:   itemName,
		X:      salvageButtonOffsetX - invW + ((invW / 2) - 4),
		Y:      i.Y,
		Height: invH,
		Width:  invW,
	}
	i.SalvageAllButton = basics.FloatRectUI{
		Name:   itemName,
		X:      salvageButtonOffsetX + (invW / 2) - 3,
		Y:      i.Y,
		Height: invH,
		Width:  invW,
	}
}
