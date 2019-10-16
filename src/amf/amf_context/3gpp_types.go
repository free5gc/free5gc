package amf_context

import (
	"free5gc/lib/openapi/models"
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
	MaxPagingRetryTime                int   = 3
	TimeT3513                         int   = 2
	MaxNotificationRetryTime          int   = 5
	TimeT3565                         int   = 6
	TimeT3560                         int   = 6
	MaxT3560RetryTimes                int   = 4
	MAxNumOfAlgorithm                 int   = 8
	TimeT3550                         int   = 6
	MaxT3550RetryTimes                int   = 4
	DefaultT3502                      int   = 720  // 12 min
	DefaultT3512                      int   = 3240 // 54 min
	DefaultNon3gppDeregistrationTimer int   = 3240 // 54 min
	TimeT3522                         int   = 6
	MaxT3522RetryTimes                int   = 4
)

type LADN struct {
	Ladn     string
	TaiLists []models.Tai
}

type CauseAll struct {
	Cause        *models.Cause
	NgapCause    *models.NgApCause
	Var5GmmCause *int32
}
