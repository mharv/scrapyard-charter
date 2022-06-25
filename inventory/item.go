package inventory

type Item struct {
	name                       string
	rawMaterials               map[string]RawMaterial
	modifiers                  map[string]float64
	depth, rarity, rarityScale float64
}

func (i *Item) Init() {
	i.name = ""
	i.rawMaterials = make(map[string]RawMaterial)
	i.modifiers = make(map[string]float64)
}

func (i *Item) SetName(name string) {
	i.name = name
}

func (i *Item) SetDepth(depth float64) {
	i.depth = depth
}

func (i *Item) SetRarity(rarity float64) {
	i.rarity = rarity
}

func (i *Item) SetRarityScale(rarityScale float64) {
	i.rarityScale = rarityScale
}

func (i *Item) GetName() string {
	return i.name
}

func (i *Item) GetDepth() float64 {
	return i.depth
}

func (i *Item) GetRarity() float64 {
	return i.rarity
}

func (i *Item) GetRarityScale() float64 {
	return i.rarityScale
}

func (i *Item) GetMaterials() map[string]RawMaterial {
	return i.rawMaterials
}

func (i *Item) AddRawMaterial(name string, min, max int) {
	r := RawMaterial{}
	r.SetMinAndMax(min, max)
	i.rawMaterials[name] = r
}

func (i *Item) AddModifier(name string, amount float64) {
	if val, ok := i.modifiers[name]; ok {
		val += amount
	} else {
		i.modifiers[name] = amount
	}
}
