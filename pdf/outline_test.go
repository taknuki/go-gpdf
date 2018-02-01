package pdf

import (
	"sort"
	"testing"
)

func TestOutlineDestinationBasic(t *testing.T) {
	d := OutlineDestinationBasic()
	actual := d.compile()
	expected := "/Fit"
	testCompillation(t, expected, actual)
}

func TestOutlineDestinationVertical(t *testing.T) {
	d := OutlineDestinationVertical(1)
	actual := d.compile()
	expected := "/FitH 1"
	testCompillation(t, expected, actual)
}

func TestOutline(t *testing.T) {
	o := newOutline()
	if o.refNo() != 0 {
		t.Error("outline is not initial state: refNo != 0")
	}
	if o.age() != 0 {
		t.Error("outline is not initial state: age != 0")
	}
	if o.first != nil {
		t.Error("outline is not initial state: first is not nil")
	}

	// AddItem first
	p1 := &mockPage{}
	exTitle1 := "first"
	exDest1 := OutlineDestinationVertical(1)
	ol1 := o.AddItem(exTitle1, p1, exDest1)
	il1, isItem := ol1.(*outlineItem)

	// validate the first child
	if !isItem {
		t.Error("AddItem: is item not outlineItem?")
	}
	if o.first != il1 || il1.parent != o {
		t.Error("AddItem: is item not added?")
	}
	if il1.title != exTitle1 {
		t.Errorf("AddItem: title: expected:%s actual:%s", exTitle1, il1.title)
	}
	if il1.destType.compile() != exDest1.compile() {
		t.Errorf("AddItem: dest: expected:%s actual:%s", exDest1.compile(), il1.destType.compile())
	}
	if il1.page != p1 {
		t.Error("AddItem: is page not added?")
	}

	// AddItem second
	p2 := &mockPage{}
	exTitle2 := "second"
	exDest2 := OutlineDestinationVertical(2)
	ol2 := o.AddItem(exTitle2, p2, exDest2)
	il2, isItem := ol2.(*outlineItem)

	// validate the second child
	if !isItem {
		t.Error("AddItem: is item not outlineItem?")
	}
	if il1.lastItem() != il2 {
		t.Error("AddItem: is item not added?")
	}
	if il2.title != exTitle2 {
		t.Errorf("AddItem: title: expected:%s actual:%s", exTitle2, il2.title)
	}
	if il2.destType.compile() != exDest2.compile() {
		t.Errorf("AddItem: dest: expected:%s actual:%s", exDest2.compile(), il2.destType.compile())
	}
	if il2.page != p2 {
		t.Error("AddItem: is page not added?")
	}

	// AddItem third
	p3 := &mockPage{}
	exTitle3 := "third"
	exDest3 := OutlineDestinationVertical(3)
	ol3 := il1.AddItem(exTitle3, p3, exDest3)
	il3, isItem := ol3.(*outlineItem)

	// validate the third child
	if !isItem {
		t.Error("AddItem: is item not outlineItem?")
	}
	if il1.first != il3 || il3.parent != il1 {
		t.Error("AddItem: is item not added?")
	}
	if il3.title != exTitle3 {
		t.Errorf("AddItem: title: expected:%s actual:%s", exTitle3, il3.title)
	}
	if il3.destType.compile() != exDest3.compile() {
		t.Errorf("AddItem: dest: expected:%s actual:%s", exDest3.compile(), il3.destType.compile())
	}
	if il3.page != p3 {
		t.Error("AddItem: is page not added?")
	}

	// AddItem forth
	p4 := &mockPage{}
	exTitle4 := "forth"
	exDest4 := OutlineDestinationVertical(4)
	ol4 := il1.AddItem(exTitle4, p4, exDest4)
	il4, isItem := ol4.(*outlineItem)

	// validate the forth child
	if !isItem {
		t.Error("AddItem: is item not outlineItem?")
	}
	if il3.next != il4 || il4.parent != il1 {
		t.Error("AddItem: is item not added?")
	}
	if il4.title != exTitle4 {
		t.Errorf("AddItem: title: expected:%s actual:%s", exTitle4, il4.title)
	}
	if il4.destType.compile() != exDest4.compile() {
		t.Errorf("AddItem: dest: expected:%s actual:%s", exDest4.compile(), il4.destType.compile())
	}
	if il4.page != p4 {
		t.Error("AddItem: is page not added?")
	}

	o.objectNumber = 1
	il1.objectNumber = 2
	il2.objectNumber = 3
	il3.objectNumber = 4
	il4.objectNumber = 5

	actual0 := o.compile()
	expected0 := "1 0 obj\n<</Type /Outlines /First 2 0 R /Last 3 0 R>>\nendobj\n"
	testCompillation(t, expected0, actual0)

	actual1 := il1.compile()
	expected1 := "2 0 obj\n<</Title (first) /Parent 1 0 R /Dest [0 0 R /FitH 1] /Next 3 0 R /First 4 0 R /Last 5 0 R>>\nendobj\n"
	testCompillation(t, expected1, actual1)

	actual2 := il2.compile()
	expected2 := "3 0 obj\n<</Title (second) /Parent 1 0 R /Dest [0 0 R /FitH 2] /Prev 2 0 R>>\nendobj\n"
	testCompillation(t, expected2, actual2)

	nums := make([]int, 0)
	o.walk(func(obj pdfObject) {
		nums = append(nums, obj.refNo())
	})
	sort.Sort(sort.IntSlice(nums))
	for i, num := range nums {
		if i+1 != num {
			t.Errorf("walking failed: expected:%d acutal:%d", i+1, num)
		}
	}
}
