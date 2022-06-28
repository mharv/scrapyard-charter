package game

import (
	_ "image/png"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/mharv/scrapyard-charter/scenes"
)

var (
	newTime, oldTime int64
)

type Game struct {
	sceneManager *scenes.SceneManager
}

func (g *Game) Update() error {
	oldTime = newTime
	newTime = time.Now().UnixNano()
	deltaTime := float64((newTime-oldTime)/1000000) * 0.001

	if g.sceneManager == nil {
		g.sceneManager = &scenes.SceneManager{}
		g.sceneManager.GoTo(&scenes.TitleScene{}, 0)
	}

	g.sceneManager.ReadInput()
	if err := g.sceneManager.Update(deltaTime); err != nil {
		return err
	}

	// globals.GetPlayerData().Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.sceneManager.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return globals.ScreenWidth, globals.ScreenHeight
}

func (g *Game) Init() {
	ebiten.SetWindowSize(globals.ScreenWidth, globals.ScreenHeight)
	ebiten.SetWindowTitle("Scrapyard Charter")
	globals.GetPlayerData().Init()
}
