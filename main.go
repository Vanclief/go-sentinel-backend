package main

import (
	"fmt"
	"log"

	"github.com/nareix/joy4/format"
	"github.com/nareix/joy4/format/rtmp"
)

func init() {
	format.RegisterAll()
}

func main() {
	server := &rtmp.Server{}

	rtmp.Debug = true

	server.HandlePlay = func(conn *rtmp.Conn) {
		// Here you can get the stream from a publisher and relay it to the player
		// The "conn" is now a source of video/audio data
		// conn.URL.Path is the path to the stream
	}

	server.HandlePublish = func(conn *rtmp.Conn) {
		fmt.Println("new connection")

		for {
			pkt, err := conn.ReadPacket()
			if err != nil {
				log.Fatal(err)
			}

			if pkt.IsKeyFrame {
				// Here you have a video packet. However, it's encoded.
				// If you want to do something with the raw frames, you'll need to decode it.
				log.Printf("Received a video packet with size %d and timestamp %d", len(pkt.Data), pkt.Time)
			}
		}

		// This is called when someone wants to publish a stream to the server
		// The "conn" is now a sink for video/audio data
		// conn.URL.Path is the path to the stream
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
