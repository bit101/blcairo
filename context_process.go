// Package cairo wraps the c cairographics library.
package cairo

// Grayscale turns the image grayscale.
func (c *Context) Grayscale() {
	data, _ := c.Surface.GetData()
	for i := 0; i < len(data); i += 4 {
		val := byte(0.2126*float64(data[i]) + 0.7152*float64(data[i+1]) + 0.0722*float64(data[i+2]))
		data[i] = val
		data[i+1] = val
		data[i+2] = val
	}
	c.Surface.SetData(data)
}
