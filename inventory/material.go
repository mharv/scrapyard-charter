package inventory

type Material struct {
	name   string
	amount int
}

func (m *Material) GetMaterial() Material {
	return *m
}
