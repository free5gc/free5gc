package adapter_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/free5gc/openapi/models"

	amffactory "github.com/free5gc/amf/pkg/factory"
	"github.com/free5gc/config/cfg-generator/adapter"
	"github.com/free5gc/config/cfg-generator/global"
)

func ptr[T any](v T) *T { return &v }

// newFullAMFConfig extends newAMFConfig with every field the global block
// can touch.
func newFullAMFConfig() *amffactory.Config {
	cfg := newAMFConfig()
	c := cfg.Configuration
	c.SupportTAIList = append(c.SupportTAIList,
		models.Tai{PlmnId: &models.PlmnId{Mcc: "208", Mnc: "93"}, Tac: "000002"})
	c.SupportDnnList = []string{"internet", "ims"}
	c.SupportLadnList = []amffactory.Ladn{{Dnn: "ims"}, {Dnn: "local-only"}}
	c.Sbi = &amffactory.Sbi{Scheme: "http"}
	c.NrfUri = "http://127.0.0.10:8000"
	c.IdleTimeout = ptr(60 * time.Second)
	c.ReadIdleTimeout = ptr(60 * time.Second)
	c.PingTimeout = ptr(15 * time.Second)
	cfg.Logger = &amffactory.Logger{Enable: true, Level: "info", ReportCaller: false}
	return cfg
}

func applyAMF(t *testing.T, g *global.Global, cfg *amffactory.Config) []adapter.Change {
	t.Helper()
	r := &adapter.Recorder{}
	require.NoError(t, adapter.AMFAdapter{}.Apply(g, cfg, r))
	return r.Changes
}

func TestAMFApplyTac(t *testing.T) {
	cfg := newFullAMFConfig()
	changes := applyAMF(t, &global.Global{Tac: ptr("00ABCD")}, cfg)

	assert.Equal(t, []adapter.Change{
		{NF: adapter.AMF, Path: "supportTaiList[0].tac", Old: "000001", New: "00ABCD"},
		{NF: adapter.AMF, Path: "supportTaiList[1].tac", Old: "000002", New: "00ABCD"},
	}, changes)
	assert.Equal(t, "208", cfg.Configuration.SupportTAIList[1].PlmnId.Mcc, "plmn untouched")
}

func TestAMFApplySnssai(t *testing.T) {
	cfg := newFullAMFConfig()
	// A second PLMN support item repeating one slice: still 2 distinct values.
	cfg.Configuration.PlmnSupportList = append(cfg.Configuration.PlmnSupportList,
		amffactory.PlmnSupportItem{
			PlmnId:     &models.PlmnId{Mcc: "208", Mnc: "93"},
			SNssaiList: []models.Snssai{{Sst: 1, Sd: "112233"}},
		})

	changes := applyAMF(t, &global.Global{Snssais: []models.Snssai{
		{Sst: 1, Sd: "000001"},
		{Sst: 2, Sd: "112233"},
	}}, cfg)

	pl := cfg.Configuration.PlmnSupportList
	assert.Equal(t, []models.Snssai{{Sst: 1, Sd: "000001"}, {Sst: 2, Sd: "112233"}}, pl[0].SNssaiList)
	assert.Equal(t, []models.Snssai{{Sst: 2, Sd: "112233"}}, pl[1].SNssaiList, "same value mapped the same way")
	assert.Equal(t, []adapter.Change{
		{NF: adapter.AMF, Path: "plmnSupportList[0].snssaiList[0].sd", Old: "010203", New: "000001"},
		{NF: adapter.AMF, Path: "plmnSupportList[0].snssaiList[1].sst", Old: "1", New: "2"},
		{NF: adapter.AMF, Path: "plmnSupportList[1].snssaiList[0].sst", Old: "1", New: "2"},
	}, changes)
}

func TestAMFApplyDnn(t *testing.T) {
	cfg := newFullAMFConfig()
	changes := applyAMF(t, &global.Global{Dnns: []string{"internet2", "ims2"}}, cfg)

	c := cfg.Configuration
	assert.Equal(t, []string{"internet2", "ims2"}, c.SupportDnnList)
	assert.Equal(t, "ims2", c.SupportLadnList[0].Dnn, "LADN follows the DNN mapping")
	assert.Equal(t, "local-only", c.SupportLadnList[1].Dnn, "LADN DNN not in supportDnnList is kept")
	assert.Equal(t, []adapter.Change{
		{NF: adapter.AMF, Path: "supportDnnList[0]", Old: "internet", New: "internet2"},
		{NF: adapter.AMF, Path: "supportDnnList[1]", Old: "ims", New: "ims2"},
		{NF: adapter.AMF, Path: "supportLadnList[0].dnn", Old: "ims", New: "ims2"},
	}, changes)
}

