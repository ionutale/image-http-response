package main

import (
	"fmt"
	"os"
	"github.com/daddye/vips"
)

func ResizeImage(image []byte, width int, height int, quality int, format string) ([]byte, error) {
	
	if format[0] == '.' {
		format = format[1:]
	}
	
	options := vips.Options{
		Width:        width,
		Height:       height,
		Crop:         true,
		Extend:       vips.EXTEND_WHITE,
		Interpolator: vips.BILINEAR,
		Gravity:      vips.CENTRE,
		Quality:      quality,
	}

	buf, err := vips.Resize(image, options)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	return buf, nil
}
