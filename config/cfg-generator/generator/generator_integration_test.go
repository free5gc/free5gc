package generator_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	amffactory "github.com/free5gc/amf/pkg/factory"
	"github.com/free5gc/config/cfg-generator/generator"
	smffactory "github.com/free5gc/smf/pkg/factory"
)

const globalPlmn46601 = "global:\n  plmn:\n    mcc: \"466\"\n    mnc: \"01\""

func run(t *testing.T, input, out string) (*generator.Result, error) {
	t.Helper()
	return generator.Run(generator.Options{
		InputPath: input,
		OutDir:    out,
		Now:       func() time.Time { return fixedNow },
	})
}

// TestEndToEndOutputIsAcceptedByNFs is the core IT: generated files must be
// readable and valid for the real NF loaders (factory.ReadConfig runs the
// NF's own Validate), carry the global PLMN, and keep every other value.
func TestEndToEndOutputIsAcceptedByNFs(t *testing.T) {
	input := writeInput(t, globalPlmn46601)
	out := t.TempDir()

	res, err := run(t, input, out)
	require.NoError(t, err)
	assert.Equal(t, []string{"amfcfg.yaml", "smfcfg.yaml"}, res.Files)
	// AMF: 1 guami + 1 tai + 1 plmnSupport; SMF: 1 plmnList entry; mcc+mnc each.
	assert.Len(t, res.Changes, 8)

	amfOut, err := amffactory.ReadConfig(filepath.Join(out, "amfcfg.yaml"))
	require.NoError(t, err, "AMF must accept the generated file")
	smfOut, err := smffactory.ReadConfig(filepath.Join(out, "smfcfg.yaml"))
	require.NoError(t, err, "SMF must accept the generated file")

	// The global PLMN is applied everywhere it belongs.
	ac := amfOut.Configuration
	assert.Equal(t, "466", ac.ServedGumaiList[0].PlmnId.Mcc)
	assert.Equal(t, "01", ac.ServedGumaiList[0].PlmnId.Mnc)
	assert.Equal(t, "466", ac.SupportTAIList[0].PlmnId.Mcc)
	assert.Equal(t, "01", ac.SupportTAIList[0].PlmnId.Mnc)
	assert.Equal(t, "466", ac.PlmnSupportList[0].PlmnId.Mcc)
	assert.Equal(t, "01", ac.PlmnSupportList[0].PlmnId.Mnc)
	assert.Equal(t, smffactory.PlmnID{Mcc: "466", Mnc: "01"}, smfOut.Configuration.PLMNList[0])

	// Everything else is identical to the input NF blocks.
	orig, err := generator.Load(writeInput(t, ""))
	require.NoError(t, err)

	// ReadConfig runs the NF's Validate, which fills an empty nfInstanceId
	// with a random UUID at load time. That is NF runtime behaviour, not file
	// content (TestEndToEndDoesNotFreezeNfInstanceId checks the file).
	amfOut.Configuration.NfInstanceId = orig.AMF.Configuration.NfInstanceId
	smfOut.Configuration.NfInstanceId = orig.SMF.Configuration.NfInstanceId

	clearAMFPlmn(orig.AMF)
	clearAMFPlmn(amfOut)
	assertSameConfig(t, orig.AMF.Info, amfOut.Info)
	assertSameConfig(t, orig.AMF.Configuration, amfOut.Configuration)
	assertSameConfig(t, orig.AMF.Logger, amfOut.Logger)

	clearSMFPlmn(orig.SMF)
	clearSMFPlmn(smfOut)
	assertSameConfig(t, orig.SMF.Info, smfOut.Info)
	assertSameConfig(t, orig.SMF.Configuration, smfOut.Configuration)
	assertSameConfig(t, orig.SMF.Logger, smfOut.Logger)
}

// assertSameConfig compares two decoded configs, treating nil and empty
// slices/maps as equal: fields without omitempty (e.g. SMF dnaiList) are
// written as `[]` and decode back as empty slices, which NFs treat the same
// as absent.
func assertSameConfig(t *testing.T, want, got any) {
	t.Helper()
	if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
		t.Errorf("config differs (-input +output):\n%s", diff)
	}
}

