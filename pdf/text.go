package pdf

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf16"

	"github.com/taknuki/go-opentype/opentype"
)

// Font interface provides a functionality of pdf font dictionary.
type Font interface {
	stringObject
	traversableObject
	// resourceName returns the name in a resource dictionary.
	resourceName() string
	// baseFont retunrs the PostScript name of the font.
	baseFont() string
	// subtype returns:
	// Type0: A composite font—a font composed of glyphs from a descendant CIDFont
	// TYpe1: A font that defines glyph shapes using Type 1 font technology
	subtype() string
	// build setup a embeded font program for writing pdf document.
	build() error
	createText(x, y int, fontSize int, text string) string
}

// defaultFont provides a common functionality of pdf font dictionary.
type defaultFont struct {
	objectIdentifier
	resName string
}

func newDefaultFont(name string) defaultFont {
	return defaultFont{
		resName: name,
	}
}

func (f *defaultFont) resourceName() string {
	return f.resName
}

func (f *defaultFont) build() error {
	return nil
}

type type1Font struct {
	defaultFont
	fontName string
}

// newFontType1 creates a type1 font.
func newFontType1(name, fontName string) Font {
	return &type1Font{
		defaultFont: newDefaultFont(name),
		fontName:    fontName,
	}
}

func (f *type1Font) baseFont() string {
	return f.fontName
}

func (f *type1Font) subtype() string {
	return "/Type1"
}

func (f *type1Font) walk(walker func(obj pdfObject)) {
	walker(f)
}

func (f *type1Font) compile() string {
	return f.bracket(fmt.Sprintf("<</Type /Font /BaseFont %s /Subtype %s>>", f.baseFont(), f.subtype()))
}

// The Tf operator identifies the font to be used.
// The Td operator adjusts the current text position to begin painting glyphs.
// The Tj operator takes a string operand and paints the corresponding glyphs.
func (f *type1Font) createText(x, y int, fontSize int, text string) string {
	texts := strings.Split(text, "\n")
	opes := make([]string, 0, len(texts))
	for _, t := range texts {
		str := ""
		for _, w := range utf16.Encode([]rune(t)) {
			str += fmt.Sprintf("%04X", w)
		}
		opes = append(opes, fmt.Sprintf("<%s> Tj", str))
	}
	return fmt.Sprintf(
		"BT\n%d %d Td %d TL \n%s %d. Tf\n%s\nET\n",
		x, y, fontSize, f.resourceName(), fontSize, strings.Join(opes, " T*\n"))
}

// compositeFont is a type0 composite font.
// A composite font is one whose glyphs are obtained from a fontlike object called a CIDFont and a character encoding defined by a CMap.
type compositeFont struct {
	defaultFont
	cmap           CMap
	descendantFont CIDFont
}

// newFontComposite creates a type0 composite font.
// A composite font is one whose glyphs are obtained from a fontlike object called a CIDFont and a character encoding defined by a CMap.
func newFontComposite(name string, cmap CMap, descendantFont CIDFont) Font {
	return &compositeFont{
		defaultFont:    newDefaultFont(name),
		cmap:           cmap,
		descendantFont: descendantFont,
	}
}

// newFontCompositeEmbeded creates a type0 composite font.
// Font created by this function embed font program.
func newFontCompositeEmbeded(name, fontFilePath string) (Font, error) {
	return newCompositeFontEmbeded(name, CMapIdentityH, fontFilePath)
}

// newFontCompositeEmbededVertical creates a type0 composite font.
// Font created by this function embed font program adjusted for vertical order.
//
// UNDER IMPLEMENTATION: go-opentype has not supported font substitution yet.
func newFontCompositeEmbededVertical(name, fontFilePath string) (Font, error) {
	return newCompositeFontEmbeded(name, CMapIdentityV, fontFilePath)
}

func newCompositeFontEmbeded(name string, cmap CMap, fontFilePath string) (Font, error) {
	fontFile, err := os.Open(fontFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create CompositeFont: %s", err)
	}
	defer fontFile.Close()
	font, err := opentype.ParseFont(fontFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create CompositeFont: %s", err)
	}
	return newFontComposite(name, cmap, newCIDFontOpenType(font)), nil
}

