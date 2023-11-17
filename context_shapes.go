// Package cairo wraps the c cairographics library.
package cairo

import (
	"fmt"
	"math"

	"github.com/bit101/bitlib/blmath"
	"github.com/bit101/bitlib/geom"
	"github.com/bit101/bitlib/random"
)

////////////////////
// ARC, ETC.
////////////////////

// FillArc draws an arc and fills it.
func (c *Context) FillArc(x, y, r, a1, a2 float64, anticlockwise bool) {
	c.Arc(x, y, r, a1, a2, anticlockwise)
	c.Fill()
}

// StrokeArc draws an arc and strokes it.
func (c *Context) StrokeArc(x, y, r, a1, a2 float64, anticlockwise bool) {
	c.Arc(x, y, r, a1, a2, anticlockwise)
	c.Stroke()
}

// Chord draws a line that represents a chord of the given circle.
func (c *Context) Chord(x, y, r, a1, a2 float64) {
	c.MoveTo(x+math.Cos(a1)*r, y+math.Sin(a1)*r)
	c.LineTo(x+math.Cos(a2)*r, y+math.Sin(a2)*r)
}

// StrokeChord draws and strokes a chord.
func (c *Context) StrokeChord(x, y, r, a1, a2 float64) {
	c.Chord(x, y, r, a1, a2)
	c.Stroke()
}

// CircleSegment draws a segment of the given circle.
func (c *Context) CircleSegment(x, y, r, a1, a2 float64, anticlockwise bool) {
	c.Arc(x, y, r, a1, a2, anticlockwise)
	c.ClosePath()
}

// StrokeCircleSegment draws and strokes a segment of the given circle.
func (c Context) StrokeCircleSegment(x, y, r, a1, a2 float64, anticlockwise bool) {
	c.CircleSegment(x, y, r, a1, a2, anticlockwise)
	c.Stroke()
}

// FillCircleSegment draws and strokes a segment of the given circle.
func (c Context) FillCircleSegment(x, y, r, a1, a2 float64, anticlockwise bool) {
	c.CircleSegment(x, y, r, a1, a2, anticlockwise)
	c.Fill()
}

// CircleSector draws a segment of the given circle.
func (c *Context) CircleSector(x, y, r, a1, a2 float64, anticlockwise bool) {
	c.MoveTo(x, y)
	c.Arc(x, y, r, a1, a2, anticlockwise)
	c.ClosePath()
}

// StrokeCircleSector draws and strokes a segment of the given circle.
func (c *Context) StrokeCircleSector(x, y, r, a1, a2 float64, anticlockwise bool) {
	c.CircleSector(x, y, r, a1, a2, anticlockwise)
	c.Stroke()
}

// FillCircleSector draws and strokes a segment of the given circle.
func (c *Context) FillCircleSector(x, y, r, a1, a2 float64, anticlockwise bool) {
	c.CircleSector(x, y, r, a1, a2, anticlockwise)
	c.Fill()
}

// CircleSectorWithHeight draws an arc of the given height connecting two points
func (c *Context) CircleSectorWithHeight(x0, y0, x1, y1, height float64) {
	circle, err := geom.CircleThroughPointsWithArcHeight(x0, y0, x1, y1, height)
	if err != nil {
		c.MoveTo(x0, y0)
		c.LineTo(x1, y1)
	}
	a0 := math.Atan2(y1-circle.Y, x1-circle.X)
	a1 := math.Atan2(y0-circle.Y, x0-circle.X)
	if circle.Radius > 0 {
		c.Arc(circle.X, circle.Y, circle.Radius, a1, a0, false)
	} else {
		c.Arc(circle.X, circle.Y, -circle.Radius, a1, a0, true)
	}
}

// StrokeCircleSectorWithHeight draws an arc of the given height connecting two points
func (c *Context) StrokeCircleSectorWithHeight(x0, y0, x1, y1, h float64) {
	c.CircleSectorWithHeight(x0, y0, x1, y1, h)
	c.Stroke()
}

// FillCircleSectorWithHeight draws an arc of the given height connecting two points
func (c *Context) FillCircleSectorWithHeight(x0, y0, x1, y1, h float64) {
	c.CircleSectorWithHeight(x0, y0, x1, y1, h)
	c.Fill()
}

////////////////////
// ARROW
////////////////////

