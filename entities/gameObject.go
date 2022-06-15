package entities

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/solarlune/resolv"
)

type GameObject struct {
	sprite  *ebiten.Image
	physObj *resolv.Object
}

func (g *GameObject) Init(ImageFilepath string) {
	// Load an image given a filepath
	img, _, err := ebitenutil.NewImageFromFile(ImageFilepath)
	if err != nil {
		log.Fatal(err)
	}

	g.sprite = img

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
