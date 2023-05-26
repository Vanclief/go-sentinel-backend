package main

import (
	"encoding/csv"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"time"

	"github.com/vanclief/go-sentinel-backend/player"
	"github.com/vanclief/go-sentinel-backend/scanner"
)

func main() {

	//
	sound := player.New()
	s := scanner.New(true, 10)

	// Open the CSV file in append mode.
	file, err := os.OpenFile("scanned.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	// Create a CSV writer.
	writer := csv.NewWriter(file)

	// Channel for listening results
	go func() {
		var lastResult string
		for result := range s.ResultStream {
			fmt.Println("Scanned:", result)
			sound.Play()

			if result.String() == lastResult {
				continue
			}

			lastResult = result.String()
			if err := writer.Write([]string{result.String()}); err != nil {
				log.Fatalf("Failed to write to CSV file: %v", err)
			}
			writer.Flush()
		}
	}()

	// Pooling to scan folder
	ticker := time.NewTicker(500 * time.Millisecond)
	for range ticker.C {
		s.ScanDir("./tmp/")
	}
}
