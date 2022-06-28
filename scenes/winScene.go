package scenes

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/tinne26/etxt"
)

type WinScene struct {
	victory *ebiten.Image
}

func (w *WinScene) Init() {
	globals.GetAudioPlayer().PlayFile("audio/victory.mp3")
	fontLib := etxt.NewFontLibrary()

	_, _, err := fontLib.ParseDirFonts("fonts")
	if err != nil {
		log.Fatal(err)
	}

	if !fontLib.HasFont("Rajdhani Regular") {
		log.Fatal("missing font Rajdhani-Regular.ttf")
	}

	w.victory = LoadImage("images/victory.png")
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
