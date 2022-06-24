package entities

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/mharv/scrapyard-charter/basics"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/mharv/scrapyard-charter/inventory"
	"github.com/solarlune/resolv"
)

type JunkObject struct {
	sprite        *ebiten.Image
	physObj       *resolv.Object
	itemData      inventory.Item
	imageFilepath string
	rot           float64
	alive         bool
}

const (
	junkPhysObjSizeDiff = 15
)

func (j *JunkObject) GetPhysObj() *resolv.Object {
	return j.physObj
}

func (j *JunkObject) GetSprite() *ebiten.Image {
	return j.sprite
}

func (j *JunkObject) Init(ImageFilepath string) {
	j.alive = true
	// Load an image given a filepath
	img, _, err := ebitenutil.NewImageFromFile(ImageFilepath)
	if err != nil {
		log.Fatal(err)
	}

	j.sprite = img

	// Setup resolv object to be size of the sprite
	j.physObj = resolv.NewObject(0, 0, float64(j.sprite.Bounds().Dx())-junkPhysObjSizeDiff, float64(j.sprite.Bounds().Dy())-junkPhysObjSizeDiff, "junk")

	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)
	j.rot = float64(rnd.Intn(360))
}

func (j *JunkObject) ReadInput() {
}

func (j *JunkObject) Update(deltaTime float64) {
	j.physObj.Update()
}

func (j *JunkObject) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(-float64(j.sprite.Bounds().Dx())/2, -float64(j.sprite.Bounds().Dy())/2)
	options.GeoM.Rotate(j.rot)
	options.GeoM.Translate(float64(j.sprite.Bounds().Dx())/2, float64(j.sprite.Bounds().Dy())/2)
	// Sprite is put over the top of the phys object
	options.GeoM.Translate(j.physObj.X-(junkPhysObjSizeDiff/2), j.physObj.Y-(junkPhysObjSizeDiff/2))

	// Debug drawing of the physics object
	if globals.Debug {
		ebitenutil.DrawRect(screen, j.physObj.X, j.physObj.Y, j.physObj.W, j.physObj.H, color.RGBA{0, 80, 255, 64})
	}

	// Draw the image (comment this out to see the above resolv rect ^^^)
	screen.DrawImage(j.sprite, options)
}

func (j *JunkObject) IsAlive() bool {
	return j.alive
}

func (j *JunkObject) Kill() {
	fmt.Println("Removing " + j.itemData.GetName() + "'s phys obj from ")
	j.RemovePhysObj(j.physObj.Space)
	fmt.Println("Seeting alive to false for " + j.itemData.GetName())
	j.alive = false
}

func (j *JunkObject) RemovePhysObj(space *resolv.Space) {
	space.Remove(j.physObj)
}

func (j *JunkObject) InitData() {
	j.itemData.Init()
}

func (j *JunkObject) SetItemDataName(name string) {
	j.itemData.SetName(name)
}

func (j *JunkObject) GetItemData() *inventory.Item {
	return &j.itemData
}

func (j *JunkObject) SetImageFilepath(filepath string) {
	j.imageFilepath = filepath
}

func (j *JunkObject) GetImageFilepath() string {
	return j.imageFilepath
}

func (j *JunkObject) AddItemDataMaterial(materialName string, minQuantity, maxQuantity int) {
	quantity := rand.Intn(maxQuantity-minQuantity) + minQuantity
	j.itemData.AddRawMaterial(materialName, quantity)
}

func (j *JunkObject) IsPhysObject(physObjectToCompare *resolv.Object) bool {
	return j.physObj == physObjectToCompare
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
