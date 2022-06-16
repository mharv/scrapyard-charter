package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// This entitiy interface needs to be updated
// once we've decided what all the game entities
// need to have
type Entity interface {
	Init(ImageFilepath string)
	ReadInput()
	Update(deltaTime float64)
	Draw(screen *ebiten.Image)
}

type EntityManager struct {
	entities []Entity
}

func (e *EntityManager) Init() {
}

func (e *EntityManager) ReadInput() {
	for i := range e.entities {
		e.entities[i].ReadInput()
	}
}

func (e *EntityManager) Update(deltaTime float64) {
	for i := range e.entities {
		e.entities[i].Update(deltaTime)
	}
}

func (e *EntityManager) Draw(screen *ebiten.Image) {
	for i := range e.entities {
		e.entities[i].Draw(screen)
	}
}

// Ensure entities are initialised before calling this method
func (e *EntityManager) AddEntity(entity Entity) {
	e.entities = append(e.entities, entity)
}
