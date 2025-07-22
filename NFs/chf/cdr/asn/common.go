package asn

import (
	"strconv"
	"strings"
)

// ASN.1 universal tag number
const (
	TagBoolean         = 1
	TagInteger         = 2
	TagBitString       = 3
	TagOctetString     = 4
	TagNull            = 5
	TagOID             = 6
	TagEnumerated      = 10
	TagUTF8String      = 12
	TagSequence        = 16
	TagSet             = 17
	TagNumericString   = 18
	TagPrintableString = 19
	TagT61String       = 20
	TagIA5String       = 22
	TagUTCTime         = 23
	TagGeneralizedTime = 24
	TagGraphicString   = 25
	TagGeneralString   = 27
	TagBMPString       = 30
)

// ASN.1 tag class
const (
	ClassUniversal       = 0
	ClassApplication     = 1
	ClassContextSpecific = 2
	ClassPrivate         = 3
)

type tagAndLen struct {
	class       int
	constructed bool
	tagNumber   uint64
	len         int64
}

// fieldParameters is the parsed representation of tag string from a structure field.
type fieldParameters struct {
	optional            bool    // true iff the type has OPTIONAL tag.
	sizeLowerBound      *int64  // a sizeLowerBound is the minimum size of type constraint(maybe nil).
	sizeUpperBound      *int64  // a sizeUpperBound is the maximum size of type constraint(maybe nil).
	valueLowerBound     *int64  // a valueLowerBound is the minimum value of type constraint(maybe nil).
	valueUpperBound     *int64  // a valueUpperBound is the maximum value of type constraint(maybe nil).
	defaultValue        *int64  // a default value for INTEGER and ENUMERATED typed fields (maybe nil).
	openType            bool    // true iff this type is opentype.
	referenceFieldName  string  // the field to get to get the corresrponding value of this type(maybe nil).
	referenceFieldValue *int64  // the field value which map to this type(maybe nil).
	tagNumber           *uint64 // the field is for ber struct type
	explicitTag         bool    // true iff the tag need to explicit encoded
	set                 bool    // true iff ASN.1 type is set
	choice              bool    // true iff ASN.1 type is choice
	stringType          int
	null                bool // true iff ASN.1 type is null
}

// Given a tag string with the format specified in the package comment,
// parseFieldParameters will parse it into a fieldParameters structure,
// ignoring unknown parts of the string. TODO:PrintableString
func parseFieldParameters(str string) (params fieldParameters) {
	for _, part := range strings.Split(str, ",") {
		switch {
		case part == "optional":
			params.optional = true
		case strings.HasPrefix(part, "sizeLB:"):
			i, err := strconv.ParseInt(part[7:], 10, 64)
			if err == nil {
				params.sizeLowerBound = new(int64)
				*params.sizeLowerBound = i
			}
		case strings.HasPrefix(part, "sizeUB:"):
			i, err := strconv.ParseInt(part[7:], 10, 64)
			if err == nil {
				params.sizeUpperBound = new(int64)
				*params.sizeUpperBound = i
			}
		case strings.HasPrefix(part, "valueLB:"):
			i, err := strconv.ParseInt(part[8:], 10, 64)
			if err == nil {
				params.valueLowerBound = new(int64)
				*params.valueLowerBound = i
			}
		case strings.HasPrefix(part, "valueUB:"):
			i, err := strconv.ParseInt(part[8:], 10, 64)
			if err == nil {
				params.valueUpperBound = new(int64)
				*params.valueUpperBound = i
			}
		case strings.HasPrefix(part, "default:"):
			i, err := strconv.ParseInt(part[8:], 10, 64)
			if err == nil {
				params.defaultValue = new(int64)
				*params.defaultValue = i
			}
		case part == "openType":
			params.openType = true
		case strings.HasPrefix(part, "referenceFieldName:"):
			params.referenceFieldName = part[19:]
		case strings.HasPrefix(part, "referenceFieldValue:"):
			i, err := strconv.ParseInt(part[20:], 10, 64)
			if err == nil {
				params.referenceFieldValue = new(int64)
				*params.referenceFieldValue = i
			}
		case strings.HasPrefix(part, "tagNum:"):
			i, err := strconv.ParseInt(part[7:], 10, 64)
			if err == nil {
				params.tagNumber = new(uint64)
				*params.tagNumber = uint64(i)
			}
		case part == "explicit":
			params.explicitTag = true
		case part == "set":
			params.set = true
		case part == "choice":
			params.choice = true
		case part == "utf8":
			params.stringType = TagUTF8String
		case part == "ia5":
			params.stringType = TagIA5String
		case part == "graphic":
			params.stringType = TagGraphicString
		case part == "null":
			params.null = true
		}
	}
	return params
}
