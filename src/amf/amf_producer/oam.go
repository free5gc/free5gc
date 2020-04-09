package amf_producer

import (
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/amf_handler/amf_message"
	"free5gc/src/amf/gmm/gmm_state"
	"free5gc/src/amf/logger"
	"net/http"
	"strconv"
)

type PduSession struct {
	PduSessionId string
	SmContextRef string
	Sst          string
	Sd           string
	Dnn          string
}

type UEContext struct {
	AccessType models.AccessType
	Supi       string
	Guti       string
	/* Tai */
	Mcc string
	Mnc string
	Tac string
	/* PDU sessions */
	PduSessions []PduSession
	/*Connection state */
	CmState models.CmState
}

type UEContexts []UEContext

func HandleOAMRegisteredUEContext(httpChannel chan amf_message.HandlerResponseMessage, supi string) {
	logger.ProducerLog.Infof("[OAM] Handle Registered UE Context")

	var response UEContexts

	amfSelf := amf_context.AMF_Self()

	if supi != "" {
		if ue, exists := amfSelf.UePool[supi]; exists {
			ueContext := buildUEContext(ue, models.AccessType__3_GPP_ACCESS)
			if ueContext != nil {
				response = append(response, *ueContext)
			}
			ueContext = buildUEContext(ue, models.AccessType_NON_3_GPP_ACCESS)
			if ueContext != nil {
				response = append(response, *ueContext)
			}
		} else {
			problem := models.ProblemDetails{
				Status: http.StatusNotFound,
				Cause:  "CONTEXT_NOT_FOUND",
			}
			amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problem)
			return
		}
	} else {
		for _, ue := range amfSelf.UePool {
			ueContext := buildUEContext(ue, models.AccessType__3_GPP_ACCESS)
			if ueContext != nil {
				response = append(response, *ueContext)
			}
			ueContext = buildUEContext(ue, models.AccessType_NON_3_GPP_ACCESS)
			if ueContext != nil {
				response = append(response, *ueContext)
			}
		}
	}

	amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, response)
}

func buildUEContext(ue *amf_context.AmfUe, accessType models.AccessType) (ueContext *UEContext) {
	if ue.Sm[accessType].Check(gmm_state.REGISTERED) {
		ueContext = &UEContext{
			AccessType: models.AccessType__3_GPP_ACCESS,
			Supi:       ue.Supi,
			Guti:       ue.Guti,
			Mcc:        ue.Tai.PlmnId.Mcc,
			Mnc:        ue.Tai.PlmnId.Mnc,
			Tac:        ue.Tai.Tac,
		}

		for _, smContext := range ue.SmContextList {
			pduSessionContext := smContext.PduSessionContext
			if pduSessionContext != nil {
				if pduSessionContext.AccessType == accessType {
					pduSession := PduSession{
						PduSessionId: strconv.Itoa(int(pduSessionContext.PduSessionId)),
						SmContextRef: pduSessionContext.SmContextRef,
						Sst:          strconv.Itoa(int(pduSessionContext.SNssai.Sst)),
						Sd:           pduSessionContext.SNssai.Sd,
						Dnn:          pduSessionContext.Dnn,
					}
					ueContext.PduSessions = append(ueContext.PduSessions, pduSession)
				}
			}
		}

		if ue.CmConnect(accessType) {
			ueContext.CmState = models.CmState_CONNECTED
		} else {
			ueContext.CmState = models.CmState_IDLE
		}
	}
	return
}
