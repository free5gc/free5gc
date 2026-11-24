package generator_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/free5gc/config/cfg-generator/internal/testfixture"
)

// fixedNow is the clock used for backup names in tests.
var fixedNow = time.Date(2026, 11, 22, 15, 30, 0, 0, time.UTC)

func writeInput(t *testing.T, globalYAML string) string {
	t.Helper()
	return writeFile(t, "input.yaml", globalYAML+"\n"+testfixture.NFBlocks(t, "amf", "smf"))
}

// writeFile writes content to a new file in a fresh temp dir.
func writeFile(t *testing.T, name, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	require.NoError(t, os.WriteFile(p, []byte(content), 0o600))
	return p
}

// snapshot returns name -> content for every entry in dir.
func snapshot(t *testing.T, dir string) map[string]string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	out := map[string]string{}
	for _, e := range entries {
		if e.IsDir() {
			out[e.Name()+"/"] = ""
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, e.Name()))
		require.NoError(t, err)
		out[e.Name()] = string(b)
	}
	return out
}
