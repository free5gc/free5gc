package util

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/free5gc/openapi/amf/Communication"
	"github.com/free5gc/openapi/bsf/Management"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/pcf/AMPolicyControl"
	"github.com/free5gc/openapi/pcf/PolicyAuthorization"
	"github.com/free5gc/openapi/pcf/SMPolicyControl"
	"github.com/free5gc/openapi/udr/DataRepository"
	"github.com/free5gc/pcf/internal/context"
	"github.com/free5gc/pcf/internal/logger"
)

const TimeFormat = time.RFC3339

var (
	PCF_BASIC_PATH                               = "https://localhost:29507"
	ERROR_REQUEST_PARAMETERS                     = "ERROR_REQUEST_PARAMETERS"
	USER_UNKNOWN                                 = "USER_UNKNOWN"
	CONTEXT_NOT_FOUND                            = "CONTEXT_NOT_FOUND"
	ERROR_INITIAL_PARAMETERS                     = "ERROR_INITIAL_PARAMETERS"
	POLICY_CONTEXT_DENIED                        = "POLICY_CONTEXT_DENIED"
	ERROR_TRIGGER_EVENT                          = "ERROR_TRIGGER_EVENT"
	ERROR_TRAFFIC_MAPPING_INFO_REJECTED          = "ERROR_TRAFFIC_MAPPING_INFO_REJECTED"
	BDT_POLICY_NOT_FOUND                         = "BDT_POLICY_NOT_FOUND"
	REQUESTED_SERVICE_NOT_AUTHORIZED             = "REQUESTED_SERVICE_NOT_AUTHORIZED"
	REQUESTED_SERVICE_TEMPORARILY_NOT_AUTHORIZED = "REQUESTED_SERVICE_TEMPORARILY_NOT_AUTHORIZED" // NWDAF
	UNAUTHORIZED_SPONSORED_DATA_CONNECTIVITY     = "UNAUTHORIZED_SPONSORED_DATA_CONNECTIVITY"
	PDU_SESSION_NOT_AVAILABLE                    = "PDU_SESSION_NOT_AVAILABLE"
	APPLICATION_SESSION_CONTEXT_NOT_FOUND        = "APPLICATION_SESSION_CONTEXT_NOT_FOUND"
	ERROR_IDGENERATOR                            = "ERROR_IDGENERATOR"
	PcpErrHttpStatusMap                          = map[string]int32{
		ERROR_REQUEST_PARAMETERS:                     http.StatusBadRequest,
		USER_UNKNOWN:                                 http.StatusBadRequest,
		ERROR_INITIAL_PARAMETERS:                     http.StatusBadRequest,
		ERROR_TRIGGER_EVENT:                          http.StatusBadRequest,
		POLICY_CONTEXT_DENIED:                        http.StatusForbidden,
		ERROR_TRAFFIC_MAPPING_INFO_REJECTED:          http.StatusForbidden,
		REQUESTED_SERVICE_NOT_AUTHORIZED:             http.StatusForbidden,
		REQUESTED_SERVICE_TEMPORARILY_NOT_AUTHORIZED: http.StatusForbidden,
		UNAUTHORIZED_SPONSORED_DATA_CONNECTIVITY:     http.StatusForbidden,
		CONTEXT_NOT_FOUND:                            http.StatusNotFound,
		BDT_POLICY_NOT_FOUND:                         http.StatusNotFound,
		APPLICATION_SESSION_CONTEXT_NOT_FOUND:        http.StatusNotFound,
		PDU_SESSION_NOT_AVAILABLE:                    http.StatusInternalServerError,
	}
)

func GetNpcfAMPolicyCallbackClient() *AMPolicyControl.APIClient {
	configuration := AMPolicyControl.NewConfiguration()
	client := AMPolicyControl.NewAPIClient(configuration)
	return client
}

func GetNpcfSMPolicyCallbackClient() *SMPolicyControl.APIClient {
	configuration := SMPolicyControl.NewConfiguration()
	client := SMPolicyControl.NewAPIClient(configuration)
	return client
}

func GetNpcfPolicyAuthorizationCallbackClient() *PolicyAuthorization.APIClient {
	configuration := PolicyAuthorization.NewConfiguration()
	client := PolicyAuthorization.NewAPIClient(configuration)
	return client
}

func GetNudrClient(uri string) *DataRepository.APIClient {
	configuration := DataRepository.NewConfiguration()
	configuration.SetBasePath(uri)
	client := DataRepository.NewAPIClient(configuration)
	return client
}

// TODO: implement Nbsf
func GetNbsfClient(uri string) *Management.APIClient {
	configuration := Management.NewConfiguration()
	configuration.SetBasePath(uri)
	client := Management.NewAPIClient(configuration)
	return client
}

