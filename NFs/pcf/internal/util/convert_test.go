package util

import (
	"testing"

	"github.com/go-playground/assert/v2"

	"github.com/free5gc/openapi/models"
)

func TestSnssaiModelsToHex(t *testing.T) {
	tests := []struct {
		name         string
		snssai       models.Snssai
		expectString string
	}{
		{
			name: "01010202",
			snssai: models.Snssai{
				Sst: 1,
				Sd:  "010203",
			},
			expectString: "01010203",
		},
		{
			name: "Empty SD",
			snssai: models.Snssai{
				Sst: 1,
			},
			expectString: "01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := SnssaiModelsToHex(tt.snssai)
			assert.Equal(t, tt.expectString, output)
		})
	}
}
