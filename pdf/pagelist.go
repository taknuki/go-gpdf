package pdf

import (
	"fmt"
	"strings"
)

// pageNode is a common feature container of a node of a pdf page tree.
type pageNode struct {
	objectIdentifier
	parent   *pageList
	mediaBox *Box
	cropBox  *Box
	resource *resource
}

// cb is the effective cropbox in this page.
func (pt *pageNode) cb() (cropBox *Box) {
	if pt.cropBox != nil {
		cropBox = pt.cropBox
	} else {
		cropBox = pt.parent.cb()
	}
	return
}

// mb is the effective mediabox in this page.
func (pt *pageNode) mb() (mediaBox *Box) {
	if pt.mediaBox != nil {
		mediaBox = pt.mediaBox
	} else {
		mediaBox = pt.parent.mb()
	}
	return
}

// pageList is a Pages node of a pdf page tree.
type pageList struct {
	pageNode
	pageLists []*pageList
	pages     []Page
}

// newRootPage creates the root Pages node.
func newRootPage(mb, cb *Box) *pageList {
	return &pageList{
		pageNode: pageNode{
			objectIdentifier: objectIdentifier{
				objectNumber:     2,
				generationNumber: 0,
			},
			mediaBox: mb,
			cropBox:  cb,
		},
	}
}

// newPage creates the child Page node of this Pages node.
func (pl *pageList) newPage(mb, cb *Box, r *resource) (p Page) {
	p = newPage(pl, mb, cb, r)
	pl.addPage(p)
	return
}

// addPageList adds the child Pages node of this Pages node.
func (pl *pageList) newPageList(mb, cb *Box, r *resource) (child *pageList) {
	child = &pageList{
		pageNode: pageNode{
			objectIdentifier: objectIdentifier{},
			parent:           pl,
			mediaBox:         mb,
			cropBox:          cb,
			resource:         r,
		},
	}
	pl.pageLists = append(pl.pageLists, child)
	return
}

// addPage adds the child Page node of this Pages node.
func (pl *pageList) addPage(p Page) {
	p.setParent(pl)
	pl.pages = append(pl.pages, p)
}

// asPDF is the pdf object expression of this Pages node.
func (pl *pageList) compile() string {
	list := make([]string, 0, 7)
	list = append(list, "/Type /Pages")
	if pl.parent != nil {
		list = append(list, fmt.Sprintf("/Parent %s", pl.parent.indirectReference()))
	}
	if pl.mediaBox != nil {
		list = append(list, fmt.Sprintf("/MediaBox %s", pl.mediaBox.compile()))
	}
	if pl.cropBox != nil {
		list = append(list, fmt.Sprintf("/CropBox %s", pl.cropBox.compile()))
	}
	if pl.resource != nil {
		list = append(list, fmt.Sprintf("/Resources %s", pl.resource.indirectReference()))
	}
	list = append(list, fmt.Sprintf("/Kids %s", pl.kidsAsPDF()))
	list = append(list, fmt.Sprintf("/Count %d", pl.count()))
	return pl.bracket(fmt.Sprintf("<<%s>>", strings.Join(list, " ")))
}

func (pl *pageList) walk(walker func(obj pdfObject)) {
	walker(pl)
	pl.resource.walk(walker)
	for _, childPL := range pl.pageLists {
		childPL.walk(walker)
	}
	for _, childP := range pl.pages {
		childP.walk(walker)
	}
}

// kidsAsPDF is the pdf object expression of the children of this Pages node.
func (pl *pageList) kidsAsPDF() string {
	list := make([]string, 0, len(pl.pageLists)+len(pl.pages))
	for _, pageList := range pl.pageLists {
		list = append(list, pageList.indirectReference())
	}
	for _, page := range pl.pages {
		list = append(list, page.indirectReference())
	}
	return fmt.Sprintf("[%s]", strings.Join(list, " "))
}

// count is The number of leaf nodes (page objects) that are descendants of this node within the page tree.
func (pl *pageList) count() (c int) {
	for _, pageList := range pl.pageLists {
		c += pageList.count()
	}
	c += len(pl.pages)
	return
}

func (pl *pageList) buildPageTree(order int) {
	listSize := len(pl.pages)
	n := listSize / order
	if listSize%order > 0 {
		n++
	}
	pageLists := make([]*pageList, 0)
	descendants := pl.pages
	pl.pageLists = pageLists
	index := 0
	for i := 0; i < n; i++ {
		child := pl.newPageList(nil, nil, nil)
		for j := 0; j < order; j++ {
			if index >= listSize {
				break
			}
			child.addPage(descendants[index])
			index++
		}
	}
	pl.pages = nil
	pl.buildSubPageTree(order)
}

func (pl *pageList) buildSubPageTree(order int) {
	listSize := len(pl.pageLists)
	if listSize <= order {
		return
	}
	n := listSize / order
	if listSize%order > 0 {
		n++
	}
	pageLists := make([]*pageList, 0)
	descendants := pl.pageLists
	pl.pageLists = pageLists
	index := 0
	for i := 0; i < n; i++ {
		child := pl.newPageList(nil, nil, nil)
		for j := 0; j < order; j++ {
			if index >= listSize {
				break
			}
			// child.addPageList(descendants[index])
			descendants[index].parent = child
			child.pageLists = append(child.pageLists, descendants[index])
			// increment
			index++
		}
	}
	pl.buildSubPageTree(order)
}
