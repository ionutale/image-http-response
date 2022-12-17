package main

import (
	"github.com/discord/lilliput"
)

var EncodeOptions = map[string]map[int]int{
	"jpeg": map[int]int{lilliput.JpegQuality: 90},
	"png":  map[int]int{lilliput.PngCompression: 9},
	".webp": map[int]int{lilliput.WebpQuality: 90},
}

func processImage(image []byte, format string, width int, height int) ([]byte, error) {
	// Decode the image
	decoder := lilliput.NewDecoder(image)
	img, err := decoder.Decode()
	if err != nil {
		return nil, err
	}

	// Resize the image
	resizer := lilliput.NewResizer(EncodeOptions)
	img, err = resizer.Resize(img, width, height)
	if err != nil {
		return nil, err
	}

	// Encode the image
	encoder := lilliput.NewEncoder(EncodeOptions)
	encodedImage, err := encoder.Encode(img, format)
	if err != nil {
		return nil, err
	}

	return encodedImage, nil
}


