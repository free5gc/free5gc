package adapter_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/free5gc/openapi/models"

	ausffactory "github.com/free5gc/ausf/pkg/factory"
	bsffactory "github.com/free5gc/bsf/pkg/factory"
	chffactory "github.com/free5gc/chf/pkg/factory"
	"github.com/free5gc/config/cfg-generator/adapter"
	"github.com/free5gc/config/cfg-generator/global"
	upffactory "github.com/free5gc/go-upf/pkg/factory"
	n3iwffactory "github.com/free5gc/n3iwf/pkg/factory"
	neffactory "github.com/free5gc/nef/pkg/factory"
	nrffactory "github.com/free5gc/nrf/pkg/factory"
	nssffactory "github.com/free5gc/nssf/pkg/factory"
	pcffactory "github.com/free5gc/pcf/pkg/factory"
	scpfactory "github.com/free5gc/scp/pkg/factory"
	tngfcontext "github.com/free5gc/tngf/pkg/context"
	tngffactory "github.com/free5gc/tngf/pkg/factory"
	udmfactory "github.com/free5gc/udm/pkg/factory"
	udrfactory "github.com/free5gc/udr/pkg/factory"
)

func TestAdapterMetadata(t *testing.T) {
	cases := []struct {
		nf, file string
		gotNF    adapter.NFName
		gotFile  string
	}{
		{"ausf", "ausfcfg.yaml", adapter.AUSFAdapter{}.NF(), adapter.AUSFAdapter{}.FileName()},
		{"bsf", "bsfcfg.yaml", adapter.BSFAdapter{}.NF(), adapter.BSFAdapter{}.FileName()},
		{"chf", "chfcfg.yaml", adapter.CHFAdapter{}.NF(), adapter.CHFAdapter{}.FileName()},
		{"n3iwf", "n3iwfcfg.yaml", adapter.N3IWFAdapter{}.NF(), adapter.N3IWFAdapter{}.FileName()},
		{"nef", "nefcfg.yaml", adapter.NEFAdapter{}.NF(), adapter.NEFAdapter{}.FileName()},
		{"nrf", "nrfcfg.yaml", adapter.NRFAdapter{}.NF(), adapter.NRFAdapter{}.FileName()},
		{"nssf", "nssfcfg.yaml", adapter.NSSFAdapter{}.NF(), adapter.NSSFAdapter{}.FileName()},
		{"pcf", "pcfcfg.yaml", adapter.PCFAdapter{}.NF(), adapter.PCFAdapter{}.FileName()},
		{"scp", "scpcfg.yaml", adapter.SCPAdapter{}.NF(), adapter.SCPAdapter{}.FileName()},
		{"tngf", "tngfcfg.yaml", adapter.TNGFAdapter{}.NF(), adapter.TNGFAdapter{}.FileName()},
		{"udm", "udmcfg.yaml", adapter.UDMAdapter{}.NF(), adapter.UDMAdapter{}.FileName()},
		{"udr", "udrcfg.yaml", adapter.UDRAdapter{}.NF(), adapter.UDRAdapter{}.FileName()},
		{"upf", "upfcfg.yaml", adapter.UPFAdapter{}.NF(), adapter.UPFAdapter{}.FileName()},
	}
	for _, tc := range cases {
		assert.Equal(t, adapter.NFName(tc.nf), tc.gotNF)
		assert.Equal(t, tc.file, tc.gotFile)
	}
}

