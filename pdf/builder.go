package pdf

import (
	"errors"
	"io"
	"strconv"
	"strings"
)

const (
	pdfVersion    = "1.7"
	pageTreeOrder = 6
)

// Builder is a pdf builder.
// Core API of this package.
type Builder struct {
	version string
	dc      *documentCatalog
	c       *counter
	fnm     *fontNameManager
	inm     *imageNameManager
	order   int
}

// NewBuilder returns a Builder.
// Arguments mb and cb specify default page size.
func NewBuilder(mb, cb *Box) *Builder {
	b := &Builder{
		version: pdfVersion,
		dc:      newDocumentCatalog(mb, cb),
		c:       newCounter(),
		fnm:     newFontNameManager(),
		inm:     newImageNameManager(),
		order:   pageTreeOrder,
	}
	b.dc.pages.resource = newResource()
	return b
}

// NewFontType1 creates a type1 font.
func (b *Builder) NewFontType1(fontName string) Font {
	return newFontType1(b.fnm.nextName(), fontName)
}

// NewFontComposite creates a type0 composite font.
// A composite font is one whose glyphs are obtained from a fontlike object called a CIDFont and a character encoding defined by a CMap.
func (b *Builder) NewFontComposite(cmap CMap, descendantFont CIDFont) Font {
	return newFontComposite(b.fnm.nextName(), cmap, descendantFont)
}

// NewFontCompositeEmbeded creates a type0 composite font.
// Font created by this function embed font program.
func (b *Builder) NewFontCompositeEmbeded(fontFilePath string) (Font, error) {
	return newFontCompositeEmbeded(b.fnm.nextName(), fontFilePath)
}

// AddFont adds the font to default resource.
func (b *Builder) AddFont(f Font) {
	b.dc.pages.resource.addFont(f)
}

// NewImageResource returns a ImageResource.
func (b *Builder) NewImageResource(width, height int, bitperComponent int, data []byte) *ImageResource {
	return newImageResource(b.inm.nextName(), width, height, bitperComponent, data)
}

// AddImage adds the image to default resource.
func (b *Builder) AddImage(i *ImageResource) {
	b.dc.pages.resource.addImage(i)
}

// AddPage adds the new Page.
func (b *Builder) AddPage() Page {
	return b.AddPageWithBox(nil, nil)
}

// AddPageWithBox adds the new Page with the specified box.
func (b *Builder) AddPageWithBox(mb, cb *Box) Page {
	return b.dc.pages.newPage(mb, cb, nil)
}

// Outline returns a document outline.
// If a outline has not been created, creates new outline and returns it.
func (b *Builder) Outline() Outline {
	return b.dc.Outline()
}

// Build creates a pdf.
func (b *Builder) Build(w io.Writer) error {
	err := b.build()
	if err != nil {
		return err
	}
	return b.write(w)
}

func (b *Builder) build() error {
	b.dc.pages.buildPageTree(b.order)
	errs := make([]string, 0)
	b.dc.walk(func(obj pdfObject) {
		obj.number(b.c)
		if font, isFont := obj.(Font); isFont {
			err := font.build()
			if err != nil {
				errs = append(errs, err.Error())
			}
		}
	})
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}
	return nil
}

func (b *Builder) write(w io.Writer) (err error) {
	return newWriter(w).
		start(b.version).
		writeTraversable(b.dc).
		finish(b.dc.objectIdentifier)
}

type fontNameManager struct {
	num int
}

func newFontNameManager() *fontNameManager {
	return &fontNameManager{0}
}

func (m *fontNameManager) nextName() (name string) {
	name = "/F" + strconv.Itoa(m.num)
	m.num++
	return
}

type imageNameManager struct {
	num int
}

func newImageNameManager() *imageNameManager {
	return &imageNameManager{0}
}

func (m *imageNameManager) nextName() (name string) {
	name = "/XI" + strconv.Itoa(m.num)
	m.num++
	return
}
