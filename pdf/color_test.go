package pdf

import "testing"

func isEqualColor(t *testing.T, cs colorSpace, sc, nsc string, actual Color) {
	t.Helper()
	if cs != actual.colorSpace() {
		t.Errorf("colorspace: expected:%s acutual:%s", cs, actual.colorSpace())
	}
	if sc != actual.strokeColor() {
		t.Errorf("strokeColor: expected:%s acutual:%s", sc, actual.strokeColor())
	}
	if nsc != actual.nonStrokeColor() {
		t.Errorf("nonStrokeColor: expected:%s acutual:%s", nsc, actual.nonStrokeColor())
	}
}

func TestColor(t *testing.T) {
	gray := NewColorGrayScale(0.5)
	isEqualColor(t, colorSpaceDeviceGray, "0.500000 G", "0.500000 g", gray)
	rgb := NewColorRGB(0.2, 0.4, 0.8)
	isEqualColor(t, colorSpaceDeviceRGB, "0.200000 0.400000 0.800000 RG", "0.200000 0.400000 0.800000 rg", rgb)
	cmyk := NewColorCMYK(0.2, 0.4, 0.6, 0.8)
	isEqualColor(t, colorSpaceDeviceCMYK, "0.200000 0.400000 0.600000 0.800000 K", "0.200000 0.400000 0.600000 0.800000 k", cmyk)
	undef := newColorUndef()
	isEqualColor(t, colorSpaceUndefined, "", "", undef)
}

func TestColorSpace(t *testing.T) {
	ex1 := "/DeviceGray"
	ac1 := colorSpaceDeviceGray
	if ex1 != ac1.String() {
		t.Errorf("expected:%s actual:%s", ex1, ac1)
	}
	ex2 := "/DeviceRGB"
	ac2 := colorSpaceDeviceRGB
	if ex2 != ac2.String() {
		t.Errorf("expected:%s actual:%s", ex2, ac2)
	}
	ex3 := "/DeviceCMYK"
	ac3 := colorSpaceDeviceCMYK
	if ex3 != ac3.String() {
		t.Errorf("expected:%s actual:%s", ex3, ac3)
	}
	ex4 := "/Undefined"
	ac4 := colorSpaceUndefined
	if ex4 != ac4.String() {
		t.Errorf("expected:%s actual:%s", ex4, ac4)
	}
	ex5 := "/Undefined"
	ac5 := colorSpace(5)
	if ex5 != ac5.String() {
		t.Errorf("expected:%s actual:%s", ex5, ac5)
	}
}
