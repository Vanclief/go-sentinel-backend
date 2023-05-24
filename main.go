package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/multi/qrcode"
)

func main() {
	// open and decode image file
	file, err := os.Open("barcode.png")
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
	qrReader := qrcode.NewQRCodeMultiReader()
	result, err := qrReader.DecodeMultipleWithoutHint(bmp)
	if err != nil {
		fmt.Println("QRReader Error", err)
	}

	fmt.Println(result)
}
