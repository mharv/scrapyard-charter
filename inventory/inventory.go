package inventory

type Inventory struct {
	items     []Item
	materials []Material
}

func (i *Inventory) AddItem(item Item) {
	i.items = append(i.items, item)
}

func (i *Inventory) AddMaterial(material Material) {
	i.materials = append(i.materials, material)
}

// func (i *Inventory) RemoveAllItemWithName(name string) {
// 	for index, element := range i.items {
// 		if element.name == name {
// 			i.items = append(i.items[:index], i.items[index+1:]...)
// 		}
// 	}
// }

// func (i *Inventory) RemoveMaterial(name string, quantity int) {
// }
