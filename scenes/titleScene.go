package scenes

import (
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/tinne26/etxt"
)

type TitleScene struct {
	play        bool
	options     bool
	esc         bool
	txtRenderer *etxt.Renderer
}

func (t *TitleScene) Init() {
	t.play = false
	t.options = false
	t.esc = false

	fontLib := etxt.NewFontLibrary()

	_, _, err := fontLib.ParseDirFonts("fonts")
	if err != nil {
		log.Fatal(err)
	}

	if !fontLib.HasFont("Blaka Regular") {
		log.Fatal("missing font Blaka-Regular.ttf")
	}

	t.txtRenderer = etxt.NewStdRenderer()
	glyphsCache := etxt.NewDefaultCache(10 * 1024 * 1024) // 10MB
	t.txtRenderer.SetCacheHandler(glyphsCache.NewHandler())
	t.txtRenderer.SetFont(fontLib.GetFont("Blaka Regular"))
	t.txtRenderer.SetAlign(etxt.YCenter, etxt.XCenter)
	t.txtRenderer.SetSizePx(24)
}

func (t *TitleScene) ReadInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		t.play = true
	} else {
		t.play = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyO) {
		t.options = true
	} else {
		t.options = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		t.esc = true
	} else {
		t.esc = false
	}
}

func (t *TitleScene) Update(state *GameState, deltaTime float64) error {
	if t.play {
		g := &OverworldScene{}
		state.SceneManager.GoTo(g, 2)
	}
	if t.options {
		// o := &OptionsScene{}
		// state.SceneManager.GoTo(o, 10)
	}
	if t.esc {
		os.Exit(0)
	}
	return nil
}

func (t *TitleScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	t.txtRenderer.SetTarget(screen)
	t.txtRenderer.SetColor(color.RGBA{255, 255, 255, 255})
	t.txtRenderer.Draw("Press SPACE to play\nPress O for options\nPress ESC to quit", globals.ScreenWidth/2, globals.ScreenHeight/2)
}