// DrawArrow draws and arrow shape.
func (c *Context) DrawArrow(x0, y0, x1, y1, pointSize float64) {
	angle := math.Atan2(y1-y0, x1-x0)
	length := math.Hypot(x1-x0, y1-y0)
	c.Save()
	c.Translate(x0, y0)
	c.Rotate(angle)
	c.MoveTo(0, 0)
	c.LineTo(length, 0)
	c.LineTo(length-pointSize, -pointSize*0.6)
	c.MoveTo(length, 0)
	c.LineTo(length-pointSize, pointSize*0.6)
	c.Restore()
}

// StrokeArrow draws and arrow and strokes it.
func (c *Context) StrokeArrow(x0, y0, x1, y1, pointSize float64) {
	c.DrawArrow(x0, y0, x1, y1, pointSize)
	c.Stroke()
}

// DrawDoubleArrow draws an arrow shape with points at both ends.
func (c *Context) DrawDoubleArrow(x0, y0, x1, y1, pointSize float64) {
	angle := math.Atan2(y1-y0, x1-x0)
	length := math.Hypot(x1-x0, y1-y0)
	c.Save()
	c.Translate(x0, y0)
	c.Rotate(angle)
	c.MoveTo(0, 0)
	c.LineTo(length, 0)
	c.LineTo(length-pointSize, -pointSize*0.6)
	c.MoveTo(length, 0)
	c.LineTo(length-pointSize, pointSize*0.6)
	c.MoveTo(pointSize, pointSize*0.6)
	c.LineTo(0, 0)
	c.LineTo(pointSize, -pointSize*0.6)
	c.Restore()
}

// StrokeDoubleArrow draws an arrow shape with points at both ends and strokes it.
func (c *Context) StrokeDoubleArrow(x0, y0, x1, y1, pointSize float64) {
	c.DrawDoubleArrow(x0, y0, x1, y1, pointSize)
	c.Stroke()
}

////////////////////
// CIRCLE
////////////////////

// Circle draws a circle
func (c *Context) Circle(x, y, r float64) {
	c.Arc(x, y, r, 0.0, blmath.TwoPi, false)
}

// FillCircle draws a circle and fills it.
func (c *Context) FillCircle(x, y, r float64) {
	c.Circle(x, y, r)
	c.Fill()
}

// StrokeCircle draws a circle and strokes it.
func (c *Context) StrokeCircle(x, y, r float64) {
	c.Circle(x, y, r)
	c.Stroke()
}

// FillCircleObject draws a circle and fills it.
func (c *Context) FillCircleObject(circle *geom.Circle) {
	c.Circle(circle.X, circle.Y, circle.Radius)
	c.Fill()
}

// StrokeCircleObject draws a circle and strokes it.
func (c *Context) StrokeCircleObject(circle *geom.Circle) {
	c.Circle(circle.X, circle.Y, circle.Radius)
	c.Stroke()
}

// StrokeCircles strokes each circle object in a list.
func (c *Context) StrokeCircles(circles []*geom.Circle) {
	for _, circle := range circles {
		c.StrokeCircleObject(circle)
	}
}

// FillCircles fills each circle object in a list.
func (c *Context) FillCircles(circles []*geom.Circle) {
	for _, circle := range circles {
		c.FillCircleObject(circle)
	}
}

////////////////////
// Crosshair
////////////////////

// Crosshair draws a crosshair a the specified location
func (c *Context) Crosshair(x, y, size float64) {
	c.MoveTo(x-size, y)
	c.LineTo(x+size, y)
	c.MoveTo(x, y-size)
	c.LineTo(x, y+size)
	c.Stroke()
}

////////////////////
// CurveTo
////////////////////

// StrokeCurveTo draws a stroked curve.
func (c *Context) StrokeCurveTo(x0, y0, x1, y1, x2, y2 float64) {
	c.CurveTo(x0, y0, x1, y1, x2, y2)
	c.Stroke()
}

////////////////////
// QuadraticCurveTo
////////////////////

// QuadraticCurveTo draws a quadratic curve to two points.
func (c *Context) QuadraticCurveTo(x0, y0, x1, y1 float64) {
	px, py := c.GetCurrentPoint()
	c.CurveTo(
		2.0/3.0*x0+1.0/3.0*px,
		2.0/3.0*y0+1.0/3.0*py,
		2.0/3.0*x0+1.0/3.0*x1,
		2.0/3.0*y0+1.0/3.0*y1,
		x1, y1,
	)
}

