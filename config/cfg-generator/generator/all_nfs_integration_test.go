package generator_test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/free5gc/openapi/models"

	amffactory "github.com/free5gc/amf/pkg/factory"
	ausffactory "github.com/free5gc/ausf/pkg/factory"
	bsffactory "github.com/free5gc/bsf/pkg/factory"
	chffactory "github.com/free5gc/chf/pkg/factory"
	"github.com/free5gc/config/cfg-generator/adapter"
	"github.com/free5gc/config/cfg-generator/generator"
	upffactory "github.com/free5gc/go-upf/pkg/factory"
	n3iwffactory "github.com/free5gc/n3iwf/pkg/factory"
	neffactory "github.com/free5gc/nef/pkg/factory"
	nrffactory "github.com/free5gc/nrf/pkg/factory"
	nssffactory "github.com/free5gc/nssf/pkg/factory"
	pcffactory "github.com/free5gc/pcf/pkg/factory"
	scpfactory "github.com/free5gc/scp/pkg/factory"
	smffactory "github.com/free5gc/smf/pkg/factory"
	tngffactory "github.com/free5gc/tngf/pkg/factory"
	udmfactory "github.com/free5gc/udm/pkg/factory"
	udrfactory "github.com/free5gc/udr/pkg/factory"
)

const fixtureAllNFs = "testdata/all-nfs.yaml"

// allFiles are the outputs for the all-NF fixture, in run.sh start order.
var allFiles = []string{
	"upfcfg.yaml", "nrfcfg.yaml", "amfcfg.yaml", "smfcfg.yaml", "udrcfg.yaml",
	"pcfcfg.yaml", "udmcfg.yaml", "nssfcfg.yaml", "ausfcfg.yaml", "chfcfg.yaml",
	"nefcfg.yaml", "scpcfg.yaml", "bsfcfg.yaml", "n3iwfcfg.yaml", "tngfcfg.yaml",
}

// fullGlobalAll sets every global field. nrf matches the sample NRF SBI
// (127.0.0.10:8000) so the NRF consistency check passes; sbiScheme and the
// nrf scheme switch to https together.
const fullGlobalAll = `global:
  plmn: {mcc: 001, mnc: 01}
  tac: 00ABCD
  snssaiList:
    - {sst: 1, sd: 000001}
    - {sst: 2, sd: 000002}
  dnnList: [ims]
  nrf: {scheme: https, host: 127.0.0.10, port: 8000}
  sbiScheme: https
  mongodb: {name: free5gc2, url: "mongodb://db.local:27017"}
  logger: {level: debug}
  timeouts: {idle: 30s}`

func writeAllInput(t *testing.T, globalYAML string) string {
	t.Helper()
	nf, err := os.ReadFile(fixtureAllNFs)
	require.NoError(t, err)
	return writeFile(t, "input.yaml", globalYAML+"\n"+string(nf))
}

// decodeFile decodes an output file the way NFs do, without running their
// Validate (which has side effects such as filling nfInstanceId).
func decodeFile[T any](t *testing.T, path string) *T {
	t.Helper()
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	v := new(T)
	require.NoError(t, yaml.Unmarshal(b, v), path)
	return v
}

// sameConfig compares configs, treating nil and empty collections as equal
// and skipping the embedded sync.RWMutex.
func sameConfig(t *testing.T, name string, want, got any) {
	t.Helper()
	if diff := cmp.Diff(want, got, cmpopts.EquateEmpty(), cmpopts.IgnoreTypes(sync.RWMutex{})); diff != "" {
		t.Errorf("%s differs (-want +got):\n%s", name, diff)
	}
}

