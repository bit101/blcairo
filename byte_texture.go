// Package cairo wraps the c cairographics library.
package cairo

import (
	"errors"
	"fmt"
)

// ByteTexture holds the pixel data from a surface.
type ByteTexture struct {
	data   []byte
	Width  int
	Height int
}

// NewByteTexture creates a new empty ByteTexture with the given width and height.
func NewByteTexture(w, h int) ByteTexture {
	return ByteTexture{
		data:   make([]byte, w*h*4),
		Width:  w,
		Height: h,
	}
}

// ByteTextureFromSurface creates a new ByteTexture from a surface.
func ByteTextureFromSurface(surface *Surface) (ByteTexture, error) {
	data, err := surface.GetData()
	if err != nil {
		return ByteTexture{}, fmt.Errorf("unable to create ByteTexture: %s", err)
	}
	return ByteTexture{
		data:   data,
		Width:  surface.GetWidth(),
		Height: surface.GetHeight(),
	}, nil
}

// ByteTextureFromPNG creates a new ByteTexture from a png file.
func ByteTextureFromPNG(filePath string) (ByteTexture, error) {
	surface, err := NewSurfaceFromPNG(filePath)
	if err != nil {
		return ByteTexture{}, fmt.Errorf("unable to create surface for ByteTexture: %s", err)
	}
	return ByteTextureFromSurface(surface)
}

// GetPixel gets the value for a single pixel.
// Values returned are normalized 0-1 for r, g, b, a.
func (b *ByteTexture) GetPixel(x, y int) (float64, float64, float64, float64) {
	if x < 0 || x >= b.Width || y < 0 || y >= b.Height {
		return 0, 0, 0, 1
	}
	index := (y*b.Width + x) * 4
	blue := float64(b.data[index])
	green := float64(b.data[index+1])
	red := float64(b.data[index+2])
	alpha := float64(b.data[index+3])
	return red / 255.0, green / 255.0, blue / 255.0, alpha / 255.0
}

// GetPixelInt gets the value for a single pixel.
// Values returned are 0-255 for r, g, b, a
func (b *ByteTexture) GetPixelInt(x, y int) (int, int, int, int) {
	if x < 0 || x >= b.Width || y < 0 || y >= b.Height {
		return 0, 0, 0, 1
	}
	index := (y*b.Width + x) * 4
	return int(b.data[index+2]), int(b.data[index+1]), int(b.data[index]), int(b.data[index+3])
}

// SetPixel sets the value for a single pixel.
// Params are normalized 0-1 for r, g, b, a.
func (b *ByteTexture) SetPixel(x, y int, red, green, blue, alpha float64) {
	if x < 0 || x >= b.Width || y < 0 || y >= b.Height {
		return
	}
	index := (y*b.Width + x) * 4
	b.data[index] = byte(blue * 255.0)
	b.data[index+1] = byte(green * 255.0)
	b.data[index+2] = byte(red * 255.0)
	b.data[index+3] = byte(alpha * 255.0)
}

// SetPixelInt sets the value for a single pixel.
// Params are 0-255 for r, g, b, a
func (b *ByteTexture) SetPixelInt(x, y int, red, green, blue, alpha int) {
	if x < 0 || x >= b.Width || y < 0 || y >= b.Height {
		return
	}
	index := (y*b.Width + x) * 4
	b.data[index] = byte(blue)
	b.data[index+1] = byte(green)
	b.data[index+2] = byte(red)
	b.data[index+3] = byte(alpha)
}

// CopyToSurface copies this ByteTexture's data into a surface
// The surface must have the same dimensions as this ByteData.
// A good use case for this is setting thousands of individual pixels.
// Create a ByteTexture, set the pixels, copy it back to the surface.
// There is a few ms overhead for getting/setting the data. 2-3ms for a 1600x800 image.
// But setting the pixels is blazingly faster than drawing 1x1 rectangles.
// 100k pixels can be set in less than 2ms with this method,
// as opposed to nearly 30 ms with drawing 1x1 rects.
// But for a small amount of pixels (<10k), it can be faster to draw 1x1 rects.
// Also note that this method will do no antialiasing ever,
// whereas you can get that when drawing rects (for better or worse).
func (b *ByteTexture) CopyToSurface(surface *Surface) error {
	if surface.GetWidth() != b.Width || surface.GetHeight() != b.Height {
		return errors.New("surface must have same dimensions as the bytedata")
	}
	err := surface.SetData(b.data)
	if err != nil {
		return fmt.Errorf("unable to copy ByteTexture: %s", err)
	}
	return nil
}
