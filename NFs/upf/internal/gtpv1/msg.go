package gtpv1

import (
	"encoding/binary"
)

type Encoder interface {
	Len() int
	Encode([]byte) (int, error)
}

// Message Type definitions.
const (
	MsgTypeTPDU uint8 = 255
)

type Message struct {
	Flags          uint8
	Type           uint8
	TEID           uint32
	SequenceNumber uint16
	NPDUNumber     uint8
	Exts           []Encoder
	Payload        []byte
}

func (m Message) HasSequence() bool {
	return m.Flags&0x2 != 0
}

func (m Message) HasNPDUNumber() bool {
	return m.Flags&0x1 != 0
}

func (m Message) Len() int {
	l := 8
	if m.HasSequence() {
		l += 2
	}
	if m.HasNPDUNumber() {
		l++
	}
	l = ((l + 4) &^ 0x3) - 1
	for _, e := range m.Exts {
		l += e.Len()
	}
	l++
	l += len(m.Payload)
	return l
}

func (m Message) Encode(b []byte) (int, error) {
	b[0] = m.Flags
	b[1] = m.Type
	l := m.Len() - 8
	binary.BigEndian.PutUint16(b[2:4], uint16(l))
	binary.BigEndian.PutUint32(b[4:8], m.TEID)
	pos := 8
	if m.HasSequence() {
		binary.BigEndian.PutUint16(b[pos:pos+2], m.SequenceNumber)
		pos += 2
	}
	if m.HasNPDUNumber() {
		b[pos] = m.NPDUNumber
		pos++
	}
	// alignment
	pos = ((pos + 4) &^ 0x3) - 1
	for _, e := range m.Exts {
		n, err := e.Encode(b[pos:])
		if err != nil {
			return n, err
		}
		pos += n
	}
	// No more extension headers
	b[pos] = 0
	pos++
	copy(b[pos:], m.Payload)
	return m.Len(), nil
}

type PDUSessionContainer struct {
	PDUType   uint8
	QoSFlowID uint8
}

func (e PDUSessionContainer) Len() int {
	return 4
}

func (e PDUSessionContainer) Encode(b []byte) (int, error) {
	b[0] = 0x85
	b[1] = 1
	b[2] = e.PDUType << 4
	b[3] = e.QoSFlowID & 0xf
	return e.Len(), nil
}
