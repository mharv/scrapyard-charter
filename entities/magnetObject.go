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

type MagnetObject struct {
	sprite             *ebiten.Image
	physObj            *resolv.Object
	magFieldPhysObj    *resolv.Object
	connectedJunk      *resolv.Object
	attractedJunk      []*resolv.Object
	attractionStrength float64
	magneticFieldSize  float64
	magneticPoint      basics.Vector2f
	linkDistance       basics.Vector2f
	targetPos          basics.Vector2f
	touch, connected   bool
	dropConnected      bool
	dropCounter        float64
}

const (
	dropReactivationTimer     = 0.5
	initialMagneticFieldSize  = 50
	initialAttractionStrength = 10
)

func (m *MagnetObject) GetPhysObj() *resolv.Object {
	return m.physObj
}

func (m *MagnetObject) GetFieldPhysObj() *resolv.Object {
	return m.magFieldPhysObj
}

func (m *MagnetObject) GetSprite() *ebiten.Image {
	return m.sprite
}

func (m *MagnetObject) Init(ImageFilepath string) {
	// Load an image given a filepath
	img, _, err := ebitenutil.NewImageFromFile(ImageFilepath)
	if err != nil {
		log.Fatal(err)
	}

	m.sprite = img

	m.physObj = resolv.NewObject(0, 0, float64(m.sprite.Bounds().Dx()), float64(m.sprite.Bounds().Dy()), "magnet")

	m.magneticFieldSize = initialMagneticFieldSize
	m.attractionStrength = initialAttractionStrength

	m.magFieldPhysObj = resolv.NewObject(-m.magneticFieldSize, -m.magneticFieldSize, float64(m.sprite.Bounds().Dx())+(m.magneticFieldSize*2), float64(m.sprite.Bounds().Dy())+(m.magneticFieldSize*2), "magneticField")

	m.magneticPoint.X = float64(m.sprite.Bounds().Dx() / 2)
	m.magneticPoint.Y = float64(m.sprite.Bounds().Dy())

	m.touch = false
	m.connected = false
	m.dropConnected = false
	m.dropCounter = 0
}

func (m *MagnetObject) ReadInput() {
	x, y := ebiten.CursorPosition()
	m.targetPos.X = float64(x)
	m.targetPos.Y = float64(y)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		m.dropConnected = true
	} else {
		m.dropConnected = false
	}
}

func (m *MagnetObject) Update(deltaTime float64) {
	dx := m.targetPos.X - m.physObj.X
	dy := m.targetPos.Y - m.physObj.Y

	if m.dropCounter > 0 {
		m.dropCounter -= deltaTime
	} else {
		if !m.connected {
			if collision := m.magFieldPhysObj.Check(dx, dy, "junk"); collision != nil {
				m.attractedJunk = collision.Objects
			}

			if collision := m.physObj.Check(dx, dy, "junk"); collision != nil {
				m.attractedJunk = nil
				m.connectedJunk = collision.Objects[0]
				m.linkDistance.X = m.connectedJunk.X - m.physObj.X
				m.linkDistance.Y = m.connectedJunk.Y - m.physObj.Y
				m.touch = true
				m.connected = true
			} else {
				m.touch = false
			}
		}
	}

	if len(m.attractedJunk) > 0 {
		for _, junk := range m.attractedJunk {
			junk.X = basics.FloatLerp(junk.X, m.magneticPoint.X+m.physObj.X, m.attractionStrength*deltaTime)
			junk.Y = basics.FloatLerp(junk.Y, m.magneticPoint.Y+m.physObj.Y, m.attractionStrength*deltaTime)
		}
	}

	m.physObj.X = m.targetPos.X
	m.physObj.Y = m.targetPos.Y

	fieldOffset := basics.Vector2f{X: -m.magneticFieldSize, Y: -m.magneticFieldSize}

	m.magFieldPhysObj.X = m.physObj.X + fieldOffset.X
	m.magFieldPhysObj.Y = m.physObj.Y + fieldOffset.Y

	if m.connected {
		m.connectedJunk.X = m.physObj.X + m.linkDistance.X
		m.connectedJunk.Y = m.physObj.Y + m.linkDistance.Y
	}

	if m.dropConnected {
		m.connectedJunk = nil
		m.linkDistance = basics.Vector2f{X: 0, Y: 0}
		m.connected = false
		m.touch = false
		m.dropCounter = dropReactivationTimer
	}

	m.physObj.Update()
	m.magFieldPhysObj.Update()
}

func (m *MagnetObject) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	// Sprite is put over the top of the phys object
	options.GeoM.Translate(m.physObj.X, m.physObj.Y)

	if m.touch {
		options.ColorM.Scale(0.5, 1, 1, 1)
	}

	// Debug drawing of the physics object
	if globals.Debug {
		ebitenutil.DrawRect(screen, m.magFieldPhysObj.X, m.magFieldPhysObj.Y, m.magFieldPhysObj.W, m.magFieldPhysObj.H, color.RGBA{255, 80, 0, 128})
		ebitenutil.DrawRect(screen, m.physObj.X, m.physObj.Y, m.physObj.W, m.physObj.H, color.RGBA{255, 80, 0, 128})
	}

	// Draw the image (comment this out to see the above resolv rect ^^^)
	screen.DrawImage(m.sprite, options)
}

func (m *MagnetObject) SetPosition(position basics.Vector2f) {
	m.physObj.X = position.X
	m.physObj.Y = position.Y
}
