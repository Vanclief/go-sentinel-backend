package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
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

	scanner := New(true)
	// scanner.scan_image("./samples/upca-2/1.png")

	ticker := time.NewTicker(500 * time.Millisecond)
	for range ticker.C {
		scanDir(scanner, "./tmp/")
	}
}

func scanDir(scanner *Scanner, dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		path := filepath.Join(dir, f.Name())
		if f.IsDir() || !isImage(path) {
			continue
		}

		go scanner.scan_image(path)
	}
}

func isImage(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif"
}

type Scanner struct {
	QRReader     gozxing.Reader
	OnedReader   gozxing.Reader
	DeleteOnScan bool
}

func New(deleteOnScan bool) *Scanner {
	return &Scanner{
		QRReader:     qrcode.NewQRCodeReader(),
		OnedReader:   oned.NewMultiFormatUPCEANReader(nil),
		DeleteOnScan: deleteOnScan,
	}
}

func (s *Scanner) scan_image(path string) error {
	const op = "scan_image"

	start := time.Now()
	defer fmt.Println("Scanned:", path, "in", time.Since(start))

	// load image
	bmp, err := s.load_image(path)
	if err != nil {
		return ez.Wrap(op, err)
	}

	result, err := s.scan_qr(bmp)
	if err != nil {
		result, err = s.scan_barcode(bmp)
		if err != nil {
			return ez.Wrap(op, err)
		}
	}

	fmt.Println("===== QR FOUND =====")
	fmt.Println(result)

	return nil
}

func (s *Scanner) scan_barcode(bmp *gozxing.BinaryBitmap) (*gozxing.Result, error) {
	const op = "scan_barcode"

	result, err := s.OnedReader.DecodeWithoutHints(bmp)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	return result, nil
}

func (s *Scanner) scan_qr(bmp *gozxing.BinaryBitmap) (*gozxing.Result, error) {
	const op = "scan_qr"

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
