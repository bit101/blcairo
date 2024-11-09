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
	RenderFunc FrameFunc
}

// newAct creates a new movie.
func newAct(parent *Movie, name string, frameCount int, renderFunc FrameFunc) *Act {
	return &Act{parent, name, frameCount, renderFunc}
}

// render renders the act.
func (a *Act) render() {
	frames := "out/" + a.Name + "_frames"
	fileName := "out/" + a.Name + ".mp4"
	Frames(float64(a.Parent.Width), float64(a.Parent.Height), a.FrameCount, frames, a.RenderFunc)
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
	fileName := "out/" + a.Name + ".png"
	percent := float64(frame) / float64(a.FrameCount)
	Image(float64(a.Parent.Width), float64(a.Parent.Height), fileName, a.RenderFunc, percent)

}

// play plays the video if it exists.
func (a *Act) play() {
	fileName := "out/" + a.Name + ".mp4"
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		fmt.Println(err)
	} else {
		PlayVideo(fileName)
	}
}

// clean deletes the video and frames.
func (a *Act) clean() {
	fileName := "out/" + a.Name + ".mp4"
	os.Remove(fileName)

	frames := "out/" + a.Name + "_frames"
	os.RemoveAll(frames)

}
