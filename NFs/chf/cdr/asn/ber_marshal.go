package asn

import (
	"fmt"
	"reflect"
)

type encoder interface {
	Len() int
	Encode(dst []byte)
}

type byteEncoder byte

func (b byteEncoder) Len() int {
	return 1
}

func (b byteEncoder) Encode(dst []byte) {
	dst[0] = byte(b)
}

type bytesEncoder []byte

func (b bytesEncoder) Len() int {
	return len(b)
}

func (b bytesEncoder) Encode(dst []byte) {
	copy(dst, b)
}

type stringEncoder string

func (s stringEncoder) Len() int {
	return len(s)
}

func (s stringEncoder) Encode(dst []byte) {
	copy(dst, s)
}

type structEncoder []encoder

func (s structEncoder) Len() int {
	var size int
	for _, e := range s {
		size += e.Len()
	}

	return size
}

func (s structEncoder) Encode(dst []byte) {
	for _, e := range s {
		len := e.Len()
		e.Encode(dst)
		dst = dst[len:]
	}
}

type berTypeEncoder struct {
	tagAndLen encoder
	value     encoder
}

func (b *berTypeEncoder) Len() int {
	return b.tagAndLen.Len() + b.value.Len()
}

func (b *berTypeEncoder) Encode(dst []byte) {
	b.tagAndLen.Encode(dst)
	b.value.Encode(dst[b.tagAndLen.Len():])
}

type int64Encoder int64

func (i int64Encoder) Len() int {
	n := 1

	for i > 127 {
		n++
		i >>= 8
	}

	for i < -128 {
		n++
		i >>= 8
	}

	return n
}

func (i int64Encoder) Encode(dst []byte) {
	n := i.Len()
	// i2 := i
	for j := 0; j < n; j++ {
		dst[n-1-j] = byte(i)
		i >>= 8
	}
	// fmt.Println("marsh: i", i2, "n:", n, "dst", dst)
}

func appendTagAndLen(dst []byte, t tagAndLen) []byte {
	var offset int

	firstByte := byte(t.class) << 6

	if t.constructed {
		firstByte |= 0x20
	}

	if t.tagNumber <= 30 {
		firstByte |= byte(t.tagNumber)
		dst = append(dst, firstByte)
		offset += 1
	} else {
		firstByte |= 31
		dst = append(dst, firstByte)
		offset += 1

		n := 1
		tmp := t.tagNumber
		for tmp > 127 {
			n++
			tmp >>= 7
		}
		dst = append(dst, make([]byte, n)...)
		for i := 0; i < n; i++ {
			dst[n-1-i+offset] = byte(t.tagNumber) | 0x80
			t.tagNumber >>= 7
		}
		dst[n-1+offset] &= 0x7f
		offset += n
	}

	if t.len <= 127 {
		dst = append(dst, byte(t.len))
	} else {
		n := 1
		tmp := t.len
		for tmp > 255 {
			n++
			tmp >>= 8
		}
		dst = append(dst, byte(n)|0x80)
		offset += 1

		dst = append(dst, make([]byte, n)...)
		for i := 0; i < n; i++ {
			dst[n-1-i+offset] = byte(t.len)
			t.len >>= 8
		}
	}

	return dst
}

type bitStringEncoder BitString

func (b bitStringEncoder) Len() int {
	return len(b.Bytes) + 1
}

func (b bitStringEncoder) Encode(dst []byte) {
	// x.690 8.6
	dst[0] = byte(8 - b.BitLength%8)
	copy(dst[1:], b.Bytes)
}

// NOTE: for managementextension field
// type oidEncoder ObjectIdentifier		// Commenting as unused

