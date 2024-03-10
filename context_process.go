// Package cairo wraps the c cairographics library.
package cairo

import (
	"math"

	"github.com/bit101/bitlib/blcolor"
	"github.com/bit101/bitlib/blmath"
)

// Grayscale turns the image grayscale.
func (c *Context) Grayscale() {
	data, _ := c.Surface.GetData()
	r := 0.299
	g := 0.587
	b := 0.113
	for i := 0; i < len(data); i += 4 {
		// note channel order: bgr
		val := byte(b*float64(data[i]) + g*float64(data[i+1]) + r*float64(data[i+2]))
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

// Gamma does gamma correction on an image.
// gamma values less than 1.0 darken the image, greater than 1.0 lighten it.
func (c *Context) Gamma(gamma float64) {
	gammaCorrection := 1.0 / gamma
	data, _ := c.Surface.GetData()
	for i := 0; i < len(data); i += 4 {
		for j := 0; j < 3; j++ {
			val := float64(data[i+j]) / 255.0
			val = math.Pow(val, gammaCorrection)
			data[i+j] = byte(val * 255)
		}
	}
	c.Surface.SetData(data)

}

// Invert inverts the colors of an image.
func (c *Context) Invert() {
	data, _ := c.Surface.GetData()
	for i := 0; i < len(data); i += 4 {
		for j := 0; j < 3; j++ {
			data[i+j] = 255 - data[i+j]
		}
	}
	c.Surface.SetData(data)
}

// Contrast inverts the colors of an image.
func (c *Context) Contrast(amt float64) {
	cont := 255.0 * amt
	f := (259.0 * (cont + 255.0)) / (255 * (259.0 - cont))
	data, _ := c.Surface.GetData()
	for i := 0; i < len(data); i += 4 {
		for j := 0; j < 3; j++ {
			val := float64(data[i+j])
			val = f*(val-128) + 128
			val = blmath.Clamp(val, 0, 255)
			data[i+j] = byte(val)
		}
	}
	c.Surface.SetData(data)
}

// Brightness adjusts the brightness of an image.
func (c *Context) Brightness(amt float64) {
	brightness := 255.0 * amt
	data, _ := c.Surface.GetData()
	for i := 0; i < len(data); i += 4 {
		for j := 0; j < 3; j++ {
			val := float64(data[i+j]) + brightness
			val = blmath.Clamp(val, 0, 255)
			data[i+j] = byte(val)
		}
	}
	c.Surface.SetData(data)
}

// Tint tints an image
func (c *Context) Tint(r, g, b, t float64) {
	data, _ := c.Surface.GetData()
	for i := 0; i < len(data); i += 4 {
		val := float64(data[i])
		val = blmath.Lerp(t, val, b*255)
		val = blmath.Clamp(val, 0, 255)
		data[i] = byte(val)

		val = float64(data[i+1])
		val = blmath.Lerp(t, val, g*255)
		val = blmath.Clamp(val, 0, 255)
		data[i+1] = byte(val)

		val = float64(data[i+2])
		val = blmath.Lerp(t, val, r*255)
		val = blmath.Clamp(val, 0, 255)
		data[i+2] = byte(val)
	}
	c.Surface.SetData(data)
}

// Hue tints an image to a given hue
func (c *Context) Hue(hue, t float64) {
	color := blcolor.HSV(hue, 1, 1)
	c.Tint(color.R, color.G, color.B, t)
}
