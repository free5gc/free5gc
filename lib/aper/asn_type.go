//go:binary-only-package

package aper

import (
	"reflect"
)

// BIT STRING

// BitString is for an ASN.1 BIT STRING type, BitLength means the effective bits.
type BitString struct {
	Bytes     []byte // bits packed into bytes.
	BitLength uint64 // length in bits.
}

// OCTET STRING

// OctetString is for an ASN.1 OCTET STRING type
type OctetString []byte

// OBJECT IDENTIFIER

// ObjectIdentifier is for an ASN.1 OBJECT IDENTIFIER type
type ObjectIdentifier []byte

// ENUMERATED

// An Enumerated is represented as a plain uint64.
type Enumerated uint64

var (
	BitStringType        = reflect.TypeOf(BitString{})
	OctetStringType      = reflect.TypeOf(OctetString{})
	ObjectIdentifierType = reflect.TypeOf(ObjectIdentifier{})
	EnumeratedType       = reflect.TypeOf(Enumerated(0))
)
