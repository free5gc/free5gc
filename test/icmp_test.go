package test_test

import (
	"encoding/hex"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func TestICMP(t *testing.T) {

	upfAddr, err := net.ResolveUDPAddr("udp", "172.0.0.1:2152")
	assert.Nil(t, err)
	gnbAddr, err := net.ResolveUDPAddr("udp", "172.0.0.2:2152")
	assert.Nil(t, err)
	conn, err := net.DialUDP("udp", gnbAddr, upfAddr)
	assert.Nil(t, err)
	gtpHdr, err := hex.DecodeString("32ff00340000000100000000")
	assert.Nil(t, err)
	icmpData, err := hex.DecodeString("8c870d0000000000101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f3031323334353637")
	assert.Nil(t, err)

	ipv4hdr := ipv4.Header{
		Version:  4,
		Len:      20,
		Protocol: 1,
		Flags:    0,
		TotalLen: 48,
		TTL:      64,
		Src:      net.ParseIP("10.60.0.1").To4(),
		Dst:      net.ParseIP("10.60.0.20").To4(),
		ID:       1,
		Checksum: 0x0240,
	}

	v4HdrBuf, err := ipv4hdr.Marshal()
	assert.Nil(t, err)
	tt := append(gtpHdr, v4HdrBuf...)
	assert.Nil(t, err)

	m := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: 12394, Seq: 1, //<< uint(seq), // TODO
			Data: icmpData,
		},
	}
	b, err := m.Marshal(nil)
	assert.Nil(t, err)
	b[2] = 0xaf
	b[3] = 0x88
	_, err = conn.Write(append(tt, b...))
	assert.Nil(t, err)
}
