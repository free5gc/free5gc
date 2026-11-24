package global_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/free5gc/openapi/models"

	"github.com/free5gc/config/cfg-generator/global"
)

func ptr[T any](v T) *T { return &v }

func TestValidate(t *testing.T) {
	cases := []struct {
		name      string
		g         *global.Global
		wantPaths []string // paths of expected errors, in order
	}{
		{name: "nil global is valid", g: nil},
		{name: "empty global is valid", g: &global.Global{}},

		// plmn
		{name: "plmn: valid 2-digit mnc", g: plmn("208", "93")},
		{name: "plmn: valid 3-digit mnc", g: plmn("310", "560")},
		{name: "plmn: valid leading-zero mnc", g: plmn("001", "01")},
		{name: "plmn: mcc too short", g: plmn("20", "93"), wantPaths: []string{"plmn.mcc"}},
		{name: "plmn: mcc too long", g: plmn("2080", "93"), wantPaths: []string{"plmn.mcc"}},
		{name: "plmn: mcc not numeric", g: plmn("2a8", "93"), wantPaths: []string{"plmn.mcc"}},
		{name: "plmn: mcc empty", g: plmn("", "93"), wantPaths: []string{"plmn.mcc"}},
		{name: "plmn: mnc too short", g: plmn("208", "9"), wantPaths: []string{"plmn.mnc"}},
		{name: "plmn: mnc too long", g: plmn("208", "9300"), wantPaths: []string{"plmn.mnc"}},
		{name: "plmn: mnc not numeric", g: plmn("208", "9x"), wantPaths: []string{"plmn.mnc"}},
		{name: "plmn: both invalid are all reported", g: plmn("x", "y"), wantPaths: []string{"plmn.mcc", "plmn.mnc"}},

		// tac
		{name: "tac: valid", g: &global.Global{Tac: ptr("000001")}},
		{name: "tac: valid upper hex", g: &global.Global{Tac: ptr("00AbCf")}},
		{name: "tac: too short", g: &global.Global{Tac: ptr("1")}, wantPaths: []string{"tac"}},
		{name: "tac: not hex", g: &global.Global{Tac: ptr("00000G")}, wantPaths: []string{"tac"}},
		{name: "tac: empty", g: &global.Global{Tac: ptr("")}, wantPaths: []string{"tac"}},

		// snssaiList
		{name: "snssai: valid", g: &global.Global{Snssais: []models.Snssai{{Sst: 1, Sd: "010203"}, {Sst: 1, Sd: "112233"}}}},
		{name: "snssai: sd optional", g: &global.Global{Snssais: []models.Snssai{{Sst: 1}}}},
		{name: "snssai: sst bounds ok", g: &global.Global{Snssais: []models.Snssai{{Sst: 0}, {Sst: 255}}}},
		{name: "snssai: empty list", g: &global.Global{Snssais: []models.Snssai{}}, wantPaths: []string{"snssaiList"}},
		{
			name:      "snssai: sst out of range",
			g:         &global.Global{Snssais: []models.Snssai{{Sst: 256}, {Sst: -1}}},
			wantPaths: []string{"snssaiList[0].sst", "snssaiList[1].sst"},
		},
		{
			name:      "snssai: bad sd",
			g:         &global.Global{Snssais: []models.Snssai{{Sst: 1, Sd: "12345"}}},
			wantPaths: []string{"snssaiList[0].sd"},
		},
		{
			name:      "snssai: duplicate",
			g:         &global.Global{Snssais: []models.Snssai{{Sst: 1, Sd: "010203"}, {Sst: 1, Sd: "010203"}}},
			wantPaths: []string{"snssaiList[1]"},
		},

		// dnnList
		{name: "dnn: valid", g: &global.Global{Dnns: []string{"internet", "ims"}}},
		{name: "dnn: empty list", g: &global.Global{Dnns: []string{}}, wantPaths: []string{"dnnList"}},
		{name: "dnn: empty name", g: &global.Global{Dnns: []string{"internet", ""}}, wantPaths: []string{"dnnList[1]"}},
		{name: "dnn: duplicate", g: &global.Global{Dnns: []string{"a", "b", "a"}}, wantPaths: []string{"dnnList[2]"}},

		// sbiScheme
		{name: "sbiScheme: http", g: &global.Global{SbiScheme: ptr("http")}},
		{name: "sbiScheme: https", g: &global.Global{SbiScheme: ptr("https")}},
		{name: "sbiScheme: invalid", g: &global.Global{SbiScheme: ptr("ftp")}, wantPaths: []string{"sbiScheme"}},
		{name: "sbiScheme: wrong case", g: &global.Global{SbiScheme: ptr("HTTP")}, wantPaths: []string{"sbiScheme"}},

		// nrf
		{name: "nrf: valid", g: &global.Global{Nrf: &global.NrfRef{Scheme: "http", Host: "127.0.0.10", Port: 8000}}},
		{
			name:      "nrf: all fields invalid",
			g:         &global.Global{Nrf: &global.NrfRef{Scheme: "tcp", Host: "", Port: 0}},
			wantPaths: []string{"nrf.scheme", "nrf.host", "nrf.port"},
		},
		{name: "nrf: port too high", g: &global.Global{Nrf: &global.NrfRef{Scheme: "http", Host: "h", Port: 65536}}, wantPaths: []string{"nrf.port"}},
		{
			name: "nrf and sbiScheme consistent",
			g:    &global.Global{SbiScheme: ptr("https"), Nrf: &global.NrfRef{Scheme: "https", Host: "h", Port: 1}},
		},
		{
			name:      "nrf and sbiScheme inconsistent",
			g:         &global.Global{SbiScheme: ptr("https"), Nrf: &global.NrfRef{Scheme: "http", Host: "h", Port: 1}},
			wantPaths: []string{"nrf.scheme"},
		},
		{
			name:      "invalid scheme reported once, not also as inconsistent",
			g:         &global.Global{SbiScheme: ptr("ftp"), Nrf: &global.NrfRef{Scheme: "http", Host: "h", Port: 1}},
			wantPaths: []string{"sbiScheme"},
		},

		// mongodb
		{name: "mongodb: valid", g: &global.Global{MongoDB: &global.MongoDB{Name: ptr("free5gc"), Url: ptr("mongodb://127.0.0.1:27017")}}},
		{name: "mongodb: srv url", g: &global.Global{MongoDB: &global.MongoDB{Url: ptr("mongodb+srv://cluster.example")}}},
		{name: "mongodb: only name", g: &global.Global{MongoDB: &global.MongoDB{Name: ptr("db2")}}},
		{
			name:      "mongodb: empty name and bad url",
			g:         &global.Global{MongoDB: &global.MongoDB{Name: ptr(""), Url: ptr("http://db:27017")}},
			wantPaths: []string{"mongodb.name", "mongodb.url"},
		},

		// logger
		{name: "logger: valid level", g: &global.Global{Logger: &global.Logger{Level: ptr("debug")}}},
		{name: "logger: only enable", g: &global.Global{Logger: &global.Logger{Enable: ptr(false)}}},
		{name: "logger: invalid level", g: &global.Global{Logger: &global.Logger{Level: ptr("verbose")}}, wantPaths: []string{"logger.level"}},

		// timeouts
		{name: "timeouts: zero allowed", g: &global.Global{Timeouts: &global.Timeouts{Idle: ptr(time.Duration(0))}}},
		{
			name: "timeouts: negative",
			g: &global.Global{Timeouts: &global.Timeouts{
				Idle: ptr(-time.Second), ReadIdle: ptr(time.Second), Ping: ptr(-time.Millisecond),
			}},
			wantPaths: []string{"timeouts.idle", "timeouts.ping"},
		},

		// several fields at once: every error is collected
		{
			name:      "errors from different fields are all reported",
			g:         &global.Global{Plmn: &models.PlmnId{Mcc: "1", Mnc: "01"}, Tac: ptr("x"), Dnns: []string{}},
			wantPaths: []string{"plmn.mcc", "tac", "dnnList"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			errs := global.Validate(tc.g)
			require.Len(t, errs, len(tc.wantPaths), "errors: %v", errs)
			for i, err := range errs {
				var fe *global.FieldError
				require.ErrorAs(t, err, &fe)
				assert.Equal(t, tc.wantPaths[i], fe.Path)
				assert.Contains(t, err.Error(), "global."+tc.wantPaths[i])
			}
		})
	}
}