// StrokeQuadraticCurveTo draws a stroked quadratic curve.
func (c *Context) StrokeQuadraticCurveTo(x0, y0, x1, y1 float64) {
	c.QuadraticCurveTo(x0, y0, x1, y1)
	c.Stroke()
}

// QuadraticCurveThrough draws a quad curve that goes through the control point.
func (c *Context) QuadraticCurveThrough(x0, y0, x1, y1 float64) {
	px, py := c.GetCurrentPoint()
	xc := x0*2 - px/2 - x1/2
	yc := y0*2 - py/2 - y1/2
	c.QuadraticCurveTo(xc, yc, x1, y1)
}

// StrokeQuadraticCurveThrough draws and strokes a quad curve that goes through the control point.
func (c *Context) StrokeQuadraticCurveThrough(x0, y0, x1, y1 float64) {
	c.QuadraticCurveThrough(x0, y0, x1, y1)
	c.Stroke()
}

////////////////////
// ELLIPSE
////////////////////

// Ellipse draws an ellipse.
func (c *Context) Ellipse(x, y, xr, yr float64) {
	if xr == 0 || yr == 0 {
		return
	}
	c.Save()
	c.Translate(x, y)
	c.Scale(xr, yr)
	c.Circle(0.0, 0.0, 1.0)
	c.Restore()
}

// FillEllipse draws an ellipse and fills it.
func (c *Context) FillEllipse(x, y, rx, ry float64) {
	c.Ellipse(x, y, rx, ry)
	c.Fill()
}

// StrokeEllipse draws an ellipse and strokes it.
func (c *Context) StrokeEllipse(x, y, rx, ry float64) {
	c.Ellipse(x, y, rx, ry)
	c.Stroke()
}

////////////////////
// FRACTAL LINE
////////////////////

// FractalLine draws a rough, fractal line between two points.
func (c *Context) FractalLine(x1, y1, x2, y2, roughness float64, iterations int) {
	dx := x2 - x1
	dy := y2 - y1
	offset := math.Sqrt(dx*dx+dy*dy) * 0.15

	var path []*geom.Point
	path = append(path, geom.NewPoint(x1, y1))
	path = append(path, geom.NewPoint(x2, y2))

	for i := 0; i < iterations; i++ {
		var newPath []*geom.Point
		for j, point := range path {
			newPath = append(newPath, geom.NewPoint(point.X, point.Y))
			if j < len(path)-1 {
				mid := geom.MidPoint(point, path[j+1])
				mid.X += random.FloatRange(-offset, offset)
				mid.Y += random.FloatRange(-offset, offset)
				newPath = append(newPath, mid)
			}
		}
		offset *= roughness
		path = newPath
	}
	c.Path(path)
}

// StrokeFractalLine draws a fractal line between two points and strokes it.
func (c *Context) StrokeFractalLine(x1, y1, x2, y2, roughness float64, iterations int) {
	c.FractalLine(x1, y1, x2, y2, roughness, iterations)
	c.Stroke()
}

////////////////////
// GRID
////////////////////

// Grid draws a grid on the context (stroked).
func (c *Context) Grid(x, y, w, h, xres, yres float64) {
	xx := x
	yy := y
	for xx <= x+w {
		c.MoveTo(xx, y)
		c.LineTo(xx, y+h)
		xx += xres
	}
	for yy <= y+h {
		c.MoveTo(x, yy)
		c.LineTo(x+w, yy)
		yy += yres
	}
	c.Stroke()
}

// GridFull draws a square grid with the given resolution and line width across the whole canvas.
func (c *Context) GridFull(res, lineWidth float64) {
	c.Save()
	c.SetMatrix(*NewMatrix())
	c.SetLineWidth(lineWidth)
	c.Grid(0, 0, c.Width, c.Height, res, res)
	c.Restore()
}

// Checkerboard fills an area with a checkerboard pattern.
func (c *Context) Checkerboard(x, y, w, h, resX, resY float64, topLeftFilled bool) {
	c.Save()
	c.Rectangle(x, y, w, h)
	c.Clip()
	offset := resX
	if topLeftFilled {
		offset = 0.0
	}
	for yy := y; yy < y+h; yy += resY {
		for xx := x; xx < x+w; xx += resX * 2 {
			c.FillRectangle(xx+offset, yy, resX, resY)
		}
		if offset == 0.0 {
			offset = resX
		} else {
			offset = 0.0
		}
	}
	c.ResetClip()
	c.Restore()
}

