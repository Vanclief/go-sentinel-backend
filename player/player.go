package player

import (
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

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
