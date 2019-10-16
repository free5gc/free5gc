package amf_producer

import (
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/amf_handler/amf_message"
	"free5gc/src/amf/logger"
	"net/http"
	"strings"
)

func HandleProvideDomainSelectionInfoRequest(httpChannel chan amf_message.HandlerResponseMessage, ueContextId string, infoClass string) {
	var response models.UeContextInfo
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
	// TODO: Error Status 307, 403 in TS29.518 Table 6.3.3.3.3.1-3
	if ue != nil {
		anType := ue.GetAnType()
		if anType != "" && infoClass != "" {
			ranUe := ue.RanUe[anType]
			response.AccessType = anType
			response.LastActTime = ranUe.LastActTime
			response.RatType = ue.RatType
			response.SupportedFeatures = ranUe.SupportedFeatures
			response.SupportVoPS = ranUe.SupportVoPS
			response.SupportVoPSn3gpp = ranUe.SupportVoPSn3gpp
		}

	} else {
		logger.ProducerLog.Errorln("ue is nil")
	}

	amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, response)
}
