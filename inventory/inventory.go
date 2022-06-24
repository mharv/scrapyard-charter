package inventory

type Inventory struct {
	items     []Item
	materials map[string]int
}

func (i *Inventory) InitMaterials() {
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

func (i *Inventory) RemoveAllItemWithName(name string) {
	for index, element := range i.items {
		if element.name == name {
			i.items = append(i.items[:index], i.items[index+1:]...)
		}
	}
}

func (i *Inventory) RemoveOneItemWithName(name string) map[string]int {
	var salvagedMaterials map[string]int
	for index, element := range i.items {
		if element.name == name {
			salvagedMaterials = i.items[index].GetMaterials()
			i.items = append(i.items[:index], i.items[index+1:]...)
			break
		}
	}
	return salvagedMaterials
}

func (i *Inventory) SalvageOneItem(name string) {
	salvagedMaterials := i.RemoveOneItemWithName(name)
	for k, v := range salvagedMaterials {
		i.AddMaterial(k, v)
	}
}

// func (i *Inventory) RemoveMaterial(name string, quantity int) {
// }
