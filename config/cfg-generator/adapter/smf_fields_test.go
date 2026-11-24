package adapter_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/free5gc/openapi/models"

	"github.com/free5gc/config/cfg-generator/adapter"
	"github.com/free5gc/config/cfg-generator/global"
	smffactory "github.com/free5gc/smf/pkg/factory"
)

var (
	sliceA = models.Snssai{Sst: 1, Sd: "010203"}
	sliceB = models.Snssai{Sst: 1, Sd: "112233"}
)

func snssaiPtr(s models.Snssai) *models.Snssai { return &s }

// newFullSMFConfig is shaped like the sample smfcfg.yaml, with a second UPF
// ("AUPF", sorted before "UPF") to exercise deterministic map traversal.
func newFullSMFConfig() *smffactory.Config {
	upfInfos := func() []*smffactory.SnssaiUpfInfoItem {
		return []*smffactory.SnssaiUpfInfoItem{
			{SNssai: snssaiPtr(sliceA), DnnUpfInfoList: []*smffactory.DnnUpfInfoItem{{Dnn: "internet"}}},
			{SNssai: snssaiPtr(sliceB), DnnUpfInfoList: []*smffactory.DnnUpfInfoItem{{Dnn: "internet"}}},
		}
	}
	return &smffactory.Config{
		Configuration: &smffactory.Configuration{
			SmfName:  "SMF",
			Sbi:      &smffactory.Sbi{Scheme: "http"},
			NrfUri:   "http://127.0.0.10:8000",
			PLMNList: []smffactory.PlmnID{{Mcc: "208", Mnc: "93"}},
			SNssaiInfo: []*smffactory.SnssaiInfoItem{
				{SNssai: snssaiPtr(sliceA), DnnInfos: []*smffactory.SnssaiDnnInfoItem{{Dnn: "internet"}}},
				{SNssai: snssaiPtr(sliceB), DnnInfos: []*smffactory.SnssaiDnnInfoItem{{Dnn: "internet"}}},
			},
			UserPlaneInformation: smffactory.UserPlaneInformation{
				UPNodes: map[string]*smffactory.UPNode{
					"gNB1": {Type: "AN"},
					"UPF": {
						Type:        "UPF",
						SNssaiInfos: upfInfos(),
						InterfaceUpfInfoList: []*smffactory.InterfaceUpfInfoItem{
							{InterfaceType: "N3", NetworkInstances: []string{"internet"}},
						},
					},
					"AUPF": {
						Type:        "UPF",
						SNssaiInfos: upfInfos(),
						InterfaceUpfInfoList: []*smffactory.InterfaceUpfInfoItem{
							{InterfaceType: "N9", NetworkInstances: []string{"internet", "transport"}},
						},
					},
				},
				Links: []*smffactory.UPLink{{A: "gNB1", B: "AUPF"}, {A: "AUPF", B: "UPF"}},
			},
			IdleTimeout:     ptr(60 * time.Second),
			ReadIdleTimeout: ptr(60 * time.Second),
			PingTimeout:     ptr(15 * time.Second),
		},
		Logger: &smffactory.Logger{Enable: true, Level: "info"},
	}
}

func applySMF(t *testing.T, g *global.Global, cfg *smffactory.Config) []adapter.Change {
	t.Helper()
	r := &adapter.Recorder{}
	require.NoError(t, adapter.SMFAdapter{}.Apply(g, cfg, r))
	return r.Changes
}

func TestSMFApplySnssaiKeepsSessionAndUserPlaneAligned(t *testing.T) {
	cfg := newFullSMFConfig()
	newA, newB := models.Snssai{Sst: 1, Sd: "000001"}, models.Snssai{Sst: 1, Sd: "000002"}
	changes := applySMF(t, &global.Global{Snssais: []models.Snssai{newA, newB}}, cfg)

	c := cfg.Configuration
	assert.Equal(t, newA, *c.SNssaiInfo[0].SNssai)
	assert.Equal(t, newB, *c.SNssaiInfo[1].SNssai)
	for _, name := range []string{"UPF", "AUPF"} {
		infos := c.UserPlaneInformation.UPNodes[name].SNssaiInfos
		assert.Equal(t, newA, *infos[0].SNssai, name)
		assert.Equal(t, newB, *infos[1].SNssai, name)
	}
	// 2 snssaiInfos + 2 nodes x 2 sNssaiUpfInfos, sd only.
	assert.Len(t, changes, 6)
	assert.Equal(t, "snssaiInfos[0].sNssai.sd", changes[0].Path)
	assert.Equal(t, "userplaneInformation.upNodes.AUPF.sNssaiUpfInfos[0].sNssai.sd", changes[2].Path,
		"nodes visited in sorted order")
	assert.Equal(t, "userplaneInformation.upNodes.UPF.sNssaiUpfInfos[1].sNssai.sd", changes[5].Path)
}

func TestSMFApplyIsDeterministic(t *testing.T) {
	g := &global.Global{
		Snssais: []models.Snssai{{Sst: 2, Sd: "000001"}, {Sst: 2, Sd: "000002"}},
		Dnns:    []string{"ims"},
	}
	want := applySMF(t, g, newFullSMFConfig())
	for i := 0; i < 20; i++ {
		assert.Equal(t, want, applySMF(t, g, newFullSMFConfig()), "run %d", i)
	}
}

