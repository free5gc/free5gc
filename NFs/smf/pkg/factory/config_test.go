package factory_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/smf/pkg/factory"
)

func TestSnssaiInfoItem(t *testing.T) {
	testcase := []struct {
		Name     string
		Snssai   *models.Snssai
		DnnInfos []*factory.SnssaiDnnInfoItem
	}{
		{
			Name: "Default",
			Snssai: &models.Snssai{
				Sst: int32(1),
				Sd:  "010203",
			},
			DnnInfos: []*factory.SnssaiDnnInfoItem{
				{
					Dnn: "internet",
					DNS: &factory.DNS{
						IPv4Addr: "8.8.8.8",
					},
				},
			},
		},
		{
			Name: "Empty SD",
			Snssai: &models.Snssai{
				Sst: int32(1),
			},
			DnnInfos: []*factory.SnssaiDnnInfoItem{
				{
					Dnn: "internet2",
					DNS: &factory.DNS{
						IPv4Addr: "1.1.1.1",
					},
				},
			},
		},
	}

	for _, tc := range testcase {
		t.Run(tc.Name, func(t *testing.T) {
			snssaiInfoItem := factory.SnssaiInfoItem{
				SNssai:   tc.Snssai,
				DnnInfos: tc.DnnInfos,
			}

			ok, err := snssaiInfoItem.Validate()
			require.True(t, ok)
			require.Nil(t, err)
		})
	}
}

func TestSnssaiUpfInfoItem(t *testing.T) {
	testcase := []struct {
		Name     string
		Snssai   *models.Snssai
		DnnInfos []*factory.DnnUpfInfoItem
	}{
		{
			Name: "Default",
			Snssai: &models.Snssai{
				Sst: int32(1),
				Sd:  "010203",
			},
			DnnInfos: []*factory.DnnUpfInfoItem{
				{
					Dnn: "internet",
				},
			},
		},
		{
			Name: "Empty SD",
			Snssai: &models.Snssai{
				Sst: int32(1),
			},
			DnnInfos: []*factory.DnnUpfInfoItem{
				{
					Dnn: "internet2",
				},
			},
		},
	}

	for _, tc := range testcase {
		t.Run(tc.Name, func(t *testing.T) {
			snssaiInfoItem := factory.SnssaiUpfInfoItem{
				SNssai:         tc.Snssai,
				DnnUpfInfoList: tc.DnnInfos,
			}

			ok, err := snssaiInfoItem.Validate()
			require.True(t, ok)
			require.Nil(t, err)
		})
	}
}
