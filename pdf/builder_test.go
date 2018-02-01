package pdf

import (
	"os"
	"testing"
)

func TestBuilder(t *testing.T) {
	mb := NewBox(1, 2, 3, 4)
	cb := NewBox(11, 12, 13, 14)
	b := NewBuilder(mb, cb)
	o := b.Outline()
	if b.dc.outline != o {
		t.Error("outline is not connected to document catalog")
	}
	o2 := b.Outline()
	if o2 != o {
		t.Error("outline is dup")
	}
}

func TestFontNameManager(t *testing.T) {
	fnm := newFontNameManager()
	if "/F0" != fnm.nextName() {
		t.Error("fontNameMager is not initialized?")
	}
	if "/F1" != fnm.nextName() {
		t.Error("fontNameMager is not incremented?")
	}
}

func TestImageNameManager(t *testing.T) {
	inm := newImageNameManager()
	if "/XI0" != inm.nextName() {
		t.Error("imageNameMager is not initialized?")
	}
	if "/XI1" != inm.nextName() {
		t.Error("imageNameMager is not incremented?")
	}
}

func ExampleBuilder() {
	// page setting
	leftMargin := 50
	rightMargin := 50
	topMargin := 50
	bottomMargin := 50
	mb := NewBoxA4()
	cb := mb.NewInnerBox(topMargin, rightMargin, bottomMargin, leftMargin)

	b := NewBuilder(mb, cb)
	f, _ := os.Open("expamle.pdf")
	defer f.Close()
	b.Build(f)
}