func TestSMFApplySnssaiOnlyInUserPlaneIsCollectedLast(t *testing.T) {
	cfg := newFullSMFConfig()
	extra := models.Snssai{Sst: 3, Sd: "333333"}
	upf := cfg.Configuration.UserPlaneInformation.UPNodes["UPF"]
	upf.SNssaiInfos = append(upf.SNssaiInfos, &smffactory.SnssaiUpfInfoItem{SNssai: snssaiPtr(extra)})

	// Distinct order: sliceA, sliceB (snssaiInfos), then extra (UPF node).
	r := &adapter.Recorder{}
	err := adapter.SMFAdapter{}.Apply(&global.Global{Snssais: []models.Snssai{sliceA, sliceB}}, cfg, r)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "has 3 distinct value(s)")

	g := &global.Global{Snssais: []models.Snssai{sliceA, sliceB, {Sst: 4, Sd: "444444"}}}
	applySMF(t, g, cfg)
	assert.Equal(t, models.Snssai{Sst: 4, Sd: "444444"}, *upf.SNssaiInfos[2].SNssai)
}

func TestSMFApplyDnn(t *testing.T) {
	cfg := newFullSMFConfig()
	changes := applySMF(t, &global.Global{Dnns: []string{"ims"}}, cfg)

	c := cfg.Configuration
	for _, info := range c.SNssaiInfo {
		assert.Equal(t, "ims", info.DnnInfos[0].Dnn)
	}
	for _, name := range []string{"UPF", "AUPF"} {
		node := c.UserPlaneInformation.UPNodes[name]
		for _, s := range node.SNssaiInfos {
			assert.Equal(t, "ims", s.DnnUpfInfoList[0].Dnn, name)
		}
	}
	assert.Equal(t, []string{"ims"}, c.UserPlaneInformation.UPNodes["UPF"].InterfaceUpfInfoList[0].NetworkInstances)
	assert.Equal(t, []string{"ims", "transport"},
		c.UserPlaneInformation.UPNodes["AUPF"].InterfaceUpfInfoList[0].NetworkInstances,
		"network instance that is not a DNN is kept")
	// 2 dnnInfos + 2 nodes x 2 dnnUpfInfoList + 2 networkInstances.
	assert.Len(t, changes, 8)
}

func TestSMFApplyValueMappingMismatchChangesNothing(t *testing.T) {
	cfg := newFullSMFConfig()
	r := &adapter.Recorder{}
	err := adapter.SMFAdapter{}.Apply(&global.Global{
		Plmn: &models.PlmnId{Mcc: "466", Mnc: "92"},
		Dnns: []string{"a", "b"},
	}, cfg, r)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "global dnnList has 2 entries but smf snssaiInfos/userplaneInformation has 1 distinct value(s) [internet]")
	assert.Empty(t, r.Changes)
	assert.Equal(t, newFullSMFConfig(), cfg)
}

func TestSMFApplyCommonFields(t *testing.T) {
	cfg := newFullSMFConfig()
	changes := applySMF(t, &global.Global{
		Nrf:       &global.NrfRef{Scheme: "http", Host: "nrf.local", Port: 8000},
		SbiScheme: ptr("http"), // same: not recorded
		Logger:    &global.Logger{Enable: ptr(false)},
		Timeouts:  &global.Timeouts{ReadIdle: ptr(2 * time.Minute)},
	}, cfg)

	assert.Equal(t, []adapter.Change{
		{NF: adapter.SMF, Path: "logger.enable", Old: "true", New: "false"},
		{NF: adapter.SMF, Path: "readIdleTimeout", Old: "1m0s", New: "2m0s"},
		{NF: adapter.SMF, Path: "nrfUri", Old: "http://127.0.0.10:8000", New: "http://nrf.local:8000"},
	}, changes)
}

func TestSMFApplyToleratesNilEntries(t *testing.T) {
	cfg := newFullSMFConfig()
	c := cfg.Configuration
	c.SNssaiInfo = append(c.SNssaiInfo, nil, &smffactory.SnssaiInfoItem{DnnInfos: []*smffactory.SnssaiDnnInfoItem{nil}})
	c.UserPlaneInformation.UPNodes["nilnode"] = nil
	upf := c.UserPlaneInformation.UPNodes["UPF"]
	upf.SNssaiInfos = append(upf.SNssaiInfos, nil, &smffactory.SnssaiUpfInfoItem{DnnUpfInfoList: []*smffactory.DnnUpfInfoItem{nil}})
	upf.InterfaceUpfInfoList = append(upf.InterfaceUpfInfoList, nil)

	changes := applySMF(t, &global.Global{
		Snssais: []models.Snssai{{Sst: 1, Sd: "000001"}, {Sst: 1, Sd: "000002"}},
		Dnns:    []string{"ims"},
	}, cfg)
	assert.NotEmpty(t, changes)
}

func TestSMFApplyMissingLoggerAndSbi(t *testing.T) {
	cfg := newFullSMFConfig()
	cfg.Logger = nil
	cfg.Configuration.Sbi = nil

	changes := applySMF(t, &global.Global{SbiScheme: ptr("https"), Logger: &global.Logger{Level: ptr("error")}}, cfg)
	assert.Nil(t, cfg.Configuration.Sbi, "no sbi block, sbiScheme skipped")
	require.NotNil(t, cfg.Logger)
	assert.Equal(t, "error", cfg.Logger.Level)
	assert.Equal(t, []adapter.Change{{NF: adapter.SMF, Path: "logger.level", Old: "", New: "error"}}, changes)
}
