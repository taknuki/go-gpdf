package pdf

import "testing"

func testFont(t *testing.T, resourceName, baseFont, subType, compillation string, actual Font) {
	t.Helper()
	if resourceName != actual.resourceName() {
		t.Errorf("resourceName:\nexpected:%s\nactual  :%s\n", resourceName, actual.resourceName())
	}
	if baseFont != actual.baseFont() {
		t.Errorf("baseFont:\nexpected:%s\nactual  :%s\n", baseFont, actual.baseFont())
	}
	if subType != actual.subtype() {
		t.Errorf("subType:\nexpected:%s\nactual  :%s\n", subType, actual.subtype())
	}
	testCompillation(t, compillation, actual.compile())
}

func TestFontDefault(t *testing.T) {
	name := "/F0"
	f := newDefaultFont(name)
	if name != f.resourceName() {
		t.Errorf("resourceName:\nexpected:%s\nactual  :%s\n", name, f.resourceName())
	}
	if f.build() != nil {
		t.Error("build: anything is done?")
	}

}

func TestFontType1(t *testing.T) {
	name := "/F0"
	fontName := "/Times-Roman"
	subType := "/Type1"
	expected := "0 0 obj\n<</Type /Font /BaseFont /Times-Roman /Subtype /Type1>>\nendobj\n"
	f := newFontType1(name, fontName)
	f.walk(func(obj pdfObject) {
		if obj != f {
			t.Error("walk: what's returned?")
		}
	})
	testFont(t, name, fontName, subType, expected, f)
}

func TestFontCompisite(t *testing.T) {
	// TODO
}

func TestPredefinedCMap(t *testing.T) {
	cm := CMapIdentityH
	expected := "Identity-H"
	if expected != cm.Name() {
		t.Errorf("name:\nexpected:%s\nactual  :%s\n", expected, cm.Name())
	}
	if expected != cm.Name() {
		t.Errorf("name:\nexpected:%s\nactual  :%s\n", expected, cm.Name())
	}
	testCompillation(t, "/"+expected, cm.compile())
}

func TestCIDSystemInfo(t *testing.T) {
	si := CIDSystemInfoAdobeJapan6
	expected := "Adobe-Japan1-6"
	actual := si.String()
	if expected != actual {
		t.Errorf("String:\nexpected:%s\nactual  :%s\n", expected, actual)
	}
	testCompillation(t, "<</Registry (Adobe) /Ordering (Japan1) /Supplement 6>>", si.compile())
}
