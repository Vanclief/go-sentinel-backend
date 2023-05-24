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
	qrReader := qrcode.NewQRCodeMultiReader()
	result, err := qrReader.DecodeMultipleWithoutHint(bmp)
	if err != nil {
		fmt.Println("QRReader Error", err)
	}

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
