package pdf

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"strconv"
	"strings"
)

// stream is a PDF stream object whose data consists of a sequence of instructions.
type stream struct {
	objectIdentifier
	dict   map[string]string
	filter streamFilter
	data   [][]byte
}

// newDeflatedStream creates Stream that is deflated.
func newDeflatedStream() *stream {
	return newStream(newDeflateEncoder())
}

// newFlatStream creates Stream that is not compressed.
func newFlatStream() *stream {
	return newStream(newFlatEncoder())
}

func newStream(f streamFilter) *stream {
	return &stream{
		objectIdentifier: objectIdentifier{},
		dict:             make(map[string]string),
		filter:           f,
		data:             make([][]byte, 0),
	}
}

func (s *stream) addStringDatum(datum string) {
	s.data = append(s.data, []byte(datum))
}

func (s *stream) addBinaryDatum(datum []byte) {
	s.data = append(s.data, datum)
}

func (s *stream) compile() (res []byte, err error) {
	if s.filter.name() != "" {
		s.dict["/Filter"] = s.filter.name()
	}
	data, err := s.filter.compress(s.data)
	if err != nil {
		return
	}
	s.dict["/Length"] = strconv.Itoa(len(data))
	dict := s.dict2pdf()
	b := bytes.NewBuffer(make([]byte, 0, len(data)+len(dict)+100))
	fmt.Fprintf(b, "%d %d obj\n<<%s>>\nstream\n", s.objectNumber, s.generationNumber, dict)
	b.Write(data)
	fmt.Fprintln(b, "\nendstream\nendobj")
	res = b.Bytes()
	return
}

// dict2pdf creates the pdf expression of the stream object dictionary.
func (s *stream) dict2pdf() string {
	dict := make([]string, 0, len(s.dict))
	for k, v := range s.dict {
		dict = append(dict, fmt.Sprintf("%s %s", k, v))
	}
	return strings.Join(dict, " ")
}

type streamFilter interface {
	name() string
	compress([][]byte) ([]byte, error)
}

type flatEncoder struct{}

func newFlatEncoder() *flatEncoder {
	return &flatEncoder{}
}

func (e *flatEncoder) name() string {
	return ""
}

func (e *flatEncoder) compress(data [][]byte) ([]byte, error) {
	size := 0
	for _, datum := range data {
		size += len(datum)
	}
	res := make([]byte, 0, size)
	for _, datum := range data {
		res = append(res, datum...)
	}
	return res, nil
}

type deflateEncoder struct{}

func newDeflateEncoder() *deflateEncoder {
	return &deflateEncoder{}
}

func (e *deflateEncoder) name() string {
	return "/FlateDecode"
}

func (e *deflateEncoder) compress(data [][]byte) ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	defer w.Close()
	for _, datum := range data {
		if _, err := w.Write(datum); err != nil {
			return []byte{}, err
		}
	}
	if err := w.Close(); err != nil {
		return []byte{}, err
	}
	return b.Bytes(), nil
}
