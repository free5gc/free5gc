package adapter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/free5gc/openapi/models"

	"github.com/free5gc/config/cfg-generator/adapter"
	"github.com/free5gc/config/cfg-generator/global"
	smffactory "github.com/free5gc/smf/pkg/factory"
)

func newSMFConfig() *smffactory.Config {
	return &smffactory.Config{
		Configuration: &smffactory.Configuration{
			SmfName:  "SMF",
			PLMNList: []smffactory.PlmnID{{Mcc: "208", Mnc: "93"}},
			SNssaiInfo: []*smffactory.SnssaiInfoItem{{
				SNssai: &models.Snssai{Sst: 1, Sd: "010203"},
			}},
		},
	}
}

func TestSMFAdapterMetadata(t *testing.T) {
	assert.Equal(t, adapter.SMF, adapter.SMFAdapter{}.NF())
	assert.Equal(t, "smfcfg.yaml", adapter.SMFAdapter{}.FileName())
}

func TestSMFAdapterApplyPlmn(t *testing.T) {
	cfg := newSMFConfig()
	cfg.Configuration.PLMNList = append(cfg.Configuration.PLMNList, smffactory.PlmnID{Mcc: "208", Mnc: "94"})
	r := &adapter.Recorder{}

	require.NoError(t, adapter.SMFAdapter{}.Apply(
		&global.Global{Plmn: &models.PlmnId{Mcc: "466", Mnc: "01"}}, cfg, r))

	assert.Equal(t, []adapter.Change{
		{NF: adapter.SMF, Path: "plmnList[0].mcc", Old: "208", New: "466"},
		{NF: adapter.SMF, Path: "plmnList[0].mnc", Old: "93", New: "01"},
		{NF: adapter.SMF, Path: "plmnList[1].mcc", Old: "208", New: "466"},
		{NF: adapter.SMF, Path: "plmnList[1].mnc", Old: "94", New: "01"},
	}, r.Changes)
	assert.Equal(t, []smffactory.PlmnID{{Mcc: "466", Mnc: "01"}, {Mcc: "466", Mnc: "01"}},
		cfg.Configuration.PLMNList)
	assert.Equal(t, "010203", cfg.Configuration.SNssaiInfo[0].SNssai.Sd, "slices must be preserved")
}

func TestSMFAdapterApplyNoOp(t *testing.T) {
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
			cfg := newSMFConfig()
			r := &adapter.Recorder{}
			require.NoError(t, adapter.SMFAdapter{}.Apply(tc.g, cfg, r))
			assert.Empty(t, r.Changes)
			assert.Equal(t, newSMFConfig(), cfg)
		})
	}
}

func TestSMFAdapterApplyEmptyPlmnListStaysEmpty(t *testing.T) {
	// plmnList is optional in SMF; the adapter keeps entries, it does not
	// invent them (docs/global-mapping.md §0, list rule (a)).
	cfg := newSMFConfig()
	cfg.Configuration.PLMNList = nil
	r := &adapter.Recorder{}
	require.NoError(t, adapter.SMFAdapter{}.Apply(
		&global.Global{Plmn: &models.PlmnId{Mcc: "466", Mnc: "01"}}, cfg, r))
	assert.Empty(t, r.Changes)
	assert.Nil(t, cfg.Configuration.PLMNList)
}

func TestSMFAdapterApplyMissingConfiguration(t *testing.T) {
	err := adapter.SMFAdapter{}.Apply(
		&global.Global{Plmn: &models.PlmnId{Mcc: "466", Mnc: "92"}},
		&smffactory.Config{}, &adapter.Recorder{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "configuration block is missing")
}

func TestSMFAdapterValidate(t *testing.T) {
	assert.Error(t, adapter.SMFAdapter{}.Validate(&smffactory.Config{}))
}
