package globals

import (
	"github.com/mharv/scrapyard-charter/basics"
	"github.com/mharv/scrapyard-charter/data"
)

const (
	ScreenWidth  = 1366
	ScreenHeight = 768
	Debug        = true
)

var playerData = &data.PlayerData{InitialOverworldPosition: basics.Vector2f{
	X: ScreenWidth / 2,
	Y: ScreenHeight / 2,
}}

func GetPlayerData() *data.PlayerData {
	return playerData
}

var MaterialNamesList []string = []string{
	"Iron",
	"Steel",
	"Copper",
	"Rubber",
	"Plastic",
	"Nickel",
	"Cobalt",
	"Titanium",
	"Gold",
}
