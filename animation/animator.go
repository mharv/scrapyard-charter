package animation

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/mharv/scrapyard-charter/basics"
)

type Animation struct {
	FrameCount         int
	FrameStartPosition basics.Vector2i
	Loop               bool
}

type Animator struct {
	spritesheet      *ebiten.Image
	animations       map[string]Animation
	frameSize        basics.Vector2i
	scale            basics.Vector2f
	position         basics.Vector2f
	currentAnimation Animation
	nextAnimations   []Animation
	speed            float64
	counter          float64
	currentFrame     int
	previousFrame    int
}

var ()

func (a *Animator) AddAnimation(anim Animation, name string) {
	a.animations[name] = anim

	if (a.currentAnimation == Animation{}) {
		a.currentAnimation = a.animations[name]
	}
}

func (a *Animator) SetAnimation(name string, queue bool) {
	if queue {
		a.nextAnimations = append(a.nextAnimations, a.animations[name])
	} else {
		a.currentFrame = 0
		a.previousFrame = 0
		a.currentAnimation = a.animations[name]
		a.nextAnimations = []Animation{}
	}
}

func (a *Animator) IsLooping() bool {
	return a.currentAnimation.Loop
}

func (a *Animator) IsAnimation(name string) bool {
	value, ok := a.animations[name]
	if !ok {
		return false
	}
	if a.currentAnimation == value {
		return true
	}
	return false
}

func (a *Animator) Init(ImageFilepath string, frameSize basics.Vector2i, scale, position basics.Vector2f, speed float64) {
	img, _, err := ebitenutil.NewImageFromFile(ImageFilepath)
	if err != nil {
		log.Fatal(err)
	}
	a.spritesheet = img
	a.frameSize = frameSize
	a.scale = scale
	a.position = position
	a.speed = speed
	a.counter = 0
	a.animations = make(map[string]Animation)
}

func (a *Animator) Update(position basics.Vector2f, deltaTime float64) {
	a.counter += deltaTime
	a.currentFrame = int(a.counter/a.speed) % a.currentAnimation.FrameCount

	a.position = position

	if !a.currentAnimation.Loop {
		if (a.nextAnimations[0] != Animation{}) {
			if a.currentFrame < a.previousFrame {
				a.currentAnimation = a.nextAnimations[0]
				a.nextAnimations = append(a.nextAnimations[:0], a.nextAnimations[1:]...)
			}
		}
	}

	a.previousFrame = a.currentFrame
}

func (a *Animator) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}

	options.GeoM.Scale(a.scale.X, a.scale.Y)
	options.GeoM.Translate(a.position.X, a.position.Y)

	sx, sy := a.currentAnimation.FrameStartPosition.X+a.currentFrame*a.frameSize.X, a.currentAnimation.FrameStartPosition.Y
	screen.DrawImage(a.spritesheet.SubImage(image.Rect(sx, sy, sx+a.frameSize.X, sy+a.frameSize.Y)).(*ebiten.Image), options)
}
