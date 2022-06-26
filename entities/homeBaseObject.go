package entities

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/mharv/scrapyard-charter/animation"
	"github.com/mharv/scrapyard-charter/basics"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/solarlune/resolv"
)

type HomeBaseObject struct {
	animator  animation.Animator
	physObj   *resolv.Object
	craftZone *resolv.Object
	alive     bool
}

const (
	homeFrameSize     = 128
	homePhysObjOffset = 32
)

func (h *HomeBaseObject) GetPhysObj() *resolv.Object {
	return h.physObj
}

func (h *HomeBaseObject) GetCraftZone() *resolv.Object {
	return h.craftZone
}

func (h *HomeBaseObject) SetPosition(position basics.Vector2f) {
	h.physObj.X = position.X
	h.physObj.Y = position.Y
	h.craftZone.X = position.X - homeFrameSize
	h.craftZone.Y = position.Y - homeFrameSize
}

func (h *HomeBaseObject) Init(ImageFilepath string) {
	h.alive = true

	h.physObj = resolv.NewObject(globals.ScreenWidth/2, globals.ScreenHeight/2, homeFrameSize-homePhysObjOffset, homeFrameSize-homePhysObjOffset)

	h.animator = animation.Animator{}
	h.animator.Init(ImageFilepath, basics.Vector2i{X: homeFrameSize, Y: homeFrameSize}, basics.Vector2f{X: 1, Y: 1}, basics.Vector2f{X: h.physObj.X, Y: h.physObj.Y}, 0.07)
	h.animator.AddAnimation(animation.Animation{
		FrameCount:         6,
		FrameStartPosition: basics.Vector2i{X: 0, Y: 0},
		Loop:               true,
	}, "idle")
	h.animator.SetAnimation("idle", false)

	h.craftZone = resolv.NewObject(
		globals.ScreenWidth/2,
		globals.ScreenHeight/2,
		float64(homeFrameSize*3),
		float64(homeFrameSize*3),
		"craft",
	)
}

func (h *HomeBaseObject) ReadInput() {
}

func (h *HomeBaseObject) Update(deltaTime float64) {
	h.animator.Update(basics.Vector2f{X: h.physObj.X - (homePhysObjOffset / 2), Y: h.physObj.Y - (homePhysObjOffset / 2)}, deltaTime)
}

func (h *HomeBaseObject) Draw(screen *ebiten.Image) {
	// Debug drawing of the physics object
	if globals.Debug {
		ebitenutil.DrawRect(screen, h.physObj.X, h.physObj.Y, h.physObj.W, h.physObj.H, color.RGBA{0, 80, 255, 128})
	}

	h.animator.Draw(screen)
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
