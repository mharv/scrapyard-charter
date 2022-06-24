package entities

import (
	"fmt"
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
	imageRotation      float64
	attractionStrength float64
	magneticFieldSize  float64
	magneticPoint      basics.Vector2f
	linkDistance       basics.Vector2f
	magnetPos          *basics.Vector2f
	magnetStartPos     *basics.Vector2f
	magnetEndPos       basics.Vector2f
	targetPos          basics.Vector2f
	touch, connected   bool
	magnetActive       bool
	dropConnected      bool
	retract            bool
	syncToRod          bool
	turnedOn           bool
	alive              bool
	dropCounter        float64
	junkLookup         map[*resolv.Object]*JunkObject
}

const (
	magPhysObjSizeDiff = 20
)

func (m *MagnetObject) GetFishingLinePoint() basics.Vector2f {
	value := basics.Vector2f{X: m.physObj.X + (m.physObj.W / 2), Y: m.physObj.Y}
	value = basics.FloatRotAroundPoint(value, basics.Vector2f{X: m.physObj.X + (m.physObj.W / 2), Y: m.physObj.Y + (m.physObj.H / 2)}, m.rotation-(float64(90)/float64(180)*math.Pi))
	return value
}

func (m *MagnetObject) GetMagnetOffset() basics.Vector2f {
	return basics.Vector2f{X: m.physObj.W / 2, Y: 0}
}

func (m *MagnetObject) SetMagnetStartPos(pos *basics.Vector2f) {
	m.magnetStartPos = pos
}

func (m *MagnetObject) GetMagnetPos() *basics.Vector2f {
	return m.magnetPos
}
func (m *MagnetObject) GetTargetPos() *basics.Vector2f {
	return &basics.Vector2f{X: m.physObj.X, Y: m.physObj.Y}
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
	m.alive = true
	// Load an image given a filepath
	img, _, err := ebitenutil.NewImageFromFile(ImageFilepath)
	if err != nil {
		log.Fatal(err)
	}

	m.sprite = img
	m.physObj = resolv.NewObject(0, 0, float64(m.sprite.Bounds().Dx())-magPhysObjSizeDiff, float64(m.sprite.Bounds().Dy())-magPhysObjSizeDiff, "magnet")
	m.magnetPos = &basics.Vector2f{X: m.physObj.X, Y: m.physObj.Y}

	tar, _, err := ebitenutil.NewImageFromFile("images/target.png")
	if err != nil {
		log.Fatal(err)
	}

	m.targetSprite = tar
	m.targetObj = resolv.NewObject(0, 0, float64(m.targetSprite.Bounds().Dx()), float64(m.targetSprite.Bounds().Dy()), "target")

	m.magneticFieldSize = globals.GetPlayerData().GetMagneticFieldSize()
	m.attractionStrength = globals.GetPlayerData().GetAttractionStrength()

	m.magFieldPhysObj = resolv.NewObject(-m.magneticFieldSize, -m.magneticFieldSize, float64(m.sprite.Bounds().Dx())+(m.magneticFieldSize*2), float64(m.sprite.Bounds().Dy())+(m.magneticFieldSize*2), "magneticField")

	m.magneticPoint.X = float64(magPhysObjSizeDiff / 2)
	m.magneticPoint.Y = float64(m.sprite.Bounds().Dy() - (magPhysObjSizeDiff))

	m.touch = false
	m.connected = false
	m.dropConnected = false
	m.magnetActive = false
	m.retract = false
	m.syncToRod = true
	m.turnedOn = true
	m.dropCounter = 0
	m.rotation = float64(float64(90) / float64(180) * math.Pi)
	m.imageRotation = m.rotation
}

func (m *MagnetObject) ReadInput() {
	x, y := ebiten.CursorPosition()
	m.targetPos.X = float64(x)
	m.targetPos.Y = float64(y)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		end := basics.Vector2f{X: m.targetPos.X - m.magnetStartPos.X, Y: m.targetPos.Y - m.magnetStartPos.Y}
		basics.FloatNormalise(end)
		end.X *= globals.GetPlayerData().GetLineLength()
		end.Y *= globals.GetPlayerData().GetLineLength()

		m.attractedJunk = nil

		m.magnetEndPos = end
		m.magnetActive = true
		m.retract = false
		m.syncToRod = false
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		m.dropConnected = true
	} else {
		m.dropConnected = false
	}

}

