package generator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/free5gc/config/cfg-generator/adapter"
	"github.com/free5gc/config/cfg-generator/generator"
)

func TestLoadFixture(t *testing.T) {
	in, err := generator.Load(writeInput(t, "global:\n  plmn: {mcc: \"466\", mnc: \"92\"}"))
	require.NoError(t, err)

	require.NotNil(t, in.Global)
	require.NotNil(t, in.Global.Plmn)
	assert.Equal(t, "466", in.Global.Plmn.Mcc)

	require.NotNil(t, in.AMF)
	require.NotNil(t, in.SMF)
	assert.Equal(t, "1.0.9", in.AMF.Info.Version)
	assert.Equal(t, "1.0.7", in.SMF.Info.Version)

	// Leading zeros in NF blocks survive because they decode into string fields.
	c := in.AMF.Configuration
	assert.Equal(t, "000001", c.SupportTAIList[0].Tac)
	assert.Equal(t, "010203", c.PlmnSupportList[0].SNssaiList[0].Sd)

	units := in.Units()
	require.Len(t, units, 2)
	assert.Equal(t, adapter.AMF, units[0].NF)
	assert.Equal(t, adapter.SMF, units[1].NF)
}

func TestLoadGlobalLeadingZeroUnquoted(t *testing.T) {
	// `mnc: 01` without quotes must stay "01", not become 1.
	in, err := generator.Load(writeInput(t, "global:\n  plmn:\n    mcc: 001\n    mnc: 01"))
	require.NoError(t, err)
	assert.Equal(t, "001", in.Global.Plmn.Mcc)
	assert.Equal(t, "01", in.Global.Plmn.Mnc)
}

func TestLoadWithoutGlobal(t *testing.T) {
	in, err := generator.Load(writeInput(t, ""))
	require.NoError(t, err)
	assert.Nil(t, in.Global)
	assert.Len(t, in.Units(), 2)
}

func TestLoadOnlyAMF(t *testing.T) {
	p := writeFile(t, "input.yaml", "amf:\n  info:\n    version: 1.0.9\n")
	in, err := generator.Load(p)
	require.NoError(t, err)
	assert.Nil(t, in.SMF)
	units := in.Units()
	require.Len(t, units, 1)
	assert.Equal(t, adapter.AMF, units[0].NF)
}

func TestLoadErrors(t *testing.T) {
	cases := []struct {
		name    string
		path    func(t *testing.T) string
		wantMsg string
	}{
		{
			name:    "missing file",
			path:    func(t *testing.T) string { return filepath.Join(t.TempDir(), "nope.yaml") },
			wantMsg: "nope.yaml",
		},
		{
			name:    "yaml syntax error",
			path:    func(t *testing.T) string { return writeFile(t, "in.yaml", "global: [unclosed\n") },
			wantMsg: "in.yaml",
		},
		{
			name: "unknown key inside global",
			path: func(t *testing.T) string {
				return writeInput(t, "global:\n  plmm: {mcc: \"466\", mnc: \"92\"}")
			},
			wantMsg: "plmm",
		},
		{
			name: "unknown key inside global.plmn",
			path: func(t *testing.T) string {
				return writeInput(t, "global:\n  plmn: {mcc: \"466\", mnx: \"92\"}")
			},
			wantMsg: "mnx",
		},
		{
			name: "unknown top-level key",
			path: func(t *testing.T) string {
				return writeFile(t, "in.yaml", "amff:\n  info: {}\nzzz: 1\n")
			},
			wantMsg: `unsupported top-level key "amff"`,
		},
		{
			name: "NF not supported yet (webui: webconsole not checked out)",
			path: func(t *testing.T) string {
				return writeFile(t, "in.yaml", "webui:\n  info: {}\n")
			},
			wantMsg: `unsupported top-level key "webui"`,
		},
		{
			name: "NF block of wrong shape",
			path: func(t *testing.T) string {
				return writeFile(t, "in.yaml", "amf: [1, 2]\n")
			},
			wantMsg: "in.yaml",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := generator.Load(tc.path(t))
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantMsg)
			assert.Contains(t, err.Error(), "[input]")
			assert.Equal(t, 1, generator.ExitCode(err))
		})
	}
}

func TestLoadUnknownTopLevelKeysAllReported(t *testing.T) {
	_, err := generator.Load(writeFile(t, "in.yaml", "bbb: 1\naaa: 2\n"))
	require.Error(t, err)
	var se *generator.StageError
	require.ErrorAs(t, err, &se)
	require.Len(t, se.Errs, 2)
	assert.Contains(t, se.Errs[0].Error(), `"aaa"`, "reported in sorted order")
	assert.Contains(t, se.Errs[1].Error(), `"bbb"`)
}

func TestLoadUnreadable(t *testing.T) {
	dir := t.TempDir()
	// A directory cannot be read as a file.
	require.NoError(t, os.Mkdir(filepath.Join(dir, "d.yaml"), 0o755))
	_, err := generator.Load(filepath.Join(dir, "d.yaml"))
	require.Error(t, err)
	assert.Equal(t, 1, generator.ExitCode(err))
}
