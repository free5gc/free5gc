package context

import (
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/free5gc/amf/internal/logger"
)

// Timer can be used for retransmission, it will manage retry times automatically
type Timer struct {
	ticker        *time.Ticker
	expireTimes   int32 // accessed atomically
	maxRetryTimes int32 // accessed atomically
	done          chan bool
}

// NewTimer will return a Timer struct and create a goroutine. Then it calls expiredFunc every time interval d until
// the user call Stop(). the number of expire event is be recorded when the timer is active. When the number of expire
// event is > maxRetryTimes, then the timer will call cancelFunc and turns off itself. Whether expiredFunc pass a
// parameter expireTimes to tell the user that the current expireTimes.
func NewTimer(d time.Duration, maxRetryTimes int,
	expiredFunc func(expireTimes int32),
	cancelFunc func(),
) *Timer {
	t := &Timer{}
	atomic.StoreInt32(&t.expireTimes, 0)
	atomic.StoreInt32(&t.maxRetryTimes, int32(maxRetryTimes))
	t.done = make(chan bool, 1)
	t.ticker = time.NewTicker(d)

	go func(ticker *time.Ticker) {
		defer func() {
			if p := recover(); p != nil {
				// Print stack for panic to log. Fatalf() will let program exit.
				logger.CtxLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
			}
		}()

		defer ticker.Stop()

		for {
			select {
			case <-t.done:
				return
			case <-ticker.C:
				atomic.AddInt32(&t.expireTimes, 1)
				if t.ExpireTimes() > t.MaxRetryTimes() {
					cancelFunc()
					return
				} else {
					expiredFunc(t.ExpireTimes())
				}
			}
		}
	}(t.ticker)

	return t
}

// MaxRetryTimes return the max retry times of the timer
func (t *Timer) MaxRetryTimes() int32 {
	return atomic.LoadInt32(&t.maxRetryTimes)
}

// ExpireTimes return the current expire times of the timer
func (t *Timer) ExpireTimes() int32 {
	return atomic.LoadInt32(&t.expireTimes)
}

// Stop turns off the timer, after Stop, no more timeout event will be triggered. User should call Stop() only once
// otherwise it may hang on writing to done channel
func (t *Timer) Stop() {
	t.done <- true
	close(t.done)
}