////////////////////
// GUIDES
////////////////////

// DrawAxes draws a vertical and horizontal line for each axis at the origin.
func (c *Context) DrawAxes(lineWidth float64) {
	origLineWidth := c.GetLineWidth()
	c.SetLineWidth(lineWidth)
	w, h := c.Size()
	c.MoveTo(0, -h)
	c.LineTo(0, h)
	c.MoveTo(-w, 0)
	c.LineTo(w, 0)
	c.Stroke()
	c.SetLineWidth(origLineWidth)
}

// DrawOrigin draws a "plus" mark at the origin, with the given size.
func (c *Context) DrawOrigin(size float64) {
	c.MoveTo(0, -size/2)
	c.LineTo(0, size/2)
	c.MoveTo(-size/2, 0)
	c.LineTo(size/2, 0)
	c.Stroke()
}

////////////////////
// GUILLOCHE
////////////////////

// Guillloche draws a Guillloche pattern.
func (c *Context) Guillloche(x, y, innerRadius, innerAmp, innerCycles, innerRotation, outerRadius, outerAmp, outerCycles, outerRotation, nodes, div float64) {
	for t := 0.0; t < blmath.Tau*div; t += 0.01 {
		r0 := innerRadius + math.Sin(t*innerCycles+innerRotation)*innerAmp
		r1 := outerRadius + math.Sin(t*outerCycles+outerRotation)*outerAmp
		rng := (r1 - r0) * 0.5
		mid := r0 + rng
		radius := mid + math.Sin(t*nodes/div)*rng
		c.LineTo(x+math.Cos(t)*radius, y+math.Sin(t)*radius)
	}
	c.Stroke()
}

////////////////////
// Hatch
////////////////////

// FillHatchAlt fills an area with grid of alternating hatch patterns.
func (c *Context) FillHatchAlt(x, y, w, h, size float64, count int, r0, r1, rrand, posRand float64) {
	c.Save()
	c.Translate(x, y)
	c.Rectangle(0, 0, w, h)
	c.Clip()
	xres := w / size
	yres := h / size
	for x1 := 0.0; x1 < xres; x1++ {
		for y1 := 0.0; y1 < yres; y1++ {
			angle := r0
			if math.Mod(x1, 2) == math.Mod(y1, 2) {
				angle = r1
			}
			angle += random.FloatRange(-rrand, rrand)
			c.Hatch(count, x1*size, y1*size, size, angle, posRand)
		}
	}
	c.ResetClip()
	c.Restore()
}

// Hatch fills a square with a hatch pattern.
func (c *Context) Hatch(count int, x, y, size, rotation, rand float64) {
	// x, y, size forms a square
	// get a line that runs through the center of the square at the given rotation.
	// this line will be perpendicular to the strokes we will create.
	rotation = blmath.WrapTau(rotation)
	cos := math.Cos(rotation - math.Pi/2)
	sin := math.Sin(rotation - math.Pi/2)
	cx := x + size/2
	cy := y + size/2
	cs := cos * size
	ss := sin * size
	line := geom.NewLine(cx-cs, cy-ss, cx+cs, cy+ss)

	// calculate what portion of that line we will use to draw perpendicular strokes (w)
	// and how far apart those strokes will be (space)
	fcount := float64(count)
	a := math.Mod(rotation, math.Pi/2) - math.Pi/4
	w := math.Cos(a) * size * math.Sqrt2
	space := w / fcount

	for i := 0.0; i < fcount; i++ {
		// find the points along the line where strokes will be drawn
		n := (space-w)*0.5 + space*i
		stroke := line.Perpendicular(geom.NewPoint(cx+cos*n, cy+sin*n))
		// find where that stroke intersects the enclosing square
		points := geom.LineOnRect(
			stroke.X0, stroke.Y0, stroke.X1, stroke.Y1,
			x, y, size, size,
		)
		// ensure it intersects in two points
		if len(points) >= 2 {
			// randomize those points and draw a line segment
			p0 := points[0]
			p1 := points[1]
			c.MoveTo(p0.X+random.FloatRange(-rand, rand), p0.Y+random.FloatRange(-rand, rand))
			c.LineTo(p1.X+random.FloatRange(-rand, rand), p1.Y+random.FloatRange(-rand, rand))
		}
	}
	c.Stroke()

}

