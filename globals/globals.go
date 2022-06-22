package globals

import "github.com/mharv/scrapyard-charter/data"

const (
	ScreenWidth  = 1366
	ScreenHeight = 768
	Debug        = true
)

var playerData = &data.PlayerData{}

func GetPlayerData() *data.PlayerData {
	return playerData
}
