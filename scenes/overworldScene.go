package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mharv/scrapyard-charter/entities"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/solarlune/resolv"
)

type OverworldScene struct {
	entityManager entities.EntityManager
	menuBtn       bool
	physSpace     *resolv.Space
}

func (o *OverworldScene) Init() {
	o.physSpace = resolv.NewSpace(globals.ScreenWidth, globals.ScreenHeight, 16, 16)

}

func (o *OverworldScene) ReadInput() {
	o.entityManager.ReadInput()

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		o.menuBtn = true
	} else {
		o.menuBtn = false
	}
}

func (o *OverworldScene) Update(state *GameState, deltaTime float64) error {
	o.entityManager.Update(deltaTime)

	if o.menuBtn {
		t := &TitleScene{}
		state.SceneManager.GoTo(t, transitionTime)
	}

	return nil
}

func (o *OverworldScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 0, 255, 255})
	o.entityManager.Draw(screen)
}
