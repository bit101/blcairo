// Package cairo wraps the c cairographics library.
package cairo

import (
	"math"

	"github.com/bit101/bitlib/blcolor"
	"github.com/bit101/bitlib/blmath"
	"github.com/bit101/bitlib/random"
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

// GrayscaleRect turns a portion of the image grayscale.
func (c *Context) GrayscaleRect(rx, ry, rw, rh float64) {
	s := NewSurface(rw, rh)
	context := NewContext(s)
	context.SetSourceSurface(c.Surface, -rx, -ry)
	context.Paint()
	context.Grayscale()
	c.SetSourceSurface(s, rx, ry)
	c.Paint()
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

// ThresholdRect performs a threshold operation on a portion of an image.
func (c *Context) ThresholdRect(rx, ry, rw, rh, t, r, g, b, a float64) {
	s := NewSurface(rw, rh)
	context := NewContext(s)
	context.SetSourceSurface(c.Surface, -rx, -ry)
	context.Paint()
	context.Threshold(t, r, g, b, a)
	c.SetSourceSurface(s, rx, ry)
	c.Paint()
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

// ReverseThresholdRect performs a reverse threshold operation on a portion of an images.
func (c *Context) ReverseThresholdRect(rx, ry, rw, rh, t, r, g, b, a float64) {
	s := NewSurface(rw, rh)
	context := NewContext(s)
	context.SetSourceSurface(c.Surface, -rx, -ry)
	context.Paint()
	context.ReverseThreshold(t, r, g, b, a)
	c.SetSourceSurface(s, rx, ry)
	c.Paint()
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

// QuantizeRect quantizes a portion of an image.
func (c *Context) QuantizeRect(rx, ry, rw, rh float64, t int) {
	s := NewSurface(rw, rh)
	context := NewContext(s)
	context.SetSourceSurface(c.Surface, -rx, -ry)
	context.Paint()
	context.Quantize(t)
	c.SetSourceSurface(s, rx, ry)
	c.Paint()
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

// GammaRect gamma corrects a portion of an image.
func (c *Context) GammaRect(rx, ry, rw, rh, gamma float64) {
	s := NewSurface(rw, rh)
	context := NewContext(s)
	context.SetSourceSurface(c.Surface, -rx, -ry)
	context.Paint()
	context.Gamma(gamma)
	c.SetSourceSurface(s, rx, ry)
	c.Paint()
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

// InvertRect inverts the colors in a portion of an image.
func (c *Context) InvertRect(rx, ry, rw, rh float64) {
	s := NewSurface(rw, rh)
	context := NewContext(s)
	context.SetSourceSurface(c.Surface, -rx, -ry)
	context.Paint()
	context.Invert()
	c.SetSourceSurface(s, rx, ry)
	c.Paint()
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

// ContrastRect adjusts the contrast in a portion of an image.
func (c *Context) ContrastRect(rx, ry, rw, rh, amt float64) {
	s := NewSurface(rw, rh)
	context := NewContext(s)
	context.SetSourceSurface(c.Surface, -rx, -ry)
	context.Paint()
	context.Contrast(amt)
	c.SetSourceSurface(s, rx, ry)
	c.Paint()
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

// BrightnessRect adjusts the brightness of a portion of an image.
func (c *Context) BrightnessRect(rx, ry, rw, rh, amt float64) {
	s := NewSurface(rw, rh)
	context := NewContext(s)
	context.SetSourceSurface(c.Surface, -rx, -ry)
	context.Paint()
	context.Brightness(amt)
	c.SetSourceSurface(s, rx, ry)
	c.Paint()
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

// TintRect tints a portion of an image.
func (c *Context) TintRect(rx, ry, rw, rh, r, g, b, t float64) {
	s := NewSurface(rw, rh)
	context := NewContext(s)
	context.SetSourceSurface(c.Surface, -rx, -ry)
	context.Paint()
	context.Tint(r, g, b, t)
	c.SetSourceSurface(s, rx, ry)
	c.Paint()
}

// Hue tints an image to a given hue
func (c *Context) Hue(hue, t float64) {
	color := blcolor.HSV(hue, 1, 1)
	c.Tint(color.R, color.G, color.B, t)
}

// HueRect tints a portion of an image to a given hue.
func (c *Context) HueRect(rx, ry, rw, rh, hue, t float64) {
	s := NewSurface(rw, rh)
	context := NewContext(s)
	context.SetSourceSurface(c.Surface, -rx, -ry)
	context.Paint()
	context.Hue(hue, t)
	c.SetSourceSurface(s, rx, ry)
	c.Paint()
}

// Blur executes a box blur.
func (c *Context) Blur(radius int) {
	if radius < 1 {
		// blur of 0 does nothing.
		// blur of less than 0 is wrong. ignore.
		return
	}

	srcIm, _ := ImageDataFromSurface(c.Surface)
	dstIm := NewImageData(srcIm.Width, srcIm.Height)
	w := int(c.Width)
	h := int(c.Height)
	// doing a two-pass (h+v) blur is O(m^2*2n) m=bitmap size, n = kernel size
	// as opposed to O(m^2*n^2) for a regular box blur
	// horizontal blur
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, t := 0, 0, 0, 0
			for j := -radius; j <= radius; j++ {
				rr, gg, bb, _ := srcIm.GetPixelIntClamped(x+j, y, w, h)
				r += rr
				g += gg
				b += bb
				t++
			}
			dstIm.SetPixelInt(x, y, r/t, g/t, b/t, 255)
		}
	}

	// vertical blur
	// switch src and dst byte textures here
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, t := 0, 0, 0, 0
			for j := -radius; j <= radius; j++ {
				rr, gg, bb, _ := dstIm.GetPixelIntClamped(x, y+j, w, h)
				r += rr
				g += gg
				b += bb
				t++
			}
			srcIm.SetPixelInt(x, y, r/t, g/t, b/t, 255)
		}
	}
	// the final pass put the bytes in src, so we copy that back.
	srcIm.CopyToSurface(c.Surface)
}

// BlurRect executes a box blur on a portion of an image.
func (c *Context) BlurRect(rx, ry, rw, rh float64, radius int) {
	s := NewSurface(rw, rh)
	context := NewContext(s)
	context.SetSourceSurface(c.Surface, -rx, -ry)
	context.Paint()
	context.Blur(radius)
	c.SetSourceSurface(s, rx, ry)
	c.Paint()
}

// GaussianBlur executes a Gaussian blur.
func (c *Context) GaussianBlur(radius int) {
	if radius < 1 {
		// blur of 0 does nothing.
		// less than 0 is just wrong. we'll ignore.
		return
	}
	kernel := getGaussKernel(radius*2 + 1)

	srcIm, _ := ImageDataFromSurface(c.Surface)
	dstIm := NewImageData(srcIm.Width, srcIm.Height)
	w := int(c.Width)
	h := int(c.Height)

	// horizontal blur
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, a := 0.0, 0.0, 0.0, 0.0
			for j := -radius; j <= radius; j++ {
				rr, gg, bb, aa := srcIm.GetPixelClamped(x+j, y, 0, 0, w, h)
				k := kernel[j+radius]
				r += rr * k
				g += gg * k
				b += bb * k
				a += aa * k
			}
			dstIm.SetPixel(x, y, r, g, b, a)
		}
	}

	// vertical blur
	// switch src and dst byte textures here
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, a := 0.0, 0.0, 0.0, 0.0
			for j := -radius; j <= radius; j++ {
				rr, gg, bb, aa := dstIm.GetPixelClamped(x, y+j, 0, 0, w, h)
				k := kernel[j+radius]
				r += rr * k
				g += gg * k
				b += bb * k
				a += aa * k
			}
			srcIm.SetPixel(x, y, r, g, b, a)
		}
	}
	// the final pass put the bytes in src, so we copy that back.
	srcIm.CopyToSurface(c.Surface)
}

// GaussianBlurRect executes a Gaussian blur on a portion of an image.
func (c *Context) GaussianBlurRect(rx, ry, rw, rh float64, radius int) {
	s := NewSurface(rw, rh)
	context := NewContext(s)
	context.SetSourceSurface(c.Surface, -rx, -ry)
	context.Paint()
	context.GaussianBlur(radius)
	c.SetSourceSurface(s, rx, ry)
	c.Paint()
}

func getGaussKernel(size int) []float64 {
	sigma := float64(size-1) / 5
	mean := size / 2
	kernel := make([]float64, size)
	sum := 0.0
	for x := 0; x < size; x++ {
		kernelValue := math.Exp(-0.5 * (math.Pow((float64(x-mean))/sigma, 2)))
		kernel[x] = kernelValue
		sum += kernelValue
	}
	for x := 0; x < size; x++ {
		kernel[x] /= sum
	}
	return kernel
}

// Pixelate pixelates an image.
func (c *Context) Pixelate(size int) {
	srcIm, _ := ImageDataFromSurface(c.Surface)
	w := int(c.Width)
	h := int(c.Height)

	for x := 0; x < w; x += size {
		for y := 0; y < h; y += size {
			rr, gg, bb, aa := srcIm.GetPixelClamped(x, y, 0, 0, w, h)
			c.SetSourceRGBA(rr, gg, bb, aa)
			c.FillRectangle(float64(x), float64(y), float64(size), float64(size))
		}
	}
}

// PixelateRect pixelates a portion of an image.
func (c *Context) PixelateRect(rx, ry, rw, rh float64, size int) {
	s := NewSurface(rw, rh)
	context := NewContext(s)
	context.SetSourceSurface(c.Surface, -rx, -ry)
	context.Paint()
	context.Pixelate(size)
	c.SetSourceSurface(s, rx, ry)
	c.Paint()
}

// Sharpen executes a sharpen filter.
func (c *Context) Sharpen() {
	srcIm, _ := ImageDataFromSurface(c.Surface)
	dstIm := NewImageData(srcIm.Width, srcIm.Height)
	w := int(c.Width)
	h := int(c.Height)
	kernel := [][]float64{{0, -1, 0}, {-1, 5, -1}, {0, -1, 0}}
	// kernel := [][]float64{{0, 0, 0}, {0, 2, 0}, {0, 0, 0}}

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b := 0.0, 0.0, 0.0
			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {
					rr, gg, bb, _ := srcIm.GetPixelClamped(x+i, y+j, 0, 0, w, h)
					r += rr * kernel[i+1][j+1]
					g += gg * kernel[i+1][j+1]
					b += bb * kernel[i+1][j+1]
				}
			}
			r = blmath.Clamp(r, 0, 1)
			g = blmath.Clamp(g, 0, 1)
			b = blmath.Clamp(b, 0, 1)
			dstIm.SetPixel(x, y, r, g, b, 1)
		}
	}
	dstIm.CopyToSurface(c.Surface)
}

