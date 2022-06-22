package inventory

type Material struct {
	name string
}

func (m *Material) GetName() string {
	return m.name
}
