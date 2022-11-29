// Package cairo wraps the c cairographics library.
package cairo

import (
	"github.com/bit101/bitlib/blcolor"
)

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

// ClearGrey clears context to white.
func (c *Context) ClearGrey(g float64) {
	c.ClearRGB(g, g, g)
}

// SetSourceColor sets the source to the given color.
func (c *Context) SetSourceColor(color blcolor.Color) {
	c.SetSourceRGBA(color.R, color.G, color.B, color.A)
}

// SetSourceHSV sets the source to a color created with the given hue, saturation and value.
func (c *Context) SetSourceHSV(hue, sat, val float64) {
	c.SetSourceColor(blcolor.HSV(hue, sat, val))
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

// Center translates the context to its center
func (c *Context) Center() {
	c.Translate(c.Width/2, c.Height/2)
}

// GetAspectRatio returns the aspect ratio of the context (width / height).
func (c *Context) GetAspectRatio() float64 {
	return c.Width / c.Height
}

// Size returns the width and height of this context.
func (c *Context) Size() (float64, float64) {
	return c.Width, c.Height
}