////////////////////
// HEART
////////////////////

// Heart draws a heart shape.
func (c *Context) Heart(x, y, w, h, r float64) {
	c.Save()
	c.Translate(x, y)
	c.Rotate(r)
	var path []*geom.Point
	res := math.Sqrt(w * h)
	for i := 0.0; i < res; i++ {
		a := math.Pi * 2 * i / res
		x := w * math.Pow(math.Sin(a), 3.0)
		y := h * (0.8125*math.Cos(a) - 0.3125*math.Cos(2.0*a) - 0.125*math.Cos(3.0*a) - 0.0625*math.Cos(4.0*a))
		path = append(path, geom.NewPoint(x, -y))
	}
	c.Path(path)
	c.Restore()
}

// FillHeart draws a heart shape and fills it.
func (c *Context) FillHeart(x, y, w, h, r float64) {
	c.Heart(x, y, w, h, r)
	c.Fill()
}

// StrokeHeart draws a heart shape and strokes it.
func (c *Context) StrokeHeart(x, y, w, h, r float64) {
	c.Heart(x, y, w, h, r)
	c.Stroke()
}

////////////////////
// HEX GRID
////////////////////

// HexGrid draws a hexagonal grid.
func (c *Context) HexGrid(x, y, w, h, res0, res1 float64) {
	sin60r := math.Sin(math.Pi/3.0) * res0
	xInc := 2.0 * sin60r
	yInc := res0 * 1.5
	offset := 0.0

	for yy := y; yy < y+h+yInc; yy += yInc {
		for xx := x; xx < x+w+xInc; xx += xInc {
			c.Polygon(xx+offset, yy, res1, 6, math.Pi/2)
		}
		if offset == 0 {
			offset = sin60r
		} else {
			offset = 0
		}
	}
}

// FillHexGrid draws a hex grid and fills it.
func (c *Context) FillHexGrid(x, y, w, h, res0, res1 float64) {
	c.Save()
	c.Rectangle(x, y, w, h)
	c.Clip()
	c.HexGrid(x, y, w, h, res0, res1)
	c.Fill()
	c.ResetClip()
	c.Restore()
}

// StrokeHexGrid draws a hex grid and strokes it.
func (c *Context) StrokeHexGrid(x, y, w, h, res0, res1 float64) {
	c.Save()
	c.Rectangle(x, y, w, h)
	c.Clip()
	c.HexGrid(x, y, w, h, res0, res1)
	c.Stroke()
	c.ResetClip()
	c.Restore()
}

////////////////////
// LINE
////////////////////

// Line draws a line between two x, y points.
func (c *Context) Line(x0, y0, x1, y1 float64) {
	c.MoveTo(x0, y0)
	c.LineTo(x1, y1)
	c.Stroke()
}

// StrokeLine strokes a line between two points.
func (c *Context) StrokeLine(x0, y0, x1, y1 float64) {
	c.Line(x0, y0, x1, y1)
	c.Stroke()
}

// StrokeLineObject strokes a line between two points.
func (c *Context) StrokeLineObject(line *geom.Line) {
	c.LineThrough(line.X0, line.Y0, line.X1, line.Y1, c.Width+c.Height)
}

// StrokeSegmentObject strokes a line segment between two points.
func (c *Context) StrokeSegmentObject(seg *geom.Segment) {
	c.StrokeLine(seg.X0, seg.Y0, seg.X1, seg.Y1)
}

// LineThrough strokes a line that extends a certain distance beyond two points.
func (c *Context) LineThrough(x0, y0, x1, y1, overlap float64) {
	c.Save()
	c.Translate(x0, y0)
	c.Rotate(math.Atan2(y1-y0, x1-x0))
	p2 := math.Hypot(x0-x1, y0-y1)

	c.MoveTo(-overlap, 0)
	c.LineTo(p2+overlap, 0)
	c.Stroke()
	c.Restore()
}

////////////////////
// MULTI CURVE
////////////////////

