package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/discord/lilliput"
)

func main() {
	handler := http.HandlerFunc(handleRequest)
	http.Handle("/photo", handler)
	http.ListenAndServe(":8080", nil)
}

type ImageFormat struct {
	Format string
	Width  int
	Height int
	Quality int
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	fileBytes, err := ioutil.ReadFile("test.png")
	if err != nil {
		panic(err)
	}

	// read query params into struct
	format := r.URL.Query().Get("fm")
	width := r.URL.Query().Get("w")
	height := r.URL.Query().Get("h")
	quality := r.URL.Query().Get("q")	

	// convert query params to int
	widthInt, err := strconv.Atoi(width)
	if err != nil {
		panic(err)
	}

	heightInt, err := strconv.Atoi(height)
	if err != nil {
		panic(err)
	}

	qualityInt, err := strconv.Atoi(quality)
	if err != nil {
		panic(err)
	}

	// process image
	newImage, err := ProcessImage(fileBytes, format, widthInt, heightInt, qualityInt)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(newImage)
}

func getEncodeOptions (format string, quality int) map[int]int {
	if format == "jpeg" || format == "jpg" {
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
