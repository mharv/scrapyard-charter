package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/mharv/scrapyard-charter/resources"
)

type WinScene struct {
	victory *ebiten.Image
}

func (w *WinScene) Init() {
	globals.GetAudioPlayer().PlayFile("audio/victory.mp3")

	w.victory = resources.LoadFileAsImage("images/victory.png")
}

func (w *WinScene) ReadInput() {
}

func (w *WinScene) Update(state *GameState, deltaTime float64) error {
	return nil
}

func (w *WinScene) Draw(screen *ebiten.Image) {
	imageOptions := &ebiten.DrawImageOptions{}
	imageOptions.GeoM.Translate(0, 0)
	screen.DrawImage(w.victory, imageOptions)
}
