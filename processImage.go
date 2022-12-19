package main

import (
	"fmt"
	"os"
	"github.com/discord/lilliput"
)

func getEncodeOptions (format string, quality int) map[int]int {
	if format == "jpeg" || format == "jpg" || format == ".jpeg" || format == ".jpg"  {
		return map[int]int{lilliput.JpegQuality: quality}
	} else if format == "png" {
		return map[int]int{lilliput.PngCompression: quality}
	} else if format == "webp" {
		return map[int]int{lilliput.WebpQuality: quality}
	} else {
		return map[int]int{}
	}
}

func ProcessImage(image []byte, format string, width int, height int, quality int) ([]byte, error) {
	decoder, err := lilliput.NewDecoder(image)
	// this error reflects very basic checks,
	// mostly just for the magic bytes of the file to match known image formats
	if err != nil {
		fmt.Printf("error decoding image, %s\n", err)
		os.Exit(1)
	}
	defer decoder.Close()

	// get ready to resize image,
	// using 8192x8192 maximum resize buffer size
	ops := lilliput.NewImageOps(8192)
	defer ops.Close()

	fmt.Println(format, width, height, quality)

	opts := &lilliput.ImageOptions{
		FileType:             format,
		Width:                width,
		Height:               height,
		ResizeMethod:         lilliput.ImageOpsFit,
		NormalizeOrientation: true,
		EncodeOptions:        getEncodeOptions(format, quality),
	}

	// create a buffer to store the output image, 50MB in this case
	outputImg := make([]byte, 50*1024*1024)
	
	// resize and transcode image
	outputImg, err = ops.Transform(decoder, opts, outputImg)
	if err != nil {
		fmt.Printf("error transforming image, %s\n", err)
		panic(err)
	}

	return outputImg, err
}
