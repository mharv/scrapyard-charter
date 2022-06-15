package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type ScavengeScene struct {
}

func (s *ScavengeScene) Init() {
}

func (s *ScavengeScene) ReadInput() {
}

func (s *ScavengeScene) Update(state *GameState, deltaTime float64) error {
	return nil
}

func (s *ScavengeScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 0, 0, 255})
}
