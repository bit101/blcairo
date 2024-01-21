// Package cairo wraps the c cairographics library.
package cairo

// Filter cairo_filter_t
type Filter int

// Filter constants
const (
	FilterFast Filter = iota
	FilterGood
	FilterBest
	FilterNearest
	FilterBilinear
)
