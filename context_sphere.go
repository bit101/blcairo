// Package cairo wraps the c cairographics library.
package cairo

import "math"

// Sphere is used to draw a realistic sphere with gradients.
type Sphere struct {
	x, y, Radius              float64
	r, g, b                   float64
	hlX, hlY                  float64
	hlSize                    float64
	shadowR, shadowG, shadowB float64
}

// NewSphere creates a new sphere.
func NewSphere(x, y, radius, r, g, b float64) *Sphere {
	return &Sphere{
		x, y, radius,
		r, g, b,
		0.33, -0.33,
		0.1,
		0, 0, 0,
	}
}

// SetColor sets the r, g, b color of the sphere
func (s *Sphere) SetColor(r, g, b float64) {
	s.r = r
	s.g = g
	s.b = b
}

// SetHighlightPosition sets the relative position of the highlight.
// 0, 0 would be center of the sphere.
// x = 1 moves the highlight to the right side of the sphere, x = -1 moves it to the left side.
// y = 1 moves the hightlight to the bottom, -1 to the top.
// default is 0.33, -0.33
func (s *Sphere) SetHighlightPosition(x, y float64) {
	s.hlX = x
	s.hlY = y
}

// SetHighlightSize sets the size of the highlight.
// 0.0 is no SetHighlightSize. 1.0 will fill the entire sphere with a highlight.
// default is 0.1
func (s *Sphere) SetHighlightSize(size float64) {
	s.hlSize = size
}

// SetShadowColor sets the r, g, b shadow color of the sphere
func (s *Sphere) SetShadowColor(r, g, b float64) {
	s.shadowR = r
	s.shadowG = g
	s.shadowB = b
}

// Draw renders a sphere to a context.
func (s *Sphere) Draw(context *Context) {
	context.Save()
	context.Translate(s.x, s.y)
	grad := CreateRadialGradient(
		s.Radius*s.hlX, s.Radius*s.hlY, 0,
		s.Radius*s.hlX, s.Radius*s.hlY, s.Radius*(1+math.Hypot(s.hlX, s.hlY)),
	)
	grad.AddColorStopRGB(0, 1, 1, 1)
	grad.AddColorStopRGB(s.hlSize, s.r, s.g, s.b)
	grad.AddColorStopRGB(1, s.shadowR, s.shadowG, s.shadowB)
	context.SetSource(grad)

	context.FillCircle(0, 0, s.Radius)
	context.Restore()
}
