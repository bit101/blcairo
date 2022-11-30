// Package logdisplay is a log/log pixel buffer for better contrast in some cases.
package logdisplay

import (
	"math"
)

// LogDisplay represents a bit array. Pixel values are incremented or set directly.
// Max value is kept track of
// If trackMins is true, than min value is kept track of. Otherwise it's 0.0.
// The value of each pixel is then log10(value - min) / log10(max - min).
// This gives more realistic contrast in some use cases.
type LogDisplay struct {
	width, height int
	values        []float64
	min           float64
	max           float64
	trackMins     bool
}

// NewLogDisplay creates a new LogDisplay struct.
func NewLogDisplay(width, height int, trackMins bool) *LogDisplay {
	min := 0.0
	if trackMins {
		min = math.MaxFloat64
	}
	return &LogDisplay{
		width:     width,
		height:    height,
		min:       min,
		max:       0.0,
		trackMins: trackMins,
		values:    make([]float64, width*height),
	}
}

// Inc increments the raw value of a pixel by 1.
// Mostly used with trackMins as false and setting only some pixels.
func (d *LogDisplay) Inc(x, y float64) {
	xx, yy := int(x), int(y)
	if xx >= 0 && xx < d.width && yy >= 0 && yy < d.height {
		index := xx + yy*d.width
		value := d.values[index] + 1
		d.values[index] = value
		if value > d.max {
			d.max = value
		}
		if d.trackMins && value < d.min {
			d.min = value
		}
	}
}

// Set sets the value for a given pixel.
// Mostly used with trackMins set to true and setting every single pixel directly.
func (d *LogDisplay) Set(value, x, y float64) {
	xx, yy := int(x), int(y)
	if xx >= 0 && xx < d.width && yy >= 0 && yy < d.height {
		index := xx + yy*d.width
		d.values[index] = value
		if value > d.max {
			d.max = value
		}
		if d.trackMins && value < d.min {
			d.min = value
		}
	}
}

// GetLog calculates the logarithmic value of the pixel.
// This should be a value from 0.0 to 1.0
// which can be mapped to a gray scale value or a range of colors.
func (d *LogDisplay) GetLog(x, y int) float64 {
	value := d.values[x+y*d.width]
	if value == d.min {
		return 0
	}
	if d.trackMins && value == 0 {
		return 0
	}
	return math.Log(value-d.min) / math.Log(d.max-d.min)
}

// GetExp is an alternate algorithm.
// Returns the base value raised to the power of the count of this pixel.
// Base should be 0.95 - 0.99... ish. Higher = more contrast.
// Tends to make higher contrast images than the log method.
func (d *LogDisplay) GetExp(x, y int, base float64) float64 {
	value := d.values[x+y*d.width]
	return 1.0 - math.Pow(base, value)
}

// GetCount returns the count for the current pixel
func (d *LogDisplay) GetCount(x, y int) float64 {
	return d.values[x+y*d.width]
}
