// Package cairo wraps the c cairographics library.
package cairo

import (
	"log"

	"github.com/bit101/bitlib/blcolor"
	"github.com/bit101/bitlib/random"
)

// SampleColors returns a pallet of colors from an external image.
func SampleColors(image string, count int) []blcolor.Color {
	surface, err := NewSurfaceFromPNG(image)
	if err != nil {
		log.Fatal(err)
	}
	data, err := surface.GetData()
	if err != nil {
		log.Fatal(err)
	}
	colors := []blcolor.Color{}
	for len(colors) < count {
		x := random.IntRange(0, surface.GetWidth())
		y := random.IntRange(0, surface.GetHeight())
		r, g, b, _ := surface.GetPixel(data, x, y)
		red := float64(r) / 255
		green := float64(g) / 255
		blue := float64(b) / 255
		c := blcolor.RGB(red, green, blue)
		colors = append(colors, c)
	}
	return colors
}
