package test

import (
	"bytes"
	"encoding/binary"
)

type gtphdr struct {
	flags uint8 //Version(3-bits), Protocal Type(1-bit), Extension Header flag(1-bit),
	// Sequence Number flag(1-bit), N-PDU number flag(1-bit)
	mtype uint8  //Message Type
	len   uint16 //Total Length
	teid  uint32 //Tunnel Endpoint Identifier
}

type gtphdropt struct {
	sn          uint16 //Sequence Number
	npdunum     uint8  //N-PDU Number
	nexthdrtype uint8  //Next Extenstion Header Type
}

func BuildGTPv1Header(snFlag bool, sn uint16, nPduFlag bool, nPduNum uint8,
	extHdrFlag bool, nExtHdrType uint8, payloadLen uint16, teID uint32) ([]byte, error) {
	var flags uint8 = (0x3 << 4) //Version=1, Protocol Type=GTP
	if extHdrFlag {
		flags |= 0x1 << 2
	}
	if snFlag {
		flags |= 0x1 << 1
	}
	if nPduFlag {
		flags |= 0x1
	}
	ghdr := gtphdr{
		flags: flags,
		mtype: 0xFF, //G-PDU
		len:   payloadLen,
		teid:  teID,
	}
	var b bytes.Buffer
	err := binary.Write(&b, binary.BigEndian, &ghdr)
	if err != nil {
		return nil, err
	}
	if snFlag || nPduFlag || extHdrFlag {
		ghdropt := gtphdropt{
			sn:          sn,
			npdunum:     nPduNum,
			nexthdrtype: nExtHdrType,
		}
		err = binary.Write(&b, binary.BigEndian, &ghdropt)
		if err != nil {
			return nil, err
		}
	}
	bb := b.Bytes()
	return bb, nil
}
