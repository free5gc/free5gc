package adapter_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ausffactory "github.com/free5gc/ausf/pkg/factory"
	bsffactory "github.com/free5gc/bsf/pkg/factory"
	chffactory "github.com/free5gc/chf/pkg/factory"
	"github.com/free5gc/config/cfg-generator/adapter"
	"github.com/free5gc/config/cfg-generator/global"
	n3iwffactory "github.com/free5gc/n3iwf/pkg/factory"
	neffactory "github.com/free5gc/nef/pkg/factory"
	nrffactory "github.com/free5gc/nrf/pkg/factory"
	pcffactory "github.com/free5gc/pcf/pkg/factory"
	scpfactory "github.com/free5gc/scp/pkg/factory"
	udmfactory "github.com/free5gc/udm/pkg/factory"
	udrfactory "github.com/free5gc/udr/pkg/factory"
)

// sbiCase describes one SBI NF adapter with the common shape. apply runs
// Apply on a fresh config whose logger block is present or absent; validate
// runs Validate on an empty config.
type sbiCase struct {
	nf       adapter.NFName
	hasMongo bool
	apply    func(g *global.Global, withLogger bool, r *adapter.Recorder) error
	validate func() error
}

func sbiCases() []sbiCase {
	idle := func() *time.Duration { d := time.Minute; return &d }
	return []sbiCase{
		{
			nf: adapter.AUSF,
			apply: func(g *global.Global, withLogger bool, r *adapter.Recorder) error {
				cfg := &ausffactory.Config{Configuration: &ausffactory.Configuration{
					Sbi: &ausffactory.Sbi{Scheme: "http"}, NrfUri: "http://nrf:8000", IdleTimeout: idle(),
				}}
				if withLogger {
					cfg.Logger = &ausffactory.Logger{Level: "info"}
				}
				return adapter.AUSFAdapter{}.Apply(g, cfg, r)
			},
			validate: func() error { return adapter.AUSFAdapter{}.Validate(&ausffactory.Config{}) },
		},
		{
			nf: adapter.NEF,
			apply: func(g *global.Global, withLogger bool, r *adapter.Recorder) error {
				cfg := &neffactory.Config{Configuration: &neffactory.Configuration{
					Sbi: &neffactory.Sbi{Scheme: "http"}, NrfUri: "http://nrf:8000", IdleTimeout: idle(),
				}}
				if withLogger {
					cfg.Logger = &neffactory.Logger{Level: "info"}
				}
				return adapter.NEFAdapter{}.Apply(g, cfg, r)
			},
			validate: func() error { return adapter.NEFAdapter{}.Validate(&neffactory.Config{}) },
		},
		{
			nf: adapter.SCP,
			apply: func(g *global.Global, withLogger bool, r *adapter.Recorder) error {
				cfg := &scpfactory.Config{Configuration: &scpfactory.Configuration{
					Sbi: &scpfactory.Sbi{Scheme: "http"}, NrfUri: "http://nrf:8000", IdleTimeout: idle(),
				}}
				if withLogger {
					cfg.Logger = &scpfactory.Logger{Level: "info"}
				}
				return adapter.SCPAdapter{}.Apply(g, cfg, r)
			},
			validate: func() error { return adapter.SCPAdapter{}.Validate(&scpfactory.Config{}) },
		},
		{
			nf: adapter.UDM,
			apply: func(g *global.Global, withLogger bool, r *adapter.Recorder) error {
				cfg := &udmfactory.Config{Configuration: &udmfactory.Configuration{
					Sbi: &udmfactory.Sbi{Scheme: "http"}, NrfUri: "http://nrf:8000", IdleTimeout: idle(),
				}}
				if withLogger {
					cfg.Logger = &udmfactory.Logger{Level: "info"}
				}
				return adapter.UDMAdapter{}.Apply(g, cfg, r)
			},
			validate: func() error { return adapter.UDMAdapter{}.Validate(&udmfactory.Config{}) },
		},
		{
			nf: adapter.CHF, hasMongo: true,
			apply: func(g *global.Global, withLogger bool, r *adapter.Recorder) error {
				cfg := &chffactory.Config{Configuration: &chffactory.Configuration{
					Sbi: &chffactory.Sbi{Scheme: "http"}, NrfUri: "http://nrf:8000", IdleTimeout: idle(),
					Mongodb: &chffactory.Mongodb{Name: "free5gc", Url: "mongodb://localhost:27017"},
				}}
				if withLogger {
					cfg.Logger = &chffactory.Logger{Level: "info"}
				}
				return adapter.CHFAdapter{}.Apply(g, cfg, r)
			},
			validate: func() error { return adapter.CHFAdapter{}.Validate(&chffactory.Config{}) },
		},
		{
			nf: adapter.UDR, hasMongo: true,
			apply: func(g *global.Global, withLogger bool, r *adapter.Recorder) error {
				cfg := &udrfactory.Config{Configuration: &udrfactory.Configuration{
					Sbi: &udrfactory.Sbi{Scheme: "http"}, NrfUri: "http://nrf:8000", IdleTimeout: idle(),
					Mongodb: &udrfactory.Mongodb{Name: "free5gc", Url: "mongodb://localhost:27017"},
				}}
				if withLogger {
					cfg.Logger = &udrfactory.Logger{Level: "info"}
				}
				return adapter.UDRAdapter{}.Apply(g, cfg, r)
			},
			validate: func() error { return adapter.UDRAdapter{}.Validate(&udrfactory.Config{}) },
		},
		{
			nf: adapter.PCF, hasMongo: true,
			apply: func(g *global.Global, withLogger bool, r *adapter.Recorder) error {
				cfg := &pcffactory.Config{Configuration: &pcffactory.Configuration{
					Sbi: &pcffactory.Sbi{Scheme: "http"}, NrfUri: "http://nrf:8000", IdleTimeout: idle(),
					Mongodb: &pcffactory.Mongodb{Name: "free5gc", Url: "mongodb://localhost:27017"},
				}}
				if withLogger {
					cfg.Logger = &pcffactory.Logger{Level: "info"}
				}
				return adapter.PCFAdapter{}.Apply(g, cfg, r)
			},
			validate: func() error { return adapter.PCFAdapter{}.Validate(&pcffactory.Config{}) },
		},
		{
			nf: adapter.BSF, hasMongo: true,
			apply: func(g *global.Global, withLogger bool, r *adapter.Recorder) error {
				cfg := &bsffactory.Config{Configuration: &bsffactory.Configuration{
					Sbi: &bsffactory.Sbi{Scheme: "http"}, NrfUri: "http://nrf:8000", IdleTimeout: idle(),
					MongoDB: &bsffactory.MongoDB{Name: "free5gc", Url: "mongodb://localhost:27017"},
				}}
				if withLogger {
					cfg.Logger = &bsffactory.Logger{Level: "info"}
				}
				return adapter.BSFAdapter{}.Apply(g, cfg, r)
			},
			validate: func() error { return adapter.BSFAdapter{}.Validate(&bsffactory.Config{}) },
		},
	}
}

