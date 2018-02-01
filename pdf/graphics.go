package pdf

import (
	"fmt"
	"strings"
)

// GraphicsObject is a graphics object in pdf.
//
// ex. Line, Rectangle
//
// Clients of GraphicsObject must call Render() after methods that set graphics style are called.
// If Render() is called, GraphicsObject ignores any method calls.
type GraphicsObject interface {
	// Render GraphicsObject on a pdf page.
	Render()
	// render returns a string expression of a graphics object.
	render(cb *Box) string
}

type graphicsObject struct {
	page     Page
	rendered bool
}

func (obj *graphicsObject) ifNotRendered(fn func()) {
	if !obj.rendered {
		fn()
	}
}

// renderSelf is a utility method for implementing Render().
func (obj *graphicsObject) renderSelf(self GraphicsObject) {
	obj.ifNotRendered(func() {
		obj.page.render(self)
		obj.rendered = true
	})
}

// Line is a graphics object that represents a line.
//
// See the comment of GraphicsObject.
type Line interface {
	// MoveTo render a line to the point.
	MoveTo(x, y int) Line
	// Color specifies the color of the line.
	Color(c Color) Line
	// DashPattern controls the pattern of dashes and gaps used to stroke paths.
	DashPattern(dash, gap, phase int) Line
	// CapStyle is the shape to be used at the ends of open subpaths.
	CapStyle(lcs LineCapStyle) Line
	// JoinStyle specifies the shape to be used at the corners of paths that are stroked.
	JoinStyle(ljs LineJoinStyle) Line
	GraphicsObject
}

// line is a implementation of a Line interface.
type line struct {
	graphicsObject
	startX    int
	startY    int
	strokes   []int
	lineWidth int
	endShape  int
	color     Color
	ldp       *lineDashPattern
	lcs       LineCapStyle
	ljs       LineJoinStyle
}

// newLine returns a new Line.
func newLine(page Page, startX, startY, endX, endY, lineWidth int) *line {
	strokes := make([]int, 2)
	strokes[0] = endX
	strokes[1] = endY
	return &line{
		graphicsObject: graphicsObject{
			page: page,
		},
		startX:    startX,
		startY:    startY,
		strokes:   strokes,
		color:     newColorUndef(),
		lineWidth: lineWidth,
		lcs:       lineCapStyleUndefined,
		ljs:       lineJoinStyleUndefined,
	}
}

func (l *line) MoveTo(x, y int) Line {
	l.ifNotRendered(func() {
		l.strokes = append(l.strokes, x, y)
	})
	return l
}

func (l *line) Color(c Color) Line {
	l.ifNotRendered(func() {
		l.color = c
	})
	return l
}
func (l *line) DashPattern(dash, gap, phase int) Line {
	l.ifNotRendered(func() {
		l.ldp = newLineDashPattern(dash, gap, phase)
	})
	return l
}

func (l *line) CapStyle(lcs LineCapStyle) Line {
	l.ifNotRendered(func() {
		l.lcs = lcs
	})
	return l
}

func (l *line) JoinStyle(ljs LineJoinStyle) Line {
	l.ifNotRendered(func() {
		l.ljs = ljs
	})
	return l
}

func (l *line) Render() {
	l.renderSelf(l)
}

func (l *line) render(cb *Box) string {
	style := make([]string, 0, 4)
	if l.color.colorSpace() != colorSpaceUndefined {
		style = append(style, l.color.strokeColor())
	}
	if l.ldp != nil {
		style = append(style, l.ldp.compile())
	}
	if l.lcs != lineCapStyleUndefined {
		style = append(style, l.lcs.compile())
	}
	if l.ljs != lineJoinStyleUndefined {
		style = append(style, l.ljs.compile())
	}
	strokes := make([]string, len(l.strokes)/2)
	for i := 0; i < len(strokes); i++ {
		strokes[i] = fmt.Sprintf("%d %d l", cb.leftBottomX+l.strokes[2*i], cb.rightTopY-l.strokes[2*i+1])
	}
	return fmt.Sprintf(
		"q %d %d m %s %d w %s S Q\n",
		cb.leftBottomX+l.startX, cb.rightTopY-l.startY, strings.Join(strokes, " "), l.lineWidth, strings.Join(style, " "))
}

// lineDashPattern controls the pattern of dashes and gaps used to stroke paths.
type lineDashPattern struct {
	dash  int
	gap   int
	phase int
}

// newLineDashPattern returns a new lineDashPattern
func newLineDashPattern(dash, gap, phase int) *lineDashPattern {
	return &lineDashPattern{dash, gap, phase}
}