func TestNrfRefURI(t *testing.T) {
	assert.Equal(t, "http://127.0.0.10:8000", global.NrfRef{Scheme: "http", Host: "127.0.0.10", Port: 8000}.URI())
	assert.Equal(t, "https://nrf.example:443", global.NrfRef{Scheme: "https", Host: "nrf.example", Port: 443}.URI())
	assert.Equal(t, "http://[fd00::10]:8000", global.NrfRef{Scheme: "http", Host: "fd00::10", Port: 8000}.URI(),
		"IPv6 host must be bracketed")
}

// TestDecode pins the YAML semantics the validator relies on.
func TestDecode(t *testing.T) {
	t.Run("full block", func(t *testing.T) {
		src := `
plmn: {mcc: 001, mnc: 01}
tac: 000001
snssaiList:
  - {sst: 1, sd: 010203}
  - {sst: 2}
dnnList: [internet, ims]
nrf: {scheme: http, host: 127.0.0.10, port: 8000}
sbiScheme: http
logger: {enable: false, level: debug}
timeouts: {idle: 60s, readIdle: 1m, ping: 15s}
`
		var g global.Global
		require.NoError(t, yaml.UnmarshalStrict([]byte(src), &g))

		assert.Equal(t, models.PlmnId{Mcc: "001", Mnc: "01"}, *g.Plmn, "leading zeros kept")
		assert.Equal(t, "000001", *g.Tac)
		assert.Equal(t, []models.Snssai{{Sst: 1, Sd: "010203"}, {Sst: 2}}, g.Snssais)
		assert.Equal(t, []string{"internet", "ims"}, g.Dnns)
		assert.Equal(t, global.NrfRef{Scheme: "http", Host: "127.0.0.10", Port: 8000}, *g.Nrf)
		assert.Equal(t, "http", *g.SbiScheme)
		assert.Equal(t, false, *g.Logger.Enable)
		assert.Equal(t, "debug", *g.Logger.Level)
		assert.Nil(t, g.Logger.ReportCaller, "not provided stays nil")
		assert.Equal(t, 60*time.Second, *g.Timeouts.Idle)
		assert.Equal(t, time.Minute, *g.Timeouts.ReadIdle)
		assert.Equal(t, 15*time.Second, *g.Timeouts.Ping)
		assert.Empty(t, global.Validate(&g))
	})

	t.Run("omitted list is nil, explicit [] is empty", func(t *testing.T) {
		var omitted, empty global.Global
		require.NoError(t, yaml.UnmarshalStrict([]byte("tac: \"000001\"\n"), &omitted))
		require.NoError(t, yaml.UnmarshalStrict([]byte("snssaiList: []\ndnnList: []\n"), &empty))

		assert.Nil(t, omitted.Snssais)
		assert.Nil(t, omitted.Dnns)
		assert.NotNil(t, empty.Snssais)
		assert.NotNil(t, empty.Dnns)
		assert.Len(t, global.Validate(&empty), 2)
	})
}

func plmn(mcc, mnc string) *global.Global {
	return &global.Global{Plmn: &models.PlmnId{Mcc: mcc, Mnc: mnc}}
}
