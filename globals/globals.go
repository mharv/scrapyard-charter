package globals

import (
	"github.com/mharv/scrapyard-charter/basics"
	"github.com/mharv/scrapyard-charter/data"
	"github.com/mharv/scrapyard-charter/gameAudio"
)

const (
	ScreenWidth  = 1366
	ScreenHeight = 768
	Debug        = false
)

var playerData = &data.PlayerData{InitialOverworldPosition: basics.Vector2f{
	X: ScreenWidth / 2,
	Y: ScreenHeight / 2,
}}

var audioPlayer = &gameAudio.Audio{}

func InitAudioPlayer() {
	audioPlayer.Init()
	audioPlayer.LoadFiles("audio")
}

func GetAudioPlayer() *gameAudio.Audio {
	return audioPlayer
}

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