// SharpenRect sharpens a portion of an image.
func (c *Context) SharpenRect(rx, ry, rw, rh float64) {
	s := NewSurface(rw, rh)
	context := NewContext(s)
	context.SetSourceSurface(c.Surface, -rx, -ry)
	context.Paint()
	context.Sharpen()
	c.SetSourceSurface(s, rx, ry)
	c.Paint()
}

// MapGradient maps the brightness values in an image to a gradient between two colors.
func (c *Context) MapGradient(col0, col1 blcolor.Color) {
	w := int(c.Width)
	h := int(c.Height)
	c.Grayscale()
	srcIm, _ := ImageDataFromSurface(c.Surface)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, _, _, a := srcIm.GetPixel(x, y)
			c := blcolor.Lerp(col0, col1, r)
			srcIm.SetPixel(x, y, c.R, c.G, c.B, a)
		}
	}
	srcIm.CopyToSurface(c.Surface)
}

// MapGradientRect performs a map gradient operation on a portion of an image.
func (c *Context) MapGradientRect(rx, ry, rw, rh float64, col0, col1 blcolor.Color) {
	s := NewSurface(rw, rh)
	context := NewContext(s)
	context.SetSourceSurface(c.Surface, -rx, -ry)
	context.Paint()
	context.MapGradient(col0, col1)
	c.SetSourceSurface(s, rx, ry)
	c.Paint()
}

