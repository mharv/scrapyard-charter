package data

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/mharv/scrapyard-charter/basics"
	"github.com/mharv/scrapyard-charter/inventory"
)

type PlayerData struct {
	inventory *inventory.Inventory
	//ScavPlayerObject
	scavMoveSpeedModifier float64
	rodStartXModifier     float64
	rodStartYModifier     float64
	rodEndXModifier       float64
	rodEndYModifier       float64
	//MagnetObject
	dropReactivationTimerModifier float64
	magneticFieldSizeModifier     float64
	attractionStrengthModifier    float64
	lineLengthModifier            float64
	magnetCastSpeedModifier       float64
	magnetReelSpeedModifier       float64
	//overworldPlayer
	overworldMoveSpeedModifier    float64
	overworldCastDistanceModifier float64
	InitialOverworldPosition      basics.Vector2f
	worldSeed                     int
	overworldIsInCraftZone        bool
	// itemSlots
	reel   inventory.KeyItem
	rod    inventory.KeyItem
	line   inventory.KeyItem
	magnet inventory.KeyItem
	boots  inventory.KeyItem
	elec   inventory.KeyItem
	rep    inventory.KeyItem

	isReelEquipped   bool
	isRodEquipped    bool
	isLineEquipped   bool
	isMagnetEquipped bool
	isBootsEquipped  bool
	isElecEquipped   bool
	isRepEquipped    bool
}

const (
	//ScavPlayerObject
	initialScavMoveSpeed = 250
	initialRodEndX       = 200
	initialRodEndY       = 25
	initialRodStartX     = 82
	initialRodStartY     = 54
	//MagnetObject
	initialDropReactivationTimer = 1
	initialMagneticFieldSize     = 100
	initialAttractionStrength    = 1
	initialLineLength            = 400
	initialMagnetCastSpeed       = 350
	initialMagnetReelSpeed       = 400
	//overworldPlayer
	initialOverworldMoveSpeed    = 200
	initialOverworldCastDistance = 200
)

func (p *PlayerData) EquipItem(item inventory.KeyItem) {
	switch item.GetKeyItemType() {
	case "Reel":
		p.reel = item
		p.isReelEquipped = true
	case "Rod":
		p.rod = item
		p.isRodEquipped = true
	case "Line":
		p.line = item
		p.isLineEquipped = true
	case "Magnet":
		p.magnet = item
		p.isMagnetEquipped = true
	case "Boots":
		p.boots = item
		p.isBootsEquipped = true
	case "Repulsor":
		p.rep = item
		p.isRepEquipped = true
	case "Electromagnet":
		p.elec = item
		p.isElecEquipped = true
	}
}

func (p *PlayerData) GetEquippedItem(slotName string) (inventory.KeyItem, error) {
	switch slotName {
	case "Reel":
		return p.reel, nil
	case "Rod":
		return p.rod, nil
	case "Line":
		return p.line, nil
	case "Magnet":
		return p.magnet, nil
	case "Boots":
		return p.boots, nil
	case "Repulsor":
		return p.rep, nil
	case "Electromagnet":
		return p.elec, nil
	}

	return inventory.KeyItem{}, errors.New("slotName does not exist")
}

func (p *PlayerData) CheckKeyItemTypeSlotIfOccupied(keyItemType string) bool {
	switch keyItemType {
	case "Reel":
		return p.isReelEquipped
	case "Rod":
		return p.isRodEquipped
	case "Line":
		return p.isLineEquipped
	case "Magnet":
		return p.isMagnetEquipped
	case "Boots":
		return p.isBootsEquipped
	case "Repulsor":
		return p.isRepEquipped
	case "Electromagnet":
		return p.isElecEquipped
	}
	return false
}

func (p *PlayerData) GetIndexOfEquippedKeyItem(keyItem inventory.KeyItem) int {
	getIndexOfEquippedItem := -1
	for i, v := range p.inventory.GetKeyItemsByType(keyItem.GetKeyItemType()) {
		itemInSlot, err := p.GetEquippedItem(keyItem.GetKeyItemType())
		if err != nil {
			fmt.Println(err)
		} else {
			if v.GetKeyItemName() == itemInSlot.GetKeyItemName() {
				getIndexOfEquippedItem = i
				break
			}
		}
	}
	return getIndexOfEquippedItem
}

func (p *PlayerData) Init() {
	p.inventory = &inventory.Inventory{}
	p.inventory.InitMaterials()
	p.worldSeed = rand.Int()
}

// func (p *PlayerData) Update() error {
// 	if p.inventory.NewKeyItemAcquired {
// 		length := len(p.inventory.GetKeyItems())
// 		if length > 0 {
// 			if p.CheckKeyItemTypeSlotIfOccupied(p.inventory.GetKeyItems()[length-1].GetKeyItemType()) {
// 				p.EquipItem(p.inventory.GetKeyItems()[length-1])
// 			}
// 		}
// 		p.inventory.NewKeyItemAcquired = false
// 	}
// 	return nil
// }

func (p *PlayerData) CheckIfInCraftZone() bool {
	return p.overworldIsInCraftZone
}

func (p *PlayerData) SetInCraftZoneStatus(status bool) {
	p.overworldIsInCraftZone = status
}

func (p *PlayerData) GetWorldSeed() int {
	return p.worldSeed
}

func (p *PlayerData) GetPlayerPosition() basics.Vector2f {
	return p.InitialOverworldPosition
}

func (p *PlayerData) SetPlayerPosition(position basics.Vector2f) {
	p.InitialOverworldPosition = position
}
func (p *PlayerData) GetInventory() *inventory.Inventory {
	return p.inventory
}

func (p *PlayerData) GetOverworldCastDistance() float64 {
	return initialOverworldCastDistance + p.overworldCastDistanceModifier
}

func (p *PlayerData) GetOverworldMoveSpeed() float64 {
	return initialOverworldMoveSpeed + p.overworldMoveSpeedModifier
}

func (p *PlayerData) GetScavMoveSpeed() float64 {
	return initialScavMoveSpeed + p.scavMoveSpeedModifier
}

func (p *PlayerData) GetRodStartX() float64 {
	return initialRodStartX + p.rodStartXModifier
}

func (p *PlayerData) GetRodStartY() float64 {
	return initialRodStartY + p.rodStartYModifier
}

func (p *PlayerData) GetRodEndX() float64 {
	return initialRodEndX + p.rodEndXModifier
}

func (p *PlayerData) GetRodEndY() float64 {
	return initialRodEndY + p.rodEndYModifier
}

func (p *PlayerData) GetDropReactivationTimer() float64 {
	return initialDropReactivationTimer + p.dropReactivationTimerModifier
}

func (p *PlayerData) GetMagneticFieldSize() float64 {
	return initialMagneticFieldSize + p.magneticFieldSizeModifier
}

func (p *PlayerData) GetAttractionStrength() float64 {
	return initialAttractionStrength + p.attractionStrengthModifier
}

func (p *PlayerData) GetLineLength() float64 {
	return initialLineLength + p.lineLengthModifier
}

func (p *PlayerData) GetMagnetCastSpeed() float64 {
	return initialMagnetCastSpeed + p.magnetCastSpeedModifier
}

func (p *PlayerData) GetMagnetReelSpeed() float64 {
	return initialMagnetReelSpeed + p.magnetReelSpeedModifier
}

func (p *PlayerData) HasElectroMagnet() bool {
	return false
}

func (p *PlayerData) HasRepulsor() bool {
	return false
}