// requireAcceptedByNFs loads every output with the NF's own ReadConfig.
func requireAcceptedByNFs(t *testing.T, dir string) {
	t.Helper()
	p := func(f string) string { return filepath.Join(dir, f) }
	loaders := map[string]func(string) error{
		"upfcfg.yaml":   func(f string) error { _, err := upffactory.ReadConfig(f); return err },
		"nrfcfg.yaml":   func(f string) error { _, err := nrffactory.ReadConfig(f); return err },
		"amfcfg.yaml":   func(f string) error { _, err := amffactory.ReadConfig(f); return err },
		"smfcfg.yaml":   func(f string) error { _, err := smffactory.ReadConfig(f); return err },
		"udrcfg.yaml":   func(f string) error { _, err := udrfactory.ReadConfig(f); return err },
		"pcfcfg.yaml":   func(f string) error { _, err := pcffactory.ReadConfig(f); return err },
		"udmcfg.yaml":   func(f string) error { _, err := udmfactory.ReadConfig(f); return err },
		"nssfcfg.yaml":  func(f string) error { _, err := nssffactory.ReadConfig(f); return err },
		"ausfcfg.yaml":  func(f string) error { _, err := ausffactory.ReadConfig(f); return err },
		"chfcfg.yaml":   func(f string) error { _, err := chffactory.ReadConfig(f); return err },
		"nefcfg.yaml":   func(f string) error { _, err := neffactory.ReadConfig(f); return err },
		"scpcfg.yaml":   func(f string) error { _, err := scpfactory.ReadConfig(f); return err },
		"bsfcfg.yaml":   func(f string) error { _, err := bsffactory.ReadConfig(f); return err },
		"n3iwfcfg.yaml": func(f string) error { _, err := n3iwffactory.ReadConfig(f); return err },
		"tngfcfg.yaml":  func(f string) error { _, err := tngffactory.ReadConfig(f); return err },
	}
	require.Len(t, loaders, len(allFiles))
	for _, f := range allFiles {
		assert.NoError(t, loaders[f](p(f)), "%s must be accepted by its NF loader", f)
	}
}

// outputs decodes every output file.
type outputs struct {
	UPF   *upffactory.Config
	NRF   *nrffactory.Config
	AMF   *amffactory.Config
	SMF   *smffactory.Config
	UDR   *udrfactory.Config
	PCF   *pcffactory.Config
	UDM   *udmfactory.Config
	NSSF  *nssffactory.Config
	AUSF  *ausffactory.Config
	CHF   *chffactory.Config
	NEF   *neffactory.Config
	SCP   *scpfactory.Config
	BSF   *bsffactory.Config
	N3IWF *n3iwffactory.Config
	TNGF  *tngffactory.Config
}

func decodeAll(t *testing.T, dir string) outputs {
	p := func(f string) string { return filepath.Join(dir, f) }
	return outputs{
		UPF:   decodeFile[upffactory.Config](t, p("upfcfg.yaml")),
		NRF:   decodeFile[nrffactory.Config](t, p("nrfcfg.yaml")),
		AMF:   decodeFile[amffactory.Config](t, p("amfcfg.yaml")),
		SMF:   decodeFile[smffactory.Config](t, p("smfcfg.yaml")),
		UDR:   decodeFile[udrfactory.Config](t, p("udrcfg.yaml")),
		PCF:   decodeFile[pcffactory.Config](t, p("pcfcfg.yaml")),
		UDM:   decodeFile[udmfactory.Config](t, p("udmcfg.yaml")),
		NSSF:  decodeFile[nssffactory.Config](t, p("nssfcfg.yaml")),
		AUSF:  decodeFile[ausffactory.Config](t, p("ausfcfg.yaml")),
		CHF:   decodeFile[chffactory.Config](t, p("chfcfg.yaml")),
		NEF:   decodeFile[neffactory.Config](t, p("nefcfg.yaml")),
		SCP:   decodeFile[scpfactory.Config](t, p("scpcfg.yaml")),
		BSF:   decodeFile[bsffactory.Config](t, p("bsfcfg.yaml")),
		N3IWF: decodeFile[n3iwffactory.Config](t, p("n3iwfcfg.yaml")),
		TNGF:  decodeFile[tngffactory.Config](t, p("tngfcfg.yaml")),
	}
}

// compareAll compares every NF of want (an Input) with got.
func compareAll(t *testing.T, want *generator.Input, got outputs) {
	t.Helper()
	sameConfig(t, "upf", want.UPF, got.UPF)
	sameConfig(t, "nrf", want.NRF, got.NRF)
	sameConfig(t, "amf", want.AMF, got.AMF)
	sameConfig(t, "smf", want.SMF, got.SMF)
	sameConfig(t, "udr", want.UDR, got.UDR)
	sameConfig(t, "pcf", want.PCF, got.PCF)
	sameConfig(t, "udm", want.UDM, got.UDM)
	sameConfig(t, "nssf", want.NSSF, got.NSSF)
	sameConfig(t, "ausf", want.AUSF, got.AUSF)
	sameConfig(t, "chf", want.CHF, got.CHF)
	sameConfig(t, "nef", want.NEF, got.NEF)
	sameConfig(t, "scp", want.SCP, got.SCP)
	sameConfig(t, "bsf", want.BSF, got.BSF)
	sameConfig(t, "n3iwf", want.N3IWF, got.N3IWF)
	sameConfig(t, "tngf", want.TNGF, got.TNGF)
}

