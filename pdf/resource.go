package pdf

import (
	"fmt"
	"strings"
)

// resource is a resource of pdf page.
type resource struct {
	objectIdentifier
	font    map[string]Font
	xobject map[string]*stream
}

func newResource() *resource {
	return &resource{
		objectIdentifier: objectIdentifier{},
		font:             make(map[string]Font),
		xobject:          make(map[string]*stream),
	}
}

// addFont adds the font to the page resource.
func (r *resource) addFont(f Font) {
	r.font[f.resourceName()] = f
}

// addImage adds the image to the page resource.
func (r *resource) addImage(i *ImageResource) {
	r.xobject[i.name] = i.asStream()
}

func (r *resource) compile() string {
	fonts := make([]string, 0, len(r.font))
	for k, f := range r.font {
		fonts = append(fonts, fmt.Sprintf("%s %s", k, f.indirectReference()))
	}
	xobjects := make([]string, 0, len(r.xobject))
	for k, xo := range r.xobject {
		xobjects = append(xobjects, fmt.Sprintf("%s %s", k, xo.indirectReference()))
	}
	dict := make([]string, 0, 2)
	if len(fonts) > 0 {
		dict = append(dict, fmt.Sprintf("/Font <<%s>>", strings.Join(fonts, " ")))
	}
	if len(xobjects) > 0 {
		dict = append(dict, fmt.Sprintf("/XObject <<%s>>", strings.Join(xobjects, " ")))
	}
	return r.bracket(fmt.Sprintf("<<%s>>", strings.Join(dict, " ")))
}

func (r *resource) walk(walker func(obj pdfObject)) {
	if r != nil {
		walker(r)
		for _, f := range r.font {
			f.walk(walker)
		}
		for _, xo := range r.xobject {
			walker(xo)
		}
	}
}
