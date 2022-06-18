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

type OverworldPlayerObject struct {
	sprite                                *ebiten.Image
	physObj                               *resolv.Object
	entityManager                         EntityManager
	moveUp, moveDown, moveRight, moveLeft bool
	moveSpeed                             float64
	CastDistanceLimit                     float64
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
	p.physObj = resolv.NewObject(globals.ScreenWidth, globals.ScreenHeight, float64(p.sprite.Bounds().Dx()), float64(p.sprite.Bounds().Dy()/2), "player")
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

	var dx, dy float64

	if p.moveLeft {
		dx = p.moveSpeed * deltaTime * -1
	}

	if p.moveRight {
		dx = p.moveSpeed * deltaTime
	}

	if p.moveDown {
		dy = p.moveSpeed * deltaTime
	}

	if p.moveUp {
		dy = p.moveSpeed * deltaTime * -1
	}

	if col := p.physObj.Check(dx, 0, "solid"); col != nil {
		dx = 0
	}

	p.physObj.X += dx

	if col := p.physObj.Check(0, dy, "solid"); col != nil {
		dy = 0
	}

	p.physObj.Y += dy

	p.physObj.Update()
}

func (p *OverworldPlayerObject) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	// Sprite is put over the top of the phys object

	// Debug drawing of the physics object
	if globals.Debug {
		ebitenutil.DrawRect(screen, p.physObj.X, p.physObj.Y, p.physObj.W, p.physObj.H, color.RGBA{0, 80, 255, 64})
	}

	options.GeoM.Translate(p.physObj.X, p.physObj.Y-p.physObj.H)
	screen.DrawImage(p.sprite, options)
}

func (p *OverworldPlayerObject) SetPosition(position basics.Vector2f) {
	p.physObj.X = position.X
	p.physObj.Y = position.Y
}

func (p *OverworldPlayerObject) GetCellPosition() (x, y int) {
	return p.physObj.CellPosition()
}
