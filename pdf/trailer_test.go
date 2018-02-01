package pdf

import "testing"

func TestTrailer(t *testing.T) {
	r := objectIdentifier{1, 2}
	tr := newTrailer(r, 3, 4)
	isEqualObjectIdentifer(t, 1, 2, tr.root)
	if tr.size != 3 {
		t.Errorf("size: expected:3, actual:%d", tr.size)
	}
	if tr.startXRef != 4 {
		t.Errorf("startXRef: expected:4, actual:%d", tr.startXRef)
	}
	actual := tr.compile()
	expected := "trailer\n<</Root 1 2 R /Size 3>>\nstartxref\n4\n%%EOF"
	testCompillation(t, expected, actual)
}
