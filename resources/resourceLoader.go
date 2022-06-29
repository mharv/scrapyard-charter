package resources

import (
	"embed"
	"image"
	"io"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/tinne26/etxt"
)

const (
	SampleRate = 44100
)

//go:embed images/*
var ImagesFS embed.FS

//go:embed audio/*
var AudioFS embed.FS

//go:embed fonts/*
var FontsFS embed.FS

func LoadFileAsImage(Filename string) *ebiten.Image {
	file, err := ImagesFS.Open(Filename)
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

func LoadFileAsFont(Filename string) *etxt.FontLibrary {
	file, err := FontsFS.Open(Filename)
	if err != nil {
		panic(err)
	}

	bs, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	fontLib := etxt.NewFontLibrary()
	fontLib.ParseFontBytes(bs)

	return fontLib
}

func LoadFileAsAudio(Filename string) *mp3.Stream {
	file, err := AudioFS.Open(Filename)
	if err != nil {
		panic(err)
	}

	decoded, err := mp3.DecodeWithSampleRate(SampleRate, file)
	if err != nil {
		panic(err)
	}

	return decoded
}
