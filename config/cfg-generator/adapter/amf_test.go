package adapter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/free5gc/openapi/models"

	amffactory "github.com/free5gc/amf/pkg/factory"
	"github.com/free5gc/config/cfg-generator/adapter"
	"github.com/free5gc/config/cfg-generator/global"
)

// newAMFConfig returns an AMF config shaped like the sample amfcfg.yaml:
// one GUAMI, one TAI, one PLMN support item with two slices.
func newAMFConfig() *amffactory.Config {
	return &amffactory.Config{
		Configuration: &amffactory.Configuration{
			ServedGumaiList: []models.Guami{{
				PlmnId: &models.PlmnIdNid{Mcc: "208", Mnc: "93", Nid: "keep-nid"},
				AmfId:  "cafe00",
			}},
			SupportTAIList: []models.Tai{{
				PlmnId: &models.PlmnId{Mcc: "208", Mnc: "93"},
				Tac:    "000001",
			}},
			PlmnSupportList: []amffactory.PlmnSupportItem{{
				PlmnId: &models.PlmnId{Mcc: "208", Mnc: "93"},
				SNssaiList: []models.Snssai{
					{Sst: 1, Sd: "010203"},
					{Sst: 1, Sd: "112233"},
				},
			}},
		},
	}
}

func TestAMFAdapterMetadata(t *testing.T) {
	assert.Equal(t, adapter.AMF, adapter.AMFAdapter{}.NF())
	assert.Equal(t, "amfcfg.yaml", adapter.AMFAdapter{}.FileName())
}

func TestAMFAdapterApplyPlmn(t *testing.T) {
	cfg := newAMFConfig()
	r := &adapter.Recorder{}

	err := adapter.AMFAdapter{}.Apply(
		&global.Global{Plmn: &models.PlmnId{Mcc: "466", Mnc: "01"}}, cfg, r)
	require.NoError(t, err)

	// Exactly these paths change, no more, no less.
	assert.Equal(t, []adapter.Change{
		{NF: adapter.AMF, Path: "servedGuamiList[0].plmnId.mcc", Old: "208", New: "466"},
		{NF: adapter.AMF, Path: "servedGuamiList[0].plmnId.mnc", Old: "93", New: "01"},
		{NF: adapter.AMF, Path: "supportTaiList[0].plmnId.mcc", Old: "208", New: "466"},
		{NF: adapter.AMF, Path: "supportTaiList[0].plmnId.mnc", Old: "93", New: "01"},
		{NF: adapter.AMF, Path: "plmnSupportList[0].plmnId.mcc", Old: "208", New: "466"},
		{NF: adapter.AMF, Path: "plmnSupportList[0].plmnId.mnc", Old: "93", New: "01"},
	}, r.Changes)

	c := cfg.Configuration
	assert.Equal(t, models.PlmnIdNid{Mcc: "466", Mnc: "01", Nid: "keep-nid"}, *c.ServedGumaiList[0].PlmnId,
		"nid must be preserved")
	assert.Equal(t, "cafe00", c.ServedGumaiList[0].AmfId, "amfId must be preserved")
	assert.Equal(t, models.PlmnId{Mcc: "466", Mnc: "01"}, *c.SupportTAIList[0].PlmnId)
	assert.Equal(t, "000001", c.SupportTAIList[0].Tac, "tac must be preserved")
	assert.Equal(t, models.PlmnId{Mcc: "466", Mnc: "01"}, *c.PlmnSupportList[0].PlmnId)
	assert.Equal(t, []models.Snssai{{Sst: 1, Sd: "010203"}, {Sst: 1, Sd: "112233"}},
		c.PlmnSupportList[0].SNssaiList, "slices must be preserved")
}

func TestAMFAdapterApplyEveryListEntry(t *testing.T) {
	cfg := newAMFConfig()
	c := cfg.Configuration
	c.ServedGumaiList = append(c.ServedGumaiList,
		models.Guami{PlmnId: &models.PlmnIdNid{Mcc: "208", Mnc: "93"}, AmfId: "cafe01"})
	c.SupportTAIList = append(c.SupportTAIList,
		models.Tai{PlmnId: &models.PlmnId{Mcc: "208", Mnc: "93"}, Tac: "000002"})

	r := &adapter.Recorder{}
	require.NoError(t, adapter.AMFAdapter{}.Apply(
		&global.Global{Plmn: &models.PlmnId{Mcc: "001", Mnc: "01"}}, cfg, r))

	assert.Len(t, r.Changes, 10, "2 guami + 2 tai + 1 plmnSupport entries, mcc+mnc each")
	for _, g := range c.ServedGumaiList {
		assert.Equal(t, "001", g.PlmnId.Mcc)
		assert.Equal(t, "01", g.PlmnId.Mnc)
	}
	for _, tai := range c.SupportTAIList {
		assert.Equal(t, "001", tai.PlmnId.Mcc)
	}
	assert.Equal(t, "cafe01", c.ServedGumaiList[1].AmfId)
	assert.Equal(t, "000002", c.SupportTAIList[1].Tac)
}

func TestAMFAdapterApplyCreatesMissingPlmnId(t *testing.T) {
	cfg := newAMFConfig()
	c := cfg.Configuration
	c.ServedGumaiList[0].PlmnId = nil
	c.SupportTAIList[0].PlmnId = nil
	c.PlmnSupportList[0].PlmnId = nil

	r := &adapter.Recorder{}
	require.NoError(t, adapter.AMFAdapter{}.Apply(
		&global.Global{Plmn: &models.PlmnId{Mcc: "466", Mnc: "92"}}, cfg, r))

	require.NotNil(t, c.ServedGumaiList[0].PlmnId)
	assert.Equal(t, "466", c.ServedGumaiList[0].PlmnId.Mcc)
	require.NotNil(t, c.SupportTAIList[0].PlmnId)
	assert.Equal(t, "92", c.SupportTAIList[0].PlmnId.Mnc)
	require.NotNil(t, c.PlmnSupportList[0].PlmnId)
	assert.Equal(t, "466", c.PlmnSupportList[0].PlmnId.Mcc)
	assert.Equal(t, `amf: servedGuamiList[0].plmnId.mcc: "" -> "466"`, r.Changes[0].String())
}

func TestAMFAdapterApplyNoOp(t *testing.T) {
	cases := []struct {
		name string
		g    *global.Global
	}{
		{name: "nil global", g: nil},
		{name: "plmn not provided", g: &global.Global{}},
		{name: "same plmn as config", g: &global.Global{Plmn: &models.PlmnId{Mcc: "208", Mnc: "93"}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := newAMFConfig()
			r := &adapter.Recorder{}
			require.NoError(t, adapter.AMFAdapter{}.Apply(tc.g, cfg, r))
			assert.Empty(t, r.Changes)
			assert.Equal(t, newAMFConfig(), cfg, "config must be untouched")
		})
	}
}

func TestAMFAdapterApplyMissingConfiguration(t *testing.T) {
	err := adapter.AMFAdapter{}.Apply(
		&global.Global{Plmn: &models.PlmnId{Mcc: "466", Mnc: "92"}},
		&amffactory.Config{}, &adapter.Recorder{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "configuration block is missing")
}

func TestAMFAdapterValidate(t *testing.T) {
	// An almost-empty config must be rejected by the AMF's own Validate.
	err := adapter.AMFAdapter{}.Validate(&amffactory.Config{})
	assert.Error(t, err)
}
