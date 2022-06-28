package crafting

import (
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/mharv/scrapyard-charter/inventory"
)

type CraftingBench struct {
	MaterialsRequiredList map[string]inventory.Material
	KeyItemsAvailable     []inventory.KeyItem
}

func (cb *CraftingBench) CraftItem() {
	keyItemsObtainedCheck := globals.GetPlayerData().GetInventory().GetKeyItems()
	materialsCollected := globals.GetPlayerData().GetInventory().GetMaterials()

	// use to store key items pool to randomly pick from
	tempKeyItems := []inventory.KeyItem{}

	// for all possible items, do a check if enough materials exist
	// and that the player has not already crafted them
	// if both pass, add them to a list that one is randomly chosen from
	// and added to the players key item inventory.

	for _, kia := range cb.KeyItemsAvailable {
		// check if required materials are found
		tempRecipe := kia.GetCraftingRecipe()

		canCraft := false
		countOfChecks := 0
		checksToPass := len(tempRecipe)
		for key := range tempRecipe {
			if tempRecipe[key] <= float64(materialsCollected[key]) {
				countOfChecks += 1
			}
		}
		if countOfChecks == checksToPass {
			canCraft = true
		}

		// check if key item has already been crafted
		previouslyCrafted := false

		for _, kio := range keyItemsObtainedCheck {
			if kio.GetKeyItemName() == kia.GetKeyItemName() {
				previouslyCrafted = true
			}
		}

		if !previouslyCrafted && canCraft {
			tempKeyItems = append(tempKeyItems, kia)
		}
	}
	amountItemsCraftable := len(tempKeyItems)

	if amountItemsCraftable > 0 {
		randomIndex := rand.Intn(amountItemsCraftable)
		// remove all materials here
		globals.GetPlayerData().GetInventory().ResetMaterials()
		globals.GetPlayerData().GetInventory().AddKeyItem(tempKeyItems[randomIndex])

		// set new key item of type in
		switch tempKeyItems[randomIndex].GetKeyItemType() {
		case "Magnet":

			globals.GetPlayerData().GetInventory().NewMagnetAcquired = true
		case "Rod":

			globals.GetPlayerData().GetInventory().NewRodAcquired = true
		case "Reel":

			globals.GetPlayerData().GetInventory().NewReelAcquired = true
		case "Electromagnet":

			globals.GetPlayerData().GetInventory().NewElecAcquired = true
		case "Repulsor":

			globals.GetPlayerData().GetInventory().NewRepAcquired = true
		case "Boots":

			globals.GetPlayerData().GetInventory().NewBootsAcquired = true
		case "Line":
			globals.GetPlayerData().GetInventory().NewLineAcquired = true
		}
	}
}

