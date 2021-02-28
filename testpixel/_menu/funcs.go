package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func loadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}

type MusicInfo struct {
	name string
	path string
}

func loadMusic() map[string]*beep.Buffer {
	res := make(map[string]*beep.Buffer)
	list := []MusicInfo{
		MusicInfo{name: "monster", path: "music\\monster.mp3"},
		MusicInfo{name: "fear", path: "music\\fear.mp3"},
		MusicInfo{name: "anxiety", path: "music\\anxiety.mp3"},
		MusicInfo{name: "laser", path: "music\\laser.mp3"},
	}

	isInited := false
	for _, mi := range list {
		f, err := os.Open(mi.path)
		if err != nil {
			log.Fatal(err)
		}

		streamer, format, err := mp3.Decode(f)
		if err != nil {
			log.Fatal(err)
		}

		if !isInited {
			speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
			isInited = true
		}

		buffer := beep.NewBuffer(format)
		buffer.Append(streamer)
		streamer.Close()

		res[mi.name] = buffer
		f.Close()
	}

	return res
}
