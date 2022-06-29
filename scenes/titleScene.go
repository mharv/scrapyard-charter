package scenes

import (
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/mharv/scrapyard-charter/resources"
	"github.com/tinne26/etxt"
)

type TitleScene struct {
	owrld            bool
	esc              bool
	image            *ebiten.Image
	txtRenderer      *etxt.Renderer
	titleText        string
	thoughtText      string
	instructionsText string
}

const (
	transitionTime      = 1
	thoughtOffsetX      = 910
	thoughtOffsetY      = 112
	instructionsOffsetX = 126
	instructionsOffsetY = 382
	titleOffsetX        = 70
	titleOffsetY        = 70
)

func (t *TitleScene) Init() {
	t.owrld = false
	t.esc = false

	t.image = resources.LoadFileAsImage("images/titlescreen.png")

	fontLib := resources.LoadFileAsFont("fonts/Rajdhani-Regular.ttf")

	t.titleText = "Scrapyard Magnate"
	t.thoughtText = "I know that golden \nmagnet is out there...\nSomewhere..."
	t.instructionsText = "Press [Enter] to play"

	t.txtRenderer = etxt.NewStdRenderer()
	glyphsCache := etxt.NewDefaultCache(10 * 1024 * 1024) // 10MB
	t.txtRenderer.SetCacheHandler(glyphsCache.NewHandler())
	t.txtRenderer.SetFont(fontLib.GetFont("Rajdhani Regular"))
	t.txtRenderer.SetAlign(etxt.Top, etxt.Left)
	t.txtRenderer.SetSizePx(24)
}

func (t *TitleScene) ReadInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		t.owrld = true
	} else {
		t.owrld = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		t.esc = true
	} else {
		t.esc = false
	}
}

func (t *TitleScene) Update(state *GameState, deltaTime float64) error {
	globals.GetAudioPlayer().PlayFile("audio/menu.mp3")

	if t.owrld {
		o := &OverworldScene{}
		state.SceneManager.GoTo(o, transitionTime)
	}
	if t.esc {
		os.Exit(0)
	}

	return nil
}

func (t *TitleScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(t.image, op)

	t.txtRenderer.SetTarget(screen)
	t.txtRenderer.SetSizePx(80)
	t.txtRenderer.SetColor(color.RGBA{157, 159, 127, 255})
	t.txtRenderer.Draw(t.titleText, titleOffsetX, titleOffsetY)

	t.txtRenderer.SetTarget(screen)
	t.txtRenderer.SetSizePx(25)
	t.txtRenderer.SetColor(color.RGBA{0, 0, 0, 255})
	t.txtRenderer.Draw(t.thoughtText, thoughtOffsetX, thoughtOffsetY)

	t.txtRenderer.SetTarget(screen)
	t.txtRenderer.SetSizePx(40)
	t.txtRenderer.SetColor(color.RGBA{197, 204, 184, 255})
	t.txtRenderer.Draw(t.instructionsText, instructionsOffsetX, instructionsOffsetY)
}
