package pdf

import "fmt"

// trailer is a trailer dictionary.
type trailer struct {
	root      objectIdentifier
	size      int
	startXRef int
}

// newTrailer returns a trailer.
func newTrailer(root objectIdentifier, size int, startXRef int) *trailer {
	return &trailer{root, size, startXRef}
}

func (t *trailer) compile() string {
	return fmt.Sprintf(
		"trailer\n<</Root %s /Size %d>>\nstartxref\n%d\n%%%%EOF",
		t.root.indirectReference(), t.size, t.startXRef)
}
