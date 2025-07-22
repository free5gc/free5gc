package asn

import (
	"fmt"
	"reflect"
)

// parse and return tag and length, also the length of two parts
func parseTagAndLength(bytes []byte) (r tagAndLen, off int, e error) {
	off++
	r.class = int(bytes[0] >> 6)
	r.constructed = (bytes[0] & 0x20) != 0
	if bytes[0]&0x1f != 0x1f {
		r.tagNumber = uint64(bytes[0] & 0x1f)
	} else {
		for off < len(bytes) {
			r.tagNumber <<= 7
			r.tagNumber |= uint64(bytes[off] & 0x7f)
			off++
			if bytes[off-1]&0x80 == 0 {
				break
			}
		}
		if off > 10 {
			e = fmt.Errorf("tag number is too large")
			return r, off, e
		}
	}

	if off >= len(bytes) {
		e = fmt.Errorf("panic: bytes: %v, off: %v", bytes, off)
		return r, off, e
	}
	if bytes[off] <= 127 {
		r.len = int64(bytes[off])
		off++
	} else {
		len := int(bytes[off] & 0x7f)
		// fmt.Println("len", len)
		if len > 3 {
			e = fmt.Errorf("length is too large")
			return r, off, e
		}
		off++
		var val int64
		val, e = parseInt64(bytes[off : off+len])
		if e != nil {
			return r, off, e
		}
		// fmt.Println("bytes[off : off+len]", bytes[off : off+len], "val", val)

		r.len = int64(val)
		off += len
	}

	return r, off, e
}

func parseBitString(bytes []byte) (r BitString, e error) {
	r.BitLength = uint64((len(bytes)-1)*8 - int(bytes[0]))
	r.Bytes = bytes[1:]
	return
}

func parseInt64(bytes []byte) (r int64, e error) {
	if len(bytes) > 8 {
		e = fmt.Errorf("out of range of int64")
		return r, e
	}

	for _, b := range bytes {
		r <<= 8
		r |= int64(b)
	}

	return r, e
}

func parseBool(b byte) (bool, error) {
	return b != 0, nil
}

// ParseField is the main parsing function. Given a byte slice containing type value,
// it will try to parse a suitable ASN.1 value out and store it
// in the given Value. TODO : ObjectIdenfier
func ParseField(v reflect.Value, bytes []byte, params fieldParameters) error {
	fieldType := v.Type()

	// If we have run out of data return error.
	if v.Kind() == reflect.Ptr {
		ptr := reflect.New(fieldType.Elem())
		v.Set(ptr)
		return ParseField(v.Elem(), bytes, params)
	}

	tal, talOff, err := parseTagAndLength(bytes)
	if err != nil {
		return err
	}
	if int64(talOff)+tal.len > int64(len(bytes)) {
		return fmt.Errorf("type value out of range")
	}

	// We deal with the structures defined in this package first.
	switch fieldType {
	case BitStringType:
		val, parse_err := parseBitString(bytes[talOff:])
		if parse_err != nil {
			return parse_err
		}

		v.Set(reflect.ValueOf(val))
		return nil
	case ObjectIdentifierType:
		return fmt.Errorf("unsppport ObjectIdenfier type")
	case OctetStringType:
		val := bytes[talOff:]
		v.Set(reflect.ValueOf(val))
		return nil
	case EnumeratedType:
		val, parse_err := parseInt64(bytes[talOff:])
		if parse_err != nil {
			return parse_err
		}

		v.Set(reflect.ValueOf(Enumerated(val)))
		return nil
	case NullType:
		val := true
		v.Set(reflect.ValueOf(val))
		return nil
	}
	switch val := v; val.Kind() {
	case reflect.Bool:
		if parsedBool, parse_err := parseBool(bytes[talOff]); parse_err != nil {
			return parse_err
		} else {
			val.SetBool(parsedBool)
			return nil
		}
	case reflect.Int, reflect.Int32, reflect.Int64:
		if parsedInt, parse_err := parseInt64(bytes[talOff:]); parse_err != nil {
			return parse_err
		} else {
			val.SetInt(parsedInt)
			return nil
		}
	case reflect.Struct:

		structType := fieldType
		var structParams []fieldParameters

		if structType.Field(0).Name == "Value" {
			// Non struct type
			// fmt.Println("Non struct type")
			return ParseField(val.Field(0), bytes, params)
		} else if structType.Field(0).Name == "List" {
			// List Type: SEQUENCE/SET OF
			// fmt.Println("List type")
			return ParseField(val.Field(0), bytes, params)
		}

		// parse parameters
		for i := 0; i < structType.NumField(); i++ {
			if structType.Field(i).PkgPath != "" {
				return fmt.Errorf("struct contains unexported fields : %s", structType.Field(i).PkgPath)
			}
			tempParams := parseFieldParameters(structType.Field(i).Tag.Get("ber"))
			structParams = append(structParams, tempParams)
		}

		// CHOICE or OpenType
		if structType.NumField() > 0 && structType.Field(0).Name == "Present" {
			present := 0

			if params.openType {
				return fmt.Errorf("openType is not implemented")
			} else {
				offset := 0
				// embed choice type
				if params.tagNumber != nil {
					tal, talOff, err = parseTagAndLength(bytes[talOff:])
					if err != nil {
						return err
					}
					if int64(talOff)+tal.len > int64(len(bytes)) {
						return fmt.Errorf("type value out of range")
					}
					offset += talOff
				}

				for i := 1; i < structType.NumField(); i++ {
					if structParams[i].tagNumber == nil {
						// TODO: choice type with a universal tag
					} else if *structParams[i].tagNumber == tal.tagNumber {
						present = i
						break
					}
				}
				val.Field(0).SetInt(int64(present))
				if present == 0 {
					return fmt.Errorf("CHOICE present is 0(present's field number)")
				} else if present >= structType.NumField() {
					return fmt.Errorf("CHOICE Present is bigger than number of struct field")
				} else {
					return ParseField(val.Field(present), bytes[offset:], structParams[present])
				}
			}
		}

		offset := int64(talOff)
		totalLen := int64(len(bytes))

		if !params.set {
			current := 0
			next := int64(0)
			for ; offset < totalLen; offset = next {
				talNow, talOffNow, parse_err := parseTagAndLength(bytes[offset:])
				if parse_err != nil {
					return parse_err
				}
				next = int64(offset) + int64(talOffNow) + talNow.len
				if next > totalLen {
					return fmt.Errorf("type value out of range")
				}
				if offset >= next {
					fmt.Println("bytes offset", offset, "next", next, "talOff", talOffNow, "tal.len", talNow.len)
					if offset > next {
						return fmt.Errorf("offset > next")
					}
				}

				for ; current < structType.NumField(); current++ {
					// for open type reference
					if params.openType {
						return fmt.Errorf("openType is not implemented")
					}
					if *structParams[current].tagNumber == talNow.tagNumber {
						if err = ParseField(val.Field(current), bytes[offset:next], structParams[current]); err != nil {
							return err
						}
						break
					}
				}
				if current >= structType.NumField() {
					return fmt.Errorf("corresponding type not found")
				}
				current++
			}
		} else {
			next := int64(0)
			for ; offset < totalLen; offset = next {
				talNow, talOffNow, parse_err := parseTagAndLength(bytes[offset:])
				if parse_err != nil {
					return parse_err
				}
				next = offset + int64(talOffNow) + talNow.len
				if next > totalLen {
					return fmt.Errorf("type value out of range")
				}

				current := 0
				for ; current < structType.NumField(); current++ {
					// for open type reference
					if params.openType {
						return fmt.Errorf("openType is not implemented")
					}
					if *structParams[current].tagNumber == talNow.tagNumber {
						if parse_err1 := ParseField(val.Field(current), bytes[offset:next], structParams[current]); parse_err1 != nil {
							return parse_err1
						}
						break
					}
				}
				if current >= structType.NumField() {
					return fmt.Errorf("corresponding type not found")
				}
			}
		}
		return nil
	case reflect.Slice:
		sliceType := fieldType
		var valArray [][]byte
		var next int64
		for offset := int64(talOff); offset < int64(len(bytes)); offset = next {
			talNow, talOffNow, errParse := parseTagAndLength(bytes[offset:])
			if errParse != nil {
				return errParse
			}
			next = offset + int64(talOffNow) + talNow.len
			if next > int64(len(bytes)) {
				return fmt.Errorf("type value out of range")
			}
			valArray = append(valArray, bytes[offset:next])
		}

		sliceLen := len(valArray)
		newSlice := reflect.MakeSlice(sliceType, sliceLen, sliceLen)
		for i := 0; i < sliceLen; i++ {
			errParse := ParseField(newSlice.Index(i), valArray[i], params)
			if errParse != nil {
				return errParse
			}
		}

		val.Set(newSlice)
		return nil
	case reflect.String:
		val.SetString(string(bytes[talOff:]))
		return nil
	}

	return fmt.Errorf("unsupported: %s", v.Type().String())
}

