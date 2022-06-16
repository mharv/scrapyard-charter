package entities

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/mharv/scrapyard-charter/basics"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/solarlune/resolv"
)

type JunkObject struct {
	sprite  *ebiten.Image
	physObj *resolv.Object
}

func (j *JunkObject) Init(ImageFilepath string) {
	// Load an image given a filepath
	img, _, err := ebitenutil.NewImageFromFile(ImageFilepath)
	if err != nil {
		log.Fatal(err)
	}

	j.sprite = img

	// Setup resolv object to be size of the sprite
	j.physObj = resolv.NewObject(0, 0, float64(j.sprite.Bounds().Dx()), float64(j.sprite.Bounds().Dy()), "junk")
}

func (j *JunkObject) ReadInput() {
}

func (j *JunkObject) Update(deltaTime float64) {
	j.SetPosition(basics.Vector2f{X: j.physObj.X + 5*deltaTime, Y: j.physObj.Y + 2*deltaTime})
}

func (j *JunkObject) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	// Sprite is put over the top of the phys object
	options.GeoM.Translate(j.physObj.X, j.physObj.Y)

	// Debug drawing of the physics object
	if globals.Debug {
		ebitenutil.DrawRect(screen, j.physObj.X, j.physObj.Y, j.physObj.W, j.physObj.H, color.RGBA{0, 80, 255, 255})
	}

	// Draw the image (comment this out to see the above resolv rect ^^^)
	screen.DrawImage(j.sprite, options)
}

func (j *JunkObject) SetPosition(position basics.Vector2f) {
	j.physObj.X = position.X
	j.physObj.Y = position.Y
}
