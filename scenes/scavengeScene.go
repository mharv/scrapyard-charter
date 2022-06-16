package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mharv/scrapyard-charter/entities"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/solarlune/resolv"
)

type ScavengeScene struct {
	entitiyManager entities.EntityManager
	physSpace      *resolv.Space
}

func (s *ScavengeScene) Init() {
	s.physSpace = resolv.NewSpace(globals.ScreenWidth, globals.ScreenHeight, 16, 16)

	s.entitiyManager.Init()

	j := &entities.JunkObject{}
	j.Init("images/oldpc.png")

	s.entitiyManager.AddEntity(j)
}

func (s *ScavengeScene) ReadInput() {
	s.entitiyManager.ReadInput()
}

func (s *ScavengeScene) Update(state *GameState, deltaTime float64) error {
	s.entitiyManager.Update(deltaTime)
	return nil
}

func (s *ScavengeScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})
	s.entitiyManager.Draw(screen)
}