func TestAllNFsWithoutGlobal(t *testing.T) {
	out := t.TempDir()
	res, err := run(t, writeAllInput(t, ""), out)
	require.NoError(t, err)

	assert.Equal(t, allFiles, res.Files)
	assert.Empty(t, res.Changes)
	assert.Empty(t, res.Warnings)
	assert.Empty(t, res.NotApplied)

	requireAcceptedByNFs(t, out)

	orig, err := generator.Load(writeAllInput(t, ""))
	require.NoError(t, err)
	compareAll(t, orig, decodeAll(t, out))
}

// TestAllNFsFullGlobal applies every global field to every NF and compares
// the output with an expectation built by hand from the input. The
// expectation independently encodes docs/global-mapping.md; it does not
// reuse adapter code.
func TestAllNFsFullGlobal(t *testing.T) {
	out := t.TempDir()
	res, err := run(t, writeAllInput(t, fullGlobalAll), out)
	require.NoError(t, err)
	requireAcceptedByNFs(t, out)

	assert.ElementsMatch(t, []adapter.NotApplied{
		{NF: adapter.PCF, Item: "NRF registration DNN list", HardCoded: "[free5gc internet]", Requested: "ims"},
		{NF: adapter.BSF, Item: "NRF registration PLMN", HardCoded: "208/93", Requested: "001/01"},
	}, res.NotApplied)
	assert.Empty(t, res.Warnings)

	w, err := generator.Load(writeAllInput(t, ""))
	require.NoError(t, err)

	plmn := models.PlmnId{Mcc: "001", Mnc: "01"}
	s1, s2 := models.Snssai{Sst: 1, Sd: "000001"}, models.Snssai{Sst: 2, Sd: "000002"}
	idle := 30 * time.Second
	nrfURI := "https://127.0.0.10:8000"
	mongoName, mongoURL := "free5gc2", "mongodb://db.local:27017"

	// UPF: dnnList names; logger.
	for i := range w.UPF.DnnList {
		w.UPF.DnnList[i].Dnn = "ims"
	}
	w.UPF.Logger.Level = "debug"

	// NRF: DefaultPlmnId, flat MongoDB fields, scheme, logger, idle. No nrfUri.
	n := w.NRF.Configuration
	n.DefaultPlmnId = plmn
	n.MongoDBName, n.MongoDBUrl = mongoName, mongoURL
	n.Sbi.Scheme = "https"
	n.IdleTimeout = &idle
	w.NRF.Logger.Level = "debug"

	// AMF.
	a := w.AMF.Configuration
	a.ServedGumaiList[0].PlmnId.Mcc, a.ServedGumaiList[0].PlmnId.Mnc = plmn.Mcc, plmn.Mnc
	*a.SupportTAIList[0].PlmnId = plmn
	a.SupportTAIList[0].Tac = "00ABCD"
	*a.PlmnSupportList[0].PlmnId = plmn
	a.PlmnSupportList[0].SNssaiList = []models.Snssai{s1, s2}
	a.SupportDnnList = []string{"ims"}
	a.NrfUri, a.Sbi.Scheme, a.IdleTimeout = nrfURI, "https", &idle
	w.AMF.Logger.Level = "debug"

	// SMF.
	s := w.SMF.Configuration
	s.PLMNList[0] = smffactory.PlmnID{Mcc: plmn.Mcc, Mnc: plmn.Mnc}
	*s.SNssaiInfo[0].SNssai, *s.SNssaiInfo[1].SNssai = s1, s2
	for _, info := range s.SNssaiInfo {
		info.DnnInfos[0].Dnn = "ims"
	}
	upf := s.UserPlaneInformation.UPNodes["UPF"]
	*upf.SNssaiInfos[0].SNssai, *upf.SNssaiInfos[1].SNssai = s1, s2
	for _, info := range upf.SNssaiInfos {
		info.DnnUpfInfoList[0].Dnn = "ims"
	}
	upf.InterfaceUpfInfoList[0].NetworkInstances = []string{"ims"}
	s.NrfUri, s.Sbi.Scheme, s.IdleTimeout = nrfURI, "https", &idle
	w.SMF.Logger.Level = "debug"

	// UDR, PCF, CHF, BSF: nrfUri, scheme, mongodb, logger, idle.
	u := w.UDR.Configuration
	u.NrfUri, u.Sbi.Scheme, u.IdleTimeout = nrfURI, "https", &idle
	u.Mongodb.Name, u.Mongodb.Url = mongoName, mongoURL
	w.UDR.Logger.Level = "debug"

	pc := w.PCF.Configuration
	pc.NrfUri, pc.Sbi.Scheme, pc.IdleTimeout = nrfURI, "https", &idle
	pc.Mongodb.Name, pc.Mongodb.Url = mongoName, mongoURL
	w.PCF.Logger.Level = "debug"

	ch := w.CHF.Configuration
	ch.NrfUri, ch.Sbi.Scheme, ch.IdleTimeout = nrfURI, "https", &idle
	ch.Mongodb.Name, ch.Mongodb.Url = mongoName, mongoURL
	w.CHF.Logger.Level = "debug"

	b := w.BSF.Configuration
	b.NrfUri, b.Sbi.Scheme, b.IdleTimeout = nrfURI, "https", &idle
	b.MongoDB.Name, b.MongoDB.Url = mongoName, mongoURL
	w.BSF.Logger.Level = "debug"

	// UDM, NEF, SCP: nrfUri, scheme, logger, idle.
	ud := w.UDM.Configuration
	ud.NrfUri, ud.Sbi.Scheme, ud.IdleTimeout = nrfURI, "https", &idle
	w.UDM.Logger.Level = "debug"

	ne := w.NEF.Configuration
	ne.NrfUri, ne.Sbi.Scheme, ne.IdleTimeout = nrfURI, "https", &idle
	w.NEF.Logger.Level = "debug"

	sc := w.SCP.Configuration
	sc.NrfUri, sc.Sbi.Scheme, sc.IdleTimeout = nrfURI, "https", &idle
	w.SCP.Logger.Level = "debug"

	// AUSF: only the first PLMN; the second (123/45) is kept.
	au := w.AUSF.Configuration
	au.PlmnSupportList[0] = plmn
	au.NrfUri, au.Sbi.Scheme, au.IdleTimeout = nrfURI, "https", &idle
	w.AUSF.Logger.Level = "debug"

	// NSSF: home PLMN lists; NRF URI prefixes; S-NSSAI not applied.
	ns := w.NSSF.Configuration
	for i := range ns.SupportedPlmnList {
		ns.SupportedPlmnList[i] = plmn
	}
	for i := range ns.SupportedNssaiInPlmnList {
		*ns.SupportedNssaiInPlmnList[i].PlmnId = plmn
	}
	const oldNrf = "http://127.0.0.10:8000"
	for i := range ns.NsiList {
		for j := range ns.NsiList[i].NsiInformationList {
			info := &ns.NsiList[i].NsiInformationList[j]
			if len(info.NrfId) >= len(oldNrf) && info.NrfId[:len(oldNrf)] == oldNrf {
				info.NrfId = nrfURI + info.NrfId[len(oldNrf):]
			}
		}
	}
	for i := range ns.AmfSetList {
		set := &ns.AmfSetList[i]
		if len(set.NrfAmfSet) >= len(oldNrf) && set.NrfAmfSet[:len(oldNrf)] == oldNrf {
			set.NrfAmfSet = nrfURI + set.NrfAmfSet[len(oldNrf):]
		}
	}
	ns.NrfUri, ns.Sbi.Scheme, ns.IdleTimeout = nrfURI, models.UriScheme("https"), &idle
	w.NSSF.Logger.Level = "debug"

	// N3IWF: PLMN, TAC, S-NSSAI (SST as int), logger, idle.
	ni := w.N3IWF.Configuration
	*ni.N3IWFInfo.GlobalN3IWFID.PLMNID = n3iwffactory.PLMNID{Mcc: plmn.Mcc, Mnc: plmn.Mnc}
	for i := range ni.N3IWFInfo.SupportedTAList {
		ta := &ni.N3IWFInfo.SupportedTAList[i]
		ta.TAC = "00ABCD"
		for j := range ta.BroadcastPLMNList {
			bp := &ta.BroadcastPLMNList[j]
			*bp.PLMNID = n3iwffactory.PLMNID{Mcc: plmn.Mcc, Mnc: plmn.Mnc}
			bp.TAISliceSupportList[0].SNSSAI = n3iwffactory.SNSSAIItem{SST: 1, SD: "000001"}
			bp.TAISliceSupportList[1].SNSSAI = n3iwffactory.SNSSAIItem{SST: 2, SD: "000002"}
		}
	}
	ni.IdleTimeout = &idle
	w.N3IWF.Logger.Level = "debug"

	// TNGF: same, SST as one hex digit.
	tg := w.TNGF.Configuration
	tg.TNGFInfo.GlobalTNGFID.PLMNID.Mcc, tg.TNGFInfo.GlobalTNGFID.PLMNID.Mnc = plmn.Mcc, plmn.Mnc
	for i := range tg.TNGFInfo.SupportedTAList {
		ta := &tg.TNGFInfo.SupportedTAList[i]
		ta.TAC = "00ABCD"
		for j := range ta.BroadcastPLMNList {
			bp := &ta.BroadcastPLMNList[j]
			bp.PLMNID.Mcc, bp.PLMNID.Mnc = plmn.Mcc, plmn.Mnc
			bp.TAISliceSupportList[0].SNSSAI.SST, bp.TAISliceSupportList[0].SNSSAI.SD = "1", "000001"
			bp.TAISliceSupportList[1].SNSSAI.SST, bp.TAISliceSupportList[1].SNSSAI.SD = "2", "000002"
		}
	}
	tg.IdleTimeout = &idle
	w.TNGF.Logger.Level = "debug"

	compareAll(t, w, decodeAll(t, out))

	// Spot-check the serialized text: leading zeros kept, NRF prefix kept for
	// the foreign NRF entry.
	nssf, err := os.ReadFile(filepath.Join(out, "nssfcfg.yaml"))
	require.NoError(t, err)
	assert.Contains(t, string(nssf), "nrfId: https://127.0.0.10:8000/nnrf-nfm/v1/nf-instances")
	assert.Contains(t, string(nssf), "http://localhost:8084/nnrf-nfm/v1/nf-instances", "foreign NRF kept")
	tngf, err := os.ReadFile(filepath.Join(out, "tngfcfg.yaml"))
	require.NoError(t, err)
	assert.Contains(t, string(tngf), `SST: "2"`)
}

