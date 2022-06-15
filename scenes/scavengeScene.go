package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mharv/scrapyard-charter/entities"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/solarlune/resolv"
)

type ScavengeScene struct {
	testObj   entities.GameObject
	physSpace *resolv.Space
}

func (s *ScavengeScene) Init() {
	s.physSpace = resolv.NewSpace(globals.ScreenWidth, globals.ScreenHeight, 16, 16)
	s.testObj.Init("images/test.png")
}

func (s *ScavengeScene) ReadInput() {
	s.testObj.ReadInput()
}

func (s *ScavengeScene) Update(state *GameState, deltaTime float64) error {
	s.testObj.Update(deltaTime)
	return nil
}

func (s *ScavengeScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})
	s.testObj.Draw(screen)
}
