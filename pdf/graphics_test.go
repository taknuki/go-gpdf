package pdf

import (
	"fmt"
	"testing"
)

type mockGraphicsObject struct {
	graphicsObject
}

func (obj *mockGraphicsObject) Render()               {}
func (obj *mockGraphicsObject) render(cb *Box) string { return "" }

type mockRectangle struct {
	mockGraphicsObject
}

func (r *mockRectangle) StrokeColor(c Color) Rectangle { return r }
func (r *mockRectangle) FillColor(c Color) Rectangle   { return r }

type mockLine struct {
	mockGraphicsObject
}

func (l *mockLine) MoveTo(x, y int) Line                  { return l }
func (l *mockLine) Color(c Color) Line                    { return l }
func (l *mockLine) DashPattern(dash, gap, phase int) Line { return l }
func (l *mockLine) CapStyle(lcs LineCapStyle) Line        { return l }
func (l *mockLine) JoinStyle(ljs LineJoinStyle) Line      { return l }
func (l *mockLine) StrokeColor(c Color) Rectangle         { return l }
func (l *mockLine) FillColor(c Color) Rectangle           { return l }

func isEqualGraphicsObject(t *testing.T, page Page, rendered bool, actual graphicsObject) {
	t.Helper()
	if actual.page != page {
		t.Errorf("graphicsObject page: expected:%p actual:%p", page, actual.page)
	}
	if actual.rendered != rendered {
		t.Errorf("graphicsObject rendered: expected:%v actual:%v", rendered, actual.rendered)
	}
}

func testRendering(t *testing.T, expected, actual string) {
	t.Helper()
	if expected != actual {
		t.Errorf("Rendering failed\nexpected:%s\nactual  :%s\n", expected, actual)
	}
}

func checkRender(t *testing.T, p *mockPage, obj GraphicsObject) {
	t.Helper()
	obj.Render()
	if p.renderResult != obj {
		t.Errorf("Render: unexpected object is rendered: expected:%p actual: %p", obj, p.renderResult)
	}
}

func TestGraphicsObject(t *testing.T) {
	p := &mockPage{}
	obj := &graphicsObject{page: p}
	isEqualGraphicsObject(t, p, false, *obj)
	called1 := false
	obj.ifNotRendered(func() {
		called1 = true
	})
	if !called1 {
		t.Error("ifNotRendered: callback is not called in spite of the not-rendered object")
	}
	self := &mockGraphicsObject{}
	obj.renderSelf(self)
	if p.renderResult != self {
		t.Errorf("renderSelf: unexpected object is rendered: expected:%p actual: %p", self, p.renderResult)
	}
	if !obj.rendered {
		t.Error("renderSelf: rendered must be true after renderSelf is called")
	}
	called2 := false
	obj.ifNotRendered(func() {
		called2 = true
	})
	if called2 {
		t.Error("ifNotRendered: callback is called in spite of the rendered object")
	}
}

func TestLine1(t *testing.T) {
	p := &mockPage{}
	cb := NewBox(10, 20, 500, 600)
	l := newLine(p, 100, 200, 300, 400, 5)
	isEqualGraphicsObject(t, p, false, l.graphicsObject)
	checkRender(t, p, l)
	expected := "q 110 400 m 310 200 l 5 w  S Q\n"
	actual := l.render(cb)
	testRendering(t, expected, actual)
}

func TestLine2(t *testing.T) {
	p := &mockPage{}
	cb := NewBox(10, 20, 500, 600)
	color := NewColorRGB(0.1, 0.2, 0.3)
	dash := newLineDashPattern(1, 2, 3)
	l := newLine(p, 100, 200, 300, 400, 5).
		MoveTo(350, 450).
		Color(color).
		DashPattern(1, 2, 3).
		CapStyle(LineCapStyleProjectingSquare).
		JoinStyle(LineJoinStyleBevel)
	expected := fmt.Sprintf("q 110 400 m 310 200 l 360 150 l 5 w %s %s %s %s S Q\n", color.strokeColor(), dash.compile(), LineCapStyleProjectingSquare.compile(), LineJoinStyleBevel.compile())
	actual := l.render(cb)
	testRendering(t, expected, actual)
}

func TestLineDashPattern(t *testing.T) {
	expected1 := "[1 2] 3 d"
	actual1 := newLineDashPattern(1, 2, 3).compile()
	testCompillation(t, expected1, actual1)
	expected2 := "[1] 3 d"
	actual2 := newLineDashPattern(1, 1, 3).compile()
	testCompillation(t, expected2, actual2)
}

func TestLineCapStyle(t *testing.T) {
	testCompillation(t, "0 J", LineCapStyleButt.compile())
	testCompillation(t, "1 J", LineCapStyleRound.compile())
	testCompillation(t, "2 J", LineCapStyleProjectingSquare.compile())
	testCompillation(t, "", lineCapStyleUndefined.compile())
}

func TestLineJoinStyle(t *testing.T) {
	testCompillation(t, "0 j", LineJoinStyleMiter.compile())
	testCompillation(t, "1 j", LineJoinStyleRound.compile())
	testCompillation(t, "2 j", LineJoinStyleBevel.compile())
	testCompillation(t, "", lineJoinStyleUndefined.compile())
}

func TestRectangle1(t *testing.T) {
	p := &mockPage{}
	cb := NewBox(10, 20, 500, 600)
	r := newRectangle(p, 10, 20, 100, 200)
	isEqualGraphicsObject(t, p, false, r.graphicsObject)
	checkRender(t, p, r)
	expected := "q  20 380 100 200 re n Q\n"
	actual := r.render(cb)
	testRendering(t, expected, actual)
}

func TestRectangle2(t *testing.T) {
	p := &mockPage{}
	cb := NewBox(10, 20, 500, 600)
	sc := NewColorRGB(0.1, 0.2, 0.3)
	r := newRectangle(p, 10, 20, 100, 200).StrokeColor(sc)
	expected := fmt.Sprintf("q %s 20 380 100 200 re S Q\n", sc.strokeColor())
	actual := r.render(cb)
	testRendering(t, expected, actual)
}

func TestRectangle3(t *testing.T) {
	p := &mockPage{}
	cb := NewBox(10, 20, 500, 600)
	fc := NewColorRGB(0.4, 0.5, 0.6)
	r := newRectangle(p, 10, 20, 100, 200).FillColor(fc)
	expected := fmt.Sprintf("q %s 20 380 100 200 re f Q\n", fc.nonStrokeColor())
	actual := r.render(cb)
	testRendering(t, expected, actual)
}

func TestRectangle4(t *testing.T) {
	p := &mockPage{}
	cb := NewBox(10, 20, 500, 600)
	sc := NewColorRGB(0.1, 0.2, 0.3)
	fc := NewColorRGB(0.4, 0.5, 0.6)
	r := newRectangle(p, 10, 20, 100, 200).StrokeColor(sc).FillColor(fc)
	expected := fmt.Sprintf("q %s %s 20 380 100 200 re B Q\n", sc.strokeColor(), fc.nonStrokeColor())
	actual := r.render(cb)
	testRendering(t, expected, actual)
}
