package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

func main() {
	// open and decode image file
	file, err := os.Open("qr.png")
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("here")
		panic(err)
	}

	// prepare BinaryBitmap
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		fmt.Println("BinaryBitmap")
		panic(err)
	}

	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		fmt.Println("QRReader")
		panic(err)
	}

	fmt.Println(result)
}
