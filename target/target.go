// Package target defines render targets
package target

const (
	// None will render nothing.
	None = iota
	// Image will render an image.
	Image
	// Gif will render an animated gif.
	Gif
	// Video will render a video.
	Video
	// SpriteSheet will render a sprite sheet.
	SpriteSheet
	// Montage will render a sprite sheet.
	Montage
)