// MultiCurve draws a piecewise bezier curve through a series of points.
func (c *Context) MultiCurve(points []*geom.Point) {
	c.MoveTo(points[0].X, points[0].Y)
	mid := geom.MidPoint(points[0], points[1])
	c.LineTo(mid.X, mid.Y)
	i := 1
	for i < len(points)-1 {
		p0 := points[i]
		p1 := points[i+1]
		mid := geom.MidPoint(p0, p1)
		c.QuadraticCurveTo(p0.X, p0.Y, mid.X, mid.Y)
		i++

	}
	p := points[len(points)-1]
	c.LineTo(p.X, p.Y)
}

// StrokeMultiCurve draws a multi curve and strokes it.
func (c *Context) StrokeMultiCurve(points []*geom.Point) {
	c.MultiCurve(points)
	c.Stroke()
}

////////////////////
// MULTI LOOP
////////////////////

// MultiLoop draws a closed piecewise bezier curve through a series of points.
func (c *Context) MultiLoop(points []*geom.Point) {
	pA := points[0]
	pZ := points[len(points)-1]
	mid1 := geom.MidPoint(pZ, pA)
	c.MoveTo(mid1.X, mid1.Y)
	for i := 0; i < len(points)-1; i++ {
		p0 := points[i]
		p1 := points[i+1]
		mid := geom.MidPoint(p0, p1)
		c.QuadraticCurveTo(p0.X, p0.Y, mid.X, mid.Y)
	}
	c.QuadraticCurveTo(pZ.X, pZ.Y, mid1.X, mid1.Y)
}

// FillMultiLoop draws a filled, smooth, closed curve between a set of points.
func (c *Context) FillMultiLoop(points []*geom.Point) {
	c.MultiLoop(points)
	c.Fill()
}

// StrokeMultiLoop draws a stroked, smooth, closed curve between a set of points.
func (c *Context) StrokeMultiLoop(points []*geom.Point) {
	c.MultiLoop(points)
	c.Stroke()
}

////////////////////
// PATH
////////////////////

// Path draws a series of lines through a set of points.
func (c *Context) Path(points []*geom.Point) {
	for _, point := range points {
		c.LineTo(point.X, point.Y)
	}
}

// FillPath draws a path and fills it.
func (c *Context) FillPath(points []*geom.Point) {
	c.Path(points)
	c.Fill()
}

// StrokePath draws a path and strokes it.
func (c *Context) StrokePath(points []*geom.Point, close bool) {
	c.Path(points)
	if close {
		c.ClosePath()
	}
	c.Stroke()
}

////////////////////
// POINT
////////////////////

// StrokePoint draws a circle at a point and strokes it.
func (c *Context) StrokePoint(x, y, r float64) {
	c.StrokeCircle(x, y, r)
}

// FillPoint  draws a cricle at a point and fills it.
func (c *Context) FillPoint(x, y, r float64) {
	c.FillCircle(x, y, r)
}

// Points draws and fills a circle at each point in a list of points.
func (c *Context) Points(points []*geom.Point, radius float64) {
	for _, point := range points {
		c.FillPoint(point.X, point.Y, radius)
	}
}

// LabelPoints draws a single char near each of a list of points, incrementing the first char by one each time.
// Makes no correction for chars higher than 0-9, a-z or A-Z at this point
func (c *Context) LabelPoints(points []*geom.Point, numeric bool) {

	if numeric {
		num := 0
		for _, p := range points {
			label := fmt.Sprintf("%d", num)
			c.FillText(label, p.X+5, p.Y-5)
			num++
		}
	} else {
		char := 'A'
		count := 0
		for _, p := range points {
			label := string(char)
			if count > 0 {
				label = fmt.Sprintf("%c%d", char, count)
			}
			c.FillText(label, p.X+5, p.Y-5)
			char++
			if char > 'Z' {
				char = 'A'
				count++
			}
		}
	}
}

////////////////////
// Polygon
////////////////////

// Polygon draws a polygon.
func (c *Context) Polygon(x, y, r float64, sides int, rotation float64) {
	c.Save()
	c.Translate(x, y)
	c.Rotate(rotation)
	c.MoveTo(r, 0.0)
	for i := 0; i < sides; i++ {
		angle := blmath.TwoPi / float64(sides) * float64(i)
		c.LineTo(math.Cos(angle)*r, math.Sin(angle)*r)
	}
	c.LineTo(r, 0.0)
	c.Restore()
}

// StrokePolygon draws a stroked polygon.
func (c *Context) StrokePolygon(x, y, r float64, sides int, rotation float64) {
	c.Polygon(x, y, r, sides, rotation)
	c.Stroke()
}

