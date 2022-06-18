package entities

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mharv/scrapyard-charter/basics"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/solarlune/resolv"
)

type OverworldPlayerObject struct {
	sprite                                *ebiten.Image
	physObj                               *resolv.Object
	entityManager                         EntityManager
	moveUp, moveDown, moveRight, moveLeft bool
	moveSpeed                             float64
}

const (
	moveSpeed = 200
)

func (p *OverworldPlayerObject) GetPhysObj() *resolv.Object {
	return p.physObj
}

func (p *OverworldPlayerObject) Init(ImageFilepath string) {
	// Load an image given a filepath
	img, _, err := ebitenutil.NewImageFromFile(ImageFilepath)
	if err != nil {
		log.Fatal(err)
	}

	p.moveSpeed = moveSpeed
	p.sprite = img

	// Setup resolv object to be size of the sprite
	p.physObj = resolv.NewObject(globals.ScreenWidth/4, globals.ScreenHeight/4, float64(p.sprite.Bounds().Dx()), float64(p.sprite.Bounds().Dy()), "player")
}

func (p *OverworldPlayerObject) ReadInput() {
	p.entityManager.ReadInput()

	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		p.moveUp = true
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyW) {
		p.moveUp = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		p.moveLeft = true
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyA) {
		p.moveLeft = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		p.moveDown = true
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyS) {
		p.moveDown = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		p.moveRight = true
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyD) {
		p.moveRight = false
	}
}

func (p *OverworldPlayerObject) Update(deltaTime float64) {
	if p.moveLeft {
		p.physObj.X -= p.moveSpeed * deltaTime
	}

	if p.moveRight {
		p.physObj.X += p.moveSpeed * deltaTime
	}

	if p.moveDown {
		p.physObj.Y += p.moveSpeed * deltaTime
	}

	if p.moveUp {
		p.physObj.Y -= p.moveSpeed * deltaTime
	}

	p.physObj.Update()
}

func (p *OverworldPlayerObject) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	// Sprite is put over the top of the phys object
	options.GeoM.Translate(p.physObj.X, p.physObj.Y)

	// Draw the image (comment this out to see the above resolv rect ^^^)
	screen.DrawImage(p.sprite, options)
}

func (p *OverworldPlayerObject) SetPosition(position basics.Vector2f) {
	p.physObj.X = position.X
	p.physObj.Y = position.Y
}
