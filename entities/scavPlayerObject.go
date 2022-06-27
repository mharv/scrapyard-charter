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

type ScavPlayerObject struct {
	animator             animation.Animator
	physObj              *resolv.Object
	magnet               *MagnetObject
	left, right          bool
	moveSpeed            float64
	fishingRodEndPoint   basics.Vector2f
	fishingRodStartPoint basics.Vector2f
	currentRodEndPoint   *basics.Vector2f
	alive                bool
}

const (
	frameSize = 128
)

func (s *ScavPlayerObject) SetMagnet(m *MagnetObject) {
	s.magnet = m
}

func (s *ScavPlayerObject) GetPhysObj() *resolv.Object {
	return s.physObj
}

func (s *ScavPlayerObject) GetFishingRodEndPoint() *basics.Vector2f {
	return &s.fishingRodEndPoint
}

func (s *ScavPlayerObject) GetFishingRodStartPoint() *basics.Vector2f {
	return &s.fishingRodStartPoint
}

func (s *ScavPlayerObject) SetFishingRodEndPoint(rodEndPoint *basics.Vector2f) {
	s.currentRodEndPoint = rodEndPoint
}

func (s *ScavPlayerObject) Init(ImageFilepath string) {
	s.alive = true

	s.physObj = resolv.NewObject(globals.ScreenWidth/2, globals.ScreenHeight/2, frameSize, frameSize, "player")

	s.animator.Init(ImageFilepath, basics.Vector2i{X: frameSize, Y: frameSize}, basics.Vector2f{X: 1, Y: 1}, basics.Vector2f{X: s.physObj.X, Y: s.physObj.Y}, 0.8)
	s.animator.AddAnimation(animation.Animation{
		FrameCount:         6,
		FrameStartPosition: basics.Vector2i{X: 0, Y: 0},
		Loop:               true,
	}, "moveLeft")
	s.animator.AddAnimation(animation.Animation{
		FrameCount:         6,
		FrameStartPosition: basics.Vector2i{X: 0, Y: frameSize},
		Loop:               true,
	}, "moveRight")
	s.animator.AddAnimation(animation.Animation{
		FrameCount:         1,
		FrameStartPosition: basics.Vector2i{X: 0, Y: frameSize * 2},
		Loop:               true,
	}, "idle")
	s.animator.SetAnimation("idle", false)

	s.left = false
	s.right = false

	s.fishingRodEndPoint.X = globals.GetPlayerData().GetRodEndX() + s.physObj.X
	s.fishingRodEndPoint.Y = globals.GetPlayerData().GetRodEndY() + s.physObj.Y

	s.fishingRodStartPoint.X = globals.GetPlayerData().GetRodStartX() + s.physObj.X
	s.fishingRodStartPoint.Y = globals.GetPlayerData().GetRodStartY() + s.physObj.Y

	s.moveSpeed = globals.GetPlayerData().GetScavMoveSpeed()
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
		if !(s.animator.IsLooping() && s.animator.IsAnimation("moveLeft")) {
			s.animator.SetAnimation("moveLeft", false)
		}
	}

	if s.right {
		s.physObj.X += s.moveSpeed * deltaTime
		if !(s.animator.IsLooping() && s.animator.IsAnimation("moveRight")) {
			s.animator.SetAnimation("moveRight", false)
		}
	}

	if !s.right && !s.left {
		if !(s.animator.IsLooping() && s.animator.IsAnimation("idle")) {
			s.animator.SetAnimation("idle", false)
		}
	}

	s.fishingRodEndPoint.X = globals.GetPlayerData().GetRodEndX() + s.physObj.X
	s.fishingRodEndPoint.Y = globals.GetPlayerData().GetRodEndY() + s.physObj.Y

	s.fishingRodStartPoint.X = globals.GetPlayerData().GetRodStartX() + s.physObj.X
	s.fishingRodStartPoint.Y = globals.GetPlayerData().GetRodStartY() + s.physObj.Y

	s.physObj.Update()
	s.animator.Update(basics.Vector2f{X: s.physObj.X, Y: s.physObj.Y}, deltaTime)
}

func (s *ScavPlayerObject) Draw(screen *ebiten.Image) {
	// Debug drawing of the physics object
	if globals.Debug {
		ebitenutil.DrawRect(screen, s.physObj.X, s.physObj.Y, s.physObj.W, s.physObj.H, color.RGBA{0, 80, 255, 64})
	}

	ebitenutil.DrawLine(screen, s.currentRodEndPoint.X, s.currentRodEndPoint.Y, s.magnet.GetFishingLinePoint().X, s.magnet.GetFishingLinePoint().Y, color.RGBA{197, 204, 184, 255})

	// Draw the image (comment this out to see the above resolv rect ^^^)
	s.animator.Draw(screen)
}

func (s *ScavPlayerObject) IsAlive() bool {
	return s.alive
}

func (s *ScavPlayerObject) Kill() {
	s.alive = false
}

func (s *ScavPlayerObject) RemovePhysObj(space *resolv.Space) {
	space.Remove(s.physObj)
}

func (s *ScavPlayerObject) SetPosition(position basics.Vector2f) {
	s.physObj.X = position.X
	s.physObj.Y = position.Y
}