func (f *compositeFont) baseFont() string {
	return f.descendantFont.parentBaseFont(f.cmap.Name())
}

func (f *compositeFont) subtype() string {
	return "/Type0"
}

func (f *compositeFont) walk(walker func(obj pdfObject)) {
	walker(f)
	f.descendantFont.walk(walker)
}

func (f *compositeFont) build() error {
	return f.descendantFont.build()
}

func (f *compositeFont) compile() string {
	return f.bracket(fmt.Sprintf(
		"<</Type /Font /BaseFont %s /Subtype %s /Encoding %s /DescendantFonts [%s]>>",
		f.baseFont(), f.subtype(), f.cmap.compile(), f.descendantFont.compile()))
}

func (f *compositeFont) createText(x, y int, fontSize int, text string) string {
	return f.descendantFont.createText(f.resourceName(), x, y, fontSize, text)
}

// CMap specifies the mapping from character codes to CIDs.
//
// Embedded CMap is not supporeted.
type CMap interface {
	Name() string
	compile() string
}

type predefinedCMap struct {
	name string
}

// newPredefinedCMap returns a predefined CMap.
func newPredefinedCMap(name string) CMap {
	return &predefinedCMap{
		name: name,
	}
}

func (cm *predefinedCMap) Name() string {
	return cm.name
}

func (cm *predefinedCMap) compile() string {
	return "/" + cm.name
}

type CIDFont interface {
	stringCompiler
	traversableObject
	// parentBaseFont returns BaseFont of Type 0 Font Dictionaries
	parentBaseFont(cmapName string) string
	// build embeded font program
	build() error
	// create text operator
	createText(fontName string, x, y int, fontSize int, text string) string
	// The PostScript name of the CIDFont.
	// For Type 0 CIDFonts, this is usually the value of the CIDFontName entry in the CIDFont program.
	// For Type 2 CIDFonts, it is derived the same way as for a simple TrueType font; see Section 5.5.2, “TrueType Fonts.” In either case, the name can have a sub- set prefix if appropriate; see Section 5.5.3, “Font Subsets.”
	BaseFont() string
	// CIDFontType0: A Type 0 CIDFont contains glyph descriptions based on the Adobe Type 1 font format
	// CIDFontType2: A Type 2 CIDFont contains glyph descriptions based on the TrueType font format
	SubType() string
}

// NewCIDFont creates a CIDFont
func NewCIDFont(baseFont, subType string, cidSystemInfo *CIDSystemInfo, fontDescriptor *FontDescriptor) (CIDFont, error) {
	switch subType {
	case "/CIDFontType0":
		return &cidFontSubType0{
			abstractCIDFont: abstractCIDFont{
				baseFont:       baseFont,
				cidSystemInfo:  cidSystemInfo,
				fontDescriptor: fontDescriptor,
			},
		}, nil
	case "/CIDFontType2":
		return nil, fmt.Errorf("Type 2 CIDFont must be embeded")
	default:
		return nil, fmt.Errorf("Illegal CIDFont Subtype: %s", subType)
	}
}

func newCIDFontOpenType(font *opentype.Font) CIDFont {
	baseFont := "unknown"
	for _, nr := range font.Name.NameRecords {
		if opentype.NameIDPostScriptName == nr.NameID {
			baseFont = "/" + nr.Value
			break
		}
	}
	fontDescriptor := NewFontDescriptor(baseFont, 4, NewBox(-538, -374, 1254, 1418), 0, 1418, -374, 763, 116)
	fontDescriptor.fontFile2 = newDeflatedStream()
	if opentype.SfntVersionCFFOpenType == font.SfntVersion {
		return &cidFontSubType0{
			abstractCIDFont: abstractCIDFont{
				baseFont:       baseFont,
				cidSystemInfo:  CIDSystemInfoAdobeIdentity0,
				fontDescriptor: fontDescriptor,
			},
		}
	}
	newGIDMap := make(map[uint16]uint16)
	newGIDMap[0] = 0
	cm := font.CMap
	var gidMap map[int32]uint16
	for _, er := range cm.EncodingRecords {
		if er.PlatformID == opentype.PlatformIDUnicode {
			gidMap = er.CMap()
			break
		}
	}
	return &cidFontSubType2{
		abstractCIDFont: abstractCIDFont{
			baseFont:       baseFont,
			cidSystemInfo:  CIDSystemInfoAdobeIdentity0,
			fontDescriptor: fontDescriptor,
		},
		gidMap:      gidMap,
		newGIDMap:   newGIDMap,
		embededFont: font,
	}
}

