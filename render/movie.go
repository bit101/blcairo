// Package render renders a single image or a number of frames
package render

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/bit101/go-ansi"
)

//////////////////////////////
// MOVIE
//////////////////////////////

const (
	Aspect16x9 = 16.0 / 9.0
	Aspect4x3  = 4.0 / 3.0
	Aspect1x1  = 1.0 / 1.0
)

// Movie represents a single movie containing multiple acts.
// This differs from a Program/Scene setup in that each act of a movie is generated as a
// separate named mp4 file, which would have to be concatenated separately.
type Movie struct {
	Name   string
	Width  int
	Height int
	FPS    int
	Acts   map[string]*Act
	List   []*Act
	Out    string
}

// NewMovie creates a new movie.
func NewMovie(name string, width, height float64, fps int) *Movie {
	return &Movie{
		name,
		int(width),
		int(height),
		fps,
		map[string]*Act{}, // for accessing acts by name.
		[]*Act{},          // for accessing acts by index or sequentially.
		fmt.Sprintf("out_%d/", int(height)),
	}
}

// NewAct adds an act to this movie.
func (m *Movie) NewAct(name string, frameCount int, renderFunc FrameFunc, render bool, play bool) {
	act := newAct(m, name, frameCount, m.Out, renderFunc)
	m.Acts[name] = act
	m.List = append(m.List, act)
	if render {
		act.render()
	}
	if play {
		act.play()
	}
}

// ReuseAct adds an existing act to a different location in this movie.
// This will re-use the already-rendered video for this act, instead of re-rendering it.
func (m *Movie) ReuseAct(name string) {
	act := m.Acts[name]
	m.List = append(m.List, act)
}

// PlayAct plays the named act if the video exists.
func (m *Movie) PlayAct(act string) {
	m.Acts[act].play()
}

// PlayIndex plays the indexed act if the video exists.
func (m *Movie) PlayIndex(actIndex int) {
	m.List[actIndex].play()
}

// PlayAll plays all acts in order of how they were added.
func (m *Movie) PlayAll() {
	for _, act := range m.List {
		act.play()
	}
}

// Clean cleans all the files created by this movie.
func (m *Movie) Clean() {
	for _, act := range m.Acts {
		act.clean()
	}
	os.Remove(m.Out + m.Name + ".manifest")
	os.Remove(m.Out + m.Name + ".mp4")
}

// CombineAll combines all rendered act videos into a single movie, and optionally plays that movie.
// This method will create a movie as long as at least one of the act videos exists,
// but will display which act videos are missing.
func (m *Movie) CombineAll(play bool) {
	ansi.ClearScreen()
	m.WriteManifest()
	cmd := exec.Command(
		"ffmpeg", "-y",
		"-f", "concat",
		"-i", m.Out+m.Name+".manifest",
		"-c", "copy",
		m.Out+m.Name+".mp4",
	)
	cmd.Run()
	fmt.Printf("Movie %q complete.\n", m.Name)
	if play {
		m.PlayCombined()
	}
}

// PlayCombined plays combined movie if it exists.
func (m *Movie) PlayCombined() {
	fileName := m.Out + m.Name + ".mp4"
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		ansi.Println(ansi.Red, err)
	} else {
		PlayVideo(fileName)
	}
}

// WriteManifest writes a file manifest used by ffmpeg to concatenate all files into one.
func (m *Movie) WriteManifest() {
	output := ""
	for _, act := range m.List {
		fileName := act.Name + ".mp4"
		output += "file " + fileName + "\n"
		if _, err := os.Stat(m.Out + fileName); errors.Is(err, os.ErrNotExist) {
			ansi.Println(ansi.Red, err)
		}
	}
	os.WriteFile(m.Out+m.Name+".manifest", []byte(output), 0777)
}
