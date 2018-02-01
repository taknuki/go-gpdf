package pdf

import (
	"fmt"
	"io"
)

// writer is a pdf writer.
type writer struct {
	w      io.Writer
	offset int
	crt    *crossRefTable
	err    error
}

func newWriter(w io.Writer) *writer {
	return &writer{
		w:      w,
		crt:    newCrossRefTable(),
		offset: 0,
		err:    nil,
	}
}

func (w *writer) hasError() bool {
	if w.err != nil {
		return true
	}
	return false
}

// start writes pdf file header.
func (w *writer) start(version string) *writer {
	// The first line of a PDF file is a header
	// identifying the version of the PDF specification
	// to which the file conforms.
	w.writeStr(fmt.Sprintf("%%PDF-%s\n", version))
	// It is recommended that the header line
	// be immediately followed by a comment line
	// containing at least four binary characters whose codes are 128 or greater.
	// This ensures proper behavior of file transfer applications
	// that inspect data near the beginning of a file to determine
	// whether to treat the file’s contents as text or as binary.
	w.write([]byte{'%', 0x80, 0x81, 0x82, 0x83, '\n'})
	return w
}

// finish writes cross reference table and trailer.
func (w *writer) finish(root objectIdentifier) error {
	startXRef := w.offset
	w.writeStr(w.crt.compile())
	t := newTrailer(root, len(w.crt.entries), startXRef)
	w.writeStr(t.compile())
	return w.err
}

func (w *writer) writeTraversable(traversable traversableObject) *writer {
	if !w.hasError() {
		traversable.walk(func(obj pdfObject) {
			w.writeObj(obj)
		})
	}
	return w
}

// writeObj writes byte expression of pdf object.
func (w *writer) writeObj(obj pdfObject) {
	if !w.crt.hasEntry(obj) {
		w.crt.addNewEntry(obj, w.offset)
		switch o := obj.(type) {
		case stringObject:
			w.writeStr(o.compile())
		case binaryObject:
			data, err := o.compile()
			if err != nil {
				w.err = fmt.Errorf("failed to write binary object: %s", err)
			} else {
				w.write(data)
			}
		default:
			if w.err == nil {
				w.err = fmt.Errorf("unknwon type of pdfobjct, no:%d", obj.refNo())
			}
		}
	}
}

// write writes data.
func (w *writer) write(data []byte) {
	if !w.hasError() {
		n, err := w.w.Write(data)
		if err != nil {
			w.err = err
		} else {
			w.offset += n
		}
	}
}

// writeStr writes byte expression of string.
func (w *writer) writeStr(s string) {
	if !w.hasError() {
		n, err := io.WriteString(w.w, s)
		if err != nil {
			w.err = err
		} else {
			w.offset += n
		}
	}
}