func TestSBIAdaptersApplyCommonFields(t *testing.T) {
	g := &global.Global{
		Nrf:       &global.NrfRef{Scheme: "https", Host: "nrf2", Port: 29510},
		SbiScheme: ptr("https"),
		MongoDB:   &global.MongoDB{Name: ptr("db2"), Url: ptr("mongodb://db:27017")},
		Logger:    &global.Logger{Enable: ptr(true), Level: ptr("debug"), ReportCaller: ptr(true)},
		Timeouts:  &global.Timeouts{Idle: ptr(30 * time.Second), ReadIdle: ptr(time.Minute), Ping: ptr(15 * time.Second)},
	}
	for _, tc := range sbiCases() {
		t.Run(string(tc.nf), func(t *testing.T) {
			for _, withLogger := range []bool{true, false} {
				r := &adapter.Recorder{}
				require.NoError(t, tc.apply(g, withLogger, r))

				got := map[string]string{}
				for _, c := range r.Changes {
					assert.Equal(t, tc.nf, c.NF)
					got[c.Path] = c.New
				}
				want := map[string]string{
					"logger.enable":       "true",
					"logger.level":        "debug",
					"logger.reportCaller": "true",
					"idleTimeout":         "30s",
					"readIdleTimeout":     "1m0s",
					"pingTimeout":         "15s",
					"sbi.scheme":          "https",
					"nrfUri":              "https://nrf2:29510",
				}
				if tc.hasMongo {
					want["mongodb.name"], want["mongodb.url"] = "db2", "mongodb://db:27017"
				}
				assert.Equal(t, want, got, "withLogger=%v", withLogger)
			}
		})
	}
}

func TestSBIAdaptersValidateRejectEmptyConfig(t *testing.T) {
	for _, tc := range sbiCases() {
		assert.Error(t, tc.validate(), tc.nf)
	}
	assert.Error(t, adapter.NRFAdapter{}.Validate(&nrffactory.Config{}))
	assert.Error(t, adapter.N3IWFAdapter{}.Validate(&n3iwffactory.Config{}))
}
