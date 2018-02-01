package pdf

import (
	"fmt"
	"math"
	"strconv"
)

// ImageResource is a image for a pdf resource.
//
// UNDER IMPLEMENTATION
type ImageResource struct {
	name             string
	width            int
	height           int
	colorSpace       colorSpace
	bitsPerComponent int
	data             []byte
}

// newImageResource returns a ImageResource.
func newImageResource(name string, width, height int, bitperComponent int, data []byte) *ImageResource {
	return &ImageResource{
		name:             name,
		width:            width,
		height:           height,
		colorSpace:       colorSpaceDeviceGray,
		bitsPerComponent: bitperComponent,
		data:             data,
	}
}

func (i *ImageResource) asStream() *stream {
	s := newFlatStream()
	s.dict["/Type"] = "/XObject"
	s.dict["/Subtype"] = "/Image"
	s.dict["/Width"] = strconv.Itoa(i.width)
	s.dict["/Height"] = strconv.Itoa(i.height)
	s.dict["/BitsPerComponent"] = strconv.Itoa(i.bitsPerComponent)
	s.dict["/ColorSpace"] = i.colorSpace.String()
	s.addBinaryDatum(i.data)
	s.addStringDatum("\n")
	return s
}

// Image is the operator for drawing a image resource.
//
// see the comment of GraphicsObject.
type Image interface {
	Width(pt float64) Image
	Height(pt float64) Image
	Rotate(radian float64) Image
	GraphicsObject
}

type image struct {
	graphicsObject
	name    string
	centerX float64
	centerY float64
	width   float64
	height  float64
	rotate  float64
}

// newImage returns the new image.
func newImage(page Page, ir *ImageResource, centerX, centerY float64) *image {
	return &image{
		graphicsObject: graphicsObject{
			page: page,
		},
		name:    ir.name,
		centerX: centerX,
		centerY: centerY,
		width:   float64(ir.width),
		height:  float64(ir.height),
	}
}

func (i *image) Width(width float64) Image {
	i.ifNotRendered(func() {
		i.width = width
	})
	return i
}
func (i *image) Height(height float64) Image {
	i.ifNotRendered(func() {
		i.height = height
	})
	return i
}
func (i *image) Rotate(rotate float64) Image {
	i.ifNotRendered(func() {
		i.rotate = rotate
	})
	return i
}

func (i *image) Render() {
	i.renderSelf(i)
}

func (i *image) render(cb *Box) string {
	cos := math.Cos(i.rotate)
	sin := math.Sin(i.rotate)
	a11 := i.width * cos
	a12 := i.width * sin
	a21 := -i.height * sin
	a22 := i.height * cos
	c1 := 0.5*i.width*cos - 0.5*i.height*sin
	c2 := 0.5*i.width*sin + 0.5*i.height*cos
	b1 := -c1 + float64(cb.leftBottomX) + i.centerX
	b2 := -c2 + float64(cb.rightTopY) - i.centerY
	return fmt.Sprintf(
		"q %f %f %f %f %f %f cm %s Do Q\n",
		a11, a12, a21, a22, b1, b2, i.name)
}
