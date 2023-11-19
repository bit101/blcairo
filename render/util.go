// Package render has functions that render gifs, images, videos, etc.
package render

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var imageRenderType = "png"

func UseBMP(value bool) {
	if value {
		imageRenderType = "bmp"
	} else {
		imageRenderType = "png"
	}
}

// MakeGIF creates an animated gif with the given tool.
func MakeGIF(tool, folder, outFileName string, fps float64) {
	os.RemoveAll(outFileName)
	if tool == "convert" {
		ConvertToGIF(folder, outFileName, fps)
	} else if tool == "ffmpeg" {
		FfmpegToGIF(folder, outFileName, fps)
	}

}

// ConvertToGIF converts a folder of pngs into an animated gif using imagemagick convert.
func ConvertToGIF(folder, outFileName string, fps float64) {
	delay := fmt.Sprintf("%f", 1000.0/fps/10.0)
	path := folder + "/*." + imageRenderType
	cmd := exec.Command("convert", "-delay", delay, "-layers", "Optimize", path, outFileName)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// FfmpegToGIF converts a folder of pngs into an animated gif using ffmpeg.
func FfmpegToGIF(folder, outFileName string, fps float64) {
	path := folder + "/frame_%04d." + imageRenderType
	fpsArg := fmt.Sprintf("%d", int(fps))

	paletteCmd := exec.Command("ffmpeg", "-y", "-i", path, "-vf", "palettegen", "palette.png")
	err := paletteCmd.Run()
	if err != nil {
		log.Fatalf("Could not create palette: %s", err)
	}

	outCmd := exec.Command("ffmpeg", "-y", "-framerate", fpsArg, "-i", path, "-i", "palette.png", "-filter_complex", "paletteuse", outFileName)
	err = outCmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// ConvertToVideo converts a folder of pngs into an mp4 video file. Requires ffmpeg.
func ConvertToVideo(folder, outFileName string, w, h, fps int) {
	os.RemoveAll(outFileName)
	path := folder + "/frame_%04d." + imageRenderType
	fpsArg := fmt.Sprintf("%d", fps)
	sizeArg := fmt.Sprintf("%dx%d", w, h)

	cmd := exec.Command("ffmpeg", "-framerate", fpsArg, "-i", path, "-s:v", sizeArg,
		"-c:v", "libx264", "-profile:v", "high", "-crf", "20",
		"-pix_fmt", "yuv420p", outFileName)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// MixAV mixes an audio and video file.
func MixAV(videoFileName, audioFileName, outFileName string) {
	cmd := exec.Command(
		"ffmpeg", "-y",
		"-i", videoFileName,
		"-i", audioFileName,
		"-c", "copy",
		"-map", "0:v:0",
		"-map", "1:a:0",
		outFileName,
	)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// ViewImage displays an image using installed image viewer.
func ViewImage(imagePath string) {
	cmd := exec.Command("bitlibImageViewer", imagePath)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// ViewGif plays an animated gif using installed gif viewer
func ViewGif(imagePath string) {
	cmd := exec.Command("bitlibGifViewer", imagePath)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// VLC launches vlc to play a video
func VLC(fileName string, loop bool) {
	loopArg := ""
	if loop {
		loopArg = "--loop"
	}
	cmd := exec.Command("vlc", loopArg, fileName)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// PlayVideo launches an app to play a video
func PlayVideo(fileName string) {
	cmd := exec.Command("bitlibVideoPlayer", fileName)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// ParentDir returns the immediated directory name of the current working directory.
func ParentDir() string {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot get directory.")
		os.Exit(1)
	}
	return filepath.Base(wd)
}
