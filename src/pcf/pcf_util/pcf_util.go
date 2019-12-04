package pcf_util

import (
	"encoding/json"
	"errors"
	"fmt"
	"free5gc/lib/Namf_Communication"
	"free5gc/lib/Npcf_AMPolicy"
	"free5gc/lib/Npcf_SMPolicyControl"
	"free5gc/lib/Nudr_DataRepository"
	"free5gc/lib/openapi/models"
	"free5gc/lib/path_util"
	"free5gc/src/pcf/pcf_context"
	"net/http"
	"reflect"
	"time"
)

// Path of HTTP2 key and log file
var (
	PCF_LOG_PATH                        = path_util.Gofree5gcPath("free5gc/pcfsslkey.log")
	PCF_PEM_PATH                        = path_util.Gofree5gcPath("free5gc/support/TLS/pcf.pem")
	PCF_KEY_PATH                        = path_util.Gofree5gcPath("free5gc/support/TLS/pcf.key")
	PCF_CONFIG_PATH                     = path_util.Gofree5gcPath("free5gc/config/pcfcfg.conf")
	PCF_BASIC_PATH                      = "https://localhost:29507"
	ERROR_REQUEST_PARAMETERS            = "ERROR_REQUEST_PARAMETERS"
	USER_UNKNOWN                        = "USER_UNKNOWN"
	CONTEXT_NOT_FOUND                   = "CONTEXT_NOT_FOUND"
	ERROR_INITIAL_PARAMETERS            = "ERROR_INITIAL_PARAMETERS"
	POLICY_CONTEXT_DENIED               = "POLICY_CONTEXT_DENIED"
	ERROR_TRIGGER_EVENT                 = "ERROR_TRIGGER_EVENT"
	ERROR_TRAFFIC_MAPPING_INFO_REJECTED = "ERROR_TRAFFIC_MAPPING_INFO_REJECTED"
	PcpErrHttpStatusMap                 = map[string]int32{
		ERROR_REQUEST_PARAMETERS:            http.StatusBadRequest,
		USER_UNKNOWN:                        http.StatusBadRequest,
		CONTEXT_NOT_FOUND:                   http.StatusNotFound,
		ERROR_INITIAL_PARAMETERS:            http.StatusBadRequest,
		POLICY_CONTEXT_DENIED:               http.StatusForbidden,
		ERROR_TRIGGER_EVENT:                 http.StatusBadRequest,
		ERROR_TRAFFIC_MAPPING_INFO_REJECTED: http.StatusForbidden,
	}
)

func GetNpcfAMPolicyCallbackClient() *Npcf_AMPolicy.APIClient {
	configuration := Npcf_AMPolicy.NewConfiguration()
	client := Npcf_AMPolicy.NewAPIClient(configuration)
	return client
}
func GetNpcfSMPolicyCallbackClient() *Npcf_SMPolicyControl.APIClient {
	configuration := Npcf_SMPolicyControl.NewConfiguration()
	client := Npcf_SMPolicyControl.NewAPIClient(configuration)
	return client
}

func GetNudrClient(uri string) *Nudr_DataRepository.APIClient {
	configuration := Nudr_DataRepository.NewConfiguration()
	configuration.SetBasePath(uri)
	client := Nudr_DataRepository.NewAPIClient(configuration)
	return client
}
func GetNamfClient(uri string) *Namf_Communication.APIClient {
	configuration := Namf_Communication.NewConfiguration()
	configuration.SetBasePath(uri)
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

func GetProblemDetail(errString, cause string) models.ProblemDetails {
	return models.ProblemDetails{
		Status: PcpErrHttpStatusMap[cause],
		Detail: errString,
		Cause:  cause,
	}
}

func GetSMPolicyDnnData(data models.SmPolicyData, snssai *models.Snssai, dnn string) (result *models.SmPolicyDnnData) {
	if snssai == nil || dnn == "" || data.SmPolicySnssaiData == nil {
		return
	}
	snssaiString := SnssaiModelsToHex(*snssai)
	if snssaiData, exist := data.SmPolicySnssaiData[snssaiString]; exist {
		if snssaiData.SmPolicyDnnData == nil {
			return
		}
		if dnnInfo, exist := snssaiData.SmPolicyDnnData[dnn]; exist {
			result = &dnnInfo
			return
		}
	}
	return

}

func MarshToJsonString(v interface{}) (result []string) {
	types := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	if types.Kind() == reflect.Slice {
		for i := 0; i < val.Len(); i++ {
			tmp, _ := json.Marshal(val.Index(i).Interface())
			result = append(result, string(tmp))

		}
	} else {
		tmp, _ := json.Marshal(v)
		result = append(result, string(tmp))
	}
	return
}
