package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func main() {
	f, err := os.Open("laser.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	buffers := make([]*beep.Buffer, 0)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	buffers = append(buffers, beep.NewBuffer(format))
	buffers[0].Append(streamer)
	streamer.Close()

	f2, err := os.Open("anxiety.mp3")
	if err != nil {
		log.Fatal(err)
	}
	streamer2, format, err := mp3.Decode(f2)
	if err != nil {
		log.Fatal(err)
	}
	buffers = append(buffers, beep.NewBuffer(format))
	buffers[1].Append(streamer2)
	streamer2.Close()

	oddEven := 0
	for {
		fmt.Println("Press [ENTER] to fire!")
		fmt.Scanln()

		oddEven++
		oddEven %= 2

		shot := buffers[oddEven].Streamer(0, buffers[oddEven].Len())
		speaker.Play(shot)
	}
}
