package iso

import (
	"math"

	cairo "github.com/bit101/blcairo"
)

var leftMatrix = cairo.NewMatrixWithValues(
	1, 0.5,
	0, 1,
	0, 0,
)

var rightMatrix = cairo.NewMatrixWithValues(
	1, -0.5,
	0, 1,
	0, 0,
)

// Box is a box
type Box struct {
	X     float64
	Y     float64
	Z     float64
	W     float64
	D     float64
	H     float64
	Top   SideRender
	Left  SideRender
	Right SideRender
}

type SideRender func(*cairo.Context, float64, float64, float64, float64)

// NewBox creates a new box of the given size
func NewBox(w, d, h float64) *Box {
	return &Box{
		X:     0,
		Y:     0,
		Z:     0,
		W:     w,
		D:     d,
		H:     h,
		Top:   TopColorRender,
		Left:  LeftColorRender,
		Right: RightColorRender,
	}
}

func TopColorRender(context *cairo.Context, x, y, w, h float64) {
	context.SetSourceRGB(1, 1, 1)
	context.FillRectangle(x, y, w, h)
}

func LeftColorRender(context *cairo.Context, x, y, w, h float64) {
	context.SetSourceRGB(0.5, 0.5, 0.5)
	context.FillRectangle(x, y, w, h)
}

func RightColorRender(context *cairo.Context, x, y, w, h float64) {
	context.SetSourceRGB(0.75, 0.75, 0.75)
	context.FillRectangle(x, y, w, h)
}

// Position positions the box
func (box *Box) Position(x, y, z float64) {
	box.X = x
	box.Y = y
	box.Z = z
}

// Size sizes the box
func (box *Box) Size(w, d, h float64) {
	box.W = w
	box.D = d
	box.H = h
}

// Render renders the box to a context
func (box *Box) Render(context *cairo.Context) {
	context.Save()
	context.Translate(box.X, box.Y-box.Z)

	box.drawBack(context)
	box.drawLeftWall(context)
	box.drawRightWall(context)
	box.drawTop(context)
	context.Restore()
}

func (box *Box) drawBack(context *cairo.Context) {
	// draw a triangle across all face seams to prevent background color leaking through
	context.Save()
	context.SetSourceRGB(0.5, 0.5, 0.5)
	context.MoveTo(0, 0)
	context.LineTo(-box.D, -box.H-box.D/2)
	context.LineTo(box.W, -box.H-box.W/2)
	context.Fill()
	context.Restore()
}

func (box *Box) drawLeftWall(context *cairo.Context) {
	context.Save()
	context.Transform(*leftMatrix)
	box.Left(context, -box.D, -box.H, box.D, box.H)
	context.Restore()
}

func (box *Box) drawRightWall(context *cairo.Context) {
	context.Save()
	context.Transform(*rightMatrix)
	box.Right(context, 0, -box.H, box.W, box.H)
	context.Restore()
}

func (box *Box) drawTop(context *cairo.Context) {
	context.Save()
	context.Translate(0, -box.H)
	context.Scale(1, 0.5)
	context.Scale(math.Sqrt2, math.Sqrt2)
	context.Rotate(math.Pi / 4)
	box.Top(context, -box.D, -box.W, box.D, box.W)
	context.Restore()
}