func (cb *CraftingBench) Init() {
	// initialize list of key items that can be crafted here

	// MAGNETS
	magnet1 := &inventory.KeyItem{}
	magnet1.Init(
		"NEW MAGNET",
		"Magnet",
		inventory.KeyItemModifiers{ModifierName: "Magnet field size", ModifierValue: 350},
		map[string]float64{"Iron": 100, "Steel": 50},
		LoadImage("images/iconmagnet1.png"),
	)
	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *magnet1)

	magnet2 := &inventory.KeyItem{}
	magnet2.Init(
		"BETTER MAGNET",
		"Magnet",
		inventory.KeyItemModifiers{ModifierName: "Magnet field size", ModifierValue: 450},
		map[string]float64{"Rubber": 20, "Iron": 20},
		LoadImage("images/iconmagnet2.png"),
	)
	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *magnet2)

	magnet3 := &inventory.KeyItem{}
	magnet3.Init(
		"BETTER MAGNET",
		"Magnet",
		inventory.KeyItemModifiers{ModifierName: "Magnet field size", ModifierValue: 450},
		map[string]float64{"Rubber": 20, "Iron": 20},
		LoadImage("images/iconmagnet3.png"),
	)
	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *magnet3)

	goldenMagnet := &inventory.KeyItem{}
	goldenMagnet.Init(
		"GOLDENMAGNET",
		"Magnet",
		inventory.KeyItemModifiers{ModifierName: "Magnet field size", ModifierValue: 450},
		map[string]float64{"Gold": 100},
		LoadImage("images/iconmagnetgold.png"),
	)
	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *goldenMagnet)

	// DAS BOOTS

	boots1 := &inventory.KeyItem{}
	boots1.Init(
		"BOOTZ",
		"Boots",
		inventory.KeyItemModifiers{ModifierName: "Move speed", ModifierValue: 450},
		map[string]float64{"Rubber": 20, "Iron": 20},
		LoadImage("images/iconboots1.png"),
	)
	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *boots1)

	boots2 := &inventory.KeyItem{}
	boots2.Init(
		"BOOTZ",
		"Boots",
		inventory.KeyItemModifiers{ModifierName: "Move speed", ModifierValue: 450},
		map[string]float64{"Rubber": 20, "Iron": 20},
		LoadImage("images/iconboots2.png"),
	)
	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *boots2)

	boots3 := &inventory.KeyItem{}
	boots3.Init(
		"BOOTZ",
		"Boots",
		inventory.KeyItemModifiers{ModifierName: "Move speed", ModifierValue: 450},
		map[string]float64{"Rubber": 20, "Iron": 20},
		LoadImage("images/iconboots3.png"),
	)
	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *boots3)

	// RODS

	rod1 := &inventory.KeyItem{}
	rod1.Init(
		"RODGER",
		"Rod",
		inventory.KeyItemModifiers{ModifierName: "Cast speed", ModifierValue: 450},
		map[string]float64{"Steel": 20, "Iron": 20},
		LoadImage("images/iconrod1.png"),
	)
	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *rod1)

	rod2 := &inventory.KeyItem{}
	rod2.Init(
		"RODGER",
		"Rod",
		inventory.KeyItemModifiers{ModifierName: "Cast speed", ModifierValue: 450},
		map[string]float64{"Steel": 20, "Iron": 20},
		LoadImage("images/iconrod2.png"),
	)
	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *rod2)

	rod3 := &inventory.KeyItem{}
	rod3.Init(
		"RODGER",
		"Rod",
		inventory.KeyItemModifiers{ModifierName: "Cast speed", ModifierValue: 450},
		map[string]float64{"Steel": 20, "Iron": 20},
		LoadImage("images/iconrod3.png"),
	)
	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *rod3)

	// REELS

	reel1 := &inventory.KeyItem{}
	reel1.Init(
		"REELY",
		"Reel",
		inventory.KeyItemModifiers{ModifierName: "Reel Speed", ModifierValue: 450},
		map[string]float64{"Rubber": 20, "Steel": 20},
		LoadImage("images/iconreel1.png"),
	)
	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *reel1)

	reel2 := &inventory.KeyItem{}
	reel2.Init(
		"REELY",
		"Reel",
		inventory.KeyItemModifiers{ModifierName: "Reel Speed", ModifierValue: 450},
		map[string]float64{"Rubber": 20, "Steel": 20},
		LoadImage("images/iconreel2.png"),
	)
	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *reel2)

	reel3 := &inventory.KeyItem{}
	reel3.Init(
		"REELY",
		"Reel",
		inventory.KeyItemModifiers{ModifierName: "Reel Speed", ModifierValue: 450},
		map[string]float64{"Rubber": 20, "Steel": 20},
		LoadImage("images/iconreel3.png"),
	)
	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *reel3)

	// LINES
	line1 := &inventory.KeyItem{}
	line1.Init(
		"linedUP",
		"Line",
		inventory.KeyItemModifiers{ModifierName: "castSpeed", ModifierValue: 450},
		map[string]float64{"Rubber": 20, "Iron": 20},
		LoadImage("images/iconline1.png"),
	)
	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *line1)

	line2 := &inventory.KeyItem{}
	line2.Init(
		"linedUP",
		"Line",
		inventory.KeyItemModifiers{ModifierName: "castSpeed", ModifierValue: 450},
		map[string]float64{"Rubber": 20, "Iron": 20},
		LoadImage("images/iconline2.png"),
	)
	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *line2)

	line3 := &inventory.KeyItem{}
	line3.Init(
		"linedUP",
		"Line",
		inventory.KeyItemModifiers{ModifierName: "castSpeed", ModifierValue: 450},
		map[string]float64{"Rubber": 20, "Iron": 20},
		LoadImage("images/iconline3.png"),
	)
	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *line3)

	// ELECTROMAGNET

	Electrom := &inventory.KeyItem{}
	Electrom.Init(
		"Electrom",
		"Electromagnet",
		inventory.KeyItemModifiers{ModifierName: "Hold down 'Space'", ModifierValue: 1337},
		map[string]float64{"Rubber": 20, "Iron": 20},
		LoadImage("images/iconelectromagnet.png"),
	)
	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *Electrom)

	// REPULSOR

	revpol := &inventory.KeyItem{}
	revpol.Init(
		"revpol",
		"Repulsor",
		inventory.KeyItemModifiers{ModifierName: "Use with 'Tab'", ModifierValue: 420},
		map[string]float64{"Rubber": 20, "Iron": 20},
		LoadImage("images/iconrepulsor.png"),
	)
	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *revpol)
}

func LoadImage(filepath string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

// func (s *ScavengeScene) InitJunkList() {
// 	cog := &entities.JunkObject{}
// 	cog.SetImageFilepath("images/cog.png")
// 	cog.InitData()
// 	cog.SetItemDataName("Cog")
// 	cog.SetItemDataDepthAndRarity(0, 80, 0.1)
// 	cog.AddItemDataMaterial("Iron", 5, 10)
// 	s.junkList = append(s.junkList, *cog)

// 	ironPipe := &entities.JunkObject{}
// 	ironPipe.SetImageFilepath("images/ironpipe.png")
// 	ironPipe.InitData()
// 	ironPipe.SetItemDataName("Iron Pipe")
// 	ironPipe.SetItemDataDepthAndRarity(1, 60, 0.2)
// 	ironPipe.AddItemDataMaterial("Iron", 15, 30)
// 	s.junkList = append(s.junkList, *ironPipe)
// }
