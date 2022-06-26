package inventory

type KeyItem struct {
	name                      string
	modifiers                 map[string]float64
	materialsRequiredForCraft map[string]float64
}

func (k *KeyItem) GetKeyItemName() string {
	return k.name
}

func (k *KeyItem) GetCraftingRecipe() map[string]float64 {
	return k.materialsRequiredForCraft
}

func (k *KeyItem) Init(name string, modifiers, materialsRequiredForCraft map[string]float64) {
	k.name = name
	k.modifiers = modifiers
	k.materialsRequiredForCraft = materialsRequiredForCraft
}