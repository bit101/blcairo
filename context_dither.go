// Package cairo has dithering algorithms
package cairo

// DitherMethod is a method that determines how to distribute the
// quantization errors across neighboring pixels.
type DitherMethod func(int, int, int, int, int, []int)

// DitherAtkinson dithers an image with the Atkinson algorithm
func (c *Context) DitherAtkinson() {
	c.dither(atkinson)
}

// DitherFloydSteinberg dithers an image with the FloydSteinberg algorithm
func (c *Context) DitherFloydSteinberg() {
	c.dither(floydSteinberg)
}

// DitherJJN dithers an image with the FloydSteinberg algorithm
func (c *Context) DitherJJN() {
	c.dither(jjn)
}

// DitherSierra dithers an image with the FloydSteinberg algorithm
func (c *Context) DitherSierra() {
	c.dither(sierra)
}

// DitherStucki dithers an image with the FloydSteinberg algorithm
func (c *Context) DitherStucki() {
	c.dither(stucki)
}

func (c *Context) dither(ditherMethod DitherMethod) {
	c.Grayscale()
	width := int(c.Width)
	height := int(c.Height)
	data, _ := c.Surface.GetData()

	// copy bytes from data to int array grays
	grays := make([]int, width*height)
	index := 0
	for i := 0; i < len(data); i += 4 {
		grays[index] = int(data[i])
		index++
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// round the value and set pixel to that
			i := y*width + x
			grayOld := grays[i]
			grayNew := 0
			if grayOld > 127 {
				grayNew = 255
			}
			setValue(x, y, width, height, grayNew, grays)

			// distribute the remainder to nearby pixels using given algorithm
			grayErr := (grayOld - grayNew)
			ditherMethod(x, y, width, height, grayErr, grays)
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
	c.Surface.SetData(data)

}

func atkinson(x, y, w, h, grayErr int, grays []int) {
	grayErr /= 8
	addValue(x+1, y, w, h, grayErr, grays)
	addValue(x+2, y, w, h, grayErr, grays)
	addValue(x-1, y+1, w, h, grayErr, grays)
	addValue(x, y+1, w, h, grayErr, grays)
	addValue(x+1, y+1, w, h, grayErr, grays)
	addValue(x, y+2, w, h, grayErr, grays)
}

func floydSteinberg(x, y, w, h, grayErr int, grays []int) {
	grayErr /= 16
	addValue(x+1, y, w, h, grayErr*7, grays)
	addValue(x-1, y+1, w, h, grayErr*3, grays)
	addValue(x, y+1, w, h, grayErr*5, grays)
	addValue(x+1, y+1, w, h, grayErr, grays)
}

func jjn(x, y, w, h, grayErr int, grays []int) {
	grayErr /= 48
	addValue(x+1, y, w, h, grayErr*7, grays)
	addValue(x+2, y, w, h, grayErr*5, grays)

	addValue(x-2, y+1, w, h, grayErr*3, grays)
	addValue(x-1, y+1, w, h, grayErr*5, grays)
	addValue(x, y+1, w, h, grayErr*7, grays)
	addValue(x+1, y+1, w, h, grayErr*5, grays)
	addValue(x+2, y+1, w, h, grayErr*3, grays)

	addValue(x-2, y+2, w, h, grayErr, grays)
	addValue(x-1, y+2, w, h, grayErr*3, grays)
	addValue(x, y+2, w, h, grayErr*5, grays)
	addValue(x+1, y+2, w, h, grayErr*3, grays)
	addValue(x+2, y+2, w, h, grayErr, grays)
}

func sierra(x, y, w, h, grayErr int, grays []int) {
	grayErr /= 32
	addValue(x+1, y, w, h, grayErr*5, grays)
	addValue(x+2, y, w, h, grayErr*3, grays)

	addValue(x-2, y+1, w, h, grayErr*2, grays)
	addValue(x-1, y+1, w, h, grayErr*4, grays)
	addValue(x, y+1, w, h, grayErr*5, grays)
	addValue(x+1, y+1, w, h, grayErr*4, grays)
	addValue(x+2, y+1, w, h, grayErr*2, grays)

	addValue(x-1, y+2, w, h, grayErr*2, grays)
	addValue(x, y+2, w, h, grayErr*3, grays)
	addValue(x+1, y+2, w, h, grayErr*2, grays)
}

func stucki(x, y, w, h, grayErr int, grays []int) {
	grayErr /= 42
	addValue(x+1, y, w, h, grayErr*8, grays)
	addValue(x+2, y, w, h, grayErr*4, grays)

	addValue(x-2, y+1, w, h, grayErr*2, grays)
	addValue(x-1, y+1, w, h, grayErr*4, grays)
	addValue(x, y+1, w, h, grayErr*8, grays)
	addValue(x+1, y+1, w, h, grayErr*4, grays)
	addValue(x+2, y+1, w, h, grayErr*2, grays)

	addValue(x-2, y+2, w, h, grayErr, grays)
	addValue(x-1, y+2, w, h, grayErr*2, grays)
	addValue(x, y+2, w, h, grayErr*4, grays)
	addValue(x+1, y+2, w, h, grayErr*2, grays)
	addValue(x+2, y+2, w, h, grayErr, grays)
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
