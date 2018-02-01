package pdf

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"
	"testing"
)

func isInitialStream(t *testing.T, s *stream) {
	t.Helper()
	if s.refNo() != 0 {
		t.Error("stream is not initial state: refNo != 0")
	}
	if s.age() != 0 {
		t.Error("stream is not initial state: age != 0")
	}
	if len(s.dict) > 0 {
		t.Error("stream is not initial state: dict is not empty")
	}
	if len(s.data) > 0 {
		t.Error("stream is not initial state: data is not empty")
	}
}

func TestFlatStream(t *testing.T) {
	s := newFlatStream()
	isInitialStream(t, s)
	if _, ok := s.filter.(*flatEncoder); !ok {
		t.Error("flatEncoder is not used.")
	}
	s.addStringDatum("abc")
	if string(s.data[0]) != "abc" {
		t.Error("addStringDatum is failed")
	}
	s.addBinaryDatum([]byte("def"))
	if string(s.data[1]) != "def" {
		t.Error("addBinaryDatum is failed")
	}
	s.dict["/G"] = "g"
	s.dict["/H"] = "h"
	s.dict["/I"] = "i"
	acDict := s.dict2pdf()
	tmpDict := strings.Split(acDict, " ")
	sort.Sort(sort.StringSlice(tmpDict))
	sortedDict := strings.Join(tmpDict, " ")
	exDict := "/G /H /I g h i"
	if exDict != sortedDict {
		t.Errorf("dict2pdf: expected:/G g /H h /I i actual:%s", acDict)
	}
	res, err := s.compile()
	if err != nil {
		t.Errorf("compillation failed: unexpected error:%s", err)
	}
	expected1 := "0 0 obj"
	expected3 := `stream
abcdef
endstream
endobj
`

	ptn := regexp.MustCompile(`([0-9a-z\s]*)\n<<(.*)>>\n([0-9a-z\s]*)`)
	actual := ptn.FindStringSubmatch(string(res))

	testCompillation(t, expected1, actual[1])
	testCompillation(t, expected3, actual[3])
	tmpAc2 := strings.Split(actual[2], " ")
	sort.Sort(sort.StringSlice(tmpAc2))
	sortedAc2 := strings.Join(tmpAc2, " ")
	expected2 := "/G /H /I /Length 6 g h i"
	testCompillation(t, expected2, sortedAc2)
}

func TestDeflatedStream(t *testing.T) {
	s := newDeflatedStream()
	isInitialStream(t, s)
	if _, ok := s.filter.(*deflateEncoder); !ok {
		t.Error("deflateEncoder is not used.")
	}
	s.addStringDatum("abc")
	res, err := s.compile()
	if err != nil {
		t.Errorf("compillation failed: unexpected error:%s", err)
	}
	acFilter := s.dict["/Filter"]
	exFilter := "/FlateDecode"
	if exFilter != acFilter {
		t.Errorf("filter name: expected:%s actual:%s", exFilter, acFilter)
	}

	ptn := regexp.MustCompile(`([0-9a-z\s]*)\n<<(.*)>>\nstream\n(.*)\nendstream\nendobj\n`)
	actual := ptn.FindStringSubmatch(string(res))
	expected1 := "0 0 obj"
	expected3 := "abc"
	reader, err := zlib.NewReader(bytes.NewReader([]byte(actual[3])))
	if err != nil {
		t.Fatalf("probably test error:%s", err)
	}
	actual3, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatalf("probably test error:%s", err)
	}
	testCompillation(t, expected1, actual[1])
	testCompillation(t, expected3, string(actual3))
}
