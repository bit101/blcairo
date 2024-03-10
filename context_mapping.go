// Package cairo wraps the c cairographics library.
package cairo

import (
	"math"

	"github.com/bit101/bitlib/blmath"
)

// SphereMap maps a ByteTexture onto a sphere.
// bt is the ByteTexture instance.
// x, y is the top left of the sphere bounding box, diameter is the diameter of the sphere.
// rotation is normalized. -1.0 is a full rotation to the right, 1.0 is a full rotation to the left.
func (c *Context) SphereMap(bt ByteTexture, x, y, diameter, rotation float64) {
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
				c.SetSourceRGBA(r, g, b, a)
				c.FillRectangle(xx, yy, 1, 1)
			}
		}
	}

}
