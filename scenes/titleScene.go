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
	owrld       bool
	scav        bool
	esc         bool
	txtRenderer *etxt.Renderer
}

func (t *TitleScene) Init() {
	t.owrld = false
	t.scav = false
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
	if inpututil.IsKeyJustPressed(ebiten.KeyO) {
		t.owrld = true
	} else {
		t.owrld = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		t.scav = true
	} else {
		t.scav = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		t.esc = true
	} else {
		t.esc = false
	}
}

func (t *TitleScene) Update(state *GameState, deltaTime float64) error {
	if t.owrld {
		o := &OverworldScene{}
		state.SceneManager.GoTo(o, 5)
	}
	if t.scav {
		s := &ScavengeScene{}
		state.SceneManager.GoTo(s, 5)
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
	t.txtRenderer.Draw("Press O to go to the overworld.\nPress S to scavenge\nPress ESC to quit", globals.ScreenWidth/2, globals.ScreenHeight/2)
}
