// Package cairo wraps the c cairographics library.
package cairo

import (
	"errors"
	"fmt"
)

// ImageData holds the pixel data from a surface.
type ImageData struct {
	data   []byte
	Width  int
	Height int
}

// NewImageData creates a new empty ImageData with the given width and height.
func NewImageData(w, h int) ImageData {
	return ImageData{
		data:   make([]byte, w*h*4),
		Width:  w,
		Height: h,
	}
}

// ImageDataFromSurface creates a new ImageData from a surface.
func ImageDataFromSurface(surface *Surface) (ImageData, error) {
	data, err := surface.GetData()
	if err != nil {
		return ImageData{}, fmt.Errorf("unable to create ImageData: %s", err)
	}
	return ImageData{
		data:   data,
		Width:  surface.GetWidth(),
		Height: surface.GetHeight(),
	}, nil
}

// ImageDataFromPNG creates a new ImageData from a png file.
func ImageDataFromPNG(filePath string) (ImageData, error) {
	surface, err := NewSurfaceFromPNG(filePath)
	if err != nil {
		return ImageData{}, fmt.Errorf("unable to create surface for ImageData: %s", err)
	}
	return ImageDataFromSurface(surface)
}

// GetPixel gets the value for a single pixel.
// Values returned are normalized 0-1 for r, g, b, a.
func (b *ImageData) GetPixel(x, y int) (float64, float64, float64, float64) {
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

// GetPixelClamped gets the value for a single pixel.
// Values returned are normalized 0-1 for r, g, b, a.
// If x and y are beyond any edge of the image, it will return the edge pixel value.
func (b *ImageData) GetPixelClamped(x, y, rx, ry, rw, rh int) (float64, float64, float64, float64) {
	if x < rx {
		x = rx
	}
	if x >= rw {
		x = rw - 1
	}
	if y < ry {
		y = ry
	}
	if y >= rh {
		y = rh - 1
	}
	return b.GetPixel(x, y)
}

// GetPixelInt gets the value for a single pixel.
// Values returned are 0-255 for r, g, b, a
func (b *ImageData) GetPixelInt(x, y int) (int, int, int, int) {
	if x < 0 || x >= b.Width || y < 0 || y >= b.Height {
		return 0, 0, 0, 0
	}
	index := (y*b.Width + x) * 4
	return int(b.data[index+2]), int(b.data[index+1]), int(b.data[index]), int(b.data[index+3])
}

// GetPixelIntClamped gets the value for a single pixel.
// Values returned are normalized 0-1 for r, g, b, a.
// If x and y are beyond any edge of the image, it will return the edge pixel value.
func (b *ImageData) GetPixelIntClamped(x, y, w, h int) (int, int, int, int) {
	if x < 0 {
		x = 0
	}
	if x >= w {
		x = w - 1
	}
	if y < 0 {
		y = 0
	}
	if y >= h {
		y = h - 1
	}
	return b.GetPixelInt(x, y)
}

// SetPixel sets the value for a single pixel.
// Params are normalized 0-1 for r, g, b, a.
func (b *ImageData) SetPixel(x, y int, red, green, blue, alpha float64) {
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
func (b *ImageData) SetPixelInt(x, y int, red, green, blue, alpha int) {
	if x < 0 || x >= b.Width || y < 0 || y >= b.Height {
		return
	}
	index := (y*b.Width + x) * 4
	b.data[index] = byte(blue)
	b.data[index+1] = byte(green)
	b.data[index+2] = byte(red)
	b.data[index+3] = byte(alpha)
}

// CopyToSurface copies this ImageData's data into a surface
// The surface must have the same dimensions as this ImageData
// A good use case for this is setting thousands of individual pixels.
// Create an ImageData, set the pixels, copy it back to the surface.
// There is a few ms overhead for getting/setting the data. 2-3ms for a 1600x800 image.
// But setting the pixels is blazingly faster than drawing 1x1 rectangles.
// 100k pixels can be set in less than 2ms with this method,
// as opposed to nearly 30 ms with drawing 1x1 rects.
// But for a small amount of pixels (<10k), it can be faster to draw 1x1 rects.
// Also note that this method will do no antialiasing ever,
// whereas you can get that when drawing rects (for better or worse).
func (b *ImageData) CopyToSurface(surface *Surface) error {
	if surface.GetWidth() != b.Width || surface.GetHeight() != b.Height {
		return errors.New("surface must have same dimensions as the imagedata")
	}
	err := surface.SetData(b.data)
	if err != nil {
		return fmt.Errorf("unable to copy ImageData: %s", err)
	}
	return nil
}
