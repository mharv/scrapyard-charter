package entities

import (
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/mharv/scrapyard-charter/basics"
	"github.com/mharv/scrapyard-charter/globals"
)

type RodSection struct {
	sprite   *ebiten.Image
	position basics.Vector2f
	rotation float64
}

type ScavRodObject struct {
	sprite            *ebiten.Image
	rodSections       []RodSection
	root              *basics.Vector2f
	tip               *basics.Vector2f
	magnetPosition    *basics.Vector2f
	rootController    basics.Vector2f
	tipController     basics.Vector2f
	maxMagnetDistance *float64
	lineOffset        float64
	rodBaseFlex       float64
	rodTipFlex        float64
	rodTipMaxSlop     float64
	rodPoints         []basics.Vector2f
	linePoints        []basics.Vector2f
}

const (
	rodResolution  = 20
	lineResolution = 5
	rodBaseFlex    = 75
	rodTipFlex     = 25
	rodTipMaxSlop  = 5
)

func (s *ScavRodObject) GetSprite() *ebiten.Image {
	return s.sprite
}

func (s *ScavRodObject) SetRoot(rootPos *basics.Vector2f) {
	s.root = rootPos
}

func (s *ScavRodObject) SetTip(tipPos *basics.Vector2f) {
	s.tip = tipPos
}

func (s *ScavRodObject) SetMagnetPosition(magnetPosition *basics.Vector2f) {
	s.magnetPosition = magnetPosition
}

func (s *ScavRodObject) GetMaxMagnetDistance(magnetDistance *float64) {
	s.maxMagnetDistance = magnetDistance
}

func (s *ScavRodObject) GetCurrentTipPos() *basics.Vector2f {
	return s.tip
}

func (s *ScavRodObject) Init(ImageFilepath string) {
	// Load an image given a filepath
	img, _, err := ebitenutil.NewImageFromFile(ImageFilepath)
	if err != nil {
		log.Fatal(err)
	}

	s.sprite = img
	s.lineOffset = float64(s.sprite.Bounds().Dy())
	s.root = &basics.Vector2f{X: 250, Y: 250}
	s.tip = &basics.Vector2f{X: s.root.X + 250, Y: s.root.Y}
	s.rootController = basics.Vector2f{X: s.root.X + rodBaseFlex, Y: s.root.Y - rodBaseFlex}
	s.tipController = basics.Vector2f{X: s.tip.X, Y: s.tip.Y - rodTipFlex}
	s.CalculatePoints()
	for i := 0; i < rodResolution; i++ {
		s.rodSections = append(s.rodSections, RodSection{})
	}
	s.UpdateRodSections()
}

func (s *ScavRodObject) ReadInput() {
}

func (s *ScavRodObject) Update(deltaTime float64) {
	s.rootController = basics.Vector2f{X: s.root.X + rodBaseFlex, Y: s.root.Y - rodBaseFlex}
	s.AngleRodTip()
	//s.tipController = basics.Vector2f{X: s.tip.X, Y: s.tip.Y - rodTipFlex}

	s.UpdatePoints()
	s.UpdateRodSections()
}

func (s *ScavRodObject) Draw(screen *ebiten.Image) {
	s.DrawRodSections(screen)

	for i := range s.linePoints {
		if i > 0 {
			ebitenutil.DrawLine(screen, s.linePoints[i-1].X, s.linePoints[i-1].Y, s.linePoints[i].X, s.linePoints[i].Y, color.RGBA{255, 255, 255, 255})
		}
	}

	if globals.Debug {
		ebitenutil.DrawLine(screen, s.root.X, s.root.Y, s.rootController.X, s.rootController.Y, color.RGBA{255, 0, 0, 255})
		ebitenutil.DrawLine(screen, s.tip.X, s.tip.Y, s.tipController.X, s.tipController.Y, color.RGBA{255, 0, 0, 255})
	}
}

