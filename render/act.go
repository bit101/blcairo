// Package render renders a single image or a number of frames
package render

import (
	"errors"
	"fmt"
	"os"
)

// Act is a single act in a movie.
type Act struct {
	Parent     *Movie
	Name       string
	FrameCount int
	Out        string
	RenderFunc FrameFunc
}

// newAct creates a new movie.
func newAct(parent *Movie, name string, frameCount int, out string, renderFunc FrameFunc) *Act {
	return &Act{parent, name, frameCount, out, renderFunc}
}

// render renders the act.
func (a *Act) render() {
	frames := a.Out + a.Name + "_frames"
	fileName := a.Out + a.Name + ".mp4"
	Frames("act: "+a.Name, float64(a.Parent.Width), float64(a.Parent.Height), a.FrameCount, frames, a.RenderFunc)
	ConvertToVideo(
		frames,
		fileName,
		float64(a.Parent.Width),
		float64(a.Parent.Height),
		a.Parent.FPS,
		a.FrameCount/a.Parent.FPS,
		false,
	)
}

// renderFrame renders a single frame of the act.
func (a *Act) renderFrame(frame int) {
	fileName := a.Out + a.Name + ".png"
	percent := float64(frame) / float64(a.FrameCount)
	Image(float64(a.Parent.Width), float64(a.Parent.Height), fileName, a.RenderFunc, percent)

}

// play plays the video if it exists.
func (a *Act) play() {
	fileName := a.Out + a.Name + ".mp4"
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		fmt.Println(err)
	} else {
		PlayVideo(fileName)
	}
}

// clean deletes the video and frames.
func (a *Act) clean() {
	fileName := a.Out + a.Name + ".mp4"
	os.Remove(fileName)

	frames := a.Out + a.Name + "_frames"
	os.RemoveAll(frames)

}
