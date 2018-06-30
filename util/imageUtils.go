package util

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/ftrvxmtrx/tga"
)

// InitImage creates a blank image of size x,y and color color
func InitImage(x, y int, color color.Color) *image.NRGBA {
	res := image.NewNRGBA(image.Rect(0, 0, x, y))
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			res.Set(i, j, color)
		}
	}
	return res
}

// LoadTexture takes a .tga image and loads it
func LoadTexture(filename string) (*image.Image, error) {
	file, err := os.Open(filename)
	texture, err := tga.Decode(file)
	if err != nil {
		return nil, err
	}
	return &texture, nil
}

// DrawFile writes image to disk as output.png
func DrawFile(i *image.NRGBA, filename string) {
	fo, err := os.Create(filename)
	// close file on exit
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	if err == nil {
		png.Encode(fo, i)
	}
}
