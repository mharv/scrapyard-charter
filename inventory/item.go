package inventory

type Item struct {
	name         string
	rawMaterials map[string]int
	modifiers    map[string]float64
}

//WHEN equip or remove an Item
//check modlist of all items, add that to the player data mod stuff
