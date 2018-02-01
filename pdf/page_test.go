package pdf

import "testing"

type mockPage struct {
	objectIdentifier
	parent       *pageList
	renderResult GraphicsObject
}

func (p *mockPage) compile() string {
	return "mock"
}
func (p *mockPage) walk(walker func(obj pdfObject)) {
	walker(p)
}
func (p *mockPage) setParent(pl *pageList) {
	p.parent = pl
}
func (p *mockPage) AddFont(f Font)                                           {}
func (p *mockPage) AddImage(i *ImageResource)                                {}
func (p *mockPage) WriteText(x, y int, font Font, fontSize int, text string) {}
func (p *mockPage) Rectangle(startX, startY, width, height int) Rectangle    { return &mockRectangle{} }
func (p *mockPage) Line(startX, startY, endX, endY, lineWidth int) Line      { return &mockLine{} }
func (p *mockPage) Image(i *ImageResource, centerX, centerY float64) Image   { return &mockImage{} }
func (p *mockPage) render(obj GraphicsObject) {
	p.renderResult = obj
}

func TestPage1(t *testing.T) {
	pmb := NewBox(1, 2, 3, 4)
	pcb := NewBox(11, 12, 13, 14)
	mb := NewBox(21, 22, 23, 24)
	cb := NewBox(31, 32, 33, 34)
	r := newResource()
	pl := newRootPage(pmb, pcb)
	p := newPage(pl, mb, cb, r)
	isEqualBox(t, 21, 22, 23, 24, p.mb())
	isEqualBox(t, 31, 32, 33, 34, p.cb())
	if p.resource != r {
		t.Error("page is not initial state: resource is not costructor argument")
	}
}
