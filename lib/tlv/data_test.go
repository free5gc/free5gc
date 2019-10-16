//go:binary-only-package

package tlv_test

import "strconv"

type TLVTest struct {
	StructTest        *StructTest         `tlv:"15"`
	BinaryMarshalTest []BinaryMarshalTest `tlv:"65535"`
	SliceTest         []uint16            `tlv:"255"`
}

type NumberTest struct {
	Int8Data   int8   `tlv:"1"`
	Int16Data  int16  `tlv:"2"`
	Int32Data  int32  `tlv:"3"`
	Int64Data  int64  `tlv:"4"`
	UInt8Data  uint8  `tlv:"8"`
	UInt16Data uint16 `tlv:"9"`
	UInt32Data uint32 `tlv:"10"`
	UInt64Data uint64 `tlv:"15"`
}

type StructTest struct {
	Name     []byte `tlv:"20"`
	Sequence uint16 `tlv:"40"`
}

type BinaryMarshalTest struct {
	Value int
}

func (mt *BinaryMarshalTest) MarshalBinary() (data []byte, err error) {}

func (mt *BinaryMarshalTest) UnmarshalBinary(data []byte) (err error) {}
