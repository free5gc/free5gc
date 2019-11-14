//go:binary-only-package

package pfcp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"free5gc/lib/pfcp/logger"
)

const PfcpVersion uint8 = 1

const (
	SEID_NOT_PRESENT = 0
	SEID_PRESENT     = 1
)

var (
	sequenceCount uint32
)

func init() {}

type Header struct {
	Version         uint8
	MP              uint8
	S               uint8
	MessageType     MessageType
	MessageLength   uint16
	SEID            uint64
	SequenceNumber  uint32
	MessagePriority uint8
}

func (h *Header) MarshalBinary() (data []byte, err error) {}

func (h *Header) UnmarshalBinary(data []byte) error {}

func (h *Header) Len() int {}