// TestAdaptersNilGlobalAndMissingConfiguration covers the two guard clauses
// every wrapped NF adapter shares.
func TestAdaptersNilGlobalAndMissingConfiguration(t *testing.T) {
	g := &global.Global{Plmn: &models.PlmnId{Mcc: "001", Mnc: "01"}}
	cases := map[string]func(*global.Global, *adapter.Recorder) error{
		"ausf": func(g *global.Global, r *adapter.Recorder) error {
			return adapter.AUSFAdapter{}.Apply(g, &ausffactory.Config{}, r)
		},
		"bsf": func(g *global.Global, r *adapter.Recorder) error {
			return adapter.BSFAdapter{}.Apply(g, &bsffactory.Config{}, r)
		},
		"chf": func(g *global.Global, r *adapter.Recorder) error {
			return adapter.CHFAdapter{}.Apply(g, &chffactory.Config{}, r)
		},
		"n3iwf": func(g *global.Global, r *adapter.Recorder) error {
			return adapter.N3IWFAdapter{}.Apply(g, &n3iwffactory.Config{}, r)
		},
		"nef": func(g *global.Global, r *adapter.Recorder) error {
			return adapter.NEFAdapter{}.Apply(g, &neffactory.Config{}, r)
		},
		"nrf": func(g *global.Global, r *adapter.Recorder) error {
			return adapter.NRFAdapter{}.Apply(g, &nrffactory.Config{}, r)
		},
		"nssf": func(g *global.Global, r *adapter.Recorder) error {
			return adapter.NSSFAdapter{}.Apply(g, &nssffactory.Config{}, r)
		},
		"pcf": func(g *global.Global, r *adapter.Recorder) error {
			return adapter.PCFAdapter{}.Apply(g, &pcffactory.Config{}, r)
		},
		"scp": func(g *global.Global, r *adapter.Recorder) error {
			return adapter.SCPAdapter{}.Apply(g, &scpfactory.Config{}, r)
		},
		"tngf": func(g *global.Global, r *adapter.Recorder) error {
			return adapter.TNGFAdapter{}.Apply(g, &tngffactory.Config{}, r)
		},
		"udm": func(g *global.Global, r *adapter.Recorder) error {
			return adapter.UDMAdapter{}.Apply(g, &udmfactory.Config{}, r)
		},
		"udr": func(g *global.Global, r *adapter.Recorder) error {
			return adapter.UDRAdapter{}.Apply(g, &udrfactory.Config{}, r)
		},
	}
	for nf, apply := range cases {
		t.Run(nf, func(t *testing.T) {
			r := &adapter.Recorder{}
			assert.NoError(t, apply(nil, r), "nil global is a no-op")
			err := apply(g, r)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "configuration block is missing")
			assert.Empty(t, r.Changes)
		})
	}
	assert.NoError(t, adapter.UPFAdapter{}.Apply(nil, &upffactory.Config{}, &adapter.Recorder{}))
}

func TestAUSFOnlyFirstPlmn(t *testing.T) {
	cfg := &ausffactory.Config{Configuration: &ausffactory.Configuration{
		PlmnSupportList: []models.PlmnId{{Mcc: "208", Mnc: "93"}, {Mcc: "123", Mnc: "45"}},
	}}
	r := &adapter.Recorder{}
	require.NoError(t, adapter.AUSFAdapter{}.Apply(&global.Global{Plmn: &models.PlmnId{Mcc: "001", Mnc: "01"}}, cfg, r))

	assert.Equal(t, []models.PlmnId{{Mcc: "001", Mnc: "01"}, {Mcc: "123", Mnc: "45"}}, cfg.Configuration.PlmnSupportList)
	assert.Len(t, r.Changes, 2)

	empty := &ausffactory.Config{Configuration: &ausffactory.Configuration{}}
	r = &adapter.Recorder{}
	require.NoError(t, adapter.AUSFAdapter{}.Apply(&global.Global{Plmn: &models.PlmnId{Mcc: "001", Mnc: "01"}}, empty, r))
	assert.Empty(t, r.Changes, "no PLMN entry to change; AUSF Validate reports it")
}

func newNRFConfig() *nrffactory.Config {
	return &nrffactory.Config{
		Configuration: &nrffactory.Configuration{
			Sbi:           &nrffactory.Sbi{Scheme: "http", RegisterIPv4: "127.0.0.10", Port: 8000},
			MongoDBName:   "free5gc",
			MongoDBUrl:    "mongodb://127.0.0.1:27017",
			DefaultPlmnId: models.PlmnId{Mcc: "208", Mnc: "93"},
		},
		Logger: &nrffactory.Logger{Level: "info"},
	}
}