// MapHue maps the brightness values in an image to a gradient between two hues.
func (c *Context) MapHue(hue0, hue1 float64) {
	w := int(c.Width)
	h := int(c.Height)
	c.Grayscale()
	srcIm, _ := ImageDataFromSurface(c.Surface)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, _, _, a := srcIm.GetPixel(x, y)
			c := blcolor.HSV(blmath.Lerp(r, hue0, hue1), 1, 1)
			srcIm.SetPixel(x, y, c.R, c.G, c.B, a)
		}
	}
	srcIm.CopyToSurface(c.Surface)
}

// MapHueRect performs a map hue operation on a portion of an image.
func (c *Context) MapHueRect(rx, ry, rw, rh, hue0, hue1 float64) {
	s := NewSurface(rw, rh)
	context := NewContext(s)
	context.SetSourceSurface(c.Surface, -rx, -ry)
	context.Paint()
	context.MapHue(hue0, hue1)
	c.SetSourceSurface(s, rx, ry)
	c.Paint()
}

// Noisify adds noise to an image.
func (c *Context) Noisify(amount float64) {
	w := int(c.Width)
	h := int(c.Height)
	srcIm, _ := ImageDataFromSurface(c.Surface)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, a := srcIm.GetPixel(x, y)
			r += random.FloatRange(-amount, amount)
			g += random.FloatRange(-amount, amount)
			g += random.FloatRange(-amount, amount)
			srcIm.SetPixel(x, y, blmath.Clamp(r, 0, 1), blmath.Clamp(g, 0, 1), blmath.Clamp(b, 0, 1), a)
		}
	}
	srcIm.CopyToSurface(c.Surface)
}

