package gameAudio

import (
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

func (a *Audio) LoadFile(filepath string) {
	f, err := ebitenutil.OpenFile(filepath)
	if err != nil {
		panic("Cannot open file: " + filepath)
	}

	d, err := mp3.DecodeWithSampleRate(sampleRate, f)
	if err != nil {
		panic("Cannot decode file: " + filepath)
	}

	p, err := a.audioContext.NewPlayer(d)
	if err != nil {
		panic("Cannot create player for: " + filepath)
	}

	a.audioPlayer[filepath] = p
}

func (a *Audio) LoadFiles(folder string) {

	f, err := os.Open(folder)
	if err != nil {
		fmt.Println(err)
		return
	}
	files, err := f.Readdir(0)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range files {
		filepath := folder + "/" + v.Name()
		a.LoadFile(filepath)
	}
}

func (a *Audio) PlayFile(filepath string) {
	if !a.audioPlayer[filepath].IsPlaying() {
		a.audioPlayer[filepath].Rewind()
		a.audioPlayer[filepath].Play()
	}
}
