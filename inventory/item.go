package inventory

type Item struct {
	name         string
	rawMaterials map[string]int
	modifiers    map[string]float64
}

func (i *Item) Init() {
	i.name = ""
	i.rawMaterials = make(map[string]int)
	i.modifiers = make(map[string]float64)
}

func (i *Item) SetName(name string) {
	i.name = name
}

func (i *Item) GetName() string {
	return i.name
}

func (i *Item) AddRawMaterial(name string, quantity int) {
	if val, ok := i.rawMaterials[name]; ok {
		val += quantity
	} else {
		i.rawMaterials[name] = quantity
	}
}

func (i *Item) AddModifier(name string, amount float64) {
	if val, ok := i.modifiers[name]; ok {
		val += amount
	} else {
		i.modifiers[name] = amount
	}
}
