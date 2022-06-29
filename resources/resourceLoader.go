package resources

import (
	"embed"
	"image"
	"io"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
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

func LoadFolderAsAudio(Folder string, context *audio.Context) map[string]*audio.Player {

	streams := map[string]*audio.Player{}
	//var files embed.FS
	folder, err := AudioFS.ReadDir(Folder)
	if err != nil {
		panic(err)
	}

	for _, v := range folder {
		name := Folder + "/" + v.Name()

		file, err := AudioFS.Open(name)
		if err != nil {
			panic(err)
		}

		d, err := mp3.DecodeWithSampleRate(SampleRate, file)
		if err != nil {
			panic(err)
		}

		p, err := context.NewPlayer(d)
		if err != nil {
			panic("Cannot create player for: " + name)
		}
		streams[name] = p
	}

	return streams
}