func GetNamfClient(uri string) *Communication.APIClient {
	configuration := Communication.NewConfiguration()
	configuration.SetBasePath(uri)
	client := Communication.NewAPIClient(configuration)
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
	startTime, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	if err != nil {
		logger.UtilLog.Errorf("startTime parsing error: %+v", err)
	}
	timeWindow.StartTime = &startTime
	lease, err := time.ParseDuration("720h")
	if err != nil {
		logger.UtilLog.Errorf("ParseDuration error: %+v", err)
	}
	stopTime, err := time.Parse(time.RFC3339, time.Now().Add(lease).Format(time.RFC3339))
	if err != nil {
		logger.UtilLog.Errorf("stopTime parsing error: %+v", err)
	}
	timeWindow.StopTime = &stopTime
	return timeWindow
}

func CheckStopTime(StopTime time.Time) bool {
	if StopTime.Before(time.Now()) {
		return false
	} else {
		return true
	}
}

// Convert int data rate bytes to string data rate bits
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

// Return ProblemDatail, errString represent Detail, cause represent Cause of the fields
func GetProblemDetail(errString, cause string) models.ProblemDetails {
	return models.ProblemDetails{
		Status: PcpErrHttpStatusMap[cause],
		Detail: errString,
		Cause:  cause,
	}
}

// GetSMPolicyDnnData returns SMPolicyDnnData derived from SmPolicy data which snssai and dnn match
func GetSMPolicyDnnData(data models.SmPolicyData, snssai *models.Snssai, dnn string) (result *models.SmPolicyDnnData) {
	if snssai == nil || dnn == "" || data.SmPolicySnssaiData == nil {
		return
	}
	snssaiStr := SnssaiModelsToHex(*snssai)
	for key, snssaiData := range data.SmPolicySnssaiData {
		if !strings.EqualFold(key, snssaiStr) {
			continue
		}
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

// MarshToJsonString returns value which can put into NewInterface()
func MarshToJsonString(v interface{}) (result []string) {
	types := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	if types.Kind() == reflect.Slice {
		for i := 0; i < val.Len(); i++ {
			tmp, err := json.Marshal(val.Index(i).Interface())
			if err != nil {
				logger.UtilLog.Errorf("Marshal error: %+v", err)
			}
			result = append(result, string(tmp))
		}
	} else {
		tmp, err := json.Marshal(v)
		if err != nil {
			logger.UtilLog.Errorf("Marshal error: %+v", err)
		}
		result = append(result, string(tmp))
	}
	return
}

// do AND on two byte array
func AndBytes(bytes1, bytes2 []byte) []byte {
	if bytes1 != nil && len(bytes1) == len(bytes2) {
		bytes3 := []byte{}
		for i, b := range bytes1 {
			bytes3 = append(bytes3, b&bytes2[i])
		}
		return bytes3
	}
	return nil
}

// Negotiate Support Feture with PCF
func GetNegotiateSuppFeat(suppFeat string, serviceSuppFeat []byte) string {
	if serviceSuppFeat == nil {
		return ""
	}
	bytes, err := hex.DecodeString(suppFeat)
	if err != nil {
		logger.UtilLog.Errorf("DecodeString error: %+v", err)
	}
	negoSuppFeat := AndBytes(bytes, serviceSuppFeat)
	return hex.EncodeToString(negoSuppFeat)
}

var serviceUriMap = map[models.ServiceName]string{
	models.ServiceName_NPCF_AM_POLICY_CONTROL:   "policies",
	models.ServiceName_NPCF_SMPOLICYCONTROL:     "sm-policies",
	models.ServiceName_NPCF_BDTPOLICYCONTROL:    "bdtpolicies",
	models.ServiceName_NPCF_POLICYAUTHORIZATION: "app-sessions",
}

// Get Resource Uri (location Header) with param id string
func GetResourceUri(name models.ServiceName, id string) string {
	return fmt.Sprintf("%s/%s/%s", context.GetUri(name), serviceUriMap[name], id)
}

// Check if Feature is Supported or not
func CheckSuppFeat(suppFeat string, number int) bool {
	bytes, err := hex.DecodeString(suppFeat)
	if err != nil || len(bytes) < 1 {
		return false
	}
	index := len(bytes) - ((number - 1) / 8) - 1
	shift := uint8((number - 1) % 8)
	if index < 0 {
		return false
	}
	if bytes[index]&(0x01<<shift) > 0 {
		return true
	}
	return false
}

func CheckPolicyControlReqTrig(
	triggers []models.PolicyControlRequestTrigger, reqTrigger models.PolicyControlRequestTrigger,
) bool {
	for _, trigger := range triggers {
		if trigger == reqTrigger {
			return true
		}
	}
	return false
}
