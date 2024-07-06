// Package cairo wraps the c cairographics library.
package cairo

import (
	"log"

	"github.com/bit101/bitlib/blcolor"
	"github.com/bit101/bitlib/random"
)

////////////////////
// CLEAR and SET
////////////////////

// BlackOnWhite clears the image to white and sets the drawing color to black.
func (c *Context) BlackOnWhite() {
	c.ClearWhite()
	c.SetSourceBlack()
}

// WhiteOnBlack clears the image to black and sets the drawing color to white.
func (c *Context) WhiteOnBlack() {
	c.ClearBlack()
	c.SetSourceWhite()
}

// Blueprint clears the image to a blueprint shade and sets the drawing color to white.
func (c *Context) Blueprint() {
	c.ClearRGB(0.0, 0.0784, 0.5176)
	c.SetSourceWhite()
}

////////////////////
// CLEAR
////////////////////

// ClearClear clears the context to full transparency.
func (c *Context) ClearClear() {
	c.Save()
	// todo: set identity transform
	c.SetOperator(OperatorClear)
	c.Paint()
	c.Restore()
}

// ClearRGB clears the context to the given rgb color.
func (c *Context) ClearRGB(r float64, g float64, b float64) {
	c.Save()
	// todo: set identity transform
	c.SetSourceRGB(r, g, b)
	c.Paint()
	c.Restore()
}

// ClearRGBA clears the context to the given rgba color.
func (c *Context) ClearRGBA(r, g, b, a float64) {
	c.Save()
	// todo: set identity transform
	c.SetSourceRGBA(r, g, b, a)
	c.Paint()
	c.Restore()
}

// ClearColor clears context to given color.
func (c *Context) ClearColor(color blcolor.Color) {
	c.ClearRGB(color.R, color.G, color.B)
}

// ClearWhite clears context to white.
func (c *Context) ClearWhite() {
	c.ClearRGB(1, 1, 1)
}

// ClearBlack clears context to white.
func (c *Context) ClearBlack() {
	c.ClearRGB(0, 0, 0)
}

// ClearGray clears context to white.
func (c *Context) ClearGray(g float64) {
	c.ClearRGB(g, g, g)
}

// ClearHSV clears the image to the given hsv value.
func (c *Context) ClearHSV(h, s, v float64) {
	c.ClearColor(blcolor.HSV(h, s, v))
}

// ClearHSVA clears the image to the given hsva value.
func (c *Context) ClearHSVA(h, s, v, a float64) {
	c.ClearColor(blcolor.HSVA(h, s, v, a))
}

// ClearRandomGray clears the image to a random shade of gray.
func (c *Context) ClearRandomGray() {
	c.ClearGray(random.Float())
}

// ClearRandomRGB clears the image to a random rgb value.
func (c *Context) ClearRandomRGB() {
	c.ClearRGB(random.Float(), random.Float(), random.Float())
}

////////////////////
// SETSOURCE
////////////////////

// SetSourceColor sets the source to the given color.
func (c *Context) SetSourceColor(color blcolor.Color) {
	c.SetSourceRGBA(color.R, color.G, color.B, color.A)
}

// SetSourceBlack sets the source to black.
func (c *Context) SetSourceBlack() {
	c.SetSourceRGB(0, 0, 0)
}

// SetSourceWhite sets the source to white.
func (c *Context) SetSourceWhite() {
	c.SetSourceRGB(1, 1, 1)
}

// SetSourceGray sets the source to the specified gray shade.
func (c *Context) SetSourceGray(gray float64) {
	c.SetSourceRGB(gray, gray, gray)
}

// SetSourceHSV sets the drawing color to the given hsv value.
func (c *Context) SetSourceHSV(h, s, v float64) {
	c.SetSourceColor(blcolor.HSV(h, s, v))
}

// SetSourceHSVA sets the drawing color to the given hsva value.
func (c *Context) SetSourceHSVA(h, s, v, a float64) {
	c.SetSourceColor(blcolor.HSVA(h, s, v, a))
}

// SetSourceRandomGray sets the drawing color to a random gray shade.
func (c *Context) SetSourceRandomGray() {
	c.SetSourceGray(random.Float())
}

// SetSourceRandomRGB sets the drawing color to a random rgb value.
func (c *Context) SetSourceRandomRGB() {
	c.SetSourceRGB(random.Float(), random.Float(), random.Float())
}

////////////////////
// MISC
////////////////////

// TranslateCenter translates the context to its center
func (c *Context) TranslateCenter() {
	c.Translate(c.GetCenter())
}

// GetCenter returns the x, y coords of the center of the context
func (c *Context) GetCenter() (float64, float64) {
	return c.Width / 2, c.Height / 2
}

// GetAspectRatio returns the aspect ratio of the context (width / height).
func (c *Context) GetAspectRatio() float64 {
	return c.Width / c.Height
}

// Size returns the width and height of this context.
func (c *Context) Size() (float64, float64) {
	return c.Width, c.Height
}

// ProcessPixels runs a function for every pixel in the context.
func (c *Context) ProcessPixels(pixelFunc func(context *Context, x, y float64)) {
	w, h := c.Size()
	for x := 0.0; x < w; x++ {
		for y := 0.0; y < h; y++ {
			pixelFunc(c, x, y)
		}
	}
}

// PaintImage loads an image from an external png file and paints the context with that image.
func (c *Context) PaintImage(imagePath string, x, y float64) {
	surface, err := NewSurfaceFromPNG(imagePath)
	if err != nil {
		log.Fatalf("could not load image: %s", err)
	}
	c.DrawSurface(surface, x, y)
}

// PaintImageCentered loads an image from an external png file and paints the context with that image.
func (c *Context) PaintImageCentered(imagePath string, x, y float64) {
	surface, err := NewSurfaceFromPNG(imagePath)
	if err != nil {
		log.Fatalf("could not load image: %s", err)
	}
	c.Save()
	c.SetSourceSurface(surface, x-surface.GetWidthF()/2, y-surface.GetHeightF()/2)
	c.Paint()
	// not sure why this next line is here, but I think it fixed something.
	// may need to add to other methods?
	c.SetSourceBlack()
	c.Restore()
}

// DrawSurface draws the given surface onto this context with default composite operation.
func (c *Context) DrawSurface(surface *Surface, x, y float64) {
	c.Save()
	c.SetSourceSurface(surface, x, y)
	c.Paint()
	c.Restore()
}

// DrawSurfaceUnder draws the given surface onto this context with the DestOver operation.
// This draws the surface UNDER any existing content on the surface.
// If the context is fully opaque, this will have no effect. The new surface only shows through
// in areas where the context has transparency.
func (c *Context) DrawSurfaceUnder(surface *Surface, x, y float64) {
	c.Save()
	c.SetOperator(OperatorDestOver)
	c.SetSourceSurface(surface, x, y)
	c.Paint()
	c.Restore()
}
