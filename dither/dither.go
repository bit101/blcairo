// Package dither has dithering algorithms
package dither

import (
	cairo "github.com/bit101/blcairo"
)

// Atkinson dithers an image with the Atkinson algorithm
func Atkinson(context *cairo.Context) {
	width := int(context.Width)
	height := int(context.Height)
	data, _ := context.Surface.GetData()

	// copy bytes from data to int array grays
	grays := make([]int, width*height)
	index := 0
	for i := 0; i < len(data); i += 4 {
		grays[index] = int(data[i])
		index++
	}

	// the atkinson algorithm
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// round the value and set pixel to that
			i := y*width + x
			grayOld := grays[i]
			grayNew := 0
			if grayOld > 127 {
				grayNew = 255
			}
			grayErr := (grayOld - grayNew) / 8
			setValue(x, y, width, height, grayNew, grays)

			// distribute the remainder to nearby pixels.
			addValue(x+1, y, width, height, grayErr, grays)
			addValue(x+2, y, width, height, grayErr, grays)
			addValue(x-1, y+1, width, height, grayErr, grays)
			addValue(x, y+1, width, height, grayErr, grays)
			addValue(x+1, y+1, width, height, grayErr, grays)
			addValue(x, y+2, width, height, grayErr, grays)
		}
	}

	// copy int data back to byte array
	index = 0
	for i := 0; i < len(grays); i++ {
		data[index] = byte(grays[i])
		data[index+1] = byte(grays[i])
		data[index+2] = byte(grays[i])
		index += 4
	}
	context.Surface.SetData(data)

}

func setValue(x, y, width, height, value int, data []int) {
	if x >= width || x < 0 || y >= height || y < 0 {
		return
	}
	index := y*width + x
	data[index] = value
}

func addValue(x, y, width, height, value int, data []int) {
	if x >= width || x < 0 || y >= height || y < 0 {
		return
	}
	index := y*width + x
	data[index] += value
}