// NoisifyRect applies noise to a portion of an image.
func (c *Context) NoisifyRect(rx, ry, rw, rh, amount float64) {
	s := NewSurface(rw, rh)
	context := NewContext(s)
	context.SetSourceSurface(c.Surface, -rx, -ry)
	context.Paint()
	context.Noisify(amount)
	c.SetSourceSurface(s, rx, ry)
	c.Paint()
}

// DrawContext sets the source of this context to the surface of another context and paints.
func (c *Context) DrawContext(src *Context, x, y float64) {
	c.SetSourceSurface(src.Surface, x, y)
	c.Paint()
}

// FilterChannels selectively filters out each rgb channel by scaling it from 0 to 1.
// Example: to get only the red channel, context.FilterChannels(1, 0, 0)
func (s *Surface) FilterChannels(r, g, b float64) {
	data, _ := s.GetData()
	for i := 0; i < len(data); i += 4 {
		data[i] = byte(float64(data[i]) * b)
		data[i+1] = byte(float64(data[i+1]) * g)
		data[i+2] = byte(float64(data[i+2]) * r)
		data[i+3] = 128
	}
	s.SetData(data)
}

// ColorFringe applies a chromatic abberation effect, seperating the rgb color channels horizontally by a given amount.
func (c *Context) ColorFringe(offset float64) {
	c.Save()

	// get image data
	data, _ := c.Surface.GetData()
	s := NewSurface(c.Width, c.Height)

	// clear image to black and set screen operator
	c.ClearBlack()
	c.SetOperator(OperatorScreen)

	// draw the red channel shifted right
	s.SetData(data)
	s.FilterChannels(1, 0, 0)
	c.SetSourceSurface(s, offset, 0)
	c.Paint()

	// draw the green channel, no shift
	s.SetData(data)
	s.FilterChannels(0, 1, 0)
	c.SetSourceSurface(s, 0, 0)
	c.Paint()

	// draw the blue channel shifted left
	s.SetData(data)
	s.FilterChannels(0, 0, 1)
	c.SetSourceSurface(s, -offset, 0)
	c.Paint()

	c.Restore()
}

// ColorFringeRect performs a color fringe operation on a portion of an image.
func (c *Context) ColorFringeRect(rx, ry, rw, rh, offset float64) {
	s := NewSurface(rw, rh)
	context := NewContext(s)
	context.SetSourceSurface(c.Surface, -rx, -ry)
	context.Paint()
	context.ColorFringe(offset)
	c.SetSourceSurface(s, rx, ry)
	c.Paint()
}
