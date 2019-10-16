//go:binary-only-package

package aper

import (
	"strconv"
	"strings"
)

// fieldParameters is the parsed representation of tag string from a structure field.
type fieldParameters struct {
	optional            bool   // true iff the type has OPTIONAL tag.
	sizeExtensible      bool   // true iff the size can be extensed.
	valueExtensible     bool   // true iff the value can be extensed.
	sizeLowerBound      *int64 // a sizeLowerBound is the minimum size of type constraint(maybe nil).
	sizeUpperBound      *int64 // a sizeUpperBound is the maximum size of type constraint(maybe nil).
	valueLowerBound     *int64 // a valueLowerBound is the minimum value of type constraint(maybe nil).
	valueUpperBound     *int64 // a valueUpperBound is the maximum value of type constraint(maybe nil).
	defaultValue        *int64 // a default value for INTEGER and ENUMERATED typed fields (maybe nil).
	openType            bool   // true iff this type is opentype.
	referenceFieldName  string // the field to get to get the corresrponding value of this type(maybe nil).
	referenceFieldValue *int64 // the field value which map to this type(maybe nil).
}

// Given a tag string with the format specified in the //go:binary-only-package

package comment,
// parseFieldParameters will parse it into a fieldParameters structure,
// ignoring unknown parts of the string. TODO:PrintableString
func parseFieldParameters(str string) (params fieldParameters) {}
