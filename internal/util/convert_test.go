package util_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/free5gc/nrf/internal/util"
)

func TestSnssaisToBsonM(t *testing.T) {
	testCases := []struct {
		Name        string
		snssais     string
		expectBsonM []bson.M
	}{
		{
			Name:    "Default 01010203",
			snssais: "{\"sst\":1,\"sd\":\"010203\"}",
			expectBsonM: []bson.M{
				{
					"sst": int32(1),
					"sd":  "010203",
				},
			},
		},
		{
			Name:    "Empty SD",
			snssais: "{\"sst\":1}",
			expectBsonM: []bson.M{
				{
					"sst": int32(1),
				},
			},
		},
		{
			Name:    "Two Slices",
			snssais: "{\"sst\":1,\"sd\":\"010203\"},{\"sst\":1,\"sd\":\"112233\"}",
			expectBsonM: []bson.M{
				{
					"sst": int32(1),
					"sd":  "010203",
				},
				{
					"sst": int32(1),
					"sd":  "112233",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			output := util.SnssaisToBsonM(tc.snssais)
			require.Equal(t, tc.expectBsonM, output)
		})
	}
}
