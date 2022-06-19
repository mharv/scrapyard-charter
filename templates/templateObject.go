package template

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/mharv/scrapyard-charter/basics"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/solarlune/resolv"
)

type TemplateObject struct {
	sprite  *ebiten.Image
	physObj *resolv.Object //Optional
}

func (t *TemplateObject) GetPhysObj() *resolv.Object {
	return t.physObj
}

func (t *TemplateObject) GetSprite() *ebiten.Image {
	return t.sprite
}

func (t *TemplateObject) Init(ImageFilepath string) {
	// Load an image given a filepath
	img, _, err := ebitenutil.NewImageFromFile(ImageFilepath)
	if err != nil {
		log.Fatal(err)
	}

	t.sprite = img

	// Setup resolv object to be size of the sprite with tags
	t.physObj = resolv.NewObject(0, 0, float64(t.sprite.Bounds().Dx()), float64(t.sprite.Bounds().Dy()), "tags go here")
}

func (t *TemplateObject) ReadInput() {
}

func (t *TemplateObject) Update(deltaTime float64) {
	t.physObj.Update()
}

func (t *TemplateObject) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	// Sprite is put over the top of the phys object
	options.GeoM.Translate(t.physObj.X, t.physObj.Y)

	// Debug drawing of the physics object
	if globals.Debug {
		ebitenutil.DrawRect(screen, t.physObj.X, t.physObj.Y, t.physObj.W, t.physObj.H, color.RGBA{255, 255, 255, 64})
	}

	// Draw the image (comment this out to see the above resolv rect ^^^)
	screen.DrawImage(t.sprite, options)
}

func (t *TemplateObject) SetPosition(position basics.Vector2f) {
	t.physObj.X = position.X
	t.physObj.Y = position.Y
}
