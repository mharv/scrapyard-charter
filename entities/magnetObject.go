package entities

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mharv/scrapyard-charter/basics"
	"github.com/mharv/scrapyard-charter/globals"
	"github.com/solarlune/resolv"
)

type MagnetObject struct {
	sprite             *ebiten.Image
	targetSprite       *ebiten.Image
	physObj            *resolv.Object
	targetObj          *resolv.Object
	magFieldPhysObj    *resolv.Object
	connectedJunk      *resolv.Object
	attractedJunk      []*resolv.Object
	rotation           float64
	attractionStrength float64
	magneticFieldSize  float64
	magneticPoint      basics.Vector2f
	linkDistance       basics.Vector2f
	magnetPos          basics.Vector2f
	magnetStartPos     *basics.Vector2f
	targetPos          basics.Vector2f
	touch, connected   bool
	magnetActive       bool
	dropConnected      bool
	dropCounter        float64
}

const (
	dropReactivationTimer     = 0.5
	initialMagneticFieldSize  = 50
	initialAttractionStrength = 10
	magPhysObjSizeDiff        = 20
)

func (m *MagnetObject) GetFishingLinePoint() basics.Vector2f {
	value := basics.Vector2f{X: m.physObj.X + (m.physObj.W / 2), Y: m.physObj.Y}
	return value
}

func (m *MagnetObject) SetMagnetStartPos(pos *basics.Vector2f) {
	m.magnetStartPos = pos
}

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
	m.physObj = resolv.NewObject(0, 0, float64(m.sprite.Bounds().Dx())-magPhysObjSizeDiff, float64(m.sprite.Bounds().Dy())-magPhysObjSizeDiff, "magnet")

	tar, _, err := ebitenutil.NewImageFromFile("images/target.png")
	if err != nil {
		log.Fatal(err)
	}

	m.targetSprite = tar
	m.targetObj = resolv.NewObject(0, 0, float64(m.targetSprite.Bounds().Dx()), float64(m.targetSprite.Bounds().Dy()), "target")

	m.magneticFieldSize = initialMagneticFieldSize
	m.attractionStrength = initialAttractionStrength

	m.magFieldPhysObj = resolv.NewObject(-m.magneticFieldSize, -m.magneticFieldSize, float64(m.sprite.Bounds().Dx())+(m.magneticFieldSize*2), float64(m.sprite.Bounds().Dy())+(m.magneticFieldSize*2), "magneticField")

	m.magneticPoint.X = float64(magPhysObjSizeDiff / 2)
	m.magneticPoint.Y = float64(m.sprite.Bounds().Dy() - (magPhysObjSizeDiff))

	m.touch = false
	m.connected = false
	m.dropConnected = false
	m.magnetActive = false
	m.dropCounter = 0
	m.rotation = float64(float64(90) / float64(180) * math.Pi)
}

func (m *MagnetObject) ReadInput() {
	x, y := ebiten.CursorPosition()
	m.targetPos.X = float64(x)
	m.targetPos.Y = float64(y)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		m.magnetActive = true
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		m.dropConnected = true
	} else {
		m.dropConnected = false
	}
}

func (m *MagnetObject) Update(deltaTime float64) {
	m.magnetPos.X = m.magnetStartPos.X - (m.physObj.W / 2)
	m.magnetPos.Y = m.magnetStartPos.Y

	dx := m.magnetPos.X - m.physObj.X
	dy := m.magnetPos.Y - m.physObj.Y

	if m.magnetActive {
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
					v := basics.Vector2f{X: m.connectedJunk.X, Y: m.connectedJunk.Y}
					m.rotation = m.RotateTo(v)
				} else {
					m.touch = false
				}
			}
		}
	} else {
		m.targetObj.X = m.targetPos.X
		m.targetObj.Y = m.targetPos.Y
	}

	r := basics.Vector2f{X: m.targetObj.X, Y: m.targetObj.Y}
	m.rotation = m.RotateTo(r)

	if len(m.attractedJunk) > 0 {
		for _, junk := range m.attractedJunk {
			junk.X = basics.FloatLerp(junk.X, m.magneticPoint.X+m.physObj.X, m.attractionStrength*deltaTime)
			junk.Y = basics.FloatLerp(junk.Y, m.magneticPoint.Y+m.physObj.Y, m.attractionStrength*deltaTime)
		}
	}

	m.physObj.X = m.magnetPos.X
	m.physObj.Y = m.magnetPos.Y

	fieldOffset := basics.Vector2f{X: -m.magneticFieldSize - (magPhysObjSizeDiff / 2), Y: -m.magneticFieldSize - (magPhysObjSizeDiff / 2)}

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
		v := basics.Vector2f{X: m.physObj.X, Y: m.physObj.Y + 1}
		m.rotation = m.RotateTo(v)
	}

	m.physObj.Update()
	m.targetObj.Update()
	m.magFieldPhysObj.Update()
}

func (m *MagnetObject) Draw(screen *ebiten.Image) {
	sop := &ebiten.DrawImageOptions{}
	sop.GeoM.Translate(-float64(m.sprite.Bounds().Dx())/2, -float64(m.sprite.Bounds().Dy())/2)
	sop.GeoM.Rotate(m.rotation - float64(float64(90)/float64(180)*math.Pi))
	sop.GeoM.Translate(float64(m.sprite.Bounds().Dx())/2, float64(m.sprite.Bounds().Dy())/2)

	top := &ebiten.DrawImageOptions{}
	top.GeoM.Translate(m.targetObj.X-(float64(m.targetSprite.Bounds().Dx())/2), m.targetObj.Y-(float64(m.targetSprite.Bounds().Dy())/2))

	// Sprite is put over the top of the phys object
	sop.GeoM.Translate(m.physObj.X-(magPhysObjSizeDiff/2), m.physObj.Y-(magPhysObjSizeDiff/2))

	if m.touch {
		sop.ColorM.Scale(0.5, 1, 1, 1)
	}

	// Debug drawing of the physics object
	if globals.Debug {
		ebitenutil.DrawRect(screen, m.magFieldPhysObj.X, m.magFieldPhysObj.Y, m.magFieldPhysObj.W, m.magFieldPhysObj.H, color.RGBA{255, 80, 0, 64})
		ebitenutil.DrawRect(screen, m.physObj.X, m.physObj.Y, m.physObj.W, m.physObj.H, color.RGBA{255, 80, 0, 64})
		if m.magnetActive {
			ebitenutil.DrawLine(screen, m.magnetStartPos.X, m.magnetStartPos.Y, m.targetObj.X, m.targetObj.Y, color.RGBA{255, 0, 0, 64})
		}
	}

	// Draw the image (comment this out to see the above resolv rect ^^^)
	screen.DrawImage(m.sprite, sop)
	if !m.magnetActive {
		screen.DrawImage(m.targetSprite, top)
	}
}

func (m *MagnetObject) SetPosition(position basics.Vector2f) {
	m.physObj.X = position.X
	m.physObj.Y = position.Y
}

func (m *MagnetObject) RotateTo(position basics.Vector2f) float64 {
	var angle = math.Atan2(position.Y-m.physObj.Y, position.X-m.physObj.X)
	return angle //* (180 / math.Pi)
}