func (s *ScavRodObject) CalculatePoints() {
	s.rodPoints = append(s.rodPoints, *s.root)
	s.linePoints = append(s.linePoints, *s.root)

	linePointCount := rodResolution / lineResolution

	for i := 1; i < rodResolution; i++ {
		t := float64(i) / float64(rodResolution)

		x0 := (math.Pow(1-t, 3) * s.root.X)
		x1 := (3 * math.Pow(1-t, 2) * t * s.rootController.X)
		x2 := (3 * (1 - t) * math.Pow(t, 2) * s.tipController.X)
		x3 := (math.Pow(t, 3) * s.tip.X)

		y0 := (math.Pow(1-t, 3) * s.root.Y)
		y1 := (3 * math.Pow(1-t, 2) * t * s.rootController.Y)
		y2 := (3 * (1 - t) * math.Pow(t, 2) * s.tipController.Y)
		y3 := (math.Pow(t, 3) * s.tip.Y)

		x := x0 + x1 + x2 + x3
		y := y0 + y1 + y2 + y3

		s.rodPoints = append(s.rodPoints, basics.Vector2f{X: x, Y: y})

		if i%linePointCount == 0 {
			s.linePoints = append(s.linePoints, basics.Vector2f{X: x, Y: y + s.lineOffset})
		}
	}

	s.rodPoints = append(s.rodPoints, *s.tip)
	s.linePoints = append(s.linePoints, *s.tip)
}

func (s *ScavRodObject) UpdatePoints() {
	linePointCount := rodResolution / lineResolution

	s.linePoints[0] = *s.root

	for i := range s.rodPoints {
		t := float64(i) / float64(rodResolution)

		x0 := (math.Pow(1-t, 3) * s.root.X)
		x1 := (3 * math.Pow(1-t, 2) * t * s.rootController.X)
		x2 := (3 * (1 - t) * math.Pow(t, 2) * s.tipController.X)
		x3 := (math.Pow(t, 3) * s.tip.X)

		y0 := (math.Pow(1-t, 3) * s.root.Y)
		y1 := (3 * math.Pow(1-t, 2) * t * s.rootController.Y)
		y2 := (3 * (1 - t) * math.Pow(t, 2) * s.tipController.Y)
		y3 := (math.Pow(t, 3) * s.tip.Y)

		x := x0 + x1 + x2 + x3
		y := y0 + y1 + y2 + y3

		s.rodPoints[i] = basics.Vector2f{X: x, Y: y}

		if i%linePointCount == 0 {
			s.linePoints[int(i/lineResolution)+1] = basics.Vector2f{X: x, Y: y + s.lineOffset}
		}
	}

	s.rodPoints[len(s.rodPoints)-1] = *s.tip
	s.linePoints[len(s.linePoints)-1] = *s.tip
}

func (s *ScavRodObject) UpdateRodSections() {
	for i := range s.rodPoints {
		if i > 0 {
			len := basics.FloatDistance(s.rodPoints[i-1], s.rodPoints[i])

			len = basics.FloatClamp(len, 1, float64(s.sprite.Bounds().Size().X))

			angleVec := basics.Vector2f{X: s.rodPoints[i].X - s.rodPoints[i-1].X, Y: s.rodPoints[i].Y - s.rodPoints[i-1].Y}

			rect := image.Rect(0, 0, int(len+2), s.sprite.Bounds().Size().Y)

			s.rodSections[i-1].sprite = s.sprite.SubImage(rect).(*ebiten.Image)
			s.rodSections[i-1].rotation = math.Atan2(-angleVec.Y, -angleVec.X)
			s.rodSections[i-1].position = s.rodPoints[i-1]
		}
	}
}

func (s *ScavRodObject) DrawRodSections(screen *ebiten.Image) {
	for _, section := range s.rodSections {
		options := &ebiten.DrawImageOptions{}
		x, y := section.sprite.Size()
		options.GeoM.Translate(-float64(x), -float64(y))
		options.GeoM.Rotate(section.rotation)
		options.GeoM.Translate(section.position.X, section.position.Y)

		screen.DrawImage(section.sprite, options)
	}
}

func (s *ScavRodObject) UpdateTipPosition(screen *ebiten.Image) {
	//dist := basics.FloatDistance(*s.magnetPosition, *s.tip)

	//percentage := dist / *s.maxMagnetDistance

}

func (s *ScavRodObject) AngleRodTip() {
	ang := basics.AngleBetweenFloatVec(*s.tip, *s.magnetPosition)

	s.tipController = basics.FloatRotAroundPoint(basics.Vector2f{X: s.tip.X, Y: s.tip.Y - rodTipFlex}, *s.tip, ang)
}