type abstractCIDFont struct {
	baseFont       string
	cidSystemInfo  *CIDSystemInfo
	fontDescriptor *FontDescriptor
	dw             int
	w              string
}

func (f *abstractCIDFont) BaseFont() string {
	return f.baseFont
}

func (f *abstractCIDFont) FontDescriptor() *FontDescriptor {
	return f.fontDescriptor
}

func (f *abstractCIDFont) walk(walker func(obj pdfObject)) {}

func (f *abstractCIDFont) compileHelper(subType string) string {
	options := make([]string, 0)
	if f.dw > 0 {
		options = append(options, fmt.Sprintf("/DW %d", f.dw))
	}
	if f.w != "" {
		options = append(options, fmt.Sprintf("/W [%s]", f.w))
	}
	return fmt.Sprintf(
		"<</Type /Font /BaseFont %s /Subtype %s /CIDSystemInfo %s /FontDescriptor %s /CIDToGIDMap /Identity %s>>",
		f.baseFont, subType, f.cidSystemInfo.compile(), f.fontDescriptor.compile(), strings.Join(options, " "))
}

// cidFontSubType0 is a Type 0 CIDFont.
// A Type 0 CIDFont contains glyph descriptions based on the Adobe Type 1 font format.
// The CIDFont program contains glyph descriptions that are identified by CIDs.
// The CIDFont program identifies the character collection by a CIDSystemInfo dictionary.
type cidFontSubType0 struct {
	abstractCIDFont
}

func (f *cidFontSubType0) SubType() string {
	return "/CIDFontType0"
}

func (f *cidFontSubType0) parentBaseFont(cmapName string) string {
	return f.baseFont + "-" + cmapName
}

func (f *cidFontSubType0) compile() string {
	return f.compileHelper(f.SubType())
}

func (f *cidFontSubType0) build() error {
	return nil
}

func (f *cidFontSubType0) createText(fontName string, x, y int, fontSize int, text string) string {
	texts := strings.Split(text, "\n")
	opes := make([]string, 0, len(texts))
	for _, t := range texts {
		str := ""
		for _, w := range utf16.Encode([]rune(t)) {
			str += fmt.Sprintf("%04X", w)
		}
		opes = append(opes, fmt.Sprintf("<%s> Tj", str))
	}
	return fmt.Sprintf(
		"BT\n%d %d Td %d TL \n%s %d. Tf\n%s\nET\n",
		x, y, fontSize, fontName, fontSize, strings.Join(opes, " T*\n"))
}

// cidFontSubType2 is a Type 2 CIDFont.
// A Type 2 CIDFont contains glyph descriptions based on the TrueType font format.
// A TrueType font program contains a “cmap” tables for predefined encoding.
// it provides mappings directly from character codes to glyph indices.
//
// Even though the CIDs are sometimes not used to select glyphs in a Type 2 CIDFont,
// they are always used to determine the glyph metrics, as described in the next section.
type cidFontSubType2 struct {
	abstractCIDFont
	gidMap      map[int32]uint16
	newGIDMap   map[uint16]uint16
	embededFont *opentype.Font
}

func (f *cidFontSubType2) SubType() string {
	return "/CIDFontType2"
}

func (f *cidFontSubType2) walk(walker func(obj pdfObject)) {
	walker(f.fontDescriptor.fontFile2)
}

func (f *cidFontSubType2) parentBaseFont(cmapName string) string {
	return f.baseFont
}

func (f *cidFontSubType2) compile() string {
	return f.compileHelper(f.SubType())
}

