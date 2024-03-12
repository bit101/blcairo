// Package cairo wraps the c cairographics library.
package cairo

import (
	"log"
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

// Threshold sets any pixel whose average value is below t to the given rgba value.
func (c *Context) Threshold(t, r, g, b, a float64) {
	data, _ := c.Surface.GetData()
	for i := 0; i < len(data); i += 4 {
		val := float64(data[i]) / 255
		val += float64(data[i+1]) / 255
		val += float64(data[i+2]) / 255
		val /= 3
		if val < t {
			data[i] = byte(b * 255)
			data[i+1] = byte(g * 255)
			data[i+2] = byte(r * 255)
			data[i+3] = byte(a * 255)
		}
	}
	c.Surface.SetData(data)
}

// ReverseThreshold sets any pixel whose average value is greater than t to the given rgba value.
func (c *Context) ReverseThreshold(t, r, g, b, a float64) {
	data, _ := c.Surface.GetData()
	for i := 0; i < len(data); i += 4 {
		val := float64(data[i]) / 255
		val += float64(data[i+1]) / 255
		val += float64(data[i+2]) / 255
		val /= 3
		if val > t {
			data[i] = byte(b * 255)
			data[i+1] = byte(g * 255)
			data[i+2] = byte(r * 255)
			data[i+3] = byte(a * 255)
		}
	}
	c.Surface.SetData(data)
}

// Quantize reduces the number of colors in an image.
// Technically, it quantizes the values of each pixel separately,
// so the result will have more than t colors.
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

// Contrast changes the balance of dark and light areas in an image.
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

// Tint tints an image.
// r, g, b, determine the color of the tint.
// t determines how much the tint is applied. t=1 will result in the entire image being a single color.
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

// Blur executes a box blur.
func (c *Context) Blur(ksize int) {
	if ksize < 0 {
		log.Fatal("blur must be greater than zero")
	}
	if ksize < 2 {
		// blur of 0 or 1 does nothing.
		return
	}

	srcBt, _ := ByteTextureFromSurface(c.Surface)
	dstBt := NewByteTexture(srcBt.Width, srcBt.Height)
	w := int(c.Width)
	h := int(c.Height)

	// doing a two-pass (h+v) blur is O(m^2*2n) m=bitmap size, n = kernel size
	// as opposed to O(m^2*n^2) for a regular box blur
	// horizontal blur
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, t := 0, 0, 0, 0
			for j := -ksize / 2; j <= ksize/2; j++ {
				rr, gg, bb, _ := srcBt.GetPixelInt(x+j, y)
				r += rr
				g += gg
				b += bb
				t++
			}
			dstBt.SetPixelInt(x, y, r/t, g/t, b/t, 255)
		}
	}

	// vertical blur
	// switch src and dst byte textures here
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, t := 0, 0, 0, 0
			for j := -ksize / 2; j <= ksize/2; j++ {
				rr, gg, bb, _ := dstBt.GetPixelInt(x, y+j)
				r += rr
				g += gg
				b += bb
				t++
			}
			srcBt.SetPixelInt(x, y, r/t, g/t, b/t, 255)
		}
	}
	// the final pass put the bytes in src, so we copy that back.
	srcBt.CopyToSurface(c.Surface)
}
