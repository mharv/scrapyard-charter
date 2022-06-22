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

type HomeBaseObject struct {
	sprite  *ebiten.Image
	physObj *resolv.Object
	alive   bool
}

func (h *HomeBaseObject) GetPhysObj() *resolv.Object {
	return h.physObj
}

func (h *HomeBaseObject) SetPosition(position basics.Vector2f) {
	h.physObj.X = position.X
	h.physObj.Y = position.Y
}

func (h *HomeBaseObject) Init(ImageFilepath string) {
	h.alive = true
	// Load an image given a filepath
	img, _, err := ebitenutil.NewImageFromFile(ImageFilepath)
	if err != nil {
		log.Fatal(err)
	}

	h.sprite = img

	// Setup resolv object to be size of the sprite
	h.physObj = resolv.NewObject(globals.ScreenWidth/2, globals.ScreenHeight/2, float64(h.sprite.Bounds().Dx()), float64(h.sprite.Bounds().Dy()))
}

func (h *HomeBaseObject) ReadInput() {
}

func (h *HomeBaseObject) Update(deltaTime float64) {
}

func (h *HomeBaseObject) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	// Sprite is put over the top of the phys object
	options.GeoM.Translate(h.physObj.X, h.physObj.Y)

	// Debug drawing of the physics object
	ebitenutil.DrawRect(screen, h.physObj.X, h.physObj.Y, h.physObj.W, h.physObj.H, color.RGBA{0, 80, 255, 255})

	// Draw the image (comment this out to see the above resolv rect ^^^)
	screen.DrawImage(h.sprite, options)
}

func (h *HomeBaseObject) IsAlive() bool {
	return h.alive
}

func (h *HomeBaseObject) Kill() {
	h.alive = false
}

func (h *HomeBaseObject) RemovePhysObj(space *resolv.Space) {
	space.Remove(h.physObj)
}
