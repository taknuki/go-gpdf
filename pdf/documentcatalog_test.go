package pdf

import "testing"

func TestDocumentCatalog(t *testing.T) {
	leftMargin := 50
	rightMargin := 50
	topMargin := 50
	bottomMargin := 50
	mb := NewBoxA4()
	cb := mb.NewInnerBox(topMargin, rightMargin, bottomMargin, leftMargin)
	dc := newDocumentCatalog(mb, cb)
	if dc.outline != nil {
		t.Error("when documentCatalog is created, outline must be nil pointer")
	}
	ac1 := dc.refNo()
	ex1 := 1
	if ex1 != ac1 {
		t.Errorf("refNo: expected:%d actual:%d", ex1, ac1)
	}
	ac2 := dc.age()
	ex2 := 0
	if ex1 != ac1 {
		t.Errorf("age: expected:%d actual:%d", ex2, ac2)
	}
	testCompillation(t, "1 0 obj\n<</Type /Catalog /Pages 2 0 R>>\nendobj\n", dc.compile())
}
