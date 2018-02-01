package pdf

import (
	"fmt"
	"strings"
)

// documentCatalog is a root object of a pdf document graph.
type documentCatalog struct {
	objectIdentifier
	pages   *pageList
	outline Outline
}

// newDocumentCatalog returns a document catalog with a root page.
// Argments mb and cb are applied to the root page.
func newDocumentCatalog(mb, cb *Box) *documentCatalog {
	return &documentCatalog{
		objectIdentifier: objectIdentifier{
			objectNumber:     1,
			generationNumber: 0,
		},
		pages:   newRootPage(mb, cb),
		outline: nil,
	}
}

// Outline returns a document outline.
// If a outline has not been created, creates new outline and returns it.
func (dc *documentCatalog) Outline() Outline {
	if dc.outline != nil {
		return dc.outline
	}
	o := newOutline()
	dc.outline = o
	return o
}

func (dc *documentCatalog) compile() string {
	options := make([]string, 1, 2)
	options[0] = fmt.Sprintf("/Pages %s", dc.pages.indirectReference())
	if dc.outline != nil {
		options = append(options, fmt.Sprintf("/Outlines %s", dc.outline.indirectReference()))
	}
	return dc.bracket(fmt.Sprintf(
		"<</Type /Catalog %s>>",
		strings.Join(options, " ")))
}

func (dc *documentCatalog) walk(walker func(obj pdfObject)) {
	walker(dc)
	dc.pages.walk(walker)
	if dc.outline != nil {
		dc.outline.walk(walker)
	}
}
