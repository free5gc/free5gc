package pcf_util

import (
	"errors"
	"fmt"
	"free5gc/lib/Namf_Communication"
	"free5gc/lib/Nudr_DataRepository"
	"free5gc/lib/openapi/models"
	"free5gc/lib/path_util"
	"free5gc/src/pcf/pcf_context"
	"time"
)

// Path of HTTP2 key and log file
var (
	PCF_LOG_PATH    = path_util.Gofree5gcPath("free5gc/pcfsslkey.log")
	PCF_PEM_PATH    = path_util.Gofree5gcPath("free5gc/support/TLS/pcf.pem")
	PCF_KEY_PATH    = path_util.Gofree5gcPath("free5gc/support/TLS/pcf.key")
	PCF_CONFIG_PATH = path_util.Gofree5gcPath("free5gc/config/pcfcfg.conf")
	PCF_BASIC_PATH  = "https://localhost:29507"
)

func GetNudrClient() *Nudr_DataRepository.APIClient {
	configuration := Nudr_DataRepository.NewConfiguration()
	BasePath := pcf_context.PCF_Self().UdrUri
	configuration.SetBasePath(BasePath)
	client := Nudr_DataRepository.NewAPIClient(configuration)
	return client
}
func GetNamfClient() *Namf_Communication.APIClient {
	configuration := Namf_Communication.NewConfiguration()
	configuration.SetBasePath("https://localhost:29518")
	client := Namf_Communication.NewAPIClient(configuration)
	return client
}

func GetDefaultDataRate() models.UsageThreshold {
	var usageThreshold models.UsageThreshold
	usageThreshold.DownlinkVolume = 1024 * 1024 / 8 // 1 Mbps
	usageThreshold.UplinkVolume = 1024 * 1024 / 8   // 1 Mbps
	return usageThreshold
}

func GetDefaultTime() models.TimeWindow {
	var timeWindow models.TimeWindow
	*timeWindow.StartTime = time.Now()
	lease, _ := time.ParseDuration("720h")
	*timeWindow.StopTime = time.Now().Add(lease)
	return timeWindow
}

func TimeParse(timeParse time.Time) (time.Time, error) {
	timeParse, err := time.Parse(pcf_context.GetTimeformat(), timeParse.Format(pcf_context.GetTimeformat()))
	if err == nil {
		return timeParse, nil
	} else {
		return timeParse, errors.New(" can't parse time ")
	}
}

func CheckStopTime(StopTime time.Time) bool {
	if StopTime.Before(time.Now()) {
		return false
	} else {
		return true
	}
}

// convert int data rate bytes to string data rate bits
func Convert(bytes int64) (DateRate string) {
	BitDateRate := float64(bytes) * 8
	if BitDateRate/1024 > 0 && BitDateRate/1024/1024 < 0 {
		DateRate = fmt.Sprintf("%.2f", BitDateRate/1024) + " Kbps"
	} else if BitDateRate/1024/1024 > 0 {
		DateRate = fmt.Sprintf("%.2f", BitDateRate/1024/1024) + " Mbps"
	} else {
		DateRate = fmt.Sprintf("%.2f", BitDateRate) + " bps"
	}
	return DateRate
}
