// Package render renders a single image or a number of frames
package render

import (
	"fmt"
	"os"

	cairo "github.com/bit101/blcairo"
)

// Program contains a collection of scenes that are rendered into a single video.
// A better paradigm is the Movie/Act system where each act is rendered separately as needed
// and then combined to a single movie with ffmpeg.
type Program struct {
	Width, Height float64
	FPS           int
	Scenes        []*Scene
}

// Scene represents a period of time where one rendering function and set of parameters are active.
type Scene struct {
	FrameCount int
	FrameFunc  FrameFunc
}

// NewProgram creates a new Program
func NewProgram(width, height float64, fps int) *Program {
	program := &Program{
		Width:  width,
		Height: height,
		FPS:    fps,
		Scenes: []*Scene{},
	}
	return program
}

// AddSceneWithFrames adds a scene to this Program with a duration in frames
func (p *Program) AddSceneWithFrames(frameFunc FrameFunc, frameCount int) {
	s := &Scene{
		FrameFunc:  frameFunc,
		FrameCount: frameCount,
	}
	p.Scenes = append(p.Scenes, s)
}

// AddSceneWithSeconds adds a scene to this Program with a duration in seconds
func (p *Program) AddSceneWithSeconds(frameFunc FrameFunc, seconds float64) {
	frameCount := int(float64(p.FPS) * seconds)
	p.AddSceneWithFrames(frameFunc, frameCount)
}

// TotalFrames returns the total number of frames in the Program.
func (p *Program) TotalFrames() int {
	total := 0
	for _, s := range p.Scenes {
		total += s.FrameCount
	}
	return total
}

// Seconds returns the time in seconds of the Program
func (p *Program) Seconds() int {
	return p.TotalFrames() / p.FPS
}

// Render coordinates the rendering of all scenes in this Program.
func (p *Program) Render(frames string) {
	initProgress()
	os.RemoveAll(frames)
	os.MkdirAll(frames, 0755)

	surface := cairo.NewSurface(int(p.Width), int(p.Height))
	context := cairo.NewContext(surface)
	currentFrame := 0
	totalFrames := float64(p.TotalFrames())

	for _, scene := range p.Scenes {
		for f := 0; f < scene.FrameCount; f++ {
			percent := float64(f) / float64(scene.FrameCount)

			totalPercent := float64(currentFrame) / totalFrames
			sceneName := fmt.Sprintf("scene %d", f)
			setProgress(sceneName, currentFrame, p.TotalFrames(), totalPercent)
			scene.FrameFunc(context, p.Width, p.Height, percent)
			surface.WriteToPNG(fmt.Sprintf("%s/frame_%04d.png", frames, currentFrame))
			currentFrame++
		}
	}
	setComplete()
}

// RenderAndPlayVideo renders the program to a video file using the given frames directory and output filename
func (p *Program) RenderAndPlayVideo(frames, fileName string) {
	p.RenderVideo(frames, fileName)
	PlayVideo(fileName)
}

// RenderVideo renders the program to a video file using the given frames directory and output filename
func (p *Program) RenderVideo(frames, fileName string) {
	p.Render(frames)
	ConvertToVideo(frames, fileName, p.Width, p.Height, p.FPS, p.Seconds(), true)
}

// RenderGif renders the program to a video file using the given frames directory and output filename
func (p *Program) RenderGif(frames, fileName string) {
	p.Render(frames)
	FfmpegToGIF(frames, fileName, p.FPS)
}