func TestNRFConsistencyCheck(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func(*nrffactory.Config)
		g       *global.Global
		wantErr string
	}{
		{name: "matches", g: &global.Global{Nrf: &global.NrfRef{Scheme: "http", Host: "127.0.0.10", Port: 8000}}},
		{
			name:   "port 0 means the NRF default 8000",
			mutate: func(c *nrffactory.Config) { c.Configuration.Sbi.Port = 0 },
			g:      &global.Global{Nrf: &global.NrfRef{Scheme: "http", Host: "127.0.0.10", Port: 8000}},
		},
		{
			name: "sbiScheme switches the effective NRF scheme",
			g: &global.Global{
				SbiScheme: ptr("https"),
				Nrf:       &global.NrfRef{Scheme: "https", Host: "127.0.0.10", Port: 8000},
			},
		},
		{
			name:    "host differs",
			g:       &global.Global{Nrf: &global.NrfRef{Scheme: "http", Host: "10.0.0.1", Port: 8000}},
			wantErr: "global nrf http://10.0.0.1:8000 does not match the NRF's own SBI http://127.0.0.10:8000",
		},
		{
			name:    "port differs",
			g:       &global.Global{Nrf: &global.NrfRef{Scheme: "http", Host: "127.0.0.10", Port: 9000}},
			wantErr: "does not match",
		},
		{
			name:    "scheme differs",
			g:       &global.Global{Nrf: &global.NrfRef{Scheme: "https", Host: "127.0.0.10", Port: 8000}},
			wantErr: "https://127.0.0.10:8000 does not match the NRF's own SBI http://127.0.0.10:8000",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := newNRFConfig()
			if tc.mutate != nil {
				tc.mutate(cfg)
			}
			g := tc.g
			g.Plmn = &models.PlmnId{Mcc: "001", Mnc: "01"} // would change if applied
			r := &adapter.Recorder{}
			err := adapter.NRFAdapter{}.Apply(g, cfg, r)
			if tc.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
			assert.Empty(t, r.Changes, "checked before any change")
		})
	}
}

func TestNRFApplyFields(t *testing.T) {
	cfg := newNRFConfig()
	r := &adapter.Recorder{}
	require.NoError(t, adapter.NRFAdapter{}.Apply(&global.Global{
		Plmn:     &models.PlmnId{Mcc: "001", Mnc: "01"},
		MongoDB:  &global.MongoDB{Name: ptr("db2"), Url: ptr("mongodb://db:27017")},
		Logger:   &global.Logger{Level: ptr("warn")},
		Timeouts: &global.Timeouts{Ping: ptr(5 * time.Second)},
	}, cfg, r))

	assert.Equal(t, []adapter.Change{
		{NF: adapter.NRF, Path: "DefaultPlmnId.mcc", Old: "208", New: "001"},
		{NF: adapter.NRF, Path: "DefaultPlmnId.mnc", Old: "93", New: "01"},
		{NF: adapter.NRF, Path: "logger.level", Old: "info", New: "warn"},
		{NF: adapter.NRF, Path: "pingTimeout", Old: "<unset>", New: "5s"},
		{NF: adapter.NRF, Path: "MongoDBName", Old: "free5gc", New: "db2"},
		{NF: adapter.NRF, Path: "MongoDBUrl", Old: "mongodb://127.0.0.1:27017", New: "mongodb://db:27017"},
	}, r.Changes)
}

func TestNSSFNrfPrefixAndPlmn(t *testing.T) {
	cfg := &nssffactory.Config{Configuration: &nssffactory.Configuration{
		Sbi:               &nssffactory.Sbi{Scheme: models.UriScheme_HTTP},
		NrfUri:            "http://127.0.0.10:8000",
		SupportedPlmnList: []models.PlmnId{{Mcc: "208", Mnc: "93"}},
		SupportedNssaiInPlmnList: []nssffactory.SupportedNssaiInPlmn{
			{PlmnId: &models.PlmnId{Mcc: "208", Mnc: "93"}},
			{PlmnId: nil},
		},
		NsiList: []nssffactory.NsiConfig{{NsiInformationList: []models.Nssf_NSSel_NsiInformation{
			{NrfId: "http://127.0.0.10:8000/nnrf-nfm/v1/nf-instances"},
			{NrfId: "http://127.0.0.10:80001/other"},
		}}},
		AmfSetList: []nssffactory.AmfSetConfig{
			{NrfAmfSet: "http://127.0.0.10:8000/nnrf-nfm/v1/nf-instances"},
			{NrfAmfSet: "http://localhost:8084/nnrf-nfm/v1/nf-instances"},
			{NrfAmfSet: ""},
		},
	}}
	r := &adapter.Recorder{}
	require.NoError(t, adapter.NSSFAdapter{}.Apply(&global.Global{
		Plmn:      &models.PlmnId{Mcc: "001", Mnc: "01"},
		Nrf:       &global.NrfRef{Scheme: "https", Host: "nrf", Port: 443},
		SbiScheme: ptr("https"),
	}, cfg, r))

	c := cfg.Configuration
	assert.Equal(t, "https://nrf:443", c.NrfUri)
	assert.Equal(t, models.UriScheme_HTTPS, c.Sbi.Scheme)
	assert.Equal(t, models.PlmnId{Mcc: "001", Mnc: "01"}, c.SupportedPlmnList[0])
	assert.Equal(t, models.PlmnId{Mcc: "001", Mnc: "01"}, *c.SupportedNssaiInPlmnList[0].PlmnId)
	assert.Nil(t, c.SupportedNssaiInPlmnList[1].PlmnId, "missing plmnId is not invented")
	assert.Equal(t, "https://nrf:443/nnrf-nfm/v1/nf-instances", c.NsiList[0].NsiInformationList[0].NrfId)
	assert.Equal(t, "http://127.0.0.10:80001/other", c.NsiList[0].NsiInformationList[1].NrfId,
		"a longer port is not a URI prefix match")
	assert.Equal(t, "https://nrf:443/nnrf-nfm/v1/nf-instances", c.AmfSetList[0].NrfAmfSet)
	assert.Equal(t, "http://localhost:8084/nnrf-nfm/v1/nf-instances", c.AmfSetList[1].NrfAmfSet, "other NRF kept")
	assert.Equal(t, "", c.AmfSetList[2].NrfAmfSet)
}

