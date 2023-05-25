package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/multi/qrcode"

	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func main() {
	// scan_image("./tmp/frame_3.png")
	watch()
}

func watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write {
					if isImage(event.Name) {
						fmt.Println("Scan frame", event.Name)
						scan_image(event.Name)
						os.Remove(event.Name)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("./tmp/")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func isImage(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif"
}

func openImage(path string) {
	cmd := exec.Command("open", path)
	err := cmd.Run()
	if err != nil {
		log.Printf("Error opening image: %v", err)
	}
}

func scan_image(path string) {

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

	fmt.Println("===== QR SCANNED =====")
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
