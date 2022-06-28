package scenes

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/tinne26/etxt"
)

type WinScene struct {
	txtRenderer *etxt.Renderer
}

func (w *WinScene) Init() {
	fontLib := etxt.NewFontLibrary()

	_, _, err := fontLib.ParseDirFonts("fonts")
	if err != nil {
		log.Fatal(err)
	}

	if !fontLib.HasFont("Rajdhani Regular") {
		log.Fatal("missing font Rajdhani-Regular.ttf")
	}

	w.txtRenderer = etxt.NewStdRenderer()
	glyphsCache := etxt.NewDefaultCache(10 * 1024 * 1024) // 10MB
	w.txtRenderer.SetCacheHandler(glyphsCache.NewHandler())
	w.txtRenderer.SetFont(fontLib.GetFont("Rajdhani Regular"))
	w.txtRenderer.SetAlign(etxt.YCenter, etxt.XCenter)
	w.txtRenderer.SetSizePx(24)
}

func (w *WinScene) ReadInput() {
}

func (w *WinScene) Update(state *GameState, deltaTime float64) error {
	return nil
}

func (w *WinScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	w.txtRenderer.SetTarget(screen)
	w.txtRenderer.SetColor(color.RGBA{255, 255, 255, 255})
	w.txtRenderer.Draw("Good work, you got a golden magnet!", globals.ScreenWidth/2, globals.ScreenHeight/2)
}
