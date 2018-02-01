package pdf

import (
	"fmt"
	"strings"
)

// Page is a pdf page.
// If AddXxx method is called even once, this page uses the original resource and the default resource is unavailable.
type Page interface {
	stringObject
	traversableObject
	setParent(pl *pageList)
	// AddFont adds the font to this page.
	AddFont(f Font)
	// AddImage adds the image to this page.
	AddImage(i *ImageResource)
	// WriteText writes the text on this page.
	WriteText(x, y int, font Font, fontSize int, text string)
	// Rectangle adds a rectangle to this page.
	Rectangle(startX, startY, width, height int) Rectangle
	// Line adds a line to this page.
	Line(startX, startY, endX, endY, lineWidth int) Line
	// Image adds the image to this page.
	Image(i *ImageResource, centerX, centerY float64) Image
	render(obj GraphicsObject)
}

type page struct {
	pageNode
	contents *stream
}

// AddFont adds the font to this page.
// If this method is called, this page uses the original resource and the default resource is unavailable.
func (p *page) AddFont(f Font) {
	if p.resource == nil {
		p.resource = newResource()
	}
	p.resource.addFont(f)
}

// AddImage adds the image to this page.
// If this method is called, this page uses the original resource and the default resource is unavailable.
func (p *page) AddImage(i *ImageResource) {
	if p.resource == nil {
		p.resource = newResource()
	}
	p.resource.addImage(i)
}

func (p *page) WriteText(x, y int, font Font, fontSize int, text string) {
	p.addStringContent(p.text(x, y, font, fontSize, text))
}

func (p *page) Rectangle(startX, startY, width, height int) Rectangle {
	return newRectangle(p, startX, startY, width, height)
}

func (p *page) Line(startX, startY, endX, endY, lineWidth int) Line {
	return newLine(p, startX, startY, endX, endY, lineWidth)
}

func (p *page) Image(i *ImageResource, centerX, centerY float64) Image {
	return newImage(p, i, centerX, centerY)
}

// newPage creates Page.
func newPage(pl *pageList, mb *Box, cb *Box, r *resource) *page {
	return &page{
		pageNode: pageNode{
			objectIdentifier: objectIdentifier{},
			parent:           pl,
			mediaBox:         mb,
			cropBox:          cb,
			resource:         r,
		},
		contents: newDeflatedStream(),
	}
}

func (p *page) setParent(pl *pageList) {
	p.parent = pl
}

func (p *page) render(obj GraphicsObject) {
	p.contents.addStringDatum(obj.render(p.cb()))
}

// addStringContent adds content whoose type is string.
func (p *page) addStringContent(content string) {
	p.contents.addStringDatum(content)
}

// asPDF is the pdf object expression of this Page node.
func (p *page) compile() string {
	list := make([]string, 0, 5)
	list = append(list, fmt.Sprintf("/Type /Page /Parent %s", p.parent.indirectReference()))
	if p.mediaBox != nil {
		list = append(list, fmt.Sprintf("/MediaBox %s", p.mediaBox.compile()))
	}
	if p.cropBox != nil {
		list = append(list, fmt.Sprintf("/CropBox %s", p.cropBox.compile()))
	}
	if p.resource != nil {
		list = append(list, fmt.Sprintf("/Resources %s", p.resource.indirectReference()))
	}
	list = append(list, fmt.Sprintf("/Contents [%s]", p.contents.indirectReference()))
	return p.bracket(fmt.Sprintf("<<%s>>", strings.Join(list, " ")))
}

func (p *page) walk(walker func(obj pdfObject)) {
	walker(p)
	p.resource.walk(walker)
	walker(p.contents)
}

func (p *page) text(x, y int, font Font, fontSize int, text string) string {
	return font.createText(p.cb().leftBottomX+x, p.cb().rightTopY-fontSize-y, fontSize, text)
}
