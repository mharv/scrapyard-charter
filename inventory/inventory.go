package inventory

type Inventory struct {
	keyItems           []KeyItem
	items              []Item
	materials          map[string]int
	NewKeyItemAcquired bool
}

func (i *Inventory) InitMaterials() {
	i.materials = make(map[string]int)
	i.materials["Iron"] = 0
	i.materials["Steel"] = 0
	i.materials["Copper"] = 0
	i.materials["Rubber"] = 0
	i.materials["Plastic"] = 0
	i.materials["Nickel"] = 0
	i.materials["Cobalt"] = 0
	i.materials["Titanium"] = 0
	i.materials["Gold"] = 0
}

func (i *Inventory) ResetMaterials() {
	i.materials["Iron"] = 0
	i.materials["Steel"] = 0
	i.materials["Copper"] = 0
	i.materials["Rubber"] = 0
	i.materials["Plastic"] = 0
	i.materials["Nickel"] = 0
	i.materials["Cobalt"] = 0
	i.materials["Titanium"] = 0
	i.materials["Gold"] = 0
}

func (i *Inventory) AddKeyItem(keyItem KeyItem) {
	i.keyItems = append(i.keyItems, keyItem)
	i.NewKeyItemAcquired = true
}

func (i *Inventory) GetKeyItems() []KeyItem {
	return i.keyItems
}

func (i *Inventory) GetKeyItemsByType(typeOfKeyItemRequired string) []KeyItem {
	keyItemsByType := []KeyItem{}
	for _, v := range i.keyItems {
		if v.keyItemType == typeOfKeyItemRequired {
			keyItemsByType = append(keyItemsByType, v)
		}
	}
	return keyItemsByType

}

func (i *Inventory) AddItem(item Item) {
	i.items = append(i.items, item)
}

func (i *Inventory) AddMaterial(name string, amount int) {
	i.materials[name] += amount
}

func (i *Inventory) RemoveMaterial(name string, amount int) {
	// need check in here to make sure enough materials when crafting
	i.materials[name] -= amount
}

func (i *Inventory) GetItems() []Item {
	return i.items
}

func (i *Inventory) GetMaterials() map[string]int {
	return i.materials
}

func (i *Inventory) RemoveAllItemWithName(name string) []map[string]RawMaterial {
	var salvagedMaterials []map[string]RawMaterial

	tempLength := len(i.items)
	for j := 0; j < tempLength; j++ {
		if i.items[j].name == name {
			salvagedMaterials = append(salvagedMaterials, i.items[j].GetMaterials())
			i.items = append(i.items[:j], i.items[j+1:]...)
			// when item is removed, reset loop variables
			j = 0
			tempLength = len(i.items)
		}
	}
	return salvagedMaterials
}

func (i *Inventory) RemoveOneItemWithName(name string) map[string]RawMaterial {
	var salvagedMaterials map[string]RawMaterial
	for index, element := range i.items {
		if element.name == name {
			salvagedMaterials = i.items[index].GetMaterials()
			i.items = append(i.items[:index], i.items[index+1:]...)
			return salvagedMaterials
		}
	}
	return salvagedMaterials
}

func (i *Inventory) SalvageOneItem(name string) {
	salvagedMaterials := i.RemoveOneItemWithName(name)
	for k, v := range salvagedMaterials {
		i.AddMaterial(k, v.GetAmount())
	}
}

func (i *Inventory) SalvageAllItems(name string) {
	salvagedMaterials := i.RemoveAllItemWithName(name)
	for _, material := range salvagedMaterials {
		for k, v := range material {
			i.AddMaterial(k, v.GetAmount())
		}
	}
}