func (m *MagnetObject) Update(deltaTime float64) {
	trackingPoint := basics.Vector2f{X: m.magnetStartPos.X - (m.physObj.W / 2), Y: m.magnetStartPos.Y}

	if m.syncToRod {
		m.magnetPos.X = trackingPoint.X
		m.magnetPos.Y = trackingPoint.Y
	}

	dx := m.magnetPos.X - m.physObj.X
	dy := m.magnetPos.Y - m.physObj.Y

	if m.dropCounter > 0 {
		m.dropCounter -= deltaTime
	} else {
		if !m.connected && m.turnedOn {
			if collision := m.physObj.Check(dx, dy, "junk"); collision != nil {
				m.connectedJunk = collision.Objects[0]
				m.linkDistance.X = m.connectedJunk.X - m.physObj.X
				m.linkDistance.Y = m.connectedJunk.Y - m.physObj.Y
				m.touch = true
				m.connected = true
				m.retract = true
				v := basics.Vector2f{X: m.connectedJunk.X, Y: m.connectedJunk.Y}
				m.rotation = m.RotateTo(v)
			} else {
				m.touch = false
			}
		}
	}

	if m.magnetActive && !m.retract {
		if basics.FloatDistance(*m.magnetPos, trackingPoint) < globals.GetPlayerData().GetLineLength() && !m.retract {
			newPos := m.MoveTowards(*m.magnetPos, m.magnetEndPos, globals.GetPlayerData().GetMagnetCastSpeed()*deltaTime)
			m.magnetPos.X += newPos.X
			m.magnetPos.Y += newPos.Y
		} else {
			m.retract = true
		}
	}
	if m.retract {
		if basics.FloatDistance(*m.magnetPos, trackingPoint) >= 5 && m.retract {
			newPos := m.MoveTowards(*m.magnetPos, trackingPoint, globals.GetPlayerData().GetMagnetReelSpeed()*deltaTime)
			m.magnetPos.X += newPos.X
			m.magnetPos.Y += newPos.Y
		} else {
			m.syncToRod = true
			m.magnetActive = false
			if m.connected {
				m.connected = false

				if val, ok := m.junkLookup[m.connectedJunk]; ok {
					if val.IsAlive() {
						fmt.Println("Adding " + val.itemData.GetName() + " to the bag")
						globals.GetPlayerData().GetInventory().AddItem(*val.GetItemData())

						val.Kill()
					}
				}

				m.connectedJunk = nil
			}
		}
	}

	if !m.connected {
		if !m.magnetActive || m.retract {
			r := basics.Vector2f{X: m.targetObj.X, Y: m.targetObj.Y}
			m.rotation = m.RotateTo(r)
		} else {
			r := basics.Vector2f{X: m.magnetEndPos.X, Y: m.magnetEndPos.Y}
			m.rotation = m.RotateTo(r)
		}
	} else {
		r := basics.Vector2f{X: m.connectedJunk.X, Y: m.connectedJunk.Y}
		m.rotation = m.RotateTo(r)
	}

	if m.magnetActive {
		if collision := m.magFieldPhysObj.Check(dx, dy, "junk"); collision != nil {
			m.attractedJunk = collision.Objects
		}
		m.MoveAttractedJunk(deltaTime)
	}

	fieldOffset := basics.Vector2f{X: -m.magneticFieldSize - (magPhysObjSizeDiff / 2), Y: -m.magneticFieldSize - (magPhysObjSizeDiff / 2)}

	m.SetObjPos(m.targetObj, m.targetPos)
	m.SetObjPos(m.physObj, *m.magnetPos)
	m.SetObjPos(m.magFieldPhysObj, *m.magnetPos, fieldOffset)

	if m.connected {
		m.SetObjPos(m.connectedJunk, *m.magnetPos, m.linkDistance)
	}

	if m.dropConnected {
		m.Drop()
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
	screen.DrawImage(m.targetSprite, top)
}

func (m *MagnetObject) IsAlive() bool {
	return m.alive
}

func (m *MagnetObject) Kill() {
	m.alive = false
}

func (m *MagnetObject) RemovePhysObj(space *resolv.Space) {
	space.Remove(m.physObj)
}

func (m *MagnetObject) Drop() {
	m.connectedJunk = nil
	m.linkDistance = basics.Vector2f{X: 0, Y: 0}
	m.connected = false
	m.touch = false
	m.dropCounter = globals.GetPlayerData().GetDropReactivationTimer()
	v := basics.Vector2f{X: m.physObj.X, Y: m.physObj.Y + 1}
	m.rotation = m.RotateTo(v)
}

func (m *MagnetObject) MoveAttractedJunk(deltaTime float64) {
	if len(m.attractedJunk) > 0 {
		for i, junk := range m.attractedJunk {
			if !junk.Overlaps(m.magFieldPhysObj) {
				if i >= len(m.attractedJunk) {
					m.attractedJunk = m.attractedJunk[:len(m.attractedJunk)-1]
				} else {
					m.attractedJunk = append(m.attractedJunk[:i], m.attractedJunk[i+1:]...)
				}
			} else {
				junk.X = basics.FloatLerp(junk.X, m.magneticPoint.X+m.physObj.X-1, m.attractionStrength*deltaTime)
				junk.Y = basics.FloatLerp(junk.Y, m.magneticPoint.Y+m.physObj.Y-1, m.attractionStrength*deltaTime)
			}
		}
	}
}

func (m *MagnetObject) SetJunkLookup(Lookup map[*resolv.Object]*JunkObject) {
	m.junkLookup = Lookup
}

func (m *MagnetObject) SetPosition(position basics.Vector2f) {
	m.physObj.X = position.X
	m.physObj.Y = position.Y
}

func (m *MagnetObject) RotateTo(position basics.Vector2f) float64 {
	var angle = math.Atan2(position.Y-m.physObj.Y, position.X-m.physObj.X)
	return angle
}

func (m *MagnetObject) MoveTowards(from, to basics.Vector2f, distance float64) basics.Vector2f {
	dir := basics.Vector2f{X: to.X - from.X, Y: to.Y - from.Y}
	dir = basics.FloatNormalise(dir)

	return basics.Vector2f{X: dir.X * distance, Y: dir.Y * distance}
}

func (m *MagnetObject) SetObjPos(obj *resolv.Object, vec ...basics.Vector2f) {
	newPos := basics.Vector2f{X: 0, Y: 0}

	for _, item := range vec {
		newPos.X += item.X
		newPos.Y += item.Y
	}

	obj.X = newPos.X
	obj.Y = newPos.Y
}