// FillPolygon draws a filled polygon.
func (c *Context) FillPolygon(x, y, r float64, sides int, rotation float64) {
	c.Polygon(x, y, r, sides, rotation)
	c.Fill()
}

////////////////////
// RAY
////////////////////

// Ray draws a line segment extending from a point at an angle.
func (c *Context) Ray(x, y, angle, offset, length float64) {
	c.Save()
	c.Translate(x, y)
	c.Rotate(angle)
	c.MoveTo(offset, 0)
	c.LineTo(offset+length, 0)
	c.Stroke()
	c.Restore()
}

////////////////////
// RECTANGLE
////////////////////

// FillRectangle draws a rectangle and fills it.
func (c *Context) FillRectangle(x, y, w, h float64) {
	c.Rectangle(x, y, w, h)
	c.Fill()
}

// StrokeRectangle draws a rectangle and strokes it.
func (c *Context) StrokeRectangle(x, y, w, h float64) {
	c.Rectangle(x, y, w, h)
	c.Stroke()
}

// StrokeRectangleObject draws and strokes a geom.Rect
func (c *Context) StrokeRectangleObject(rect *geom.Rect) {
	c.Rectangle(rect.X, rect.Y, rect.W, rect.H)
	c.Stroke()
}

// FillRectangleObject draws and fills a geom.Rect
func (c *Context) FillRectangleObject(rect *geom.Rect) {
	c.Rectangle(rect.X, rect.Y, rect.W, rect.H)
	c.Fill()
}

////////////////////
// RIGHT TRIANGLE
////////////////////

// DrawRightTriangle draws a right triangle.
func (c *Context) DrawRightTriangle(x, y, w, h, r float64) {
	c.Save()
	c.Translate(x, y)
	c.Rotate(r)
	c.MoveTo(0, 0)
	c.LineTo(w, 0)
	c.LineTo(0, h)
	c.LineTo(0, 0)
	c.Restore()
}

// StrokeRightTriangle draws a right triangle and strokes it.
func (c *Context) StrokeRightTriangle(x, y, w, h, r float64) {
	c.DrawRightTriangle(x, y, w, h, r)
	c.Stroke()
}

// FillRightTriangle draws a right triangle and fills it.
func (c *Context) FillRightTriangle(x, y, w, h, r float64) {
	c.DrawRightTriangle(x, y, w, h, r)
	c.Fill()
}

////////////////////
// ROUNDED RECTANGLE
////////////////////

// RoundRectangle draws a rounded rectangle.
func (c *Context) RoundRectangle(x, y, w, h, r float64) {
	c.MoveTo(x+r, y)
	c.LineTo(x+w-r, y)
	c.Arc(x+w-r, y+r, r, -blmath.HalfPi, 0.0, false)
	c.LineTo(x+w, y+h-r)
	c.Arc(x+w-r, y+h-r, r, 0.0, blmath.HalfPi, false)
	c.LineTo(x+r, y+h)
	c.Arc(x+r, y+h-r, r, blmath.HalfPi, math.Pi, false)
	c.LineTo(x, y+r)
	c.Arc(x+r, y+r, r, math.Pi, -blmath.HalfPi, false)
}

// StrokeRoundRectangle draws a stroked, rounded rectangle.
func (c *Context) StrokeRoundRectangle(x, y, w, h, r float64) {
	c.RoundRectangle(x, y, w, h, r)
	c.Stroke()
}

// FillRoundRectangle draws a filled, rounded rectangle.
func (c *Context) FillRoundRectangle(x, y, w, h, r float64) {
	c.RoundRectangle(x, y, w, h, r)
	c.Fill()
}

// sign is a helper function to determine the sign of a value.
func sign(v float64) float64 {
	if v < 0 {
		return -1.0
	}
	if v > 0 {
		return 1.0
	}
	return 0
}

// Superellipse draws a rounded rectangle with a superellipse formula.
// A p of 1 makes a diamond. A p of 2 makes an ellipse. Higher p's get closer to square corners.
// This shape is drawn from the center of xc, yc.
func (c *Context) Superellipse(xc, yc, rx, ry, p float64) {
	for t := 0.0; t < blmath.Tau; t += 0.01 {
		cos := math.Cos(t)
		sin := math.Sin(t)

		x := math.Pow(math.Abs(cos), 2.0/p) * rx * sign(cos)
		y := math.Pow(math.Abs(sin), 2.0/p) * ry * sign(sin)
		c.LineTo(xc+x, yc+y)
	}
	c.ClosePath()
}

