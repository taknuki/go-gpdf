package pdf

import "testing"

func isEqualObjectIdentifer(t *testing.T, exRefNo, exAge int, actual objectIdentifier) {
	t.Helper()
	if exRefNo != actual.refNo() {
		t.Errorf("refNo: expected:%d actual:%d", exRefNo, actual.refNo())
	}
	if exAge != actual.age() {
		t.Errorf("age: expected:%d actual:%d", exAge, actual.age())
	}
}

func testCompillation(t *testing.T, expected, actual string) {
	t.Helper()
	if expected != actual {
		t.Errorf("compillation failed\nexpected: %s\nactual  :%s\n", expected, actual)
	}
}

func TestObjectIdentifier(t *testing.T) {
	obj := objectIdentifier{
		objectNumber:     1,
		generationNumber: 2,
	}
	isEqualObjectIdentifer(t, 1, 2, obj)
	ac3 := obj.indirectReference()
	ex3 := "1 2 R"
	if ex3 != ac3 {
		t.Errorf("indirectReference: expected:%s actual:%s", ex3, ac3)
	}
	ac4 := obj.bracket("foo")
	ex4 := "1 2 obj\nfoo\nendobj\n"
	if ex4 != ac4 {
		t.Errorf("bracket: expected:%s actual:%s", ex4, ac4)
	}
}

func TestCounter(t *testing.T) {
	c := newCounter()
	ac1 := c.count
	ex1 := 2
	if ex1 != ac1 {
		t.Error("counter does not start with 3")
	}
	ac2 := c.next()
	ex2 := 3
	if ex2 != ac2 {
		t.Errorf("next: expected: %d actual: %d", ex2, ac2)
	}
	ac3 := c.next()
	ex3 := 4
	if ex3 != ac3 {
		t.Errorf("next: expected: %d actual: %d", ex3, ac3)
	}
}

func TestObjectIdentifierNumber(t *testing.T) {
	c := newCounter()
	obj1 := objectIdentifier{
		objectNumber:     1,
		generationNumber: 1,
	}
	obj1.number(c)
	obj2 := objectIdentifier{
		objectNumber:     0,
		generationNumber: 1,
	}
	obj2.number(c)
	isEqualObjectIdentifer(t, 1, 1, obj1)
	isEqualObjectIdentifer(t, 3, 0, obj2)
}

type mockPDFObject struct {
	objectIdentifier
	res    string
	child1 stringObject
	child2 binaryObject
}

func (obj *mockPDFObject) compile() string {
	return obj.res
}

func (obj *mockPDFObject) walk(walker func(obj pdfObject)) {
	walker(obj)
	if obj.child1 != nil {
		walker(obj.child1)
	}
	if obj.child2 != nil {
		walker(obj.child2)
	}
}

type mockBinaryObject struct {
	objectIdentifier
	res []byte
	err error
}

func (obj *mockBinaryObject) compile() ([]byte, error) {
	return obj.res, obj.err
}