// Unmarshal parses the BER-encoded ASN.1 data structure b
// and uses the reflect package to fill in an arbitrary value pointed at by value.
// Because Unmarshal uses the reflect package, the structs
// being written to must use upper case field names.
//
// An ASN.1 INTEGER can be written to an int, int32, int64,
// If the encoded value does not fit in the Go type,
// Unmarshal returns a parse error.
//
// An ASN.1 BIT STRING can be written to a BitString.
//
// An ASN.1 OCTET STRING can be written to a []byte.
//
// An ASN.1 OBJECT IDENTIFIER can be written to an
// ObjectIdentifier.
//
// An ASN.1 ENUMERATED can be written to an Enumerated.
//
// Any of the above ASN.1 values can be written to an interface{}.
// The value stored in the interface has the corresponding Go type.
// For integers, that type is int64.
//
// An ASN.1 SEQUENCE OF x can be written
// to a slice if an x can be written to the slice's element type.
//
// An ASN.1 SEQUENCE can be written to a struct
// if each of the elements in the sequence can be
// written to the corresponding element in the struct.
//
// The following tags on struct fields have special meaning to Unmarshal:
// optional        	OPTIONAL tag in SEQUENCE
// sizeLB		        set the minimum value of size constraint
// sizeUB              set the maximum value of value constraint
// valueLB		        set the minimum value of size constraint
// valueUB             set the maximum value of value constraint
// default             sets the default value
// openType            specifies the open Type
// referenceFieldName	the string of the reference field for this type (only if openType used)
// referenceFieldValue	the corresponding value of the reference field for this type (only if openType used)
//
// Other ASN.1 types are not supported; if it encounters them,
// Unmarshal returns a parse error.
func Unmarshal(b []byte, value interface{}) error {
	return UnmarshalWithParams(b, value, "")
}

// UnmarshalWithParams allows field parameters to be specified for the
// top-level element. The form of the params is the same as the field tags.
func UnmarshalWithParams(b []byte, value interface{}, params string) error {
	v := reflect.ValueOf(value).Elem()
	return ParseField(v, b, parseFieldParameters(params))
}