// StrokeSuperellipse draws and strokes a superellipse.
func (c *Context) StrokeSuperellipse(xc, yc, w, h, p float64) {
	c.Superellipse(xc, yc, w, h, p)
	c.Stroke()
}

// FillSuperellipse draws and fills a superellipse.
func (c *Context) FillSuperellipse(xc, yc, w, h, p float64) {
	c.Superellipse(xc, yc, w, h, p)
	c.Fill()
}

// Squircle draws a blend of a circle and square - a rounded square. It is a superellipse with a p of 4.
// This shape is drawn from the center of xc, yc.
func (c *Context) Squircle(xc, yc, radius float64) {
	c.Superellipse(xc, yc, radius, radius, 4)
}

// StrokeSquircle draws a squircle and strokes it.
func (c *Context) StrokeSquircle(xc, yc, radius float64) {
	c.StrokeSuperellipse(xc, yc, radius, radius, 4)
}

// FillSquircle draws a squircle and fills it.
func (c *Context) FillSquircle(xc, yc, radius float64) {
	c.FillSuperellipse(xc, yc, radius, radius, 4)
}

// Superformula draws a superformula shape.
func (c *Context) Superformula(x, y, radius, symmetry, n1, n2, n3 float64) {
	for t := 0.0; t < blmath.Tau; t += 0.01 {
		angle := symmetry * t / 4.0
		term1 := math.Pow(math.Abs(math.Cos(angle)), n2)
		term2 := math.Pow(math.Abs(math.Sin(angle)), n3)
		r := math.Pow(term1+term2, -1/n1) * radius
		c.LineTo(x+math.Cos(t)*r, y+math.Sin(t)*r)
	}
	c.ClosePath()
}

// StrokeSuperformula draws and strokes a superformula.
func (c *Context) StrokeSuperformula(xc, yc, radius, symmetry, n1, n2, n3 float64) {
	c.Superformula(xc, yc, radius, symmetry, n1, n2, n3)
	c.Stroke()
}

// FillSuperformula draws and fills a superformula.
func (c *Context) FillSuperformula(xc, yc, radius, symmetry, n1, n2, n3 float64) {
	c.Superformula(xc, yc, radius, symmetry, n1, n2, n3)
	c.Fill()
}

////////////////////
// STAR
////////////////////

// Star draws a star shape.
func (c *Context) Star(x, y, r0, r1 float64, points int, rotation float64) {
	c.Save()
	c.Translate(x, y)
	c.Rotate(rotation)
	for i := 0; i < points*2; i++ {
		r := r1
		if i%2 == 1 {
			r = r0
		}
		angle := math.Pi / float64(points) * float64(i)
		c.LineTo(math.Cos(angle)*r, math.Sin(angle)*r)
	}
	c.ClosePath()
	c.Restore()
}

// StrokeStar draws a star and strokes it.
func (c *Context) StrokeStar(x, y, r0, r1 float64, points int, rotation float64) {
	c.Star(x, y, r0, r1, points, rotation)
	c.Stroke()
}

// FillStar draws a star and fills it.
func (c *Context) FillStar(x, y, r0, r1 float64, points int, rotation float64) {
	c.Star(x, y, r0, r1, points, rotation)
	c.Fill()
}

////////////////////
// Text
////////////////////

// FillText draws text
func (c *Context) FillText(text string, x, y float64) {
	c.Save()
	c.Translate(x, y)
	c.ShowText(text)
	c.Fill()
	c.Restore()
}

// StrokeText draws text
func (c *Context) StrokeText(text string, x, y float64) {
	c.Save()
	c.Translate(x, y)
	c.ShowText(text)
	c.Stroke()
	c.Restore()
}

////////////////////
// Triangle
////////////////////

// StrokeTriangleObject strokes a triangle.
func (c *Context) StrokeTriangleObject(t *geom.Triangle) {
	c.StrokePath(t.Points(), true)
}

// FillTriangleObject fills a triangle.
func (c *Context) FillTriangleObject(t *geom.Triangle) {
	c.FillPath(t.Points())
}
