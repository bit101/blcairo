// Package cairo wraps the c cairographics library.
package cairo

import (
	"log"
	"math"

	"github.com/bit101/bitlib/blmath"
)

// SphereMap maps a ByteTexture onto a sphere.
// bt is the ByteTexture instance.
// x, y is the top left of the sphere bounding box, diameter is the diameter of the sphere.
// rotation is normalized. -1.0 is a full rotation to the right, 1.0 is a full rotation to the left.
func (c *Context) SphereMap(bt ByteTexture, x, y, diameter, rotation float64) {
	buffer, err := ByteTextureFromSurface(c.Surface)
	if err != nil {
		log.Fatalf("unable to create buffer: %s", err)
	}
	mapHeightF := float64(bt.Height)
	mapWidthF := float64(bt.Width)
	for xx := x; xx < diameter+x; xx++ {
		for yy := 0.0; yy < diameter+y; yy++ {
			px := blmath.Map(xx-x, 0, diameter, -1, 1)
			py := blmath.Map(yy-y, 0, diameter, -1, 1)
			magSq := px*px + py*py
			if magSq <= 1 {

				widthAtHeight := math.Sqrt(1 - py*py)
				if widthAtHeight != 0 {
					px = math.Asin(px/widthAtHeight) * 2 / math.Pi
				}
				py = math.Asin(py) * 2 / math.Pi
				u := int(blmath.Map(px+mapWidthF+rotation*4, -1, 1, 0, mapWidthF/2)) % bt.Width
				v := int(blmath.Map(py, -1, 1, 0, mapHeightF))
				r, g, b, a := bt.GetPixel(u, v)
				buffer.SetPixel(int(xx), int(yy), r, g, b, a)
			}
		}
	}
	buffer.CopyToSurface(c.Surface)
}

// PolarMap maps an image to polar coordinates - essentially wraps it around in a donut shape.
// bt is a ByteTexture containing the image data.
// cx, cy is the center of the donut shape.
// outerRadius and innerRadius are the radii of the donut shape.
// cropTop and cropBottom discard that many pixels from the top and bottom of the source image.
// rotation determines the starting point of the wrapped image. If rotation is 0, the image will wrap clockwise from the bottom of the donut.
// mirror will mirror the image, allowing it to blend perfectly all the way around.
func (c *Context) PolarMap(bt ByteTexture, cx, cy, outerRadius, innerRadius, cropTop, cropBottom, rotation float64, mirror bool) {
	buffer, err := ByteTextureFromSurface(c.Surface)
	if err != nil {
		log.Fatalf("unable to create buffer: %s", err)
	}
	for x := cx - outerRadius; x < cx+outerRadius; x++ {
		for y := cy - outerRadius; y < cy+outerRadius; y++ {
			angle := math.Atan2(y-cy, x-cx) - rotation - math.Pi/2
			angle = blmath.WrapTau(angle)
			dist := math.Hypot(x-cx, y-cy)
			if dist <= outerRadius && dist > innerRadius {
				if mirror {
					angle *= 2
					if angle > blmath.Tau {
						angle -= blmath.Tau
						angle = blmath.Tau - angle
					}
				}
				px := int(blmath.Map(angle, 0, blmath.Tau, 0, float64(bt.Width-1)))
				py := int(blmath.Map(dist, innerRadius, outerRadius, float64(bt.Height-1)-cropBottom, cropTop))
				r, g, b, a := bt.GetPixel(px, py)
				// context.SetSourceRGB(r, g, b)
				// context.FillRectangle(x, y, 1, 1)
				buffer.SetPixel(int(x), int(y), r, g, b, a)
			}
		}
	}
	buffer.CopyToSurface(c.Surface)
}
