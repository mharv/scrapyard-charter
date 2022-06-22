package data

type PlayerData struct {
	//ScavPlayerObject
	scavMoveSpeedModifier float64
	rodStartXModifier     float64
	rodStartYModifier     float64
	rodEndXModifier       float64
	rodEndYModifier       float64
	//MagnetObject
	dropReactivationTimerModifier float64
	magneticFieldSizeModifier     float64
	attractionStrengthModifier    float64
	lineLengthModifier            float64
	magnetCastSpeedModifier       float64
	magnetReelSpeedModifier       float64
	//overworldPlayer
	overworldMoveSpeedModifier float64
}

const (
	//ScavPlayerObject
	initialScavMoveSpeed = 250
	initialRodEndX       = 150
	initialRodEndY       = 25
	initialRodStartX     = 52
	initialRodStartY     = 52
	//MagnetObject
	initialDropReactivationTimer = 0
	initialMagneticFieldSize     = 50
	initialAttractionStrength    = 20
	initialLineLength            = 200
	initialMagnetCastSpeed       = 350
	initialMagnetReelSpeed       = 400
	//overworldPlayer
	initialOverworldMoveSpeed = 200
)

func (p *PlayerData) GetOverworldMoveSpeed() float64 {
	return initialOverworldMoveSpeed + p.overworldMoveSpeedModifier
}

func (p *PlayerData) GetScavMoveSpeed() float64 {
	return initialScavMoveSpeed + p.scavMoveSpeedModifier
}

func (p *PlayerData) GetRodStartX() float64 {
	return initialRodStartX + p.rodStartXModifier
}

func (p *PlayerData) GetRodStartY() float64 {
	return initialRodStartY + p.rodStartYModifier
}

func (p *PlayerData) GetRodEndX() float64 {
	return initialRodEndX + p.rodEndXModifier
}

func (p *PlayerData) GetRodEndY() float64 {
	return initialRodEndY + p.rodEndYModifier
}

func (p *PlayerData) GetDropReactivationTimer() float64 {
	return initialDropReactivationTimer + p.dropReactivationTimerModifier
}

func (p *PlayerData) GetMagneticFieldSize() float64 {
	return initialMagneticFieldSize + p.magneticFieldSizeModifier
}

func (p *PlayerData) GetAttractionStrength() float64 {
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
