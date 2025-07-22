package processor_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/smf/internal/context"
	"github.com/free5gc/util/idgenerator"
)

func TestBuildNASPacketFilterFromPacketFilterInfo(t *testing.T) {
	testCases := []struct {
		name         string
		packetFilter []nasType.PacketFilter
		flowInfo     models.FlowInformation
	}{
		{
			name: "MatchAll",
			packetFilter: []nasType.PacketFilter{
				{
					Direction: nasType.PacketFilterDirectionBidirectional,
					Components: nasType.PacketFilterComponentList{
						&nasType.PacketFilterMatchAll{},
					},
				},
			},
			flowInfo: models.FlowInformation{
				FlowDirection:   models.FlowDirection_BIDIRECTIONAL,
				FlowDescription: "permit out ip from any to assigned",
			},
		},
		{
			name: "MatchIPNet1",
			packetFilter: []nasType.PacketFilter{
				{
					Direction: nasType.PacketFilterDirectionUplink,
					Components: nasType.PacketFilterComponentList{
						&nasType.PacketFilterIPv4LocalAddress{
							Address: net.ParseIP("192.168.0.0").To4(),
							Mask:    net.IPv4Mask(255, 255, 0, 0),
						},
					},
				},
			},
			flowInfo: models.FlowInformation{
				FlowDirection:   models.FlowDirection_UPLINK,
				FlowDescription: "permit out ip from any to 192.168.0.0/16",
			},
		},
		{
			name: "MatchIPNet2",
			packetFilter: []nasType.PacketFilter{
				{
					Direction: nasType.PacketFilterDirectionBidirectional,
					Components: nasType.PacketFilterComponentList{
						&nasType.PacketFilterIPv4LocalAddress{
							Address: net.ParseIP("192.168.0.0").To4(),
							Mask:    net.IPv4Mask(255, 255, 0, 0),
						},
						&nasType.PacketFilterIPv4RemoteAddress{
							Address: net.ParseIP("10.160.20.0").To4(),
							Mask:    net.IPv4Mask(255, 255, 255, 0),
						},
					},
				},
			},
			flowInfo: models.FlowInformation{
				FlowDirection:   models.FlowDirection_BIDIRECTIONAL,
				FlowDescription: "permit out ip from 10.160.20.0/24 to 192.168.0.0/16",
			},
		},
		{
			name: "MatchIPNetPort",
			packetFilter: []nasType.PacketFilter{
				{
					Direction: nasType.PacketFilterDirectionBidirectional,
					Components: nasType.PacketFilterComponentList{
						&nasType.PacketFilterIPv4LocalAddress{
							Address: net.ParseIP("192.168.0.0").To4(),
							Mask:    net.IPv4Mask(255, 255, 0, 0),
						},
						&nasType.PacketFilterSingleLocalPort{
							Value: 8000,
						},
						&nasType.PacketFilterIPv4RemoteAddress{
							Address: net.ParseIP("10.160.20.0").To4(),
							Mask:    net.IPv4Mask(255, 255, 255, 0),
						},
					},
				},
			},
			flowInfo: models.FlowInformation{
				FlowDirection:   models.FlowDirection_BIDIRECTIONAL,
				FlowDescription: "permit out ip from 10.160.20.0/24 to 192.168.0.0/16 8000",
			},
		},
		{
			name: "MatchIPNetPortRanges",
			packetFilter: []nasType.PacketFilter{
				{
					Direction: nasType.PacketFilterDirectionDownlink,
					Components: nasType.PacketFilterComponentList{
						&nasType.PacketFilterIPv4LocalAddress{
							Address: net.ParseIP("192.168.0.0").To4(),
							Mask:    net.IPv4Mask(255, 255, 0, 0),
						},
						&nasType.PacketFilterLocalPortRange{
							LowLimit:  3000,
							HighLimit: 8000,
						},
						&nasType.PacketFilterIPv4RemoteAddress{
							Address: net.ParseIP("10.160.20.0").To4(),
							Mask:    net.IPv4Mask(255, 255, 255, 0),
						},
					},
				},
			},
			flowInfo: models.FlowInformation{
				FlowDirection:   models.FlowDirection_DOWNLINK,
				FlowDescription: "permit out ip from 10.160.20.0/24 to 192.168.0.0/16 3000-8000",
			},
		},
		{
			name: "MatchIPNetPortRanges2",
			packetFilter: []nasType.PacketFilter{
				{
					Direction: nasType.PacketFilterDirectionDownlink,
					Components: nasType.PacketFilterComponentList{
						&nasType.PacketFilterIPv4LocalAddress{
							Address: net.ParseIP("192.168.0.0").To4(),
							Mask:    net.IPv4Mask(255, 255, 0, 0),
						},
						&nasType.PacketFilterLocalPortRange{
							LowLimit:  6000,
							HighLimit: 8000,
						},
						&nasType.PacketFilterIPv4RemoteAddress{
							Address: net.ParseIP("10.160.20.0").To4(),
							Mask:    net.IPv4Mask(255, 255, 255, 0),
						},
						&nasType.PacketFilterRemotePortRange{
							LowLimit:  3000,
							HighLimit: 4000,
						},
					},
				},
			},
			flowInfo: models.FlowInformation{
				FlowDirection:   models.FlowDirection_DOWNLINK,
				FlowDescription: "permit out ip from 10.160.20.0/24 3000-4000 to 192.168.0.0/16 6000-8000",
			},
		},
		{
			name: "MatchIPNetPortRanges3",
			packetFilter: []nasType.PacketFilter{
				{
					Direction: nasType.PacketFilterDirectionDownlink,
					Components: nasType.PacketFilterComponentList{
						&nasType.PacketFilterIPv4LocalAddress{
							Address: net.ParseIP("192.168.0.0").To4(),
							Mask:    net.IPv4Mask(255, 255, 0, 0),
						},
						&nasType.PacketFilterLocalPortRange{
							LowLimit:  6000,
							HighLimit: 7000,
						},
						&nasType.PacketFilterIPv4RemoteAddress{
							Address: net.ParseIP("10.160.20.0").To4(),
							Mask:    net.IPv4Mask(255, 255, 255, 0),
						},
						&nasType.PacketFilterRemotePortRange{
							LowLimit:  3000,
							HighLimit: 4000,
						},
					},
				},
				{
					Direction: nasType.PacketFilterDirectionDownlink,
					Components: nasType.PacketFilterComponentList{
						&nasType.PacketFilterIPv4LocalAddress{
							Address: net.ParseIP("192.168.0.0").To4(),
							Mask:    net.IPv4Mask(255, 255, 0, 0),
						},
						&nasType.PacketFilterSingleLocalPort{
							Value: 8000,
						},
						&nasType.PacketFilterIPv4RemoteAddress{
							Address: net.ParseIP("10.160.20.0").To4(),
							Mask:    net.IPv4Mask(255, 255, 255, 0),
						},
						&nasType.PacketFilterRemotePortRange{
							LowLimit:  3000,
							HighLimit: 4000,
						},
					},
				},
			},
			flowInfo: models.FlowInformation{
				FlowDirection:   models.FlowDirection_DOWNLINK,
				FlowDescription: "permit out ip from 10.160.20.0/24 3000-4000 to 192.168.0.0/16 6000-7000,8000",
			},
		},
		{
			name: "MatchIPNetPortRanges4",
			packetFilter: []nasType.PacketFilter{
				{
					Direction: nasType.PacketFilterDirectionDownlink,
					Components: nasType.PacketFilterComponentList{
						&nasType.PacketFilterIPv4LocalAddress{
							Address: net.ParseIP("192.168.0.0").To4(),
							Mask:    net.IPv4Mask(255, 255, 0, 0),
						},
						&nasType.PacketFilterLocalPortRange{
							LowLimit:  6000,
							HighLimit: 7000,
						},
						&nasType.PacketFilterIPv4RemoteAddress{
							Address: net.ParseIP("10.160.20.0").To4(),
							Mask:    net.IPv4Mask(255, 255, 255, 0),
						},
						&nasType.PacketFilterRemotePortRange{
							LowLimit:  3000,
							HighLimit: 4000,
						},
					},
				},
				{
					Direction: nasType.PacketFilterDirectionDownlink,
					Components: nasType.PacketFilterComponentList{
						&nasType.PacketFilterIPv4LocalAddress{
							Address: net.ParseIP("192.168.0.0").To4(),
							Mask:    net.IPv4Mask(255, 255, 0, 0),
						},
						&nasType.PacketFilterLocalPortRange{
							LowLimit:  6000,
							HighLimit: 7000,
						},
						&nasType.PacketFilterIPv4RemoteAddress{
							Address: net.ParseIP("10.160.20.0").To4(),
							Mask:    net.IPv4Mask(255, 255, 255, 0),
						},
						&nasType.PacketFilterSingleRemotePort{
							Value: 5000,
						},
					},
				},
			},
			flowInfo: models.FlowInformation{
				FlowDirection:   models.FlowDirection_DOWNLINK,
				FlowDescription: "permit out ip from 10.160.20.0/24 3000-4000,5000 to 192.168.0.0/16 6000-7000",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			smCtx := &context.SMContext{
				PacketFilterIDGenerator: idgenerator.NewGenerator(1, 255),
				PacketFilterIDToNASPFID: make(map[string]uint8),
			}
			packetFilters, err := context.BuildNASPacketFiltersFromFlowInformation(&tc.flowInfo, smCtx)
			require.NoError(t, err)

			for i, pf := range packetFilters {
				require.Equal(t, tc.packetFilter[i].Direction, pf.Direction)
				require.Equal(t, tc.packetFilter[i].Components, pf.Components)
			}
		})
	}
}
