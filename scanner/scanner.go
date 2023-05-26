package scanner

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/oned"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/vanclief/ez"
)

type Scanner struct {
	QRReader     gozxing.Reader
	OnedReader   gozxing.Reader
	DeleteOnScan bool
	Semaphore    chan bool
	ResultStream chan *gozxing.Result
	Debug        bool
}

func New(deleteOnScan bool, maxGoroutines int) *Scanner {

	return &Scanner{
		Semaphore:    make(chan bool, maxGoroutines),
		ResultStream: make(chan *gozxing.Result),
		QRReader:     qrcode.NewQRCodeReader(),
		OnedReader:   oned.NewMultiFormatUPCEANReader(nil),
		DeleteOnScan: deleteOnScan,
	}
}

func (s *Scanner) ScanDir(dir string) error {
	const op = "Scanner.ScanDir"

	files, err := os.ReadDir(dir)
	if err != nil {
		return ez.Wrap(op, err)
	}

	for _, f := range files {
		path := filepath.Join(dir, f.Name())
		if f.IsDir() || !isImage(path) {
			continue
		}

		s.Semaphore <- true

		go func(path string) {
			defer func() { <-s.Semaphore }()
			res, err := s.ScanImage(path)
			if err == nil && res != nil {
				s.ResultStream <- res
			}
		}(path)
	}

	return nil
}

func isImage(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif"
}

func (s *Scanner) ScanImage(path string) (*gozxing.Result, error) {
	const op = "scanImage"

	if s.Debug {
		start := time.Now()
		defer fmt.Println("Scanned:", path, "in", time.Since(start))
	}

	// load image
	bmp, err := s.load_image(path)
	if err != nil {
		return nil, ez.Wrap(op, err)
	}

	result, err := s.scanQR(bmp)
	if err != nil {
		result, err = s.scanBarcode(bmp)
		if err != nil {
			return nil, ez.Wrap(op, err)
		}
	}

	return result, nil
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
