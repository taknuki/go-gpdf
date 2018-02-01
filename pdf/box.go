package pdf

import "fmt"

// Box is a rectangle box that represents a region.
type Box struct {
	leftBottomX int
	leftBottomY int
	rightTopX   int
	rightTopY   int
}

// NewBox creates the Box.
func NewBox(leftBottomX, leftBottomY, rightTopX, rightTopY int) *Box {
	return &Box{leftBottomX, leftBottomY, rightTopX, rightTopY}
}

// NewBoxA4 creates A4 size Box
func NewBoxA4() *Box {
	return NewBox(0, 0, 595, 842)
}

// NewBoxA3 creates A3 size Box
func NewBoxA3() *Box {
	return NewBox(0, 0, 842, 1191)
}

// NewBoxA2 creates A2 size Box
func NewBoxA2() *Box {
	return NewBox(0, 0, 1191, 1684)
}

// NewBoxA1 creates A1 size Box
func NewBoxA1() *Box {
	return NewBox(0, 0, 1684, 2384)
}

// NewBoxA0 creates A0 size Box
func NewBoxA0() *Box {
	return NewBox(0, 0, 2384, 3370)
}

// NewInnerBox creates the inner Box.
func (r *Box) NewInnerBox(topMargin, rightMargin, bottomMargin, leftMargin int) *Box {
	return NewBox(r.leftBottomX+leftMargin, r.leftBottomY+bottomMargin, r.rightTopX-rightMargin, r.rightTopY-topMargin)
}

func (r *Box) compile() string {
	return fmt.Sprintf("[%d %d %d %d]", r.leftBottomX, r.leftBottomY, r.rightTopX, r.rightTopY)
}
