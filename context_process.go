// Package cairo wraps the c cairographics library.
package cairo

import "github.com/bit101/bitlib/blmath"

// Grayscale turns the image grayscale.
func (c *Context) Grayscale() {
	data, _ := c.Surface.GetData()
	for i := 0; i < len(data); i += 4 {
		val := byte(0.2126*float64(data[i]) + 0.7152*float64(data[i+1]) + 0.0722*float64(data[i+2]))
		data[i] = val
		data[i+1] = val
		data[i+2] = val
	}
	c.Surface.SetData(data)
}

// Threshold turns the image grayscale.
func (c *Context) Threshold(t float64) {
	data, _ := c.Surface.GetData()
	for i := 0; i < len(data); i += 4 {
		for j := 0; j < 3; j++ {
			val := float64(data[i+j]) / 255
			if val < t {
				data[i+j] = 0
			} else {
				data[i+j] = 255
			}
		}
	}
	c.Surface.SetData(data)
}

// Quantize turns the image grayscale.
func (c *Context) Quantize(t int) {
	data, _ := c.Surface.GetData()
	for i := 0; i < len(data); i += 4 {
		for j := 0; j < 3; j++ {
			val := float64(data[i+j])
			val = blmath.Quantize(val, 0, 255, t)
			data[i+j] = byte(val)
		}
	}
	c.Surface.SetData(data)
}
