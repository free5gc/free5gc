package context

import (
	"time"

	"github.com/free5gc/openapi/models"
)

const (
	MaxNumOfTAI                       int   = 16
	MaxNumOfBroadcastPLMNs            int   = 12
	MaxNumOfPLMNs                     int   = 12
	MaxNumOfSlice                     int   = 1024
	MaxNumOfAllowedSnssais            int   = 8
	MaxValueOfAmfUeNgapId             int64 = 1099511627775
	MaxNumOfServedGuamiList           int   = 256
	MaxNumOfPDUSessions               int   = 256
	MaxNumOfDRBs                      int   = 32
	MaxNumOfAOI                       int   = 64
	MaxT3513RetryTimes                int   = 4
	MaxT3522RetryTimes                int   = 4
	MaxT3550RetryTimes                int   = 4
	MaxT3560RetryTimes                int   = 4
	MaxT3565RetryTimes                int   = 4
	MAxNumOfAlgorithm                 int   = 8
	DefaultT3502                      int   = 720  // 12 min
	DefaultT3512                      int   = 3240 // 54 min
	DefaultNon3gppDeregistrationTimer int   = 3240 // 54 min
)

// timers at AMF side, defined in TS 24.501 table 10.2.2
const (
	TimeT3513 time.Duration = 6 * time.Second
	TimeT3522 time.Duration = 6 * time.Second
	TimeT3550 time.Duration = 6 * time.Second
	TimeT3560 time.Duration = 6 * time.Second
	TimeT3565 time.Duration = 6 * time.Second
)

type CauseAll struct {
	Cause        *models.SmfPduSessionCause
	NgapCause    *models.NgApCause
	Var5GmmCause *int32
}
