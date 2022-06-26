package entities

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mharv/scrapyard-charter/animation"
	"github.com/mharv/scrapyard-charter/basics"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/solarlune/resolv"
)

type OverworldPlayerObject struct {
	animator                              animation.Animator
	physObj                               *resolv.Object
	entityManager                         EntityManager
	moveUp, moveDown, moveRight, moveLeft bool
	moveSpeed                             float64
	CastDistanceLimit                     float64
	alive, move, flip                     bool
}

const (
	frameSizeX = 48
	frameSizeY = 80
)

func (p *OverworldPlayerObject) GetPhysObj() *resolv.Object {
	return p.physObj
}

func (p *OverworldPlayerObject) Init(ImageFilepath string) {
	p.alive = true
	p.move = false
	p.flip = false
	// Load an image given a filepath
	p.physObj = resolv.NewObject(globals.ScreenWidth, globals.ScreenHeight, frameSizeX, frameSizeY/2, "player")

	p.animator = animation.Animator{}
	p.animator.Init(ImageFilepath, basics.Vector2i{X: frameSizeX, Y: frameSizeY}, basics.Vector2f{X: 1, Y: 1}, basics.Vector2f{X: p.physObj.X, Y: p.physObj.Y}, 0.1)
	p.animator.AddAnimation(animation.Animation{
		FrameCount:         1,
		FrameStartPosition: basics.Vector2i{X: 0, Y: 0},
		Loop:               true,
	}, "idleRight")
	p.animator.AddAnimation(animation.Animation{
		FrameCount:         1,
		FrameStartPosition: basics.Vector2i{X: frameSizeX, Y: 0},
		Loop:               true,
	}, "idleLeft")
	p.animator.AddAnimation(animation.Animation{
		FrameCount:         6,
		FrameStartPosition: basics.Vector2i{X: 0, Y: frameSizeY},
		Loop:               true,
	}, "moveRight")
	p.animator.AddAnimation(animation.Animation{
		FrameCount:         6,
		FrameStartPosition: basics.Vector2i{X: 0, Y: frameSizeY * 2},
		Loop:               true,
	}, "moveLeft")
	p.animator.SetAnimation("idleRight", false)

	playerData := globals.GetPlayerData()

	p.moveSpeed = playerData.GetOverworldMoveSpeed()
	p.CastDistanceLimit = playerData.GetOverworldCastDistance()

	// Setup resolv object to be size of the sprite
}

func (p *OverworldPlayerObject) ReadInput() {
	p.entityManager.ReadInput()

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		p.moveSpeed *= 2
	}
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		p.moveSpeed = globals.GetPlayerData().GetOverworldMoveSpeed()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		p.move = true
		p.moveUp = true
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyW) {
		p.moveUp = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		p.flip = false
		p.move = true
		p.moveLeft = true
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyA) {
		p.moveLeft = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		p.move = true
		p.moveDown = true
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyS) {
		p.moveDown = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		p.flip = true
		p.move = true
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

	if !p.moveRight && !p.moveLeft && !p.moveDown && !p.moveUp {
		p.move = false
	}

	if p.moveDown {
		dy = p.moveSpeed * deltaTime
	}

	if p.moveUp {
		dy = p.moveSpeed * deltaTime * -1
	}

	if p.move {
		if p.flip {
			if !(p.animator.IsLooping() && p.animator.IsAnimation("moveRight")) {
				p.animator.SetAnimation("moveRight", false)
			}
		} else {
			if !(p.animator.IsLooping() && p.animator.IsAnimation("moveLeft")) {
				p.animator.SetAnimation("moveLeft", false)
			}
		}
	} else {
		if p.flip {
			if !(p.animator.IsLooping() && p.animator.IsAnimation("idleRight")) {
				p.animator.SetAnimation("idleRight", false)
			}
		} else {
			if !(p.animator.IsLooping() && p.animator.IsAnimation("idleLeft")) {
				p.animator.SetAnimation("idleLeft", false)
			}
		}
	}

	if col := p.physObj.Check(dx, 0, "craft"); col != nil {
		globals.GetPlayerData().SetInCraftZoneStatus(true)
	} else {
		globals.GetPlayerData().SetInCraftZoneStatus(false)
	}

	if col := p.physObj.Check(dy, 0, "craft"); col != nil {
		globals.GetPlayerData().SetInCraftZoneStatus(true)
	} else {
		globals.GetPlayerData().SetInCraftZoneStatus(false)
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
	p.animator.Update(basics.Vector2f{X: p.physObj.X, Y: p.physObj.Y - p.physObj.H}, deltaTime)
}

func (p *OverworldPlayerObject) Draw(screen *ebiten.Image) {
	// Debug drawing of the physics object
	if globals.Debug {
		ebitenutil.DrawRect(screen, p.physObj.X, p.physObj.Y, p.physObj.W, p.physObj.H, color.RGBA{0, 80, 255, 64})
	}

	p.animator.Draw(screen)
}

func (p *OverworldPlayerObject) SetPosition(position basics.Vector2f) {
	p.physObj.X = position.X
	p.physObj.Y = position.Y
}

func (p *OverworldPlayerObject) GetCellPosition() (x, y int) {
	return p.physObj.CellPosition()
}

func (p *OverworldPlayerObject) IsAlive() bool {
	return p.alive
}

func (p *OverworldPlayerObject) Kill() {
	p.alive = false
}

func (p *OverworldPlayerObject) RemovePhysObj(space *resolv.Space) {
	space.Remove(p.physObj)
}