func TestAMFApplyValueMappingMismatchChangesNothing(t *testing.T) {
	cases := []struct {
		name    string
		g       *global.Global
		wantMsg string
	}{
		{
			name: "snssai count differs",
			g: &global.Global{
				Plmn:    &models.PlmnId{Mcc: "466", Mnc: "92"}, // would change if applied
				Snssais: []models.Snssai{{Sst: 1, Sd: "000001"}},
			},
			wantMsg: "global snssaiList has 1 entries but amf plmnSupportList[*].snssaiList has 2 distinct value(s)",
		},
		{
			name: "dnn count differs",
			g: &global.Global{
				Tac:  ptr("000009"),
				Dnns: []string{"a", "b", "c"},
			},
			wantMsg: "global dnnList has 3 entries but amf supportDnnList has 2 distinct value(s)",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := newFullAMFConfig()
			r := &adapter.Recorder{}
			err := adapter.AMFAdapter{}.Apply(tc.g, cfg, r)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantMsg)
			assert.Empty(t, r.Changes)
			assert.Equal(t, newFullAMFConfig(), cfg, "no partial changes")
		})
	}
}

func TestAMFApplyCommonFields(t *testing.T) {
	cfg := newFullAMFConfig()
	changes := applyAMF(t, &global.Global{
		Nrf:       &global.NrfRef{Scheme: "https", Host: "10.0.0.10", Port: 29510},
		SbiScheme: ptr("https"),
		Logger:    &global.Logger{Level: ptr("debug"), ReportCaller: ptr(true)},
		Timeouts:  &global.Timeouts{Idle: ptr(30 * time.Second), Ping: ptr(15 * time.Second)},
	}, cfg)

	assert.Equal(t, []adapter.Change{
		{NF: adapter.AMF, Path: "logger.level", Old: "info", New: "debug"},
		{NF: adapter.AMF, Path: "logger.reportCaller", Old: "false", New: "true"},
		{NF: adapter.AMF, Path: "idleTimeout", Old: "1m0s", New: "30s"},
		{NF: adapter.AMF, Path: "sbi.scheme", Old: "http", New: "https"},
		{NF: adapter.AMF, Path: "nrfUri", Old: "http://127.0.0.10:8000", New: "https://10.0.0.10:29510"},
	}, changes, "ping is unchanged (same value), logger.enable not provided")

	c := cfg.Configuration
	assert.Equal(t, "https://10.0.0.10:29510", c.NrfUri)
	assert.Equal(t, "https", c.Sbi.Scheme)
	assert.Equal(t, 30*time.Second, *c.IdleTimeout)
	assert.Equal(t, 60*time.Second, *c.ReadIdleTimeout, "not provided, kept")
	assert.Equal(t, amffactory.Logger{Enable: true, Level: "debug", ReportCaller: true}, *cfg.Logger)
}

func TestAMFApplyCommonFieldsOnSparseConfig(t *testing.T) {
	t.Run("missing logger is created only when global logger is given", func(t *testing.T) {
		cfg := newFullAMFConfig()
		cfg.Logger = nil
		changes := applyAMF(t, &global.Global{Logger: &global.Logger{Enable: ptr(true), Level: ptr("warn")}}, cfg)
		require.NotNil(t, cfg.Logger)
		assert.Equal(t, amffactory.Logger{Enable: true, Level: "warn"}, *cfg.Logger)
		assert.Len(t, changes, 2)
	})

	t.Run("missing logger stays missing without global logger", func(t *testing.T) {
		cfg := newFullAMFConfig()
		cfg.Logger = nil
		applyAMF(t, &global.Global{Tac: ptr("000001")}, cfg)
		assert.Nil(t, cfg.Logger)
	})

	t.Run("missing sbi block skips sbiScheme", func(t *testing.T) {
		cfg := newFullAMFConfig()
		cfg.Configuration.Sbi = nil
		changes := applyAMF(t, &global.Global{SbiScheme: ptr("https")}, cfg)
		assert.Nil(t, cfg.Configuration.Sbi)
		assert.Empty(t, changes)
	})

	t.Run("unset timeouts are set", func(t *testing.T) {
		cfg := newFullAMFConfig()
		cfg.Configuration.ReadIdleTimeout = nil
		changes := applyAMF(t, &global.Global{Timeouts: &global.Timeouts{ReadIdle: ptr(time.Minute)}}, cfg)
		require.NotNil(t, cfg.Configuration.ReadIdleTimeout)
		assert.Equal(t, time.Minute, *cfg.Configuration.ReadIdleTimeout)
		assert.Equal(t, []adapter.Change{{NF: adapter.AMF, Path: "readIdleTimeout", Old: "<unset>", New: "1m0s"}}, changes)
	})
}

func TestAMFApplyEveryFieldNoOpWhenSameValues(t *testing.T) {
	cfg := newFullAMFConfig()
	changes := applyAMF(t, &global.Global{
		Plmn:      &models.PlmnId{Mcc: "208", Mnc: "93"},
		Tac:       nil, // TAI entries differ, so TAC cannot be a no-op here
		Snssais:   []models.Snssai{{Sst: 1, Sd: "010203"}, {Sst: 1, Sd: "112233"}},
		Dnns:      []string{"internet", "ims"},
		Nrf:       &global.NrfRef{Scheme: "http", Host: "127.0.0.10", Port: 8000},
		SbiScheme: ptr("http"),
		Logger:    &global.Logger{Enable: ptr(true), Level: ptr("info"), ReportCaller: ptr(false)},
		Timeouts:  &global.Timeouts{Idle: ptr(time.Minute), ReadIdle: ptr(time.Minute), Ping: ptr(15 * time.Second)},
	}, cfg)
	assert.Empty(t, changes)
	assert.Equal(t, newFullAMFConfig(), cfg)
}
