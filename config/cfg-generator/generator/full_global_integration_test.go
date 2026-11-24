package generator_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/free5gc/openapi/models"

	amffactory "github.com/free5gc/amf/pkg/factory"
	"github.com/free5gc/config/cfg-generator/generator"
	smffactory "github.com/free5gc/smf/pkg/factory"
)

// fullGlobal sets every field supported for AMF/SMF. Values are chosen so
// that each one differs from the sample configs, including leading zeros.
const fullGlobal = `global:
  plmn: {mcc: 001, mnc: 01}
  tac: 00ABCD
  snssaiList:
    - {sst: 1, sd: 000001}
    - {sst: 2, sd: 000002}
  dnnList: [ims]
  nrf: {scheme: https, host: 10.0.0.10, port: 29510}
  sbiScheme: https
  logger: {level: debug}
  timeouts: {idle: 30s}`

// TestEndToEndFullGlobal applies every supported global field and checks the
// output, as loaded by the real NF loaders, against an expectation built by
// hand from the input. The expectation does not reuse adapter code, so it
// independently encodes docs/global-mapping.md.
func TestEndToEndFullGlobal(t *testing.T) {
	out := t.TempDir()
	res, err := run(t, writeInput(t, fullGlobal), out)
	require.NoError(t, err)

	amfOut, err := amffactory.ReadConfig(filepath.Join(out, "amfcfg.yaml"))
	require.NoError(t, err, "AMF must accept the generated file")
	smfOut, err := smffactory.ReadConfig(filepath.Join(out, "smfcfg.yaml"))
	require.NoError(t, err, "SMF must accept the generated file")

	orig, err := generator.Load(writeInput(t, ""))
	require.NoError(t, err)
	wantAMF, wantSMF := orig.AMF, orig.SMF

	plmn := models.PlmnId{Mcc: "001", Mnc: "01"}
	s1, s2 := models.Snssai{Sst: 1, Sd: "000001"}, models.Snssai{Sst: 2, Sd: "000002"}
	idle := 30 * time.Second

	// AMF expectation (global-mapping §1).
	a := wantAMF.Configuration
	a.ServedGumaiList[0].PlmnId.Mcc, a.ServedGumaiList[0].PlmnId.Mnc = plmn.Mcc, plmn.Mnc
	*a.SupportTAIList[0].PlmnId = plmn
	a.SupportTAIList[0].Tac = "00ABCD"
	*a.PlmnSupportList[0].PlmnId = plmn
	a.PlmnSupportList[0].SNssaiList = []models.Snssai{s1, s2} // sample: 010203, 112233
	a.SupportDnnList = []string{"ims"}                        // sample: internet
	a.NrfUri = "https://10.0.0.10:29510"
	a.Sbi.Scheme = "https"
	a.IdleTimeout = &idle
	wantAMF.Logger.Level = "debug"

	// SMF expectation: session and user-plane slices stay aligned.
	s := wantSMF.Configuration
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
	s.NrfUri = "https://10.0.0.10:29510"
	s.Sbi.Scheme = "https"
	s.IdleTimeout = &idle
	wantSMF.Logger.Level = "debug"

	// nfInstanceId is filled by ReadConfig at load time, not by the file.
	amfOut.Configuration.NfInstanceId = ""
	smfOut.Configuration.NfInstanceId = ""

	assertSameConfig(t, wantAMF.Info, amfOut.Info)
	assertSameConfig(t, wantAMF.Configuration, amfOut.Configuration)
	assertSameConfig(t, wantAMF.Logger, amfOut.Logger)
	assertSameConfig(t, wantSMF.Info, smfOut.Info)
	assertSameConfig(t, wantSMF.Configuration, smfOut.Configuration)
	assertSameConfig(t, wantSMF.Logger, smfOut.Logger)

	// AMF 15 = plmn 6 + tac 1 + snssai 3 + dnn 1 + nrf/scheme/level/idle 4.
	// SMF 17 = plmn 2 + snssai 3+3 + dnn 2+2 + networkInstance 1 + 4.
	assert.Len(t, res.Changes, 32)

	// Leading zeros survive serialization in the file itself.
	amfFile, err := os.ReadFile(filepath.Join(out, "amfcfg.yaml"))
	require.NoError(t, err)
	assert.Contains(t, string(amfFile), `mcc: "001"`)
	assert.Contains(t, string(amfFile), `mnc: "01"`)
	assert.Contains(t, string(amfFile), `sd: "000001"`)
}

func TestEndToEndValueMappingMismatchExitCode(t *testing.T) {
	out := t.TempDir()
	before := snapshot(t, out)

	// The sample has 2 distinct S-NSSAIs; 3 given.
	_, err := run(t, writeInput(t, `global:
  snssaiList:
    - {sst: 1, sd: "000001"}
    - {sst: 1, sd: "000002"}
    - {sst: 1, sd: "000003"}`), out)
	require.Error(t, err)
	assert.Equal(t, 3, generator.ExitCode(err))
	// Both NFs are reported in one run.
	assert.Contains(t, err.Error(), "[adapter] amf: global snssaiList has 3 entries")
	assert.Contains(t, err.Error(), "[adapter] smf: global snssaiList has 3 entries")
	assert.Equal(t, before, snapshot(t, out))
}

func TestEndToEndGlobalErrorsAllReported(t *testing.T) {
	_, err := run(t, writeInput(t, `global:
  tac: "1"
  dnnList: []
  sbiScheme: ftp
  logger: {level: loud}`), t.TempDir())
	require.Error(t, err)
	assert.Equal(t, 2, generator.ExitCode(err))
	for _, p := range []string{"global.tac", "global.dnnList", "global.sbiScheme", "global.logger.level"} {
		assert.Contains(t, err.Error(), p)
	}
}
