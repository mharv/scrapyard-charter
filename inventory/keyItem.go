package inventory

type KeyItem struct {
	name string
	// keyItemTypeIndex          int
	keyItemType               string
	modifiers                 map[string]float64
	materialsRequiredForCraft map[string]float64
}

func (k *KeyItem) GetKeyItemName() string {
	return k.name
}

func (k *KeyItem) GetKeyItemType() string {
	return k.keyItemType
}

// func (k *KeyItem) GetKeyItemTypeIndex() int {
// 	return k.keyItemTypeIndex
// }

// func (k *KeyItem) SetKeyItemTypeIndex(index int) {
// 	k.keyItemTypeIndex = index
// }

func (k *KeyItem) GetCraftingRecipe() map[string]float64 {
	return k.materialsRequiredForCraft
}

func (k *KeyItem) Init(name, keyItemType string, modifiers, materialsRequiredForCraft map[string]float64) {
	k.name = name
	k.keyItemType = keyItemType
	k.modifiers = modifiers
	k.materialsRequiredForCraft = materialsRequiredForCraft
}
