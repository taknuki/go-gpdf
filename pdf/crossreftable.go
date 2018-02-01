package pdf

import (
	"fmt"
	"sort"
)

// crossRefTable is The cross-reference table.
// it contains information that permits random access to indirect objects within the file so that the entire file need not be read to locate any particular object.
type crossRefTable struct {
	entries map[int]*crossRefEntry
}

func newCrossRefTable() *crossRefTable {
	entries := make(map[int]*crossRefEntry)
	// The first entry in the table (object number 0) is always free
	// and has a generation number of 65,535.
	entries[0] = newCrossRefEntry(0, 65535, false)
	return &crossRefTable{entries}
}

func (crt *crossRefTable) addNewEntry(obj pdfObject, offset int) {
	crt.entries[obj.refNo()] = newCrossRefEntry(offset, obj.age(), true)
}

func (crt *crossRefTable) hasEntry(obj pdfObject) bool {
	_, ok := crt.entries[obj.refNo()]
	return ok
}

func (crt *crossRefTable) compile() string {
	nums := make([]int, 0, len(crt.entries))
	for i := range crt.entries {
		nums = append(nums, i)
	}
	sort.Sort(sort.IntSlice(nums))
	subSections := make([][]int, 0)
	curSec := 0
	for i, num := range nums {
		if i == 0 {
			subSections = append(subSections, make([]int, 1))
			subSections[0][0] = 0
		} else {
			if num > nums[i-1]+1 {
				subSections = append(subSections, make([]int, 1))
				curSec++
				subSections[curSec][0] = num
			} else {
				subSections[curSec] = append(subSections[curSec], num)
			}
		}
	}
	result := "xref\n"
	for _, section := range subSections {
		t := make([]byte, 0, 20*len(section))
		for _, num := range section {
			t = append(t, crt.entries[num].compile()...)
		}
		result += fmt.Sprintf("%d %d\n%s", section[0], len(section), string(t))
	}
	return result
}

// crossRefEntry is a entry of a crossRefTable.
type crossRefEntry struct {
	offset int
	age    int
	inuse  bool
}

func newCrossRefEntry(offset int, age int, inuse bool) *crossRefEntry {
	return &crossRefEntry{offset, age, inuse}
}

func (cre *crossRefEntry) compile() string {
	u := "n"
	if !cre.inuse {
		u = "f"
	}
	return fmt.Sprintf("%010d %05d %s \n", cre.offset, cre.age, u)
}
