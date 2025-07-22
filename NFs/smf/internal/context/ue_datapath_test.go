package context_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/free5gc/smf/internal/context"
	"github.com/free5gc/smf/pkg/factory"
)

var config = configuration

// smfContext.UserPlaneInformation = NewUserPlaneInformation(config)

func TestNewUEPreConfigPaths(t *testing.T) {
	smfContext := context.GetSelf()
	smfContext.UserPlaneInformation = context.NewUserPlaneInformation(config)
	fmt.Println("Start")
	testcases := []struct {
		name                  string
		inPaths               []factory.SpecificPath
		expectedDataPathNodes [][]*context.UPF
	}{
		{
			name: "singlePath-singleUPF",
			inPaths: []factory.SpecificPath{
				{
					DestinationIP:   "10.60.0.101/32",
					DestinationPort: "12345",
					Path: []string{
						"UPF1",
					},
				},
			},
			expectedDataPathNodes: [][]*context.UPF{
				{
					getUpf("UPF1"),
				},
			},
		},
		{
			name: "singlePath-multiUPF",
			inPaths: []factory.SpecificPath{
				{
					DestinationIP:   "10.60.0.101/32",
					DestinationPort: "12345",
					Path: []string{
						"UPF1",
						"UPF2",
					},
				},
			},
			expectedDataPathNodes: [][]*context.UPF{
				{
					getUpf("UPF1"),
					getUpf("UPF2"),
				},
			},
		},
		{
			name: "multiPath-singleUPF",
			inPaths: []factory.SpecificPath{
				{
					DestinationIP:   "10.60.0.101/32",
					DestinationPort: "12345",
					Path: []string{
						"UPF1",
					},
				},
				{
					DestinationIP:   "10.60.0.103/32",
					DestinationPort: "12345",
					Path: []string{
						"UPF2",
					},
				},
			},
			expectedDataPathNodes: [][]*context.UPF{
				{
					getUpf("UPF1"),
				},
				{
					getUpf("UPF2"),
				},
			},
		},
		{
			name: "multiPath-multiUPF",
			inPaths: []factory.SpecificPath{
				{
					DestinationIP:   "10.60.0.101/32",
					DestinationPort: "12345",
					Path: []string{
						"UPF1",
						"UPF2",
					},
				},
				{
					DestinationIP:   "10.60.0.103/32",
					DestinationPort: "12345",
					Path: []string{
						"UPF1",
						"UPF3",
					},
				},
			},
			expectedDataPathNodes: [][]*context.UPF{
				{
					getUpf("UPF1"),
					getUpf("UPF2"),
				},
				{
					getUpf("UPF1"),
					getUpf("UPF3"),
				},
			},
		},
		{
			name: "multiPath-single&multiUPF",
			inPaths: []factory.SpecificPath{
				{
					DestinationIP:   "10.60.0.101/32",
					DestinationPort: "12345",
					Path: []string{
						"UPF1",
					},
				},
				{
					DestinationIP:   "10.60.0.103/32",
					DestinationPort: "12345",
					Path: []string{
						"UPF1",
						"UPF3",
					},
				},
			},
			expectedDataPathNodes: [][]*context.UPF{
				{
					getUpf("UPF1"),
				},
				{
					getUpf("UPF1"),
					getUpf("UPF3"),
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			retUePreConfigPaths, err := context.NewUEPreConfigPaths(tc.inPaths)
			require.Nil(t, err)
			require.NotNil(t, retUePreConfigPaths.PathIDGenerator)
			for pathIndex, path := range tc.inPaths {
				retDataPath := retUePreConfigPaths.DataPathPool[int64(pathIndex+1)]
				require.Equal(t, path.DestinationIP, retDataPath.Destination.DestinationIP)
				require.Equal(t, path.DestinationPort, retDataPath.Destination.DestinationPort)
				retNode := retDataPath.FirstDPNode
				for _, expectedUpf := range tc.expectedDataPathNodes[pathIndex] {
					require.NotNil(t, retNode.UPF)
					require.Equal(t, retNode.UPF, expectedUpf)
					retNode = retNode.DownLinkTunnel.SrcEndPoint
				}
				require.Nil(t, retNode)
			}
		})
	}
}

func getUpf(name string) *context.UPF {
	newUeNode, err := context.NewUEDataPathNode(name)
	if err != nil {
		return nil
	}

	Upf := newUeNode.UPF

	return Upf
}
