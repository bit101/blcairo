// Package render renders a single image or a number of frames
package render

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/bit101/bitlib/blcolor"
	cairo "github.com/bit101/blcairo"
	"github.com/bit101/go-ansi"
)

var startTime time.Time

// FrameFunc is the interface for a function that renders a single frame.
type FrameFunc func(*cairo.Context, float64, float64, float64)

// Image sets up the rendering of a single image.
func Image(width, height float64, path string, frameFunc FrameFunc, percent float64) {
	surface := cairo.NewSurface(int(width), int(height))
	context := cairo.NewContext(surface)
	frameFunc(context, width, height, percent)
	surface.WriteToPNG(path)
}

// Frames sets up the renderin of a series of frames.
func Frames(width, height float64, numFrames int, frames string, frameFunc FrameFunc) {
	initProgress()
	os.RemoveAll(frames)
	os.MkdirAll(frames, 0755)
	surface := cairo.NewSurface(int(width), int(height))
	context := cairo.NewContext(surface)
	for frame := 0; frame < numFrames; frame++ {
		percent := float64(frame) / float64(numFrames)
		setProgress(percent)
		frameFunc(context, width, height, percent)
		surface.WriteToPNG(fmt.Sprintf("%s/frame_%04d.png", frames, frame))
	}
	setComplete()
}

// CleanFrames cleans out frames.
func CleanFrames(frames string) {
	os.RemoveAll(frames)
	os.MkdirAll(frames, 0755)
}

// FrameRange renders a range of frames
func FrameRange(width, height float64, numFrames, start, end int, frames string, frameFunc FrameFunc) {
	initProgress()
	surface := cairo.NewSurface(int(width), int(height))
	context := cairo.NewContext(surface)
	for frame := start; frame <= end; frame++ {
		percent := float64(frame) / float64(numFrames)
		setProgress(percent)
		frameFunc(context, width, height, percent)
		surface.WriteToPNG(fmt.Sprintf("%s/frame_%04d.png", frames, frame))
	}
	setComplete()
}

// SpriteSheet sets up the rendering of a sprite sheet.
func SpriteSheet(width, height float64, bg blcolor.Color, path string, numFrames int, frameFunc FrameFunc) {
	x := 0.0
	y := 0.0
	nf := float64(numFrames)
	size := math.Ceil(math.Sqrt(nf))
	surface := cairo.NewSurface(int(width*size), int(height*size))
	context := cairo.NewContext(surface)
	context.ClearColor(bg)

	for i := 0.0; i < nf; i++ {
		context.Save()
		context.Translate(x, y)
		percent := i / float64(numFrames)
		frameFunc(context, width, height, percent)
		context.Restore()

		x += width
		if x >= size*width {
			x = 0.0
			y += height
		}
	}
	surface.WriteToPNG(path)
}

func initProgress() {
	ansi.ClearScreen()
	ansi.SetScrollRegion(4, 1000)
	ansi.MoveTo(0, 4)
	startTime = time.Now()
}

func setProgress(percent float64) {
	ansi.Save()
	ansi.MoveTo(1, 1)
	ansi.ClearLine()

	count := 40.0
	fmt.Print("[")
	for i := 0.0; i < count; i++ {
		if i/count >= percent {
			fmt.Print(" ")
		} else {
			ansi.Print(ansi.BoldYellow, "#")
		}
	}
	fmt.Print("]")
	fmt.Printf(" %0.2f%%", percent*100)

	endTime := time.Now()
	seconds := int(endTime.Sub(startTime).Seconds() / percent * (1 - percent))
	minutes := seconds / 60
	seconds = seconds % 60

	ansi.MoveTo(1, 2)
	ansi.ClearLine()
	if percent > 0 {
		fmt.Printf("Estimated %d:%02d left", minutes, seconds)
	}
	ansi.Restore()
}

func setComplete() {
	endTime := time.Now()
	seconds := int(endTime.Sub(startTime).Seconds())
	minutes := seconds / 60
	seconds = seconds % 60
	ansi.ClearScreen()
	ansi.ResetScrollRegion()
	fmt.Println("Frames render complete!")
	fmt.Printf("Elapsed time: %d:%02d.\n\n\n", minutes, seconds)
}
