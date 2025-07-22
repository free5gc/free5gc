package forwarder

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// IPFilterRule <- action s dir s proto s 'from' s src 'to' s dst
// action <- 'permit' / 'deny'
// dir <- 'in' / 'out'
// proto <- 'ip' / digit
// src <- addr s ports?
// dst <- addr s ports?
// addr <- 'any' / 'assigned' / cidr
// cidr <- ipv4addr ('/' digit)?
// ipv4addr <- digit ('.' digit){3}
// ports <- port (',' port)*
// port <- (digit '-' digit) / digit
// digit <- [1-9][0-9]+
// s <- ' '+

type FlowDesc struct {
	Action   string
	Dir      string
	Proto    uint8
	Src      *net.IPNet
	Dst      *net.IPNet
	SrcPorts [][]uint16
	DstPorts [][]uint16
}

func ParseFlowDesc(s string) (*FlowDesc, error) {
	fd := new(FlowDesc)
	token := strings.Fields(s)
	pos := 0

	if pos >= len(token) {
		return nil, fmt.Errorf("too few fields %v", len(token))
	}
	switch token[pos] {
	case "permit":
		fd.Action = token[pos]
	default:
		return nil, fmt.Errorf("unknown action %v", token[pos])
	}
	pos++

	if pos >= len(token) {
		return nil, fmt.Errorf("too few fields %v", len(token))
	}
	switch token[pos] {
	case "in", "out":
		fd.Dir = token[pos]
	default:
		return nil, fmt.Errorf("unknown direction %v", token[pos])
	}
	pos++

	if pos >= len(token) {
		return nil, fmt.Errorf("too few fields %v", len(token))
	}
	switch token[pos] {
	case "ip":
		fd.Proto = 0xff
	default:
		v, err := strconv.ParseUint(token[pos], 10, 8)
		if err != nil {
			return nil, err
		}
		fd.Proto = uint8(v)
	}
	pos++

	if pos >= len(token) {
		return nil, fmt.Errorf("too few fields %v", len(token))
	}
	if token[pos] != "from" {
		return nil, fmt.Errorf("not match 'from'")
	}
	pos++

	if pos >= len(token) {
		return nil, fmt.Errorf("too few fields %v", len(token))
	}
	src, err := ParseFlowDescIPNet(token[pos])
	if err != nil {
		return nil, err
	}
	pos++
	fd.Src = src

	if pos >= len(token) {
		return nil, fmt.Errorf("too few fields %v", len(token))
	}
	sports, err := ParseFlowDescPorts(token[pos])
	if err == nil {
		fd.SrcPorts = sports
		pos++
	}

	if pos >= len(token) {
		return nil, fmt.Errorf("too few fields %v", len(token))
	}
	if token[pos] != "to" {
		return nil, fmt.Errorf("not match 'to'")
	}
	pos++

	if pos >= len(token) {
		return nil, fmt.Errorf("too few fields %v", len(token))
	}
	dst, err := ParseFlowDescIPNet(token[pos])
	if err != nil {
		return nil, err
	}
	pos++
	fd.Dst = dst

	if pos < len(token) {
		dports, err := ParseFlowDescPorts(token[pos])
		if err == nil {
			fd.DstPorts = dports
		}
	}

	return fd, nil
}

func ParseFlowDescIPNet(s string) (*net.IPNet, error) {
	if s == "any" || s == "assigned" {
		return &net.IPNet{
			IP:   net.IPv6zero,
			Mask: net.CIDRMask(0, 128),
		}, nil
	}
	_, ipnet, err := net.ParseCIDR(s)
	if err == nil {
		return ipnet, nil
	}
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, fmt.Errorf("invalid address %v", s)
	}
	v4 := ip.To4()
	if v4 != nil {
		ip = v4
	}
	n := len(ip) * 8
	return &net.IPNet{
		IP:   ip,
		Mask: net.CIDRMask(n, n),
	}, nil
}

func ParseFlowDescPorts(s string) ([][]uint16, error) {
	var vals [][]uint16
	for _, port := range strings.Split(s, ",") {
		digit := strings.SplitN(port, "-", 2)
		switch len(digit) {
		case 1:
			v, err := strconv.ParseUint(digit[0], 10, 16)
			if err != nil {
				return nil, err
			}
			vals = append(vals, []uint16{uint16(v)})
		case 2:
			start, err := strconv.ParseUint(digit[0], 10, 16)
			if err != nil {
				return nil, err
			}
			end, err := strconv.ParseUint(digit[1], 10, 16)
			if err != nil {
				return nil, err
			}
			vals = append(vals, []uint16{uint16(start), uint16(end)})
		default:
			return nil, fmt.Errorf("invalid port: %q", port)
		}
	}
	return vals, nil
}
