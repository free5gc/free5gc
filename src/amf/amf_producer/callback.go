package amf_producer

import (
	"fmt"
	"github.com/mohae/deepcopy"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_consumer"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/amf_handler/amf_message"
	"free5gc/src/amf/gmm/gmm_message"
	"free5gc/src/amf/logger"
	"net/http"
	"strconv"
)

func HandleSmContextStatusNotify(httpChannel chan amf_message.HandlerResponseMessage, guti, pduSessionIdString string, body models.SmContextStatusNotification) {
	var problem models.ProblemDetails
	amfSelf := amf_context.AMF_Self()
	ue := amfSelf.AmfUeFindByGuti(guti)
	if ue == nil {
		problem.Status = 404
		problem.Cause = "CONTEXT_NOT_FOUND"
		problem.Detail = fmt.Sprintf("Guti[%s] Not Found", guti)
		amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problem)
		return
	}
	pduSessionID, _ := strconv.Atoi(pduSessionIdString)
	_, ok := ue.SmContextList[int32(pduSessionID)]
	if !ok {
		problem.Status = 404
		problem.Cause = "CONTEXT_NOT_FOUND"
		problem.Detail = fmt.Sprintf("PDUSessionID[%d] Not Found", pduSessionID)
		amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problem)
		return
	}
	logger.CallbackLog.Debugf("Release PDUSessionId[%d] of UE[%s] By SmContextStatus Notification because of %s", pduSessionID, ue.Supi, body.StatusInfo.Cause)
	pduSessionId := int32(pduSessionID)
	delete(ue.SmContextList, pduSessionId)
	amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNoContent, nil)
	if storedSmContext, exist := ue.StoredSmContext[pduSessionId]; exist {

		smContextCreateData := amf_consumer.BuildCreateSmContextRequest(ue, *storedSmContext.PduSessionContext, models.RequestType_INITIAL_REQUEST)

		response, smContextRef, errResponse, problemDetail, err := amf_consumer.SendCreateSmContextRequest(ue, storedSmContext.SmfUri, storedSmContext.Payload, smContextCreateData)
		if response != nil {
			var smContext amf_context.SmContext
			smContext.PduSessionContext = storedSmContext.PduSessionContext
			smContext.PduSessionContext.SmContextRef = smContextRef
			smContext.UserLocation = deepcopy.Copy(ue.Location).(models.UserLocation)
			smContext.SmfUri = storedSmContext.SmfUri
			smContext.SmfId = storedSmContext.SmfId
			ue.SmContextList[pduSessionId] = &smContext
			logger.CallbackLog.Infof("Http create smContext[pduSessionID: %d] Success", pduSessionId)
			// TODO: handle response(response N2SmInfo to RAN if exists)
		} else if errResponse != nil {
			logger.CallbackLog.Warnf("PDU Session Establishment Request is rejected by SMF[pduSessionId:%d]\n", pduSessionId)
			gmm_message.SendDLNASTransport(ue.RanUe[storedSmContext.AnType], nasMessage.PayloadContainerTypeN1SMInfo, errResponse.BinaryDataN1SmInfoToUe, &pduSessionId, 0, nil, 0)
		} else if err != nil {
			logger.CallbackLog.Errorf("Failed to Create smContext[pduSessionID: %d], Error[%s]\n", pduSessionID, err.Error())
		} else {
			logger.CallbackLog.Errorf("Failed to Create smContext[pduSessionID: %d], Error[%v]\n", pduSessionID, problemDetail)
		}
		delete(ue.StoredSmContext, pduSessionId)

	}
}
