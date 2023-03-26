package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/discord/lilliput"
)

var ops = lilliput.NewImageOps(8192)

func getEncodeOptions(format string, quality int) map[int]int {
	if format == "jpeg" || format == "jpg" {
		return map[int]int{lilliput.JpegQuality: quality, lilliput.JpegProgressive: 1}
	} else if format == "png" {
		return map[int]int{lilliput.PngCompression: quality}
	} else if format == "webp" {
		return map[int]int{lilliput.WebpQuality: quality}
	} else {
		return map[int]int{}
	}
}

func ProcessImage(image []byte, format string, width int, height int, quality int) ([]byte, error) {
	// check if the format has a dot and remove the dot from the format
	if format[0] == '.' {
		format = format[1:]
	}

	decoder, err := lilliput.NewDecoder(image)
	// this error reflects very basic checks,
	// mostly just for the magic bytes of the file to match known image formats
	if err != nil {
		fmt.Printf("error decoding image, %s\n", err)
		os.Exit(1)
	}
	defer decoder.Close()

	header, err := decoder.Header()
	if err != nil {
		fmt.Printf("error getting image header, %s\n", err)
		os.Exit(1)
	}

	fmt.Println(len(image), header.Width(), header.Height(), header.Orientation())
	// get ready to resize image,
	// using 8192x8192 maximum resize buffer size
	// ops := lilliput.NewImageOps(8192)
	// defer ops.Close()

	fmt.Println(format, width, height, quality)

	opts := &lilliput.ImageOptions{
		FileType:             "." + format,
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
	// PrintMemUsage()
	// runtime.GC()

	PrintMemUsage("After image processing")
	return outputImg, err
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage(location string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats

	fmt.Print(location, " =>")
	fmt.Printf("\tAlloc = %v MiB,", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB,", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB,", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
