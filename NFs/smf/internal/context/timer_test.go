package context_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/free5gc/smf/internal/context"
)

func TestTimerNewTimer(t *testing.T) {
	timer := context.NewTimer(100*time.Millisecond, 3, func(expireTimes int32) {
		t.Logf("expire %d times", expireTimes)
	}, func() {
		t.Log("exceed max retry times (3)")
	})
	assert.NotNil(t, timer)
}

func TestTimerStartAndStop(t *testing.T) {
	timer := context.NewTimer(100*time.Millisecond, 3,
		func(expireTimes int32) {
			t.Logf("expire %d times", expireTimes)
		},
		func() {
			t.Log("exceed max retry times (3)")
		})
	assert.NotNil(t, timer)

	time.Sleep(350 * time.Millisecond)
	timer.Stop()
	assert.EqualValues(t, 3, timer.ExpireTimes())
}

func TestTimerExceedMaxRetryTimes(t *testing.T) {
	timer := context.NewTimer(100*time.Millisecond, 3,
		func(expireTimes int32) {
			t.Logf("expire %d times", expireTimes)
		},
		func() {
			t.Log("exceed max retry times (3)")
		})
	assert.NotNil(t, timer)

	time.Sleep(450 * time.Millisecond)
}
