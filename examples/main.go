// Package main contains an executable
package main

// Simple test for cairo package

import (
	"math"

	"github.com/bit101/bitlib/random"
	cairo "github.com/bit101/blcairo"
)

func main() {
	surface := cairo.NewSurface(800, 800)
	context := cairo.NewContext(surface)
	context.BlackOnWhite()
	text(context)
	shapes(context)
	surface.WriteToPNG("out.png")
	surface.Finish()
}

func text(context *cairo.Context) {
	context.SelectFontFace("Arial", cairo.FontSlantNormal, cairo.FontWeightBold)
	context.SetFontSize(32.0)
	context.FillText("Hello World", 10, 50)
}

func shapes(context *cairo.Context) {
	context.Translate(0, 50)

	context.FillRectangle(10, 10, 100, 100)

	context.StrokeRectangle(120, 10, 100, 100)

	context.FillCircle(280, 60, 50)

	context.StrokeCircle(390, 60, 50)

	for i := 0; i < 50; i++ {
		context.MoveTo(random.FloatRange(450, 550), random.FloatRange(10, 110))
		context.LineTo(random.FloatRange(450, 550), random.FloatRange(10, 110))
		context.Stroke()
	}

	context.MoveTo(10, 120)
	context.StrokeCurveTo(590, 120, 10, 220, 590, 220)

	context.SetLineWidth(0.15)
	for i := 1; i < 6; i++ {
		random.Seed(0)
		context.StrokeFractalLine(10, 350, 790, 350, 0.6, i)
	}
	random.Seed(0)
	context.SetLineWidth(0.8)
	context.StrokeFractalLine(10, 350, 790, 350, 0.6, 6)
}

func colors() {
	surface := cairo.NewSurface(600, 600)
	context := cairo.NewContext(surface)
	for i := 0.0; i < 100; i++ {
		for j := 0.0; j < 100; j++ {
			dist := math.Hypot(i*6-300, j*6-300)
			red := i / 100
			green := math.Max(0, 1.0-dist/200)
			blue := j / 100
			context.SetSourceRGB(red, green, blue)
			context.Rectangle(i*6, j*6, 6, 6)
			context.Fill()
		}
	}

	surface.WriteToPNG("out.png")
	surface.Finish()
}

func gradients() {
	surface := cairo.NewSurface(600, 300)
	context := cairo.NewContext(surface)
	radialPattern := cairo.CreateRadialGradient(150, 150, 0, 150, 150, 150)
	radialPattern.AddColorStopRGB(0, 1, 0, 0)
	radialPattern.AddColorStopRGB(1, 0, 0, 1)
	context.SetSource(radialPattern)
	context.Rectangle(0, 0, 300, 300)
	context.Fill()

	linearPattern := cairo.CreateLinearGradient(300, 0, 600, 300)
	linearPattern.AddColorStopRGB(0, 1, 0, 0)
	linearPattern.AddColorStopRGB(1, 0, 0, 1)
	context.SetSource(linearPattern)
	context.Rectangle(300, 0, 300, 300)
	context.Fill()

	surface.WriteToPNG("out.png")
	surface.Finish()
}

func mesh() {
	surface := cairo.NewSurface(600, 600)
	context := cairo.NewContext(surface)
	pattern := cairo.CreateMesh()

	pattern.BeginPatch()
	pattern.MoveTo(100, 100)
	pattern.LineTo(500, 100)
	pattern.LineTo(500, 500)
	pattern.LineTo(100, 500)

	pattern.SetCornerColorRGB(0, 1, 0, 0)
	pattern.SetCornerColorRGB(1, 0, 1, 0)
	pattern.SetCornerColorRGB(2, 0, 0, 1)
	pattern.SetCornerColorRGB(3, 1, 1, 0)
	pattern.EndPatch()

	context.SetSource(pattern)
	context.Rectangle(0, 0, 600, 600)
	context.Fill()

	surface.WriteToPNG("out.png")
	surface.Finish()
}
