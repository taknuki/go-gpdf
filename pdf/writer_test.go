package pdf

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"
)

type mockWriter struct {
	w   io.Writer
	err error
}

func (w *mockWriter) Write(p []byte) (n int, err error) {
	if w.err == nil {
		return w.w.Write(p)
	}
	return 0, w.err
}

func TestWriterStart(t *testing.T) {
	version := "1.7"
	b := bytes.NewBuffer([]byte{})
	w := newWriter(b)
	if w.hasError() {
		t.Error("writer is not initial state?: hasError")
	}
	w.start(version)
	expected := []byte{37, 80, 68, 70, 45, 49, 46, 55, 10, 37, 128, 129, 130, 131, 10}
	actual := b.Bytes()
	if bytes.Compare(expected, actual) != 0 {
		t.Errorf("start: written bytes are unexpected\nexpected:%b\nactual  :%b", expected, actual)
	}
	if w.offset != len(actual) {
		t.Errorf("start: offset is not expected:%d actual:%d", len(actual), w.offset)
	}
}

func TestWriterFinish(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	w := newWriter(b)
	w.offset = 100
	root := objectIdentifier{1, 0}
	trailer := newTrailer(root, len(w.crt.entries), 100)
	w.finish(root)
	expected := fmt.Sprintf("%s%s", w.crt.compile(), trailer.compile())
	actual := string(b.Bytes())
	if expected != actual {
		t.Errorf("written string is unexpected\nexpected:%s\nactual  :%s\n", expected, actual)
	}
}

func testWriterResult(t *testing.T, expected string, offset int, b *bytes.Buffer) {
	t.Helper()
	actual := string(b.Bytes())
	if expected != actual {
		t.Errorf("written string is unexpected\nexpected:%s\nactual  :%s\n", expected, actual)
	}
	len := b.Len()
	if offset != len {
		t.Errorf("offset is unexpected\nexpected:%d\nactual  :%d\n", len, offset)
	}
}

func TestWriter1(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	w := newWriter(b)
	obj := &mockPDFObject{
		objectIdentifier: objectIdentifier{11, 0},
		res:              "first",
		child1: &mockPDFObject{
			objectIdentifier: objectIdentifier{12, 0},
			res:              "second",
		},
		child2: &mockBinaryObject{
			objectIdentifier: objectIdentifier{13, 0},
			res:              []byte("third"),
		},
	}
	w.writeTraversable(obj)
	testWriterResult(t, "firstsecondthird", w.offset, b)
	w.writeTraversable(obj)
	testWriterResult(t, "firstsecondthird", w.offset, b)
}

func testWriterError(t *testing.T, b *bytes.Buffer, w *writer) {
	t.Helper()
	if !w.hasError() {
		t.Error("writer should have error")
	}
	if b.Len() > 0 || w.offset > 0 {
		t.Error("is illegal data written?")
	}
}

func TestWriterError1(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	w := newWriter(b)
	obj1 := &mockBinaryObject{
		objectIdentifier: objectIdentifier{1, 0},
		res:              []byte("first"),
		err:              errors.New("error"),
	}
	w.writeObj(obj1)
	testWriterError(t, b, w)
	obj2 := &objectIdentifier{2, 0}
	w.err = nil
	w.writeObj(obj2)
	testWriterError(t, b, w)
}

func TestWriterError2(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	w := newWriter(&mockWriter{
		w:   b,
		err: errors.New("error"),
	})
	w.write([]byte{1})
	testWriterError(t, b, w)
	w.err = nil
	w.writeStr("a")
	testWriterError(t, b, w)
}
