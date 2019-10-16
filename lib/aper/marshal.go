//go:binary-only-package

package aper

import (
	"fmt"
	"log"
	"reflect"
)

type perRawBitData struct {
	bytes      []byte
	bitsOffset uint
}

func perRawBitLog(numBits uint64, byteLen int, bitsOffset uint, value interface{}) string {}

func (pd *perRawBitData) bitCarry() {}
func (pd *perRawBitData) appendAlignBits() {}

func (pd *perRawBitData) putBitString(bytes []byte, numBits uint) (err error) {}

func (pd *perRawBitData) putBitsValue(value uint64, numBits uint) (err error) {}

func (pd *perRawBitData) appendConstraintValue(valueRange int64, value uint64) (err error) {}

func (pd *perRawBitData) appendLength(sizeRange int64, value uint64) (err error) {}

func (pd *perRawBitData) appendBitString(bytes []byte, bitsLength uint64, extensive bool, lowerBoundPtr *int64, upperBoundPtr *int64) (err error) {}

func (pd *perRawBitData) appendOctetString(bytes []byte, extensive bool, lowerBoundPtr *int64, upperBoundPtr *int64) (err error) {}

func (pd *perRawBitData) appendBool(value bool) (err error) {}

func (pd *perRawBitData) appendInteger(value int64, extensive bool, lowerBoundPtr *int64, upperBoundPtr *int64) (err error) {}

// append ENUMERATED type but do not implement extensive value and different value with index
func (pd *perRawBitData) appendEnumerated(value uint64, extensive bool, lowerBoundPtr *int64, upperBoundPtr *int64) (err error) {}

func (pd *perRawBitData) parseSequenceOf(v reflect.Value, params fieldParameters) (err error) {}

func (pd *perRawBitData) appendChoiceIndex(present int, extensive bool, upperBoundPtr *int64) (err error) {}

func (pd *perRawBitData) appendOpenType(v reflect.Value, params fieldParameters) (err error) {}
func (pd *perRawBitData) makeField(v reflect.Value, params fieldParameters) (err error) {}

// Marshal returns the ASN.1 encoding of val.
func Marshal(val interface{}) ([]byte, error) {}

// MarshalWithParams allows field parameters to be specified for the
// top-level element. The form of the params is the same as the field tags.
func MarshalWithParams(val interface{}, params string) ([]byte, error) {}
