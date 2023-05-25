package main

import (
	"fmt"
	"image"
	_ "image/png"
	"io/ioutil"
	"os"
	"time"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/multi/qrcode"

	"log"
	"path/filepath"
	"strings"
)

func main() {
	// scan_image("./tmp/frame_3.png")
	// watch()
	ticker := time.NewTicker(500 * time.Millisecond)
	for range ticker.C {
		scanDir("./tmp/")
	}
}

func scanDir(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		path := filepath.Join(dir, f.Name())
		if f.IsDir() || !isImage(path) {
			continue
		}

		start := time.Now()
		scan_image(path)
		elapsed := time.Since(start)
		fmt.Println("Scanned:", path, "in", elapsed)
	}
}

// func watch() {
// 	watcher, err := fsnotify.NewWatcher()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer watcher.Close()

// 	done := make(chan bool)
// 	go func() {
// 		for {
// 			select {
// 			case event, ok := <-watcher.Events:
// 				if !ok {
// 					return
// 				}

// 				if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write {
// 					if isImage(event.Name) {
// 						fmt.Println("Scan:", event.Name)
// 						scan_image(event.Name)
// 						os.Remove(event.Name)
// 					}
// 				}
// 			case err, ok := <-watcher.Errors:
// 				if !ok {
// 					return
// 				}
// 				log.Println("error:", err)
// 			}
// 		}
// 	}()

// 	err = watcher.Add("./tmp/")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	<-done
// }

func isImage(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif"
}

func scan_image(path string) {

	defer os.Remove(path)

	// fmt.Println("path", path)

	file, err := os.Open(path)
	if err != nil {
		// fmt.Println("Open err", err)
		return
	}

	img, _, err := image.Decode(file)
	if err != nil {
		// fmt.Println("Image err", err)
		return
	}

	// prepare BinaryBitmap
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		// fmt.Println("BinaryBitmap err", err)
		return
	}

	// decode image
	qrReader := qrcode.NewQRCodeMultiReader()
	result, err := qrReader.DecodeMultipleWithoutHint(bmp)
	if err != nil {
		// fmt.Println("QRReader Error", err)
		return
	}

	fmt.Println("===== QR FOUND =====")
	fmt.Println(result)

	// hints := make(map[gozxing.DecodeHintType]interface{})
	// possibleFormats := []gozxing.BarcodeFormat{
	// 	gozxing.BarcodeFormat_EAN_13,
	// 	gozxing.BarcodeFormat_EAN_8,
	// 	gozxing.BarcodeFormat_UPC_A,
	// 	gozxing.BarcodeFormat_UPC_E,
	// 	gozxing.BarcodeFormat_CODE_39,
	// 	gozxing.BarcodeFormat_CODE_93,
	// 	gozxing.BarcodeFormat_CODE_128,
	// }
	//
	// hints[gozxing.DecodeHintType_POSSIBLE_FORMATS] = possibleFormats

	// onedReader := oned.NewMultiFormatUPCEANReader(hints)
	// onedReader := oned.NewMultiFormatUPCEANReader(nil)
	// res2, err := onedReader.DecodeWithoutHints(bmp)
	// if err != nil {
	// 	fmt.Println("OnedReader Error", err)
	// }
	//
	// fmt.Println(res2)
}
