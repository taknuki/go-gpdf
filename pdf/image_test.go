package pdf

import (
	"fmt"
	"math"
	"testing"
)

type mockImage struct {
	mockGraphicsObject
}

func (i *mockImage) Width(width float64) Image   { return i }
func (i *mockImage) Height(height float64) Image { return i }
func (i *mockImage) Rotate(roate float64) Image  { return i }

func TestImage1(t *testing.T) {
	p := &mockPage{}
	cb := NewBox(10, 20, 500, 600)
	ir := newImageResource("/I1", 30, 40, 1, []byte{1})
	i := newImage(p, ir, 100, 200)
	isEqualGraphicsObject(t, p, false, i.graphicsObject)
	checkRender(t, p, i)
	expected := fmt.Sprintf("q 30.000000 0.000000 -0.000000 40.000000 %f %f cm /I1 Do Q\n", float64(100+10-30/2), float64(600-200-40/2))
	actual := i.render(cb)
	testRendering(t, expected, actual)
}

func TestImage2(t *testing.T) {
	p := &mockPage{}
	cb := NewBox(10, 20, 500, 600)
	ir := newImageResource("/I1", 30, 40, 1, []byte{1})
	i := newImage(p, ir, 100, 200).Width(50).Height(60).Rotate(math.Pi / 2)
	expected := fmt.Sprintf("q 0.000000 50.000000 -60.000000 0.000000 %f %f cm /I1 Do Q\n", float64(100+10+60/2), float64(600-200-50/2))
	actual := i.render(cb)
	testRendering(t, expected, actual)
}
