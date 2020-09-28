package version_test

import (
	"fmt"
	"free5gc/src/nrf/version"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	t.Run("VERSION not specified", func(t *testing.T) {
		var expected = fmt.Sprintf(
			"\n\tNot specify ldflags (which link version) during go build\n\tgo version: %s %s/%s",
			runtime.Version(),
			runtime.GOOS,
			runtime.GOARCH)
		assert.Equal(t, expected, version.GetVersion())
	})

	t.Run("VERSION specified", func(t *testing.T) {
		var stdout []byte
		stdout, _ = exec.Command("bash", "-c", "cd ../../.. && git describe --tags").Output()
		version.VERSION = strings.TrimSuffix(string(stdout), "\n")
		stdout, _ = exec.Command("bash", "-c", "date -u +\"%Y-%m-%dT%H:%M:%SZ\"").Output()
		version.BUILD_TIME = strings.TrimSuffix(string(stdout), "\n")
		stdout, _ = exec.Command("bash", "-c", "git log --pretty=\"%H\" -1 | cut -c1-8").Output()
		version.COMMIT_HASH = strings.TrimSuffix(string(stdout), "\n")
		stdout, _ = exec.Command("bash", "-c", "git log --pretty=\"%ai\" -1 | awk '{time=$1\"T\"$2\"Z\"; print time}'").Output()
		version.COMMIT_TIME = strings.TrimSuffix(string(stdout), "\n")

		var expected = fmt.Sprintf(
			"\n\tversion: %s\n\tbuild time: %s\n\tcommit hash: %s\n\tcommit time: %s\n\tgo version: %s %s/%s",
			version.VERSION,
			version.BUILD_TIME,
			version.COMMIT_HASH,
			version.COMMIT_TIME,
			runtime.Version(),
			runtime.GOOS,
			runtime.GOARCH)

		assert.Equal(t, expected, version.GetVersion())
	})
}
