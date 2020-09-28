package test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"runtime"

	"github.com/calee0219/fatal"
	"golang.org/x/sys/unix"
)

type iphdr struct {
	vhl   uint8
	tos   uint8
	iplen uint16
	// id    uint16
	off   uint16
	ttl   uint8
	proto uint8
	csum  uint16
	src   [4]byte
	dst   [4]byte
}

type udphdr struct {
	src  uint16
	dst  uint16
	ulen uint16
	csum uint16
}

// pseudo header used for checksum calculation
type pseudohdr struct {
	ipsrc   [4]byte
	ipdst   [4]byte
	zero    uint8
	ipproto uint8
	plen    uint16
}

func checksum(buf []byte) uint16 {
	sum := uint32(0)

	for ; len(buf) >= 2; buf = buf[2:] {
		sum += uint32(buf[0])<<8 | uint32(buf[1])
	}
	if len(buf) > 0 {
		sum += uint32(buf[0]) << 8
	}
	for sum > 0xffff {
		sum = (sum >> 16) + (sum & 0xffff)
	}
	csum := ^uint16(sum)
	/*
	 * From RFC 768:
	 * If the computed checksum is zero, it is transmitted as all ones (the
	 * equivalent in one's complement arithmetic). An all zero transmitted
	 * checksum value means that the transmitter generated no checksum (for
	 * debugging or for higher level protocols that don't care).
	 */
	if csum == 0 {
		csum = 0xffff
	}
	return csum
}

func (h *iphdr) checksum() {
	h.csum = 0
	var b bytes.Buffer
	err := binary.Write(&b, binary.BigEndian, h)
	if err != nil {
		fatal.Fatalf("binary Write error in checksum: %+v", err)
	}
	h.csum = checksum(b.Bytes())
}

func (u *udphdr) checksum(ip *iphdr, payload []byte) {
	u.csum = 0
	phdr := pseudohdr{
		ipsrc:   ip.src,
		ipdst:   ip.dst,
		zero:    0,
		ipproto: ip.proto,
		plen:    u.ulen,
	}
	var b bytes.Buffer
	err := binary.Write(&b, binary.BigEndian, &phdr)
	if err != nil {
		fatal.Fatalf("binary Write error in checksum: %+v", err)
	}
	err = binary.Write(&b, binary.BigEndian, u)
	if err != nil {
		fatal.Fatalf("binary Write error in checksum: %+v", err)
	}
	err = binary.Write(&b, binary.BigEndian, &payload)
	if err != nil {
		fatal.Fatalf("binary Write error in checksum: %+v", err)
	}
	u.csum = checksum(b.Bytes())
}

func BuildRawUdpIp(srcIP, dstIP string, srcPort, dstPort uint16, payload []byte) ([]byte, error) {
	var err error

	ip := iphdr{
		vhl:   0x45,
		tos:   0,
		off:   0,
		ttl:   64,
		proto: unix.IPPROTO_UDP,
	}
	copy(ip.src[:], net.ParseIP(srcIP).To4())
	copy(ip.dst[:], net.ParseIP(dstIP).To4())
	// iplen and csum set later

	udp := udphdr{
		src: srcPort,
		dst: dstPort,
	}
	// ulen and csum set later

	udplen := 8 + len(payload)
	totalLen := 20 + udplen
	if totalLen > 0xffff {
		err = fmt.Errorf("message is too large to fit into a packet: %v > %v\n", totalLen, 0xffff)
		return nil, err
	}

	// the kernel will overwrite the IP checksum, so this is included just for
	// completeness
	ip.iplen = uint16(totalLen)
	ip.checksum()

	// the kernel doesn't touch the UDP checksum, so we can either set it
	// correctly or leave it zero to indicate that we didn't use a checksum
	udp.ulen = uint16(udplen)
	udp.checksum(&ip, payload)

	var b bytes.Buffer
	err = binary.Write(&b, binary.BigEndian, &ip)
	if err != nil {
		return nil, err
	}
	err = binary.Write(&b, binary.BigEndian, &udp)
	if err != nil {
		return nil, err
	}
	err = binary.Write(&b, binary.BigEndian, &payload)
	if err != nil {
		return nil, err
	}
	bb := b.Bytes()

	/*
	 * For some reason, the IP header's length field needs to be in host byte order
	 * in OS X.
	 */
	if runtime.GOOS == "darwin" {
		bb[2], bb[3] = bb[3], bb[2]
	}
	return bb, nil
}