func (f *cidFontSubType2) build() error {
	list := make([]uint16, len(f.newGIDMap))
	for base, new := range f.newGIDMap {
		list[new] = base
	}
	newFont, err := f.embededFont.FilterGlyf(list)
	if err != nil {
		return err
	}
	// TODO
	// w/dw builing
	unitsPerEm := int(newFont.Head.UnitsPerEm)
	wArray := make([]string, len(newFont.Hmtx.HMetrics))
	for i, hm := range newFont.Hmtx.HMetrics {
		wArray[i] = strconv.Itoa(int(hm.AdvanceWidth) * 1000 / unitsPerEm)
	}
	f.w = fmt.Sprintf("0 [%s]", strings.Join(wArray, " "))
	f.dw = int(newFont.Hmtx.HMetrics[len(wArray)-1].AdvanceWidth) * 1000 / unitsPerEm

	var buf bytes.Buffer
	if err := opentype.NewBuilder(f.embededFont.SfntVersion).WithTables(newFont.Tables()).Build(&buf); err != nil {
		return err
	}
	f.fontDescriptor.fontFile2.dict["/Length1"] = strconv.Itoa(buf.Len())
	f.fontDescriptor.fontFile2.addBinaryDatum(buf.Bytes())
	return nil
}

func (f *cidFontSubType2) createText(fontName string, x, y int, fontSize int, text string) string {
	texts := strings.Split(text, "\n")
	opes := make([]string, 0, len(texts))
	for _, t := range texts {
		str := ""
		for _, w := range utf16.Encode([]rune(t)) {
			fileGID := f.gidMap[int32(w)]
			newGID, ok := f.newGIDMap[fileGID]
			if !ok {
				newGID = uint16(len(f.newGIDMap))
				f.newGIDMap[fileGID] = newGID
			}
			str += fmt.Sprintf("%04X", newGID)
		}
		opes = append(opes, fmt.Sprintf("<%s> Tj", str))
	}
	return fmt.Sprintf(
		"BT\n%d %d Td %d TL \n%s %d. Tf\n%s\nET\n",
		x, y, fontSize, fontName, fontSize, strings.Join(opes, " T*\n"))
}

// CIDSystemInfo specifies the character collection.
// CIDFont and CMap dictionaries contain a CIDSystemInfo entry specifying the character collection assumed by the CIDFont associated with the CMap.
// Character collections whose Registry and Ordering values are the same are compatible.
type CIDSystemInfo struct {
	registry   string
	ordering   string
	supplement int
}

// newCIDSystemInfo returns a CID system infomation.
func newCIDSystemInfo(registry, ordering string, supplement int) *CIDSystemInfo {
	return &CIDSystemInfo{
		registry:   registry,
		ordering:   ordering,
		supplement: supplement,
	}
}

func (c *CIDSystemInfo) String() string {
	return fmt.Sprintf("%s-%s-%d", c.registry, c.ordering, c.supplement)
}

func (c *CIDSystemInfo) compile() string {
	return fmt.Sprintf(
		"<</Registry (%s) /Ordering (%s) /Supplement %d>>",
		c.registry, c.ordering, c.supplement)
}

// FontDescriptor specifies metrics and other attributes of a simple font or a CIDFont as a whole, as distinct from the metrics of individual glyphs.
type FontDescriptor struct {
	fontName    string
	flags       int
	fontBBox    *Box
	italicAngle int
	ascent      int
	descent     int
	capHeight   int
	stemV       int
	fontFile2   *stream
}

// NewFontDescriptor creates a Font
func NewFontDescriptor(fontName string, flags int, fontBBox *Box, italicAngle, ascent, descent, capHeight, stemV int) *FontDescriptor {
	return &FontDescriptor{fontName, flags, fontBBox, italicAngle, ascent, descent, capHeight, stemV, nil}
}

func (f *FontDescriptor) compile() string {
	options := make([]string, 0)
	if f.fontFile2 != nil {
		options = append(options, fmt.Sprintf("/FontFile2 %s", f.fontFile2.indirectReference()))
	}
	return fmt.Sprintf(
		"<</Type /FontDescriptor /FontName %s /Flags %d /FontBBox %s /ItalicAngle %d /Ascent %d /Descent %d /CapHeight %d /StemV %d %s>>",
		f.fontName, f.flags, f.fontBBox.compile(), f.italicAngle, f.ascent, f.descent, f.capHeight, f.stemV, strings.Join(options, " "))
}
