package testfixture

import (
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var sampleSection = regexp.MustCompile(`(?m)^([a-z][a-z0-9]*):[ \t]*(?:#[^\r\n]*)?\r?$`)

// NFBlocks copies the sample's top-level NF blocks verbatim so that
// unquoted leading zeros reach the generator's decoder unchanged.
func NFBlocks(t *testing.T, names ...string) string {
	t.Helper()
	data, err := os.ReadFile("../../global-config.yaml")
	require.NoError(t, err)
	sections := sampleSection.FindAllSubmatchIndex(data, -1)
	var blocks strings.Builder
	count := 0
	for i, section := range sections {
		name := string(data[section[2]:section[3]])
		if name == "global" || (len(names) > 0 && !slices.Contains(names, name)) {
			continue
		}
		end := len(data)
		if i+1 < len(sections) {
			end = sections[i+1][0]
		}
		blocks.Write(data[section[0]:end])
		count++
	}
	if len(names) > 0 {
		require.Equal(t, len(names), count, "missing blocks in sample: %v", names)
	}
	return blocks.String()
}