func clearAMFPlmn(c *amffactory.Config) {
	for i := range c.Configuration.ServedGumaiList {
		c.Configuration.ServedGumaiList[i].PlmnId.Mcc, c.Configuration.ServedGumaiList[i].PlmnId.Mnc = "", ""
	}
	for i := range c.Configuration.SupportTAIList {
		c.Configuration.SupportTAIList[i].PlmnId.Mcc, c.Configuration.SupportTAIList[i].PlmnId.Mnc = "", ""
	}
	for i := range c.Configuration.PlmnSupportList {
		c.Configuration.PlmnSupportList[i].PlmnId.Mcc, c.Configuration.PlmnSupportList[i].PlmnId.Mnc = "", ""
	}
}

func clearSMFPlmn(c *smffactory.Config) {
	for i := range c.Configuration.PLMNList {
		c.Configuration.PLMNList[i] = smffactory.PlmnID{}
	}
}

// TestEndToEndDoesNotFreezeNfInstanceId guards against NF Validate side
// effects: AMF/SMF fill an empty nfInstanceId with a random UUID. The sample
// configs leave it empty so each NF start generates its own; the generated
// file must keep it empty.
func TestEndToEndDoesNotFreezeNfInstanceId(t *testing.T) {
	out := t.TempDir()
	_, err := run(t, writeInput(t, globalPlmn46601), out)
	require.NoError(t, err)

	for name, content := range snapshot(t, out) {
		assert.NotContains(t, content, "nfInstanceId", name)
	}
}

func TestEndToEndWithoutGlobalOnlyNormalizes(t *testing.T) {
	out := t.TempDir()
	res, err := run(t, writeInput(t, ""), out)
	require.NoError(t, err)
	assert.Empty(t, res.Changes)

	amfOut, err := amffactory.ReadConfig(filepath.Join(out, "amfcfg.yaml"))
	require.NoError(t, err)
	assert.Equal(t, "208", amfOut.Configuration.PlmnSupportList[0].PlmnId.Mcc)
}

func TestEndToEndIsIdempotent(t *testing.T) {
	input := writeInput(t, globalPlmn46601)
	out := t.TempDir()

	_, err := run(t, input, out)
	require.NoError(t, err)
	first := snapshot(t, out)

	_, err = run(t, input, out)
	require.NoError(t, err)
	assert.Equal(t, first, snapshot(t, out), "second run: same content, no backups")
}

func TestEndToEndMigratesHandWrittenConfigs(t *testing.T) {
	out := t.TempDir()
	amfSample, err := os.ReadFile("../../amfcfg.yaml")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(out, "amfcfg.yaml"), amfSample, 0o600))

	_, err = run(t, writeInput(t, globalPlmn46601), out)
	require.NoError(t, err)

	files := snapshot(t, out)
	assert.Equal(t, string(amfSample), files["amfcfg-20261122T153000.bak.yaml"], "original kept as backup")
	assert.Contains(t, files["amfcfg.yaml"], generator.GeneratedMarker)
	assert.NotContains(t, files, "smfcfg-20261122T153000.bak.yaml", "no smfcfg existed, nothing to back up")
}

func TestEndToEndDryRunWritesNothing(t *testing.T) {
	out := t.TempDir()
	res, err := generator.Run(generator.Options{
		InputPath: writeInput(t, globalPlmn46601),
		OutDir:    out,
		DryRun:    true,
	})
	require.NoError(t, err)
	assert.Len(t, res.Changes, 8)
	assert.Equal(t, []string{"amfcfg.yaml", "smfcfg.yaml"}, res.Files)
	assert.Empty(t, snapshot(t, out))
}

