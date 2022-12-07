// Package main renders an image, gif or video
package main

import (
	"github.com/bit101/bitlib/blcolor"
	"github.com/bit101/bitlib/blmath"
	cairo "github.com/bit101/blcairo"
	"github.com/bit101/blcairo/render"
	"github.com/bit101/blcairo/target"
)

func main() {
	renderTarget := target.Image

	switch renderTarget {
	case target.Image:
		render.Image(800, 800, "out.png", renderFrame, 0.0)
		render.ViewImage("out.png")
		break

	case target.SpriteSheet:
		render.SpriteSheet(40, 40, blcolor.White(), "out.png", 25, renderSpriteSheetFrame)
		render.ViewImage("out.png")
		break

	case target.Gif:
		render.Frames(400, 400, 60, "frames", renderFrame)
		render.MakeGIF("ffmpeg", "frames", "out.gif", 30)
		render.ViewImage("out.gif")
		break

	case target.Video:
		render.Frames(1280, 800, 60, "frames", renderFrame)
		render.ConvertToYoutube("frames", "out.mp4", 60)
		render.VLC("out.mp4", true)
		break
	}
}

func renderFrame(context *cairo.Context, width, height, percent float64) {
	context.BlackOnWhite()
	context.Save()
	context.TranslateCenter()
	context.DrawAxes(0.25)
	r := blmath.LerpSin(percent, 0, width/2)
	context.FillCircle(0, 0, r)
	context.Restore()
}

func renderSpriteSheetFrame(context *cairo.Context, width, height, percent float64) {
	context.SetSourceBlack()
	r := blmath.LerpSin(percent, 2, width*0.45)
	context.FillCircle(width/2, height/2, r)
}
