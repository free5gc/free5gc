package forwarder

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlowDesc(t *testing.T) {
	cases := []struct {
		name string
		s    string
		fd   FlowDesc
		err  error
	}{
		{
			name: "host addr",
			s:    "permit out ip from 10.20.30.40 to 50.60.70.80",
			fd: FlowDesc{
				Action: "permit",
				Dir:    "out",
				Proto:  0xff,
				Src: &net.IPNet{
					IP:   net.IPv4(10, 20, 30, 40).To4(),
					Mask: net.CIDRMask(32, 32),
				},
				Dst: &net.IPNet{
					IP:   net.IPv4(50, 60, 70, 80).To4(),
					Mask: net.CIDRMask(32, 32),
				},
			},
		},
		{
			name: "proto",
			s:    "permit out 210 from 10.20.30.40 to 50.60.70.80",
			fd: FlowDesc{
				Action: "permit",
				Dir:    "out",
				Proto:  210,
				Src: &net.IPNet{
					IP:   net.IPv4(10, 20, 30, 40).To4(),
					Mask: net.CIDRMask(32, 32),
				},
				Dst: &net.IPNet{
					IP:   net.IPv4(50, 60, 70, 80).To4(),
					Mask: net.CIDRMask(32, 32),
				},
			},
		},
		{
			name: "network addr",
			s:    "permit out ip from 10.20.30.40/24 to 50.60.70.80/16",
			fd: FlowDesc{
				Action: "permit",
				Dir:    "out",
				Proto:  0xff,
				Src: &net.IPNet{
					IP:   net.IPv4(10, 20, 30, 0).To4(),
					Mask: net.CIDRMask(24, 32),
				},
				Dst: &net.IPNet{
					IP:   net.IPv4(50, 60, 0, 0).To4(),
					Mask: net.CIDRMask(16, 32),
				},
			},
		},
		{
			name: "source port",
			s:    "permit out ip from 10.20.30.0/24 345,789-792,1023-1026 to 50.60.0.0/16",
			fd: FlowDesc{
				Action: "permit",
				Dir:    "out",
				Proto:  0xff,
				Src: &net.IPNet{
					IP:   net.IPv4(10, 20, 30, 0).To4(),
					Mask: net.CIDRMask(24, 32),
				},
				Dst: &net.IPNet{
					IP:   net.IPv4(50, 60, 0, 0).To4(),
					Mask: net.CIDRMask(16, 32),
				},
				SrcPorts: [][]uint16{
					{
						345,
					},
					{
						789,
						792,
					},
					{
						1023,
						1026,
					},
				},
			},
		},
		{
			name: "dst port",
			s:    "permit out ip from 10.20.30.0/24 to 50.60.0.0/16 345,789-792,1023-1026",
			fd: FlowDesc{
				Action: "permit",
				Dir:    "out",
				Proto:  0xff,
				Src: &net.IPNet{
					IP:   net.IPv4(10, 20, 30, 0).To4(),
					Mask: net.CIDRMask(24, 32),
				},
				Dst: &net.IPNet{
					IP:   net.IPv4(50, 60, 0, 0).To4(),
					Mask: net.CIDRMask(16, 32),
				},
				DstPorts: [][]uint16{
					{
						345,
					},
					{
						789,
						792,
					},
					{
						1023,
						1026,
					},
				},
			},
		},
		{
			name: "any to assign",
			s:    "permit out ip from any to assigned",
			fd: FlowDesc{
				Action: "permit",
				Dir:    "out",
				Proto:  0xff,
				Src: &net.IPNet{
					IP:   net.IPv6zero,
					Mask: net.CIDRMask(0, 128),
				},
				Dst: &net.IPNet{
					IP:   net.IPv6zero,
					Mask: net.CIDRMask(0, 128),
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			fd, err := ParseFlowDesc(tt.s)
			if tt.err == nil {
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, &tt.fd, fd)
			} else if err != tt.err {
				t.Errorf("wantErr %v; but got %v", tt.err, err)
			}
		})
	}
}