// TestEndToEndFailuresLeaveOutputUntouched covers exit codes 1-4: each must
// report the right stage and must not modify the output directory.
func TestEndToEndFailuresLeaveOutputUntouched(t *testing.T) {
	nfBlocks, err := os.ReadFile(fixtureNFBlocks)
	require.NoError(t, err)

	cases := []struct {
		name     string
		input    func(t *testing.T) string
		wantCode int
		wantMsg  []string
	}{
		{
			name:     "input: unknown global key",
			input:    func(t *testing.T) string { return writeInput(t, "global:\n  plmnn: {}") },
			wantCode: 1,
			wantMsg:  []string{"[input]", "plmnn"},
		},
		{
			name: "global: invalid mcc and mnc, both reported",
			input: func(t *testing.T) string {
				return writeInput(t, "global:\n  plmn: {mcc: \"46\", mnc: \"x\"}")
			},
			wantCode: 2,
			wantMsg:  []string{"[global] global.plmn.mcc", "[global] global.plmn.mnc"},
		},
		{
			name: "adapter: amf block without configuration",
			input: func(t *testing.T) string {
				return writeFile(t, "in.yaml", globalPlmn46601+"\namf:\n  info:\n    version: 1.0.9\n")
			},
			wantCode: 3,
			wantMsg:  []string{"[adapter] amf: configuration block is missing"},
		},
		{
			name: "validate: NF rejects its config",
			input: func(t *testing.T) string {
				// Break the AMF version; AMF Validate only accepts 1.0.9.
				broken := []byte(string(nfBlocks))
				broken = replaceOnce(t, broken, "version: 1.0.9", "version: 9.9.9")
				return writeFile(t, "in.yaml", globalPlmn46601+"\n"+string(broken))
			},
			wantCode: 4,
			wantMsg:  []string{"[validate] amf:"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := t.TempDir()
			existing := "hand written\n"
			require.NoError(t, os.WriteFile(filepath.Join(out, "amfcfg.yaml"), []byte(existing), 0o600))
			before := snapshot(t, out)

			_, err := run(t, tc.input(t), out)
			require.Error(t, err)
			assert.Equal(t, tc.wantCode, generator.ExitCode(err))
			for _, m := range tc.wantMsg {
				assert.Contains(t, err.Error(), m)
			}
			assert.Equal(t, before, snapshot(t, out), "output directory must be untouched")
		})
	}
}

func TestEndToEndPublishFailureExitCode(t *testing.T) {
	outFile := writeFile(t, "config", "not a dir")
	_, err := run(t, writeInput(t, globalPlmn46601), outFile)
	require.Error(t, err)
	assert.Equal(t, 5, generator.ExitCode(err))
	assert.Contains(t, err.Error(), "[publish]")
}

// TestEndToEndStrict covers exit code 6: with Strict, a not-applied global
// value fails the run before publish, and the result still lists it.
func TestEndToEndStrict(t *testing.T) {
	const (
		plmnNotApplied = "global:\n  plmn: {mcc: \"001\", mnc: \"01\"}" // BSF hard-codes 208/93
		plmnApplied    = "global:\n  plmn: {mcc: \"208\", mnc: \"93\"}"
	)
	cases := []struct {
		name       string
		global     string
		strict     bool
		dryRun     bool
		wantCode   int
		wantN      int
		wantWrites bool
	}{
		{"strict, not applied", plmnNotApplied, true, false, 6, 1, false},
		{"strict dry run, not applied", plmnNotApplied, true, true, 6, 1, false},
		{"strict, everything applied", plmnApplied, true, false, 0, 0, true},
		{"not strict, not applied", plmnNotApplied, false, false, 0, 1, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := t.TempDir()
			require.NoError(t, os.WriteFile(filepath.Join(out, "bsfcfg.yaml"), []byte("hand written\n"), 0o600))
			before := snapshot(t, out)

			res, err := generator.Run(generator.Options{
				InputPath: writeAllInput(t, tc.global),
				OutDir:    out,
				Strict:    tc.strict,
				DryRun:    tc.dryRun,
				Now:       func() time.Time { return fixedNow },
			})
			assert.Equal(t, tc.wantCode, generator.ExitCode(err))
			require.NotNil(t, res, "the result is reported even on a strict failure")
			assert.Len(t, res.NotApplied, tc.wantN)
			if tc.wantCode != 0 {
				assert.Contains(t, err.Error(), "[strict] 1 global value(s) not applied; nothing written")
			}
			if tc.wantWrites {
				assert.NotEqual(t, before, snapshot(t, out))
			} else {
				assert.Equal(t, before, snapshot(t, out), "output directory must be untouched")
			}
		})
	}
}

func replaceOnce(t *testing.T, b []byte, old, repl string) []byte {
	t.Helper()
	s := string(b)
	i := indexOf(s, old)
	require.GreaterOrEqual(t, i, 0, "fixture must contain %q", old)
	return []byte(s[:i] + repl + s[i+len(old):])
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
