package pdf

import (
	"fmt"
	"strings"
)

// Outline is a document-level navigation.
// Structure of Outline is hierarchical.
type Outline interface {
	stringObject
	traversableObject
	// AddItem adds a new Outline to its children.
	AddItem(title string, page Page, destType OutlineDestination) (o Outline)
}

// outline is a root node of a document outline tree.
type outline struct {
	objectIdentifier
	first *outlineItem
}

func newOutline() *outline {
	return &outline{
		objectIdentifier: objectIdentifier{},
		first:            nil,
	}
}

// AddItem adds the new OutlineItem to Outline
func (o *outline) AddItem(title string, page Page, destType OutlineDestination) Outline {
	if o.first != nil {
		return newOutlineItem(title, o, o.first.lastItem(), page, destType)
	}
	ol := newOutlineItem(title, o, nil, page, destType)
	o.first = ol
	return ol
}

func (o *outline) compile() string {
	return o.bracket(fmt.Sprintf(
		"<</Type /Outlines /First %s /Last %s>>",
		o.first.indirectReference(), o.first.lastItem().indirectReference()))
}

func (o *outline) walk(walker func(obj pdfObject)) {
	if o != nil {
		if o.first != nil {
			walker(o)
			o.first.walk(walker)
		}
	}
}

// outlineItem is a leaf node of a document outline tree.
type outlineItem struct {
	objectIdentifier
	title    string
	parent   Outline
	prev     *outlineItem
	next     *outlineItem
	first    *outlineItem
	page     Page
	destType OutlineDestination
}

func newOutlineItem(title string, parent Outline, prev *outlineItem, page Page, destType OutlineDestination) *outlineItem {
	oi := &outlineItem{
		objectIdentifier: objectIdentifier{},
		title:            title,
		parent:           parent,
		prev:             prev,
		page:             page,
		destType:         destType,
	}
	if prev != nil {
		prev.next = oi
	}
	return oi
}

// AddItem adds the new child OutlineItem to OutlineItem
func (oi *outlineItem) AddItem(title string, page Page, destType OutlineDestination) Outline {
	if oi.first != nil {
		return newOutlineItem(title, oi, oi.first.lastItem(), page, destType)
	}
	child := newOutlineItem(title, oi, nil, page, destType)
	oi.first = child
	return child
}

func (oi *outlineItem) lastItem() *outlineItem {
	if oi.next != nil {
		return oi.next.lastItem()
	}
	return oi
}

func (oi *outlineItem) compile() string {
	dict := make([]string, 3, 7)
	dict[0] = fmt.Sprintf("/Title (%s)", oi.title)
	dict[1] = fmt.Sprintf("/Parent %s", oi.parent.indirectReference())
	dict[2] = fmt.Sprintf("/Dest [%s %s]", oi.page.indirectReference(), oi.destType.compile())
	if oi.prev != nil {
		dict = append(dict, fmt.Sprintf("/Prev %s", oi.prev.indirectReference()))
	}
	if oi.next != nil {
		dict = append(dict, fmt.Sprintf("/Next %s", oi.next.indirectReference()))
	}
	if oi.first != nil {
		dict = append(dict, fmt.Sprintf("/First %s", oi.first.indirectReference()))
		dict = append(dict, fmt.Sprintf("/Last %s", oi.first.lastItem().indirectReference()))
	}
	return oi.bracket(fmt.Sprintf("<<%s>>", strings.Join(dict, " ")))
}

func (oi *outlineItem) walk(walker func(obj pdfObject)) {
	walker(oi)
	if oi.first != nil {
		oi.first.walk(walker)
	}
	if oi.next != nil {
		oi.next.walk(walker)
	}
}

// OutlineDestination is a destination type of a outline.
type OutlineDestination interface {
	compile() string
}

// singleton
var basicOutlineDestinationInst = &basicOutlineDestination{}

type basicOutlineDestination struct{}

func (d *basicOutlineDestination) compile() string {
	return "/Fit"
}

// OutlineDestinationBasic returns a basic OutlineDestination.
func OutlineDestinationBasic() OutlineDestination {
	return basicOutlineDestinationInst
}

type verticalOutlineDestination struct {
	top int
}

// OutlineDestinationVertical returns a OutlineDestination that locates a page vertical coordinate top positioned at the top edge of the window.
func OutlineDestinationVertical(top int) OutlineDestination {
	return &verticalOutlineDestination{top}
}

func (d *verticalOutlineDestination) compile() string {
	return fmt.Sprintf("/FitH %d", d.top)
}
