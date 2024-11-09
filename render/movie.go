// Package render renders a single image or a number of frames
package render

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

//////////////////////////////
// MOVIE
//////////////////////////////

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
}

// NewMovie creates a new movie.
func NewMovie(name string, width, height, fps int) *Movie {
	return &Movie{
		name,
		width,
		height,
		fps,
		map[string]*Act{}, // for accessing acts by name.
		[]*Act{},          // for accessing acts by index or sequentially.
	}
}

// NewAct adds an act to this movie.
func (m *Movie) NewAct(name string, frameCount int, renderFunc FrameFunc, render bool) {
	act := newAct(m, name, frameCount, renderFunc)
	m.Acts[name] = act
	m.List = append(m.List, act)
	if render {
		act.render()
	}
}

// PlayAct plays a given act if the video exists.
func (m *Movie) PlayAct(act string) {
	m.Acts[act].play()
}

// PlayIndex plays a given act if the video exists.
func (m *Movie) PlayIndex(actIndex int) {
	m.List[actIndex].play()
}

// PlayAll plays all acts in order.
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
	os.Remove("out/" + m.Name + ".manifest")
	os.Remove("out/" + m.Name + ".mp4")
}

// CombineAll combines all rendered act videos into a single movie, and optionally plays that movie.
// This relies on all acts having been rendered. It will fail if there is a missing act video.
func (m *Movie) CombineAll(play bool) {
	m.WriteManifest()
	cmd := exec.Command(
		"ffmpeg", "-y",
		"-f", "concat",
		"-i", "out/"+m.Name+".manifest",
		"-c", "copy",
		"out/"+m.Name+".mp4",
	)
	cmd.Run()
	if play {
		m.PlayCombined()
	}
}

// PlayCombined plays combined movie if it exists.
func (m *Movie) PlayCombined() {
	fileName := "out/" + m.Name + ".mp4"
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		fmt.Println(err)
	} else {
		PlayVideo(fileName)
	}
}

// WriteManifest writes a file manifest used by ffmpeg to concatenate all files into one.
func (m *Movie) WriteManifest() {
	output := ""
	for _, act := range m.List {
		output += "file " + act.Name + ".mp4\n"
	}
	os.WriteFile("out/"+m.Name+".manifest", []byte(output), 0777)
}
