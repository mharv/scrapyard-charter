package data

type PlayerData struct {
	//ScavPlayerObject
	moveSpeedModifier float64
	rodStartXModifier int
	rodStartYModifier int
	rodEndXModifier   int
	rodEndYModifier   int
	//MagnetObject
	dropReactivationTimerModifier float64
	magneticFieldSizeModifier     int
	attractionStrengthModifier    int
	lineLengthModifier            float64
	magnetCastSpeedModifier       float64
	magnetReelSpeedModifier       float64
}

const (
	//ScavPlayerObject
	initialMoveSpeed = 300
	initialRodEndX   = 150
	initialRodEndY   = 25
	initialRodStartX = 52
	initialRodStartY = 52
	//MagnetObject
	initialDropReactivationTimer = 0.5
	initialMagneticFieldSize     = 50
	initialAttractionStrength    = 20
	initialLineLength            = 800.0
	initialMagnetCastSpeed       = 400.0
	initialMagnetReelSpeed       = 600.0
)

func (p *PlayerData) GetMoveSpeed() float64 {
	return initialMoveSpeed + p.moveSpeedModifier
}

func (p *PlayerData) GetRodStartX() int {
	return initialRodStartX + p.rodStartXModifier
}

func (p *PlayerData) GetRodStartY() int {
	return initialRodStartY + p.rodStartYModifier
}

func (p *PlayerData) GetRodEndX() int {
	return initialRodEndX + p.rodEndXModifier
}

func (p *PlayerData) GetRodEndY() int {
	return initialRodEndY + p.rodEndYModifier
}

func (p *PlayerData) GetDropReactivationTimer() float64 {
	return initialDropReactivationTimer + p.dropReactivationTimerModifier
}

func (p *PlayerData) GetMagneticFieldSize() int {
	return initialMagneticFieldSize + p.magneticFieldSizeModifier
}

func (p *PlayerData) GetAttractionStrength() int {
	return initialAttractionStrength + p.attractionStrengthModifier
}

func (p *PlayerData) GetLineLength() float64 {
	return initialLineLength + p.lineLengthModifier
}

func (p *PlayerData) GetMagnetCastSpeed() float64 {
	return initialMagnetCastSpeed + p.magnetCastSpeedModifier
}

func (p *PlayerData) GetMagnetReelSpeed() float64 {
	return initialMagnetReelSpeed + p.magnetReelSpeedModifier
}
