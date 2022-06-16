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

func (j *JunkObject) GetPhysObj() *resolv.Object {
	return j.physObj
}

func (j *JunkObject) GetSprite() *ebiten.Image {
	return j.sprite
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
	j.physObj.Update()
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
	x := position.X
	y := position.Y

	j.physObj.X = x
	j.physObj.Y = y

	if len(j.physObj.Space.Objects()) > 1 {
		for _, obj := range j.physObj.Space.Objects() {
			if obj != j.physObj {
				if j.physObj.Overlaps(obj) {
					if j.physObj.X > obj.X {
						j.physObj.X += ((obj.X + obj.W) - j.physObj.X)
					} else {
						j.physObj.X += (obj.X - (j.physObj.X + j.physObj.W))
					}

					if j.physObj.Y > obj.Y {
						j.physObj.Y += ((obj.Y + obj.H) - j.physObj.Y)
					} else {
						j.physObj.Y += (obj.Y - (j.physObj.Y + j.physObj.H))
					}
				}
			}
		}
	}
}
