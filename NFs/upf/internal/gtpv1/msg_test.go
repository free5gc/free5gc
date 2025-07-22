package gtpv1

import (
	"bytes"
	"testing"
)

func TestMessage(t *testing.T) {
	pkt := []byte{
		0x34, 0xff, 0x00, 0x0c, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x85, 0x01, 0x10, 0x09, 0x00,
		0xde, 0xad, 0xbe, 0xef,
	}
	msg := Message{
		Flags: 0x34,
		Type:  MsgTypeTPDU,
		TEID:  1,
		Exts: []Encoder{
			PDUSessionContainer{
				PDUType:   1,
				QoSFlowID: 9,
			},
		},
		Payload: []byte{0xde, 0xad, 0xbe, 0xef},
	}
	l := msg.Len()
	b := make([]byte, l)
	n, err := msg.Encode(b)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(pkt) {
		t.Errorf("want %v; but got %v\n", len(pkt), n)
	}
	if !bytes.Equal(b, pkt) {
		t.Errorf("want %x; but got %x\n", pkt, b)
	}
}
