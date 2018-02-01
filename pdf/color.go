package pdf

import "fmt"

// Color is a color container.
type Color interface {
	colorSpace() colorSpace
	strokeColor() string
	nonStrokeColor() string
}

type grayScaleColor struct {
	scale float32
}

// NewColorGrayScale returns a gray scale color
func NewColorGrayScale(scale float32) Color {
	return &grayScaleColor{scale}
}

func (c *grayScaleColor) colorSpace() colorSpace {
	return colorSpaceDeviceGray
}

func (c *grayScaleColor) strokeColor() string {
	return fmt.Sprintf("%f G", c.scale)
}

func (c *grayScaleColor) nonStrokeColor() string {
	return fmt.Sprintf("%f g", c.scale)
}

type rgbColor struct {
	red   float32
	green float32
	blue  float32
}

// NewColorRGB returns a rgb color
func NewColorRGB(red, green, blue float32) Color {
	return &rgbColor{red, green, blue}
}

func (c *rgbColor) colorSpace() colorSpace {
	return colorSpaceDeviceRGB
}

func (c *rgbColor) strokeColor() string {
	return fmt.Sprintf("%f %f %f RG", c.red, c.green, c.blue)
}

func (c *rgbColor) nonStrokeColor() string {
	return fmt.Sprintf("%f %f %f rg", c.red, c.green, c.blue)
}

type cmykColor struct {
	cyan    float32
	magenta float32
	yellow  float32
	key     float32
}

// NewColorCMYK returns a cmyk color
func NewColorCMYK(cyan, magenta, yellow, key float32) Color {
	return &cmykColor{cyan, magenta, yellow, key}
}

func (c *cmykColor) colorSpace() colorSpace {
	return colorSpaceDeviceCMYK
}

func (c *cmykColor) strokeColor() string {
	return fmt.Sprintf("%f %f %f %f K", c.cyan, c.magenta, c.yellow, c.key)
}

func (c *cmykColor) nonStrokeColor() string {
	return fmt.Sprintf("%f %f %f %f k", c.cyan, c.magenta, c.yellow, c.key)
}

type undefColor struct{}

// newColorUndef returns a dummy Color for no filling or no stroking.
func newColorUndef() Color {
	return &undefColor{}
}

func (c *undefColor) colorSpace() colorSpace {
	return colorSpaceUndefined
}

func (c *undefColor) strokeColor() string {
	return ""
}

func (c *undefColor) nonStrokeColor() string {
	return ""
}

// colorSpace is color space
type colorSpace int

const (
	// colorSpaceUndefined is dummy color space for no filling or no stroking.
	colorSpaceUndefined colorSpace = iota
	// colorSpaceDeviceGray is gray scale.
	colorSpaceDeviceGray
	// colorSpaceDeviceRGB is rgb color.
	colorSpaceDeviceRGB
	// colorSpaceDeviceCMYK is cmyk color.
	colorSpaceDeviceCMYK
)

func (cs colorSpace) String() string {
	switch cs {
	case colorSpaceDeviceGray:
		return "/DeviceGray"
	case colorSpaceDeviceRGB:
		return "/DeviceRGB"
	case colorSpaceDeviceCMYK:
		return "/DeviceCMYK"
	default:
		return "/Undefined"
	}
}