func (ldp *lineDashPattern) compile() string {
	if ldp.dash != ldp.gap {
		return fmt.Sprintf("[%d %d] %d d", ldp.dash, ldp.gap, ldp.phase)
	}
	return fmt.Sprintf("[%d] %d d", ldp.dash, ldp.phase)
}

// LineCapStyle is the shape to be used at the ends of open subpaths.
type LineCapStyle int

const (
	// lineCapStyleUndefinedCap is the default LineCapStyle.
	lineCapStyleUndefined LineCapStyle = iota
	// LineCapStyleButt : The stroke is squared off at the endpoint of the path. There is no projection beyond the end of the path.
	LineCapStyleButt
	// LineCapStyleRound : A semicircular arc with a diameter equal to the line width is drawn around the endpoint and filled in.
	LineCapStyleRound
	// LineCapStyleProjectingSquare : The stroke continues beyond the endpoint of the path for a distance equal to half the line width and is squared off.
	LineCapStyleProjectingSquare
)

func (lcs LineCapStyle) compile() string {
	switch lcs {
	case LineCapStyleButt:
		return "0 J"
	case LineCapStyleRound:
		return "1 J"
	case LineCapStyleProjectingSquare:
		return "2 J"
	default:
		return ""
	}
}

// LineJoinStyle specifies the shape to be used at the corners of paths that are stroked.
type LineJoinStyle int

const (
	// lineJoinStyleUndefinedJoin is the default LineJoinStyle.
	lineJoinStyleUndefined LineJoinStyle = iota
	// LineJoinStyleMiter : The outer edges of the strokes for the two segments are extended until they meet at an angle.
	LineJoinStyleMiter
	// LineJoinStyleRound : An arc of a circle with a diameter equal to the line width is drawn around the point where the two segments meet, connecting the outer edges of the strokes for the two segments.
	LineJoinStyleRound
	// LineJoinStyleBevel : The two segments are finished with butt caps and the resulting notch beyond the ends of the segments is filled with a triangle.
	LineJoinStyleBevel
)

func (ljs LineJoinStyle) compile() string {
	switch ljs {
	case LineJoinStyleMiter:
		return "0 j"
	case LineJoinStyleRound:
		return "1 j"
	case LineJoinStyleBevel:
		return "2 j"
	default:
		return ""
	}

}

// Rectangle is a graphics object that represents a rectable.
//
// See the comment of GraphicsObject.
type Rectangle interface {
	// StrokeColor specifies the color of edge.
	StrokeColor(c Color) Rectangle
	// FillColor specifies the color of surface.
	FillColor(c Color) Rectangle
	GraphicsObject
}

// rectangle is a implementation of a Rectangle interface.
type rectangle struct {
	graphicsObject
	startX      int
	startY      int
	width       int
	height      int
	strokeColor Color
	fillColor   Color
}

// newRectangle returns a rectangle
func newRectangle(page Page, startX, startY, width, height int) *rectangle {
	return &rectangle{
		graphicsObject: graphicsObject{
			page: page,
		},
		startX:      startX,
		startY:      startY,
		width:       width,
		height:      height,
		strokeColor: newColorUndef(),
		fillColor:   newColorUndef(),
	}
}

func (r *rectangle) StrokeColor(sc Color) Rectangle {
	r.ifNotRendered(func() {
		r.strokeColor = sc
	})
	return r
}

func (r *rectangle) FillColor(fc Color) Rectangle {
	r.ifNotRendered(func() {
		r.fillColor = fc
	})
	return r
}

func (r *rectangle) Render() {
	r.renderSelf(r)
}

func (r *rectangle) render(cb *Box) string {
	flag := 0
	colors := make([]string, 0, 2)
	if r.strokeColor.colorSpace() != colorSpaceUndefined {
		colors = append(colors, r.strokeColor.strokeColor())
		flag += 2
	}
	if r.fillColor.colorSpace() != colorSpaceUndefined {
		colors = append(colors, r.fillColor.nonStrokeColor())
		flag++
	}
	var draw string
	switch flag {
	case 3:
		// B: Fill and then stroke the path.
		draw = "B"
	case 2:
		// S: Stroke the path.
		draw = "S"
	case 1:
		// f: Fill the path.
		draw = "f"
	default:
		// n: End the path object without filling or stroking it.
		draw = "n"
	}
	return fmt.Sprintf(
		"q %s %d %d %d %d re %s Q\n",
		strings.Join(colors, " "), cb.leftBottomX+r.startX, cb.rightTopY-r.startY-r.height, r.width, r.height, draw)
}
