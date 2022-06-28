package inventory

import "github.com/hajimehoshi/ebiten/v2"

type KeyItem struct {
	name string
	// keyItemTypeIndex          int
	keyItemType               string
	modifiers                 KeyItemModifiers
	materialsRequiredForCraft map[string]float64
	keyItemImage              *ebiten.Image
}

type KeyItemModifiers struct {
	ModifierName  string
	ModifierValue float64
}

func (k *KeyItem) GetKeyItemName() string {
	return k.name
}

func (k *KeyItem) GetKeyItemImage() *ebiten.Image {
	return k.keyItemImage
}

func (k *KeyItem) GetKeyItemType() string {
	return k.keyItemType
}

func (k *KeyItem) GetKeyItemModifiers() KeyItemModifiers {
	return k.modifiers
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

func (k *KeyItem) Init(name, keyItemType string, modifiers KeyItemModifiers, materialsRequiredForCraft map[string]float64, keyItemImage *ebiten.Image) {
	k.name = name
	k.keyItemType = keyItemType
	k.modifiers = modifiers
	k.materialsRequiredForCraft = materialsRequiredForCraft
	k.keyItemImage = keyItemImage
}
