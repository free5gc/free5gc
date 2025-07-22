package context

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimerNewTimer(t *testing.T) {
	timer := NewTimer(100*time.Millisecond, 3, func(expireTimes int32) {
		t.Logf("expire %d times", expireTimes)
	}, func() {
		t.Log("exceed max retry times (3)")
	})
	assert.NotNil(t, timer)
}

func TestTimerStartAndStop(t *testing.T) {
	timer := NewTimer(100*time.Millisecond, 3,
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
	timer := NewTimer(100*time.Millisecond, 3,
		func(expireTimes int32) {
			t.Logf("expire %d times", expireTimes)
		},
		func() {
			t.Log("exceed max retry times (3)")
		})
	assert.NotNil(t, timer)

	time.Sleep(450 * time.Millisecond)
}

/*
func TestTimerRestartTimerWithoutExceedMaxRetry(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("timer %d", i), func(t *testing.T) {
			t.Parallel()
			for j := 0; j < 10; j++ {
				t.Logf("timer start %d-th", j)
				timer := NewTimer(100*time.Millisecond, 10, func(expireTimes int32) {
					t.Logf("expire %d times", expireTimes)
				}, func() {
					t.Log("exceed max retry times")
				})
				time.Sleep(50 * time.Millisecond)
				timer.Stop()
			}
		})
	}
}

func TestTimerRestartTimerWithExceedMaxRetry(t *testing.T) {
	for i := 0; i < 50; i++ {
		t.Run(fmt.Sprintf("timer %d", i), func(t *testing.T) {
			t.Parallel()
			for j := 0; j < 10; j++ {
				t.Logf("timer start %d-th", j)
				timer := NewTimer(50*time.Millisecond, 3, func(expireTimes int32) {
					t.Logf("expire %d times", expireTimes)
				}, func() {
					t.Log("exceed max retry times")
				})
				time.Sleep(200 * time.Millisecond)
				timer.Stop()
			}
		})
	}
}
*/
