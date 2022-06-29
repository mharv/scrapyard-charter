package entities

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/mharv/scrapyard-charter/resources"
	"github.com/solarlune/resolv"
)

type GameObject struct {
	sprite  *ebiten.Image
	physObj *resolv.Object
	alive   bool
}

func (g *GameObject) Init(ImageFilepath string) {
	g.alive = true
	// Load an image given a filepath
	g.sprite = resources.LoadFileAsImage(ImageFilepath)

	// Setup resolv object to be size of the sprite
	g.physObj = resolv.NewObject(globals.ScreenWidth/2, globals.ScreenHeight/2, float64(g.sprite.Bounds().Dx()), float64(g.sprite.Bounds().Dy()))
}

func (g *GameObject) ReadInput() {
}

func (g *GameObject) Update(deltaTime float64) {
}

func (g *GameObject) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	// Sprite is put over the top of the phys object
	options.GeoM.Translate(g.physObj.X, g.physObj.Y)

	// Debug drawing of the physics object
	ebitenutil.DrawRect(screen, g.physObj.X, g.physObj.Y, g.physObj.W, g.physObj.H, color.RGBA{0, 80, 255, 255})

	// Draw the image (comment this out to see the above resolv rect ^^^)
	screen.DrawImage(g.sprite, options)
}

func (g *GameObject) IsAlive() bool {
	return g.alive
}

func (g *GameObject) Kill() {
	g.alive = false
}

func (g *GameObject) RemovePhysObj(space *resolv.Space) {
	space.Remove(g.physObj)
}
