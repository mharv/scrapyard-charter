package crafting

import (
	"math/rand"

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
	}

}

func (cb *CraftingBench) Init() {
	// initialize list of key items that can be crafted here
	// cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, TEMPITEM)

	NEWMAGNET := &inventory.KeyItem{}
	NEWMAGNET.Init(
		"NEW MAGNET",
		map[string]float64{"castSpeed": 350},
		map[string]float64{"Iron": 100, "Steel": 100},
	)

	cb.KeyItemsAvailable = append(cb.KeyItemsAvailable, *NEWMAGNET)

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
