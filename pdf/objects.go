package pdf

import (
	"fmt"
)

// pdfObject has a common functionality of pdf objects.
type pdfObject interface {
	// refNo returns the object number.
	refNo() int
	// age returns the age number.
	age() int
	// indirectReference retunrs the expression for indirect refernce.
	indirectReference() string
	// number sets the own object number by using object number counter.
	number(c *counter)
}

// stringCompiler compiles itself and returns string.
type stringCompiler interface {
	compile() string
}

type stringObject interface {
	pdfObject
	stringCompiler
}

// binaryCompiler compiles itself and returns byte array.
type binaryCompiler interface {
	compile() ([]byte, error)
}

type binaryObject interface {
	pdfObject
	binaryCompiler
}

// traversableObject returns itself and descendant objects recursively.
// walker function traverse it.
type traversableObject interface {
	// traversableObject call walker each time it find descendant objects.
	walk(walker func(obj pdfObject))
}

// objectIdentifier is a minimal pdfObject.
type objectIdentifier struct {
	objectNumber     int
	generationNumber int
}

func (obj *objectIdentifier) refNo() int {
	return obj.objectNumber
}

func (obj *objectIdentifier) age() int {
	return obj.generationNumber
}

func (obj *objectIdentifier) indirectReference() string {
	return fmt.Sprintf("%d %d R", obj.objectNumber, obj.generationNumber)
}

func (obj *objectIdentifier) number(c *counter) {
	// If obj has not been numbered yet, do.
	if 0 == obj.objectNumber {
		obj.objectNumber = c.next()
		obj.generationNumber = 0
	}
}

// bracket is a utility function to compile own object expression.
func (obj *objectIdentifier) bracket(value string) string {
	return fmt.Sprintf("%d %d obj\n%s\nendobj\n", obj.objectNumber, obj.generationNumber, value)
}

// counter is a object number counter.
type counter struct {
	count int
}

// newCounter returns object number counter that starts with 3.
// 1: document catalog
// 2: root page
func newCounter() *counter {
	return &counter{2}
}

// next returns next object number
func (c *counter) next() int {
	c.count++
	return c.count
}
