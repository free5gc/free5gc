//go:binary-only-package

package aper

import (
	"fmt"
	"free5gc/lib/aper/logger"
	"path"
	"reflect"
	"runtime"
)

type perBitData struct {
	bytes      []byte
	byteOffset uint64
	bitsOffset uint
}

func perTrace(level int, s string) {}

func perBitLog(numBits uint64, byteOffset uint64, bitsOffset uint, value interface{}) string {}

// GetBitString is to get BitString with desire size from source byte array with bit offset
func GetBitString(srcBytes []byte, bitsOffset uint, numBits uint) (dstBytes []byte, err error) {}

// GetFewBits is to get Value with desire few bits from source byte with bit offset
// func GetFewBits(srcByte byte, bitsOffset uint, numBits uint) (value uint64, err error) {}

func (pd *perBitData) bitCarry() {}

func (pd *perBitData) getBitString(numBits uint) (dstBytes []byte, err error) {}

func (pd *perBitData) getBitsValue(numBits uint) (value uint64, err error) {}

func (pd *perBitData) parseAlignBits() error {}

func (pd *perBitData) parseConstraintValue(valueRange int64) (value uint64, err error) {}

func (pd *perBitData) parseLength(sizeRange int64, repeat *bool) (value uint64, err error) {}

func (pd *perBitData) parseBitString(extensed bool, lowerBoundPtr *int64, upperBoundPtr *int64) (bitString BitString, err error) {}
func (pd *perBitData) parseOctetString(extensed bool, lowerBoundPtr *int64, upperBoundPtr *int64) (octetString OctetString, err error) {}

func (pd *perBitData) parseBool() (value bool, err error) {}

func (pd *perBitData) parseInteger(extensed bool, lowerBoundPtr *int64, upperBoundPtr *int64) (value int64, err error) {}

// parse ENUMERATED type but do not implement extensive value and different value with index
func (pd *perBitData) parseEnumerated(extensed bool, lowerBoundPtr *int64, upperBoundPtr *int64) (value uint64, err error) {}
func (pd *perBitData) parseSequenceOf(sizeExtensed bool, params fieldParameters, sliceType reflect.Type) (sliceContent reflect.Value, err error) {}

func (pd *perBitData) getChoiceIndex(extensed bool, upperBoundPtr *int64) (present int, err error) {}
func getReferenceFieldValue(v reflect.Value) (value int64, err error) {}

func (pd *perBitData) parseOpenType(v reflect.Value, params fieldParameters) (err error) {}

// parseField is the main parsing function. Given a byte slice and an offset
// into the array, it will try to parse a suitable ASN.1 value out and store it
// in the given Value. TODO : ObjectIdenfier, handle extension Field
func parseField(v reflect.Value, pd *perBitData, params fieldParameters) (err error) {}

// Unmarshal parses the APER-encoded ASN.1 data structure b
// and uses the reflect //go:binary-only-package

package to fill in an arbitrary value pointed at by value.
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
//
//	optional        	OPTIONAL tag in SEQUENCE
//	sizeExt             specifies that size  is extensible
//	valueExt            specifies that value is extensible
//	sizeLB		        set the minimum value of size constraint
//	sizeUB              set the maximum value of value constraint
//	valueLB		        set the minimum value of size constraint
//	valueUB             set the maximum value of value constraint
//	default             sets the default value
//	openType            specifies the open Type
//  referenceFieldName	the string of the reference field for this type (only if openType used)
//  referenceFieldValue	the corresponding value of the reference field for this type (only if openType used)
//
// Other ASN.1 types are not supported; if it encounters them,
// Unmarshal returns a parse error.
func Unmarshal(b []byte, value interface{}) error {}

// UnmarshalWithParams allows field parameters to be specified for the
// top-level element. The form of the params is the same as the field tags.
func UnmarshalWithParams(b []byte, value interface{}, params string) error {}
