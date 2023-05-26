package main

import (
	"encoding/csv"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net"
	"os"
	"time"

	"github.com/vanclief/go-sentinel-backend/player"
	"github.com/vanclief/go-sentinel-backend/scanner"
)

func main() {

	fmt.Println("GeoScanner 0.1.0")
	print_local_ip()

	// Instantiate a new player and scanner
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

func print_local_ip() {
	// Fetch all network interfaces.
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Print(err)
		return
	}

	for _, i := range interfaces {
		// Fetch all associated addresses for this interface.
		addresses, err := i.Addrs()
		if err != nil {
			fmt.Print(err)
			return
		}

		// Iterate over all addresses and print the first non-loopback IPv4 address.
		for _, address := range addresses {
			// Check if the address is an IP net address type and is not a loopback address.
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				// Check if this is an IPv4 address.
				if ipnet.IP.To4() != nil {
					fmt.Printf("Connect to: rtmp://%s:1935/dji\n", ipnet.IP.String())
					return
				}
			}
		}
	}
}
