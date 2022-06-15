package template

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mharv/scrapyard-charter/scenes"
)

type TemplateScene struct {
}

func (t *TemplateScene) Init() {
}

func (t *TemplateScene) ReadInput() {
}

func (t *TemplateScene) Update(state *scenes.GameState, deltaTime float64) error {
	return nil
}

func (t *TemplateScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})
}
