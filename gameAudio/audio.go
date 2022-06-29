package gameAudio

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/mharv/scrapyard-charter/resources"
)

const (
	sampleRate = 44100
)

type Audio struct {
	audioContext *audio.Context
	audioPlayer  map[string]*audio.Player
}

func (a *Audio) Init() {
	a.audioContext = audio.NewContext(sampleRate)
	a.audioPlayer = make(map[string]*audio.Player)
}

func (a *Audio) LoadFiles(folder string) {
	a.audioPlayer = resources.LoadFolderAsAudio(folder, a.audioContext)
}

func (a *Audio) PlayFile(filepath string) {
	if !a.audioPlayer[filepath].IsPlaying() {
		a.audioPlayer[filepath].Rewind()
		a.audioPlayer[filepath].Play()
	}
}

func (a *Audio) StopAllAudio() {
	for _, v := range a.audioPlayer {
		if v.IsPlaying() {
			v.Pause()
			v.Rewind()
		}
	}
}