func TestAllNFsAdapterErrors(t *testing.T) {
	cases := []struct {
		name    string
		global  string
		wantMsg []string
	}{
		{
			name:    "global nrf does not match the NRF's own SBI",
			global:  "global:\n  nrf: {scheme: http, host: 10.0.0.10, port: 8000}",
			wantMsg: []string{"[adapter] nrf: global nrf http://10.0.0.10:8000 does not match the NRF's own SBI http://127.0.0.10:8000"},
		},
		{
			name:    "sbiScheme changes the NRF scheme, nrf keeps the old one",
			global:  "global:\n  sbiScheme: http\n  nrf: {scheme: http, host: 127.0.0.10, port: 8001}",
			wantMsg: []string{"[adapter] nrf:", "8001"},
		},
		{
			name:   "SST beyond one hex digit",
			global: "global:\n  snssaiList:\n    - {sst: 1, sd: \"000001\"}\n    - {sst: 16, sd: \"000002\"}",
			wantMsg: []string{
				"[adapter] tngf: global SST 16 cannot be represented by TNGF",
			},
		},
		{
			name:    "UPF DNN count differs",
			global:  "global:\n  dnnList: [a, b]",
			wantMsg: []string{"[adapter] upf: global dnnList has 2 entries but upf dnnList has 1 distinct value(s) [internet]"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := t.TempDir()
			_, err := run(t, writeAllInput(t, tc.global), out)
			require.Error(t, err)
			assert.Equal(t, 3, generator.ExitCode(err))
			for _, m := range tc.wantMsg {
				assert.Contains(t, err.Error(), m)
			}
			assert.Empty(t, snapshot(t, out), "nothing written")
		})
	}
}

func TestAllNFsMongoDBPartialOverride(t *testing.T) {
	out := t.TempDir()
	res, err := run(t, writeAllInput(t, "global:\n  mongodb: {url: \"mongodb://db:27017\"}"), out)
	require.NoError(t, err)

	got := decodeAll(t, out)
	assert.Equal(t, "mongodb://db:27017", got.NRF.Configuration.MongoDBUrl)
	assert.Equal(t, "free5gc", got.NRF.Configuration.MongoDBName, "name not provided, kept")
	assert.Equal(t, "mongodb://db:27017", got.UDR.Configuration.Mongodb.Url)
	assert.Equal(t, "mongodb://db:27017", got.BSF.Configuration.MongoDB.Url)

	var paths []string
	for _, c := range res.Changes {
		paths = append(paths, string(c.NF)+":"+c.Path)
	}
	// Exactly the NFs with a database: nrf (flat key), udr, pcf, chf, bsf.
	assert.ElementsMatch(t, []string{
		"nrf:MongoDBUrl", "udr:mongodb.url", "pcf:mongodb.url", "chf:mongodb.url", "bsf:mongodb.url",
	}, paths)
}
