package generator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/free5gc/config/cfg-generator/generator"
)

func generated(body string) []byte {
	return []byte(generator.GeneratedMarker + "\n" + body)
}

func TestBackupName(t *testing.T) {
	assert.Equal(t, "amfcfg-20261122T153000.bak.yaml", generator.BackupName("amfcfg.yaml", fixedNow))
	assert.Equal(t, "smfcfg-20261122T153000.bak.yaml", generator.BackupName("smfcfg.yaml", fixedNow))
	assert.NotContains(t, generator.BackupName("amfcfg.yaml", fixedNow), ":", "must be valid on Windows")
}

func TestPublishNewFiles(t *testing.T) {
	out := filepath.Join(t.TempDir(), "config") // does not exist yet
	files := []generator.File{
		{Name: "amfcfg.yaml", Data: generated("a: 1\n")},
		{Name: "smfcfg.yaml", Data: generated("s: 1\n")},
	}
	require.NoError(t, generator.Publish(out, files, fixedNow))

	assert.Equal(t, map[string]string{
		"amfcfg.yaml": string(files[0].Data),
		"smfcfg.yaml": string(files[1].Data),
	}, snapshot(t, out), "only the outputs, temp dir removed")
}

func TestPublishBacksUpHandWrittenFile(t *testing.T) {
	out := t.TempDir()
	original := "# free5gc sample\ninfo:\n  version: 1.0.9\n"
	require.NoError(t, os.WriteFile(filepath.Join(out, "amfcfg.yaml"), []byte(original), 0o600))

	newData := generated("a: 2\n")
	require.NoError(t, generator.Publish(out, []generator.File{{Name: "amfcfg.yaml", Data: newData}}, fixedNow))

	assert.Equal(t, map[string]string{
		"amfcfg.yaml":                     string(newData),
		"amfcfg-20261122T153000.bak.yaml": original,
	}, snapshot(t, out))
}

func TestPublishOverwritesGeneratedFileWithoutBackup(t *testing.T) {
	out := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(out, "amfcfg.yaml"), generated("a: 1\n"), 0o600))

	newData := generated("a: 2\n")
	require.NoError(t, generator.Publish(out, []generator.File{{Name: "amfcfg.yaml", Data: newData}}, fixedNow))

	assert.Equal(t, map[string]string{"amfcfg.yaml": string(newData)}, snapshot(t, out))
}

func TestPublishMarkerWithCRLF(t *testing.T) {
	// A generated file checked out with Windows line endings is still generated.
	out := t.TempDir()
	crlf := []byte(generator.GeneratedMarker + "\r\na: 1\r\n")
	require.NoError(t, os.WriteFile(filepath.Join(out, "amfcfg.yaml"), crlf, 0o600))

	require.NoError(t, generator.Publish(out, []generator.File{{Name: "amfcfg.yaml", Data: generated("a: 2\n")}}, fixedNow))
	assert.Len(t, snapshot(t, out), 1, "no backup")
}

func TestPublishEmptyExistingFileIsBackedUp(t *testing.T) {
	out := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(out, "amfcfg.yaml"), nil, 0o600))

	require.NoError(t, generator.Publish(out, []generator.File{{Name: "amfcfg.yaml", Data: generated("a: 2\n")}}, fixedNow))
	assert.Contains(t, snapshot(t, out), "amfcfg-20261122T153000.bak.yaml")
}

func TestPublishBackupCollisionFails(t *testing.T) {
	out := t.TempDir()
	original := "hand written\n"
	oldBackup := "older backup\n"
	require.NoError(t, os.WriteFile(filepath.Join(out, "amfcfg.yaml"), []byte(original), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(out, "amfcfg-20261122T153000.bak.yaml"), []byte(oldBackup), 0o600))
	before := snapshot(t, out)

	err := generator.Publish(out, []generator.File{{Name: "amfcfg.yaml", Data: generated("a: 2\n")}}, fixedNow)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	assert.Contains(t, err.Error(), "already published: []")
	assert.Equal(t, before, snapshot(t, out), "nothing overwritten, temp dir removed")
}

func TestPublishPartialFailureReportsPublishedFiles(t *testing.T) {
	out := t.TempDir()
	// smfcfg.yaml is hand-written and its backup name is already taken, so
	// publishing fails on the second file after amfcfg.yaml is published.
	require.NoError(t, os.WriteFile(filepath.Join(out, "smfcfg.yaml"), []byte("hand written\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(out, "smfcfg-20261122T153000.bak.yaml"), []byte("old\n"), 0o600))

	err := generator.Publish(out, []generator.File{
		{Name: "amfcfg.yaml", Data: generated("a: 1\n")},
		{Name: "smfcfg.yaml", Data: generated("s: 1\n")},
	}, fixedNow)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "publish smfcfg.yaml")
	assert.Contains(t, err.Error(), "already published: [amfcfg.yaml]")

	files := snapshot(t, out)
	assert.Equal(t, "hand written\n", files["smfcfg.yaml"], "failed file left as it was")
	assert.Contains(t, files["amfcfg.yaml"], generator.GeneratedMarker)
	assert.NotContains(t, files, ".cfg-generator-tmp-", "temp dir removed")
}

func TestPublishTargetIsDirectory(t *testing.T) {
	out := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(out, "amfcfg.yaml", "x"), 0o755))

	err := generator.Publish(out, []generator.File{{Name: "amfcfg.yaml", Data: generated("a: 1\n")}}, fixedNow)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "is a directory")
	assert.DirExists(t, filepath.Join(out, "amfcfg.yaml", "x"), "directory left untouched")
}

func TestPublishOutDirIsAFile(t *testing.T) {
	f := writeFile(t, "config", "not a dir")
	err := generator.Publish(f, []generator.File{{Name: "amfcfg.yaml", Data: generated("a: 1\n")}}, fixedNow)
	assert.Error(t, err)
	b, readErr := os.ReadFile(f)
	require.NoError(t, readErr)
	assert.Equal(t, "not a dir", string(b))
}
