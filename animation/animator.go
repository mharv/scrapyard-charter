package animation

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/mharv/scrapyard-charter/basics"
)

type Animation struct {
	frameCount         int
	frameStartPosition basics.Vector2i
	loop               bool
}

type Animator struct {
	spritesheet      *ebiten.Image
	animations       map[string]Animation
	frameSize        basics.Vector2i
	scale            basics.Vector2f
	position         basics.Vector2f
	currentAnimation Animation
	nextAnimations   []Animation
	speed            int
	counter          int
}

var (
	currentFrame  int
	previousFrame int
)

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
		currentFrame = 0
		previousFrame = 0
		a.currentAnimation = a.animations[name]
		a.nextAnimations = []Animation{}
	}
}

func (a *Animator) IsLooping() bool {
	return a.currentAnimation.loop
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

func (a *Animator) Init(ImageFilepath string, frameSize basics.Vector2i, scale, position basics.Vector2f, speed int) {
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

func (a *Animator) Update(deltaTime float64) {
	a.counter++
	currentFrame = (a.counter / a.speed) % a.currentAnimation.frameCount

	if !a.currentAnimation.loop {
		if (a.nextAnimations[0] != Animation{}) {
			if currentFrame < previousFrame {
				a.currentAnimation = a.nextAnimations[0]
				a.nextAnimations = append(a.nextAnimations[:0], a.nextAnimations[1:]...)
			}
		}
	}

	previousFrame = currentFrame
}

func (a *Animator) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}

	options.GeoM.Scale(a.scale.X, a.scale.Y)
	options.GeoM.Translate((-float64(a.frameSize.X)*a.scale.X)/2, (-float64(a.frameSize.Y)*a.scale.Y)/2)
	options.GeoM.Translate(a.position.X, a.position.Y)

	sx, sy := a.currentAnimation.frameStartPosition.X+currentFrame*a.frameSize.X, a.currentAnimation.frameStartPosition.Y
	screen.DrawImage(a.spritesheet.SubImage(image.Rect(sx, sy, sx+a.frameSize.X, sy+a.frameSize.Y)).(*ebiten.Image), options)
}
