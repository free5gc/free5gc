package context

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/free5gc/amf/internal/logger"
)

func TestRemoveAndRemoveAllRanUeRaceCondition(t *testing.T) {
	ran := &AmfRan{
		Log: logger.NgapLog.WithField("", ""),
	}

	// create ranUe & store in RanUeList
	for i := 1; i <= 10000; i++ {
		ranUe, err := ran.NewRanUe(int64(i))
		require.NoError(t, err)
		ran.RanUeList.Store(i, ranUe)
	}

	require.NotPanics(t, func() { runRanUeRemove(ran) })
}

func runRanUeRemove(ran *AmfRan) {
	for i := 1; i <= 10000; i++ {
		go ran.RanUeList.Delete(i)
	}
	ran.RemoveAllRanUe(true)
}
