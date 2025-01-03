// Package render renders a single image or a number of frames
package render

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/bit101/bitlib/blcolor"
	cairo "github.com/bit101/blcairo"
	"github.com/bit101/go-ansi"
)

var startTime time.Time

// FrameFunc is the interface for a function that renders a single frame.
type FrameFunc func(*cairo.Context, float64, float64, float64)

// CreateAndViewImage sets up the rendering of a single image.
func CreateAndViewImage(width, height float64, path string, frameFunc FrameFunc, percent float64) {
	Image(width, height, path, frameFunc, percent)
	ViewImage(path)
}

// Image sets up the rendering of a single image.
func Image(width, height float64, path string, frameFunc FrameFunc, percent float64) {
	fmt.Println("Generating image...")
	surface := cairo.NewSurface(int(width), int(height))
	context := cairo.NewContext(surface)
	frameFunc(context, width, height, percent)
	checkOutDir(path)
	surface.WriteToPNG(path)
	fmt.Println("Image complete!")
	data, _ := os.Stat(path)
	fmt.Println("File:", path)
	fmt.Printf("Resolution: %dx%d\n", int(width), int(height))
	fmt.Printf("Size: %0.2f kb\n", float64(data.Size())/1000)
}

func checkOutDir(path string) {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			os.MkdirAll(dir, 0755)
		}
	}
}

// Frames sets up the renderin of a series of frames.
func Frames(renderName string, width, height float64, numFrames int, frames string, frameFunc FrameFunc) {
	initProgress()
	os.RemoveAll(frames)
	os.MkdirAll(frames, 0755)
	surface := cairo.NewSurface(int(width), int(height))
	context := cairo.NewContext(surface)
	for frame := 0; frame < numFrames; frame++ {
		percent := float64(frame) / float64(numFrames)
		setProgress(renderName, frame, numFrames, percent)
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
		fr := fmt.Sprintf("range: %d-%d", start, end)
		setProgress(fr, frame, numFrames, percent)
		frameFunc(context, width, height, percent)
		surface.WriteToPNG(fmt.Sprintf("%s/frame_%04d.png", frames, frame))
	}
	setComplete()
}

// SpriteSheet sets up the rendering of a sprite sheet.
func SpriteSheet(width, height float64, bg blcolor.Color, path string, numFrames int, frameFunc FrameFunc) {
	// todo: create an intermediate surface of sizw width*height. Pass its context to the frameFunc.
	// Then draw it to the correct location in the main surface.
	// This will allow for functions that affect the whole context to work in frameFunc, such as:
	// clearing functions, translateCenter
	// as well as limiting the scope of functions that draw outside bounds.
	initProgress()
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
		setProgress("sprite sheet", int(i), numFrames, percent)
		frameFunc(context, width, height, percent)
		context.Restore()

		x += width
		if x >= size*width {
			x = 0.0
			y += height
		}
	}
	surface.WriteToPNG(path)
	setComplete()
}

func initProgress() {
	ansi.ClearScreen()
	ansi.SetScrollRegion(4, 1000)
	ansi.MoveTo(0, 4)
	startTime = time.Now()
}

func setProgress(renderName string, frame, total int, percent float64) {
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
	fmt.Println("]")
	fmt.Println(renderName)
	fmt.Printf("Frame %d of %d (%0.1f%%)\n", frame, total, percent*100)

	endTime := time.Now()
	seconds := int(endTime.Sub(startTime).Seconds() / percent * (1 - percent))
	minutes := seconds / 60
	seconds = seconds % 60

	// ansi.MoveTo(1, 3)
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
	ansi.MoveTo(1, 1)
	ansi.ClearLine()
	ansi.MoveTo(1, 2)
	ansi.ClearLine()
	ansi.ResetScrollRegion()
	fmt.Println("Frames render complete!")
	fmt.Printf("Elapsed time: %d:%02d.\n\n", minutes, seconds)
}
