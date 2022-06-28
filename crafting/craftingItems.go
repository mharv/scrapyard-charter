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
	// cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, TEMPITEM)

	NEWMAGNET := &inventory.KeyItem{}
	NEWMAGNET.Init(
		"NEW MAGNET",
		"Magnet",
		inventory.KeyItemModifiers{ModifierName: "castSpeed", ModifierValue: 350},
		map[string]float64{"Iron": 100, "Steel": 50},
		LoadImage("images/iconmagnet1.png"),
	)

	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *NEWMAGNET)

	NEWMAGNET2 := &inventory.KeyItem{}
	NEWMAGNET2.Init(
		"BETTER MAGNET",
		"Magnet",
		inventory.KeyItemModifiers{ModifierName: "castSpeed", ModifierValue: 450},
		map[string]float64{"Rubber": 20, "Iron": 20},
		LoadImage("images/iconmagnet2.png"),
	)

	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *NEWMAGNET2)

	BOOTZ := &inventory.KeyItem{}
	BOOTZ.Init(
		"BOOTZ",
		"Boots",
		inventory.KeyItemModifiers{ModifierName: "castSpeed", ModifierValue: 450},
		map[string]float64{"Rubber": 20, "Iron": 20},
		LoadImage("images/iconboots1.png"),
	)

	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *BOOTZ)

	RODGER := &inventory.KeyItem{}
	RODGER.Init(
		"RODGER",
		"Rod",
		inventory.KeyItemModifiers{ModifierName: "castSpeed", ModifierValue: 450},
		map[string]float64{"Steel": 20, "Iron": 20},
		LoadImage("images/iconrod1.png"),
	)

	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *RODGER)

	REELY := &inventory.KeyItem{}
	REELY.Init(
		"REELY",
		"Reel",
		inventory.KeyItemModifiers{ModifierName: "castSpeed", ModifierValue: 450},
		map[string]float64{"Rubber": 20, "Steel": 20},
		LoadImage("images/iconreel1.png"),
	)

	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *REELY)

	Electrom := &inventory.KeyItem{}
	Electrom.Init(
		"Electrom",
		"Electromagnet",
		inventory.KeyItemModifiers{ModifierName: "castSpeed", ModifierValue: 450},
		map[string]float64{"Rubber": 20, "Iron": 20},
		LoadImage("images/iconelectromagnet.png"),
	)

	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *Electrom)

	revpol := &inventory.KeyItem{}
	revpol.Init(
		"revpol",
		"Repulsor",
		inventory.KeyItemModifiers{ModifierName: "castSpeed", ModifierValue: 450},
		map[string]float64{"Rubber": 20, "Iron": 20},
		LoadImage("images/iconrepulsor.png"),
	)

	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *revpol)

	linedUP := &inventory.KeyItem{}
	linedUP.Init(
		"linedUP",
		"Line",
		inventory.KeyItemModifiers{ModifierName: "castSpeed", ModifierValue: 450},
		map[string]float64{"Rubber": 20, "Iron": 20},
		LoadImage("images/iconline1.png"),
	)

	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *linedUP)

	GOLDENMAGNET := &inventory.KeyItem{}
	GOLDENMAGNET.Init(
		"GOLDENMAGNET",
		"Magnet",
		inventory.KeyItemModifiers{ModifierName: "castSpeed", ModifierValue: 450},
		map[string]float64{"Gold": 100},
		LoadImage("images/icongoldenmagnet.png"),
	)

	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *GOLDENMAGNET)
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
