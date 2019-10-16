package amf_producer

import (
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/amf_handler/amf_message"
	"free5gc/src/amf/logger"
	"net/http"
	"strings"
)

func HandleProvideLocationInfoRequest(httpChannel chan amf_message.HandlerResponseMessage, ueContextId string, body models.RequestLocInfo) {
	var response models.ProvideLocInfo
	var problem models.ProblemDetails
	var ue *amf_context.AmfUe
	var ok bool
	amfSelf := amf_context.AMF_Self()
	if strings.HasPrefix(ueContextId, "imsi") {

		if ue, ok = amfSelf.UePool[ueContextId]; !ok {
			problem.Status = 404
			problem.Cause = "CONTEXT_NOT_FOUND"
			amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problem)
			return
		}
	} else if strings.HasPrefix(ueContextId, "imei") {
		for _, ue1 := range amfSelf.UePool {
			if ue1.Pei == ueContextId {
				ue = ue1
				break
			}
		}
		if ue == nil {
			problem.Status = 404
			problem.Cause = "CONTEXT_NOT_FOUND"
			amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problem)
			return
		}
	}

	requestData := body
	anType := ue.GetAnType()
	if anType == "" {
		problem.Status = 404
		problem.Cause = "CONTEXT_NOT_FOUND"
		amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problem)
		return
	}

	if ue != nil {
		ranUe := ue.RanUe[anType]
		if requestData.Req5gsLoc || requestData.ReqCurrentLoc {
			response.CurrentLoc = true
			response.Location = &ue.Location
		}

		if requestData.ReqRatType {
			response.RatType = ue.RatType
		}

		if requestData.ReqTimeZone {
			response.Timezone = ue.TimeZone
		}

		if requestData.SupportedFeatures != "" {
			response.SupportedFeatures = ranUe.SupportedFeatures
		}
	} else {
		logger.LocationLog.Errorln("ue is nil")
	}
	amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, response)
}
