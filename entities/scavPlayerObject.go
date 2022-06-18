package entities

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mharv/scrapyard-charter/basics"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/solarlune/resolv"
)

type ScavPlayerObject struct {
	sprite           *ebiten.Image
	physObj          *resolv.Object
	magnet           *MagnetObject
	left, right      bool
	moveSpeed        float64
	fishingLinePoint basics.Vector2f
}

const (
	initialMoveSpeed = 100
	rodEndX          = 121
	rodEndY          = 58
)

func (s *ScavPlayerObject) SetMagnet(m *MagnetObject) {
	s.magnet = m
}

func (s *ScavPlayerObject) GetPhysObj() *resolv.Object {
	return s.physObj
}

func (s *ScavPlayerObject) GetFishingLinePoint() *basics.Vector2f {
	return &s.fishingLinePoint
}

func (s *ScavPlayerObject) Init(ImageFilepath string) {
	// Load an image given a filepath
	img, _, err := ebitenutil.NewImageFromFile(ImageFilepath)
	if err != nil {
		log.Fatal(err)
	}

	s.sprite = img

	// Setup resolv object to be size of the sprite
	s.physObj = resolv.NewObject(globals.ScreenWidth/2, globals.ScreenHeight/2, float64(s.sprite.Bounds().Dx()), float64(s.sprite.Bounds().Dy()), "player")

	s.left = false
	s.right = false

	s.fishingLinePoint.X = rodEndX + s.physObj.X
	s.fishingLinePoint.Y = rodEndY + s.physObj.Y

	s.moveSpeed = initialMoveSpeed
}

func (s *ScavPlayerObject) ReadInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		s.left = true
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyA) {
		s.left = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		s.right = true
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyD) {
		s.right = false
	}
}

func (s *ScavPlayerObject) Update(deltaTime float64) {
	if s.left {
		s.physObj.X -= s.moveSpeed * deltaTime
	}

	if s.right {
		s.physObj.X += s.moveSpeed * deltaTime
	}

	s.fishingLinePoint.X = rodEndX + s.physObj.X
	s.fishingLinePoint.Y = rodEndY + s.physObj.Y

	s.physObj.Update()
}

func (s *ScavPlayerObject) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	// Sprite is put over the top of the phys object
	options.GeoM.Translate(s.physObj.X, s.physObj.Y)

	// Debug drawing of the physics object
	if globals.Debug {
		ebitenutil.DrawRect(screen, s.physObj.X, s.physObj.Y, s.physObj.W, s.physObj.H, color.RGBA{0, 80, 255, 64})
	}

	ebitenutil.DrawLine(screen, s.fishingLinePoint.X, s.fishingLinePoint.Y, s.magnet.GetFishingLinePoint().X, s.magnet.GetFishingLinePoint().Y, color.RGBA{197, 204, 184, 255})

	// Draw the image (comment this out to see the above resolv rect ^^^)
	screen.DrawImage(s.sprite, options)
}

func (s *ScavPlayerObject) SetPosition(position basics.Vector2f) {
	s.physObj.X = position.X
	s.physObj.Y = position.Y
}