package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

// This entitiy interface needs to be updated
// once we've decided what all the game entities
// need to have
type Entity interface {
	Init(ImageFilepath string)
	ReadInput()
	Update(deltaTime float64)
	Draw(screen *ebiten.Image)
	IsAlive() bool
	Kill()
	RemovePhysObj(space *resolv.Space)
}

type EntityManager struct {
	entities     []Entity
	deadEntities []Entity
}

func (e *EntityManager) Init() {
}

func (e *EntityManager) ReadInput() {
	for _, entity := range e.entities {
		entity.ReadInput()
	}
}

func (e *EntityManager) Update(deltaTime float64) {
	for i, entity := range e.entities {
		entity.Update(deltaTime)
		if !entity.IsAlive() {
			e.deadEntities = append(e.deadEntities, entity)
			e.entities = append(e.entities[:i], e.entities[i+1:]...)
		}
	}
}

func (e *EntityManager) RemoveDead(space *resolv.Space) {
	if len(e.deadEntities) > 0 {
		newDeadList := &[]Entity{}
		e.deadEntities = *newDeadList
	}
}

func (e *EntityManager) Draw(screen *ebiten.Image) {
	for _, entity := range e.entities {
		entity.Draw(screen)
	}
}

// Ensure entities are initialised before calling this method
func (e *EntityManager) AddEntity(entity Entity) {
	e.entities = append(e.entities, entity)
}
