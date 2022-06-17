package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mharv/scrapyard-charter/entities"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/solarlune/resolv"
)

type ScavengeScene struct {
	entityManager entities.EntityManager
	physSpace     *resolv.Space
}

func (s *ScavengeScene) Init() {
	s.physSpace = resolv.NewSpace(globals.ScreenWidth, globals.ScreenHeight, 16, 16)

	s.entityManager.Init()

	j := &entities.JunkObject{}
	j.Init("images/oldpc.png")

	s.entityManager.AddEntity(j)
}

func (s *ScavengeScene) ReadInput() {
	s.entityManager.ReadInput()
}

func (s *ScavengeScene) Update(state *GameState, deltaTime float64) error {
	s.entityManager.Update(deltaTime)
	return nil
}

func (s *ScavengeScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})
	s.entityManager.Draw(screen)
}
