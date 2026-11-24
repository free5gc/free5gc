package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/free5gc/config/cfg-generator/cmd"
	"github.com/free5gc/config/cfg-generator/generator"
)

func input(t *testing.T, global string) string {
	t.Helper()
	nf, err := os.ReadFile("../generator/testdata/amf-smf.yaml")
	require.NoError(t, err)
	p := filepath.Join(t.TempDir(), "input.yaml")
	require.NoError(t, os.WriteFile(p, []byte(global+"\n"+string(nf)), 0o600))
	return p
}

func runCLI(args ...string) (code int, stdout, stderr string) {
	var o, e bytes.Buffer
	code = cmd.Run(args, &o, &e)
	return code, o.String(), e.String()
}

func TestCLISuccess(t *testing.T) {
	out := t.TempDir()
	code, stdout, stderr := runCLI("-i", input(t, "global:\n  plmn: {mcc: \"466\", mnc: \"01\"}"), "-o", out)

	assert.Equal(t, 0, code, "stderr: %s", stderr)
	assert.Empty(t, stderr)
	assert.Contains(t, stdout, `amf: servedGuamiList[0].plmnId.mcc: "208" -> "466"`)
	assert.Contains(t, stdout, `smf: plmnList[0].mnc: "93" -> "01"`)
	assert.Contains(t, stdout, "8 change(s), 0 not applied, wrote [amfcfg.yaml smfcfg.yaml]")
	assert.FileExists(t, filepath.Join(out, "amfcfg.yaml"))
	assert.FileExists(t, filepath.Join(out, "smfcfg.yaml"))
}

// allInput writes an integrated input with every NF, so BSF and PCF are
// present and can report values they do not apply.
func allInput(t *testing.T, global string) string {
	t.Helper()
	nf, err := os.ReadFile("../generator/testdata/all-nfs.yaml")
	require.NoError(t, err)
	p := filepath.Join(t.TempDir(), "input.yaml")
	require.NoError(t, os.WriteFile(p, []byte(global+"\n"+string(nf)), 0o600))
	return p
}

const (
	// globalNotApplied sets a PLMN that BSF hard-codes differently.
	globalNotApplied = "global:\n  plmn: {mcc: \"466\", mnc: \"01\"}"
	// globalAllApplied matches BSF's hard-coded PLMN.
	globalAllApplied = "global:\n  plmn: {mcc: \"208\", mnc: \"93\"}"
	bsfNotAppliedOut = "[not applied] bsf: NRF registration PLMN stays 208/93 (hard-coded); global 466/01 not applied\n"
	strictErrOut     = "[strict] 1 global value(s) not applied; nothing written (run without --strict to accept)\n"
)

func TestCLINotAppliedGoesToStderrAndDoesNotFail(t *testing.T) {
	code, stdout, stderr := runCLI("-i", allInput(t, globalNotApplied), "-o", t.TempDir())
	assert.Equal(t, 0, code, "stderr: %s", stderr)
	assert.Equal(t, bsfNotAppliedOut, stderr)
	assert.Regexp(t, `\n\d+ change\(s\), 1 not applied, wrote \[upfcfg.yaml nrfcfg.yaml amfcfg.yaml`, stdout)
}

func TestCLIStrict(t *testing.T) {
	cases := []struct {
		name        string
		global      string
		extra       []string
		wantCode    int
		wantStderr  string
		wantSummary string // "" means stdout must be empty
		wantWrites  bool
	}{
		{
			name:       "strict with not-applied value fails before publish",
			global:     globalNotApplied,
			extra:      []string{"--strict"},
			wantCode:   6,
			wantStderr: bsfNotAppliedOut + strictErrOut,
		},
		{
			name:       "strict dry run with not-applied value reports and fails",
			global:     globalNotApplied,
			extra:      []string{"--strict", "--dry-run"},
			wantCode:   6,
			wantStderr: bsfNotAppliedOut + strictErrOut,
		},
		{
			name:        "strict with everything applied succeeds",
			global:      globalAllApplied,
			extra:       []string{"--strict"},
			wantSummary: " change(s), 0 not applied, wrote [upfcfg.yaml",
			wantWrites:  true,
		},
		{
			name:        "not strict with not-applied value succeeds",
			global:      globalNotApplied,
			wantStderr:  bsfNotAppliedOut,
			wantSummary: " change(s), 1 not applied, wrote [upfcfg.yaml",
			wantWrites:  true,
		},
		{
			name:        "dry run counts not-applied values",
			global:      globalNotApplied,
			extra:       []string{"--dry-run"},
			wantStderr:  bsfNotAppliedOut,
			wantSummary: " change(s), 1 not applied, would write [upfcfg.yaml",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := t.TempDir()
			require.NoError(t, os.WriteFile(filepath.Join(out, "bsfcfg.yaml"), []byte("hand written\n"), 0o600))
			before := dirSnapshot(t, out)

			code, stdout, stderr := runCLI(append([]string{"-i", allInput(t, tc.global), "-o", out}, tc.extra...)...)
			assert.Equal(t, tc.wantCode, code, "stderr: %s", stderr)
			assert.Equal(t, tc.wantStderr, stderr)
			if tc.wantSummary == "" {
				assert.Empty(t, stdout, "nothing reported on stdout on failure")
			} else {
				assert.Contains(t, stdout, tc.wantSummary)
			}
			if tc.wantWrites {
				assert.NotEqual(t, before, dirSnapshot(t, out))
			} else {
				assert.Equal(t, before, dirSnapshot(t, out), "output directory must be untouched")
			}
		})
	}
}

