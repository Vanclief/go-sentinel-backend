package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"time"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/oned"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/vanclief/ez"

	"log"
	"path/filepath"
	"strings"
)

func main() {

	scanner := New(true, 10)
	// scanner.scanImage("./samples/upca-2/1.png")

	ticker := time.NewTicker(500 * time.Millisecond)
	for range ticker.C {
		scanner.scanDir(scanner, "./tmp/")
	}
}

func (s *Scanner) scanDir(scanner *Scanner, dir string) {

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		path := filepath.Join(dir, f.Name())
		if f.IsDir() || !isImage(path) {
			continue
		}

		s.Semaphore <- true

		go func(path string) {
			defer func() { <-s.Semaphore }()
			scanner.scanImage(path)
		}(path)
	}
}

func isImage(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif"
}

type Scanner struct {
	Semaphore    chan bool
	QRReader     gozxing.Reader
	OnedReader   gozxing.Reader
	DeleteOnScan bool
}

func New(deleteOnScan bool, maxGoroutines int) *Scanner {

	return &Scanner{
		Semaphore:    make(chan bool, maxGoroutines),
		QRReader:     qrcode.NewQRCodeReader(),
		OnedReader:   oned.NewMultiFormatUPCEANReader(nil),
		DeleteOnScan: deleteOnScan,
	}
}

func (s *Scanner) scanImage(path string) error {
	const op = "scanImage"

	start := time.Now()
	defer fmt.Println("Scanned:", path, "in", time.Since(start))

	// load image
	bmp, err := s.load_image(path)
	if err != nil {
		return ez.Wrap(op, err)
	}

	result, err := s.scanQR(bmp)
	if err != nil {
		result, err = s.scanBarcode(bmp)
		if err != nil {
			return ez.Wrap(op, err)
		}
	}

	fmt.Println("===== QR FOUND =====")
	fmt.Println(result)

	return nil
}

func (s *Scanner) scanBarcode(bmp *gozxing.BinaryBitmap) (*gozxing.Result, error) {
	const op = "scanBarcode"

	result, err := s.OnedReader.DecodeWithoutHints(bmp)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return result, nil
}

func (s *Scanner) scanQR(bmp *gozxing.BinaryBitmap) (*gozxing.Result, error) {
	const op = "scanQR"

	// decode image
	result, err := s.QRReader.DecodeWithoutHints(bmp)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return result, nil
}

func (s *Scanner) load_image(path string) (*gozxing.BinaryBitmap, error) {
	const op = "load_image"

	if s.DeleteOnScan {
		defer os.Remove(path)
	}

	// fmt.Println("path", path)

	file, err := os.Open(path)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	// prepare BinaryBitmap
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return bmp, nil
}
