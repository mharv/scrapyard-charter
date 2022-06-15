package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type OverworldScene struct {
}

func (o *OverworldScene) Init() {
}

func (o *OverworldScene) ReadInput() {
}

func (o *OverworldScene) Update(state *GameState, deltaTime float64) error {
	return nil
}

func (o *OverworldScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 0, 0, 255})
}