func makeField(v reflect.Value, params fieldParameters) (encoder, error) {
	if !v.IsValid() {
		return nil, fmt.Errorf("ber: cannot marshal nil value")
	}
	// If the field is an interface{} then recurse into it.
	if v.Kind() == reflect.Interface && v.Type().NumMethod() == 0 {
		return makeField(v.Elem(), params)
	}
	if v.Kind() == reflect.Ptr {
		return makeField(v.Elem(), params)
	}
	fieldType := v.Type()

	var berType berTypeEncoder
	var tag tagAndLen

	// We deal with the structures defined in this package first.
	switch fieldType {
	case BitStringType:
		tag.class = ClassUniversal
		tag.constructed = false
		tag.tagNumber = TagBitString
		berType.value = bitStringEncoder(v.Interface().(BitString))
	case ObjectIdentifierType:
		err := fmt.Errorf("unsupport ObjectIdenfier type")
		return bytesEncoder(nil), err
	case OctetStringType:
		tag.class = ClassUniversal
		tag.constructed = false
		tag.tagNumber = TagOctetString
		berType.value = bytesEncoder(v.Interface().(OctetString))
	case EnumeratedType:
		tag.class = ClassUniversal
		tag.constructed = false
		tag.tagNumber = TagEnumerated
		berType.value = int64Encoder(v.Interface().(Enumerated))
	case NullType:
		tag.class = ClassUniversal
		tag.constructed = false
		tag.tagNumber = TagNull
		berType.value = bytesEncoder(nil)
	default:
		switch val := v; val.Kind() {
		case reflect.Bool:
			tag.class = ClassUniversal
			tag.constructed = false
			tag.tagNumber = TagBoolean
			if v.Bool() {
				berType.value = byteEncoder(0xff)
			} else {
				berType.value = byteEncoder(0)
			}
		case reflect.Int, reflect.Int32, reflect.Int64:
			tag.class = ClassUniversal
			tag.constructed = false
			tag.tagNumber = TagInteger
			berType.value = int64Encoder(v.Int())

		case reflect.Struct:
			structType := fieldType
			if structType.Field(0).Name == "Value" {
				// Non struct type
				// fmt.Println("Non struct type")
				return makeField(val.Field(0), params)
			} else if structType.Field(0).Name == "List" {
				// List Type: SEQUENCE/SET OF
				// fmt.Println("List type")
				return makeField(val.Field(0), params)
			} else if structType.Field(0).Name == "Present" {
				// Open type or CHOICE type
				present := int(v.Field(0).Int())
				tempParams := parseFieldParameters(structType.Field(present).Tag.Get("ber"))
				if present == 0 {
					return nil, fmt.Errorf("CHOICE or OpenType present is 0(present's field number)")
				} else if present >= structType.NumField() {
					return nil, fmt.Errorf("present is bigger than number of struct field")
				} else if params.openType {
					// TODO openType
					return nil, fmt.Errorf("open Type is not implemented")
				} else {
					// Chioce type
					// fmt.Println("Chioce type")
					if params.tagNumber == nil {
						return makeField(val.Field(present), tempParams)
					}
					tag.constructed = true
					var err error
					berType.value, err = makeField(val.Field(present), tempParams)
					if err != nil {
						fmt.Println(err)
					}
				}
			} else {
				// Struct type: SEQUENCE, SET
				// fmt.Println("Struct type")
				tag.class = ClassUniversal
				tag.constructed = true
				if params.set {
					tag.tagNumber = TagSet
				} else {
					tag.tagNumber = TagSequence
				}
				s := make([]encoder, structType.NumField())
				for i := 0; i < structType.NumField(); i++ {
					tempParams := parseFieldParameters(structType.Field(i).Tag.Get("ber"))
					if tempParams.optional {
						if v.Field(i).IsNil() {
							// berTrace(
							// 3, fmt.Sprintf("Field \"%s\" in %s is OPTIONAL and not present", structType.Field(i).Name, structType)
							// )
							s[i] = bytesEncoder(nil)
							continue
						} /*else {
							berTrace(3, fmt.Sprintf("Field \"%s\" in %s is OPTIONAL and present", structType.Field(i).Name, structType))
						}*/
					}

					if tempParams.openType {
						// TODO
						return nil, fmt.Errorf("open Type is not implemented")
					}

					var err error
					s[i], err = makeField(val.Field(i), tempParams)
					if err != nil {
						return nil, fmt.Errorf("iterate subtype error")
					}

					berType.value = structEncoder(s)
				}
			}
		case reflect.Slice:
			tag.class = ClassUniversal
			tag.constructed = true
			if params.set {
				tag.tagNumber = TagSet
			} else {
				tag.tagNumber = TagSequence
			}
			s := make([]encoder, v.Len())
			var err error
			tempParams := params
			tempParams.tagNumber = nil
			for i := 0; i < v.Len(); i++ {
				s[i], err = makeField(val.Index(i), tempParams)
				if err != nil {
					return nil, fmt.Errorf("iterate subtype error")
				}
			}

			berType.value = structEncoder(s)
		case reflect.String:
			tag.class = ClassUniversal
			tag.constructed = false
			tag.tagNumber = uint64(params.stringType)

			berType.value = stringEncoder(v.String())
		}
	}
	tag.len = int64(berType.value.Len())

	if params.tagNumber != nil {
		if params.explicitTag {
			t := berType
			t.tagAndLen = bytesEncoder(appendTagAndLen(make([]byte, 8)[:0], tag))
			berType.value = &t
			tag.constructed = true
			tag.len = int64(berType.value.Len())
		}
		tag.class = ClassContextSpecific
		tag.tagNumber = *params.tagNumber
	}

	berType.tagAndLen = bytesEncoder(appendTagAndLen(make([]byte, 8)[:0], tag))

	return &berType, nil
}

// Marshal returns the ASN.1 encoding of val.
func BerMarshal(val interface{}) ([]byte, error) {
	return BerMarshalWithParams(val, "")
}

// MarshalWithParams allows field parameters to be specified for the
// top-level element. The form of the params is the same as the field tags.
func BerMarshalWithParams(val interface{}, params string) ([]byte, error) {
	e, err := makeField(reflect.ValueOf(val), parseFieldParameters(params))
	if err != nil {
		return nil, err
	}
	b := make([]byte, e.Len())
	e.Encode(b)
	return b, nil
}
