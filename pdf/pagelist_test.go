package pdf

import "testing"

func TestPageList1(t *testing.T) {
	rmb := NewBox(1, 2, 3, 4)
	rcb := NewBox(11, 12, 13, 14)
	rp := newRootPage(rmb, rcb)
	if rp.refNo() != 2 {
		t.Error("root page is not initial state: refNo != 2")
	}
	if rp.age() != 0 {
		t.Error("root page is not initial state: age != 0")
	}
	if rp.parent != nil {
		t.Error("root page is not root")
	}
	if rp.resource != nil {
		t.Error("root page is not initial state: resource is not nil")
	}
	if len(rp.pageLists) > 0 {
		t.Error("initial root page has any pageLists")
	}
	if len(rp.pages) > 0 {
		t.Error("initial root page has any pages")
	}
	isEqualBox(t, 1, 2, 3, 4, rp.mb())
	isEqualBox(t, 11, 12, 13, 14, rp.cb())

	pl1 := rp.newPageList(nil, nil, nil)
	if pl1.refNo() != 0 {
		t.Error("pageList is not initial state: refNo != 0")
	}
	if pl1.age() != 0 {
		t.Error("pageList is not initial state: age != 0")
	}
	if len(pl1.pageLists) > 0 {
		t.Error("initial pageList has any pageLists")
	}
	if len(pl1.pages) > 0 {
		t.Error("initial pageList has any pages")
	}
	if rp.refNo() != 2 {
		t.Error("root page is not initial state: refNo != 2")
	}
	if rp.age() != 0 {
		t.Error("root page is not initial state: age != 0")
	}
	isEqualBox(t, 1, 2, 3, 4, pl1.mb())
	isEqualBox(t, 11, 12, 13, 14, pl1.cb())
	pl1.objectNumber = 3

	pmb := NewBox(21, 22, 23, 24)
	pcb := NewBox(31, 32, 33, 34)
	pr := newResource()
	pr.objectNumber = 4
	pl2 := pl1.newPageList(pmb, pcb, pr)
	pl2.objectNumber = 5
	if pl2.parent != pl1 {
		t.Error("pageList has unexpected parent")
	}
	isEqualBox(t, 21, 22, 23, 24, pl2.mb())
	isEqualBox(t, 31, 32, 33, 34, pl2.cb())

	expectedRP := "2 0 obj\n<</Type /Pages /MediaBox [1 2 3 4] /CropBox [11 12 13 14] /Kids [3 0 R] /Count 0>>\nendobj\n"
	actualRP := rp.compile()
	testCompillation(t, expectedRP, actualRP)
	expectedPL := "5 0 obj\n<</Type /Pages /Parent 3 0 R /MediaBox [21 22 23 24] /CropBox [31 32 33 34] /Resources 4 0 R /Kids [] /Count 0>>\nendobj\n"
	actualPL := pl2.compile()
	testCompillation(t, expectedPL, actualPL)
}

func TestPageList2(t *testing.T) {
	rmb := NewBox(1, 2, 3, 4)
	rcb := NewBox(11, 12, 13, 14)
	rp := newRootPage(rmb, rcb)
	p1 := &mockPage{}
	p1.objectNumber = 11
	rp.addPage(p1)
	p2 := &mockPage{}
	p2.objectNumber = 12
	rp.addPage(p2)
	p3 := &mockPage{}
	p3.objectNumber = 13
	rp.addPage(p3)
	p4 := &mockPage{}
	p4.objectNumber = 14
	rp.addPage(p4)
	p5 := &mockPage{}
	p5.objectNumber = 15
	rp.addPage(p5)
	p6 := &mockPage{}
	p6.objectNumber = 16
	rp.addPage(p6)
	p7 := &mockPage{}
	p7.objectNumber = 17
	rp.addPage(p7)
	p8 := &mockPage{}
	p8.objectNumber = 18
	rp.addPage(p8)
	p9 := &mockPage{}
	p9.objectNumber = 19
	rp.addPage(p9)
	pmb := NewBox(21, 22, 23, 24)
	pcb := NewBox(31, 32, 33, 34)
	pr := newResource()
	pr.objectNumber = 4
	rp.newPage(pmb, pcb, pr)

	expected := "2 0 obj\n<</Type /Pages /MediaBox [1 2 3 4] /CropBox [11 12 13 14] /Kids [11 0 R 12 0 R 13 0 R 14 0 R 15 0 R 16 0 R 17 0 R 18 0 R 19 0 R 0 0 R] /Count 10>>\nendobj\n"
	actual := rp.compile()
	testCompillation(t, expected, actual)

	rp.buildPageTree(3)
	actualNum := 0
	rp.walk(func(obj pdfObject) {
		actualNum++
	})
	expectedNum := 19
	if expectedNum != actualNum {
		t.Errorf("walking failed: expected:%d actual:%d", expectedNum, actualNum)
	}
}
