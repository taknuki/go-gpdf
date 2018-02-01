package pdf

import "testing"

func isEqualBox(t *testing.T, leftBottomX, leftBottomY, rightTopX, rightTopY int, actual *Box) {
	t.Helper()
	if actual.leftBottomX != leftBottomX {
		t.Errorf("leftBottomX: expected:%d actual:%d", leftBottomX, actual.leftBottomX)
	}
	if actual.leftBottomY != leftBottomY {
		t.Errorf("leftBottomY: expected:%d actual:%d", leftBottomY, actual.leftBottomY)
	}
	if actual.rightTopX != rightTopX {
		t.Errorf("rightTopX: expected:%d actual:%d", rightTopX, actual.rightTopX)
	}
	if actual.rightTopY != rightTopY {
		t.Errorf("rightTopY: expected:%d actual:%d", rightTopY, actual.rightTopY)
	}
}

func TestBox(t *testing.T) {
	b := NewBox(10, 20, 30, 40)
	isEqualBox(t, 10, 20, 30, 40, b)
	a0 := NewBoxA0()
	isEqualBox(t, 0, 0, 2384, 3370, a0)
	a1 := NewBoxA1()
	isEqualBox(t, 0, 0, 1684, 2384, a1)
	a2 := NewBoxA2()
	isEqualBox(t, 0, 0, 1191, 1684, a2)
	a3 := NewBoxA3()
	isEqualBox(t, 0, 0, 842, 1191, a3)
	a4 := NewBoxA4()
	isEqualBox(t, 0, 0, 595, 842, a4)
	inner := b.NewInnerBox(1, 2, 3, 4)
	isEqualBox(t, 14, 23, 28, 39, inner)
	expected := "[14 23 28 39]"
	actual := inner.compile()
	testCompillation(t, expected, actual)
}
