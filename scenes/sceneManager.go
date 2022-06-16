package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mharv/scrapyard-charter/globals"
)

var (
	transitionFrom = ebiten.NewImage(globals.ScreenWidth, globals.ScreenHeight)
	transitionTo   = ebiten.NewImage(globals.ScreenWidth, globals.ScreenHeight)
)

type Scene interface {
	Init()
	ReadInput()
	Update(state *GameState, deltaTime float64) error
	Draw(screen *ebiten.Image)
}

type SceneManager struct {
	current            Scene
	next               Scene
	transitionCount    float64
	transitionMaxCount float64
}

type GameState struct {
	SceneManager *SceneManager
}

func (s *SceneManager) ReadInput() {
	if s.transitionCount <= 0 {
		s.current.ReadInput()
		return
	}
}

func (s *SceneManager) Update(deltaTime float64) error {
	if s.transitionCount <= 0 {
		return s.current.Update(&GameState{
			SceneManager: s,
		}, deltaTime)
	}

	s.transitionCount -= 1 * deltaTime
	if s.transitionCount > 0 {
		return nil
	}

	s.current = s.next
	s.next = nil
	return nil
}

func (s *SceneManager) Draw(screen *ebiten.Image) {
	if s.transitionCount <= 0 {
		s.current.Draw(screen)
		return
	}

	transitionFrom.Clear()
	s.current.Draw(transitionFrom)

	transitionTo.Clear()
	s.next.Draw(transitionTo)

	screen.DrawImage(transitionFrom, nil)

	alpha := 1 - float64(s.transitionCount)/float64(s.transitionMaxCount)
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, alpha)
	screen.DrawImage(transitionTo, op)
}

func (s *SceneManager) GoTo(scene Scene, fadeTime float64) {
	scene.Init()

	if s.current == nil {
		s.current = scene
	} else {
		s.next = scene
		s.transitionCount = fadeTime
		s.transitionMaxCount = fadeTime
	}
}
