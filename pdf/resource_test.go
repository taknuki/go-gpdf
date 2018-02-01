package pdf

import "testing"

type mockFont struct {
}

func TestResource(t *testing.T) {
	r := newResource()
	if r.refNo() != 0 {
		t.Error("resource is not initial state: refNo != 0")
	}
	if r.age() != 0 {
		t.Error("resource is not initial state: age != 0")
	}
	if len(r.font) > 0 {
		t.Error("resource is not initial state: font is not empty")
	}
	if len(r.xobject) > 0 {
		t.Error("resource is not initial state: xobject is not empty")
	}
}
