package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"time"

	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/vanclief/go-sentinel-backend/scanner"

	"log"
)

func main() {

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

	defer streamer.Close()

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	s := scanner.New(true, 10)

	go func() {
		for res := range s.ResultStream {
			fmt.Println("Got result:", res)
			speaker.Play(streamer)
			// Handle the result as needed
		}
	}()

	ticker := time.NewTicker(500 * time.Millisecond)
	for range ticker.C {
		fmt.Println("Scan dir")
		s.ScanDir("./tmp/")
	}
}
