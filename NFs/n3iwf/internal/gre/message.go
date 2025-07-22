package gre

import (
	"encoding/binary"
	"math"

	"github.com/pkg/errors"
)

// [TS 24.502] 9.3.3 GRE encapsulated user data packet
const (
	GREHeaderFieldLength    = 8
	GREHeaderKeyFieldLength = 4
)

// Ethertypes Specified by the IETF
const (
	IPv4 uint16 = 0x0800
	IPv6 uint16 = 0x86DD
)

type GREPacket struct {
	flags        uint8
	version      uint8
	protocolType uint16
	key          uint32
	payload      []byte
}

func (p *GREPacket) Marshal() []byte {
	packet := make([]byte, GREHeaderFieldLength+len(p.payload))

	packet[0] = p.flags
	packet[1] = p.version
	binary.BigEndian.PutUint16(packet[2:4], p.protocolType)
	binary.BigEndian.PutUint32(packet[4:8], p.key)
	copy(packet[GREHeaderFieldLength:], p.payload)
	return packet
}

func (p *GREPacket) Unmarshal(b []byte) error {
	p.flags = b[0]
	p.version = b[1]

	p.protocolType = binary.BigEndian.Uint16(b[2:4])

	offset := 4

	if p.GetKeyFlag() {
		p.key = binary.BigEndian.Uint32(b[offset : offset+GREHeaderKeyFieldLength])
		offset += GREHeaderKeyFieldLength
	}

	p.payload = append(p.payload, b[offset:]...)
	return nil
}

func (p *GREPacket) SetPayload(payload []byte, protocolType uint16) {
	p.payload = payload
	p.protocolType = protocolType
}

func (p *GREPacket) GetPayload() ([]byte, uint16) {
	return p.payload, p.protocolType
}

func (p *GREPacket) setKeyFlag() {
	p.flags |= 0x20
}

func (p *GREPacket) GetKeyFlag() bool {
	return (p.flags & 0x20) > 0
}

func (p *GREPacket) setQFI(qfi uint8) {
	p.key |= (uint32(qfi) & 0x3F) << 24
}

func (p *GREPacket) setRQI(rqi bool) {
	if rqi {
		p.key |= 0x80
	}
}

func (p *GREPacket) GetQFI() (uint8, error) {
	value := (p.key >> 24) & 0x3F

	if value > math.MaxUint8 {
		return 0, errors.Errorf("GetQFI() value exceeds uint8: %d", value)
	} else {
		return uint8(value), nil
	}
}

func (p *GREPacket) GetRQI() bool {
	return (p.key & 0x80) > 0
}

func (p *GREPacket) GetKeyField() uint32 {
	return p.key
}

func (p *GREPacket) SetQoS(qfi uint8, rqi bool) {
	p.setQFI(qfi)
	p.setRQI(rqi)
	p.setKeyFlag()
}