func TestNSSFNrfPrefixSkippedWhenNrfUriEmpty(t *testing.T) {
	cfg := &nssffactory.Config{Configuration: &nssffactory.Configuration{
		AmfSetList: []nssffactory.AmfSetConfig{{NrfAmfSet: "http://x/nnrf"}},
	}}
	require.NoError(t, adapter.NSSFAdapter{}.Apply(&global.Global{
		Nrf: &global.NrfRef{Scheme: "http", Host: "nrf", Port: 1},
	}, cfg, &adapter.Recorder{}))
	assert.Equal(t, "http://x/nnrf", cfg.Configuration.AmfSetList[0].NrfAmfSet)
	assert.Equal(t, "http://nrf:1", cfg.Configuration.NrfUri)
}

func TestNSSFValidateAddsMissingChecks(t *testing.T) {
	err := adapter.NSSFAdapter{}.Validate(&nssffactory.Config{
		Info:          &nssffactory.Info{Version: "1.0.2"},
		Configuration: &nssffactory.Configuration{NrfUri: "not a url"},
		Logger:        &nssffactory.Logger{Level: "info"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "configuration.sbi is required")
	assert.Contains(t, err.Error(), `configuration.nrfUri "not a url" must be an absolute URL`)
	assert.Contains(t, err.Error(), "configuration.serviceNameList is required")

	err = adapter.NSSFAdapter{}.Validate(&nssffactory.Config{
		Configuration: &nssffactory.Configuration{Sbi: &nssffactory.Sbi{Scheme: "ftp"}},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), `configuration.sbi.scheme "ftp" must be http or https`)
}

func TestPCFNotAppliedForUnadvertisedDnn(t *testing.T) {
	notApplied := func(dnn string) adapter.NotApplied {
		return adapter.NotApplied{
			NF: adapter.PCF, Item: "NRF registration DNN list",
			HardCoded: "[free5gc internet]", Requested: dnn,
		}
	}
	cases := []struct {
		name string
		dnns []string
		want []adapter.NotApplied
	}{
		{"no global DNN", nil, nil},
		{"covered: internet", []string{"internet"}, nil},
		{"covered: free5gc and internet", []string{"free5gc", "internet"}, nil},
		{"two outside the list", []string{"internet", "ims", "iot"}, []adapter.NotApplied{notApplied("ims"), notApplied("iot")}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := &adapter.Recorder{}
			cfg := &pcffactory.Config{Configuration: &pcffactory.Configuration{}}
			require.NoError(t, adapter.PCFAdapter{}.Apply(&global.Global{Dnns: tc.dnns}, cfg, r))
			assert.Equal(t, tc.want, r.NotApplied)
			assert.Empty(t, r.Warnings, "not-applied values are not generic warnings")
			assert.Empty(t, r.Changes)
		})
	}
	assert.Equal(t,
		"pcf: NRF registration DNN list stays [free5gc internet] (hard-coded); global ims not applied",
		notApplied("ims").String())
}

func TestMongoBlockCreatedOnlyWhenGiven(t *testing.T) {
	g := &global.Global{MongoDB: &global.MongoDB{Url: ptr("mongodb://db:27017")}}

	pcf := &pcffactory.Config{Configuration: &pcffactory.Configuration{}}
	r := &adapter.Recorder{}
	require.NoError(t, adapter.PCFAdapter{}.Apply(g, pcf, r))
	require.NotNil(t, pcf.Configuration.Mongodb)
	assert.Equal(t, "mongodb://db:27017", pcf.Configuration.Mongodb.Url)
	assert.Equal(t, []adapter.Change{{NF: adapter.PCF, Path: "mongodb.url", Old: "", New: "mongodb://db:27017"}}, r.Changes)

	udr := &udrfactory.Config{Configuration: &udrfactory.Configuration{}}
	require.NoError(t, adapter.UDRAdapter{}.Apply(&global.Global{Tac: ptr("000001")}, udr, &adapter.Recorder{}))
	assert.Nil(t, udr.Configuration.Mongodb, "not created without global mongodb")

	chf := &chffactory.Config{Configuration: &chffactory.Configuration{}}
	require.NoError(t, adapter.CHFAdapter{}.Apply(g, chf, &adapter.Recorder{}))
	assert.Equal(t, "mongodb://db:27017", chf.Configuration.Mongodb.Url)

	bsf := &bsffactory.Config{Configuration: &bsffactory.Configuration{}}
	require.NoError(t, adapter.BSFAdapter{}.Apply(g, bsf, &adapter.Recorder{}))
	assert.Equal(t, "mongodb://db:27017", bsf.Configuration.MongoDB.Url)
}

func TestBSFNotAppliedForDifferentPlmn(t *testing.T) {
	cases := []struct {
		name string
		plmn *models.PlmnId
		want []adapter.NotApplied
	}{
		{"no global PLMN", nil, nil},
		{"same as the hard-coded PLMN", &models.PlmnId{Mcc: "208", Mnc: "93"}, nil},
		{"different MCC and MNC", &models.PlmnId{Mcc: "466", Mnc: "01"}, []adapter.NotApplied{{
			NF: adapter.BSF, Item: "NRF registration PLMN", HardCoded: "208/93", Requested: "466/01",
		}}},
		{"different MNC only", &models.PlmnId{Mcc: "208", Mnc: "95"}, []adapter.NotApplied{{
			NF: adapter.BSF, Item: "NRF registration PLMN", HardCoded: "208/93", Requested: "208/95",
		}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := &adapter.Recorder{}
			cfg := &bsffactory.Config{Configuration: &bsffactory.Configuration{}}
			require.NoError(t, adapter.BSFAdapter{}.Apply(&global.Global{Plmn: tc.plmn}, cfg, r))
			assert.Equal(t, tc.want, r.NotApplied)
			assert.Empty(t, r.Warnings, "not-applied values are not generic warnings")
		})
	}
	assert.Equal(t,
		"bsf: NRF registration PLMN stays 208/93 (hard-coded); global 466/01 not applied",
		adapter.NotApplied{NF: adapter.BSF, Item: "NRF registration PLMN", HardCoded: "208/93", Requested: "466/01"}.String())
}

func TestRecorderWarnIsSeparateFromNotApplied(t *testing.T) {
	r := &adapter.Recorder{}
	r.Warn(adapter.AMF, "something %s", "odd")
	r.NotApply(adapter.BSF, "item", "a", "b")
	assert.Equal(t, []adapter.Warning{{NF: adapter.AMF, Msg: "something odd"}}, r.Warnings)
	assert.Equal(t, "amf: something odd", r.Warnings[0].String())
	assert.Equal(t, []adapter.NotApplied{{NF: adapter.BSF, Item: "item", HardCoded: "a", Requested: "b"}}, r.NotApplied)
}

func validBSF() *bsffactory.Config {
	return &bsffactory.Config{
		Info: &bsffactory.Info{Version: "1.0.0"},
		Configuration: &bsffactory.Configuration{
			Sbi:     &bsffactory.Sbi{Scheme: "http", RegisterIPv4: "127.0.0.15", BindingIPv4: "127.0.0.15", Port: 8000},
			NrfUri:  "http://127.0.0.10:8000",
			MongoDB: &bsffactory.MongoDB{Name: "free5gc", Url: "mongodb://127.0.0.1:27017"},
		},
		Logger: &bsffactory.Logger{Level: "info"},
	}
}

func TestBSFValidate(t *testing.T) {
	require.NoError(t, adapter.BSFAdapter{}.Validate(validBSF()))

	cases := []struct {
		name    string
		mutate  func(*bsffactory.Config)
		wantErr string
	}{
		{"wrong version", func(c *bsffactory.Config) { c.Info.Version = "1.0.9" }, `info.version must be "1.0.0"`},
		{"no info", func(c *bsffactory.Config) { c.Info = nil }, "info.version"},
		{"no logger", func(c *bsffactory.Config) { c.Logger = nil }, "logger is required"},
		{"bad level", func(c *bsffactory.Config) { c.Logger.Level = "loud" }, `logger.level "loud" is invalid`},
		{"no configuration", func(c *bsffactory.Config) { c.Configuration = nil }, "configuration is required"},
		{"no sbi", func(c *bsffactory.Config) { c.Configuration.Sbi = nil }, "configuration.sbi is required"},
		{"bad scheme", func(c *bsffactory.Config) { c.Configuration.Sbi.Scheme = "ftp" }, "must be http or https"},
		{"no register ip", func(c *bsffactory.Config) { c.Configuration.Sbi.RegisterIPv4 = "" }, "registerIPv4 or registerIPv6"},
		{"no binding ip", func(c *bsffactory.Config) { c.Configuration.Sbi.BindingIPv4 = "" }, "bindingIPv4 or bindingIPv6"},
		{"bad port", func(c *bsffactory.Config) { c.Configuration.Sbi.Port = 70000 }, "port 70000"},
		{"bad nrfUri", func(c *bsffactory.Config) { c.Configuration.NrfUri = "127.0.0.10" }, "nrfUri"},
		{"no mongodb", func(c *bsffactory.Config) { c.Configuration.MongoDB = nil }, "mongodb.name and url are required"},
		{"bad mongodb url", func(c *bsffactory.Config) { c.Configuration.MongoDB.Url = "http://db" }, "must start with mongodb://"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := validBSF()
			tc.mutate(cfg)
			err := adapter.BSFAdapter{}.Validate(cfg)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

func newUPFConfig() *upffactory.Config {
	return &upffactory.Config{
		Version: "1.0.3",
		Pfcp:    &upffactory.Pfcp{Addr: "127.0.0.8", NodeID: "127.0.0.8", RetransTimeout: time.Second},
		Gtpu:    &upffactory.Gtpu{Forwarder: "gtp5g"},
		DnnList: []upffactory.DnnList{
			{Dnn: "internet", Cidr: "10.60.0.0/16"},
			{Dnn: "internet", Cidr: "10.61.0.0/16"},
		},
		Logger: &upffactory.Logger{Level: "info"},
	}
}

func TestUPFApply(t *testing.T) {
	cfg := newUPFConfig()
	r := &adapter.Recorder{}
	require.NoError(t, adapter.UPFAdapter{}.Apply(&global.Global{
		Dnns:     []string{"ims"},
		Logger:   &global.Logger{Level: ptr("debug")},
		Timeouts: &global.Timeouts{Idle: ptr(time.Second)}, // UPF has none: ignored
		Plmn:     &models.PlmnId{Mcc: "001", Mnc: "01"},    // UPF has none: ignored
	}, cfg, r))
	assert.Equal(t, []adapter.Change{
		{NF: adapter.UPF, Path: "dnnList[0].dnn", Old: "internet", New: "ims"},
		{NF: adapter.UPF, Path: "dnnList[1].dnn", Old: "internet", New: "ims"},
		{NF: adapter.UPF, Path: "logger.level", Old: "info", New: "debug"},
	}, r.Changes)

	r = &adapter.Recorder{}
	err := adapter.UPFAdapter{}.Apply(&global.Global{Dnns: []string{"a", "b"}}, newUPFConfig(), r)
	require.Error(t, err)
	assert.Empty(t, r.Changes)

	nolog := newUPFConfig()
	nolog.Logger = nil
	require.NoError(t, adapter.UPFAdapter{}.Apply(&global.Global{Logger: &global.Logger{Level: ptr("warn")}}, nolog, &adapter.Recorder{}))
	assert.Equal(t, "warn", nolog.Logger.Level)
}

func TestUPFValidate(t *testing.T) {
	require.NoError(t, adapter.UPFAdapter{}.Validate(newUPFConfig()))

	bad := newUPFConfig()
	bad.DnnList[0].Cidr = "10.60.0.0"
	assert.Error(t, adapter.UPFAdapter{}.Validate(bad), "cidr tag must be enforced")

	badVersion := newUPFConfig()
	badVersion.Version = "9.9.9"
	assert.Error(t, adapter.UPFAdapter{}.Validate(badVersion))
}

func newN3IWFConfig() *n3iwffactory.Config {
	return &n3iwffactory.Config{Configuration: &n3iwffactory.Configuration{
		N3IWFInfo: &n3iwffactory.N3IWFNFInfo{
			GlobalN3IWFID: &n3iwffactory.GlobalN3IWFID{PLMNID: &n3iwffactory.PLMNID{Mcc: "208", Mnc: "93"}, N3IWFID: 135},
			SupportedTAList: []n3iwffactory.SupportedTAItem{{
				TAC: "000001",
				BroadcastPLMNList: []n3iwffactory.BroadcastPLMNItem{{
					PLMNID: &n3iwffactory.PLMNID{Mcc: "208", Mnc: "93"},
					TAISliceSupportList: []n3iwffactory.SliceSupportItem{
						{SNSSAI: n3iwffactory.SNSSAIItem{SST: 1, SD: "010203"}},
						{SNSSAI: n3iwffactory.SNSSAIItem{SST: 1, SD: "112233"}},
					},
				}},
			}},
		},
	}}
}

func TestN3IWFApply(t *testing.T) {
	cfg := newN3IWFConfig()
	cfg.Configuration.N3IWFInfo.SupportedTAList[0].BroadcastPLMNList[0].PLMNID = nil
	r := &adapter.Recorder{}
	require.NoError(t, adapter.N3IWFAdapter{}.Apply(&global.Global{
		Plmn:    &models.PlmnId{Mcc: "001", Mnc: "01"},
		Tac:     ptr("00ABCD"),
		Snssais: []models.Snssai{{Sst: 2, Sd: "000001"}, {Sst: 1, Sd: "112233"}},
	}, cfg, r))

	info := cfg.Configuration.N3IWFInfo
	assert.Equal(t, n3iwffactory.PLMNID{Mcc: "001", Mnc: "01"}, *info.GlobalN3IWFID.PLMNID)
	assert.Equal(t, uint16(135), info.GlobalN3IWFID.N3IWFID, "n3iwfID kept")
	ta := info.SupportedTAList[0]
	assert.Equal(t, "00ABCD", ta.TAC)
	assert.Equal(t, n3iwffactory.PLMNID{Mcc: "001", Mnc: "01"}, *ta.BroadcastPLMNList[0].PLMNID, "missing plmnID created")
	assert.Equal(t, n3iwffactory.SNSSAIItem{SST: 2, SD: "000001"}, ta.BroadcastPLMNList[0].TAISliceSupportList[0].SNSSAI)
	assert.Equal(t, n3iwffactory.SNSSAIItem{SST: 1, SD: "112233"}, ta.BroadcastPLMNList[0].TAISliceSupportList[1].SNSSAI)
	assert.Contains(t, r.Changes, adapter.Change{
		NF:   adapter.N3IWF,
		Path: "n3iwfInformation.supportedTAList[0].broadcastPlmnList[0].taiSliceSupportList[0].snssai.sst",
		Old:  "1", New: "2",
	})
}

func TestN3IWFApplyEdgeCases(t *testing.T) {
	t.Run("snssai count mismatch changes nothing", func(t *testing.T) {
		cfg := newN3IWFConfig()
		r := &adapter.Recorder{}
		err := adapter.N3IWFAdapter{}.Apply(&global.Global{
			Plmn:    &models.PlmnId{Mcc: "001", Mnc: "01"},
			Snssais: []models.Snssai{{Sst: 1}},
		}, cfg, r)
		require.Error(t, err)
		assert.Empty(t, r.Changes)
		assert.Equal(t, newN3IWFConfig(), cfg)
	})

	t.Run("missing n3iwfInformation is skipped", func(t *testing.T) {
		cfg := &n3iwffactory.Config{Configuration: &n3iwffactory.Configuration{}}
		r := &adapter.Recorder{}
		require.NoError(t, adapter.N3IWFAdapter{}.Apply(&global.Global{
			Plmn:     &models.PlmnId{Mcc: "001", Mnc: "01"},
			Timeouts: &global.Timeouts{Idle: ptr(time.Minute), Ping: ptr(time.Second)},
		}, cfg, r))
		assert.Equal(t, []adapter.Change{{NF: adapter.N3IWF, Path: "idleTimeout", Old: "<unset>", New: "1m0s"}}, r.Changes,
			"only idleTimeout exists on N3IWF")
	})

	t.Run("missing global plmnID is created", func(t *testing.T) {
		cfg := newN3IWFConfig()
		cfg.Configuration.N3IWFInfo.GlobalN3IWFID.PLMNID = nil
		require.NoError(t, adapter.N3IWFAdapter{}.Apply(&global.Global{Plmn: &models.PlmnId{Mcc: "001", Mnc: "01"}}, cfg, &adapter.Recorder{}))
		assert.Equal(t, "001", cfg.Configuration.N3IWFInfo.GlobalN3IWFID.PLMNID.Mcc)
	})
}

func newTNGFConfig(ssts ...string) *tngffactory.Config {
	slices := []tngfcontext.SliceSupportItem{}
	sds := []string{"010203", "112233", "445566"}
	for i, sst := range ssts {
		slices = append(slices, tngfcontext.SliceSupportItem{SNSSAI: tngfcontext.SNSSAIItem{SST: sst, SD: sds[i]}})
	}
	return &tngffactory.Config{Configuration: &tngffactory.Configuration{
		TNGFInfo: tngfcontext.TNGFNFInfo{
			GlobalTNGFID: tngfcontext.GlobalTNGFID{PLMNID: tngfcontext.PLMNID{Mcc: "208", Mnc: "93"}, TNGFID: 135},
			SupportedTAList: []tngfcontext.SupportedTAItem{{
				TAC: "000001",
				BroadcastPLMNList: []tngfcontext.BroadcastPLMNItem{{
					PLMNID:              tngfcontext.PLMNID{Mcc: "208", Mnc: "93"},
					TAISliceSupportList: slices,
				}},
			}},
		},
	}}
}

func TestTNGFApplySstConversion(t *testing.T) {
	cfg := newTNGFConfig("1", "1")
	r := &adapter.Recorder{}
	require.NoError(t, adapter.TNGFAdapter{}.Apply(&global.Global{
		Plmn:    &models.PlmnId{Mcc: "001", Mnc: "01"},
		Tac:     ptr("00ABCD"),
		Snssais: []models.Snssai{{Sst: 1, Sd: "000001"}, {Sst: 10, Sd: "000002"}},
	}, cfg, r))

	info := cfg.Configuration.TNGFInfo
	assert.Equal(t, tngfcontext.PLMNID{Mcc: "001", Mnc: "01"}, info.GlobalTNGFID.PLMNID)
	ta := info.SupportedTAList[0]
	assert.Equal(t, "00ABCD", ta.TAC)
	assert.Equal(t, tngfcontext.PLMNID{Mcc: "001", Mnc: "01"}, ta.BroadcastPLMNList[0].PLMNID)
	assert.Equal(t, tngfcontext.SNSSAIItem{SST: "1", SD: "000001"}, ta.BroadcastPLMNList[0].TAISliceSupportList[0].SNSSAI)
	assert.Equal(t, tngfcontext.SNSSAIItem{SST: "A", SD: "000002"}, ta.BroadcastPLMNList[0].TAISliceSupportList[1].SNSSAI,
		"SST 10 is written as one uppercase hex digit")
	assert.Contains(t, r.Changes, adapter.Change{
		NF:   adapter.TNGF,
		Path: "TNGFInformation.SupportedTAList[0].BroadcastPLMNList[0].TAISliceSupportList[1].SNSSAI.SST",
		Old:  "1", New: "A",
	})
}

func TestTNGFApplyErrorsChangeNothing(t *testing.T) {
	cases := []struct {
		name    string
		cfg     *tngffactory.Config
		g       *global.Global
		wantErr string
	}{
		{
			name:    "global SST above 15",
			cfg:     newTNGFConfig("1", "1"),
			g:       &global.Global{Snssais: []models.Snssai{{Sst: 1, Sd: "000001"}, {Sst: 16, Sd: "000002"}}},
			wantErr: "global SST 16 cannot be represented by TNGF",
		},
		{
			name:    "existing SST is not one hex digit",
			cfg:     newTNGFConfig("1", "10"),
			g:       &global.Global{},
			wantErr: `invalid TNGF SST "10"`,
		},
		{
			name:    "existing SST is not hex",
			cfg:     newTNGFConfig("G"),
			g:       &global.Global{},
			wantErr: `invalid TNGF SST "G"`,
		},
		{
			name:    "count mismatch",
			cfg:     newTNGFConfig("1", "1"),
			g:       &global.Global{Snssais: []models.Snssai{{Sst: 1}}},
			wantErr: "global snssaiList has 1 entries but tngf SupportedTAList has 2",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.g.Plmn = &models.PlmnId{Mcc: "001", Mnc: "01"} // would change if applied
			r := &adapter.Recorder{}
			err := adapter.TNGFAdapter{}.Apply(tc.g, tc.cfg, r)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
			assert.Empty(t, r.Changes)
			assert.Equal(t, "208", tc.cfg.Configuration.TNGFInfo.GlobalTNGFID.PLMNID.Mcc)
		})
	}
}

func TestTNGFValidateChecksVersion(t *testing.T) {
	cfg := newTNGFConfig("1")
	cfg.Info = &tngffactory.Info{Version: "9.9.9"}
	err := adapter.TNGFAdapter{}.Validate(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), `info.version must be "1.0.3"`)
}
