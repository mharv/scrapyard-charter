package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mharv/scrapyard-charter/entities"
)

type OverworldScene struct {
	entitiyManager entities.EntityManager
	menuBtn        bool
}

func (o *OverworldScene) Init() {
}

func (o *OverworldScene) ReadInput() {
	o.entitiyManager.ReadInput()

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		o.menuBtn = true
	} else {
		o.menuBtn = false
	}
}

func (o *OverworldScene) Update(state *GameState, deltaTime float64) error {
	o.entitiyManager.Update(deltaTime)

	if o.menuBtn {
		t := &TitleScene{}
		state.SceneManager.GoTo(t, transitionTime)
	}

	return nil
}

func (o *OverworldScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 0, 255, 255})
}
