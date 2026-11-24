package generator_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/free5gc/openapi/models"

	amffactory "github.com/free5gc/amf/pkg/factory"
	"github.com/free5gc/config/cfg-generator/generator"
)

func TestMarshal(t *testing.T) {
	cfg := &amffactory.Config{
		Info: &amffactory.Info{Version: "1.0.9"},
		Configuration: &amffactory.Configuration{
			SupportTAIList: []models.Tai{{PlmnId: &models.PlmnId{Mcc: "001", Mnc: "01"}, Tac: "000001"}},
			PlmnSupportList: []amffactory.PlmnSupportItem{{
				PlmnId:     &models.PlmnId{Mcc: "001", Mnc: "01"},
				SNssaiList: []models.Snssai{{Sst: 1, Sd: "010203"}},
			}},
		},
	}

	out, err := generator.Marshal(cfg)
	require.NoError(t, err)
	s := string(out)

	t.Run("first line is the generated marker", func(t *testing.T) {
		first, _, _ := strings.Cut(s, "\n")
		assert.Equal(t, generator.GeneratedMarker, first)
	})

	t.Run("embedded mutex is not serialized", func(t *testing.T) {
		assert.NotContains(t, s, "rwmutex")
	})

	t.Run("leading-zero strings are quoted", func(t *testing.T) {
		assert.Contains(t, s, `sd: "010203"`)
		assert.Contains(t, s, `tac: "000001"`)
		assert.Contains(t, s, `mnc: "01"`)
		assert.Contains(t, s, `mcc: "001"`)
	})

	t.Run("round trip preserves values", func(t *testing.T) {
		back := &amffactory.Config{}
		require.NoError(t, yaml.Unmarshal(out, back))
		assert.Equal(t, "010203", back.Configuration.PlmnSupportList[0].SNssaiList[0].Sd)
		assert.Equal(t, "000001", back.Configuration.SupportTAIList[0].Tac)
		assert.Equal(t, "01", back.Configuration.SupportTAIList[0].PlmnId.Mnc)
		assert.Equal(t, "001", back.Configuration.SupportTAIList[0].PlmnId.Mcc)
	})

	t.Run("top-level key order follows the struct", func(t *testing.T) {
		assert.Less(t, strings.Index(s, "\ninfo:"), strings.Index(s, "\nconfiguration:"))
	})
}

func TestMarshalError(t *testing.T) {
	// yaml.v2 panics on some types and returns errors on others; a channel
	// is rejected with a panic that yaml.Marshal recovers into an error.
	_, err := generator.Marshal(map[string]any{"c": make(chan int)})
	assert.Error(t, err)
}
