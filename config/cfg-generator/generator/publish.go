package generator

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// File is one generated output file.
type File struct {
	Name string // e.g. "amfcfg.yaml"
	Data []byte
}

// backupTimeLayout has no ':' so backup names are valid on Windows too.
const backupTimeLayout = "20060102T150405"

// BackupName returns the backup file name for name at time t, e.g.
// "amfcfg.yaml" -> "amfcfg-20261122T153000.bak.yaml".
func BackupName(name string, t time.Time) string {
	ext := filepath.Ext(name)
	return strings.TrimSuffix(name, ext) + "-" + t.Format(backupTimeLayout) + ".bak" + ext
}

// Publish writes files into outDir.
//
//  1. All files are written to a temporary directory inside outDir (same
//     filesystem, so the final rename is atomic per file).
//  2. For each file, an existing target without GeneratedMarker on its first
//     line is treated as hand-written and renamed to BackupName first; a
//     target with the marker is simply replaced.
//  3. The temporary directory is removed.
//
// If step 1 fails, outDir is left untouched. If step 2 fails part-way, the
// error lists which files were already published.
func Publish(outDir string, files []File, now time.Time) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	tmp, err := os.MkdirTemp(outDir, ".cfg-generator-tmp-")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(tmp) }()

	for _, f := range files {
		if err = os.WriteFile(filepath.Join(tmp, f.Name), f.Data, 0o644); err != nil {
			return fmt.Errorf("stage %s: %w", f.Name, err)
		}
	}

	var published []string
	for _, f := range files {
		target := filepath.Join(outDir, f.Name)
		if err = backupIfHandWritten(target, now); err != nil {
			return publishErr(f.Name, published, err)
		}
		if err = os.Rename(filepath.Join(tmp, f.Name), target); err != nil {
			return publishErr(f.Name, published, err)
		}
		published = append(published, f.Name)
	}
	return nil
}

func publishErr(name string, published []string, err error) error {
	return fmt.Errorf("publish %s: %w (already published: [%s])",
		name, err, strings.Join(published, ", "))
}

func backupIfHandWritten(target string, now time.Time) error {
	if fi, err := os.Stat(target); err == nil && fi.IsDir() {
		return fmt.Errorf("%s is a directory", target)
	}
	generated, err := hasMarker(target)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if generated {
		return nil
	}

	backup := filepath.Join(filepath.Dir(target), BackupName(filepath.Base(target), now))
	if _, err = os.Stat(backup); err == nil {
		return fmt.Errorf("backup %s already exists", backup)
	}
	return os.Rename(target, backup)
}

func hasMarker(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer func() { _ = f.Close() }()

	line, err := bufio.NewReader(f).ReadString('\n')
	if err != nil && line == "" {
		return false, nil // empty file: not generated
	}
	return strings.TrimRight(line, "\r\n") == GeneratedMarker, nil
}
