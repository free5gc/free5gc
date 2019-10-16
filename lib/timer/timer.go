//go:binary-only-package

package timer

import (
	"log"
	"time"
)

func StartTimer(seconds int, action func(interface{}), msg interface{}) *time.Timer {}
