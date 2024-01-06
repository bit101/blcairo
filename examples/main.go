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
	renderTarget := target.Video

	switch renderTarget {
	case target.Image:
		render.Image(800, 800, "out.png", scene1, 0.0)
		render.ViewImage("out.png")
		break

	case target.SpriteSheet:
		render.SpriteSheet(50, 50, blcolor.White(), "out.png", 36, renderSpriteSheetFrame)
		render.ViewImage("out.png")
		break

	case target.Gif:
		program := render.NewProgram(400, 400, 30)
		program.AddSceneWithFrames(scene1, 60)

		program.Render("frames")
		render.MakeGIF("ffmpeg", "frames", "out.gif", program.Width, program.Height, program.FPS, program.Seconds())
		render.ViewGif("out.gif")
		break

	case target.Video:
		program := render.NewProgram(400, 400, 30)
		program.AddSceneWithFrames(scene1, 60)

		program.Render("frames")
		render.ConvertToVideo("frames", "out.mp4", program.Width, program.Height, program.FPS, program.Seconds())
		render.PlayVideo("out.mp4")
		break
	}
}

func scene1(context *cairo.Context, width, height, percent float64) {
	context.BlackOnWhite()
	context.Save()
	context.TranslateCenter()
	context.DrawAxes(0.25)
	r := blmath.LerpSin(percent, 50, width/2)
	sphere := cairo.NewSphere(0, 0, r, 1, 0, 0)
	sphere.SetShadowColor(0.2, 0, 0)
	sphere.Draw(context)
	context.Restore()
}

func renderSpriteSheetFrame(context *cairo.Context, width, height, percent float64) {
	r := blmath.LerpSin(percent, 5, width/2)
	sphere := cairo.NewSphere(width/2, height/2, r, 1, 0, 0)
	sphere.SetShadowColor(0.2, 0, 0)
	sphere.Draw(context)
}