// dirSnapshot returns name -> content for every file in dir.
func dirSnapshot(t *testing.T, dir string) map[string]string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	out := map[string]string{}
	for _, e := range entries {
		b, err := os.ReadFile(filepath.Join(dir, e.Name()))
		require.NoError(t, err)
		out[e.Name()] = string(b)
	}
	return out
}

func TestCLIHelpDocumentsStrictAndExitCodes(t *testing.T) {
	code, stdout, _ := runCLI("--help")
	assert.Equal(t, 0, code)
	assert.Contains(t, stdout, "--strict")
	assert.Contains(t, stdout, "6  --strict and at least one global value was not applied")
}

func TestCLILongFlags(t *testing.T) {
	out := t.TempDir()
	code, _, stderr := runCLI("--input", input(t, ""), "--output", out)
	assert.Equal(t, 0, code, "stderr: %s", stderr)
	assert.FileExists(t, filepath.Join(out, "amfcfg.yaml"))
}

func TestCLIDryRun(t *testing.T) {
	out := t.TempDir()
	code, stdout, _ := runCLI("-i", input(t, "global:\n  plmn: {mcc: \"466\", mnc: \"01\"}"), "-o", out, "--dry-run")

	assert.Equal(t, 0, code)
	assert.Contains(t, stdout, "dry run: 8 change(s), 0 not applied, would write [amfcfg.yaml smfcfg.yaml]")
	entries, err := os.ReadDir(out)
	require.NoError(t, err)
	assert.Empty(t, entries)
}

func TestCLIExitCodes(t *testing.T) {
	cases := []struct {
		name     string
		args     func(t *testing.T) []string
		wantCode int
		wantErr  string
	}{
		{
			name:     "missing required -i",
			args:     func(*testing.T) []string { return nil },
			wantCode: 1,
			wantErr:  `required flag(s) "input" not set`,
		},
		{
			name:     "unknown flag",
			args:     func(t *testing.T) []string { return []string{"-i", input(t, ""), "--bogus"} },
			wantCode: 1,
			wantErr:  "unknown flag",
		},
		{
			name:     "positional argument rejected",
			args:     func(t *testing.T) []string { return []string{"-i", input(t, ""), "extra"} },
			wantCode: 1,
			wantErr:  "unknown command",
		},
		{
			name: "input file missing",
			args: func(t *testing.T) []string {
				return []string{"-i", filepath.Join(t.TempDir(), "none.yaml"), "-o", t.TempDir()}
			},
			wantCode: 1,
			wantErr:  "[input]",
		},
		{
			name: "invalid global",
			args: func(t *testing.T) []string {
				return []string{"-i", input(t, "global:\n  plmn: {mcc: \"4\", mnc: \"01\"}"), "-o", t.TempDir()}
			},
			wantCode: 2,
			wantErr:  "[global] global.plmn.mcc",
		},
		{
			name: "output dir is a file",
			args: func(t *testing.T) []string {
				f := filepath.Join(t.TempDir(), "f")
				require.NoError(t, os.WriteFile(f, nil, 0o600))
				return []string{"-i", input(t, ""), "-o", f}
			},
			wantCode: 5,
			wantErr:  "[publish]",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			code, stdout, stderr := runCLI(tc.args(t)...)
			assert.Equal(t, tc.wantCode, code)
			assert.Contains(t, stderr, tc.wantErr)
			assert.Empty(t, stdout, "nothing reported on failure")
		})
	}
}

func TestExitCodeMapping(t *testing.T) {
	// The CLI relies on generator.ExitCode; keep the documented values fixed.
	assert.Equal(t, 1, generator.StageInput.ExitCode())
	assert.Equal(t, 2, generator.StageGlobal.ExitCode())
	assert.Equal(t, 3, generator.StageAdapter.ExitCode())
	assert.Equal(t, 4, generator.StageValidate.ExitCode())
	assert.Equal(t, 5, generator.StagePublish.ExitCode())
	assert.Equal(t, 6, generator.StageStrict.ExitCode())
}
