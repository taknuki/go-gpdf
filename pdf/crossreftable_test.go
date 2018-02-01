package pdf

import "testing"

func isEqualCrossRefEntry(t *testing.T, expected, actual *crossRefEntry) {
	t.Helper()
	if actual.offset != expected.offset {
		t.Errorf("offset: expected:%d actual:%d", expected.offset, actual.offset)
	}
	if actual.age != expected.age {
		t.Errorf("age: expected:%d actual:%d", expected.age, actual.age)
	}
	if actual.inuse != expected.inuse {
		t.Errorf("inuse: expected:%t actual:%t", expected.inuse, actual.inuse)
	}
}

func TestCrossRefTable(t *testing.T) {
	crt := newCrossRefTable()
	if len(crt.entries) != 1 {
		t.Errorf("entry size: expected:1 actual:%d", len(crt.entries))
	}
	actual1 := crt.entries[0]
	expected1 := &crossRefEntry{
		offset: 0,
		age:    65535,
		inuse:  false,
	}
	isEqualCrossRefEntry(t, expected1, actual1)
	obj1 := &objectIdentifier{1, 1}
	crt.addNewEntry(obj1, 123)
	if len(crt.entries) != 2 {
		t.Errorf("entry size: expected:2 actual:%d", len(crt.entries))
	}
	actual2, ok := crt.entries[1]
	if !ok {
		t.Error("entry is not added?")
	}
	expected2 := &crossRefEntry{
		offset: 123,
		age:    1,
		inuse:  true,
	}
	isEqualCrossRefEntry(t, expected2, actual2)
	if !crt.hasEntry(obj1) {
		t.Fail()
	}
	obj2 := &objectIdentifier{3, 0}
	crt.addNewEntry(obj2, 200)
	obj3 := &objectIdentifier{4, 0}
	crt.addNewEntry(obj3, 300)
	obj4 := &objectIdentifier{5, 0}
	crt.addNewEntry(obj4, 400)
	obj5 := &objectIdentifier{10, 0}
	crt.addNewEntry(obj5, 500)
	obj6 := &objectIdentifier{11, 0}
	crt.addNewEntry(obj6, 600)
	actual3 := crt.compile()
	expected3 := "xref\n"
	expected3 += "0 2\n"
	expected3 += "0000000000 65535 f \n"
	expected3 += "0000000123 00001 n \n"
	expected3 += "3 3\n"
	expected3 += "0000000200 00000 n \n"
	expected3 += "0000000300 00000 n \n"
	expected3 += "0000000400 00000 n \n"
	expected3 += "10 2\n"
	expected3 += "0000000500 00000 n \n"
	expected3 += "0000000600 00000 n \n"
	testCompillation(t, expected3, actual3)
}

func TestNewCrossRefEntry(t *testing.T) {
	actual := newCrossRefEntry(123, 1, true)
	expected := &crossRefEntry{
		offset: 123,
		age:    1,
		inuse:  true,
	}
	isEqualCrossRefEntry(t, expected, actual)
}

func TestCrossRefEntryCompile(t *testing.T) {
	entry1 := newCrossRefEntry(123, 1, true)
	actual1 := entry1.compile()
	expected1 := "0000000123 00001 n \n"
	if actual1 != expected1 {
		t.Errorf("expected:'%s' actual:'%s'", expected1, actual1)
	}
	entry2 := newCrossRefEntry(123, 1, false)
	actual2 := entry2.compile()
	expected2 := "0000000123 00001 f \n"
	if actual2 != expected2 {
		t.Errorf("expected:'%s' actual:'%s'", expected2, actual2)
	}
}
