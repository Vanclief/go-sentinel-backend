package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/vanclief/go-sentinel-backend/scanner"

	"log"
)

func main() {

	player := New()

	s := scanner.New(true, 10)

	go func() {
		for res := range s.ResultStream {
			fmt.Println("Got result:", res)
			player.Play()
			// scanSound := buffer.Streamer(0, buffer.Len())
			// speaker.Play(scanSound)
			// Handle the result as needed
		}
	}()

	ticker := time.NewTicker(500 * time.Millisecond)
	for range ticker.C {
		s.ScanDir("./tmp/")
	}
}

type SoundPlayer struct {
	Beep  beep.StreamSeeker
	Queue chan bool
}

func New() *SoundPlayer {

	f, err := os.Open("sounds/beep2.mp3")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	buffer := beep.NewBuffer(format)

	buffer.Append(streamer)
	streamer.Close()

	queue := make(chan bool)

	go func() {
		var isPlaying bool
		for range queue {
			if !isPlaying {
				isPlaying = true
				streamSeeker := buffer.Streamer(0, buffer.Len())
				speaker.Play(beep.Seq(streamSeeker, beep.Callback(func() {
					isPlaying = false
				})))
			}

		}
	}()

	return &SoundPlayer{
		Queue: queue,
	}
}

func (s *SoundPlayer) Play() {
	s.Queue <- true
}
